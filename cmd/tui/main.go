package main

import (
	"context"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nerdminertui/nerdtui/internal/config"
	"github.com/nerdminertui/nerdtui/internal/model"
	"github.com/nerdminertui/nerdtui/internal/store"
	"github.com/nerdminertui/nerdtui/internal/ui"
	"github.com/nerdminertui/nerdtui/internal/worker"
	"github.com/nerdminertui/nerdtui/pkg/mining"
)

var version = "dev"

func main() {
	configPath := flag.String("config", "", "Path to config file")
	mockMode := flag.Bool("mock", false, "Enable mock mining mode")
	noStore := flag.Bool("no-store", false, "Disable SQLite store")
	cpuTargetFlag := flag.Float64("cpu", 0.5, "Target CPU utilization (0.05 - 1.0)")
	flag.Parse()

	// Load configuration
	// Note: config.Load() does not take arguments, it reads from Viper/env.
	// If a config file was specified, we'd need to modify config.Load to accept it,
	// but per internal/config spec it doesn't take arguments.
	_ = configPath // Avoid unused variable
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Override config with flags if provided
	if *mockMode {
		cfg.MockMining = true
	}
	if *noStore {
		cfg.StorePath = ""
	}

	configValid := true
	if !cfg.MockMining && cfg.BTCAddress == "" {
		configValid = false
	}

	// Initialize Storage
	var st store.Store
	if cfg.StorePath == "" {
		st = store.NewNilStore()
	} else {
		storeDir := filepath.Dir(cfg.StorePath)
		if err := os.MkdirAll(storeDir, 0755); err != nil {
			log.Fatalf("Failed to create store directory: %v", err)
		}
		sqlStore, err := store.NewSQLiteStore(cfg.StorePath)
		if err != nil {
			log.Fatalf("Failed to initialize SQLite store: %v", err)
		}
		defer sqlStore.Close()
		st = sqlStore
	}
	_ = st // Used in background or future tasks

	// Initialize UI Channels
	outCh := make(chan tea.Msg, 100)
	throttleCh := make(chan float64, 10)

	jobCh := make(chan mining.Job, 10)

	// Initialize Pool Client factory function
	createClient := func(mock bool, addr string, port int, btcAddr string, workerName string) worker.PoolClient {
		if mock {
			return &worker.MockPoolClient{}
		}
		return worker.NewStratumClient(addr, port, btcAddr, workerName, outCh, jobCh)
	}

	// Setup Initial Job (Dummy for now)
	initialJob := mining.Job{
		Header: make([]byte, 80),
		Target: [32]byte{0x00},
	}

	configUpdateCh := make(chan *model.AppState, 1)

	// Build initial state
	initialState := model.AppState{
		CPUTarget:       *cpuTargetFlag,
		Screen:          model.ScreenDashboard,
		StartedAt:       time.Now(),
		PoolAddress:     cfg.PoolAddress,
		PoolPort:        cfg.PoolPort,
		WorkerName:      cfg.WorkerName,
		BTCAddress:      cfg.BTCAddress,
		MockMining:      cfg.MockMining,
		DebugMode:       cfg.DebugMode,
		ConfigValid:     configValid,
	}

	if st != nil {
		if history, err := st.QueryHashRateHistory(context.Background(), model.HashHistoryLen); err == nil {
			offset := model.HashHistoryLen - len(history)
			for i, val := range history {
				initialState.HashRateHistory[offset+i] = val
			}
		}
	}

	if cfg.MockMining {
		initialState.PoolConnected = true
		initialState.ConnectionStatus = "Mock"
	} else {
		initialState.PoolConnected = false
		initialState.ConnectionStatus = "Desconectado"
	}

	if !configValid {
		initialState.Screen = model.ScreenSettings
	}

	// Build App Model
	app := ui.NewAppModel(initialState, throttleCh, configUpdateCh)

	// Setup Debug Logging
	var logFile *os.File
	if cfg.DebugMode {
		var err error
		logFile, err = os.OpenFile("debug.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Fatal: could not open debug log file: %v", err)
		}
		log.SetOutput(logFile)
		log.Printf("--- NerdMiner TUI Started ---")
	} else {
		log.SetOutput(io.Discard)
	}
	if logFile != nil {
		defer logFile.Close()
	}

	// Initialize Bubbletea Program
	popts := []tea.ProgramOption{tea.WithAltScreen()}
	if flag.Lookup("test.v") != nil {
		popts = append(popts, tea.WithInput(os.Stdin))
	}
	p := tea.NewProgram(app, popts...)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Route background messages to Bubbletea program
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-outCh:
				if hrMsg, ok := msg.(worker.HashRateMsg); ok {
					if st != nil {
						_ = st.AppendHashRate(ctx, hrMsg.HPS, time.Now())
					}
				}
				p.Send(msg)
			}
		}
	}()

	// Mempool Network Stats Poller
	go func() {
		if cfg.MockMining {
			outCh <- worker.PoolStatsMsg{
				GlobalHashRate:    600000000000,
				BlockHeight:       850000,
				NetworkDifficulty: 83000000000000,
			}
			return
		}

		mempoolClient := &worker.MempoolClient{BaseURL: "https://mempool.space"}
		// Fetch immediately at startup
		if stats, err := mempoolClient.FetchStats(ctx); err == nil {
			outCh <- stats
		}
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if stats, err := mempoolClient.FetchStats(ctx); err == nil {
					outCh <- stats
				}
			}
		}
	}()

	// Worker manager loop
	go func() {
		var workerCtx context.Context
		var workerCancel context.CancelFunc

		startWorker := func(s *model.AppState) {
			if workerCancel != nil {
				workerCancel()
			}
			if !s.ConfigValid {
				return
			}
			workerCtx, workerCancel = context.WithCancel(ctx)
			client := createClient(s.MockMining, s.PoolAddress, s.PoolPort, s.BTCAddress, s.WorkerName)
			miner := worker.NewMinerWorker(client, s.CPUTarget, initialJob, outCh, throttleCh, jobCh)
			go client.Run(workerCtx)
			go miner.Run(workerCtx)
		}

		startWorker(&initialState)

		for {
			select {
			case <-ctx.Done():
				if workerCancel != nil {
					workerCancel()
				}
				return
			case newState := <-configUpdateCh:
				// Validate config
				valid := true
				if !newState.MockMining && newState.BTCAddress == "" {
					valid = false
				}
				newState.ConfigValid = valid
				
				// Save config
				cfg.PoolAddress = newState.PoolAddress
				cfg.PoolPort = newState.PoolPort
				cfg.WorkerName = newState.WorkerName
				cfg.MockMining = newState.MockMining
				cfg.BTCAddress = newState.BTCAddress
				cfg.DebugMode = newState.DebugMode
				_ = config.Save(cfg)

				// Update logging mode
				if cfg.DebugMode {
					if logFile == nil {
						logFile, _ = os.OpenFile("debug.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
					}
					if logFile != nil {
						log.SetOutput(logFile)
					}
				} else {
					log.SetOutput(io.Discard)
				}

				startWorker(newState)
			}
		}
	}()

	// Execute TUI
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error running program: %v", err)
	}
}

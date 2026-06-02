package worker

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nerdminertui/nerdtui/pkg/mining"
)

// PoolClient abstracts communication with the pool.
type PoolClient interface {
	FetchStats(ctx context.Context) (PoolStatsMsg, error)
	SubmitShare(ctx context.Context, nonce uint32, hash [32]byte) (bool, error)
	Run(ctx context.Context)
}

// MockPoolClient simulates a pool connection for offline execution.
type MockPoolClient struct{}

// FetchStats returns simulated global statistics.
func (c *MockPoolClient) FetchStats(ctx context.Context) (PoolStatsMsg, error) {
	return PoolStatsMsg{
		GlobalHashRate:    4.5e18,
		NetworkDifficulty: 88.1e12,
		BlockHeight:       850000,
	}, nil
}

// SubmitShare unconditionally accepts the mock share.
func (c *MockPoolClient) SubmitShare(ctx context.Context, nonce uint32, hash [32]byte) (bool, error) {
	return true, nil
}

// Run simulates background processing.
func (c *MockPoolClient) Run(ctx context.Context) {
	<-ctx.Done()
}

// MempoolClient implements REST polling for global network stats via mempool.space.
type MempoolClient struct {
	BaseURL string
}

// FetchStats retrieves global stats over HTTP from mempool.space.
func (c *MempoolClient) FetchStats(ctx context.Context) (PoolStatsMsg, error) {
	var stats PoolStatsMsg

	// Create a custom HTTP client with timeout
	client := &http.Client{Timeout: 10 * time.Second}

	// 1. Fetch Block Height
	reqHeight, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/api/blocks/tip/height", nil)
	if err == nil {
		if resp, err := client.Do(reqHeight); err == nil {
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			heightStr := strings.TrimSpace(string(body))
			if h, err := strconv.Atoi(heightStr); err == nil {
				stats.BlockHeight = h
			}
		}
	}

	// 2. Fetch Hashrate & Difficulty
	reqHashrate, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/api/v1/mining/hashrate/3d", nil)
	if err == nil {
		if resp, err := client.Do(reqHashrate); err == nil {
			defer resp.Body.Close()
			var hrData struct {
				CurrentHashrate   float64 `json:"currentHashrate"`
				CurrentDifficulty float64 `json:"currentDifficulty"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&hrData); err == nil {
				stats.GlobalHashRate = hrData.CurrentHashrate
				stats.NetworkDifficulty = hrData.CurrentDifficulty
			}
		}
	}

	return stats, nil
}

// SubmitShare is not implemented for HTTP client.
func (c *MempoolClient) SubmitShare(ctx context.Context, nonce uint32, hash [32]byte) (bool, error) {
	return false, nil
}

// Run is a no-op for HTTP client.
func (c *MempoolClient) Run(ctx context.Context) {
	<-ctx.Done()
}

// StratumPoolClient implements JSON-RPC over TCP for mining.
type StratumPoolClient struct {
	Address    string
	Port       int
	BTCAddress string
	WorkerName string
	OutCh      chan<- tea.Msg
	JobCh      chan<- mining.Job
	
	conn            net.Conn
	mu              sync.Mutex
	reqID           int
	extranonce1     string
	extranonce2Size int
	extranonce2     uint32
	difficulty      float64
	lastJob         *mining.Job
	pendingRequests map[int]chan JSONRPCResponse
	lastRxTime      time.Time
}

// NewStratumClient initializes a new Stratum client.
func NewStratumClient(addr string, port int, btcAddr, workerName string, outCh chan<- tea.Msg, jobCh chan<- mining.Job) *StratumPoolClient {
	return &StratumPoolClient{
		Address:    addr,
		Port:       port,
		BTCAddress: btcAddr,
		WorkerName: workerName,
		OutCh:      outCh,
		JobCh:      jobCh,
	}
}

// Run connects to the stratum server and handles the IO loop.
func (c *StratumPoolClient) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := c.connectAndLoop(ctx)
			if err != nil {
				c.OutCh <- PoolErrorMsg{Err: err}
			}
			c.OutCh <- ConnectionStatusMsg{Status: "Desconectado"}
			
			// Backoff before reconnecting
			select {
			case <-time.After(5 * time.Second):
			case <-ctx.Done():
				return
			}
		}
	}
}

func (c *StratumPoolClient) connectAndLoop(ctx context.Context) error {
	c.OutCh <- ConnectionStatusMsg{Status: "Conectando..."}
	dialer := net.Dialer{Timeout: 10 * time.Second}
	addr := fmt.Sprintf("%s:%d", c.Address, c.Port)
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return err
	}
	
	c.mu.Lock()
	c.conn = conn
	c.lastRxTime = time.Now()
	c.mu.Unlock()
	
	defer func() {
		c.mu.Lock()
		if c.conn != nil {
			c.conn.Close()
			c.conn = nil
		}
		c.mu.Unlock()
	}()

	c.OutCh <- ConnectionStatusMsg{Status: "Conectado"}

	if err := c.sendSubscribe(); err != nil {
		return err
	}
	if err := c.sendAuthorize(); err != nil {
		return err
	}

	// Keep-alive watchdog
	watchdogCtx, cancelWatchdog := context.WithCancel(ctx)
	defer cancelWatchdog()
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-watchdogCtx.Done():
				return
			case <-ticker.C:
				c.mu.Lock()
				last := c.lastRxTime
				c.mu.Unlock()
				
				if time.Since(last) > 3*time.Minute {
					// Pool has been silent for 3 minutes, send a ping
					_, err := c.sendAndWait("mining.suggest_difficulty", []interface{}{0.1}, 10*time.Second)
					if err != nil {
						// Connection is dead, close it to break the scanner loop
						c.mu.Lock()
						if c.conn != nil {
							c.conn.Close()
						}
						c.mu.Unlock()
						return
					}
				}
			}
		}
	}()

	scanner := bufio.NewScanner(conn)
	for {
		// Set a hard deadline to prevent eternal blocking
		conn.SetReadDeadline(time.Now().Add(10 * time.Minute))
		if !scanner.Scan() {
			break
		}
		
		c.mu.Lock()
		c.lastRxTime = time.Now()
		c.mu.Unlock()
		
		line := scanner.Bytes()
		var notif JSONRPCNotification
		if err := json.Unmarshal(line, &notif); err == nil && notif.Method != "" {
			c.handleNotification(notif)
			continue
		}
		var resp JSONRPCResponse
		if err := json.Unmarshal(line, &resp); err == nil {
			c.handleResponse(resp)
		}
	}
	return scanner.Err()
}

func (c *StratumPoolClient) nextID() int {
	// MUST be called while c.mu is already held.
	c.reqID++
	return c.reqID
}

func (c *StratumPoolClient) sendAndWait(method string, params []interface{}, timeout time.Duration) (JSONRPCResponse, error) {
	c.mu.Lock()
	if c.conn == nil {
		c.mu.Unlock()
		return JSONRPCResponse{}, fmt.Errorf("not connected")
	}
	id := c.reqID + 1
	c.reqID = id
	
	respCh := make(chan JSONRPCResponse, 1)
	if c.pendingRequests == nil {
		c.pendingRequests = make(map[int]chan JSONRPCResponse)
	}
	c.pendingRequests[id] = respCh
	
	req := JSONRPCRequest{
		ID:     id,
		Method: method,
		Params: params,
	}
	data, _ := json.Marshal(req)
	data = append(data, '\n')
	_, err := c.conn.Write(data)
	c.mu.Unlock()
	
	if err != nil {
		c.mu.Lock()
		delete(c.pendingRequests, id)
		c.mu.Unlock()
		return JSONRPCResponse{}, err
	}
	
	select {
	case resp := <-respCh:
		return resp, nil
	case <-time.After(timeout):
		c.mu.Lock()
		delete(c.pendingRequests, id)
		c.mu.Unlock()
		return JSONRPCResponse{}, fmt.Errorf("timeout waiting for response")
	}
}

func (c *StratumPoolClient) send(method string, params []interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn == nil {
		return fmt.Errorf("not connected")
	}
	req := JSONRPCRequest{
		ID:     c.nextID(), // safe: mutex already held
		Method: method,
		Params: params,
	}
	data, _ := json.Marshal(req)
	data = append(data, '\n')
	_, err := c.conn.Write(data)
	return err
}

func (c *StratumPoolClient) sendSubscribe() error {
	return c.send("mining.subscribe", []interface{}{"NerdMinerTUI/1.0", nil})
}

func (c *StratumPoolClient) sendAuthorize() error {
	worker := c.BTCAddress
	if c.WorkerName != "" {
		worker = worker + "." + c.WorkerName
	}
	return c.send("mining.authorize", []interface{}{worker, "x"})
}

func (c *StratumPoolClient) handleNotification(notif JSONRPCNotification) {
	if notif.Method == "mining.set_difficulty" {
		var params []interface{}
		if err := json.Unmarshal(notif.Params, &params); err == nil && len(params) > 0 {
			if diff, ok := params[0].(float64); ok {
				c.mu.Lock()
				c.difficulty = diff
				c.mu.Unlock()
			}
		}
		return
	}

	if notif.Method == "mining.notify" {
		var params []interface{}
		if err := json.Unmarshal(notif.Params, &params); err == nil && len(params) >= 9 {
			jobID, _ := params[0].(string)
			prevhashHex, _ := params[1].(string)
			coinb1Hex, _ := params[2].(string)
			coinb2Hex, _ := params[3].(string)
			
			merkleBranchIfaces, _ := params[4].([]interface{})
			var merkleBranchHex []string
			for _, mb := range merkleBranchIfaces {
				if s, ok := mb.(string); ok {
					merkleBranchHex = append(merkleBranchHex, s)
				}
			}
			
			versionHex, _ := params[5].(string)
			nbitsHex, _ := params[6].(string)
			ntimeHex, _ := params[7].(string)
			
			c.mu.Lock()
			en1 := c.extranonce1
			en2Size := c.extranonce2Size
			en2 := c.extranonce2
			c.extranonce2++ // Increment for next job
			poolDiff := c.difficulty
			if poolDiff <= 0 {
				poolDiff = 1.0
			}
			c.mu.Unlock()

			job, err := mining.ParseStratumJob(jobID, prevhashHex, coinb1Hex, coinb2Hex, versionHex, nbitsHex, ntimeHex, en1, en2, en2Size, merkleBranchHex, poolDiff)
			if err == nil && job != nil {
				c.mu.Lock()
				c.lastJob = job
				c.mu.Unlock()
				select {
				case c.JobCh <- *job:
				default:
				}
			}
		}
	}
}

func (c *StratumPoolClient) handleResponse(resp JSONRPCResponse) {
	c.mu.Lock()
	if ch, ok := c.pendingRequests[resp.ID]; ok {
		ch <- resp
		delete(c.pendingRequests, resp.ID)
	}
	c.mu.Unlock()

	// Parse extranonce1 from subscribe response
	if resp.Result != nil {
		var resultArr []interface{}
		if err := json.Unmarshal(resp.Result, &resultArr); err == nil && len(resultArr) >= 3 {
			// Typical subscribe response: [ subscriptions, extranonce1, extranonce2_size ]
			if en1, ok := resultArr[1].(string); ok {
				if en2Size, ok := resultArr[2].(float64); ok {
					c.mu.Lock()
					c.extranonce1 = en1
					c.extranonce2Size = int(en2Size)
					c.mu.Unlock()
				}
			}
		}
	}
}

// FetchStats is not natively implemented by Stratum (some pools support an ext).
func (c *StratumPoolClient) FetchStats(ctx context.Context) (PoolStatsMsg, error) {
	return PoolStatsMsg{}, nil
}

// SubmitShare submits mining results over TCP Stratum.
func (c *StratumPoolClient) SubmitShare(ctx context.Context, nonce uint32, hash [32]byte) (bool, error) {
	worker := c.BTCAddress
	if c.WorkerName != "" {
		worker = worker + "." + c.WorkerName
	}
	
	c.mu.Lock()
	job := c.lastJob
	c.mu.Unlock()

	if job == nil {
		return false, fmt.Errorf("no active job to submit")
	}

	nonceHexLE := fmt.Sprintf("%02x%02x%02x%02x", byte(nonce), byte(nonce>>8), byte(nonce>>16), byte(nonce>>24))

	resp, err := c.sendAndWait("mining.submit", []interface{}{worker, job.JobID, job.Extranonce2Hex, job.NtimeHex, nonceHexLE}, 10*time.Second)
	if err != nil {
		return false, err
	}
	if resp.Error != nil {
		return false, fmt.Errorf("pool rejected share: %v", resp.Error)
	}
	var resultBool bool
	if err := json.Unmarshal(resp.Result, &resultBool); err == nil && resultBool {
		return true, nil
	}
	return false, fmt.Errorf("pool returned unknown result: %s", string(resp.Result))
}

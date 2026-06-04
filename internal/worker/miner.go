package worker

import (
	"context"
	"encoding/hex"
	"log"
	"math/rand/v2"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nerdminertui/nerdtui/pkg/mining"
	"github.com/nerdminertui/nerdtui/pkg/trivia"
)

// MinerWorker manages the background hashing loop and CPU throttling.
type MinerWorker struct {
	client     PoolClient
	cpuTarget  atomic.Value // float64
	job        atomic.Value // mining.Job
	outCh      chan<- tea.Msg
	throttleCh <-chan float64
	jobCh      <-chan mining.Job
}

// NewMinerWorker initializes a new miner worker.
func NewMinerWorker(client PoolClient, initialCPUTarget float64, initialJob mining.Job, outCh chan<- tea.Msg, throttleCh <-chan float64, jobCh <-chan mining.Job) *MinerWorker {
	w := &MinerWorker{
		client:     client,
		outCh:      outCh,
		throttleCh: throttleCh,
		jobCh:      jobCh,
	}
	w.cpuTarget.Store(initialCPUTarget)
	w.job.Store(initialJob)
	return w
}

// Run executes the hashing loop. It blocks until the context is cancelled.
func (w *MinerWorker) Run(ctx context.Context) {
	numCPU := runtime.NumCPU()
	if numCPU < 1 {
		numCPU = 1
	}

	type localJob struct {
		mining.Job
		StartNonce uint32
	}
	var currentJob atomic.Value
	
	j, ok := w.job.Load().(mining.Job)
	if !ok {
		j = mining.Job{}
	}
	currentJob.Store(localJob{Job: j, StartNonce: rand.Uint32()})

	var totalHashes atomic.Uint64
	var totalWorkTime atomic.Int64
	
	var wg sync.WaitGroup
	wg.Add(numCPU)
	
	for i := 0; i < numCPU; i++ {
		go func(workerID int) {
			defer wg.Done()
			
			const BatchSize = 10000
			
			var lastJobID string
			var hashState *mining.MinerHashState
			var currentNonce uint32
			var currentSpaceWord string
			
			for {
				select {
				case <-ctx.Done():
					return
				default:
				}
				
				start := time.Now()
				lj := currentJob.Load().(localJob)
				
				if lj.JobID == "" {
					time.Sleep(100 * time.Millisecond)
					continue
				}
				
				if lj.JobID != lastJobID {
					lastJobID = lj.JobID
					
					// Generate a unique extra nonce 2 for this specific worker using astronomy trivia
					currentSpaceWord = trivia.GetRandomSpaceWordHex(lj.Extranonce2Size)
					if currentSpaceWord == "" {
						currentSpaceWord = lj.Extranonce2Hex // fallback to pool's default
					} else {
						// Inject worker ID into the last byte of the extranonce2 to guarantee
						// unique merkle roots across local workers (L4)
						if b, _ := hex.DecodeString(currentSpaceWord); len(b) > 0 {
							b[len(b)-1] = byte(workerID)
							currentSpaceWord = hex.EncodeToString(b)
						}
					}
					
					// Rebuild the header so this worker has a completely unique Merkle Root
					newHeader, err := mining.RebuildHeaderWithExtraNonce2(&lj.Job, currentSpaceWord)
					if err != nil {
						log.Printf("[Miner] Failed to rebuild header: %v", err)
						newHeader = lj.Header
						currentSpaceWord = lj.Extranonce2Hex
					}
					
					hashState = mining.NewMinerHashState(newHeader)
					currentNonce = rand.Uint32() // fully random start since roots are unique!
				}
				
				cpuTarget := w.cpuTarget.Load().(float64)
				
				// Perform batch
				for b := 0; b < BatchSize; b++ {
					hash := hashState.HashNonce(currentNonce)
					if mining.MeetsTarget(hash, lj.Target) {
						log.Printf("[Miner] Found valid share! JobID=%s, Word=%s, Nonce=%d, Hash=%x", lj.JobID, currentSpaceWord, currentNonce, hash)
						go func(sw string, n uint32, h [32]byte) {
							submitCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
							defer cancel()
							accepted, err := w.client.SubmitShare(submitCtx, sw, n, h)
							if err != nil {
								select {
								case w.outCh <- PoolErrorMsg{Err: err}:
								case <-ctx.Done():
								}
							} else {
								select {
								case w.outCh <- ShareFoundMsg{Accepted: accepted}:
								case <-ctx.Done():
								}
							}
						}(currentSpaceWord, currentNonce, hash)
					}
					currentNonce++ // Sequential increment is safe since Merkle Root is unique
				}
				
				workDur := time.Since(start)
				totalHashes.Add(BatchSize)
				totalWorkTime.Add(int64(workDur))
				
				if cpuTarget < 0.05 {
					cpuTarget = 0.05
				} else if cpuTarget > 1.0 {
					cpuTarget = 1.0
				}
				
				if cpuTarget < 1.0 {
					sleepDur := time.Duration(float64(workDur) * (1.0 - cpuTarget) / cpuTarget)
					time.Sleep(sleepDur)
				}
			}
		}(i)
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	lastTick := time.Now()
	for {
		select {
		case <-ctx.Done():
			wg.Wait()
			return
		case target := <-w.throttleCh:
			w.cpuTarget.Store(target)
		case newJob := <-w.jobCh:
			currentJob.Store(localJob{Job: newJob, StartNonce: rand.Uint32()})
		case <-ticker.C:
			now := time.Now()
			elapsed := now.Sub(lastTick).Seconds()
			lastTick = now
			
			hashes := totalHashes.Swap(0)
			wt := totalWorkTime.Swap(0)
			
			if hashes > 0 && elapsed > 0 {
				hps := float64(hashes) / elapsed
				cpuActual := (float64(wt) / float64(time.Second)) / float64(numCPU)
				if cpuActual > 1.0 {
					cpuActual = 1.0
				}
				
				select {
				case w.outCh <- HashRateMsg{HPS: hps, CPUActual: cpuActual}:
				case <-ctx.Done():
				default:
				}
			}
		}
	}
}

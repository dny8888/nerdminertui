package worker

import (
	"context"
	"sync/atomic"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nerdminertui/nerdtui/pkg/mining"
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
	const BatchSize = 50000

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	var currentNonce uint32
	var intervalHashes uint64
	var intervalWorkTime time.Duration
	var intervalSleepTime time.Duration

	for {
		select {
		case <-ctx.Done():
			return
		case target := <-w.throttleCh:
			w.cpuTarget.Store(target)
		case newJob := <-w.jobCh:
			w.job.Store(newJob)
		case <-ticker.C:
			// Emit metrics
			if intervalWorkTime+intervalSleepTime > 0 {
				hps := float64(intervalHashes)
				cpuActual := float64(intervalWorkTime) / float64(intervalWorkTime+intervalSleepTime)
				select {
				case w.outCh <- HashRateMsg{HPS: hps, CPUActual: cpuActual}:
				case <-ctx.Done():
					return
				}
			}
			intervalHashes = 0
			intervalWorkTime = 0
			intervalSleepTime = 0
		default:
			start := time.Now()

			j, ok := w.job.Load().(mining.Job)
			if !ok {
				// Fallback generic job if cast fails
				j = mining.Job{}
			}

			// Perform a batch of hashes
			for i := 0; i < BatchSize; i++ {
				hash := mining.HashHeader(j.Header, currentNonce)
				if mining.MeetsTarget(hash, j.Target) {
					// 10s timeout per requirement for Share submission
					go func(n uint32, h [32]byte) {
						submitCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
						defer cancel()
						accepted, err := w.client.SubmitShare(submitCtx, n, h)
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
					}(currentNonce, hash)
				}
				currentNonce++
			}

			workDur := time.Since(start)
			intervalHashes += BatchSize
			intervalWorkTime += workDur

			cpuTarget := w.cpuTarget.Load().(float64)
			if cpuTarget < 0.05 {
				cpuTarget = 0.05
			} else if cpuTarget > 1.0 {
				cpuTarget = 1.0
			}

			if cpuTarget < 1.0 {
				sleepDur := time.Duration(float64(workDur) * (1.0 - cpuTarget) / cpuTarget)
				time.Sleep(sleepDur)
				intervalSleepTime += sleepDur
			}
		}
	}
}

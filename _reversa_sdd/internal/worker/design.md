# internal/worker, Technical Design

> Design specification for the `internal/worker` module. Focuses on HOW the workers are designed.

## Interface

| Symbol | Signature | Return | Observation |
|---------|-----------|---------|------------|
| `MinerWorker.Run` | `func (w *MinerWorker) Run(ctx context.Context)` | `void` | Infinite hashing and throttling loop (blocks until cancelled). 🟢 |
| `HTTPPoolClient.FetchStats` | `func (c *HTTPPoolClient) FetchStats(ctx context.Context)` | `(PoolStats, error)` | Queries REST API for Screen 3 global network stats. 🟢 |
| `StratumPoolClient.SubmitShare` | `func (c *StratumPoolClient) SubmitShare(ctx context.Context, nonce uint32, hash [32]byte) (bool, error)` | `(bool, error)` | Submits share via JSON-RPC `mining.submit` over the persistent TCP socket. 🟢 |

## Main Flow
1. **Initiate Miner Run Loop**:
   - Ticker in main schedules execution of `MinerWorker.Run(ctx)` in separate goroutine. 🟢
2. **Execute Hashing Batch**:
   - Check context cancellation. If `ctx.Done()`, close channels and terminate. 🟢
   - Read active Job atomic value. 🟢
   - Start stopwatch (`time.Now()`). 🟢
   - Calculate 50,000 hashes in loop calling `pkg.mining.HashHeader`. 🟢
   - If hash is below Job target:
     - Invoke `PoolClient.SubmitShare` over the persistent Stratum TCP connection, blocking to await pool validation response up to 10 seconds. 🟢
     - Construct `ShareFoundMsg` with `Accepted` boolean status and write to `outCh`. 🟢
     - If socket write errors or timeout is triggered, emit `PoolErrorMsg` to statusbar without interrupting miner thread. 🟢
3. **Calculate Throttle Soneca**:
   - Measure time elapsed (`workDuration`). 🟢
   - Compute `sleepDuration = workDuration * (1.0 - CPUTarget) / CPUTarget`. 🟢
   - Execute `time.Sleep(sleepDuration)`. 🟢
   - Calculate `CPUActual = workDuration / (workDuration + sleepDuration)`. 🟢
4. **Emit metrics**:
   - 1s ticker in worker compiles `HashRateMsg` sending HPS and `CPUActual` to UI. 🟢

## Alternative Flows
- **Stratum TCP Connection Failure**:
  - `StratumPoolClient` TCP write throws socket error.
  - Return sentinela `ErrPoolUnreachable`.
  - Poller triggers retry exponentially (e.g., 2s, 4s, 8s, up to max limit). 🟢
- **Mock Mode Active (`MockMining=true`)**:
  - `MockPoolClient` replaces both HTTP and Stratum connections. `SubmitShare` returns `true` immediately without I/O. 🟢

## Dependencies
- `pkg/mining`: Cryptographic functions. 🟢
- `github.com/charmbracelet/bubbletea`: Sychronizes states. 🟢
- `context`: Propagates timeouts and cancellations. 🟢

## Design Decisions Identified

| Decision | Evidence | Confidence |
|---------|---------------------|-----------|
| Atomic Value Job | `miner.go:w.job.Load()` | 🟢 |
| Batch Size of 50K | `miner.go:const BatchSize = 50000` | 🟢 |
| Concurrent Clients | `cmd/tui/main.go:concurrent fetch/stratum` | 🟢 |
| 10s Submit Timeout | `fetcher.go:10s ctx timeout` | 🟢 |

## Internal State
- **MinerWorker**:
  - `throttleCh`: float64 receiver channel.
  - `outCh`: Bubbletea message dispatcher.
  - `job`: thread-safe atomic placeholder for incoming pool headers. 🟢

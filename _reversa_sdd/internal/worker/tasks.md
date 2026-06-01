# internal/worker, Implementation Tasks

> Implementation checklists for concurrency loops and throttling scheduler.

## Prerequisites
- [ ] Cryptographic algorithms in `pkg/mining` compile and pass tests. 🟢

## Tasks

- [ ] **T-01: Implement Messages definition**
  - Origin in Legacy: `nerdtui-spec.md:356` (§6)
  - Criteria of Done: Create `internal/worker/messages.go` declaring all typed `tea.Msg` structs.
  - Confidence: 🟢 CONFIRMADO

- [ ] **T-02: Implement Miner Worker Loop**
  - Origin in Legacy: `nerdtui-spec.md:234` (§5.5)
  - Criteria of Done: Write `MinerWorker` structure with channel synchronization and batch execution loop.
  - Confidence: 🟢 CONFIRMADO

- [ ] **T-03: Implement CPU limit scheduler logic**
  - Origin in Legacy: `nerdtui-spec.md:249` (§5.5)
  - Criteria of Done: Implement micro-sleep math calculations using time duration measurements.
  - Confidence: 🟢 CONFIRMADO

- [ ] **T-04: Implement Pool fetcher HTTP and Stratum Clients**
  - Origin in Legacy: `nerdtui-spec.md:212` (§5.4)
  - Criteria of Done: Complete clients supporting API endpoint parsing and socket handshakes.
  - Confidence: 🟢 CONFIRMADO

- [ ] **T-05: Implement Ticker Poller with Exponential Backoff**
  - Origin in Legacy: `nerdtui-spec.md:90` (§3)
  - Criteria of Done: Write reconnect poller executing retry intervals correctly.
  - Confidence: 🟢 CONFIRMADO

## Test Tasks

- [ ] **TT-01: Throttling mathematical accuracy**
  - Assert that sleep calculations are strictly correct at 25%, 50%, 75% and 100% CPUTarget inputs.
- [ ] **TT-02: Miner worker thread leaks check**
  - Verify `goleak` asserts no active goroutine leaks after miner context is cancelled.

# internal/worker

> Requirements specification for the `internal/worker` module. Focuses on WHAT the concurrent workers do, not how.

## Overview
The `internal/worker` module manages all background concurrent operations in the NerdTUI, executing the core SHA256d hashing loop with CPU limit controls, TCP Stratum socket handlers, and REST API pollers. 🟢

## Responsibilities
- Execute the hashing miner process in a safe background goroutine. 🟢
- Schedul micro-sleeps to maintain real-time CPU throttling limits. 🟢
- Handle TCP Stratum JSON-RPC connections and HTTP REST statistics queries. 🟢
- Sychronize background alerts and events to the Bubbletea main loop using typed messages (`tea.Msg`). 🟢

## Business Rules
- **BR-01: Hashing Cycle Throttling**: The hashing loop executes in batches of `50,000` hashes, measuring exact execution duration to schedule throttling sleep as:
  $$\text{sleep} = \text{workDuration} \times \frac{1 - P}{P}$$
  where $P$ is the target CPU percentage. 🟢
- **BR-02: Hashrate Discrepancy Alert**: The system alerts the operator if the absolute discrepancy between real and target CPU exceeds `0.05` ($5\%$). 🟢
- **BR-03: Stratum Reconnection Ticker**: Stratum connection losses trigger a reconnection poller with exponential backoff retry. 🟢
- **BR-04: Concurrent Network Clients**: When `MockMining` is false, `HTTPPoolClient` and `StratumPoolClient` run concurrently. `StratumPoolClient` maintains a persistent TCP socket connection for mining jobs, while `HTTPPoolClient` performs periodic REST polling exclusively for global network statistics displayed on Screen 3. When `MockMining` is true, a single `MockPoolClient` replaces both. 🟢
- **BR-05: Protocol Share Submission (REQ-WORKER-SUBMIT-01)**: When a valid share (satisfying the local target difficulty) is found, `StratumPoolClient` must transmit a `mining.submit` JSON-RPC message over the active TCP connection to the pool and wait for an accepted/rejected response. Submission timeout is set to 10 seconds. In `MockMining` mode, mock shares are accepted unconditionally without any network writes. 🟢

## Functional Requirements

| ID | Requirement | Priority | Acceptance Criteria |
|----|-----------|-----------|-------------------|
| RF-01 | Background Hashing Loop | Must | Execute continuous SHA256d rounds until context is cancelled. 🟢 |
| RF-02 | Channel Throttling Updates | Must | Update CPUTarget dynamically mid-run on receiving updates from `throttleCh`. 🟢 |
| RF-03 | Bubbletea Event Dispatch | Must | Emit `HashRateMsg` every 1s, and `ShareFoundMsg` immediately when a share passes the target. 🟢 |
| RF-04 | Stratum Share Submission (REQ-WORKER-SUBMIT-01) | Must | Execute JSON-RPC `mining.submit` over the TCP connection and process response within 10s timeout, updating `ShareFoundMsg` with acceptance status. 🟢 |

## Non-Functional Requirements

| Type | Inferred Requirement | Code Evidence | Confidence |
|------|--------------------|---------------------|-----------|
| Concurrency | Fully thread-safe variables via atomics/channels | `miner.go` | 🟢 |
| Availability | Exponential backoff reconnect poller | `poller.go` | 🟢 |

## Acceptance Criteria

```gherkin
Given a miner thread running with a CPUTarget of 0.5 (50%)
When a batch of hashes completes execution in 10ms
Then the thread sleeps for exactly 10ms before commencing the next batch
And emits a HashRateMsg showing ~50% actual CPU utilization

Given a Stratum connection failure
When the poller catches the TCP socket error
Then it drops connection state and triggers the exponential backoff reconnection timer

Given a miner thread that finds a valid share meeting the target
When the ShareFoundMsg is generated
Then StratumPoolClient sends a mining.submit JSON-RPC request to the pool
And blocks up to 10 seconds awaiting the pool's accepted/rejected confirmation
```

## Priority (MoSCoW)

| Requirement | MoSCoW | Justification |
|-----------|--------|---------------|
| Multi-thread Hashing engine | Must | Core feature of the mining dashboard. 🟢 |
| CPU limit scheduler | Must | Avoids device overheating, satisfying design constraint. 🟢 |
| TCP Share Submission | Must | Mandated by pools to maintain connection and record mining work. 🟢 |

## Code Traceability

| File | Function / Class | Coverage |
|---------|-----------------|-----------|
| `internal/worker/miner.go` | `MinerWorker.Run` | 🟢 |
| `internal/worker/fetcher.go` | `PoolClient.FetchStats` | 🟢 |
| `internal/worker/fetcher.go` | `PoolClient.SubmitShare` | 🟢 |
| `internal/worker/poller.go` | `PollCmd` | 🟢 |
| `internal/worker/messages.go` | `HashRateMsg` | 🟢 |

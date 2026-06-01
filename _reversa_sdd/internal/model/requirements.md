# internal/model

> Requirements specification for the `internal/model` module. Focuses on WHAT the data model guarantees, not how.

## Overview
The `internal/model` module defines the immutable global state `AppState` and constants for the NerdTUI, acting as the thread-safe state container of the Elm UI loop. 🟢

## Responsibilities
- Define structural state fields to render screens and components. 🟢
- Guarantee thread-safety through immutability conventions. 🟢
- Define global domain bounds for TUI screens, limits, and histories. 🟢

## Business Rules
- **BR-01: Structural Immutability**: AppState values must not contain mutable pointers to internal fields that could trigger data races. 🟢
- **BR-02: FIFO Hashrate History**: The hashrate history is a fixed-size `[60]float64` array acting as a FIFO queue, rotating metrics cleanly upon every new update interval. 🟢

## Functional Requirements

| ID | Requirement | Priority | Acceptance Criteria |
|----|-----------|-----------|-------------------|
| RF-01 | State Cloning Updates | Must | Provide a pure method to shift hashrate values and return a modified copy of the state. 🟢 |
| RF-02 | Screen Identifiers | Must | Support enum identification for exactly 3 distinct TUI views (Dashboard, Clock, Global Stats). 🟢 |

## Non-Functional Requirements

| Type | Inferred Requirement | Code Evidence | Confidence |
|------|--------------------|---------------------|-----------|
| Safety | 100% thread safety without CPU mutex overhead | `state.go` | 🟢 |
| Memory | Very small allocations during copying | `state.go` | 🟢 |

## Acceptance Criteria

```gherkin
Given an active AppState with an empty hashrate history
When WithHashRate(5000.0) is called
Then the returned copy has HashRate set to 5000.0, and the first index in the array is 5000.0
And the original receiver state remains unchanged (asserting immutability)
```

## Priority (MoSCoW)

| Requirement | MoSCoW | Justification |
|-----------|--------|---------------|
| AppState Immutability | Must | Vital core design constraint preventing data races in TUI threads. 🟢 |
| Hashrate FIFO circular array | Must | Essential data container to feed the sparkline component visualizer. 🟢 |

## Code Traceability

| File | Function / Class | Coverage |
|---------|-----------------|-----------|
| `internal/model/state.go` | `AppState` | 🟢 |
| `internal/model/state.go` | `WithHashRate` | 🟢 |

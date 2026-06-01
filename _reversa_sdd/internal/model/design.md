# internal/model, Technical Design

> Design specification for the `internal/model` module. Focuses on HOW the state model is structured.

## Interface

### Classes / Functions

| Symbol | Signature | Retorno | Observação |
|---------|-----------|---------|------------|
| `WithHashRate` | `func (s AppState) WithHashRate(hps float64) AppState` | `AppState` | Returns a modified state copy with HPS and rolled array. 🟢 |

## Main Flow
1. **Define state.go Consts**:
   - `ScreenID` as `uint8` with:
     - `ScreenDashboard = 0`
     - `ScreenClock = 1`
     - `ScreenGlobalStats = 2`
   - `NumScreens = 3`, `MinCPUTarget = 0.05`, `MaxCPUTarget = 1.00`, `CPUStep = 0.05`, `HashHistoryLen = 60` 🟢
2. **Implement WithHashRate**:
   - Create a copy of `AppState`.
   - Assign new `HashRate`.
   - Shift array values to the left by 1 element, placing `hps` at index `59` of `HashRateHistory`. 🟢
   - Return copy. 🟢

## Dependencies
- None. This package has zero external and internal package dependencies. 🟢

## Design Decisions Identified

| Decision | Evidence | Confidence |
|---------|---------------------|-----------|
| Value Receiver | `state.go:func (s AppState)` (no pointer `*`) | 🟢 |
| Fixed Circular FIFO | `state.go:[60]float64` array type | 🟢 |

## Internal State
- Although this is the state model, the package defines only the structure and operations. The actual state instance resides on Bubbletea's app model in the UI layer. 🟢

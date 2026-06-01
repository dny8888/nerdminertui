# internal/model, Implementation Tasks

> Implementation checklists for data struct and constant bounds.

## Prerequisites
- [ ] Go SDK standard library is set up. 🟢

## Tasks

- [ ] **T-01: Define State Constants**
  - Origin in Legacy: `nerdtui-spec.md:200` (§5.3)
  - Criteria of Done: Declare `ScreenID` enums, `HashHistoryLen = 60`, and CPU bounds.
  - Confidence: 🟢 CONFIRMADO

- [ ] **T-02: Implement AppState Struct**
  - Origin in Legacy: `nerdtui-spec.md:128` (§4)
  - Criteria of Done: Create `AppState` structure containing only value types (no reference pointers).
  - Confidence: 🟢 CONFIRMADO

- [ ] **T-03: Implement WithHashRate Method**
  - Origin in Legacy: `nerdtui-spec.md:208` (§5.3)
  - Criteria of Done: Write `WithHashRate` on `AppState` shifting the hashrate history FIFO array without modifying the receiver.
  - Confidence: 🟢 CONFIRMADO

## Test Tasks

- [ ] **TT-01: Verify AppState Immutability**
  - Call `WithHashRate` on an instance and verify the original instance's history array remains unchanged.
- [ ] **TT-02: Verify FIFO FIFO circular rotation**
  - Add 61 items consecutively and assert older values roll off and history maintains last 60 entries.

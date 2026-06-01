# cmd/tui, Implementation Tasks

> Task checklist to implement the `cmd/tui` module based on spec descriptions.

## Prerequisites
- [ ] Dependencies `internal/config`, `internal/store`, `internal/worker`, and `internal/ui` are compiled and fully tested. 🟢
- [ ] Standard terminal UNIX libraries (`flag`, `os`, `log`) are available in the local Go SDK runtime. 🟢

## Tasks

- [ ] **T-01: Parse Command-Line Flags**
  - Origin in Legacy: `nerdtui-spec.md:165` (§5.1)
  - Criteria of Done: Run bin with `--config`, `--no-mine`, `--cpu` and ensure they are parsed with correct defaults.
  - Confidence: 🟢 CONFIRMADO

- [ ] **T-02: Load and Validate Configuration**
  - Origin in Legacy: `nerdtui-spec.md:166` (§5.1)
  - Criteria of Done: Trigger validator and assert program logs fatal and exits if validation fails.
  - Confidence: 🟢 CONFIRMADO

- [ ] **T-03: Initialize Persistence Handler**
  - Origin in Legacy: `nerdtui-spec.md:167` (§5.1)
  - Criteria of Done: Verify the SQLite file `metrics.db` is correctly instantiated in store folder.
  - Confidence: 🟢 CONFIRMADO

- [ ] **T-04: Worker and Channel Instantiation**
  - Origin in Legacy: `nerdtui-spec.md:168` (§5.1)
  - Criteria of Done: Start channels and wire miner threads with valid context cancellation loops.
  - Confidence: 🟢 CONFIRMADO

- [ ] **T-05: Launch Bubbletea Elm Loop**
  - Origin in Legacy: `nerdtui-spec.md:170` (§5.1)
  - Criteria of Done: Execute the app program using `tea.WithAltscreen()` and confirm terminal switches context smoothly.
  - Confidence: 🟢 CONFIRMADO

## Test Tasks

- [ ] **TT-01: Smoke test binary execution**
  - Verify program compiles static (`CGO_ENABLED=0`) and launches successfully under `--mock` flag.
- [ ] **TT-02: Assert fatal exit on wrong configs**
  - Verify loading returns fatal exit code when BTC address is missing and mock mining is inactive.

## Suggested Order
1. Parse flags (T-01) and wire configs validation (T-02) first.
2. Initialize SQLite mock/real store handlers (T-03).
3. Connect concurrency workers channels (T-04) before executing the UI program (T-05).

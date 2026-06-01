# cmd/tui, Technical Design

> Design specification for the `cmd/tui` module. Focuses on HOW the unit is constructed based on the target specifications.

## Interface

### Entry Point
- **Symbol**: `main`
- **Signature**: `func main()`
- **Return**: `void`
- **Behavior**: Thin bootstrapping routine. Performs error logging and exits on failures.

## Main Flow
1. **Initialize CLI Flags**: Parse program inputs (`--config`, `--no-mine`, `--cpu`) using standard Go `flag` or `pflag` package. 🟢
2. **Load Configuration**: Invoke `config.Load()` to load active configs. If validation fails, call `log.Fatalf` with the validation error. 🟢
3. **Instantiate Storage**: Initialize the SQLite persistent store `store.New(cfg.StorePath)`. If error occurs, fail fatally. 🟢
4. **Wire workers and channels**:
   - Create channels: `throttleCh` (chan float64) and `outCh` (chan tea.Msg). 🟢
   - Initialize `MinerWorker` injecting channels. 🟢
5. **Create AppModel**: Wire `store`, `throttleCh`, and initial `AppState` into Bubbletea `AppModel`. 🟢
6. **Execute Bubbletea**:
   - Run `tea.NewProgram(appModel, tea.WithAltscreen()).Run()`. 🟢
   - Gracefully handle terminal buffer teardown on stop/interrupt. 🟢

## Alternative Flows
- **Mock Mode Active (`--mock` / `cfg.MockMining = true`)**:
  - Main wires the `MockPoolClient` instead of the physical `StratumPoolClient`/`HTTPPoolClient`. 🟢
  - BTCAddress validation is bypassed. 🟢
- **No Store Mode Active (`--no-store` / `cfg.StorePath = ""`)**:
  - Main wires `NilStore` to bypass SQLite metrics writes. 🟢

## Dependencies
- `internal/config`: Loads Viper-based config properties. 🟢
- `internal/store`: Provides persistent database queries. 🟢
- `internal/worker`: Starts MinerWorker and Stratum Fetchers. 🟢
- `internal/ui`: Provides Elm program loop app container. 🟢

## Design Decisions Identified

| Decision | Evidence | Confidence |
|---------|---------------------|-----------|
| Alternative Screen | `nerdtui-spec.md:tea.WithAltscreen()` | 🟢 |
| No-Panic Runtime | `nerdtui-spec.md:log.Fatal only during init` | 🟢 |

## Internal State
- This module is stateless. It only manages bootstrapping local variables on stack memory during the startup phase. 🟢

## Observability
- Emits fatal messages to standard error (`os.Stderr`) on startup failures before exit. 🟢
- No additional metrics are logged here to keep stdout clean. 🟢

## Risks and Gaps
- None. The bootstrapping routine is fully deterministic and specified. 🟢

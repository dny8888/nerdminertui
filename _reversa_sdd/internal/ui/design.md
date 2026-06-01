# internal/ui, Technical Design

> Design specification for the `internal/ui` module. Focuses on HOW the terminal UI is designed.

## Interface

### Classes / Functions

| Symbol | Signature | Return | Observation |
|---------|-----------|---------|------------|
| `AppModel.Init` | `func (m *AppModel) Init()` | `tea.Cmd` | Initiates the first poll and starts the 1s uptime ticker. 🟢 |
| `AppModel.Update` | `func (m *AppModel) Update(msg tea.Msg)` | `(tea.Model, tea.Cmd)` | Elm state machine update cycle. 🟢 |
| `AppModel.View` | `func (m *AppModel) View()` | `string` | Composites the screens and rodapé statusbar to terminal string. 🟢 |

## Main Flow
1. **Initialize AppModel (`Init`)**:
   - Return standard ticket commands to begin counting Uptime. 🟢
2. **Dispatch Messages (`Update`)**:
   - `worker.HashRateMsg`: Copy state, calculate and rotate hashrate history circular array. 🟢
   - `worker.ShareFoundMsg`: Increment `SharesFound`, update `BestDifficulty`. 🟢
   - `worker.PoolStatsMsg`: Update `BlockHeight`, check network state connection bullet status. 🟢
   - `worker.MinerErrorMsg`: Write error message to the StatusBar state. 🟢
   - `tea.KeyMsg`:
     - `tab` / `→`: Rotate active `Screen` by incrementing index modulus `NumScreens`. 🟢
     - `+` / `=`: `CPUTarget = min(1.0, CPUTarget + 0.05)`, dispatch on `throttleCh`. 🟢
     - `-` / `_`: `CPUTarget = max(0.05, CPUTarget - 0.05)`, dispatch on `throttleCh`. 🟢
     - `q` / `ctrl+c`: Emit `tea.Quit`. 🟢
3. **Composite Terminal Display (`View`)**:
   - Query active `Screen` enum. Call pure screen renderers:
     - `0`: `screens.RenderDashboard(s, width, height)`
     - `1`: `screens.RenderClock(s, width, height)`
     - `2`: `screens.RenderGlobalStats(s, width, height)` 🟢
   - Embed `components.RenderStatusBar(s, width)` to the bottom line of terminal display. 🟢

## Dependencies
- `github.com/charmbracelet/bubbletea`: State machine loop. 🟢
- `github.com/charmbracelet/lipgloss`: CSS-like terminal styles. 🟢
- `internal/model`: State structures. 🟢
- `pkg/format`: Unit format strings. 🟢

## Design Decisions Identified

| Decision | Evidence | Confidence |
|---------|---------------------|-----------|
| Stateless Screens | `screens/*.go: Render*(state AppState)` | 🟢 |
| Adaptive Resizing | `app.go: tea.WindowSizeMsg handler` | 🟢 |

## Internal State
- **AppModel**:
  - `state`: Active `AppState` struct (copied by value).
  - `width`, `height`: Viewport boundary properties. 🟢

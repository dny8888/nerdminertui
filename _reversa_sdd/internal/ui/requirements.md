# internal/ui

> Requirements specification for the `internal/ui` module. Focuses on WHAT the terminal view does, not how.

## Overview
The `internal/ui` module constructs the graphical interface of the NerdTUI, utilizing standard Bubbletea loops and Lipgloss styling to display live charts, CPU statistics, and clock views. 🟢

## Responsibilities
- Render three distinct screens (Dashboard, Clock, and Global Stats) in the terminal. 🟢
- Consume user keyboard interactions to adjust settings and rotate screens. 🟢
- Handle display scaling and resizing events safely without causing panic. 🟢
- Guarantee that all view renderers are mathematical pure functions. 🟢

## Business Rules
- **BR-01: Terminal View Rotation**: The `tab` or `→` keyboard inputs must cyclically rotate the view:
  $$\text{Dashboard (0)} \longrightarrow \text{Clock (1)} \longrightarrow \text{Global Stats (2)} \longrightarrow \text{Dashboard (0)}$$ 🟢
- **BR-02: User Throttling Command**: The `+` and `-` inputs adjust `CPUTarget` dynamically by a step value of `0.05` ($5\%$), clamping values inside `[0.05, 1.0]`. 🟢
- **BR-03: UI Purity**: All components (sparkline, gauge, statusbar) are designed as pure functional blocks that consume local properties and output a single stylized string. 🟢

## Functional Requirements

| ID | Requirement | Priority | Acceptance Criteria |
|----|-----------|-----------|-------------------|
| RF-01 | live screen rendering | Must | Support drawing dashboard, clock, and global stats layout. 🟢 |
| RF-02 | Key Mapping | Must | Map keyboard events `tab`, `+`, `-`, `q` and `ctrl+c` to respective messages or updates. 🟢 |
| RF-03 | Resizing Support | Must | Adapt all screen and component widths to standard `tea.WindowSizeMsg` dimensions. 🟢 |

## Non-Functional Requirements

| Type | Inferred Requirement | Code Evidence | Confidence |
|------|--------------------|---------------------|-----------|
| Performance | Render views in < 16ms to avoid screen flickering | `app.go` | 🟢 |
| Styling | No generic primary colors (Harmonious Theme) | `nerdtui-spec.md:Theme` | 🟢 |

## Acceptance Criteria

```gherkin
Given a terminal of size 80x24 showing the Dashboard Screen
When the user presses the 'tab' key
Then the state transitions to ScreenClock (1) and the view renders the large ASCII clock

Given the user presses the '+' key with a CPUTarget of 0.95
When the update cycle evaluates
Then CPUTarget clamps to exactly 1.00 (MaxCPUTarget) and launches a throttle update
```

## Priority (MoSCoW)

| Requirement | MoSCoW | Justification |
|-----------|--------|---------------|
| AppModel Elm cycle | Must | Crucial loop driving all updates, keyboard bindings, and views. 🟢 |
| Pure Screen renderers | Must | Vital to prevent rendering inconsistencies or memory leaks. 🟢 |

## Code Traceability

| File | Function / Class | Coverage |
|---------|-----------------|-----------|
| `internal/ui/app.go` | `AppModel.Update` | 🟢 |
| `internal/ui/app.go` | `AppModel.View` | 🟢 |
| `internal/ui/screens/dashboard.go` | `RenderDashboard` | 🟢 |
| `internal/ui/components/cpubar.go` | `RenderCPUBar` | 🟢 |

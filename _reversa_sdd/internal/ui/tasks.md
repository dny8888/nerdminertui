# internal/ui, Implementation Tasks

> Implementation checklists for Elm app model cycles and screen layouts.

## Prerequisites
- [ ] Bubbletea, lipgloss, and bubbles library imports are registered. 🟢
- [ ] Mappings in `internal/model` and `pkg/format` compile and test successfully. 🟢

## Tasks

- [ ] **T-01: Setup AppModel Elm Structure**
  - Origin in Legacy: `nerdtui-spec.md:264` (§5.6)
  - Criteria of Done: Create `AppModel` satisfying Bubbletea `tea.Model` interface contract methods.
  - Confidence: 🟢 CONFIRMADO

- [ ] **T-02: Implement Keyboard Bindings mappings**
  - Origin in Legacy: `nerdtui-spec.md:282` (§5.6)
  - Criteria of Done: Map all key inputs (`tab`, `+`, `-`, `q`, `ctrl+c`) in `keys.go` or `Update()`.
  - Confidence: 🟢 CONFIRMADO

- [ ] **T-03: Implement Pure Screens and Components**
  - Origin in Legacy: `nerdtui-spec.md:291` (§5.6)
  - Criteria of Done: Create functional, stateless renderers in `screens/` and `components/` folders.
  - Confidence: 🟢 CONFIRMADO

- [ ] **T-04: Implement Adaptive Terminal Viewport resizing**
  - Origin in Legacy: `nerdtui-spec.md:286` (§5.6)
  - Criteria of Done: Write `tea.WindowSizeMsg` dispatcher resizing layout dynamically without clipping widgets.
  - Confidence: 🟢 CONFIRMADO

## Test Tasks

- [ ] **TT-01: Screens Purity verification**
  - Assert that multiple calls to screen renderers with exact identical states return identical strings.
- [ ] **TT-02: Bubbletea Interactive Loop teatest verification**
  - Run teatest loop, click `tab` key 3 times, and verify active view rotates in a full circle.

# pkg/format, Implementation Tasks

> Implementation checklists for display unit formatters.

## Prerequisites
- [ ] Go SDK formatting standard packages are configured. 🟢

## Tasks

- [ ] **T-01: Implement FormatHashRate conversion**
  - Origin in Legacy: `nerdtui-spec.md:349` (§5.9)
  - Criteria of Done: Write `FormatHashRate` converting float inputs to correct scaled sufixes.
  - Confidence: 🟢 CONFIRMADO

- [ ] **T-02: Implement FormatUptime conversion**
  - Origin in Legacy: `nerdtui-spec.md:350` (§5.9)
  - Criteria of Done: Write `FormatUptime` converting durations into structured time metrics.
  - Confidence: 🟢 CONFIRMADO

- [ ] **T-03: Implement FormatBlockHeight parsing**
  - Origin in Legacy: `nerdtui-spec.md:352` (§5.9)
  - Criteria of Done: Write block parser inserting dot separators.
  - Confidence: 🟢 CONFIRMADO

- [ ] **T-04: Implement FormatDifficulty parsing**
  - Origin in Legacy: `nerdtui-spec.md:351` (§5.9)
  - Criteria of Done: Complete difficulty string formatting statement.
  - Confidence: 🟢 CONFIRMADO

## Test Tasks

- [ ] **TT-01: Hashrate units scaling checks**
  - Verify boundary assertions: `999.0` -> `"999 H/s"`, `1000.0` -> `"1.0 KH/s"`, `1500000.0` -> `"1.5 MH/s"`.
- [ ] **TT-02: Uptime durations parsing checks**
  - Assert that durations under 1 minute display `"0m XXs"`, and durations over a day display `"Xd XXh XXm"`.

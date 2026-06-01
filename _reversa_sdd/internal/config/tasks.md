# internal/config, Implementation Tasks

> Checklists to implement the config parsing and validation.

## Prerequisites
- [ ] Go viper dependency `github.com/spf13/viper` is added to `go.mod`. 🟢

## Tasks

- [ ] **T-01: Define Config Struct**
  - Origin in Legacy: `nerdtui-spec.md:180` (§5.2)
  - Criteria of Done: Define `Config` structure with required fields matching environment overrides.
  - Confidence: 🟢 CONFIRMADO

- [ ] **T-02: Setup Viper and Defaults**
  - Origin in Legacy: `nerdtui-spec.md:181` (§5.2)
  - Criteria of Done: Map defaults in code and enable env bind prefixed with `NM_`.
  - Confidence: 🟢 CONFIRMADO

- [ ] **T-03: Implement Validation Routine**
  - Origin in Legacy: `nerdtui-spec.md:191` (§5.2)
  - Criteria of Done: Implement `Config.Validate() error` throwing explicit bounds errors on target checks.
  - Confidence: 🟢 CONFIRMADO

## Test Tasks

- [ ] **TT-01: Config boundary tests**
  - Verify CPUTarget values less than 0.05 or greater than 1.0 return non-nil error.
- [ ] **TT-02: Missing Address checks**
  - Verify `Validate` throws a missing BTC address error when `MockMining` is false.

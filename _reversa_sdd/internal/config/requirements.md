# internal/config

> Requirements specification for the `internal/config` module. Focuses on WHAT the config management does, not how.

## Overview
The `internal/config` module is responsible for loading, parsing, overriding, and validating the NerdTUI execution parameters using environment variables and configuration files. 🟢

## Responsibilities
- Parse configurations from active environment variables or target files. 🟢
- Populate and structure the execution `Config` struct. 🟢
- Perform strict domain parameter validations before handing config over to downstream dependencies. 🟢

## Business Rules
- **BR-01: BTC Address Requirement**: In real mining mode (`MockMining == false`), a valid, non-empty Bitcoin Address (`BTCAddress`) must be provided. 🟢
- **BR-02: CPU Throttling Bounds**: The target CPU utilization (`CPUTarget`) must fall strictly within the inclusive boundary range of `[0.05, 1.00]`. 🟢
- **BR-03: Variable Overrides**: Environment variables prefixed with `NM_` override equivalent settings defined in static configuration files. 🟢
- **BR-04: Tilde Path Expansion (REQ-CONFIG-PATH-01)**: If any path-typed field starts with a tilde (`~/`), it must be resolved dynamically by replacing the tilde with the user's home directory retrieved via `os.UserHomeDir()`. Any failure to fetch the home directory must result in a fatal initialization error during `config.Load()`. Relative paths without `~` pass through unchanged and are documented as unsupported for production use. 🟢

## Functional Requirements

| ID | Requirement | Priority | Acceptance Criteria |
|----|-----------|-----------|-------------------|
| RF-01 | Environment Parsing | Must | Bind environment variables like `NM_POOL_URL` and `NM_CPU_TARGET` automatically to Viper. 🟢 |
| RF-02 | Boundary Validations | Must | Assert `CPUTarget` validates within `[0.05, 1.0]` constraints. 🟢 |
| RF-03 | Missing Config Defaults | Must | Fallback to sensible defaults (e.g. `public-pool.io` for Pool URL, `0.5` for CPUTarget) when fields are omitted. 🟢 |
| RF-04 | Tilde Expansion (REQ-CONFIG-PATH-01) | Must | Resolve leading `~/` to active OS home directory and fail fatally if user directory is inaccessible. 🟢 |

## Non-Functional Requirements

| Type | Inferred Requirement | Code Evidence | Confidence |
|------|--------------------|---------------------|-----------|
| Performance | Fast loading without network queries | `config.go` | 🟢 |
| Security | Secure environment mapping | `config.go` | 🟢 |

## Acceptance Criteria

```gherkin
Given a configuration file with CPUTarget set to 0.75 and BTCAddress filled
When c.Validate() is evaluated
Then it returns nil (validation succeeds)

Given a configuration file with CPUTarget set to 0.02
When c.Validate() is evaluated
Then it returns an out of bounds validation error

Given a configured storage path starting with "~/.nerdtui/"
When config.Load() is executed
Then the path is expanded dynamically to the active user's absolute home directory
```

## Priority (MoSCoW)

| Requirement | MoSCoW | Justification |
|-----------|--------|---------------|
| Config Struct Validation | Must | Ensures downstream workers do not run with corrupt or unsafe properties. 🟢 |
| Environment Overrides | Must | Highly critical for containerized or automated environment launches. 🟢 |
| Tilde Path Expansion | Must | Prevents runtime storage file creation errors and ensures system portability. 🟢 |

## Code Traceability

| File | Function / Class | Coverage |
|---------|-----------------|-----------|
| `internal/config/config.go` | `Config.Validate` | 🟢 |
| `internal/config/config.go` | `Load` | 🟢 |
| `internal/config/paths.go` | `ExpandPath` | 🟢 |

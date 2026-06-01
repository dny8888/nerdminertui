# pkg/format

> Requirements specification for the `pkg/format` module. Focuses on WHAT the formatting utilities do, not how.

## Overview
The `pkg/format` module provides pure utility functions to convert raw floating-point and integer metrics into human-readable strings structured for terminal display layouts. 🟢

## Responsibilities
- Convert raw HPS numbers into scaled hashrate units (e.g. H/s, KH/s, MH/s). 🟢
- Format runtime time durations into unified Uptime displays. 🟢
- Format large block numbers using clear digit separators. 🟢
- Format difficulty ratios into readable scientific notation formats. 🟢

## Business Rules
- **BR-01: Zero/NaN Protection**: Formatting zero hashrates or timestamps must never produce NaN or infinite values, returning clean default outputs (e.g. `"0 H/s"`). 🟢
- **BR-02: Strict Mathematical Purity**: All formatting routines must be stateless, having zero internal states or side effects. 🟢

## Functional Requirements

| ID | Requirement | Priority | Acceptance Criteria |
|----|-----------|-----------|-------------------|
| RF-01 | Hashrate Suffix Scaling | Must | Convert floats into rounded, single-decimal scaled values with unit suffix. 🟢 |
| RF-02 | Uptime String Compilation | Must | Compile time durations into `"Xd Xh Xm"` or `"Xm Xs"` strings. 🟢 |
| RF-03 | Thousand Separators | Must | Format block heights inserting dot separators (e.g. `892441` -> `"#892.441"`). 🟢 |

## Non-Functional Requirements

| Type | Inferred Requirement | Code Evidence | Confidence |
|------|--------------------|---------------------|-----------|
| Performance | Execute formatting under ~50ns to prevent display latency | `hashrate.go` | 🟢 |
| Safety | 100% thread safety (pure functions) | `duration.go` | 🟢 |

## Acceptance Criteria

```gherkin
Given a raw hashrate float of 1520.0
When FormatHashRate(1520.0) is evaluated
Then it returns exactly "1.5 KH/s"

Given a time duration of 2 days, 3 hours, and 14 minutes
When FormatUptime(duration) is called
Then it returns exactly "2d 03h 14m"
```

## Priority (MoSCoW)

| Requirement | MoSCoW | Justification |
|-----------|--------|---------------|
| Hashrate unit scaling | Must | Critical for displaying readable values on the dashboard gauge. 🟢 |
| Uptime duration formatting | Must | Essential string formatting for the Clock and Dashboard rodapé. 🟢 |

## Code Traceability

| File | Function / Class | Coverage |
|---------|-----------------|-----------|
| `pkg/format/hashrate.go` | `FormatHashRate` | 🟢 |
| `pkg/format/duration.go` | `FormatUptime` | 🟢 |
| `pkg/format/difficulty.go` | `FormatDifficulty` | 🟢 |

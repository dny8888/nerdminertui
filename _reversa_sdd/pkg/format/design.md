# pkg/format, Technical Design

> Design specification for the `pkg/format` module. Focuses on HOW the formats are structured.

## Interface

### Classes / Functions

| Symbol | Signature | Return | Observation |
|---------|-----------|---------|------------|
| `FormatHashRate` | `func FormatHashRate(hps float64) string` | `string` | Converts hashrate float to formatted scale. 🟢 |
| `FormatUptime` | `func FormatUptime(d time.Duration) string` | `string` | Formats uptime to string output. 🟢 |
| `FormatDifficulty` | `func FormatDifficulty(d float64) string` | `string` | Scientific notation string formats. 🟢 |
| `FormatBlockHeight` | `func FormatBlockHeight(h uint32) string` | `string` | Dots separation block display string. 🟢 |

## Main Flow
1. **Format Hashrate (`FormatHashRate`)**:
   - If `hps == 0`, return `"0 H/s"`. 🟢
   - If `hps >= 1000000.0`, return float divided by 1,000,000 formatted with one decimal and sufix `" MH/s"`. 🟢
   - If `hps >= 1000.0`, return float divided by 1,000 formatted with one decimal and sufix `" KH/s"`. 🟢
   - Else, return float formatted to integer and sufix `" H/s"`. 🟢
2. **Format Uptime (`FormatUptime`)**:
   - Extract elapsed days, hours, and minutes from the time duration. 🟢
   - If `days > 0`, return string in the format `"Xd XXh XXm"`. 🟢
   - Else, return string in the format `"Xm XXs"`. 🟢
3. **Format Block Height (`FormatBlockHeight`)**:
   - Construct string using standard print formats inserting a dot separator at thousands thresholds (e.g. `"#892.441"`). 🟢

## Dependencies
- `time`: Native duration utility. 🟢
- `fmt`: Standard string compilation. 🟢

## Design Decisions Identified

| Decision | Evidence | Confidence |
|---------|---------------------|-----------|
| Stateless Purity | Functions have zero local configurations or shared mutables. | 🟢 |
| Floating point safety | Explicit checks for zero values to block Nan outputs. | 🟢 |

## Internal State
- This package is completely stateless. 🟢

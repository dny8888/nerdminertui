# cmd/tui

> Requirements specification for the `cmd/tui` module. Focuses on WHAT the unit does, not how.

## Overview
The `cmd/tui` module serves as the exclusive entry point and dependency injector (wiring) for the NerdTUI application, parsing command-line parameters and initiating the Bubbletea program runtime. 🟢

## Responsibilities
- Parse incoming command line flags correctly. 🟢
- Instantiate global dependencies (Config load, SQLite metrics store, Worker initialization). 🟢
- Launch the main Bubbletea application loop using an alternate terminal buffer (`AltScreen`). 🟢

## Business Rules
- **BR-01: Alternate Screen Execution**: The application must run in an alternate terminal screen buffer to prevent visual pollution of the operator's active shell history. 🟢
- **BR-02: Exclusive Wiring**: No domain or business logic (like hashing calculations, UI sizing math, or DB queries) is allowed in this package. It must remain a thin coordinator shell. 🟢

## Functional Requirements

| ID | Requirement | Priority | Acceptance Criteria |
|----|-----------|-----------|-------------------|
| RF-01 | CLI Flags Parsing | Must | Parse `--config`, `--no-mine`, and `--cpu` values on startup. 🟢 |
| RF-02 | AltScreen Bootstrapping | Must | Execute Bubbletea program in AltScreen mode, handling exit codes gracefully. 🟢 |
| RF-03 | Fatal Error Handover | Must | Terminate with non-zero exit code if fatal loading issues occur during startup. 🟢 |

## Non-Functional Requirements

| Type | Inferred Requirement | Code Evidence | Confidence |
|------|--------------------|---------------------|-----------|
| Portability | Zero GCC Dependency (`CGO_ENABLED=0`) | `Makefile:872` | 🟢 |
| Safety | No Panic on Runtime execution | `nerdtui-spec.md:R6` | 🟢 |

## Acceptance Criteria

```gherkin
Given the system is compiled and ready to execute
When I launch the bin/nerdtui binary with a valid BTC address config
Then the terminal turns into an altscreen displaying the mining stats dashboard

Given the system is launched without required configurations and mock mode is off
When the program initializes
Then it writes a fatal error to log and exits with code 1
```

## Priority (MoSCoW)

| Requirement | MoSCoW | Justification |
|-----------|--------|---------------|
| Program Bootstrapping | Must | Entry point of execution, nothing runs without it. 🟢 |
| CLI Flag Mapping | Must | Handles configuration parameters required for other components. 🟢 |

## Code Traceability

| File | Function / Class | Coverage |
|---------|-----------------|-----------|
| `cmd/tui/main.go` | `main` | 🟢 |
| `Makefile` | `build` | 🟢 |

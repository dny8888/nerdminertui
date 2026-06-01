# internal/store

> Requirements specification for the `internal/store` module. Focuses on WHAT the store layer does, not how.

## Overview
The `internal/store` module provides persistent SQL-based storage for the NerdTUI hashrate records, utilizing pure Go SQLite bindings to avoid CGO runtime compilation requirements. 🟢

## Responsibilities
- Open and close local metrics database connections cleanly. 🟢
- Append HPS samples to local storage. 🟢
- Query the most recent hashrate entries (up to limit of 60) to feed sparklines. 🟢
- Provide a silent no-op client placeholder (`NilStore`) when storage is deactivated. 🟢

## Business Rules
- **BR-01: WAL Performance Mode**: The SQLite connection must establish Write-Ahead Logging (WAL) pragmas to allow highly concurrent, non-blocking metrics updates from background threads. 🟢
- **BR-02: Silent No-Op Fallback**: If standard local storage deactivation `--no-store` is mapped on start, the application routes calls to `NilStore`, returning nil values without throwing errors. 🟢
- **BR-03: Transient Statistics Scope (REQ-STORE-SCOPE-01)**: The database schema only persists hashrate metrics (`hashrate_history`). Cumulative per-session variables, namely `SharesFound` and `BestDifficulty`, are strictly transient and reset to zero on every restart. They are not stored in SQLite. 🟢

## Functional Requirements

| ID | Requirement | Priority | Acceptance Criteria |
|----|-----------|-----------|-------------------|
| RF-01 | SQLite SQL Migrations | Must | Embed SQL migrations into binary to initialize required hashrate tables and indexes automatically. 🟢 |
| RF-02 | HPS Append Queries | Must | Append float64 hashrate metrics and int64 Unix timestamps securely to storage. 🟢 |
| RF-03 | Metrics History Retrieve | Must | Query recent hashrate history sorted descending by timestamp up to a specified limit. 🟢 |

## Non-Functional Requirements

| Type | Inferred Requirement | Code Evidence | Confidence |
|------|--------------------|---------------------|-----------|
| Portability | Standard Go SQL driver compatibility without CGO | `store.go` | 🟢 |
| Safety | Multithreading protection via connection pool thresholds | `store.go` | 🟢 |

## Acceptance Criteria

```gherkin
Given a store connected in WAL mode
When a concurrent thread appends hashrate records during busy locks
Then the execution respects a busy timeout of 5000ms and returns success without data races

Given NilStore is wired
When QueryHashRateHistory is evaluated
Then it returns an empty slice and nil error immediately
```

## Priority (MoSCoW)

| Requirement | MoSCoW | Justification |
|-----------|--------|---------------|
| Pure Go SQLite compilation | Must | Strict compliance with zero CGO design limits. 🟢 |
| WAL concurrency optimization | Should | Highly important to prevent lock contention in worker threads. 🟢 |

## Code Traceability

| File | Function / Class | Coverage |
|---------|-----------------|-----------|
| `internal/store/store.go` | `SQLiteStore.AppendHashRate` | 🟢 |
| `internal/store/store.go` | `SQLiteStore.QueryHashRateHistory` | 🟢 |
| `internal/store/migrations.go` | `migrations` DDL | 🟢 |

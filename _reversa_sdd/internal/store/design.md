# internal/store, Technical Design

> Design specification for the `internal/store` module. Focuses on HOW the storage is designed.

## Interface

### Classes / Functions

| Symbol | Signature | Return | Observation |
|---------|-----------|---------|------------|
| `SQLiteStore.AppendHashRate` | `(s *SQLiteStore) AppendHashRate(ctx context.Context, hps float64, at time.Time)` | `error` | Writes HPS metric to database. 🟢 |
| `SQLiteStore.QueryHashRateHistory` | `(s *SQLiteStore) QueryHashRateHistory(ctx context.Context, limit int)` | `([]float64, error)` | Retrieves slice of recent metrics. 🟢 |

## Main Flow
1. **Initialize Connection**:
   - Open DSN using Go driver: `modernc.org/sqlite`. 🟢
   - Append connection query pragmas: `_busy_timeout=5000` and `_journal_mode=WAL`. 🟢
2. **Execute Migrations**:
   - Parse embedded DDL schema from `migrations.go`.
   - Run `CREATE TABLE IF NOT EXISTS hashrate_history` and index bindings. 🟢
3. **Write / Read Operations**:
   - `AppendHashRate`: Build SQL query and call `ExecContext`. 🟢
   - `QueryHashRateHistory`: Build SQL query selecting the latest rows sorted descending by `recorded_at`. Shift values back to chronological order (FIFO) before returning. 🟢

## Dependencies
- `modernc.org/sqlite`: Pure Go SQL driver library. 🟢

## Design Decisions Identified

| Decision | Evidence | Confidence |
|---------|---------------------|-----------|
| Embedded SQL DDL | `migrations.go:embed` directives | 🟢 |
| Busy Timeout Lock | `store.go:_busy_timeout=5000` | 🟢 |

## Internal State
- **SQLiteStore**:
  - `db`: Pointer to native `*sql.DB` connection pool handler. 🟢

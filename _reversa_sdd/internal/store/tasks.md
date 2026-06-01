# internal/store, Implementation Tasks

> Implementation checklists for SQLite migrations and queries.

## Prerequisites
- [ ] Go sqlite driver package `modernc.org/sqlite` is added to `go.mod`. 🟢

## Tasks

- [ ] **T-01: Embed SQL Migrations**
  - Origin in Legacy: `nerdtui-spec.md:315` (§5.7)
  - Criteria of Done: Create `internal/store/migrations.go` containing embedded DDL SQL strings.
  - Confidence: 🟢 CONFIRMADO

- [ ] **T-02: Implement Store Interface and SQLiteStore**
  - Origin in Legacy: `nerdtui-spec.md:307` (§5.7)
  - Criteria of Done: Write `SQLiteStore` struct and connection loop supporting WAL mode.
  - Confidence: 🟢 CONFIRMADO

- [ ] **T-03: Implement Metrics database queries**
  - Origin in Legacy: `nerdtui-spec.md:309` (§5.7)
  - Criteria of Done: Write `AppendHashRate` and `QueryHashRateHistory` execution statements.
  - Confidence: 🟢 CONFIRMADO

- [ ] **T-04: Implement NilStore fallback**
  - Origin in Legacy: `nerdtui-spec.md:325` (§5.7)
  - Criteria of Done: Write no-op `NilStore` satisfying `Store` interface methods.
  - Confidence: 🟢 CONFIRMADO

## Test Tasks

- [ ] **TT-01: SQLite memory database roundtrip**
  - Launch `SQLiteStore` with `:memory:` endpoint, write 3 records, query them back, and assert consistency.
- [ ] **TT-02: Multi-thread SQLite busy lock test**
  - Execute concurrent writes under `-race` flag to ensure WAL pragmas prevent deadlock states.

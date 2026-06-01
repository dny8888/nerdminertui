package store

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	_ "modernc.org/sqlite"
)

func TestSchemaCompilesAndExecutes(t *testing.T) {
	// 1. Open an in-memory database with the pragmas specified in the design.
	// WAL mode might not be strictly necessary for :memory: but we test the connection string.
	db, err := sql.Open("sqlite", ":memory:?_busy_timeout=5000&_journal_mode=WAL")
	assert.NoError(t, err)
	defer db.Close()

	// 2. Execute the embedded schema
	_, err = db.Exec(SchemaDDL)
	assert.NoError(t, err, "schema should execute without errors")

	// 3. Verify the table exists by inserting a test row
	res, err := db.Exec(`INSERT INTO hashrate_history (hps, recorded_at) VALUES (123.45, 1716912000)`)
	assert.NoError(t, err, "should be able to insert into hashrate_history")

	rowsAffected, err := res.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), rowsAffected)

	// 4. Verify types and reading
	var id int
	var hps float64
	var recordedAt int64

	row := db.QueryRow(`SELECT id, hps, recorded_at FROM hashrate_history LIMIT 1`)
	err = row.Scan(&id, &hps, &recordedAt)
	assert.NoError(t, err)
	assert.Equal(t, 123.45, hps)
	assert.Equal(t, int64(1716912000), recordedAt)
}

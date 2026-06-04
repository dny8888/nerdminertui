package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

// Store defines the interface for persisting application metrics.
type Store interface {
	AppendHashRate(ctx context.Context, hps float64, at time.Time) error
	QueryHashRateHistory(ctx context.Context, limit int) ([]float64, error)
	Close() error
}

// SQLiteStore implements Store using a local SQLite database in WAL mode.
type SQLiteStore struct {
	db *sql.DB
}

// NewSQLiteStore opens a connection to the SQLite database and runs migrations.
func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	dsn := dbPath
	if !strings.Contains(dsn, "?") {
		dsn += "?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)"
	}

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	if err := Migrate(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	return &SQLiteStore{db: db}, nil
}

// AppendHashRate writes a new hashrate metric to the database.
func (s *SQLiteStore) AppendHashRate(ctx context.Context, hps float64, at time.Time) error {
	_, err := s.db.ExecContext(ctx, "INSERT INTO hashrate_history (hps, recorded_at) VALUES (?, ?)", hps, at.Unix())
	return err
}

// QueryHashRateHistory retrieves the most recent hashrate metrics up to limit.
// It returns the metrics in chronological order (oldest first).
func (s *SQLiteStore) QueryHashRateHistory(ctx context.Context, limit int) ([]float64, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT hps FROM hashrate_history ORDER BY recorded_at DESC LIMIT ?", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var desc []float64
	for rows.Next() {
		var hps float64
		if err := rows.Scan(&hps); err != nil {
			return nil, err
		}
		desc = append(desc, hps)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Reverse to get chronological order (oldest to newest)
	for i := 0; i < len(desc)/2; i++ {
		desc[i], desc[len(desc)-1-i] = desc[len(desc)-1-i], desc[i]
	}
	return desc, nil
}

// Close gracefully closes the database connection pool.
func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

// NilStore is a silent no-op implementation of Store for mock execution.
type NilStore struct{}

// NewNilStore creates a new no-op store.
func NewNilStore() *NilStore {
	return &NilStore{}
}

// AppendHashRate does nothing and returns nil.
func (n *NilStore) AppendHashRate(ctx context.Context, hps float64, at time.Time) error {
	return nil
}

// QueryHashRateHistory returns an empty slice and nil.
func (n *NilStore) QueryHashRateHistory(ctx context.Context, limit int) ([]float64, error) {
	return []float64{}, nil
}

// Close does nothing and returns nil.
func (n *NilStore) Close() error {
	return nil
}

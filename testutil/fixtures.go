package testutil

import (
	"fmt"
	"testing"

	"github.com/nerdminertui/nerdtui/internal/config"
	"github.com/nerdminertui/nerdtui/internal/model"
	"github.com/nerdminertui/nerdtui/internal/store"
)

// DefaultAppState returns a baseline application state for use in testing.
func DefaultAppState() model.AppState {
	return model.AppState{
		CPUTarget:       0.5,
		HashRateHistory: [model.HashHistoryLen]float64{},
		Screen:          model.ScreenDashboard,
	}
}

// DummyJobHeader returns a fixed 80-byte header for testing hashing logic.
func DummyJobHeader() []byte {
	b := make([]byte, 80)
	for i := range b {
		b[i] = 0xAA
	}
	return b
}

// NewTestStore creates a new in-memory SQLite store for isolated testing.
func NewTestStore(t testing.TB) store.Store {
	dbName := fmt.Sprintf("file:test_%s?mode=memory&cache=shared", t.Name())
	s, err := store.NewSQLiteStore(dbName)
	if err != nil {
		t.Fatalf("Failed to create test store: %v", err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

// NewTestConfig returns a basic valid configuration.
func NewTestConfig(t testing.TB) *config.Config {
	return &config.Config{
		PoolAddress: "public-pool.io",
		PoolPort:    21496,
		CPUTarget:   0.5,
		MockMining:  true,
	}
}

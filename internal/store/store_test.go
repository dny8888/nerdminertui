package store

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSQLiteStore_AppendAndQuery(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	// Use an isolated in-memory DB per connection for SQLite testing
	dbName := fmt.Sprintf("file:test_%s?mode=memory&cache=shared", t.Name())
	s, err := NewSQLiteStore(dbName)
	require.NoError(t, err)
	defer s.Close()

	// Append some entries
	err = s.AppendHashRate(ctx, 100.0, time.Unix(1000, 0))
	assert.NoError(t, err)
	err = s.AppendHashRate(ctx, 200.0, time.Unix(1010, 0))
	assert.NoError(t, err)
	err = s.AppendHashRate(ctx, 300.0, time.Unix(1020, 0))
	assert.NoError(t, err)

	// Query last 2 entries
	history, err := s.QueryHashRateHistory(ctx, 2)
	assert.NoError(t, err)
	
	// Should be chronologically ordered (oldest to newest of the last 2)
	assert.Len(t, history, 2)
	assert.Equal(t, 200.0, history[0])
	assert.Equal(t, 300.0, history[1])
}

func TestSQLiteStore_ConcurrentAppend(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	dbName := fmt.Sprintf("file:test_%s?mode=memory&cache=shared", t.Name())
	s, err := NewSQLiteStore(dbName)
	require.NoError(t, err)
	defer s.Close()

	var wg sync.WaitGroup
	workers := 10
	appendsPerWorker := 20

	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func(w int) {
			defer wg.Done()
			for j := 0; j < appendsPerWorker; j++ {
				_ = s.AppendHashRate(ctx, float64(w*100+j), time.Now())
			}
		}(i)
	}

	wg.Wait()

	history, err := s.QueryHashRateHistory(ctx, workers*appendsPerWorker)
	assert.NoError(t, err)
	assert.Len(t, history, workers*appendsPerWorker)
}

func TestNilStore(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	n := NewNilStore()

	err := n.AppendHashRate(ctx, 100.0, time.Now())
	assert.NoError(t, err)

	history, err := n.QueryHashRateHistory(ctx, 10)
	assert.NoError(t, err)
	assert.Empty(t, history)

	assert.NoError(t, n.Close())
}

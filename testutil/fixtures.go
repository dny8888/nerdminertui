package testutil

import (
	"github.com/nerdminertui/nerdtui/internal/model"
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

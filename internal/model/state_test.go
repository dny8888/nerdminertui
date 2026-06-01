package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithHashRate_ImmutabilityAndFIFO(t *testing.T) {
	initial := AppState{
		HashRate: 100,
	}
	// Pre-fill history to test shifting
	for i := 0; i < HashHistoryLen; i++ {
		initial.HashRateHistory[i] = float64(i)
	}

	updated := initial.WithHashRate(999.0)

	// Verify the original state is unchanged
	assert.Equal(t, 100.0, initial.HashRate)
	assert.Equal(t, float64(HashHistoryLen-1), initial.HashRateHistory[HashHistoryLen-1])

	// Verify the updated state has the new HashRate
	assert.Equal(t, 999.0, updated.HashRate)

	// Verify FIFO shift: index 0 should be the old index 1, index 59 should be 999.0
	assert.Equal(t, 1.0, updated.HashRateHistory[0])
	assert.Equal(t, 999.0, updated.HashRateHistory[HashHistoryLen-1])
}

func TestNextScreen_Cyclic(t *testing.T) {
	state := AppState{Screen: ScreenDashboard}

	state = state.NextScreen()
	assert.Equal(t, ScreenClock, state.Screen)

	state = state.NextScreen()
	assert.Equal(t, ScreenGlobalStats, state.Screen)

	state = state.NextScreen()
	assert.Equal(t, ScreenSettings, state.Screen)

	state = state.NextScreen()
	assert.Equal(t, ScreenDashboard, state.Screen)
}

func TestWithCPUTarget_Clamp(t *testing.T) {
	state := AppState{CPUTarget: 0.5}

	// Add step
	state = state.WithCPUTarget(0.05)
	assert.InDelta(t, 0.55, state.CPUTarget, 0.0001)

	// Clamp max
	state = state.WithCPUTarget(10.0)
	assert.InDelta(t, MaxCPUTarget, state.CPUTarget, 0.0001)

	// Clamp min
	state = state.WithCPUTarget(-10.0)
	assert.InDelta(t, MinCPUTarget, state.CPUTarget, 0.0001)
}

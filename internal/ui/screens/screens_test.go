package screens

import (
	"testing"

	"github.com/nerdminertui/nerdtui/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestRenderDashboard(t *testing.T) {
	state := model.AppState{
		HashRateHistory: [model.HashHistoryLen]float64{},
		BestDifficulty:  123.45,
	}
	state.HashRateHistory[model.HashHistoryLen-1] = 1000.0

	out := RenderDashboard(state, "", 80, 24)
	assert.Contains(t, out, "DASHBOARD")
	assert.Contains(t, out, "Hash Rate")
}

func TestRenderClock(t *testing.T) {
	state := model.AppState{}
	out := RenderClock(state, 80, 24)
	assert.Contains(t, out, ":") // simple clock check
}

func TestRenderGlobalStats(t *testing.T) {
	state := model.AppState{BlockHeight: 123456}
	out := RenderGlobalStats(state, 80, 24)
	assert.Contains(t, out, "123456")
	assert.Contains(t, out, "GLOBAL")
}

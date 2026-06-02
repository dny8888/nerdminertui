package components

import (
	"testing"

	"github.com/nerdminertui/nerdtui/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestRenderStatusBar(t *testing.T) {
	state := model.AppState{
		SharesFound: 5,
		Screen:      0,
	}
	out := RenderStatusBar(state, 80)
	assert.Contains(t, out, "shares: 5")
	assert.Contains(t, out, "tela 1/3")
}

func TestRenderCPUBar(t *testing.T) {
	state := model.AppState{
		CPUTarget: 0.5,
	}
	out := RenderCPUBar(state, 80)
	assert.Contains(t, out, "cpu [")
	assert.Contains(t, out, "50%")
}

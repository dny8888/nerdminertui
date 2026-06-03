package ui

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nerdminertui/nerdtui/internal/model"
	"github.com/nerdminertui/nerdtui/internal/worker"
	"github.com/stretchr/testify/assert"
)

func TestAppModel_UpdateAndScreens(t *testing.T) {
	throttleCh := make(chan float64, 10)
	state := model.AppState{
		CPUTarget: 0.5,
	}
	app := NewAppModel(state, throttleCh, nil)

	// Simulate Window Size
	m, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	app = m.(AppModel)

	// Tab rotates screens
	for i := 0; i < 4; i++ {
		m, _ = app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t', 'a', 'b'}, Alt: false})
		app = m.(AppModel)
		view := app.View()
		assert.NotEmpty(t, view)
	}

	// Plus key increases CPU Target
	m, _ = app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'+'}, Alt: false})
	app = m.(AppModel)
	assert.Equal(t, 0.55, app.state.CPUTarget)

	// Minus key decreases CPU Target
	m, _ = app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'-'}, Alt: false})
	app = m.(AppModel)
	assert.Equal(t, 0.50, app.state.CPUTarget)

	// Test messages
	app.state.PoolConnected = true
	m, _ = app.Update(worker.HashRateMsg{HPS: 1000, CPUActual: 0.5})
	app = m.(AppModel)
	assert.Equal(t, 1000.0, app.state.HashRateHistory[len(app.state.HashRateHistory)-1])

	m, _ = app.Update(worker.ShareFoundMsg{Accepted: true})
	app = m.(AppModel)
	assert.Equal(t, uint64(1), app.state.SharesFound)

	m, _ = app.Update(worker.PoolStatsMsg{GlobalHashRate: 4000})
	app = m.(AppModel)
	assert.Equal(t, float64(4000), app.state.NetworkHashRate)

	// Test pure rendering of screens
	app.state.Screen = 0
	view0 := app.View()
	app.state.Screen = 1
	view1 := app.View()
	app.state.Screen = 2
	view2 := app.View()
	assert.NotEqual(t, view0, view1)
	assert.NotEqual(t, view1, view2)
	assert.NotEqual(t, view0, view2)

	// Quit
	app.state.Screen = model.ScreenDashboard // Ensure we are not in settings where q is text input
	_, cmd := app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}, Alt: false})
	assert.Equal(t, tea.Quit(), cmd())
}

func TestInteractiveLoop(t *testing.T) {
	state := model.AppState{CPUTarget: 0.5}
	app := NewAppModel(state, nil, nil)

	// Just test Init and a tick
	cmd := app.Init()
	assert.NotNil(t, cmd)

	msg := cmd()
	m, _ := app.Update(msg)
	app = m.(AppModel)
	assert.Equal(t, time.Second, app.state.Uptime)
}

func TestFormatThreeDigits(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{0, "0.00"},
		{-5, "0.00"},
		{1.234, "1.23"},
		{9.999, "10.0"}, // formatting rounds up, which is fine
		{12.34, "12.3"},
		{99.99, "100"},
		{123.4, "123"},
		{999.9, "1000"},
	}

	for _, tc := range tests {
		actual := formatThreeDigits(tc.input)
		assert.Equal(t, tc.expected, actual, "For input %v", tc.input)
	}
}

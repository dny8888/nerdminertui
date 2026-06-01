package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/nerdminertui/nerdtui/internal/model"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).MarginTop(1)
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).MarginTop(1).Bold(true)
)

// SettingsModel holds the state for the settings form
type SettingsModel struct {
	Inputs     []textinput.Model
	FocusIndex int
}

// NewSettingsModel creates a new settings form model
func NewSettingsModel(state model.AppState) SettingsModel {
	m := SettingsModel{
		Inputs: make([]textinput.Model, 6),
	}

	var t textinput.Model
	for i := range m.Inputs {
		t = textinput.New()
		t.CharLimit = 128

		switch i {
		case 0:
			t.Placeholder = "BTC Address"
			t.SetValue(state.BTCAddress)
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Worker Name (e.g. .nerdtui)"
			t.SetValue(state.WorkerName)
		case 2:
			t.Placeholder = "Pool URL (e.g. public-pool.io)"
			t.SetValue(state.PoolAddress)
		case 3:
			t.Placeholder = "Pool Port"
			t.SetValue(fmt.Sprintf("%d", state.PoolPort))
		case 4:
			t.Placeholder = "Mock Mining (true/false)"
			t.SetValue(fmt.Sprintf("%t", state.MockMining))
		case 5:
			t.Placeholder = "CPU Target % (Max 75)"
			t.SetValue(fmt.Sprintf("%.0f", state.CPUTarget*100))
		}

		m.Inputs[i] = t
	}

	return m
}

// RenderSettings renders the configuration settings form.
func RenderSettings(state model.AppState, inputs []textinput.Model, width, height int) string {
	var b strings.Builder
	b.WriteString(lipgloss.NewStyle().Bold(true).Render("NERDMINER TUI - SETTINGS"))
	b.WriteString("\n\n")

	if !state.ConfigValid {
		b.WriteString(errorStyle.Render("Configuration is missing or invalid. Please fix and press Ctrl+S to save."))
		b.WriteString("\n\n")
	}

	for i := range inputs {
		b.WriteString(inputs[i].View())
		b.WriteString("\n")
	}

	b.WriteString(helpStyle.Render("UP: Prev Field  |  DOWN: Next Field  |  CTRL+S: Save and Apply"))

	content := b.String()

	// Center content
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}

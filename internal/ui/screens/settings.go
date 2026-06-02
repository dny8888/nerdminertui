package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/nerdminertui/nerdtui/internal/model"
	"github.com/nerdminertui/nerdtui/internal/ui/components"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(components.ColorOrange)
	helpStyle    = lipgloss.NewStyle().Foreground(components.ColorGray).MarginTop(1)
	errorStyle   = lipgloss.NewStyle().Foreground(components.ColorRed).MarginTop(1).Bold(true)
)

// SettingsModel holds the state for the settings form
type SettingsModel struct {
	Inputs     []textinput.Model
	FocusIndex int
	MockMining bool
}

// NewSettingsModel creates a new settings form model
func NewSettingsModel(state model.AppState) SettingsModel {
	m := SettingsModel{
		Inputs:     make([]textinput.Model, 5),
		MockMining: state.MockMining,
	}

	var t textinput.Model
	for i := range m.Inputs {
		t = textinput.New()
		t.CharLimit = 128
		t.Prompt = " "

		switch i {
		case 0:
			// t.Placeholder = "bc1q..."
			t.SetValue(state.BTCAddress)
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			// t.Placeholder = ".nerdtui"
			t.SetValue(state.WorkerName)
		case 2:
			// t.Placeholder = "public-pool.io"
			t.SetValue(state.PoolAddress)
		case 3:
			// t.Placeholder = "21496"
			t.SetValue(fmt.Sprintf("%d", state.PoolPort))
		case 4:
			// t.Placeholder = "75"
			t.SetValue(fmt.Sprintf("%.0f", state.CPUTarget*100))
		}

		m.Inputs[i] = t
	}

	return m
}

func renderInput(label string, t textinput.Model, width int, focused bool) string {
	lbl := components.StyleDim.Render(label)
	
	borderStyle := components.BorderStyle.Width(width)
	if focused {
		borderStyle = borderStyle.BorderForeground(components.ColorYellow)
	}

	// Make text input take full inner width
	t.Width = width - 2 
	
	box := borderStyle.Render(t.View())
	return fmt.Sprintf("%s\n%s", lbl, box)
}

// RenderSettings renders the configuration settings form.
func RenderSettings(state model.AppState, sm SettingsModel, width, height int) string {
	var b strings.Builder

	b.WriteString(components.RenderHeader(state, "SETTINGS", width))

	if !state.ConfigValid {
		b.WriteString(errorStyle.Render("Configuration is missing or invalid. Please fix and press Ctrl+S to save."))
		b.WriteString("\n\n")
	}

	// Calculate if focused
	f0 := sm.FocusIndex == 0
	f1 := sm.FocusIndex == 1
	f2 := sm.FocusIndex == 2
	f3 := sm.FocusIndex == 3
	fMock := sm.FocusIndex == 4 // Mock is 4 according to app.go
	f4 := sm.FocusIndex == 5 // CPU Target is 5 according to app.go

	// Use lipgloss to layout the main grid
	gridStyle := lipgloss.NewStyle().MarginLeft(4)

	// Row 1: BTC Address (width 66)
	b.WriteString(gridStyle.Render(renderInput("BTC Address", sm.Inputs[0], 66, f0)))
	b.WriteString("\n\n")

	// Row 2: Worker Name (16) | Pool URL (34) | Port (12)
	w1 := renderInput("Worker Name", sm.Inputs[1], 16, f1)
	w2 := renderInput("Pool URL", sm.Inputs[2], 34, f2)
	w3 := renderInput("Port", sm.Inputs[3], 12, f3)
	row2 := lipgloss.JoinHorizontal(lipgloss.Top, w1, "  ", w2, "  ", w3)
	b.WriteString(gridStyle.Render(row2))
	b.WriteString("\n\n")

	// Mock mining radio
	var radio string
	if sm.MockMining {
		radio = "(•) true   ( ) false"
	} else {
		radio = "( ) true   (•) false"
	}
	
	radioStyle := components.StyleDim
	borderStyle := components.BorderStyle.Width(16).PaddingLeft(1)
	if fMock {
		borderStyle = borderStyle.BorderForeground(components.ColorYellow)
		radioStyle = focusedStyle
	}
	
	mockLabel := components.StyleDim.Render("Mock Mode")
	mockBox := borderStyle.Render(radioStyle.Render(radio))
	mockWidget := fmt.Sprintf("%s\n%s", mockLabel, mockBox)

	// Row 3: CPU % (16) | Mock Mode (16)
	w4 := renderInput("CPU % (max 75)", sm.Inputs[4], 16, f4)
	row3 := lipgloss.JoinHorizontal(lipgloss.Top, w4, "  ", mockWidget)
	b.WriteString(gridStyle.Render(row3))
	b.WriteString("\n\n\n\n")

	// Footer
	footer := helpStyle.Render("↑ ↓ campo  ·  ctrl+s salvar e reiniciar  ·  tab outras telas")
	b.WriteString(gridStyle.Render(footer))

	content := b.String()

	// Optionally center content, but the mockup shows it left-aligned inside the window
	return content
}

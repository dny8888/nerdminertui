package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/nerdminertui/nerdtui/internal/model"
	"github.com/nerdminertui/nerdtui/internal/ui/components"
	"github.com/nerdminertui/nerdtui/pkg/format"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFA500")).
			MarginBottom(1)
)

// RenderDashboard renders Screen 0: the main hashing dashboard.
func RenderDashboard(state model.AppState, chartView string, width, height int) string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("NERDMINER TUI - DASHBOARD"))
	b.WriteString("\n")

	// Current HPS is at the end of the circular buffer after WithHashRate
	currentHPS := state.HashRateHistory[len(state.HashRateHistory)-1]

	b.WriteString(fmt.Sprintf("Hash Rate: %s\n", format.FormatHashRate(currentHPS)))
	b.WriteString(fmt.Sprintf("Best Diff: %s\n", format.FormatDifficulty(state.BestDifficulty)))
	b.WriteString("\n")
	
	b.WriteString(components.RenderCPUBar(state, width))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("Use +/- to adjust CPU limit"))
	b.WriteString("\n\n")

	// Render the ntcharts timeseries graph
	b.WriteString(chartView)

	// Pad to height
	lines := strings.Split(b.String(), "\n")
	padding := height - len(lines)
	if padding > 0 {
		b.WriteString(strings.Repeat("\n", padding))
	}

	return b.String()
}

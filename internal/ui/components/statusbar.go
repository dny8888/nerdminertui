package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/nerdminertui/nerdtui/internal/model"
)

var (
	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#333333")).
			Padding(0, 1)
)

// RenderStatusBar returns the bottom status bar string.
func RenderStatusBar(state model.AppState, width int) string {
	uptime := state.Uptime.String()
	shares := fmt.Sprintf("Shares: %d", state.SharesFound)
	connStatus := state.ConnectionStatus
	if connStatus == "" {
		connStatus = "Desconectado"
	}
	screen := fmt.Sprintf("Screen: %d/%d", state.Screen+1, model.NumScreens)
	
	left := uptime + " | " + shares + " | " + connStatus
	right := screen

	spaces := width - lipgloss.Width(left) - lipgloss.Width(right) - 2
	if spaces < 0 {
		spaces = 0
	}

	content := left + strings.Repeat(" ", spaces) + right
	return statusBarStyle.Render(content)
}

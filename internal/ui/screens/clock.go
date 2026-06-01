package screens

import (
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/nerdminertui/nerdtui/internal/model"
)

// RenderClock renders Screen 1: a large clock.
func RenderClock(state model.AppState, width, height int) string {
	now := time.Now().Format("15:04:05")
	
	clockStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00FFFF")).
		Align(lipgloss.Center).
		Width(width)

	content := clockStyle.Render(now)
	
	// vertically center
	lines := strings.Split(content, "\n")
	padTop := (height - len(lines)) / 2
	if padTop < 0 {
		padTop = 0
	}
	
	padBottom := height - len(lines) - padTop
	if padBottom < 0 {
		padBottom = 0
	}

	return strings.Repeat("\n", padTop) + content + strings.Repeat("\n", padBottom)
}

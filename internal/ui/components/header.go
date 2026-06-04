package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/nerdminertui/nerdtui/internal/model"
)

var (
	headerTitleStyle  = lipgloss.NewStyle().Foreground(ColorOrange).Bold(true)
	headerLineStyle   = lipgloss.NewStyle().Foreground(ColorGray)
	headerScreenStyle = lipgloss.NewStyle().Foreground(ColorWhite).Bold(true)
	headerBlockStyle  = lipgloss.NewStyle().Foreground(ColorGreen)
)

// RenderHeader renders the common top header for all screens.
func RenderHeader(state model.AppState, screenName string, width int) string {
	titleStr := "▶ NERDMINER TUI"
	
	var blockStr string
	if state.BlockHeight > 0 {
		blockStr = fmt.Sprintf("#%d", state.BlockHeight)
	} else {
		blockStr = "--"
	}

	// Render parts
	titlePart := headerTitleStyle.Render(titleStr)
	screenPart := headerScreenStyle.Render(strings.ToUpper(screenName))
	blockPart := headerBlockStyle.Render(blockStr)

	// Left section is: title + line + screenName + blockStr
	// We want to calculate the line length dynamically to fill the screen
	// width = len(title) + len(line) + len(screenName) + len(blockStr) + 3 spaces
	
	staticLen := lipgloss.Width(titlePart) + lipgloss.Width(screenPart) + lipgloss.Width(blockPart) + 3
	lineLen := width - staticLen - 2 // -2 for safety margin
	if lineLen < 4 {
		lineLen = 4
	}
	
	linePart := headerLineStyle.Render(strings.Repeat("─", lineLen))

	leftContent := fmt.Sprintf("%s %s %s %s", titlePart, linePart, screenPart, blockPart)

	return leftContent + "\n\n"
}

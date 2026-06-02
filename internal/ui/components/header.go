package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/nerdminertui/nerdtui/internal/model"
)

// RenderHeader renders the common top header for all screens.
func RenderHeader(state model.AppState, screenName string, width int) string {
	titleStr := "▶ NERDMINER TUI"
	blockStr := fmt.Sprintf("#%d", state.BlockHeight)

	// Styles
	titleStyle := lipgloss.NewStyle().Foreground(ColorOrange).Bold(true)
	lineStyle := lipgloss.NewStyle().Foreground(ColorGray)
	screenStyle := lipgloss.NewStyle().Foreground(ColorWhite).Bold(true)
	blockStyle := lipgloss.NewStyle().Foreground(ColorGreen)

	// Render parts
	titlePart := titleStyle.Render(titleStr)
	screenPart := screenStyle.Render(strings.ToUpper(screenName))
	blockPart := blockStyle.Render(blockStr)

	// Left section is: title + line + screenName + blockStr
	// e.g. "▶ NERDMINER TUI -------- DASHBOARD #892441"
	lineLen := 8
	linePart := lineStyle.Render(strings.Repeat("─", lineLen))

	leftContent := fmt.Sprintf("%s %s %s %s", titlePart, linePart, screenPart, blockPart)

	// Optional padding to full width if needed, but for now we just return the string
	return leftContent + "\n\n"
}

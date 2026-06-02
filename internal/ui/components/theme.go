package components

import "github.com/charmbracelet/lipgloss"

var (
	// Color Palette
	ColorOrange  = lipgloss.Color("#d79921")
	ColorGreen   = lipgloss.Color("#98c379")
	ColorCyan    = lipgloss.Color("#56b6c2")
	ColorMagenta = lipgloss.Color("#c678dd")
	ColorRed     = lipgloss.Color("#e06c75")
	ColorGray    = lipgloss.Color("#5c6370")
	ColorYellow  = lipgloss.Color("#e5c07b")
	ColorWhite   = lipgloss.Color("#abb2bf")
	ColorDark    = lipgloss.Color("#282c34")

	// Global Styles
	StyleLabel  = lipgloss.NewStyle().Foreground(ColorOrange)
	StyleValue  = lipgloss.NewStyle().Foreground(ColorGreen)
	StyleValue2 = lipgloss.NewStyle().Foreground(ColorCyan)
	StyleDim    = lipgloss.NewStyle().Foreground(ColorGray)

	// Border Style
	BorderStyle = lipgloss.NewStyle().
		BorderForeground(ColorGray).
		Border(lipgloss.NormalBorder())
)

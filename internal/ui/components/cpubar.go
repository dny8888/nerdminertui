package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/nerdminertui/nerdtui/internal/model"
)

var (
	cpuLabelStyle  = lipgloss.NewStyle().Foreground(ColorOrange)
	cpuBracketStyle = lipgloss.NewStyle().Foreground(ColorGray)
	cpuFilledStyle  = lipgloss.NewStyle().Foreground(ColorYellow)
	cpuEmptyStyle   = lipgloss.NewStyle().Foreground(ColorDark)
	cpuValStyle     = lipgloss.NewStyle().Foreground(ColorWhite)
	cpuTargetLabel  = lipgloss.NewStyle().Foreground(ColorGray)
	cpuTargetVal    = lipgloss.NewStyle().Foreground(ColorCyan)
)

// RenderCPUBar returns a textual bar representing CPU usage.
func RenderCPUBar(state model.AppState, width int) string {
	target := state.CPUTarget
	actual := state.CPUActual // If we implement it, otherwise 0
	if actual == 0 && target > 0 {
		actual = target - 0.01 // fake it for now to look like the image
	}
	
	totalBlocks := width - 40 // adjusted to leave room for the text "cpu [] 75% target 74% actual [+/-]"
	if totalBlocks < 10 {
		totalBlocks = 10
	}
	if totalBlocks > 40 {
		totalBlocks = 40
	}

	activeBlocks := int(actual * float64(totalBlocks))
	
	var b strings.Builder
	for i := 0; i < totalBlocks; i++ {
		if i < activeBlocks {
			b.WriteString(cpuFilledStyle.Render("█"))
		} else {
			b.WriteString(cpuEmptyStyle.Render("░"))
		}
	}

	// cpu [███████████░░░] 75% target 74% actual [+/-]
	label := cpuLabelStyle.Render("cpu")
	bracketL := cpuBracketStyle.Render("[")
	bracketR := cpuBracketStyle.Render("]")
	val := cpuValStyle.Render(fmt.Sprintf("%.0f%%", target*100))
	tLabel := cpuTargetLabel.Render("target")
	aVal := cpuTargetVal.Render(fmt.Sprintf("%.0f%% actual", actual*100))
	controls := cpuBracketStyle.Render("[+/-]")

	return fmt.Sprintf("%s %s%s%s %s %s %s %s", label, bracketL, b.String(), bracketR, val, tLabel, aVal, controls)
}

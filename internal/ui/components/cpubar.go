package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/nerdminertui/nerdtui/internal/model"
)

var (
	cpuBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00"))
	cpuTargetStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFF00"))
)

// RenderCPUBar returns a textual bar representing CPU usage.
func RenderCPUBar(state model.AppState, width int) string {
	target := state.CPUTarget
	actual := 0.0
	// If we wanted to parse HashRateHistory we could, but let's just make a dummy actual for now
	
	totalBlocks := width - 15
	if totalBlocks < 1 {
		totalBlocks = 1
	}

	activeBlocks := int(actual * float64(totalBlocks))
	targetBlock := int(target * float64(totalBlocks))

	var b strings.Builder
	for i := 0; i < totalBlocks; i++ {
		if i == targetBlock {
			b.WriteString(cpuTargetStyle.Render("|"))
		} else if i < activeBlocks {
			b.WriteString(cpuBarStyle.Render("█"))
		} else {
			b.WriteString(" ")
		}
	}

	return fmt.Sprintf("CPU [%s] %.0f%%", b.String(), target*100)
}

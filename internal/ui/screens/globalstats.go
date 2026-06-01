package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/nerdminertui/nerdtui/internal/model"
)

// RenderGlobalStats renders Screen 2: global network statistics.
func RenderGlobalStats(state model.AppState, width, height int) string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("GLOBAL NETWORK STATS"))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("Block Height: %d\n", state.BlockHeight))
	b.WriteString("Network Hashrate: (fetched from pool)\n")
	
	// Pad to height
	lines := strings.Split(b.String(), "\n")
	padding := height - len(lines)
	if padding > 0 {
		b.WriteString(strings.Repeat("\n", padding))
	}

	return lipgloss.NewStyle().Width(width).Render(b.String())
}

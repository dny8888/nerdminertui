package screens

import (
	"fmt"
	"strings"

	"github.com/nerdminertui/nerdtui/internal/model"
	"github.com/nerdminertui/nerdtui/internal/ui/components"
	"github.com/nerdminertui/nerdtui/pkg/format"
	"github.com/nerdminertui/nerdtui/pkg/i18n"
)

func RenderDashboard(state model.AppState, chartView string, currentUnit string, width, height int) string {
	var b strings.Builder

	b.WriteString(components.RenderHeader(state, i18n.DashboardTitle, width))

	// Current HPS is at the end of the circular buffer after WithHashRate
	currentHPS := state.HashRateHistory[len(state.HashRateHistory)-1]

	hashStr := components.StyleValue.Render(format.FormatHashRate(currentHPS))
	diffStr := components.StyleValue.Render(format.FormatDifficulty(state.BestDifficulty))
	sharesStr := components.StyleValue.Render(fmt.Sprintf("%d", state.SharesFound))
	
	// Uptime
	uptimeStr := components.StyleValue2.Render(format.FormatUptime(state.Uptime))

	// Layout Row 1: hash rate <val>   best diff <val>
	r1col1 := fmt.Sprintf("%-10s %-12s", components.StyleLabel.Render("hash rate"), hashStr)
	r1col2 := fmt.Sprintf("%-10s %s", components.StyleLabel.Render("best diff"), diffStr)
	b.WriteString(fmt.Sprintf("%s   %s\n", r1col1, r1col2))

	// Layout Row 2: shares <val>      uptime <val>
	r2col1 := fmt.Sprintf("%-10s %-12s", components.StyleLabel.Render("shares"), sharesStr)
	r2col2 := fmt.Sprintf("%-10s %s", components.StyleLabel.Render("uptime"), uptimeStr)
	b.WriteString(fmt.Sprintf("%s   %s\n\n", r2col1, r2col2))
	
	b.WriteString(components.RenderCPUBar(state, width))
	b.WriteString("\n\n")

	// Print the unit right above the graph, like a physics graph Y-axis
	b.WriteString(components.StyleLabel.Render(fmt.Sprintf("(%s)", currentUnit)) + "\n")

	// Render the ntcharts timeseries graph
	// Add borders to the graph
	graphView := components.BorderStyle.
		BorderTop(true).BorderLeft(true).BorderRight(true).BorderBottom(true).
		Width(width - 2).
		Render(chartView)
		
	b.WriteString(graphView)

	// Pad to height
	lines := strings.Split(b.String(), "\n")
	padding := height - len(lines)
	if padding > 0 {
		b.WriteString(strings.Repeat("\n", padding))
	}

	return b.String()
}

package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/nerdminertui/nerdtui/internal/model"
	"github.com/nerdminertui/nerdtui/pkg/format"
	"github.com/nerdminertui/nerdtui/pkg/i18n"
)

var (
	sbBaseStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("#1e1e1e")).
		Foreground(ColorGray).
		Padding(0, 1)

	sbAddrStyle   = lipgloss.NewStyle().Foreground(ColorGray)
	sbDotStyle    = lipgloss.NewStyle().Foreground(ColorGreen)
	sbStatusStyle = lipgloss.NewStyle().Foreground(ColorGreen)
	sbDotWarnStyle    = lipgloss.NewStyle().Foreground(ColorYellow)
	sbStatusWarnStyle = lipgloss.NewStyle().Foreground(ColorYellow)
	sbDotErrStyle     = lipgloss.NewStyle().Foreground(ColorRed)
	sbStatusErrStyle  = lipgloss.NewStyle().Foreground(ColorRed)
	sbValStyle    = lipgloss.NewStyle().Foreground(ColorGreen)
	sbUptimeStyle = lipgloss.NewStyle().Foreground(ColorCyan)
	sbTelaStyle   = lipgloss.NewStyle().Foreground(ColorOrange)
)

// RenderStatusBar returns the dense bottom status bar.
func RenderStatusBar(state model.AppState, width int) string {
	// public-pool.io:21496 ● Conectado | shares: 84 | up 2d 06h 36m | tela 1/3

	addr := sbAddrStyle.Render(fmt.Sprintf("%s:%d", state.PoolAddress, state.PoolPort))
	
	statusText := state.ConnectionStatus
	if statusText == "" {
		statusText = i18n.StatusDisconnected
	}
	
	var statusDot string	
	if !state.PoolConnected {
		if statusText == i18n.StatusConnecting {
			statusDot = sbDotWarnStyle.Render("●")
			statusText = sbStatusWarnStyle.Render(statusText)
		} else {
			statusDot = sbDotErrStyle.Render("●")
			statusText = sbStatusErrStyle.Render(statusText)
		}
	} else {
		statusDot = sbDotStyle.Render("●")
		statusText = sbStatusStyle.Render(statusText)
	}

	shares := fmt.Sprintf("%s %s", i18n.StatusShares, sbValStyle.Render(fmt.Sprintf("%d", state.SharesFound)))
	
	uptimeStr := format.FormatUptime(state.Uptime)
	up := fmt.Sprintf("%s %s", i18n.StatusUp, sbUptimeStyle.Render(uptimeStr))

	tela := sbTelaStyle.Render(fmt.Sprintf("%s %d/%d", i18n.StatusScreen, state.Screen+1, model.NumScreens))

	content := fmt.Sprintf("%s %s %s  |  %s  |  %s  |  %s", addr, statusDot, statusText, shares, up, tela)

	return sbBaseStyle.Width(width).Render(content)
}

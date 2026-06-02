package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/nerdminertui/nerdtui/internal/model"
)

var (
	sbBaseStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("#1e1e1e")).
		Foreground(ColorGray).
		Padding(0, 1)

	sbAddrStyle   = lipgloss.NewStyle().Foreground(ColorGray)
	sbDotStyle    = lipgloss.NewStyle().Foreground(ColorGreen)
	sbStatusStyle = lipgloss.NewStyle().Foreground(ColorGreen)
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
		statusText = "Desconectado"
	}
	
	var statusDot string	
	if !state.PoolConnected {
		if statusText == "Conectando..." {
			statusDot = lipgloss.NewStyle().Foreground(ColorYellow).Render("●")
			statusText = lipgloss.NewStyle().Foreground(ColorYellow).Render(statusText)
		} else {
			statusDot = lipgloss.NewStyle().Foreground(ColorRed).Render("●")
			statusText = lipgloss.NewStyle().Foreground(ColorRed).Render(statusText)
		}
	} else {
		statusDot = sbDotStyle.Render("●")
		statusText = sbStatusStyle.Render(statusText)
	}

	shares := fmt.Sprintf("shares: %s", sbValStyle.Render(fmt.Sprintf("%d", state.SharesFound)))
	
	uptimeStr := ""
	if state.Uptime > 0 {
		days := int(state.Uptime.Hours()) / 24
		hours := int(state.Uptime.Hours()) % 24
		minutes := int(state.Uptime.Minutes()) % 60
		if days > 0 {
			uptimeStr = fmt.Sprintf("%dd %02dh %02dm", days, hours, minutes)
		} else {
			uptimeStr = fmt.Sprintf("%02dh %02dm", hours, minutes)
		}
	} else {
		uptimeStr = "00h 00m"
	}
	up := fmt.Sprintf("up %s", sbUptimeStyle.Render(uptimeStr))

	tela := sbTelaStyle.Render(fmt.Sprintf("tela %d/%d", state.Screen+1, model.NumScreens))

	content := fmt.Sprintf("%s %s %s  |  %s  |  %s  |  %s", addr, statusDot, statusText, shares, up, tela)

	return sbBaseStyle.Width(width).Render(content)
}

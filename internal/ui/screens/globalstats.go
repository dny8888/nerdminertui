package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/nerdminertui/nerdtui/internal/model"
	"github.com/nerdminertui/nerdtui/internal/ui/components"
)

// RenderGlobalStats renders Screen 2: global network statistics.
func RenderGlobalStats(state model.AppState, width, height int) string {
	var b strings.Builder

	b.WriteString(components.RenderHeader(state, "GLOBAL STATS", width))

	// Card Styles
	cardStyle := components.BorderStyle.
		Padding(0, 1).
		Width(30)

	// Mock Data for now (Network Hashrate, Difficulty, Miners Online)
	netHash := components.StyleValue.Render("4.50 EH/s")
	diff := components.StyleValue.Render("88.1 T")
	blockHeight := components.StyleValue.Render(fmt.Sprintf("#%d", state.BlockHeight))
	minersOnline := components.StyleValue.Render("12.500")

	// Render Cards
	card1 := cardStyle.Render(fmt.Sprintf("%s\n\n%s", components.StyleDim.Render("Network Hashrate"), netHash))
	card2 := cardStyle.Render(fmt.Sprintf("%s\n\n%s", components.StyleDim.Render("Difficulty"), diff))
	card3 := cardStyle.Render(fmt.Sprintf("%s\n\n%s", components.StyleDim.Render("Block Height"), blockHeight))
	card4 := cardStyle.Render(fmt.Sprintf("%s\n\n%s", components.StyleDim.Render("Miners Online"), minersOnline))

	row1 := lipgloss.JoinHorizontal(lipgloss.Top, card1, "    ", card2)
	row2 := lipgloss.JoinHorizontal(lipgloss.Top, card3, "    ", card4)

	// Indent grid to match mockup using lipgloss
	gridStyle := lipgloss.NewStyle().MarginLeft(4)
	
	b.WriteString(gridStyle.Render(row1) + "\n\n")
	b.WriteString(gridStyle.Render(row2) + "\n\n\n")

	// Lottery Section
	lotteryCardStyle := lipgloss.NewStyle().
		BorderForeground(lipgloss.Color("#5c6370")).
		Border(lipgloss.NormalBorder(), true, true, true, true).
		Padding(0, 1).
		Width(64)

	lotteryTitle := components.StyleLabel.Render("Loteria")

	oddsLabel := components.StyleDim.Render("odds hoje:   ")
	oddsVal1 := components.StyleLabel.Render("1 em 577.000.000.000")
	oddsVal2 := lipgloss.NewStyle().Foreground(components.ColorRed).Render("≈ 0.00000000017%")
	oddsLine := fmt.Sprintf("%s%s   %s", oddsLabel, oddsVal1, oddsVal2)

	esperadoLabel := components.StyleDim.Render("esperado:    ")
	esperadoVal1 := components.StyleValue2.Render("~18.000 anos")
	esperadoVal2 := components.StyleDim.Render("com o hashrate atual")
	esperadoLine := fmt.Sprintf("%s%s %s", esperadoLabel, esperadoVal1, esperadoVal2)

	ultimaLabel := components.StyleDim.Render("última vez:  ")
	ultimaVal := components.StyleValue.Render("nenhum bloco encontrado ainda")
	ultimaLine := fmt.Sprintf("%s%s", ultimaLabel, ultimaVal)

	lotteryContent := fmt.Sprintf("%s\n\n%s\n%s\n%s", lotteryTitle, oddsLine, esperadoLine, ultimaLine)
	b.WriteString(gridStyle.Render(lotteryCardStyle.Render(lotteryContent)))

	// Pad to height
	lines := strings.Split(b.String(), "\n")
	padding := height - len(lines)
	if padding > 0 {
		b.WriteString(strings.Repeat("\n", padding))
	}

	return lipgloss.NewStyle().Width(width).Render(b.String())
}

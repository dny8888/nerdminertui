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
	// Since we omitted Miners Online, we have 3 cards. Let's make them Width(24) and put all 3 in one row.
	cardStyle3 := components.BorderStyle.
		Padding(0, 1).
		Width(24)

	var netHashStr, diffStr, blockHeightStr string

	if state.NetworkHashRate > 0 {
		netHashStr = fmt.Sprintf("%.2f EH/s", state.NetworkHashRate/1e18)
	} else {
		netHashStr = "N/A"
	}

	if state.NetworkDifficulty > 0 {
		diffStr = fmt.Sprintf("%.1f T", state.NetworkDifficulty/1e12)
	} else {
		diffStr = "N/A"
	}

	if state.BlockHeight > 0 {
		blockHeightStr = fmt.Sprintf("#%d", state.BlockHeight)
	} else {
		blockHeightStr = "N/A"
	}

	// Render Cards
	card1 := cardStyle3.Render(fmt.Sprintf("%s\n\n%s", components.StyleDim.Render("Network Hashrate"), components.StyleValue.Render(netHashStr)))
	card2 := cardStyle3.Render(fmt.Sprintf("%s\n\n%s", components.StyleDim.Render("Difficulty"), components.StyleValue.Render(diffStr)))
	card3 := cardStyle3.Render(fmt.Sprintf("%s\n\n%s", components.StyleDim.Render("Block Height"), components.StyleValue.Render(blockHeightStr)))

	row1 := lipgloss.JoinHorizontal(lipgloss.Top, card1, "  ", card2, "  ", card3)

	// Indent grid to match mockup using lipgloss
	gridStyle := lipgloss.NewStyle().MarginLeft(4)
	
	b.WriteString(gridStyle.Render(row1) + "\n\n\n")

	// Lottery Section
	lotteryCardStyle := lipgloss.NewStyle().
		BorderForeground(lipgloss.Color("#5c6370")).
		Border(lipgloss.NormalBorder(), true, true, true, true).
		Padding(0, 1).
		Width(76)

	lotteryTitle := components.StyleLabel.Render("Loteria")

	var oddsVal1, oddsVal2, esperadoVal1 string
	if state.NetworkDifficulty > 0 {
		totalHashesForBlock := state.NetworkDifficulty * 4294967296.0 // Diff * 2^32
		
		// Create a readable representation for odds
		oddsStr := fmt.Sprintf("%.0f", totalHashesForBlock)
		if len(oddsStr) > 3 {
			// Basic thousands separator (dot) for readability
			var parts []string
			for i := len(oddsStr); i > 0; i -= 3 {
				start := i - 3
				if start < 0 {
					start = 0
				}
				parts = append([]string{oddsStr[start:i]}, parts...)
			}
			oddsStr = strings.Join(parts, ".")
		}
		
		oddsVal1 = components.StyleLabel.Render(fmt.Sprintf("1 em %s", oddsStr))
		oddsVal2 = lipgloss.NewStyle().Foreground(components.ColorRed).Render(fmt.Sprintf("≈ %.13f%%", 100.0/totalHashesForBlock))

		if state.HashRate > 0 {
			expectedSeconds := totalHashesForBlock / state.HashRate
			expectedYears := expectedSeconds / (365.25 * 24 * 3600)
			
			// Format years nicely
			yearsStr := fmt.Sprintf("%.0f", expectedYears)
			if len(yearsStr) > 3 {
				var parts []string
				for i := len(yearsStr); i > 0; i -= 3 {
					start := i - 3
					if start < 0 {
						start = 0
					}
					parts = append([]string{yearsStr[start:i]}, parts...)
				}
				yearsStr = strings.Join(parts, ".")
			}
			
			esperadoVal1 = components.StyleValue2.Render(fmt.Sprintf("~%s anos", yearsStr))
		} else {
			esperadoVal1 = components.StyleValue2.Render("∞ anos")
		}
	} else {
		oddsVal1 = components.StyleLabel.Render("N/A")
		oddsVal2 = ""
		esperadoVal1 = components.StyleValue2.Render("N/A")
	}

	oddsLabel := components.StyleDim.Render("odds hoje:   ")
	oddsLine := fmt.Sprintf("%s%s   %s", oddsLabel, oddsVal1, oddsVal2)

	esperadoLabel := components.StyleDim.Render("esperado:    ")
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

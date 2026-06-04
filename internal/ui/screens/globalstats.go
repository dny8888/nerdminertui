package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/nerdminertui/nerdtui/internal/model"
	"github.com/nerdminertui/nerdtui/internal/ui/components"
	"github.com/nerdminertui/nerdtui/pkg/format"
	"github.com/nerdminertui/nerdtui/pkg/i18n"
)

var (
	globalStatsGridStyle = lipgloss.NewStyle().
				MarginTop(1).
				MarginBottom(1)
	lotteryCardStyle = lipgloss.NewStyle().
				BorderForeground(lipgloss.Color("#5c6370")).
				Border(lipgloss.NormalBorder(), true, true, true, true).
				Padding(0, 1)
)

// RenderGlobalStats renders Screen 2: global network statistics.
func RenderGlobalStats(state model.AppState, width, height int) string {
	var b strings.Builder

	b.WriteString(components.RenderHeader(state, i18n.GlobalStatsTitle, width))

	// Card Styles
	cardWidth := (width - 16) / 3 // 3 cards, margin left 4, padding and gaps
	if cardWidth < 20 {
		cardWidth = 20
	}
	cardStyle3 := components.BorderStyle.
		Padding(0, 1).
		Width(cardWidth)

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
	card1 := cardStyle3.Render(fmt.Sprintf("%s\n\n%s", components.StyleDim.Render(i18n.NetworkHashrate), components.StyleValue.Render(netHashStr)))
	card2 := cardStyle3.Render(fmt.Sprintf("%s\n\n%s", components.StyleDim.Render(i18n.Difficulty), components.StyleValue.Render(diffStr)))
	card3 := cardStyle3.Render(fmt.Sprintf("%s\n\n%s", components.StyleDim.Render(i18n.BlockHeight), components.StyleValue.Render(blockHeightStr)))

	row1 := lipgloss.JoinHorizontal(lipgloss.Top, card1, "  ", card2, "  ", card3)

	b.WriteString(globalStatsGridStyle.Render(row1) + "\n\n\n")

	// Lottery Section
	lotteryWidth := width - 8
	if lotteryWidth < 50 {
		lotteryWidth = 50
	}
	lotteryStyle := lotteryCardStyle.Width(lotteryWidth)

	lotteryTitle := components.StyleLabel.Render(i18n.LotteryTitle)

	var oddsVal1, oddsVal2, esperadoVal1 string
	if state.NetworkDifficulty > 0 {
		totalHashesForBlock := state.NetworkDifficulty * 4294967296.0 // Diff * 2^32
		
		// Create a readable representation for odds
		oddsStr := format.FormatWithThousandSeparator(totalHashesForBlock)
		
		oddsVal1 = components.StyleLabel.Render(fmt.Sprintf(i18n.OneIn, oddsStr))
		oddsVal2 = lipgloss.NewStyle().Foreground(components.ColorRed).Render(fmt.Sprintf("≈ %.13f%%", 100.0/totalHashesForBlock))

		if state.HashRate > 0 {
			expectedSeconds := totalHashesForBlock / state.HashRate
			expectedYears := expectedSeconds / (365.25 * 24 * 3600)
			
			// Format years nicely
			yearsStr := format.FormatWithThousandSeparator(expectedYears)
			
			esperadoVal1 = components.StyleValue2.Render(fmt.Sprintf("~%s %s", yearsStr, i18n.Years))
		} else {
			esperadoVal1 = components.StyleValue2.Render(fmt.Sprintf("∞ %s", i18n.Years))
		}
	} else {
		oddsVal1 = components.StyleLabel.Render("N/A")
		oddsVal2 = ""
		esperadoVal1 = components.StyleValue2.Render("N/A")
	}

	oddsLabel := components.StyleDim.Render(i18n.OddsToday)
	oddsLine := fmt.Sprintf("%s%s   %s", oddsLabel, oddsVal1, oddsVal2)

	esperadoLabel := components.StyleDim.Render(i18n.Expected)
	esperadoVal2 := components.StyleDim.Render(i18n.WithCurrentHash)
	esperadoLine := fmt.Sprintf("%s%s %s", esperadoLabel, esperadoVal1, esperadoVal2)

	ultimaLabel := components.StyleDim.Render(i18n.LastTime)
	ultimaVal := components.StyleValue.Render(i18n.NoBlockFoundYet)
	ultimaLine := fmt.Sprintf("%s%s", ultimaLabel, ultimaVal)

	lotteryContent := fmt.Sprintf("%s\n\n%s\n%s\n%s", lotteryTitle, oddsLine, esperadoLine, ultimaLine)

	b.WriteString(globalStatsGridStyle.Render(lotteryStyle.Render(lotteryContent)))

	// Pad to height
	lines := strings.Split(b.String(), "\n")
	padding := height - len(lines)
	if padding > 0 {
		b.WriteString(strings.Repeat("\n", padding))
	}

	return lipgloss.NewStyle().Width(width).Render(b.String())
}

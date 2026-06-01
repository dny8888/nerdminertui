package worker

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
)

// PollCmd returns a tea.Cmd that fetches global stats from the pool client.
func PollCmd(ctx context.Context, client PoolClient) tea.Cmd {
	return func() tea.Msg {
		stats, err := client.FetchStats(ctx)
		if err != nil {
			return PoolErrorMsg{Err: err}
		}
		return stats
	}
}

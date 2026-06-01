package ui

import (
	"fmt"
	"time"

	"github.com/NimbleMarkets/ntcharts/linechart/timeserieslinechart"
	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nerdminertui/nerdtui/internal/model"
	"github.com/nerdminertui/nerdtui/internal/ui/components"
	"github.com/nerdminertui/nerdtui/internal/ui/screens"
	"github.com/nerdminertui/nerdtui/internal/worker"
)

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// SaveConfigMsg is sent when the user saves the settings
type SaveConfigMsg struct {
	Config *model.AppState
}

// AppModel implements tea.Model for the main application.
type AppModel struct {
	state          model.AppState
	throttleCh     chan<- float64
	configUpdateCh chan<- *model.AppState
	settings       screens.SettingsModel
	hashChart      timeserieslinechart.Model
	width          int
	height         int
}

// NewAppModel initializes a new AppModel.
func NewAppModel(initialState model.AppState, throttleCh chan<- float64, configUpdateCh chan<- *model.AppState) AppModel {
	chart := timeserieslinechart.New(40, 15)
	chart.AutoMinY = true
	chart.AutoMaxY = true
	chart.XLabelFormatter = func(i int, v float64) string {
		return time.Unix(int64(v), 0).Local().Format("15:04:05")
	}
	
	return AppModel{
		state:          initialState,
		throttleCh:     throttleCh,
		configUpdateCh: configUpdateCh,
		settings:       screens.NewSettingsModel(initialState),
		hashChart:      chart,
	}
}

// Init sets up the initial commands.
func (m AppModel) Init() tea.Cmd {
	return tickCmd()
}

// Update processes incoming messages and updates state.
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			if msg.String() == "ctrl+c" || m.state.Screen != model.ScreenSettings {
				return m, tea.Quit
			}
		case "ctrl+s":
			// Update state from inputs before saving
			m.state.BTCAddress = m.settings.Inputs[0].Value()
			m.state.WorkerName = m.settings.Inputs[1].Value()
			m.state.PoolAddress = m.settings.Inputs[2].Value()
			_, _ = fmt.Sscanf(m.settings.Inputs[3].Value(), "%d", &m.state.PoolPort)
			m.state.MockMining = m.settings.Inputs[4].Value() == "true"
			
			var cpuTargetInt int
			if _, err := fmt.Sscanf(m.settings.Inputs[5].Value(), "%d", &cpuTargetInt); err == nil {
				if cpuTargetInt >= 5 && cpuTargetInt <= 75 {
					m.state.CPUTarget = float64(cpuTargetInt) / 100.0
				}
			}
			
			// Notify main loop to save config and restart worker
			cfgCopy := m.state
			return m, func() tea.Msg { return SaveConfigMsg{Config: &cfgCopy} }
		case "up", "down":
			if m.state.Screen == model.ScreenSettings {
				s := msg.String()
				
				// Handle focus switching
				if s == "up" {
					m.settings.FocusIndex--
				} else {
					m.settings.FocusIndex++
				}

				if m.settings.FocusIndex > len(m.settings.Inputs)-1 {
					m.settings.FocusIndex = 0
				} else if m.settings.FocusIndex < 0 {
					m.settings.FocusIndex = len(m.settings.Inputs) - 1
				}

				var cmds []tea.Cmd
				for i := 0; i <= len(m.settings.Inputs)-1; i++ {
					if i == m.settings.FocusIndex {
						cmds = append(cmds, m.settings.Inputs[i].Focus())
						m.settings.Inputs[i].PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
						m.settings.Inputs[i].TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
					} else {
						m.settings.Inputs[i].Blur()
						m.settings.Inputs[i].PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
						m.settings.Inputs[i].TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
					}
				}
				return m, tea.Batch(cmds...)
			}
		case "tab", "right":
			m.state = m.state.NextScreen()
		case "+", "=":
			m.state = m.state.WithCPUTarget(0.05)
			if m.throttleCh != nil {
				m.throttleCh <- m.state.CPUTarget
			}
		case "-", "_":
			m.state = m.state.WithCPUTarget(-0.05)
			if m.throttleCh != nil {
				m.throttleCh <- m.state.CPUTarget
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		
		chartWidth := m.width - 4
		chartHeight := m.height - 13 // Reserve space for title, CPU bar, status bar, padding
		if chartWidth < 10 {
			chartWidth = 10
		}
		if chartHeight < 5 {
			chartHeight = 5
		}
		m.hashChart.Resize(chartWidth, chartHeight)

	case tickMsg:
		m.state.Uptime += time.Second
		cmds = append(cmds, tickCmd())

	case worker.HashRateMsg:
		now := time.Now()
		m.state = m.state.WithHashRate(msg.HPS)
		m.state.CPUActual = msg.CPUActual
		m.hashChart.Push(timeserieslinechart.TimePoint{Time: now, Value: msg.HPS})
		m.hashChart.SetTimeRange(now.Add(-30*time.Second), now)
		m.hashChart.SetViewTimeRange(now.Add(-30*time.Second), now)
		m.hashChart.DrawBraille()

	case worker.ShareFoundMsg:
		if msg.Accepted {
			m.state.SharesFound++
		}

	case worker.PoolStatsMsg:
		m.state.BlockHeight = uint32(msg.GlobalHashRate)

	case worker.MinerErrorMsg:
		// Ignoring for now or handle appropriately
	case worker.PoolErrorMsg:
		// Ignoring or handle
	case worker.ConnectionStatusMsg:
		m.state.ConnectionStatus = msg.Status
	case SaveConfigMsg:
		msg.Config.ConfigValid = true
		m.state = *msg.Config
		
		// Use a goroutine to not block the UI thread
		go func(stateCopy model.AppState) {
			m.configUpdateCh <- &stateCopy
		}(m.state)
	}

	// Handle input updates if on settings screen
	if m.state.Screen == model.ScreenSettings {
		var cmd tea.Cmd
		for i := range m.settings.Inputs {
			m.settings.Inputs[i], cmd = m.settings.Inputs[i].Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

// View renders the terminal output.
func (m AppModel) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	var content string
	switch m.state.Screen {
	case 0:
		content = screens.RenderDashboard(m.state, m.hashChart.View(), m.width, m.height-1)
	case 1:
		content = screens.RenderClock(m.state, m.width, m.height-1)
	case 2:
		content = screens.RenderGlobalStats(m.state, m.width, m.height-1)
	case 3:
		content = screens.RenderSettings(m.state, m.settings.Inputs, m.width, m.height-1)
	default:
		content = "Unknown Screen"
	}

	statusBar := components.RenderStatusBar(m.state, m.width)
	return fmt.Sprintf("%s\n%s", content, statusBar)
}

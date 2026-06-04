package model

import (
	"time"
)

// ScreenID represents the current visible screen in the TUI.
type ScreenID uint8

const (
	ScreenDashboard ScreenID = 0
	ScreenGlobalStats ScreenID = 1
	ScreenSettings    ScreenID = 2
)

// NumScreens is the total number of screens.
const NumScreens = 3

// CPU throttle constants.
const (
	MinCPUTarget float64 = 0.05
	MaxCPUTarget float64 = 0.75
	CPUStep      float64 = 0.05
)

// HashHistoryLen is the length of the circular FIFO hashrate history buffer.
const HashHistoryLen = 60

// AppState represents the complete state of the TUI application.
// It is a value type used in Bubbletea's MUV pattern.
type AppState struct {
	HashRate        float64
	HashRateHistory [HashHistoryLen]float64
	SharesFound     uint64
	BestDifficulty  float64
	BlockHeight     uint32
	NetworkHashRate float64
	NetworkDifficulty float64
	CPUTarget       float64
	CPUActual       float64
	ConnectionStatus string
	PoolConnected   bool
	PoolAddress     string
	PoolPort        int
	WorkerName      string
	BTCAddress      string
	MockMining      bool
	DebugMode       bool
	ConfigValid     bool
	Uptime          time.Duration
	StartedAt       time.Time
	Screen          ScreenID
	Error           string
}

// WithHashRate returns a new copy of the AppState with the updated HashRate
// and the HashRateHistory rotated by 1 position (FIFO), inserting the new value
// at the end of the history array.
func (s AppState) WithHashRate(hps float64) AppState {
	s.HashRate = hps
	// Shift elements to the left by 1
	for i := 1; i < HashHistoryLen; i++ {
		s.HashRateHistory[i-1] = s.HashRateHistory[i]
	}
	// Insert new value at the end
	s.HashRateHistory[HashHistoryLen-1] = hps
	return s
}

// NextScreen returns a new copy of the AppState with the ScreenID advanced
// to the next screen cyclically (Dashboard -> Clock -> GlobalStats -> Dashboard).
func (s AppState) NextScreen() AppState {
	s.Screen = (s.Screen + 1) % NumScreens
	return s
}

// PrevScreen returns a new copy of the AppState with the ScreenID retreated
// to the previous screen cyclically.
func (s AppState) PrevScreen() AppState {
	if s.Screen == 0 {
		s.Screen = NumScreens - 1
	} else {
		s.Screen--
	}
	return s
}

// WithCPUTarget returns a new copy of the AppState with the CPUTarget adjusted
// by adding the specified delta, clamped to the [MinCPUTarget, MaxCPUTarget] range.
func (s AppState) WithCPUTarget(delta float64) AppState {
	newTarget := s.CPUTarget + delta
	if newTarget > MaxCPUTarget {
		newTarget = MaxCPUTarget
	} else if newTarget < MinCPUTarget {
		newTarget = MinCPUTarget
	}
	s.CPUTarget = newTarget
	return s
}

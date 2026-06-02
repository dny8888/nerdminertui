package worker

// HashRateMsg is emitted by the MinerWorker to report current mining speed and CPU usage.
type HashRateMsg struct {
	HPS       float64
	CPUActual float64
}

// ShareFoundMsg is emitted when a valid share is found and processed by the pool client.
type ShareFoundMsg struct {
	Accepted bool
}

// PoolStatsMsg is emitted by the PoolClient after fetching global network statistics.
type PoolStatsMsg struct {
	GlobalHashRate    float64
	NetworkDifficulty float64
	BlockHeight       int
}

// MinerErrorMsg is emitted when an internal error occurs within the MinerWorker.
type MinerErrorMsg struct {
	Err error
}

// PoolErrorMsg is emitted when there is a network error communicating with the pool.
type PoolErrorMsg struct {
	Err error
}

// ConnectionStatusMsg is emitted when TCP pool connection changes state.
type ConnectionStatusMsg struct {
	Status string
}

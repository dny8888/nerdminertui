package mining

// Job represents a mining job dispatched by the pool or locally mocked.
type Job struct {
	Header         []byte
	Target         [32]byte
	ExtraNonce     uint32
	Height         uint32
	// Stratum parameters needed for submitting the share
	JobID          string
	Extranonce2Hex string
	NtimeHex       string
}

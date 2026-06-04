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

	// Raw stratum parameters needed to rebuild the header for multi-worker entropy
	Coinb1Hex       string
	Coinb2Hex       string
	Extranonce1Hex  string
	MerkleBranchHex []string
	Extranonce2Size int
	VersionLE       []byte
	PrevhashLE      []byte
	NtimeLE         []byte
	NbitsLE         []byte
}

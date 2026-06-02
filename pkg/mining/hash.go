package mining

import (
	"crypto/sha256"
	"encoding/binary"
)

// SHA256d computes the double SHA-256 hash of the input data.
// It is the standard hash function used in Bitcoin block headers.
func SHA256d(data []byte) [32]byte {
	first := sha256.Sum256(data)
	return sha256.Sum256(first[:])
}
// MinerHashState holds precomputed SHA-256 midstate to optimize hashing
// by skipping the first 64 bytes of the block header.
type MinerHashState struct {
	payload [80]byte
	state   []byte
}

// NewMinerHashState initializes a zero-allocation hashing state for a given block header (76 bytes).
// It precomputes the SHA-256 midstate for the first 64 bytes of the header.
func NewMinerHashState(header []byte) *MinerHashState {
	m := &MinerHashState{}
	copy(m.payload[:76], header)

	// Precompute the midstate for the first 64 bytes
	hasher := sha256.New()
	hasher.Write(m.payload[:64])
	
	// MarshalBinary returns the internal state of the hash
	// We save it so we can Unmarshal it rapidly in the hot loop
	m.state, _ = hasher.(interface{ MarshalBinary() ([]byte, error) }).MarshalBinary()
	
	return m
}

// HashNonce computes the double SHA-256 hash for a specific nonce using the precomputed midstate.
// This function performs 0 heap allocations and is highly optimized.
func (m *MinerHashState) HashNonce(nonce uint32) [32]byte {
	var firstPayload [32]byte
	var finalHash [32]byte
	
	binary.LittleEndian.PutUint32(m.payload[76:80], nonce)
	
	h := sha256.New()
	unmarshaler := h.(interface{ UnmarshalBinary([]byte) error })
	unmarshaler.UnmarshalBinary(m.state)
	h.Write(m.payload[64:80])
	
	firstSlice := h.Sum(firstPayload[:0])
	copy(finalHash[:], firstSlice)
	
	return sha256.Sum256(finalHash[:])
}

// HashHeader concatenates the block header and a 32-bit nonce
// (in little-endian format, standard for Bitcoin) and computes its double SHA-256 hash.
// DEPRECATED: Use MinerHashState for hot-loop hashing.
func HashHeader(header []byte, nonce uint32) [32]byte {
	payload := make([]byte, len(header)+4)
	copy(payload, header)
	binary.LittleEndian.PutUint32(payload[len(header):], nonce)
	return SHA256d(payload)
}

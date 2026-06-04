package mining

import (
	"crypto/sha256"
	"encoding/binary"
	"hash"
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
	payload     [80]byte
	state       []byte
	hasher      hash.Hash
	hasher2     hash.Hash
	unmarshaler interface{ UnmarshalBinary([]byte) error }
	sumBuf      []byte
	sumBuf2     []byte
}

// NewMinerHashState initializes a zero-allocation hashing state for a given block header (76 bytes).
// It precomputes the SHA-256 midstate for the first 64 bytes of the header.
func NewMinerHashState(header []byte) *MinerHashState {
	m := &MinerHashState{}
	copy(m.payload[:76], header)

	// Precompute the midstate for the first 64 bytes
	h := sha256.New()
	h.Write(m.payload[:64])
	
	// MarshalBinary returns the internal state of the hash
	// We save it so we can Unmarshal it rapidly in the hot loop
	state, err := h.(interface{ MarshalBinary() ([]byte, error) }).MarshalBinary()
	if err != nil {
		panic("crypto/sha256: MarshalBinary failed: " + err.Error())
	}
	m.state = state
	
	// Create a dedicated hasher instance for the hot loop
	m.hasher = sha256.New()
	m.hasher2 = sha256.New()
	m.unmarshaler = m.hasher.(interface{ UnmarshalBinary([]byte) error })
	m.sumBuf = make([]byte, 0, 32)
	m.sumBuf2 = make([]byte, 0, 32)
	
	return m
}

// Reset re-initializes the hash state for a new block header, avoiding
// the allocation of a new MinerHashState object.
func (m *MinerHashState) Reset(header []byte) {
	copy(m.payload[:76], header)
	h := sha256.New()
	h.Write(m.payload[:64])
	state, err := h.(interface{ MarshalBinary() ([]byte, error) }).MarshalBinary()
	if err != nil {
		panic("crypto/sha256: MarshalBinary failed: " + err.Error())
	}
	m.state = state
}

// HashNonce computes the double SHA-256 hash for a specific nonce using the precomputed midstate.
// This function performs 0 heap allocations and is highly optimized.
func (m *MinerHashState) HashNonce(nonce uint32) [32]byte {
	var finalHash [32]byte
	
	binary.LittleEndian.PutUint32(m.payload[76:80], nonce)
	
	// Zero-allocation: reuse the pre-allocated hasher and unmarshaler
	m.unmarshaler.UnmarshalBinary(m.state)
	m.hasher.Write(m.payload[64:80])
	
	firstSlice := m.hasher.Sum(m.sumBuf[:0])
	
	m.hasher2.Reset()
	m.hasher2.Write(firstSlice)
	secondSlice := m.hasher2.Sum(m.sumBuf2[:0])
	
	copy(finalHash[:], secondSlice)
	
	return finalHash
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

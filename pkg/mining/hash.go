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

// HashHeader concatenates the block header and a 32-bit nonce
// (in little-endian format, standard for Bitcoin) and computes its double SHA-256 hash.
func HashHeader(header []byte, nonce uint32) [32]byte {
	// Bitcoin block headers are typically 80 bytes.
	// We append the 4-byte nonce to the header.
	payload := make([]byte, len(header)+4)
	copy(payload, header)
	binary.LittleEndian.PutUint32(payload[len(header):], nonce)
	return SHA256d(payload)
}

package mining

import (
	"encoding/hex"
	"fmt"
)

// ParseStratumJob creates a Job from mining.notify parameters.
func ParseStratumJob(jobID, prevhashHex, coinb1Hex, coinb2Hex, versionHex, nbitsHex, ntimeHex, extranonce1Hex string, extranonce2 uint32, extranonce2Size int, merkleBranchHex []string, poolDifficulty float64) (*Job, error) {
	// Format Extranonce2
	formatStr := fmt.Sprintf("%%0%dx", extranonce2Size*2)
	extranonce2Hex := fmt.Sprintf(formatStr, extranonce2)

	// Decode Header Fields
	versionBytes, err := hex.DecodeString(versionHex)
	if err != nil {
		return nil, fmt.Errorf("invalid version hex: %v", err)
	}
	prevhashBytes, err := hex.DecodeString(prevhashHex)
	if err != nil {
		return nil, fmt.Errorf("invalid prevhash hex: %v", err)
	}
	ntimeBytes, err := hex.DecodeString(ntimeHex)
	if err != nil {
		return nil, fmt.Errorf("invalid ntime hex: %v", err)
	}
	nbitsBytes, err := hex.DecodeString(nbitsHex)
	if err != nil {
		return nil, fmt.Errorf("invalid nbits hex: %v", err)
	}

	// Endianness handling for Stratum
	// Stratum sends prevhash, version, ntime, nbits in big-endian hex (usually),
	// but Bitcoin block headers require little-endian.
	versionLE := reverseBytes(versionBytes)
	ntimeLE := reverseBytes(ntimeBytes)
	nbitsLE := reverseBytes(nbitsBytes)
	
	// Prevhash from stratum is in "word-swapped big-endian" format:
	// 8 groups of 4 bytes, each group in big-endian order, groups in header order.
	// We need to reverse bytes WITHIN each 4-byte word (like cgminer's flip32/swab256),
	// NOT reverse all 32 bytes.
	prevhashLE := swapWords(prevhashBytes, 4)

	header, err := buildHeader(versionLE, prevhashLE, ntimeLE, nbitsLE, coinb1Hex, extranonce1Hex, extranonce2Hex, coinb2Hex, merkleBranchHex)
	if err != nil {
		return nil, err
	}

	// Target calculation from nBits (Simplified target logic if needed)
	// For now we will use the difficulty set by the pool.
	target := TargetFromDifficulty(poolDifficulty)

	return &Job{
		Header:          header,
		Target:          target,
		ExtraNonce:      extranonce2,
		Height:          0,
		JobID:           jobID,
		Extranonce2Hex:  extranonce2Hex,
		NtimeHex:        ntimeHex,
		Coinb1Hex:       coinb1Hex,
		Coinb2Hex:       coinb2Hex,
		Extranonce1Hex:  extranonce1Hex,
		MerkleBranchHex: merkleBranchHex,
		Extranonce2Size: extranonce2Size,
		VersionLE:       versionLE,
		PrevhashLE:      prevhashLE,
		NtimeLE:         ntimeLE,
		NbitsLE:         nbitsLE,
	}, nil
}

// RebuildHeaderWithExtraNonce2 rebuilds the block header with a new extranonce2.
// This allows different workers to have completely different Merkle Roots,
// avoiding any nonce collision without requiring range coordination.
func RebuildHeaderWithExtraNonce2(job *Job, newExtraNonce2Hex string) ([]byte, error) {
	return buildHeader(job.VersionLE, job.PrevhashLE, job.NtimeLE, job.NbitsLE, job.Coinb1Hex, job.Extranonce1Hex, newExtraNonce2Hex, job.Coinb2Hex, job.MerkleBranchHex)
}

func buildHeader(versionLE, prevhashLE, ntimeLE, nbitsLE []byte, coinb1Hex, extranonce1Hex, extranonce2Hex, coinb2Hex string, merkleBranchHex []string) ([]byte, error) {
	// Build Coinbase
	coinbaseHex := coinb1Hex + extranonce1Hex + extranonce2Hex + coinb2Hex
	coinbaseBytes, err := hex.DecodeString(coinbaseHex)
	if err != nil {
		return nil, fmt.Errorf("invalid coinbase hex: %v", err)
	}

	// Coinbase Hash
	coinbaseHash := SHA256d(coinbaseBytes)

	// Merkle Root
	merkleRoot := coinbaseHash
	payload := make([]byte, 64)
	for _, branchHex := range merkleBranchHex {
		branchBytes, err := hex.DecodeString(branchHex)
		if err != nil {
			return nil, fmt.Errorf("invalid merkle branch hex: %v", err)
		}
		// Hash(merkleRoot + branchBytes)
		copy(payload[0:32], merkleRoot[:])
		copy(payload[32:64], branchBytes)
		merkleRoot = SHA256d(payload)
	}

	// Build the 76-byte Header (without nonce)
	header := make([]byte, 76)
	copy(header[0:4], versionLE)
	copy(header[4:36], prevhashLE)
	copy(header[36:68], merkleRoot[:])
	copy(header[68:72], ntimeLE)
	copy(header[72:76], nbitsLE)

	return header, nil
}

func reverseBytes(data []byte) []byte {
	out := make([]byte, len(data))
	for i := 0; i < len(data); i++ {
		out[i] = data[len(data)-1-i]
	}
	return out
}

// swapWords reverses the bytes within each word of the given word size.
// This is equivalent to cgminer's flip32/swab256 when wordSize=4.
// Example (wordSize=4): [A1 B1 C1 D1 A2 B2 C2 D2] → [D1 C1 B1 A1 D2 C2 B2 A2]
func swapWords(data []byte, wordSize int) []byte {
	out := make([]byte, len(data))
	for i := 0; i < len(data); i += wordSize {
		end := i + wordSize
		if end > len(data) {
			end = len(data)
		}
		word := data[i:end]
		for j := 0; j < len(word); j++ {
			out[i+j] = word[len(word)-1-j]
		}
	}
	return out
}

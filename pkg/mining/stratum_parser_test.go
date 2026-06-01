package mining

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// swapWords vs reverseBytes — the root cause of Bug #2
// ---------------------------------------------------------------------------

func TestSwapWords_4ByteWords(t *testing.T) {
	// Input: 2 words of 4 bytes each, in big-endian order
	// Word1: A1 B1 C1 D1   Word2: A2 B2 C2 D2
	input := []byte{0xA1, 0xB1, 0xC1, 0xD1, 0xA2, 0xB2, 0xC2, 0xD2}

	got := swapWords(input, 4)

	// Each word reversed individually, but word ORDER preserved
	expected := []byte{0xD1, 0xC1, 0xB1, 0xA1, 0xD2, 0xC2, 0xB2, 0xA2}
	assert.Equal(t, expected, got)
}

func TestSwapWords_DiffersFromReverseBytes(t *testing.T) {
	// This is the critical test: swapWords and reverseBytes produce DIFFERENT
	// results for data longer than 4 bytes. Using reverseBytes for prevhash
	// was the bug.
	input := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}

	swapped := swapWords(input, 4)
	reversed := reverseBytes(input)

	// swapWords: [04 03 02 01] [08 07 06 05]  (each word reversed)
	// reverse:   [08 07 06 05] [04 03 02 01]  (entire array reversed)
	assert.Equal(t, []byte{0x04, 0x03, 0x02, 0x01, 0x08, 0x07, 0x06, 0x05}, swapped)
	assert.Equal(t, []byte{0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01}, reversed)
	assert.NotEqual(t, swapped, reversed, "swapWords and reverseBytes must differ for >4 bytes")
}

func TestSwapWords_IdenticalToReverseBytesForSingleWord(t *testing.T) {
	// For exactly 4 bytes (a single word), swapWords and reverseBytes
	// produce the same result. This is why reverseBytes is correct for
	// version, ntime, and nbits (all 4-byte fields).
	input := []byte{0x20, 0x00, 0x00, 0x00}

	swapped := swapWords(input, 4)
	reversed := reverseBytes(input)

	assert.Equal(t, swapped, reversed, "For 4 bytes, swapWords == reverseBytes")
	assert.Equal(t, []byte{0x00, 0x00, 0x00, 0x20}, swapped)
}

func TestSwapWords_32BytePrevHash(t *testing.T) {
	// Real-world example: Stratum prevhash hex → header bytes
	// cgminer's flip32 operation on a 32-byte prevhash.
	prevhashHex := "852ab3acf6baeb51e883cc88f49ef03ae17ed8110009a5fb0000000000000000"
	prevhashBytes, _ := hex.DecodeString(prevhashHex)

	got := swapWords(prevhashBytes, 4)

	// Each 4-byte word reversed:
	//   852ab3ac → acb32a85
	//   f6baeb51 → 51ebbaf6
	//   e883cc88 → 88cc83e8
	//   f49ef03a → 3af09ef4
	//   e17ed811 → 11d87ee1
	//   0009a5fb → fba50900
	//   00000000 → 00000000
	//   00000000 → 00000000
	expectedHex := "acb32a8551ebbaf688cc83e83af09ef411d87ee1fba509000000000000000000"
	expected, _ := hex.DecodeString(expectedHex)
	assert.Equal(t, expected, got)
}

func TestSwapWords_EmptyInput(t *testing.T) {
	got := swapWords([]byte{}, 4)
	assert.Equal(t, []byte{}, got)
}

// ---------------------------------------------------------------------------
// reverseBytes — used for version, ntime, nbits (4-byte fields)
// ---------------------------------------------------------------------------

func TestReverseBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		{"4 bytes", []byte{0x20, 0x00, 0x00, 0x04}, []byte{0x04, 0x00, 0x00, 0x20}},
		{"1 byte", []byte{0xAA}, []byte{0xAA}},
		{"empty", []byte{}, []byte{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, reverseBytes(tt.input))
		})
	}
}

// ---------------------------------------------------------------------------
// TargetFromDifficulty
// ---------------------------------------------------------------------------

func TestTargetFromDifficulty(t *testing.T) {
	t.Run("difficulty 1 equals genesis target", func(t *testing.T) {
		target := TargetFromDifficulty(1.0)
		assert.Equal(t, genesisTargetBytes, target)
	})

	t.Run("difficulty 2 is half of genesis", func(t *testing.T) {
		target := TargetFromDifficulty(2.0)
		// genesis = 0x00000000FFFF...
		// diff 2  = 0x000000007FFF8000... (approximately half)
		assert.Equal(t, byte(0x00), target[0])
		assert.Equal(t, byte(0x00), target[1])
		assert.Equal(t, byte(0x00), target[2])
		assert.Equal(t, byte(0x00), target[3])
		assert.Equal(t, byte(0x7F), target[4])
		assert.Equal(t, byte(0xFF), target[5])
	})

	t.Run("difficulty 0 returns genesis", func(t *testing.T) {
		target := TargetFromDifficulty(0)
		assert.Equal(t, genesisTargetBytes, target)
	})

	t.Run("negative difficulty returns genesis", func(t *testing.T) {
		target := TargetFromDifficulty(-5.0)
		assert.Equal(t, genesisTargetBytes, target)
	})

	t.Run("roundtrip difficulty→target→difficulty", func(t *testing.T) {
		for _, diff := range []float64{1.0, 10.0, 100.0, 4096.0} {
			target := TargetFromDifficulty(diff)
			// A hash equal to the target should have roughly this difficulty
			gotDiff := DifficultyFromHash(target)
			assert.InDelta(t, diff, gotDiff, diff*0.01, "diff=%f", diff)
		}
	})
}

// ---------------------------------------------------------------------------
// ParseStratumJob — full integration from raw stratum hex to header
// ---------------------------------------------------------------------------

func TestParseStratumJob_HeaderLength(t *testing.T) {
	// Minimal valid inputs
	job, err := ParseStratumJob(
		"test_job_1",
		"0000000000000000000000000000000000000000000000000000000000000000", // prevhash (64 hex = 32 bytes)
		"01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff", // coinb1
		"ffffffff0100f2052a0100000043410496b538e853519c726a2c91e61ec11600ae1390813a627c66fb8be7947be63c52da7589379515d4e0a604f8141781e62294721166bf621e73a82cbf2342c858eeac00000000", // coinb2
		"20000000",    // version
		"1d00ffff",    // nbits
		"5bfc2e56",    // ntime
		"aabbccdd",    // extranonce1
		0,             // extranonce2
		4,             // extranonce2Size
		[]string{},    // no merkle branches
		1.0,           // difficulty
	)

	require.NoError(t, err)
	require.NotNil(t, job)
	assert.Len(t, job.Header, 76, "header must be exactly 76 bytes (nonce added during hashing)")
}

func TestParseStratumJob_VersionEndianness(t *testing.T) {
	job, err := ParseStratumJob(
		"version_test",
		"0000000000000000000000000000000000000000000000000000000000000000",
		"01",          // minimal coinb1
		"01",          // minimal coinb2
		"20000004",    // version in big-endian hex from stratum
		"1d00ffff",
		"5bfc2e56",
		"aa",
		0, 1, []string{}, 1.0,
	)
	require.NoError(t, err)

	// Version 0x20000004 in BE hex → LE bytes: 04 00 00 20
	assert.Equal(t, byte(0x04), job.Header[0])
	assert.Equal(t, byte(0x00), job.Header[1])
	assert.Equal(t, byte(0x00), job.Header[2])
	assert.Equal(t, byte(0x20), job.Header[3])
}

func TestParseStratumJob_PrevHashWordSwap(t *testing.T) {
	// Use a recognizable pattern: 8 words where we can verify word-swap
	// Word1=AABBCCDD Word2=11223344 ... rest zeros
	prevhash := "AABBCCDD11223344" + "0000000000000000" + "0000000000000000" + "0000000000000000"

	job, err := ParseStratumJob(
		"prevhash_test", prevhash,
		"01", "01",
		"00000001", "1d00ffff", "00000001",
		"aa", 0, 1, []string{}, 1.0,
	)
	require.NoError(t, err)

	// PrevHash in header is at bytes [4..35]
	// Word1 AABBCCDD → word-swapped → DDCCBBAA
	assert.Equal(t, byte(0xDD), job.Header[4])
	assert.Equal(t, byte(0xCC), job.Header[5])
	assert.Equal(t, byte(0xBB), job.Header[6])
	assert.Equal(t, byte(0xAA), job.Header[7])
	// Word2 11223344 → word-swapped → 44332211
	assert.Equal(t, byte(0x44), job.Header[8])
	assert.Equal(t, byte(0x33), job.Header[9])
	assert.Equal(t, byte(0x22), job.Header[10])
	assert.Equal(t, byte(0x11), job.Header[11])
}

func TestParseStratumJob_NtimeNbitsEndianness(t *testing.T) {
	job, err := ParseStratumJob(
		"ntime_test",
		"0000000000000000000000000000000000000000000000000000000000000000",
		"01", "01",
		"00000001",
		"1a0575ef",    // nbits in BE
		"5bfc2e56",    // ntime in BE
		"aa", 0, 1, []string{}, 1.0,
	)
	require.NoError(t, err)

	// nTime at offset 68..71: 5bfc2e56 → LE: 562efc5b
	assert.Equal(t, byte(0x56), job.Header[68])
	assert.Equal(t, byte(0x2e), job.Header[69])
	assert.Equal(t, byte(0xfc), job.Header[70])
	assert.Equal(t, byte(0x5b), job.Header[71])

	// nBits at offset 72..75: 1a0575ef → LE: ef75051a
	assert.Equal(t, byte(0xef), job.Header[72])
	assert.Equal(t, byte(0x75), job.Header[73])
	assert.Equal(t, byte(0x05), job.Header[74])
	assert.Equal(t, byte(0x1a), job.Header[75])
}

func TestParseStratumJob_CoinbaseMerkleRoot(t *testing.T) {
	// When there are no merkle branches, the merkle root is just SHA256d(coinbase).
	// coinbase = coinb1 + extranonce1 + extranonce2 + coinb2
	job1, err := ParseStratumJob(
		"merkle_test",
		"0000000000000000000000000000000000000000000000000000000000000000",
		"aabb",     // coinb1
		"ccdd",     // coinb2
		"00000001", "1d00ffff", "00000001",
		"ee",       // extranonce1
		0x01,       // extranonce2 = 1
		2,          // extranonce2Size = 2 → "0001"
		[]string{}, 1.0,
	)
	require.NoError(t, err)

	// Manually compute: coinbase = "aabb" + "ee" + "0001" + "ccdd"
	coinbaseHex := "aabbee0001ccdd"
	coinbaseBytes, _ := hex.DecodeString(coinbaseHex)
	expectedRoot := SHA256d(coinbaseBytes)

	// Merkle root sits at header[36..67]
	var gotRoot [32]byte
	copy(gotRoot[:], job1.Header[36:68])
	assert.Equal(t, expectedRoot, gotRoot, "merkle root must match SHA256d of coinbase when no branches")
}

func TestParseStratumJob_MerkleBranch(t *testing.T) {
	// With one merkle branch, the root should be SHA256d(coinbaseHash + branchHash)
	branchHex := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

	job, err := ParseStratumJob(
		"branch_test",
		"0000000000000000000000000000000000000000000000000000000000000000",
		"01", "02",
		"00000001", "1d00ffff", "00000001",
		"ff", 0, 1,
		[]string{branchHex}, 1.0,
	)
	require.NoError(t, err)

	// Manual: coinbase = "01" + "ff" + "00" + "02" = "01ff0002"
	coinbaseBytes, _ := hex.DecodeString("01ff0002")
	coinbaseHash := SHA256d(coinbaseBytes)

	branchBytes, _ := hex.DecodeString(branchHex)
	payload := make([]byte, 64)
	copy(payload[0:32], coinbaseHash[:])
	copy(payload[32:64], branchBytes)
	expectedRoot := SHA256d(payload)

	var gotRoot [32]byte
	copy(gotRoot[:], job.Header[36:68])
	assert.Equal(t, expectedRoot, gotRoot)
}

func TestParseStratumJob_Extranonce2Formatting(t *testing.T) {
	// extranonce2 = 1, size = 4 → should be "00000001"
	job, err := ParseStratumJob(
		"en2_test",
		"0000000000000000000000000000000000000000000000000000000000000000",
		"01", "02",
		"00000001", "1d00ffff", "00000001",
		"aa", 1, 4, []string{}, 1.0,
	)
	require.NoError(t, err)
	assert.Equal(t, "00000001", job.Extranonce2Hex)

	// extranonce2 = 255, size = 2 → should be "00ff"
	job2, err := ParseStratumJob(
		"en2_test2",
		"0000000000000000000000000000000000000000000000000000000000000000",
		"01", "02",
		"00000001", "1d00ffff", "00000001",
		"aa", 255, 2, []string{}, 1.0,
	)
	require.NoError(t, err)
	assert.Equal(t, "00ff", job2.Extranonce2Hex)
}

func TestParseStratumJob_InvalidHex(t *testing.T) {
	_, err := ParseStratumJob(
		"bad_hex",
		"0000000000000000000000000000000000000000000000000000000000000000",
		"ZZZZ",    // invalid hex
		"01",
		"00000001", "1d00ffff", "00000001",
		"aa", 0, 1, []string{}, 1.0,
	)
	assert.Error(t, err, "should fail on invalid hex in coinbase")
}

func TestParseStratumJob_PoolDifficulty(t *testing.T) {
	job, err := ParseStratumJob(
		"diff_test",
		"0000000000000000000000000000000000000000000000000000000000000000",
		"01", "02",
		"00000001", "1d00ffff", "00000001",
		"aa", 0, 1, []string{}, 4096.0,
	)
	require.NoError(t, err)

	expectedTarget := TargetFromDifficulty(4096.0)
	assert.Equal(t, expectedTarget, job.Target)
}

func TestParseStratumJob_HashHeaderProduces32Bytes(t *testing.T) {
	job, err := ParseStratumJob(
		"hash_test",
		"0000000000000000000000000000000000000000000000000000000000000000",
		"01", "02",
		"00000001", "1d00ffff", "00000001",
		"aa", 0, 1, []string{}, 1.0,
	)
	require.NoError(t, err)

	// Hash with several nonces — all should produce valid 32-byte hashes
	for nonce := uint32(0); nonce < 10; nonce++ {
		hash := HashHeader(job.Header, nonce)
		assert.Len(t, hash, 32)
		assert.NotEqual(t, [32]byte{}, hash, "hash should not be all zeros")
	}
}

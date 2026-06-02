package mining

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSHA256d(t *testing.T) {
	// Simple test vector for double sha256
	data := []byte("hello world")
	hash := SHA256d(data)
	
	// Double SHA256 of "hello world"
	// 1st: b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9
	// 2nd: bc62d4b80d9e36da29c16c5d4d9f11731f36052c72401a76c23c0fb5a9b74423
	expectedHex := "bc62d4b80d9e36da29c16c5d4d9f11731f36052c72401a76c23c0fb5a9b74423"
	expectedBytes, _ := hex.DecodeString(expectedHex)
	
	var expected [32]byte
	copy(expected[:], expectedBytes)
	
	assert.Equal(t, expected, hash)
}

func TestHashHeader(t *testing.T) {
	header := make([]byte, 80) // 80 bytes of 0
	nonce := uint32(1)
	
	hash := HashHeader(header, nonce)
	// Just verify it doesn't panic and returns a 32 byte array
	assert.Len(t, hash, 32)
	assert.NotEqual(t, [32]byte{}, hash)
}

func TestMeetsTarget(t *testing.T) {
	target := [32]byte{
		0x00, 0x00, 0xFF, 0x00,
	} // rest is 0

	t.Run("hash is strictly less", func(t *testing.T) {
		var hash [32]byte
		hash[31] = 0x00
		hash[30] = 0x00
		hash[29] = 0xFE
		hash[28] = 0xFF
		assert.True(t, MeetsTarget(hash, target))
	})

	t.Run("hash is exactly equal", func(t *testing.T) {
		var hash [32]byte
		for i := 0; i < 32; i++ {
			hash[31-i] = target[i]
		}
		assert.False(t, MeetsTarget(hash, target))
	})

	t.Run("hash is greater", func(t *testing.T) {
		var hash [32]byte
		hash[31] = 0x00
		hash[30] = 0x01
		hash[29] = 0x00
		hash[28] = 0x00
		assert.False(t, MeetsTarget(hash, target))
	})
	
	t.Run("all zeros", func(t *testing.T) {
		assert.False(t, MeetsTarget([32]byte{}, [32]byte{}))
	})
}

func TestDifficultyFromHash(t *testing.T) {
	// A hash that is exactly genesisTarget should be diff 1.0
	diff := DifficultyFromHash(genesisTargetBytes)
	assert.InDelta(t, 1.0, diff, 0.0001)

	// A hash that is exactly half of genesisTarget should be diff 2.0
	halfGenesis := genesisTargetBytes
	halfGenesis[4] = 127 // Half of 255
	halfGenesis[5] = 255
	diff2 := DifficultyFromHash(halfGenesis)
	assert.InDelta(t, 2.0, diff2, 0.1)

	// Zero hash should return 0 (avoid panic)
	diffZero := DifficultyFromHash([32]byte{})
	assert.Equal(t, 0.0, diffZero)
}

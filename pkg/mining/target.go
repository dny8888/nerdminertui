package mining

import (
	"math/big"
)

// MeetsTarget compares a hash against a target byte-by-byte (Big-Endian).
// A hash meets the target if its numerical value is strictly less than the target.
func MeetsTarget(hash, target [32]byte) bool {
	// The block hash (result of SHA256d) has its most significant byte at the end (hash[31])
	// when treated as a 256-bit integer, while the target array is strictly Big-Endian
	// (target[0] is the most significant byte).
	// We compare hash[31-i] with target[i].
	for i := 0; i < 32; i++ {
		if hash[31-i] < target[i] {
			return true
		}
		if hash[31-i] > target[i] {
			return false
		}
	}
	// If all bytes are equal, it's not strictly less
	return false
}

// GenesisTarget represents the standard Bitcoin Genesis Block target
// (Difficulty 1), which is 0x00000000FFFF0000000000000000000000000000000000000000000000000000.
var genesisTargetBytes = [32]byte{
	0, 0, 0, 0, 255, 255, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var genesisTargetInt *big.Int

func init() {
	genesisTargetInt = new(big.Int).SetBytes(genesisTargetBytes[:])
}

// DifficultyFromHash calculates the relative difficulty of a given hash
// compared to the Genesis Target. If hash is 0, it returns a maximum float64 value.
func DifficultyFromHash(hash [32]byte) float64 {
	hashInt := new(big.Int).SetBytes(hash[:])
	
	if hashInt.Sign() == 0 {
		return 0 // Avoid division by zero, though 0 hash is infinite diff.
	}

	// diff = genesisTarget / hash
	// Since we need float precision, we convert to big.Float
	genFloat := new(big.Float).SetInt(genesisTargetInt)
	hashFloat := new(big.Float).SetInt(hashInt)

	diffFloat := new(big.Float).Quo(genFloat, hashFloat)
	diff, _ := diffFloat.Float64()
	return diff
}

// TargetFromDifficulty calculates a target [32]byte given a difficulty.
func TargetFromDifficulty(diff float64) [32]byte {
	if diff <= 0 {
		return genesisTargetBytes
	}
	
	// target = genesisTarget / diff
	genFloat := new(big.Float).SetInt(genesisTargetInt)
	diffFloat := big.NewFloat(diff)
	
	targetFloat := new(big.Float).Quo(genFloat, diffFloat)
	targetInt, _ := targetFloat.Int(nil)
	
	var targetBytes [32]byte
	b := targetInt.Bytes()
	
	// Ensure the bytes are copied to the end of the 32-byte array (Big-Endian)
	start := 32 - len(b)
	if start >= 0 {
		copy(targetBytes[start:], b)
	} else {
		// If length > 32, copy the last 32 bytes (shouldn't happen for valid diff)
		copy(targetBytes[:], b[len(b)-32:])
	}
	
	return targetBytes
}

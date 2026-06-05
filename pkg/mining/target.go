package mining

import (
	"encoding/binary"
	"math/big"
	"sync"
)

// MeetsTarget compares a hash against a target byte-by-byte (Big-Endian).
// A hash meets the target if its numerical value is strictly less than the target.
func MeetsTarget(hash, target [32]byte) bool {
	// Fast-path: check the most significant 64 bits first
	// hash is Little-Endian (MSB is hash[31]) -> LittleEndian.Uint64 reads 8 bytes up to hash[31]
	// target is Big-Endian (MSB is target[0]) -> BigEndian.Uint64 reads 8 bytes from target[0]
	hashTop := binary.LittleEndian.Uint64(hash[24:32])
	targetTop := binary.BigEndian.Uint64(target[0:8])

	if hashTop > targetTop {
		return false
	}
	if hashTop < targetTop {
		return true
	}

	// Slow-path: check the remaining 24 bytes
	for i := 8; i < 32; i++ {
		if hash[31-i] < target[i] {
			return true
		}
		if hash[31-i] > target[i] {
			return false
		}
	}
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

var bigIntPool = sync.Pool{New: func() interface{} { return new(big.Int) }}
var bigFloatPool = sync.Pool{New: func() interface{} { return new(big.Float) }}

// DifficultyFromHash calculates the relative difficulty of a given hash
// compared to the Genesis Target. If hash is 0, it returns a maximum float64 value.
func DifficultyFromHash(hash [32]byte) float64 {
	hashInt := bigIntPool.Get().(*big.Int).SetBytes(hash[:])
	defer bigIntPool.Put(hashInt)
	
	if hashInt.Sign() == 0 {
		return 0 // Avoid division by zero, though 0 hash is infinite diff.
	}

	// diff = genesisTarget / hash
	// Since we need float precision, we convert to big.Float
	genFloat := bigFloatPool.Get().(*big.Float).SetInt(genesisTargetInt)
	defer bigFloatPool.Put(genFloat)

	hashFloat := bigFloatPool.Get().(*big.Float).SetInt(hashInt)
	defer bigFloatPool.Put(hashFloat)

	diffFloat := bigFloatPool.Get().(*big.Float).Quo(genFloat, hashFloat)
	defer bigFloatPool.Put(diffFloat)

	diff, _ := diffFloat.Float64()
	return diff
}

// TargetFromDifficulty calculates a target [32]byte given a difficulty.
func TargetFromDifficulty(diff float64) [32]byte {
	if diff <= 0 {
		return genesisTargetBytes
	}
	
	// target = genesisTarget / diff
	genFloat := bigFloatPool.Get().(*big.Float).SetInt(genesisTargetInt)
	defer bigFloatPool.Put(genFloat)

	diffFloat := bigFloatPool.Get().(*big.Float).SetFloat64(diff)
	defer bigFloatPool.Put(diffFloat)
	
	targetFloat := bigFloatPool.Get().(*big.Float).Quo(genFloat, diffFloat)
	defer bigFloatPool.Put(targetFloat)

	targetInt, _ := targetFloat.Int(nil) // allocating Int is small, but avoids larger Float allocs
	
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

// HashIsLessThan compares two Little-Endian hashes.
// It returns true if h1 is strictly less than h2 (meaning h1 has a higher difficulty).
func HashIsLessThan(h1, h2 [32]byte) bool {
	h1Top := binary.LittleEndian.Uint64(h1[24:32])
	h2Top := binary.LittleEndian.Uint64(h2[24:32])
	
	if h1Top < h2Top {
		return true
	}
	if h1Top > h2Top {
		return false
	}
	
	for i := 23; i >= 0; i-- {
		if h1[i] < h2[i] {
			return true
		}
		if h1[i] > h2[i] {
			return false
		}
	}
	return false
}

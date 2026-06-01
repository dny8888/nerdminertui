# pkg/mining, Technical Design

> Design specification for the `pkg/mining` module. Focuses on HOW the algorithms are designed.

## Interface

### Classes / Functions

| Symbol | Signature | Return | Observation |
|---------|-----------|---------|------------|
| `SHA256d` | `func SHA256d(data []byte) [32]byte` | `[32]byte` | Computes double SHA-256 rounds. 🟢 |
| `HashHeader` | `func HashHeader(header []byte, nonce uint32) [32]byte` | `[32]byte` | Wires header with nonce and hashes. 🟢 |
| `MeetsTarget` | `func MeetsTarget(hash, target [32]byte) bool` | `bool` | Byte-wise big-endian comparison. 🟢 |
| `DifficultyFromHash` | `func DifficultyFromHash(hash [32]byte) float64` | `float64` | Target genesis / hash_as_bigint. 🟢 |

## Main Flow
1. **Double Hash Loop (`SHA256d`)**:
   - Call standard library `sha256.Sum256(data)`. 🟢
   - Call standard library `sha256.Sum256` over the resulting first hash. 🟢
   - Return second `[32]byte` array. 🟢
2. **Wire Header and Nonce (`HashHeader`)**:
   - Merge `header` bytes and `nonce` (encoded as 4-byte little-endian or big-endian). 🟢
   - Execute `SHA256d` over the merged slice. 🟢
3. **Compare target (`MeetsTarget`)**:
   - Compare `hash` and `target` byte-by-byte from index `0` to `31` (Big-Endian). If `hash[i] < target[i]`, return true immediately. If `hash[i] > target[i]`, return false immediately. 🟢
4. **Compute Difficulty (`DifficultyFromHash`)**:
   - Convert target genesis limit to big integer.
   - Convert hash to big integer.
   - Return floating point division `diff = genesis / hash`. 🟢

## Dependencies
- `crypto/sha256`: Standard cryptographic library. 🟢
- `math/big`: Used for large number conversions in difficulty. 🟢

## Design Decisions Identified

| Decision | Evidence | Confidence |
|---------|---------------------|-----------|
| No state allocation | Functions are completely stateless. | 🟢 |
| Byte comparison loop | Direct array iteration instead of string parses. | 🟢 |

## Internal State
- This package is completely stateless. 🟢

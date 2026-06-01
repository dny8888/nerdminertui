# pkg/mining

> Requirements specification for the `pkg/mining` module. Focuses on WHAT the core cryptographic structures do, not how.

## Overview
The `pkg/mining` module provides pure cryptographic primitives and data structures to perform Bitcoin-style hashing and block header target checking. 🟢

## Responsibilities
- Calculate double-round SHA256 hashes (`SHA256d`). 🟢
- Determine if a given block candidate hash satisfies target difficulty boundaries. 🟢
- Calculate target difficulty scaling ratios relative to the original Bitcoin Genesis Block target. 🟢

## Business Rules
- **BR-01: SHA256d Double Rounding**: The hashing function must perform two consecutive iterations of SHA256 calculations over the block header payload to match Bitcoin protocol specification. 🟢
- **BR-02: Target Bound Verification**: A candidate hash satisfies the target difficulty bounds if its numeric value is strictly less than the target value (evaluated byte-by-byte big-endian).
  $$\text{MeetsTarget}(\text{hash}, \text{target}) \iff \text{hash} < \text{target}$$ 🟢
- **BR-03: Pure Math Purity**: All methods in this package must remain mathematically pure, having no structural dependencies or I/O side-effects. 🟢

## Functional Requirements

| ID | Requirement | Priority | Acceptance Criteria |
|----|-----------|-----------|-------------------|
| RF-01 | SHA256d Calculations | Must | Process byte slices and return double-hashed `[32]byte` arrays. 🟢 |
| RF-02 | Byte-Wise Big-Endian Comparison | Must | Validate `hash < target` boundaries properly. 🟢 |
| RF-03 | Difficulty Calculation Ratio | Must | Calculate relative difficulty accurately matching Bitcoin rules. 🟢 |

## Non-Functional Requirements

| Type | Inferred Requirement | Code Evidence | Confidence |
|------|--------------------|---------------------|-----------|
| Performance | Perform double hashing rounds under ~200ns | `hash.go` | 🟢 |
| Safety | 100% pure thread-safety (stateless) | `target.go` | 🟢 |

## Acceptance Criteria

```gherkin
Given a block candidate hash of value 0x00...01 and target of 0x00...02
When MeetsTarget(hash, target) is called
Then it returns true

Given a block candidate hash of value 0x00...03 and target of 0x00...02
When MeetsTarget(hash, target) is called
Then it returns false
```

## Priority (MoSCoW)

| Requirement | MoSCoW | Justification |
|-----------|--------|---------------|
| Criptographic SHA256d | Must | Critical engine, block processing cannot occur without hashing. 🟢 |
| Target Comparison logic | Must | Ensures found shares are structurally valid before submission. 🟢 |

## Code Traceability

| File | Function / Class | Coverage |
|---------|-----------------|-----------|
| `pkg/mining/hash.go` | `SHA256d` | 🟢 |
| `pkg/mining/target.go` | `MeetsTarget` | 🟢 |
| `pkg/mining/target.go` | `DifficultyFromHash` | 🟢 |

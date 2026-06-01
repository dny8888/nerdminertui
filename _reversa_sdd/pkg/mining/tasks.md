# pkg/mining, Implementation Tasks

> Implementation checklists for cryptographic primitives and targets.

## Prerequisites
- [ ] Go cryptography library standard imports are available. 🟢

## Tasks

- [ ] **T-01: Implement SHA256d Double Hashing**
  - Origin in Legacy: `nerdtui-spec.md:336` (§5.8)
  - Criteria of Done: Write `SHA256d` function returning double SHA-256 rounds.
  - Confidence: 🟢 CONFIRMADO

- [ ] **T-02: Implement HashHeader Wiring**
  - Origin in Legacy: `nerdtui-spec.md:337` (§5.8)
  - Criteria of Done: Map nonce into byte slice and execute double hashing.
  - Confidence: 🟢 CONFIRMADO

- [ ] **T-03: Implement MeetsTarget comparison**
  - Origin in Legacy: `nerdtui-spec.md:338` (§5.8)
  - Criteria of Done: Write byte-wise loop comparing big-endian arrays.
  - Confidence: 🟢 CONFIRMADO

- [ ] **T-04: Implement Difficulty Math calculation**
  - Origin in Legacy: `nerdtui-spec.md:339` (§5.8)
  - Criteria of Done: Execute ratio division using `math/big` packages.
  - Confidence: 🟢 CONFIRMADO

## Test Tasks

- [ ] **TT-01: Known NIST Double SHA256 vectors**
  - Verify `SHA256d` output matches official NIST double-hashing test vectors.
- [ ] **TT-02: MeetsTarget edge cases**
  - Assert that equal values return false, and hashes starting with leading zeros evaluate correctly.

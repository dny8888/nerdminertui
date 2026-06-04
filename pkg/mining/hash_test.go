package mining

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestMinerHashState_HashNonce(t *testing.T) {
	// A sample header (76 bytes)
	headerHex := "0100000081cd02ab7e569e8bcd9317e2fe99f2de44d49ab2b8851eb4ba29000000000000e320b6c2fffc8d750423db8b1eb942ae710e951ed797f7affc8892b0f1fc122bc7f5d74df2b9441a"
	header, _ := hex.DecodeString(headerHex)

	nonce := uint32(0x42a14695) // Some nonce

	expected := HashHeader(header, nonce)
	
	state := NewMinerHashState(header)
	actual := state.HashNonce(nonce)

	if !bytes.Equal(expected[:], actual[:]) {
		t.Errorf("MinerHashState produced different hash. Expected %x, got %x", expected, actual)
	}
}

func BenchmarkHashHeader(b *testing.B) {
	headerHex := "0100000081cd02ab7e569e8bcd9317e2fe99f2de44d49ab2b8851eb4ba29000000000000e320b6c2fffc8d750423db8b1eb942ae710e951ed797f7affc8892b0f1fc122bc7f5d74df2b9441a"
	header, _ := hex.DecodeString(headerHex)
	nonce := uint32(0x42a14695)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = HashHeader(header, nonce)
	}
}

func BenchmarkMinerHashState(b *testing.B) {
	headerHex := "0100000081cd02ab7e569e8bcd9317e2fe99f2de44d49ab2b8851eb4ba29000000000000e320b6c2fffc8d750423db8b1eb942ae710e951ed797f7affc8892b0f1fc122bc7f5d74df2b9441a"
	header, _ := hex.DecodeString(headerHex)
	nonce := uint32(0x42a14695)
	
	state := NewMinerHashState(header)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = state.HashNonce(nonce + uint32(i))
	}
}

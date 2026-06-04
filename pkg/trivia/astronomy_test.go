package trivia

import (
	"encoding/hex"
	"testing"
)

func TestGetRandomSpaceWordHex_ReturnsCorrectByteLength(t *testing.T) {
	lengths := []int{4, 8, 12}
	for _, l := range lengths {
		res := GetRandomSpaceWordHex(l)
		decoded, err := hex.DecodeString(res)
		if err != nil {
			t.Fatalf("expected valid hex, got err: %v", err)
		}
		if len(decoded) != l {
			t.Errorf("expected %d bytes, got %d bytes: '%s'", l, len(decoded), string(decoded))
		}
	}
}

func TestGetRandomSpaceWordHex_ZeroOrNegativeLength(t *testing.T) {
	if res := GetRandomSpaceWordHex(0); res != "" {
		t.Errorf("expected empty string for length 0, got %s", res)
	}
	if res := GetRandomSpaceWordHex(-5); res != "" {
		t.Errorf("expected empty string for negative length, got %s", res)
	}
}

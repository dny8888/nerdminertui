package trivia

import (
	"encoding/hex"
	"math/rand"
	"strings"
)

// SpaceWords contains a list of short astronomy-related words (4-8 letters)
// that fit perfectly inside the ExtraNonce2 field of the Stratum protocol.
var SpaceWords = []string{
	"MARS", "NOVA", "MOON", "STAR",
	"ORION", "SUN", "VENUS", "PLUTO",
	"COMET", "EARTH", "LUNA", "HALO",
	"APOLLO", "ORBIT", "ZENITH", "NEBULA",
	"PULSAR", "QUASAR", "SATURN", "URANUS",
	"COSMOS", "GALAXY", "METEOR", "ASTRO",
}

// GetRandomSpaceWordHex returns a random space-related word encoded as hex.
// It truncates the string to maxBytes length so it strictly adheres to
// the pool's ExtraNonce2 size requirements.
func GetRandomSpaceWordHex(maxBytes int) string {
	if maxBytes <= 0 {
		return ""
	}

	word := SpaceWords[rand.Intn(len(SpaceWords))]
	
	// Truncate or pad with space character to match exactly maxBytes
	if len(word) > maxBytes {
		word = word[:maxBytes]
	} else if len(word) < maxBytes {
		word += strings.Repeat(" ", maxBytes-len(word))
	}

	return hex.EncodeToString([]byte(word))
}

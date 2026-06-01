package format

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormatHashRate(t *testing.T) {
	assert.Equal(t, "0 H/s", FormatHashRate(0))
	assert.Equal(t, "500 H/s", FormatHashRate(500.4))
	assert.Equal(t, "1.5 KH/s", FormatHashRate(1520.0))
	assert.Equal(t, "2.0 MH/s", FormatHashRate(2000000.0))
}

func TestFormatUptime(t *testing.T) {
	// Negative edge case
	assert.Equal(t, "0m 00s", FormatUptime(-5*time.Second))

	// Less than an hour
	assert.Equal(t, "14m 05s", FormatUptime(14*time.Minute+5*time.Second))

	// 2 days, 3 hours, 14 mins (from acceptance criteria)
	d := 2*24*time.Hour + 3*time.Hour + 14*time.Minute
	assert.Equal(t, "2d 03h 14m", FormatUptime(d))
}

func TestFormatDifficulty(t *testing.T) {
	assert.Equal(t, "0", FormatDifficulty(0))
	assert.Equal(t, "12.35", FormatDifficulty(12.345))
	assert.Equal(t, "1.23e+03", FormatDifficulty(1234.5))
}

func TestFormatBlockHeight(t *testing.T) {
	assert.Equal(t, "#0", FormatBlockHeight(0))
	assert.Equal(t, "#123", FormatBlockHeight(123))
	assert.Equal(t, "#1.234", FormatBlockHeight(1234))
	assert.Equal(t, "#892.441", FormatBlockHeight(892441))
	assert.Equal(t, "#1.234.567", FormatBlockHeight(1234567))
}

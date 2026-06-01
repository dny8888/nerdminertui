package format

import (
	"fmt"
)

// FormatDifficulty formats a difficulty float to a readable scientific or decimal string.
func FormatDifficulty(d float64) string {
	if d == 0 {
		return "0"
	}
	if d >= 1000.0 {
		return fmt.Sprintf("%.2e", d)
	}
	return fmt.Sprintf("%.2f", d)
}

// FormatBlockHeight formats a block height with dot separators for thousands.
func FormatBlockHeight(h uint32) string {
	s := fmt.Sprintf("%d", h)
	var result []byte
	for i := 0; i < len(s); i++ {
		// Insert dot before every 3rd digit from the end
		if i > 0 && (len(s)-i)%3 == 0 {
			result = append(result, '.')
		}
		result = append(result, s[i])
	}
	return "#" + string(result)
}

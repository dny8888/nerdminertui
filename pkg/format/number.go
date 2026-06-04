package format

import (
	"fmt"
	"strings"
)

// FormatWithThousandSeparator formats a large float64 value into a string with
// a dot '.' as the thousand separator, removing decimal places.
// Useful for representing large numbers like total hashes or years.
func FormatWithThousandSeparator(value float64) string {
	str := fmt.Sprintf("%.0f", value)
	if len(str) <= 3 {
		return str
	}

	var parts []string
	for i := len(str); i > 0; i -= 3 {
		start := i - 3
		if start < 0 {
			start = 0
		}
		parts = append([]string{str[start:i]}, parts...)
	}
	return strings.Join(parts, ".")
}

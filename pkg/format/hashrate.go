package format

import "fmt"

// FormatHashRate converts a raw hashrate float into a scaled human-readable string.
func FormatHashRate(hps float64) string {
	if hps == 0 {
		return "0 H/s"
	}
	if hps >= 1000000.0 {
		return fmt.Sprintf("%.1f MH/s", hps/1000000.0)
	}
	if hps >= 1000.0 {
		return fmt.Sprintf("%.1f KH/s", hps/1000.0)
	}
	return fmt.Sprintf("%.0f H/s", hps)
}

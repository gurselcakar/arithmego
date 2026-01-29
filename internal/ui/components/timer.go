package components

import (
	"fmt"
	"time"
)

// FormatTimer formats a duration as MM:SS for display.
func FormatTimer(remaining time.Duration) string {
	if remaining < 0 {
		remaining = 0
	}
	minutes := int(remaining.Minutes())
	seconds := int(remaining.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

package components

import (
	"strings"

	"github.com/gurselcakar/arithmego/internal/ui/styles"
)

const (
	// ProgressBarFilledChar is the character for filled portion of progress bar.
	ProgressBarFilledChar = "█"
	// ProgressBarEmptyChar is the character for empty portion of progress bar.
	ProgressBarEmptyChar = "░"
)

// renderProgressBar renders a progress bar at the given width.
// Accuracy should be 0-100.
func renderProgressBar(accuracy float64, width int) string {
	if width <= 0 {
		return ""
	}

	filled := int(accuracy / 100 * float64(width))
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}
	empty := width - filled

	return strings.Repeat(ProgressBarFilledChar, filled) + strings.Repeat(ProgressBarEmptyChar, empty)
}

// RenderProgressBarColored renders a progress bar with color based on accuracy.
// High accuracy (>=80%) is green, medium (60-79%) is default, low (<60%) is red.
func RenderProgressBarColored(accuracy float64, width int) string {
	bar := renderProgressBar(accuracy, width)

	if accuracy >= 80 {
		return styles.Correct.Render(bar)
	} else if accuracy < 60 {
		return styles.Incorrect.Render(bar)
	}
	return bar // Default color for medium accuracy
}

// ProgressBarWidth returns the appropriate progress bar width based on terminal width.
func ProgressBarWidth(termWidth int) int {
	switch {
	case termWidth < 60:
		return 10
	case termWidth < 80:
		return 16
	default:
		return 20
	}
}

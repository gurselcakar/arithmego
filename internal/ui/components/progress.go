package components

import (
	"strings"

	"github.com/gurselcakar/arithmego/internal/ui/styles"
)

// ProgressDots returns a visual representation of progress using filled and empty dots.
// Example: ProgressDots(2, 4) returns "● ● ○ ○"
func ProgressDots(current, total int) string {
	if total <= 0 {
		return ""
	}
	if current < 0 {
		current = 0
	}
	if current > total {
		current = total
	}

	var parts []string
	for i := 0; i < total; i++ {
		if i < current {
			parts = append(parts, "●")
		} else {
			parts = append(parts, "○")
		}
	}
	return strings.Join(parts, " ")
}

// ProgressDotsColored returns progress dots with active dots colored.
// Active (filled) dots use the logo color, inactive dots are dim.
func ProgressDotsColored(current, total int) string {
	if total <= 0 {
		return ""
	}
	if current < 0 {
		current = 0
	}
	if current > total {
		current = total
	}

	var parts []string
	for i := 0; i < total; i++ {
		if i < current {
			parts = append(parts, styles.Logo.Render("●"))
		} else {
			parts = append(parts, styles.Dim.Render("○"))
		}
	}
	return strings.Join(parts, " ")
}

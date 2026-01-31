package components

import (
	"strings"

	"github.com/gurselcakar/arithmego/internal/ui/styles"
)

// RenderHints renders navigation hints at the bottom of screens.
func RenderHints(hints []string) string {
	return styles.Dim.Render(strings.Join(hints, "  "))
}

// Hint represents a key-action pair for navigation hints.
type Hint struct {
	Key    string // The key to press (e.g., "Enter", "S", "↑↓")
	Action string // The action description (e.g., "Continue", "Skip")
}

// RenderHintsStructured renders hints with keys inside brackets.
// Hints are rendered in the order provided, so caller should order them:
// back actions on left, navigation in middle, forward actions on right.
func RenderHintsStructured(hints []Hint) string {
	var parts []string
	for _, h := range hints {
		parts = append(parts, "["+h.Key+"] "+h.Action)
	}
	return styles.Dim.Render(strings.Join(parts, "    "))
}

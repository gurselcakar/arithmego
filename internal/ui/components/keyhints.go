package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/gurselcakar/arithmego/internal/ui/styles"
)

// HintsHeight is the height reserved for the hints section at the bottom of viewport screens.
const HintsHeight = 3

// Hint represents a key-action pair for navigation hints.
type Hint struct {
	Key    string // The key to press (e.g., "Enter", "S", "↑↓")
	Action string // The action description (e.g., "Continue", "Skip")
}

// RenderHintsResponsive renders hints with width-aware wrapping.
// If all hints fit on one line within the given width, they are rendered
// on a single line. Otherwise, they are split into two centered rows
// at the most balanced split point.
func RenderHintsResponsive(hints []Hint, width int) string {
	gap := "    "
	parts := make([]string, len(hints))
	for i, h := range hints {
		parts[i] = "[" + h.Key + "] " + h.Action
	}

	single := strings.Join(parts, gap)
	if lipgloss.Width(single) <= width {
		return styles.Dim.Render(single)
	}

	// Find the split point that minimizes the width difference between two rows.
	bestSplit := 1
	bestDiff := -1
	for i := 1; i < len(parts); i++ {
		row1 := strings.Join(parts[:i], gap)
		row2 := strings.Join(parts[i:], gap)
		diff := lipgloss.Width(row1) - lipgloss.Width(row2)
		if diff < 0 {
			diff = -diff
		}
		if bestDiff < 0 || diff < bestDiff {
			bestDiff = diff
			bestSplit = i
		}
	}

	row1 := strings.Join(parts[:bestSplit], gap)
	row2 := strings.Join(parts[bestSplit:], gap)

	centered1 := lipgloss.PlaceHorizontal(width, lipgloss.Center, row1)
	centered2 := lipgloss.PlaceHorizontal(width, lipgloss.Center, row2)

	return styles.Dim.Render(centered1 + "\n" + centered2)
}

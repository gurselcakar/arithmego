package screens

import (
	"github.com/gurselcakar/arithmego/internal/game"
)

// maxLen returns the length of the longest string in the slice.
func maxLen(items []string) int {
	max := 0
	for _, s := range items {
		if len(s) > max {
			max = len(s)
		}
	}
	return max
}

// difficultyNames returns the display names for all difficulties.
func difficultyNames(diffs []game.Difficulty) []string {
	names := make([]string, len(diffs))
	for i, d := range diffs {
		names[i] = d.String()
	}
	return names
}

// findDifficultyIndex finds the index of a difficulty by name.
// Falls back to Medium, then to middle index if not found.
func findDifficultyIndex(name string) int {
	diffs := game.AllDifficulties()
	for i, d := range diffs {
		if d.String() == name {
			return i
		}
	}
	// Fallback: find Medium explicitly
	for i, d := range diffs {
		if d == game.Medium {
			return i
		}
	}
	// Ultimate fallback: middle index
	if len(diffs) > 0 {
		return len(diffs) / 2
	}
	return 0
}

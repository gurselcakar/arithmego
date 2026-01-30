package components

import (
	"fmt"

	"github.com/gurselcakar/arithmego/internal/ui/styles"
)

// SelectorOptions configures the horizontal selector rendering.
type SelectorOptions struct {
	Label   string // Optional label (e.g., "Difficulty")
	Prefix  string // Prefix before arrows (e.g., "  " for indentation)
	Focused bool   // Whether the selector is focused
}

// RenderSelector renders a horizontal selector with arrows: ◀ value ▶
// When Label is provided: "Label: ◀ Value ▶"
// When only Prefix is provided: "  ◀ Value ▶"
// Panics if options is empty or index is out of bounds.
func RenderSelector(index int, options []string, opts SelectorOptions) string {
	if len(options) == 0 {
		panic("RenderSelector: options cannot be empty")
	}
	if index < 0 || index >= len(options) {
		panic(fmt.Sprintf("RenderSelector: index %d out of bounds for %d options", index, len(options)))
	}

	leftArrow := "◀"
	rightArrow := "▶"

	if index == 0 {
		leftArrow = " "
	}
	if index >= len(options)-1 {
		rightArrow = " "
	}

	value := options[index]

	if opts.Label != "" {
		// Format: Label: ◀ Value ▶
		if opts.Focused {
			return fmt.Sprintf("%s: %s %s %s",
				styles.Bold.Render(opts.Label),
				styles.Accent.Render(leftArrow),
				styles.Selected.Render(value),
				styles.Accent.Render(rightArrow),
			)
		}
		return fmt.Sprintf("%s: %s %s %s",
			styles.Subtle.Render(opts.Label),
			styles.Subtle.Render(leftArrow),
			styles.Unselected.Render(value),
			styles.Subtle.Render(rightArrow),
		)
	}

	// Format: [Prefix]◀ Value ▶
	if opts.Focused {
		return fmt.Sprintf("%s%s %s %s",
			opts.Prefix,
			styles.Accent.Render(leftArrow),
			styles.Selected.Render(value),
			styles.Accent.Render(rightArrow),
		)
	}
	return fmt.Sprintf("%s%s %s %s",
		opts.Prefix,
		styles.Dim.Render(leftArrow),
		styles.Unselected.Render(value),
		styles.Dim.Render(rightArrow),
	)
}

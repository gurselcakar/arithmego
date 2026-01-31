package components

import (
	"fmt"

	"github.com/gurselcakar/arithmego/internal/ui/styles"
)

// SelectorOptions configures the horizontal selector rendering.
type SelectorOptions struct {
	Label      string // Optional label (e.g., "Difficulty")
	LabelWidth int    // Width for label alignment (0 = no padding)
	ValueWidth int    // Width for value alignment (0 = no padding)
	Prefix     string // Prefix before arrows (e.g., "  " for indentation)
	Focused    bool   // Whether the selector is focused
}

// RenderSelector renders a horizontal selector with arrows: ◀ value ▶
// Arrows are always visible but dimmed at boundaries to indicate navigation limits.
// When Label is provided: "Label  ◀ Value ▶"
// When only Prefix is provided: "  ◀ Value ▶"
func RenderSelector(index int, options []string, opts SelectorOptions) string {
	if len(options) == 0 {
		return styles.Dim.Render("[no options]")
	}
	if index < 0 || index >= len(options) {
		// Clamp index to valid range for graceful degradation
		if index < 0 {
			index = 0
		} else {
			index = len(options) - 1
		}
	}

	// Track if arrows are active (can navigate in that direction)
	leftActive := index > 0
	rightActive := index < len(options)-1

	value := options[index]

	// Helper to render arrow with appropriate style
	renderArrow := func(arrow string, active bool, focused bool) string {
		if focused {
			if active {
				return styles.Accent.Render(arrow)
			}
			return styles.Dim.Render(arrow)
		}
		if active {
			return styles.Subtle.Render(arrow)
		}
		return styles.Dim.Render(arrow)
	}

	if opts.Label != "" {
		label := opts.Label
		if opts.LabelWidth > 0 {
			label = fmt.Sprintf("%-*s", opts.LabelWidth, opts.Label)
		}

		// Calculate padding to add after the selector for fixed total width
		// Format: ◀ Value ▶ then padding to reach ValueWidth
		padding := ""
		if opts.ValueWidth > 0 {
			// We want the value area to be ValueWidth, so pad after ▶
			padLen := opts.ValueWidth - len(value)
			if padLen > 0 {
				padding = fmt.Sprintf("%*s", padLen, "")
			}
		}

		// Format: Label  ◀ Value ▶[padding] (arrows hug the value)
		if opts.Focused {
			return fmt.Sprintf("%s  %s %s %s%s",
				styles.Bold.Render(label),
				renderArrow("◀", leftActive, true),
				styles.Selected.Render(value),
				renderArrow("▶", rightActive, true),
				padding,
			)
		}
		return fmt.Sprintf("%s  %s %s %s%s",
			styles.Subtle.Render(label),
			renderArrow("◀", leftActive, false),
			styles.Unselected.Render(value),
			renderArrow("▶", rightActive, false),
			padding,
		)
	}

	// Format: [Prefix]◀ Value ▶
	if opts.Focused {
		return fmt.Sprintf("%s%s %s %s",
			opts.Prefix,
			renderArrow("◀", leftActive, true),
			styles.Selected.Render(value),
			renderArrow("▶", rightActive, true),
		)
	}
	return fmt.Sprintf("%s%s %s %s",
		opts.Prefix,
		renderArrow("◀", leftActive, false),
		styles.Unselected.Render(value),
		renderArrow("▶", rightActive, false),
	)
}

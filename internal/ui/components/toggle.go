package components

import (
	"fmt"

	"github.com/gurselcakar/arithmego/internal/ui/styles"
)

// ToggleOptions configures the toggle rendering.
type ToggleOptions struct {
	Label      string // Label text (e.g., "Auto-update")
	LabelWidth int    // Width for label alignment (0 = no padding)
	Focused    bool   // Whether the toggle is focused
}

// RenderToggle renders a toggle switch: Label    On · Off
// The active state is highlighted, inactive is dimmed.
func RenderToggle(value bool, opts ToggleOptions) string {
	label := opts.Label
	if opts.LabelWidth > 0 {
		label = fmt.Sprintf("%-*s", opts.LabelWidth, opts.Label)
	}

	var onText, offText, separator string
	if value {
		// On is active
		if opts.Focused {
			onText = styles.Selected.Render("On")
		} else {
			onText = styles.Unselected.Render("On")
		}
		offText = styles.Dim.Render("Off")
	} else {
		// Off is active
		onText = styles.Dim.Render("On")
		if opts.Focused {
			offText = styles.Selected.Render("Off")
		} else {
			offText = styles.Unselected.Render("Off")
		}
	}
	separator = styles.Dim.Render(" · ")

	if opts.Focused {
		return fmt.Sprintf("%s  %s%s%s",
			styles.Bold.Render(label),
			onText, separator, offText,
		)
	}
	return fmt.Sprintf("%s  %s%s%s",
		styles.Subtle.Render(label),
		onText, separator, offText,
	)
}

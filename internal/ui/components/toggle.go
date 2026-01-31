package components

import (
	"fmt"

	"github.com/gurselcakar/arithmego/internal/ui/styles"
)

// ToggleOptions configures the toggle rendering.
type ToggleOptions struct {
	Label   string // Label text (e.g., "Auto-update")
	Focused bool   // Whether the toggle is focused
}

// RenderToggle renders a toggle switch: Label [ON] or Label [OFF]
func RenderToggle(value bool, opts ToggleOptions) string {
	var stateText string
	if value {
		stateText = "[ON]"
	} else {
		stateText = "[OFF]"
	}

	if opts.Focused {
		return fmt.Sprintf("%s %s",
			styles.Bold.Render(opts.Label),
			styles.Selected.Render(stateText),
		)
	}
	return fmt.Sprintf("%s %s",
		styles.Subtle.Render(opts.Label),
		styles.Unselected.Render(stateText),
	)
}

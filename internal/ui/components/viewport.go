package components

import "github.com/charmbracelet/bubbles/viewport"

// SetViewportSize initializes or resizes a viewport.
func SetViewportSize(vp *viewport.Model, ready *bool, width, height int) {
	if !*ready {
		*vp = viewport.New(width, height)
		vp.YPosition = 0
		*ready = true
	} else {
		vp.Width = width
		vp.Height = height
	}
}

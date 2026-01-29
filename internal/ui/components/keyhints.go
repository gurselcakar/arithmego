package components

import (
	"strings"

	"github.com/gurselcakar/arithmego/internal/ui/styles"
)

// RenderHints renders navigation hints at the bottom of screens.
func RenderHints(hints []string) string {
	return styles.Dim.Render(strings.Join(hints, "  "))
}

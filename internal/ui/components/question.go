package components

import "github.com/gurselcakar/arithmego/internal/ui/styles"

// RenderQuestion renders a question prominently for display.
func RenderQuestion(display string) string {
	return styles.Bold.Render(display + " = ")
}

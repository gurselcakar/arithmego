package components

import "github.com/gurselcakar/arithmego/internal/ui/styles"

// RenderQuestion renders a question prominently for display.
// For typing mode, the "=" is shown as the input prompt.
// For multiple choice mode, use RenderQuestionWithAnswer to append "= ?".
func RenderQuestion(display string) string {
	return styles.Bold.Render(display)
}

// RenderQuestionWithAnswer renders a question with "= ?" suffix for multiple choice mode.
func RenderQuestionWithAnswer(display string) string {
	return styles.Bold.Render(display + " = ?")
}

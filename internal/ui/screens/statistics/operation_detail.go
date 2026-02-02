package statistics

import (
	"fmt"
	"strings"

	"github.com/gurselcakar/arithmego/internal/analytics"
	"github.com/gurselcakar/arithmego/internal/ui/components"
	"github.com/gurselcakar/arithmego/internal/ui/styles"
)

// DifficultyOrder defines the order for displaying difficulties.
var DifficultyOrder = []string{"Beginner", "Easy", "Medium", "Hard", "Expert"}

// RenderOperationDetailContent renders the operation detail view content for viewport.
func RenderOperationDetailContent(
	operation string,
	extStats analytics.ExtendedOperationStats,
	mistakes []analytics.RecentMistake,
	difficultyFilter string,
	timePeriod analytics.TimePeriod,
	width int,
) string {
	var b strings.Builder

	// Title
	title := fmt.Sprintf("STATISTICS · %s", strings.ToUpper(operation))
	b.WriteString(styles.Bold.Render(title))
	b.WriteString("\n\n")

	// Inline filters
	diffDisplay := "All Difficulties"
	if difficultyFilter != "" {
		diffDisplay = difficultyFilter
	}
	filterLine := fmt.Sprintf("Filter: ◀ %s ▶       ◀ %s ▶", diffDisplay, timePeriod.String())
	b.WriteString(styles.Dim.Render(filterLine))
	b.WriteString("\n\n")

	// Separator
	separatorWidth := 56
	if width > 0 && width-10 < separatorWidth {
		separatorWidth = width - 10
	}
	b.WriteString(styles.Dim.Render(strings.Repeat("─", separatorWidth)))
	b.WriteString("\n\n")

	// Summary section
	b.WriteString(styles.Bold.Render("SUMMARY"))
	b.WriteString("\n")
	b.WriteString(styles.Dim.Render("───────"))
	b.WriteString("\n")

	labelWidth := 16
	valueWidth := 18 // Fixed width for consistent alignment
	b.WriteString(fmt.Sprintf("%-*s %-*d\n", labelWidth, "Total Questions", valueWidth, extStats.Total))
	b.WriteString(fmt.Sprintf("%-*s %-*s\n", labelWidth, "Correct", valueWidth,
		fmt.Sprintf("%d  (%s)", extStats.Correct, FormatAccuracyPlain(extStats.Accuracy))))

	if extStats.AvgResponseTimeMs > 0 {
		b.WriteString(fmt.Sprintf("%-*s %-*s\n", labelWidth, "Avg Response", valueWidth, FormatResponseTime(extStats.AvgResponseTimeMs)))
	}
	if extStats.FastestTimeMs > 0 {
		b.WriteString(fmt.Sprintf("%-*s %-*s\n", labelWidth, "Fastest", valueWidth, FormatResponseTime(extStats.FastestTimeMs)))
	}
	b.WriteString("\n")

	// By difficulty section
	if len(extStats.ByDifficulty) > 0 {
		b.WriteString(styles.Bold.Render("BY DIFFICULTY"))
		b.WriteString("\n")
		b.WriteString(styles.Dim.Render("─────────────"))
		b.WriteString("\n")

		barWidth := components.ProgressBarWidth(width)

		for _, diff := range DifficultyOrder {
			if stats, ok := extStats.ByDifficulty[diff]; ok && stats.Total > 0 {
				accStr := fmt.Sprintf("%3.0f%%", stats.Accuracy)
				bar := components.RenderProgressBarColored(stats.Accuracy, barWidth)
				counts := fmt.Sprintf("(%d/%d)", stats.Correct, stats.Total)

				b.WriteString(fmt.Sprintf("%-10s %s  %s  %s\n", diff, accStr, bar, counts))
			}
		}
		b.WriteString("\n")
	}

	// Recent mistakes section
	b.WriteString(styles.Bold.Render(fmt.Sprintf("RECENT MISTAKES (%d)", len(mistakes))))
	b.WriteString("\n")
	b.WriteString(styles.Dim.Render("───────────────────"))
	b.WriteString("\n")

	if len(mistakes) == 0 {
		b.WriteString(styles.Correct.Render("No mistakes yet - perfect!"))
		b.WriteString("\n")
	} else {
		for _, m := range mistakes {
			// Format: "15 + 8 = 23   →  You: 22    1.8s ago"
			mistakeLine := fmt.Sprintf("%s   →  You: %d    %s",
				m.Question,
				m.UserAnswer,
				FormatRelativeTime(m.SessionDate),
			)
			b.WriteString(styles.Incorrect.Render(mistakeLine))
			b.WriteString("\n")
		}
	}

	return b.String()
}

// RenderOperationReviewContent renders the review all mistakes view content for viewport.
func RenderOperationReviewContent(
	operation string,
	allMistakes []analytics.RecentMistake,
	width int,
) string {
	var b strings.Builder

	// Title
	title := fmt.Sprintf("STATISTICS · %s · REVIEW", strings.ToUpper(operation))
	b.WriteString(styles.Bold.Render(title))
	b.WriteString("\n")
	b.WriteString(styles.Dim.Render("All Mistakes"))
	b.WriteString("\n\n")

	// Separator
	separatorWidth := 56
	if width > 0 && width-10 < separatorWidth {
		separatorWidth = width - 10
	}
	b.WriteString(styles.Dim.Render(strings.Repeat("─", separatorWidth)))
	b.WriteString("\n\n")

	if len(allMistakes) == 0 {
		b.WriteString(styles.Correct.Render("No mistakes - perfect!"))
		b.WriteString("\n")
	} else {
		// Column header
		headerLine := fmt.Sprintf("  #   %-25s  %-8s  %-8s  %s", "Question", "You", "Correct", "When")
		b.WriteString(styles.Dim.Render(headerLine))
		b.WriteString("\n")
		b.WriteString(styles.Dim.Render(strings.Repeat("─", separatorWidth)))
		b.WriteString("\n")

		// Render all mistakes - viewport handles scrolling
		for i, m := range allMistakes {
			// Truncate question if needed
			question := m.Question
			if len(question) > 22 {
				question = question[:19] + "..."
			}

			line := fmt.Sprintf(" %3d  %-25s  %-8d  %-8d  %s",
				i+1,
				question,
				m.UserAnswer,
				m.CorrectAnswer,
				FormatRelativeTime(m.SessionDate),
			)
			b.WriteString(styles.Incorrect.Render(line))
			b.WriteString("\n")
		}

		// Total info
		b.WriteString("\n")
		b.WriteString(styles.Dim.Render(fmt.Sprintf("Total: %d mistakes", len(allMistakes))))
	}

	return b.String()
}

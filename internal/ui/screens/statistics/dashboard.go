package statistics

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/gurselcakar/arithmego/internal/analytics"
	"github.com/gurselcakar/arithmego/internal/ui/components"
	"github.com/gurselcakar/arithmego/internal/ui/styles"
)


// RenderDashboardContent generates the dashboard content without layout wrapper.
// This is used by the viewport to set its content.
func RenderDashboardContent(agg analytics.ExtendedAggregates, width int) string {
	if agg.TotalSessions == 0 {
		return renderEmptyDashboardContent(width)
	}

	var sections []string

	// Title
	sections = append(sections, styles.Bold.Render("STATISTICS"))
	sections = append(sections, "")

	// Stats row: points • sessions • accuracy
	statsRow := fmt.Sprintf("%d points  •  %d sessions  •  %.0f%% accuracy",
		agg.TotalPoints, agg.TotalSessions, agg.OverallAccuracy)
	sections = append(sections, statsRow)
	sections = append(sections, "")

	// Separator
	sections = append(sections, renderSeparator(52))
	sections = append(sections, "")

	// Operations section (only show played operations)
	opSection := renderOperationsSection(agg)
	if opSection != "" {
		sections = append(sections, opSection)
		sections = append(sections, "")
		sections = append(sections, renderSeparator(52))
		sections = append(sections, "")
	}

	// Records section
	sections = append(sections, renderRecordsSection(agg))
	sections = append(sections, "")
	sections = append(sections, renderSeparator(52))
	sections = append(sections, "")

	// Footer: thinking time (prominent) and last played (dimmed)
	if agg.TotalResponseTimeMs > 0 {
		thinkingTime := FormatThinkingTime(agg.TotalResponseTimeMs)
		sections = append(sections, thinkingTime+" thinking")
		sections = append(sections, "")
	}
	if !agg.LastPlayedAt.IsZero() {
		sections = append(sections, styles.Dim.Render("Last played: "+FormatRelativeTime(agg.LastPlayedAt)))
	}

	return lipgloss.JoinVertical(lipgloss.Center, sections...)
}

// renderSeparator renders a horizontal separator line.
func renderSeparator(width int) string {
	return styles.Dim.Render(strings.Repeat("━", width))
}

// renderOperationsSection renders the operations with progress bars.
// Only shows operations that have been played.
func renderOperationsSection(agg analytics.ExtendedAggregates) string {
	// Get played operations sorted by name
	var ops []string
	for op, stats := range agg.ByOperation {
		if stats.Total > 0 {
			ops = append(ops, op)
		}
	}

	if len(ops) == 0 {
		return ""
	}

	// Sort operations in a logical order
	opOrder := map[string]int{
		"Addition":       1,
		"Subtraction":    2,
		"Multiplication": 3,
		"Division":       4,
		"Square":         5,
		"Cube":           6,
		"Square Root":    7,
		"Cube Root":      8,
		"Modulo":         9,
		"Power":          10,
		"Percentage":     11,
		"Factorial":      12,
	}
	sort.Slice(ops, func(i, j int) bool {
		return opOrder[ops[i]] < opOrder[ops[j]]
	})

	var lines []string
	lines = append(lines, styles.Bold.Render("OPERATIONS"))
	lines = append(lines, "")

	// Find max operation name length for alignment
	maxNameLen := 0
	for _, op := range ops {
		if len(op) > maxNameLen {
			maxNameLen = len(op)
		}
	}

	barWidth := 10

	for _, op := range ops {
		stats := agg.ByOperation[op]
		symbol := OperationSymbol(op)

		// Format: "symbol Name        XX%  ████░░░░░░"
		accStr := fmt.Sprintf("%3.0f%%", stats.Accuracy)

		// Color the accuracy
		if stats.Accuracy >= 80 {
			accStr = styles.Correct.Render(accStr)
		} else if stats.Accuracy < 60 {
			accStr = styles.Incorrect.Render(accStr)
		}

		// Progress bar
		bar := components.RenderProgressBarColored(stats.Accuracy, barWidth)

		// Pad operation name
		paddedName := op + strings.Repeat(" ", maxNameLen-len(op))

		line := fmt.Sprintf("%s %s   %s  %s", symbol, paddedName, accStr, bar)
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

// renderRecordsSection renders the records in a 2x2 grid.
func renderRecordsSection(agg analytics.ExtendedAggregates) string {
	var lines []string
	lines = append(lines, styles.Bold.Render("RECORDS"))
	lines = append(lines, "")

	// Build 2x2 grid: Best Streak | High Score
	//                 Best Acc    | Fastest
	colWidth := 24

	// Row 1: Best Streak and High Score
	streak := ""
	if agg.PersonalBests.BestStreak > 0 {
		streak = fmt.Sprintf("Best Streak   %-6d", agg.PersonalBests.BestStreak)
	}
	score := ""
	if agg.PersonalBests.BestScore > 0 {
		score = fmt.Sprintf("High Score   %-6d", agg.PersonalBests.BestScore)
	}
	if streak != "" || score != "" {
		row1 := fmt.Sprintf("%-*s%s", colWidth, streak, score)
		lines = append(lines, row1)
	}

	// Row 2: Best Accuracy and Fastest
	acc := ""
	if agg.PersonalBests.BestAccuracy > 0 {
		acc = fmt.Sprintf("Best Acc      %-6s", fmt.Sprintf("%.0f%%", agg.PersonalBests.BestAccuracy))
	}
	fastest := ""
	if agg.FastestResponseMs > 0 {
		fastest = fmt.Sprintf("Fastest      %-6s", FormatResponseTime(agg.FastestResponseMs))
	}
	if acc != "" || fastest != "" {
		row2 := fmt.Sprintf("%-*s%s", colWidth, acc, fastest)
		lines = append(lines, row2)
	}

	// If no records yet
	if len(lines) == 2 {
		lines = append(lines, styles.Dim.Render("Play more to unlock!"))
	}

	return strings.Join(lines, "\n")
}

// renderEmptyDashboardContent renders the empty state content for dashboard.
func renderEmptyDashboardContent(width int) string {
	content := lipgloss.JoinVertical(lipgloss.Center,
		styles.Bold.Render("STATISTICS"),
		"",
		"",
		"Play your first game!",
		"",
		styles.Dim.Render("Complete a session to see your stats."),
	)

	// Center the content
	if width > 0 {
		return lipgloss.PlaceHorizontal(width, lipgloss.Center, content)
	}
	return content
}


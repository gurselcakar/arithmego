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

	// Hero metric - total points in a centered box
	sections = append(sections, renderHeroPoints(agg.TotalPoints))
	sections = append(sections, "")

	// Quick stats row
	sessions := fmt.Sprintf("%d sessions", agg.TotalSessions)
	accuracy := fmt.Sprintf("%.0f%% accuracy", agg.OverallAccuracy)
	quickStats := sessions + "   •   " + accuracy
	sections = append(sections, quickStats)
	sections = append(sections, "")

	// Separator
	sections = append(sections, renderSeparator(38))
	sections = append(sections, "")

	// Operations section (only show played operations)
	opSection := renderOperationsSection(agg)
	if opSection != "" {
		sections = append(sections, opSection)
		sections = append(sections, "")
		sections = append(sections, renderSeparator(38))
		sections = append(sections, "")
	}

	// Personal Bests section
	sections = append(sections, renderPersonalBestsSection(agg))
	sections = append(sections, "")
	sections = append(sections, renderSeparator(38))
	sections = append(sections, "")

	// Time thinking
	if agg.TotalResponseTimeMs > 0 {
		thinkingTime := FormatThinkingTime(agg.TotalResponseTimeMs)
		sections = append(sections, fmt.Sprintf("⏱  %s time thinking", thinkingTime))
		sections = append(sections, "")
	}

	// Last played
	if !agg.LastPlayedAt.IsZero() {
		sections = append(sections, styles.Dim.Render("Last played: "+FormatRelativeTime(agg.LastPlayedAt)))
	}

	return lipgloss.JoinVertical(lipgloss.Center, sections...)
}

// renderHeroPoints renders the total points in a centered box.
func renderHeroPoints(points int) string {
	pointsStr := fmt.Sprintf("%d", points)
	label := "total points"

	// Calculate box width based on content
	contentWidth := len(pointsStr)
	if len(label) > contentWidth {
		contentWidth = len(label)
	}
	boxWidth := contentWidth + 6 // padding on each side

	// Build the box
	var lines []string

	// Top border
	lines = append(lines, "╭"+strings.Repeat("─", boxWidth)+"╮")

	// Empty line
	lines = append(lines, "│"+strings.Repeat(" ", boxWidth)+"│")

	// Points value (centered)
	pointsPadding := (boxWidth - len(pointsStr)) / 2
	pointsLine := "│" + strings.Repeat(" ", pointsPadding) + styles.ScoreLarge.Render(pointsStr) + strings.Repeat(" ", boxWidth-pointsPadding-len(pointsStr)) + "│"
	lines = append(lines, pointsLine)

	// Label (centered)
	labelPadding := (boxWidth - len(label)) / 2
	labelLine := "│" + strings.Repeat(" ", labelPadding) + styles.Dim.Render(label) + strings.Repeat(" ", boxWidth-labelPadding-len(label)) + "│"
	lines = append(lines, labelLine)

	// Empty line
	lines = append(lines, "│"+strings.Repeat(" ", boxWidth)+"│")

	// Bottom border
	lines = append(lines, "╰"+strings.Repeat("─", boxWidth)+"╯")

	return strings.Join(lines, "\n")
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
	lines = append(lines, styles.Bold.Render("YOUR OPERATIONS"))
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

// renderPersonalBestsSection renders the personal bests.
func renderPersonalBestsSection(agg analytics.ExtendedAggregates) string {
	var lines []string
	lines = append(lines, styles.Bold.Render("PERSONAL BESTS"))
	lines = append(lines, "")

	labelWidth := 16
	valueWidth := 6
	totalWidth := labelWidth + valueWidth

	// Best Streak
	if agg.PersonalBests.BestStreak > 0 {
		lines = append(lines, formatBestLine("Best Streak", fmt.Sprintf("%d", agg.PersonalBests.BestStreak), labelWidth, totalWidth))
	}

	// High Score
	if agg.PersonalBests.BestScore > 0 {
		lines = append(lines, formatBestLine("High Score", fmt.Sprintf("%d", agg.PersonalBests.BestScore), labelWidth, totalWidth))
	}

	// Best Accuracy (only show if meaningful - requires min 10 questions)
	if agg.PersonalBests.BestAccuracy > 0 {
		lines = append(lines, formatBestLine("Best Accuracy", fmt.Sprintf("%.0f%%", agg.PersonalBests.BestAccuracy), labelWidth, totalWidth))
	}

	// Fastest Answer
	if agg.FastestResponseMs > 0 {
		lines = append(lines, formatBestLine("Fastest Answer", FormatResponseTime(agg.FastestResponseMs), labelWidth, totalWidth))
	}

	// If no personal bests yet
	if len(lines) == 2 {
		lines = append(lines, styles.Dim.Render("Play more to unlock!"))
	}

	return strings.Join(lines, "\n")
}

// formatBestLine formats a personal best line with label and value.
// All lines are padded to totalWidth for consistent alignment when centered.
func formatBestLine(label, value string, labelWidth, totalWidth int) string {
	paddedLabel := label + strings.Repeat(" ", labelWidth-len(label))
	line := paddedLabel + value
	// Pad to total width for consistent centering
	if len(line) < totalWidth {
		line += strings.Repeat(" ", totalWidth-len(line))
	}
	return line
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


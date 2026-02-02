package statistics

import (
	"fmt"
	"strings"

	"github.com/gurselcakar/arithmego/internal/analytics"
	"github.com/gurselcakar/arithmego/internal/storage"
	"github.com/gurselcakar/arithmego/internal/ui/components"
	"github.com/gurselcakar/arithmego/internal/ui/styles"
)

// OperationRow represents a row in the operations list.
type OperationRow struct {
	Symbol    string
	Name      string
	Category  string
	Accuracy  float64
	Correct   int
	Total     int
	AvgTimeMs int64
}

// BuildOperationList builds the list of operations from aggregates.
func BuildOperationList(agg analytics.ExtendedAggregates, stats *storage.Statistics, filter analytics.AggregateFilter) []OperationRow {
	var rows []OperationRow

	// Get operations that have data, filtered by category
	ops := analytics.GetOperationsByCategory(stats, filter.Category)

	for _, op := range ops {
		extStats, ok := agg.ByOperationExtended[op]
		if !ok || extStats.Total == 0 {
			continue
		}

		rows = append(rows, OperationRow{
			Symbol:    OperationSymbol(op),
			Name:      op,
			Category:  analytics.GetOperationCategory(op),
			Accuracy:  extStats.Accuracy,
			Correct:   extStats.Correct,
			Total:     extStats.Total,
			AvgTimeMs: extStats.AvgResponseTimeMs,
		})
	}

	return rows
}

// RenderOperationsContent renders the operations view content for viewport.
func RenderOperationsContent(
	rows []OperationRow,
	selectedIdx int,
	filterPanel FilterPanelModel,
	width int,
) string {
	var b strings.Builder

	// Title
	b.WriteString(styles.Bold.Render("STATISTICS · OPERATIONS"))
	b.WriteString("\n\n")

	// Filter summary
	b.WriteString(styles.Dim.Render("Showing: " + filterPanel.FilterSummary()))
	b.WriteString("\n\n")

	// Separator
	separatorWidth := 56
	if width > 0 && width-10 < separatorWidth {
		separatorWidth = width - 10
	}
	b.WriteString(styles.Dim.Render(strings.Repeat("─", separatorWidth)))
	b.WriteString("\n\n")

	// Empty state
	if len(rows) == 0 {
		b.WriteString("\n")
		b.WriteString(styles.Dim.Render("No data for these filters yet."))
		b.WriteString("\n\n")
		b.WriteString(styles.Dim.Render("Try different filters or play more games!"))
		b.WriteString("\n")
	} else {
		// Group by category
		currentCategory := ""
		for i, row := range rows {
			// Category header
			if row.Category != currentCategory {
				if currentCategory != "" {
					b.WriteString("\n")
				}
				currentCategory = row.Category
				b.WriteString(styles.Bold.Render(strings.ToUpper(currentCategory)))
				b.WriteString("\n")
				b.WriteString(styles.Dim.Render(strings.Repeat("─", len(currentCategory)+2)))
				b.WriteString("\n")
			}

			// Operation row
			b.WriteString(renderOperationRow(row, i == selectedIdx, width))
			b.WriteString("\n")
		}
	}

	mainContent := b.String()

	// Layout with filter panel overlay if open
	if filterPanel.IsOpen() {
		mainContent = filterPanel.OverlayView(mainContent, width)
	}

	return mainContent
}

// renderOperationRow renders a single operation row.
func renderOperationRow(row OperationRow, selected bool, termWidth int) string {
	// Calculate progress bar width
	barWidth := components.ProgressBarWidth(termWidth)

	// Format: "▸ +  Addition        95%  ████████████████████░░░░  (142/150)  1.2s"
	symbol := row.Symbol
	name := fmt.Sprintf("%-14s", row.Name)
	accStr := fmt.Sprintf("%3.0f%%", row.Accuracy)
	bar := components.RenderProgressBarColored(row.Accuracy, barWidth)
	counts := fmt.Sprintf("(%d/%d)", row.Correct, row.Total)
	timeStr := FormatResponseTime(row.AvgTimeMs)

	// Selection indicator
	prefix := "  "
	if selected {
		prefix = styles.Accent.Render("▸ ")
	}

	// Color the accuracy
	if row.Accuracy >= 80 {
		accStr = styles.Correct.Render(accStr)
	} else if row.Accuracy < 60 {
		accStr = styles.Incorrect.Render(accStr)
	}

	line := fmt.Sprintf("%s%s  %s  %s  %s  %-10s  %s",
		prefix, symbol, name, accStr, bar, counts, timeStr)

	if selected {
		return styles.Bold.Render(line)
	}
	return line
}

// GetSelectedOperation returns the operation name at the given index.
func GetSelectedOperation(rows []OperationRow, idx int) string {
	if idx >= 0 && idx < len(rows) {
		return rows[idx].Name
	}
	return ""
}

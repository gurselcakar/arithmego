package statistics

import (
	"fmt"
	"strings"

	"github.com/gurselcakar/arithmego/internal/analytics"
	"github.com/gurselcakar/arithmego/internal/ui/styles"
)

// ChartPoint is the character used for data points in line charts.
const ChartPoint = "∗"

// RenderLineChart renders an ASCII line chart from trend points.
// width and height specify the chart dimensions (not including labels).
func RenderLineChart(points []analytics.TrendPoint, metric analytics.TrendMetric, width, height int) string {
	if len(points) == 0 || width < 10 || height < 3 {
		return ""
	}

	// Get values based on metric
	values := make([]float64, len(points))
	for i, p := range points {
		switch metric {
		case analytics.TrendMetricAccuracy:
			values[i] = p.Accuracy
		case analytics.TrendMetricSessions:
			values[i] = float64(p.Sessions)
		case analytics.TrendMetricScore:
			values[i] = float64(p.TotalScore)
		case analytics.TrendMetricResponseTime:
			values[i] = float64(p.AvgResponseTime) / 1000.0 // Convert to seconds
		}
	}

	// Find min and max for scaling
	minVal, maxVal := values[0], values[0]
	for _, v := range values {
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
	}

	// For accuracy, use fixed 60-100% range if all values are above 60%
	if metric == analytics.TrendMetricAccuracy {
		if minVal >= 60 {
			minVal = 60
			maxVal = 100
		} else {
			minVal = 0
			maxVal = 100
		}
	}

	// Add padding to range
	valRange := maxVal - minVal
	if valRange == 0 {
		valRange = 1
	}

	// Build the chart grid
	grid := make([][]rune, height)
	for i := range grid {
		grid[i] = make([]rune, width)
		for j := range grid[i] {
			grid[i][j] = ' '
		}
	}

	// Plot points
	for i, v := range values {
		// Map value to y position (inverted: 0 is top)
		y := height - 1 - int((v-minVal)/valRange*float64(height-1))
		if y < 0 {
			y = 0
		}
		if y >= height {
			y = height - 1
		}

		// Map index to x position (guard against single point division by zero)
		var x int
		if len(values) == 1 {
			x = width / 2
		} else {
			x = i * (width - 1) / (len(values) - 1)
		}
		if x >= width {
			x = width - 1
		}

		grid[y][x] = []rune(ChartPoint)[0]
	}

	// Build output with Y-axis labels
	var b strings.Builder

	for row := 0; row < height; row++ {
		// Y-axis label
		yVal := maxVal - float64(row)*(maxVal-minVal)/float64(height-1)
		var label string
		switch metric {
		case analytics.TrendMetricAccuracy:
			label = fmt.Sprintf("%3.0f%% │", yVal)
		case analytics.TrendMetricSessions:
			label = fmt.Sprintf("%3.0f │", yVal)
		case analytics.TrendMetricScore:
			label = fmt.Sprintf("%4.0f │", yVal)
		case analytics.TrendMetricResponseTime:
			label = fmt.Sprintf("%.1fs │", yVal)
		default:
			label = fmt.Sprintf("%4.0f │", yVal)
		}

		b.WriteString(styles.Dim.Render(label))
		b.WriteString(string(grid[row]))
		b.WriteString("\n")
	}

	// X-axis
	b.WriteString(styles.Dim.Render("     └" + strings.Repeat("─", width)))
	b.WriteString("\n")

	// X-axis labels (start and end dates)
	if len(points) > 0 {
		startDate := points[0].Date.Format("Jan 2")
		endDate := points[len(points)-1].Date.Format("Jan 2")
		labelLine := "      " + startDate
		padding := width - len(startDate) - len(endDate)
		if padding > 0 {
			labelLine += strings.Repeat(" ", padding) + endDate
		}
		b.WriteString(styles.Dim.Render(labelLine))
	}

	return b.String()
}

// RenderBarChart renders a horizontal bar chart from weekly data.
func RenderBarChart(data []analytics.WeeklyData, maxWidth int) string {
	if len(data) == 0 || maxWidth < 10 {
		return ""
	}

	// Find max sessions for scaling
	maxSessions := 0
	for _, d := range data {
		if d.Sessions > maxSessions {
			maxSessions = d.Sessions
		}
	}
	if maxSessions == 0 {
		maxSessions = 1
	}

	// Calculate label width
	labelWidth := 8 // "Week N: " or date range

	var b strings.Builder
	for _, d := range data {
		// Render label
		label := fmt.Sprintf("%-8s ", d.Label)
		b.WriteString(styles.Dim.Render(label))

		// Render bar
		barWidth := maxWidth - labelWidth - 15 // Leave room for count
		if barWidth < 5 {
			barWidth = 5
		}
		filled := d.Sessions * barWidth / maxSessions
		if filled < 1 && d.Sessions > 0 {
			filled = 1
		}
		bar := strings.Repeat("█", filled)
		b.WriteString(styles.Correct.Render(bar))

		// Render count
		b.WriteString(fmt.Sprintf("  %d sessions\n", d.Sessions))
	}

	return strings.TrimSuffix(b.String(), "\n")
}

// RenderMiniBarChart renders a compact bar chart for operation accuracy.
// Returns a single line with the bar and percentage.
func RenderMiniBarChart(accuracy float64, width int) string {
	filled := int(accuracy / 100 * float64(width))
	if filled > width {
		filled = width
	}
	empty := width - filled

	bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)

	// Color based on accuracy
	if accuracy >= 80 {
		return styles.Correct.Render(bar)
	} else if accuracy < 60 {
		return styles.Incorrect.Render(bar)
	}
	return bar
}

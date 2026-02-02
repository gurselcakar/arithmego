package statistics

import (
	"strings"

	"github.com/gurselcakar/arithmego/internal/analytics"
	"github.com/gurselcakar/arithmego/internal/ui/styles"
)

// TrendsFocus represents which selector is focused in trends view.
type TrendsFocus int

const (
	TrendsFocusMetric TrendsFocus = iota
	TrendsFocusPeriod
)

// TrendsState holds the state for the trends view.
type TrendsState struct {
	Metric    analytics.TrendMetric
	Period    analytics.TimePeriod
	Focus     TrendsFocus
	TrendData analytics.TrendData
	Insights  []analytics.Insight
}

// NewTrendsState creates a new trends state.
func NewTrendsState() TrendsState {
	return TrendsState{
		Metric: analytics.TrendMetricAccuracy,
		Period: analytics.TimePeriod30Days,
	}
}

// NextMetric cycles to the next metric.
func (s *TrendsState) NextMetric() {
	metrics := analytics.AllTrendMetrics()
	for i, m := range metrics {
		if m == s.Metric {
			s.Metric = metrics[(i+1)%len(metrics)]
			return
		}
	}
}

// PrevMetric cycles to the previous metric.
func (s *TrendsState) PrevMetric() {
	metrics := analytics.AllTrendMetrics()
	for i, m := range metrics {
		if m == s.Metric {
			if i == 0 {
				s.Metric = metrics[len(metrics)-1]
			} else {
				s.Metric = metrics[i-1]
			}
			return
		}
	}
}

// NextPeriod cycles to the next period.
func (s *TrendsState) NextPeriod() {
	periods := analytics.AllTimePeriods()
	for i, p := range periods {
		if p == s.Period {
			s.Period = periods[(i+1)%len(periods)]
			return
		}
	}
}

// PrevPeriod cycles to the previous period.
func (s *TrendsState) PrevPeriod() {
	periods := analytics.AllTimePeriods()
	for i, p := range periods {
		if p == s.Period {
			if i == 0 {
				s.Period = periods[len(periods)-1]
			} else {
				s.Period = periods[i-1]
			}
			return
		}
	}
}

// RenderTrendsContent renders the trends view content for viewport.
func RenderTrendsContent(
	state TrendsState,
	agg analytics.ExtendedAggregates,
	width int,
) string {
	var b strings.Builder

	// Title
	b.WriteString(styles.Bold.Render("STATISTICS · TRENDS"))
	b.WriteString("\n\n")

	// Selectors
	metricSelector := formatTrendSelector(state.Metric.String(), state.Focus == TrendsFocusMetric)
	periodSelector := formatTrendSelector(state.Period.String(), state.Focus == TrendsFocusPeriod)

	selectorLine := "Metric: " + metricSelector + "           Period: " + periodSelector
	b.WriteString(selectorLine)
	b.WriteString("\n\n")

	// Separator
	separatorWidth := 56
	if width > 0 && width-10 < separatorWidth {
		separatorWidth = width - 10
	}
	b.WriteString(styles.Dim.Render(strings.Repeat("─", separatorWidth)))
	b.WriteString("\n\n")

	// Check if enough data
	if agg.TotalSessions < 3 || len(state.TrendData.Points) < 2 {
		b.WriteString("\n")
		b.WriteString(styles.Dim.Render("Play a few more games to see trends!"))
		b.WriteString("\n\n")
		b.WriteString(styles.Dim.Render("Trends require at least 3 sessions over"))
		b.WriteString("\n")
		b.WriteString(styles.Dim.Render("multiple days to be meaningful."))
		b.WriteString("\n")

		return b.String()
	}

	// Chart title
	chartTitle := state.Metric.String() + " OVER TIME"
	b.WriteString(styles.Bold.Render(chartTitle))
	b.WriteString("\n")
	b.WriteString(styles.Dim.Render(strings.Repeat("─", len(chartTitle)+2)))
	b.WriteString("\n")

	// Render the chart
	chartWidth := 40
	if width > 0 && width-20 < chartWidth {
		chartWidth = width - 20
	}
	if chartWidth < 20 {
		chartWidth = 20
	}

	chart := RenderLineChart(state.TrendData.Points, state.Metric, chartWidth, 5)
	b.WriteString(chart)
	b.WriteString("\n\n")

	// Sessions per week bar chart
	if len(state.TrendData.SessionsPerWeek) > 0 {
		b.WriteString(styles.Bold.Render("SESSIONS PER WEEK"))
		b.WriteString("\n")
		b.WriteString(styles.Dim.Render("─────────────────"))
		b.WriteString("\n")

		barChart := RenderBarChart(state.TrendData.SessionsPerWeek, chartWidth+10)
		b.WriteString(barChart)
		b.WriteString("\n\n")
	}

	// Insights
	if len(state.Insights) > 0 {
		b.WriteString(styles.Bold.Render("INSIGHTS"))
		b.WriteString("\n")
		b.WriteString(styles.Dim.Render("────────"))
		b.WriteString("\n")

		for _, insight := range state.Insights {
			b.WriteString(insight.Icon + " " + insight.Message)
			b.WriteString("\n")
		}
	}

	return b.String()
}

// formatTrendSelector formats a selector value with arrows.
func formatTrendSelector(value string, focused bool) string {
	if focused {
		return styles.Accent.Render("◀") + " " + styles.Bold.Render(value) + " " + styles.Accent.Render("▶")
	}
	return styles.Dim.Render("◀") + " " + value + " " + styles.Dim.Render("▶")
}

// Package analytics computes aggregate statistics, filters, and trends
// from game session data stored by the storage package.
//
// The package provides three main capabilities:
//
//   - Aggregation: Compute summary statistics across sessions
//   - Filtering: Query sessions by time, difficulty, mode, or operation
//   - Trends: Track performance changes over time with generated insights
//
// # Aggregates
//
// The primary type is [ExtendedAggregates], which contains comprehensive
// statistics computed from session data:
//
//	stats, _ := storage.Load()
//	agg := analytics.ComputeExtendedAggregates(stats)
//	fmt.Printf("Total sessions: %d\n", agg.TotalSessions)
//	fmt.Printf("Overall accuracy: %.1f%%\n", agg.OverallAccuracy)
//
// Aggregates include basic counts, per-operation breakdowns, personal bests,
// and response time statistics.
//
// # Filtering
//
// Use [AggregateFilter] to compute statistics for a subset of data:
//
//	filter := analytics.AggregateFilter{
//	    TimePeriod: analytics.TimePeriod7Days,
//	    Difficulty: "Hard",
//	    Category:   "Basic",
//	}
//	agg := analytics.ComputeFilteredAggregates(stats, filter)
//
// Filters can be combined. An empty string means "all" for string filters,
// and [TimePeriodAllTime] means no time restriction.
//
// # Time Periods
//
// The [TimePeriod] type provides predefined time ranges:
//
//   - [TimePeriodAllTime]: No time restriction
//   - [TimePeriod7Days]: Last 7 days
//   - [TimePeriod14Days]: Last 14 days
//   - [TimePeriod30Days]: Last 30 days
//   - [TimePeriod90Days]: Last 90 days
//
// # Operation Categories
//
// Operations are grouped into categories for filtering:
//
//   - Basic: Addition, Subtraction, Multiplication, Division
//   - Power: Square, Cube, Square Root, Cube Root
//   - Advanced: Modulo, Power, Percentage, Factorial
//
// Use [GetOperationCategory] to look up an operation's category.
//
// # Trends
//
// The [ComputeTrendData] function generates time-series data for charting:
//
//	data := analytics.ComputeTrendData(stats, analytics.TimePeriod30Days)
//	for _, point := range data.Points {
//	    fmt.Printf("%s: %.1f%% accuracy\n", point.Date.Format("Jan 2"), point.Accuracy)
//	}
//
// # Insights
//
// The [GenerateInsights] function produces human-readable observations:
//
//	insights := analytics.GenerateInsights(stats, analytics.TimePeriod7Days)
//	for _, insight := range insights {
//	    fmt.Printf("%s %s\n", insight.Icon, insight.Message)
//	}
//
// Insights include accuracy trends, personal bests, strongest/weakest
// operations, and play streaks.
package analytics

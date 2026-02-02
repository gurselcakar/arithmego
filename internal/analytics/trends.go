package analytics

import (
	"fmt"
	"sort"
	"time"

	"github.com/gurselcakar/arithmego/internal/storage"
)

// TrendMetric represents what metric to track over time.
type TrendMetric int

const (
	TrendMetricAccuracy TrendMetric = iota
	TrendMetricSessions
	TrendMetricScore
	TrendMetricResponseTime
)

// String returns the display name for this metric.
func (m TrendMetric) String() string {
	switch m {
	case TrendMetricAccuracy:
		return "Accuracy"
	case TrendMetricSessions:
		return "Sessions"
	case TrendMetricScore:
		return "Score"
	case TrendMetricResponseTime:
		return "Response Time"
	default:
		return "Unknown"
	}
}

// AllTrendMetrics returns all available trend metrics.
func AllTrendMetrics() []TrendMetric {
	return []TrendMetric{
		TrendMetricAccuracy,
		TrendMetricSessions,
		TrendMetricScore,
		TrendMetricResponseTime,
	}
}

// TrendData contains computed trend information for charting.
type TrendData struct {
	// Data points for charting (sorted by date)
	Points []TrendPoint

	// Computed insights
	AccuracyChange float64 // Percentage change over period (positive = improvement)

	// Sessions per week for bar chart
	SessionsPerWeek []WeeklyData
}

// TrendPoint represents a single data point for trend charting.
type TrendPoint struct {
	Date            time.Time
	Accuracy        float64
	Sessions        int
	TotalScore      int
	AvgResponseTime int64 // milliseconds
}

// WeeklyData holds aggregated data for a week.
type WeeklyData struct {
	WeekStart time.Time
	Sessions  int
	Label     string // e.g., "Week 1", "Jan 1-7"
}

// Insight represents a generated insight about user progress.
type Insight struct {
	Icon    string // "•", "↑", "↓", "★"
	Message string
}

// ComputeTrendData computes trend data points for charting.
// All metrics are computed regardless of which metric is specified,
// as the TrendData structure contains all values for flexible display.
func ComputeTrendData(stats *storage.Statistics, period TimePeriod) TrendData {
	data := TrendData{}

	if len(stats.Sessions) == 0 {
		return data
	}

	// Get cutoff time
	cutoff := period.Cutoff()

	// Filter sessions by period
	var sessions []storage.SessionRecord
	for _, s := range stats.Sessions {
		if cutoff.IsZero() || s.Timestamp.After(cutoff) {
			sessions = append(sessions, s)
		}
	}

	if len(sessions) == 0 {
		return data
	}

	// Sort by timestamp
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].Timestamp.Before(sessions[j].Timestamp)
	})

	// Group by day
	dayGroups := make(map[string][]storage.SessionRecord)
	for _, s := range sessions {
		dayKey := s.Timestamp.Format("2006-01-02")
		dayGroups[dayKey] = append(dayGroups[dayKey], s)
	}

	// Get sorted days
	days := make([]string, 0, len(dayGroups))
	for day := range dayGroups {
		days = append(days, day)
	}
	sort.Strings(days)

	// Compute data points
	for _, day := range days {
		daySessions := dayGroups[day]
		dayTime, _ := time.Parse("2006-01-02", day)

		var totalCorrect, totalQuestions, totalScore int
		var totalResponseTime int64
		var questionsWithTime int

		for _, s := range daySessions {
			totalCorrect += s.QuestionsCorrect
			totalQuestions += s.QuestionsAttempted
			totalScore += s.Score

			if s.AvgResponseTimeMs > 0 {
				totalResponseTime += s.AvgResponseTimeMs * int64(s.QuestionsAttempted)
				questionsWithTime += s.QuestionsAttempted
			}
		}

		point := TrendPoint{
			Date:     dayTime,
			Sessions: len(daySessions),
		}

		if totalQuestions > 0 {
			point.Accuracy = float64(totalCorrect) / float64(totalQuestions) * 100
			point.TotalScore = totalScore
		}

		if questionsWithTime > 0 {
			point.AvgResponseTime = totalResponseTime / int64(questionsWithTime)
		}

		data.Points = append(data.Points, point)
	}

	// Compute accuracy change (first half vs second half of period)
	if len(data.Points) >= 2 {
		mid := len(data.Points) / 2
		var firstHalfAcc, secondHalfAcc float64
		var firstCount, secondCount int

		for i, p := range data.Points {
			if p.Accuracy > 0 {
				if i < mid {
					firstHalfAcc += p.Accuracy
					firstCount++
				} else {
					secondHalfAcc += p.Accuracy
					secondCount++
				}
			}
		}

		if firstCount > 0 && secondCount > 0 {
			avgFirst := firstHalfAcc / float64(firstCount)
			avgSecond := secondHalfAcc / float64(secondCount)
			data.AccuracyChange = avgSecond - avgFirst
		}
	}

	// Compute weekly data
	data.SessionsPerWeek = computeWeeklyData(sessions)

	return data
}

// computeWeeklyData groups sessions into weeks and computes per-week stats.
func computeWeeklyData(sessions []storage.SessionRecord) []WeeklyData {
	if len(sessions) == 0 {
		return nil
	}

	// Find the week boundaries
	weekGroups := make(map[string][]storage.SessionRecord)
	for _, s := range sessions {
		// Get the Monday of the week containing this session
		weekStart := s.Timestamp.Truncate(24 * time.Hour)
		for weekStart.Weekday() != time.Monday {
			weekStart = weekStart.AddDate(0, 0, -1)
		}
		weekKey := weekStart.Format("2006-01-02")
		weekGroups[weekKey] = append(weekGroups[weekKey], s)
	}

	// Get sorted week starts
	weeks := make([]string, 0, len(weekGroups))
	for week := range weekGroups {
		weeks = append(weeks, week)
	}
	sort.Strings(weeks)

	var result []WeeklyData
	for i, week := range weeks {
		weekStart, _ := time.Parse("2006-01-02", week)
		weekEnd := weekStart.AddDate(0, 0, 6)

		label := fmt.Sprintf("Week %d", i+1)
		if len(weeks) <= 4 {
			// Use date range for fewer weeks
			label = fmt.Sprintf("%s-%s", weekStart.Format("Jan 2"), weekEnd.Format("2"))
		}

		result = append(result, WeeklyData{
			WeekStart: weekStart,
			Sessions:  len(weekGroups[week]),
			Label:     label,
		})
	}

	return result
}

// GenerateInsights generates dynamic insights based on statistics.
func GenerateInsights(stats *storage.Statistics, period TimePeriod) []Insight {
	var insights []Insight

	if len(stats.Sessions) == 0 {
		return insights
	}

	// Get cutoff time
	cutoff := period.Cutoff()

	// Filter sessions by period
	var sessions []storage.SessionRecord
	for _, s := range stats.Sessions {
		if cutoff.IsZero() || s.Timestamp.After(cutoff) {
			sessions = append(sessions, s)
		}
	}

	if len(sessions) == 0 {
		return insights
	}

	// Compute aggregates for the period
	agg := ComputeFilteredAggregates(stats, AggregateFilter{TimePeriod: period})

	// Insight 1: Accuracy trend
	trendData := ComputeTrendData(stats, period)
	if len(trendData.Points) >= 2 {
		change := trendData.AccuracyChange
		if change > 2 {
			insights = append(insights, Insight{
				Icon:    "↑",
				Message: fmt.Sprintf("Your accuracy improved %.0f%% over this period", change),
			})
		} else if change < -2 {
			insights = append(insights, Insight{
				Icon:    "↓",
				Message: fmt.Sprintf("Your accuracy decreased %.0f%% over this period", -change),
			})
		}
	}

	// Insight 2: Best streak
	if agg.PersonalBests.BestStreak > 0 {
		insights = append(insights, Insight{
			Icon:    "★",
			Message: fmt.Sprintf("Best streak: %d correct in a row", agg.PersonalBests.BestStreak),
		})
	}

	// Insight 3: Most played mode
	if len(agg.ByMode) > 0 {
		var topMode string
		var topCount int
		for mode, count := range agg.ByMode {
			if count > topCount {
				topMode = mode
				topCount = count
			}
		}
		if topMode != "" && topCount > 1 {
			insights = append(insights, Insight{
				Icon:    "•",
				Message: fmt.Sprintf("Most played: %s (%d sessions)", topMode, topCount),
			})
		}
	}

	// Insight 4: Strongest operation
	if len(agg.ByOperation) > 0 {
		var bestOp string
		var bestAccuracy float64
		for op, opStats := range agg.ByOperation {
			if opStats.Total >= 10 && opStats.Accuracy > bestAccuracy {
				bestOp = op
				bestAccuracy = opStats.Accuracy
			}
		}
		if bestOp != "" && bestAccuracy >= 80 {
			insights = append(insights, Insight{
				Icon:    "•",
				Message: fmt.Sprintf("Strongest: %s (%.0f%% accuracy)", bestOp, bestAccuracy),
			})
		}
	}

	// Insight 5: Weakest operation (needs improvement)
	if len(agg.ByOperation) > 0 {
		var worstOp string
		var worstAccuracy float64 = 100
		for op, opStats := range agg.ByOperation {
			if opStats.Total >= 10 && opStats.Accuracy < worstAccuracy {
				worstOp = op
				worstAccuracy = opStats.Accuracy
			}
		}
		if worstOp != "" && worstAccuracy < 70 {
			insights = append(insights, Insight{
				Icon:    "•",
				Message: fmt.Sprintf("Needs practice: %s (%.0f%% accuracy)", worstOp, worstAccuracy),
			})
		}
	}

	// Insight 6: Play streak (consecutive days)
	playStreak := computePlayStreak(sessions)
	if playStreak >= 3 {
		insights = append(insights, Insight{
			Icon:    "★",
			Message: fmt.Sprintf("You're on a %d-day play streak!", playStreak),
		})
	}

	// Limit to 4 insights
	if len(insights) > 4 {
		insights = insights[:4]
	}

	return insights
}

// computePlayStreak returns the number of consecutive days with at least one session,
// counting back from today.
func computePlayStreak(sessions []storage.SessionRecord) int {
	if len(sessions) == 0 {
		return 0
	}

	// Get unique days with sessions
	daySet := make(map[string]bool)
	for _, s := range sessions {
		dayKey := s.Timestamp.Format("2006-01-02")
		daySet[dayKey] = true
	}

	// Count consecutive days from today
	// Maximum 365 days to prevent infinite loops
	const maxDays = 365
	streak := 0
	today := time.Now().Truncate(24 * time.Hour)

	for i := 0; i < maxDays; i++ {
		dayKey := today.Format("2006-01-02")
		if daySet[dayKey] {
			streak++
			today = today.AddDate(0, 0, -1)
		} else {
			break
		}
	}

	return streak
}

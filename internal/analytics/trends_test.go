package analytics

import (
	"testing"
	"time"

	"github.com/gurselcakar/arithmego/internal/storage"
)

func TestTrendMetric_String(t *testing.T) {
	tests := []struct {
		metric TrendMetric
		want   string
	}{
		{TrendMetricAccuracy, "Accuracy"},
		{TrendMetricSessions, "Sessions"},
		{TrendMetricScore, "Score"},
		{TrendMetricResponseTime, "Response Time"},
		{TrendMetric(99), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.metric.String(); got != tt.want {
				t.Errorf("String() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestAllTrendMetrics(t *testing.T) {
	metrics := AllTrendMetrics()

	if len(metrics) != 4 {
		t.Errorf("len(AllTrendMetrics()) = %d, want 4", len(metrics))
	}

	expected := []TrendMetric{
		TrendMetricAccuracy,
		TrendMetricSessions,
		TrendMetricScore,
		TrendMetricResponseTime,
	}

	for i, m := range metrics {
		if m != expected[i] {
			t.Errorf("metrics[%d] = %v, want %v", i, m, expected[i])
		}
	}
}

func TestComputeTrendData_Empty(t *testing.T) {
	stats := &storage.Statistics{Sessions: []storage.SessionRecord{}}
	data := ComputeTrendData(stats, TimePeriodAllTime)

	if len(data.Points) != 0 {
		t.Errorf("len(Points) = %d, want 0", len(data.Points))
	}
	if len(data.SessionsPerWeek) != 0 {
		t.Errorf("len(SessionsPerWeek) = %d, want 0", len(data.SessionsPerWeek))
	}
}

func TestComputeTrendData_SingleDay(t *testing.T) {
	// Use a fixed time at noon to avoid midnight boundary issues in CI
	now := time.Date(2025, 6, 15, 12, 0, 0, 0, time.UTC)
	stats := &storage.Statistics{
		Sessions: []storage.SessionRecord{
			{
				Timestamp:          now,
				QuestionsAttempted: 10,
				QuestionsCorrect:   8,
				Score:              500,
				AvgResponseTimeMs:  2000,
			},
			{
				Timestamp:          now.Add(time.Hour), // Same day
				QuestionsAttempted: 10,
				QuestionsCorrect:   6,
				Score:              300,
				AvgResponseTimeMs:  3000,
			},
		},
	}

	data := ComputeTrendData(stats, TimePeriodAllTime)

	if len(data.Points) != 1 {
		t.Fatalf("len(Points) = %d, want 1 (single day)", len(data.Points))
	}

	point := data.Points[0]
	if point.Sessions != 2 {
		t.Errorf("Sessions = %d, want 2", point.Sessions)
	}
	// Combined accuracy: 14/20 = 70%
	if point.Accuracy != 70.0 {
		t.Errorf("Accuracy = %f, want 70.0", point.Accuracy)
	}
	// Total score
	if point.TotalScore != 800 {
		t.Errorf("TotalScore = %d, want 800", point.TotalScore)
	}
}

func TestComputeTrendData_MultipleDays(t *testing.T) {
	now := time.Now()
	stats := &storage.Statistics{
		Sessions: []storage.SessionRecord{
			{
				Timestamp:          now.AddDate(0, 0, -2),
				QuestionsAttempted: 10,
				QuestionsCorrect:   5,
				Score:              200,
			},
			{
				Timestamp:          now.AddDate(0, 0, -1),
				QuestionsAttempted: 10,
				QuestionsCorrect:   7,
				Score:              350,
			},
			{
				Timestamp:          now,
				QuestionsAttempted: 10,
				QuestionsCorrect:   9,
				Score:              500,
			},
		},
	}

	data := ComputeTrendData(stats, TimePeriodAllTime)

	if len(data.Points) != 3 {
		t.Fatalf("len(Points) = %d, want 3", len(data.Points))
	}

	// Points should be sorted by date (oldest first)
	if data.Points[0].Accuracy != 50.0 {
		t.Errorf("First point accuracy = %f, want 50.0", data.Points[0].Accuracy)
	}
	if data.Points[2].Accuracy != 90.0 {
		t.Errorf("Last point accuracy = %f, want 90.0", data.Points[2].Accuracy)
	}
}

func TestComputeTrendData_AccuracyChange(t *testing.T) {
	now := time.Now()
	stats := &storage.Statistics{
		Sessions: []storage.SessionRecord{
			// First half: low accuracy
			{Timestamp: now.AddDate(0, 0, -4), QuestionsAttempted: 10, QuestionsCorrect: 5},
			{Timestamp: now.AddDate(0, 0, -3), QuestionsAttempted: 10, QuestionsCorrect: 6},
			// Second half: high accuracy
			{Timestamp: now.AddDate(0, 0, -1), QuestionsAttempted: 10, QuestionsCorrect: 9},
			{Timestamp: now, QuestionsAttempted: 10, QuestionsCorrect: 10},
		},
	}

	data := ComputeTrendData(stats, TimePeriodAllTime)

	// First half avg: (50+60)/2 = 55%
	// Second half avg: (90+100)/2 = 95%
	// Change: 95 - 55 = 40
	expectedChange := 40.0
	if data.AccuracyChange < expectedChange-1 || data.AccuracyChange > expectedChange+1 {
		t.Errorf("AccuracyChange = %f, want ~%f", data.AccuracyChange, expectedChange)
	}
}

func TestComputeTrendData_FiltersByPeriod(t *testing.T) {
	now := time.Now()
	stats := &storage.Statistics{
		Sessions: []storage.SessionRecord{
			{Timestamp: now.AddDate(0, 0, -3), QuestionsAttempted: 10, QuestionsCorrect: 5},
			{Timestamp: now.AddDate(0, 0, -10), QuestionsAttempted: 10, QuestionsCorrect: 8}, // Outside 7 days
		},
	}

	data := ComputeTrendData(stats, TimePeriod7Days)

	if len(data.Points) != 1 {
		t.Errorf("len(Points) = %d, want 1 (only sessions within 7 days)", len(data.Points))
	}
}

func TestComputeTrendData_SessionsPerWeek(t *testing.T) {
	now := time.Now()
	// Find the Monday of this week
	monday := now
	for monday.Weekday() != time.Monday {
		monday = monday.AddDate(0, 0, -1)
	}

	stats := &storage.Statistics{
		Sessions: []storage.SessionRecord{
			{Timestamp: monday},
			{Timestamp: monday.AddDate(0, 0, 1)},        // Tuesday same week
			{Timestamp: monday.AddDate(0, 0, -7)},       // Previous week
			{Timestamp: monday.AddDate(0, 0, -7)},       // Previous week
			{Timestamp: monday.AddDate(0, 0, -7)},       // Previous week
		},
	}

	data := ComputeTrendData(stats, TimePeriodAllTime)

	if len(data.SessionsPerWeek) != 2 {
		t.Fatalf("len(SessionsPerWeek) = %d, want 2", len(data.SessionsPerWeek))
	}

	// First week (older) should have 3 sessions
	if data.SessionsPerWeek[0].Sessions != 3 {
		t.Errorf("First week sessions = %d, want 3", data.SessionsPerWeek[0].Sessions)
	}
	// Second week should have 2 sessions
	if data.SessionsPerWeek[1].Sessions != 2 {
		t.Errorf("Second week sessions = %d, want 2", data.SessionsPerWeek[1].Sessions)
	}
}

func TestGenerateInsights_Empty(t *testing.T) {
	stats := &storage.Statistics{Sessions: []storage.SessionRecord{}}
	insights := GenerateInsights(stats, TimePeriodAllTime)

	if len(insights) != 0 {
		t.Errorf("len(insights) = %d, want 0 for empty stats", len(insights))
	}
}

func TestGenerateInsights_BestStreak(t *testing.T) {
	stats := &storage.Statistics{
		Sessions: []storage.SessionRecord{
			{
				Timestamp:  time.Now(),
				BestStreak: 15,
				Questions:  []storage.QuestionRecord{},
			},
		},
	}

	insights := GenerateInsights(stats, TimePeriodAllTime)

	// Should have best streak insight
	found := false
	for _, insight := range insights {
		if insight.Icon == "★" && insight.Message == "Best streak: 15 correct in a row" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected best streak insight not found")
	}
}

func TestGenerateInsights_MostPlayedMode(t *testing.T) {
	stats := &storage.Statistics{
		Sessions: []storage.SessionRecord{
			{Timestamp: time.Now(), Mode: "Addition"},
			{Timestamp: time.Now(), Mode: "Addition"},
			{Timestamp: time.Now(), Mode: "Addition"},
			{Timestamp: time.Now(), Mode: "Multiplication"},
		},
	}

	insights := GenerateInsights(stats, TimePeriodAllTime)

	found := false
	for _, insight := range insights {
		if insight.Message == "Most played: Addition (3 sessions)" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected most played mode insight not found")
	}
}

func TestGenerateInsights_StrongestOperation(t *testing.T) {
	stats := &storage.Statistics{
		Sessions: []storage.SessionRecord{
			{
				Timestamp: time.Now(),
				Questions: []storage.QuestionRecord{
					// Addition: 9/10 = 90%
					{Operation: "Addition", Correct: true},
					{Operation: "Addition", Correct: true},
					{Operation: "Addition", Correct: true},
					{Operation: "Addition", Correct: true},
					{Operation: "Addition", Correct: true},
					{Operation: "Addition", Correct: true},
					{Operation: "Addition", Correct: true},
					{Operation: "Addition", Correct: true},
					{Operation: "Addition", Correct: true},
					{Operation: "Addition", Correct: false},
					// Multiplication: 5/10 = 50%
					{Operation: "Multiplication", Correct: true},
					{Operation: "Multiplication", Correct: true},
					{Operation: "Multiplication", Correct: true},
					{Operation: "Multiplication", Correct: true},
					{Operation: "Multiplication", Correct: true},
					{Operation: "Multiplication", Correct: false},
					{Operation: "Multiplication", Correct: false},
					{Operation: "Multiplication", Correct: false},
					{Operation: "Multiplication", Correct: false},
					{Operation: "Multiplication", Correct: false},
				},
			},
		},
	}

	insights := GenerateInsights(stats, TimePeriodAllTime)

	foundStrong := false
	foundWeak := false
	for _, insight := range insights {
		if insight.Message == "Strongest: Addition (90% accuracy)" {
			foundStrong = true
		}
		if insight.Message == "Needs practice: Multiplication (50% accuracy)" {
			foundWeak = true
		}
	}
	if !foundStrong {
		t.Error("Expected strongest operation insight not found")
	}
	if !foundWeak {
		t.Error("Expected needs practice insight not found")
	}
}

func TestGenerateInsights_MaxFour(t *testing.T) {
	// Create stats that would generate many insights
	now := time.Now()
	stats := &storage.Statistics{
		Sessions: []storage.SessionRecord{
			{
				Timestamp:  now,
				Mode:       "Addition",
				BestStreak: 20,
				Questions: []storage.QuestionRecord{
					{Operation: "Addition", Correct: true},
					{Operation: "Addition", Correct: true},
					{Operation: "Addition", Correct: true},
					{Operation: "Addition", Correct: true},
					{Operation: "Addition", Correct: true},
					{Operation: "Addition", Correct: true},
					{Operation: "Addition", Correct: true},
					{Operation: "Addition", Correct: true},
					{Operation: "Addition", Correct: true},
					{Operation: "Addition", Correct: true},
				},
			},
			{Timestamp: now.AddDate(0, 0, -1), Mode: "Addition"},
			{Timestamp: now.AddDate(0, 0, -2), Mode: "Addition"},
		},
	}

	insights := GenerateInsights(stats, TimePeriodAllTime)

	if len(insights) > 4 {
		t.Errorf("len(insights) = %d, want <= 4", len(insights))
	}
}

func TestGenerateInsights_AccuracyImprovement(t *testing.T) {
	now := time.Now()
	stats := &storage.Statistics{
		Sessions: []storage.SessionRecord{
			// Earlier: low accuracy
			{Timestamp: now.AddDate(0, 0, -5), QuestionsAttempted: 10, QuestionsCorrect: 5},
			// Later: high accuracy
			{Timestamp: now, QuestionsAttempted: 10, QuestionsCorrect: 10},
		},
	}

	insights := GenerateInsights(stats, TimePeriodAllTime)

	found := false
	for _, insight := range insights {
		if insight.Icon == "↑" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected accuracy improvement insight not found")
	}
}

func TestGenerateInsights_AccuracyDecline(t *testing.T) {
	now := time.Now()
	stats := &storage.Statistics{
		Sessions: []storage.SessionRecord{
			// Earlier: high accuracy
			{Timestamp: now.AddDate(0, 0, -5), QuestionsAttempted: 10, QuestionsCorrect: 10},
			// Later: low accuracy
			{Timestamp: now, QuestionsAttempted: 10, QuestionsCorrect: 5},
		},
	}

	insights := GenerateInsights(stats, TimePeriodAllTime)

	found := false
	for _, insight := range insights {
		if insight.Icon == "↓" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected accuracy decline insight not found")
	}
}

func TestComputePlayStreak_Empty(t *testing.T) {
	sessions := []storage.SessionRecord{}
	streak := computePlayStreak(sessions)

	if streak != 0 {
		t.Errorf("streak = %d, want 0 for empty sessions", streak)
	}
}

func TestComputePlayStreak_ConsecutiveDays(t *testing.T) {
	now := time.Now()
	sessions := []storage.SessionRecord{
		{Timestamp: now},
		{Timestamp: now.AddDate(0, 0, -1)},
		{Timestamp: now.AddDate(0, 0, -2)},
		{Timestamp: now.AddDate(0, 0, -3)},
	}

	streak := computePlayStreak(sessions)

	if streak != 4 {
		t.Errorf("streak = %d, want 4", streak)
	}
}

func TestComputePlayStreak_GapBreaksStreak(t *testing.T) {
	now := time.Now()
	sessions := []storage.SessionRecord{
		{Timestamp: now},
		{Timestamp: now.AddDate(0, 0, -1)},
		// Gap on day -2
		{Timestamp: now.AddDate(0, 0, -3)},
		{Timestamp: now.AddDate(0, 0, -4)},
	}

	streak := computePlayStreak(sessions)

	if streak != 2 {
		t.Errorf("streak = %d, want 2 (gap breaks streak)", streak)
	}
}

func TestComputePlayStreak_NoSessionToday(t *testing.T) {
	now := time.Now()
	sessions := []storage.SessionRecord{
		// No session today
		{Timestamp: now.AddDate(0, 0, -1)},
		{Timestamp: now.AddDate(0, 0, -2)},
	}

	streak := computePlayStreak(sessions)

	if streak != 0 {
		t.Errorf("streak = %d, want 0 (no session today)", streak)
	}
}

func TestComputePlayStreak_MultipleSameDay(t *testing.T) {
	now := time.Now()
	sessions := []storage.SessionRecord{
		{Timestamp: now},
		{Timestamp: now.Add(-1 * time.Minute)}, // Same day
		{Timestamp: now.AddDate(0, 0, -1)},
		{Timestamp: now.AddDate(0, 0, -1).Add(-1 * time.Minute)}, // Same day
	}

	streak := computePlayStreak(sessions)

	if streak != 2 {
		t.Errorf("streak = %d, want 2 (multiple sessions per day count as 1)", streak)
	}
}

func TestComputeWeeklyData_Empty(t *testing.T) {
	sessions := []storage.SessionRecord{}
	weekly := computeWeeklyData(sessions)

	if weekly != nil {
		t.Errorf("weekly = %v, want nil for empty sessions", weekly)
	}
}

func TestComputeWeeklyData_SingleWeek(t *testing.T) {
	// Find a Monday
	now := time.Now()
	monday := now
	for monday.Weekday() != time.Monday {
		monday = monday.AddDate(0, 0, -1)
	}

	sessions := []storage.SessionRecord{
		{Timestamp: monday},
		{Timestamp: monday.AddDate(0, 0, 1)},
		{Timestamp: monday.AddDate(0, 0, 2)},
	}

	weekly := computeWeeklyData(sessions)

	if len(weekly) != 1 {
		t.Fatalf("len(weekly) = %d, want 1", len(weekly))
	}
	if weekly[0].Sessions != 3 {
		t.Errorf("weekly[0].Sessions = %d, want 3", weekly[0].Sessions)
	}
}

func TestComputeWeeklyData_MultipleWeeks(t *testing.T) {
	// Find a Monday
	now := time.Now()
	monday := now
	for monday.Weekday() != time.Monday {
		monday = monday.AddDate(0, 0, -1)
	}

	sessions := []storage.SessionRecord{
		{Timestamp: monday},
		{Timestamp: monday.AddDate(0, 0, -7)},  // Previous week
		{Timestamp: monday.AddDate(0, 0, -14)}, // Two weeks ago
	}

	weekly := computeWeeklyData(sessions)

	if len(weekly) != 3 {
		t.Fatalf("len(weekly) = %d, want 3", len(weekly))
	}

	// Each week should have 1 session
	for i, w := range weekly {
		if w.Sessions != 1 {
			t.Errorf("weekly[%d].Sessions = %d, want 1", i, w.Sessions)
		}
	}
}

func TestComputeWeeklyData_Labels(t *testing.T) {
	now := time.Now()
	monday := now
	for monday.Weekday() != time.Monday {
		monday = monday.AddDate(0, 0, -1)
	}

	// More than 4 weeks should use "Week N" format
	sessions := []storage.SessionRecord{
		{Timestamp: monday},
		{Timestamp: monday.AddDate(0, 0, -7)},
		{Timestamp: monday.AddDate(0, 0, -14)},
		{Timestamp: monday.AddDate(0, 0, -21)},
		{Timestamp: monday.AddDate(0, 0, -28)},
	}

	weekly := computeWeeklyData(sessions)

	if len(weekly) != 5 {
		t.Fatalf("len(weekly) = %d, want 5", len(weekly))
	}

	// Should use "Week N" format for >4 weeks
	if weekly[0].Label != "Week 1" {
		t.Errorf("weekly[0].Label = %s, want 'Week 1'", weekly[0].Label)
	}
	if weekly[4].Label != "Week 5" {
		t.Errorf("weekly[4].Label = %s, want 'Week 5'", weekly[4].Label)
	}
}

package analytics

import (
	"testing"
	"time"

	"github.com/gurselcakar/arithmego/internal/storage"
)

func TestComputeExtendedAggregates_Empty(t *testing.T) {
	stats := &storage.Statistics{Sessions: []storage.SessionRecord{}}
	agg := ComputeExtendedAggregates(stats)

	if agg.TotalSessions != 0 {
		t.Errorf("TotalSessions = %d, want 0", agg.TotalSessions)
	}
	if agg.TotalQuestions != 0 {
		t.Errorf("TotalQuestions = %d, want 0", agg.TotalQuestions)
	}
	if agg.TotalPoints != 0 {
		t.Errorf("TotalPoints = %d, want 0", agg.TotalPoints)
	}
	if agg.OverallAccuracy != 0 {
		t.Errorf("OverallAccuracy = %f, want 0", agg.OverallAccuracy)
	}
	if !agg.LastPlayedAt.IsZero() {
		t.Errorf("LastPlayedAt should be zero, got %v", agg.LastPlayedAt)
	}
}

func TestComputeExtendedAggregates_SingleSession(t *testing.T) {
	now := time.Now()
	stats := &storage.Statistics{
		Sessions: []storage.SessionRecord{
			{
				ID:                 "test-1",
				Timestamp:          now,
				Mode:               "Addition",
				Difficulty:         "Medium",
				DurationSeconds:    60,
				QuestionsAttempted: 10,
				QuestionsCorrect:   8,
				QuestionsWrong:     2,
				Score:              500,
				BestStreak:         5,
				AvgResponseTimeMs:  2500,
				Questions: []storage.QuestionRecord{
					{Operation: "Addition", Correct: true, ResponseTimeMs: 2000},
					{Operation: "Addition", Correct: true, ResponseTimeMs: 3000},
					{Operation: "Addition", Correct: false, ResponseTimeMs: 2500},
				},
			},
		},
	}

	agg := ComputeExtendedAggregates(stats)

	if agg.TotalSessions != 1 {
		t.Errorf("TotalSessions = %d, want 1", agg.TotalSessions)
	}
	if agg.TotalQuestions != 3 {
		t.Errorf("TotalQuestions = %d, want 3", agg.TotalQuestions)
	}
	if agg.TotalCorrect != 2 {
		t.Errorf("TotalCorrect = %d, want 2", agg.TotalCorrect)
	}
	if agg.TotalPoints != 500 {
		t.Errorf("TotalPoints = %d, want 500", agg.TotalPoints)
	}
	if agg.BestStreakEver != 5 {
		t.Errorf("BestStreakEver = %d, want 5", agg.BestStreakEver)
	}
	if !agg.LastPlayedAt.Equal(now) {
		t.Errorf("LastPlayedAt = %v, want %v", agg.LastPlayedAt, now)
	}

	// Check accuracy (2/3 ≈ 66.67%)
	expectedAccuracy := float64(2) / float64(3) * 100
	if agg.OverallAccuracy < expectedAccuracy-0.1 || agg.OverallAccuracy > expectedAccuracy+0.1 {
		t.Errorf("OverallAccuracy = %f, want ~%f", agg.OverallAccuracy, expectedAccuracy)
	}

	// Check operation stats
	addStats, ok := agg.ByOperation["Addition"]
	if !ok {
		t.Fatal("Addition stats not found")
	}
	if addStats.Total != 3 {
		t.Errorf("Addition.Total = %d, want 3", addStats.Total)
	}
	if addStats.Correct != 2 {
		t.Errorf("Addition.Correct = %d, want 2", addStats.Correct)
	}

	// Check mode stats
	if agg.ByMode["Addition"] != 1 {
		t.Errorf("ByMode[Addition] = %d, want 1", agg.ByMode["Addition"])
	}
}

func TestComputeExtendedAggregates_PersonalBests(t *testing.T) {
	stats := &storage.Statistics{
		Sessions: []storage.SessionRecord{
			{
				ID:                 "test-1",
				Score:              300,
				BestStreak:         5,
				QuestionsAttempted: 10,
				QuestionsCorrect:   7,
				AvgResponseTimeMs:  3000,
				Questions:          []storage.QuestionRecord{},
			},
			{
				ID:                 "test-2",
				Score:              800,
				BestStreak:         12,
				QuestionsAttempted: 15,
				QuestionsCorrect:   14,
				AvgResponseTimeMs:  2000,
				Questions:          []storage.QuestionRecord{},
			},
			{
				ID:                 "test-3",
				Score:              500,
				BestStreak:         8,
				QuestionsAttempted: 20,
				QuestionsCorrect:   19,
				AvgResponseTimeMs:  2500,
				Questions:          []storage.QuestionRecord{},
			},
		},
	}

	agg := ComputeExtendedAggregates(stats)

	if agg.PersonalBests.BestScore != 800 {
		t.Errorf("BestScore = %d, want 800", agg.PersonalBests.BestScore)
	}
	if agg.PersonalBests.BestStreak != 12 {
		t.Errorf("BestStreak = %d, want 12", agg.PersonalBests.BestStreak)
	}
	if agg.PersonalBests.FastestAvgTime != 2000 {
		t.Errorf("FastestAvgTime = %d, want 2000", agg.PersonalBests.FastestAvgTime)
	}
	// Best accuracy: 19/20 = 95%
	if agg.PersonalBests.BestAccuracy != 95.0 {
		t.Errorf("BestAccuracy = %f, want 95.0", agg.PersonalBests.BestAccuracy)
	}
}

func TestComputeExtendedAggregates_BestAccuracyRequiresMinQuestions(t *testing.T) {
	stats := &storage.Statistics{
		Sessions: []storage.SessionRecord{
			{
				ID:                 "test-1",
				QuestionsAttempted: 5, // Less than 10, should not count
				QuestionsCorrect:   5, // 100% but too few questions
			},
			{
				ID:                 "test-2",
				QuestionsAttempted: 10, // Exactly 10, should count
				QuestionsCorrect:   8,  // 80%
			},
		},
	}

	agg := ComputeExtendedAggregates(stats)

	// Should be 80%, not 100% (first session doesn't qualify)
	if agg.PersonalBests.BestAccuracy != 80.0 {
		t.Errorf("BestAccuracy = %f, want 80.0 (100%% session has < 10 questions)", agg.PersonalBests.BestAccuracy)
	}
}

func TestComputeExtendedAggregates_ExtendedOperationStats(t *testing.T) {
	stats := &storage.Statistics{
		Sessions: []storage.SessionRecord{
			{
				ID:         "test-1",
				Difficulty: "Easy",
				Questions: []storage.QuestionRecord{
					{Operation: "Addition", Correct: true, ResponseTimeMs: 1000},
					{Operation: "Addition", Correct: true, ResponseTimeMs: 2000},
					{Operation: "Addition", Correct: false, ResponseTimeMs: 3000},
				},
			},
			{
				ID:         "test-2",
				Difficulty: "Hard",
				Questions: []storage.QuestionRecord{
					{Operation: "Addition", Correct: true, ResponseTimeMs: 1500},
					{Operation: "Addition", Correct: false, ResponseTimeMs: 2500},
				},
			},
		},
	}

	agg := ComputeExtendedAggregates(stats)

	extStats := agg.ByOperationExtended["Addition"]
	if extStats.Total != 5 {
		t.Errorf("Extended Addition.Total = %d, want 5", extStats.Total)
	}
	if extStats.Correct != 3 {
		t.Errorf("Extended Addition.Correct = %d, want 3", extStats.Correct)
	}
	if extStats.FastestTimeMs != 1000 {
		t.Errorf("FastestTimeMs = %d, want 1000", extStats.FastestTimeMs)
	}
	// Avg: (1000+2000+3000+1500+2500)/5 = 2000
	if extStats.AvgResponseTimeMs != 2000 {
		t.Errorf("AvgResponseTimeMs = %d, want 2000", extStats.AvgResponseTimeMs)
	}

	// Check difficulty breakdown
	easyStats := extStats.ByDifficulty["Easy"]
	if easyStats.Total != 3 {
		t.Errorf("Easy.Total = %d, want 3", easyStats.Total)
	}
	if easyStats.Correct != 2 {
		t.Errorf("Easy.Correct = %d, want 2", easyStats.Correct)
	}

	hardStats := extStats.ByDifficulty["Hard"]
	if hardStats.Total != 2 {
		t.Errorf("Hard.Total = %d, want 2", hardStats.Total)
	}
	if hardStats.Correct != 1 {
		t.Errorf("Hard.Correct = %d, want 1", hardStats.Correct)
	}
}

func TestComputeExtendedAggregates_SkippedQuestionsExcluded(t *testing.T) {
	stats := &storage.Statistics{
		Sessions: []storage.SessionRecord{
			{
				Questions: []storage.QuestionRecord{
					{Operation: "Addition", Correct: true, Skipped: false},
					{Operation: "Addition", Correct: false, Skipped: true}, // Should be excluded
					{Operation: "Addition", Correct: true, Skipped: false},
				},
			},
		},
	}

	agg := ComputeExtendedAggregates(stats)

	if agg.TotalQuestions != 2 {
		t.Errorf("TotalQuestions = %d, want 2 (excluding skipped)", agg.TotalQuestions)
	}
	if agg.TotalCorrect != 2 {
		t.Errorf("TotalCorrect = %d, want 2", agg.TotalCorrect)
	}
}

func TestComputeFilteredAggregates_ByTimePeriod(t *testing.T) {
	now := time.Now()
	stats := &storage.Statistics{
		Sessions: []storage.SessionRecord{
			{
				ID:        "recent",
				Timestamp: now.AddDate(0, 0, -3), // 3 days ago
				Score:     100,
				Questions: []storage.QuestionRecord{
					{Operation: "Addition", Correct: true},
				},
			},
			{
				ID:        "old",
				Timestamp: now.AddDate(0, 0, -30), // 30 days ago
				Score:     200,
				Questions: []storage.QuestionRecord{
					{Operation: "Addition", Correct: true},
				},
			},
		},
	}

	// Filter to last 7 days
	agg := ComputeFilteredAggregates(stats, AggregateFilter{TimePeriod: TimePeriod7Days})

	if agg.TotalSessions != 1 {
		t.Errorf("TotalSessions = %d, want 1 (only recent session)", agg.TotalSessions)
	}
	if agg.TotalPoints != 100 {
		t.Errorf("TotalPoints = %d, want 100", agg.TotalPoints)
	}
}

func TestComputeFilteredAggregates_ByDifficulty(t *testing.T) {
	stats := &storage.Statistics{
		Sessions: []storage.SessionRecord{
			{
				ID:         "easy-session",
				Difficulty: "Easy",
				Score:      100,
				Questions: []storage.QuestionRecord{
					{Operation: "Addition", Correct: true},
				},
			},
			{
				ID:         "hard-session",
				Difficulty: "Hard",
				Score:      300,
				Questions: []storage.QuestionRecord{
					{Operation: "Addition", Correct: true},
				},
			},
		},
	}

	agg := ComputeFilteredAggregates(stats, AggregateFilter{Difficulty: "Hard"})

	if agg.TotalSessions != 1 {
		t.Errorf("TotalSessions = %d, want 1", agg.TotalSessions)
	}
	if agg.TotalPoints != 300 {
		t.Errorf("TotalPoints = %d, want 300", agg.TotalPoints)
	}
}

func TestComputeFilteredAggregates_ByCategory(t *testing.T) {
	stats := &storage.Statistics{
		Sessions: []storage.SessionRecord{
			{
				Questions: []storage.QuestionRecord{
					{Operation: "Addition", Correct: true},      // Basic
					{Operation: "Subtraction", Correct: true},   // Basic
					{Operation: "Square", Correct: false},       // Power
					{Operation: "Modulo", Correct: true},        // Advanced
				},
			},
		},
	}

	agg := ComputeFilteredAggregates(stats, AggregateFilter{Category: "Basic"})

	if agg.TotalQuestions != 2 {
		t.Errorf("TotalQuestions = %d, want 2 (only Basic operations)", agg.TotalQuestions)
	}
	if agg.TotalCorrect != 2 {
		t.Errorf("TotalCorrect = %d, want 2", agg.TotalCorrect)
	}
}

func TestGetRecentMistakes(t *testing.T) {
	now := time.Now()
	stats := &storage.Statistics{
		Sessions: []storage.SessionRecord{
			{
				Timestamp: now.AddDate(0, 0, -2), // Older
				Questions: []storage.QuestionRecord{
					{Question: "1+1", Operation: "Addition", Correct: false, UserAnswer: 3, CorrectAnswer: 2},
					{Question: "2+2", Operation: "Addition", Correct: true},
				},
			},
			{
				Timestamp: now, // More recent
				Questions: []storage.QuestionRecord{
					{Question: "3×3", Operation: "Multiplication", Correct: false, UserAnswer: 8, CorrectAnswer: 9},
					{Question: "4×4", Operation: "Multiplication", Correct: false, UserAnswer: 15, CorrectAnswer: 16},
					{Question: "5×5", Operation: "Multiplication", Correct: true},
				},
			},
		},
	}

	mistakes := GetRecentMistakes(stats, "", 3)

	if len(mistakes) != 3 {
		t.Fatalf("len(mistakes) = %d, want 3", len(mistakes))
	}

	// Most recent first (second session, reverse order within session)
	if mistakes[0].Question != "3×3" {
		t.Errorf("First mistake = %s, want 3×3", mistakes[0].Question)
	}
	if mistakes[1].Question != "4×4" {
		t.Errorf("Second mistake = %s, want 4×4", mistakes[1].Question)
	}
	if mistakes[2].Question != "1+1" {
		t.Errorf("Third mistake = %s, want 1+1", mistakes[2].Question)
	}
}

func TestGetRecentMistakes_FilterByOperation(t *testing.T) {
	stats := &storage.Statistics{
		Sessions: []storage.SessionRecord{
			{
				Questions: []storage.QuestionRecord{
					{Question: "1+1", Operation: "Addition", Correct: false},
					{Question: "2×2", Operation: "Multiplication", Correct: false},
					{Question: "3+3", Operation: "Addition", Correct: false},
				},
			},
		},
	}

	mistakes := GetRecentMistakes(stats, "Addition", 10)

	if len(mistakes) != 2 {
		t.Fatalf("len(mistakes) = %d, want 2", len(mistakes))
	}
	for _, m := range mistakes {
		if m.Operation != "Addition" {
			t.Errorf("Got mistake with operation %s, want only Addition", m.Operation)
		}
	}
}

func TestGetRecentMistakes_ExcludesSkipped(t *testing.T) {
	stats := &storage.Statistics{
		Sessions: []storage.SessionRecord{
			{
				Questions: []storage.QuestionRecord{
					{Question: "1+1", Correct: false, Skipped: false},
					{Question: "2+2", Correct: false, Skipped: true}, // Should be excluded
				},
			},
		},
	}

	mistakes := GetRecentMistakes(stats, "", 10)

	if len(mistakes) != 1 {
		t.Errorf("len(mistakes) = %d, want 1 (excluding skipped)", len(mistakes))
	}
}

func TestGetSessionsByFilter(t *testing.T) {
	now := time.Now()
	stats := &storage.Statistics{
		Sessions: []storage.SessionRecord{
			{ID: "1", Timestamp: now.AddDate(0, 0, -1), Difficulty: "Easy"},
			{ID: "2", Timestamp: now.AddDate(0, 0, -5), Difficulty: "Hard"},
			{ID: "3", Timestamp: now, Difficulty: "Easy"},
		},
	}

	// No filter - should return all, sorted by timestamp (most recent first)
	sessions := GetSessionsByFilter(stats, AggregateFilter{})

	if len(sessions) != 3 {
		t.Fatalf("len(sessions) = %d, want 3", len(sessions))
	}
	if sessions[0].ID != "3" {
		t.Errorf("First session ID = %s, want 3 (most recent)", sessions[0].ID)
	}
	if sessions[1].ID != "1" {
		t.Errorf("Second session ID = %s, want 1", sessions[1].ID)
	}

	// Filter by difficulty
	sessions = GetSessionsByFilter(stats, AggregateFilter{Difficulty: "Easy"})

	if len(sessions) != 2 {
		t.Errorf("len(sessions) = %d, want 2 (only Easy)", len(sessions))
	}
}

func TestGetAllModes(t *testing.T) {
	stats := &storage.Statistics{
		Sessions: []storage.SessionRecord{
			{Mode: "Addition"},
			{Mode: "Multiplication"},
			{Mode: "Addition"},
			{Mode: "Mixed Basics"},
		},
	}

	modes := GetAllModes(stats)

	if len(modes) != 3 {
		t.Errorf("len(modes) = %d, want 3", len(modes))
	}

	// Should be sorted
	expected := []string{"Addition", "Mixed Basics", "Multiplication"}
	for i, m := range modes {
		if m != expected[i] {
			t.Errorf("modes[%d] = %s, want %s", i, m, expected[i])
		}
	}
}

func TestGetAllOperations(t *testing.T) {
	stats := &storage.Statistics{
		Sessions: []storage.SessionRecord{
			{
				Questions: []storage.QuestionRecord{
					{Operation: "Addition"},
					{Operation: "Multiplication"},
				},
			},
			{
				Questions: []storage.QuestionRecord{
					{Operation: "Addition"},
					{Operation: "Division"},
				},
			},
		},
	}

	ops := GetAllOperations(stats)

	if len(ops) != 3 {
		t.Errorf("len(ops) = %d, want 3", len(ops))
	}

	// Should be sorted
	expected := []string{"Addition", "Division", "Multiplication"}
	for i, op := range ops {
		if op != expected[i] {
			t.Errorf("ops[%d] = %s, want %s", i, op, expected[i])
		}
	}
}

func TestGetOperationsByCategory(t *testing.T) {
	stats := &storage.Statistics{
		Sessions: []storage.SessionRecord{
			{
				Questions: []storage.QuestionRecord{
					{Operation: "Addition"},
					{Operation: "Square"},
					{Operation: "Modulo"},
					{Operation: "Subtraction"},
				},
			},
		},
	}

	// Get only Basic operations
	ops := GetOperationsByCategory(stats, "Basic")

	if len(ops) != 2 {
		t.Errorf("len(ops) = %d, want 2", len(ops))
	}

	for _, op := range ops {
		cat := GetOperationCategory(op)
		if cat != "Basic" {
			t.Errorf("Operation %s has category %s, want Basic", op, cat)
		}
	}

	// Empty category returns all
	allOps := GetOperationsByCategory(stats, "")
	if len(allOps) != 4 {
		t.Errorf("len(allOps) = %d, want 4", len(allOps))
	}
}

func TestComputeExtendedAggregates_FastestResponse(t *testing.T) {
	stats := &storage.Statistics{
		Sessions: []storage.SessionRecord{
			{
				Questions: []storage.QuestionRecord{
					{Operation: "Addition", ResponseTimeMs: 3000},
					{Operation: "Addition", ResponseTimeMs: 1500},
					{Operation: "Multiplication", ResponseTimeMs: 2000},
				},
			},
		},
	}

	agg := ComputeExtendedAggregates(stats)

	if agg.FastestResponseMs != 1500 {
		t.Errorf("FastestResponseMs = %d, want 1500", agg.FastestResponseMs)
	}
}

func TestComputeExtendedAggregates_TotalResponseTime(t *testing.T) {
	stats := &storage.Statistics{
		Sessions: []storage.SessionRecord{
			{
				Questions: []storage.QuestionRecord{
					{ResponseTimeMs: 1000},
					{ResponseTimeMs: 2000},
					{ResponseTimeMs: 3000},
				},
			},
		},
	}

	agg := ComputeExtendedAggregates(stats)

	if agg.TotalResponseTimeMs != 6000 {
		t.Errorf("TotalResponseTimeMs = %d, want 6000", agg.TotalResponseTimeMs)
	}
	if agg.AvgResponseTimeMs != 2000 {
		t.Errorf("AvgResponseTimeMs = %d, want 2000", agg.AvgResponseTimeMs)
	}
}

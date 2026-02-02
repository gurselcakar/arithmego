package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestComputeAggregates_Empty(t *testing.T) {
	stats := &Statistics{Sessions: []SessionRecord{}}
	agg := ComputeAggregates(stats)

	if agg.TotalSessions != 0 {
		t.Errorf("TotalSessions = %d, want 0", agg.TotalSessions)
	}
	if agg.TotalQuestions != 0 {
		t.Errorf("TotalQuestions = %d, want 0", agg.TotalQuestions)
	}
	if agg.OverallAccuracy != 0 {
		t.Errorf("OverallAccuracy = %f, want 0", agg.OverallAccuracy)
	}
}

func TestComputeAggregates_SingleSession(t *testing.T) {
	stats := &Statistics{
		Sessions: []SessionRecord{
			{
				ID:                 "test-1",
				Mode:               "Addition",
				Difficulty:         "Medium",
				DurationSeconds:    60,
				QuestionsAttempted: 10,
				QuestionsCorrect:   8,
				QuestionsWrong:     2,
				Score:              500,
				BestStreak:         5,
				Questions: []QuestionRecord{
					{Operation: "Addition", Correct: true, ResponseTimeMs: 2000},
					{Operation: "Addition", Correct: true, ResponseTimeMs: 3000},
					{Operation: "Addition", Correct: false, ResponseTimeMs: 4000},
				},
			},
		},
	}

	agg := ComputeAggregates(stats)

	if agg.TotalSessions != 1 {
		t.Errorf("TotalSessions = %d, want 1", agg.TotalSessions)
	}
	if agg.TotalQuestions != 10 {
		t.Errorf("TotalQuestions = %d, want 10", agg.TotalQuestions)
	}
	if agg.TotalCorrect != 8 {
		t.Errorf("TotalCorrect = %d, want 8", agg.TotalCorrect)
	}
	if agg.OverallAccuracy != 80.0 {
		t.Errorf("OverallAccuracy = %f, want 80.0", agg.OverallAccuracy)
	}
	if agg.BestStreakEver != 5 {
		t.Errorf("BestStreakEver = %d, want 5", agg.BestStreakEver)
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
	// 2/3 ≈ 66.67%
	if addStats.Accuracy < 66.6 || addStats.Accuracy > 66.7 {
		t.Errorf("Addition.Accuracy = %f, want ~66.67", addStats.Accuracy)
	}

	// Check mode stats
	modeCount, ok := agg.ByMode["Addition"]
	if !ok {
		t.Fatal("Addition mode not found")
	}
	if modeCount != 1 {
		t.Errorf("Addition count = %d, want 1", modeCount)
	}
}

func TestComputeAggregates_MultipleSessions(t *testing.T) {
	stats := &Statistics{
		Sessions: []SessionRecord{
			{
				ID:                 "test-1",
				Mode:               "Addition",
				QuestionsAttempted: 10,
				QuestionsCorrect:   8,
				BestStreak:         5,
				Questions: []QuestionRecord{
					{Operation: "Addition", Correct: true},
					{Operation: "Addition", Correct: true},
				},
			},
			{
				ID:                 "test-2",
				Mode:               "Mixed Basics",
				QuestionsAttempted: 20,
				QuestionsCorrect:   15,
				BestStreak:         12,
				Questions: []QuestionRecord{
					{Operation: "Addition", Correct: true},
					{Operation: "Multiplication", Correct: false},
					{Operation: "Multiplication", Correct: true},
				},
			},
			{
				ID:                 "test-3",
				Mode:               "Addition",
				QuestionsAttempted: 15,
				QuestionsCorrect:   10,
				BestStreak:         7,
				Questions: []QuestionRecord{
					{Operation: "Addition", Correct: false},
				},
			},
		},
	}

	agg := ComputeAggregates(stats)

	if agg.TotalSessions != 3 {
		t.Errorf("TotalSessions = %d, want 3", agg.TotalSessions)
	}
	if agg.TotalQuestions != 45 {
		t.Errorf("TotalQuestions = %d, want 45", agg.TotalQuestions)
	}
	if agg.TotalCorrect != 33 {
		t.Errorf("TotalCorrect = %d, want 33", agg.TotalCorrect)
	}
	if agg.BestStreakEver != 12 {
		t.Errorf("BestStreakEver = %d, want 12", agg.BestStreakEver)
	}

	// Check mode counts
	if agg.ByMode["Addition"] != 2 {
		t.Errorf("Addition count = %d, want 2", agg.ByMode["Addition"])
	}
	if agg.ByMode["Mixed Basics"] != 1 {
		t.Errorf("Mixed Basics count = %d, want 1", agg.ByMode["Mixed Basics"])
	}

	// Check operation aggregation across sessions
	addStats := agg.ByOperation["Addition"]
	if addStats.Total != 4 {
		t.Errorf("Addition.Total = %d, want 4", addStats.Total)
	}
	if addStats.Correct != 3 {
		t.Errorf("Addition.Correct = %d, want 3", addStats.Correct)
	}

	mulStats := agg.ByOperation["Multiplication"]
	if mulStats.Total != 2 {
		t.Errorf("Multiplication.Total = %d, want 2", mulStats.Total)
	}
	if mulStats.Correct != 1 {
		t.Errorf("Multiplication.Correct = %d, want 1", mulStats.Correct)
	}
}

func TestComputeAggregates_SkippedQuestionsExcluded(t *testing.T) {
	stats := &Statistics{
		Sessions: []SessionRecord{
			{
				ID:                 "test-1",
				QuestionsAttempted: 5,
				QuestionsCorrect:   3,
				QuestionsSkipped:   2,
				Questions: []QuestionRecord{
					{Operation: "Addition", Correct: true, Skipped: false},
					{Operation: "Addition", Correct: false, Skipped: false},
					{Operation: "Addition", Correct: false, Skipped: true}, // Should be excluded
					{Operation: "Addition", Correct: false, Skipped: true}, // Should be excluded
					{Operation: "Addition", Correct: true, Skipped: false},
				},
			},
		},
	}

	agg := ComputeAggregates(stats)

	addStats := agg.ByOperation["Addition"]
	// Only non-skipped questions should be counted
	if addStats.Total != 3 {
		t.Errorf("Addition.Total = %d, want 3 (excluding skipped)", addStats.Total)
	}
	if addStats.Correct != 2 {
		t.Errorf("Addition.Correct = %d, want 2", addStats.Correct)
	}
}

func TestComputeAggregates_AvgResponseTime(t *testing.T) {
	stats := &Statistics{
		Sessions: []SessionRecord{
			{
				Questions: []QuestionRecord{
					{ResponseTimeMs: 2000},
					{ResponseTimeMs: 3000},
					{ResponseTimeMs: 4000},
				},
			},
		},
	}

	agg := ComputeAggregates(stats)

	if agg.AvgResponseTimeMs != 3000 {
		t.Errorf("AvgResponseTimeMs = %d, want 3000", agg.AvgResponseTimeMs)
	}
}

func TestNewSessionRecord(t *testing.T) {
	before := time.Now()
	record, err := NewSessionRecord("Test Mode", "Hard", 90)
	after := time.Now()

	if err != nil {
		t.Fatalf("NewSessionRecord() error = %v", err)
	}
	if record.ID == "" {
		t.Error("ID should not be empty")
	}
	if record.Mode != "Test Mode" {
		t.Errorf("Mode = %s, want Test Mode", record.Mode)
	}
	if record.Difficulty != "Hard" {
		t.Errorf("Difficulty = %s, want Hard", record.Difficulty)
	}
	if record.DurationSeconds != 90 {
		t.Errorf("DurationSeconds = %d, want 90", record.DurationSeconds)
	}
	if record.Timestamp.Before(before) || record.Timestamp.After(after) {
		t.Errorf("Timestamp %v not in expected range", record.Timestamp)
	}
	if record.Questions == nil {
		t.Error("Questions should be initialized")
	}
}

func TestNewSessionRecord_Validation(t *testing.T) {
	tests := []struct {
		name       string
		mode       string
		difficulty string
		duration   int
		wantErr    bool
	}{
		{"valid", "Test Mode", "Easy", 60, false},
		{"empty mode", "", "Easy", 60, true},
		{"empty difficulty", "Test Mode", "", 60, true},
		{"negative duration", "Test Mode", "Easy", -1, true},
		{"zero duration", "Test Mode", "Easy", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewSessionRecord(tt.mode, tt.difficulty, tt.duration)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSessionRecord() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoadSave(t *testing.T) {
	// Use a temporary directory for test isolation
	tempDir := t.TempDir()
	SetConfigDirForTesting(tempDir)
	defer SetConfigDirForTesting("")

	path, err := StatisticsPath()
	if err != nil {
		t.Fatalf("StatisticsPath() error = %v", err)
	}

	// Test loading non-existent file returns empty stats
	stats, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if len(stats.Sessions) != 0 {
		t.Errorf("Expected empty sessions, got %d", len(stats.Sessions))
	}

	// Add a session and save
	record, err := NewSessionRecord("Test Mode", "Easy", 60)
	if err != nil {
		t.Fatalf("NewSessionRecord() error = %v", err)
	}
	record.QuestionsCorrect = 10
	record.Score = 500

	err = AddSession(record)
	if err != nil {
		t.Fatalf("AddSession() error = %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("Statistics file should exist after save")
	}

	// Load and verify
	stats, err = Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if len(stats.Sessions) != 1 {
		t.Fatalf("Expected 1 session, got %d", len(stats.Sessions))
	}
	if stats.Sessions[0].Mode != "Test Mode" {
		t.Errorf("Mode = %s, want Test Mode", stats.Sessions[0].Mode)
	}
	if stats.Sessions[0].Score != 500 {
		t.Errorf("Score = %d, want 500", stats.Sessions[0].Score)
	}
}

func TestLoad_CorruptedJSON(t *testing.T) {
	// Use a temporary directory for test isolation
	tempDir := t.TempDir()
	SetConfigDirForTesting(tempDir)
	defer SetConfigDirForTesting("")

	path, err := StatisticsPath()
	if err != nil {
		t.Fatalf("StatisticsPath() error = %v", err)
	}

	// Write corrupted JSON to the file
	corruptedData := []byte(`{"sessions": [{"id": "test", "mode": `)
	if err := os.WriteFile(path, corruptedData, 0600); err != nil {
		t.Fatalf("Failed to write corrupted data: %v", err)
	}

	// Load should return an error
	stats, err := Load()
	if err == nil {
		t.Error("Load() should return error for corrupted JSON")
	}
	if stats != nil {
		t.Error("Load() should return nil stats for corrupted JSON")
	}
}

func TestConfigDir(t *testing.T) {
	dir, err := ConfigDir()
	if err != nil {
		t.Fatalf("ConfigDir() error = %v", err)
	}

	// Should end with "arithmego"
	if filepath.Base(dir) != "arithmego" {
		t.Errorf("ConfigDir() should end with 'arithmego', got %s", dir)
	}

	// Verify directory was created
	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("Directory not created: %v", err)
	}
	if !info.IsDir() {
		t.Error("Expected a directory")
	}
}

package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

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

func TestLoad_NullSessions(t *testing.T) {
	// Use a temporary directory for test isolation
	tempDir := t.TempDir()
	SetConfigDirForTesting(tempDir)
	defer SetConfigDirForTesting("")

	path, err := StatisticsPath()
	if err != nil {
		t.Fatalf("StatisticsPath() error = %v", err)
	}

	// Write JSON with null sessions
	nullData := []byte(`{"sessions": null}`)
	if err := os.WriteFile(path, nullData, 0600); err != nil {
		t.Fatalf("Failed to write null sessions data: %v", err)
	}

	// Load should return empty slice, not nil
	stats, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if stats.Sessions == nil {
		t.Error("Load() should return empty slice for null sessions, not nil")
	}
	if len(stats.Sessions) != 0 {
		t.Errorf("Load() should return empty slice, got %d sessions", len(stats.Sessions))
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

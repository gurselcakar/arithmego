package storage

import (
	"os"
	"testing"
)

func TestSettingsHasLastPlayed(t *testing.T) {
	tests := []struct {
		name     string
		settings Settings
		want     bool
	}{
		{
			name:     "empty settings",
			settings: Settings{},
			want:     false,
		},
		{
			name: "missing mode ID",
			settings: Settings{
				LastPlayedDifficulty: "Medium",
				LastPlayedDurationMs: 60000,
			},
			want: false,
		},
		{
			name: "missing difficulty",
			settings: Settings{
				LastPlayedModeID:     "addition",
				LastPlayedDurationMs: 60000,
			},
			want: false,
		},
		{
			name: "missing duration",
			settings: Settings{
				LastPlayedModeID:     "addition",
				LastPlayedDifficulty: "Medium",
			},
			want: false,
		},
		{
			name: "zero duration",
			settings: Settings{
				LastPlayedModeID:     "addition",
				LastPlayedDifficulty: "Medium",
				LastPlayedDurationMs: 0,
			},
			want: false,
		},
		{
			name: "negative duration",
			settings: Settings{
				LastPlayedModeID:     "addition",
				LastPlayedDifficulty: "Medium",
				LastPlayedDurationMs: -60000,
			},
			want: false,
		},
		{
			name: "complete settings",
			settings: Settings{
				LastPlayedModeID:     "addition",
				LastPlayedDifficulty: "Medium",
				LastPlayedDurationMs: 60000,
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.settings.HasLastPlayed(); got != tt.want {
				t.Errorf("HasLastPlayed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadSaveSettings(t *testing.T) {
	// Use a temporary directory for test isolation
	tempDir := t.TempDir()
	SetConfigDirForTesting(tempDir)
	defer SetConfigDirForTesting("")

	// Test loading non-existent file returns empty settings
	settings, err := LoadSettings()
	if err != nil {
		t.Fatalf("LoadSettings() error = %v", err)
	}
	if settings.HasLastPlayed() {
		t.Error("Expected empty settings for non-existent file")
	}

	// Save settings
	settings.LastPlayedModeID = "addition"
	settings.LastPlayedDifficulty = "Hard"
	settings.LastPlayedDurationMs = 120000

	err = SaveSettings(settings)
	if err != nil {
		t.Fatalf("SaveSettings() error = %v", err)
	}

	// Verify file exists
	path, err := SettingsPath()
	if err != nil {
		t.Fatalf("SettingsPath() error = %v", err)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("Settings file should exist after save")
	}

	// Load and verify
	loaded, err := LoadSettings()
	if err != nil {
		t.Fatalf("LoadSettings() error = %v", err)
	}
	if loaded.LastPlayedModeID != "addition" {
		t.Errorf("LastPlayedModeID = %s, want addition", loaded.LastPlayedModeID)
	}
	if loaded.LastPlayedDifficulty != "Hard" {
		t.Errorf("LastPlayedDifficulty = %s, want Hard", loaded.LastPlayedDifficulty)
	}
	if loaded.LastPlayedDurationMs != 120000 {
		t.Errorf("LastPlayedDurationMs = %d, want 120000", loaded.LastPlayedDurationMs)
	}
	if !loaded.HasLastPlayed() {
		t.Error("HasLastPlayed() should return true after loading valid settings")
	}
}

func TestLoadSettings_CorruptedJSON(t *testing.T) {
	// Use a temporary directory for test isolation
	tempDir := t.TempDir()
	SetConfigDirForTesting(tempDir)
	defer SetConfigDirForTesting("")

	path, err := SettingsPath()
	if err != nil {
		t.Fatalf("SettingsPath() error = %v", err)
	}

	// Write corrupted JSON to the file
	corruptedData := []byte(`{"last_played_mode_id": "test", "last_played`)
	if err := os.WriteFile(path, corruptedData, 0600); err != nil {
		t.Fatalf("Failed to write corrupted data: %v", err)
	}

	// LoadSettings should return empty settings on parse error (graceful handling)
	settings, err := LoadSettings()
	if err != nil {
		t.Errorf("LoadSettings() should not return error for corrupted JSON, got %v", err)
	}
	if settings.HasLastPlayed() {
		t.Error("LoadSettings() should return empty settings for corrupted JSON")
	}
}

func TestSettingsPath(t *testing.T) {
	path, err := SettingsPath()
	if err != nil {
		t.Fatalf("SettingsPath() error = %v", err)
	}

	// Should contain settings.json
	if path == "" {
		t.Error("SettingsPath() should not be empty")
	}
}

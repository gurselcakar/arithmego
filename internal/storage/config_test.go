package storage

import (
	"os"
	"testing"
)

func TestNewConfig(t *testing.T) {
	config := NewConfig()

	if config.DefaultDifficulty != DefaultDifficulty {
		t.Errorf("DefaultDifficulty = %s, want %s", config.DefaultDifficulty, DefaultDifficulty)
	}
	if config.DefaultDurationMs != DefaultDurationMs {
		t.Errorf("DefaultDurationMs = %d, want %d", config.DefaultDurationMs, DefaultDurationMs)
	}
	if config.AutoUpdate != DefaultAutoUpdate {
		t.Errorf("AutoUpdate = %v, want %v", config.AutoUpdate, DefaultAutoUpdate)
	}
}

func TestConfigHasLastPlayed(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		want   bool
	}{
		{
			name:   "empty config",
			config: Config{},
			want:   false,
		},
		{
			name: "missing mode ID",
			config: Config{
				LastPlayedDifficulty: "Normal",
				LastPlayedDurationMs: 60000,
			},
			want: false,
		},
		{
			name: "missing difficulty",
			config: Config{
				LastPlayedModeID:     "addition",
				LastPlayedDurationMs: 60000,
			},
			want: false,
		},
		{
			name: "missing duration",
			config: Config{
				LastPlayedModeID:     "addition",
				LastPlayedDifficulty: "Normal",
			},
			want: false,
		},
		{
			name: "zero duration",
			config: Config{
				LastPlayedModeID:     "addition",
				LastPlayedDifficulty: "Normal",
				LastPlayedDurationMs: 0,
			},
			want: false,
		},
		{
			name: "negative duration",
			config: Config{
				LastPlayedModeID:     "addition",
				LastPlayedDifficulty: "Normal",
				LastPlayedDurationMs: -60000,
			},
			want: false,
		},
		{
			name: "complete last played data",
			config: Config{
				LastPlayedModeID:     "addition",
				LastPlayedDifficulty: "Normal",
				LastPlayedDurationMs: 60000,
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.HasLastPlayed(); got != tt.want {
				t.Errorf("HasLastPlayed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadSaveConfig(t *testing.T) {
	tempDir := t.TempDir()
	SetConfigDirForTesting(tempDir)
	defer SetConfigDirForTesting("")

	// Test loading non-existent file returns default config
	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	if config.DefaultDifficulty != DefaultDifficulty {
		t.Errorf("Expected default difficulty for non-existent file")
	}
	if config.AutoUpdate != DefaultAutoUpdate {
		t.Errorf("Expected default auto-update for non-existent file")
	}

	// Modify and save config
	config.DefaultDifficulty = "Hard"
	config.DefaultDurationMs = 90000
	config.AutoUpdate = false
	config.LastPlayedModeID = "multiplication"
	config.LastPlayedDifficulty = "Expert"
	config.LastPlayedDurationMs = 120000

	err = SaveConfig(config)
	if err != nil {
		t.Fatalf("SaveConfig() error = %v", err)
	}

	// Verify file exists
	path, err := ConfigPath()
	if err != nil {
		t.Fatalf("ConfigPath() error = %v", err)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("Config file should exist after save")
	}

	// Load and verify
	loaded, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	if loaded.DefaultDifficulty != "Hard" {
		t.Errorf("DefaultDifficulty = %s, want Hard", loaded.DefaultDifficulty)
	}
	if loaded.DefaultDurationMs != 90000 {
		t.Errorf("DefaultDurationMs = %d, want 90000", loaded.DefaultDurationMs)
	}
	if loaded.AutoUpdate != false {
		t.Errorf("AutoUpdate = %v, want false", loaded.AutoUpdate)
	}
	if loaded.LastPlayedModeID != "multiplication" {
		t.Errorf("LastPlayedModeID = %s, want multiplication", loaded.LastPlayedModeID)
	}
	if loaded.LastPlayedDifficulty != "Expert" {
		t.Errorf("LastPlayedDifficulty = %s, want Expert", loaded.LastPlayedDifficulty)
	}
	if loaded.LastPlayedDurationMs != 120000 {
		t.Errorf("LastPlayedDurationMs = %d, want 120000", loaded.LastPlayedDurationMs)
	}
	if !loaded.HasLastPlayed() {
		t.Error("HasLastPlayed() should return true after loading valid config")
	}
}

func TestLoadConfig_CorruptedJSON(t *testing.T) {
	tempDir := t.TempDir()
	SetConfigDirForTesting(tempDir)
	defer SetConfigDirForTesting("")

	path, err := ConfigPath()
	if err != nil {
		t.Fatalf("ConfigPath() error = %v", err)
	}

	// Write corrupted JSON to the file
	corruptedData := []byte(`{"default_difficulty": "Hard", "auto_update`)
	if err := os.WriteFile(path, corruptedData, 0600); err != nil {
		t.Fatalf("Failed to write corrupted data: %v", err)
	}

	// LoadConfig should return default config on parse error
	config, err := LoadConfig()
	if err != nil {
		t.Errorf("LoadConfig() should not return error for corrupted JSON, got %v", err)
	}
	if config.DefaultDifficulty != DefaultDifficulty {
		t.Error("LoadConfig() should return default config for corrupted JSON")
	}
}

func TestLoadConfig_AppliesDefaults(t *testing.T) {
	tempDir := t.TempDir()
	SetConfigDirForTesting(tempDir)
	defer SetConfigDirForTesting("")

	path, err := ConfigPath()
	if err != nil {
		t.Fatalf("ConfigPath() error = %v", err)
	}

	// Write partial config (missing defaults)
	partialData := []byte(`{"auto_update": false, "last_played_mode_id": "test"}`)
	if err := os.WriteFile(path, partialData, 0600); err != nil {
		t.Fatalf("Failed to write partial data: %v", err)
	}

	// LoadConfig should apply defaults for missing fields
	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	if config.DefaultDifficulty != DefaultDifficulty {
		t.Errorf("DefaultDifficulty should be set to default, got %s", config.DefaultDifficulty)
	}
	if config.DefaultDurationMs != DefaultDurationMs {
		t.Errorf("DefaultDurationMs should be set to default, got %d", config.DefaultDurationMs)
	}
	// Explicit values should be preserved
	if config.AutoUpdate != false {
		t.Error("AutoUpdate should preserve explicit false value")
	}
	if config.LastPlayedModeID != "test" {
		t.Error("LastPlayedModeID should be preserved")
	}
}

func TestConfigPath(t *testing.T) {
	path, err := ConfigPath()
	if err != nil {
		t.Fatalf("ConfigPath() error = %v", err)
	}

	if path == "" {
		t.Error("ConfigPath() should not be empty")
	}
}

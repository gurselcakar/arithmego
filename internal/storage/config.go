package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Default values for new configurations.
const (
	DefaultDifficulty   = "Normal"
	DefaultDurationMs   = 60000 // 60 seconds
	DefaultAutoUpdate   = true
)

// Config stores user preferences and Quick Play state.
type Config struct {
	// Defaults (applied to Launch screen)
	DefaultDifficulty string `json:"default_difficulty,omitempty"`
	DefaultDurationMs int64  `json:"default_duration_ms,omitempty"`

	// Quick Play state (auto-saved after each game)
	LastPlayedModeID     string `json:"last_played_mode_id,omitempty"`
	LastPlayedDifficulty string `json:"last_played_difficulty,omitempty"`
	LastPlayedDurationMs int64  `json:"last_played_duration_ms,omitempty"`

	// Preferences
	AutoUpdate bool `json:"auto_update"`
}

// NewConfig creates a new Config with default values.
func NewConfig() *Config {
	return &Config{
		DefaultDifficulty: DefaultDifficulty,
		DefaultDurationMs: DefaultDurationMs,
		AutoUpdate:        DefaultAutoUpdate,
	}
}

// HasLastPlayed returns true if config contains valid last played data.
func (c *Config) HasLastPlayed() bool {
	return c.LastPlayedModeID != "" && c.LastPlayedDifficulty != "" && c.LastPlayedDurationMs > 0
}

// LoadConfig reads config from the JSON file.
// Returns default config if the file doesn't exist or is invalid.
func LoadConfig() (*Config, error) {
	path, err := ConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return NewConfig(), nil
		}
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		// Return default config on parse error.
		// Config is non-critical and can be regenerated.
		return NewConfig(), nil
	}

	// Apply defaults for any missing fields
	if config.DefaultDifficulty == "" {
		config.DefaultDifficulty = DefaultDifficulty
	}
	if config.DefaultDurationMs == 0 {
		config.DefaultDurationMs = DefaultDurationMs
	}

	return &config, nil
}

// SaveConfig writes config to the JSON file using atomic write.
func SaveConfig(config *Config) error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	// Write to temp file first for atomic operation.
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, "config-*.tmp")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()

	// Clean up temp file on any error
	shouldCleanup := true
	defer func() {
		if shouldCleanup {
			os.Remove(tmpPath)
		}
	}()

	if err := tmp.Chmod(0600); err != nil {
		_ = tmp.Close()
		return err
	}

	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return err
	}

	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return err
	}

	if err := tmp.Close(); err != nil {
		return err
	}

	if err := os.Rename(tmpPath, path); err != nil {
		return err
	}

	shouldCleanup = false
	return nil
}

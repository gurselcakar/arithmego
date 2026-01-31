package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gurselcakar/arithmego/internal/storage"
	"github.com/gurselcakar/arithmego/internal/ui"
)

func TestDeterminePlayStartMode(t *testing.T) {
	// Save original config path and restore after tests
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	t.Run("returns onboarding when no config exists", func(t *testing.T) {
		// Use temp dir as HOME to ensure no config exists
		tmpDir := t.TempDir()
		os.Setenv("HOME", tmpDir)

		mode := determinePlayStartMode()
		if mode != ui.StartModeOnboarding {
			t.Errorf("expected StartModeOnboarding, got %v", mode)
		}
	})

	t.Run("returns onboarding when config has no last played data", func(t *testing.T) {
		tmpDir := t.TempDir()
		os.Setenv("HOME", tmpDir)

		// Create config dir and empty config
		configDir := filepath.Join(tmpDir, ".config", "arithmego")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			t.Fatal(err)
		}

		config := storage.NewConfig()
		if err := storage.SaveConfig(config); err != nil {
			t.Fatal(err)
		}

		mode := determinePlayStartMode()
		if mode != ui.StartModeOnboarding {
			t.Errorf("expected StartModeOnboarding, got %v", mode)
		}
	})

	t.Run("returns quick play when config has last played data", func(t *testing.T) {
		tmpDir := t.TempDir()
		os.Setenv("HOME", tmpDir)

		// Create config dir
		configDir := filepath.Join(tmpDir, ".config", "arithmego")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			t.Fatal(err)
		}

		// Create config with last played data
		config := storage.NewConfig()
		config.LastPlayedModeID = "addition"
		config.LastPlayedDifficulty = "Easy"
		config.LastPlayedDurationMs = 60000
		if err := storage.SaveConfig(config); err != nil {
			t.Fatal(err)
		}

		mode := determinePlayStartMode()
		if mode != ui.StartModeQuickPlay {
			t.Errorf("expected StartModeQuickPlay, got %v", mode)
		}
	})

	t.Run("returns onboarding when last played data is incomplete", func(t *testing.T) {
		tmpDir := t.TempDir()
		os.Setenv("HOME", tmpDir)

		configDir := filepath.Join(tmpDir, ".config", "arithmego")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			t.Fatal(err)
		}

		// Config with partial last played data (missing duration)
		config := storage.NewConfig()
		config.LastPlayedModeID = "addition"
		config.LastPlayedDifficulty = "Easy"
		// LastPlayedDurationMs is 0 (invalid)
		if err := storage.SaveConfig(config); err != nil {
			t.Fatal(err)
		}

		mode := determinePlayStartMode()
		if mode != ui.StartModeOnboarding {
			t.Errorf("expected StartModeOnboarding for incomplete data, got %v", mode)
		}
	})
}

func TestVersionVariables(t *testing.T) {
	// Test that version variables have default values
	t.Run("default version is dev", func(t *testing.T) {
		if Version != "dev" {
			t.Errorf("expected default Version to be 'dev', got %q", Version)
		}
	})

	t.Run("default commit is unknown", func(t *testing.T) {
		if CommitSHA != "unknown" {
			t.Errorf("expected default CommitSHA to be 'unknown', got %q", CommitSHA)
		}
	})

	t.Run("default build date is unknown", func(t *testing.T) {
		if BuildDate != "unknown" {
			t.Errorf("expected default BuildDate to be 'unknown', got %q", BuildDate)
		}
	})
}

func TestRootCommandSetup(t *testing.T) {
	t.Run("root command has correct use", func(t *testing.T) {
		if rootCmd.Use != "arithmego" {
			t.Errorf("expected Use to be 'arithmego', got %q", rootCmd.Use)
		}
	})

	t.Run("subcommands are registered", func(t *testing.T) {
		expectedCommands := []string{"play", "statistics", "update", "version"}
		commands := rootCmd.Commands()

		for _, expected := range expectedCommands {
			found := false
			for _, cmd := range commands {
				if cmd.Use == expected || cmd.Name() == expected {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("expected subcommand %q to be registered", expected)
			}
		}
	})

	t.Run("statistics has stats alias", func(t *testing.T) {
		for _, cmd := range rootCmd.Commands() {
			if cmd.Use == "statistics" {
				hasAlias := false
				for _, alias := range cmd.Aliases {
					if alias == "stats" {
						hasAlias = true
						break
					}
				}
				if !hasAlias {
					t.Error("statistics command should have 'stats' alias")
				}
				return
			}
		}
		t.Error("statistics command not found")
	})
}

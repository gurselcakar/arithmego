package cli

import (
	"testing"

	"github.com/spf13/cobra"

	"github.com/gurselcakar/arithmego/internal/modes"
)

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

func TestPlayCommandSetup(t *testing.T) {
	t.Run("play command accepts mode argument", func(t *testing.T) {
		// Find play command
		var playCommand *cobra.Command
		for _, cmd := range rootCmd.Commands() {
			if cmd.Name() == "play" {
				playCommand = cmd
				break
			}
		}
		if playCommand == nil {
			t.Fatal("play command not found")
		}

		// Check Use includes [mode]
		if playCommand.Use != "play [mode]" {
			t.Errorf("expected Use to be 'play [mode]', got %q", playCommand.Use)
		}
	})

	t.Run("play command has valid args function for completion", func(t *testing.T) {
		var playCommand *cobra.Command
		for _, cmd := range rootCmd.Commands() {
			if cmd.Name() == "play" {
				playCommand = cmd
				break
			}
		}
		if playCommand == nil {
			t.Fatal("play command not found")
		}

		if playCommand.ValidArgsFunction == nil {
			t.Error("play command should have ValidArgsFunction for tab completion")
		}
	})
}

func TestPlayCommandModeValidation(t *testing.T) {
	// Ensure modes are registered for validation
	modes.RegisterPresets()

	tests := []struct {
		modeID string
		valid  bool
	}{
		// Basic modes
		{"addition", true},
		{"subtraction", true},
		{"multiplication", true},
		{"division", true},
		// Power modes
		{"squares", true},
		{"cubes", true},
		{"square-roots", true},
		{"cube-roots", true},
		// Advanced modes
		{"exponents", true},
		{"remainders", true},
		{"percentages", true},
		{"factorials", true},
		// Mixed modes
		{"mixed-basics", true},
		{"mixed-powers", true},
		{"mixed-advanced", true},
		{"anything-goes", true},
		// Invalid modes
		{"invalid-mode", false},
		{"", false},
		{"add", false},
	}

	for _, tt := range tests {
		t.Run(tt.modeID, func(t *testing.T) {
			_, ok := modes.Get(tt.modeID)
			if ok != tt.valid {
				t.Errorf("modes.Get(%q) = %v, want %v", tt.modeID, ok, tt.valid)
			}
		})
	}
}

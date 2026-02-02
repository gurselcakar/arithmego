package cli

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/gurselcakar/arithmego/internal/modes"
	"github.com/gurselcakar/arithmego/internal/ui"

	// Register operations (must come before modes.RegisterPresets)
	_ "github.com/gurselcakar/arithmego/internal/game/operations"
)

// Version info - set via ldflags during build
var (
	Version   = "dev"
	CommitSHA = "unknown"
	BuildDate = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "arithmego",
	Short: "Terminal-based arithmetic game for developers",
	Long: `ArithmeGo is a terminal-based arithmetic game designed for developers.
Short sessions. Minimal friction. Never leave the terminal.

Run without arguments to open the main menu.`,
	Run: func(cmd *cobra.Command, args []string) {
		runTUI(ui.StartModeMenu)
	},
}

// Execute runs the root command.
func Execute() {
	// Initialize modes
	modes.RegisterPresets()

	// Note: Update check is now handled within the TUI (see ui/app.go)
	// This allows the notification to be displayed in the menu screen.

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// runTUI starts the Bubble Tea application with the specified start mode.
func runTUI(startMode ui.StartMode) {
	// Recover from panics to ensure terminal is restored properly.
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "\nUnexpected error: %v\n", r)
			os.Exit(1)
		}
	}()

	// Set version for update checking within TUI
	ui.Version = Version

	app := ui.NewWithStartMode(startMode)
	p := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Disable Cobra's default completion command
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

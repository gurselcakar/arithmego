package cli

import (
	"github.com/spf13/cobra"

	"github.com/gurselcakar/arithmego/internal/storage"
	"github.com/gurselcakar/arithmego/internal/ui"
)

var playCmd = &cobra.Command{
	Use:   "play",
	Short: "Start a quick play session",
	Long: `Start a game immediately with your last played settings.

If you haven't played before, the onboarding flow will guide you
through selecting your preferences.`,
	Run: func(cmd *cobra.Command, args []string) {
		runTUI(determinePlayStartMode())
	},
}

// determinePlayStartMode decides whether to start quick play or onboarding.
// Returns StartModeQuickPlay if user has last played data, otherwise StartModeOnboarding.
func determinePlayStartMode() ui.StartMode {
	config, _ := storage.LoadConfig()
	if config != nil && config.HasLastPlayed() {
		return ui.StartModeQuickPlay
	}
	// No last played data - launch onboarding
	return ui.StartModeOnboarding
}

func init() {
	rootCmd.AddCommand(playCmd)
}

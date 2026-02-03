package cli

import (
	"github.com/spf13/cobra"

	"github.com/gurselcakar/arithmego/internal/ui"
)

var settingsCmd = &cobra.Command{
	Use:   "settings",
	Short: "Open settings",
	Long: `Open the settings screen to configure your preferences.

Adjust default difficulty, duration, input method, and auto-update behavior.
Press Esc to exit.`,
	Run: func(cmd *cobra.Command, args []string) {
		runTUI(ui.StartModeSettings)
	},
}

func init() {
	rootCmd.AddCommand(settingsCmd)
}

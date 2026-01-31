package cli

import (
	"github.com/spf13/cobra"

	"github.com/gurselcakar/arithmego/internal/ui"
)

var statisticsCmd = &cobra.Command{
	Use:     "statistics",
	Aliases: []string{"stats"},
	Short:   "View your performance statistics",
	Long: `Open the statistics screen to view your game history and performance.

Shows overall accuracy, per-operation breakdown, and best streaks.
Press 'D' for detailed view, 'S' for summary, Esc to exit.`,
	Run: func(cmd *cobra.Command, args []string) {
		runTUI(ui.StartModeStatistics)
	},
}

func init() {
	rootCmd.AddCommand(statisticsCmd)
}

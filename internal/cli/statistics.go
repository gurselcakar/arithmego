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

The dashboard shows overall accuracy, personal bests, and insights.
Navigate between views using: O (Operations), H (History), T (Trends).
Press Esc to return to menu.`,
	Run: func(cmd *cobra.Command, args []string) {
		runTUI(ui.StartModeStatistics)
	},
}

func init() {
	rootCmd.AddCommand(statisticsCmd)
}

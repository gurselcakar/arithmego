package cli

import (
	"github.com/spf13/cobra"

	"github.com/gurselcakar/arithmego/internal/ui"
)

var practiceCmd = &cobra.Command{
	Use:   "practice",
	Short: "Start practice mode",
	Long: `Open the practice screen to drill arithmetic problems at your own pace.

Choose a category, operation, and difficulty level.
Press Esc to return to menu.`,
	Run: func(cmd *cobra.Command, args []string) {
		runTUI(ui.StartModePractice)
	},
}

func init() {
	rootCmd.AddCommand(practiceCmd)
}

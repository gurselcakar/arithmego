package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/gurselcakar/arithmego/internal/modes"
	"github.com/gurselcakar/arithmego/internal/ui"
)

var playCmd = &cobra.Command{
	Use:   "play [mode]",
	Short: "Start a game session",
	Long: `Open the play screen to browse and select a game mode.

If a mode is specified, opens the configuration screen for that mode directly.

Available modes:
  Basic:    addition, subtraction, multiplication, division
  Powers:   squares, cubes, square-roots, cube-roots
  Advanced: exponents, remainders, percentages, factorials
  Mixed:    mixed-basics, mixed-powers, mixed-advanced, anything-goes

Examples:
  arithmego play              # Browse all modes
  arithmego play addition     # Configure Addition mode
  arithmego play mixed-basics # Configure Mixed Basics mode`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			// No mode specified - open play browse
			runTUI(ui.StartModePlayBrowse)
			return
		}

		// Mode specified - validate and open play config
		modeID := args[0]
		mode, ok := modes.Get(modeID)
		if !ok {
			fmt.Fprintf(os.Stderr, "Unknown mode: %s\n\n", modeID)
			fmt.Fprintln(os.Stderr, "Available modes:")
			fmt.Fprintln(os.Stderr, "  Basic:    addition, subtraction, multiplication, division")
			fmt.Fprintln(os.Stderr, "  Powers:   squares, cubes, square-roots, cube-roots")
			fmt.Fprintln(os.Stderr, "  Advanced: exponents, remainders, percentages, factorials")
			fmt.Fprintln(os.Stderr, "  Mixed:    mixed-basics, mixed-powers, mixed-advanced, anything-goes")
			os.Exit(1)
		}

		ui.CLIModeID = mode.ID
		runTUI(ui.StartModePlayConfig)
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		// Return all mode IDs for tab completion
		var completions []string
		for _, mode := range modes.All() {
			if strings.HasPrefix(mode.ID, toComplete) {
				completions = append(completions, mode.ID+"\t"+mode.Name)
			}
		}
		return completions, cobra.ShellCompDirectiveNoFileComp
	},
}

func init() {
	rootCmd.AddCommand(playCmd)
}

package cli

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  `Display the current version of ArithmeGo along with build information.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("ArithmeGo %s\n", Version)
		if CommitSHA != "unknown" {
			fmt.Printf("Commit:    %s\n", CommitSHA)
		}
		if BuildDate != "unknown" {
			fmt.Printf("Built:     %s\n", BuildDate)
		}
		fmt.Printf("Go:        %s\n", runtime.Version())
		fmt.Printf("OS/Arch:   %s/%s\n", runtime.GOOS, runtime.GOARCH)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

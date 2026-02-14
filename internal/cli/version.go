package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  `Display the current version of ArithmeGo.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("ArithmeGo %s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

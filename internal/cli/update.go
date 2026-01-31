package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gurselcakar/arithmego/internal/update"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Check for updates",
	Long: `Check if a newer version of ArithmeGo is available.

If an update is available, instructions for upgrading will be displayed.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Checking for updates...")
		fmt.Println()

		info, err := update.Check(Version)
		if err != nil {
			fmt.Printf("Error checking for updates: %v\n", err)
			fmt.Println()
			fmt.Println("You can manually check for updates at:")
			fmt.Println("  https://github.com/gurselcakar/arithmego/releases")
			return
		}

		fmt.Printf("Current version: %s\n", info.CurrentVersion)
		fmt.Printf("Latest version:  %s\n", info.LatestVersion)
		fmt.Println()

		if info.UpdateAvailable {
			fmt.Println("A new version is available!")
			fmt.Println()
			fmt.Println("To update, run:")
			fmt.Println("  curl -fsSL https://arithmego.com/install.sh | bash")
			fmt.Println()
			fmt.Println("Or download directly from:")
			fmt.Printf("  %s\n", info.ReleaseURL)
		} else {
			fmt.Println("You're running the latest version.")
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

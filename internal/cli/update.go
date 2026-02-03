package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gurselcakar/arithmego/internal/update"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Check for and install updates",
	Long: `Check if a newer version of ArithmeGo is available and install it.

If an update is available, it will be downloaded and installed automatically.
If automatic installation fails, manual instructions will be displayed.`,
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

		if !info.UpdateAvailable {
			fmt.Println("You're running the latest version.")
			return
		}

		fmt.Printf("Downloading %s...\n", info.LatestVersion)
		if err := update.DownloadAndApply(info.LatestVersion); err != nil {
			fmt.Printf("Automatic update failed: %v\n", err)
			fmt.Println()
			fmt.Println("To update manually, run:")
			fmt.Println("  curl -fsSL https://arithmego.com/install.sh | bash")
			fmt.Println()
			fmt.Println("Or download directly from:")
			fmt.Printf("  %s\n", info.ReleaseURL)
			return
		}

		fmt.Printf("Updated to %s. Restart arithmego to use the new version.\n", info.LatestVersion)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

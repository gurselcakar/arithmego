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

		info, err := update.Check(Version)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			fmt.Println("Check manually: https://github.com/gurselcakar/arithmego/releases")
			return
		}

		if !info.UpdateAvailable {
			fmt.Printf("ArithmeGo %s is up to date.\n", info.LatestVersion)
			return
		}

		fmt.Printf("Updating %s → %s...\n", info.CurrentVersion, info.LatestVersion)
		if err := update.DownloadAndApply(info.LatestVersion); err != nil {
			fmt.Printf("Automatic update failed: %v\n", err)
			fmt.Println()
			fmt.Println("To update manually, run:")
			fmt.Println("  curl -fsSL https://arithmego.com/install.sh | bash")
			return
		}

		fmt.Printf("Updated to %s. Restart arithmego to use the new version.\n", info.LatestVersion)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

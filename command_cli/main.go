package commandcli

import (
	"fmt"
	"os"
	"strings"

	"github.com/Netsocs-Team/netsocs-manager-cli/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func UpdateCLICommand(cmd *cobra.Command, args []string) {
	var version string
	if len(args) > 0 {
		version = args[0]
		pterm.Info.Printfln("Updating CLI to version: %s", version)
	} else {
		pterm.Info.Println("Updating CLI to the latest version available")
	}

	if err := utils.DownloadAndReplaceCLI(version); err != nil {
		pterm.Error.Printfln("Error updating CLI: %v", err)
		os.Exit(1)
	}

	pterm.Success.Println("CLI updated successfully! Please restart the CLI.")
}

func ListCLIVersionsCommand(cmd *cobra.Command, args []string, currentVer string) {
	// Try to extract only the version part from e.g. netsocs-helm-chart-1.0.1
	if idx := len(currentVer) - 1 - len("-"); idx > 0 {
		parts := strings.Split(currentVer, "-")
		if len(parts) > 0 {
			currentVer = parts[len(parts)-1]
		}
	}
	versions, err := utils.ListAvailableCLIVersions()
	if err != nil {
		fmt.Printf("Error fetching CLI versions: %v\n", err)
		return
	}
	fmt.Println("Available CLI versions:")
	for _, v := range versions {
		if v == currentVer {
			fmt.Printf("* %s (in use)\n", v)
		} else {
			fmt.Printf("  %s\n", v)
		}
	}
}

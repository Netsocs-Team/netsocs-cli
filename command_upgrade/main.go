package commandupgrade

import (
	"os"

	"github.com/Netsocs-Team/netsocs-manager-cli/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func UpgradeCommand(cmd *cobra.Command, args []string) {
	var version string
	if len(args) > 0 {
		version = args[0]
		pterm.Info.Printfln("Upgrading to version: %s", version)
	} else {
		pterm.Info.Println("Upgrading to the latest version available")
	}

	if err := utils.RunHelmUpgradeWithVersion(version); err != nil {
		pterm.Error.Printfln("Error upgrading application: %v", err)
		os.Exit(1)
	}

	pterm.Success.Println("Upgrade completed successfully!")
}

func RollbackCommand(cmd *cobra.Command, args []string) {
	var revision string
	if len(args) > 0 {
		revision = args[0]
		pterm.Info.Printfln("Rolling back to revision: %s", revision)
	} else {
		pterm.Error.Println("You must specify a revision number. Example: my-cli rollback 2")
		cmd.Usage()
		os.Exit(1)
	}

	if err := utils.RunHelmRollback(revision); err != nil {
		pterm.Error.Printfln("Error during rollback: %v", err)
		os.Exit(1)
	}

	pterm.Success.Println("Rollback completed successfully!")
}

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

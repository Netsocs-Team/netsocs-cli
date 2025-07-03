package main

import (
	"fmt"
	"os"
	"strings"

	commandcli "github.com/Netsocs-Team/netsocs-manager-cli/command_cli"
	commandconfig "github.com/Netsocs-Team/netsocs-manager-cli/command_config"
	commandinit "github.com/Netsocs-Team/netsocs-manager-cli/command_init"
	commandstatus "github.com/Netsocs-Team/netsocs-manager-cli/command_status"
	commandupgrade "github.com/Netsocs-Team/netsocs-manager-cli/command_upgrade"
	"github.com/Netsocs-Team/netsocs-manager-cli/utils"

	"github.com/spf13/cobra"
)

var version = "v0.0.0"

var rootCmd = &cobra.Command{
	Use:     "netsocs-manager-cli",
	Short:   "Server configuration tool",
	Version: version,
}

type ChartValues struct {
	HttpHostname string `yaml:"httpHostname"`
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes Helm configuration for Netsocs",
	Run:   commandinit.InitCommand,
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Base configuration for Netsocs",
	Run:   commandconfig.ConfigCommand,
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Shows the status of NETSOCS",
	Run:   commandstatus.StatusHandler,
}

var upgradeCmd = &cobra.Command{
	Use:   "upgrade [version]",
	Short: "Upgrade the application to a specific version or the latest if not specified",
	Args:  cobra.MaximumNArgs(1),
	Run:   commandupgrade.UpgradeCommand,
}

var rollbackCmd = &cobra.Command{
	Use:   "rollback <revision>",
	Short: "Rollback the application to a previous revision",
	Args:  cobra.MaximumNArgs(1),
	Run:   commandupgrade.RollbackCommand,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show CLI and netsocs version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("CLI version: %s\n", version)
		fmt.Printf("Netsocs version: %s\n", utils.GetCurrentAppVersion())
	},
}

var listVersionsCmd = &cobra.Command{
	Use:   "list-versions",
	Short: "Show the 10 latest available versions and mark the one in use",
	Run: func(cmd *cobra.Command, args []string) {
		current := utils.GetCurrentAppVersion()
		currentVer := current
		if idx := len(current) - 1 - len("-"); idx > 0 {
			parts := strings.Split(current, "-")
			if len(parts) > 0 {
				currentVer = parts[len(parts)-1]
			}
		}
		versions, err := utils.ListAvailableAppVersions()
		if err != nil {
			fmt.Printf("Error fetching versions: %v\n", err)
			return
		}
		fmt.Println("Available versions:")
		for _, v := range versions {
			if v == currentVer {
				fmt.Printf("* %s (in use)\n", v)
			} else {
				fmt.Printf("  %s\n", v)
			}
		}
	},
}

var cliCmd = &cobra.Command{
	Use:   "cli",
	Short: "CLI management commands",
}

var cliUpdateCmd = &cobra.Command{
	Use:   "update [version]",
	Short: "Update the CLI to a specific version or the latest if not specified",
	Args:  cobra.MaximumNArgs(1),
	Run:   commandcli.UpdateCLICommand,
}

var cliListVersionsCmd = &cobra.Command{
	Use:   "list-versions",
	Short: "Show the 10 latest available CLI versions and mark the one in use",
	Run: func(cmd *cobra.Command, args []string) {
		commandcli.ListCLIVersionsCommand(cmd, args, version)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(configCmd)
	statusCmd.Flags().BoolP("verbose", "v", false, "Show full pod details")
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(upgradeCmd)
	rootCmd.AddCommand(rollbackCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(listVersionsCmd)
	// CLI group
	cliCmd.AddCommand(cliUpdateCmd)
	cliCmd.AddCommand(cliListVersionsCmd)
	rootCmd.AddCommand(cliCmd)
}

func main() {
	if utils.GetCurrentAppVersion() == "not installed" {
		fmt.Println("Netsocs is not installed. Please run 'netsocs init' to install it.")
		os.Exit(1)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

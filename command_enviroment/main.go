package commandenviroment

import (
	"github.com/Netsocs-Team/netsocs-manager-cli/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func EnvironmentCommand(cmd *cobra.Command, args []string) {
	if !utils.CheckNetworkConnection() {
		pterm.Error.Println("🚨 Network connection is not working. Please check your internet connection.")
		return
	}
}

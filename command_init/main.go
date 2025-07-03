package commandinit

import (
	"github.com/Netsocs-Team/netsocs-manager-cli/utils"
	"github.com/spf13/cobra"
)

func InitCommand(cmd *cobra.Command, args []string) {
	utils.ShowBannerArt()

	if err := utils.InitializeHelmSetup(); err != nil {
		cmd.PrintErrf("Helm configuration error: %v\n", err)
		return
	}
}

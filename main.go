package main

import (
	"fmt"
	"os"

	commandconfig "github.com/Netsocs-Team/netsocs-manager-cli/command_config"
	commandstatus "github.com/Netsocs-Team/netsocs-manager-cli/command_status"
	"github.com/pterm/pterm"

	"github.com/spf13/cobra"
)

var version = "v0.0.0"

var rootCmd = &cobra.Command{
	Use:   "netsocs-manager-cli",
	Short: "Herramienta de configuración de servidor",
}

type ChartValues struct {
	HttpHostname string `yaml:"httpHostname"`
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuracion base de Netsocs",
	Run:   commandconfig.ConfigCommand,
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Muestra el estado de NETSOCS",
	Run:   commandstatus.StatusHandler,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Muestra la version de Netsocs",
	Run: func(cmd *cobra.Command, args []string) {
		pterm.DefaultHeader.
			WithBackgroundStyle(pterm.NewStyle(pterm.BgGreen)).
			WithTextStyle(pterm.NewStyle(pterm.FgBlack)).
			Println("Netsocs Manager CLI " + version)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	statusCmd.Flags().BoolP("verbose", "v", false, "Mostrar detalles completos de los pods")
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(versionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

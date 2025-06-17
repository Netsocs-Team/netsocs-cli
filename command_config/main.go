package commandconfig

import (
	"fmt"
	"os"
	"regexp"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Netsocs-Team/netsocs-manager-cli/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func ConfigCommand(cmd *cobra.Command, args []string) {
	utils.ShowBannerArt()
	address := promptAddress()
	// Actualizar el campo en values.yaml
	if err := utils.UpdateChartConfig("httpHostname", "https://"+address); err != nil {
		pterm.Error.Printfln("Error al actualizar configuración: %v", err)
		os.Exit(1)
	}

	// Ejecutar Helm upgrade
	if err := utils.RunHelmUpgrade(); err != nil {
		pterm.Error.Printfln("Error al ejecutar Helm: %v", err)
		os.Exit(1)
	}

	pterm.Success.Printfln("¡Configuración completada!")
	// pterm.Info.Printfln("Dirección configurada: %s", pterm.LightGreen(address))
}

func promptAddress() string {
	address := ""
	prompt := &survey.Input{
		Message: "Dirección de NETSOCS:",
		Help:    "Ejemplo: 192.168.1.1 o servidor.netsocs.com",
		Suggest: func(toComplete string) []string {
			return []string{
				"192.168.1.",
				"10.0.0.",
				"servidor.netsocs.com",
				"dns.netsocs.local",
			}
		},
	}

	validation := func(input interface{}) error {
		str, ok := input.(string)
		if !ok {
			return fmt.Errorf("tipo de dato inválido")
		}
		if isValidAddress(str) {
			return nil
		}
		return fmt.Errorf("formato inválido: debe ser IP (XXX.XXX.XXX.XXX) o dominio (ej: dns.netsocs.com)")
	}

	// survey.AskOne(prompt, &ip, survey.WithValidator(validation))
	err := survey.AskOne(prompt, &address,
		survey.WithIcons(func(icons *survey.IconSet) {
			icons.Question.Text = pterm.Green(">")
			icons.Question.Format = "green"
		}),
		survey.WithValidator(validation),
	)
	if err != nil {
		pterm.Error.Printfln("Error al obtener la dirección: %v", err)
		os.Exit(1)
	}
	return address
}

func isValidAddress(address string) bool {
	// Validación para dirección IP
	ipRegex := regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)
	if ipRegex.MatchString(address) {
		return true
	}

	// Validación para dominio
	domainRegex := regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9](\.[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9])*(\.[a-zA-Z]{2,})$`)
	if domainRegex.MatchString(address) {
		return true
	}

	// Permitir dominios locales sin TLD
	localDomainRegex := regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9](\.[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9])+$`)
	return localDomainRegex.MatchString(address)
}

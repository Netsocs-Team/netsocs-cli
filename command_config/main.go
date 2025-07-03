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
	// Update the field in values.yaml
	if err := utils.UpdateChartConfig("httpHostname", "https://"+address); err != nil {
		pterm.Error.Printfln("Error updating configuration: %v", err)
		os.Exit(1)
	}

	// Run Helm upgrade
	if err := utils.RunHelmUpgrade(); err != nil {
		pterm.Error.Printfln("Error running Helm: %v", err)
		os.Exit(1)
	}

	pterm.Success.Printfln("¡Configuration completed!")
	// pterm.Info.Printfln("Address configured: %s", pterm.LightGreen(address))
}

func promptAddress() string {
	address := ""
	prompt := &survey.Input{
		Message: "NETSOCS Address:",
		Help:    "Example: 192.168.1.1 or server.netsocs.com",
		Suggest: func(toComplete string) []string {
			return []string{
				"192.168.1.",
				"10.0.0.",
				"server.netsocs.com",
				"dns.netsocs.local",
			}
		},
	}

	validation := func(input interface{}) error {
		str, ok := input.(string)
		if !ok {
			return fmt.Errorf("invalid data type")
		}
		if isValidAddress(str) {
			return nil
		}
		return fmt.Errorf("invalid format: must be IP (XXX.XXX.XXX.XXX) or domain (e.g: dns.netsocs.com)")
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
		pterm.Error.Printfln("Error getting address: %v", err)
		os.Exit(1)
	}
	return address
}

func isValidAddress(address string) bool {
	// Validation for IP address
	ipRegex := regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)
	if ipRegex.MatchString(address) {
		return true
	}

	// Validation for domain
	domainRegex := regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9](\.[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9])*(\.[a-zA-Z]{2,})$`)
	if domainRegex.MatchString(address) {
		return true
	}

	// Allow local domains without TLD
	localDomainRegex := regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9](\.[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9])+$`)
	return localDomainRegex.MatchString(address)
}

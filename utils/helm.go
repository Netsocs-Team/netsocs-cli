package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pterm/pterm"
)

const (
	HelmRepoURL = "https://netsocs-team.github.io/netsocs-helm-chart/"
	AppName     = "netsocs"
)

type helmRelease struct {
	Name  string `json:"name"`
	Chart string `json:"chart"`
}

type helmChartVersion struct {
	Version string `json:"version"`
}

func InitializeHelmSetup() error {
	pterm.Info.Println("Initializing Helm configuration for Netsocs...")

	if err := checkHelmInstalled(); err != nil {
		return fmt.Errorf("Helm is not installed or not accessible: %w", err)
	}

	if err := addHelmRepo(); err != nil {
		return fmt.Errorf("error adding Helm repository: %w", err)
	}

	isInstalled, err := checkIfAppInstalled()
	if err != nil {
		return fmt.Errorf("error checking installation: %w", err)
	}

	if isInstalled {
		pterm.Success.Println("The netsocs application is already installed")
		return nil
	}

	valuesExists, err := checkValuesFileOrCreate()
	if err != nil {
		return fmt.Errorf("error checking/creating values.yaml file: %w", err)
	}

	if err := installNetsocsApp(valuesExists); err != nil {
		return fmt.Errorf("error installing application: %w", err)
	}

	pterm.Success.Println("Helm configuration completed successfully")
	return nil
}

func checkHelmInstalled() error {
	pterm.Info.Println("Checking Helm installation...")

	cmd := exec.Command("helm", "version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Helm is not installed or not accessible. Please install Helm first")
	}

	pterm.Success.Println("Helm is installed correctly")
	return nil
}

func addHelmRepo() error {
	pterm.Info.Println("Adding Helm repository...")

	cmd := exec.Command("helm", "repo", "list")
	output, _ := cmd.Output()

	if strings.Contains(string(output), "netsocs") {
		pterm.Info.Println("netsocs repository already exists")
		return nil
	}

	cmd = exec.Command("helm", "repo", "add", "netsocs", HelmRepoURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error adding repository: %w", err)
	}

	cmd = exec.Command("helm", "repo", "update")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error updating repositories: %w", err)
	}

	pterm.Success.Println("Repository added and updated successfully")
	return nil
}

func checkIfAppInstalled() (bool, error) {
	pterm.Info.Println("Checking if netsocs is installed...")

	cmd := exec.Command("helm", "list", "--output", "json")
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("error listing Helm applications: %w", err)
	}

	return strings.Contains(string(output), AppName), nil
}

func checkValuesFileOrCreate() (bool, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false, fmt.Errorf("error getting home directory: %w", err)
	}

	valuesDir := filepath.Join(homeDir, "netsocs")
	valuesPath := filepath.Join(valuesDir, "values.yaml")

	if _, err := os.Stat(valuesPath); os.IsNotExist(err) {
		pterm.Warning.Println("values.yaml file not found at ~/netsocs/values.yaml. Creating default file...")
		if err := os.MkdirAll(valuesDir, 0755); err != nil {
			return false, fmt.Errorf("error creating ~/netsocs directory: %w", err)
		}
		cmd := exec.Command("helm", "show", "values", "netsocs/netsocs-helm-chart ")
		output, err := cmd.Output()
		if err != nil {
			return false, fmt.Errorf("error getting default Helm values: %w", err)
		}
		if err := os.WriteFile(valuesPath, output, 0644); err != nil {
			return false, fmt.Errorf("error writing values.yaml: %w", err)
		}
		pterm.Success.Println("values.yaml file created successfully at ~/netsocs/values.yaml")
		return true, nil
	}

	pterm.Success.Println("values.yaml file found")
	return true, nil
}

func installNetsocsApp(hasValuesFile bool) error {
	pterm.Info.Println("Installing netsocs application...")

	var cmd *exec.Cmd

	if hasValuesFile {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("error getting home directory: %w", err)
		}
		valuesPath := filepath.Join(homeDir, "netsocs", "values.yaml")

		cmd = exec.Command("helm", "install", AppName, "netsocs/netsocs-helm-chart ",
			"--values", valuesPath)
	} else {
		cmd = exec.Command("helm", "install", AppName, "netsocs/netsocs-helm-chart ")
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error installing application: %w", err)
	}

	pterm.Success.Println("netsocs application installed successfully")
	return nil
}

func GetCurrentAppVersion() string {
	cmd := exec.Command("helm", "list", "--output", "json")
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	var releases []helmRelease
	if err := json.Unmarshal(output, &releases); err != nil {
		return "unknown"
	}
	for _, rel := range releases {
		if rel.Name == AppName {
			return rel.Chart
		}
	}
	return "not installed"
}

func ListAvailableAppVersions() ([]string, error) {
	cmd := exec.Command("helm", "search", "repo", "netsocs/netsocs-helm-chart", "--versions", "--output", "json")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	var versions []helmChartVersion
	if err := json.Unmarshal(output, &versions); err != nil {
		return nil, err
	}
	var result []string
	for i, v := range versions {
		if i >= 10 {
			break
		}
		result = append(result, v.Version)
	}
	return result, nil
}

func ListAvailableCLIVersions() ([]string, error) {
	const repo = "Netsocs-Team/netsocs-cli"
	resp, err := http.Get("https://api.github.com/repos/" + repo + "/releases")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var data []struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	var result []string
	for i, rel := range data {
		if i >= 10 {
			break
		}
		result = append(result, rel.TagName)
	}
	return result, nil
}

func DownloadAndReplaceCLI(version string) error {
	const repo = "Netsocs-Team/netsocs-cli"
	var downloadURL string
	if version == "" {
		// Get latest release info
		resp, err := http.Get("https://api.github.com/repos/" + repo + "/releases/latest")
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		var data struct {
			TagName string `json:"tag_name"`
			Assets  []struct {
				Name               string `json:"name"`
				BrowserDownloadURL string `json:"browser_download_url"`
			}
		}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return err
		}
		for _, asset := range data.Assets {
			if asset.Name == "netsocs" {
				downloadURL = asset.BrowserDownloadURL
				break
			}
		}
		if downloadURL == "" {
			return fmt.Errorf("no suitable binary found in latest release")
		}
	} else {
		// Specific version
		resp, err := http.Get("https://api.github.com/repos/" + repo + "/releases/tags/" + version)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		var data struct {
			Assets []struct {
				BrowserDownloadURL string `json:"browser_download_url"`
				Name               string `json:"name"`
			}
		}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return err
		}
		for _, asset := range data.Assets {
			if asset.Name == "netsocs" {
				downloadURL = asset.BrowserDownloadURL
				break
			}
		}
		if downloadURL == "" {
			return fmt.Errorf("no suitable binary found for version %s", version)
		}
	}

	// Download binary to $HOME/netsocs/netsocs.new
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	netsocsDir := filepath.Join(homeDir, "netsocs")
	if err := os.MkdirAll(netsocsDir, 0755); err != nil {
		return err
	}
	newBinPath := filepath.Join(netsocsDir, "netsocs.new")
	resp, err := http.Get(downloadURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	f, err := os.OpenFile(newBinPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	if _, err := io.Copy(f, resp.Body); err != nil {
		f.Close()
		return err
	}
	f.Close()

	// Create update.sh script
	updateScriptPath := filepath.Join(netsocsDir, "update.sh")
	updateScript := fmt.Sprintf(`#!/bin/bash
set -e
echo "Updating CLI..."
sudo rm /usr/local/bin/netsocs
sudo cp "%s" /usr/local/bin/netsocs
sudo chmod +x /usr/local/bin/netsocs
echo "Update complete!"
`, newBinPath)
	if err := os.WriteFile(updateScriptPath, []byte(updateScript), 0755); err != nil {
		return err
	}

	pterm.Info.Println("Nuevo binario descargado en:", newBinPath)
	pterm.Info.Println("Ejecutando script de actualización en segundo plano:", updateScriptPath)

	// Ejecutar el script de forma asíncrona
	cmd := exec.Command("bash", updateScriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Iniciar el comando en segundo plano sin esperar
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error iniciando el script de actualización: %w", err)
	}

	// No esperar a que termine, solo informar que se inició
	pterm.Success.Println("Actualización iniciada en segundo plano. El CLI se actualizará automáticamente.")

	return nil
}

func isCurrentPlatformAsset(name string) bool {
	osName := runtime.GOOS
	arch := runtime.GOARCH
	return strings.Contains(name, osName) && strings.Contains(name, arch)
}

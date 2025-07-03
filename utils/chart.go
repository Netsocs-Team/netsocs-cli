package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pterm/pterm"
	"gopkg.in/yaml.v3"
)

func UpdateChartConfig(fieldPath string, value interface{}) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not get home directory: %w", err)
	}

	valuesPath := filepath.Join(homeDir, "netsocs", "values.yaml")
	pterm.Debug.Printfln("Buscando archivo en: %s", valuesPath)

	if _, err := os.Stat(valuesPath); os.IsNotExist(err) {
		return fmt.Errorf("the values.yaml file does not exist at %s", valuesPath)
	}

	yamlFile, err := os.ReadFile(valuesPath)
	if err != nil {
		return fmt.Errorf("error reading YAML file: %w", err)
	}

	var data map[string]interface{}
	if err := yaml.Unmarshal(yamlFile, &data); err != nil {
		return fmt.Errorf("error decoding YAML: %w", err)
	}

	pterm.Info.Printfln("Updating field '%s' with value: %v", fieldPath, value)

	fields := strings.Split(fieldPath, ".")
	updateNestedField(data, fields, value)

	updatedYaml, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("error generating YAML: %w", err)
	}

	if err := os.WriteFile(valuesPath, updatedYaml, 0644); err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	pterm.Success.Printfln("values.yaml file updated successfully")
	return nil
}

func updateNestedField(data map[string]interface{}, fields []string, value interface{}) {
	currentField := fields[0]

	if len(fields) == 1 {
		data[currentField] = value
		return
	}

	if _, exists := data[currentField]; !exists {
		data[currentField] = make(map[string]interface{})
	}

	if nestedMap, ok := data[currentField].(map[string]interface{}); ok {
		updateNestedField(nestedMap, fields[1:], value)
	} else {
		newMap := make(map[string]interface{})
		data[currentField] = newMap
		updateNestedField(newMap, fields[1:], value)
	}
}

func RunHelmUpgrade() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	valuesPath := filepath.Join(homeDir, "netsocs", "values.yaml")
	cmd := exec.Command("helm", "upgrade", "netsocs", "netsocs/netsocs-helm-chart", "--values", valuesPath)

	pterm.Info.Printfln("Ejecutando: %s", strings.Join(cmd.Args, " "))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func RunHelmUpgradeWithVersion(version string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	valuesPath := filepath.Join(homeDir, "netsocs", "values.yaml")
	args := []string{"upgrade", "netsocs", "netsocs/netsocs-helm-chart", "--values", valuesPath}
	if version != "" {
		args = append(args, "--version", version)
	}
	cmd := exec.Command("helm", args...)

	pterm.Info.Printfln("Running: %s", strings.Join(cmd.Args, " "))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func RunHelmRollback(revision string) error {
	args := []string{"rollback", "netsocs"}
	if revision != "" {
		args = append(args, revision)
	}
	cmd := exec.Command("helm", args...)

	pterm.Info.Printfln("Running: %s", strings.Join(cmd.Args, " "))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

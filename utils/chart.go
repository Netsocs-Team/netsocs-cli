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

// updateChartConfig actualiza un campo específico en values.yaml manteniendo el resto de la configuración
func UpdateChartConfig(fieldPath string, value interface{}) error {
	// Obtener directorio home
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("no se pudo obtener el directorio home: %w", err)
	}

	// Construir ruta al archivo values.yaml
	valuesPath := filepath.Join(homeDir, "netsocs", "values.yaml")
	pterm.Debug.Printfln("Buscando archivo en: %s", valuesPath)

	// Verificar si el archivo existe
	if _, err := os.Stat(valuesPath); os.IsNotExist(err) {
		return fmt.Errorf("el archivo values.yaml no existe en %s", valuesPath)
	}

	// Leer el archivo YAML existente
	yamlFile, err := os.ReadFile(valuesPath)
	if err != nil {
		return fmt.Errorf("error al leer el archivo YAML: %w", err)
	}

	// Decodificar el YAML a un map genérico
	var data map[string]interface{}
	if err := yaml.Unmarshal(yamlFile, &data); err != nil {
		return fmt.Errorf("error al decodificar YAML: %w", err)
	}

	pterm.Info.Printfln("Actualizando campo '%s' con valor: %v", fieldPath, value)

	// Actualizar el campo específico usando notación de puntos (field.path.nested)
	fields := strings.Split(fieldPath, ".")
	updateNestedField(data, fields, value)

	// Codificar de vuelta a YAML
	updatedYaml, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("error al generar YAML: %w", err)
	}

	// Escribir el archivo actualizado
	if err := os.WriteFile(valuesPath, updatedYaml, 0644); err != nil {
		return fmt.Errorf("error al escribir archivo: %w", err)
	}

	pterm.Success.Printfln("Archivo values.yaml actualizado correctamente")
	return nil
}

// updateNestedField actualiza campos anidados en un map
func updateNestedField(data map[string]interface{}, fields []string, value interface{}) {
	currentField := fields[0]

	// Si es el último campo, asignamos el valor
	if len(fields) == 1 {
		data[currentField] = value
		return
	}

	// Si el campo no existe o no es un map, lo creamos
	if _, exists := data[currentField]; !exists {
		data[currentField] = make(map[string]interface{})
	}

	// Convertimos el valor a map para seguir navegando
	if nestedMap, ok := data[currentField].(map[string]interface{}); ok {
		updateNestedField(nestedMap, fields[1:], value)
	} else {
		// Si el campo existe pero no es un map, lo reemplazamos
		newMap := make(map[string]interface{})
		data[currentField] = newMap
		updateNestedField(newMap, fields[1:], value)
	}
}

// Función para ejecutar Helm upgrade
func RunHelmUpgrade() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	chartPath := filepath.Join(homeDir, "netsocs")
	cmd := exec.Command("helm", "upgrade", "netsocs", chartPath)

	pterm.Info.Printfln("Ejecutando: %s", strings.Join(cmd.Args, " "))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

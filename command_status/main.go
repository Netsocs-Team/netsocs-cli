package commandstatus

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const NETSOCS_VERSION = "3.0.0"

func StatusHandler(cmd *cobra.Command, args []string) {
	verbose, _ := cmd.Flags().GetBool("verbose")

	// Mostrar versión
	pterm.DefaultHeader.
		WithBackgroundStyle(pterm.NewStyle(pterm.BgGreen)).
		WithTextStyle(pterm.NewStyle(pterm.FgBlack)).
		Println(" Estado de NETSOCS " + NETSOCS_VERSION)

	// Verificar pods
	pods, err := GetNetsocsPods()
	if err != nil {
		pterm.Error.Printfln("Error al verificar pods: %v", err)
		os.Exit(1)
	}

	// Analizar estado de los pods
	allHealthy := true
	var problemPods []string

	for _, pod := range pods {
		if !IsPodHealthy(pod) {
			allHealthy = false
			problemPods = append(problemPods, pod.Name)
		}
	}

	// Mostrar resumen
	if allHealthy {
		pterm.Success.Println("Todos los servicios de NETSOCS están operativos")
	} else {
		pterm.Error.Printfln("Problemas detectados en los siguientes pods: %s", strings.Join(problemPods, ", "))
	}

	// Mostrar detalles si es verbose o hay errores
	if verbose || !allHealthy {
		DisplayPodsStatus(pods)
	}
}

type PodStatus struct {
	Name     string
	Ready    string
	Status   string
	Restarts string
	Age      string
}

func GetNetsocsPods() ([]PodStatus, error) {
	// Ejecutar kubectl para obtener pods de NETSOCS
	cmd := exec.Command("kubectl", "get", "pods", "-o=wide")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("error al ejecutar kubectl: %w", err)
	}

	// Parsear salida
	lines := strings.Split(out.String(), "\n")
	var pods []PodStatus

	// Saltar la línea de encabezado
	for i, line := range lines {
		if i == 0 || strings.TrimSpace(line) == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 6 {
			pods = append(pods, PodStatus{
				Name:     fields[0],
				Ready:    fields[1],
				Status:   fields[2],
				Restarts: fields[3],
				Age:      fields[4],
			})
		}
	}

	return pods, nil
}

func IsPodHealthy(pod PodStatus) bool {
	// Verificar que el pod esté en estado Running y todas las containers estén ready
	readyParts := strings.Split(pod.Ready, "/")
	if len(readyParts) != 2 {
		return false
	}

	return readyParts[0] == readyParts[1] && pod.Status == "Running"
}

func DisplayPodsStatus(pods []PodStatus) {
	// Crear tabla para mostrar los pods
	tableData := pterm.TableData{
		{"Pod", "Ready", "Status", "Restarts", "Age"},
	}

	for _, pod := range pods {
		statusColor := pterm.FgGreen
		if !IsPodHealthy(pod) {
			statusColor = pterm.FgRed
		}

		tableData = append(tableData, []string{
			pod.Name,
			statusColor.Sprint(pod.Ready),
			statusColor.Sprint(pod.Status),
			statusColor.Sprint(pod.Restarts),
			pod.Age,
		})
	}

	pterm.DefaultSection.Println("Estado detallado de los pods")
	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
}

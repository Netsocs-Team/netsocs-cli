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

	// Show version
	pterm.DefaultHeader.
		WithBackgroundStyle(pterm.NewStyle(pterm.BgGreen)).
		WithTextStyle(pterm.NewStyle(pterm.FgBlack)).
		Println(" NETSOCS Status " + NETSOCS_VERSION)

	// Check pods
	pods, err := GetNetsocsPods()
	if err != nil {
		pterm.Error.Printfln("Error checking pods: %v", err)
		os.Exit(1)
	}

	// Analyze pod status
	allHealthy := true
	var problemPods []string

	for _, pod := range pods {
		if !IsPodHealthy(pod) {
			allHealthy = false
			problemPods = append(problemPods, pod.Name)
		}
	}

	// Show summary
	if allHealthy {
		pterm.Success.Println("All NETSOCS services are operational")
	} else {
		pterm.Error.Printfln("Problems detected in the following pods: %s", strings.Join(problemPods, ", "))
	}

	// Show details if verbose or there are errors
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
	// Run kubectl to get NETSOCS pods
	cmd := exec.Command("kubectl", "get", "pods", "-o=wide")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("error running kubectl: %w", err)
	}

	// Parse output
	lines := strings.Split(out.String(), "\n")
	var pods []PodStatus

	// Skip header line
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
	// Check that the pod is in Running state and all containers are ready
	readyParts := strings.Split(pod.Ready, "/")
	if len(readyParts) != 2 {
		return false
	}

	return readyParts[0] == readyParts[1] && pod.Status == "Running"
}

func DisplayPodsStatus(pods []PodStatus) {
	// Create table to display pods
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

	pterm.DefaultSection.Println("Detailed pod status")
	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
}

package utils

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pterm/pterm"
)

var urlsList = []string{
	"https://netsocs.com",
	"https://netsocs-team.github.io/netsocs-helm-chart",
	"https://ghcr.io",
	"https://plugins.traefik.io",
	"http://github.com/",
	"https://hub.docker.com/",
}

func CheckNetworkConnection() bool {
	// Configure pterm style
	pterm.Info.Println("🔍 Checking enviroment connectivity...")

	// Create a table to show results
	tableData := pterm.TableData{
		{"URL", "Status", "Response Time", "HTTP Code"},
	}

	// Configure HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create progress bar
	progress, _ := pterm.DefaultProgressbar.WithTotal(len(urlsList)).WithTitle("Checking URLs").Start()

	successCount := 0
	failedCount := 0

	for _, url := range urlsList {
		progress.UpdateTitle(fmt.Sprintf("Checking: %s", url))

		start := time.Now()
		resp, err := client.Get(url)
		duration := time.Since(start)

		var status, statusCode string

		if err != nil {
			status = "❌ Error"
			statusCode = "N/A"
			failedCount++
		} else {
			defer resp.Body.Close()
			statusCode = fmt.Sprintf("%d", resp.StatusCode)

			if resp.StatusCode >= 200 && resp.StatusCode < 400 {
				status = "✅ Connected"
				successCount++
			} else {
				status = "⚠️  HTTP Error"
				failedCount++
			}
		}

		// Agregar fila a la tabla
		tableData = append(tableData, []string{
			url,
			status,
			fmt.Sprintf("%.2fs", duration.Seconds()),
			statusCode,
		})

		progress.Increment()
	}

	progress.Stop()

	// Show summary
	pterm.Println()
	pterm.DefaultSection.Println("📊 Verification Results")

	// Show statistics
	stats := fmt.Sprintf("✅ Connected: %d | ❌ Failed: %d | 📊 Total: %d",
		successCount, failedCount, len(urlsList))
	pterm.Info.Println(stats)

	// Show table
	pterm.Println()
	_ = pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()

	// Show recommendations
	pterm.Println()
	if failedCount == 0 {
		pterm.Success.Println("🎉 Excellent! All URLs are accessible.")
		return true
	} else if failedCount < len(urlsList)/2 {
		pterm.Warning.Println("⚠️  Some URLs are not accessible. Check the enviroment network connection.")
		return false
	} else {
		pterm.Error.Println("🚨 Many URLs are not accessible. Possible network blocking detected.")
		return false
	}
}

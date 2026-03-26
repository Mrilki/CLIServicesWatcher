package reporter

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Mrilki/CLIServicesWatcher/internal/models"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type Report struct {
	Timestamp    string          `json:"timestamp"`
	TotalTargets int             `json:"total_targets"`
	TotalSuccess int             `json:"total_success"`
	TotalFail    int             `json:"total_fail"`
	AvgLatency   string          `json:"avg_latency"`
	Results      []models.Result `json:"results"`
}

func SaveToJSON(results []models.Result, fileName string) error {
	successCount := 0
	var totalLatency models.Duration

	for _, result := range results {
		if result.Success {
			successCount++
		}
		totalLatency += result.Latency
	}
	totalTargets := len(results)

	var avgLatency string
	if totalTargets > 0 {
		avgLatency = (totalLatency / models.Duration(totalTargets)).String()
	} else {
		avgLatency = "0ms"
	}

	report := Report{
		Timestamp:    time.Now().Format(time.RFC3339),
		TotalTargets: totalTargets,
		TotalSuccess: successCount,
		TotalFail:    totalTargets - successCount,
		AvgLatency:   avgLatency,
		Results:      results,
	}

	file, err := os.Create(fileName)

	if err != nil {
		return fmt.Errorf("error creating report file %s: %w", fileName, err)
	}

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	if err := enc.Encode(report); err != nil {
		return fmt.Errorf("error encoding report file %s: %w", fileName, err)
	}
	fmt.Printf("Report file saved to %s\n", fileName)
	if err := file.Close(); err != nil {
		return fmt.Errorf("error closing report file %s: %w", fileName, err)
	}
	return nil
}

func PrintStats(results []models.Result) {
	successCount := 0
	var totalLatency models.Duration
	totalTargets := len(results)

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleRounded)
	t.SetTitle("RESULTS")

	t.AppendHeader(table.Row{
		"Success",
		"Type",
		"Name",
		"Code",
		"Error",
		"Latency",
	})
	for _, result := range results {
		if result.Success {
			successCount++
		}
		totalLatency += result.Latency

		status := getStatusString(result.Success)

		t.AppendRow(table.Row{
			status,
			result.Type,
			result.Name,
			getStatusCodeString(result.StatusCode),
			result.Error,
			result.Latency,
		})
	}

	var avgLatency string
	if totalTargets > 0 {
		avgLatency = (totalLatency / models.Duration(totalTargets)).String()
	} else {
		avgLatency = "0ms"
	}

	t.AppendFooter(table.Row{
		fmt.Sprintf("Total success: %d/%d", successCount, totalTargets),
		"",
		"",
		"",
		"",
		avgLatency,
	})
	t.SetColumnConfigs([]table.ColumnConfig{
		{Name: "Success", WidthMax: 14},
		{Name: "Name", WidthMax: 60, WidthMaxEnforcer: text.WrapText},
		{Name: "Error", WidthMax: 90, WidthMaxEnforcer: text.WrapText},
		{Name: "Type", WidthMax: 4},
		{Name: "Code", WidthMax: 4},
		{Name: "Latency", WidthMax: 10},
	})

	fmt.Println()
	t.Render()
	fmt.Println()
}

func getStatusString(success bool) string {
	if success {
		return color.New(color.FgGreen).SprintFunc()("OK")
	}
	return color.New(color.FgRed).SprintFunc()("FAIL")
}

func getStatusCodeString(statusCode *int) string {
	if statusCode == nil {
		return "N/A"
	}
	return fmt.Sprintf("%d", *statusCode)
}

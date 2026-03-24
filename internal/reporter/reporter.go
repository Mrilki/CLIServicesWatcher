package reporter

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Mrilki/CLIServicesWatcher/internal/models"
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

	report := Report{
		Timestamp:    time.Now().Format(time.RFC3339),
		TotalTargets: totalTargets,
		TotalSuccess: successCount,
		TotalFail:    totalTargets - successCount,
		AvgLatency:   fmt.Sprintf("%.2f", float32(totalLatency)/float32(totalTargets)),
		Results:      results,
	}

	file, err := os.Create(fileName)

	if err != nil {
		return fmt.Errorf("error creating report file %s: %w", fileName, err)
	}

	defer file.Close()
	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	if err := enc.Encode(report); err != nil {
		return fmt.Errorf("error encoding report file %s: %w", fileName, err)
	}
	fmt.Printf("Report file saved to %s\n", fileName)
	return nil
}

func PrintStats(results []models.Result) {
	successCount := 0
	for _, result := range results {
		if result.Success {
			successCount++
		}
		fmt.Println(result)
	}
	fmt.Printf("Total success: %d/%d", successCount, len(results))
}

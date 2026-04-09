package reporter

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/Mrilki/CLIServicesWatcher/internal/models"
	"github.com/Mrilki/CLIServicesWatcher/internal/testutil"
)

func TestSaveToJSON(t *testing.T) {
	tmp, err := os.CreateTemp("", "report-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmpPath := tmp.Name()
	tmp.Close()              //nolint:errcheck
	defer os.Remove(tmpPath) //nolint:errcheck

	results := []models.Result{
		{
			Name:       "Google",
			Address:    "http://google.com",
			StatusCode: testutil.IntPtr(200),
			Type:       models.CheckTypeHTTP,
			Latency:    models.Duration(50 * time.Millisecond),
			Success:    true,
		},
		{
			Name:    "Broken",
			Address: "http://broken.com",
			Type:    models.CheckTypeHTTP,
			Success: false,
			Error:   "connection refused",
		},
	}

	err = SaveToJSON(results, tmpPath, testutil.DiscardLogger())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if _, err := os.Stat(tmpPath); os.IsNotExist(err) {
		t.Fatal("expected file to exist")
	}

	data, err := os.ReadFile(tmpPath)
	if err != nil {
		t.Fatalf("failed to read file %v", err)
	}

	var report Report
	err = json.Unmarshal(data, &report)
	if err != nil {
		t.Errorf("expected valid json, got %v", err)
	}

	if report.TotalTargets != 2 {
		t.Errorf("expected TotalTargets=2, got %d", report.TotalTargets)
	}

	if report.TotalSuccess != 1 {
		t.Errorf("expected TotalSuccess=1, got %d", report.TotalSuccess)
	}
	if report.TotalFail != 1 {
		t.Errorf("expected TotalFail=1, got %d", report.TotalFail)
	}
}

func TestSaveToJSON_InvalidPath(t *testing.T) {
	invalidPath := "/nonexistent/report.json"
	results := []models.Result{
		{Name: "test", Success: true},
	}

	err := SaveToJSON(results, invalidPath, testutil.DiscardLogger())

	if err == nil {
		t.Error("expected error invalid path")
	}
}

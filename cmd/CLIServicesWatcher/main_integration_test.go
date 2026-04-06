//go:build integration

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func createTempConfig(t *testing.T, content string) string {
	t.Helper()

	tmp, err := os.CreateTemp("", "test-config-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	_, err = tmp.WriteString(content)
	if err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}
	tmp.Close()

	return tmp.Name()
}

func runCLI(t *testing.T, args ...string) (string, error) {
	t.Helper()

	_, filename, _, _ := runtime.Caller(0)
	testDir := filepath.Dir(filename)
	projectRoot := filepath.Join(testDir, "..", "..")
	binaryName := "watcher"
	if runtime.GOOS == "windows" {
		binaryName = "watcher.exe"
	}
	binary := filepath.Join(projectRoot, binaryName)

	if _, err := os.Stat(binary); os.IsNotExist(err) {
		return "", fmt.Errorf("binary not found at %s (cwd: %s)", binary, filepath.Join(testDir))
	}

	cmd := exec.Command(binary, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func TestCLI_Run_Success(t *testing.T) {
	if testing.Short() {
		t.Skip("skipped integration test in short mode")
	}

	configContent := `{
		"timeout": 2,
		"targets": [
			{
				"name": "TestHTTP",
				"address": "http://httpbin.org/status/200",
				"type": "http"
			}
		]
	}`

	configPath := createTempConfig(t, configContent)
	defer os.Remove(configPath)

	reportPath := filepath.Join(t.TempDir(), "report.json")

	output, err := runCLI(t,
		"--config", configPath,
		"--output", reportPath,
	)
	if err != nil {
		t.Fatalf("CLI failed: %v\nOutput: %s", err, output)
	}

	if _, err := os.Stat(reportPath); os.IsNotExist(err) {
		t.Fatal("expected report.json to be create")
	}

	data, err := os.ReadFile(reportPath)
	if err != nil {
		t.Fatalf("failed to read report file: %v", err)
	}

	var report map[string]interface{}
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatalf("expected valid JSON, got error: %v", err)
	}

	if report["total_targets"] != 1.0 {
		t.Errorf("expected total_targets=1, got %v", report["total_targets"])
	}

	if results, ok := report["results"].([]interface{}); ok {
		if len(results) != 1 {
			t.Errorf("expected 1 result, got %d", len(results))
		}
	} else {
		t.Errorf("expected results to be an array")
	}
}

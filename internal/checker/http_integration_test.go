//go:build integration

package checker

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Mrilki/CLIServicesWatcher/internal/models"
)

func TestHTTPChecker_Integration_Success(t *testing.T) {
	if testing.Short() {
		t.Skip("skipped integration test in short mode")
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	checker := NewHTTPChecker(5 * time.Second)

	result := checker.Check(context.Background(), models.Target{
		Name:    "TestServer",
		Address: server.URL,
		Type:    models.CheckTypeHTTP,
	})

	if !result.Success {
		t.Errorf("expected success, got error: %v", result.Error)
	}

	if result.StatusCode == nil {
		t.Error("expected status code, got nil")
	} else if *result.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, *result.StatusCode)
	}
	if result.Latency <= 0 {
		t.Error("expected positive latency")
	}

	if result.Error != "" {
		t.Errorf("expected empty error, got %q", result.Error)
	}

}

func TestHTTPChecker_Integration_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipped integration test in short mode")
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	checker := NewHTTPChecker(5 * time.Second)

	result := checker.Check(context.Background(), models.Target{
		Name:    "Test404",
		Address: server.URL,
		Type:    models.CheckTypeHTTP,
	})

	if result.Success {
		t.Errorf("expected success=false, got true")
	}

	if result.StatusCode == nil {
		t.Error("expected status code, got nil")
	} else if *result.StatusCode != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, *result.StatusCode)
	}
}

func TestHTTPChecker_Integration_Timeout(t *testing.T) {
	if testing.Short() {
		t.Skip("skipped integration test in short mode")
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	defer server.Close()

	checker := NewHTTPChecker(100 * time.Millisecond)

	result := checker.Check(context.Background(), models.Target{
		Name:    "TestTimeout",
		Address: server.URL,
		Type:    models.CheckTypeHTTP,
	})

	if result.Success {
		t.Error("expected failure due to timeout, got success")
	}

	if result.Error == "" {
		t.Error("expected error message, got empty")
	}
}

func TestHTTPChecker_Integration_Unreachable(t *testing.T) {
	if testing.Short() {
		t.Skip("skipped integration test in short mode")
	}

	unreachableURL := "http://192.0.2.1:9999/nonexistent"

	checker := NewHTTPChecker(1 * time.Second)

	result := checker.Check(context.Background(), models.Target{
		Name:    "TestUnreachable",
		Address: unreachableURL,
		Type:    models.CheckTypeHTTP,
	})

	if result.Success {
		t.Error("expected failure for unreachable host, got success")
	}
	if result.Error == "" {
		t.Error("expected error message, got empty")
	}
}

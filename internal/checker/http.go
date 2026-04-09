package checker

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Mrilki/CLIServicesWatcher/internal/models"
)

type HTTPChecker struct {
	Client        *http.Client
	GlobalTimeout time.Duration
}

func NewHTTPChecker(globalTimeout time.Duration) *HTTPChecker {
	return &HTTPChecker{
		Client: &http.Client{
			Timeout: globalTimeout,
		},
		GlobalTimeout: globalTimeout,
	}
}

func (c *HTTPChecker) Check(ctx context.Context, target models.Target) models.Result {
	result := models.Result{
		Name:    target.Name,
		Address: target.Address,
		Success: false,
		Type:    target.Type,
	}

	timeout := target.GetTimeoutDuration(c.GlobalTimeout)

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	start := time.Now()
	req, err := http.NewRequestWithContext(ctx, "GET", target.Address, nil)
	if err != nil {
		result.Error = fmt.Sprintf("failed to create request: %v", err)
		result.Success = false
		result.SetLatency(0)
		return result
	}

	resp, err := c.Client.Do(req)
	latency := time.Since(start)
	result.SetLatency(latency)

	if err != nil {
		result.Error = classifyError(err, timeout)
		return result
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		result.Success = true
	}
	result.StatusCode = &resp.StatusCode
	return result
}

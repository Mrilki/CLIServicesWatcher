package checker

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Mrilki/CLIServicesWatcher/internal/models"
)

type Checker interface {
	Check(target models.Target) models.Result
}

type HTTPChecker struct {
	Client        *http.Client
	GlobalTimeout time.Duration
}

func NewHttpChecker(GlobalTimeout time.Duration) *HTTPChecker {
	return &HTTPChecker{
		Client:        &http.Client{},
		GlobalTimeout: GlobalTimeout,
	}
}

func (c *HTTPChecker) Check(target models.Target) models.Result {
	result := models.Result{
		Name:    target.Name,
		URL:     target.URL,
		Success: false,
	}

	timeout := target.GetTimeoutDuration(c.GlobalTimeout)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", target.URL, nil)
	if err != nil {
		result.Error = fmt.Sprintf("failed to create request: %v", err)
		result.Success = false
		result.Latency = time.Duration(0)
		return result
	}
	start := time.Now()
	resp, err := c.Client.Do(req)
	latency := time.Since(start)
	result.Latency = latency

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			result.Error = fmt.Sprintf("timeout after %d seconds", timeout)
		} else {
			result.Error = err.Error()
		}

	} else {
		result.StatusCode = resp.StatusCode
		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			result.Success = true
		}
	}
	return result
}

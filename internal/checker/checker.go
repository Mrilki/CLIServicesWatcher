package checker

import (
	"net/http"
	"time"

	"github.com/Mrilki/CLIServicesWatcher/internal/models"
)

type Checker interface {
	Check(target models.Target) models.Result
}

type HTTPChecker struct {
	Client *http.Client
}

func NewHttpChecker(timeout time.Duration) *HTTPChecker {
	return &HTTPChecker{
		Client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *HTTPChecker) Check(target models.Target) models.Result {
	result := models.Result{
		Name:    target.Name,
		URL:     target.URL,
		Success: false,
	}

	start := time.Now()

	resp, err := c.Client.Get(target.URL)
	latency := time.Since(start)
	result.Latency = latency

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		result.Error = err.Error()
	} else {
		result.StatusCode = resp.StatusCode
		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			result.Success = true
		}
	}
	return result
}

package checker

import (
	"context"
	"net"
	"time"

	"github.com/Mrilki/CLIServicesWatcher/internal/models"
)

type TCPChecker struct {
	Dialer        *net.Dialer
	GlobalTimeout time.Duration
}

func NewTCPChecker(globalTimeout time.Duration) *TCPChecker {
	return &TCPChecker{
		Dialer:        &net.Dialer{},
		GlobalTimeout: globalTimeout,
	}
}

func (c *TCPChecker) Check(ctx context.Context, target models.Target) models.Result {
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
	conn, err := c.Dialer.DialContext(ctx, "tcp", target.Address)
	latency := time.Since(start)
	result.SetLatency(latency)

	if err != nil {
		result.Error = classifyError(err, timeout)
		return result
	}
	defer conn.Close()
	result.Success = true
	return result
}

package checker

import (
	"context"
	"net"
	"time"

	"github.com/Mrilki/CLIServicesWatcher/internal/models"
)

type DNSChecker struct {
	Resolver      *net.Resolver
	GlobalTimeout time.Duration
}

func NewDNSChecker(globalTimeout time.Duration) *DNSChecker {
	return &DNSChecker{
		Resolver:      &net.Resolver{},
		GlobalTimeout: globalTimeout,
	}
}

func (c *DNSChecker) Check(ctx context.Context, target models.Target) models.Result {
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
	_, err := c.Resolver.LookupHost(ctx, target.Address)
	latency := time.Since(start)
	result.SetLatency(latency)
	if err != nil {
		result.Error = classifyError(err, timeout)
	} else {
		result.Success = true
	}
	return result

}

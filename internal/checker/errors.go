package checker

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"
)

var (
	ErrUnknownType = errors.New("checker: unknown check type")
	ErrTimeout     = errors.New("checker: timeout")
	ErrNetwork     = errors.New("checker: network error")
)

func classifyError(err error, timeout time.Duration) string {
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return fmt.Sprintf("%s: timeout after %v", ErrTimeout, timeout)
	}

	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return fmt.Sprintf("%s: timeout after %v", ErrTimeout, timeout)
	}

	return fmt.Sprintf("%s: %v", ErrNetwork, err)
}

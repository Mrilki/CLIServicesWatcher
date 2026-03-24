package checker

import (
	"fmt"
	"time"

	"github.com/Mrilki/CLIServicesWatcher/internal/models"
)

type Factory struct {
	GlobalTimeout time.Duration
}

func NewCheckerFactory(globalTimeout time.Duration) *Factory {
	return &Factory{
		GlobalTimeout: globalTimeout,
	}
}

func (f *Factory) New(checkType models.CheckType) (Checker, error) {
	switch checkType {
	case models.CheckTypeTCP:
		return NewTCPChecker(f.GlobalTimeout), nil
	case models.CheckTypeDNS:
		return NewDNSChecker(f.GlobalTimeout), nil
	case models.CheckTypeHTTP:
		return NewHTTPChecker(f.GlobalTimeout), nil
	default:
		return nil, fmt.Errorf("%w: got=%s", ErrUnknownType, checkType)
	}

}

package checker

import (
	"context"

	"github.com/Mrilki/CLIServicesWatcher/internal/models"
)

type Checker interface {
	Check(ctx context.Context, target models.Target) models.Result
}

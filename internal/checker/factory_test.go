package checker

import (
	"errors"
	"testing"
	"time"

	"github.com/Mrilki/CLIServicesWatcher/internal/models"
)

func TestFactory_New(t *testing.T) {
	factory := Factory{GlobalTimeout: 5 * time.Second}

	tests := []struct {
		name        string
		checkType   models.CheckType
		expectError bool
		errType     error
	}{
		{
			name:        "http checker",
			checkType:   models.CheckTypeHTTP,
			expectError: false,
			errType:     nil,
		},
		{
			name:        "TCP checker",
			checkType:   models.CheckTypeTCP,
			expectError: false,
			errType:     nil,
		},
		{
			name:        "DNS checker",
			checkType:   models.CheckTypeDNS,
			expectError: false,
			errType:     nil,
		},
		{
			name:        "invalid type",
			checkType:   "invalid",
			expectError: true,
			errType:     ErrUnknownType,
		},
		{
			name:        "empty type",
			checkType:   "",
			expectError: true,
			errType:     ErrUnknownType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chkr, err := factory.New(tt.checkType)

			if tt.expectError && err == nil {
				t.Error("expect error, got nil")
			}

			if !tt.expectError && err != nil {
				t.Errorf("expect no error, got %v", err)
			}

			if !tt.expectError && chkr == nil {
				t.Error("expect checker, got nil")
			}
			if tt.errType != nil && !errors.Is(err, tt.errType) {
				t.Errorf("expected error %v, got %v", tt.errType, err)
			}
		})
	}
}

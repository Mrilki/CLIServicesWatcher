package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/Mrilki/CLIServicesWatcher/internal/testutil"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		expectError bool
		errType     error
	}{
		{
			name:        "file not found",
			path:        "nonexistent.json",
			expectError: true,
			errType:     ErrNotFound,
		},
		{
			name:        "read error (directory)",
			path:        os.TempDir(),
			expectError: true,
			errType:     ErrRead,
		},
		{
			name:        "valid config",
			path:        filepath.Join("testdata", "valid-config.json"),
			expectError: false,
			errType:     nil,
		},
		{
			name:        "invalid config",
			path:        filepath.Join("testdata", "invalid-config.json"),
			expectError: true,
			errType:     ErrParse,
		},
		{
			name:        "empty targets",
			path:        filepath.Join("testdata", "empty-targets.json"),
			expectError: true,
			errType:     ErrValidate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := Load(tt.path, testutil.DiscardLogger())

			if tt.expectError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if !tt.expectError && cfg == nil {
				t.Error("expected config, got nil")
			}
			if tt.errType != nil && !errors.Is(err, tt.errType) {
				t.Errorf("expected error %v, got %v", tt.errType, err)
			}
		})
	}
}

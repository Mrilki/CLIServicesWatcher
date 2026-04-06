package models

import (
	"testing"
	"time"

	"github.com/Mrilki/CLIServicesWatcher/internal/testutil"
)

func TestCheckType_IsValid(t *testing.T) {
	tests := []struct {
		name      string
		checkType CheckType
		want      bool
	}{
		{name: "valid http", checkType: "http", want: true},
		{name: "valid tcp", checkType: "tcp", want: true},
		{name: "valid dns", checkType: "dns", want: true},
		{name: "invalid empty", checkType: "", want: false},
		{name: "invalid unknown", checkType: "udp", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.checkType.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectError bool
	}{
		{
			name: "valid config",
			config: &Config{
				Targets: []Target{
					{Name: "Google", Address: "google.com", Type: CheckTypeHTTP},
				},
				Timeout: 10,
			},
			expectError: false,
		},
		{
			name:        "empty target list",
			config:      &Config{},
			expectError: true,
		},
		{
			name: "empty address",
			config: &Config{
				Targets: []Target{
					{Name: "Google", Address: "", Type: CheckTypeHTTP},
				},
				Timeout: 10,
			},
			expectError: true,
		},
		{
			name: "invalid type",
			config: &Config{
				Targets: []Target{
					{Name: "Google", Address: "google.com", Type: "fdf"},
				},
				Timeout: 10,
			},
			expectError: true,
		},
		{
			name: "empty name auto-filled",
			config: &Config{
				Targets: []Target{
					{Name: "", Address: "google.com", Type: CheckTypeHTTP},
				},
				Timeout: 10,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if err == nil && tt.expectError {
				t.Error("expected error, got nil")
			}
			if err != nil && !tt.expectError {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestTarget_GetTimeoutDuration(t *testing.T) {
	globalTimeout := 10 * time.Second
	tests := []struct {
		name   string
		target *Target
		want   time.Duration
	}{
		{
			name: "nil timeout uses global",
			target: &Target{
				Name:    "Test",
				Address: "http://test.com",
				Type:    CheckTypeHTTP,
				Timeout: nil,
			},
			want: 10 * time.Second,
		},
		{
			name: "zero timeout uses global",
			target: &Target{
				Name:    "Test",
				Address: "http://test.com",
				Type:    CheckTypeHTTP,
				Timeout: testutil.IntPtr(0),
			},
			want: globalTimeout,
		},
		{
			name: "negative timeout uses global",
			target: &Target{
				Name:    "Test",
				Address: "http://test.com",
				Type:    CheckTypeHTTP,
				Timeout: testutil.IntPtr(-2),
			},
			want: globalTimeout,
		},
		{
			name: "custom timeout overrides global",
			target: &Target{
				Name:    "Test",
				Address: "http://test.com",
				Type:    CheckTypeHTTP,
				Timeout: testutil.IntPtr(5),
			},
			want: 5 * time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.target.GetTimeoutDuration(globalTimeout)
			if got != tt.want {
				t.Errorf("GetTimeoutDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

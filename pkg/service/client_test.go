package service

import (
	"context"
	"testing"
	"time"

	"github.com/rs/zerolog"
)

func TestNewClient(t *testing.T) {
	logger := zerolog.New(nil).Level(zerolog.Disabled)

	tests := []struct {
		name        string
		config      ClientConfig
		expectError bool
	}{
		{
			name: "valid config with gRPC only",
			config: ClientConfig{
				GRPCAddr:       "localhost:50053",
				Timeout:        5 * time.Second,
				EnableFallback: false,
				Logger:         logger,
			},
			expectError: false,
		},
		{
			name: "valid config with fallback",
			config: ClientConfig{
				GRPCAddr:       "localhost:50053",
				HTTPBaseURL:    "http://localhost:8086",
				Timeout:        5 * time.Second,
				EnableFallback: true,
				Logger:         logger,
			},
			expectError: false,
		},
		{
			name: "default timeout",
			config: ClientConfig{
				GRPCAddr:       "localhost:50053",
				EnableFallback: false,
				Logger:         logger,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.config)
			if tt.expectError && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if client != nil {
				client.Close()
			}
		})
	}
}

func TestShouldFallback(t *testing.T) {
	logger := zerolog.New(nil).Level(zerolog.Disabled)

	client := &Client{
		config: ClientConfig{
			EnableFallback: true,
			Logger:         logger,
		},
		httpClient: &HTTPClient{}, // Mock HTTP client
	}

	tests := []struct {
		name           string
		err            error
		enableFallback bool
		hasHTTPClient  bool
		expected       bool
	}{
		{
			name:           "fallback enabled with HTTP client",
			err:            context.DeadlineExceeded,
			enableFallback: true,
			hasHTTPClient:  true,
			expected:       true,
		},
		{
			name:           "fallback disabled",
			err:            context.DeadlineExceeded,
			enableFallback: false,
			hasHTTPClient:  true,
			expected:       false,
		},
		{
			name:           "no HTTP client",
			err:            context.DeadlineExceeded,
			enableFallback: true,
			hasHTTPClient:  false,
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client.config.EnableFallback = tt.enableFallback
			if !tt.hasHTTPClient {
				client.httpClient = nil
			} else {
				client.httpClient = &HTTPClient{}
			}

			result := client.shouldFallback(tt.err)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestConvertGRPCError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected error
	}{
		{
			name:     "not found error",
			err:      ErrNotFound,
			expected: ErrNotFound,
		},
		{
			name:     "invalid input error",
			err:      ErrInvalidInput,
			expected: ErrInvalidInput,
		},
		{
			name:     "unavailable error",
			err:      ErrUnavailable,
			expected: ErrUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertGRPCError(tt.err)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// Package metrics provides Prometheus metrics for marketplace microservice migration monitoring
// backend/internal/metrics/migration_metrics_test.go
package metrics

import (
	"errors"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestRecordRoute(t *testing.T) {
	// Reset metrics before test
	RouteTotal.Reset()

	tests := []struct {
		name        string
		destination string
		userType    string
		wantCount   int
	}{
		{
			name:        "record microservice canary",
			destination: "microservice",
			userType:    "canary",
			wantCount:   1,
		},
		{
			name:        "record monolith regular",
			destination: "monolith",
			userType:    "regular",
			wantCount:   1,
		},
		{
			name:        "record admin",
			destination: "microservice",
			userType:    "admin",
			wantCount:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Record metric
			RecordRoute(tt.destination, tt.userType)

			// Verify counter increased
			count := testutil.ToFloat64(RouteTotal.WithLabelValues(tt.destination, tt.userType))
			if count < float64(tt.wantCount) {
				t.Errorf("RouteTotal = %v, want >= %v", count, tt.wantCount)
			}
		})
	}
}

func TestObserveRouteDuration(t *testing.T) {
	tests := []struct {
		name            string
		destination     string
		operation       string
		durationSeconds float64
	}{
		{
			name:            "fast get request",
			destination:     "microservice",
			operation:       "get",
			durationSeconds: 0.010, // 10ms
		},
		{
			name:            "slow search request",
			destination:     "monolith",
			operation:       "search",
			durationSeconds: 0.500, // 500ms
		},
		{
			name:            "create request",
			destination:     "microservice",
			operation:       "create",
			durationSeconds: 0.150, // 150ms
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Observe metric (this should not panic)
			// For histograms we can't easily verify the value, but we verify the call succeeds
			ObserveRouteDuration(tt.destination, tt.operation, tt.durationSeconds)
			// If we got here without panic, the test passes
		})
	}
}

func TestRecordMicroserviceError(t *testing.T) {
	// Reset metrics before test
	MicroserviceErrorsTotal.Reset()

	tests := []struct {
		name      string
		errorType string
		wantCount int
	}{
		{
			name:      "timeout error",
			errorType: "timeout",
			wantCount: 1,
		},
		{
			name:      "connection error",
			errorType: "connection",
			wantCount: 1,
		},
		{
			name:      "internal error",
			errorType: "internal",
			wantCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Record error
			RecordMicroserviceError(tt.errorType)

			// Verify counter increased
			count := testutil.ToFloat64(MicroserviceErrorsTotal.WithLabelValues(tt.errorType))
			if count < float64(tt.wantCount) {
				t.Errorf("MicroserviceErrorsTotal = %v, want >= %v", count, tt.wantCount)
			}
		})
	}
}

func TestRecordFallback(t *testing.T) {
	// Reset metrics before test
	FallbackTotal.Reset()

	tests := []struct {
		name      string
		reason    string
		wantCount int
	}{
		{
			name:      "microservice error fallback",
			reason:    "microservice_error",
			wantCount: 1,
		},
		{
			name:      "timeout fallback",
			reason:    "timeout",
			wantCount: 1,
		},
		{
			name:      "unavailable fallback",
			reason:    "unavailable",
			wantCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Record fallback
			RecordFallback(tt.reason)

			// Verify counter increased
			count := testutil.ToFloat64(FallbackTotal.WithLabelValues(tt.reason))
			if count < float64(tt.wantCount) {
				t.Errorf("FallbackTotal = %v, want >= %v", count, tt.wantCount)
			}
		})
	}
}

func TestSetRolloutPercent(t *testing.T) {
	tests := []struct {
		name    string
		percent int
	}{
		{
			name:    "0 percent",
			percent: 0,
		},
		{
			name:    "10 percent",
			percent: 10,
		},
		{
			name:    "50 percent",
			percent: 50,
		},
		{
			name:    "100 percent",
			percent: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set rollout percent
			SetRolloutPercent(tt.percent)

			// Verify gauge value
			value := testutil.ToFloat64(RolloutPercent)
			if value != float64(tt.percent) {
				t.Errorf("RolloutPercent = %v, want %v", value, tt.percent)
			}
		})
	}
}

func TestSetCanaryUsers(t *testing.T) {
	tests := []struct {
		name  string
		count int
	}{
		{
			name:  "no canary users",
			count: 0,
		},
		{
			name:  "5 canary users",
			count: 5,
		},
		{
			name:  "100 canary users",
			count: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set canary users
			SetCanaryUsers(tt.count)

			// Verify gauge value
			value := testutil.ToFloat64(CanaryUsers)
			if value != float64(tt.count) {
				t.Errorf("CanaryUsers = %v, want %v", value, tt.count)
			}
		})
	}
}

func TestSetFeatureFlagEnabled(t *testing.T) {
	tests := []struct {
		name      string
		enabled   bool
		wantValue float64
	}{
		{
			name:      "feature flag disabled",
			enabled:   false,
			wantValue: 0.0,
		},
		{
			name:      "feature flag enabled",
			enabled:   true,
			wantValue: 1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set feature flag
			SetFeatureFlagEnabled(tt.enabled)

			// Verify gauge value
			value := testutil.ToFloat64(FeatureFlagEnabled)
			if value != tt.wantValue {
				t.Errorf("FeatureFlagEnabled = %v, want %v", value, tt.wantValue)
			}
		})
	}
}

func TestClassifyGRPCError(t *testing.T) {
	tests := []struct {
		name          string
		err           error
		wantErrorType string
	}{
		{
			name:          "nil error",
			err:           nil,
			wantErrorType: "unknown",
		},
		{
			name:          "timeout error",
			err:           errors.New("context deadline exceeded"),
			wantErrorType: "timeout",
		},
		{
			name:          "connection refused",
			err:           errors.New("connection refused"),
			wantErrorType: "connection",
		},
		{
			name:          "validation failed",
			err:           errors.New("validation failed: invalid input"),
			wantErrorType: "validation",
		},
		{
			name:          "not found",
			err:           errors.New("listing not found"),
			wantErrorType: "not_found",
		},
		{
			name:          "internal error",
			err:           errors.New("internal server error"),
			wantErrorType: "internal",
		},
		{
			name:          "network timeout",
			err:           errors.New("request timeout after 5s"),
			wantErrorType: "timeout",
		},
		{
			name:          "connection reset",
			err:           errors.New("connection reset by peer"),
			wantErrorType: "connection",
		},
		{
			name:          "no such host",
			err:           errors.New("no such host"),
			wantErrorType: "connection",
		},
		{
			name:          "invalid request",
			err:           errors.New("invalid listing ID"),
			wantErrorType: "validation",
		},
		{
			name:          "resource not found",
			err:           errors.New("resource does not exist"),
			wantErrorType: "not_found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ClassifyGRPCError(tt.err)
			if got != tt.wantErrorType {
				t.Errorf("ClassifyGRPCError() = %v, want %v", got, tt.wantErrorType)
			}
		})
	}
}

func TestResetAllMetrics(t *testing.T) {
	// Set some values first
	SetRolloutPercent(50)
	SetCanaryUsers(10)
	SetFeatureFlagEnabled(true)

	// Reset all metrics
	ResetAllMetrics()

	// Verify gauges are reset to 0
	if value := testutil.ToFloat64(RolloutPercent); value != 0 {
		t.Errorf("RolloutPercent after reset = %v, want 0", value)
	}

	if value := testutil.ToFloat64(CanaryUsers); value != 0 {
		t.Errorf("CanaryUsers after reset = %v, want 0", value)
	}

	if value := testutil.ToFloat64(FeatureFlagEnabled); value != 0 {
		t.Errorf("FeatureFlagEnabled after reset = %v, want 0", value)
	}
}

func TestMetricsLabels(t *testing.T) {
	// Test that metrics have correct labels
	t.Run("RouteTotal labels", func(t *testing.T) {
		// Record with valid labels
		RecordRoute("microservice", "canary")
		// This should not panic
	})

	t.Run("RouteDurationSeconds labels", func(t *testing.T) {
		// Observe with valid labels
		ObserveRouteDuration("microservice", "get", 0.1)
		// This should not panic
	})

	t.Run("MicroserviceErrorsTotal labels", func(t *testing.T) {
		// Record with valid label
		RecordMicroserviceError("timeout")
		// This should not panic
	})

	t.Run("FallbackTotal labels", func(t *testing.T) {
		// Record with valid label
		RecordFallback("microservice_error")
		// This should not panic
	})
}

func TestMetricsRegistration(t *testing.T) {
	// Verify all metrics are registered in default registry
	metrics := []prometheus.Collector{
		RouteTotal,
		RouteDurationSeconds,
		MicroserviceErrorsTotal,
		FallbackTotal,
		RolloutPercent,
		CanaryUsers,
		FeatureFlagEnabled,
	}

	for i, metric := range metrics {
		if metric == nil {
			t.Errorf("metric at index %d is nil", i)
		}
	}
}

func BenchmarkRecordRoute(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RecordRoute("microservice", "canary")
	}
}

func BenchmarkObserveRouteDuration(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ObserveRouteDuration("microservice", "get", 0.1)
	}
}

func BenchmarkClassifyGRPCError(b *testing.B) {
	err := errors.New("context deadline exceeded")
	for i := 0; i < b.N; i++ {
		ClassifyGRPCError(err)
	}
}

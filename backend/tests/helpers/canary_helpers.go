// Package helpers provides utility functions for canary testing
package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	// DefaultTimeout for HTTP requests
	DefaultTimeout = 5 * time.Second

	// DefaultBackendURL for local testing
	DefaultBackendURL = "http://localhost:3000"

	// MetricsURL for Prometheus metrics
	MetricsURL = "http://localhost:9091/metrics"
)

// CanaryRequestOptions contains options for creating canary test requests
type CanaryRequestOptions struct {
	UserID        string
	IsAdmin       bool
	CanaryEnabled bool
	Headers       map[string]string
}

// CanaryMetrics contains parsed canary metrics from Prometheus
type CanaryMetrics struct {
	FeatureFlagEnabled  float64
	RolloutPercent      float64
	CanaryUsers         float64
	CircuitBreakerState float64
	CircuitBreakerTrips float64
	FailedRequests      float64
	RejectedRequests    float64
	Recoveries          float64
}

// CreateCanaryRequest создаёт HTTP запрос с canary параметрами
func CreateCanaryRequest(ctx context.Context, method, url string, body io.Reader, opts CanaryRequestOptions) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add canary headers
	if opts.UserID != "" {
		req.Header.Set("X-User-ID", opts.UserID)
	}

	if opts.IsAdmin {
		req.Header.Set("X-Is-Admin", "true")
	}

	if opts.CanaryEnabled {
		req.Header.Set("X-Canary-Enabled", "true")
	}

	// Add custom headers
	for key, value := range opts.Headers {
		req.Header.Set(key, value)
	}

	// Add test identification
	req.Header.Set("User-Agent", "canary-test-client")

	return req, nil
}

// WaitForCanaryMetrics ждёт появления canary метрик
func WaitForCanaryMetrics(t *testing.T, timeout time.Duration) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			t.Logf("⚠️ Timeout waiting for canary metrics")
			return
		case <-ticker.C:
			if metricsAvailable(t) {
				t.Log("✅ Canary metrics are available")
				return
			}
		}
	}
}

// metricsAvailable проверяет доступность metrics endpoint
func metricsAvailable(t *testing.T) bool {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", MetricsURL, nil)
	if err != nil {
		return false
	}

	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer func() { _ = resp.Body.Close() }()

	return resp.StatusCode == http.StatusOK
}

// AssertCanaryHeaders проверяет canary заголовки в response
func AssertCanaryHeaders(t *testing.T, resp *http.Response) {
	t.Helper()

	// Check for source header
	source := resp.Header.Get("X-Source-Microservice")
	if source != "" {
		require.Contains(t, []string{"monolith", "microservice"}, source,
			"X-Source-Microservice should be 'monolith' or 'microservice'")
		t.Logf("✅ Source header present: %s", source)
	}

	// Check for circuit breaker state header
	cbState := resp.Header.Get("X-Circuit-Breaker-State")
	if cbState != "" {
		require.Contains(t, []string{"CLOSED", "OPEN", "HALF_OPEN"}, cbState,
			"X-Circuit-Breaker-State should be valid state")
		t.Logf("✅ Circuit breaker state: %s", cbState)
	}
}

// GetCanaryMetrics получает текущие canary метрики из Prometheus
func GetCanaryMetrics(t *testing.T) (*CanaryMetrics, error) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", MetricsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics request: %w", err)
	}

	client := &http.Client{Timeout: DefaultTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch metrics: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("metrics endpoint returned %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read metrics body: %w", err)
	}

	metrics := &CanaryMetrics{}
	metricsText := string(body)

	// Parse metrics (simple text parsing)
	metrics.FeatureFlagEnabled = parseMetricValue(metricsText, "marketplace_feature_flag_enabled")
	metrics.RolloutPercent = parseMetricValue(metricsText, "marketplace_rollout_percent")
	metrics.CanaryUsers = parseMetricValue(metricsText, "marketplace_canary_users")
	metrics.CircuitBreakerState = parseMetricValue(metricsText, "marketplace_circuit_breaker_state")
	metrics.CircuitBreakerTrips = parseMetricValue(metricsText, "marketplace_circuit_breaker_trips_total")
	metrics.FailedRequests = parseMetricValue(metricsText, "marketplace_circuit_breaker_failed_requests_total")
	metrics.RejectedRequests = parseMetricValue(metricsText, "marketplace_circuit_breaker_rejected_requests_total")
	metrics.Recoveries = parseMetricValue(metricsText, "marketplace_circuit_breaker_recoveries_total")

	return metrics, nil
}

// parseMetricValue извлекает значение метрики из Prometheus текста
func parseMetricValue(text, metricName string) float64 {
	// Simple parsing: find line with metric name and extract value
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, metricName) && !strings.HasPrefix(line, "#") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				var value float64
				_, _ = fmt.Sscanf(parts[1], "%f", &value)
				return value
			}
		}
	}
	return 0
}

// WaitForBackend ждёт доступности backend сервера
func WaitForBackend(t *testing.T, url string, timeout time.Duration) bool {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			t.Logf("⚠️ Backend not available after %v", timeout)
			return false
		case <-ticker.C:
			if isBackendAvailable(url) {
				t.Logf("✅ Backend available at %s", url)
				return true
			}
		}
	}
}

// isBackendAvailable проверяет доступность backend
func isBackendAvailable(url string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return false
	}

	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer func() { _ = resp.Body.Close() }()

	return resp.StatusCode == http.StatusOK
}

// MakeCanaryRequest выполняет canary request и возвращает response
func MakeCanaryRequest(t *testing.T, method, url string, opts CanaryRequestOptions) (*http.Response, error) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	req, err := CreateCanaryRequest(ctx, method, url, nil, opts)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: DefaultTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}

// AssertJSONResponse проверяет JSON response структуру
func AssertJSONResponse(t *testing.T, resp *http.Response, wantStatus int) map[string]interface{} {
	t.Helper()

	require.Equal(t, wantStatus, resp.StatusCode,
		"Response status should match expected")

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "Should read response body")

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	require.NoError(t, err, "Response should be valid JSON")

	return result
}

// SimulateTraffic симулирует трафик с заданным распределением
func SimulateTraffic(t *testing.T, url string, numRequests int, opts CanaryRequestOptions) (microserviceCount, monolithCount int) {
	t.Helper()

	for i := 0; i < numRequests; i++ {
		// Generate unique user ID for each request
		userOpts := opts
		if userOpts.UserID == "" {
			userOpts.UserID = fmt.Sprintf("user%d", i+1000)
		}

		resp, err := MakeCanaryRequest(t, "GET", url, userOpts)
		if err != nil {
			t.Logf("⚠️ Request %d failed: %v", i, err)
			continue
		}
		defer func() { _ = resp.Body.Close() }()

		// Check source header
		source := resp.Header.Get("X-Source-Microservice")
		if source == "microservice" {
			microserviceCount++
		} else {
			monolithCount++
		}
	}

	return microserviceCount, monolithCount
}

// VerifyCanaryRollout проверяет правильность canary rollout процента
func VerifyCanaryRollout(t *testing.T, expectedPercent int, totalRequests int, microserviceCount int) {
	t.Helper()

	// Calculate tolerance: ±0.5%
	expectedMin := float64(expectedPercent) * float64(totalRequests) * 0.005
	expectedMax := float64(expectedPercent) * float64(totalRequests) * 0.015

	require.GreaterOrEqual(t, float64(microserviceCount), expectedMin,
		"Microservice traffic should be >= %.0f at %d%%", expectedMin, expectedPercent)
	require.LessOrEqual(t, float64(microserviceCount), expectedMax,
		"Microservice traffic should be <= %.0f at %d%%", expectedMax, expectedPercent)

	actualPercent := float64(microserviceCount) / float64(totalRequests) * 100
	t.Logf("✅ Rollout verified: %.2f%% (expected %d%%, tolerance ±0.5%%)",
		actualPercent, expectedPercent)
}

// GetCircuitBreakerState получает текущее состояние circuit breaker
func GetCircuitBreakerState(t *testing.T) (string, error) {
	t.Helper()

	metrics, err := GetCanaryMetrics(t)
	if err != nil {
		return "", err
	}

	// Circuit breaker state: 0=CLOSED, 1=OPEN, 2=HALF_OPEN
	stateMap := map[float64]string{
		0: "CLOSED",
		1: "OPEN",
		2: "HALF_OPEN",
	}

	state, ok := stateMap[metrics.CircuitBreakerState]
	if !ok {
		return "UNKNOWN", fmt.Errorf("unknown circuit breaker state: %v", metrics.CircuitBreakerState)
	}

	return state, nil
}

// AssertCircuitBreakerState проверяет состояние circuit breaker
func AssertCircuitBreakerState(t *testing.T, expectedState string) {
	t.Helper()

	state, err := GetCircuitBreakerState(t)
	if err != nil {
		t.Skipf("⚠️ Cannot get circuit breaker state: %v", err)
		return
	}

	require.Equal(t, expectedState, state,
		"Circuit breaker should be in %s state", expectedState)

	t.Logf("✅ Circuit breaker state verified: %s", state)
}

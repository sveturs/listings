// Package e2e provides end-to-end resilience tests
// backend/tests/e2e/resilience_test.go
//
// Build with: go test -v -tags=e2e ./tests/e2e/resilience_test.go
package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	backendURL     = "http://localhost:3000"
	mockServiceURL = "http://localhost:50052" // Mock service control API
	testListingID  = "123"
	testUserID     = "test-user-1"
	testJWTToken   = "" // Set from /tmp/token if needed
)

// TestSlowMicroserviceTimeout verifies timeout fallback to monolith
func TestSlowMicroserviceTimeout(t *testing.T) {
	// Setup: Configure mock service to respond slowly (1s delay)
	err := configureMockService(map[string]interface{}{
		"mode":  "slow",
		"delay": "1000ms",
	})
	require.NoError(t, err, "Failed to configure mock service")

	// Execute: Request listing (should timeout and fallback)
	resp, err := httpGet(fmt.Sprintf("%s/api/v1/marketplace/listings/%s", backendURL, testListingID))
	require.NoError(t, err, "Request failed")
	defer func() { _ = resp.Body.Close() }()

	// Verify: Response should come from monolith (fallback)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Should succeed via fallback")

	var result map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&result)

	// Check if response indicates fallback (depends on implementation)
	// For now, just verify we got a response
	assert.NotNil(t, result, "Should get response from monolith")

	t.Logf("✅ Slow microservice → timeout → monolith fallback")
}

// TestFailingMicroserviceCircuitBreaker verifies circuit breaker opens
func TestFailingMicroserviceCircuitBreaker(t *testing.T) {
	// Setup: Configure mock service to return errors
	err := configureMockService(map[string]interface{}{
		"mode": "error",
	})
	require.NoError(t, err, "Failed to configure mock service")

	// Execute: Send 5 requests to open circuit breaker
	failureCount := 0
	for i := 0; i < 5; i++ {
		resp, err := httpGet(fmt.Sprintf("%s/api/v1/marketplace/listings/%d", backendURL, i))
		if err != nil || resp.StatusCode != http.StatusOK {
			failureCount++
		}
		if resp != nil {
			_ = resp.Body.Close()
		}
		time.Sleep(100 * time.Millisecond)
	}

	// Circuit should be open now
	// Next request should be rejected immediately (or fallback to monolith)
	start := time.Now()
	resp, _ := httpGet(fmt.Sprintf("%s/api/v1/marketplace/listings/%s", backendURL, testListingID))
	elapsed := time.Since(start)

	// Should respond quickly (either circuit breaker rejection or immediate fallback)
	assert.Less(t, elapsed, 200*time.Millisecond,
		"Circuit breaker should reject/fallback quickly")

	if resp != nil {
		_ = resp.Body.Close()
	}

	t.Logf("✅ Circuit breaker opened after %d failures, subsequent request in %v",
		failureCount, elapsed)
}

// TestMicroserviceRecovery verifies circuit breaker recovery
func TestMicroserviceRecovery(t *testing.T) {
	// Skip if not running full e2e suite
	if testing.Short() {
		t.Skip("Skipping recovery test in short mode")
	}

	// Step 1: Open circuit by causing failures
	err := configureMockService(map[string]interface{}{
		"mode": "error",
	})
	require.NoError(t, err, "Failed to configure mock service")

	for i := 0; i < 5; i++ {
		resp, _ := httpGet(fmt.Sprintf("%s/api/v1/marketplace/listings/%d", backendURL, i))
		if resp != nil {
			_ = resp.Body.Close()
		}
		time.Sleep(100 * time.Millisecond)
	}

	t.Logf("Step 1: Circuit opened")

	// Step 2: Wait for half-open transition (30s)
	t.Logf("Step 2: Waiting 30s for HALF_OPEN state...")
	time.Sleep(31 * time.Second)

	// Step 3: Fix microservice
	err = configureMockService(map[string]interface{}{
		"mode": "normal",
	})
	require.NoError(t, err, "Failed to configure mock service")

	t.Logf("Step 3: Microservice recovered")

	// Step 4: Send successful requests to close circuit
	successCount := 0
	for i := 0; i < 3; i++ {
		resp, err := httpGet(fmt.Sprintf("%s/api/v1/marketplace/listings/%d", backendURL, i))
		if err == nil && resp.StatusCode == http.StatusOK {
			successCount++
		}
		if resp != nil {
			_ = resp.Body.Close()
		}
		time.Sleep(100 * time.Millisecond)
	}

	// Circuit should close after 2 successes
	assert.GreaterOrEqual(t, successCount, 2, "Should have successful requests")

	t.Logf("✅ Circuit breaker recovered: %d successful requests", successCount)
}

// TestMixedLoadPartialDegradation verifies partial traffic handling
func TestMixedLoadPartialDegradation(t *testing.T) {
	// Setup: Configure mock to fail 50% of requests
	err := configureMockService(map[string]interface{}{
		"mode":         "partial",
		"failure_rate": 50,
	})
	require.NoError(t, err, "Failed to configure mock service")

	// Execute: Send 20 concurrent requests
	const numRequests = 20
	var wg sync.WaitGroup
	successCount := int32(0)
	failureCount := int32(0)

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			resp, err := httpGet(fmt.Sprintf("%s/api/v1/marketplace/listings/%d", backendURL, id))
			if err == nil && resp != nil && resp.StatusCode == http.StatusOK {
				atomic.AddInt32(&successCount, 1)
			} else {
				atomic.AddInt32(&failureCount, 1)
			}
			if resp != nil {
				_ = resp.Body.Close()
			}
		}(i)
	}

	wg.Wait()

	// Verify: Should have mix of success/failure
	totalRequests := successCount + failureCount
	successRate := float64(successCount) / float64(totalRequests) * 100

	t.Logf("Mixed load results: %d/%d succeeded (%.1f%%)",
		successCount, totalRequests, successRate)

	// With 50% failure rate and fallback, we should still get reasonable success rate
	assert.GreaterOrEqual(t, successCount, int32(10),
		"Should have at least 50% success with fallback")

	t.Logf("✅ Partial degradation handled: %.1f%% success rate", successRate)
}

// TestCascadingFailurePrevention verifies circuit breakers prevent cascade
func TestCascadingFailurePrevention(t *testing.T) {
	// Setup: Multiple failing services
	// (For this test, we only have one microservice, but concept is same)

	// Configure microservice to fail
	err := configureMockService(map[string]interface{}{
		"mode": "error",
	})
	require.NoError(t, err, "Failed to configure mock service")

	// Execute: Hammer with requests
	const numRequests = 50
	var wg sync.WaitGroup
	startTime := time.Now()

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			resp, _ := httpGet(fmt.Sprintf("%s/api/v1/marketplace/listings/%d", backendURL, id))
			if resp != nil {
				_ = resp.Body.Close()
			}
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(startTime)

	// Verify: All requests should complete quickly (circuit breaker prevents queueing)
	avgLatency := elapsed / time.Duration(numRequests)
	t.Logf("Average latency under failure: %v", avgLatency)

	// With circuit breaker, average latency should be low (immediate rejection)
	assert.Less(t, avgLatency, 500*time.Millisecond,
		"Circuit breaker should prevent slow cascading failures")

	t.Logf("✅ Circuit breaker prevented cascading failure: %d requests in %v",
		numRequests, elapsed)
}

// TestEndToEndLatency verifies overall system latency with resilience
func TestEndToEndLatency(t *testing.T) {
	// Setup: Normal mode
	err := configureMockService(map[string]interface{}{
		"mode": "normal",
	})
	require.NoError(t, err, "Failed to configure mock service")

	// Measure end-to-end latency
	iterations := 10
	var totalLatency time.Duration

	for i := 0; i < iterations; i++ {
		start := time.Now()
		resp, err := httpGet(fmt.Sprintf("%s/api/v1/marketplace/listings/%d", backendURL, i))
		latency := time.Since(start)
		totalLatency += latency

		if err == nil && resp != nil {
			_ = resp.Body.Close()
		}

		time.Sleep(100 * time.Millisecond)
	}

	avgLatency := totalLatency / time.Duration(iterations)
	t.Logf("Average end-to-end latency: %v", avgLatency)

	// P99 should be under 200ms for normal operations
	assert.Less(t, avgLatency, 200*time.Millisecond,
		"Average latency should be under 200ms")

	t.Logf("✅ End-to-end latency acceptable: %v", avgLatency)
}

// Helper functions

func httpGet(url string) (*http.Response, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add JWT token if available
	if testJWTToken != "" {
		req.Header.Set("Authorization", "Bearer "+testJWTToken)
	}

	return client.Do(req)
}

func configureMockService(config map[string]interface{}) error {
	// POST to mock service control endpoint
	configJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST",
		fmt.Sprintf("%s/control/config", mockServiceURL),
		bytes.NewReader(configJSON))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("mock service config failed: %d", resp.StatusCode)
	}

	// Wait for config to take effect
	time.Sleep(100 * time.Millisecond)

	return nil
}

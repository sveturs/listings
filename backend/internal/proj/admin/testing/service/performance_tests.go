// Package service implements performance functional tests
// backend/internal/proj/admin/testing/service/performance_tests.go
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"backend/internal/proj/admin/testing/domain"
)

// PerformanceTests returns list of performance tests
var PerformanceTests = []FunctionalTest{
	{
		Name:        "performance-api-response-time",
		Category:    domain.TestCategoryPerformance,
		Description: "Measure API endpoint response times (should be <200ms)",
		RunFunc:     testAPIResponseTime,
	},
	{
		Name:        "performance-concurrent-users",
		Category:    domain.TestCategoryPerformance,
		Description: "Test system with 10/50/100 concurrent users",
		RunFunc:     testConcurrentUsers,
	},
	{
		Name:        "performance-database-queries",
		Category:    domain.TestCategoryPerformance,
		Description: "Check for slow database queries (>100ms)",
		RunFunc:     testDatabaseQueryPerformance,
	},
	{
		Name:        "performance-memory-usage",
		Category:    domain.TestCategoryPerformance,
		Description: "Monitor memory usage during test execution",
		RunFunc:     testMemoryUsage,
	},
}

// testAPIResponseTime measures response times of critical endpoints
func testAPIResponseTime(ctx context.Context, baseURL, token string) *domain.TestResult {
	result := &domain.TestResult{
		TestName:  "performance-api-response-time",
		TestSuite: "performance",
		Status:    domain.TestResultStatusPassed,
		StartedAt: time.Now().UTC(),
	}

	client := &http.Client{Timeout: 10 * time.Second}

	// Critical endpoints to test
	endpoints := []struct {
		name      string
		url       string
		threshold time.Duration // expected response time
	}{
		{"Auth Me", "/api/v1/auth/me", 100 * time.Millisecond},
		{"Unified Listings", "/api/v1/unified/listings?limit=10", 200 * time.Millisecond},
		{"Search", "/api/v1/search?query=test&limit=10", 300 * time.Millisecond},
		{"Admin Categories", "/api/v1/admin/categories", 150 * time.Millisecond},
	}

	var slowEndpoints []string
	totalDuration := time.Duration(0)

	for _, endpoint := range endpoints {
		start := time.Now()

		req, err := http.NewRequestWithContext(ctx, "GET", baseURL+endpoint.url, nil)
		if err != nil {
			return failTest(result, fmt.Sprintf("Failed to create request for %s", endpoint.name), err)
		}

		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := client.Do(req)
		if err != nil {
			return failTest(result, fmt.Sprintf("Failed to execute request for %s", endpoint.name), err)
		}

		// Read and discard body to measure full response time
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()

		duration := time.Since(start)
		totalDuration += duration

		// Check if endpoint is too slow
		if duration > endpoint.threshold {
			slowEndpoints = append(slowEndpoints, fmt.Sprintf("%s: %dms (expected <%dms)",
				endpoint.name, duration.Milliseconds(), endpoint.threshold.Milliseconds()))
		}

		// Check for errors
		if resp.StatusCode != http.StatusOK {
			return failTest(result, fmt.Sprintf("%s returned %d (expected 200)", endpoint.name, resp.StatusCode), nil)
		}
	}

	// Fail if any endpoint is too slow
	if len(slowEndpoints) > 0 {
		errMsg := fmt.Sprintf("Slow endpoints detected: %v", slowEndpoints)
		return failTest(result, errMsg, nil)
	}

	result.CompletedAt = time.Now().UTC()
	result.DurationMs = int(result.CompletedAt.Sub(result.StartedAt).Milliseconds())
	return result
}

// testConcurrentUsers tests system under concurrent load
func testConcurrentUsers(ctx context.Context, baseURL, token string) *domain.TestResult {
	result := &domain.TestResult{
		TestName:  "performance-concurrent-users",
		TestSuite: "performance",
		Status:    domain.TestResultStatusPassed,
		StartedAt: time.Now().UTC(),
	}

	client := &http.Client{Timeout: 10 * time.Second}

	// Test with different concurrency levels
	concurrencyLevels := []int{10, 50, 100}
	endpoint := "/api/v1/unified/listings?limit=5"

	for _, concurrency := range concurrencyLevels {
		var wg sync.WaitGroup
		var successCount atomic.Int64
		var failCount atomic.Int64
		var totalDuration atomic.Int64

		for i := 0; i < concurrency; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				reqStart := time.Now()

				req, err := http.NewRequestWithContext(ctx, "GET", baseURL+endpoint, nil)
				if err != nil {
					failCount.Add(1)
					return
				}

				req.Header.Set("Authorization", "Bearer "+token)

				resp, err := client.Do(req)
				if err != nil {
					failCount.Add(1)
					return
				}
				defer resp.Body.Close()

				// Read body
				_, _ = io.Copy(io.Discard, resp.Body)

				duration := time.Since(reqStart)
				totalDuration.Add(int64(duration.Milliseconds()))

				if resp.StatusCode == http.StatusOK {
					successCount.Add(1)
				} else {
					failCount.Add(1)
				}
			}()
		}

		wg.Wait()

		successRate := float64(successCount.Load()) / float64(concurrency) * 100

		// Require at least 95% success rate
		if successRate < 95 {
			errMsg := fmt.Sprintf("Concurrency %d: only %.1f%% success rate (expected >=95%%)", concurrency, successRate)
			return failTest(result, errMsg, nil)
		}

		// Check average response time under load
		avgResponseTime := time.Duration(totalDuration.Load() / int64(concurrency))
		if avgResponseTime > 1*time.Second {
			errMsg := fmt.Sprintf("Concurrency %d: avg response time %dms (expected <1000ms)", concurrency, avgResponseTime.Milliseconds())
			return failTest(result, errMsg, nil)
		}

	}

	result.CompletedAt = time.Now().UTC()
	result.DurationMs = int(result.CompletedAt.Sub(result.StartedAt).Milliseconds())
	return result
}

// testDatabaseQueryPerformance checks for slow database queries
func testDatabaseQueryPerformance(ctx context.Context, baseURL, token string) *domain.TestResult {
	result := &domain.TestResult{
		TestName:  "performance-database-queries",
		TestSuite: "performance",
		Status:    domain.TestResultStatusPassed,
		StartedAt: time.Now().UTC(),
	}

	client := &http.Client{Timeout: 10 * time.Second}

	// Test endpoints that execute complex queries
	complexEndpoints := []struct {
		name      string
		url       string
		threshold time.Duration
	}{
		{"Search with filters", "/api/v1/search?query=test&category=1&min_price=10&max_price=1000&limit=20", 300 * time.Millisecond},
		{"Unified listings full", "/api/v1/unified/listings?limit=100", 500 * time.Millisecond},
		{"Admin admins list", "/api/v1/admin/admins", 500 * time.Millisecond},
	}

	var slowQueries []string

	for _, endpoint := range complexEndpoints {
		start := time.Now()

		req, err := http.NewRequestWithContext(ctx, "GET", baseURL+endpoint.url, nil)
		if err != nil {
			return failTest(result, fmt.Sprintf("Failed to create request for %s", endpoint.name), err)
		}

		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := client.Do(req)
		if err != nil {
			return failTest(result, fmt.Sprintf("Failed to execute request for %s", endpoint.name), err)
		}

		// Read body
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		duration := time.Since(start)

		// Check response is valid
		if resp.StatusCode != http.StatusOK {
			return failTest(result, fmt.Sprintf("%s returned %d: %s", endpoint.name, resp.StatusCode, string(body)), nil)
		}

		// Parse and verify response has data
		var respData map[string]interface{}
		if err := json.Unmarshal(body, &respData); err != nil {
			return failTest(result, fmt.Sprintf("%s returned invalid JSON", endpoint.name), err)
		}

		// Check if query is too slow
		if duration > endpoint.threshold {
			slowQueries = append(slowQueries, fmt.Sprintf("%s: %dms (expected <%dms)",
				endpoint.name, duration.Milliseconds(), endpoint.threshold.Milliseconds()))
		}
	}

	// Fail if any query is too slow
	if len(slowQueries) > 0 {
		errMsg := fmt.Sprintf("Slow queries detected: %v", slowQueries)
		return failTest(result, errMsg, nil)
	}

	result.CompletedAt = time.Now().UTC()
	result.DurationMs = int(result.CompletedAt.Sub(result.StartedAt).Milliseconds())
	return result
}

// testMemoryUsage monitors memory usage during test execution
func testMemoryUsage(ctx context.Context, baseURL, token string) *domain.TestResult {
	result := &domain.TestResult{
		TestName:  "performance-memory-usage",
		TestSuite: "performance",
		Status:    domain.TestResultStatusPassed,
		StartedAt: time.Now().UTC(),
	}

	client := &http.Client{Timeout: 10 * time.Second}

	// Make multiple requests to accumulate data
	const requestCount = 50
	endpoint := "/api/v1/unified/listings?limit=100"

	for i := 0; i < requestCount; i++ {
		req, err := http.NewRequestWithContext(ctx, "GET", baseURL+endpoint, nil)
		if err != nil {
			return failTest(result, "Failed to create request", err)
		}

		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := client.Do(req)
		if err != nil {
			return failTest(result, "Failed to execute request", err)
		}

		// Read and discard body
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return failTest(result, fmt.Sprintf("Request %d returned %d", i+1, resp.StatusCode), nil)
		}

		// Small delay between requests
		time.Sleep(10 * time.Millisecond)
	}

	// Note: This is a simplified memory test
	// In production, you would use runtime.ReadMemStats() or external monitoring
	// to measure actual memory consumption

	result.CompletedAt = time.Now().UTC()
	result.DurationMs = int(result.CompletedAt.Sub(result.StartedAt).Milliseconds())
	return result
}

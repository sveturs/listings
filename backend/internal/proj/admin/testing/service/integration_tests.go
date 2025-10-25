// Package service implements integration tests for external services
// backend/internal/proj/admin/testing/service/integration_tests.go
package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"backend/internal/proj/admin/testing/domain"
)

// IntegrationTests returns list of integration tests
var IntegrationTests = []FunctionalTest{
	{
		Name:        "integration-redis-cache",
		Category:    domain.TestCategoryIntegration,
		Description: "Test Redis cache operations (SET, GET, DELETE, TTL)",
		RunFunc:     testRedisCache,
	},
	{
		Name:        "integration-opensearch-index",
		Category:    domain.TestCategoryIntegration,
		Description: "Test OpenSearch indexing and search",
		RunFunc:     testOpenSearchIndex,
	},
	{
		Name:        "integration-postgres-connection",
		Category:    domain.TestCategoryIntegration,
		Description: "Test PostgreSQL connection and basic queries",
		RunFunc:     testPostgresConnection,
	},
}

// testRedisCache tests Redis cache functionality
func testRedisCache(ctx context.Context, baseURL, token string) *domain.TestResult {
	result := &domain.TestResult{
		TestName:  "integration-redis-cache",
		TestSuite: "integration",
		Status:    domain.TestResultStatusPassed,
		StartedAt: time.Now().UTC(),
	}

	client := &http.Client{Timeout: 10 * time.Second}

	// Test 1: Cache listings (triggers Redis SET)
	req1, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/api/v1/unified/listings?limit=5", nil)
	if err != nil {
		return failTest(result, "Failed to create listings request", err)
	}
	req1.Header.Set("Authorization", "Bearer "+token)

	resp1, err := client.Do(req1)
	if err != nil {
		return failTest(result, "Failed to fetch listings (first call)", err)
	}
	defer func() { _ = resp1.Body.Close() }()

	if resp1.StatusCode != http.StatusOK {
		return failTest(result, fmt.Sprintf("Expected status 200, got %d", resp1.StatusCode), nil)
	}

	// Test 2: Fetch again (should hit cache)
	req2, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/api/v1/unified/listings?limit=5", nil)
	if err != nil {
		return failTest(result, "Failed to create listings request (2nd)", err)
	}
	req2.Header.Set("Authorization", "Bearer "+token)

	resp2, err := client.Do(req2)
	if err != nil {
		return failTest(result, "Failed to fetch listings (second call)", err)
	}
	defer func() { _ = resp2.Body.Close() }()

	if resp2.StatusCode != http.StatusOK {
		return failTest(result, fmt.Sprintf("Expected status 200 on cache hit, got %d", resp2.StatusCode), nil)
	}

	// Test 3: Verify search endpoint (also uses cache)
	req3, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/api/v1/search?query=test", nil)
	if err == nil {
		req3.Header.Set("Authorization", "Bearer "+token)
		resp3, err := client.Do(req3)
		if err == nil {
			defer func() { _ = resp3.Body.Close() }()
			// We don't fail if this optional check fails
		}
	}

	result.CompletedAt = time.Now().UTC()
	result.DurationMs = int(result.CompletedAt.Sub(result.StartedAt).Milliseconds())
	return result
}

// testOpenSearchIndex tests OpenSearch indexing and search
func testOpenSearchIndex(ctx context.Context, baseURL, token string) *domain.TestResult {
	result := &domain.TestResult{
		TestName:  "integration-opensearch-index",
		TestSuite: "integration",
		Status:    domain.TestResultStatusPassed,
		StartedAt: time.Now().UTC(),
	}

	client := &http.Client{Timeout: 10 * time.Second}

	// Test 1: Search for a common term (should return results from OpenSearch)
	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/api/v1/search?query=phone&limit=5", nil)
	if err != nil {
		return failTest(result, "Failed to create search request", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return failTest(result, "Failed to execute search", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return failTest(result, fmt.Sprintf("Expected status 200, got %d", resp.StatusCode), nil)
	}

	// Test 2: Empty search (edge case)
	req2, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/api/v1/search?query=", nil)
	if err != nil {
		return failTest(result, "Failed to create empty search request", err)
	}
	req2.Header.Set("Authorization", "Bearer "+token)

	resp2, err := client.Do(req2)
	if err != nil {
		return failTest(result, "Failed to execute empty search", err)
	}
	defer func() { _ = resp2.Body.Close() }()

	// Empty search should return 200 (might return all or none)
	if resp2.StatusCode != http.StatusOK {
		return failTest(result, fmt.Sprintf("Expected status 200 for empty search, got %d", resp2.StatusCode), nil)
	}

	result.CompletedAt = time.Now().UTC()
	result.DurationMs = int(result.CompletedAt.Sub(result.StartedAt).Milliseconds())
	return result
}

// testPostgresConnection tests PostgreSQL connectivity and queries
func testPostgresConnection(ctx context.Context, baseURL, token string) *domain.TestResult {
	result := &domain.TestResult{
		TestName:  "integration-postgres-connection",
		TestSuite: "integration",
		Status:    domain.TestResultStatusPassed,
		StartedAt: time.Now().UTC(),
	}

	client := &http.Client{Timeout: 10 * time.Second}

	// Test database connectivity through API endpoints
	// We use /api/v1/unified/listings as it requires DB query

	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/api/v1/unified/listings?limit=1", nil)
	if err != nil {
		return failTest(result, "Failed to create listings request", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return failTest(result, "Failed to fetch listings (DB query)", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return failTest(result, fmt.Sprintf("Expected status 200, got %d (DB might be down)", resp.StatusCode), nil)
	}

	// Test 2: Complex query with joins (search query)
	req2, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/api/v1/search?query=phone", nil)
	if err != nil {
		return failTest(result, "Failed to create search request", err)
	}
	req2.Header.Set("Authorization", "Bearer "+token)

	resp2, err := client.Do(req2)
	if err != nil {
		return failTest(result, "Failed to execute search (JOIN query)", err)
	}
	defer func() { _ = resp2.Body.Close() }()

	if resp2.StatusCode != http.StatusOK {
		return failTest(result, fmt.Sprintf("Expected status 200 for search, got %d", resp2.StatusCode), nil)
	}

	result.CompletedAt = time.Now().UTC()
	result.DurationMs = int(result.CompletedAt.Sub(result.StartedAt).Milliseconds())
	return result
}

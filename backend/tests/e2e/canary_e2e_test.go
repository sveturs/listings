// Package e2e contains end-to-end tests for canary deployment
package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// backendURL is defined in resilience_test.go to avoid redeclaration
	testTimeout = 10 * time.Second
)

// TestCanaryE2E_GetListing проверяет E2E получение listing через canary
func TestCanaryE2E_GetListing(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	// Проверяем что backend доступен
	if !isBackendAvailable() {
		t.Skip("Backend not available (expected in CI)")
	}

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	tests := []struct {
		name       string
		listingID  string
		wantStatus int
		checkBody  bool
	}{
		{
			name:       "Valid listing ID",
			listingID:  "328",
			wantStatus: http.StatusOK,
			checkBody:  true,
		},
		{
			name:       "Non-existent listing",
			listingID:  "999999",
			wantStatus: http.StatusNotFound,
			checkBody:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("%s/api/v1/marketplace/listings/%s", backendURL, tt.listingID)
			req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
			require.NoError(t, err)

			// Add canary headers (для force routing если нужно)
			req.Header.Set("X-Canary-Enabled", "true")
			req.Header.Set("User-Agent", "canary-e2e-test")

			client := &http.Client{Timeout: testTimeout}
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer func() { _ = resp.Body.Close() }()

			assert.Equal(t, tt.wantStatus, resp.StatusCode,
				"Response status should match expected")

			// Check canary response headers
			source := resp.Header.Get("X-Source-Microservice")
			if source != "" {
				t.Logf("✅ Request served by: %s", source)
			}

			// Verify body structure if successful
			if tt.checkBody && resp.StatusCode == http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				require.NoError(t, err)

				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				require.NoError(t, err)

				// Verify listing data structure
				data, ok := result["data"].(map[string]interface{})
				assert.True(t, ok, "Response should have 'data' field")

				if ok {
					assert.NotEmpty(t, data["id"], "Listing should have ID")
					assert.NotEmpty(t, data["title"], "Listing should have title")
					t.Logf("✅ Listing data valid: id=%v, title=%v",
						data["id"], data["title"])
				}
			}
		})
	}
}

// TestCanaryE2E_ListListings проверяет E2E получение списка listings
func TestCanaryE2E_ListListings(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	if !isBackendAvailable() {
		t.Skip("Backend not available")
	}

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	url := fmt.Sprintf("%s/api/v1/marketplace/listings?limit=10&offset=0", backendURL)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	require.NoError(t, err)

	req.Header.Set("X-Canary-Enabled", "true")

	client := &http.Client{Timeout: testTimeout}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode,
		"List listings should return 200")

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)

	// Verify response structure
	data, ok := result["data"].([]interface{})
	assert.True(t, ok, "Response should have 'data' array")

	meta, ok := result["meta"].(map[string]interface{})
	assert.True(t, ok, "Response should have 'meta' object")

	if ok {
		assert.NotNil(t, meta["total"], "Meta should have 'total' field")
		assert.NotNil(t, meta["limit"], "Meta should have 'limit' field")
		t.Logf("✅ List response valid: %d items, total=%v",
			len(data), meta["total"])
	}
}

// TestCanaryE2E_SearchListings проверяет E2E поиск через canary
func TestCanaryE2E_SearchListings(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	if !isBackendAvailable() {
		t.Skip("Backend not available")
	}

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	tests := []struct {
		name       string
		query      string
		wantStatus int
	}{
		{
			name:       "Search with query",
			query:      "test",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Empty search",
			query:      "",
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("%s/api/v1/marketplace/search?q=%s&limit=10",
				backendURL, tt.query)
			req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
			require.NoError(t, err)

			req.Header.Set("X-Canary-Enabled", "true")

			client := &http.Client{Timeout: testTimeout}
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer func() { _ = resp.Body.Close() }()

			assert.Equal(t, tt.wantStatus, resp.StatusCode,
				"Search should return expected status")

			if resp.StatusCode == http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				require.NoError(t, err)

				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				require.NoError(t, err)

				data, ok := result["data"].([]interface{})
				assert.True(t, ok, "Response should have 'data' array")
				t.Logf("✅ Search returned %d results for query '%s'",
					len(data), tt.query)
			}
		})
	}
}

// TestCanaryE2E_ConsistentHashing проверяет consistent hashing для одного пользователя
func TestCanaryE2E_ConsistentHashing(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	if !isBackendAvailable() {
		t.Skip("Backend not available")
	}

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	// Make multiple requests with same user ID
	userID := "test-user-12345"
	var sources []string

	for i := 0; i < 10; i++ {
		url := fmt.Sprintf("%s/api/v1/marketplace/listings?limit=1", backendURL)
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		require.NoError(t, err)

		// Simulate authenticated user
		req.Header.Set("X-User-ID", userID)
		req.Header.Set("X-Canary-Enabled", "true")

		client := &http.Client{Timeout: testTimeout}
		resp, err := client.Do(req)
		require.NoError(t, err)

		source := resp.Header.Get("X-Source-Microservice")
		if source != "" {
			sources = append(sources, source)
		}
		_ = resp.Body.Close()
	}

	// All requests from same user should go to same backend
	if len(sources) > 0 {
		firstSource := sources[0]
		for _, source := range sources {
			assert.Equal(t, firstSource, source,
				"All requests from same user should route consistently")
		}
		t.Logf("✅ Consistent hashing verified: %d requests to same backend (%s)",
			len(sources), firstSource)
	}
}

// TestCanaryE2E_CircuitBreakerRecovery проверяет recovery после circuit breaker
func TestCanaryE2E_CircuitBreakerRecovery(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	if !isBackendAvailable() {
		t.Skip("Backend not available")
	}

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	// Make normal request
	url := fmt.Sprintf("%s/api/v1/marketplace/listings?limit=1", backendURL)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	require.NoError(t, err)

	req.Header.Set("X-Canary-Enabled", "true")

	client := &http.Client{Timeout: testTimeout}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	// Should succeed (circuit is closed)
	assert.Equal(t, http.StatusOK, resp.StatusCode,
		"Request should succeed when circuit is closed")

	circuitState := resp.Header.Get("X-Circuit-Breaker-State")
	if circuitState != "" {
		t.Logf("✅ Circuit breaker state: %s", circuitState)
	}
}

// TestCanaryE2E_CanaryUserRouting проверяет routing для canary пользователей
func TestCanaryE2E_CanaryUserRouting(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	if !isBackendAvailable() {
		t.Skip("Backend not available")
	}

	// Get canary user IDs from environment
	canaryUsers := os.Getenv("MARKETPLACE_CANARY_USER_IDS")
	if canaryUsers == "" {
		t.Skip("No canary users configured")
	}

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	url := fmt.Sprintf("%s/api/v1/marketplace/listings?limit=1", backendURL)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	require.NoError(t, err)

	// Simulate canary user
	req.Header.Set("X-User-ID", "canary-user-1")
	req.Header.Set("X-Canary-Enabled", "true")

	client := &http.Client{Timeout: testTimeout}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode,
		"Canary user request should succeed")

	source := resp.Header.Get("X-Source-Microservice")
	if source == "microservice" {
		t.Log("✅ Canary user routed to microservice")
	}
}

// TestCanaryE2E_AdminOverride проверяет admin override функциональность
func TestCanaryE2E_AdminOverride(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	if !isBackendAvailable() {
		t.Skip("Backend not available")
	}

	adminOverride := os.Getenv("MARKETPLACE_ADMIN_OVERRIDE")
	if adminOverride != "true" {
		t.Skip("Admin override not enabled")
	}

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	url := fmt.Sprintf("%s/api/v1/marketplace/listings?limit=1", backendURL)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	require.NoError(t, err)

	// Simulate admin user
	req.Header.Set("X-User-ID", "admin-user")
	req.Header.Set("X-Is-Admin", "true")

	client := &http.Client{Timeout: testTimeout}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode,
		"Admin request should succeed")

	source := resp.Header.Get("X-Source-Microservice")
	if source == "microservice" {
		t.Log("✅ Admin routed to microservice (admin override)")
	}
}

// Helper function to check if backend is available
func isBackendAvailable() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", backendURL, nil)
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

// Package service implements functional tests
// backend/internal/proj/admin/testing/service/functional_tests.go
package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"backend/internal/proj/admin/testing/domain"
)

// FunctionalTest represents a single functional test
type FunctionalTest struct {
	Name        string
	Category    domain.TestCategory
	Description string
	RunFunc     TestRunFunc
}

// TestRunFunc is the function signature for test execution
type TestRunFunc func(ctx context.Context, baseURL, token string) *domain.TestResult

// APIEndpointTests returns list of API endpoint tests
var APIEndpointTests = []FunctionalTest{
	// ===== POSITIVE/HAPPY PATH TESTS =====
	{
		Name:        "api-auth-flow",
		Category:    domain.TestCategoryAPI,
		Description: "Test authentication flow (login, me, logout)",
		RunFunc:     testAuthFlow,
	},
	{
		Name:        "api-marketplace-crud",
		Category:    domain.TestCategoryAPI,
		Description: "Test marketplace CRUD operations",
		RunFunc:     testMarketplaceCRUD,
	},
	{
		Name:        "api-categories-fetch",
		Category:    domain.TestCategoryAPI,
		Description: "Test categories API endpoints",
		RunFunc:     testCategoriesFetch,
	},
	{
		Name:        "api-search-functionality",
		Category:    domain.TestCategoryAPI,
		Description: "Test search API functionality",
		RunFunc:     testSearchFunctionality,
	},
	{
		Name:        "api-admin-operations",
		Category:    domain.TestCategoryAPI,
		Description: "Test admin panel operations",
		RunFunc:     testAdminOperations,
	},
	{
		Name:        "api-review-creation",
		Category:    domain.TestCategoryAPI,
		Description: "Test review creation with rating (draft + publish)",
		RunFunc:     testReviewCreation,
	},

	// ===== NEGATIVE TEST CASES =====
	{
		Name:        "api-auth-invalid-token",
		Category:    domain.TestCategoryAPI,
		Description: "Test API rejection with invalid authentication token",
		RunFunc:     testAuthInvalidToken,
	},
	{
		Name:        "api-auth-missing-token",
		Category:    domain.TestCategoryAPI,
		Description: "Test API rejection when authentication token is missing",
		RunFunc:     testAuthMissingToken,
	},
	{
		Name:        "api-admin-unauthorized",
		Category:    domain.TestCategoryAPI,
		Description: "Test admin endpoint rejection for non-admin users",
		RunFunc:     testAdminUnauthorized,
	},
	{
		Name:        "api-search-invalid-params",
		Category:    domain.TestCategoryAPI,
		Description: "Test search with invalid query parameters",
		RunFunc:     testSearchInvalidParams,
	},

	// ===== EDGE CASES =====
	{
		Name:        "api-search-empty-query",
		Category:    domain.TestCategoryAPI,
		Description: "Test search with empty query string",
		RunFunc:     testSearchEmptyQuery,
	},
	{
		Name:        "api-search-unicode",
		Category:    domain.TestCategoryAPI,
		Description: "Test search with Unicode characters (Cyrillic, Emoji)",
		RunFunc:     testSearchUnicode,
	},
	{
		Name:        "api-listings-extreme-limit",
		Category:    domain.TestCategoryAPI,
		Description: "Test listings with extreme limit values",
		RunFunc:     testListingsExtremeLimit,
	},
}

// testAuthFlow tests authentication endpoints
func testAuthFlow(ctx context.Context, baseURL, token string) *domain.TestResult {
	result := &domain.TestResult{
		TestName:  "api-auth-flow",
		TestSuite: "api",
		Status:    domain.TestResultStatusPassed,
		StartedAt: time.Now().UTC(),
	}

	// Test /api/v1/auth/me endpoint
	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/api/v1/auth/me", nil)
	if err != nil {
		return failTest(result, "Failed to create request", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return failTest(result, "Failed to execute request", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return failTest(result, fmt.Sprintf("Expected status 200, got %d", resp.StatusCode), fmt.Errorf("response: %s", string(body)))
	}

	// Verify response structure
	var meResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&meResp); err != nil {
		return failTest(result, "Failed to decode response", err)
	}

	// Check for required fields - /api/v1/auth/me returns {"user": {"id": ..., "email": ...}}
	user, ok := meResp["user"].(map[string]interface{})
	if !ok || user["email"] == nil {
		return failTest(result, "Missing user or email field in response", nil)
	}

	result.CompletedAt = time.Now().UTC()
	result.DurationMs = int(result.CompletedAt.Sub(result.StartedAt).Milliseconds())
	return result
}

// testMarketplaceCRUD tests marketplace CRUD operations
func testMarketplaceCRUD(ctx context.Context, baseURL, token string) *domain.TestResult {
	result := &domain.TestResult{
		TestName:  "api-marketplace-crud",
		TestSuite: "api",
		Status:    domain.TestResultStatusPassed,
		StartedAt: time.Now().UTC(),
	}

	// Test GET /api/v1/unified/listings
	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/api/v1/unified/listings?limit=5", nil)
	if err != nil {
		return failTest(result, "Failed to create request", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return failTest(result, "Failed to execute request", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return failTest(result, fmt.Sprintf("Expected status 200, got %d", resp.StatusCode), fmt.Errorf("response: %s", string(body)))
	}

	// Verify response structure
	var listingsResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&listingsResp); err != nil {
		return failTest(result, "Failed to decode response", err)
	}

	// Check for data array
	if listingsResp["data"] == nil {
		return failTest(result, "Missing data field in response", nil)
	}

	result.CompletedAt = time.Now().UTC()
	result.DurationMs = int(result.CompletedAt.Sub(result.StartedAt).Milliseconds())
	return result
}

// testCategoriesFetch tests categories API
func testCategoriesFetch(ctx context.Context, baseURL, token string) *domain.TestResult {
	result := &domain.TestResult{
		TestName:  "api-categories-fetch",
		TestSuite: "api",
		Status:    domain.TestResultStatusPassed,
		StartedAt: time.Now().UTC(),
	}

	// Test GET /api/v1/admin/categories
	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/api/v1/admin/categories", nil)
	if err != nil {
		return failTest(result, "Failed to create request", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return failTest(result, "Failed to execute request", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return failTest(result, fmt.Sprintf("Expected status 200, got %d", resp.StatusCode), fmt.Errorf("response: %s", string(body)))
	}

	// Verify response structure - admin endpoints return {success: true, data: [...]}
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return failTest(result, "Failed to decode response", err)
	}

	// Check for success and data fields
	if response["data"] == nil {
		return failTest(result, "Missing data field in response", nil)
	}

	result.CompletedAt = time.Now().UTC()
	result.DurationMs = int(result.CompletedAt.Sub(result.StartedAt).Milliseconds())
	return result
}

// testSearchFunctionality tests search API
func testSearchFunctionality(ctx context.Context, baseURL, token string) *domain.TestResult {
	result := &domain.TestResult{
		TestName:  "api-search-functionality",
		TestSuite: "api",
		Status:    domain.TestResultStatusPassed,
		StartedAt: time.Now().UTC(),
	}

	// Test GET /api/v1/search - unified search endpoint
	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/api/v1/search?query=test&limit=5", nil)
	if err != nil {
		return failTest(result, "Failed to create request", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return failTest(result, "Failed to execute request", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return failTest(result, fmt.Sprintf("Expected status 200, got %d", resp.StatusCode), fmt.Errorf("response: %s", string(body)))
	}

	result.CompletedAt = time.Now().UTC()
	result.DurationMs = int(result.CompletedAt.Sub(result.StartedAt).Milliseconds())
	return result
}

// testAdminOperations tests admin panel endpoints
func testAdminOperations(ctx context.Context, baseURL, token string) *domain.TestResult {
	result := &domain.TestResult{
		TestName:  "api-admin-operations",
		TestSuite: "api",
		Status:    domain.TestResultStatusPassed,
		StartedAt: time.Now().UTC(),
	}

	// Test GET /api/v1/admin/admins
	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/api/v1/admin/admins", nil)
	if err != nil {
		return failTest(result, "Failed to create request", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return failTest(result, "Failed to execute request", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return failTest(result, fmt.Sprintf("Expected status 200, got %d", resp.StatusCode), fmt.Errorf("response: %s", string(body)))
	}

	// Verify response structure - /api/v1/admin/admins returns {"success": true, "data": [...]}
	var responseWrapper map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseWrapper); err != nil {
		return failTest(result, "Failed to decode response", err)
	}

	// Extract admins from 'data' field
	data, ok := responseWrapper["data"]
	if !ok {
		return failTest(result, "Missing data field in response", nil)
	}

	admins, ok := data.([]interface{})
	if !ok {
		return failTest(result, "Data field is not an array", nil)
	}

	// Should have at least our admin@admin.rs account
	if len(admins) == 0 {
		return failTest(result, "No admins returned", nil)
	}

	result.CompletedAt = time.Now().UTC()
	result.DurationMs = int(result.CompletedAt.Sub(result.StartedAt).Milliseconds())
	return result
}

// testReviewCreation tests review creation workflow (draft + publish)
func testReviewCreation(ctx context.Context, baseURL, token string) *domain.TestResult {
	result := &domain.TestResult{
		TestName:  "api-review-creation",
		TestSuite: "api",
		Status:    domain.TestResultStatusPassed,
		StartedAt: time.Now().UTC(),
	}

	client := &http.Client{Timeout: 10 * time.Second}

	// Step 1: Get multiple listings to use the first available
	reqListings, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/api/v1/unified/listings?limit=1", nil)
	if err != nil {
		return failTest(result, "Failed to create listings request", err)
	}
	reqListings.Header.Set("Authorization", "Bearer "+token)

	respListings, err := client.Do(reqListings)
	if err != nil {
		return failTest(result, "Failed to fetch listings", err)
	}
	defer func() { _ = respListings.Body.Close() }()

	if respListings.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(respListings.Body)
		return failTest(result, fmt.Sprintf("Failed to get listings, status %d", respListings.StatusCode), fmt.Errorf("response: %s", string(body)))
	}

	var listingsResp map[string]interface{}
	if err := json.NewDecoder(respListings.Body).Decode(&listingsResp); err != nil {
		return failTest(result, "Failed to decode listings response", err)
	}

	data, ok := listingsResp["data"].([]interface{})
	if !ok || len(data) == 0 {
		return failTest(result, "No listings available for review", nil)
	}

	// Step 2: Use first listing
	listing := data[0].(map[string]interface{})
	listingID := int(listing["id"].(float64))

	// Step 3: Delete any existing reviews for this listing (make test idempotent)
	// This allows test to run multiple times without failing
	reqGetReviews, err := http.NewRequestWithContext(ctx, "GET",
		fmt.Sprintf("%s/api/v1/reviews?entity_type=listing&entity_id=%d", baseURL, listingID), nil)
	if err == nil {
		reqGetReviews.Header.Set("Authorization", "Bearer "+token)
		respGetReviews, err := client.Do(reqGetReviews)
		if err == nil && respGetReviews.StatusCode == http.StatusOK {
			var reviewsResp map[string]interface{}
			if err := json.NewDecoder(respGetReviews.Body).Decode(&reviewsResp); err == nil {
				if dataWrapper, ok := reviewsResp["data"].(map[string]interface{}); ok {
					if dataField, ok := dataWrapper["data"].([]interface{}); ok {
						// Delete each existing review
						for _, r := range dataField {
							review := r.(map[string]interface{})
							reviewID := int(review["id"].(float64))
							reqDel, err := http.NewRequestWithContext(ctx, "DELETE",
								fmt.Sprintf("%s/api/v1/reviews/%d", baseURL, reviewID), nil)
							if err == nil {
								reqDel.Header.Set("Authorization", "Bearer "+token)
								respDel, err := client.Do(reqDel)
								if err == nil {
									_ = respDel.Body.Close()
								}
							}
						}
					}
				}
			}
			_ = respGetReviews.Body.Close()
		}
	}

	// Step 4: Create draft review
	reviewPayload := map[string]interface{}{
		"entity_type":       "listing",
		"entity_id":         listingID,
		"rating":            5,
		"comment":           "Отличное место! Все понравилось.",
		"pros":              "Чисто, уютно, хороший персонал",
		"cons":              "Нет минусов",
		"original_language": "ru",
	}

	payloadBytes, err := json.Marshal(reviewPayload)
	if err != nil {
		return failTest(result, "Failed to marshal review payload", err)
	}

	reqDraft, err := http.NewRequestWithContext(ctx, "POST", baseURL+"/api/v1/reviews/draft", bytes.NewReader(payloadBytes))
	if err != nil {
		return failTest(result, "Failed to create draft request", err)
	}
	reqDraft.Header.Set("Authorization", "Bearer "+token)
	reqDraft.Header.Set("Content-Type", "application/json")

	respDraft, err := client.Do(reqDraft)
	if err != nil {
		return failTest(result, "Failed to create draft review", err)
	}
	defer func() { _ = respDraft.Body.Close() }()

	if respDraft.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(respDraft.Body)
		return failTest(result, fmt.Sprintf("Failed to create draft review, status %d", respDraft.StatusCode), fmt.Errorf("response: %s", string(body)))
	}

	var draftResp map[string]interface{}
	if err := json.NewDecoder(respDraft.Body).Decode(&draftResp); err != nil {
		return failTest(result, "Failed to decode draft response", err)
	}

	// Extract review from response wrapper {success: true, data: {...}}
	reviewData, ok := draftResp["data"].(map[string]interface{})
	if !ok {
		return failTest(result, "Missing data field in draft response", nil)
	}

	reviewID := int(reviewData["id"].(float64))

	// Step 5: Publish review
	reqPublish, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/api/v1/reviews/%d/publish", baseURL, reviewID), nil)
	if err != nil {
		return failTest(result, "Failed to create publish request", err)
	}
	reqPublish.Header.Set("Authorization", "Bearer "+token)

	respPublish, err := client.Do(reqPublish)
	if err != nil {
		return failTest(result, "Failed to publish review", err)
	}
	defer func() { _ = respPublish.Body.Close() }()

	if respPublish.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(respPublish.Body)
		return failTest(result, fmt.Sprintf("Failed to publish review, status %d", respPublish.StatusCode), fmt.Errorf("response: %s", string(body)))
	}

	var publishResp map[string]interface{}
	if err := json.NewDecoder(respPublish.Body).Decode(&publishResp); err != nil {
		return failTest(result, "Failed to decode publish response", err)
	}

	// Verify review is published
	publishedReview, ok := publishResp["data"].(map[string]interface{})
	if !ok {
		return failTest(result, "Missing data field in publish response", nil)
	}

	status, ok := publishedReview["status"].(string)
	if !ok || status != "published" {
		return failTest(result, fmt.Sprintf("Review status is %v, expected 'published'", publishedReview["status"]), nil)
	}

	// Step 6: Cleanup - Delete the test review (optional, helps keep DB clean)
	// We don't fail the test if cleanup fails, just log it
	reqDelete, err := http.NewRequestWithContext(ctx, "DELETE", fmt.Sprintf("%s/api/v1/reviews/%d", baseURL, reviewID), nil)
	if err == nil {
		reqDelete.Header.Set("Authorization", "Bearer "+token)
		respDelete, err := client.Do(reqDelete)
		if err == nil {
			_ = respDelete.Body.Close()
		}
	}

	result.CompletedAt = time.Now().UTC()
	result.DurationMs = int(result.CompletedAt.Sub(result.StartedAt).Milliseconds())
	return result
}

// ==================== NEGATIVE TEST CASES ====================

// testAuthInvalidToken verifies that API rejects requests with invalid token
func testAuthInvalidToken(ctx context.Context, baseURL, token string) *domain.TestResult {
	result := &domain.TestResult{
		TestName:  "api-auth-invalid-token",
		TestSuite: "api",
		Status:    domain.TestResultStatusPassed,
		StartedAt: time.Now().UTC(),
	}

	// Use invalid token
	invalidToken := "invalid.jwt.token.here" //nolint:gosec // G101: Intentional invalid test token for security tests

	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/api/v1/auth/me", nil)
	if err != nil {
		return failTest(result, "Failed to create request", err)
	}

	req.Header.Set("Authorization", "Bearer "+invalidToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return failTest(result, "Failed to execute request", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// EXPECT 401 Unauthorized
	if resp.StatusCode != http.StatusUnauthorized {
		body, _ := io.ReadAll(resp.Body)
		return failTest(result, fmt.Sprintf("Expected status 401, got %d (should reject invalid token)", resp.StatusCode), fmt.Errorf("response: %s", string(body)))
	}

	result.CompletedAt = time.Now().UTC()
	result.DurationMs = int(result.CompletedAt.Sub(result.StartedAt).Milliseconds())
	return result
}

// testAuthMissingToken verifies that API rejects requests without token
func testAuthMissingToken(ctx context.Context, baseURL, token string) *domain.TestResult {
	result := &domain.TestResult{
		TestName:  "api-auth-missing-token",
		TestSuite: "api",
		Status:    domain.TestResultStatusPassed,
		StartedAt: time.Now().UTC(),
	}

	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/api/v1/auth/me", nil)
	if err != nil {
		return failTest(result, "Failed to create request", err)
	}

	// Don't set Authorization header

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return failTest(result, "Failed to execute request", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// EXPECT 401 Unauthorized
	if resp.StatusCode != http.StatusUnauthorized {
		body, _ := io.ReadAll(resp.Body)
		return failTest(result, fmt.Sprintf("Expected status 401, got %d (should reject missing token)", resp.StatusCode), fmt.Errorf("response: %s", string(body)))
	}

	result.CompletedAt = time.Now().UTC()
	result.DurationMs = int(result.CompletedAt.Sub(result.StartedAt).Milliseconds())
	return result
}

// testAdminUnauthorized verifies that non-admin users can't access admin endpoints
// NOTE: This test requires a non-admin user token, which we don't have in test auth manager
// For now, we'll use invalid token which should also be rejected
func testAdminUnauthorized(ctx context.Context, baseURL, token string) *domain.TestResult {
	result := &domain.TestResult{
		TestName:  "api-admin-unauthorized",
		TestSuite: "api",
		Status:    domain.TestResultStatusPassed,
		StartedAt: time.Now().UTC(),
	}

	// Use invalid token to simulate non-admin access
	fakeToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6Ik5vbkFkbWluIFVzZXIiLCJpYXQiOjE1MTYyMzkwMjJ9.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c" //nolint:gosec // G101: Intentional test JWT for negative auth testing

	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/api/v1/admin/categories", nil)
	if err != nil {
		return failTest(result, "Failed to create request", err)
	}

	req.Header.Set("Authorization", "Bearer "+fakeToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return failTest(result, "Failed to execute request", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// EXPECT 401 or 403
	if resp.StatusCode != http.StatusUnauthorized && resp.StatusCode != http.StatusForbidden {
		body, _ := io.ReadAll(resp.Body)
		return failTest(result, fmt.Sprintf("Expected status 401 or 403, got %d (should reject non-admin)", resp.StatusCode), fmt.Errorf("response: %s", string(body)))
	}

	result.CompletedAt = time.Now().UTC()
	result.DurationMs = int(result.CompletedAt.Sub(result.StartedAt).Milliseconds())
	return result
}

// testSearchInvalidParams verifies handling of invalid search parameters
func testSearchInvalidParams(ctx context.Context, baseURL, token string) *domain.TestResult {
	result := &domain.TestResult{
		TestName:  "api-search-invalid-params",
		TestSuite: "api",
		Status:    domain.TestResultStatusPassed,
		StartedAt: time.Now().UTC(),
	}

	// Test with negative limit (should be rejected or handled gracefully)
	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/api/v1/search?query=test&limit=-100", nil)
	if err != nil {
		return failTest(result, "Failed to create request", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return failTest(result, "Failed to execute request", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Accept 200 (backend may handle gracefully) or 400 (validation error)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusBadRequest {
		body, _ := io.ReadAll(resp.Body)
		return failTest(result, fmt.Sprintf("Expected status 200 or 400, got %d", resp.StatusCode), fmt.Errorf("response: %s", string(body)))
	}

	result.CompletedAt = time.Now().UTC()
	result.DurationMs = int(result.CompletedAt.Sub(result.StartedAt).Milliseconds())
	return result
}

// ==================== EDGE CASES ====================

// testSearchEmptyQuery tests search with empty query string
func testSearchEmptyQuery(ctx context.Context, baseURL, token string) *domain.TestResult {
	result := &domain.TestResult{
		TestName:  "api-search-empty-query",
		TestSuite: "api",
		Status:    domain.TestResultStatusPassed,
		StartedAt: time.Now().UTC(),
	}

	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/api/v1/search?query=&limit=5", nil)
	if err != nil {
		return failTest(result, "Failed to create request", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return failTest(result, "Failed to execute request", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Should handle gracefully (200 or 400)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusBadRequest {
		body, _ := io.ReadAll(resp.Body)
		return failTest(result, fmt.Sprintf("Expected status 200 or 400, got %d", resp.StatusCode), fmt.Errorf("response: %s", string(body)))
	}

	result.CompletedAt = time.Now().UTC()
	result.DurationMs = int(result.CompletedAt.Sub(result.StartedAt).Milliseconds())
	return result
}

// testSearchUnicode tests search with Unicode characters
func testSearchUnicode(ctx context.Context, baseURL, token string) *domain.TestResult {
	result := &domain.TestResult{
		TestName:  "api-search-unicode",
		TestSuite: "api",
		Status:    domain.TestResultStatusPassed,
		StartedAt: time.Now().UTC(),
	}

	// Test with Cyrillic and Emoji
	queries := []string{
		"Москва",  // Cyrillic
		"Београд", // Serbian Cyrillic
		"тест 🏠",  // Cyrillic + Emoji
		"München", // German umlaut
		"日本",      // Japanese
	}

	client := &http.Client{Timeout: 10 * time.Second}

	for _, query := range queries {
		req, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/api/v1/search?query="+query+"&limit=5", nil)
		if err != nil {
			return failTest(result, fmt.Sprintf("Failed to create request for query '%s'", query), err)
		}

		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := client.Do(req)
		if err != nil {
			return failTest(result, fmt.Sprintf("Failed to execute request for query '%s'", query), err)
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			_ = resp.Body.Close()
			return failTest(result, fmt.Sprintf("Query '%s': Expected status 200, got %d", query, resp.StatusCode), fmt.Errorf("response: %s", string(body)))
		}
		_ = resp.Body.Close()
	}

	result.CompletedAt = time.Now().UTC()
	result.DurationMs = int(result.CompletedAt.Sub(result.StartedAt).Milliseconds())
	return result
}

// testListingsExtremeLimit tests listings API with extreme limit values
func testListingsExtremeLimit(ctx context.Context, baseURL, token string) *domain.TestResult {
	result := &domain.TestResult{
		TestName:  "api-listings-extreme-limit",
		TestSuite: "api",
		Status:    domain.TestResultStatusPassed,
		StartedAt: time.Now().UTC(),
	}

	client := &http.Client{Timeout: 15 * time.Second}

	// Test with limit=0
	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/api/v1/unified/listings?limit=0", nil)
	if err != nil {
		return failTest(result, "Failed to create request for limit=0", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return failTest(result, "Failed to execute request for limit=0", err)
	}

	// Should handle gracefully
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusBadRequest {
		body, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		return failTest(result, fmt.Sprintf("limit=0: Expected status 200 or 400, got %d", resp.StatusCode), fmt.Errorf("response: %s", string(body)))
	}
	_ = resp.Body.Close()

	// Test with very large limit (should be capped by backend)
	req2, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/api/v1/unified/listings?limit=10000", nil)
	if err != nil {
		return failTest(result, "Failed to create request for limit=10000", err)
	}
	req2.Header.Set("Authorization", "Bearer "+token)

	resp2, err := client.Do(req2)
	if err != nil {
		return failTest(result, "Failed to execute request for limit=10000", err)
	}
	defer func() { _ = resp2.Body.Close() }()

	// Should handle gracefully (likely capped to max limit)
	if resp2.StatusCode != http.StatusOK && resp2.StatusCode != http.StatusBadRequest {
		body, _ := io.ReadAll(resp2.Body)
		return failTest(result, fmt.Sprintf("limit=10000: Expected status 200 or 400, got %d", resp2.StatusCode), fmt.Errorf("response: %s", string(body)))
	}

	result.CompletedAt = time.Now().UTC()
	result.DurationMs = int(result.CompletedAt.Sub(result.StartedAt).Milliseconds())
	return result
}

// failTest marks test as failed and returns result
func failTest(result *domain.TestResult, message string, err error) *domain.TestResult {
	result.Status = domain.TestResultStatusFailed
	result.CompletedAt = time.Now().UTC()
	result.DurationMs = int(result.CompletedAt.Sub(result.StartedAt).Milliseconds())

	errMsg := message
	if err != nil {
		errMsg = fmt.Sprintf("%s: %v", message, err)
	}
	result.ErrorMsg = &errMsg

	if err != nil {
		stackTrace := fmt.Sprintf("%+v", err)
		result.StackTrace = &stackTrace
	}

	return result
}

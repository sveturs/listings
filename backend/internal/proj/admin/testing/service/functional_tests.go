// Package service implements functional tests
// backend/internal/proj/admin/testing/service/functional_tests.go
package service

import (
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
}

// testAuthFlow tests authentication endpoints
func testAuthFlow(ctx context.Context, baseURL, token string) *domain.TestResult {
	result := &domain.TestResult{
		TestName:   "api-auth-flow",
		TestSuite:  "api",
		Status:     domain.TestResultStatusPassed,
		StartedAt:  time.Now(),
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
	defer resp.Body.Close()

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

	result.CompletedAt = time.Now()
	result.DurationMs = int(result.CompletedAt.Sub(result.StartedAt).Milliseconds())
	return result
}

// testMarketplaceCRUD tests marketplace CRUD operations
func testMarketplaceCRUD(ctx context.Context, baseURL, token string) *domain.TestResult {
	result := &domain.TestResult{
		TestName:   "api-marketplace-crud",
		TestSuite:  "api",
		Status:     domain.TestResultStatusPassed,
		StartedAt:  time.Now(),
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
	defer resp.Body.Close()

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

	result.CompletedAt = time.Now()
	result.DurationMs = int(result.CompletedAt.Sub(result.StartedAt).Milliseconds())
	return result
}

// testCategoriesFetch tests categories API
func testCategoriesFetch(ctx context.Context, baseURL, token string) *domain.TestResult {
	result := &domain.TestResult{
		TestName:   "api-categories-fetch",
		TestSuite:  "api",
		Status:     domain.TestResultStatusPassed,
		StartedAt:  time.Now(),
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
	defer resp.Body.Close()

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

	result.CompletedAt = time.Now()
	result.DurationMs = int(result.CompletedAt.Sub(result.StartedAt).Milliseconds())
	return result
}

// testSearchFunctionality tests search API
func testSearchFunctionality(ctx context.Context, baseURL, token string) *domain.TestResult {
	result := &domain.TestResult{
		TestName:   "api-search-functionality",
		TestSuite:  "api",
		Status:     domain.TestResultStatusPassed,
		StartedAt:  time.Now(),
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
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return failTest(result, fmt.Sprintf("Expected status 200, got %d", resp.StatusCode), fmt.Errorf("response: %s", string(body)))
	}

	result.CompletedAt = time.Now()
	result.DurationMs = int(result.CompletedAt.Sub(result.StartedAt).Milliseconds())
	return result
}

// testAdminOperations tests admin panel endpoints
func testAdminOperations(ctx context.Context, baseURL, token string) *domain.TestResult {
	result := &domain.TestResult{
		TestName:   "api-admin-operations",
		TestSuite:  "api",
		Status:     domain.TestResultStatusPassed,
		StartedAt:  time.Now(),
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
	defer resp.Body.Close()

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

	result.CompletedAt = time.Now()
	result.DurationMs = int(result.CompletedAt.Sub(result.StartedAt).Milliseconds())
	return result
}

// failTest marks test as failed and returns result
func failTest(result *domain.TestResult, message string, err error) *domain.TestResult {
	result.Status = domain.TestResultStatusFailed
	result.CompletedAt = time.Now()
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

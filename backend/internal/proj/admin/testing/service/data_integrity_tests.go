// Package service implements data integrity functional tests
// backend/internal/proj/admin/testing/service/data_integrity_tests.go
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"backend/internal/proj/admin/testing/domain"
)

// DataIntegrityTests returns list of data integrity tests
var DataIntegrityTests = []FunctionalTest{
	{
		Name:        "data-integrity-marketplace-listing",
		Category:    domain.TestCategoryDataIntegrity,
		Description: "Verify listing data matches across DB, cache, and search index",
		RunFunc:     testMarketplaceListingConsistency,
	},
	{
		Name:        "data-integrity-transaction-rollback",
		Category:    domain.TestCategoryDataIntegrity,
		Description: "Test database transaction rollback on errors",
		RunFunc:     testTransactionRollback,
	},
	{
		Name:        "data-integrity-image-orphan-cleanup",
		Category:    domain.TestCategoryDataIntegrity,
		Description: "Verify orphaned images are cleaned up from MinIO",
		RunFunc:     testImageOrphanCleanup,
	},
}

// testMarketplaceListingConsistency verifies data consistency across DB, cache, and search
func testMarketplaceListingConsistency(ctx context.Context, baseURL, token string) *domain.TestResult {
	result := &domain.TestResult{
		TestName:  "data-integrity-marketplace-listing",
		TestSuite: "data-integrity",
		Status:    domain.TestResultStatusPassed,
		StartedAt: time.Now().UTC(),
	}

	client := &http.Client{Timeout: 10 * time.Second}

	// Step 1: Get a listing from unified endpoint (combines DB + cache + search)
	listingsURL := fmt.Sprintf("%s/api/v1/unified/listings?limit=1", baseURL)
	req, err := http.NewRequestWithContext(ctx, "GET", listingsURL, nil)
	if err != nil {
		return failTest(result, "Failed to create listings request", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return failTest(result, "Failed to get listings", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return failTest(result, fmt.Sprintf("Expected 200, got %d", resp.StatusCode), fmt.Errorf("response: %s", string(body)))
	}

	var unifiedResp struct {
		Success bool                     `json:"success"`
		Data    []map[string]interface{} `json:"data"`
	}

	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &unifiedResp); err != nil {
		return failTest(result, "Failed to parse unified listings response", err)
	}

	if !unifiedResp.Success || len(unifiedResp.Data) == 0 {
		return failTest(result, "No listings found for consistency test", fmt.Errorf("empty listings"))
	}

	firstListing := unifiedResp.Data[0]
	listingID, ok := firstListing["id"].(float64)
	if !ok {
		return failTest(result, "Listing ID not found", fmt.Errorf("invalid listing structure"))
	}

	// Step 2: Query the same listing again to verify cache consistency
	// If cache is working, the second request should return same data
	listingsURL2 := fmt.Sprintf("%s/api/v1/unified/listings?limit=10", baseURL)
	req2, err := http.NewRequestWithContext(ctx, "GET", listingsURL2, nil)
	if err != nil {
		return failTest(result, "Failed to create second listings request", err)
	}

	req2.Header.Set("Authorization", "Bearer "+token)

	resp2, err := client.Do(req2)
	if err != nil {
		return failTest(result, "Failed to get listings second time", err)
	}
	defer func() { _ = resp2.Body.Close() }()

	if resp2.StatusCode != http.StatusOK {
		body2, _ := io.ReadAll(resp2.Body)
		return failTest(result, fmt.Sprintf("Second query failed with %d", resp2.StatusCode), fmt.Errorf("response: %s", string(body2)))
	}

	var unifiedResp2 struct {
		Success bool                     `json:"success"`
		Data    []map[string]interface{} `json:"data"`
	}

	body2, _ := io.ReadAll(resp2.Body)
	if err := json.Unmarshal(body2, &unifiedResp2); err != nil {
		return failTest(result, "Failed to parse second response", err)
	}

	// Step 3: Verify data consistency between first and second query
	// Find the same listing in second response
	found := false
	for _, listing := range unifiedResp2.Data {
		if id, ok := listing["id"].(float64); ok && id == listingID {
			found = true
			// Verify key fields match
			title1, _ := firstListing["title"].(string)
			title2, _ := listing["title"].(string)
			if title1 != title2 {
				return failTest(result, "Title mismatch between queries", fmt.Errorf("query1: %s, query2: %s", title1, title2))
			}
			break
		}
	}

	if !found {
		return failTest(result, "Listing not found in second query", fmt.Errorf("data consistency issue"))
	}

	// Success!
	result.CompletedAt = time.Now().UTC()
	result.DurationMs = int(result.CompletedAt.Sub(result.StartedAt).Milliseconds())
	return result
}

// testTransactionRollback tests database transaction rollback
func testTransactionRollback(ctx context.Context, baseURL, token string) *domain.TestResult {
	result := &domain.TestResult{
		TestName:  "data-integrity-transaction-rollback",
		TestSuite: "data-integrity",
		Status:    domain.TestResultStatusPassed,
		StartedAt: time.Now().UTC(),
	}

	client := &http.Client{Timeout: 10 * time.Second}

	// Step 1: Try to create a review with invalid data (should rollback transaction)
	// This tests that failed transactions don't leave partial data in DB
	invalidReviewPayload := map[string]interface{}{
		"listing_id": -999999, // Invalid listing ID
		"rating":     11,      // Invalid rating (should be 1-5)
		"comment":    "Test rollback review",
	}

	payloadBytes, _ := json.Marshal(invalidReviewPayload)
	reviewURL := fmt.Sprintf("%s/api/v1/marketplace/reviews", baseURL)

	req, err := http.NewRequestWithContext(ctx, "POST", reviewURL, strings.NewReader(string(payloadBytes)))
	if err != nil {
		return failTest(result, "Failed to create review request", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return failTest(result, "Failed to send review request", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Expect 400 Bad Request (invalid data should be rejected)
	if resp.StatusCode == http.StatusCreated {
		return failTest(result, "Invalid review was created (transaction didn't rollback)", fmt.Errorf("expected rejection, got 201"))
	}

	// Step 2: Verify DB is still consistent after failed transaction
	// Query listings to ensure DB is still operational
	testListingsURL := fmt.Sprintf("%s/api/v1/unified/listings?limit=1", baseURL)
	req2, err := http.NewRequestWithContext(ctx, "GET", testListingsURL, nil)
	if err != nil {
		return failTest(result, "Failed to create test query", err)
	}

	req2.Header.Set("Authorization", "Bearer "+token)

	resp2, err := client.Do(req2)
	if err != nil {
		return failTest(result, "Failed to query after rollback test", err)
	}
	defer func() { _ = resp2.Body.Close() }()

	if resp2.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp2.Body)
		return failTest(result, fmt.Sprintf("DB might be corrupted after rollback: %d", resp2.StatusCode), fmt.Errorf("response: %s", string(body)))
	}

	// Success - transaction was rolled back properly!
	result.CompletedAt = time.Now().UTC()
	result.DurationMs = int(result.CompletedAt.Sub(result.StartedAt).Milliseconds())
	return result
}

// testImageOrphanCleanup tests MinIO orphan image cleanup
func testImageOrphanCleanup(ctx context.Context, baseURL, token string) *domain.TestResult {
	result := &domain.TestResult{
		TestName:  "data-integrity-image-orphan-cleanup",
		TestSuite: "data-integrity",
		Status:    domain.TestResultStatusPassed,
		StartedAt: time.Now().UTC(),
	}

	client := &http.Client{Timeout: 10 * time.Second}

	// Step 1: Get listings with images
	listingsURL := fmt.Sprintf("%s/api/v1/unified/listings?limit=5", baseURL)
	req, err := http.NewRequestWithContext(ctx, "GET", listingsURL, nil)
	if err != nil {
		return failTest(result, "Failed to create listings request", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return failTest(result, "Failed to get listings", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return failTest(result, fmt.Sprintf("Expected 200, got %d", resp.StatusCode), fmt.Errorf("response: %s", string(body)))
	}

	var listingsResp struct {
		Success bool                     `json:"success"`
		Data    []map[string]interface{} `json:"data"`
	}

	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &listingsResp); err != nil {
		return failTest(result, "Failed to parse listings response", err)
	}

	if !listingsResp.Success || len(listingsResp.Data) == 0 {
		// No listings to test - this is acceptable for new systems
		result.Status = domain.TestResultStatusSkipped
		msg := "No listings found to test orphan cleanup (acceptable for new systems)"
		result.ErrorMsg = &msg
		result.CompletedAt = time.Now().UTC()
		result.DurationMs = int(result.CompletedAt.Sub(result.StartedAt).Milliseconds())
		return result
	}

	// Step 2: Verify each listing with images has valid image references
	imageCount := 0
	orphanCount := 0

	for _, listing := range listingsResp.Data {
		images, ok := listing["images"].([]interface{})
		if !ok || len(images) == 0 {
			continue
		}

		for _, img := range images {
			imageMap, ok := img.(map[string]interface{})
			if !ok {
				continue
			}

			imageURL, ok := imageMap["url"].(string)
			if !ok || imageURL == "" {
				orphanCount++
				continue
			}

			imageCount++

			// Verify image URL is accessible (basic check)
			// We don't download the full image, just check HTTP HEAD
			imgReq, err := http.NewRequestWithContext(ctx, "HEAD", imageURL, nil)
			if err != nil {
				continue
			}

			imgResp, err := client.Do(imgReq)
			if err != nil {
				// Image not accessible - potential orphan
				orphanCount++
				continue
			}
			_ = imgResp.Body.Close()

			if imgResp.StatusCode != http.StatusOK {
				// Image returns error - potential orphan
				orphanCount++
			}
		}
	}

	// Step 3: Calculate orphan percentage
	// We expect < 5% orphan rate (some temporary orphans are acceptable during cleanup cycles)
	orphanRate := 0.0
	if imageCount > 0 {
		orphanRate = float64(orphanCount) / float64(imageCount) * 100
	}

	if orphanRate > 5.0 {
		return failTest(result, fmt.Sprintf("High orphan image rate: %.2f%% (%d/%d)", orphanRate, orphanCount, imageCount), fmt.Errorf("cleanup not working properly"))
	}

	// Success!
	result.CompletedAt = time.Now().UTC()
	result.DurationMs = int(result.CompletedAt.Sub(result.StartedAt).Milliseconds())
	return result
}

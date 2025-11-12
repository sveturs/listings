package integration

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	backendURL         = "http://localhost:3000"
	openSearchURL      = "http://localhost:9200"
	listingsMicroDBDSN = "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable"
	testListingID      = 74
)

// TestGetListingWithImages checks that GetListing endpoint returns images array
func TestGetListingWithImages(t *testing.T) {
	// Get JWT token
	token, err := os.ReadFile("/tmp/token")
	require.NoError(t, err, "Failed to read JWT token")

	// Make request to GetListing
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/marketplace/listings/%d?lang=ru", backendURL, testListingID), nil)
	require.NoError(t, err)
	// Trim whitespace and newlines from token
	tokenStr := string(token)
	tokenStr = tokenStr[:len(tokenStr)-1] // Remove trailing newline
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenStr))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Logf("Failed to close response body: %v", err)
		}
	}()

	require.Equal(t, http.StatusOK, resp.StatusCode, "Expected 200 OK")

	// Parse response
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var result struct {
		Data struct {
			ID     int    `json:"id"`
			Title  string `json:"title"`
			Images []struct {
				ID           int    `json:"id"`
				URL          string `json:"url"`
				ThumbnailURL string `json:"thumbnail_url"`
				IsPrimary    bool   `json:"is_primary"`
				DisplayOrder int    `json:"display_order"`
			} `json:"images"`
		} `json:"data"`
		Success bool `json:"success"`
	}
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)

	// Assertions
	assert.True(t, result.Success, "Response should be successful")
	assert.Equal(t, testListingID, result.Data.ID, "Listing ID should match")
	assert.NotEmpty(t, result.Data.Title, "Listing should have title")

	// CRITICAL: Check images array
	assert.NotNil(t, result.Data.Images, "Images field should not be nil")
	assert.NotEmpty(t, result.Data.Images, "Images array should not be empty")

	if len(result.Data.Images) > 0 {
		img := result.Data.Images[0]
		assert.Greater(t, img.ID, 0, "Image should have valid ID")
		assert.NotEmpty(t, img.URL, "Image should have URL")
		assert.NotEmpty(t, img.ThumbnailURL, "Image should have thumbnail URL")
	}
}

// TestSearchAPI checks that search endpoint works and returns listing 74
func TestSearchAPI(t *testing.T) {
	// Make request to Search API (public endpoint)
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/marketplace/search?limit=10", backendURL), nil)
	require.NoError(t, err)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Logf("Failed to close response body: %v", err)
		}
	}()

	require.Equal(t, http.StatusOK, resp.StatusCode, "Expected 200 OK")

	// Parse response
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var result struct {
		Data struct {
			Listings []struct {
				ID     int    `json:"id"`
				Title  string `json:"title"`
				Status string `json:"status"`
			} `json:"listings"`
			Total int `json:"total"`
		} `json:"data"`
		Success bool `json:"success"`
	}
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)

	// Assertions
	assert.True(t, result.Success, "Response should be successful")
	assert.Greater(t, result.Data.Total, 0, "Should have at least 1 listing")
	assert.NotEmpty(t, result.Data.Listings, "Listings array should not be empty")

	// Check if listing 74 is present
	found := false
	for _, listing := range result.Data.Listings {
		if listing.ID == testListingID {
			found = true
			assert.Equal(t, "active", listing.Status, "Listing 74 should be active")
			break
		}
	}
	assert.True(t, found, "Listing 74 should be in search results")
}

// TestOpenSearchSync checks that OpenSearch data matches database
func TestOpenSearchSync(t *testing.T) {
	// 1. Get data from listings microservice database
	db, err := sqlx.Connect("postgres", listingsMicroDBDSN)
	require.NoError(t, err)
	defer func() {
		if err := db.Close(); err != nil {
			t.Logf("Failed to close database connection: %v", err)
		}
	}()

	var dbData struct {
		ID          int       `db:"id"`
		Title       string    `db:"title"`
		Status      string    `db:"status"`
		PublishedAt time.Time `db:"published_at"`
		ImagesCount int       `db:"images_count"`
	}

	err = db.QueryRowx(`
		SELECT
			l.id,
			l.title,
			l.status,
			l.published_at,
			(SELECT COUNT(*) FROM listing_images WHERE listing_id = $1) as images_count
		FROM listings l
		WHERE l.id = $1
	`, testListingID).StructScan(&dbData)
	require.NoError(t, err, "Failed to fetch from database")

	// Assertions on DB data
	assert.Equal(t, "active", dbData.Status, "Listing should be active in DB")
	assert.False(t, dbData.PublishedAt.IsZero(), "Listing should have published_at")
	assert.Greater(t, dbData.ImagesCount, 0, "Listing should have images in DB")

	// 2. Get data from OpenSearch
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/marketplace_listings/_doc/c2c_%d", openSearchURL, testListingID), nil)
	require.NoError(t, err)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Logf("Failed to close response body: %v", err)
		}
	}()

	require.Equal(t, http.StatusOK, resp.StatusCode, "OpenSearch document should exist")

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var osData struct {
		Source struct {
			ID          int    `json:"id"`
			Title       string `json:"title"`
			Status      string `json:"status"`
			PublishedAt string `json:"published_at"`
			Images      []struct {
				ID  int    `json:"id"`
				URL string `json:"url"`
			} `json:"images"`
		} `json:"_source"`
	}
	err = json.Unmarshal(body, &osData)
	require.NoError(t, err)

	// Assertions on OpenSearch data
	assert.Equal(t, testListingID, osData.Source.ID, "IDs should match")
	assert.Equal(t, dbData.Title, osData.Source.Title, "Titles should match")
	assert.Equal(t, "active", osData.Source.Status, "Status should be active in OpenSearch")
	assert.NotEmpty(t, osData.Source.PublishedAt, "OpenSearch should have published_at")
	assert.Greater(t, len(osData.Source.Images), 0, "OpenSearch should have images")

	// Check data consistency
	assert.Equal(t, dbData.ImagesCount, len(osData.Source.Images), "Images count should match between DB and OpenSearch")
}

// TestListingImagesInDatabase checks that listing has images in microservice database
func TestListingImagesInDatabase(t *testing.T) {
	db, err := sqlx.Connect("postgres", listingsMicroDBDSN)
	require.NoError(t, err)
	defer func() {
		if err := db.Close(); err != nil {
			t.Logf("Failed to close database connection: %v", err)
		}
	}()

	var images []struct {
		ID           int    `db:"id"`
		ListingID    int    `db:"listing_id"`
		URL          string `db:"url"`
		ThumbnailURL string `db:"thumbnail_url"`
		DisplayOrder int    `db:"display_order"`
		IsPrimary    bool   `db:"is_primary"`
	}

	err = db.Select(&images, `
		SELECT id, listing_id, url, thumbnail_url, display_order, is_primary
		FROM listing_images
		WHERE listing_id = $1
		ORDER BY display_order
	`, testListingID)
	require.NoError(t, err)

	assert.NotEmpty(t, images, "Listing should have images in database")

	if len(images) > 0 {
		img := images[0]
		assert.Equal(t, testListingID, img.ListingID, "Image should belong to correct listing")
		assert.NotEmpty(t, img.URL, "Image should have URL")
		assert.Greater(t, img.DisplayOrder, 0, "Image should have display order")
	}
}

// TestBackendLogsForErrors checks backend logs for errors related to listings
func TestBackendLogsForErrors(t *testing.T) {
	// Read last 100 lines from backend log
	content, err := os.ReadFile("/tmp/backend.log")
	if err != nil {
		t.Skip("Backend log not available")
		return
	}

	logContent := string(content)

	// Check for error keywords
	errorKeywords := []string{
		"failed to get listing",
		"failed to load images",
		"error loading images",
		"panic",
	}

	for _, keyword := range errorKeywords {
		assert.NotContains(t, logContent, keyword, fmt.Sprintf("Backend log should not contain '%s'", keyword))
	}
}

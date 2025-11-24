//go:build integration
// +build integration

package indexer

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// getTestDB returns a database connection for integration testing
func getTestDB(t *testing.T) *sqlx.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Use correct credentials from .env
		dsn = "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable"
	}

	db, err := sqlx.Connect("postgres", dsn)
	require.NoError(t, err, "failed to connect to test database")

	return db
}

func TestAttributeIndexer_BuildAttributesForIndex_Integration(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	indexer := NewAttributeIndexer(db, logger)

	ctx := context.Background()

	// Test with a real listing that has attributes (ID 106)
	t.Run("real listing with attributes", func(t *testing.T) {
		attributes, searchable, filterable, err := indexer.BuildAttributesForIndex(ctx, 106)

		require.NoError(t, err)
		assert.NotEmpty(t, attributes, "should have attributes")
		assert.NotEmpty(t, searchable, "should have searchable text")
		assert.NotEmpty(t, filterable, "should have filterable data")

		t.Logf("Attributes count: %d", len(attributes))
		t.Logf("Searchable text: %s", searchable)
		t.Logf("Filterable data: %+v", filterable)
	})

	t.Run("listing without attributes", func(t *testing.T) {
		// Use a listing ID that likely has no attributes
		attributes, searchable, filterable, err := indexer.BuildAttributesForIndex(ctx, 999999)

		require.NoError(t, err)
		assert.Empty(t, attributes, "should have no attributes")
		assert.Empty(t, searchable, "should have empty searchable text")
		assert.Empty(t, filterable, "should have empty filterable data")
	})

	t.Run("non-existent listing", func(t *testing.T) {
		// This should not error - just return empty results
		attributes, searchable, filterable, err := indexer.BuildAttributesForIndex(ctx, -1)

		require.NoError(t, err)
		assert.Empty(t, attributes)
		assert.Empty(t, searchable)
		assert.Empty(t, filterable)
	})
}

func TestAttributeIndexer_UpdateListingAttributeCache_Integration(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	indexer := NewAttributeIndexer(db, logger)

	ctx := context.Background()

	t.Run("update cache for listing with attributes", func(t *testing.T) {
		listingID := int32(106)

		// Update cache
		err := indexer.UpdateListingAttributeCache(ctx, listingID)
		require.NoError(t, err)

		// Verify cache was created
		var count int
		err = db.QueryRowContext(ctx,
			"SELECT COUNT(*) FROM attribute_search_cache WHERE listing_id = $1",
			listingID).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 1, count, "cache entry should exist")

		// Verify data structure
		var attributesFlat, attributesFilterable json.RawMessage
		var attributesSearchable string

		err = db.QueryRowContext(ctx, `
			SELECT attributes_flat, attributes_searchable, attributes_filterable
			FROM attribute_search_cache WHERE listing_id = $1
		`, listingID).Scan(&attributesFlat, &attributesSearchable, &attributesFilterable)

		require.NoError(t, err)
		assert.NotEmpty(t, attributesFlat)
		assert.NotEmpty(t, attributesSearchable)
		assert.NotEmpty(t, attributesFilterable)

		// Verify JSON is valid
		var flatArray []AttributeForIndex
		err = json.Unmarshal(attributesFlat, &flatArray)
		require.NoError(t, err)
		assert.NotEmpty(t, flatArray)

		var filterableMap map[string]interface{}
		err = json.Unmarshal(attributesFilterable, &filterableMap)
		require.NoError(t, err)

		t.Logf("Flat attributes: %d items", len(flatArray))
		t.Logf("Searchable text length: %d chars", len(attributesSearchable))
		t.Logf("Filterable keys: %v", getKeys(filterableMap))
	})

	t.Run("upsert behavior - update existing cache", func(t *testing.T) {
		listingID := int32(106)

		// Update twice - should not create duplicate
		err := indexer.UpdateListingAttributeCache(ctx, listingID)
		require.NoError(t, err)

		err = indexer.UpdateListingAttributeCache(ctx, listingID)
		require.NoError(t, err)

		// Verify still only one entry
		var count int
		err = db.QueryRowContext(ctx,
			"SELECT COUNT(*) FROM attribute_search_cache WHERE listing_id = $1",
			listingID).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 1, count, "should still be only one entry")
	})
}

func TestAttributeIndexer_GetListingAttributeCache_Integration(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	indexer := NewAttributeIndexer(db, logger)

	ctx := context.Background()
	listingID := int32(106)

	// Ensure cache exists
	err := indexer.UpdateListingAttributeCache(ctx, listingID)
	require.NoError(t, err)

	// Retrieve cache
	attributes, err := indexer.GetListingAttributeCache(ctx, listingID)
	require.NoError(t, err)
	assert.NotEmpty(t, attributes)

	// Verify attribute structure
	for _, attr := range attributes {
		assert.NotZero(t, attr.ID)
		assert.NotEmpty(t, attr.Code)
		assert.NotEmpty(t, attr.Name)

		// At least one value should be set
		hasValue := attr.ValueText != nil || attr.ValueNumber != nil || attr.ValueBoolean != nil
		assert.True(t, hasValue, "attribute should have at least one value")
	}
}

func TestAttributeIndexer_DeleteListingAttributeCache_Integration(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	indexer := NewAttributeIndexer(db, logger)

	ctx := context.Background()
	listingID := int32(106)

	// Ensure cache exists
	err := indexer.UpdateListingAttributeCache(ctx, listingID)
	require.NoError(t, err)

	// Delete cache
	err = indexer.DeleteListingAttributeCache(ctx, listingID)
	require.NoError(t, err)

	// Verify deletion
	var count int
	err = db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM attribute_search_cache WHERE listing_id = $1",
		listingID).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count, "cache should be deleted")
}

func TestAttributeIndexer_PopulateAttributeSearchCache_Integration(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	indexer := NewAttributeIndexer(db, logger)

	ctx := context.Background()

	// Clear cache first
	_, err := db.ExecContext(ctx, "DELETE FROM attribute_search_cache")
	require.NoError(t, err)

	// Populate cache
	start := time.Now()
	err = indexer.PopulateAttributeSearchCache(ctx, 100)
	duration := time.Since(start)

	require.NoError(t, err)

	// Verify cache was populated
	var count int
	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM attribute_search_cache").Scan(&count)
	require.NoError(t, err)
	assert.Greater(t, count, 0, "cache should have entries")

	// Performance check: should be fast (target: <5ms per listing on average)
	avgTimePerListing := duration / time.Duration(count)
	t.Logf("Populated %d listings in %v (avg: %v per listing)", count, duration, avgTimePerListing)

	// Allow up to 10ms per listing for CI/slower systems
	assert.Less(t, avgTimePerListing, 10*time.Millisecond, "should be reasonably fast")
}

func TestAttributeIndexer_BulkUpdateCache_Integration(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	indexer := NewAttributeIndexer(db, logger)

	ctx := context.Background()

	t.Run("bulk update multiple listings", func(t *testing.T) {
		listingIDs := []int32{106, 200, 201}

		// Clear cache first
		for _, id := range listingIDs {
			_ = indexer.DeleteListingAttributeCache(ctx, id)
		}

		// Bulk update
		err := indexer.BulkUpdateCache(ctx, listingIDs)
		require.NoError(t, err)

		// Verify all were updated (check individually to avoid array conversion issues)
		for _, id := range listingIDs {
			var exists bool
			err = db.QueryRowContext(ctx,
				"SELECT EXISTS(SELECT 1 FROM attribute_search_cache WHERE listing_id = $1)",
				id).Scan(&exists)
			require.NoError(t, err)
			assert.True(t, exists, "listing %d should have cache entry", id)
		}
	})

	t.Run("empty list should not error", func(t *testing.T) {
		err := indexer.BulkUpdateCache(ctx, []int32{})
		require.NoError(t, err)
	})
}

func TestAttributeIndexer_CascadeDelete_Integration(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	indexer := NewAttributeIndexer(db, logger)

	ctx := context.Background()

	// Get a valid category ID first
	var categoryID int32
	err := db.QueryRowContext(ctx, "SELECT id FROM categories LIMIT 1").Scan(&categoryID)
	require.NoError(t, err, "need at least one category in test database")

	// Create a test listing (with all required NOT NULL fields)
	var testListingID int32
	err = db.QueryRowContext(ctx, `
		INSERT INTO listings (title, description, category_id, status, visibility, user_id, price, currency)
		VALUES ('Test Listing', 'Test Description', $1, 'active', 'public', 1, 100.00, 'USD')
		RETURNING id
	`, categoryID).Scan(&testListingID)
	require.NoError(t, err)
	defer func() {
		// Cleanup
		_, _ = db.ExecContext(ctx, "DELETE FROM listings WHERE id = $1", testListingID)
	}()

	// Create cache entry
	err = indexer.UpdateListingAttributeCache(ctx, testListingID)
	require.NoError(t, err)

	// Verify cache exists
	var count int
	err = db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM attribute_search_cache WHERE listing_id = $1",
		testListingID).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)

	// Delete listing - should cascade to cache
	_, err = db.ExecContext(ctx, "DELETE FROM listings WHERE id = $1", testListingID)
	require.NoError(t, err)

	// Verify cache was deleted via CASCADE
	err = db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM attribute_search_cache WHERE listing_id = $1",
		testListingID).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count, "cache should be deleted via CASCADE")
}

// Helper function
func getKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

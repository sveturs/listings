//go:build integration
// +build integration

package indexer

import (
	"context"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

// BenchmarkUpdateListingAttributeCache measures cache update performance
func BenchmarkUpdateListingAttributeCache(b *testing.B) {
	db := getBenchDB(b)
	defer db.Close()

	logger := zerolog.New(os.Stdout).Level(zerolog.Disabled)
	indexer := NewAttributeIndexer(db, logger)

	ctx := context.Background()
	listingID := int32(106) // Listing with known attributes

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := indexer.UpdateListingAttributeCache(ctx, listingID)
		if err != nil {
			b.Fatalf("failed to update cache: %v", err)
		}
	}
}

// BenchmarkBuildAttributesForIndex measures attribute building performance
func BenchmarkBuildAttributesForIndex(b *testing.B) {
	db := getBenchDB(b)
	defer db.Close()

	logger := zerolog.New(os.Stdout).Level(zerolog.Disabled)
	indexer := NewAttributeIndexer(db, logger)

	ctx := context.Background()
	listingID := int32(106)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _, err := indexer.BuildAttributesForIndex(ctx, listingID)
		if err != nil {
			b.Fatalf("failed to build attributes: %v", err)
		}
	}
}

// BenchmarkGetListingAttributeCache measures cache retrieval performance
func BenchmarkGetListingAttributeCache(b *testing.B) {
	db := getBenchDB(b)
	defer db.Close()

	logger := zerolog.New(os.Stdout).Level(zerolog.Disabled)
	indexer := NewAttributeIndexer(db, logger)

	ctx := context.Background()
	listingID := int32(106)

	// Ensure cache exists
	_ = indexer.UpdateListingAttributeCache(ctx, listingID)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := indexer.GetListingAttributeCache(ctx, listingID)
		if err != nil {
			b.Fatalf("failed to get cache: %v", err)
		}
	}
}

// BenchmarkPopulateAttributeSearchCache measures full population performance
func BenchmarkPopulateAttributeSearchCache(b *testing.B) {
	db := getBenchDB(b)
	defer db.Close()

	logger := zerolog.New(os.Stdout).Level(zerolog.Disabled)
	indexer := NewAttributeIndexer(db, logger)

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Clear cache first
		_, _ = db.ExecContext(ctx, "DELETE FROM attribute_search_cache")

		err := indexer.PopulateAttributeSearchCache(ctx, 100)
		if err != nil {
			b.Fatalf("failed to populate cache: %v", err)
		}
	}
}

// getBenchDB returns database connection for benchmarks
func getBenchDB(b *testing.B) *sqlx.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable"
	}

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		b.Fatalf("failed to connect to test database: %v", err)
	}

	return db
}

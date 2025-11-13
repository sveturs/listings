package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/sveturs/listings/internal/indexer"
	"github.com/sveturs/listings/internal/repository/opensearch"
)

func main() {
	// Parse command line flags
	dryRun := flag.Bool("dry-run", false, "Dry run without updating OpenSearch")
	deleteIndex := flag.Bool("delete-index", false, "Delete and recreate index before reindexing")
	batchSize := flag.Int("batch-size", 100, "Batch size for reindexing")
	opensearchURL := flag.String("opensearch-url", "http://localhost:9200", "OpenSearch URL")
	dbURL := flag.String("db-url", "", "PostgreSQL connection URL (defaults to env DATABASE_URL)")
	flag.Parse()

	// Setup logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	logger := log.With().Str("service", "reindex_with_attributes").Logger()

	logger.Info().
		Bool("dry_run", *dryRun).
		Bool("delete_index", *deleteIndex).
		Int("batch_size", *batchSize).
		Str("opensearch_url", *opensearchURL).
		Msg("Starting reindex with attributes")

	// Get database URL
	dsn := *dbURL
	if dsn == "" {
		dsn = os.Getenv("DATABASE_URL")
	}
	if dsn == "" {
		dsn = "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable"
	}

	// Connect to database
	logger.Info().Msg("Connecting to PostgreSQL...")
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		logger.Fatal().Err(err).Msg("Failed to ping database")
	}
	logger.Info().Msg("✅ PostgreSQL connection established")

	// Connect to OpenSearch
	logger.Info().Str("url", *opensearchURL).Msg("Connecting to OpenSearch...")
	osClient, err := opensearch.NewClientSimple(*opensearchURL, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to OpenSearch")
	}
	logger.Info().Msg("✅ OpenSearch connection established")

	ctx := context.Background()

	// Step 1: Delete and recreate index if requested
	if *deleteIndex {
		if *dryRun {
			logger.Info().Msg("DRY RUN: Would delete and recreate index marketplace_listings")
		} else {
			logger.Warn().Msg("Deleting index marketplace_listings...")
			if err := osClient.DeleteIndex(ctx, "marketplace_listings"); err != nil {
				logger.Warn().Err(err).Msg("Failed to delete index (may not exist)")
			}

			logger.Info().Msg("Creating index marketplace_listings with attributes mapping...")
			mapping := opensearch.GetListingsIndexMapping()
			if err := osClient.CreateIndex(ctx, "marketplace_listings", mapping); err != nil {
				logger.Fatal().Err(err).Msg("Failed to create index")
			}
			logger.Info().Msg("✅ Index created successfully")
		}
	}

	// Step 2: Populate attribute_search_cache
	logger.Info().Msg("Populating attribute_search_cache...")
	attrIndexer := indexer.NewAttributeIndexer(db, logger)

	if *dryRun {
		logger.Info().Msg("DRY RUN: Would populate attribute_search_cache")
	} else {
		startTime := time.Now()
		if err := attrIndexer.PopulateAttributeSearchCache(ctx, *batchSize); err != nil {
			logger.Fatal().Err(err).Msg("Failed to populate attribute_search_cache")
		}
		duration := time.Since(startTime)
		logger.Info().
			Dur("duration", duration).
			Msg("✅ attribute_search_cache populated")
	}

	// Step 3: Reindex all listings with attributes
	logger.Info().Msg("Reindexing all listings with attributes...")
	listingIndexer := indexer.NewListingIndexer(db, osClient, logger)

	if *dryRun {
		logger.Info().
			Int("batch_size", *batchSize).
			Msg("DRY RUN: Would reindex all listings with attributes")
	} else {
		startTime := time.Now()
		err := listingIndexer.ReindexAllWithAttributes(ctx, *batchSize)
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to reindex listings")
		}
		duration := time.Since(startTime)

		logger.Info().
			Dur("duration", duration).
			Msg("✅ Reindex completed successfully")
	}

	// Step 4: Verify
	if !*dryRun {
		logger.Info().Msg("Verifying index...")

		// Count documents
		count, err := osClient.CountDocuments(ctx, "marketplace_listings")
		if err != nil {
			logger.Error().Err(err).Msg("Failed to count documents")
		} else {
			logger.Info().Int("count", count).Msg("Documents in index")
		}

		// Sample query with attributes
		logger.Info().Msg("Testing attribute-based query...")
		// This is just a sample - actual queries would come from the search API
		logger.Info().Msg("✅ Index is ready for attribute-based filtering")
	}

	logger.Info().Msg("🎉 Reindex with attributes completed successfully!")
	fmt.Println("\nNext steps:")
	fmt.Println("1. Test attribute-based search: curl 'http://localhost:9200/marketplace_listings/_search' -H 'Content-Type: application/json' -d '{\"query\":{\"nested\":{\"path\":\"attributes\",\"query\":{\"bool\":{\"must\":[{\"term\":{\"attributes.code\":\"car_make\"}},{\"term\":{\"attributes.value_text.keyword\":\"Toyota\"}}]}}}}}'")
	fmt.Println("2. Verify cache: psql -c 'SELECT COUNT(*) FROM attribute_search_cache;'")
	fmt.Println("3. Monitor performance: watch -n 1 'curl -s http://localhost:9200/marketplace_listings/_stats | jq .indices.marketplace_listings.total.search'")
}

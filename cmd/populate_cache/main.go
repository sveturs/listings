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

	"github.com/vondi-global/listings/internal/config"
	"github.com/vondi-global/listings/internal/indexer"
)

func main() {
	// Parse flags
	batchSize := flag.Int("batch", 100, "Batch size for processing listings")
	dryRun := flag.Bool("dry-run", false, "Dry run mode (don't write to database)")
	flag.Parse()

	// Setup logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load configuration")
	}

	log.Info().
		Str("env", cfg.App.Env).
		Int("batch_size", *batchSize).
		Bool("dry_run", *dryRun).
		Msg("starting attribute search cache population")

	// Connect to database
	db, err := sqlx.Connect("postgres", cfg.DB.DSN())
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer db.Close()

	// Configure connection pool
	db.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	db.SetMaxIdleConns(cfg.DB.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.DB.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.DB.ConnMaxIdleTime)

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatal().Err(err).Msg("failed to ping database")
	}

	log.Info().Msg("database connection established")

	// Create indexer
	attrIndexer := indexer.NewAttributeIndexer(db, log.Logger)

	ctx := context.Background()

	if *dryRun {
		log.Info().Msg("dry run mode: showing stats only")

		// Get count of listings with attributes
		var count int
		err := db.QueryRowContext(ctx, "SELECT COUNT(DISTINCT listing_id) FROM listing_attribute_values").Scan(&count)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to count listings")
		}

		log.Info().Int("total_listings_with_attributes", count).Msg("dry run completed")
		return
	}

	// Populate cache
	startTime := time.Now()

	err = attrIndexer.PopulateAttributeSearchCache(ctx, *batchSize)
	if err != nil {
		log.Error().Err(err).Msg("failed to populate attribute search cache")
		os.Exit(1)
	}

	elapsed := time.Since(startTime)

	// Verify results
	var cacheCount int
	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM attribute_search_cache").Scan(&cacheCount)
	if err != nil {
		log.Error().Err(err).Msg("failed to count cache entries")
	} else {
		log.Info().Int("cache_entries", cacheCount).Msg("cache entries created")
	}

	log.Info().
		Dur("elapsed", elapsed).
		Msg("attribute search cache population completed successfully")

	fmt.Println("\n=== Summary ===")
	fmt.Printf("Batch size: %d\n", *batchSize)
	fmt.Printf("Cache entries: %d\n", cacheCount)
	fmt.Printf("Time elapsed: %s\n", elapsed)
	fmt.Println("\nNext steps:")
	fmt.Println("1. Update OpenSearch mapping to support attributes")
	fmt.Println("2. Reindex listings with attributes data")
	fmt.Println("3. Test search and filtering")
}

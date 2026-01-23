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

	"github.com/vondi-global/listings/internal/config"
	"github.com/vondi-global/listings/internal/indexer"
	"github.com/vondi-global/listings/internal/reindexer"
	"github.com/vondi-global/listings/internal/repository/opensearch"
)

func main() {
	// Parse command line flags
	verify := flag.Bool("verify", false, "Only verify current index without reindexing")
	rollback := flag.String("rollback", "", "Rollback to specified index version (v1 or v2)")
	batchSize := flag.Int("batch", 500, "Batch size for reindexing")
	flag.Parse()

	// Setup logger
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	logger.Info().Msg("🚀 Vondi Listings Reindex Tool")

	// Load configuration
	if err := config.LoadEnv(); err != nil {
		logger.Warn().Err(err).Msg("failed to load .env file (continuing with env vars)")
	}

	cfg, err := config.Load()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load configuration")
	}

	logger.Info().
		Str("env", cfg.App.Env).
		Str("db", cfg.DB.Name).
		Str("opensearch", cfg.Search.Addresses[0]).
		Msg("configuration loaded")

	// Connect to database
	db, err := sqlx.Connect("postgres", cfg.DB.DSN())
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer db.Close()

	logger.Info().Msg("✅ database connected")

	// Connect to OpenSearch
	osClient, err := opensearch.NewClient(
		cfg.Search.Addresses,
		cfg.Search.Username,
		cfg.Search.Password,
		cfg.Search.Index,
		logger,
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to OpenSearch")
	}

	logger.Info().Msg("✅ OpenSearch connected")

	// Create indexer
	listingIndexer := indexer.NewListingIndexer(db, osClient, logger)

	// Create reindex manager
	manager := reindexer.NewReindexManager(osClient, listingIndexer, logger)

	ctx := context.Background()

	// Handle verify mode
	if *verify {
		handleVerify(ctx, manager, logger)
		return
	}

	// Handle rollback mode
	if *rollback != "" {
		handleRollback(ctx, manager, *rollback, logger)
		return
	}

	// Handle full reindex (blue-green)
	handleReindex(ctx, manager, *batchSize, logger)
}

func handleVerify(ctx context.Context, manager *reindexer.ReindexManager, logger zerolog.Logger) {
	logger.Info().Msg("🔍 Verifying current index...")

	// Get current version
	version, err := manager.GetCurrentIndexVersion(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to get current version")
	}

	indexName := fmt.Sprintf("marketplace_listings_%s", version)
	logger.Info().Str("index", indexName).Msg("verifying index")

	result, err := manager.VerifyReindex(ctx, indexName)
	if err != nil {
		logger.Fatal().Err(err).Msg("verification failed")
	}

	// Print results
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📊 VERIFICATION RESULTS")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Index:          %s\n", indexName)
	fmt.Printf("Valid:          %v\n", result.Valid)
	fmt.Printf("Total Docs:     %d\n", result.TotalDocs)
	fmt.Printf("Expected Docs:  %d\n", result.ExpectedDocs)
	fmt.Printf("Mismatch Count: %d\n", result.MismatchedCount)
	fmt.Println()

	if len(result.FieldCoverage) > 0 {
		fmt.Println("Field Coverage:")
		for field, coverage := range result.FieldCoverage {
			status := "✅"
			if coverage < 99.0 {
				status = "❌"
			}
			fmt.Printf("  %s %s: %.2f%%\n", status, field, coverage)
		}
		fmt.Println()
	}

	if len(result.SampleErrors) > 0 {
		fmt.Println("⚠️  Errors:")
		for i, err := range result.SampleErrors {
			fmt.Printf("  %d. %s\n", i+1, err)
		}
		fmt.Println()
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if result.Valid {
		logger.Info().Msg("✅ Verification PASSED")
		os.Exit(0)
	} else {
		logger.Error().Msg("❌ Verification FAILED")
		os.Exit(1)
	}
}

func handleRollback(ctx context.Context, manager *reindexer.ReindexManager, targetVersion string, logger zerolog.Logger) {
	if targetVersion != "v1" && targetVersion != "v2" {
		logger.Fatal().Str("version", targetVersion).Msg("invalid rollback version - must be v1 or v2")
	}

	logger.Warn().Str("target", targetVersion).Msg("🔄 Rolling back to previous index...")

	targetIndex := fmt.Sprintf("marketplace_listings_%s", targetVersion)

	if err := manager.RollbackToOldIndex(ctx, targetIndex); err != nil {
		logger.Fatal().Err(err).Msg("rollback failed")
	}

	logger.Info().Str("index", targetIndex).Msg("✅ Rollback completed successfully")
}

func handleReindex(ctx context.Context, manager *reindexer.ReindexManager, batchSize int, logger zerolog.Logger) {
	logger.Info().
		Int("batch_size", batchSize).
		Msg("🔄 Starting blue-green reindex...")

	startTime := time.Now()

	if err := manager.StartBlueGreenReindex(ctx, batchSize); err != nil {
		logger.Fatal().Err(err).Msg("❌ Reindex failed")
	}

	elapsed := time.Since(startTime)

	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("✅ REINDEX COMPLETED")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Duration:  %.2f seconds (%.2f minutes)\n", elapsed.Seconds(), elapsed.Minutes())
	fmt.Println()
	fmt.Println("💡 Old index kept for 24h rollback period")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	logger.Info().
		Str("duration", elapsed.String()).
		Msg("✅ Blue-green reindex completed successfully")
}

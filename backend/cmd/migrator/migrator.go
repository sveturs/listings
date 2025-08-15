package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"backend/internal/app/migrator"
	"backend/internal/logger"
)

// Build information set by ldflags
var (
	gitCommit = "unknown"
	buildTime = "unknown"
)

// Command line flags
var (
	withFixtures = flag.Bool("with-fixtures", false, "Run fixtures after migrations")
	onlyFixtures = flag.Bool("only-fixtures", false, "Run only fixtures without migrations")
)

func main() {
	// Parse command line flags
	flag.Parse()

	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}

	// Initialize logger
	if err := logger.Init(os.Getenv("APP_MODE"), os.Getenv("LOG_LEVEL")); err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize logger")
	}

	// Log build information
	logger.Info().
		Str("gitCommit", gitCommit).
		Str("buildTime", buildTime).
		Msg("Migrator version")

	// Load configuration from environment
	config := migrator.LoadConfigFromEnv()

	// Override direction from command line argument for backward compatibility
	args := flag.Args()
	if len(args) > 0 {
		config.Direction = args[0]
	}

	// Set flags from command line
	config.WithFixtures = *withFixtures
	config.OnlyFixtures = *onlyFixtures

	// Run migrations using the unified logic
	if err := migrator.RunMigrations(config); err != nil {
		logger.Fatal().Err(err).Msg("Migration failed")
	}

	logger.Info().Msg("Migration completed successfully")
}

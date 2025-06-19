package main

import (
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

func main() {
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

	// Get migration settings from environment variables
	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	if migrationsPath == "" {
		migrationsPath = "migrations" // default path
	}

	direction := os.Getenv("MIGRATION_DIRECTION")
	if direction == "" {
		// Fallback to command line argument for backward compatibility
		if len(os.Args) > 1 {
			direction = os.Args[1]
		} else {
			direction = "up" // default direction
		}
	}

	targetVersion := os.Getenv("MIGRATION_TARGET")

	logger.Info().
		Str("direction", direction).
		Str("path", migrationsPath).
		Str("target", targetVersion).
		Msg("Starting migration")

	// Load configuration
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		logger.Fatal().Str("DATABASE_URL", dbURL).Msgf("Error loading env")
	}

	migrtr := migrator.NewMigrator(migrationsPath, dbURL)

	// Execute the requested command
	switch direction {
	case "down":
		if targetVersion != "" {
			// Down to specific version
			if err := migrtr.DownTo(targetVersion); err != nil {
				logger.Fatal().Err(err).Msgf("Error running migrations down to version %s", targetVersion)
			}
		} else {
			// Down all migrations
			if err := migrtr.Down(); err != nil {
				logger.Fatal().Err(err).Msgf("Error running migrations down")
			}
		}
	case "up":
		if targetVersion != "" {
			// Up to specific version
			if err := migrtr.UpTo(targetVersion); err != nil {
				logger.Fatal().Err(err).Msgf("Error running migrations up to version %s", targetVersion)
			}
		} else {
			// Up all pending migrations
			if err := migrtr.Up(); err != nil {
				logger.Fatal().Err(err).Msgf("Error running migrations up")
			}
		}
	default:
		logger.Fatal().Msgf("Unknown migration direction: %s", direction)
	}
}

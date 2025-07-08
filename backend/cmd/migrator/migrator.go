package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

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

	// Get migration settings from environment variables
	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	if migrationsPath == "" {
		migrationsPath = "migrations" // default path
	}

	// Get fixtures settings from environment variables
	fixturesPath := os.Getenv("MIGRATIONS_FIXTURES_PATH")
	if fixturesPath == "" {
		fixturesPath = "fixtures" // default path
	}

	fixturesTable := os.Getenv("MIGRATIONS_FIXTURES_TABLE")
	if fixturesTable == "" {
		fixturesTable = "schema_fixtures" // default table name
	}

	direction := os.Getenv("MIGRATION_DIRECTION")
	if direction == "" {
		// Fallback to command line argument for backward compatibility
		args := flag.Args()
		if len(args) > 0 {
			direction = args[0]
		} else {
			direction = "up" // default direction
		}
	}

	targetVersion := os.Getenv("MIGRATION_TARGET")

	logger.Info().
		Str("direction", direction).
		Str("migrationsPath", migrationsPath).
		Str("fixturesPath", fixturesPath).
		Str("fixturesTable", fixturesTable).
		Str("target", targetVersion).
		Bool("withFixtures", *withFixtures).
		Bool("onlyFixtures", *onlyFixtures).
		Msg("Starting migration")

	// Load configuration
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		logger.Fatal().Str("DATABASE_URL", dbURL).Msgf("Error loading env")
	}

	// Check flags conflict
	if *withFixtures && *onlyFixtures {
		logger.Fatal().Msg("Cannot use --with-fixtures and --only-fixtures together")
	}

	// Structure to hold migrator info
	type migratorInfo struct {
		migr    *migrator.Migrator
		name    string
		errName string
	}

	// Create slice of migrators based on flags
	var migrators []migratorInfo

	// Add schema migrations unless --only-fixtures is specified
	if !*onlyFixtures {
		migrators = append(migrators, migratorInfo{
			migr:    migrator.NewMigrator(migrationsPath, dbURL),
			name:    "migrations",
			errName: "migrations",
		})
	}

	// Add fixtures if requested (--with-fixtures or --only-fixtures)
	if *withFixtures || *onlyFixtures {
		logger.Info().Msg("Running fixtures")
		// Check if dbURL already has query parameters
		var fixturesDSN string
		if strings.Contains(dbURL, "?") {
			fixturesDSN = dbURL + "&x-migrations-table=" + fixturesTable
		} else {
			fixturesDSN = dbURL + "?x-migrations-table=" + fixturesTable
		}
		migrators = append(migrators, migratorInfo{
			migr:    migrator.NewMigrator(fixturesPath, fixturesDSN),
			name:    "fixtures",
			errName: "fixtures",
		})
	}

	// Execute migration for all migrators
	for _, m := range migrators {
		switch direction {
		case "down":
			if targetVersion != "" {
				// Down to specific version
				if err := m.migr.DownTo(targetVersion); err != nil {
					logger.Fatal().Err(err).Msgf("Error running %s down to version %s", m.errName, targetVersion)
				}
			} else {
				// Down all
				if err := m.migr.Down(); err != nil {
					logger.Fatal().Err(err).Msgf("Error running %s down", m.errName)
				}
			}
		case "up":
			if targetVersion != "" {
				// Up to specific version
				if err := m.migr.UpTo(targetVersion); err != nil {
					logger.Fatal().Err(err).Msgf("Error running %s up to version %s", m.errName, targetVersion)
				}
			} else {
				// Up all pending
				if err := m.migr.Up(); err != nil {
					logger.Fatal().Err(err).Msgf("Error running %s up", m.errName)
				}
			}
		default:
			logger.Fatal().Msgf("Unknown migration direction: %s", direction)
		}
	}
}

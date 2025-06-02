package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"backend/internal/app/migrator"
	"backend/internal/config"
	"backend/internal/logger"
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

	arg := "up"
	if len(os.Args) > 1 {
		arg = os.Args[1]
	}

	logger.Info().Msg("Starting migration")

	// Load configuration
	cfg, err := config.NewConfig()
	if err != nil {
		logger.Fatal().Err(err).Msgf("Error loading config")
	}

	migrtr := migrator.NewMigrator("migrations", cfg.DatabaseURL)

	// Execute the requested command
	switch arg {
	case "down":
		if err := migrtr.Down(); err != nil {
			logger.Fatal().Err(err).Msgf("Error running migrations down")
		}
	default:
		if err := migrtr.Up(); err != nil {
			logger.Fatal().Err(err).Msgf("Error running migrations up")
		}
	}
}

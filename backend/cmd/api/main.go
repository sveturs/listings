//	@title			Sve Tu Marketplace API
//	@version		1.0
//	@description	API для платформы Sve Tu - маркетплейс объявлений
//	@termsOfService	https://svetu.rs/terms

//	@contact.name	API Support
//	@contact.url	https://svetu.rs/support
//	@contact.email	support@svetu.rs

//	@license.name	MIT
//	@license.url	https://opensource.org/licenses/MIT

//	@host		localhost:3000
//	@BasePath	/api/v1

//	@securityDefinitions.apikey	Bearer
//	@in							header
//	@name						Authorization
//	@description				Bearer token для авторизации

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"backend/internal/app/migrator"
	"backend/internal/config"
	"backend/internal/logger"
	"backend/internal/server"
	"backend/internal/version"
)

// Build information set by ldflags
var (
	gitCommit = "unknown"
	buildTime = "unknown"
)

func main() {
	// Загрузка конфигурации из файла окружения
	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = ".env"
	}

	if err := godotenv.Load(envFile); err != nil {
		logger.Info().Str("envFile", envFile).Msgf("Warning: Could not load .env file: %s", envFile)
	}

	// Инициализация конфигурации
	cfg, err := config.NewConfig()
	if err != nil {
		logger.Fatal().Err(err).Msgf("Failed to load config: %v", err)
	}

	// Initialize logger с конфигурацией
	if err := logger.Init(cfg.Environment, cfg.LogLevel, version.GetVersion()); err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize logger")
	}
	logger.Info().
		Str("gitCommit", gitCommit).
		Str("buildTime", buildTime).
		Any("config", cfg).
		Msg("Config loaded successfully")

	// Выполнение миграций при старте API (если включено)
	if cfg.MigrationsOnAPI != "off" {
		logger.Info().
			Str("migrationsMode", cfg.MigrationsOnAPI).
			Msg("Running migrations on API startup")

		switch cfg.MigrationsOnAPI {
		case "schema":
			if err := migrator.RunMigrationsSchema(cfg.DatabaseURL); err != nil {
				logger.Fatal().Err(err).Msg("Failed to run schema migrations")
			}
		case "full":
			if err := migrator.RunMigrationsFull(cfg.DatabaseURL); err != nil {
				logger.Fatal().Err(err).Msg("Failed to run full migrations")
			}
		default:
			logger.Warn().
				Str("migrationsMode", cfg.MigrationsOnAPI).
				Msg("Unknown migrations mode, skipping")
		}

		logger.Info().Msg("Migrations completed successfully")
	}

	// Создание и запуск сервера
	// Удаляем второй аргумент fileStorage, так как NewServer инициализирует его внутри себя
	ctx := context.Background()
	srv, err := server.NewServer(ctx, cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create server")
	}

	// Graceful shutdown
	go func() {
		if err := srv.Start(); err != nil {
			logger.Error().Err(err).Msg("Server error")
		}
	}()

	// Ожидание сигнала для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal().Err(err).Msgf("Server forced to shutdown")
	}

	logger.Info().Msg("Server exited properly")
}

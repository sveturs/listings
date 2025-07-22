// backend/cmd/reindex/main.go
package main

import (
	"context"
	"flag"
	"os"
	"strings"
	"time"

	"backend/internal/config"
	"backend/internal/logger"
	"backend/internal/storage/filestorage"
	"backend/internal/storage/opensearch"
	"backend/internal/storage/postgres"

	"github.com/joho/godotenv"
)

// Build information set by ldflags
var (
	gitCommit = "unknown"
	buildTime = "unknown"
)

func main() {
	// Загрузка переменных окружения
	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = ".env"
	}

	if err := godotenv.Load(envFile); err != nil {
		logger.Info().Msgf("Warning: Could not load .env file: %s", err)
	}
	// Initialize logger
	if err := logger.Init(os.Getenv("APP_MODE"), os.Getenv("LOG_LEVEL")); err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize logger")
	}

	// Чтение конфигурации
	cfg, err := config.NewConfig()
	if err != nil {
		logger.Fatal().Err(err).Msgf("Failed to load config: %v", err)
	}
	logger.Info().
		Str("gitCommit", gitCommit).
		Str("buildTime", buildTime).
		Any("config", cfg).
		Msg("Config loaded successfully")

	// Парсинг аргументов командной строки
	entityTypeFlag := flag.String("type", "listings", "Type of entities to reindex (listings, reviews)")
	batchSizeFlag := flag.Int("batch-size", 500, "Batch size for reindexing")
	flag.Parse()

	// Подключение к OpenSearch
	var osClient *opensearch.OpenSearchClient
	if cfg.OpenSearch.URL != "" {
		osClient, err = opensearch.NewOpenSearchClient(opensearch.Config{
			URL:      cfg.OpenSearch.URL,
			Username: cfg.OpenSearch.Username,
			Password: cfg.OpenSearch.Password,
		})
		if err != nil {
			logger.Info().Msgf("Error connecting to OpenSearch: %v", err)
		}
	}

	// Инициализация файлового хранилища
	fileStorage, err := filestorage.NewFileStorage(cfg.FileStorage)
	if err != nil {
		logger.Info().Msgf("Warning: Failed to initialize file storage: %v. Proceeding without file storage.", err)
		// Не прерываем выполнение программы, так как для индексации это не критично
	}

	// Подключение к базе данных
	db, err := postgres.NewDatabase(cfg.DatabaseURL, osClient, cfg.OpenSearch.MarketplaceIndex, fileStorage, cfg.SearchWeights)
	if err != nil {
		logger.Fatal().Err(err).Msgf("Error connecting to database: %v", err)
	}
	defer func() {
		db.Close()
	}()

	// Проверка подключения к БД
	if err := db.Ping(context.Background()); err != nil {
		logger.Fatal().Err(err).Msgf("Error pinging database: %v", err)
	}

	logger.Info().Msgf("Starting reindexing of %s with batch size %d", *entityTypeFlag, *batchSizeFlag)
	start := time.Now()

	// Выполнение операции реиндексации в зависимости от типа сущности
	switch strings.ToLower(*entityTypeFlag) {
	case "listings", "marketplace":
		if err := db.ReindexAllListings(context.Background()); err != nil {
			logger.Fatal().Err(err).Msgf("Error reindexing listings: %v", err)
		}
	// Можно добавить другие типы реиндексации здесь
	default:
		logger.Fatal().Err(err).Msgf("Unknown entity type: %s", *entityTypeFlag)
	}

	logger.Info().Msgf("Reindexing completed in %v", time.Since(start))
}

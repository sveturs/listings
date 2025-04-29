// backend/cmd/reindex/main.go
package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strings"
	"time"

	"backend/internal/config"
	"backend/internal/storage/filestorage"
	"backend/internal/storage/opensearch"
	"backend/internal/storage/postgres"

	"github.com/joho/godotenv"
)

func main() {
	// Загрузка переменных окружения
	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = ".env"
	}

	if err := godotenv.Load(envFile); err != nil {
		log.Printf("Warning: Could not load .env file: %s", err)
	}

	// Чтение конфигурации
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

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
			log.Printf("Error connecting to OpenSearch: %v", err)
		}
	}

	// Инициализация файлового хранилища
	fileStorage, err := filestorage.NewFileStorage(cfg.FileStorage)
	if err != nil {
		log.Printf("Warning: Failed to initialize file storage: %v. Proceeding without file storage.", err)
		// Не прерываем выполнение программы, так как для индексации это не критично
	}

	// Подключение к базе данных
	db, err := postgres.NewDatabase(cfg.DatabaseURL, osClient, cfg.OpenSearch.MarketplaceIndex, fileStorage)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Проверка подключения к БД
	if err := db.Ping(context.Background()); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	log.Printf("Starting reindexing of %s with batch size %d", *entityTypeFlag, *batchSizeFlag)
	start := time.Now()

	// Выполнение операции реиндексации в зависимости от типа сущности
	switch strings.ToLower(*entityTypeFlag) {
	case "listings", "marketplace":
		if err := db.ReindexAllListings(context.Background()); err != nil {
			log.Fatalf("Error reindexing listings: %v", err)
		}
	// Можно добавить другие типы реиндексации здесь
	default:
		log.Fatalf("Unknown entity type: %s", *entityTypeFlag)
	}

	log.Printf("Reindexing completed in %v", time.Since(start))
}
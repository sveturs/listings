package main

import (
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"

	"backend/internal/config"
	"backend/internal/logger"
	"backend/internal/proj/marketplace/storage/opensearch"
	"backend/internal/storage"
	osClient "backend/internal/storage/opensearch"
)

func main() {
	ctx := context.Background()
	
	// Загрузка конфигурации
	if err := godotenv.Load(".env"); err != nil {
		fmt.Printf("Warning: Could not load .env file\n")
	}
	
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	if err := logger.Init(cfg.Environment, cfg.LogLevel); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// Подключение к БД
	db, err := storage.NewPostgresConnection(cfg.Database)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()
	
	storageRepo := storage.NewStorage(db)
	
	// Подключение к OpenSearch
	osConfig := &osClient.Config{
		URL:      cfg.OpenSearch.URL,
		Username: cfg.OpenSearch.Username,
		Password: cfg.OpenSearch.Password,
	}
	
	client, err := osClient.NewClient(osConfig)
	if err != nil {
		log.Fatalf("Ошибка подключения к OpenSearch: %v", err)
	}
	
	// Создание репозитория OpenSearch
	repo := opensearch.NewRepository(client, "marketplace", storageRepo, &cfg.SearchWeights)
	
	// Тестируем загрузку переводов для объявления с ID 98
	testListingID := 98
	translations, err := repo.GetListingTranslationsFromDB(ctx, testListingID)
	if err != nil {
		log.Fatalf("Ошибка загрузки переводов: %v", err)
	}
	
	fmt.Printf("✅ Найдено %d переводов для объявления %d:\n", len(translations), testListingID)
	for _, t := range translations {
		fmt.Printf("  - %s.%s: %s\n", t.Language, t.FieldName, t.TranslatedText[:50]+"...")
	}
	
	// Тестируем извлечение поддерживаемых языков
	languages := repo.ExtractSupportedLanguages(translations)
	fmt.Printf("✅ Поддерживаемые языки: %v\n", languages)
}
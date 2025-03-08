// backend/cmd/reindex/main.go
package main

import (
    "context"
    "backend/internal/config"
    "backend/internal/storage/postgres"
    "backend/internal/storage/opensearch"
    "log"
    "os"
    "time"

    "github.com/joho/godotenv"
)

func main() {
    // Загрузка .env файла
    if err := godotenv.Load(); err != nil {
        log.Printf("Warning: Could not load .env file")
    }

    // Загружаем конфигурацию
    cfg, err := config.NewConfig()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    
    // Проверяем и устанавливаем URL для OpenSearch
    if cfg.OpenSearch.URL == "" {
        cfg.OpenSearch.URL = os.Getenv("OPENSEARCH_URL") // Читаем из переменной окружения
        if cfg.OpenSearch.URL == "" {
            cfg.OpenSearch.URL = "http://opensearch:9200" // или используем значение по умолчанию
        }
        log.Printf("Using OpenSearch URL: %s", cfg.OpenSearch.URL)
    }
    
    if cfg.OpenSearch.MarketplaceIndex == "" {
        cfg.OpenSearch.MarketplaceIndex = os.Getenv("OPENSEARCH_MARKETPLACE_INDEX")
        if cfg.OpenSearch.MarketplaceIndex == "" {
            cfg.OpenSearch.MarketplaceIndex = "marketplace"
        }
        log.Printf("Using OpenSearch index: %s", cfg.OpenSearch.MarketplaceIndex)
    }
    
    // Инициализируем клиент OpenSearch
    osClient, err := opensearch.NewOpenSearchClient(opensearch.Config{
        URL:      cfg.OpenSearch.URL,
        Username: cfg.OpenSearch.Username,
        Password: cfg.OpenSearch.Password,
    })
    if err != nil {
        log.Fatalf("Failed to create OpenSearch client: %v", err)
    }

    // Используем ваш DatabaseURL из конфигурации или переменной окружения
    dbUrl := cfg.DatabaseURL
    if dbUrl == "" {
        dbUrl = os.Getenv("DATABASE_URL")
    }

    // Инициализируем базу данных с OpenSearch
    db, err := postgres.NewDatabase(dbUrl, osClient, cfg.OpenSearch.MarketplaceIndex)
    if err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }
    defer db.Close()

    // Создаем контекст с таймаутом
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
    defer cancel()

    // Запускаем переиндексацию
    log.Println("Starting reindexing...")
    start := time.Now()

    if err := db.ReindexAllListings(ctx); err != nil {
        log.Fatalf("Error during reindexing: %v", err)
    }

    log.Printf("Reindexing completed in %v", time.Since(start))
}
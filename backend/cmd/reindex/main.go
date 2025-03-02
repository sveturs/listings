package main

import (
    "context"
    "backend/internal/config"
    "backend/internal/storage/postgres"
    "backend/internal/storage/opensearch"
    "log"
    "os"
    "time"
	"fmt"

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
        cfg.OpenSearch.URL = "http://localhost:9200" // или http://opensearch:9200 в Docker
        log.Printf("Using default OpenSearch URL: %s", cfg.OpenSearch.URL)
    }
    
    if cfg.OpenSearch.MarketplaceIndex == "" {
        cfg.OpenSearch.MarketplaceIndex = "marketplace"
        log.Printf("Using default OpenSearch index: %s", cfg.OpenSearch.MarketplaceIndex)
    }
    
    // Проверка DATABASE_URL - важно для работы вне Docker
    if cfg.DatabaseURL == "" {
        host := os.Getenv("DB_HOST")
        if host == "" {
            host = "localhost" // Установка локального хоста по умолчанию
        }
        
        cfg.DatabaseURL = fmt.Sprintf("postgres://postgres:password@%s:5432/hostel_db?sslmode=disable", host)
        log.Printf("Using default database URL: %s", cfg.DatabaseURL)
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

    // Инициализируем базу данных
    db, err := postgres.NewDatabase(cfg.DatabaseURL, osClient, cfg.OpenSearch.MarketplaceIndex)
    if err != nil {
        log.Fatalf("Failed to create database: %v", err)
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
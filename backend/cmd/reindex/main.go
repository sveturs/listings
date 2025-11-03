//go:build ignore
// +build ignore

// DEPRECATED: This tool uses outdated Database API (ReindexAllProducts method removed)
// Use reindex_unified.py script instead
package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"backend/internal/config"
	"backend/internal/logger"
	"backend/internal/storage/filestorage"
	"backend/internal/storage/opensearch"
	"backend/internal/storage/postgres"
	"backend/internal/version"
)

func main() {
	log.Println("Запуск переиндексации OpenSearch...")
	startTime := time.Now()

	if err := godotenv.Load(); err != nil {
		logger.Info().Msgf("Warning: Could not load .env file: %s", err)
	}
	// Initialize logger
	if err := logger.Init(os.Getenv("APP_MODE"), os.Getenv("LOG_LEVEL"), version.GetVersion()); err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize logger")
	}

	// Чтение конфигурации
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Создаем контекст
	ctx := context.Background()

	// Инициализируем OpenSearch
	osConfig := opensearch.Config{
		URL:      cfg.OpenSearch.URL,
		Username: cfg.OpenSearch.Username,
		Password: cfg.OpenSearch.Password,
	}
	osClient, err := opensearch.NewOpenSearchClient(osConfig)
	indexName := cfg.OpenSearch.C2CIndex
	if err != nil {
		log.Fatalf("Failed to create OpenSearch client: %v", err)
	}

	// Инициализируем файловое хранилище
	fileStorage, err := filestorage.NewFileStorage(ctx, cfg.FileStorage)
	if err != nil {
		log.Fatalf("Failed to create file storage: %v", err)
	}

	// Инициализируем базу данных
	db, err := postgres.NewDatabase(ctx, cfg.DatabaseURL, osClient, indexName, cfg.OpenSearch.B2CIndex, fileStorage, cfg.SearchWeights)
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}

	// Используем правильный метод реиндексации через storage layer
	log.Println("Запуск реиндексации marketplace listings...")
	if err := db.ReindexAllListings(ctx); err != nil {
		log.Fatalf("Error reindexing listings: %v", err)
	}

	log.Println("Запуск реиндексации b2c_stores...")
	if err := db.ReindexAllStorefronts(ctx); err != nil {
		log.Fatalf("Error reindexing b2c_stores: %v", err)
	}

	log.Println("Запуск реиндексации товаров витрин...")
	if err := db.ReindexAllProducts(ctx); err != nil {
		log.Fatalf("Error reindexing products: %v", err)
	}

	duration := time.Since(startTime)
	log.Printf("Переиндексация завершена успешно за %v", duration)
}

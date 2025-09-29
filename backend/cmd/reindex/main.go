package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"backend/internal/config"
	"backend/internal/storage/postgres"
	"backend/internal/storage/opensearch"
	"backend/internal/storage/filestorage"
)

func main() {
	log.Println("Запуск переиндексации OpenSearch...")
	startTime := time.Now()

	// Загружаем конфигурацию
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
	indexName := cfg.OpenSearch.MarketplaceIndex
	if err != nil {
		log.Fatalf("Failed to create OpenSearch client: %v", err)
	}

	// Инициализируем файловое хранилище
	fileStorage, err := filestorage.NewFileStorage(ctx, cfg.FileStorage)
	if err != nil {
		log.Fatalf("Failed to create file storage: %v", err)
	}

	// Инициализируем базу данных
	db, err := postgres.NewDatabase(ctx, cfg.DatabaseURL, osClient, indexName, fileStorage, cfg.SearchWeights)
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}

	// Параметр для фильтрации по категориям (автомобили)
	categoryFilter := ""
	if len(os.Args) > 1 && os.Args[1] == "cars" {
		categoryFilter = "1301" // ID категории автомобилей
		log.Println("Фильтрация по категории: Автомобили (ID 1301)")
	}

	// Получаем все объявления из базы данных
	rows, err := db.GetPool().Query(ctx, `
		SELECT
			ml.id,
			ml.title,
			ml.description,
			ml.price,
			ml.category_id,
			ml.location,
			ml.latitude,
			ml.longitude,
			ml.user_id,
			ml.created_at,
			ml.updated_at,
			ml.status,
			ml.attributes,
			mc.name as category_name
		FROM marketplace_listings ml
		LEFT JOIN marketplace_categories mc ON ml.category_id = mc.id
		WHERE ($1 = '' OR ml.category_id::text = $1)
	`, categoryFilter)
	if err != nil {
		log.Fatalf("Failed to query listings: %v", err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var listing struct {
			ID           int
			Title        string
			Description  *string
			Price        float64
			CategoryID   *int
			Location     *string
			Latitude     *float64
			Longitude    *float64
			UserID       int
			CreatedAt    time.Time
			UpdatedAt    time.Time
			Status       string
			Attributes   json.RawMessage
			CategoryName *string
		}

		err := rows.Scan(
			&listing.ID,
			&listing.Title,
			&listing.Description,
			&listing.Price,
			&listing.CategoryID,
			&listing.Location,
			&listing.Latitude,
			&listing.Longitude,
			&listing.UserID,
			&listing.CreatedAt,
			&listing.UpdatedAt,
			&listing.Status,
			&listing.Attributes,
			&listing.CategoryName,
		)
		if err != nil {
			log.Printf("Failed to scan listing: %v", err)
			continue
		}

		// Формируем атрибуты для поиска
		var searchAttributes []map[string]interface{}
		if len(listing.Attributes) > 0 && string(listing.Attributes) != "null" {
			var attrs map[string]interface{}
			if err := json.Unmarshal(listing.Attributes, &attrs); err == nil {
				for key, value := range attrs {
					searchAttr := map[string]interface{}{
						"key":   key,
						"value": fmt.Sprintf("%v", value),
					}
					searchAttributes = append(searchAttributes, searchAttr)
				}
			}
		}

		// Формируем документ для индексации
		doc := map[string]interface{}{
			"id":            listing.ID,
			"title":         listing.Title,
			"description":   listing.Description,
			"price":         listing.Price,
			"category_id":   listing.CategoryID,
			"category_name": listing.CategoryName,
			"location":      listing.Location,
			"latitude":      listing.Latitude,
			"longitude":     listing.Longitude,
			"user_id":       listing.UserID,
			"created_at":    listing.CreatedAt,
			"updated_at":    listing.UpdatedAt,
			"status":        listing.Status,
			"attributes":    searchAttributes,
		}

		// Индексируем в OpenSearch
		err = osClient.IndexDocument(ctx, indexName, fmt.Sprintf("%d", listing.ID), doc)
		if err != nil {
			log.Printf("Failed to index listing %d: %v", listing.ID, err)
			continue
		}

		count++
		if count%100 == 0 {
			log.Printf("Проиндексировано %d объявлений...", count)
		}
	}

	duration := time.Since(startTime)
	log.Printf("Переиндексация завершена успешно! Проиндексировано %d объявлений за %v", count, duration)
}
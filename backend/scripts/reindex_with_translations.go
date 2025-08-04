package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend/internal/config"
	"backend/internal/domain/models"
	"backend/internal/logger"
	"backend/internal/storage"
	marketplaceOS "backend/internal/proj/marketplace/storage/opensearch"
	osClient "backend/internal/storage/opensearch"
)

func main() {
	ctx := context.Background()
	
	// Настройка логгера
	logger.Init("info")
	
	log.Println("🌍 Начинаем реиндексацию с мультиязычной поддержкой...")
	
	// Загрузка конфигурации
	cfg := config.LoadConfig()
	
	// Подключение к базе данных
	db, err := storage.NewPostgresConnection(cfg.Database)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()
	
	storageRepo := storage.NewStorage(db)
	
	// Подключение к OpenSearch
	osConfig := &osClient.Config{
		Addresses: []string{cfg.OpenSearch.Host},
		Username:  cfg.OpenSearch.Username,
		Password:  cfg.OpenSearch.Password,
	}
	
	client, err := osClient.NewClient(osConfig)
	if err != nil {
		log.Fatalf("Ошибка подключения к OpenSearch: %v", err)
	}
	
	// Создание репозитория OpenSearch
	searchWeights := &config.SearchWeights{
		// Используем дефолтные веса
		OpenSearchBoosts: config.OpenSearchBoostWeights{
			Title:               5.0,
			TitleNgram:         2.0,
			Description:        1.5,
			TranslationTitle:   5.0,
			TranslationDesc:    2.0,
			AttributeTextValue: 4.0,
			AttributeDisplayValue: 4.0,
			AttributeTextValueKeyword: 6.0,
			AttributeGeneralBoost: 2.0,
			RealEstateAttributesCombined: 3.0,
			PropertyType: 4.0,
			RoomsText: 3.0,
			CarMake: 5.0,
			CarModel: 4.0,
		},
	}
	
	repo := marketplaceOS.NewRepository(client, "marketplace_v2", storageRepo, searchWeights)
	
	// Обработка сигналов для graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	// Получение всех активных листингов
	log.Println("📋 Получение списка всех объявлений...")
	
	query := `
		SELECT 
			id, title, description, price, condition, status, location, city, country,
			views_count, created_at, updated_at, show_on_map, original_language,
			category_id, user_id, average_rating, review_count, old_price,
			coordinates_lat, coordinates_lon, storefront_id
		FROM marketplace_listings 
		WHERE status = 'active'
		ORDER BY id
	`
	
	rows, err := storageRepo.Query(ctx, query)
	if err != nil {
		log.Fatalf("Ошибка получения листингов: %v", err)
	}
	defer rows.Close()
	
	var listings []*models.MarketplaceListing
	
	for rows.Next() {
		listing := &models.MarketplaceListing{}
		var coordsLat, coordsLon *float64
		
		err := rows.Scan(
			&listing.ID, &listing.Title, &listing.Description, &listing.Price,
			&listing.Condition, &listing.Status, &listing.Location, 
			&listing.City, &listing.Country, &listing.ViewsCount,
			&listing.CreatedAt, &listing.UpdatedAt, &listing.ShowOnMap,
			&listing.OriginalLanguage, &listing.CategoryID, &listing.UserID,
			&listing.AverageRating, &listing.ReviewCount, &listing.OldPrice,
			&coordsLat, &coordsLon, &listing.StorefrontID,
		)
		if err != nil {
			log.Printf("Ошибка сканирования листинга: %v", err)
			continue
		}
		
		// Установка координат
		if coordsLat != nil && coordsLon != nil {
			listing.CoordinatesLat = coordsLat
			listing.CoordinatesLon = coordsLon
		}
		
		listings = append(listings, listing)
	}
	
	if err = rows.Err(); err != nil {
		log.Fatalf("Ошибка итерации по листингам: %v", err)
	}
	
	log.Printf("📊 Найдено %d объявлений для реиндексации", len(listings))
	
	if len(listings) == 0 {
		log.Println("⚠️ Нет объявлений для индексации")
		return
	}
	
	// Загрузка атрибутов, переводов и других данных для каждого листинга
	log.Println("🔄 Загрузка дополнительных данных (атрибуты, переводы, изображения)...")
	
	for i, listing := range listings {
		select {
		case <-sigChan:
			log.Println("🛑 Получен сигнал остановки, завершаем работу...")
			return
		default:
		}
		
		// Загрузка атрибутов
		attributes, err := storageRepo.GetListingAttributes(ctx, listing.ID)
		if err == nil {
			listing.Attributes = attributes
		}
		
		// Загрузка переводов
		translations, err := storageRepo.GetTranslationsForEntity(ctx, "listing", listing.ID)
		if err == nil && len(translations) > 0 {
			transMap := make(models.TranslationMap)
			for _, t := range translations {
				if _, ok := transMap[t.Language]; !ok {
					transMap[t.Language] = make(map[string]string)
				}
				transMap[t.Language][t.FieldName] = t.TranslatedText
			}
			listing.Translations = transMap
		}
		
		// Загрузка изображений
		images, err := storageRepo.GetListingImages(ctx, listing.ID)
		if err == nil {
			listing.Images = images
		}
		
		// Прогресс
		if (i+1)%100 == 0 {
			log.Printf("📈 Обработано %d/%d объявлений (%.1f%%)", 
				i+1, len(listings), float64(i+1)/float64(len(listings))*100)
		}
	}
	
	log.Println("🚀 Начинаем массовую индексацию в OpenSearch...")
	
	// Индексация батчами по 50 объявлений
	batchSize := 50
	totalBatches := (len(listings) + batchSize - 1) / batchSize
	
	for i := 0; i < len(listings); i += batchSize {
		select {
		case <-sigChan:
			log.Println("🛑 Получен сигнал остановки, завершаем работу...")
			return
		default:
		}
		
		end := i + batchSize
		if end > len(listings) {
			end = len(listings)
		}
		
		batch := listings[i:end]
		batchNum := i/batchSize + 1
		
		log.Printf("📦 Индексируем батч %d/%d (%d объявлений)", 
			batchNum, totalBatches, len(batch))
		
		start := time.Now()
		err := repo.BulkIndexListings(ctx, batch)
		if err != nil {
			log.Printf("❌ Ошибка индексации батча %d: %v", batchNum, err)
			// Попробуем индексировать по одному
			log.Printf("🔄 Повторная индексация по одному для батча %d", batchNum)
			for _, listing := range batch {
				if err := repo.IndexListing(ctx, listing); err != nil {
					log.Printf("❌ Ошибка индексации объявления %d: %v", listing.ID, err)
				}
			}
		} else {
			duration := time.Since(start)
			log.Printf("✅ Батч %d успешно проиндексирован за %v", batchNum, duration)
		}
		
		// Небольшая пауза между батчами
		time.Sleep(500 * time.Millisecond)
	}
	
	log.Println("🎉 Реиндексация завершена! Теперь можно переключить алиас...")
	log.Println("")
	log.Println("📋 Для переключения на новый индекс выполните:")
	log.Println("curl -X POST \"http://localhost:9200/_aliases\" -H \"Content-Type: application/json\" -d '{")
	log.Println("  \"actions\": [")
	log.Println("    {\"remove\": {\"index\": \"marketplace\", \"alias\": \"marketplace_current\"}},")
	log.Println("    {\"add\": {\"index\": \"marketplace_v2\", \"alias\": \"marketplace_current\"}}")
	log.Println("  ]")
	log.Println("}'")
	log.Println("")
	log.Println("🔍 Для проверки нового индекса:")
	log.Printf("curl -X GET \"http://localhost:9200/marketplace_v2/_search?size=1&pretty\"")
}
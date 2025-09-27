package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"backend/internal/config"
	"backend/internal/server"
	"backend/pkg/utils"
)

func main() {
	log.Println("Запуск переиндексации OpenSearch...")
	startTime := time.Now()

	// Загружаем конфигурацию
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Создаем контекст
	ctx := context.Background()

	// Создаем сервер для инициализации всех сервисов
	srv, err := server.NewServer(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Получаем доступ к базе данных и сервисам через сервер
	dbPool := srv.Storage().PostgresPool()
	marketplaceService := srv.MarketplaceService()

	// Параметр для фильтрации по категориям (автомобили)
	categoryFilter := ""
	if len(os.Args) > 1 && os.Args[1] == "cars" {
		categoryFilter = "AND category_id IN (1301, 1303)"
		log.Println("Переиндексация только автомобилей...")
	} else {
		log.Println("Переиндексация всех объявлений...")
	}

	// Получаем все активные объявления
	query := fmt.Sprintf(`
		SELECT id FROM marketplace_listings
		WHERE status = 'active'
		%s
		ORDER BY id
	`, categoryFilter)

	rows, err := dbPool.Query(ctx, query)
	if err != nil {
		log.Fatalf("Ошибка получения списка объявлений: %v", err)
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			log.Printf("Ошибка сканирования ID: %v", err)
			continue
		}
		ids = append(ids, id)
	}

	log.Printf("Найдено %d объявлений для переиндексации", len(ids))

	// Переиндексируем каждое объявление
	success := 0
	errors := 0

	for i, id := range ids {
		// Получаем полное объявление с атрибутами
		listing, err := marketplaceService.GetListingByID(ctx, id)
		if err != nil {
			log.Printf("[%d/%d] ❌ Ошибка получения объявления %d: %v", i+1, len(ids), id, err)
			errors++
			continue
		}

		// Проверяем наличие атрибутов
		if len(listing.Attributes) > 0 {
			log.Printf("[%d/%d] Объявление %d имеет %d атрибутов", i+1, len(ids), id, len(listing.Attributes))
			// Выводим первые несколько атрибутов для проверки
			for j, attr := range listing.Attributes {
				if j < 3 {
					log.Printf("  - %s: %s", attr.AttributeName, utils.GetAttributeDisplayValue(attr))
				}
			}
		}

		// Теперь индексируем в OpenSearch
		if err := srv.MarketplaceService().IndexListing(ctx, id, false); err != nil {
			log.Printf("[%d/%d] ❌ Ошибка индексации объявления %d: %v", i+1, len(ids), id, err)
			errors++
		} else {
			log.Printf("[%d/%d] ✅ Объявление %d успешно переиндексировано", i+1, len(ids), id)
			success++
		}
	}

	// Закрываем соединения через сервер при завершении
	duration := time.Since(startTime)
	log.Printf("\n=== Переиндексация завершена ===")
	log.Printf("Успешно: %d", success)
	log.Printf("Ошибки: %d", errors)
	log.Printf("Время выполнения: %s", duration)
}

package main

import (
	"backend/internal/config"
	"backend/internal/domain/models"
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Подключаемся к базе данных
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Получаем объявления категории автомобилей
	rows, err := pool.Query(context.Background(), `
		SELECT id FROM marketplace_listings
		WHERE category_id = 2000 OR category_id IN (
			SELECT id FROM marketplace_categories WHERE parent_id = 2000
		)
	`)
	if err != nil {
		log.Fatalf("Failed to query listings: %v", err)
	}
	defer rows.Close()

	var listingIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		listingIDs = append(listingIDs, id)
	}

	log.Printf("Found %d vehicle listings", len(listingIDs))

	// Годы выпуска для случайной генерации
	years := []int{2010, 2012, 2015, 2018, 2020, 2022}
	makes := []string{"Audi", "BMW", "Mercedes", "Toyota", "Honda", "Ford"}
	models := []string{"A4", "X5", "C-Class", "Corolla", "Civic", "Focus"}

	// Добавляем атрибуты для каждого объявления
	for _, id := range listingIDs {
		// Добавляем атрибут года выпуска
		year := years[rand.Intn(len(years))]
		yearValue := float64(year)
		
		// Добавляем атрибут марки
		make := makes[rand.Intn(len(makes))]
		
		// Добавляем атрибут модели
		model := models[rand.Intn(len(models))]

		// Сначала удаляем существующие атрибуты
		_, err := pool.Exec(context.Background(), 
			"DELETE FROM listing_attribute_values WHERE listing_id = $1", id)
		if err != nil {
			log.Printf("Error deleting attributes for listing %d: %v", id, err)
			continue
		}

		// Добавляем атрибуты
		_, err = pool.Exec(context.Background(), `
			INSERT INTO listing_attribute_values 
			(listing_id, attribute_id, attribute_name, numeric_value, text_value)
			VALUES 
			($1, 2, 'year', $2, NULL),
			($1, 0, 'make', NULL, $3),
			($1, 1, 'model', NULL, $4)
		`, id, yearValue, make, model)

		if err != nil {
			log.Printf("Error adding attributes for listing %d: %v", id, err)
			continue
		}

		log.Printf("Added attributes to listing %d: Year=%d, Make=%s, Model=%s", 
			id, year, make, model)
	}

	log.Printf("Completed adding attributes to %d listings", len(listingIDs))
}
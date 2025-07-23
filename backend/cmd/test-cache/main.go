package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"backend/internal/cache"
	"backend/internal/config"
)

// Category представляет категорию для тестирования
type Category struct {
	ID          int               `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	ParentID    *int              `json:"parent_id"`
	IsActive    bool              `json:"is_active"`
	Attributes  map[string]string `json:"attributes"`
}

func main() {
	ctx := context.Background()

	// Загружаем конфигурацию
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Создаем логгер
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Создаем Redis кеш
	redisCache, err := cache.NewRedisCache(
		ctx,
		cfg.Redis.URL,
		cfg.Redis.Password,
		cfg.Redis.DB,
		cfg.Redis.PoolSize,
		logger,
	)
	if err != nil {
		log.Fatalf("Failed to create Redis cache: %v", err)
	}
	defer func() {
		if err := redisCache.Close(); err != nil {
			log.Printf("Failed to close Redis connection: %v", err)
		}
	}()

	// Создаем адаптер
	cacheAdapter := cache.NewAdapter(redisCache)

	fmt.Println("=== Redis Cache Test ===")
	fmt.Printf("Connected to Redis at: %s\n\n", cfg.Redis.URL)

	// 1. Тест базовых операций
	fmt.Println("1. Testing basic Set/Get operations:")
	testKey := "test:simple:key"
	testData := map[string]string{
		"message": "Hello from Redis!",
		"time":    time.Now().Format(time.RFC3339),
	}

	// Сохраняем
	err = cacheAdapter.Set(ctx, testKey, testData, 1*time.Minute)
	if err != nil {
		log.Printf("Error setting value: %v", err)
	} else {
		fmt.Printf("   ✓ Set data with key: %s\n", testKey)
	}

	// Получаем
	var result map[string]string
	err = cacheAdapter.Get(ctx, testKey, &result)
	if err != nil {
		log.Printf("Error getting value: %v", err)
	} else {
		fmt.Printf("   ✓ Got data: %v\n", result)
	}

	// 2. Тест кеширования категорий
	fmt.Println("\n2. Testing category caching:")
	categories := []Category{
		{
			ID:          1,
			Name:        "Electronics",
			Description: "Electronic devices",
			IsActive:    true,
			Attributes: map[string]string{
				"icon": "📱",
				"slug": "electronics",
			},
		},
		{
			ID:          2,
			Name:        "Clothing",
			Description: "Clothes and accessories",
			IsActive:    true,
			Attributes: map[string]string{
				"icon": "👕",
				"slug": "clothing",
			},
		},
	}

	// Кешируем категории для разных локалей
	locales := []string{"en", "ru", "sr"}
	for _, locale := range locales {
		key := cache.BuildCategoriesKey(locale)
		err = cacheAdapter.Set(ctx, key, categories, 6*time.Hour)
		if err != nil {
			log.Printf("Error caching categories for %s: %v", locale, err)
		} else {
			fmt.Printf("   ✓ Cached categories for locale: %s\n", locale)
		}
	}

	// 3. Тест GetOrSet
	fmt.Println("\n3. Testing GetOrSet functionality:")
	loadCount := 0
	getOrSetKey := "test:getorset:data"

	// Первый вызов - загружает данные
	var data1 map[string]interface{}
	err = cacheAdapter.GetOrSet(ctx, getOrSetKey, &data1, 5*time.Minute, func() (interface{}, error) {
		loadCount++
		fmt.Printf("   → Loader called (count: %d), simulating DB query...\n", loadCount)
		return map[string]interface{}{
			"id":        100,
			"name":      "Loaded from DB",
			"timestamp": time.Now().Unix(),
		}, nil
	})
	if err != nil {
		log.Printf("Error in GetOrSet: %v", err)
	} else {
		fmt.Printf("   ✓ First call result: %v\n", data1)
	}

	// Небольшая задержка для завершения асинхронного сохранения
	time.Sleep(100 * time.Millisecond)

	// Второй вызов - должен получить из кеша
	var data2 map[string]interface{}
	err = cacheAdapter.GetOrSet(ctx, getOrSetKey, &data2, 5*time.Minute, func() (interface{}, error) {
		loadCount++
		fmt.Printf("   → Loader called (count: %d), simulating DB query...\n", loadCount)
		return map[string]interface{}{
			"id":        100,
			"name":      "Loaded from DB",
			"timestamp": time.Now().Unix(),
		}, nil
	})
	if err != nil {
		log.Printf("Error in GetOrSet: %v", err)
	} else {
		fmt.Printf("   ✓ Second call result: %v\n", data2)
		fmt.Printf("   ✓ Loader was called %d time(s) (should be 1)\n", loadCount)
	}

	// 4. Тест инвалидации кеша
	fmt.Println("\n4. Testing cache invalidation:")
	pattern := cache.BuildAllCategoriesInvalidationPattern()
	err = cacheAdapter.DeletePattern(ctx, pattern)
	if err != nil {
		log.Printf("Error deleting pattern: %v", err)
	} else {
		fmt.Printf("   ✓ Deleted all category keys with pattern: %s\n", pattern)
	}

	// Проверяем, что ключи удалены
	for _, locale := range locales {
		key := cache.BuildCategoriesKey(locale)
		var temp []Category
		err = cacheAdapter.Get(ctx, key, &temp)
		switch {
		case errors.Is(err, cache.ErrCacheMiss):
			fmt.Printf("   ✓ Key %s successfully deleted\n", key)
		case err != nil:
			fmt.Printf("   × Error checking key %s: %v\n", key, err)
		default:
			fmt.Printf("   × Key %s still exists!\n", key)
		}
	}

	// 5. Запускаем простой HTTP сервер для демонстрации
	fmt.Println("\n5. Starting HTTP server for cache demonstration...")
	fmt.Println("   Available endpoints:")
	fmt.Println("   - GET  /categories?locale=en")
	fmt.Println("   - POST /categories/invalidate")
	fmt.Println("   - GET  /stats")

	http.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		locale := r.URL.Query().Get("locale")
		if locale == "" {
			locale = "en"
		}

		key := cache.BuildCategoriesKey(locale)
		var cats []Category

		// Используем GetOrSet для автоматической загрузки
		err := cacheAdapter.GetOrSet(ctx, key, &cats, 6*time.Hour, func() (interface{}, error) {
			logger.Info("Loading categories from 'database' for locale: ", locale)
			// Симулируем загрузку из БД
			return categories, nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"locale":     locale,
			"categories": cats,
			"cached":     true,
		}); err != nil {
			log.Printf("Failed to encode categories response: %v", err)
		}
	})

	http.HandleFunc("/categories/invalidate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		pattern := cache.BuildAllCategoriesInvalidationPattern()
		err := cacheAdapter.DeletePattern(ctx, pattern)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "success",
			"pattern": pattern,
			"message": "All category cache entries invalidated",
		}); err != nil {
			log.Printf("Failed to encode invalidate response: %v", err)
		}
	})

	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		// Показываем статистику использования кеша
		stats := map[string]interface{}{
			"redis_url":    cfg.Redis.URL,
			"redis_db":     cfg.Redis.DB,
			"pool_size":    cfg.Redis.PoolSize,
			"current_time": time.Now().Format(time.RFC3339),
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(stats); err != nil {
			log.Printf("Failed to encode stats response: %v", err)
		}
	})

	fmt.Println("\n   Server starting on http://localhost:8080")
	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Printf("Server failed: %v", err)
	}
}

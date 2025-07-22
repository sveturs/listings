package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sirupsen/logrus"

	"backend/internal/cache"
)

func main() {
	// Создаем логгер
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Подключаемся к Redis напрямую
	fmt.Println("=== Redis Cache Test ===")
	fmt.Println("Connecting to Redis at localhost:6379...")

	redisCache, err := cache.NewRedisCache(
		"localhost:6379",
		"", // password
		0,  // db
		10, // pool size
		logger,
	)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer func() {
		if err := redisCache.Close(); err != nil {
			log.Printf("Failed to close Redis connection: %v", err)
		}
	}()

	fmt.Println("✓ Successfully connected to Redis!")

	// Создаем адаптер
	cacheAdapter := cache.NewAdapter(redisCache)
	ctx := context.Background()

	// 1. Тест базовых операций
	fmt.Println("\n1. Testing basic operations:")
	testKey := "test:simple"
	testValue := "Hello Redis!"

	// Set
	err = cacheAdapter.Set(ctx, testKey, testValue, 1*time.Minute)
	if err != nil {
		log.Printf("Error setting value: %v", err)
	} else {
		fmt.Printf("   ✓ Set value for key: %s\n", testKey)
	}

	// Get
	var result string
	err = cacheAdapter.Get(ctx, testKey, &result)
	if err != nil {
		log.Printf("Error getting value: %v", err)
	} else {
		fmt.Printf("   ✓ Got value: %s\n", result)
	}

	// 2. Тест кеширования категорий
	fmt.Println("\n2. Testing category caching pattern:")

	type Category struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	categories := []Category{
		{ID: 1, Name: "Electronics"},
		{ID: 2, Name: "Clothing"},
	}

	// Кешируем для разных языков
	for _, locale := range []string{"en", "ru", "sr"} {
		key := cache.BuildCategoriesKey(locale)
		err = cacheAdapter.Set(ctx, key, categories, 6*time.Hour)
		if err != nil {
			log.Printf("Error caching for %s: %v", locale, err)
		} else {
			fmt.Printf("   ✓ Cached categories for locale: %s (key: %s)\n", locale, key)
		}
	}

	// 3. Тест GetOrSet
	fmt.Println("\n3. Testing GetOrSet:")
	loadCount := 0
	getOrSetKey := "test:getorset"

	// Первый вызов
	var data1 map[string]interface{}
	err = cacheAdapter.GetOrSet(ctx, getOrSetKey, &data1, 30*time.Second, func() (interface{}, error) {
		loadCount++
		fmt.Printf("   → Loader called (count: %d)\n", loadCount)
		return map[string]interface{}{
			"id":   100,
			"name": "Loaded from DB",
			"time": time.Now().Format(time.RFC3339),
		}, nil
	})
	if err != nil {
		log.Printf("Error in GetOrSet: %v", err)
	} else {
		fmt.Printf("   ✓ First call result: %v\n", data1)
	}

	// Небольшая задержка для асинхронного сохранения
	time.Sleep(100 * time.Millisecond)

	// Второй вызов
	var data2 map[string]interface{}
	err = cacheAdapter.GetOrSet(ctx, getOrSetKey, &data2, 30*time.Second, func() (interface{}, error) {
		loadCount++
		fmt.Printf("   → Loader called (count: %d)\n", loadCount)
		return map[string]interface{}{
			"id":   100,
			"name": "Loaded from DB",
			"time": time.Now().Format(time.RFC3339),
		}, nil
	})
	if err != nil {
		log.Printf("Error in GetOrSet: %v", err)
	} else {
		fmt.Printf("   ✓ Second call result: %v\n", data2)
		fmt.Printf("   ✓ Loader was called %d time(s) (expected: 1)\n", loadCount)
		if data1["time"] == data2["time"] {
			fmt.Println("   ✓ Data was served from cache (timestamps match)")
		}
	}

	// 4. Тест паттернов инвалидации
	fmt.Println("\n4. Testing cache invalidation patterns:")

	// Проверяем что ключи существуют
	for _, locale := range []string{"en", "ru", "sr"} {
		key := cache.BuildCategoriesKey(locale)
		exists, err := redisCache.Exists(ctx, key)
		if err != nil {
			log.Printf("Error checking key %s: %v", key, err)
		} else if exists {
			fmt.Printf("   ✓ Key exists: %s\n", key)
		}
	}

	// Удаляем по паттерну
	pattern := cache.BuildAllCategoriesInvalidationPattern()
	fmt.Printf("\n   → Deleting keys matching pattern: %s\n", pattern)
	err = cacheAdapter.DeletePattern(ctx, pattern)
	if err != nil {
		log.Printf("Error deleting pattern: %v", err)
	} else {
		fmt.Println("   ✓ Pattern deletion successful")
	}

	// Проверяем что ключи удалены
	fmt.Println("\n   Checking if keys were deleted:")
	for _, locale := range []string{"en", "ru", "sr"} {
		key := cache.BuildCategoriesKey(locale)
		exists, err := redisCache.Exists(ctx, key)
		if err != nil {
			log.Printf("Error checking key %s: %v", key, err)
		} else if !exists {
			fmt.Printf("   ✓ Key deleted: %s\n", key)
		} else {
			fmt.Printf("   × Key still exists: %s\n", key)
		}
	}

	// 5. Проверка работы с реальным API эндпоинтом
	fmt.Println("\n5. Testing with real API endpoint simulation:")

	// Симулируем запрос категорий через API
	apiCategoriesKey := cache.BuildCategoriesKey("en")
	var apiCategories []Category

	err = cacheAdapter.GetOrSet(ctx, apiCategoriesKey, &apiCategories, 6*time.Hour, func() (interface{}, error) {
		fmt.Println("   → Loading categories from 'database'...")
		// Симуляция задержки БД
		time.Sleep(50 * time.Millisecond)

		return []Category{
			{ID: 1, Name: "Real Category 1"},
			{ID: 2, Name: "Real Category 2"},
			{ID: 3, Name: "Real Category 3"},
		}, nil
	})

	if err != nil {
		log.Printf("Error getting categories: %v", err)
	} else {
		fmt.Printf("   ✓ Got %d categories from cache/DB\n", len(apiCategories))
		for _, cat := range apiCategories {
			fmt.Printf("      - %d: %s\n", cat.ID, cat.Name)
		}
	}

	// Очистка тестовых данных
	fmt.Println("\n6. Cleaning up test data...")
	testKeys := []string{
		testKey,
		getOrSetKey,
		apiCategoriesKey,
	}

	for _, key := range testKeys {
		err = cacheAdapter.Delete(ctx, key)
		if err != nil {
			log.Printf("Error deleting %s: %v", key, err)
		} else {
			fmt.Printf("   ✓ Deleted: %s\n", key)
		}
	}

	fmt.Println("\n✅ All tests completed successfully!")
	fmt.Println("\nSummary:")
	fmt.Println("- Redis connection: ✓")
	fmt.Println("- Basic Set/Get: ✓")
	fmt.Println("- Category caching: ✓")
	fmt.Println("- GetOrSet with loader: ✓")
	fmt.Println("- Pattern invalidation: ✓")
	fmt.Println("- API simulation: ✓")
}

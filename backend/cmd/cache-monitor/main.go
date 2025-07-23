package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"

	"backend/internal/cache"
)

func main() {
	// Создаем логгер
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	// Подключаемся к Redis
	fmt.Println("=== Redis Cache Monitor ===")
	fmt.Println("Connecting to Redis at localhost:6379...")

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer func() {
		if err := client.Close(); err != nil {
			log.Printf("Failed to close Redis client: %v", err)
		}
	}()

	ctx := context.Background()

	// Проверяем подключение
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	fmt.Println("✓ Connected to Redis!")
	fmt.Println()

	// Мониторинг ключей кеша
	fmt.Println("Cache Keys Monitoring:")
	fmt.Println("---------------------")

	// Проверяем ключи категорий
	fmt.Println("\n1. Category Cache Keys:")
	categoryPattern := cache.BuildAllCategoriesInvalidationPattern()
	scanAndDisplay(ctx, client, categoryPattern, "Category")

	// Проверяем ключи атрибутов
	fmt.Println("\n2. Attribute Cache Keys:")
	scanAndDisplay(ctx, client, cache.PrefixCategoryAttrs+"*", "Attribute")

	// Проверяем ключи групп атрибутов
	fmt.Println("\n3. Attribute Group Cache Keys:")
	scanAndDisplay(ctx, client, cache.PrefixAttributeGroups+"*", "Attribute Group")

	// Проверяем дерево категорий
	fmt.Println("\n4. Category Tree Cache Keys:")
	scanAndDisplay(ctx, client, cache.PrefixCategoryTree+"*", "Category Tree")

	// Статистика
	fmt.Println("\n5. Cache Statistics:")
	fmt.Println("-------------------")

	// Информация о БД
	info, err := client.Info(ctx, "keyspace").Result()
	if err == nil {
		fmt.Println("Database info:")
		fmt.Println(info)
	}

	// Размер БД
	dbSize, err := client.DBSize(ctx).Result()
	if err == nil {
		fmt.Printf("Total keys in database: %d\n", dbSize)
	}

	// Использование памяти
	memInfo, err := client.Info(ctx, "memory").Result()
	if err == nil {
		fmt.Println("\nMemory usage (excerpt):")
		// Парсим только нужные строки
		for _, line := range []string{"used_memory_human", "used_memory_peak_human"} {
			if idx := findLine(memInfo, line); idx != -1 {
				// Найти конец строки
				end := idx
				for end < len(memInfo) && memInfo[end] != '\n' {
					end++
				}
				fmt.Printf("  %s\n", memInfo[idx:end])
			}
		}
	}

	// Тест производительности
	fmt.Println("\n6. Performance Test:")
	fmt.Println("-------------------")
	testCachePerformance(ctx, client)
}

func scanAndDisplay(ctx context.Context, client *redis.Client, pattern string, keyType string) {
	var cursor uint64
	var totalKeys int
	keyExamples := make([]string, 0, 5)

	for {
		keys, nextCursor, err := client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			log.Printf("Error scanning %s keys: %v", keyType, err)
			return
		}

		totalKeys += len(keys)

		// Сохраняем первые 5 ключей как примеры
		for _, key := range keys {
			if len(keyExamples) < 5 {
				keyExamples = append(keyExamples, key)
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	fmt.Printf("Pattern: %s\n", pattern)
	fmt.Printf("Found %d %s keys\n", totalKeys, keyType)

	if len(keyExamples) > 0 {
		fmt.Println("Examples:")
		for _, key := range keyExamples {
			// Получаем TTL для ключа
			ttl, err := client.TTL(ctx, key).Result()
			ttlStr := "no TTL"
			if err == nil && ttl > 0 {
				ttlStr = fmt.Sprintf("TTL: %v", ttl)
			}

			// Получаем тип данных
			dataType, _ := client.Type(ctx, key).Result()

			fmt.Printf("  - %s (type: %s, %s)\n", key, dataType, ttlStr)
		}
	}
}

func findLine(text string, prefix string) int {
	// Простой поиск строки в тексте
	start := 0
	for i := 0; i < len(text); i++ {
		if text[i] == '\n' || i == 0 {
			if i > 0 {
				start = i + 1
			}
			if start < len(text) && len(text) >= start+len(prefix) {
				if text[start:start+len(prefix)] == prefix {
					end := start
					for end < len(text) && text[end] != '\n' {
						end++
					}
					return start
				}
			}
		}
	}
	return -1
}

func testCachePerformance(ctx context.Context, client *redis.Client) {
	// Создаем тестовый Redis кеш
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	redisCache, err := cache.NewRedisCache(ctx, "localhost:6379", "", 0, 10, logger)
	if err != nil {
		log.Printf("Failed to create cache for performance test: %v", err)
		return
	}
	defer func() {
		if err := redisCache.Close(); err != nil {
			log.Printf("Failed to close Redis cache: %v", err)
		}
	}()

	adapter := cache.NewAdapter(redisCache)

	// Тест записи
	fmt.Println("Write performance:")
	testData := map[string]interface{}{
		"id":   1,
		"name": "Test Category",
		"desc": "Performance test data",
	}

	iterations := 100
	start := time.Now()

	for i := 0; i < iterations; i++ {
		key := fmt.Sprintf("perf:test:%d", i)
		_ = adapter.Set(ctx, key, testData, 1*time.Minute)
	}

	writeTime := time.Since(start)
	fmt.Printf("  Written %d keys in %v (avg: %v/key)\n",
		iterations, writeTime, writeTime/time.Duration(iterations))

	// Тест чтения
	fmt.Println("\nRead performance:")
	start = time.Now()

	for i := 0; i < iterations; i++ {
		key := fmt.Sprintf("perf:test:%d", i)
		var result map[string]interface{}
		_ = adapter.Get(ctx, key, &result)
	}

	readTime := time.Since(start)
	fmt.Printf("  Read %d keys in %v (avg: %v/key)\n",
		iterations, readTime, readTime/time.Duration(iterations))

	// Очистка тестовых данных
	_ = adapter.DeletePattern(ctx, "perf:test:*")

	// Проверка GetOrSet
	fmt.Println("\nGetOrSet performance:")
	loadCount := 0
	start = time.Now()

	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("getorset:test:%d", i)
		var result map[string]interface{}
		_ = adapter.GetOrSet(ctx, key, &result, 1*time.Minute, func() (interface{}, error) {
			loadCount++
			return testData, nil
		})
	}

	// Второй проход - все должно быть из кеша
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("getorset:test:%d", i)
		var result map[string]interface{}
		_ = adapter.GetOrSet(ctx, key, &result, 1*time.Minute, func() (interface{}, error) {
			loadCount++
			return testData, nil
		})
	}

	getOrSetTime := time.Since(start)
	fmt.Printf("  20 GetOrSet calls in %v (loader called %d times)\n", getOrSetTime, loadCount)

	// Очистка
	_ = adapter.DeletePattern(ctx, "getorset:test:*")
}

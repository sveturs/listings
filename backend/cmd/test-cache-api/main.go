package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"backend/internal/cache"
	"backend/internal/common"
)

func main() {
	ctx := context.Background()
	
	// Создаем логгер
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	// Подключаемся к Redis
	redisCache, err := cache.NewRedisCache(ctx, "localhost:6379", "", 0, 10, logger)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer func() {
		if err := redisCache.Close(); err != nil {
			log.Printf("Failed to close redis cache: %v", err)
		}
	}()

	cacheAdapter := cache.NewAdapter(redisCache)

	fmt.Println("=== Testing API Cache Behavior ===")
	fmt.Println()

	// Тестируем разные сценарии
	testScenarios := []struct {
		name   string
		locale string
	}{
		{"English", "en"},
		{"Russian", "ru"},
		{"Serbian", "sr"},
		{"Default (no locale)", ""},
	}

	for _, scenario := range testScenarios {
		fmt.Printf("\nTesting scenario: %s\n", scenario.name)
		fmt.Println(strings.Repeat("-", 40))

		// Создаем контекст с locale
		ctx := context.Background()
		if scenario.locale != "" {
			ctx = context.WithValue(ctx, common.ContextKeyLocale, scenario.locale)
		}

		// Формируем ключ кеша
		locale := "en" // по умолчанию
		if l, ok := ctx.Value(common.ContextKeyLocale).(string); ok && l != "" {
			locale = l
		}
		cacheKey := cache.BuildCategoriesKey(locale)

		fmt.Printf("Context locale: %v\n", ctx.Value(common.ContextKeyLocale))
		fmt.Printf("Effective locale: %s\n", locale)
		fmt.Printf("Cache key: %s\n", cacheKey)

		// Проверяем, есть ли уже данные в кеше
		var cachedCategories []map[string]interface{}
		err := cacheAdapter.Get(ctx, cacheKey, &cachedCategories)
		if err == cache.ErrCacheMiss {
			fmt.Println("Status: NOT in cache")
		} else if err != nil {
			fmt.Printf("Status: Error checking cache: %v\n", err)
		} else {
			fmt.Printf("Status: FOUND in cache (%d categories)\n", len(cachedCategories))
		}

		// Симулируем вызов API
		fmt.Println("\nSimulating API call...")
		callAPI(scenario.locale)

		// Проверяем кеш снова
		time.Sleep(100 * time.Millisecond) // Даем время на обработку
		err = cacheAdapter.Get(ctx, cacheKey, &cachedCategories)
		if err == cache.ErrCacheMiss {
			fmt.Println("After API call: Still NOT in cache")
		} else if err != nil {
			fmt.Printf("After API call: Error: %v\n", err)
		} else {
			fmt.Printf("After API call: NOW in cache (%d categories)\n", len(cachedCategories))
		}
	}

	// Проверяем все ключи категорий в Redis
	fmt.Println("\n\nFinal cache state:")
	fmt.Println(strings.Repeat("=", 50))

	ctx := context.Background()

	// Проверяем ключи категорий
	foundKeys := []string{}

	// Пока просто проверим несколько ключей напрямую
	for _, locale := range []string{"en", "ru", "sr"} {
		key := cache.BuildCategoriesKey(locale)
		exists, err := redisCache.Exists(ctx, key)
		if err == nil && exists {
			foundKeys = append(foundKeys, key)
		}
	}

	if len(foundKeys) == 0 {
		fmt.Println("No category keys found in cache")
	} else {
		fmt.Printf("Found %d category keys:\n", len(foundKeys))
		for _, key := range foundKeys {
			fmt.Printf("  - %s\n", key)
		}
	}

	// Тестируем прямое сохранение в кеш
	fmt.Println("\n\nDirect cache test:")
	fmt.Println(strings.Repeat("=", 50))

	testCategories := []map[string]interface{}{
		{"id": 1, "name": "Test Category 1"},
		{"id": 2, "name": "Test Category 2"},
	}

	for _, locale := range []string{"en", "ru", "sr"} {
		key := cache.BuildCategoriesKey(locale)
		err := cacheAdapter.Set(ctx, key, testCategories, 1*time.Hour)
		if err != nil {
			fmt.Printf("Failed to set %s: %v\n", key, err)
		} else {
			fmt.Printf("Successfully cached: %s\n", key)
		}
	}

	// Проверяем что сохранилось
	fmt.Println("\nVerifying direct cache writes:")
	for _, locale := range []string{"en", "ru", "sr"} {
		key := cache.BuildCategoriesKey(locale)
		var result []map[string]interface{}
		err := cacheAdapter.Get(ctx, key, &result)
		if err != nil {
			fmt.Printf("  %s: Error: %v\n", key, err)
		} else {
			fmt.Printf("  %s: OK (%d items)\n", key, len(result))
		}
	}
}

func callAPI(locale string) {
	url := "http://localhost:3000/api/v1/marketplace/categories"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return
	}

	if locale != "" {
		req.Header.Set("Accept-Language", locale)
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error calling API: %v", err)
		return
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Failed to close response body: %v", err)
		}
	}()

	fmt.Printf("  API Response: %s\n", resp.Status)
}

var strings = struct {
	Repeat func(s string, count int) string
}{
	Repeat: func(s string, count int) string {
		result := ""
		for i := 0; i < count; i++ {
			result += s
		}
		return result
	},
}

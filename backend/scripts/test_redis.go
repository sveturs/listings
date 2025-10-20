//go:build ignore
// +build ignore

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"backend/internal/cache"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("=== Redis Cache Integration Test ===")
	fmt.Println()

	// Загружаем .env файл
	if err := godotenv.Load("../.env"); err != nil {
		fmt.Printf("⚠️  Warning: .env file not found, using defaults\n\n")
	}

	// Получаем конфигурацию из переменных окружения
	redisURL := getEnv("REDIS_URL", "localhost:6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")
	redisDB := 0 // Default DB

	fmt.Printf("📋 Configuration:\n")
	fmt.Printf("   Redis URL: %s\n", redisURL)
	fmt.Printf("   Redis DB: %d\n", redisDB)
	fmt.Println()

	// Создаем логгер
	logger := logrus.New()
	logger.SetLevel(logrus.WarnLevel) // Только warnings и errors

	ctx := context.Background()

	// Создаем Redis клиент
	fmt.Println("🔌 Connecting to Redis...")
	redisCache, err := cache.NewRedisCache(ctx, redisURL, redisPassword, redisDB, 10, logger)
	if err != nil {
		fmt.Printf("❌ Failed to connect to Redis: %v\n", err)
		os.Exit(1)
	}
	defer redisCache.Close()
	fmt.Println("✅ Connected to Redis successfully!")
	fmt.Println()

	// Test 1: SET and GET
	fmt.Println("📝 Test 1: SET and GET operations")
	fmt.Println("----------------------------------")

	type TestData struct {
		Message string `json:"message"`
		Count   int    `json:"count"`
	}

	testKey := "test:redis:integration"
	testValue := TestData{
		Message: "Hello from Redis test!",
		Count:   42,
	}

	// SET
	fmt.Printf("Setting key '%s' with TTL 60s...\n", testKey)
	if err := redisCache.Set(ctx, testKey, testValue, 60*time.Second); err != nil {
		fmt.Printf("❌ SET failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ SET successful")

	// GET
	var retrieved TestData
	fmt.Printf("Getting key '%s'...\n", testKey)
	if err := redisCache.Get(ctx, testKey, &retrieved); err != nil {
		fmt.Printf("❌ GET failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ GET successful: %+v\n", retrieved)

	// Verify data
	if retrieved.Message != testValue.Message || retrieved.Count != testValue.Count {
		fmt.Printf("❌ Data mismatch! Expected %+v, got %+v\n", testValue, retrieved)
		os.Exit(1)
	}
	fmt.Println("✅ Data verified successfully")
	fmt.Println()

	// Test 2: EXISTS
	fmt.Println("🔍 Test 2: EXISTS operation")
	fmt.Println("---------------------------")

	exists, err := redisCache.Exists(ctx, testKey)
	if err != nil {
		fmt.Printf("❌ EXISTS failed: %v\n", err)
		os.Exit(1)
	}
	if !exists {
		fmt.Printf("❌ Key should exist but doesn't\n")
		os.Exit(1)
	}
	fmt.Printf("✅ Key exists: %v\n", exists)
	fmt.Println()

	// Test 3: TTL verification (via GetClient)
	fmt.Println("⏰ Test 3: TTL verification")
	fmt.Println("---------------------------")

	client := redisCache.GetClient()
	ttl := client.TTL(ctx, testKey).Val()
	fmt.Printf("✅ TTL: %v (should be ~60s)\n", ttl)
	if ttl <= 0 || ttl > 61*time.Second {
		fmt.Printf("❌ TTL is incorrect: %v\n", ttl)
		os.Exit(1)
	}
	fmt.Println()

	// Test 4: DELETE
	fmt.Println("🗑️  Test 4: DELETE operation")
	fmt.Println("---------------------------")

	fmt.Printf("Deleting key '%s'...\n", testKey)
	if err := redisCache.Delete(ctx, testKey); err != nil {
		fmt.Printf("❌ DELETE failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ DELETE successful")

	// Verify deletion
	exists, err = redisCache.Exists(ctx, testKey)
	if err != nil {
		fmt.Printf("❌ EXISTS check failed: %v\n", err)
		os.Exit(1)
	}
	if exists {
		fmt.Printf("❌ Key should not exist but does\n")
		os.Exit(1)
	}
	fmt.Printf("✅ Key deleted successfully (exists: %v)\n", exists)
	fmt.Println()

	// Test 5: Pattern deletion
	fmt.Println("🔥 Test 5: DELETE by pattern")
	fmt.Println("----------------------------")

	// Create multiple keys
	for i := 1; i <= 3; i++ {
		key := fmt.Sprintf("test:pattern:%d", i)
		data := TestData{Message: fmt.Sprintf("Test %d", i), Count: i}
		if err := redisCache.Set(ctx, key, data, 60*time.Second); err != nil {
			fmt.Printf("❌ Failed to set key %s: %v\n", key, err)
			os.Exit(1)
		}
		fmt.Printf("  ✓ Created key: %s\n", key)
	}

	// Delete by pattern
	pattern := "test:pattern:*"
	fmt.Printf("Deleting keys matching pattern '%s'...\n", pattern)
	if err := redisCache.DeletePattern(ctx, pattern); err != nil {
		fmt.Printf("❌ DeletePattern failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Pattern deletion successful")

	// Verify all deleted
	for i := 1; i <= 3; i++ {
		key := fmt.Sprintf("test:pattern:%d", i)
		exists, err := redisCache.Exists(ctx, key)
		if err != nil {
			fmt.Printf("❌ EXISTS check failed for %s: %v\n", key, err)
			os.Exit(1)
		}
		if exists {
			fmt.Printf("❌ Key %s should not exist\n", key)
			os.Exit(1)
		}
	}
	fmt.Println("✅ All pattern keys deleted successfully")
	fmt.Println()

	// Test 6: Cache miss
	fmt.Println("❓ Test 6: Cache miss handling")
	fmt.Println("------------------------------")

	var missData TestData
	err = redisCache.Get(ctx, "nonexistent:key", &missData)
	if err == nil {
		fmt.Println("❌ Should return error for nonexistent key")
		os.Exit(1)
	}
	if err != cache.ErrCacheMiss {
		fmt.Printf("❌ Wrong error type: %v (expected ErrCacheMiss)\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ Cache miss handled correctly: %v\n", err)
	fmt.Println()

	fmt.Println("╔════════════════════════════════════╗")
	fmt.Println("║  ✅ All Redis tests passed! 🎉    ║")
	fmt.Println("╚════════════════════════════════════╝")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

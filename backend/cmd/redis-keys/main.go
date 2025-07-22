package main

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/redis/go-redis/v9"
)

func main() {
	// Подключаемся к Redis
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

	fmt.Println("=== Redis Keys Analysis ===")
	fmt.Println()

	// Сканируем все ключи
	var cursor uint64
	keyPrefixes := make(map[string]int)
	var allKeys []string

	for {
		keys, nextCursor, err := client.Scan(ctx, cursor, "*", 100).Result()
		if err != nil {
			log.Fatalf("Error scanning keys: %v", err)
		}

		allKeys = append(allKeys, keys...)

		// Анализируем префиксы
		for _, key := range keys {
			parts := strings.Split(key, ":")
			if len(parts) > 0 {
				prefix := parts[0]
				keyPrefixes[prefix]++
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	// Сортируем префиксы по количеству
	type prefixCount struct {
		Prefix string
		Count  int
	}

	var prefixList []prefixCount
	for prefix, count := range keyPrefixes {
		prefixList = append(prefixList, prefixCount{prefix, count})
	}

	sort.Slice(prefixList, func(i, j int) bool {
		return prefixList[i].Count > prefixList[j].Count
	})

	// Выводим статистику
	fmt.Printf("Total keys in Redis: %d\n", len(allKeys))
	fmt.Println("\nKey prefixes by count:")
	fmt.Println("---------------------")

	for i, pc := range prefixList {
		if i < 10 { // Показываем топ-10
			fmt.Printf("%-20s: %d keys\n", pc.Prefix, pc.Count)
		}
	}

	// Показываем примеры ключей для каждого префикса
	fmt.Println("\nSample keys by prefix:")
	fmt.Println("---------------------")

	prefixSamples := make(map[string][]string)
	for _, key := range allKeys {
		parts := strings.Split(key, ":")
		if len(parts) > 0 {
			prefix := parts[0]
			if len(prefixSamples[prefix]) < 3 {
				prefixSamples[prefix] = append(prefixSamples[prefix], key)
			}
		}
	}

	for i, pc := range prefixList {
		if i < 5 { // Показываем примеры для топ-5
			fmt.Printf("\n%s (%d keys):\n", pc.Prefix, pc.Count)
			samples := prefixSamples[pc.Prefix]
			for _, sample := range samples {
				// Получаем TTL
				ttl, _ := client.TTL(ctx, sample).Result()
				ttlStr := "no TTL"
				if ttl > 0 {
					ttlStr = fmt.Sprintf("TTL: %v", ttl)
				}

				// Получаем размер значения
				size, _ := client.MemoryUsage(ctx, sample).Result()
				sizeStr := fmt.Sprintf("size: %d bytes", size)

				fmt.Printf("  - %s (%s, %s)\n", sample, ttlStr, sizeStr)
			}
		}
	}

	// Проверяем специфичные паттерны для marketplace
	fmt.Println("\n\nMarketplace-specific patterns:")
	fmt.Println("------------------------------")

	patterns := []string{
		"categories:*",
		"category_*",
		"attribute*",
		"listing*",
		"user*",
		"search*",
	}

	for _, pattern := range patterns {
		count := 0
		var examples []string

		for _, key := range allKeys {
			matched, _ := filepath.Match(pattern, key)
			if matched || strings.HasPrefix(key, strings.TrimSuffix(pattern, "*")) {
				count++
				if len(examples) < 3 {
					examples = append(examples, key)
				}
			}
		}

		if count > 0 {
			fmt.Printf("\n%s: %d keys\n", pattern, count)
			for _, ex := range examples {
				fmt.Printf("  - %s\n", ex)
			}
		}
	}
}

// Простая реализация match для паттернов
var filepath = struct {
	Match func(pattern, name string) (matched bool, err error)
}{
	Match: func(pattern, name string) (bool, error) {
		// Простая реализация для * паттернов
		if strings.HasSuffix(pattern, "*") {
			prefix := strings.TrimSuffix(pattern, "*")
			return strings.HasPrefix(name, prefix), nil
		}
		return pattern == name, nil
	},
}

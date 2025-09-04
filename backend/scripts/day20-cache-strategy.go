package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

// CacheStrategy управляет стратегией кеширования для unified attributes
type CacheStrategy struct {
	redis  *redis.Client
	ctx    context.Context
	prefix string
}

// CacheConfig конфигурация для различных типов кеша
type CacheConfig struct {
	AttributesTTL         time.Duration // TTL для списка атрибутов
	CategoryAttributesTTL time.Duration // TTL для атрибутов категории
	AttributeValuesTTL    time.Duration // TTL для значений атрибутов
	PopularCacheWarmup    bool          // Прогревание кеша популярных данных
	CompressionEnabled    bool          // Сжатие данных в кеше
}

// DefaultCacheConfig возвращает оптимальную конфигурацию кеширования
func DefaultCacheConfig() CacheConfig {
	return CacheConfig{
		AttributesTTL:         24 * time.Hour, // Атрибуты редко меняются
		CategoryAttributesTTL: 12 * time.Hour, // Связи категорий с атрибутами
		AttributeValuesTTL:    6 * time.Hour,  // Значения атрибутов
		PopularCacheWarmup:    true,
		CompressionEnabled:    true,
	}
}

// NewCacheStrategy создает новую стратегию кеширования
func NewCacheStrategy(redisAddr string, config CacheConfig) (*CacheStrategy, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	ctx := context.Background()

	// Проверяем подключение
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &CacheStrategy{
		redis:  rdb,
		ctx:    ctx,
		prefix: "unified_attrs:",
	}, nil
}

// CacheKeys содержит шаблоны ключей для кеширования
type CacheKeys struct {
	AllAttributes      string // "unified_attrs:all"
	CategoryAttributes string // "unified_attrs:category:{category_id}"
	AttributeValues    string // "unified_attrs:values:{entity_type}:{entity_id}"
	AttributeStats     string // "unified_attrs:stats:{attribute_id}"
	PopularCategories  string // "unified_attrs:popular_categories"
	SearchIndex        string // "unified_attrs:search_index"
}

func (cs *CacheStrategy) Keys() CacheKeys {
	return CacheKeys{
		AllAttributes:      cs.prefix + "all",
		CategoryAttributes: cs.prefix + "category:%d",
		AttributeValues:    cs.prefix + "values:%s:%d",
		AttributeStats:     cs.prefix + "stats:%d",
		PopularCategories:  cs.prefix + "popular_categories",
		SearchIndex:        cs.prefix + "search_index",
	}
}

// WarmupCache прогревает кеш наиболее популярными данными
func (cs *CacheStrategy) WarmupCache() error {
	log.Println("Starting cache warmup...")

	// 1. Кешируем все активные атрибуты
	if err := cs.cacheAllAttributes(); err != nil {
		log.Printf("Failed to cache all attributes: %v", err)
	}

	// 2. Кешируем атрибуты для популярных категорий
	popularCategories := []int{10101, 10102, 10103, 10110, 10120} // Автомобильные категории
	for _, categoryID := range popularCategories {
		if err := cs.cacheCategoryAttributes(categoryID); err != nil {
			log.Printf("Failed to cache attributes for category %d: %v", categoryID, err)
		}
	}

	// 3. Кешируем материализованное представление популярных атрибутов
	if err := cs.cachePopularCategoryAttributes(); err != nil {
		log.Printf("Failed to cache popular category attributes: %v", err)
	}

	log.Println("Cache warmup completed")
	return nil
}

// cacheAllAttributes кеширует список всех активных атрибутов
func (cs *CacheStrategy) cacheAllAttributes() error {
	// В реальном приложении здесь был бы запрос к БД
	// Для демонстрации создаем мок-данные
	attributes := map[string]interface{}{
		"data": []map[string]interface{}{
			{
				"id":             1,
				"code":           "car_make_id",
				"name":           "Марка автомобиля",
				"attribute_type": "select",
				"is_active":      true,
			},
			{
				"id":             2,
				"code":           "car_model_id",
				"name":           "Модель автомобиля",
				"attribute_type": "select",
				"is_active":      true,
			},
		},
		"cached_at": time.Now().Unix(),
		"ttl":       24 * 60 * 60, // 24 hours
	}

	data, err := json.Marshal(attributes)
	if err != nil {
		return err
	}

	key := cs.Keys().AllAttributes
	return cs.redis.Set(cs.ctx, key, data, 24*time.Hour).Err()
}

// cacheCategoryAttributes кеширует атрибуты для конкретной категории
func (cs *CacheStrategy) cacheCategoryAttributes(categoryID int) error {
	// Мок-данные для атрибутов категории
	categoryAttributes := map[string]interface{}{
		"category_id": categoryID,
		"attributes": []map[string]interface{}{
			{
				"id":          1,
				"name":        "Марка автомобиля",
				"is_required": true,
				"sort_order":  1,
			},
			{
				"id":          2,
				"name":        "Модель автомобиля",
				"is_required": true,
				"sort_order":  2,
			},
		},
		"cached_at": time.Now().Unix(),
	}

	data, err := json.Marshal(categoryAttributes)
	if err != nil {
		return err
	}

	key := fmt.Sprintf(cs.Keys().CategoryAttributes, categoryID)
	return cs.redis.Set(cs.ctx, key, data, 12*time.Hour).Err()
}

// cachePopularCategoryAttributes кеширует материализованное представление
func (cs *CacheStrategy) cachePopularCategoryAttributes() error {
	// Данные из материализованного представления mv_popular_category_attributes
	popularData := map[string]interface{}{
		"popular_combinations": []map[string]interface{}{
			{
				"category_id":    10101,
				"category_name":  "Легковые автомобили",
				"attribute_id":   1,
				"attribute_name": "Марка",
				"usage_count":    150,
				"avg_query_time": 0.12,
			},
		},
		"last_updated": time.Now().Unix(),
		"total_count":  25,
	}

	data, err := json.Marshal(popularData)
	if err != nil {
		return err
	}

	key := cs.Keys().PopularCategories
	return cs.redis.Set(cs.ctx, key, data, 6*time.Hour).Err()
}

// GetCachedAttributes получает кешированные атрибуты
func (cs *CacheStrategy) GetCachedAttributes(key string) (map[string]interface{}, error) {
	val, err := cs.redis.Get(cs.ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("cache miss")
	}
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal([]byte(val), &result)
	return result, err
}

// SetCacheWithTTL устанавливает значение в кеш с TTL
func (cs *CacheStrategy) SetCacheWithTTL(key string, data interface{}, ttl time.Duration) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return cs.redis.Set(cs.ctx, key, jsonData, ttl).Err()
}

// InvalidateCache инвалидирует кеш по паттерну
func (cs *CacheStrategy) InvalidateCache(pattern string) error {
	keys, err := cs.redis.Keys(cs.ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return cs.redis.Del(cs.ctx, keys...).Err()
	}

	return nil
}

// GetCacheStats возвращает статистику использования кеша
func (cs *CacheStrategy) GetCacheStats() (map[string]interface{}, error) {
	info, err := cs.redis.Info(cs.ctx, "stats").Result()
	if err != nil {
		return nil, err
	}

	// Получаем количество ключей в нашем namespace
	keys, err := cs.redis.Keys(cs.ctx, cs.prefix+"*").Result()
	if err != nil {
		return nil, err
	}

	memory, err := cs.redis.Info(cs.ctx, "memory").Result()
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_keys":          len(keys),
		"namespace":           cs.prefix,
		"redis_info_stats":    info,
		"redis_memory_info":   memory,
		"cache_warmup_status": "active",
		"timestamp":           time.Now().Unix(),
	}

	return stats, nil
}

// OptimizeCacheSettings настраивает оптимальные параметры Redis
func (cs *CacheStrategy) OptimizeCacheSettings() error {
	// Настройка политики вытеснения данных
	if err := cs.redis.ConfigSet(cs.ctx, "maxmemory-policy", "allkeys-lru").Err(); err != nil {
		log.Printf("Warning: Could not set maxmemory-policy: %v", err)
	}

	// Настройка оптимального размера hash-table
	if err := cs.redis.ConfigSet(cs.ctx, "hash-max-ziplist-entries", "512").Err(); err != nil {
		log.Printf("Warning: Could not set hash-max-ziplist-entries: %v", err)
	}

	// Настройка компрессии
	if err := cs.redis.ConfigSet(cs.ctx, "rdbcompression", "yes").Err(); err != nil {
		log.Printf("Warning: Could not set rdbcompression: %v", err)
	}

	log.Println("Redis cache settings optimized")
	return nil
}

// ScheduledCacheRefresh запускает периодическое обновление кеша
func (cs *CacheStrategy) ScheduledCacheRefresh(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("Running scheduled cache refresh...")
			if err := cs.WarmupCache(); err != nil {
				log.Printf("Scheduled cache refresh failed: %v", err)
			}
		}
	}
}

func main() {
	// Получаем адрес Redis из переменной окружения или используем дефолтный
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	// Создаем стратегию кеширования
	config := DefaultCacheConfig()
	cacheStrategy, err := NewCacheStrategy(redisAddr, config)
	if err != nil {
		log.Fatalf("Failed to initialize cache strategy: %v", err)
	}

	// Проверяем аргументы командной строки
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run day20-cache-strategy.go <command>")
		fmt.Println("Commands:")
		fmt.Println("  warmup     - Прогреть кеш популярными данными")
		fmt.Println("  stats      - Показать статистику кеша")
		fmt.Println("  optimize   - Оптимизировать настройки Redis")
		fmt.Println("  clear      - Очистить кеш unified attributes")
		fmt.Println("  refresh    - Запустить периодическое обновление (daemon)")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "warmup":
		log.Println("Warming up cache...")
		if err := cacheStrategy.WarmupCache(); err != nil {
			log.Fatalf("Cache warmup failed: %v", err)
		}
		log.Println("Cache warmup completed successfully")

	case "stats":
		stats, err := cacheStrategy.GetCacheStats()
		if err != nil {
			log.Fatalf("Failed to get cache stats: %v", err)
		}

		fmt.Println("=== UNIFIED ATTRIBUTES CACHE STATISTICS ===")
		fmt.Printf("Total keys in namespace: %d\n", stats["total_keys"])
		fmt.Printf("Namespace: %s\n", stats["namespace"])
		fmt.Printf("Timestamp: %d\n", stats["timestamp"])
		fmt.Printf("Cache warmup status: %s\n", stats["cache_warmup_status"])

		// Подробная статистика Redis
		fmt.Println("\n=== REDIS STATISTICS ===")
		fmt.Printf("Stats info: %s\n", stats["redis_info_stats"])

	case "optimize":
		log.Println("Optimizing Redis settings...")
		if err := cacheStrategy.OptimizeCacheSettings(); err != nil {
			log.Fatalf("Cache optimization failed: %v", err)
		}
		log.Println("Redis settings optimized successfully")

	case "clear":
		log.Println("Clearing unified attributes cache...")
		if err := cacheStrategy.InvalidateCache(cacheStrategy.prefix + "*"); err != nil {
			log.Fatalf("Cache clear failed: %v", err)
		}
		log.Println("Cache cleared successfully")

	case "refresh":
		log.Println("Starting scheduled cache refresh daemon...")
		refreshInterval := 6 * time.Hour
		if intervalStr := os.Getenv("CACHE_REFRESH_INTERVAL"); intervalStr != "" {
			if hours, err := strconv.Atoi(intervalStr); err == nil {
				refreshInterval = time.Duration(hours) * time.Hour
			}
		}
		fmt.Printf("Refresh interval: %v\n", refreshInterval)

		// Первоначальный прогрев
		if err := cacheStrategy.WarmupCache(); err != nil {
			log.Fatalf("Initial cache warmup failed: %v", err)
		}

		// Запуск daemon'а
		cacheStrategy.ScheduledCacheRefresh(refreshInterval)

	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}

package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var (
	ErrCacheMiss = errors.New("cache miss")
	ErrCacheNil  = redis.Nil
)

// CacheConfig конфигурация для разных типов кеша
type CacheConfig struct {
	SearchTTL          time.Duration // TTL для результатов поиска
	ListingDetailsTTL  time.Duration // TTL для деталей объявлений
	RecommendationsTTL time.Duration // TTL для рекомендаций
	CategoryStatsTTL   time.Duration // TTL для статистики категорий
	ViewHistoryTTL     time.Duration // TTL для истории просмотров
}

// DefaultCacheConfig возвращает дефолтную конфигурацию
func DefaultCacheConfig() *CacheConfig {
	return &CacheConfig{
		SearchTTL:          5 * time.Minute,
		ListingDetailsTTL:  1 * time.Hour,
		RecommendationsTTL: 30 * time.Minute,
		CategoryStatsTTL:   10 * time.Minute,
		ViewHistoryTTL:     24 * time.Hour,
	}
}

type UniversalCache struct {
	client *redis.Client
	logger *zap.Logger
	config *CacheConfig
}

func NewUniversalCache(ctx context.Context, addr string, logger *zap.Logger, config *CacheConfig) (*UniversalCache, error) {
	if config == nil {
		config = DefaultCacheConfig()
	}

	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     "",
		DB:           1, // Используем DB 1 для маркетплейса
		PoolSize:     50,
		MinIdleConns: 10,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolTimeout:  4 * time.Second,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logger.Info("Connected to Redis for universal cache",
		zap.String("addr", addr),
		zap.Int("poolSize", 50))

	return &UniversalCache{
		client: client,
		logger: logger,
		config: config,
	}, nil
}

// SearchCache методы для кеширования поисковых запросов

func (c *UniversalCache) GetSearchResults(ctx context.Context, key string) (interface{}, error) {
	return c.get(ctx, fmt.Sprintf("search:%s", key))
}

func (c *UniversalCache) SetSearchResults(ctx context.Context, key string, data interface{}) error {
	return c.set(ctx, fmt.Sprintf("search:%s", key), data, c.config.SearchTTL)
}

func (c *UniversalCache) InvalidateSearchCache(ctx context.Context) error {
	return c.deletePattern(ctx, "search:*")
}

// ListingCache методы для кеширования деталей объявлений

func (c *UniversalCache) GetListingDetails(ctx context.Context, listingID int) (interface{}, error) {
	return c.get(ctx, fmt.Sprintf("listing:%d", listingID))
}

func (c *UniversalCache) SetListingDetails(ctx context.Context, listingID int, data interface{}) error {
	return c.set(ctx, fmt.Sprintf("listing:%d", listingID), data, c.config.ListingDetailsTTL)
}

func (c *UniversalCache) InvalidateListingCache(ctx context.Context, listingID int) error {
	keys := []string{
		fmt.Sprintf("listing:%d", listingID),
		// Также инвалидируем связанные кеши
		fmt.Sprintf("recommendations:similar:%d", listingID),
	}

	for _, key := range keys {
		if err := c.delete(ctx, key); err != nil {
			c.logger.Warn("Failed to delete cache key", zap.String("key", key), zap.Error(err))
		}
	}

	// Инвалидируем поисковый кеш, так как объявление изменилось
	return c.InvalidateSearchCache(ctx)
}

// RecommendationsCache методы для кеширования рекомендаций

func (c *UniversalCache) GetRecommendations(ctx context.Context, recType string, params map[string]interface{}) (interface{}, error) {
	key := c.buildRecommendationKey(recType, params)
	return c.get(ctx, key)
}

func (c *UniversalCache) SetRecommendations(ctx context.Context, recType string, params map[string]interface{}, data interface{}) error {
	key := c.buildRecommendationKey(recType, params)
	return c.set(ctx, key, data, c.config.RecommendationsTTL)
}

func (c *UniversalCache) buildRecommendationKey(recType string, params map[string]interface{}) string {
	// Создаем уникальный ключ на основе типа и параметров
	paramStr := ""
	if categoryID, ok := params["category_id"].(int); ok {
		paramStr = fmt.Sprintf(":%d", categoryID)
	}
	if itemID, ok := params["item_id"].(int); ok {
		paramStr = fmt.Sprintf("%s:%d", paramStr, itemID)
	}
	if userID, ok := params["user_id"].(int); ok {
		paramStr = fmt.Sprintf("%s:u%d", paramStr, userID)
	}
	return fmt.Sprintf("recommendations:%s%s", recType, paramStr)
}

// CategoryStatsCache методы для кеширования статистики категорий

func (c *UniversalCache) GetCategoryStats(ctx context.Context, categoryID int) (interface{}, error) {
	return c.get(ctx, fmt.Sprintf("category:stats:%d", categoryID))
}

func (c *UniversalCache) SetCategoryStats(ctx context.Context, categoryID int, data interface{}) error {
	return c.set(ctx, fmt.Sprintf("category:stats:%d", categoryID), data, c.config.CategoryStatsTTL)
}

// ViewHistoryCache методы для кеширования истории просмотров

func (c *UniversalCache) GetUserViewHistory(ctx context.Context, userID int) ([]int, error) {
	key := fmt.Sprintf("viewhistory:user:%d", userID)

	// Получаем список ID из sorted set
	results, err := c.client.ZRevRange(ctx, key, 0, 99).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return []int{}, nil
		}
		return nil, err
	}

	ids := make([]int, 0, len(results))
	for _, idStr := range results {
		var id int
		if _, err := fmt.Sscanf(idStr, "%d", &id); err == nil {
			ids = append(ids, id)
		}
	}

	return ids, nil
}

func (c *UniversalCache) AddToViewHistory(ctx context.Context, userID int, listingID int) error {
	key := fmt.Sprintf("viewhistory:user:%d", userID)

	// Добавляем в sorted set с текущим временем как score
	score := float64(time.Now().Unix())
	if err := c.client.ZAdd(ctx, key, &redis.Z{
		Score:  score,
		Member: fmt.Sprintf("%d", listingID),
	}).Err(); err != nil {
		return err
	}

	// Устанавливаем TTL
	return c.client.Expire(ctx, key, c.config.ViewHistoryTTL).Err()
}

// Helper methods

func (c *UniversalCache) get(ctx context.Context, key string) (interface{}, error) {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrCacheMiss
		}
		c.logger.Debug("Cache get error", zap.String("key", key), zap.Error(err))
		return nil, err
	}

	var result interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		c.logger.Warn("Failed to unmarshal cached data", zap.String("key", key), zap.Error(err))
		return nil, err
	}

	// Обновляем счетчик попаданий
	c.incrementCounter(ctx, "hits")

	return result, nil
}

func (c *UniversalCache) set(ctx context.Context, key string, data interface{}, ttl time.Duration) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		c.logger.Warn("Failed to marshal data for cache", zap.String("key", key), zap.Error(err))
		return err
	}

	if err := c.client.Set(ctx, key, jsonData, ttl).Err(); err != nil {
		c.logger.Warn("Failed to set cache", zap.String("key", key), zap.Error(err))
		return err
	}

	// Обновляем счетчик промахов (это был промах, поэтому мы сохраняем)
	c.incrementCounter(ctx, "misses")

	return nil
}

func (c *UniversalCache) delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *UniversalCache) deletePattern(ctx context.Context, pattern string) error {
	var cursor uint64
	var err error

	for {
		var keys []string
		keys, cursor, err = c.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			if err := c.client.Del(ctx, keys...).Err(); err != nil {
				c.logger.Warn("Failed to delete keys", zap.Int("count", len(keys)), zap.Error(err))
			}
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}

func (c *UniversalCache) incrementCounter(ctx context.Context, counterType string) {
	key := fmt.Sprintf("cache:stats:%s", counterType)
	c.client.Incr(ctx, key)

	// Daily stats
	today := time.Now().Format("2006-01-02")
	dailyKey := fmt.Sprintf("cache:stats:daily:%s:%s", today, counterType)
	c.client.Incr(ctx, dailyKey)
	c.client.Expire(ctx, dailyKey, 7*24*time.Hour) // Keep for 7 days
}

// GetStats возвращает статистику кеша
func (c *UniversalCache) GetStats(ctx context.Context) (map[string]interface{}, error) {
	hits, _ := c.client.Get(ctx, "cache:stats:hits").Int64()
	misses, _ := c.client.Get(ctx, "cache:stats:misses").Int64()

	hitRate := float64(0)
	if hits+misses > 0 {
		hitRate = float64(hits) / float64(hits+misses) * 100
	}

	// Получаем информацию о памяти
	info := c.client.Info(ctx, "memory").Val()

	// Получаем размер БД
	dbSize, _ := c.client.DBSize(ctx).Result()

	return map[string]interface{}{
		"hits":       hits,
		"misses":     misses,
		"hitRate":    fmt.Sprintf("%.2f%%", hitRate),
		"totalKeys":  dbSize,
		"memoryInfo": info,
	}, nil
}

// WarmUp прогревает кеш популярными данными
func (c *UniversalCache) WarmUp(ctx context.Context) error {
	c.logger.Info("Starting cache warm-up")

	// Здесь можно добавить логику для прогрева кеша
	// Например, загрузка популярных категорий, последних объявлений и т.д.

	return nil
}

// Close закрывает соединение с Redis
func (c *UniversalCache) Close() error {
	return c.client.Close()
}

// Ping проверяет соединение с Redis
func (c *UniversalCache) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

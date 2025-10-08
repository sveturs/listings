package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

// ErrCacheKeyNotFound возвращается когда ключ не найден в кеше
var ErrCacheKeyNotFound = errors.New("cache key not found")

type RedisCache struct {
	client *redis.Client
	logger *zap.Logger
	prefix string
	ttl    time.Duration
}

func NewRedisCache(ctx context.Context, addr string, logger *zap.Logger) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     "", // no password set
		DB:           0,  // use default DB
		PoolSize:     10,
		MinIdleConns: 5,
		MaxRetries:   3,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisCache{
		client: client,
		logger: logger,
		prefix: "ai:category:",
		ttl:    15 * time.Minute,
	}, nil
}

// Get получает результат из кэша
func (r *RedisCache) Get(ctx context.Context, key string) (*AIDetectionResult, error) {
	fullKey := r.prefix + key

	data, err := r.client.Get(ctx, fullKey).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrCacheKeyNotFound
		}
		r.logger.Warn("Failed to get from Redis cache",
			zap.String("key", fullKey),
			zap.Error(err))
		return nil, err
	}

	var result AIDetectionResult
	if err := json.Unmarshal(data, &result); err != nil {
		r.logger.Warn("Failed to unmarshal cached data",
			zap.String("key", fullKey),
			zap.Error(err))
		return nil, err
	}

	// Обновляем счетчик попаданий в кэш
	r.incrementHitCounter(ctx)

	return &result, nil
}

// Set сохраняет результат в кэше
func (r *RedisCache) Set(ctx context.Context, key string, result *AIDetectionResult) error {
	fullKey := r.prefix + key

	data, err := json.Marshal(result)
	if err != nil {
		r.logger.Warn("Failed to marshal data for cache",
			zap.String("key", fullKey),
			zap.Error(err))
		return err
	}

	if err := r.client.Set(ctx, fullKey, data, r.ttl).Err(); err != nil {
		r.logger.Warn("Failed to set Redis cache",
			zap.String("key", fullKey),
			zap.Error(err))
		return err
	}

	// Обновляем счетчик промахов в кэш
	r.incrementMissCounter(ctx)

	return nil
}

// Delete удаляет запись из кэша
func (r *RedisCache) Delete(ctx context.Context, key string) error {
	fullKey := r.prefix + key
	return r.client.Del(ctx, fullKey).Err()
}

// Clear очищает весь кэш AI категорий
func (r *RedisCache) Clear(ctx context.Context) error {
	pattern := r.prefix + "*"

	// Получаем все ключи по паттерну
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		return nil
	}

	// Удаляем все найденные ключи
	return r.client.Del(ctx, keys...).Err()
}

// GetStats возвращает статистику кэша
func (r *RedisCache) GetStats(ctx context.Context) (map[string]interface{}, error) {
	hits, _ := r.client.Get(ctx, "ai:stats:hits").Int64()
	misses, _ := r.client.Get(ctx, "ai:stats:misses").Int64()

	// Получаем информацию о памяти
	info := r.client.Info(ctx, "memory").Val()

	// Подсчитываем количество ключей
	keys, _ := r.client.Keys(ctx, r.prefix+"*").Result()

	hitRate := float64(0)
	if hits+misses > 0 {
		hitRate = float64(hits) / float64(hits+misses) * 100
	}

	return map[string]interface{}{
		"hits":       hits,
		"misses":     misses,
		"hitRate":    hitRate,
		"keysCount":  len(keys),
		"memoryInfo": info,
	}, nil
}

// SetWithTTL сохраняет результат с кастомным TTL
func (r *RedisCache) SetWithTTL(ctx context.Context, key string, result *AIDetectionResult, ttl time.Duration) error {
	fullKey := r.prefix + key

	data, err := json.Marshal(result)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, fullKey, data, ttl).Err()
}

// BatchGet получает множество результатов за один запрос
func (r *RedisCache) BatchGet(ctx context.Context, keys []string) (map[string]*AIDetectionResult, error) {
	if len(keys) == 0 {
		return make(map[string]*AIDetectionResult), nil
	}

	// Подготавливаем полные ключи
	fullKeys := make([]string, len(keys))
	for i, key := range keys {
		fullKeys[i] = r.prefix + key
	}

	// Получаем все значения за один запрос
	values, err := r.client.MGet(ctx, fullKeys...).Result()
	if err != nil {
		return nil, err
	}

	results := make(map[string]*AIDetectionResult)
	for i, val := range values {
		if val == nil {
			continue
		}

		data, ok := val.(string)
		if !ok {
			continue
		}

		var result AIDetectionResult
		if err := json.Unmarshal([]byte(data), &result); err != nil {
			r.logger.Warn("Failed to unmarshal batch data",
				zap.String("key", keys[i]),
				zap.Error(err))
			continue
		}

		results[keys[i]] = &result
	}

	return results, nil
}

// WarmUp прогревает кэш популярными категориями
func (r *RedisCache) WarmUp(ctx context.Context, popularCategories []AIDetectionResult) error {
	for _, cat := range popularCategories {
		// Генерируем ключ на основе категории
		key := fmt.Sprintf("warmup:%d", cat.CategoryID)
		if err := r.Set(ctx, key, &cat); err != nil {
			r.logger.Warn("Failed to warm up cache",
				zap.Int32("categoryID", cat.CategoryID),
				zap.Error(err))
		}
	}
	return nil
}

// incrementHitCounter увеличивает счетчик попаданий
func (r *RedisCache) incrementHitCounter(ctx context.Context) {
	r.client.Incr(ctx, "ai:stats:hits")

	// Обновляем статистику за сегодня
	today := time.Now().Format("2006-01-02")
	r.client.Incr(ctx, fmt.Sprintf("ai:stats:daily:%s:hits", today))
}

// incrementMissCounter увеличивает счетчик промахов
func (r *RedisCache) incrementMissCounter(ctx context.Context) {
	r.client.Incr(ctx, "ai:stats:misses")

	// Обновляем статистику за сегодня
	today := time.Now().Format("2006-01-02")
	r.client.Incr(ctx, fmt.Sprintf("ai:stats:daily:%s:misses", today))
}

// GetDailyStats возвращает статистику за день
func (r *RedisCache) GetDailyStats(ctx context.Context, date string) (map[string]interface{}, error) {
	hits, _ := r.client.Get(ctx, fmt.Sprintf("ai:stats:daily:%s:hits", date)).Int64()
	misses, _ := r.client.Get(ctx, fmt.Sprintf("ai:stats:daily:%s:misses", date)).Int64()

	hitRate := float64(0)
	if hits+misses > 0 {
		hitRate = float64(hits) / float64(hits+misses) * 100
	}

	return map[string]interface{}{
		"date":    date,
		"hits":    hits,
		"misses":  misses,
		"hitRate": hitRate,
		"total":   hits + misses,
	}, nil
}

// Close закрывает соединение с Redis
func (r *RedisCache) Close() error {
	return r.client.Close()
}

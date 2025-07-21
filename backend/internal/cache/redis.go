package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

// RedisCache представляет клиент для работы с Redis кешем
type RedisCache struct {
	client *redis.Client
	logger *logrus.Logger
}

// NewRedisCache создает новый экземпляр Redis кеша
func NewRedisCache(url string, password string, db int, poolSize int, logger *logrus.Logger) (*RedisCache, error) {
	options := &redis.Options{
		Addr:     url,
		Password: password,
		DB:       db,
		PoolSize: poolSize,
	}

	client := redis.NewClient(options)

	// Проверяем подключение
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisCache{
		client: client,
		logger: logger,
	}, nil
}

// Get получает значение из кеша
func (r *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return ErrCacheMiss
	}
	if err != nil {
		r.logger.WithError(err).WithField("key", key).Error("Failed to get value from cache")
		return err
	}

	if err := json.Unmarshal([]byte(val), dest); err != nil {
		r.logger.WithError(err).WithField("key", key).Error("Failed to unmarshal cached value")
		return err
	}

	return nil
}

// Set сохраняет значение в кеш
func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		r.logger.WithError(err).WithField("key", key).Error("Failed to marshal value for cache")
		return err
	}

	if err := r.client.Set(ctx, key, data, ttl).Err(); err != nil {
		r.logger.WithError(err).WithField("key", key).Error("Failed to set value in cache")
		return err
	}

	return nil
}

// Delete удаляет значение из кеша
func (r *RedisCache) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	if err := r.client.Del(ctx, keys...).Err(); err != nil {
		r.logger.WithError(err).WithField("keys", keys).Error("Failed to delete keys from cache")
		return err
	}

	return nil
}

// DeletePattern удаляет все ключи по паттерну
func (r *RedisCache) DeletePattern(ctx context.Context, pattern string) error {
	var cursor uint64
	var deletedCount int

	for {
		keys, nextCursor, err := r.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			r.logger.WithError(err).WithField("pattern", pattern).Error("Failed to scan keys")
			return err
		}

		if len(keys) > 0 {
			if err := r.client.Del(ctx, keys...).Err(); err != nil {
				r.logger.WithError(err).WithField("keys", keys).Error("Failed to delete keys")
				return err
			}
			deletedCount += len(keys)
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	r.logger.WithFields(logrus.Fields{
		"pattern": pattern,
		"deleted": deletedCount,
	}).Debug("Deleted keys by pattern")

	return nil
}

// Close закрывает соединение с Redis
func (r *RedisCache) Close() error {
	return r.client.Close()
}

// Exists проверяет существование ключа
func (r *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	n, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		r.logger.WithError(err).WithField("key", key).Error("Failed to check key existence")
		return false, err
	}
	return n > 0, nil
}

// GetOrSet получает значение из кеша или устанавливает новое, если ключ не найден
func (r *RedisCache) GetOrSet(ctx context.Context, key string, dest interface{}, ttl time.Duration, loader func() (interface{}, error)) error {
	// Сначала пытаемся получить из кеша
	err := r.Get(ctx, key, dest)
	if err == nil {
		return nil
	}

	if err != ErrCacheMiss {
		// Если ошибка не cache miss, логируем но продолжаем
		r.logger.WithError(err).WithField("key", key).Warn("Cache get error, loading fresh data")
	}

	// Загружаем данные
	data, err := loader()
	if err != nil {
		return err
	}

	// Копируем загруженные данные в dest
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(dataBytes, dest); err != nil {
		return err
	}

	// Асинхронно сохраняем в кеш
	go func() {
		if err := r.Set(context.Background(), key, data, ttl); err != nil {
			r.logger.WithError(err).WithField("key", key).Warn("Failed to cache value")
		}
	}()

	return nil
}

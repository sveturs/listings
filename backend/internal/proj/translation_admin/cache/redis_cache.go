package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

// TranslationCache интерфейс для кеширования переводов
type TranslationCache interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	DeletePattern(ctx context.Context, pattern string) error
	Flush(ctx context.Context) error
}

// RedisTranslationCache реализация кеша переводов на Redis
type RedisTranslationCache struct {
	client *redis.Client
	prefix string
}

// NewRedisTranslationCache создает новый экземпляр кеша
func NewRedisTranslationCache(client *redis.Client) *RedisTranslationCache {
	return &RedisTranslationCache{
		client: client,
		prefix: "translations:",
	}
}

// BuildKey строит ключ для кеша
func (c *RedisTranslationCache) BuildKey(entityType string, entityID int64, language, fieldName string) string {
	return fmt.Sprintf("%s%s:%d:%s:%s", c.prefix, entityType, entityID, language, fieldName)
}

// BuildPatternKey строит паттерн для удаления группы ключей
func (c *RedisTranslationCache) BuildPatternKey(entityType string, entityID int64) string {
	return fmt.Sprintf("%s%s:%d:*", c.prefix, entityType, entityID)
}

// Get получает значение из кеша
func (c *RedisTranslationCache) Get(ctx context.Context, key string) (interface{}, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // Ключ не найден
	}
	if err != nil {
		log.Error().Err(err).Str("key", key).Msg("Failed to get from cache")
		return nil, err
	}

	var result interface{}
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		log.Error().Err(err).Str("key", key).Msg("Failed to unmarshal cached value")
		return nil, err
	}

	return result, nil
}

// GetTranslation получает перевод из кеша
func (c *RedisTranslationCache) GetTranslation(ctx context.Context, entityType string, entityID int64, language, fieldName string) (string, bool) {
	key := c.BuildKey(entityType, entityID, language, fieldName)
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", false
	}
	if err != nil {
		log.Error().Err(err).Str("key", key).Msg("Failed to get translation from cache")
		return "", false
	}
	return val, true
}

// Set сохраняет значение в кеш
func (c *RedisTranslationCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		log.Error().Err(err).Str("key", key).Msg("Failed to marshal value for cache")
		return err
	}

	if err := c.client.Set(ctx, key, data, expiration).Err(); err != nil {
		log.Error().Err(err).Str("key", key).Msg("Failed to set cache")
		return err
	}

	return nil
}

// SetTranslation сохраняет перевод в кеш
func (c *RedisTranslationCache) SetTranslation(ctx context.Context, entityType string, entityID int64, language, fieldName, translation string, expiration time.Duration) error {
	key := c.BuildKey(entityType, entityID, language, fieldName)
	if err := c.client.Set(ctx, key, translation, expiration).Err(); err != nil {
		log.Error().Err(err).Str("key", key).Msg("Failed to set translation cache")
		return err
	}
	return nil
}

// Delete удаляет значение из кеша
func (c *RedisTranslationCache) Delete(ctx context.Context, key string) error {
	if err := c.client.Del(ctx, key).Err(); err != nil {
		log.Error().Err(err).Str("key", key).Msg("Failed to delete from cache")
		return err
	}
	return nil
}

// DeletePattern удаляет все ключи по паттерну
func (c *RedisTranslationCache) DeletePattern(ctx context.Context, pattern string) error {
	iter := c.client.Scan(ctx, 0, pattern, 0).Iterator()
	var keys []string
	
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	
	if err := iter.Err(); err != nil {
		log.Error().Err(err).Str("pattern", pattern).Msg("Failed to scan keys")
		return err
	}

	if len(keys) > 0 {
		if err := c.client.Del(ctx, keys...).Err(); err != nil {
			log.Error().Err(err).Str("pattern", pattern).Msg("Failed to delete keys by pattern")
			return err
		}
		log.Info().Int("count", len(keys)).Str("pattern", pattern).Msg("Deleted cache keys")
	}

	return nil
}

// InvalidateEntity удаляет все переводы для сущности
func (c *RedisTranslationCache) InvalidateEntity(ctx context.Context, entityType string, entityID int64) error {
	pattern := c.BuildPatternKey(entityType, entityID)
	return c.DeletePattern(ctx, pattern)
}

// Flush очищает весь кеш переводов
func (c *RedisTranslationCache) Flush(ctx context.Context) error {
	pattern := fmt.Sprintf("%s*", c.prefix)
	return c.DeletePattern(ctx, pattern)
}

// BatchGet получает несколько переводов одним запросом
func (c *RedisTranslationCache) BatchGet(ctx context.Context, keys []string) (map[string]string, error) {
	if len(keys) == 0 {
		return make(map[string]string), nil
	}

	values, err := c.client.MGet(ctx, keys...).Result()
	if err != nil {
		log.Error().Err(err).Msg("Failed to batch get from cache")
		return nil, err
	}

	result := make(map[string]string)
	for i, val := range values {
		if val != nil {
			if str, ok := val.(string); ok {
				result[keys[i]] = str
			}
		}
	}

	return result, nil
}

// BatchSet сохраняет несколько переводов одним запросом
func (c *RedisTranslationCache) BatchSet(ctx context.Context, data map[string]string, expiration time.Duration) error {
	pipe := c.client.Pipeline()
	
	for key, value := range data {
		pipe.Set(ctx, key, value, expiration)
	}
	
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to batch set cache")
		return err
	}
	
	return nil
}

// WarmUp предзагружает переводы в кеш
func (c *RedisTranslationCache) WarmUp(ctx context.Context, translations map[string]map[string]string) error {
	data := make(map[string]string)
	
	for key, langs := range translations {
		for lang, text := range langs {
			cacheKey := fmt.Sprintf("%s%s:%s", c.prefix, key, lang)
			data[cacheKey] = text
		}
	}
	
	if len(data) > 0 {
		return c.BatchSet(ctx, data, 24*time.Hour)
	}
	
	return nil
}

// GetStats возвращает статистику кеша
func (c *RedisTranslationCache) GetStats(ctx context.Context) (map[string]interface{}, error) {
	pattern := fmt.Sprintf("%s*", c.prefix)
	iter := c.client.Scan(ctx, 0, pattern, 0).Iterator()
	
	count := 0
	for iter.Next(ctx) {
		count++
	}
	
	if err := iter.Err(); err != nil {
		return nil, err
	}
	
	info, err := c.client.Info(ctx, "memory").Result()
	if err != nil {
		return nil, err
	}
	
	return map[string]interface{}{
		"translation_keys": count,
		"memory_info":      info,
	}, nil
}
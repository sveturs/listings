package cache

import (
	"context"
	"time"
)

// Adapter адаптирует RedisCache к интерфейсу CacheInterface
type Adapter struct {
	cache *RedisCache
}

// NewAdapter создает новый адаптер
func NewAdapter(cache *RedisCache) *Adapter {
	return &Adapter{cache: cache}
}

// Get получает значение из кеша
func (a *Adapter) Get(ctx context.Context, key string, dest interface{}) error {
	return a.cache.Get(ctx, key, dest)
}

// Set сохраняет значение в кеш
func (a *Adapter) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return a.cache.Set(ctx, key, value, ttl)
}

// Delete удаляет значение из кеша
func (a *Adapter) Delete(ctx context.Context, keys ...string) error {
	return a.cache.Delete(ctx, keys...)
}

// DeletePattern удаляет все ключи по паттерну
func (a *Adapter) DeletePattern(ctx context.Context, pattern string) error {
	return a.cache.DeletePattern(ctx, pattern)
}

// GetOrSet получает значение из кеша или устанавливает новое
func (a *Adapter) GetOrSet(ctx context.Context, key string, dest interface{}, ttl time.Duration, loader func() (interface{}, error)) error {
	return a.cache.GetOrSet(ctx, key, dest, ttl, loader)
}

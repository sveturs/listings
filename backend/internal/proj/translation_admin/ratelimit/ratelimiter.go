package ratelimit

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
)

// RateLimiter интерфейс для ограничения частоты запросов
type RateLimiter interface {
	Allow(ctx context.Context, key string) (bool, error)
	AllowN(ctx context.Context, key string, n int) (bool, error)
	Reset(ctx context.Context, key string) error
	GetStatus(ctx context.Context, key string) (*Status, error)
}

// Status статус rate limiter'а
type Status struct {
	Limit     int       `json:"limit"`
	Remaining int       `json:"remaining"`
	Reset     time.Time `json:"reset"`
}

// Config конфигурация для rate limiter'а
type Config struct {
	// Requests per minute для разных провайдеров
	OpenAI    int `json:"openai"`
	Google    int `json:"google"`
	DeepL     int `json:"deepl"`
	Claude    int `json:"claude"`
	
	// Глобальные лимиты
	GlobalRPM int `json:"global_rpm"` // Requests per minute
	GlobalRPH int `json:"global_rph"` // Requests per hour
	GlobalRPD int `json:"global_rpd"` // Requests per day
	
	// Burst размер для token bucket
	BurstSize int `json:"burst_size"`
}

// DefaultConfig возвращает конфигурацию по умолчанию
func DefaultConfig() *Config {
	return &Config{
		OpenAI:    20,   // 20 запросов в минуту для OpenAI
		Google:    100,  // 100 запросов в минуту для Google
		DeepL:     50,   // 50 запросов в минуту для DeepL
		Claude:    30,   // 30 запросов в минуту для Claude
		GlobalRPM: 100,  // 100 запросов в минуту глобально
		GlobalRPH: 3000, // 3000 запросов в час
		GlobalRPD: 50000, // 50000 запросов в день
		BurstSize: 10,   // Burst до 10 запросов
	}
}

// InMemoryRateLimiter реализация rate limiter'а в памяти
type InMemoryRateLimiter struct {
	limiters map[string]*rate.Limiter
	mutex    sync.RWMutex
	config   *Config
}

// NewInMemoryRateLimiter создает новый in-memory rate limiter
func NewInMemoryRateLimiter(cfg *Config) *InMemoryRateLimiter {
	if cfg == nil {
		cfg = DefaultConfig()
	}
	
	return &InMemoryRateLimiter{
		limiters: make(map[string]*rate.Limiter),
		config:   cfg,
	}
}

// getLimiter получает или создает limiter для ключа
func (r *InMemoryRateLimiter) getLimiter(key string) *rate.Limiter {
	r.mutex.RLock()
	limiter, exists := r.limiters[key]
	r.mutex.RUnlock()
	
	if exists {
		return limiter
	}
	
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	// Проверяем еще раз после блокировки
	if limiter, exists := r.limiters[key]; exists {
		return limiter
	}
	
	// Определяем лимит для провайдера
	var limit int
	switch key {
	case "openai":
		limit = r.config.OpenAI
	case "google":
		limit = r.config.Google
	case "deepl":
		limit = r.config.DeepL
	case "claude":
		limit = r.config.Claude
	case "global":
		limit = r.config.GlobalRPM
	default:
		limit = r.config.GlobalRPM
	}
	
	// Создаем новый limiter (limit per minute)
	limiter = rate.NewLimiter(rate.Limit(float64(limit)/60.0), r.config.BurstSize)
	r.limiters[key] = limiter
	
	return limiter
}

// Allow проверяет, можно ли выполнить запрос
func (r *InMemoryRateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	limiter := r.getLimiter(key)
	return limiter.Allow(), nil
}

// AllowN проверяет, можно ли выполнить N запросов
func (r *InMemoryRateLimiter) AllowN(ctx context.Context, key string, n int) (bool, error) {
	limiter := r.getLimiter(key)
	return limiter.AllowN(time.Now(), n), nil
}

// Reset сбрасывает лимит для ключа
func (r *InMemoryRateLimiter) Reset(ctx context.Context, key string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	delete(r.limiters, key)
	return nil
}

// GetStatus возвращает статус лимита
func (r *InMemoryRateLimiter) GetStatus(ctx context.Context, key string) (*Status, error) {
	limiter := r.getLimiter(key)
	
	// Приблизительный расчет оставшихся запросов
	tokens := int(limiter.Tokens())
	
	var limit int
	switch key {
	case "openai":
		limit = r.config.OpenAI
	case "google":
		limit = r.config.Google
	case "deepl":
		limit = r.config.DeepL
	case "claude":
		limit = r.config.Claude
	default:
		limit = r.config.GlobalRPM
	}
	
	return &Status{
		Limit:     limit,
		Remaining: tokens,
		Reset:     time.Now().Add(time.Minute),
	}, nil
}

// RedisRateLimiter реализация rate limiter'а через Redis
type RedisRateLimiter struct {
	client *redis.Client
	config *Config
}

// NewRedisRateLimiter создает новый Redis-based rate limiter
func NewRedisRateLimiter(client *redis.Client, cfg *Config) *RedisRateLimiter {
	if cfg == nil {
		cfg = DefaultConfig()
	}
	
	return &RedisRateLimiter{
		client: client,
		config: cfg,
	}
}

// Allow проверяет и увеличивает счетчик
func (r *RedisRateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	return r.AllowN(ctx, key, 1)
}

// AllowN проверяет и увеличивает счетчик на N
func (r *RedisRateLimiter) AllowN(ctx context.Context, key string, n int) (bool, error) {
	// Определяем лимит для провайдера
	var limit int
	switch key {
	case "openai":
		limit = r.config.OpenAI
	case "google":
		limit = r.config.Google
	case "deepl":
		limit = r.config.DeepL
	case "claude":
		limit = r.config.Claude
	default:
		limit = r.config.GlobalRPM
	}
	
	// Используем sliding window algorithm
	now := time.Now()
	windowStart := now.Add(-time.Minute).Unix()
	windowKey := fmt.Sprintf("ratelimit:%s:window", key)
	
	pipe := r.client.Pipeline()
	
	// Удаляем старые записи
	pipe.ZRemRangeByScore(ctx, windowKey, "0", fmt.Sprintf("%d", windowStart))
	
	// Считаем текущие записи в окне
	countCmd := pipe.ZCard(ctx, windowKey)
	
	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}
	
	currentCount := countCmd.Val()
	
	// Проверяем лимит
	if int(currentCount)+n > limit {
		return false, nil
	}
	
	// Добавляем новые записи
	members := make([]redis.Z, n)
	for i := 0; i < n; i++ {
		members[i] = redis.Z{
			Score:  float64(now.UnixNano()),
			Member: fmt.Sprintf("%d:%d", now.UnixNano(), i),
		}
	}
	
	pipe = r.client.Pipeline()
	pipe.ZAdd(ctx, windowKey, members...)
	pipe.Expire(ctx, windowKey, time.Minute*2)
	
	_, err = pipe.Exec(ctx)
	if err != nil {
		return false, err
	}
	
	return true, nil
}

// Reset сбрасывает счетчик для ключа
func (r *RedisRateLimiter) Reset(ctx context.Context, key string) error {
	windowKey := fmt.Sprintf("ratelimit:%s:window", key)
	return r.client.Del(ctx, windowKey).Err()
}

// GetStatus возвращает статус лимита
func (r *RedisRateLimiter) GetStatus(ctx context.Context, key string) (*Status, error) {
	// Определяем лимит для провайдера
	var limit int
	switch key {
	case "openai":
		limit = r.config.OpenAI
	case "google":
		limit = r.config.Google
	case "deepl":
		limit = r.config.DeepL
	case "claude":
		limit = r.config.Claude
	default:
		limit = r.config.GlobalRPM
	}
	
	// Считаем текущие записи в окне
	now := time.Now()
	windowStart := now.Add(-time.Minute).Unix()
	windowKey := fmt.Sprintf("ratelimit:%s:window", key)
	
	// Удаляем старые и считаем оставшиеся
	pipe := r.client.Pipeline()
	pipe.ZRemRangeByScore(ctx, windowKey, "0", fmt.Sprintf("%d", windowStart))
	countCmd := pipe.ZCard(ctx, windowKey)
	
	_, err := pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}
	
	currentCount := int(countCmd.Val())
	remaining := limit - currentCount
	if remaining < 0 {
		remaining = 0
	}
	
	return &Status{
		Limit:     limit,
		Remaining: remaining,
		Reset:     now.Add(time.Minute),
	}, nil
}

// MultiProviderRateLimiter объединяет лимиты для разных провайдеров
type MultiProviderRateLimiter struct {
	limiter RateLimiter
}

// NewMultiProviderRateLimiter создает новый multi-provider rate limiter
func NewMultiProviderRateLimiter(redisClient *redis.Client, cfg *Config) *MultiProviderRateLimiter {
	var limiter RateLimiter
	
	if redisClient != nil {
		limiter = NewRedisRateLimiter(redisClient, cfg)
		log.Info().Msg("Using Redis-based rate limiter for AI translations")
	} else {
		limiter = NewInMemoryRateLimiter(cfg)
		log.Info().Msg("Using in-memory rate limiter for AI translations")
	}
	
	return &MultiProviderRateLimiter{
		limiter: limiter,
	}
}

// CheckProviderLimit проверяет лимит для конкретного провайдера
func (m *MultiProviderRateLimiter) CheckProviderLimit(ctx context.Context, provider string) (bool, error) {
	// Проверяем лимит провайдера
	providerAllowed, err := m.limiter.Allow(ctx, provider)
	if err != nil {
		log.Error().Err(err).Str("provider", provider).Msg("Failed to check provider rate limit")
		return false, err
	}
	
	if !providerAllowed {
		log.Warn().Str("provider", provider).Msg("Provider rate limit exceeded")
		return false, nil
	}
	
	// Проверяем глобальный лимит
	globalAllowed, err := m.limiter.Allow(ctx, "global")
	if err != nil {
		log.Error().Err(err).Msg("Failed to check global rate limit")
		// Откатываем счетчик провайдера
		_ = m.limiter.Reset(ctx, provider)
		return false, err
	}
	
	if !globalAllowed {
		log.Warn().Msg("Global rate limit exceeded")
		// Откатываем счетчик провайдера
		_ = m.limiter.Reset(ctx, provider)
		return false, nil
	}
	
	return true, nil
}

// GetProviderStatus возвращает статус лимитов для провайдера
func (m *MultiProviderRateLimiter) GetProviderStatus(ctx context.Context, provider string) (map[string]*Status, error) {
	result := make(map[string]*Status)
	
	// Статус провайдера
	providerStatus, err := m.limiter.GetStatus(ctx, provider)
	if err != nil {
		return nil, err
	}
	result[provider] = providerStatus
	
	// Глобальный статус
	globalStatus, err := m.limiter.GetStatus(ctx, "global")
	if err != nil {
		return nil, err
	}
	result["global"] = globalStatus
	
	return result, nil
}

// GetAllProvidersStatus возвращает статус всех провайдеров
func (m *MultiProviderRateLimiter) GetAllProvidersStatus(ctx context.Context) (map[string]*Status, error) {
	providers := []string{"openai", "google", "deepl", "claude", "global"}
	result := make(map[string]*Status)
	
	for _, provider := range providers {
		status, err := m.limiter.GetStatus(ctx, provider)
		if err != nil {
			log.Error().Err(err).Str("provider", provider).Msg("Failed to get provider status")
			continue
		}
		result[provider] = status
	}
	
	return result, nil
}
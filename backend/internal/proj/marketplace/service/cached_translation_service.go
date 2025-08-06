package service

import (
	"context"
	"crypto/sha256"

	//	"encoding/json"
	"fmt"
	"sync"
	"time"

	"backend/internal/cache"
	"backend/internal/logger"
)

// CachedTranslationService обертка над существующим сервисом перевода с улучшенным кэшированием
type CachedTranslationService struct {
	underlying    TranslationServiceInterface
	redisCache    *cache.RedisCache
	localCache    *sync.Map // быстрый локальный кэш для горячих переводов
	localExpiry   *sync.Map // время истечения для локального кэша
	cachePrefix   string
	cacheTTL      time.Duration
	localCacheTTL time.Duration
}

// CachedTranslationConfig конфигурация для кэшированного сервиса переводов
type CachedTranslationConfig struct {
	CachePrefix      string        // префикс для ключей кэша
	RedisCacheTTL    time.Duration // время жизни в Redis (долгосрочное хранение)
	LocalCacheTTL    time.Duration // время жизни в локальном кэше (быстрый доступ)
	LocalCacheSize   int           // размер локального кэша
	EnableLocalCache bool          // включить локальный кэш
}

// DefaultCachedTranslationConfig возвращает конфигурацию по умолчанию
func DefaultCachedTranslationConfig() CachedTranslationConfig {
	return CachedTranslationConfig{
		CachePrefix:      "translation:",
		RedisCacheTTL:    24 * time.Hour,   // 24 часа в Redis
		LocalCacheTTL:    30 * time.Minute, // 30 минут в локальном кэше
		LocalCacheSize:   1000,             // 1000 записей в локальном кэше
		EnableLocalCache: true,
	}
}

// NewCachedTranslationService создает новый кэшированный сервис переводов
func NewCachedTranslationService(
	underlying TranslationServiceInterface,
	redisCache *cache.RedisCache,
	config ...CachedTranslationConfig,
) *CachedTranslationService {
	cfg := DefaultCachedTranslationConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	service := &CachedTranslationService{
		underlying:    underlying,
		redisCache:    redisCache,
		cachePrefix:   cfg.CachePrefix,
		cacheTTL:      cfg.RedisCacheTTL,
		localCacheTTL: cfg.LocalCacheTTL,
		localCache:    &sync.Map{},
		localExpiry:   &sync.Map{},
	}

	return service
}

// generateCacheKey создает уникальный ключ для кэширования
func (s *CachedTranslationService) generateCacheKey(text, sourceLanguage, targetLanguage, context, fieldName string) string {
	// Используем SHA256 для создания безопасного ключа
	input := fmt.Sprintf("%s|%s|%s|%s|%s", text, sourceLanguage, targetLanguage, context, fieldName)
	hash := sha256.Sum256([]byte(input))
	return fmt.Sprintf("%s%x", s.cachePrefix, hash)
}

// TranslateWithContext реализует интерфейс с улучшенным кэшированием
func (s *CachedTranslationService) TranslateWithContext(ctx context.Context, text, sourceLanguage, targetLanguage, context, fieldName string) (string, error) {
	cacheKey := s.generateCacheKey(text, sourceLanguage, targetLanguage, context, fieldName)

	// 1. Проверяем локальный кэш (самый быстрый)
	if s.localCache != nil {
		if cached, found := s.getFromLocalCache(cacheKey); found {
			logger.Debug().
				Str("cache_hit", "local").
				Str("key", cacheKey).
				Msg("Translation found in local cache")
			return cached, nil
		}
	}

	// 2. Проверяем Redis кэш
	if s.redisCache != nil {
		var cached string
		err := s.redisCache.Get(ctx, cacheKey, &cached)
		if err == nil && cached != "" {
			// Сохраняем в локальный кэш для быстрого доступа
			if s.localCache != nil {
				s.setToLocalCache(cacheKey, cached)
			}

			logger.Debug().
				Str("cache_hit", "redis").
				Str("key", cacheKey).
				Msg("Translation found in Redis cache")
			return cached, nil
		}
		if err != nil {
			logger.Warn().Err(err).Msg("Redis cache error, proceeding without cache")
		}
	}

	// 3. Кэш промах - выполняем перевод
	logger.Debug().
		Str("cache_miss", "all").
		Str("text", text[:min(50, len(text))]).
		Str("source", sourceLanguage).
		Str("target", targetLanguage).
		Msg("Cache miss, translating")

	translated, err := s.underlying.TranslateWithContext(ctx, text, sourceLanguage, targetLanguage, context, fieldName)
	if err != nil {
		return "", fmt.Errorf("translation failed: %w", err)
	}

	// 4. Сохраняем результат в кэши асинхронно
	go s.cacheTranslation(ctx, cacheKey, translated)

	return translated, nil
}

// cacheTranslation сохраняет перевод в кэши асинхронно
func (s *CachedTranslationService) cacheTranslation(ctx context.Context, key, value string) {
	// Устанавливаем таймаут для операций кэширования
	cacheCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Сохраняем в локальный кэш (быстро)
	if s.localCache != nil {
		s.setToLocalCache(key, value)
		logger.Debug().Str("key", key).Msg("Translation cached locally")
	}

	// Сохраняем в Redis (может быть медленно)
	if s.redisCache != nil {
		if err := s.redisCache.Set(cacheCtx, key, value, s.cacheTTL); err != nil {
			logger.Error().Err(err).Str("key", key).Msg("Failed to cache translation in Redis")
		} else {
			logger.Debug().Str("key", key).Msg("Translation cached in Redis")
		}
	}
}

// Translate реализует базовый интерфейс перевода
func (s *CachedTranslationService) Translate(ctx context.Context, text, sourceLanguage, targetLanguage string) (string, error) {
	return s.TranslateWithContext(ctx, text, sourceLanguage, targetLanguage, "", "")
}

// DetectLanguage проксирует вызов к underlying сервису
func (s *CachedTranslationService) DetectLanguage(ctx context.Context, text string) (string, float64, error) {
	return s.underlying.DetectLanguage(ctx, text)
}

// ClearCache очищает все кэши
func (s *CachedTranslationService) ClearCache(ctx context.Context) error {
	// Очищаем локальный кэш
	if s.localCache != nil {
		s.clearLocalCache()
		logger.Info().Msg("Local translation cache cleared")
	}

	// Очищаем Redis кэш (только ключи с нашим префиксом)
	if s.redisCache != nil {
		pattern := s.cachePrefix + "*"
		client := s.redisCache.GetClient()

		keys, err := client.Keys(ctx, pattern).Result()
		if err != nil {
			return fmt.Errorf("failed to get cache keys: %w", err)
		}

		if len(keys) > 0 {
			if err := client.Del(ctx, keys...).Err(); err != nil {
				return fmt.Errorf("failed to delete cache keys: %w", err)
			}
			logger.Info().Int("count", len(keys)).Msg("Redis translation cache cleared")
		}
	}

	return nil
}

// GetCacheStats возвращает статистику кэша
func (s *CachedTranslationService) GetCacheStats(ctx context.Context) map[string]interface{} {
	stats := make(map[string]interface{})

	// Статистика локального кэша
	if s.localCache != nil {
		// Подсчитываем размер локального кэша
		size := 0
		s.localCache.Range(func(key, value interface{}) bool {
			size++
			return true
		})

		stats["local_cache"] = map[string]interface{}{
			"enabled": true,
			"size":    size,
			"ttl":     s.localCacheTTL.String(),
		}
	} else {
		stats["local_cache"] = map[string]interface{}{
			"enabled": false,
		}
	}

	// Статистика Redis кэша
	if s.redisCache != nil {
		client := s.redisCache.GetClient()
		pattern := s.cachePrefix + "*"

		keys, err := client.Keys(ctx, pattern).Result()
		redisStats := map[string]interface{}{
			"enabled": true,
			"ttl":     s.cacheTTL.String(),
		}

		if err == nil {
			redisStats["keys_count"] = len(keys)
		} else {
			redisStats["error"] = err.Error()
		}

		stats["redis_cache"] = redisStats
	} else {
		stats["redis_cache"] = map[string]interface{}{
			"enabled": false,
		}
	}

	return stats
}

// WarmupCache предварительно загружает популярные переводы в кэш
func (s *CachedTranslationService) WarmupCache(ctx context.Context, translations []struct {
	Text           string
	SourceLanguage string
	TargetLanguage string
	Context        string
	FieldName      string
},
) error {
	logger.Info().Int("count", len(translations)).Msg("Starting translation cache warmup")

	for i, t := range translations {
		// Выполняем перевод (что автоматически сохранит его в кэш)
		_, err := s.TranslateWithContext(ctx, t.Text, t.SourceLanguage, t.TargetLanguage, t.Context, t.FieldName)
		if err != nil {
			logger.Error().
				Err(err).
				Int("index", i).
				Str("text", t.Text[:min(50, len(t.Text))]).
				Msg("Failed to warmup translation")
			continue
		}

		// Небольшая пауза между переводами чтобы не перегружать API
		if i%10 == 0 && i > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(100 * time.Millisecond):
			}
		}
	}

	logger.Info().Int("count", len(translations)).Msg("Translation cache warmup completed")
	return nil
}

// getFromLocalCache получает значение из локального кэша с проверкой TTL
func (s *CachedTranslationService) getFromLocalCache(key string) (string, bool) {
	if value, ok := s.localCache.Load(key); ok {
		// Проверяем не истекло ли время жизни
		if expiry, ok := s.localExpiry.Load(key); ok {
			if time.Now().After(expiry.(time.Time)) {
				// Время истекло, удаляем
				s.localCache.Delete(key)
				s.localExpiry.Delete(key)
				return "", false
			}
		}
		return value.(string), true
	}
	return "", false
}

// setToLocalCache сохраняет значение в локальный кэш с TTL
func (s *CachedTranslationService) setToLocalCache(key, value string) {
	s.localCache.Store(key, value)
	s.localExpiry.Store(key, time.Now().Add(s.localCacheTTL))
}

// clearLocalCache очищает локальный кэш
func (s *CachedTranslationService) clearLocalCache() {
	s.localCache = &sync.Map{}
	s.localExpiry = &sync.Map{}
}

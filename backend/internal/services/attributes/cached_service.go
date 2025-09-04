package attributes

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"backend/internal/domain/models"
)

// CachedAttributeService обеспечивает кэширование для сервиса атрибутов
type CachedAttributeService struct {
	redis       *redis.Client
	repository  Repository
	logger      *zap.Logger
	cacheTTL    time.Duration
	cachePrefix string
}

// NewCachedAttributeService создает новый экземпляр кэшированного сервиса атрибутов
func NewCachedAttributeService(
	redis *redis.Client,
	repository Repository,
	logger *zap.Logger,
) *CachedAttributeService {
	return &CachedAttributeService{
		redis:       redis,
		repository:  repository,
		logger:      logger,
		cacheTTL:    1 * time.Hour, // Кэш на 1 час
		cachePrefix: "unified_attr:",
	}
}

// Repository интерфейс для репозитория атрибутов
type Repository interface {
	GetCategoryAttributes(ctx context.Context, categoryID int64) ([]models.UnifiedAttribute, error)
	GetAttributeById(ctx context.Context, attributeID int64) (*models.UnifiedAttribute, error)
	GetPopularValues(ctx context.Context, attributeID int64, limit int) ([]string, error)
}

// GetCategoryAttributes получает атрибуты категории с кэшированием
func (s *CachedAttributeService) GetCategoryAttributes(ctx context.Context, categoryID int64) ([]models.UnifiedAttribute, error) {
	// Формируем ключ кэша
	cacheKey := fmt.Sprintf("%scategory:%d", s.cachePrefix, categoryID)

	// Пробуем получить из кэша
	cachedData, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		// Десериализуем из кэша
		var attributes []models.UnifiedAttribute
		if err := json.Unmarshal([]byte(cachedData), &attributes); err == nil {
			s.logger.Debug("Cache hit for category attributes",
				zap.Int64("category_id", categoryID),
				zap.Int("attributes_count", len(attributes)))
			return attributes, nil
		}
		// Если ошибка десериализации, продолжаем с БД
		s.logger.Warn("Failed to deserialize cached category attributes",
			zap.Error(err), zap.Int64("category_id", categoryID))
	}

	// Если кэш пустой или ошибка, получаем из БД
	attributes, err := s.repository.GetCategoryAttributes(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get category attributes from repository: %w", err)
	}

	// Сохраняем в кэш (асинхронно, чтобы не замедлять ответ)
	go s.cacheAttributesList(context.Background(), cacheKey, attributes)

	s.logger.Debug("Cache miss for category attributes, loaded from DB",
		zap.Int64("category_id", categoryID),
		zap.Int("attributes_count", len(attributes)))

	return attributes, nil
}

// GetAttributeById получает атрибут по ID с кэшированием
func (s *CachedAttributeService) GetAttributeById(ctx context.Context, attributeID int64) (*models.UnifiedAttribute, error) {
	// Формируем ключ кэша
	cacheKey := fmt.Sprintf("%sattr:%d", s.cachePrefix, attributeID)

	// Пробуем получить из кэша
	cachedData, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		// Десериализуем из кэша
		var attribute models.UnifiedAttribute
		if err := json.Unmarshal([]byte(cachedData), &attribute); err == nil {
			s.logger.Debug("Cache hit for attribute",
				zap.Int64("attribute_id", attributeID))
			return &attribute, nil
		}
		// Если ошибка десериализации, продолжаем с БД
		s.logger.Warn("Failed to deserialize cached attribute",
			zap.Error(err), zap.Int64("attribute_id", attributeID))
	}

	// Если кэш пустой или ошибка, получаем из БД
	attribute, err := s.repository.GetAttributeById(ctx, attributeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get attribute from repository: %w", err)
	}

	// Сохраняем в кэш (асинхронно)
	go s.cacheAttribute(context.Background(), cacheKey, *attribute)

	s.logger.Debug("Cache miss for attribute, loaded from DB",
		zap.Int64("attribute_id", attributeID))

	return attribute, nil
}

// GetPopularValues получает популярные значения атрибута с кэшированием
func (s *CachedAttributeService) GetPopularValues(ctx context.Context, attributeID int64, limit int) ([]string, error) {
	// Формируем ключ кэша
	cacheKey := fmt.Sprintf("%spopular:%d:%d", s.cachePrefix, attributeID, limit)

	// Пробуем получить из кэша
	cachedData, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		// Десериализуем из кэша
		var values []string
		if err := json.Unmarshal([]byte(cachedData), &values); err == nil {
			s.logger.Debug("Cache hit for popular values",
				zap.Int64("attribute_id", attributeID),
				zap.Int("limit", limit),
				zap.Int("values_count", len(values)))
			return values, nil
		}
		// Если ошибка десериализации, продолжаем с БД
		s.logger.Warn("Failed to deserialize cached popular values",
			zap.Error(err), zap.Int64("attribute_id", attributeID))
	}

	// Если кэш пустой или ошибка, получаем из БД
	values, err := s.repository.GetPopularValues(ctx, attributeID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get popular values from repository: %w", err)
	}

	// Сохраняем в кэш (асинхронно) с более коротким TTL для популярных значений
	go s.cachePopularValues(context.Background(), cacheKey, values)

	s.logger.Debug("Cache miss for popular values, loaded from DB",
		zap.Int64("attribute_id", attributeID),
		zap.Int("limit", limit),
		zap.Int("values_count", len(values)))

	return values, nil
}

// InvalidateCategory инвалидирует кэш для категории
func (s *CachedAttributeService) InvalidateCategory(categoryID int64) error {
	cacheKey := fmt.Sprintf("%scategory:%d", s.cachePrefix, categoryID)
	return s.redis.Del(context.Background(), cacheKey).Err()
}

// InvalidateAttribute инвалидирует кэш для атрибута
func (s *CachedAttributeService) InvalidateAttribute(attributeID int64) error {
	// Инвалидируем кэш атрибута
	attrKey := fmt.Sprintf("%sattr:%d", s.cachePrefix, attributeID)

	// Инвалидируем кэш популярных значений (паттерн поиск)
	popularPattern := fmt.Sprintf("%spopular:%d:*", s.cachePrefix, attributeID)

	// Удаляем ключи по паттерну
	ctx := context.Background()
	keys, err := s.redis.Keys(ctx, popularPattern).Result()
	if err != nil {
		return fmt.Errorf("failed to get keys for pattern %s: %w", popularPattern, err)
	}

	// Добавляем ключ атрибута к списку для удаления
	keys = append(keys, attrKey)

	if len(keys) > 0 {
		return s.redis.Del(ctx, keys...).Err()
	}

	return nil
}

// CacheStats возвращает статистику кэша для мониторинга
func (s *CachedAttributeService) CacheStats() map[string]interface{} {
	ctx := context.Background()

	// Получаем общую информацию о Redis
	info, err := s.redis.Info(ctx, "memory").Result()
	if err != nil {
		s.logger.Error("Failed to get Redis info", zap.Error(err))
		return map[string]interface{}{
			"error": err.Error(),
		}
	}

	// Подсчитываем количество ключей с нашим префиксом
	keys, err := s.redis.Keys(ctx, s.cachePrefix+"*").Result()
	if err != nil {
		s.logger.Error("Failed to count cache keys", zap.Error(err))
		return map[string]interface{}{
			"error": err.Error(),
		}
	}

	return map[string]interface{}{
		"redis_memory_info": info,
		"cached_keys_count": len(keys),
		"cache_prefix":      s.cachePrefix,
		"cache_ttl":         s.cacheTTL.String(),
	}
}

// Вспомогательные методы для асинхронного кэширования

func (s *CachedAttributeService) cacheAttributesList(ctx context.Context, key string, attributes []models.UnifiedAttribute) {
	data, err := json.Marshal(attributes)
	if err != nil {
		s.logger.Error("Failed to serialize attributes for cache", zap.Error(err))
		return
	}

	if err := s.redis.Set(ctx, key, data, s.cacheTTL).Err(); err != nil {
		s.logger.Error("Failed to cache attributes list",
			zap.Error(err), zap.String("key", key))
	}
}

func (s *CachedAttributeService) cacheAttribute(ctx context.Context, key string, attribute models.UnifiedAttribute) {
	data, err := json.Marshal(attribute)
	if err != nil {
		s.logger.Error("Failed to serialize attribute for cache", zap.Error(err))
		return
	}

	if err := s.redis.Set(ctx, key, data, s.cacheTTL).Err(); err != nil {
		s.logger.Error("Failed to cache attribute",
			zap.Error(err), zap.String("key", key))
	}
}

func (s *CachedAttributeService) cachePopularValues(ctx context.Context, key string, values []string) {
	data, err := json.Marshal(values)
	if err != nil {
		s.logger.Error("Failed to serialize popular values for cache", zap.Error(err))
		return
	}

	// Популярные значения кэшируем на меньшее время (30 минут)
	ttl := 30 * time.Minute
	if err := s.redis.Set(ctx, key, data, ttl).Err(); err != nil {
		s.logger.Error("Failed to cache popular values",
			zap.Error(err), zap.String("key", key))
	}
}

package translation_admin

import (
	"context"
	"fmt"
	"time"

	"backend/internal/domain/models"
)

const (
	defaultCacheExpiration = 24 * time.Hour
	shortCacheExpiration   = 1 * time.Hour
)

// GetTranslationWithCache получает перевод с использованием кеша
func (s *Service) GetTranslationWithCache(ctx context.Context, entityType string, entityID int64, language, fieldName string) (string, error) {
	// Если кеш не настроен, получаем напрямую из БД
	if s.cache == nil {
		return s.getTranslationFromDB(ctx, entityType, entityID, language, fieldName)
	}

	// Пытаемся получить из кеша
	if translation, found := s.cache.GetTranslation(ctx, entityType, entityID, language, fieldName); found {
		s.logger.Debug().
			Str("entity_type", entityType).
			Int64("entity_id", entityID).
			Str("language", language).
			Str("field", fieldName).
			Msg("Translation cache hit")
		return translation, nil
	}

	// Не найдено в кеше, получаем из БД
	translation, err := s.getTranslationFromDB(ctx, entityType, entityID, language, fieldName)
	if err != nil {
		return "", err
	}

	// Сохраняем в кеш для будущих запросов
	if translation != "" {
		if err := s.cache.SetTranslation(ctx, entityType, entityID, language, fieldName, translation, defaultCacheExpiration); err != nil {
			s.logger.Warn().Err(err).Msg("Failed to cache translation")
		}
	}

	return translation, nil
}

// getTranslationFromDB получает перевод из базы данных
func (s *Service) getTranslationFromDB(ctx context.Context, entityType string, entityID int64, language, fieldName string) (string, error) {
	filters := map[string]interface{}{
		"entity_type": entityType,
		"entity_id":   entityID,
		"language":    language,
		"field_name":  fieldName,
	}

	translations, err := s.translationRepo.GetTranslations(ctx, filters)
	if err != nil {
		return "", err
	}

	if len(translations) > 0 {
		return translations[0].TranslatedText, nil
	}

	return "", nil
}

// SaveTranslationWithCache сохраняет перевод и обновляет кеш
func (s *Service) SaveTranslationWithCache(ctx context.Context, translation *models.Translation) error {
	// Сохраняем в БД
	var err error
	if translation.ID > 0 {
		err = s.translationRepo.UpdateTranslation(ctx, translation)
	} else {
		err = s.translationRepo.CreateTranslation(ctx, translation)
	}

	if err != nil {
		return err
	}

	// Обновляем кеш если он настроен
	if s.cache != nil {
		cacheErr := s.cache.SetTranslation(
			ctx,
			translation.EntityType,
			int64(translation.EntityID),
			translation.Language,
			translation.FieldName,
			translation.TranslatedText,
			defaultCacheExpiration,
		)
		if cacheErr != nil {
			s.logger.Warn().Err(cacheErr).Msg("Failed to update translation cache")
		}
	}

	return nil
}

// InvalidateTranslationCache инвалидирует кеш для сущности
func (s *Service) InvalidateTranslationCache(ctx context.Context, entityType string, entityID int64) error {
	if s.cache == nil {
		return nil
	}

	return s.cache.InvalidateEntity(ctx, entityType, entityID)
}

// GetBatchTranslationsWithCache получает несколько переводов с использованием кеша
func (s *Service) GetBatchTranslationsWithCache(ctx context.Context, requests []TranslationRequest) (map[string]string, error) {
	result := make(map[string]string)
	
	if s.cache == nil {
		// Если кеш не настроен, получаем все из БД
		return s.getBatchTranslationsFromDB(ctx, requests)
	}

	// Формируем ключи для кеша
	var cacheKeys []string
	keyToRequest := make(map[string]TranslationRequest)
	
	for _, req := range requests {
		key := s.cache.BuildKey(req.EntityType, req.EntityID, req.Language, req.FieldName)
		cacheKeys = append(cacheKeys, key)
		keyToRequest[key] = req
	}

	// Пытаемся получить из кеша
	cached, err := s.cache.BatchGet(ctx, cacheKeys)
	if err != nil {
		s.logger.Warn().Err(err).Msg("Failed to batch get from cache")
	}

	// Определяем какие переводы нужно получить из БД
	var missingRequests []TranslationRequest
	for key, req := range keyToRequest {
		if translation, found := cached[key]; found {
			resultKey := fmt.Sprintf("%s:%d:%s:%s", req.EntityType, req.EntityID, req.Language, req.FieldName)
			result[resultKey] = translation
		} else {
			missingRequests = append(missingRequests, req)
		}
	}

	// Получаем недостающие переводы из БД
	if len(missingRequests) > 0 {
		dbTranslations, err := s.getBatchTranslationsFromDB(ctx, missingRequests)
		if err != nil {
			return result, err
		}

		// Добавляем в результат и кешируем
		toCache := make(map[string]string)
		for key, translation := range dbTranslations {
			result[key] = translation
			
			// Парсим ключ для кеширования
			var entityType, fieldName, language string
			var entityID int64
			fmt.Sscanf(key, "%[^:]:%d:%[^:]:%s", &entityType, &entityID, &language, &fieldName)
			
			cacheKey := s.cache.BuildKey(entityType, entityID, language, fieldName)
			toCache[cacheKey] = translation
		}

		// Сохраняем в кеш
		if len(toCache) > 0 {
			if err := s.cache.BatchSet(ctx, toCache, defaultCacheExpiration); err != nil {
				s.logger.Warn().Err(err).Msg("Failed to batch set cache")
			}
		}
	}

	return result, nil
}

// TranslationRequest структура запроса перевода
type TranslationRequest struct {
	EntityType string
	EntityID   int64
	Language   string
	FieldName  string
}

// getBatchTranslationsFromDB получает несколько переводов из БД
func (s *Service) getBatchTranslationsFromDB(ctx context.Context, requests []TranslationRequest) (map[string]string, error) {
	result := make(map[string]string)
	
	// Группируем запросы для оптимизации
	for _, req := range requests {
		filters := map[string]interface{}{
			"entity_type": req.EntityType,
			"entity_id":   req.EntityID,
			"language":    req.Language,
			"field_name":  req.FieldName,
		}
		
		translations, err := s.translationRepo.GetTranslations(ctx, filters)
		if err != nil {
			s.logger.Error().Err(err).Interface("filters", filters).Msg("Failed to get translation from DB")
			continue
		}
		
		if len(translations) > 0 {
			key := fmt.Sprintf("%s:%d:%s:%s", req.EntityType, req.EntityID, req.Language, req.FieldName)
			result[key] = translations[0].TranslatedText
		}
	}
	
	return result, nil
}

// WarmUpCache предзагружает часто используемые переводы в кеш
func (s *Service) WarmUpCache(ctx context.Context) error {
	if s.cache == nil {
		return nil
	}

	s.logger.Info().Msg("Starting translation cache warm-up")

	// Получаем самые популярные переводы (например, категории)
	filters := map[string]interface{}{
		"entity_type": "category",
		"is_verified": true,
	}

	translations, err := s.translationRepo.GetTranslations(ctx, filters)
	if err != nil {
		return fmt.Errorf("failed to get translations for warm-up: %w", err)
	}

	// Группируем переводы для batch set
	toCache := make(map[string]string)
	for _, t := range translations {
		key := s.cache.BuildKey(t.EntityType, int64(t.EntityID), t.Language, t.FieldName)
		toCache[key] = t.TranslatedText
	}

	if len(toCache) > 0 {
		if err := s.cache.BatchSet(ctx, toCache, defaultCacheExpiration); err != nil {
			return fmt.Errorf("failed to warm-up cache: %w", err)
		}
	}

	s.logger.Info().Int("count", len(toCache)).Msg("Translation cache warmed up")
	return nil
}

// GetCacheStats возвращает статистику кеша
func (s *Service) GetCacheStats(ctx context.Context) (map[string]interface{}, error) {
	if s.cache == nil {
		return map[string]interface{}{
			"enabled": false,
		}, nil
	}

	stats, err := s.cache.GetStats(ctx)
	if err != nil {
		return nil, err
	}

	stats["enabled"] = true
	return stats, nil
}

// FlushCache очищает весь кеш переводов
func (s *Service) FlushCache(ctx context.Context) error {
	if s.cache == nil {
		return nil
	}

	return s.cache.Flush(ctx)
}
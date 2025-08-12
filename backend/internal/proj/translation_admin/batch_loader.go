package translation_admin

import (
	"context"
	"fmt"
	"sync"
	"time"

	"backend/internal/domain/models"
	"backend/internal/proj/translation_admin/cache"

	"github.com/rs/zerolog/log"
)

const (
	defaultCacheExpiration = 24 * time.Hour
)

// BatchLoader загрузчик для batch операций с переводами
type BatchLoader struct {
	repo  TranslationRepository
	cache *cache.RedisTranslationCache
	mu    sync.Mutex
}

// NewBatchLoader создает новый batch loader
func NewBatchLoader(repo TranslationRepository, cache *cache.RedisTranslationCache) *BatchLoader {
	return &BatchLoader{
		repo:  repo,
		cache: cache,
	}
}

// TranslationKey ключ для загрузки перевода
type TranslationKey struct {
	EntityType string
	EntityID   int
	Language   string
	FieldName  string
}

// LoadTranslations загружает переводы для множества сущностей одним запросом
func (b *BatchLoader) LoadTranslations(ctx context.Context, keys []TranslationKey) (map[string]string, error) {
	if len(keys) == 0 {
		return make(map[string]string), nil
	}

	result := make(map[string]string)
	var missingKeys []TranslationKey

	// Сначала проверяем кеш если он доступен
	if b.cache != nil {
		for _, key := range keys {
			if translation, found := b.cache.GetTranslation(ctx, key.EntityType, int64(key.EntityID), key.Language, key.FieldName); found {
				resultKey := b.makeKey(key)
				result[resultKey] = translation
			} else {
				missingKeys = append(missingKeys, key)
			}
		}
	} else {
		missingKeys = keys
	}

	// Если все найдено в кеше, возвращаем
	if len(missingKeys) == 0 {
		return result, nil
	}

	// Группируем ключи по entity_type для оптимизации запросов
	groupedKeys := b.groupKeysByType(missingKeys)

	// Загружаем каждую группу одним запросом
	for entityType, typeKeys := range groupedKeys {
		translations, err := b.loadGroupTranslations(ctx, entityType, typeKeys)
		if err != nil {
			log.Error().Err(err).Str("entity_type", entityType).Msg("Failed to load translations group")
			continue
		}

		// Добавляем в результат и кешируем
		for _, trans := range translations {
			key := TranslationKey{
				EntityType: trans.EntityType,
				EntityID:   trans.EntityID,
				Language:   trans.Language,
				FieldName:  trans.FieldName,
			}
			resultKey := b.makeKey(key)
			result[resultKey] = trans.TranslatedText

			// Кешируем если кеш доступен
			if b.cache != nil {
				_ = b.cache.SetTranslation(ctx, trans.EntityType, int64(trans.EntityID), trans.Language, trans.FieldName, trans.TranslatedText, defaultCacheExpiration)
			}
		}
	}

	return result, nil
}

// LoadEntityTranslations загружает все переводы для одной сущности
func (b *BatchLoader) LoadEntityTranslations(ctx context.Context, entityType string, entityID int, languages []string) (map[string]map[string]string, error) {
	// Результат: map[field_name]map[language]translation
	result := make(map[string]map[string]string)

	// Загружаем все переводы для сущности одним запросом
	filters := map[string]interface{}{
		"entity_type": entityType,
		"entity_id":   entityID,
	}

	translations, err := b.repo.GetTranslations(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to load entity translations: %w", err)
	}

	// Группируем по полям и языкам
	for _, trans := range translations {
		if result[trans.FieldName] == nil {
			result[trans.FieldName] = make(map[string]string)
		}
		result[trans.FieldName][trans.Language] = trans.TranslatedText

		// Кешируем если кеш доступен
		if b.cache != nil {
			_ = b.cache.SetTranslation(ctx, trans.EntityType, int64(trans.EntityID), trans.Language, trans.FieldName, trans.TranslatedText, defaultCacheExpiration)
		}
	}

	return result, nil
}

// LoadMultipleEntitiesTranslations загружает переводы для множества сущностей
func (b *BatchLoader) LoadMultipleEntitiesTranslations(ctx context.Context, entityType string, entityIDs []int, language string, fields []string) (map[int]map[string]string, error) {
	// Результат: map[entity_id]map[field_name]translation
	result := make(map[int]map[string]string)

	if len(entityIDs) == 0 {
		return result, nil
	}

	// Инициализируем результат для всех ID
	for _, id := range entityIDs {
		result[id] = make(map[string]string)
	}

	// Создаем batch запрос для всех сущностей
	// В реальной реализации здесь нужен специальный метод в репозитории
	// который поддерживает WHERE entity_id IN (...)
	for _, entityID := range entityIDs {
		filters := map[string]interface{}{
			"entity_type": entityType,
			"entity_id":   entityID,
			"language":    language,
		}

		translations, err := b.repo.GetTranslations(ctx, filters)
		if err != nil {
			log.Error().Err(err).Int("entity_id", entityID).Msg("Failed to load translations for entity")
			continue
		}

		for _, trans := range translations {
			// Фильтруем только нужные поля если указаны
			if len(fields) > 0 && !b.contains(fields, trans.FieldName) {
				continue
			}

			result[trans.EntityID][trans.FieldName] = trans.TranslatedText

			// Кешируем
			if b.cache != nil {
				_ = b.cache.SetTranslation(ctx, trans.EntityType, int64(trans.EntityID), trans.Language, trans.FieldName, trans.TranslatedText, defaultCacheExpiration)
			}
		}
	}

	return result, nil
}

// PreloadTranslationsForListings предзагружает переводы для списка объявлений
func (b *BatchLoader) PreloadTranslationsForListings(ctx context.Context, listings []models.MarketplaceListing, language string) error {
	if len(listings) == 0 {
		return nil
	}

	// Собираем все ID
	listingIDs := make([]int, 0, len(listings))
	for _, listing := range listings {
		listingIDs = append(listingIDs, listing.ID)
	}

	// Загружаем переводы одним запросом
	translations, err := b.LoadMultipleEntitiesTranslations(ctx, "listing", listingIDs, language, []string{"title", "description"})
	if err != nil {
		return fmt.Errorf("failed to preload translations: %w", err)
	}

	// Применяем переводы к объявлениям
	for i := range listings {
		if trans, ok := translations[listings[i].ID]; ok {
			if listings[i].Translations == nil {
				listings[i].Translations = make(models.TranslationMap)
			}

			if title, ok := trans["title"]; ok {
				if listings[i].Translations["title"] == nil {
					listings[i].Translations["title"] = make(map[string]string)
				}
				listings[i].Translations["title"][language] = title
			}

			if description, ok := trans["description"]; ok {
				if listings[i].Translations["description"] == nil {
					listings[i].Translations["description"] = make(map[string]string)
				}
				listings[i].Translations["description"][language] = description
			}
		}
	}

	return nil
}

// PreloadTranslationsForCategories предзагружает переводы для списка категорий
func (b *BatchLoader) PreloadTranslationsForCategories(ctx context.Context, categories []models.MarketplaceCategory, languages []string) error {
	if len(categories) == 0 || len(languages) == 0 {
		return nil
	}

	// Для каждого языка загружаем переводы
	for _, lang := range languages {
		categoryIDs := make([]int, 0, len(categories))
		for _, cat := range categories {
			categoryIDs = append(categoryIDs, cat.ID)
		}

		translations, err := b.LoadMultipleEntitiesTranslations(ctx, "category", categoryIDs, lang, []string{"name", "description"})
		if err != nil {
			log.Error().Err(err).Str("language", lang).Msg("Failed to preload category translations")
			continue
		}

		// Применяем переводы
		for i := range categories {
			if trans, ok := translations[categories[i].ID]; ok {
				if categories[i].Translations == nil {
					categories[i].Translations = make(map[string]string)
				}

				if name, ok := trans["name"]; ok {
					categories[i].Translations[lang+"_name"] = name
				}

				if description, ok := trans["description"]; ok {
					categories[i].Translations[lang+"_description"] = description
				}
			}
		}
	}

	return nil
}

// Вспомогательные методы

func (b *BatchLoader) makeKey(key TranslationKey) string {
	return fmt.Sprintf("%s:%d:%s:%s", key.EntityType, key.EntityID, key.Language, key.FieldName)
}

func (b *BatchLoader) groupKeysByType(keys []TranslationKey) map[string][]TranslationKey {
	grouped := make(map[string][]TranslationKey)
	for _, key := range keys {
		grouped[key.EntityType] = append(grouped[key.EntityType], key)
	}
	return grouped
}

func (b *BatchLoader) loadGroupTranslations(ctx context.Context, entityType string, keys []TranslationKey) ([]models.Translation, error) {
	// Собираем уникальные entity_id
	entityIDMap := make(map[int]bool)
	for _, key := range keys {
		entityIDMap[key.EntityID] = true
	}

	entityIDs := make([]int, 0, len(entityIDMap))
	for id := range entityIDMap {
		entityIDs = append(entityIDs, id)
	}

	// В идеале здесь должен быть batch запрос типа:
	// SELECT * FROM translations WHERE entity_type = ? AND entity_id IN (?, ?, ...)
	// Но пока используем отдельные запросы
	var allTranslations []models.Translation

	for _, entityID := range entityIDs {
		filters := map[string]interface{}{
			"entity_type": entityType,
			"entity_id":   entityID,
		}

		translations, err := b.repo.GetTranslations(ctx, filters)
		if err != nil {
			return nil, err
		}

		allTranslations = append(allTranslations, translations...)
	}

	return allTranslations, nil
}

func (b *BatchLoader) contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

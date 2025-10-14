// backend/internal/proj/c2c/service/marketplace_categories.go
package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"backend/internal/cache"
	"backend/internal/common"
	"backend/internal/domain/models"
)

// getParentCategoryID получает ID родительской категории
func (s *MarketplaceService) getParentCategoryID(ctx context.Context, categoryID int) (int, error) {
	// ВРЕМЕННОЕ РЕШЕНИЕ: хардкод для известных категорий
	// TODO: сделать полноценную функцию через storage
	parentMap := map[int]int{
		10176: 1301, // Градски automobili -> Lični automobili
		1301:  1003, // Lični automobili -> Automobili
	}

	if parentID, exists := parentMap[categoryID]; exists {
		return parentID, nil
	}

	return 0, fmt.Errorf("parent category not found for category %d", categoryID)
}

// GetCategorySuggestions получает предложения категорий по запросу
func (s *MarketplaceService) GetCategorySuggestions(ctx context.Context, query string, size int) ([]models.CategorySuggestion, error) {
	log.Printf("Запрос предложений категорий: '%s'", query)

	// Проверка входных параметров
	if query == "" {
		return []models.CategorySuggestion{}, nil
	}

	// Выполняем SQL-запрос для поиска категорий, связанных с запросом
	sqlQuery := `
        WITH RECURSIVE category_tree AS (
            SELECT c.id, c.name, c.parent_id
            FROM c2c_categories c
            WHERE 1=1
            
            UNION
            
            SELECT c.id, c.name, c.parent_id
            FROM c2c_categories c
            JOIN category_tree t ON c.parent_id = t.id
        ),
        matching_categories AS (
            SELECT 
                c.id,
                c.name,
                (SELECT COUNT(*) FROM c2c_listings ml 
                 WHERE ml.category_id = c.id 
                 AND ml.status = 'active') as listing_count,
                CASE WHEN LOWER(c.name) LIKE LOWER($1) THEN 100 ELSE 0 END +
                (SELECT COUNT(*) FROM c2c_listings ml 
                 WHERE ml.category_id = c.id 
                 AND (LOWER(ml.title) LIKE LOWER($1) OR LOWER(ml.description) LIKE LOWER($1)) 
                 AND ml.status = 'active') as relevance
            FROM c2c_categories c
            WHERE LOWER(c.name) LIKE LOWER($1)
            OR EXISTS (
                SELECT 1 FROM c2c_listings ml 
                WHERE ml.category_id = c.id 
                AND (LOWER(ml.title) LIKE LOWER($1) OR LOWER(ml.description) LIKE LOWER($1))
                AND ml.status = 'active'
            )
        )
        SELECT id, name, listing_count
        FROM matching_categories
        WHERE listing_count > 0
        ORDER BY relevance DESC, listing_count DESC
        LIMIT $2
    `

	rows, err := s.storage.Query(ctx, sqlQuery, "%"+query+"%", size)
	if err != nil {
		log.Printf("Ошибка при выполнении запроса категорий: %v", err)
		return []models.CategorySuggestion{}, nil
	}
	defer func() {
		if err := rows.Close(); err != nil {
			// Логирование ошибки закрытия rows
			_ = err // Explicitly ignore error
		}
	}()

	var results []models.CategorySuggestion
	for rows.Next() {
		var suggestion models.CategorySuggestion

		if err := rows.Scan(&suggestion.ID, &suggestion.Name, &suggestion.ListingCount); err != nil {
			log.Printf("Ошибка сканирования категории: %v", err)
			continue
		}

		results = append(results, suggestion)
	}

	log.Printf("Найдено %d релевантных категорий для запроса '%s'", len(results), query)

	return results, nil
}

// GetCategoryTree получает дерево категорий с кэшированием
func (s *MarketplaceService) GetCategoryTree(ctx context.Context) ([]models.CategoryTreeNode, error) {
	// Если кеш не настроен, работаем напрямую со storage
	if s.cache == nil {
		return s.storage.GetCategoryTree(ctx)
	}

	// Получаем язык из контекста (по умолчанию "en")
	locale := "en"
	if lang, ok := ctx.Value(common.ContextKeyLocale).(string); ok && lang != "" {
		locale = lang
	}

	// Формируем ключ кеша
	cacheKey := cache.BuildCategoryTreeKey(locale, true)

	// Пытаемся получить из кеша
	var result []models.CategoryTreeNode
	err := s.cache.GetOrSet(ctx, cacheKey, &result, 6*time.Hour, func() (interface{}, error) {
		// Загружаем данные из БД
		return s.storage.GetCategoryTree(ctx)
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// RefreshCategoryListingCounts обновляет материализованное представление счетчиков категорий
func (s *MarketplaceService) RefreshCategoryListingCounts(ctx context.Context) error {
	_, err := s.storage.Exec(ctx, "REFRESH MATERIALIZED VIEW CONCURRENTLY category_listing_counts")
	return err
}

// GetCategories получает список активных категорий с кэшированием
func (s *MarketplaceService) GetCategories(ctx context.Context) ([]models.MarketplaceCategory, error) {
	// Если кеш не настроен, работаем напрямую со storage
	if s.cache == nil {
		return s.storage.GetCategories(ctx)
	}

	// Получаем язык из контекста (по умолчанию "en")
	locale := "en"
	if lang, ok := ctx.Value(common.ContextKeyLocale).(string); ok && lang != "" {
		locale = lang
	}

	// Формируем ключ кеша
	cacheKey := cache.BuildCategoriesKey(locale)

	// Пытаемся получить из кеша
	var result []models.MarketplaceCategory
	err := s.cache.GetOrSet(ctx, cacheKey, &result, 6*time.Hour, func() (interface{}, error) {
		// Загружаем данные из БД
		return s.storage.GetCategories(ctx)
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetAllCategories получает все категории (включая неактивные) с кэшированием
func (s *MarketplaceService) GetAllCategories(ctx context.Context) ([]models.MarketplaceCategory, error) {
	// Если кеш не настроен, работаем напрямую со storage
	if s.cache == nil {
		return s.storage.GetAllCategories(ctx)
	}

	// Получаем язык из контекста (по умолчанию "en")
	locale := "en"
	if lang, ok := ctx.Value(common.ContextKeyLocale).(string); ok && lang != "" {
		locale = lang
	}

	// Формируем ключ кеша для всех категорий (включая неактивные)
	cacheKey := cache.BuildCategoryTreeKey(locale, false)

	// Пытаемся получить из кеша
	var result []models.MarketplaceCategory
	err := s.cache.GetOrSet(ctx, cacheKey, &result, 6*time.Hour, func() (interface{}, error) {
		// Загружаем данные из БД
		return s.storage.GetAllCategories(ctx)
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetPopularCategories возвращает самые популярные категории по количеству активных объявлений
func (s *MarketplaceService) GetPopularCategories(ctx context.Context, limit int) ([]models.MarketplaceCategory, error) {
	// Получаем язык из контекста
	locale := "en"
	if lang, ok := ctx.Value(common.ContextKeyLocale).(string); ok && lang != "" {
		locale = lang
	}

	// Получаем популярные категории из хранилища
	categories, err := s.storage.GetPopularCategories(ctx, limit)
	if err != nil {
		return nil, err
	}

	// Переводим названия категорий на нужный язык
	for i := range categories {
		// Проверяем переводы для названия категории
		translations, err := s.storage.GetTranslationsForEntity(ctx, "category", categories[i].ID)
		if err == nil && len(translations) > 0 {
			for _, t := range translations {
				if t.Language == locale && t.FieldName == "name" && t.TranslatedText != "" {
					categories[i].Name = t.TranslatedText
				}
			}
		}
	}

	return categories, nil
}

// getCategorySuggestionsUnified получает унифицированные предложения категорий для автодополнения
func (s *MarketplaceService) getCategorySuggestionsUnified(ctx context.Context, query string, limit int) []models.UnifiedSuggestion {
	// Поиск категорий по:
	// 1. Названию категории и переводам
	// 2. Товарам в категории (если товар соответствует запросу, показываем его категорию)
	// 3. Ключевым словам категории
	sqlQuery := `
		WITH relevant_categories AS (
			-- Категории по прямому совпадению названия
			SELECT DISTINCT mc.id, mc.name, mc.slug, 1 as relevance
			FROM c2c_categories mc
			LEFT JOIN translations t ON t.entity_type = 'MarketplaceCategory' 
									 AND t.entity_id = mc.id 
									 AND t.field_name = 'name'
			WHERE (LOWER(mc.name) LIKE LOWER($1) 
				OR LOWER(t.translated_text) LIKE LOWER($1))
			  AND mc.is_active = true
			
			UNION
			
			-- Категории товаров, которые соответствуют запросу
			SELECT DISTINCT mc.id, mc.name, mc.slug, 2 as relevance
			FROM c2c_categories mc
			INNER JOIN c2c_listings ml ON mc.id = ml.category_id
			WHERE LOWER(ml.title) LIKE LOWER($1)
			  AND ml.status = 'active'
			  AND mc.is_active = true
			
			UNION
			
			-- Категории по ключевым словам
			SELECT DISTINCT mc.id, mc.name, mc.slug, 3 as relevance
			FROM c2c_categories mc
			INNER JOIN category_keywords ck ON mc.id = ck.category_id
			WHERE LOWER(ck.keyword) LIKE LOWER($1)
			  AND mc.is_active = true
		)
		SELECT rc.id, rc.name, rc.slug, 
		       COUNT(DISTINCT ml.id) as listing_count,
		       MIN(rc.relevance) as relevance
		FROM relevant_categories rc
		LEFT JOIN c2c_listings ml ON rc.id = ml.category_id AND ml.status = 'active'
		GROUP BY rc.id, rc.name, rc.slug
		HAVING COUNT(DISTINCT ml.id) > 0  -- Только категории с товарами
		ORDER BY MIN(rc.relevance) ASC, COUNT(DISTINCT ml.id) DESC, LENGTH(rc.name) ASC
		LIMIT $2`

	rows, err := s.storage.Query(ctx, sqlQuery, "%"+query+"%", limit)
	if err != nil {
		log.Printf("Ошибка получения категорий: %v", err)
		return []models.UnifiedSuggestion{}
	}
	defer func() { _ = rows.Close() }()

	var suggestions []models.UnifiedSuggestion
	for rows.Next() {
		var id int
		var name, slug string
		var count int
		var relevance int
		if err := rows.Scan(&id, &name, &slug, &count, &relevance); err != nil {
			log.Printf("Ошибка сканирования категории: %v", err)
			continue
		}

		suggestions = append(suggestions, models.UnifiedSuggestion{
			Type:       "category",
			Value:      name,
			Label:      name,
			CategoryID: &id,
			Count:      &count,
			Metadata: &models.UnifiedSuggestionMeta{
				Category: &name,
			},
		})
	}

	return suggestions
}

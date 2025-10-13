// backend/internal/proj/c2c/storage/postgres/categories.go
package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"backend/internal/common"
	"backend/internal/domain/models"
)

// GetCategoryTree возвращает дерево категорий
func (s *Storage) GetCategoryTree(ctx context.Context) ([]models.CategoryTreeNode, error) {
	log.Printf("GetCategoryTree in storage called")

	// Получаем язык из контекста (по умолчанию "sr")
	locale := "sr"
	if lang, ok := ctx.Value(common.ContextKeyLocale).(string); ok && lang != "" {
		locale = lang
	}
	log.Printf("GetCategoryTree: using locale: %s", locale)

	query := `
WITH RECURSIVE category_tree AS (
    SELECT
        c.id,
        c.name,
        c.slug,
        c.icon,
        c.parent_id,
        to_char(c.created_at, 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"') as created_at,
        ARRAY[c.id] as category_path,
        1 as level,
        COALESCE(clc.listing_count, 0) as listing_count,
        (SELECT COUNT(*) FROM c2c_categories sc WHERE sc.parent_id = c.id) as children_count
    FROM c2c_categories c
    LEFT JOIN category_listing_counts clc ON clc.category_id = c.id
    WHERE c.parent_id IS NULL

    UNION ALL

    SELECT
        c.id,
        c.name,
        c.slug,
        c.icon,
        c.parent_id,
        to_char(c.created_at, 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"') as created_at,
        ct.category_path || c.id,
        ct.level + 1,
        COALESCE(clc.listing_count, 0),
        (SELECT COUNT(*) FROM c2c_categories sc WHERE sc.parent_id = c.id)
    FROM c2c_categories c
    LEFT JOIN category_listing_counts clc ON clc.category_id = c.id
    INNER JOIN category_tree ct ON ct.id = c.parent_id
    WHERE ct.level < 10
),
categories_with_translations AS (
    SELECT
        ct.*,
        COALESCE(
            jsonb_object_agg(
                t.language,
                t.translated_text
            ) FILTER (WHERE t.language IS NOT NULL),
            '{}'::jsonb
        ) as translations
    FROM category_tree ct
    LEFT JOIN translations t ON
        t.entity_type = 'c2c_category'
        AND t.entity_id = ct.id
        AND t.field_name = 'name'
    GROUP BY
        ct.id, ct.name, ct.slug, ct.icon, ct.parent_id,
        ct.created_at, ct.category_path, ct.level, ct.listing_count,
        ct.children_count
)
SELECT
    c1.id,
    c1.name,
    c1.slug,
    c1.icon,
    c1.parent_id,
    c1.created_at,
    c1.level,
    array_to_string(c1.category_path, ',') as path,
    c1.listing_count,
    c1.children_count,
    c1.translations,
    COALESCE(
        json_agg(
            json_build_object(
                'id', c2.id,
                'name', c2.name,
                'slug', c2.slug,
                'icon', c2.icon,
                'parent_id', c2.parent_id,
                'created_at', c2.created_at,
                'level', c2.level,
                'path', array_to_string(c2.category_path, ','),
                'listing_count', c2.listing_count,
                'children_count', c2.children_count,
                'translations', c2.translations
            ) ORDER BY c2.name ASC
        ) FILTER (WHERE c2.id IS NOT NULL),
        '[]'::json
    ) as children
FROM categories_with_translations c1
LEFT JOIN categories_with_translations c2 ON c2.parent_id = c1.id
GROUP BY
    c1.id, c1.name, c1.slug, c1.icon, c1.parent_id,
    c1.created_at, c1.level, c1.category_path, c1.listing_count,
    c1.children_count, c1.translations
ORDER BY c1.name ASC;
`

	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, fmt.Errorf("error querying categories: %w", err)
	}
	defer rows.Close()

	var rootCategories []models.CategoryTreeNode

	for rows.Next() {
		var node models.CategoryTreeNode
		var translationsJson, childrenJson []byte
		var pathStr string
		var icon sql.NullString

		err := rows.Scan(
			&node.ID,
			&node.Name,
			&node.Slug,
			&icon,
			&node.ParentID,
			&node.CreatedAt,
			&node.Level,
			&pathStr,
			&node.ListingCount,
			&node.ChildrenCount,
			&translationsJson,
			&childrenJson,
		)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, fmt.Errorf("error scanning category: %w", err)
		}

		// Обработка NULL icon
		if icon.Valid {
			node.Icon = icon.String
		}

		// Парсим переводы
		if err := json.Unmarshal(translationsJson, &node.Translations); err != nil {
			log.Printf("Error unmarshaling translations for category %d: %v", node.ID, err)
			node.Translations = make(map[string]string)
		}

		// Применяем перевод к названию категории
		if translatedName, ok := node.Translations[locale]; ok && translatedName != "" {
			log.Printf("GetCategoryTree: Applying translation for category %d: %s -> %s (locale: %s)",
				node.ID, node.Name, translatedName, locale)
			node.Name = translatedName
		}

		// Парсим дочерние категории
		var children []models.CategoryTreeNode
		if err := json.Unmarshal(childrenJson, &children); err != nil {
			log.Printf("Error unmarshaling children for category %d: %v", node.ID, err)
			node.Children = make([]models.CategoryTreeNode, 0)
		} else {
			// Применяем переводы к дочерним категориям
			for i := range children {
				if translatedName, ok := children[i].Translations[locale]; ok && translatedName != "" {
					log.Printf("GetCategoryTree: Applying translation for child category %d: %s -> %s (locale: %s)",
						children[i].ID, children[i].Name, translatedName, locale)
					children[i].Name = translatedName
				}
			}
			node.Children = children
		}

		rootCategories = append(rootCategories, node)
	}

	log.Printf("Returning %d root categories with tree", len(rootCategories))
	return rootCategories, nil
}

// GetCategories возвращает список активных категорий
func (s *Storage) GetCategories(ctx context.Context) ([]models.MarketplaceCategory, error) {
	log.Printf("GetCategories: starting to fetch categories")

	// Получаем язык из контекста (по умолчанию "sr")
	locale := "sr"
	if lang, ok := ctx.Value(common.ContextKeyLocale).(string); ok && lang != "" {
		locale = lang
	}
	log.Printf("GetCategories: using locale: %s", locale)

	// Сначала проверим подключение к базе данных
	if err := s.pool.Ping(ctx); err != nil {
		log.Printf("GetCategories: Database ping failed: %v", err)
		return nil, err
	}
	log.Printf("GetCategories: Database ping successful")

	query := `
        WITH category_translations AS (
            SELECT
                t.entity_id,
                jsonb_object_agg(
                    t.language,
                    t.translated_text
                ) as translations
            FROM translations t
            WHERE t.entity_type = 'category'
            AND t.field_name = 'name'
            GROUP BY t.entity_id
        ),
        category_counts AS (
            SELECT
                c.id as category_id,
                COUNT(DISTINCT l.id) + COUNT(DISTINCT sp.id) as total_count
            FROM c2c_categories c
            LEFT JOIN c2c_listings l ON l.category_id = c.id AND l.status = 'active'
            LEFT JOIN b2c_products sp ON sp.category_id = c.id AND sp.is_active = true
            GROUP BY c.id
        )
        SELECT
            c.id, c.name, c.slug, c.parent_id, c.icon, c.description, c.is_active, c.created_at,
            c.seo_title, c.seo_description, c.seo_keywords,
            COALESCE(ct.translations, '{}'::jsonb) as translations,
            COALESCE(cc.total_count, 0) as listing_count
        FROM c2c_categories c
        LEFT JOIN category_translations ct ON c.id = ct.entity_id
        LEFT JOIN category_counts cc ON cc.category_id = c.id
        WHERE c.is_active = true
    `

	log.Printf("GetCategories: Executing query")
	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		log.Printf("GetCategories: Error querying categories: %v", err)
		return nil, err
	}
	defer rows.Close()

	var categories []models.MarketplaceCategory
	for rows.Next() {
		var cat models.MarketplaceCategory
		var translationsJson []byte
		var icon, description, seoTitle, seoDescription, seoKeywords sql.NullString

		var listingCount int
		err := rows.Scan(
			&cat.ID,
			&cat.Name,
			&cat.Slug,
			&cat.ParentID,
			&icon,
			&description,
			&cat.IsActive,
			&cat.CreatedAt,
			&seoTitle,
			&seoDescription,
			&seoKeywords,
			&translationsJson,
			&listingCount,
		)
		if err != nil {
			log.Printf("GetCategories: Error scanning category: %v", err)
			continue
		}

		// Обрабатываем NULL значения
		if icon.Valid {
			cat.Icon = &icon.String
		}
		if description.Valid {
			cat.Description = description.String
		}
		if seoTitle.Valid {
			cat.SEOTitle = seoTitle.String
		}
		if seoDescription.Valid {
			cat.SEODescription = seoDescription.String
		}
		if seoKeywords.Valid {
			cat.SEOKeywords = seoKeywords.String
		}

		// Добавляем количество объявлений
		cat.Count = listingCount

		translations := make(map[string]string)
		if err := json.Unmarshal(translationsJson, &translations); err != nil {
			log.Printf("GetCategories: Error unmarshaling translations for category %d: %v", cat.ID, err)
		} else {
			cat.Translations = translations

			// Применяем перевод к названию категории
			if translatedName, ok := translations[locale]; ok && translatedName != "" {
				cat.Name = translatedName
			}
		}

		categories = append(categories, cat)
	}

	log.Printf("GetCategories: returning %d categories", len(categories))
	return categories, rows.Err()
}

// GetAllCategories returns all categories including inactive ones (for admin panel)
func (s *Storage) GetAllCategories(ctx context.Context) ([]models.MarketplaceCategory, error) {
	log.Printf("GetAllCategories: Starting to fetch all categories (including inactive)")

	// Получаем язык из контекста (по умолчанию "sr")
	locale := "sr"
	if lang, ok := ctx.Value(common.ContextKeyLocale).(string); ok && lang != "" {
		locale = lang
	}
	log.Printf("GetAllCategories: using locale: %s", locale)

	query := `
        WITH category_translations AS (
            SELECT
                c.id AS entity_id,
                jsonb_object_agg(
                    COALESCE(t.language, 'ru'),
                    t.translated_text
                ) AS translations
            FROM c2c_categories c
            LEFT JOIN translations t ON t.entity_id = c.id AND t.entity_type = 'category' AND t.field_name = 'name'
            GROUP BY c.id
        )
        SELECT
            c.id, c.name, c.slug, c.parent_id, c.icon, c.description, c.is_active, c.created_at,
            c.seo_title, c.seo_description, c.seo_keywords,
            COALESCE(ct.translations, '{}'::jsonb) as translations
        FROM c2c_categories c
        LEFT JOIN category_translations ct ON c.id = ct.entity_id
    `

	log.Printf("GetAllCategories: Executing query")
	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		log.Printf("GetAllCategories: Error querying categories: %v", err)
		return nil, err
	}
	defer rows.Close()

	var categories []models.MarketplaceCategory
	for rows.Next() {
		var cat models.MarketplaceCategory
		var parentID sql.NullInt32
		var icon, description, seoTitle, seoDescription, seoKeywords sql.NullString
		var translationsJson []byte

		err := rows.Scan(
			&cat.ID,
			&cat.Name,
			&cat.Slug,
			&parentID,
			&icon,
			&description,
			&cat.IsActive,
			&cat.CreatedAt,
			&seoTitle,
			&seoDescription,
			&seoKeywords,
			&translationsJson,
		)
		if err != nil {
			log.Printf("GetAllCategories: Error scanning row: %v", err)
			return nil, err
		}

		if parentID.Valid {
			pid := int(parentID.Int32)
			cat.ParentID = &pid
		}
		if icon.Valid {
			cat.Icon = &icon.String
		}
		if description.Valid {
			cat.Description = description.String
		}
		if seoTitle.Valid {
			cat.SEOTitle = seoTitle.String
		}
		if seoDescription.Valid {
			cat.SEODescription = seoDescription.String
		}
		if seoKeywords.Valid {
			cat.SEOKeywords = seoKeywords.String
		}

		translations := make(map[string]string)
		if err := json.Unmarshal(translationsJson, &translations); err != nil {
			log.Printf("GetAllCategories: Error unmarshaling translations for category %d: %v", cat.ID, err)
		} else {
			cat.Translations = translations

			// Применяем перевод к названию категории
			if translatedName, ok := translations[locale]; ok && translatedName != "" {
				log.Printf("GetAllCategories: Applying translation for category %d: %s -> %s (locale: %s)",
					cat.ID, cat.Name, translatedName, locale)
				cat.Name = translatedName
			}
		}

		categories = append(categories, cat)
	}

	log.Printf("GetAllCategories: returning %d categories (including inactive)", len(categories))
	return categories, rows.Err()
}

// GetPopularCategories возвращает самые популярные категории по количеству активных объявлений
func (s *Storage) GetPopularCategories(ctx context.Context, limit int) ([]models.MarketplaceCategory, error) {
	log.Printf("GetPopularCategories: fetching top %d categories", limit)

	query := `
		WITH category_counts AS (
			SELECT
				c.id,
				c.name,
				c.slug,
				c.parent_id,
				c.icon,
				c.description,
				c.is_active,
				c.created_at,
				c.seo_title,
				c.seo_description,
				c.seo_keywords,
				COUNT(DISTINCT l.id) as listing_count
			FROM c2c_categories c
			LEFT JOIN c2c_listings l ON l.category_id = c.id AND l.status = 'active'
			WHERE c.is_active = true AND c.parent_id IS NULL
			GROUP BY c.id, c.name, c.slug, c.parent_id, c.icon, c.description,
					 c.is_active, c.created_at, c.seo_title, c.seo_description, c.seo_keywords
		)
		SELECT
			id, name, slug, parent_id, icon, description, is_active, created_at,
			seo_title, seo_description, seo_keywords, listing_count
		FROM category_counts
		ORDER BY listing_count DESC, name ASC
		LIMIT $1
	`

	rows, err := s.pool.Query(ctx, query, limit)
	if err != nil {
		log.Printf("GetPopularCategories: Error querying categories: %v", err)
		return nil, err
	}
	defer rows.Close()

	var categories []models.MarketplaceCategory
	for rows.Next() {
		var cat models.MarketplaceCategory
		var parentID sql.NullInt32
		var icon, description, seoTitle, seoDescription, seoKeywords sql.NullString
		var listingCount int

		err := rows.Scan(
			&cat.ID,
			&cat.Name,
			&cat.Slug,
			&parentID,
			&icon,
			&description,
			&cat.IsActive,
			&cat.CreatedAt,
			&seoTitle,
			&seoDescription,
			&seoKeywords,
			&listingCount,
		)
		if err != nil {
			log.Printf("GetPopularCategories: Error scanning row: %v", err)
			return nil, err
		}

		if parentID.Valid {
			pid := int(parentID.Int32)
			cat.ParentID = &pid
		}
		if icon.Valid {
			cat.Icon = &icon.String
		}
		if description.Valid {
			cat.Description = description.String
		}
		if seoTitle.Valid {
			cat.SEOTitle = seoTitle.String
		}
		if seoDescription.Valid {
			cat.SEODescription = seoDescription.String
		}
		if seoKeywords.Valid {
			cat.SEOKeywords = seoKeywords.String
		}

		// Добавляем количество объявлений в поле Count
		cat.Count = listingCount

		categories = append(categories, cat)
	}

	log.Printf("GetPopularCategories: returning %d popular categories", len(categories))
	return categories, rows.Err()
}

// GetCategoryByID возвращает категорию по ID
func (s *Storage) GetCategoryByID(ctx context.Context, id int) (*models.MarketplaceCategory, error) {
	cat := &models.MarketplaceCategory{}
	var icon, description sql.NullString

	var seoTitle, seoDescription, seoKeywords sql.NullString
	err := s.pool.QueryRow(ctx, `
        SELECT
            id, name, slug, parent_id, icon, description, is_active, created_at,
            seo_title, seo_description, seo_keywords
        FROM c2c_categories
        WHERE id = $1
    `, id).Scan(
		&cat.ID,
		&cat.Name,
		&cat.Slug,
		&cat.ParentID,
		&icon,
		&description,
		&cat.IsActive,
		&cat.CreatedAt,
		&seoTitle,
		&seoDescription,
		&seoKeywords,
	)
	if err != nil {
		return nil, err
	}

	// Обрабатываем NULL значения
	if icon.Valid {
		cat.Icon = &icon.String
	}
	if description.Valid {
		cat.Description = description.String
	}
	if seoTitle.Valid {
		cat.SEOTitle = seoTitle.String
	}
	if seoDescription.Valid {
		cat.SEODescription = seoDescription.String
	}
	if seoKeywords.Valid {
		cat.SEOKeywords = seoKeywords.String
	}

	return cat, nil
}

// SearchCategories ищет категории по названию
func (s *Storage) SearchCategories(ctx context.Context, query string, limit int) ([]models.MarketplaceCategory, error) {
	searchPattern := "%" + strings.ToLower(strings.TrimSpace(query)) + "%"

	// Получаем язык из контекста (по умолчанию "sr")
	locale := "sr"
	if lang, ok := ctx.Value(common.ContextKeyLocale).(string); ok && lang != "" {
		locale = lang
	}
	log.Printf("SearchCategories: using locale: %s", locale)

	sqlQuery := `
		WITH category_counts AS (
			SELECT
				category_id,
				COUNT(*) as listing_count
			FROM c2c_listings
			WHERE status = 'active'
			GROUP BY category_id
		),
		category_translations AS (
			SELECT
				entity_id,
				jsonb_object_agg(
					language,
					translated_text
				) as translations
			FROM translations
			WHERE entity_type = 'c2c_category'
			AND field_name = 'name'
			GROUP BY entity_id
		)
		SELECT
			c.id,
			c.name,
			c.slug,
			c.parent_id,
			c.icon,
			c.created_at,
			COALESCE(ct.translations, '{}'::jsonb) as translations,
			COALESCE(cc.listing_count, 0) as listing_count
		FROM c2c_categories c
		LEFT JOIN category_counts cc ON c.id = cc.category_id
		LEFT JOIN category_translations ct ON c.id = ct.entity_id
		WHERE LOWER(c.name) LIKE $1
			OR EXISTS (
				SELECT 1
				FROM translations t
				WHERE t.entity_type = 'c2c_category'
				AND t.entity_id = c.id
				AND t.field_name = 'name'
				AND LOWER(t.translated_text) LIKE $1
			)
		ORDER BY
			CASE WHEN LOWER(c.name) = LOWER($2) THEN 0 ELSE 1 END,
			cc.listing_count DESC NULLS LAST,
			c.name
		LIMIT $3
	`

	rows, err := s.pool.Query(ctx, sqlQuery, searchPattern, strings.TrimSpace(query), limit)
	if err != nil {
		return nil, fmt.Errorf("error searching categories: %w", err)
	}
	defer rows.Close()

	var categories []models.MarketplaceCategory
	for rows.Next() {
		var cat models.MarketplaceCategory
		var icon sql.NullString
		var translationsJson []byte
		var listingCount int

		err := rows.Scan(
			&cat.ID,
			&cat.Name,
			&cat.Slug,
			&cat.ParentID,
			&icon,
			&cat.CreatedAt,
			&translationsJson,
			&listingCount,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning category: %w", err)
		}

		// Обработка NULL значения для icon
		if icon.Valid {
			cat.Icon = &icon.String
		}

		// Парсим переводы
		translations := make(map[string]string)
		if err := json.Unmarshal(translationsJson, &translations); err != nil {
			log.Printf("Error unmarshaling translations for category %d: %v", cat.ID, err)
		} else {
			cat.Translations = translations

			// Применяем перевод к названию категории
			if translatedName, ok := translations[locale]; ok && translatedName != "" {
				log.Printf("SearchCategories: Applying translation for category %d: %s -> %s (locale: %s)",
					cat.ID, cat.Name, translatedName, locale)
				cat.Name = translatedName
			}
		}

		// Добавляем количество объявлений
		cat.ListingCount = listingCount

		categories = append(categories, cat)
	}

	return categories, nil
}

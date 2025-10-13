// backend/internal/proj/c2c/storage/postgres/listings_crud.go
package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"backend/internal/domain/models"

	"github.com/jackc/pgx/v5"
)

// CreateListing создает новое объявление с переводами и атрибутами
func (s *Storage) CreateListing(ctx context.Context, listing *models.MarketplaceListing) (int, error) {
	var listingID int

	// Если не указан язык, берем значение из контекста или используем по умолчанию
	if listing.OriginalLanguage == "" {
		// Пытаемся получить язык из контекста
		if userLang, ok := ctx.Value("language").(string); ok && userLang != "" {
			listing.OriginalLanguage = userLang
			log.Printf("Using language from context: %s", userLang)
		} else if userLang, ok := ctx.Value("Accept-Language").(string); ok && userLang != "" {
			listing.OriginalLanguage = userLang
			log.Printf("Using language from Accept-Language header: %s", userLang)
		} else {
			// Используем русский по умолчанию, т.к. большинство пользователей русскоговорящие
			listing.OriginalLanguage = "ru"
			log.Printf("Using default language (ru)")
		}
	}

	// Проверяем уникальность slug, если он есть в metadata
	if listing.Metadata != nil {
		if seoData, ok := listing.Metadata["seo"].(map[string]interface{}); ok {
			if slug, ok := seoData["slug"].(string); ok && slug != "" {
				// Генерируем уникальный slug
				uniqueSlug, err := s.GenerateUniqueSlug(ctx, slug, 0)
				if err != nil {
					return 0, fmt.Errorf("error generating unique slug: %w", err)
				}
				// Обновляем slug на уникальный
				seoData["slug"] = uniqueSlug
			}
		}
	}

	// Конвертируем metadata в JSON
	var metadataJSON []byte
	if listing.Metadata != nil {
		var err error
		metadataJSON, err = json.Marshal(listing.Metadata)
		if err != nil {
			return 0, fmt.Errorf("error marshaling metadata: %w", err)
		}
	}

	// Конвертируем мультиязычные адреса в JSON
	var addressMultilingualJSON []byte
	if len(listing.AddressMultilingual) > 0 {
		var err error
		addressMultilingualJSON, err = json.Marshal(listing.AddressMultilingual)
		if err != nil {
			log.Printf("Error marshaling multilingual addresses: %v", err)
			// Не прерываем процесс, продолжаем без мультиязычных адресов
		}
	}

	// Вставляем основные данные объявления
	err := s.pool.QueryRow(ctx, `
        INSERT INTO c2c_listings (
            user_id, category_id, title, description, price,
            condition, status, location, latitude, longitude,
            address_city, address_country, show_on_map, original_language,
            storefront_id, external_id, metadata, address_multilingual
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
        RETURNING id
    `,
		listing.UserID, listing.CategoryID, listing.Title, listing.Description,
		listing.Price, listing.Condition, listing.Status, listing.Location,
		listing.Latitude, listing.Longitude, listing.City, listing.Country,
		listing.ShowOnMap, listing.OriginalLanguage, listing.StorefrontID, listing.ExternalID,
		metadataJSON, addressMultilingualJSON,
	).Scan(&listingID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert listing: %w", err)
	}

	// Сохраняем оригинальный текст как перевод для исходного языка
	_, err = s.pool.Exec(ctx, `
        INSERT INTO translations (
            entity_type, entity_id, language, field_name,
            translated_text, is_machine_translated, is_verified
        ) VALUES ($1, $2, $3, $4, $5, $6, $7)
    `,
		"listing", listingID, listing.OriginalLanguage, "title",
		listing.Title, false, true)
	if err != nil {
		log.Printf("Error saving original title translation: %v", err)
	}

	_, err = s.pool.Exec(ctx, `
        INSERT INTO translations (
            entity_type, entity_id, language, field_name,
            translated_text, is_machine_translated, is_verified
        ) VALUES ($1, $2, $3, $4, $5, $6, $7)
    `,
		"listing", listingID, listing.OriginalLanguage, "description",
		listing.Description, false, true)
	if err != nil {
		log.Printf("Error saving original description translation: %v", err)
	}

	// Сохраняем переводы title и description, если они есть
	if len(listing.Translations) > 0 {
		for lang, fields := range listing.Translations {
			// Пропускаем исходный язык, так как мы уже сохранили его выше
			if lang == listing.OriginalLanguage {
				continue
			}

			// Сохраняем перевод title
			if title, ok := fields["title"]; ok && title != "" {
				_, err = s.pool.Exec(ctx, `
					INSERT INTO translations (
						entity_type, entity_id, language, field_name,
						translated_text, is_machine_translated, is_verified
					) VALUES ($1, $2, $3, $4, $5, $6, $7)
					ON CONFLICT (entity_type, entity_id, language, field_name)
					DO UPDATE SET
						translated_text = EXCLUDED.translated_text,
						is_machine_translated = EXCLUDED.is_machine_translated
				`,
					"listing", listingID, lang, "title",
					title, true, false)
				if err != nil {
					log.Printf("Error saving title translation for lang %s: %v", lang, err)
				}
			}

			// Сохраняем перевод description
			if description, ok := fields["description"]; ok && description != "" {
				_, err = s.pool.Exec(ctx, `
					INSERT INTO translations (
						entity_type, entity_id, language, field_name,
						translated_text, is_machine_translated, is_verified
					) VALUES ($1, $2, $3, $4, $5, $6, $7)
					ON CONFLICT (entity_type, entity_id, language, field_name)
					DO UPDATE SET
						translated_text = EXCLUDED.translated_text,
						is_machine_translated = EXCLUDED.is_machine_translated
				`,
					"listing", listingID, lang, "description",
					description, true, false)
				if err != nil {
					log.Printf("Error saving description translation for lang %s: %v", lang, err)
				}
			}
		}
	}

	if len(listing.Attributes) > 0 {
		// Устанавливаем ID объявления для каждого атрибута
		for i := range listing.Attributes {
			listing.Attributes[i].ListingID = listingID
		}

		if err := s.SaveListingAttributes(ctx, listingID, listing.Attributes); err != nil {
			log.Printf("Error saving attributes for listing %d: %v", listingID, err)
			// Не прерываем создание объявления из-за ошибки с атрибутами
		}
	}

	// Создаем запись в unified_geo для поддержки GIS поиска
	if listing.Latitude != nil && listing.Longitude != nil {
		// Определяем privacy_level, по умолчанию 'exact'
		privacyLevel := "exact"
		if listing.LocationPrivacy != "" {
			privacyLevel = listing.LocationPrivacy
		}

		_, err = s.pool.Exec(ctx, `
			INSERT INTO unified_geo (
				source_type, source_id, location, geohash,
				formatted_address, privacy_level
			) VALUES (
				'c2c_listing', $1,
				ST_SetSRID(ST_MakePoint($2, $3), 4326)::geography,
				substring(ST_GeoHash(ST_SetSRID(ST_MakePoint($2, $3), 4326)) from 1 for 12),
				$4, $5::location_privacy_level
			)
			ON CONFLICT (source_type, source_id)
			DO UPDATE SET
				location = EXCLUDED.location,
				geohash = EXCLUDED.geohash,
				formatted_address = EXCLUDED.formatted_address,
				privacy_level = EXCLUDED.privacy_level,
				updated_at = CURRENT_TIMESTAMP
		`, listingID, listing.Longitude, listing.Latitude, listing.Location, privacyLevel)

		if err != nil {
			log.Printf("Error creating unified_geo entry for listing %d: %v", listingID, err)
			// Не прерываем создание объявления из-за ошибки с geo
		} else {
			// Обновляем materialized view после успешного создания geo записи
			_, err = s.pool.Exec(ctx, "SELECT refresh_map_items_cache()")
			if err != nil {
				log.Printf("Error refreshing map_items_cache: %v", err)
			}
		}
	}

	// Сохраняем варианты товара, если они есть
	if len(listing.Variants) > 0 {
		err = s.CreateListingVariants(ctx, listingID, listing.Variants)
		if err != nil {
			log.Printf("Error saving variants for listing %d: %v", listingID, err)
			// Не прерываем создание объявления из-за ошибки с вариантами
		}
	}

	return listingID, nil
}

// GetListings возвращает список объявлений с фильтрацией и пагинацией
func (s *Storage) GetListings(ctx context.Context, filters map[string]string, limit int, offset int) ([]models.MarketplaceListing, int64, error) {
	userID, _ := ctx.Value("user_id").(int)
	if userID == 0 {
		userID = -1
	}

	// Проверяем существование столбца storefront_id
	hasStorefrontID, err := s.checkStorefrontIDColumn(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("error checking storefront_id column: %w", err)
	}

	// Формируем запрос
	query, args := s.buildListingsQuery(filters, userID, limit, offset, hasStorefrontID)

	// Выполняем запрос
	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying listings: %w", err)
	}
	defer rows.Close()

	// Обрабатываем результаты
	listings, totalCount, err := s.processListingsRows(rows, hasStorefrontID)
	if err != nil {
		return nil, 0, err
	}

	return listings, totalCount, nil
}

// checkStorefrontIDColumn проверяет существование столбца storefront_id
func (s *Storage) checkStorefrontIDColumn(ctx context.Context) (bool, error) {
	var hasStorefrontID bool
	err := s.pool.QueryRow(ctx, `
        SELECT EXISTS (
            SELECT 1
            FROM information_schema.columns
            WHERE table_name = 'c2c_listings' AND column_name = 'storefront_id'
        )
    `).Scan(&hasStorefrontID)
	return hasStorefrontID, err
}

// buildListingsQuery строит SQL запрос для получения листингов
func (s *Storage) buildListingsQuery(filters map[string]string, userID, limit, offset int, hasStorefrontID bool) (string, []interface{}) {
	baseQuery := s.buildListingsBaseQuery(hasStorefrontID)

	args := []interface{}{
		filters["category_id"], // $1
		userID,                 // $2
	}

	conditions, conditionArgs := s.buildListingsConditions(filters, hasStorefrontID, 2)
	args = append(args, conditionArgs...)

	if len(conditions) > 0 {
		baseQuery += " " + strings.Join(conditions, " ")
	}

	// Сортировка
	baseQuery += s.buildListingsSortClause(filters["sort_by"])

	// Пагинация
	argCount := len(args)
	baseQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCount+1, argCount+2)
	args = append(args, limit, offset)

	return baseQuery, args
}

// buildListingsBaseQuery формирует базовую часть SQL запроса
func (s *Storage) buildListingsBaseQuery(hasStorefrontID bool) string {
	baseQuery := `WITH RECURSIVE category_tree AS (
        SELECT c.id, c.parent_id, c.name
        FROM c2c_categories c
        WHERE CASE
            WHEN $1::text = '' OR $1::text IS NULL THEN parent_id IS NULL
            ELSE id = CAST($1 AS INT)
        END
        UNION ALL
        SELECT c.id, c.parent_id, c.name
        FROM c2c_categories c
        INNER JOIN category_tree ct ON c.parent_id = ct.id
    ),
    translations_agg AS (
        SELECT
            entity_id,
            jsonb_object_agg(
                t2.language || '_' || t2.field_name,
                t2.translated_text
            ) as translations
        FROM translations t1
        CROSS JOIN LATERAL (
            SELECT language, field_name, translated_text
            FROM translations t2
            WHERE t2.entity_type = 'listing'
            AND t2.entity_id = t1.entity_id
        ) t2
        WHERE t1.entity_type = 'listing'
        GROUP BY entity_id
    ),
    listing_images AS (
        SELECT
            listing_id,
            COALESCE(
                jsonb_agg(
                    jsonb_build_object(
                        'id', id,
                        'listing_id', listing_id,
                        'file_path', file_path,
                        'file_name', file_name,
                        'file_size', file_size,
                        'content_type', content_type,
                        'is_main', is_main,
                        'created_at', to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS.US"Z"'),
                        'storage_type', storage_type,
                        'storage_bucket', storage_bucket,
                        'public_url', public_url
                    ) ORDER BY is_main DESC, id ASC
                ), '[]'::jsonb
            ) as images
        FROM c2c_images
        GROUP BY listing_id
    )
    SELECT
        l.id, l.user_id, l.category_id, l.title, l.description, l.price,
        l.condition, l.status, l.location, l.latitude, l.longitude,
        l.address_city as city, l.address_country as country, l.views_count,
        l.created_at, l.updated_at, l.show_on_map, l.original_language,`

	if hasStorefrontID {
		baseQuery += ` l.storefront_id,`
	} else {
		baseQuery += ` NULL as storefront_id,`
	}

	baseQuery += `
        l.metadata,
        c.name as category_name, c.slug as category_slug,
        COALESCE(t.translations, '{}'::jsonb) as translations,
        COALESCE(li.images, '[]'::jsonb) as images,
        EXISTS (
            SELECT 1 FROM c2c_favorites mf
            WHERE mf.listing_id = l.id AND mf.user_id = $2
        ) as is_favorite,
        COUNT(*) OVER() as total_count
    FROM c2c_listings l
    JOIN c2c_categories c ON l.category_id = c.id
    LEFT JOIN translations_agg t ON t.entity_id = l.id
    LEFT JOIN listing_images li ON li.listing_id = l.id
    WHERE 1=1
        AND CASE
            WHEN $1::text = '' OR $1::text IS NULL THEN true
            ELSE l.category_id IN (SELECT id FROM category_tree)
        END`

	return baseQuery
}

// buildListingsConditions формирует условия WHERE для фильтрации
func (s *Storage) buildListingsConditions(filters map[string]string, hasStorefrontID bool, startArgNum int) ([]string, []interface{}) {
	conditions := []string{}
	args := []interface{}{}
	argCount := startArgNum

	// Поиск по тексту
	if v, ok := filters["query"]; ok && v != "" {
		argCount++
		conditions = append(conditions, fmt.Sprintf(`
        AND (
            LOWER(l.title) LIKE LOWER($%d)
            OR LOWER(l.description) LIKE LOWER($%d)
            OR EXISTS (
                SELECT 1 FROM translations t
                WHERE t.entity_type = 'listing' AND t.entity_id = l.id
                AND t.field_name IN ('title', 'description')
                AND LOWER(t.translated_text) LIKE LOWER($%d)
            )
        )`, argCount, argCount, argCount))
		args = append(args, "%"+v+"%")
	}

	// Фильтр по цене
	if v, ok := filters["min_price"]; ok && v != "" {
		argCount++
		conditions = append(conditions, fmt.Sprintf("AND l.price >= $%d", argCount))
		args = append(args, v)
	}
	if v, ok := filters["max_price"]; ok && v != "" {
		argCount++
		conditions = append(conditions, fmt.Sprintf("AND l.price <= $%d", argCount))
		args = append(args, v)
	}

	// Фильтр по состоянию
	if v, ok := filters["condition"]; ok && v != "" {
		argCount++
		conditions = append(conditions, fmt.Sprintf("AND l.condition = $%d", argCount))
		args = append(args, v)
	}

	// Фильтр по user_id
	if v, ok := filters["user_id"]; ok && v != "" {
		argCount++
		conditions = append(conditions, fmt.Sprintf("AND l.user_id = $%d", argCount))
		args = append(args, v)
	}

	// Фильтры для storefront_id (только если колонка существует)
	if hasStorefrontID {
		if v, ok := filters["storefront_id"]; ok && v != "" {
			argCount++
			conditions = append(conditions, fmt.Sprintf("AND l.storefront_id = $%d", argCount))
			args = append(args, v)
		}
		if v, ok := filters["exclude_b2c_stores"]; ok && v == "true" {
			conditions = append(conditions, "AND l.storefront_id IS NULL")
		}
	}

	return conditions, args
}

// buildListingsSortClause формирует ORDER BY клаузу
func (s *Storage) buildListingsSortClause(sortBy string) string {
	switch sortBy {
	case "price_asc":
		return " ORDER BY l.price ASC"
	case "price_desc":
		return " ORDER BY l.price DESC"
	default:
		return " ORDER BY l.created_at DESC"
	}
}

// processListingsRows обрабатывает результаты запроса
func (s *Storage) processListingsRows(rows pgx.Rows, hasStorefrontID bool) ([]models.MarketplaceListing, int64, error) {
	var listings []models.MarketplaceListing
	var totalCount int64

	for rows.Next() {
		listing, count, err := s.scanListingRow(rows, hasStorefrontID)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning listing: %w", err)
		}
		totalCount = count
		listings = append(listings, listing)
	}

	return listings, totalCount, nil
}

// scanListingRow сканирует одну строку результата
func (s *Storage) scanListingRow(rows pgx.Rows, hasStorefrontID bool) (models.MarketplaceListing, int64, error) {
	var listing models.MarketplaceListing
	var translationsJSON, imagesJSON, metadataJSON []byte
	var totalCount int64

	listing.User = &models.User{}
	listing.Category = &models.MarketplaceCategory{}

	var (
		tempLocation     sql.NullString
		tempLatitude     sql.NullFloat64
		tempLongitude    sql.NullFloat64
		tempCity         sql.NullString
		tempCountry      sql.NullString
		tempCategoryName sql.NullString
		tempCategorySlug sql.NullString
		tempStorefrontID sql.NullInt32
		tempStatus       sql.NullString
		tempCondition    sql.NullString
		tempDescription  sql.NullString
	)

	err := rows.Scan(
		&listing.ID, &listing.UserID, &listing.CategoryID, &listing.Title,
		&tempDescription, &listing.Price, &tempCondition, &tempStatus,
		&tempLocation, &tempLatitude, &tempLongitude, &tempCity, &tempCountry,
		&listing.ViewsCount, &listing.CreatedAt, &listing.UpdatedAt,
		&listing.ShowOnMap, &listing.OriginalLanguage, &tempStorefrontID,
		&metadataJSON, &tempCategoryName, &tempCategorySlug,
		&translationsJSON, &imagesJSON, &listing.IsFavorite, &totalCount,
	)
	if err != nil {
		return listing, 0, err
	}

	// Обработка метаданных и скидок
	s.processListingMetadata(&listing, metadataJSON)

	// Обработка изображений
	s.processListingImages(&listing, imagesJSON)

	// Обработка NULL значений
	s.processListingNullables(&listing, tempDescription, tempCondition, tempStatus,
		tempLocation, tempLatitude, tempLongitude, tempCity, tempCountry,
		tempCategoryName, tempCategorySlug, tempStorefrontID)

	// Обработка переводов
	if err := json.Unmarshal(translationsJSON, &listing.RawTranslations); err != nil {
		listing.RawTranslations = make(map[string]interface{})
	}
	if listing.RawTranslations != nil {
		listing.Translations = s.processTranslations(listing.RawTranslations)
	}

	// User info будет загружена в handler через auth-service
	listing.User.ID = listing.UserID

	return listing, totalCount, nil
}

// processListingMetadata обрабатывает метаданные и скидки
func (s *Storage) processListingMetadata(listing *models.MarketplaceListing, metadataJSON []byte) {
	if len(metadataJSON) > 0 {
		if err := json.Unmarshal(metadataJSON, &listing.Metadata); err != nil {
			log.Printf("Error unmarshaling metadata for listing %d: %v", listing.ID, err)
		} else if listing.Metadata != nil {
			if discount, ok := listing.Metadata["discount"].(map[string]interface{}); ok {
				listing.HasDiscount = true
				if prevPrice, ok := discount["previous_price"].(float64); ok {
					listing.OldPrice = &prevPrice
				}
				if discountPercent, ok := discount["discount_percent"].(float64); ok {
					percent := int(discountPercent)
					listing.DiscountPercentage = &percent
				}
			}
		}
	}
}

// processListingImages обрабатывает изображения
func (s *Storage) processListingImages(listing *models.MarketplaceListing, imagesJSON []byte) {
	listing.Images = []models.MarketplaceImage{}
	if len(imagesJSON) > 0 {
		var images []models.MarketplaceImage
		if err := json.Unmarshal(imagesJSON, &images); err != nil {
			log.Printf("Error unmarshaling images for listing %d: %v", listing.ID, err)
		} else {
			for i := range images {
				images[i].PublicURL = buildFullImageURL(images[i].PublicURL)
				images[i].ImageURL = buildFullImageURL(images[i].ImageURL)
			}
			listing.Images = images
		}
	}
}

// processListingNullables обрабатывает nullable значения
func (s *Storage) processListingNullables(listing *models.MarketplaceListing,
	description, condition, status, location sql.NullString,
	latitude, longitude sql.NullFloat64,
	city, country, categoryName, categorySlug sql.NullString,
	storefrontID sql.NullInt32) {

	if description.Valid {
		listing.Description = description.String
	}
	if condition.Valid {
		listing.Condition = condition.String
	}
	if status.Valid {
		listing.Status = status.String
	}
	if location.Valid {
		listing.Location = location.String
	}
	if latitude.Valid {
		listing.Latitude = &latitude.Float64
	}
	if longitude.Valid {
		listing.Longitude = &longitude.Float64
	}
	if city.Valid {
		listing.City = city.String
	}
	if country.Valid {
		listing.Country = country.String
	}
	if categoryName.Valid {
		listing.Category.Name = categoryName.String
	}
	if categorySlug.Valid {
		listing.Category.Slug = categorySlug.String
	}
	if storefrontID.Valid {
		sfID := int(storefrontID.Int32)
		listing.StorefrontID = &sfID
	}
}

// deleteListing - приватный метод для удаления листинга (общая логика для user и admin)
func (s *Storage) deleteListing(ctx context.Context, id int, userID int, isAdmin bool) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Проверка владельца для обычных пользователей
	if !isAdmin {
		var listingUserID int
		err = tx.QueryRow(ctx, `SELECT user_id FROM c2c_listings WHERE id = $1`, id).Scan(&listingUserID)
		if err != nil {
			return fmt.Errorf("listing not found: %w", err)
		}
		if listingUserID != userID {
			return fmt.Errorf("you don't have permission to delete this listing")
		}
	} else {
		// Для админа просто проверяем существование
		var exists bool
		err = tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM c2c_listings WHERE id = $1)`, id).Scan(&exists)
		if err != nil || !exists {
			return fmt.Errorf("listing not found")
		}
	}

	// Получаем список изображений для удаления из MinIO
	rows, err := tx.Query(ctx, `SELECT file_path FROM c2c_images WHERE listing_id = $1`, id)
	if err != nil {
		return fmt.Errorf("error getting images: %w", err)
	}
	var imagePaths []string
	for rows.Next() {
		var path string
		if err := rows.Scan(&path); err != nil {
			log.Printf("Error scanning image path: %v", err)
			continue
		}
		imagePaths = append(imagePaths, path)
	}
	rows.Close()

	// Удаляем связанные данные
	if err := s.deleteListingRelatedData(ctx, tx, id); err != nil {
		return err
	}

	// Удаляем само объявление
	result, err := tx.Exec(ctx, `DELETE FROM c2c_listings WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("error deleting listing: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("listing not found")
	}

	// Коммитим транзакцию
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	// Логируем информацию об изображениях
	if len(imagePaths) > 0 {
		prefix := ""
		if isAdmin {
			prefix = "Admin: "
		}
		log.Printf("%sImages marked for deletion from MinIO for listing %d: %v", prefix, id, imagePaths)
	}

	log.Printf("Successfully deleted listing %d with all related data", id)
	return nil
}

// deleteListingRelatedData удаляет все связанные данные листинга
func (s *Storage) deleteListingRelatedData(ctx context.Context, tx pgx.Tx, id int) error {
	// Список таблиц и их условий для удаления
	deletions := []struct {
		table     string
		where     string
		critical  bool // критичная ли ошибка
	}{
		{"c2c_favorites", "listing_id = $1", true},
		{"translations", "entity_type = 'listing' AND entity_id = $1", true},
		{"reviews", "entity_type = 'c2c_listing' AND entity_id = $1", true},
		{"listing_attribute_values", "listing_id = $1", false},
		{"listings_geo", "listing_id = $1", false},
		{"c2c_messages", "listing_id = $1", false},
		{"c2c_chats", "listing_id = $1", false},
		{"c2c_images", "listing_id = $1", true},
	}

	for _, del := range deletions {
		_, err := tx.Exec(ctx, fmt.Sprintf("DELETE FROM %s WHERE %s", del.table, del.where), id)
		if err != nil {
			if del.critical {
				return fmt.Errorf("error removing from %s: %w", del.table, err)
			}
			log.Printf("Error removing from %s: %v", del.table, err)
		}
	}

	return nil
}

// DeleteListing удаляет объявление (для владельца)
func (s *Storage) DeleteListing(ctx context.Context, id int, userID int) error {
	return s.deleteListing(ctx, id, userID, false)
}

// DeleteListingAdmin удаляет объявление без проверки владельца (для администраторов)
func (s *Storage) DeleteListingAdmin(ctx context.Context, id int) error {
	return s.deleteListing(ctx, id, 0, true)
}

// UpdateListing обновляет существующее объявление
func (s *Storage) UpdateListing(ctx context.Context, listing *models.MarketplaceListing) error {
	// Проверяем, не равен ли category_id нулю
	if listing.CategoryID == 0 {
		var currentCategoryID int
		err := s.pool.QueryRow(ctx, `SELECT category_id FROM c2c_listings WHERE id = $1`, listing.ID).Scan(&currentCategoryID)
		if err != nil {
			log.Printf("Ошибка при получении текущей категории: %v", err)
		} else if currentCategoryID > 0 {
			log.Printf("Заменяем нулевую категорию текущей категорией %d для объявления %d", currentCategoryID, listing.ID)
			listing.CategoryID = currentCategoryID
		}
	}

	// Проверяем уникальность slug, если он есть в metadata
	if listing.Metadata != nil {
		if seoData, ok := listing.Metadata["seo"].(map[string]interface{}); ok {
			if slug, ok := seoData["slug"].(string); ok && slug != "" {
				uniqueSlug, err := s.GenerateUniqueSlug(ctx, slug, listing.ID)
				if err != nil {
					return fmt.Errorf("error generating unique slug: %w", err)
				}
				seoData["slug"] = uniqueSlug
			}
		}
	}

	// Конвертируем metadata в JSON
	var metadataJSON []byte
	if listing.Metadata != nil {
		var err error
		metadataJSON, err = json.Marshal(listing.Metadata)
		if err != nil {
			return fmt.Errorf("error marshaling metadata: %w", err)
		}
	}

	result, err := s.pool.Exec(ctx, `
        UPDATE c2c_listings
        SET title = $1, description = $2, price = $3, condition = $4, status = $5,
            location = $6, latitude = $7, longitude = $8, address_city = $9,
            address_country = $10, show_on_map = $11, category_id = $12,
            original_language = $13, metadata = $14, updated_at = CURRENT_TIMESTAMP
        WHERE id = $15 AND user_id = $16
    `,
		listing.Title, listing.Description, listing.Price, listing.Condition, listing.Status,
		listing.Location, listing.Latitude, listing.Longitude, listing.City, listing.Country,
		listing.ShowOnMap, listing.CategoryID, listing.OriginalLanguage, metadataJSON,
		listing.ID, listing.UserID,
	)
	if err != nil {
		return fmt.Errorf("error updating listing: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("listing not found or you don't have permission to update it")
	}

	// Обновляем атрибуты, если они переданы
	if listing.Attributes != nil {
		if err := s.SaveListingAttributes(ctx, listing.ID, listing.Attributes); err != nil {
			log.Printf("Error updating attributes for listing %d: %v", listing.ID, err)
		}
	} else {
		log.Printf("No attributes provided in update for listing %d, existing attributes preserved", listing.ID)
	}

	return nil
}

// GetListingByID получает полную информацию об объявлении по ID
func (s *Storage) GetListingByID(ctx context.Context, id int) (*models.MarketplaceListing, error) {
	listing := &models.MarketplaceListing{
		User:     &models.User{},
		Category: &models.MarketplaceCategory{},
	}

	// Загружаем базовую информацию
	if err := s.fetchListingBase(ctx, listing, id); err != nil {
		// Пробуем найти в storefront_products
		if errors.Is(err, pgx.ErrNoRows) {
			return s.getStorefrontProductAsListing(ctx, id)
		}
		return nil, fmt.Errorf("error getting listing: %w", err)
	}

	// Загружаем дополнительные данные
	s.enrichListingData(ctx, listing)

	return listing, nil
}

// fetchListingBase загружает базовую информацию об объявлении
func (s *Storage) fetchListingBase(ctx context.Context, listing *models.MarketplaceListing, id int) error {
	var (
		description         sql.NullString
		condition           sql.NullString
		status              sql.NullString
		location            sql.NullString
		latitude            sql.NullFloat64
		longitude           sql.NullFloat64
		city                sql.NullString
		country             sql.NullString
		originalLang        sql.NullString
		categoryName        sql.NullString
		categorySlug        sql.NullString
		storefrontID        sql.NullInt32
		locationPrivacy     sql.NullString
		addressMultilingual []byte
	)

	err := s.pool.QueryRow(ctx, `
        SELECT l.id, l.user_id, l.category_id, l.title, l.description, l.price,
               l.condition, l.status, l.location, l.latitude, l.longitude,
               l.address_city as city, l.address_country as country, l.views_count,
               l.created_at, l.updated_at, l.show_on_map, l.original_language,
               c.name as category_name, c.slug as category_slug, l.metadata, l.storefront_id,
               COALESCE(ug.privacy_level::text, 'exact') as location_privacy, l.address_multilingual
        FROM c2c_listings l
        LEFT JOIN c2c_categories c ON l.category_id = c.id
        LEFT JOIN unified_geo ug ON ug.source_type = 'c2c_listing' AND ug.source_id = l.id
        WHERE l.id = $1
    `, id).Scan(
		&listing.ID, &listing.UserID, &listing.CategoryID, &listing.Title, &description,
		&listing.Price, &condition, &status, &location, &latitude, &longitude,
		&city, &country, &listing.ViewsCount, &listing.CreatedAt, &listing.UpdatedAt,
		&listing.ShowOnMap, &originalLang, &categoryName, &categorySlug, &listing.Metadata,
		&storefrontID, &locationPrivacy, &addressMultilingual,
	)
	if err != nil {
		return err
	}

	// Обработка nullable значений
	s.processListingBaseNullables(listing, description, condition, status, location,
		latitude, longitude, city, country, originalLang, categoryName, categorySlug,
		storefrontID, locationPrivacy, addressMultilingual)

	// Обработка метаданных
	s.processDiscountMetadata(listing)

	return nil
}

// processListingBaseNullables обрабатывает nullable поля базовой информации
func (s *Storage) processListingBaseNullables(listing *models.MarketplaceListing,
	description, condition, status, location sql.NullString,
	latitude, longitude sql.NullFloat64,
	city, country, originalLang, categoryName, categorySlug sql.NullString,
	storefrontID sql.NullInt32, locationPrivacy sql.NullString,
	addressMultilingual []byte) {

	if description.Valid {
		listing.Description = description.String
	}
	if condition.Valid {
		listing.Condition = condition.String
	}
	if status.Valid {
		listing.Status = status.String
	}
	if location.Valid {
		listing.Location = location.String
	}
	if latitude.Valid {
		listing.Latitude = &latitude.Float64
	}
	if longitude.Valid {
		listing.Longitude = &longitude.Float64
	}
	if city.Valid {
		listing.City = city.String
	}
	if country.Valid {
		listing.Country = country.String
	}
	if originalLang.Valid {
		listing.OriginalLanguage = originalLang.String
	}
	if categoryName.Valid {
		listing.Category.Name = categoryName.String
	}
	if categorySlug.Valid {
		listing.Category.Slug = categorySlug.String
	}
	if storefrontID.Valid {
		sfID := int(storefrontID.Int32)
		listing.StorefrontID = &sfID
	}
	if locationPrivacy.Valid {
		listing.LocationPrivacy = locationPrivacy.String
	}

	// Парсим мультиязычные адреса
	if len(addressMultilingual) > 0 {
		var multilingualMap map[string]string
		if err := json.Unmarshal(addressMultilingual, &multilingualMap); err != nil {
			log.Printf("Error unmarshaling multilingual addresses for listing %d: %v", listing.ID, err)
		} else {
			listing.AddressMultilingual = multilingualMap
		}
	}
}

// processDiscountMetadata обрабатывает метаданные о скидках
func (s *Storage) processDiscountMetadata(listing *models.MarketplaceListing) {
	if listing.Metadata != nil {
		if discount, ok := listing.Metadata["discount"].(map[string]interface{}); ok {
			listing.HasDiscount = true
			if prevPrice, ok := discount["previous_price"].(float64); ok {
				listing.OldPrice = &prevPrice
				if prevPrice > listing.Price {
					discountPercent := int((prevPrice - listing.Price) / prevPrice * 100)
					discount["discount_percent"] = discountPercent
					listing.DiscountPercentage = &discountPercent
					log.Printf("Обновлен процент скидки для просмотра объявления %d: %d%%", listing.ID, discountPercent)
				}
			}
		}
	}
}

// enrichListingData загружает дополнительные данные для объявления
func (s *Storage) enrichListingData(ctx context.Context, listing *models.MarketplaceListing) {
	// Загружаем переводы
	s.loadListingTranslations(ctx, listing)

	// Загружаем изображения
	s.loadListingImages(ctx, listing)

	// Загружаем путь категории
	s.loadCategoryPath(ctx, listing)

	// Загружаем атрибуты
	s.loadListingAttributes(ctx, listing)

	// Загружаем варианты
	s.loadListingVariants(ctx, listing)
}

// loadListingTranslations загружает переводы объявления
func (s *Storage) loadListingTranslations(ctx context.Context, listing *models.MarketplaceListing) {
	translations := make(map[string]map[string]string)
	rows, err := s.pool.Query(ctx, `
        SELECT language, field_name, translated_text
        FROM translations
        WHERE entity_type = 'listing' AND entity_id = $1
    `, listing.ID)
	if err != nil {
		log.Printf("Error loading translations for listing %d: %v", listing.ID, err)
		listing.Translations = translations
		return
	}
	defer rows.Close()

	for rows.Next() {
		var lang, field, text string
		if err := rows.Scan(&lang, &field, &text); err != nil {
			log.Printf("Error scanning translation: %v", err)
			continue
		}
		if translations[lang] == nil {
			translations[lang] = make(map[string]string)
		}
		translations[lang][field] = text
	}

	listing.Translations = translations
}

// loadListingImages загружает изображения объявления
func (s *Storage) loadListingImages(ctx context.Context, listing *models.MarketplaceListing) {
	var images []models.MarketplaceImage
	var err error

	if listing.StorefrontID != nil && *listing.StorefrontID > 0 {
		images, err = s.GetB2CProductImages(ctx, listing.ID)
	} else {
		images, err = s.GetListingImages(ctx, fmt.Sprintf("%d", listing.ID))
	}

	if err != nil {
		log.Printf("Error loading images for listing %d: %v", listing.ID, err)
		listing.Images = []models.MarketplaceImage{}
	} else {
		listing.Images = images
	}
}

// loadCategoryPath загружает путь категории
func (s *Storage) loadCategoryPath(ctx context.Context, listing *models.MarketplaceListing) {
	if listing.CategoryID <= 0 {
		return
	}

	lang := "en"
	if ctxLang, ok := ctx.Value("locale").(string); ok && ctxLang != "" {
		lang = ctxLang
	}

	query := `
        WITH RECURSIVE category_path AS (
            SELECT id, name, slug, parent_id, 1 as level
            FROM c2c_categories WHERE id = $1
            UNION ALL
            SELECT c.id, c.name, c.slug, c.parent_id, cp.level + 1
            FROM c2c_categories c
            JOIN category_path cp ON c.id = cp.parent_id
        )
        SELECT cp.id, COALESCE(t.translated_text, cp.name) as name, cp.slug
        FROM category_path cp
        LEFT JOIN translations t ON t.entity_type = 'category'
            AND t.entity_id = cp.id AND t.field_name = 'name' AND t.language = $2
        ORDER BY cp.level DESC
    `
	rows, err := s.pool.Query(ctx, query, listing.CategoryID, lang)
	if err != nil {
		log.Printf("Error loading category path for listing %d: %v", listing.ID, err)
		return
	}
	defer rows.Close()

	var categoryIds []int
	var categoryNames []string
	var categorySlugs []string

	for rows.Next() {
		var id int
		var name, slug string
		if err := rows.Scan(&id, &name, &slug); err == nil {
			categoryIds = append(categoryIds, id)
			categoryNames = append(categoryNames, name)
			categorySlugs = append(categorySlugs, slug)
		}
	}

	listing.CategoryPathIds = categoryIds
	listing.CategoryPathNames = categoryNames
	listing.CategoryPathSlugs = categorySlugs
}

// loadListingAttributes загружает атрибуты объявления
func (s *Storage) loadListingAttributes(ctx context.Context, listing *models.MarketplaceListing) {
	attributes, err := s.GetListingAttributes(ctx, listing.ID)
	if err != nil {
		log.Printf("Error loading attributes for listing %d: %v", listing.ID, err)
		listing.Attributes = []models.ListingAttributeValue{}
		return
	}

	// Обработка и форматирование атрибутов
	for i, attr := range attributes {
		attributes[i] = s.processAttributeDisplayValue(attr)
	}

	listing.Attributes = attributes
}

// processAttributeDisplayValue обрабатывает и форматирует отображаемое значение атрибута
func (s *Storage) processAttributeDisplayValue(attr models.ListingAttributeValue) models.ListingAttributeValue {
	hasValue := false

	if attr.TextValue != nil && *attr.TextValue != "" {
		hasValue = true
		attr.DisplayValue = *attr.TextValue
	}

	if attr.NumericValue != nil {
		hasValue = true
		attr.DisplayValue = s.formatNumericAttribute(attr.AttributeName, *attr.NumericValue)
	}

	if attr.BooleanValue != nil {
		hasValue = true
		if *attr.BooleanValue {
			attr.DisplayValue = "Да"
		} else {
			attr.DisplayValue = "Нет"
		}
	}

	if len(attr.JSONValue) > 0 {
		hasValue = true
		if attr.DisplayValue == "" {
			attr.DisplayValue = string(attr.JSONValue)
		}
	}

	// Восстановление типизированного значения из отображаемого
	if !hasValue && attr.DisplayValue != "" {
		s.restoreTypedValue(&attr)
	}

	return attr
}

// formatNumericAttribute форматирует числовое значение атрибута
func (s *Storage) formatNumericAttribute(attrName string, value float64) string {
	switch attrName {
	case attrNameYear:
		return fmt.Sprintf("%d", int(value))
	case attrNameEngineCapacity:
		return fmt.Sprintf("%.1f л", value)
	case attrNameMileage:
		return fmt.Sprintf("%d км", int(value))
	case attrNamePower:
		return fmt.Sprintf("%d л.с.", int(value))
	default:
		return fmt.Sprintf("%g", value)
	}
}

// restoreTypedValue восстанавливает типизированное значение из отображаемого
func (s *Storage) restoreTypedValue(attr *models.ListingAttributeValue) {
	switch attr.AttributeType {
	case "text", "select":
		attr.TextValue = &attr.DisplayValue
	case "number":
		clean := regexp.MustCompile(`[^\d\.-]`).ReplaceAllString(attr.DisplayValue, "")
		if numVal, err := strconv.ParseFloat(clean, 64); err == nil {
			attr.NumericValue = &numVal
		}
	case "boolean":
		boolVal := strings.ToLower(attr.DisplayValue) == "да" ||
			strings.ToLower(attr.DisplayValue) == "true" ||
			attr.DisplayValue == "1"
		attr.BooleanValue = &boolVal
	}
}

// loadListingVariants загружает варианты товара
func (s *Storage) loadListingVariants(ctx context.Context, listing *models.MarketplaceListing) {
	variants, err := s.GetListingVariants(ctx, listing.ID)
	if err != nil {
		log.Printf("Error loading variants for listing %d: %v", listing.ID, err)
		listing.Variants = []models.MarketplaceListingVariant{}
	} else {
		listing.Variants = variants
	}
}

// GetListingBySlug получает объявление по slug
func (s *Storage) GetListingBySlug(ctx context.Context, slug string) (*models.MarketplaceListing, error) {
	var id int
	err := s.pool.QueryRow(ctx, `
		SELECT id FROM c2c_listings WHERE metadata->'seo'->>'slug' = $1 LIMIT 1
	`, slug).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("listing with slug %s not found", slug)
		}
		return nil, fmt.Errorf("error getting listing by slug: %w", err)
	}

	return s.GetListingByID(ctx, id)
}

// IsSlugUnique проверяет уникальность slug
func (s *Storage) IsSlugUnique(ctx context.Context, slug string, excludeID int) (bool, error) {
	var count int
	err := s.pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM c2c_listings
		WHERE metadata->'seo'->>'slug' = $1 AND id != $2
	`, slug, excludeID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking slug uniqueness: %w", err)
	}

	return count == 0, nil
}

// GenerateUniqueSlug генерирует уникальный slug на основе базового
func (s *Storage) GenerateUniqueSlug(ctx context.Context, baseSlug string, excludeID int) (string, error) {
	isUnique, err := s.IsSlugUnique(ctx, baseSlug, excludeID)
	if err != nil {
		return "", err
	}
	if isUnique {
		return baseSlug, nil
	}

	// Пробуем с числовыми суффиксами
	for i := 2; i <= 99; i++ {
		candidateSlug := fmt.Sprintf("%s-%d", baseSlug, i)
		isUnique, err := s.IsSlugUnique(ctx, candidateSlug, excludeID)
		if err != nil {
			return "", err
		}
		if isUnique {
			return candidateSlug, nil
		}
	}

	// Если все числа заняты, используем короткий хеш
	shortHash := fmt.Sprintf("%x", time.Now().Unix())[:6]
	return fmt.Sprintf("%s-%s", baseSlug, shortHash), nil
}

// getStorefrontProductAsListing получает товар из b2c_products и возвращает как MarketplaceListing
func (s *Storage) getStorefrontProductAsListing(ctx context.Context, id int) (*models.MarketplaceListing, error) {
	listing := &models.MarketplaceListing{
		User:     &models.User{},
		Category: &models.MarketplaceCategory{},
	}

	var categoryName, categorySlug sql.NullString

	err := s.pool.QueryRow(ctx, `
        SELECT sp.id, sp.storefront_id, sf.user_id, sp.category_id, sp.name, sp.description,
               sp.price, 'new' as condition, 'active' as status, '' as location,
               0 as latitude, 0 as longitude, '' as city, '' as country,
               sp.view_count, sp.created_at, sp.updated_at, false as show_on_map, 'sr' as original_language,
               c.name as category_name, c.slug as category_slug, '{}'::jsonb as metadata
        FROM b2c_products sp
        LEFT JOIN b2c_stores sf ON sp.storefront_id = sf.id
        LEFT JOIN c2c_categories c ON sp.category_id = c.id
        WHERE sp.id = $1 AND sp.is_active = true
    `, id).Scan(
		&listing.ID, &listing.StorefrontID, &listing.UserID, &listing.CategoryID, &listing.Title,
		&listing.Description, &listing.Price, &listing.Condition, &listing.Status,
		&listing.Location, &listing.Latitude, &listing.Longitude, &listing.City,
		&listing.Country, &listing.ViewsCount, &listing.CreatedAt, &listing.UpdatedAt,
		&listing.ShowOnMap, &listing.OriginalLanguage,
		&categoryName, &categorySlug, &listing.Metadata,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting storefront product: %w", err)
	}

	if categoryName.Valid {
		listing.Category.Name = categoryName.String
	}
	if categorySlug.Valid {
		listing.Category.Slug = categorySlug.String
	}

	listing.User.ID = listing.UserID

	// Загружаем изображения для B2C продукта
	storefrontImages, err := s.GetB2CProductImages(ctx, listing.ID)
	if err != nil {
		log.Printf("Error loading storefront images for product %d: %v", listing.ID, err)
		listing.Images = []models.MarketplaceImage{}
	} else {
		listing.Images = storefrontImages
	}

	listing.Attributes = []models.ListingAttributeValue{}
	listing.Translations = make(map[string]map[string]string)

	return listing, nil
}

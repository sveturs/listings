// backend/internal/proj/c2c/storage/postgres/marketplace.go
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

// Storage struct, конструктор и utility функции перенесены в:
// - storage.go (Storage struct, конструктор, константы, кэш-поля)
// - storage_utils.go (processTranslations, buildFullImageURL)

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

// AddListingImage, GetListingImages, DeleteListingImage, GetB2CProductImages
// перенесены в listings_images.go

// Функция GetListings с проверкой наличия поля storefront_id
func (s *Storage) GetListings(ctx context.Context, filters map[string]string, limit int, offset int) ([]models.MarketplaceListing, int64, error) {
	userID, _ := ctx.Value("user_id").(int)
	if userID == 0 {
		userID = -1
	}

	// Проверяем существование столбца storefront_id
	var hasStorefrontID bool
	err := s.pool.QueryRow(ctx, `
        SELECT EXISTS (
            SELECT 1
            FROM information_schema.columns
            WHERE table_name = 'c2c_listings' AND column_name = 'storefront_id'
        )
    `).Scan(&hasStorefrontID)
	if err != nil {
		return nil, 0, fmt.Errorf("error checking storefront_id column: %w", err)
	}

	// Формируем базовый запрос в зависимости от наличия столбца storefront_id
	baseQuery := `WITH RECURSIVE category_tree AS (
        -- Базовый случай: корневые категории или конкретная категория
        SELECT c.id, c.parent_id, c.name
        FROM c2c_categories c
        WHERE CASE
            WHEN $1::text = '' OR $1::text IS NULL THEN parent_id IS NULL
            ELSE id = CAST($1 AS INT)
        END

        UNION ALL

        -- Рекурсивный случай: все подкатегории
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
        l.id,
        l.user_id,
        l.category_id,
        l.title,
        l.description,
        l.price,
        l.condition,
        l.status,
        l.location,
        l.latitude,
        l.longitude,
        l.address_city as city,
        l.address_country as country,
        l.views_count,
        l.created_at,
        l.updated_at,
        l.show_on_map,
    l.original_language,`

	// Добавляем поле storefront_id, если оно существует
	if hasStorefrontID {
		baseQuery += `
    l.storefront_id,`
	} else {
		baseQuery += `
    NULL as storefront_id,`
	}

	// Добавляем поле для метаданных
	baseQuery += `
    l.metadata,`

	baseQuery += `
    c.name as category_name,
        c.slug as category_slug,
        COALESCE(t.translations, '{}'::jsonb) as translations,
        COALESCE(li.images, '[]'::jsonb) as images,
        EXISTS (
            SELECT 1
            FROM c2c_favorites mf
            WHERE mf.listing_id = l.id
            AND mf.user_id = $2
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

	args := []interface{}{
		filters["category_id"], // $1
		userID,                 // $2
	}

	conditions := []string{}
	argCount := 2

	// Добавляем остальные фильтры
	if v, ok := filters["query"]; ok && v != "" {
		argCount++
		conditions = append(conditions, fmt.Sprintf(`
        AND (
            LOWER(l.title) LIKE LOWER($%d)
            OR LOWER(l.description) LIKE LOWER($%d)
            OR EXISTS (
                SELECT 1
                FROM translations t
                WHERE t.entity_type = 'listing'
                AND t.entity_id = l.id
                AND t.field_name IN ('title', 'description')
                AND LOWER(t.translated_text) LIKE LOWER($%d)
            )
        )`,
			argCount, argCount, argCount))
		args = append(args, "%"+v+"%")
	}

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

	if v, ok := filters["condition"]; ok && v != "" {
		argCount++
		conditions = append(conditions, fmt.Sprintf("AND l.condition = $%d", argCount))
		args = append(args, v)
	}

	// Добавляем фильтр по storefront_id только если колонка существует
	if hasStorefrontID {
		if v, ok := filters["storefront_id"]; ok && v != "" {
			argCount++
			conditions = append(conditions, fmt.Sprintf("AND l.storefront_id = $%d", argCount))
			args = append(args, v)
		}

		// Исключить товары витрин (для админки P2P листингов)
		if v, ok := filters["exclude_b2c_stores"]; ok && v == "true" {
			conditions = append(conditions, "AND l.storefront_id IS NULL")
		}
	}

	// Добавляем фильтр по user_id
	if v, ok := filters["user_id"]; ok && v != "" {
		argCount++
		conditions = append(conditions, fmt.Sprintf("AND l.user_id = $%d", argCount))
		args = append(args, v)
	}

	if len(conditions) > 0 {
		baseQuery += " " + strings.Join(conditions, " ")
	}

	// Сортировка
	switch filters["sort_by"] {
	case "price_asc":
		baseQuery += " ORDER BY l.price ASC"
	case "price_desc":
		baseQuery += " ORDER BY l.price DESC"
	default:
		baseQuery += " ORDER BY l.created_at DESC"
	}

	// Пагинация
	baseQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCount+1, argCount+2)
	args = append(args, limit, offset)

	// Выполнение запроса
	rows, err := s.pool.Query(ctx, baseQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying listings: %w", err)
	}
	defer rows.Close()

	var listings []models.MarketplaceListing
	var totalCount int64

	for rows.Next() {
		var listing models.MarketplaceListing
		var translationsJSON []byte
		var imagesJSON []byte

		listing.User = &models.User{}
		listing.Category = &models.MarketplaceCategory{}
		var metadataJSON []byte

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
			&listing.ID,
			&listing.UserID,
			&listing.CategoryID,
			&listing.Title,
			&tempDescription,
			&listing.Price,
			&tempCondition,
			&tempStatus,
			&tempLocation,
			&tempLatitude,
			&tempLongitude,
			&tempCity,
			&tempCountry,
			&listing.ViewsCount,
			&listing.CreatedAt,
			&listing.UpdatedAt,
			&listing.ShowOnMap,
			&listing.OriginalLanguage,
			&tempStorefrontID,
			&metadataJSON,
			&tempCategoryName,
			&tempCategorySlug,
			&translationsJSON,
			&imagesJSON,
			&listing.IsFavorite,
			&totalCount,
		)
		// После Scan добавьте обработку метаданных
		if len(metadataJSON) > 0 {
			if err := json.Unmarshal(metadataJSON, &listing.Metadata); err != nil {
				log.Printf("Error unmarshaling metadata for listing %d: %v", listing.ID, err)
			} else if listing.Metadata != nil {
				// Обработка метаданных о скидках
				if discount, ok := listing.Metadata["discount"].(map[string]interface{}); ok {
					listing.HasDiscount = true
					if prevPrice, ok := discount["previous_price"].(float64); ok {
						listing.OldPrice = &prevPrice
					}
					// Добавляем процент скидки если его нет
					if discountPercent, ok := discount["discount_percent"].(float64); ok {
						percent := int(discountPercent)
						listing.DiscountPercentage = &percent
					}
				}
			}
		}
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning listing: %w", err)
		}

		// Парсим изображения из JSON
		listing.Images = []models.MarketplaceImage{}
		if len(imagesJSON) > 0 {
			var images []models.MarketplaceImage
			if err := json.Unmarshal(imagesJSON, &images); err != nil {
				log.Printf("Error unmarshaling images for listing %d: %v", listing.ID, err)
			} else {
				// Преобразуем относительные URL в полные
				for i := range images {
					images[i].PublicURL = buildFullImageURL(images[i].PublicURL)
					images[i].ImageURL = buildFullImageURL(images[i].ImageURL)
				}
				listing.Images = images
			}
		}

		// Обработка NULL значений
		if tempDescription.Valid {
			listing.Description = tempDescription.String
		}
		if tempLocation.Valid {
			listing.Location = tempLocation.String
		}
		if tempCondition.Valid {
			listing.Condition = tempCondition.String
		}
		if tempStatus.Valid {
			listing.Status = tempStatus.String
		}
		// User info будет загружена в handler через auth-service
		listing.User.ID = listing.UserID
		if tempStorefrontID.Valid {
			sfID := int(tempStorefrontID.Int32)
			listing.StorefrontID = &sfID
		}
		if tempLatitude.Valid {
			listing.Latitude = &tempLatitude.Float64
		}
		if tempLongitude.Valid {
			listing.Longitude = &tempLongitude.Float64
		}
		if tempCity.Valid {
			listing.City = tempCity.String
		}
		if tempCountry.Valid {
			listing.Country = tempCountry.String
		}
		// User email and picture будет загружено в handler через auth-service
		if tempCategoryName.Valid {
			listing.Category.Name = tempCategoryName.String
		}
		if tempCategorySlug.Valid {
			listing.Category.Slug = tempCategorySlug.String
		}

		// Обработка переводов
		if err := json.Unmarshal(translationsJSON, &listing.RawTranslations); err != nil {
			listing.RawTranslations = make(map[string]interface{})
		}

		if listing.RawTranslations != nil {
			listing.Translations = s.processTranslations(listing.RawTranslations)
		}

		listings = append(listings, listing)
	}

	return listings, totalCount, nil
}


func (s *Storage) DeleteListing(ctx context.Context, id int, userID int) error {
	// Начинаем транзакцию для атомарного удаления
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Проверяем владельца и получаем изображения для удаления
	var listingUserID int
	err = tx.QueryRow(ctx, `
		SELECT user_id FROM c2c_listings WHERE id = $1
	`, id).Scan(&listingUserID)
	if err != nil {
		return fmt.Errorf("listing not found: %w", err)
	}
	if listingUserID != userID {
		return fmt.Errorf("you don't have permission to delete this listing")
	}

	// Получаем список изображений для удаления из MinIO
	rows, err := tx.Query(ctx, `
		SELECT file_path FROM c2c_images WHERE listing_id = $1
	`, id)
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

	// Удаляем записи из избранного
	_, err = tx.Exec(ctx, `
		DELETE FROM c2c_favorites WHERE listing_id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("error removing from favorites: %w", err)
	}

	// Удаляем переводы
	_, err = tx.Exec(ctx, `
		DELETE FROM translations 
		WHERE entity_type = 'listing' AND entity_id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("error removing translations: %w", err)
	}

	// Удаляем отзывы
	_, err = tx.Exec(ctx, `
		DELETE FROM reviews 
		WHERE entity_type = 'c2c_listing' AND entity_id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("error removing reviews: %w", err)
	}

	// Удаляем геоданные
	_, err = tx.Exec(ctx, `
		DELETE FROM listings_geo WHERE listing_id = $1
	`, id)
	if err != nil {
		log.Printf("Error removing geo data: %v", err)
		// Не прерываем, т.к. таблица может не существовать
	}

	// Удаляем атрибуты объявления
	_, err = tx.Exec(ctx, `
		DELETE FROM listing_attribute_values WHERE listing_id = $1
	`, id)
	if err != nil {
		log.Printf("Error removing attribute values: %v", err)
		// Не прерываем, т.к. таблица может не существовать
	}

	// Удаляем сообщения чата
	_, err = tx.Exec(ctx, `
		DELETE FROM c2c_messages WHERE listing_id = $1
	`, id)
	if err != nil {
		log.Printf("Error removing chat messages: %v", err)
	}

	// Удаляем чаты
	_, err = tx.Exec(ctx, `
		DELETE FROM c2c_chats WHERE listing_id = $1
	`, id)
	if err != nil {
		log.Printf("Error removing chats: %v", err)
	}

	// Удаляем изображения из БД (каскадное удаление должно сработать для c2c_listings)
	_, err = tx.Exec(ctx, `
		DELETE FROM c2c_images WHERE listing_id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("error removing images: %w", err)
	}

	// Удаляем само объявление
	result, err := tx.Exec(ctx, `
		DELETE FROM c2c_listings WHERE id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("error deleting listing: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("listing not found")
	}

	// Коммитим транзакцию
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	// Логируем информацию об изображениях для удаления из MinIO
	// Фактическое удаление изображений из MinIO происходит в сервисном слое
	if len(imagePaths) > 0 {
		log.Printf("Images marked for deletion from MinIO for listing %d: %v", id, imagePaths)
		// Изображения будут удалены в методе сервиса после успешного удаления из БД
	}

	log.Printf("Successfully deleted listing %d with all related data from database", id)
	return nil
}

// DeleteListingAdmin удаляет объявление без проверки владельца (для администраторов)
func (s *Storage) DeleteListingAdmin(ctx context.Context, id int) error {
	// Начинаем транзакцию для атомарного удаления
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Проверяем существование объявления
	var exists bool
	err = tx.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM c2c_listings WHERE id = $1)
	`, id).Scan(&exists)
	if err != nil || !exists {
		return fmt.Errorf("listing not found")
	}

	// Получаем список изображений для удаления из MinIO
	rows, err := tx.Query(ctx, `
		SELECT file_path FROM c2c_images WHERE listing_id = $1
	`, id)
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

	// Удаляем записи из избранного
	_, err = tx.Exec(ctx, `
		DELETE FROM c2c_favorites WHERE listing_id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("error removing from favorites: %w", err)
	}

	// Удаляем переводы
	_, err = tx.Exec(ctx, `
		DELETE FROM translations 
		WHERE entity_type = 'listing' AND entity_id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("error removing translations: %w", err)
	}

	// Удаляем отзывы
	_, err = tx.Exec(ctx, `
		DELETE FROM reviews 
		WHERE entity_type = 'c2c_listing' AND entity_id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("error removing reviews: %w", err)
	}

	// Удаляем геоданные
	_, err = tx.Exec(ctx, `
		DELETE FROM listings_geo WHERE listing_id = $1
	`, id)
	if err != nil {
		log.Printf("Error removing geo data: %v", err)
		// Не прерываем, т.к. таблица может не существовать
	}

	// Удаляем атрибуты объявления
	_, err = tx.Exec(ctx, `
		DELETE FROM listing_attribute_values WHERE listing_id = $1
	`, id)
	if err != nil {
		log.Printf("Error removing attribute values: %v", err)
		// Не прерываем, т.к. таблица может не существовать
	}

	// Удаляем сообщения чата
	_, err = tx.Exec(ctx, `
		DELETE FROM c2c_messages WHERE listing_id = $1
	`, id)
	if err != nil {
		log.Printf("Error removing chat messages: %v", err)
	}

	// Удаляем чаты
	_, err = tx.Exec(ctx, `
		DELETE FROM c2c_chats WHERE listing_id = $1
	`, id)
	if err != nil {
		log.Printf("Error removing chats: %v", err)
	}

	// Удаляем изображения из БД (каскадное удаление должно сработать для c2c_listings)
	_, err = tx.Exec(ctx, `
		DELETE FROM c2c_images WHERE listing_id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("error removing images: %w", err)
	}

	// Удаляем само объявление (без проверки user_id для админа)
	result, err := tx.Exec(ctx, `
		DELETE FROM c2c_listings WHERE id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("error deleting listing: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("listing not found")
	}

	// Коммитим транзакцию
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	// Логируем информацию об изображениях для удаления из MinIO
	// Фактическое удаление изображений из MinIO происходит в сервисном слое
	if len(imagePaths) > 0 {
		log.Printf("Admin: Images marked for deletion from MinIO for listing %d: %v", id, imagePaths)
		// Изображения будут удалены в методе сервиса после успешного удаления из БД
	}

	log.Printf("Admin successfully deleted listing %d with all related data from database", id)
	return nil
}

func (s *Storage) UpdateListing(ctx context.Context, listing *models.MarketplaceListing) error {
	// Проверяем, не равен ли category_id нулю
	if listing.CategoryID == 0 {
		// Если category_id = 0, запрашиваем текущее значение из базы
		var currentCategoryID int
		err := s.pool.QueryRow(ctx, `
            SELECT category_id FROM c2c_listings WHERE id = $1
        `, listing.ID).Scan(&currentCategoryID)

		if err != nil {
			log.Printf("Ошибка при получении текущей категории: %v", err)
		} else if currentCategoryID > 0 {
			// Используем текущую категорию, если она не нулевая
			log.Printf("Заменяем нулевую категорию текущей категорией %d для объявления %d", currentCategoryID, listing.ID)
			listing.CategoryID = currentCategoryID
		}
	}

	// Проверяем уникальность slug, если он есть в metadata
	if listing.Metadata != nil {
		if seoData, ok := listing.Metadata["seo"].(map[string]interface{}); ok {
			if slug, ok := seoData["slug"].(string); ok && slug != "" {
				// Генерируем уникальный slug
				uniqueSlug, err := s.GenerateUniqueSlug(ctx, slug, listing.ID)
				if err != nil {
					return fmt.Errorf("error generating unique slug: %w", err)
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
			return fmt.Errorf("error marshaling metadata: %w", err)
		}
	}

	result, err := s.pool.Exec(ctx, `
        UPDATE c2c_listings
        SET
            title = $1,
            description = $2,
            price = $3,
            condition = $4,
            status = $5,
            location = $6,
            latitude = $7,
            longitude = $8,
            address_city = $9,
            address_country = $10,
            show_on_map = $11,
            category_id = $12,
            original_language = $13,
            metadata = $14,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = $15 AND user_id = $16
    `,
		listing.Title,
		listing.Description,
		listing.Price,
		listing.Condition,
		listing.Status,
		listing.Location,
		listing.Latitude,
		listing.Longitude,
		listing.City,
		listing.Country,
		listing.ShowOnMap,
		listing.CategoryID,
		listing.OriginalLanguage,
		metadataJSON,
		listing.ID,
		listing.UserID,
	)
	if err != nil {
		return fmt.Errorf("error updating listing: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("listing not found or you don't have permission to update it")
	}

	// Проверяем, переданы ли атрибуты в запросе
	if listing.Attributes != nil {
		// Атрибуты переданы, обновляем их
		if err := s.SaveListingAttributes(ctx, listing.ID, listing.Attributes); err != nil {
			log.Printf("Error updating attributes for listing %d: %v", listing.ID, err)
			// Продолжаем выполнение даже при ошибке с атрибутами
		}
	} else {
		// Атрибуты не переданы, логируем информацию
		log.Printf("No attributes provided in update for listing %d, existing attributes preserved", listing.ID)
	}

	return nil
}

// Добавить эту функцию перед SaveListingAttributes

func (s *Storage) GetListingByID(ctx context.Context, id int) (*models.MarketplaceListing, error) {
	listing := &models.MarketplaceListing{
		User:     &models.User{},
		Category: &models.MarketplaceCategory{},
	}

	// Получаем основные данные объявления с original_language
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
        SELECT
            l.id, l.user_id, l.category_id, l.title, l.description,
            l.price, l.condition, l.status, l.location, l.latitude,
            l.longitude, l.address_city as city, l.address_country as country, l.views_count,
            l.created_at, l.updated_at, l.show_on_map, l.original_language,
            c.name as category_name, c.slug as category_slug, l.metadata, l.storefront_id,
            COALESCE(ug.privacy_level::text, 'exact') as location_privacy, l.address_multilingual
        FROM c2c_listings l
        LEFT JOIN c2c_categories c ON l.category_id = c.id
        LEFT JOIN unified_geo ug ON ug.source_type = 'c2c_listing' AND ug.source_id = l.id
        WHERE l.id = $1
    `, id).Scan(
		&listing.ID, &listing.UserID, &listing.CategoryID, &listing.Title,
		&description, &listing.Price, &condition, &status,
		&location, &latitude, &longitude, &city,
		&country, &listing.ViewsCount, &listing.CreatedAt, &listing.UpdatedAt,
		&listing.ShowOnMap, &originalLang,
		&categoryName, &categorySlug, &listing.Metadata, &storefrontID, &locationPrivacy, &addressMultilingual,
	)
	if err != nil {
		// Если не найдено в c2c_listings, попробуем найти в storefront_products
		if err.Error() == "no rows in result set" {
			return s.getStorefrontProductAsListing(ctx, id)
		}
		return nil, fmt.Errorf("error getting listing: %w", err)
	}

	// Обрабатываем nullable значения
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
	// Парсим мультиязычные адреса из JSON
	if len(addressMultilingual) > 0 {
		var multilingualMap map[string]string
		if err := json.Unmarshal(addressMultilingual, &multilingualMap); err != nil {
			log.Printf("Error unmarshaling multilingual addresses for listing %d: %v", listing.ID, err)
		} else {
			listing.AddressMultilingual = multilingualMap
		}
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

	// Обработка метаданных и скидок для согласованного отображения
	if listing.Metadata != nil {
		if discount, ok := listing.Metadata["discount"].(map[string]interface{}); ok {
			listing.HasDiscount = true
			if prevPrice, ok := discount["previous_price"].(float64); ok {
				listing.OldPrice = &prevPrice

				// Пересчитываем актуальный процент скидки, чтобы он всегда соответствовал
				// текущей разнице между старой и новой ценой
				if prevPrice > listing.Price {
					discountPercent := int((prevPrice - listing.Price) / prevPrice * 100)
					discount["discount_percent"] = discountPercent
					listing.DiscountPercentage = &discountPercent

					log.Printf("Обновлен процент скидки для просмотра объявления %d: %d%%",
						id, discountPercent)
				}
			}
		}
	}

	// Загружаем переводы
	translations := make(map[string]map[string]string)
	rows, err := s.pool.Query(ctx, `
    SELECT language, field_name, translated_text
    FROM translations
    WHERE entity_type = 'listing' AND entity_id = $1
`, id)
	if err != nil {
		// Просто логируем ошибку, но продолжаем выполнение
		log.Printf("Error loading translations for listing %d: %v", id, err)
	} else {
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
	}

	// Получаем изображения
	var images []models.MarketplaceImage
	if listing.StorefrontID != nil && *listing.StorefrontID > 0 {
		// Для B2C товаров загружаем изображения из storefront_product_images
		storefrontImages, err := s.GetB2CProductImages(ctx, listing.ID)
		if err != nil {
			log.Printf("Error loading storefront images for listing %d: %v", listing.ID, err)
			// Не прерываем выполнение, просто оставляем пустой массив изображений
			images = []models.MarketplaceImage{}
		} else {
			images = storefrontImages
		}
	} else {
		// Для обычных marketplace товаров загружаем изображения из c2c_images
		marketplaceImages, err := s.GetListingImages(ctx, fmt.Sprintf("%d", listing.ID))
		if err != nil {
			log.Printf("Error loading marketplace images for listing %d: %v", listing.ID, err)
			// Не прерываем выполнение, просто оставляем пустой массив изображений
			images = []models.MarketplaceImage{}
		} else {
			images = marketplaceImages
		}
	}

	// Важно! Присваиваем изображения объявлению
	listing.Images = images
	listing.Translations = translations

	// Достаем информацию о пути категории
	if listing.CategoryID > 0 {
		// Получаем язык из контекста
		lang := "en" // значение по умолчанию
		if ctxLang, ok := ctx.Value("locale").(string); ok && ctxLang != "" {
			lang = ctxLang
		}

		// Запрос с поддержкой переводов категорий
		query := `
        WITH RECURSIVE category_path AS (
            SELECT id, name, slug, parent_id, 1 as level
            FROM c2c_categories
            WHERE id = $1

            UNION ALL

            SELECT c.id, c.name, c.slug, c.parent_id, cp.level + 1
            FROM c2c_categories c
            JOIN category_path cp ON c.id = cp.parent_id
        )
        SELECT
            cp.id,
            COALESCE(t.translated_text, cp.name) as name,
            cp.slug
        FROM category_path cp
        LEFT JOIN translations t ON
            t.entity_type = 'category' AND
            t.entity_id = cp.id AND
            t.field_name = 'name' AND
            t.language = $2
        ORDER BY cp.level DESC
    `
		rows, err := s.pool.Query(ctx, query, listing.CategoryID, lang)
		if err == nil {
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
	}
	attributes, err := s.GetListingAttributes(ctx, listing.ID)
	if err != nil {
		log.Printf("Error loading attributes for listing %d: %v", id, err)
	} else {
		log.Printf("INFO: Loaded %d attributes for listing %d", len(attributes), id)

		// Проверка и обработка значений атрибутов
		for i, attr := range attributes {
			// Проверка на пустые значения и установка отображаемого значения
			hasValue := false

			if attr.TextValue != nil && *attr.TextValue != "" {
				hasValue = true
				attr.DisplayValue = *attr.TextValue
			}

			if attr.NumericValue != nil {
				hasValue = true
				// Форматирование числовых значений
				switch attr.AttributeName {
				case attrNameYear:
					attr.DisplayValue = fmt.Sprintf("%d", int(*attr.NumericValue))
				case "engine_capacity":
					attr.DisplayValue = fmt.Sprintf("%.1f л", *attr.NumericValue)
				case attrNameMileage:
					attr.DisplayValue = fmt.Sprintf("%d км", int(*attr.NumericValue))
				case attrNamePower:
					attr.DisplayValue = fmt.Sprintf("%d л.с.", int(*attr.NumericValue))
				default:
					attr.DisplayValue = fmt.Sprintf("%g", *attr.NumericValue)
				}
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

			// Если нет значения, но есть отображаемое значение
			if !hasValue && attr.DisplayValue != "" {
				// Попытка восстановить типизированное значение из отображаемого
				switch attr.AttributeType {
				case "text", "select":
					strVal := attr.DisplayValue
					attr.TextValue = &strVal
				case "number":
					// Удаляем неожиданные символы (буквы, единицы измерения)
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

			// Обновляем атрибут в массиве
			attributes[i] = attr
		}

		listing.Attributes = attributes
	}

	// Загружаем варианты товара, если они есть
	variants, err := s.GetListingVariants(ctx, listing.ID)
	if err != nil {
		log.Printf("Error loading variants for listing %d: %v", listing.ID, err)
		// Не прерываем загрузку объявления из-за ошибки с вариантами
		listing.Variants = []models.MarketplaceListingVariant{}
	} else {
		listing.Variants = variants
	}

	return listing, nil
}

// GetListingBySlug получает объявление по slug
func (s *Storage) GetListingBySlug(ctx context.Context, slug string) (*models.MarketplaceListing, error) {
	// Сначала получаем ID объявления по slug из metadata
	var id int
	err := s.pool.QueryRow(ctx, `
		SELECT id 
		FROM c2c_listings 
		WHERE metadata->'seo'->>'slug' = $1
		LIMIT 1
	`, slug).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("listing with slug %s not found", slug)
		}
		return nil, fmt.Errorf("error getting listing by slug: %w", err)
	}

	// Используем существующий метод GetListingByID
	return s.GetListingByID(ctx, id)
}

// IsSlugUnique проверяет уникальность slug
func (s *Storage) IsSlugUnique(ctx context.Context, slug string, excludeID int) (bool, error) {
	var count int
	err := s.pool.QueryRow(ctx, `
		SELECT COUNT(*) 
		FROM c2c_listings 
		WHERE metadata->'seo'->>'slug' = $1 AND id != $2
	`, slug, excludeID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking slug uniqueness: %w", err)
	}

	return count == 0, nil
}

// GenerateUniqueSlug генерирует уникальный slug на основе базового
func (s *Storage) GenerateUniqueSlug(ctx context.Context, baseSlug string, excludeID int) (string, error) {
	// Сначала проверяем исходный slug
	isUnique, err := s.IsSlugUnique(ctx, baseSlug, excludeID)
	if err != nil {
		return "", err
	}

	if isUnique {
		return baseSlug, nil
	}

	// Если не уникален, пробуем с числовыми суффиксами
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

	// Если все числа от 2 до 99 заняты, используем короткий хеш
	shortHash := fmt.Sprintf("%x", time.Now().Unix())[:6]
	return fmt.Sprintf("%s-%s", baseSlug, shortHash), nil
}

// getStorefrontProductAsListing получает товар из b2c_products и возвращает как MarketplaceListing
func (s *Storage) getStorefrontProductAsListing(ctx context.Context, id int) (*models.MarketplaceListing, error) {
	listing := &models.MarketplaceListing{
		User:     &models.User{},
		Category: &models.MarketplaceCategory{},
	}

	// Получаем данные товара из b2c_products
	var categoryName, categorySlug sql.NullString

	err := s.pool.QueryRow(ctx, `
        SELECT
            sp.id, sp.storefront_id, sf.user_id, sp.category_id, sp.name, sp.description,
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

	// Обрабатываем nullable значения категории
	if categoryName.Valid {
		listing.Category.Name = categoryName.String
	}
	if categorySlug.Valid {
		listing.Category.Slug = categorySlug.String
	}

	// User info будет загружена в handler через auth-service
	listing.User.ID = listing.UserID

	// Загружаем изображения для B2C продукта
	storefrontImages, err := s.GetB2CProductImages(ctx, listing.ID)
	if err != nil {
		log.Printf("Error loading storefront images for product %d: %v", listing.ID, err)
		listing.Images = []models.MarketplaceImage{}
	} else {
		listing.Images = storefrontImages
	}

	// Для storefront продуктов нет атрибутов пока что
	listing.Attributes = []models.ListingAttributeValue{}
	listing.Translations = make(map[string]map[string]string)

	return listing, nil
}

// SearchCategories ищет категории по названию

// GetB2CProductImages перенесена в listings_images.go

// CreateListingVariants, GetListingVariants, UpdateListingVariant, DeleteListingVariant
// перенесены в listings_variants.go

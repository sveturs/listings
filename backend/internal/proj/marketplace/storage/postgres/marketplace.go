// backend/internal/proj/marketplace/storage/postgres/marketplace.go
package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"backend/internal/common"
	"backend/internal/domain/models"
	"backend/internal/proj/marketplace/service"

	"github.com/jackc/pgx/v5/pgxpool"
	// "time"
	// "github.com/jackc/pgx/v5"
)

const (
	// Attribute names
	attrNameArea           = "area"
	attrNameLandArea       = "land_area"
	attrNameMileage        = "mileage"
	attrNameEngineCapacity = "engine_capacity"
	attrNamePower          = "power"
	attrNameYear           = "year"
)

var (
	attributeCacheMutex sync.RWMutex
	attributeCache      map[int][]models.CategoryAttribute
	attributeCacheTime  map[int]time.Time

	rangesCacheMutex sync.RWMutex
	rangesCache      map[int]map[string]map[string]interface{}
	rangesCacheTime  map[int]time.Time

	cacheTTL = 30 * time.Minute
)

func init() {
	attributeCache = make(map[int][]models.CategoryAttribute)
	attributeCacheTime = make(map[int]time.Time)
	rangesCache = make(map[int]map[string]map[string]interface{})
	rangesCacheTime = make(map[int]time.Time)
}

type Storage struct {
	pool               *pgxpool.Pool
	translationService service.TranslationServiceInterface
}

func NewStorage(pool *pgxpool.Pool, translationService service.TranslationServiceInterface) *Storage {
	return &Storage{
		pool:               pool,
		translationService: translationService,
	}
}

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

	// Вставляем основные данные объявления
	err := s.pool.QueryRow(ctx, `
        INSERT INTO marketplace_listings (
            user_id, category_id, title, description, price,
            condition, status, location, latitude, longitude,
            address_city, address_country, show_on_map, original_language,
            storefront_id, external_id
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
        RETURNING id
    `,
		listing.UserID, listing.CategoryID, listing.Title, listing.Description,
		listing.Price, listing.Condition, listing.Status, listing.Location,
		listing.Latitude, listing.Longitude, listing.City, listing.Country,
		listing.ShowOnMap, listing.OriginalLanguage, listing.StorefrontID, listing.ExternalID,
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
		_, err = s.pool.Exec(ctx, `
			INSERT INTO unified_geo (
				source_type, source_id, location, geohash,
				formatted_address
			) VALUES (
				'marketplace_listing', $1, 
				ST_SetSRID(ST_MakePoint($2, $3), 4326)::geography, 
				substring(ST_GeoHash(ST_SetSRID(ST_MakePoint($2, $3), 4326)) from 1 for 12),
				$4
			)
			ON CONFLICT (source_type, source_id) 
			DO UPDATE SET
				location = EXCLUDED.location,
				geohash = EXCLUDED.geohash,
				formatted_address = EXCLUDED.formatted_address,
				updated_at = CURRENT_TIMESTAMP
		`, listingID, listing.Longitude, listing.Latitude, listing.Location)

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

	return listingID, nil
}

func (s *Storage) AddListingImage(ctx context.Context, image *models.MarketplaceImage) (int, error) {
	var id int
	err := s.pool.QueryRow(ctx, `
        INSERT INTO marketplace_images
        (listing_id, file_path, file_name, file_size, content_type, is_main, storage_type, storage_bucket, public_url, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
        RETURNING id
    `, image.ListingID, image.FilePath, image.FileName, image.FileSize, image.ContentType, image.IsMain,
		image.StorageType, image.StorageBucket, image.PublicURL).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Storage) GetListingImages(ctx context.Context, listingID string) ([]models.MarketplaceImage, error) {
	query := `
        SELECT
            id, listing_id, file_path, file_name, file_size,
            content_type, is_main, created_at,
            storage_type, storage_bucket, public_url  -- Эти поля обязательно должны быть в запросе
        FROM marketplace_images
        WHERE listing_id = $1
        ORDER BY is_main DESC, id ASC
    `

	rows, err := s.pool.Query(ctx, query, listingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []models.MarketplaceImage
	for rows.Next() {
		var image models.MarketplaceImage
		err := rows.Scan(
			&image.ID, &image.ListingID, &image.FilePath, &image.FileName,
			&image.FileSize, &image.ContentType, &image.IsMain, &image.CreatedAt,
			&image.StorageType, &image.StorageBucket, &image.PublicURL,
		)
		if err != nil {
			return nil, err
		}
		images = append(images, image)
	}

	return images, nil
}

func (s *Storage) DeleteListingImage(ctx context.Context, imageID string) (string, error) {
	var filePath string
	err := s.pool.QueryRow(ctx,
		"SELECT file_path FROM marketplace_images WHERE id = $1",
		imageID,
	).Scan(&filePath)
	if err != nil {
		return "", err
	}

	_, err = s.pool.Exec(ctx,
		"DELETE FROM marketplace_images WHERE id = $1",
		imageID,
	)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func (s *Storage) processTranslations(rawTranslations interface{}) models.TranslationMap {
	translations := make(models.TranslationMap)

	if rawMap, ok := rawTranslations.(map[string]interface{}); ok {
		for key, value := range rawMap {
			parts := strings.Split(key, "_")
			if len(parts) != 2 {
				continue
			}

			lang, field := parts[0], parts[1]
			if translations[lang] == nil {
				translations[lang] = make(map[string]string)
			}

			if strValue, ok := value.(string); ok {
				translations[lang][field] = strValue
			}
		}
	}

	return translations
}

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
            WHERE table_name = 'marketplace_listings' AND column_name = 'storefront_id'
        )
    `).Scan(&hasStorefrontID)
	if err != nil {
		return nil, 0, fmt.Errorf("error checking storefront_id column: %w", err)
	}

	// Формируем базовый запрос в зависимости от наличия столбца storefront_id
	baseQuery := `WITH RECURSIVE category_tree AS (
        -- Базовый случай: корневые категории или конкретная категория
        SELECT c.id, c.parent_id, c.name
        FROM marketplace_categories c
        WHERE CASE
            WHEN $1::text = '' OR $1::text IS NULL THEN parent_id IS NULL
            ELSE id = CAST($1 AS INT)
        END

        UNION ALL

        -- Рекурсивный случай: все подкатегории
        SELECT c.id, c.parent_id, c.name
        FROM marketplace_categories c
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
        FROM marketplace_images
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
    u.name as user_name,
        u.email as user_email,
        u.created_at as user_created_at,
        u.picture_url as user_picture_url,
        c.name as category_name,
        c.slug as category_slug,
        COALESCE(t.translations, '{}'::jsonb) as translations,
        COALESCE(li.images, '[]'::jsonb) as images,
        EXISTS (
            SELECT 1
            FROM marketplace_favorites mf
            WHERE mf.listing_id = l.id
            AND mf.user_id = $2
        ) as is_favorite,
        COUNT(*) OVER() as total_count
    FROM marketplace_listings l
    JOIN users u ON l.user_id = u.id
    JOIN marketplace_categories c ON l.category_id = c.id
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
			tempEmail        sql.NullString
			tempPictureURL   sql.NullString
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
			tempUserName     sql.NullString
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
			&tempUserName,
			&tempEmail,
			&listing.User.CreatedAt,
			&tempPictureURL,
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
						listing.OldPrice = prevPrice
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
		if tempUserName.Valid {
			listing.User.Name = tempUserName.String
		}
		if tempStorefrontID.Valid {
			sfID := int(tempStorefrontID.Int32)
			listing.StorefrontID = &sfID
			log.Printf("DEBUG: Listing %d has storefront_id: %d", listing.ID, *listing.StorefrontID)
		} else {
			log.Printf("DEBUG: Listing %d has no storefront_id", listing.ID)
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
		if tempEmail.Valid {
			listing.User.Email = tempEmail.String
		}
		if tempPictureURL.Valid {
			listing.User.PictureURL = tempPictureURL.String
		}
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
        (SELECT COUNT(*) FROM marketplace_categories sc WHERE sc.parent_id = c.id) as children_count
    FROM marketplace_categories c
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
        (SELECT COUNT(*) FROM marketplace_categories sc WHERE sc.parent_id = c.id)
    FROM marketplace_categories c
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
        t.entity_type = 'marketplace_category'
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

		// Добавляем логирование переводов и детей
		//		log.Printf("Raw translations JSON for category %d: %s", node.ID, string(translationsJson))
		//		log.Printf("Raw children JSON for category %d: %s", node.ID, string(childrenJson))

		if err := json.Unmarshal(translationsJson, &node.Translations); err != nil {
			log.Printf("Error unmarshaling translations for category %d: %v", node.ID, err)
			node.Translations = make(map[string]string)
		}

		// Применяем перевод к названию категории если он есть для запрашиваемого языка
		if translatedName, ok := node.Translations[locale]; ok && translatedName != "" {
			log.Printf("GetCategoryTree: Applying translation for category %d: %s -> %s (locale: %s)",
				node.ID, node.Name, translatedName, locale)
			node.Name = translatedName
		}

		var children []models.CategoryTreeNode
		if err := json.Unmarshal(childrenJson, &children); err != nil {
			log.Printf("Error unmarshaling children for category %d: %v", node.ID, err)
			node.Children = make([]models.CategoryTreeNode, 0)
		} else {
			// Применяем переводы к дочерним категориям рекурсивно
			for i := range children {
				if translatedName, ok := children[i].Translations[locale]; ok && translatedName != "" {
					log.Printf("GetCategoryTree: Applying translation for child category %d: %s -> %s (locale: %s)",
						children[i].ID, children[i].Name, translatedName, locale)
					children[i].Name = translatedName
				}
			}
			node.Children = children
			//			log.Printf("Category %d has %d children", node.ID, len(children))
		}

		rootCategories = append(rootCategories, node)
	}

	// Добавляем логирование результата
	//	for _, cat := range rootCategories {
	//		log.Printf("Root category %d (%s) has %d children and children_count=%d",
	//			cat.ID, cat.Name, len(cat.Children), cat.ChildrenCount)
	//	}

	log.Printf("Returning %d root categories with tree", len(rootCategories))
	return rootCategories, nil
}

func (s *Storage) AddToFavorites(ctx context.Context, userID int, listingID int) error {
	_, err := s.pool.Exec(ctx, `
        INSERT INTO marketplace_favorites (user_id, listing_id)
        VALUES ($1, $2)
        ON CONFLICT (user_id, listing_id) DO NOTHING
    `, userID, listingID)
	return err
}

func (s *Storage) RemoveFromFavorites(ctx context.Context, userID int, listingID int) error {
	_, err := s.pool.Exec(ctx, `
        DELETE FROM marketplace_favorites
        WHERE user_id = $1 AND listing_id = $2
    `, userID, listingID)
	return err
}

func (s *Storage) GetUserFavorites(ctx context.Context, userID int) ([]models.MarketplaceListing, error) {
	query := `
        WITH listing_images AS (
            SELECT
                listing_id,
                jsonb_agg(
                    jsonb_build_object(
                        'id', id,
                        'listing_id', listing_id,
                        'file_path', file_path,
                        'file_name', file_name,
                        'file_size', file_size,
                        'content_type', content_type,
                        'is_main', is_main,
                        'storage_type', storage_type,
                        'storage_bucket', storage_bucket,
                        'public_url', public_url,
                        'created_at', created_at
                    ) ORDER BY is_main DESC, id ASC
                ) as images
            FROM marketplace_images
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
            u.name,
            u.email,
			u.created_at as user_created_at,
            COALESCE(u.picture_url, ''),
            c.name as category_name,
            c.slug as category_slug,
            true as is_favorite,
            COALESCE(li.images, '[]'::jsonb) as listing_images
        FROM marketplace_listings l
        JOIN marketplace_favorites f ON l.id = f.listing_id
        LEFT JOIN users u ON l.user_id = u.id
        LEFT JOIN marketplace_categories c ON l.category_id = c.id
        LEFT JOIN listing_images li ON li.listing_id = l.id
        WHERE f.user_id = $1
        ORDER BY f.created_at DESC
    `

	rows, err := s.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying favorites: %w", err)
	}
	defer rows.Close()

	var listings []models.MarketplaceListing
	for rows.Next() {
		listing := models.MarketplaceListing{
			User:     &models.User{},
			Category: &models.MarketplaceCategory{},
		}
		var userPictureURL string
		var imagesJSON json.RawMessage

		err := rows.Scan(
			&listing.ID,
			&listing.UserID,
			&listing.CategoryID,
			&listing.Title,
			&listing.Description,
			&listing.Price,
			&listing.Condition,
			&listing.Status,
			&listing.Location,
			&listing.Latitude,
			&listing.Longitude,
			&listing.City,
			&listing.Country,
			&listing.ViewsCount,
			&listing.CreatedAt,
			&listing.UpdatedAt,
			&listing.User.Name,
			&listing.User.Email,
			&listing.User.CreatedAt,
			&userPictureURL,
			&listing.Category.Name,
			&listing.Category.Slug,
			&listing.IsFavorite,
			&imagesJSON,
		)
		if err != nil {
			log.Printf("Error scanning listing: %v", err)
			continue
		}

		// Присваиваем отдельно
		listing.User.PictureURL = userPictureURL
		listing.User.ID = listing.UserID

		// Парсим изображения из JSON
		var images []models.MarketplaceImage
		if err := json.Unmarshal(imagesJSON, &images); err != nil {
			log.Printf("Error unmarshalling images for listing %d: %v", listing.ID, err)
			images = []models.MarketplaceImage{}
		}
		listing.Images = images

		listings = append(listings, listing)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return listings, nil
}

func (s *Storage) GetFavoritedUsers(ctx context.Context, listingID int) ([]int, error) {
	query := `
        SELECT user_id
        FROM marketplace_favorites
        WHERE listing_id = $1
    `
	rows, err := s.pool.Query(ctx, query, listingID)
	if err != nil {
		return nil, fmt.Errorf("error querying favorited users: %w", err)
	}
	defer rows.Close()

	var userIDs []int
	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			return nil, fmt.Errorf("error scanning user ID: %w", err)
		}
		userIDs = append(userIDs, userID)
	}

	return userIDs, nil
}

func (s *Storage) DeleteListing(ctx context.Context, id int, userID int) error {
	// Сначала удаляем записи из избранного
	_, err := s.pool.Exec(ctx, `
        DELETE FROM marketplace_favorites
        WHERE listing_id = $1
    `, id)
	if err != nil {
		return fmt.Errorf("error removing listing from favorites: %w", err)
	}

	// Удаляем объявление
	result, err := s.pool.Exec(ctx, `
        DELETE FROM marketplace_listings
        WHERE id = $1 AND user_id = $2
    `, id, userID)
	if err != nil {
		return fmt.Errorf("error deleting listing: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("listing not found or you don't have permission to delete it")
	}

	return nil
}

func (s *Storage) UpdateListing(ctx context.Context, listing *models.MarketplaceListing) error {
	// Проверяем, не равен ли category_id нулю
	if listing.CategoryID == 0 {
		// Если category_id = 0, запрашиваем текущее значение из базы
		var currentCategoryID int
		err := s.pool.QueryRow(ctx, `
            SELECT category_id FROM marketplace_listings WHERE id = $1
        `, listing.ID).Scan(&currentCategoryID)

		if err != nil {
			log.Printf("Ошибка при получении текущей категории: %v", err)
		} else if currentCategoryID > 0 {
			// Используем текущую категорию, если она не нулевая
			log.Printf("Заменяем нулевую категорию текущей категорией %d для объявления %d", currentCategoryID, listing.ID)
			listing.CategoryID = currentCategoryID
		}
	}

	result, err := s.pool.Exec(ctx, `
        UPDATE marketplace_listings
        SET
            title = $1,
            description = $2,
            price = $3,
            condition = $4,
            status = $5,
            location = $6,
            latitude = $7,
            longitude = $8,
            city = $9,
            country = $10,
            show_on_map = $11,
            category_id = $12,
            original_language = $13,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = $14 AND user_id = $15
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
func sanitizeAttributeValue(attr *models.ListingAttributeValue) {
	// Ограничение длины текстовых атрибутов
	if attr.TextValue != nil {
		if len(*attr.TextValue) > 1000 {
			truncated := (*attr.TextValue)[:1000]
			attr.TextValue = &truncated
			log.Printf("Attribute value truncated for attribute %s (ID: %d)",
				attr.AttributeName, attr.AttributeID)
		}
	}

	// Проверка на NaN и Inf для числовых атрибутов
	if attr.NumericValue != nil {
		numVal := *attr.NumericValue
		if math.IsNaN(numVal) || math.IsInf(numVal, 0) {
			defaultVal := 0.0
			attr.NumericValue = &defaultVal
			log.Printf("Invalid numeric value (NaN/Inf) replaced with 0 for attribute %s (ID: %d)",
				attr.AttributeName, attr.AttributeID)
		}
	}

	// Стандартизация обработки пустых значений
	if attr.TextValue != nil && *attr.TextValue == "" {
		attr.TextValue = nil // Пустые строки -> NULL
	}

	if attr.NumericValue != nil && *attr.NumericValue == 0 {
		// Для некоторых атрибутов нуль может быть валидным значением
		// Проверяем название атрибута
		if !isZeroValidValue(attr.AttributeName) {
			attr.NumericValue = nil
		}
	}

	// Если все значения NULL, устанавливаем DisplayValue в пустую строку
	if attr.TextValue == nil && attr.NumericValue == nil &&
		attr.BooleanValue == nil && attr.JSONValue == nil {
		attr.DisplayValue = ""
	}
}

// Функция определяет, является ли нулевое значение допустимым для атрибута
func isZeroValidValue(attrName string) bool {
	// Для этих атрибутов ноль - допустимое значение
	zeroValidAttrs := map[string]bool{
		"floor":         true, // Например, цокольный этаж
		attrNameMileage: true, // Для новых автомобилей
		"price":         true, // Для бесплатных объявлений
	}
	return zeroValidAttrs[attrName]
}

// SaveListingAttributes сохраняет значения атрибутов для объявления
func (s *Storage) SaveListingAttributes(ctx context.Context, listingID int, attributes []models.ListingAttributeValue) error {
	log.Printf("Saving %d attributes for listing %d", len(attributes), listingID)

	// Начинаем транзакцию
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			// Игнорируем ошибку если транзакция уже была завершена
			_ = err // Explicitly ignore error
		}
	}()

	// Удаляем старые атрибуты
	_, err = tx.Exec(ctx, `DELETE FROM listing_attribute_values WHERE listing_id = $1`, listingID)
	if err != nil {
		return fmt.Errorf("error deleting old attributes: %w", err)
	}

	// Проверяем, есть ли атрибуты для сохранения
	if len(attributes) == 0 {
		log.Printf("Storage: No attributes to save for listing %d", listingID)
		return tx.Commit(ctx)
	}

	// Карта для отслеживания уникальных attribute_id
	seen := make(map[int]bool)
	valueStrings := make([]string, 0, len(attributes))
	valueArgs := make([]interface{}, 0, len(attributes)*7) // 7 параметров включая unit
	counter := 1

	for i, attr := range attributes {
		// Санитизация значений атрибутов
		sanitizeAttributeValue(&attr)
		attributes[i] = attr

		// Проверка на нулевые или некорректные attribute_id
		if attr.AttributeID <= 0 {
			log.Printf("Storage: Invalid attribute ID: %d, skipping", attr.AttributeID)
			continue
		}

		// Определяем единицу измерения на основе имени атрибута или используем указанную
		var unit string
		if attr.Unit != "" {
			unit = attr.Unit
		} else {
			switch attr.AttributeName {
			case attrNameArea:
				unit = "m²"
			case attrNameLandArea:
				unit = "ar"
			case attrNameMileage:
				unit = "km"
			case attrNameEngineCapacity:
				unit = "l"
			case attrNamePower:
				unit = "ks"
			case "screen_size":
				unit = "inč"
			case "rooms":
				unit = "soba"
			case "floor", "total_floors":
				unit = "sprat"
			}
		}

		// Числовые атрибуты - дополнительная обработка
		if attr.NumericValue == nil && attr.TextValue != nil && *attr.TextValue != "" {
			// Список числовых атрибутов для конвертации
			numericAttrs := map[string]bool{
				"rooms": true, "floor": true, "total_floors": true, attrNameArea: true,
				attrNameLandArea: true, attrNameMileage: true, attrNameYear: true, attrNameEngineCapacity: true,
				attrNamePower: true, "screen_size": true,
			}

			if numericAttrs[attr.AttributeName] {
				// Преобразуем текст в число
				clean := regexp.MustCompile(`[^\d\.-]`).ReplaceAllString(*attr.TextValue, "")
				if numVal, err := strconv.ParseFloat(clean, 64); err == nil {
					attr.NumericValue = &numVal
					log.Printf("Converted text value '%s' to numeric: %f for attribute %s",
						*attr.TextValue, numVal, attr.AttributeName)
				}
			}
		}

		// Проверка на дубликаты по attribute_id
		if seen[attr.AttributeID] {
			log.Printf("Storage: Duplicate attribute ID %d for listing %d, skipping", attr.AttributeID, listingID)
			continue
		}
		seen[attr.AttributeID] = true

		// Проверяем, что есть хотя бы одно значение для сохранения
		hasValue := attr.TextValue != nil || attr.NumericValue != nil ||
			attr.BooleanValue != nil || attr.JSONValue != nil ||
			attr.DisplayValue != ""
		if !hasValue {
			log.Printf("Storage: No value provided for attribute %d, skipping", attr.AttributeID)
			continue
		}

		// Подготавливаем часть запроса для этого атрибута, добавляя unit
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			counter, counter+1, counter+2, counter+3, counter+4, counter+5, counter+6))

		// Добавляем параметры
		valueArgs = append(valueArgs, listingID, attr.AttributeID)

		// Текстовое значение
		if attr.TextValue != nil && *attr.TextValue != "" {
			valueArgs = append(valueArgs, *attr.TextValue)
		} else {
			valueArgs = append(valueArgs, nil)
		}

		// Числовое значение с проверками
		if attr.NumericValue != nil {
			numericVal := *attr.NumericValue
			// Проверка на NaN или Inf
			if math.IsNaN(numericVal) || math.IsInf(numericVal, 0) {
				log.Printf("Storage: Invalid numeric value (NaN/Inf) for attribute %d, using 0", attr.AttributeID)
				numericVal = 0.0
			}
			valueArgs = append(valueArgs, numericVal)
		} else {
			valueArgs = append(valueArgs, nil)
		}

		// Логическое значение
		if attr.BooleanValue != nil {
			valueArgs = append(valueArgs, *attr.BooleanValue)
		} else {
			valueArgs = append(valueArgs, nil)
		}

		// JSON значение
		if len(attr.JSONValue) > 0 {
			valueArgs = append(valueArgs, string(attr.JSONValue))
		} else {
			valueArgs = append(valueArgs, nil)
		}

		// Добавляем единицу измерения
		valueArgs = append(valueArgs, unit)

		counter += 7
	}

	// Если нет атрибутов для вставки, завершаем транзакцию
	if len(valueStrings) == 0 {
		log.Printf("Storage: No valid attributes found for listing %d after filtering", listingID)
		return tx.Commit(ctx)
	}

	// Составляем запрос для множественной вставки
	query := fmt.Sprintf(`
        INSERT INTO listing_attribute_values (
            listing_id, attribute_id, text_value, numeric_value, boolean_value, json_value, unit
        ) VALUES %s
        ON CONFLICT (listing_id, attribute_id) DO UPDATE SET
            text_value = EXCLUDED.text_value,
            numeric_value = EXCLUDED.numeric_value,
            boolean_value = EXCLUDED.boolean_value,
            json_value = EXCLUDED.json_value,
            unit = EXCLUDED.unit
    `, strings.Join(valueStrings, ","))

	// Выполняем запрос
	_, err = tx.Exec(ctx, query, valueArgs...)
	if err != nil {
		log.Printf("Storage: Error executing bulk insert: %v", err)
		log.Printf("Storage: Query: %s", query)
		log.Printf("Storage: Args: %+v", valueArgs)
		return fmt.Errorf("error inserting attribute values: %w", err)
	}

	// Фиксируем транзакцию
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	log.Printf("Storage: Successfully saved %d unique attributes for listing %d", len(valueStrings), listingID)
	return nil
}

// Добавить в файле marketplace.go
func (s *Storage) GetFormattedAttributeValue(ctx context.Context, attr models.ListingAttributeValue, language string) string {
	// Для числовых атрибутов с единицей измерения
	if attr.NumericValue != nil && attr.Unit != "" {
		// Получаем перевод единицы измерения
		var displayFormat string
		err := s.pool.QueryRow(ctx, `
            SELECT display_format FROM unit_translations
            WHERE unit = $1 AND language = $2
        `, attr.Unit, language).Scan(&displayFormat)

		if err == nil && displayFormat != "" {
			// Используем формат для отображения
			return fmt.Sprintf(displayFormat, *attr.NumericValue)
		}

		// Если не нашли перевод, используем стандартный формат
		return fmt.Sprintf("%g %s", *attr.NumericValue, attr.Unit)
	}

	// Для других типов атрибутов возвращаем DisplayValue
	return attr.DisplayValue
}

// GetListingAttributes получает значения атрибутов для объявления
func (s *Storage) GetListingAttributes(ctx context.Context, listingID int) ([]models.ListingAttributeValue, error) {
	query := `
        SELECT DISTINCT ON (a.id)
            v.listing_id,
            v.attribute_id,
            a.name AS attribute_name,
            a.display_name,
            a.attribute_type,
            v.text_value,
            v.numeric_value,
            v.boolean_value,
            v.json_value,
            v.unit,
            -- Настройки отображения из основной таблицы атрибутов или переопределения в маппинге
            COALESCE(cam.is_required, a.is_required) as is_required,
            COALESCE(cam.show_in_card, a.show_in_card) as show_in_card,
            COALESCE(cam.show_in_list, a.show_in_list) as show_in_list,
            COALESCE(
                jsonb_object_agg(
                    t.language,
                    t.translated_text
                ) FILTER (WHERE t.language IS NOT NULL),
                '{}'::jsonb
            ) as translations,
            COALESCE(
                (SELECT jsonb_object_agg(
                    o.language,
                    o.field_translations
                )
                FROM (
                    SELECT
                        language,
                        jsonb_object_agg(field_name, translated_text) as field_translations
                    FROM translations
                    WHERE entity_type = 'attribute_option'
                    AND entity_id = a.id
                    GROUP BY language
                ) o),
                '{}'::jsonb
            ) as option_translations
        FROM listing_attribute_values v
        JOIN category_attributes a ON v.attribute_id = a.id
        JOIN marketplace_listings ml ON ml.id = v.listing_id
        LEFT JOIN category_attribute_mapping cam ON
            cam.category_id = ml.category_id AND
            cam.attribute_id = a.id
        LEFT JOIN translations t ON
            t.entity_type = 'attribute'
            AND t.entity_id = a.id
            AND t.field_name = 'display_name'
        WHERE v.listing_id = $1
        GROUP BY
            v.listing_id, v.attribute_id, a.id, a.name, a.display_name,
            a.attribute_type, v.text_value, v.numeric_value, v.boolean_value,
            v.json_value, v.unit, a.is_required, a.show_in_card, a.show_in_list,
            cam.is_required, cam.show_in_card, cam.show_in_list
        ORDER BY a.id, a.sort_order, a.display_name
    `

	rows, err := s.pool.Query(ctx, query, listingID)
	if err != nil {
		return nil, fmt.Errorf("error querying listing attributes: %w", err)
	}
	defer rows.Close()

	log.Printf("Запрос атрибутов для объявления %d", listingID)
	var allAttributes []models.ListingAttributeValue

	// Для защиты от дубликатов по ID атрибута используем карту
	seen := make(map[int]bool)

	for rows.Next() {
		var attr models.ListingAttributeValue
		var textValue sql.NullString
		var numericValue sql.NullFloat64
		var boolValue sql.NullBool
		var jsonValue sql.NullString
		var unit sql.NullString
		var translationsJson []byte
		var optionTranslationsJson []byte

		if err := rows.Scan(
			&attr.ListingID,
			&attr.AttributeID,
			&attr.AttributeName,
			&attr.DisplayName,
			&attr.AttributeType,
			&textValue,
			&numericValue,
			&boolValue,
			&jsonValue,
			&unit,
			&attr.IsRequired,
			&attr.ShowInCard,
			&attr.ShowInList,
			&translationsJson,
			&optionTranslationsJson,
		); err != nil {
			log.Printf("Error scanning attribute: %v", err)
			return nil, fmt.Errorf("error scanning listing attribute: %w", err)
		}

		// Добавляем переводы атрибута
		if err := json.Unmarshal(translationsJson, &attr.Translations); err != nil {
			log.Printf("Error unmarshal attribute translations: %v", err)
			attr.Translations = make(map[string]string)
		}

		// Добавляем переводы опций атрибута
		if err := json.Unmarshal(optionTranslationsJson, &attr.OptionTranslations); err != nil {
			log.Printf("Error unmarshal option translations: %v", err)
			attr.OptionTranslations = make(map[string]map[string]string)
		}

		// Проверяем, не добавляли ли мы уже этот атрибут
		if seen[attr.AttributeID] {
			log.Printf("WARNING: Skipping duplicate attribute ID=%d, Name=%s",
				attr.AttributeID, attr.AttributeName)
			continue
		}
		seen[attr.AttributeID] = true

		// Сохраняем единицу измерения
		if unit.Valid {
			attr.Unit = unit.String
		}

		// Заполняем значения в зависимости от типа
		if textValue.Valid {
			attr.TextValue = &textValue.String
			attr.DisplayValue = textValue.String
			log.Printf("DEBUG: Attribute %d (%s) has text value: %s",
				attr.AttributeID, attr.AttributeName, textValue.String)
		}

		if numericValue.Valid {
			attr.NumericValue = &numericValue.Float64

			// Создаем отображаемое значение с учетом единиц измерения
			unitStr := attr.Unit // Используем сохраненную единицу или определяем по имени атрибута
			if unitStr == "" {
				// Определяем единицу измерения на основе имени атрибута
				switch attr.AttributeName {
				case attrNameArea:
					unitStr = "m²"
				case attrNameLandArea:
					unitStr = "ar"
				case attrNameMileage:
					unitStr = "km"
				case attrNameEngineCapacity:
					unitStr = "l"
				case attrNamePower:
					unitStr = "ks"
				case "screen_size":
					unitStr = "inč"
				}
			}

			// Форматируем отображаемое значение с учетом типа
			if attr.AttributeName == attrNameYear {
				attr.DisplayValue = fmt.Sprintf("%d", int(numericValue.Float64))
			} else if unitStr != "" {
				attr.DisplayValue = fmt.Sprintf("%g %s", numericValue.Float64, unitStr)
			} else {
				attr.DisplayValue = fmt.Sprintf("%g", numericValue.Float64)
			}

			log.Printf("DEBUG: Attribute %d (%s) has numeric value: %f, display: %s",
				attr.AttributeID, attr.AttributeName, numericValue.Float64, attr.DisplayValue)
		}

		if boolValue.Valid {
			attr.BooleanValue = &boolValue.Bool
			if boolValue.Bool {
				attr.DisplayValue = "Да"
			} else {
				attr.DisplayValue = "Нет"
			}
			log.Printf("DEBUG: Attribute %d (%s) has boolean value: %t",
				attr.AttributeID, attr.AttributeName, boolValue.Bool)
		}

		if jsonValue.Valid {
			attr.JSONValue = json.RawMessage(jsonValue.String)
			// Для multiselect можно форматировать массив значений
			if attr.AttributeType == "multiselect" {
				var values []string
				if err := json.Unmarshal(attr.JSONValue, &values); err == nil {
					attr.DisplayValue = strings.Join(values, ", ")
				}
			} else {
				attr.DisplayValue = jsonValue.String
			}
			log.Printf("DEBUG: Attribute %d (%s) has JSON value: %s",
				attr.AttributeID, attr.AttributeName, jsonValue.String)
		}

		allAttributes = append(allAttributes, attr)
	}

	log.Printf("DEBUG: Found %d unique attributes for listing %d", len(allAttributes), listingID)

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating listing attributes: %w", err)
	}

	return allAttributes, nil
}

// GetAttributeRanges получает минимальные и максимальные значения для числовых атрибутов
func (s *Storage) GetAttributeRanges(ctx context.Context, categoryID int) (map[string]map[string]interface{}, error) {
	// Проверяем наличие в кеше
	rangesCacheMutex.RLock()
	cachedRanges, hasCached := rangesCache[categoryID]
	cacheTime, hasTime := rangesCacheTime[categoryID]
	rangesCacheMutex.RUnlock()

	// Если данные в кеше и они не устарели
	if hasCached && hasTime && time.Since(cacheTime) < cacheTTL {
		log.Printf("Using cached attribute ranges for category %d", categoryID)
		return cachedRanges, nil
	}

	// Получаем ID всех подкатегорий заданной категории
	query := `
    WITH RECURSIVE category_tree AS (
        SELECT id FROM marketplace_categories WHERE id = $1
        UNION ALL
        SELECT c.id FROM marketplace_categories c
        JOIN category_tree t ON c.parent_id = t.id
    )
    SELECT string_agg(id::text, ',') FROM category_tree
    `

	var categoryIDs string
	err := s.pool.QueryRow(ctx, query, categoryID).Scan(&categoryIDs)
	if err != nil {
		return nil, fmt.Errorf("error getting category tree: %w", err)
	}

	if categoryIDs == "" {
		categoryIDs = strconv.Itoa(categoryID)
	}

	// Запрос для получения границ числовых атрибутов
	rangesQuery := `
    SELECT
        a.name,
        MIN(v.numeric_value) as min_value,
        MAX(v.numeric_value) as max_value,
        COUNT(DISTINCT v.numeric_value) as value_count
    FROM listing_attribute_values v
    JOIN category_attributes a ON v.attribute_id = a.id
    JOIN marketplace_listings l ON v.listing_id = l.id
    WHERE
        l.category_id IN (` + categoryIDs + `)
        AND l.status = 'active'
        AND v.numeric_value IS NOT NULL
        AND a.attribute_type = 'number'
    GROUP BY a.name
    `

	rows, err := s.pool.Query(ctx, rangesQuery)
	if err != nil {
		return nil, fmt.Errorf("error querying attribute ranges: %w", err)
	}
	defer rows.Close()

	// Формируем результат
	ranges := make(map[string]map[string]interface{})

	for rows.Next() {
		var attrName string
		var minValue, maxValue float64
		var valueCount int

		if err := rows.Scan(&attrName, &minValue, &maxValue, &valueCount); err != nil {
			return nil, fmt.Errorf("error scanning attribute range: %w", err)
		}

		// Для года и других целочисленных параметров округляем границы
		if attrName == attrNameYear || attrName == "rooms" || attrName == "floor" || attrName == "total_floors" {
			minValue = float64(int(minValue))
			maxValue = float64(int(maxValue))
		}

		// Для "year" добавляем запас +1 год для новых автомобилей
		if attrName == attrNameYear && maxValue >= float64(time.Now().Year()-1) {
			maxValue = float64(time.Now().Year() + 1)
		}

		// Установка разумных шагов в зависимости от диапазона
		var step float64
		switch attrName {
		case "engine_capacity":
			step = 0.1
		case "area", "land_area":
			step = 0.5
		default:
			step = 1.0
		}

		// Создаем информацию о границах
		ranges[attrName] = map[string]interface{}{
			"min":   minValue,
			"max":   maxValue,
			"step":  step,
			"count": valueCount,
		}

		log.Printf("Attribute %s range: min=%.2f, max=%.2f, values=%d",
			attrName, minValue, maxValue, valueCount)
	}

	// Если нет данных, устанавливаем разумные значения по умолчанию
	defaultRanges := map[string]map[string]interface{}{
		attrNameYear:      {"min": float64(time.Now().Year() - 30), "max": float64(time.Now().Year() + 1), "step": 1.0},
		attrNameMileage:   {"min": 0.0, "max": 500000.0, "step": 1000.0},
		"engine_capacity": {"min": 0.5, "max": 8.0, "step": 0.1},
		attrNamePower:     {"min": 50.0, "max": 500.0, "step": 10.0},
		"rooms":           {"min": 1.0, "max": 10.0, "step": 1.0},
		"floor":           {"min": 1.0, "max": 25.0, "step": 1.0},
		"total_floors":    {"min": 1.0, "max": 30.0, "step": 1.0},
		"area":            {"min": 10.0, "max": 300.0, "step": 0.5},
		"land_area":       {"min": 1.0, "max": 100.0, "step": 0.5},
	}

	// Заполняем отсутствующие атрибуты значениями по умолчанию
	for attr, defaultRange := range defaultRanges {
		if _, exists := ranges[attr]; !exists {
			ranges[attr] = defaultRange
			log.Printf("No data for attribute %s, using defaults: min=%.2f, max=%.2f",
				attr, defaultRange["min"], defaultRange["max"])
		}
	}

	// Кешируем результат
	rangesCacheMutex.Lock()
	rangesCache[categoryID] = ranges
	rangesCacheTime[categoryID] = time.Now()
	rangesCacheMutex.Unlock()

	return ranges, nil
}

// Добавим метод для очистки кеша атрибутов
func (s *Storage) InvalidateAttributesCache(categoryID int) {
	attributeCacheMutex.Lock()
	delete(attributeCache, categoryID)
	delete(attributeCacheTime, categoryID)
	attributeCacheMutex.Unlock()

	rangesCacheMutex.Lock()
	delete(rangesCache, categoryID)
	delete(rangesCacheTime, categoryID)
	rangesCacheMutex.Unlock()

	log.Printf("Invalidated attributes cache for category %d", categoryID)
}

// Исправленная версия функции GetCategoryAttributes
func (s *Storage) GetCategoryAttributes(ctx context.Context, categoryID int) ([]models.CategoryAttribute, error) {
	// Проверяем наличие в кеше
	attributeCacheMutex.RLock()
	cachedAttrs, hasCached := attributeCache[categoryID]
	cacheTime, hasTime := attributeCacheTime[categoryID]
	attributeCacheMutex.RUnlock()

	// Если данные в кеше и они не устарели
	if hasCached && hasTime && time.Since(cacheTime) < cacheTTL {
		log.Printf("Using cached attributes for category %d", categoryID)
		return cachedAttrs, nil
	}

	// Добавляем логирование для отладки
	log.Printf("GetCategoryAttributes: Получение атрибутов для категории %d", categoryID)

	query := `
    WITH RECURSIVE category_hierarchy AS (
        -- Находим все родительские категории (включая текущую)
        WITH RECURSIVE parents AS (
            SELECT id, parent_id
            FROM marketplace_categories
            WHERE id = $1

            UNION

            SELECT c.id, c.parent_id
            FROM marketplace_categories c
            INNER JOIN parents p ON c.id = p.parent_id
        )
        SELECT id FROM parents
    ),
    attribute_translations AS (
        SELECT
            entity_id,
            jsonb_object_agg(language, translated_text) as translations
        FROM translations
        WHERE entity_type = 'attribute'
        AND field_name = 'display_name'
        GROUP BY entity_id
    ),
    option_translations AS (
        SELECT
            entity_id,
            language,
            jsonb_object_agg(field_name, translated_text) as field_translations
        FROM translations
        WHERE entity_type = 'attribute_option'
        GROUP BY entity_id, language
    ),
    option_lang_agg AS (
        SELECT
            entity_id,
            jsonb_object_agg(language, field_translations) as option_translations
        FROM option_translations
        GROUP BY entity_id
    )
    SELECT DISTINCT ON (a.id)
        a.id,
        a.name,
        a.display_name,
        a.icon,
        a.attribute_type,
        a.options,
        a.validation_rules,
        a.is_searchable,
        a.is_filterable,
        COALESCE(m.is_required, a.is_required) as is_required,
        a.sort_order,
        a.created_at,
        COALESCE(m.custom_component, a.custom_component) as custom_component,
        COALESCE(at.translations, '{}'::jsonb) as translations,
        COALESCE(ol.option_translations, '{}'::jsonb) as option_translations
    FROM category_attribute_mapping m
    JOIN category_attributes a ON m.attribute_id = a.id
    JOIN category_hierarchy h ON m.category_id = h.id
    LEFT JOIN attribute_translations at ON a.id = at.entity_id
    LEFT JOIN option_lang_agg ol ON a.id = ol.entity_id
    WHERE m.is_enabled = true
    ORDER BY a.id, m.category_id = $1 DESC, a.sort_order, a.display_name
    `

	// Выполняем запрос
	rows, err := s.pool.Query(ctx, query, categoryID)
	if err != nil {
		log.Printf("GetCategoryAttributes: Ошибка запроса: %v", err)
		return nil, fmt.Errorf("error querying category attributes: %w", err)
	}
	defer rows.Close()

	var attributes []models.CategoryAttribute
	for rows.Next() {
		var attr models.CategoryAttribute
		var options, validRules, customComponent sql.NullString
		var translationsJson, optionTranslationsJson []byte

		if err := rows.Scan(
			&attr.ID,
			&attr.Name,
			&attr.DisplayName,
			&attr.Icon,
			&attr.AttributeType,
			&options,
			&validRules,
			&attr.IsSearchable,
			&attr.IsFilterable,
			&attr.IsRequired,
			&attr.SortOrder,
			&attr.CreatedAt,
			&customComponent,
			&translationsJson,
			&optionTranslationsJson,
		); err != nil {
			log.Printf("GetCategoryAttributes: Ошибка при сканировании результата: %v", err)
			return nil, fmt.Errorf("error scanning category attribute: %w", err)
		}

		// Обработка опциональных JSON полей
		if options.Valid && len(options.String) > 0 {
			attr.Options = json.RawMessage(options.String)
		} else {
			// Если options пустой или не валидный, устанавливаем пустой JSON объект
			attr.Options = json.RawMessage(`{}`)
		}

		if validRules.Valid && len(validRules.String) > 0 {
			attr.ValidRules = json.RawMessage(validRules.String)
		} else {
			attr.ValidRules = json.RawMessage(`{}`)
		}
		// Всегда инициализируем CustomComponent пустой строкой
		attr.CustomComponent = ""
		if customComponent.Valid {
			attr.CustomComponent = customComponent.String
		}

		// Обработка переводов
		attr.Translations = make(map[string]string)
		if err := json.Unmarshal(translationsJson, &attr.Translations); err != nil {
			log.Printf("GetCategoryAttributes: Ошибка парсинга переводов для атрибута %d: %v", attr.ID, err)
		}

		// Обработка переводов опций
		attr.OptionTranslations = make(map[string]map[string]string)
		if err := json.Unmarshal(optionTranslationsJson, &attr.OptionTranslations); err != nil {
			log.Printf("GetCategoryAttributes: Ошибка парсинга переводов опций для атрибута %d: %v", attr.ID, err)
		} else {
			log.Printf("GetCategoryAttributes: Получены переводы опций для атрибута %d: %v", attr.ID, attr.OptionTranslations)
		}

		attributes = append(attributes, attr)
	}

	if err := rows.Err(); err != nil {
		log.Printf("GetCategoryAttributes: Ошибка при итерации результатов: %v", err)
		return nil, fmt.Errorf("error iterating category attributes: %w", err)
	}

	// Получаем актуальные диапазоны для атрибутов
	attributeRanges, err := s.GetAttributeRanges(ctx, categoryID)
	if err != nil {
		log.Printf("GetCategoryAttributes: Ошибка получения диапазонов атрибутов: %v", err)
		// Продолжаем работу без диапазонов
	} else {
		// Обновляем опции атрибутов с учетом реальных диапазонов
		for i, attr := range attributes {
			// Обрабатываем только числовые атрибуты
			if attr.AttributeType == "number" {
				if ranges, ok := attributeRanges[attr.Name]; ok {
					// Создаем или обновляем объект options
					var options map[string]interface{}

					// Пытаемся использовать существующие options
					if len(attr.Options) > 0 {
						err := json.Unmarshal(attr.Options, &options)
						if err != nil {
							options = make(map[string]interface{})
						}
					} else {
						options = make(map[string]interface{})
					}

					// Обновляем значения диапазонов
					options["min"] = ranges["min"]
					options["max"] = ranges["max"]
					options["step"] = ranges["step"]
					options["real_data"] = true

					// Сериализуем обратно в JSON
					optionsJSON, err := json.Marshal(options)
					if err == nil {
						attributes[i].Options = optionsJSON
						log.Printf("GetCategoryAttributes: Обновлены диапазоны для атрибута %s: min=%.2f, max=%.2f",
							attr.Name, ranges["min"], ranges["max"])
					}
				}
			}
		}
	}

	// Кешируем результат
	attributeCacheMutex.Lock()
	attributeCache[categoryID] = attributes
	attributeCacheTime[categoryID] = time.Now()
	attributeCacheMutex.Unlock()

	log.Printf("GetCategoryAttributes: Успешно получено %d атрибутов для категории %d", len(attributes), categoryID)
	return attributes, nil
}

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
                    t.language,  -- Изменено: убираем конкатенацию с field_name
                    t.translated_text
                ) as translations
            FROM translations t
            WHERE t.entity_type = 'marketplace_category'
            AND t.field_name = 'name'
            GROUP BY t.entity_id
        )
        SELECT
            c.id, c.name, c.slug, c.parent_id, c.icon, c.description, c.is_active, c.created_at,
            c.seo_title, c.seo_description, c.seo_keywords,
            COALESCE(ct.translations, '{}'::jsonb) as translations
        FROM marketplace_categories c
        LEFT JOIN category_translations ct ON c.id = ct.entity_id
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
		)
		if err != nil {
			log.Printf("GetCategories: Error scanning category: %v", err)
			continue
		}

		// Обрабатываем NULL значения
		if icon.Valid {
			cat.Icon = icon.String
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

		//    log.Printf("Raw translations for category %d: %s", cat.ID, string(translationsJson))

		translations := make(map[string]string)
		if err := json.Unmarshal(translationsJson, &translations); err != nil {
			log.Printf("GetCategories: Error unmarshaling translations for category %d: %v", cat.ID, err)
		} else {
			cat.Translations = translations

			// Применяем перевод к названию категории если он есть для запрашиваемого языка
			if translatedName, ok := translations[locale]; ok && translatedName != "" {
				log.Printf("GetCategories: Applying translation for category %d: %s -> %s (locale: %s)",
					cat.ID, cat.Name, translatedName, locale)
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
            FROM marketplace_categories c
            LEFT JOIN translations t ON t.entity_id = c.id AND t.entity_type = 'marketplace_category' AND t.field_name = 'name'
            GROUP BY c.id
        )
        SELECT
            c.id, c.name, c.slug, c.parent_id, c.icon, c.description, c.is_active, c.created_at,
            c.seo_title, c.seo_description, c.seo_keywords,
            COALESCE(ct.translations, '{}'::jsonb) as translations
        FROM marketplace_categories c
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
			cat.Icon = icon.String
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

			// Применяем перевод к названию категории если он есть для запрашиваемого языка
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

func (s *Storage) GetCategoryByID(ctx context.Context, id int) (*models.MarketplaceCategory, error) {
	cat := &models.MarketplaceCategory{}
	var icon, description sql.NullString

	var seoTitle, seoDescription, seoKeywords sql.NullString
	err := s.pool.QueryRow(ctx, `
        SELECT
            id, name, slug, parent_id, icon, description, is_active, created_at,
            seo_title, seo_description, seo_keywords
        FROM marketplace_categories
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
		cat.Icon = icon.String
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

func (s *Storage) GetListingByID(ctx context.Context, id int) (*models.MarketplaceListing, error) {
	listing := &models.MarketplaceListing{
		User:     &models.User{},
		Category: &models.MarketplaceCategory{},
	}

	// Получаем основные данные объявления с original_language
	var (
		description    sql.NullString
		condition      sql.NullString
		status         sql.NullString
		location       sql.NullString
		latitude       sql.NullFloat64
		longitude      sql.NullFloat64
		city           sql.NullString
		country        sql.NullString
		originalLang   sql.NullString
		userName       sql.NullString
		userEmail      sql.NullString
		userPictureURL sql.NullString
		userPhone      sql.NullString
		categoryName   sql.NullString
		categorySlug   sql.NullString
		storefrontID   sql.NullInt32
	)

	err := s.pool.QueryRow(ctx, `
        SELECT
            l.id, l.user_id, l.category_id, l.title, l.description,
            l.price, l.condition, l.status, l.location, l.latitude,
            l.longitude, l.address_city as city, l.address_country as country, l.views_count,
            l.created_at, l.updated_at, l.show_on_map, l.original_language,
            u.name, u.email, u.created_at as user_created_at,
            u.picture_url, u.phone,
            c.name as category_name, c.slug as category_slug, l.metadata, l.storefront_id
        FROM marketplace_listings l
        LEFT JOIN users u ON l.user_id = u.id
        LEFT JOIN marketplace_categories c ON l.category_id = c.id
        WHERE l.id = $1
    `, id).Scan(
		&listing.ID, &listing.UserID, &listing.CategoryID, &listing.Title,
		&description, &listing.Price, &condition, &status,
		&location, &latitude, &longitude, &city,
		&country, &listing.ViewsCount, &listing.CreatedAt, &listing.UpdatedAt,
		&listing.ShowOnMap, &originalLang,
		&userName, &userEmail, &listing.User.CreatedAt,
		&userPictureURL, &userPhone,
		&categoryName, &categorySlug, &listing.Metadata, &storefrontID,
	)
	log.Printf("999 DEBUG: Listing %d metadata: %+v", id, listing.Metadata)
	log.Printf("DEBUG: err after query = %v", err)

	if err != nil {
		// Если не найдено в marketplace_listings, попробуем найти в storefront_products
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
	if userName.Valid {
		listing.User.Name = userName.String
	}
	if userEmail.Valid {
		listing.User.Email = userEmail.String
	}
	if userPictureURL.Valid {
		listing.User.PictureURL = userPictureURL.String
	}
	if userPhone.Valid {
		listing.User.Phone = &userPhone.String
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

	// Обработка метаданных и скидок для согласованного отображения
	if listing.Metadata != nil {
		if discount, ok := listing.Metadata["discount"].(map[string]interface{}); ok {
			listing.HasDiscount = true
			if prevPrice, ok := discount["previous_price"].(float64); ok {
				listing.OldPrice = prevPrice

				// Пересчитываем актуальный процент скидки, чтобы он всегда соответствовал
				// текущей разнице между старой и новой ценой
				if prevPrice > listing.Price {
					discountPercent := int((prevPrice - listing.Price) / prevPrice * 100)
					discount["discount_percent"] = discountPercent

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
	log.Printf("DEBUG GetListingByID: listing.ID=%d, StorefrontID=%v", listing.ID, listing.StorefrontID)
	if listing.StorefrontID != nil && *listing.StorefrontID > 0 {
		log.Printf("DEBUG GetListingByID: Loading storefront images for listing %d, storefront_id=%d", listing.ID, *listing.StorefrontID)
		// Для storefront товаров загружаем изображения из storefront_product_images
		storefrontImages, err := s.GetStorefrontProductImages(ctx, listing.ID)
		if err != nil {
			log.Printf("Error loading storefront images for listing %d: %v", listing.ID, err)
			// Не прерываем выполнение, просто оставляем пустой массив изображений
			images = []models.MarketplaceImage{}
		} else {
			log.Printf("DEBUG GetListingByID: Successfully loaded %d storefront images for listing %d", len(storefrontImages), listing.ID)
			images = storefrontImages
		}
	} else {
		log.Printf("DEBUG GetListingByID: Loading marketplace images for listing %d (no storefront_id)", listing.ID)
		// Для обычных marketplace товаров загружаем изображения из marketplace_images
		marketplaceImages, err := s.GetListingImages(ctx, fmt.Sprintf("%d", listing.ID))
		if err != nil {
			log.Printf("Error loading marketplace images for listing %d: %v", listing.ID, err)
			// Не прерываем выполнение, просто оставляем пустой массив изображений
			images = []models.MarketplaceImage{}
		} else {
			log.Printf("DEBUG GetListingByID: Successfully loaded %d marketplace images for listing %d", len(marketplaceImages), listing.ID)
			images = marketplaceImages
		}
	}

	// Важно! Присваиваем изображения объявлению
	listing.Images = images
	listing.Translations = translations

	// Достаем информацию о пути категории
	if listing.CategoryID > 0 {
		query := `
        WITH RECURSIVE category_path AS (
            SELECT id, name, slug, parent_id, 1 as level
            FROM marketplace_categories
            WHERE id = $1

            UNION ALL

            SELECT c.id, c.name, c.slug, c.parent_id, cp.level + 1
            FROM marketplace_categories c
            JOIN category_path cp ON c.id = cp.parent_id
        )
        SELECT id, name, slug
        FROM category_path
        ORDER BY level DESC
    `
		rows, err := s.pool.Query(ctx, query, listing.CategoryID)
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

	return listing, nil
}

// getStorefrontProductAsListing получает товар из storefront_products и возвращает как MarketplaceListing
func (s *Storage) getStorefrontProductAsListing(ctx context.Context, id int) (*models.MarketplaceListing, error) {
	listing := &models.MarketplaceListing{
		User:     &models.User{},
		Category: &models.MarketplaceCategory{},
	}

	// Получаем данные товара из storefront_products
	var categoryName, categorySlug sql.NullString

	err := s.pool.QueryRow(ctx, `
        SELECT
            sp.id, sp.storefront_id, 0 as user_id, sp.category_id, sp.name, sp.description,
            sp.price, 'new' as condition, 'active' as status, '' as location,
            0 as latitude, 0 as longitude, '' as city, '' as country,
            sp.view_count, sp.created_at, sp.updated_at, false as show_on_map, 'sr' as original_language,
            '' as user_name, '' as user_email, sp.created_at as user_created_at,
            '' as user_picture_url, '' as user_phone,
            c.name as category_name, c.slug as category_slug, '{}'::jsonb as metadata
        FROM storefront_products sp
        LEFT JOIN marketplace_categories c ON sp.category_id = c.id
        WHERE sp.id = $1 AND sp.is_active = true
    `, id).Scan(
		&listing.ID, &listing.StorefrontID, &listing.UserID, &listing.CategoryID, &listing.Title,
		&listing.Description, &listing.Price, &listing.Condition, &listing.Status,
		&listing.Location, &listing.Latitude, &listing.Longitude, &listing.City,
		&listing.Country, &listing.ViewsCount, &listing.CreatedAt, &listing.UpdatedAt,
		&listing.ShowOnMap, &listing.OriginalLanguage,
		&listing.User.Name, &listing.User.Email, &listing.User.CreatedAt,
		&listing.User.PictureURL, &listing.User.Phone,
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

	// Загружаем изображения для storefront продукта
	log.Printf("DEBUG getStorefrontProductAsListing: Loading images for storefront product %d", listing.ID)
	storefrontImages, err := s.GetStorefrontProductImages(ctx, listing.ID)
	if err != nil {
		log.Printf("Error loading storefront images for product %d: %v", listing.ID, err)
		listing.Images = []models.MarketplaceImage{}
	} else {
		log.Printf("DEBUG getStorefrontProductAsListing: Successfully loaded %d images for storefront product %d", len(storefrontImages), listing.ID)
		listing.Images = storefrontImages
	}

	// Для storefront продуктов нет атрибутов пока что
	listing.Attributes = []models.ListingAttributeValue{}
	listing.Translations = make(map[string]map[string]string)

	return listing, nil
}

// GetPopularSearchQueries возвращает популярные поисковые запросы
func (s *Storage) GetPopularSearchQueries(ctx context.Context, query string, limit int) ([]service.SearchQuery, error) {
	normalizedQuery := strings.ToLower(strings.TrimSpace(query))

	sqlQuery := `
		SELECT
			id,
			query,
			normalized_query,
			search_count,
			to_char(last_searched, 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"') as last_searched,
			language,
			results_count
		FROM search_queries
		WHERE normalized_query LIKE '%' || $1 || '%'
		ORDER BY search_count DESC
		LIMIT $2
	`

	rows, err := s.pool.Query(ctx, sqlQuery, normalizedQuery, limit)
	if err != nil {
		return nil, fmt.Errorf("error querying popular searches: %w", err)
	}
	defer rows.Close()

	var queries []service.SearchQuery
	for rows.Next() {
		var q service.SearchQuery
		if err := rows.Scan(
			&q.ID,
			&q.Query,
			&q.NormalizedQuery,
			&q.SearchCount,
			&q.LastSearched,
			&q.Language,
			&q.ResultsCount,
		); err != nil {
			return nil, fmt.Errorf("error scanning search query: %w", err)
		}
		queries = append(queries, q)
	}

	return queries, nil
}

// SaveSearchQuery сохраняет или обновляет поисковый запрос
func (s *Storage) SaveSearchQuery(ctx context.Context, query, normalizedQuery string, resultsCount int, language string) error {
	if normalizedQuery == "" {
		normalizedQuery = strings.ToLower(strings.TrimSpace(query))
	}

	if normalizedQuery == "" {
		return nil // Не сохраняем пустые запросы
	}

	// Используем UPSERT для обновления существующих записей
	sqlQuery := `
		INSERT INTO search_queries (
			query, normalized_query, search_count, last_searched,
			language, results_count
		) VALUES ($1, $2, 1, NOW(), $3, $4)
		ON CONFLICT (normalized_query, language)
		DO UPDATE SET
			query = EXCLUDED.query,
			search_count = search_queries.search_count + 1,
			last_searched = NOW(),
			results_count = EXCLUDED.results_count
	`

	_, err := s.pool.Exec(ctx, sqlQuery, query, normalizedQuery, language, resultsCount)
	if err != nil {
		return fmt.Errorf("error saving search query: %w", err)
	}

	return nil
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
			FROM marketplace_listings
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
			WHERE entity_type = 'marketplace_category'
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
		FROM marketplace_categories c
		LEFT JOIN category_counts cc ON c.id = cc.category_id
		LEFT JOIN category_translations ct ON c.id = ct.entity_id
		WHERE LOWER(c.name) LIKE $1
			OR EXISTS (
				SELECT 1
				FROM translations t
				WHERE t.entity_type = 'marketplace_category'
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
			cat.Icon = icon.String
		}

		// Парсим переводы
		translations := make(map[string]string)
		if err := json.Unmarshal(translationsJson, &translations); err != nil {
			log.Printf("Error unmarshaling translations for category %d: %v", cat.ID, err)
		} else {
			cat.Translations = translations

			// Применяем перевод к названию категории если он есть для запрашиваемого языка
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

// GetStorefrontProductImages загружает изображения для storefront товара и конвертирует их в MarketplaceImage
func (s *Storage) GetStorefrontProductImages(ctx context.Context, productID int) ([]models.MarketplaceImage, error) {
	query := `
        SELECT
            id, storefront_product_id, image_url, thumbnail_url, 
            display_order, is_default, created_at
        FROM storefront_product_images
        WHERE storefront_product_id = $1
        ORDER BY display_order ASC, id ASC
    `

	rows, err := s.pool.Query(ctx, query, productID)
	if err != nil {
		return nil, fmt.Errorf("error querying storefront product images: %w", err)
	}
	defer rows.Close()

	var images []models.MarketplaceImage
	for rows.Next() {
		var img struct {
			ID                  int
			StorefrontProductID int
			ImageURL            string
			ThumbnailURL        string
			DisplayOrder        int
			IsDefault           bool
			CreatedAt           time.Time
		}

		err := rows.Scan(
			&img.ID, &img.StorefrontProductID, &img.ImageURL, &img.ThumbnailURL,
			&img.DisplayOrder, &img.IsDefault, &img.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning storefront product image: %w", err)
		}

		// Конвертируем в MarketplaceImage
		marketplaceImage := models.MarketplaceImage{
			ID:          img.ID,
			ListingID:   img.StorefrontProductID,
			PublicURL:   img.ImageURL,
			IsMain:      img.IsDefault,
			StorageType: "minio",
			CreatedAt:   img.CreatedAt,
		}

		images = append(images, marketplaceImage)
	}

	return images, nil
}

// backend/internal/proj/marketplace/storage/postgres/marketplace.go
package postgres

import (
	"backend/internal/domain/models"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"regexp"
	"backend/internal/proj/marketplace/service"
	"log"
	"strings"
	"math"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	//"time"
	// "github.com/jackc/pgx/v5"
)

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
func (s *Storage) GetSubcategories(ctx context.Context, parentID *int, limit int, offset int) ([]models.CategoryTreeNode, error) {
	query := `
    WITH RECURSIVE category_path AS (
        SELECT 
            c.id,
            c.name,
            c.slug,
            c.icon,
            c.parent_id,
            c.created_at,
            ARRAY[c.id] as path,
            1 as level,
            (SELECT COUNT(*) FROM marketplace_listings ml WHERE ml.category_id = c.id AND ml.status = 'active') as listing_count,
            (SELECT COUNT(*) FROM marketplace_categories sc WHERE sc.parent_id = c.id) as children_count
        FROM marketplace_categories c
        WHERE CASE 
            WHEN $1::int IS NULL THEN c.parent_id IS NULL
            ELSE c.parent_id = $1
        END
    )
    SELECT 
        cp.*,
        COALESCE(
            jsonb_object_agg(
                t.language, 
                t.translated_text
            ) FILTER (WHERE t.language IS NOT NULL),
            '{}'::jsonb
        ) as translations
    FROM category_path cp
    LEFT JOIN translations t ON 
        t.entity_type = 'category' 
        AND t.entity_id = cp.id 
        AND t.field_name = 'name'
    GROUP BY 
        cp.id, cp.name, cp.slug, cp.icon, cp.parent_id, 
        cp.created_at, cp.path, cp.level, cp.listing_count, 
        cp.children_count
    ORDER BY cp.name`

	rows, err := s.pool.Query(ctx, query, parentID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error querying subcategories: %w", err)
	}
	defer rows.Close()

	var categories []models.CategoryTreeNode
	for rows.Next() {
		var cat models.CategoryTreeNode
		var pathArray []int
		var translationsJson []byte
		var directCount int // Добавляем эту переменную

		err := rows.Scan(
			&cat.ID,
			&cat.Name,
			&cat.Slug,
			&cat.Icon,
			&cat.ParentID,
			&cat.CreatedAt,
			&pathArray,
			&cat.Level,
			&directCount,
			&cat.ChildrenCount,
			&cat.ListingCount,
			&translationsJson,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning category: %w", err)
		}

		if err := json.Unmarshal(translationsJson, &cat.Translations); err != nil {
			cat.Translations = make(map[string]string)
		}

		categories = append(categories, cat)
	}

	return categories, nil
}

func (s *Storage) CreateListing(ctx context.Context, listing *models.MarketplaceListing) (int, error) {
	var listingID int

	// Если не указан язык, берем значение из контекста или используем по умолчанию
	if listing.OriginalLanguage == "" {
		if userLang, ok := ctx.Value("language").(string); ok && userLang != "" {
			listing.OriginalLanguage = userLang
		} else {
			listing.OriginalLanguage = "ru" // Значение по умолчанию
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

	if listing.Attributes != nil && len(listing.Attributes) > 0 {
		// Устанавливаем ID объявления для каждого атрибута
		for i := range listing.Attributes {
			listing.Attributes[i].ListingID = listingID
		}

		if err := s.SaveListingAttributes(ctx, listingID, listing.Attributes); err != nil {
			log.Printf("Error saving attributes for listing %d: %v", listingID, err)
			// Не прерываем создание объявления из-за ошибки с атрибутами
		}
	}

	return listingID, nil
}

func (s *Storage) AddListingImage(ctx context.Context, image *models.MarketplaceImage) (int, error) {
	var imageID int
	err := s.pool.QueryRow(ctx, `
        INSERT INTO marketplace_images (
            listing_id, file_path, file_name, file_size, 
            content_type, is_main
        ) VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `,
		image.ListingID, image.FilePath, image.FileName,
		image.FileSize, image.ContentType, image.IsMain,
	).Scan(&imageID)

	return imageID, err
}

func (s *Storage) GetListingImages(ctx context.Context, listingID string) ([]models.MarketplaceImage, error) {
	// Преобразуем listingID в int, так как принимаем строку
	id, err := strconv.Atoi(listingID)
	if err != nil {
		return nil, fmt.Errorf("invalid listing ID: %w", err)
	}

	query := `
        SELECT 
            id, 
            listing_id,
            file_path,
            file_name,
            file_size,
            content_type,
            is_main,
            created_at
        FROM marketplace_images
        WHERE listing_id = $1
        ORDER BY is_main DESC, created_at DESC
    `

	rows, err := s.pool.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("error querying images: %w", err)
	}
	defer rows.Close()

	var images []models.MarketplaceImage
	for rows.Next() {
		var img models.MarketplaceImage
		err := rows.Scan(
			&img.ID,
			&img.ListingID,
			&img.FilePath,
			&img.FileName,
			&img.FileSize,
			&img.ContentType,
			&img.IsMain,
			&img.CreatedAt,
		)
		if err != nil {
			log.Printf("Error scanning image row: %v", err)
			continue
		}
		images = append(images, img)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over images: %w", err)
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
		)

		err := rows.Scan(
			&listing.ID,
			&listing.UserID,
			&listing.CategoryID,
			&listing.Title,
			&listing.Description,
			&listing.Price,
			&listing.Condition,
			&listing.Status,
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
			&listing.User.Name,
			&tempEmail,
			&listing.User.CreatedAt,
			&tempPictureURL,
			&tempCategoryName,
			&tempCategorySlug,
			&translationsJSON,
			&listing.IsFavorite,
			&totalCount,
		)
		// После Scan добавьте обработку метаданных
		if len(metadataJSON) > 0 {
			if err := json.Unmarshal(metadataJSON, &listing.Metadata); err != nil {
				log.Printf("Error unmarshaling metadata for listing %d: %v", listing.ID, err)
			} else {
				// Обработка метаданных о скидках
				if listing.Metadata != nil {
					if discount, ok := listing.Metadata["discount"].(map[string]interface{}); ok {
						listing.HasDiscount = true
						if prevPrice, ok := discount["previous_price"].(float64); ok {
							listing.OldPrice = prevPrice
						}
					}
				}
			}
		}
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning listing: %w", err)
		}

		// Получаем изображения
		images, err := s.GetListingImages(ctx, fmt.Sprintf("%d", listing.ID))
		if err != nil {
			log.Printf("Error getting images for listing %d: %v", listing.ID, err)
			listing.Images = []models.MarketplaceImage{}
		} else {
			listing.Images = images
		}

		// Обработка NULL значений
		if tempLocation.Valid {
			listing.Location = tempLocation.String
		}
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
        t.entity_type = 'category' 
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

		err := rows.Scan(
			&node.ID,
			&node.Name,
			&node.Slug,
			&node.Icon,
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

		// Добавляем логирование переводов и детей
		//		log.Printf("Raw translations JSON for category %d: %s", node.ID, string(translationsJson))
		//		log.Printf("Raw children JSON for category %d: %s", node.ID, string(childrenJson))

		if err := json.Unmarshal(translationsJson, &node.Translations); err != nil {
			log.Printf("Error unmarshaling translations for category %d: %v", node.ID, err)
			node.Translations = make(map[string]string)
		}

		var children []models.CategoryTreeNode
		if err := json.Unmarshal(childrenJson, &children); err != nil {
			log.Printf("Error unmarshaling children for category %d: %v", node.ID, err)
			node.Children = make([]models.CategoryTreeNode, 0)
		} else {
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
            l.address_city, 
            l.address_country, 
            l.views_count,
            l.created_at, 
            l.updated_at,
            u.name, 
            u.email, 
			u.created_at as user_created_at, 
            COALESCE(u.picture_url, ''),
            c.name as category_name, 
            c.slug as category_slug,
            true as is_favorite
        FROM marketplace_listings l
        JOIN marketplace_favorites f ON l.id = f.listing_id
        LEFT JOIN users u ON l.user_id = u.id
        LEFT JOIN marketplace_categories c ON l.category_id = c.id
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
		)
		if err != nil {
			log.Printf("Error scanning listing: %v", err)
			continue
		}

		// Присваиваем отдельно
		listing.User.PictureURL = userPictureURL
		listing.User.ID = listing.UserID

		// Получаем изображения для каждого объявления
		images, err := s.GetListingImages(ctx, fmt.Sprintf("%d", listing.ID))
		if err != nil {
			log.Printf("Error getting images for listing %d: %v", listing.ID, err)
			listing.Images = []models.MarketplaceImage{} // Пустой массив вместо nil
		} else {
			listing.Images = images
		}

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
            address_city = $9,
            address_country = $10,
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

	if listing.Attributes != nil {
		if err := s.SaveListingAttributes(ctx, listing.ID, listing.Attributes); err != nil {
			log.Printf("Error updating attributes for listing %d: %v", listing.ID, err)
			// Не прерываем обновление объявления из-за ошибки с атрибутами
		}
	}

	return nil
}

// SaveListingAttributes сохраняет значения атрибутов для объявления
func (s *Storage) SaveListingAttributes(ctx context.Context, listingID int, attributes []models.ListingAttributeValue) error {
    log.Printf("Saving %d attributes for listing %d", len(attributes), listingID)

    // Начинаем транзакцию
    tx, err := s.pool.Begin(ctx)
    if err != nil {
        return fmt.Errorf("error starting transaction: %w", err)
    }
    defer tx.Rollback(ctx)

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

    for _, attr := range attributes {
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
            case "area":
                unit = "m²"
            case "land_area":
                unit = "ar"
            case "mileage":
                unit = "km"
            case "engine_capacity":
                unit = "l"
            case "power":
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
                "rooms": true, "floor": true, "total_floors": true, "area": true, 
                "land_area": true, "mileage": true, "year": true, "engine_capacity": true,
                "power": true, "screen_size": true,
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
        if attr.JSONValue != nil && len(attr.JSONValue) > 0 {
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
            v.unit
        FROM listing_attribute_values v
        JOIN category_attributes a ON v.attribute_id = a.id
        WHERE v.listing_id = $1
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
        ); err != nil {
            log.Printf("Error scanning attribute: %v", err)
            return nil, fmt.Errorf("error scanning listing attribute: %w", err)
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
                case "area":
                    unitStr = "m²"
                case "land_area":
                    unitStr = "ar"
                case "mileage":
                    unitStr = "km"
                case "engine_capacity":
                    unitStr = "l"
                case "power":
                    unitStr = "ks"
                case "screen_size":
                    unitStr = "inč"
                }
            }
            
            // Форматируем отображаемое значение с учетом типа
            if attr.AttributeName == "year" {
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
        if attrName == "year" || attrName == "rooms" || attrName == "floor" || attrName == "total_floors" {
            minValue = float64(int(minValue))
            maxValue = float64(int(maxValue))
        }
        
        // Для "year" добавляем запас +1 год для новых автомобилей
        if attrName == "year" && maxValue >= float64(time.Now().Year()-1) {
            maxValue = float64(time.Now().Year() + 1)
        }
        
        // Установка разумных шагов в зависимости от диапазона
        var step float64 = 1.0
        if attrName == "engine_capacity" {
            step = 0.1
        } else if attrName == "area" || attrName == "land_area" {
            step = 0.5
        }
        
        // Создаем информацию о границах
        ranges[attrName] = map[string]interface{}{
            "min": minValue,
            "max": maxValue,
            "step": step,
            "count": valueCount,
        }
        
        log.Printf("Attribute %s range: min=%.2f, max=%.2f, values=%d", 
                  attrName, minValue, maxValue, valueCount)
    }
    
    // Если нет данных, устанавливаем разумные значения по умолчанию
    defaultRanges := map[string]map[string]interface{}{
        "year": {"min": float64(time.Now().Year() - 30), "max": float64(time.Now().Year() + 1), "step": 1.0},
        "mileage": {"min": 0.0, "max": 500000.0, "step": 1000.0},
        "engine_capacity": {"min": 0.5, "max": 8.0, "step": 0.1},
        "power": {"min": 50.0, "max": 500.0, "step": 10.0},
        "rooms": {"min": 1.0, "max": 10.0, "step": 1.0},
        "floor": {"min": 1.0, "max": 25.0, "step": 1.0},
        "total_floors": {"min": 1.0, "max": 30.0, "step": 1.0},
        "area": {"min": 10.0, "max": 300.0, "step": 0.5},
        "land_area": {"min": 1.0, "max": 100.0, "step": 0.5},
    }
    
    // Заполняем отсутствующие атрибуты значениями по умолчанию
    for attr, defaultRange := range defaultRanges {
        if _, exists := ranges[attr]; !exists {
            ranges[attr] = defaultRange
            log.Printf("No data for attribute %s, using defaults: min=%.2f, max=%.2f",
                      attr, defaultRange["min"], defaultRange["max"])
        }
    }
    
    return ranges, nil
}

// Исправленная версия функции GetCategoryAttributes
func (s *Storage) GetCategoryAttributes(ctx context.Context, categoryID int) ([]models.CategoryAttribute, error) {
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
        a.attribute_type, 
        a.options, 
        a.validation_rules,
        a.is_searchable, 
        a.is_filterable, 
        COALESCE(m.is_required, a.is_required) as is_required,
        a.sort_order,
        a.created_at,
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
        var options, validRules sql.NullString
        var translationsJson, optionTranslationsJson []byte

        if err := rows.Scan(
            &attr.ID,
            &attr.Name,
            &attr.DisplayName,
            &attr.AttributeType,
            &options,
            &validRules,
            &attr.IsSearchable,
            &attr.IsFilterable,
            &attr.IsRequired,
            &attr.SortOrder,
            &attr.CreatedAt,
            &translationsJson,
            &optionTranslationsJson,
        ); err != nil {
            log.Printf("GetCategoryAttributes: Ошибка при сканировании результата: %v", err)
            return nil, fmt.Errorf("error scanning category attribute: %w", err)
        }

        // Обработка опциональных JSON полей
        if options.Valid {
            attr.Options = json.RawMessage(options.String)
        }
        if validRules.Valid {
            attr.ValidRules = json.RawMessage(validRules.String)
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
                    if attr.Options != nil && len(attr.Options) > 0 {
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

    log.Printf("GetCategoryAttributes: Успешно получено %d атрибутов для категории %d", len(attributes), categoryID)
    return attributes, nil
}

func (s *Storage) GetCategories(ctx context.Context) ([]models.MarketplaceCategory, error) {
	log.Printf("GetCategories: starting to fetch categories")
	query := `
        WITH category_translations AS (
            SELECT 
                t.entity_id,
                jsonb_object_agg(
                    t.language,  -- Изменено: убираем конкатенацию с field_name
                    t.translated_text
                ) as translations
            FROM translations t
            WHERE t.entity_type = 'category' 
            AND t.field_name = 'name'
            GROUP BY t.entity_id
        )
        SELECT 
            c.id, c.name, c.slug, c.parent_id, c.icon, c.created_at,
            COALESCE(ct.translations, '{}'::jsonb) as translations
        FROM marketplace_categories c
        LEFT JOIN category_translations ct ON c.id = ct.entity_id
    `

	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		log.Printf("Error querying categories: %v", err)
		return nil, err
	}
	defer rows.Close()

	var categories []models.MarketplaceCategory
	for rows.Next() {
		var cat models.MarketplaceCategory
		var translationsJson []byte

		err := rows.Scan(
			&cat.ID,
			&cat.Name,
			&cat.Slug,
			&cat.ParentID,
			&cat.Icon,
			&cat.CreatedAt,
			&translationsJson,
		)

		if err != nil {
			log.Printf("Error scanning category: %v", err)
			continue
		}

		//    log.Printf("Raw translations for category %d: %s", cat.ID, string(translationsJson))

		translations := make(map[string]string)
		if err := json.Unmarshal(translationsJson, &translations); err != nil {
			log.Printf("Error unmarshaling translations for category %d: %v", cat.ID, err)
		} else {
			cat.Translations = translations
			//     log.Printf("Processed translations for category %d: %+v", cat.ID, translations)
		}

		categories = append(categories, cat)
	}

	log.Printf("GetCategories: returning %d categories", len(categories))
	return categories, rows.Err()
}

func (s *Storage) GetCategoryByID(ctx context.Context, id int) (*models.MarketplaceCategory, error) {
	cat := &models.MarketplaceCategory{}
	err := s.pool.QueryRow(ctx, `
        SELECT 
            id, name, slug, parent_id, icon, created_at
        FROM marketplace_categories
        WHERE id = $1
    `, id).Scan(
		&cat.ID,
		&cat.Name,
		&cat.Slug,
		&cat.ParentID,
		&cat.Icon,
		&cat.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return cat, nil
}
func (s *Storage) GetListingByID(ctx context.Context, id int) (*models.MarketplaceListing, error) {
	listing := &models.MarketplaceListing{
		User:     &models.User{},
		Category: &models.MarketplaceCategory{},
	}

	// Получаем основные данные объявления с original_language
	err := s.pool.QueryRow(ctx, `
        SELECT 
            l.id, l.user_id, l.category_id, l.title, l.description,
            l.price, l.condition, l.status, l.location, l.latitude,
            l.longitude, l.address_city, l.address_country, l.views_count,
            l.created_at, l.updated_at, l.show_on_map, l.original_language,
            u.name, u.email, u.created_at as user_created_at, 
            COALESCE(u.picture_url, ''), u.phone,
            c.name as category_name, c.slug as category_slug, l.metadata
        FROM marketplace_listings l
        LEFT JOIN users u ON l.user_id = u.id
        LEFT JOIN marketplace_categories c ON l.category_id = c.id
        WHERE l.id = $1
    `, id).Scan(
		&listing.ID, &listing.UserID, &listing.CategoryID, &listing.Title,
		&listing.Description, &listing.Price, &listing.Condition, &listing.Status,
		&listing.Location, &listing.Latitude, &listing.Longitude, &listing.City,
		&listing.Country, &listing.ViewsCount, &listing.CreatedAt, &listing.UpdatedAt,
		&listing.ShowOnMap, &listing.OriginalLanguage,
		&listing.User.Name, &listing.User.Email, &listing.User.CreatedAt,
		&listing.User.PictureURL, &listing.User.Phone,
		&listing.Category.Name, &listing.Category.Slug, &listing.Metadata,
	)
	log.Printf("999 DEBUG: Listing %d metadata: %+v", id, listing.Metadata)

	if err != nil {
		return nil, fmt.Errorf("error getting listing: %w", err)
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
	images, err := s.GetListingImages(ctx, fmt.Sprintf("%d", listing.ID))
	if err != nil {
		log.Printf("Error loading images: %v", err) // Добавляем лог
		return nil, err
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
				if attr.AttributeName == "year" {
					attr.DisplayValue = fmt.Sprintf("%d", int(*attr.NumericValue))
				} else if attr.AttributeName == "engine_capacity" {
					attr.DisplayValue = fmt.Sprintf("%.1f л", *attr.NumericValue)
				} else if attr.AttributeName == "mileage" {
					attr.DisplayValue = fmt.Sprintf("%d км", int(*attr.NumericValue))
				} else if attr.AttributeName == "power" {
					attr.DisplayValue = fmt.Sprintf("%d л.с.", int(*attr.NumericValue))
				} else {
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

			if attr.JSONValue != nil && len(attr.JSONValue) > 0 {
				hasValue = true
				if attr.DisplayValue == "" {
					attr.DisplayValue = string(attr.JSONValue)
				}
			}

			// Если нет значения, но есть отображаемое значение
			if !hasValue && attr.DisplayValue != "" {
				// Попытка восстановить типизированное значение из отображаемого
				if attr.AttributeType == "text" || attr.AttributeType == "select" {
					strVal := attr.DisplayValue
					attr.TextValue = &strVal
				} else if attr.AttributeType == "number" {
					// Удаляем неожиданные символы (буквы, единицы измерения)
					clean := regexp.MustCompile(`[^\d\.-]`).ReplaceAllString(attr.DisplayValue, "")
					if numVal, err := strconv.ParseFloat(clean, 64); err == nil {
						attr.NumericValue = &numVal
					}
				} else if attr.AttributeType == "boolean" {
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
// backend/internal/storage/postgres/marketplace.go
package postgres

import (
	"backend/internal/domain/models"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	"log"
	"strings"
	//"time"
	// "github.com/jackc/pgx/v5"
)

func (s *Storage) CreateListing(ctx context.Context, listing *models.MarketplaceListing) (int, error) {
    var listingID int
    
    // Определяем язык, если не указан
    if listing.OriginalLanguage == "" {
        sourceLanguage, _, err := s.translationService.DetectLanguage(ctx, listing.Title + " " + listing.Description)
        if err != nil {
            // Если не удалось определить язык, используем английский по умолчанию
            listing.OriginalLanguage = "en"
        } else {
            listing.OriginalLanguage = sourceLanguage
        }
    }

    // Вставляем основные данные объявления
    err := s.pool.QueryRow(ctx, `
        INSERT INTO marketplace_listings (
            user_id, category_id, title, description, price,
            condition, status, location, latitude, longitude,
            address_city, address_country, show_on_map, original_language
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
        RETURNING id
    `,
        listing.UserID, listing.CategoryID, listing.Title, listing.Description,
        listing.Price, listing.Condition, listing.Status, listing.Location,
        listing.Latitude, listing.Longitude, listing.City, listing.Country,
        listing.ShowOnMap, listing.OriginalLanguage,
    ).Scan(&listingID)

    if err != nil {
        return 0, fmt.Errorf("failed to insert listing: %w", err)
    }

    // Сохраняем оригинальный текст как перевод для исходного языка
    err = s.pool.QueryRow(ctx, `
        INSERT INTO translations (
            entity_type, entity_id, language, field_name, 
            translated_text, is_machine_translated, is_verified
        ) VALUES ($1, $2, $3, $4, $5, $6, $7)
    `,
        "listing", listingID, listing.OriginalLanguage, "title",
        listing.Title, false, true,
    ).Scan()

    err = s.pool.QueryRow(ctx, `
        INSERT INTO translations (
            entity_type, entity_id, language, field_name, 
            translated_text, is_machine_translated, is_verified
        ) VALUES ($1, $2, $3, $4, $5, $6, $7)
    `,
        "listing", listingID, listing.OriginalLanguage, "description",
        listing.Description, false, true,
    ).Scan()

    // Переводим на другие языки
    targetLanguages := []string{"en", "ru", "sr"}
    for _, targetLang := range targetLanguages {
        if targetLang == listing.OriginalLanguage {
            continue
        }

        // Переводим заголовок
        translatedTitle, err := s.translationService.Translate(ctx, listing.Title, listing.OriginalLanguage, targetLang)
        if err == nil {
            s.pool.QueryRow(ctx, `
                INSERT INTO translations (
                    entity_type, entity_id, language, field_name, 
                    translated_text, is_machine_translated, is_verified
                ) VALUES ($1, $2, $3, $4, $5, $6, $7)
            `,
                "listing", listingID, targetLang, "title",
                translatedTitle, true, false,
            ).Scan()
        }

        // Переводим описание
        translatedDesc, err := s.translationService.Translate(ctx, listing.Description, listing.OriginalLanguage, targetLang)
        if err == nil {
            s.pool.QueryRow(ctx, `
                INSERT INTO translations (
                    entity_type, entity_id, language, field_name, 
                    translated_text, is_machine_translated, is_verified
                ) VALUES ($1, $2, $3, $4, $5, $6, $7)
            `,
                "listing", listingID, targetLang, "description",
                translatedDesc, true, false,
            ).Scan()
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
func (s *Storage) GetListings(ctx context.Context, filters map[string]string, limit int, offset int) ([]models.MarketplaceListing, int64, error) {
    userID, _ := ctx.Value("user_id").(int)
    if userID == 0 {
        userID = -1
    }

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
        l.original_language,
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
        userID,                // $2
    }

    conditions := []string{}
    argCount := 2

    // Добавляем остальные фильтры
    if v, ok := filters["query"]; ok && v != "" {
        argCount++
        conditions = append(conditions, fmt.Sprintf("AND (LOWER(l.title) LIKE LOWER($%d) OR LOWER(l.description) LIKE LOWER($%d))", 
            argCount, argCount))
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

	// Добавляем условия фильтрации
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
	query := `
        WITH RECURSIVE category_tree AS (
            -- Корневые категории
            SELECT 
                c.id,
                c.name,
                c.slug,
                c.icon,
                c.parent_id,
                c.created_at,
                0 as level,
                c.name::text as path  -- Добавляем явное приведение типа
            FROM marketplace_categories c
            WHERE parent_id IS NULL

            UNION ALL

            -- Дочерние категории
            SELECT 
                c.id,
                c.name,
                c.slug,
                c.icon,
                c.parent_id,
                c.created_at,
                ct.level + 1,
                (ct.path || ' > ' || c.name)::text -- Добавляем явное приведение типа
            FROM marketplace_categories c
            JOIN category_tree ct ON c.parent_id = ct.id
        )
        SELECT 
            id, name, slug, icon, parent_id, created_at, 
            level, path,
            (SELECT COUNT(*) FROM marketplace_listings WHERE category_id = ct.id) as listing_count
        FROM category_tree ct
        ORDER BY path;
    `

	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying category tree: %w", err)
	}
	defer rows.Close()

	var categories []models.CategoryTreeNode
	for rows.Next() {
		var cat models.CategoryTreeNode
		err := rows.Scan(
			&cat.ID, &cat.Name, &cat.Slug, &cat.Icon, &cat.ParentID,
			&cat.CreatedAt, &cat.Level, &cat.Path, &cat.ListingCount,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning category: %w", err)
		}
		categories = append(categories, cat)
	}

	return categories, rows.Err()
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
	result, err := s.pool.Exec(ctx, `
        DELETE FROM marketplace_listings
        WHERE id = $1 AND user_id = $2
    `, id, userID)

	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("listing not found or you don't have permission to delete it")
	}

	return nil
}

func (s *Storage) UpdateListing(ctx context.Context, listing *models.MarketplaceListing) error {
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

    return nil
}

func (s *Storage) GetCategories(ctx context.Context) ([]models.MarketplaceCategory, error) {
	rows, err := s.pool.Query(ctx, `
        SELECT 
            id, name, slug, parent_id, icon, created_at
        FROM marketplace_categories
        ORDER BY 
            CASE WHEN parent_id IS NULL THEN 0 ELSE 1 END,
            name
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.MarketplaceCategory
	for rows.Next() {
		var cat models.MarketplaceCategory
		err := rows.Scan(
			&cat.ID,
			&cat.Name,
			&cat.Slug,
			&cat.ParentID,
			&cat.Icon,
			&cat.CreatedAt,
		)
		if err != nil {
			continue
		}
		categories = append(categories, cat)
	}

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
            c.name as category_name, c.slug as category_slug
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
        &listing.Category.Name, &listing.Category.Slug,
    )

    if err != nil {
        return nil, fmt.Errorf("error getting listing: %w", err)
    }

    // Загружаем переводы
    translations := make(map[string]map[string]string)
    rows, err := s.pool.Query(ctx, `
        SELECT language, field_name, translated_text
        FROM translations
        WHERE entity_type = 'listing' AND entity_id = $1
    `, id)
    if err != nil {
        return listing, nil
    }
    defer rows.Close()

    for rows.Next() {
        var lang, field, text string
        if err := rows.Scan(&lang, &field, &text); err != nil {
            continue
        }
        if translations[lang] == nil {
            translations[lang] = make(map[string]string)
        }
        translations[lang][field] = text
        log.Printf("Added translation for %s.%s: %s", lang, field, text)
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
    log.Printf("Final translations for listing %d: %+v", id, translations)
    log.Printf("Original language: %s", listing.OriginalLanguage)

    return listing, nil
}
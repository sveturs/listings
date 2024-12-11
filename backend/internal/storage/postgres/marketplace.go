// backend/internal/storage/postgres/marketplace.go
package postgres

import (
	"backend/internal/domain/models"
	"context"
	"encoding/json"
	"fmt"

	"strings"
	"log"
)

func (db *Database) CreateListing(ctx context.Context, listing *models.MarketplaceListing) (int, error) {
	var listingID int
	err := db.pool.QueryRow(ctx, `
        INSERT INTO marketplace_listings (
            user_id, category_id, title, description, price,
            condition, status, location, latitude, longitude,
            address_city, address_country
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
        RETURNING id
    `,
		listing.UserID, listing.CategoryID, listing.Title, listing.Description,
		listing.Price, listing.Condition, listing.Status, listing.Location,
		listing.Latitude, listing.Longitude, listing.City, listing.Country,
	).Scan(&listingID)

	return listingID, err
}
func (db *Database) AddListingImage(ctx context.Context, image *models.MarketplaceImage) (int, error) {
	var imageID int
	err := db.pool.QueryRow(ctx, `
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

func (db *Database) GetListingImages(ctx context.Context, listingID string) ([]models.MarketplaceImage, error) {
	rows, err := db.pool.Query(ctx, `
        SELECT id, listing_id, file_path, file_name, 
               file_size, content_type, is_main, created_at
        FROM marketplace_images
        WHERE listing_id = $1
        ORDER BY is_main DESC, created_at DESC
    `, listingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []models.MarketplaceImage
	for rows.Next() {
		var img models.MarketplaceImage
		err := rows.Scan(
			&img.ID, &img.ListingID, &img.FilePath,
			&img.FileName, &img.FileSize, &img.ContentType,
			&img.IsMain, &img.CreatedAt,
		)
		if err != nil {
			continue
		}
		images = append(images, img)
	}

	return images, rows.Err()
}

func (db *Database) DeleteListingImage(ctx context.Context, imageID string) (string, error) {
	var filePath string
	err := db.pool.QueryRow(ctx,
		"SELECT file_path FROM marketplace_images WHERE id = $1",
		imageID,
	).Scan(&filePath)
	if err != nil {
		return "", err
	}

	_, err = db.pool.Exec(ctx,
		"DELETE FROM marketplace_images WHERE id = $1",
		imageID,
	)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func (db *Database) GetListings(ctx context.Context, filters map[string]string, limit int, offset int) ([]models.MarketplaceListing, int64, error) {
    // Сначала создаем CTE для получения ID всех подкатегорий
    baseQuery := `
        WITH RECURSIVE category_tree AS (
            -- Базовый случай: выбранная категория
            SELECT id, parent_id, name
            FROM marketplace_categories
            WHERE id = COALESCE($1::int, id)
            
            UNION ALL
            
            -- Рекурсивная часть: все дочерние категории
            SELECT c.id, c.parent_id, c.name
            FROM marketplace_categories c
            INNER JOIN category_tree ct ON c.parent_id = ct.id
        ),
        listing_data AS (
            SELECT 
                l.*,
                u.name as user_name, 
                u.email as user_email,
                c.name as category_name, 
                c.slug as category_slug,
                COALESCE(
                    (SELECT json_agg(
                        json_build_object(
                            'id', mi.id,
                            'listing_id', mi.listing_id,
                            'file_path', mi.file_path,
                            'file_name', mi.file_name,
                            'file_size', mi.file_size,
                            'content_type', mi.content_type,
                            'is_main', mi.is_main,
                            'created_at', to_char(mi.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
                        )
                    ) 
                    FROM marketplace_images mi 
                    WHERE mi.listing_id = l.id
                    ),
                    '[]'::json
                ) as images
            FROM marketplace_listings l
            JOIN users u ON l.user_id = u.id
            JOIN marketplace_categories c ON l.category_id = c.id
            WHERE 1=1
    `

    var conditions []string
    var args []interface{}
    argCount := 1 // Начинаем с 1 для category_id

    // Добавляем фильтр по категории
    if v := filters["category_id"]; v != "" {
        args = append(args, v)
        conditions = append(conditions, "AND l.category_id IN (SELECT id FROM category_tree)")
    } else {
        args = append(args, nil) // nil для COALESCE в CTE
    }

    // Добавляем условия фильтрации
    if v := filters["min_price"]; v != "" {
        argCount++
        conditions = append(conditions, fmt.Sprintf("AND l.price >= $%d", argCount))
        args = append(args, v)
    }
    if v := filters["max_price"]; v != "" {
        argCount++
        conditions = append(conditions, fmt.Sprintf("AND l.price <= $%d", argCount))
        args = append(args, v)
    }
    if v := filters["query"]; v != "" {
        argCount++
        conditions = append(conditions, fmt.Sprintf("AND (LOWER(l.title) LIKE LOWER($%d) OR LOWER(l.description) LIKE LOWER($%d))", argCount, argCount+1))
        args = append(args, "%"+v+"%", "%"+v+"%")
        argCount++
    }

    // Добавляем условия в запрос
    if len(conditions) > 0 {
        baseQuery += " " + strings.Join(conditions, " ")
    }

    // Закрываем CTE и добавляем основной запрос
    baseQuery += `)
        SELECT *, COUNT(*) OVER() as total_count 
        FROM listing_data 
        ORDER BY created_at DESC 
        LIMIT $` + fmt.Sprintf("%d", argCount+1) + ` OFFSET $` + fmt.Sprintf("%d", argCount+2)

    args = append(args, limit, offset)


    // Выполняем запрос
    rows, err := db.pool.Query(ctx, baseQuery, args...)
    if err != nil {
        return nil, 0, fmt.Errorf("error querying listings: %w", err)
    }
    defer rows.Close()

	var listings []models.MarketplaceListing
	var totalCount int64

	for rows.Next() {
		var listing models.MarketplaceListing
		// Инициализируем вложенные структуры
		listing.User = &models.User{}
		listing.Category = &models.MarketplaceCategory{}
		var imagesJSON []byte

		err := rows.Scan(
			&listing.ID, &listing.UserID, &listing.CategoryID, &listing.Title,
			&listing.Description, &listing.Price, &listing.Condition, &listing.Status,
			&listing.Location, &listing.Latitude, &listing.Longitude, &listing.City,
			&listing.Country, &listing.ViewsCount, &listing.CreatedAt, &listing.UpdatedAt,
			&listing.User.Name, &listing.User.Email,
			&listing.Category.Name, &listing.Category.Slug,
			&imagesJSON, &totalCount,
		)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, 0, fmt.Errorf("error scanning listing: %w", err)
		}

		if len(imagesJSON) > 0 {
			if err := json.Unmarshal(imagesJSON, &listing.Images); err != nil {
				log.Printf("Error parsing images JSON for listing %d: %v", listing.ID, err)
				log.Printf("Raw JSON: %s", string(imagesJSON))
				listing.Images = []models.MarketplaceImage{}
			} else {
				log.Printf("Successfully parsed %d images for listing %d", len(listing.Images), listing.ID)
			}
		} else {
			listing.Images = []models.MarketplaceImage{}
		}

		listings = append(listings, listing)
	}

	return listings, totalCount, nil
}
func (db *Database) GetCategoryTree(ctx context.Context) ([]models.CategoryTreeNode, error) {
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


    rows, err := db.pool.Query(ctx, query)
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
func (db *Database) AddToFavorites(ctx context.Context, userID int, listingID int) error {
	_, err := db.pool.Exec(ctx, `
        INSERT INTO marketplace_favorites (user_id, listing_id)
        VALUES ($1, $2)
        ON CONFLICT (user_id, listing_id) DO NOTHING
    `, userID, listingID)
	return err
}

func (db *Database) RemoveFromFavorites(ctx context.Context, userID int, listingID int) error {
	_, err := db.pool.Exec(ctx, `
        DELETE FROM marketplace_favorites
        WHERE user_id = $1 AND listing_id = $2
    `, userID, listingID)
	return err
}

func (db *Database) GetUserFavorites(ctx context.Context, userID int) ([]models.MarketplaceListing, error) {
	query := `
        SELECT 
            l.id, l.user_id, l.category_id, l.title, l.description,
            l.price, l.condition, l.status, l.location, l.latitude,
            l.longitude, l.address_city, l.address_country, l.views_count,
            l.created_at, l.updated_at
        FROM marketplace_listings l
        JOIN marketplace_favorites f ON l.id = f.listing_id
        WHERE f.user_id = $1
    `

	rows, err := db.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var listings []models.MarketplaceListing
	for rows.Next() {
		var l models.MarketplaceListing
		err := rows.Scan(
			&l.ID, &l.UserID, &l.CategoryID, &l.Title, &l.Description,
			&l.Price, &l.Condition, &l.Status, &l.Location, &l.Latitude,
			&l.Longitude, &l.City, &l.Country, &l.ViewsCount,
			&l.CreatedAt, &l.UpdatedAt,
		)
		if err != nil {
			continue
		}
		listings = append(listings, l)
	}

	return listings, rows.Err()
}
func (db *Database) DeleteListing(ctx context.Context, id int, userID int) error {
	result, err := db.pool.Exec(ctx, `
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

func (db *Database) UpdateListing(ctx context.Context, listing *models.MarketplaceListing) error {
	result, err := db.pool.Exec(ctx, `
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
            updated_at = CURRENT_TIMESTAMP
        WHERE id = $11 AND user_id = $12
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
		listing.ID,
		listing.UserID,
	)

	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("listing not found or you don't have permission to update it")
	}

	return nil
}
func (db *Database) GetCategories(ctx context.Context) ([]models.MarketplaceCategory, error) {
	rows, err := db.pool.Query(ctx, `
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

func (db *Database) GetCategoryByID(ctx context.Context, id int) (*models.MarketplaceCategory, error) {
	cat := &models.MarketplaceCategory{}
	err := db.pool.QueryRow(ctx, `
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
func (db *Database) GetListingByID(ctx context.Context, id int) (*models.MarketplaceListing, error) {
	listing := &models.MarketplaceListing{}

	// Получаем основные данные объявления вместе с информацией о пользователе и категории
	err := db.pool.QueryRow(ctx, `
        SELECT 
            l.id, l.user_id, l.category_id, l.title, l.description,
            l.price, l.condition, l.status, l.location, l.latitude,
            l.longitude, l.address_city, l.address_country, l.views_count,
            l.created_at, l.updated_at,
            u.name as user_name, u.email as user_email,
            c.name as category_name, c.slug as category_slug
        FROM marketplace_listings l
        JOIN users u ON l.user_id = u.id
        JOIN marketplace_categories c ON l.category_id = c.id
        WHERE l.id = $1
    `, id).Scan(
		&listing.ID, &listing.UserID, &listing.CategoryID, &listing.Title, &listing.Description,
		&listing.Price, &listing.Condition, &listing.Status, &listing.Location, &listing.Latitude,
		&listing.Longitude, &listing.City, &listing.Country, &listing.ViewsCount,
		&listing.CreatedAt, &listing.UpdatedAt,
		&listing.User.Name, &listing.User.Email,
		&listing.Category.Name, &listing.Category.Slug,
	)

	if err != nil {
		return nil, err
	}

	// Увеличиваем счетчик просмотров
	_, err = db.pool.Exec(ctx, `
        UPDATE marketplace_listings 
        SET views_count = views_count + 1 
        WHERE id = $1
    `, id)
	if err != nil {
		// Логируем ошибку, но не прерываем выполнение
		log.Printf("Error updating views count: %v", err)
	}

	// Получаем изображения для объявления
	images, err := db.GetListingImages(ctx, fmt.Sprintf("%d", id))
	if err != nil {
		// Логируем ошибку, но продолжаем без изображений
		log.Printf("Error getting images for listing %d: %v", id, err)
	} else {
		listing.Images = images
	}

	return listing, nil
}

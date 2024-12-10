package postgres

import (
    "backend/internal/domain/models"
    "context"
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
    // Базовый запрос с JOIN'ами
    query := `
        SELECT 
            l.id, l.user_id, l.category_id, l.title, l.description,
            l.price, l.condition, l.status, l.location, l.latitude,
            l.longitude, l.address_city, l.address_country, l.views_count,
            l.created_at, l.updated_at,
            u.name as user_name, u.email as user_email,
            c.name as category_name, c.slug as category_slug,
            COALESCE(i.image_count, 0) as image_count,
            COUNT(*) OVER() as total_count
        FROM marketplace_listings l
        JOIN users u ON l.user_id = u.id
        JOIN marketplace_categories c ON l.category_id = c.id
        LEFT JOIN (
            SELECT listing_id, COUNT(*) as image_count
            FROM marketplace_images
            GROUP BY listing_id
        ) i ON i.listing_id = l.id
        WHERE l.status = 'active'
    `

    var conditions []string
    var args []interface{}
    argCount := 0

    // Поиск по тексту (в заголовке и описании)
    if v := filters["query"]; v != "" {
        argCount++
        conditions = append(conditions, fmt.Sprintf(
            "(LOWER(l.title) LIKE LOWER($%d) OR LOWER(l.description) LIKE LOWER($%d))",
            argCount, argCount,
        ))
        args = append(args, "%"+v+"%")
    }

    // Фильтр по категории
    if v := filters["category_id"]; v != "" {
        argCount++
        conditions = append(conditions, fmt.Sprintf(
            "(l.category_id = $%d OR c.parent_id = $%d)", 
            argCount, argCount,
        ))
        args = append(args, v)
    }

    // Фильтр по диапазону цен
    if v := filters["min_price"]; v != "" {
        argCount++
        conditions = append(conditions, fmt.Sprintf("l.price >= $%d", argCount))
        args = append(args, v)
    }
    if v := filters["max_price"]; v != "" {
        argCount++
        conditions = append(conditions, fmt.Sprintf("l.price <= $%d", argCount))
        args = append(args, v)
    }

    // Фильтр по состоянию товара
    if v := filters["condition"]; v != "" {
        argCount++
        conditions = append(conditions, fmt.Sprintf("l.condition = $%d", argCount))
        args = append(args, v)
    }

    // Фильтр по городу
    if v := filters["city"]; v != "" {
        argCount++
        conditions = append(conditions, fmt.Sprintf("LOWER(l.address_city) LIKE LOWER($%d)", argCount))
        args = append(args, "%"+v+"%")
    }

    // Фильтр по стране
    if v := filters["country"]; v != "" {
        argCount++
        conditions = append(conditions, fmt.Sprintf("LOWER(l.address_country) LIKE LOWER($%d)", argCount))
        args = append(args, "%"+v+"%")
    }

    // Фильтр только с фото
    if v := filters["with_photos"]; v == "true" {
        conditions = append(conditions, "i.image_count > 0")
    }

    // Фильтр по дате публикации
    if v := filters["date_from"]; v != "" {
        argCount++
        conditions = append(conditions, fmt.Sprintf("l.created_at >= $%d", argCount))
        args = append(args, v)
    }
    if v := filters["date_to"]; v != "" {
        argCount++
        conditions = append(conditions, fmt.Sprintf("l.created_at <= $%d", argCount))
        args = append(args, v)
    }

    // Фильтр по радиусу от точки (если указаны координаты)
    if lat := filters["latitude"]; lat != "" {
        if lon := filters["longitude"]; lon != "" {
            if radius := filters["radius"]; radius != "" {
                argCount += 3
                // Используем формулу гаверсинусов для расчета расстояния
                conditions = append(conditions, fmt.Sprintf(`
                    (6371 * acos(cos(radians($%d)) * cos(radians(l.latitude)) * 
                    cos(radians(l.longitude) - radians($%d)) + 
                    sin(radians($%d)) * sin(radians(l.latitude)))) <= $%d`,
                    argCount-2, argCount-1, argCount-2, argCount))
                args = append(args, lat, lon, radius)
            }
        }
    }

    // Добавляем условия в запрос
    if len(conditions) > 0 {
        query += " AND " + strings.Join(conditions, " AND ")
    }

    // Сортировка
    switch filters["sort_by"] {
    case "price_asc":
        query += " ORDER BY l.price ASC"
    case "price_desc":
        query += " ORDER BY l.price DESC"
    case "date_asc":
        query += " ORDER BY l.created_at ASC"
    case "views":
        query += " ORDER BY l.views_count DESC"
    default:
        query += " ORDER BY l.created_at DESC"
    }

    // Добавляем пагинацию
    argCount += 2
    query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCount-1, argCount)
    args = append(args, limit, offset)

    // Выполняем запрос
    rows, err := db.pool.Query(ctx, query, args...)
    if err != nil {
        return nil, 0, fmt.Errorf("error querying listings: %w", err)
    }
    defer rows.Close()

    var listings []models.MarketplaceListing
    var total int64

    for rows.Next() {
        var l models.MarketplaceListing
        var u models.User
        var c models.MarketplaceCategory
        var imageCount int

        err := rows.Scan(
            &l.ID, &l.UserID, &l.CategoryID, &l.Title, &l.Description,
            &l.Price, &l.Condition, &l.Status, &l.Location, &l.Latitude,
            &l.Longitude, &l.City, &l.Country, &l.ViewsCount,
            &l.CreatedAt, &l.UpdatedAt,
            &u.Name, &u.Email,
            &c.Name, &c.Slug,
            &imageCount,
            &total,
        )
        if err != nil {
            return nil, 0, fmt.Errorf("error scanning listing: %w", err)
        }

        l.User = &u
        l.Category = &c

        // Получаем изображения для объявления
        if imageCount > 0 {
            images, err := db.GetListingImages(ctx, fmt.Sprintf("%d", l.ID))
            if err != nil {
                log.Printf("Error getting images for listing %d: %v", l.ID, err)
            } else {
                l.Images = images
            }
        }

        listings = append(listings, l)
    }

    return listings, total, rows.Err()
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
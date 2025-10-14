// backend/internal/proj/c2c/storage/postgres/listings_favorites.go
package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"backend/internal/domain/models"
)

// AddToFavorites добавляет C2C листинг в избранное пользователя
func (s *Storage) AddToFavorites(ctx context.Context, userID int, listingID int) error {
	_, err := s.pool.Exec(ctx, `
        INSERT INTO c2c_favorites (user_id, listing_id)
        VALUES ($1, $2)
        ON CONFLICT (user_id, listing_id) DO NOTHING
    `, userID, listingID)
	return err
}

// RemoveFromFavorites удаляет C2C листинг из избранного пользователя
func (s *Storage) RemoveFromFavorites(ctx context.Context, userID int, listingID int) error {
	_, err := s.pool.Exec(ctx, `
        DELETE FROM c2c_favorites
        WHERE user_id = $1 AND listing_id = $2
    `, userID, listingID)
	return err
}

// AddStorefrontToFavorites добавляет товар витрины в избранное
func (s *Storage) AddStorefrontToFavorites(ctx context.Context, userID int, productID int) error {
	_, err := s.pool.Exec(ctx, `
        INSERT INTO b2c_favorites (user_id, product_id)
        VALUES ($1, $2)
        ON CONFLICT (user_id, product_id) DO NOTHING
    `, userID, productID)
	return err
}

// RemoveStorefrontFromFavorites удаляет товар витрины из избранного
func (s *Storage) RemoveStorefrontFromFavorites(ctx context.Context, userID int, productID int) error {
	_, err := s.pool.Exec(ctx, `
        DELETE FROM b2c_favorites
        WHERE user_id = $1 AND product_id = $2
    `, userID, productID)
	return err
}

// GetUserStorefrontFavorites получает избранные товары витрин пользователя
func (s *Storage) GetUserStorefrontFavorites(ctx context.Context, userID int) ([]models.MarketplaceListing, error) {
	query := `
        WITH product_images AS (
            SELECT
                storefront_product_id as product_id,
                jsonb_agg(
                    jsonb_build_object(
                        'id', id,
                        'public_url', image_url,
                        'image_url', image_url,
                        'is_main', is_default,
                        'display_order', display_order
                    ) ORDER BY is_default DESC, display_order ASC, id ASC
                ) as images
            FROM b2c_product_images
            GROUP BY storefront_product_id
        )
        SELECT
            p.id,
            p.storefront_id as user_id,
            p.category_id,
            p.name as title,
            p.description,
            p.price,
            'new' as condition,
            'active' as status,
            COALESCE(s.address, '') as location,
            s.latitude,
            s.longitude,
            COALESCE(s.city, '') as city,
            COALESCE(s.country, 'Serbia') as country,
            p.view_count as views_count,
            p.created_at,
            p.updated_at,
            s.name as store_name,
            s.email as store_email,
            s.created_at as store_created_at,
            COALESCE(s.logo_url, ''),
            COALESCE(c.name, '') as category_name,
            COALESCE(c.slug, '') as category_slug,
            true as is_favorite,
            COALESCE(product_images.images, '[]'::jsonb) as product_images
        FROM b2c_products p
        JOIN b2c_favorites f ON p.id = f.product_id
        LEFT JOIN b2c_stores s ON p.storefront_id = s.id
        LEFT JOIN c2c_categories c ON p.category_id = c.id
        LEFT JOIN product_images ON p.id = product_images.product_id
        WHERE f.user_id = $1
        ORDER BY f.created_at DESC
    `

	rows, err := s.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying storefront favorites: %w", err)
	}
	defer rows.Close()

	var listings []models.MarketplaceListing

	for rows.Next() {
		var listing models.MarketplaceListing
		var imagesJSON []byte
		var userPictureURL sql.NullString
		var storeName, storeEmail sql.NullString
		var storeCreatedAt sql.NullTime

		// Инициализируем Category и User чтобы избежать nil pointer
		listing.Category = &models.MarketplaceCategory{}
		listing.User = &models.User{}

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
			&storeName,
			&storeEmail,
			&storeCreatedAt,
			&userPictureURL,
			&listing.Category.Name,
			&listing.Category.Slug,
			&listing.IsFavorite,
			&imagesJSON,
		)
		if err != nil {
			log.Printf("Error scanning storefront favorite row: %v", err)
			continue
		}

		// Устанавливаем флаг что это товар витрины
		listing.IsStorefrontProduct = true

		// Используем имя магазина вместо имени пользователя
		if storeName.Valid {
			listing.User.Name = storeName.String
		}
		if storeEmail.Valid {
			listing.User.Email = storeEmail.String
		}
		if userPictureURL.Valid {
			listing.User.PictureURL = userPictureURL.String
		}
		listing.User.ID = listing.UserID

		// Парсим изображения из JSON
		var images []models.MarketplaceImage
		if err := json.Unmarshal(imagesJSON, &images); err != nil {
			log.Printf("Error unmarshalling images for storefront product %d: %v", listing.ID, err)
			images = []models.MarketplaceImage{}
		}
		// Преобразуем относительные URL в полные
		for i := range images {
			images[i].PublicURL = buildFullImageURL(images[i].PublicURL)
			images[i].ImageURL = buildFullImageURL(images[i].ImageURL)
		}
		listing.Images = images

		listings = append(listings, listing)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating storefront favorites rows: %w", err)
	}

	return listings, nil
}

// GetUserFavorites получает избранные C2C листинги пользователя
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
                        'created_at', to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS.US"Z"')
                    ) ORDER BY is_main DESC, id ASC
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
            COALESCE(c.name, '') as category_name,
            COALESCE(c.slug, '') as category_slug,
            true as is_favorite,
            COALESCE(li.images, '[]'::jsonb) as listing_images
        FROM c2c_listings l
        JOIN c2c_favorites f ON l.id = f.listing_id
        LEFT JOIN c2c_categories c ON l.category_id = c.id
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
			&listing.Category.Name,
			&listing.Category.Slug,
			&listing.IsFavorite,
			&imagesJSON,
		)
		if err != nil {
			log.Printf("Error scanning listing: %v", err)
			continue
		}

		// User info будет загружена в handler через auth-service
		listing.User.ID = listing.UserID

		// Парсим изображения из JSON
		var images []models.MarketplaceImage
		if err := json.Unmarshal(imagesJSON, &images); err != nil {
			log.Printf("Error unmarshalling images for listing %d: %v", listing.ID, err)
			images = []models.MarketplaceImage{}
		}
		// Преобразуем относительные URL в полные
		for i := range images {
			images[i].PublicURL = buildFullImageURL(images[i].PublicURL)
			images[i].ImageURL = buildFullImageURL(images[i].ImageURL)
		}
		listing.Images = images

		listings = append(listings, listing)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return listings, nil
}

// GetFavoritedUsers получает список пользователей, добавивших листинг в избранное
func (s *Storage) GetFavoritedUsers(ctx context.Context, listingID int) ([]int, error) {
	query := `
        SELECT user_id
        FROM c2c_favorites
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

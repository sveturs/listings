// backend/internal/proj/c2c/storage/postgres/listings_images.go
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"backend/internal/domain/models"
)

// AddListingImage добавляет изображение к листингу
func (s *Storage) AddListingImage(ctx context.Context, image *models.MarketplaceImage) (int, error) {
	var id int
	err := s.pool.QueryRow(ctx, `
        INSERT INTO c2c_images
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

// GetListingImages получает все изображения листинга
func (s *Storage) GetListingImages(ctx context.Context, listingID string) ([]models.MarketplaceImage, error) {
	query := `
        SELECT
            id, listing_id, file_path, file_name, file_size,
            content_type, is_main, created_at,
            storage_type, storage_bucket, public_url
        FROM c2c_images
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
		var storageBucket sql.NullString
		var publicURL sql.NullString

		err := rows.Scan(
			&image.ID, &image.ListingID, &image.FilePath, &image.FileName,
			&image.FileSize, &image.ContentType, &image.IsMain, &image.CreatedAt,
			&image.StorageType, &storageBucket, &publicURL,
		)
		if err != nil {
			return nil, err
		}

		// Обработка NULL значений
		if storageBucket.Valid {
			image.StorageBucket = storageBucket.String
		}
		if publicURL.Valid {
			// Преобразуем относительный URL в полный
			fullURL := buildFullImageURL(publicURL.String)
			image.PublicURL = fullURL
			// Заполняем ImageURL из PublicURL для API
			image.ImageURL = fullURL
		}

		images = append(images, image)
	}

	return images, nil
}

// DeleteListingImage удаляет изображение и возвращает путь к файлу
func (s *Storage) DeleteListingImage(ctx context.Context, imageID string) (string, error) {
	var filePath string
	err := s.pool.QueryRow(ctx,
		"SELECT file_path FROM c2c_images WHERE id = $1",
		imageID,
	).Scan(&filePath)
	if err != nil {
		return "", err
	}

	_, err = s.pool.Exec(ctx,
		"DELETE FROM c2c_images WHERE id = $1",
		imageID,
	)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

// GetB2CProductImages загружает изображения для B2C товара и конвертирует их в MarketplaceImage
func (s *Storage) GetB2CProductImages(ctx context.Context, productID int) ([]models.MarketplaceImage, error) {
	query := `
        SELECT
            id, storefront_product_id, image_url, thumbnail_url,
            display_order, is_default, created_at
        FROM b2c_product_images
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
			ID:           img.ID,
			ListingID:    img.StorefrontProductID,
			PublicURL:    img.ImageURL,
			ImageURL:     img.ImageURL,     // Заполняем ImageURL для API
			ThumbnailURL: img.ThumbnailURL, // Добавляем ThumbnailURL
			IsMain:       img.IsDefault,
			DisplayOrder: img.DisplayOrder, // Добавляем DisplayOrder
			StorageType:  "minio",
			CreatedAt:    img.CreatedAt,
		}

		images = append(images, marketplaceImage)
	}

	return images, nil
}

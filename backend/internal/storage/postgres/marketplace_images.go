// backend/internal/storage/postgres/marketplace_images.go
package postgres

import (
	"context"
	"errors"
	"fmt"

	"backend/internal/domain/models"
)

// AddListingImage добавляет изображение к листингу
func (db *Database) AddListingImage(ctx context.Context, image *models.MarketplaceImage) (int, error) {
	result, err := db.grpcClient.AddListingImage(ctx, image)
	if err != nil {
		return 0, err
	}
	return result.ID, nil
}

// GetListingImages возвращает все изображения для листинга
func (db *Database) GetListingImages(ctx context.Context, listingID string) ([]models.MarketplaceImage, error) {
	// Convert string listingID to int64
	var listingIDInt64 int64
	_, err := fmt.Sscanf(listingID, "%d", &listingIDInt64)
	if err != nil {
		return nil, fmt.Errorf("invalid listing ID format: %s", listingID)
	}

	imagesPtrs, err := db.grpcClient.GetListingImages(ctx, listingIDInt64)
	if err != nil {
		return nil, err
	}

	// Convert []*models.MarketplaceImage to []models.MarketplaceImage
	images := make([]models.MarketplaceImage, len(imagesPtrs))
	for i, img := range imagesPtrs {
		images[i] = *img
	}

	return images, nil
}

// UpdateImageMainStatus обновляет статус is_main для изображения
func (db *Database) UpdateImageMainStatus(ctx context.Context, imageID int, isMain bool) error {
	query := `
		UPDATE c2c_images
		SET is_main = $1
		WHERE id = $2
	`

	result, err := db.pool.Exec(ctx, query, isMain, imageID)
	if err != nil {
		return fmt.Errorf("failed to update image main status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("image not found")
	}

	return nil
}

// SetMainImage устанавливает основное изображение для листинга (сбрасывает is_main у других)
func (db *Database) SetMainImage(ctx context.Context, listingID int, imageID int) error {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx) // Rollback is safe to call even if transaction already committed
	}()

	// 1. Сначала сбрасываем is_main у всех изображений этого листинга
	_, err = tx.Exec(ctx, `
		UPDATE c2c_images
		SET is_main = false
		WHERE listing_id = $1
	`, listingID)
	if err != nil {
		return fmt.Errorf("failed to reset main images: %w", err)
	}

	// 2. Устанавливаем is_main для выбранного изображения
	result, err := tx.Exec(ctx, `
		UPDATE c2c_images
		SET is_main = true
		WHERE id = $1 AND listing_id = $2
	`, imageID, listingID)
	if err != nil {
		return fmt.Errorf("failed to set main image: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("image not found or belongs to different listing")
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetListingImagesCount возвращает количество изображений для листинга
func (db *Database) GetListingImagesCount(ctx context.Context, listingID int) (int, error) {
	var count int
	err := db.pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM c2c_images
		WHERE listing_id = $1
	`, listingID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count listing images: %w", err)
	}

	return count, nil
}

package postgres

import (
	"context"
	"fmt"

	"backend/internal/domain/models"
	"backend/internal/storage/interfaces"

	"github.com/jmoiron/sqlx"
)

// ImageRepository - репозиторий для работы с изображениями
type ImageRepository struct {
	db *sqlx.DB
}

// NewImageRepository создает новый ImageRepository
func NewImageRepository(db *sqlx.DB) interfaces.ImageRepositoryInterface {
	return &ImageRepository{db: db}
}

// CreateImage создает новое изображение
func (r *ImageRepository) CreateImage(ctx context.Context, image models.ImageInterface) (models.ImageInterface, error) {
	// Определяем тип изображения и выбираем соответствующий запрос
	switch img := image.(type) {
	case *models.StorefrontProductImage:
		return r.createStorefrontProductImage(ctx, img)
	case *models.MarketplaceImage:
		return r.createMarketplaceImage(ctx, img)
	default:
		return nil, fmt.Errorf("unsupported image type: %T", image)
	}
}

// GetImageByID получает изображение по ID
func (r *ImageRepository) GetImageByID(ctx context.Context, imageID int) (models.ImageInterface, error) {
	// Сначала пытаемся найти в storefront_product_images
	storefrontImage, err := r.getStorefrontProductImageByID(ctx, imageID)
	if err == nil {
		return storefrontImage, nil
	}

	// Если не найдено, пытаемся найти в marketplace_images
	marketplaceImage, err := r.getMarketplaceImageByID(ctx, imageID)
	if err == nil {
		return marketplaceImage, nil
	}

	return nil, fmt.Errorf("image not found with ID: %d", imageID)
}

// GetImagesByEntity получает все изображения для сущности
func (r *ImageRepository) GetImagesByEntity(ctx context.Context, entityType string, entityID int) ([]models.ImageInterface, error) {
	switch entityType {
	case "storefront_product":
		return r.getStorefrontProductImages(ctx, entityID)
	case "marketplace_listing":
		return r.getMarketplaceImages(ctx, entityID)
	default:
		return nil, fmt.Errorf("unsupported entity type: %s", entityType)
	}
}

// DeleteImage удаляет изображение
func (r *ImageRepository) DeleteImage(ctx context.Context, imageID int) error {
	// Пытаемся удалить из storefront_product_images
	err := r.deleteStorefrontProductImage(ctx, imageID)
	if err == nil {
		return nil
	}

	// Пытаемся удалить из marketplace_images
	err = r.deleteMarketplaceImage(ctx, imageID)
	if err == nil {
		return nil
	}

	return fmt.Errorf("failed to delete image with ID: %d", imageID)
}

// UnsetMainImages сбрасывает флаг главного изображения для всех изображений сущности
func (r *ImageRepository) UnsetMainImages(ctx context.Context, entityType string, entityID int) error {
	switch entityType {
	case "storefront_product":
		query := `UPDATE b2c_product_images SET is_default = false WHERE storefront_product_id = $1`
		_, err := r.db.ExecContext(ctx, query, entityID)
		return err
	case "marketplace_listing":
		query := `UPDATE c2c_images SET is_main = false WHERE listing_id = $1`
		_, err := r.db.ExecContext(ctx, query, entityID)
		return err
	default:
		return fmt.Errorf("unsupported entity type: %s", entityType)
	}
}

// SetMainImage устанавливает изображение как главное
func (r *ImageRepository) SetMainImage(ctx context.Context, imageID int, isMain bool) error {
	// Пытаемся обновить в b2c_product_images
	query := `UPDATE b2c_product_images SET is_default = $1 WHERE id = $2`
	result, err := r.db.ExecContext(ctx, query, isMain, imageID)
	if err == nil {
		rowsAffected, _ := result.RowsAffected()
		if rowsAffected > 0 {
			return nil
		}
	}

	// Пытаемся обновить в c2c_images
	query = `UPDATE c2c_images SET is_main = $1 WHERE id = $2`
	result, err = r.db.ExecContext(ctx, query, isMain, imageID)
	if err == nil {
		rowsAffected, _ := result.RowsAffected()
		if rowsAffected > 0 {
			return nil
		}
	}

	return fmt.Errorf("failed to set main image for ID: %d", imageID)
}

// UpdateDisplayOrder обновляет порядок отображения изображений
func (r *ImageRepository) UpdateDisplayOrder(ctx context.Context, imageID int, displayOrder int) error {
	// Пытаемся обновить в b2c_product_images
	query := `UPDATE b2c_product_images SET display_order = $1 WHERE id = $2`
	result, err := r.db.ExecContext(ctx, query, displayOrder, imageID)
	if err == nil {
		rowsAffected, _ := result.RowsAffected()
		if rowsAffected > 0 {
			return nil
		}
	}

	// Пытаемся обновить в c2c_images
	query = `UPDATE c2c_images SET display_order = $1 WHERE id = $2`
	result, err = r.db.ExecContext(ctx, query, displayOrder, imageID)
	if err == nil {
		rowsAffected, _ := result.RowsAffected()
		if rowsAffected > 0 {
			return nil
		}
	}

	return fmt.Errorf("failed to update display order for image ID: %d", imageID)
}

// Методы для работы с b2c_product_images
func (r *ImageRepository) createStorefrontProductImage(ctx context.Context, img *models.StorefrontProductImage) (models.ImageInterface, error) {
	query := `
		INSERT INTO b2c_product_images (storefront_product_id, image_url, thumbnail_url, display_order, is_default)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`

	err := r.db.QueryRowContext(ctx, query, img.StorefrontProductID, img.ImageURL, img.ThumbnailURL, img.DisplayOrder, img.IsDefault).
		Scan(&img.ID, &img.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create storefront product image: %w", err)
	}

	return img, nil
}

func (r *ImageRepository) getStorefrontProductImageByID(ctx context.Context, imageID int) (*models.StorefrontProductImage, error) {
	query := `
		SELECT id, storefront_product_id, image_url, thumbnail_url, display_order, is_default, created_at
		FROM b2c_product_images
		WHERE id = $1`

	var img models.StorefrontProductImage
	err := r.db.GetContext(ctx, &img, query, imageID)
	if err != nil {
		return nil, fmt.Errorf("failed to get storefront product image: %w", err)
	}

	return &img, nil
}

func (r *ImageRepository) getStorefrontProductImages(ctx context.Context, productID int) ([]models.ImageInterface, error) {
	query := `
		SELECT id, storefront_product_id, image_url, thumbnail_url, display_order, is_default, created_at
		FROM b2c_product_images
		WHERE storefront_product_id = $1
		ORDER BY is_default DESC, display_order ASC, created_at ASC`

	var images []models.StorefrontProductImage
	err := r.db.SelectContext(ctx, &images, query, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get storefront product images: %w", err)
	}

	// Конвертируем в интерфейс
	result := make([]models.ImageInterface, len(images))
	for i := range images {
		result[i] = &images[i]
	}

	return result, nil
}

func (r *ImageRepository) deleteStorefrontProductImage(ctx context.Context, imageID int) error {
	query := `DELETE FROM b2c_product_images WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, imageID)
	if err != nil {
		return fmt.Errorf("failed to delete storefront product image: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no storefront product image found with ID: %d", imageID)
	}

	return nil
}

// Методы для работы с marketplace_images
func (r *ImageRepository) createMarketplaceImage(ctx context.Context, img *models.MarketplaceImage) (models.ImageInterface, error) {
	query := `
		INSERT INTO c2c_images (
			id, listing_id, file_path, file_name, file_size, content_type,
			is_main, storage_type, storage_bucket, public_url
		)
		VALUES (nextval('c2c_images_id_seq'), $1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at`

	err := r.db.QueryRowContext(
		ctx, query,
		img.ListingID,
		img.FilePath,
		img.FileName,
		img.FileSize,
		img.ContentType,
		img.IsMain,
		img.StorageType,
		img.StorageBucket,
		img.PublicURL,
	).Scan(&img.ID, &img.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create marketplace image: %w", err)
	}

	// Устанавливаем URL'ы из PublicURL
	if img.ImageURL == "" {
		img.ImageURL = img.PublicURL
	}

	return img, nil
}

func (r *ImageRepository) getMarketplaceImageByID(ctx context.Context, imageID int) (*models.MarketplaceImage, error) {
	query := `
		SELECT id, listing_id, file_path, file_name, file_size, content_type,
		       is_main, storage_type, storage_bucket, public_url, created_at
		FROM c2c_images
		WHERE id = $1`

	var img models.MarketplaceImage
	err := r.db.QueryRowContext(ctx, query, imageID).Scan(
		&img.ID,
		&img.ListingID,
		&img.FilePath,
		&img.FileName,
		&img.FileSize,
		&img.ContentType,
		&img.IsMain,
		&img.StorageType,
		&img.StorageBucket,
		&img.PublicURL,
		&img.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get marketplace image: %w", err)
	}

	// Устанавливаем URL'ы
	img.ImageURL = img.PublicURL
	// DisplayOrder не хранится в таблице c2c_images
	img.DisplayOrder = 0

	return &img, nil
}

func (r *ImageRepository) getMarketplaceImages(ctx context.Context, listingID int) ([]models.ImageInterface, error) {
	query := `
		SELECT id, listing_id, file_path, file_name, file_size, content_type,
		       is_main, storage_type, storage_bucket, public_url, created_at
		FROM c2c_images
		WHERE listing_id = $1
		ORDER BY is_main DESC, created_at ASC`

	rows, err := r.db.QueryContext(ctx, query, listingID)
	if err != nil {
		return nil, fmt.Errorf("failed to query marketplace images: %w", err)
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
			&img.StorageType,
			&img.StorageBucket,
			&img.PublicURL,
			&img.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan marketplace image: %w", err)
		}

		// Устанавливаем URL'ы
		img.ImageURL = img.PublicURL
		// DisplayOrder не хранится в таблице c2c_images
		img.DisplayOrder = 0

		images = append(images, img)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	// Конвертируем в интерфейс
	result := make([]models.ImageInterface, len(images))
	for i := range images {
		result[i] = &images[i]
	}

	return result, nil
}

func (r *ImageRepository) deleteMarketplaceImage(ctx context.Context, imageID int) error {
	query := `DELETE FROM c2c_images WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, imageID)
	if err != nil {
		return fmt.Errorf("failed to delete marketplace image: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no marketplace image found with ID: %d", imageID)
	}

	return nil
}

// GetFirstImageByListingID получает первое изображение для listing по дате создания
// Используется для автоматического назначения нового главного изображения
func (r *ImageRepository) GetFirstImageByListingID(ctx context.Context, listingID int) (*models.MarketplaceImage, error) {
	query := `
		SELECT id, listing_id, file_path, file_name, file_size, content_type,
		       is_main, storage_type, storage_bucket, public_url, created_at
		FROM c2c_images
		WHERE listing_id = $1
		ORDER BY created_at ASC
		LIMIT 1`

	var img models.MarketplaceImage
	err := r.db.QueryRowContext(ctx, query, listingID).Scan(
		&img.ID,
		&img.ListingID,
		&img.FilePath,
		&img.FileName,
		&img.FileSize,
		&img.ContentType,
		&img.IsMain,
		&img.StorageType,
		&img.StorageBucket,
		&img.PublicURL,
		&img.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get first image: %w", err)
	}

	// Устанавливаем URL'ы
	img.ImageURL = img.PublicURL
	img.DisplayOrder = 0

	return &img, nil
}

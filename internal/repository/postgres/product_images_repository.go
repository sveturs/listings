package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/sveturs/listings/internal/domain"
)

// GetProductImageByID retrieves a single product image by ID
// Uses listing_images table (unified for both C2C and B2C products)
func (r *Repository) GetProductImageByID(ctx context.Context, imageID int64) (*domain.ProductImage, error) {
	query := `
		SELECT id, listing_id, url, storage_path,
		       thumbnail_url, display_order, is_primary,
		       width, height, file_size,
		       mime_type, created_at, updated_at
		FROM listing_images
		WHERE id = $1
	`

	var image domain.ProductImage
	var productID sql.NullInt64
	var width, height sql.NullInt32
	var fileSize sql.NullInt64
	var mimeType, storagePath, thumbnailURL sql.NullString

	err := r.db.QueryRowxContext(ctx, query, imageID).Scan(
		&image.ID,
		&productID,
		&image.URL,
		&storagePath,
		&thumbnailURL,
		&image.DisplayOrder,
		&image.IsPrimary,
		&width,
		&height,
		&fileSize,
		&mimeType,
		&image.CreatedAt,
		&image.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product image not found: %w", err)
		}
		r.logger.Error().Err(err).Int64("image_id", imageID).Msg("failed to get product image")
		return nil, fmt.Errorf("failed to get product image: %w", err)
	}

	// Handle nullable fields
	if productID.Valid {
		image.ProductID = &productID.Int64
	}
	if storagePath.Valid {
		image.StoragePath = &storagePath.String
	}
	if thumbnailURL.Valid {
		image.ThumbnailURL = &thumbnailURL.String
	}
	if width.Valid {
		image.Width = &width.Int32
	}
	if height.Valid {
		image.Height = &height.Int32
	}
	if fileSize.Valid {
		image.FileSize = &fileSize.Int64
	}
	if mimeType.Valid {
		image.MimeType = &mimeType.String
	}

	return &image, nil
}

// AddProductImage adds a new image to a product (B2C or C2C)
// Uses listing_images table (unified for both C2C and B2C products)
func (r *Repository) AddProductImage(ctx context.Context, image *domain.ProductImage) (*domain.ProductImage, error) {
	query := `
		INSERT INTO listing_images (
			listing_id, url, storage_path, thumbnail_url, display_order,
			is_primary, width, height, file_size, mime_type
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`

	var newImage domain.ProductImage
	err := r.db.QueryRowxContext(ctx, query,
		image.ProductID,
		image.URL,
		image.StoragePath,
		image.ThumbnailURL,
		image.DisplayOrder,
		image.IsPrimary,
		image.Width,
		image.Height,
		image.FileSize,
		image.MimeType,
	).Scan(&newImage.ID, &newImage.CreatedAt, &newImage.UpdatedAt)

	if err != nil {
		r.logger.Error().Err(err).Msg("failed to add product image")
		return nil, fmt.Errorf("failed to add product image: %w", err)
	}

	// Populate the returned image with input data
	newImage.ProductID = image.ProductID
	newImage.URL = image.URL
	newImage.StoragePath = image.StoragePath
	newImage.ThumbnailURL = image.ThumbnailURL
	newImage.DisplayOrder = image.DisplayOrder
	newImage.IsPrimary = image.IsPrimary
	newImage.Width = image.Width
	newImage.Height = image.Height
	newImage.FileSize = image.FileSize
	newImage.MimeType = image.MimeType
	newImage.UpdatedAt = newImage.CreatedAt

	r.logger.Info().
		Int64("image_id", newImage.ID).
		Msg("product image added")
	return &newImage, nil
}

// GetProductImages retrieves all images for a product (B2C or C2C)
// Uses listing_images table (unified for both C2C and B2C products)
func (r *Repository) GetProductImages(ctx context.Context, productID int64) ([]*domain.ProductImage, error) {
	query := `
		SELECT id, listing_id, url, storage_path,
		       thumbnail_url, display_order, is_primary,
		       width, height, file_size,
		       mime_type, created_at, updated_at
		FROM listing_images
		WHERE listing_id = $1
		ORDER BY is_primary DESC, display_order ASC, id ASC
	`

	rows, err := r.db.QueryxContext(ctx, query, productID)
	if err != nil {
		r.logger.Error().Err(err).Int64("product_id", productID).Msg("failed to get product images")
		return nil, fmt.Errorf("failed to get product images: %w", err)
	}
	defer rows.Close()

	var images []*domain.ProductImage
	for rows.Next() {
		var image domain.ProductImage
		var prodID sql.NullInt64
		var width, height sql.NullInt32
		var fileSize sql.NullInt64
		var mimeType, storagePath, thumbnailURL sql.NullString

		err := rows.Scan(
			&image.ID,
			&prodID,
			&image.URL,
			&storagePath,
			&thumbnailURL,
			&image.DisplayOrder,
			&image.IsPrimary,
			&width,
			&height,
			&fileSize,
			&mimeType,
			&image.CreatedAt,
			&image.UpdatedAt,
		)
		if err != nil {
			r.logger.Error().Err(err).Msg("failed to scan product image")
			return nil, fmt.Errorf("failed to scan product image: %w", err)
		}

		// Handle nullable fields
		if prodID.Valid {
			image.ProductID = &prodID.Int64
		}
		if storagePath.Valid {
			image.StoragePath = &storagePath.String
		}
		if thumbnailURL.Valid {
			image.ThumbnailURL = &thumbnailURL.String
		}
		if width.Valid {
			image.Width = &width.Int32
		}
		if height.Valid {
			image.Height = &height.Int32
		}
		if fileSize.Valid {
			image.FileSize = &fileSize.Int64
		}
		if mimeType.Valid {
			image.MimeType = &mimeType.String
		}

		images = append(images, &image)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating product images: %w", err)
	}

	return images, nil
}

// GetProductImagesBatch retrieves images for multiple products at once
// Returns a map of productID -> images slice
func (r *Repository) GetProductImagesBatch(ctx context.Context, productIDs []int64) (map[int64][]*domain.ProductImage, error) {
	if len(productIDs) == 0 {
		return make(map[int64][]*domain.ProductImage), nil
	}

	query := `
		SELECT id, listing_id, url, storage_path,
		       thumbnail_url, display_order, is_primary,
		       width, height, file_size,
		       mime_type, created_at, updated_at
		FROM listing_images
		WHERE listing_id = ANY($1)
		ORDER BY listing_id, is_primary DESC, display_order ASC, id ASC
	`

	rows, err := r.db.QueryxContext(ctx, query, pq.Array(productIDs))
	if err != nil {
		r.logger.Error().Err(err).Int("product_count", len(productIDs)).Msg("failed to get product images batch")
		return nil, fmt.Errorf("failed to get product images batch: %w", err)
	}
	defer rows.Close()

	result := make(map[int64][]*domain.ProductImage)
	for rows.Next() {
		var image domain.ProductImage
		var prodID sql.NullInt64
		var width, height sql.NullInt32
		var fileSize sql.NullInt64
		var mimeType, storagePath, thumbnailURL sql.NullString

		err := rows.Scan(
			&image.ID,
			&prodID,
			&image.URL,
			&storagePath,
			&thumbnailURL,
			&image.DisplayOrder,
			&image.IsPrimary,
			&width,
			&height,
			&fileSize,
			&mimeType,
			&image.CreatedAt,
			&image.UpdatedAt,
		)
		if err != nil {
			r.logger.Error().Err(err).Msg("failed to scan product image in batch")
			return nil, fmt.Errorf("failed to scan product image: %w", err)
		}

		// Handle nullable fields
		if prodID.Valid {
			image.ProductID = &prodID.Int64
		}
		if storagePath.Valid {
			image.StoragePath = &storagePath.String
		}
		if thumbnailURL.Valid {
			image.ThumbnailURL = &thumbnailURL.String
		}
		if width.Valid {
			image.Width = &width.Int32
		}
		if height.Valid {
			image.Height = &height.Int32
		}
		if fileSize.Valid {
			image.FileSize = &fileSize.Int64
		}
		if mimeType.Valid {
			image.MimeType = &mimeType.String
		}

		if prodID.Valid {
			result[prodID.Int64] = append(result[prodID.Int64], &image)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating product images batch: %w", err)
	}

	return result, nil
}

// DeleteProductImage removes a product image from the database
// Uses listing_images table (unified for both C2C and B2C products)
func (r *Repository) DeleteProductImage(ctx context.Context, imageID int64) error {
	query := `DELETE FROM listing_images WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, imageID)
	if err != nil {
		r.logger.Error().Err(err).Int64("image_id", imageID).Msg("failed to delete product image")
		return fmt.Errorf("failed to delete product image: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product image not found")
	}

	r.logger.Info().Int64("image_id", imageID).Msg("product image deleted")
	return nil
}

// ProductImageOrder represents display order update for a single product image
type ProductImageOrder struct {
	ImageID      int64
	DisplayOrder int32
}

// ReorderProductImages updates display order for multiple product images in a single transaction
// Uses listing_images table (unified for both C2C and B2C products)
func (r *Repository) ReorderProductImages(ctx context.Context, productID int64, orders []ProductImageOrder) error {
	if len(orders) == 0 {
		return nil
	}

	// Start transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to start transaction")
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Build CASE statement for batch update
	query := `UPDATE listing_images SET display_order = CASE `
	args := make([]interface{}, 0, len(orders)*3+1)
	imageIDsForIN := make([]int64, 0, len(orders))

	argIdx := 1
	for _, order := range orders {
		query += fmt.Sprintf("WHEN id = $%d THEN $%d::integer ", argIdx, argIdx+1)
		args = append(args, order.ImageID)
		args = append(args, int(order.DisplayOrder))
		imageIDsForIN = append(imageIDsForIN, order.ImageID)
		argIdx += 2
	}

	query += fmt.Sprintf("END WHERE listing_id = $%d AND id IN (", argIdx)
	args = append(args, productID)

	// Add placeholders for IN clause
	for i, imageID := range imageIDsForIN {
		if i > 0 {
			query += ", "
		}
		query += fmt.Sprintf("$%d", argIdx+1+i)
		args = append(args, imageID)
	}
	query += ")"

	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		r.logger.Error().Err(err).Int64("product_id", productID).Msg("failed to reorder product images")
		return fmt.Errorf("failed to reorder product images: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no product images updated, check product_id and image_ids")
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		r.logger.Error().Err(err).Msg("failed to commit transaction")
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	r.logger.Info().
		Int64("product_id", productID).
		Int("count", len(orders)).
		Msg("product images reordered")

	return nil
}

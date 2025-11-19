package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sveturs/listings/internal/domain"
)

// GetImageByID retrieves a single image by ID
func (r *Repository) GetImageByID(ctx context.Context, imageID int64) (*domain.ListingImage, error) {
	query := `
		SELECT id, listing_id, url, storage_path,
		       thumbnail_url, display_order, is_primary,
		       width, height, file_size,
		       mime_type, created_at, updated_at
		FROM listing_images
		WHERE id = $1
	`

	var image domain.ListingImage
	var width, height sql.NullInt32
	var fileSize sql.NullInt64
	var mimeType, storagePath, thumbnailURL sql.NullString

	err := r.db.QueryRowxContext(ctx, query, imageID).Scan(
		&image.ID,
		&image.ListingID,
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
			return nil, fmt.Errorf("image not found: %w", err)
		}
		r.logger.Error().Err(err).Int64("image_id", imageID).Msg("failed to get image")
		return nil, fmt.Errorf("failed to get image: %w", err)
	}

	// Handle nullable fields
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

// DeleteImage removes an image from the database
func (r *Repository) DeleteImage(ctx context.Context, imageID int64) error {
	query := `
		DELETE FROM listing_images
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, imageID)
	if err != nil {
		r.logger.Error().Err(err).Int64("image_id", imageID).Msg("failed to delete image")
		return fmt.Errorf("failed to delete image: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("image not found")
	}

	r.logger.Info().Int64("image_id", imageID).Msg("image deleted")
	return nil
}

// AddImage adds a new image to a listing
func (r *Repository) AddImage(ctx context.Context, image *domain.ListingImage) (*domain.ListingImage, error) {
	query := `
		INSERT INTO listing_images (
			listing_id, url, storage_path, thumbnail_url, display_order,
			is_primary, width, height, file_size, mime_type
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`

	var newImage domain.ListingImage
	err := r.db.QueryRowxContext(ctx, query,
		image.ListingID,
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
		r.logger.Error().Err(err).Int64("listing_id", image.ListingID).Msg("failed to add image")
		return nil, fmt.Errorf("failed to add image: %w", err)
	}

	// Populate the returned image with input data
	newImage.ListingID = image.ListingID
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

	r.logger.Info().Int64("image_id", newImage.ID).Int64("listing_id", image.ListingID).Msg("image added")
	return &newImage, nil
}

// GetImages retrieves all images for a listing
func (r *Repository) GetImages(ctx context.Context, listingID int64) ([]*domain.ListingImage, error) {
	query := `
		SELECT id, listing_id, url, storage_path,
		       thumbnail_url, display_order, is_primary,
		       width, height, file_size,
		       mime_type, created_at, updated_at
		FROM listing_images
		WHERE listing_id = $1
		ORDER BY is_primary DESC, display_order ASC, id ASC
	`

	rows, err := r.db.QueryxContext(ctx, query, listingID)
	if err != nil {
		r.logger.Error().Err(err).Int64("listing_id", listingID).Msg("failed to get images")
		return nil, fmt.Errorf("failed to get images: %w", err)
	}
	defer rows.Close()

	var images []*domain.ListingImage
	for rows.Next() {
		var image domain.ListingImage
		var width, height sql.NullInt32
		var fileSize sql.NullInt64
		var mimeType, storagePath, thumbnailURL sql.NullString

		err := rows.Scan(
			&image.ID,
			&image.ListingID,
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
			r.logger.Error().Err(err).Msg("failed to scan image")
			return nil, fmt.Errorf("failed to scan image: %w", err)
		}

		// Handle nullable fields
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
		return nil, fmt.Errorf("error iterating images: %w", err)
	}

	return images, nil
}

// ImageOrder represents display order update for a single image
type ImageOrder struct {
	ImageID      int64
	DisplayOrder int32
}

// ReorderImages updates display order for multiple images in a single transaction
func (r *Repository) ReorderImages(ctx context.Context, listingID int64, orders []ImageOrder) error {
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
	// UPDATE listing_images SET display_order = CASE
	//   WHEN id = $1 THEN $2
	//   WHEN id = $3 THEN $4
	//   ...
	// END
	// WHERE listing_id = $N AND id IN ($1, $3, ...)

	query := `UPDATE listing_images SET display_order = CASE `
	args := make([]interface{}, 0, len(orders)*3+1) // Preallocate for CASE + WHERE + IN
	imageIDsForIN := make([]int64, 0, len(orders))

	argIdx := 1
	for _, order := range orders {
		// Use explicit ::integer cast in SQL to force correct type
		query += fmt.Sprintf("WHEN id = $%d THEN $%d::integer ", argIdx, argIdx+1)
		// Append typed values, NOT interface{} wrapping
		args = append(args, order.ImageID)          // int64
		args = append(args, int(order.DisplayOrder)) // int (cast from int32)
		imageIDsForIN = append(imageIDsForIN, order.ImageID)
		argIdx += 2
	}

	query += fmt.Sprintf("END WHERE listing_id = $%d AND id IN (", argIdx)
	args = append(args, listingID) // int64

	// Add placeholders for IN clause
	for i, imageID := range imageIDsForIN {
		if i > 0 {
			query += ", "
		}
		query += fmt.Sprintf("$%d", argIdx+1+i)
		args = append(args, imageID) // int64 directly, not interface{}
	}
	query += ")"

	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		r.logger.Error().Err(err).Int64("listing_id", listingID).Msg("failed to reorder images")
		return fmt.Errorf("failed to reorder images: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no images updated, check listing_id and image_ids")
	}

	if int64(rowsAffected) != int64(len(orders)) {
		r.logger.Warn().
			Int64("listing_id", listingID).
			Int64("expected", int64(len(orders))).
			Int64("actual", rowsAffected).
			Msg("not all images were updated")
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		r.logger.Error().Err(err).Msg("failed to commit transaction")
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	r.logger.Info().
		Int64("listing_id", listingID).
		Int("count", len(orders)).
		Msg("images reordered")

	return nil
}

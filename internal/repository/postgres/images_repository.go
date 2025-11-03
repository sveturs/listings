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
		SELECT id, listing_id, public_url as url, file_path as storage_path,
		       NULL as thumbnail_url, 1 as display_order, is_main as is_primary,
		       NULL::integer as width, NULL::integer as height, file_size,
		       content_type as mime_type, created_at, created_at as updated_at
		FROM c2c_images
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
		DELETE FROM c2c_images
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
	// Convert domain image to c2c_images schema
	query := `
		INSERT INTO c2c_images (
			listing_id, file_path, file_name, file_size, content_type,
			is_main, storage_type, storage_bucket, public_url
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at
	`

	// Extract filename from storage_path or use default
	fileName := "image.jpg"
	if image.StoragePath != nil && *image.StoragePath != "" {
		// Simple filename extraction
		fileName = *image.StoragePath
	}

	// Default storage type to minio
	storageType := "minio"
	storageBucket := sql.NullString{String: "listings", Valid: true}
	publicURL := sql.NullString{String: image.URL, Valid: image.URL != ""}

	var fileSize sql.NullInt64
	if image.FileSize != nil {
		fileSize = sql.NullInt64{Int64: *image.FileSize, Valid: true}
	}

	var mimeType sql.NullString
	if image.MimeType != nil {
		mimeType = sql.NullString{String: *image.MimeType, Valid: true}
	}

	var storagePath string
	if image.StoragePath != nil {
		storagePath = *image.StoragePath
	} else {
		storagePath = ""
	}

	var newImage domain.ListingImage
	err := r.db.QueryRowxContext(ctx, query,
		image.ListingID,
		storagePath,
		fileName,
		fileSize,
		mimeType,
		image.IsPrimary,
		storageType,
		storageBucket,
		publicURL,
	).Scan(&newImage.ID, &newImage.CreatedAt)

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
		SELECT id, listing_id, public_url as url, file_path as storage_path,
		       NULL as thumbnail_url, 1 as display_order, is_main as is_primary,
		       NULL::integer as width, NULL::integer as height, file_size,
		       content_type as mime_type, created_at, created_at as updated_at
		FROM c2c_images
		WHERE listing_id = $1
		ORDER BY is_main DESC, id ASC
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

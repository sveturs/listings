package postgres

import (
	"context"
	"fmt"

	"backend/internal/domain/models"
	"backend/internal/logger"
)

// CreateProductVariant creates a new product variant
func (s *Database) CreateProductVariant(ctx context.Context, variant *models.CreateProductVariantRequest) (*models.StorefrontProductVariant, error) {
	query := `
		INSERT INTO storefront_product_variants (
			product_id, sku, barcode, price, compare_at_price, cost_price,
			stock_quantity, stock_status, low_stock_threshold,
			variant_attributes, weight, dimensions,
			is_active, is_default
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9,
			$10, $11, $12,
			$13, $14
		) RETURNING id, product_id, sku, barcode, price, compare_at_price, cost_price,
			stock_quantity, stock_status, low_stock_threshold,
			variant_attributes, weight, dimensions,
			is_active, is_default, view_count, sold_count, created_at, updated_at`

	var created models.StorefrontProductVariant
	err := s.pool.QueryRow(ctx, query,
		variant.ProductID, variant.SKU, variant.Barcode, variant.Price, variant.CompareAtPrice, variant.CostPrice,
		variant.StockQuantity, variant.StockStatus, variant.LowStockThreshold,
		variant.VariantAttributes, variant.Weight, variant.Dimensions,
		variant.IsActive, variant.IsDefault,
	).Scan(
		&created.ID, &created.ProductID, &created.SKU, &created.Barcode, &created.Price, &created.CompareAtPrice, &created.CostPrice,
		&created.StockQuantity, &created.StockStatus, &created.LowStockThreshold,
		&created.VariantAttributes, &created.Weight, &created.Dimensions,
		&created.IsActive, &created.IsDefault, &created.ViewCount, &created.SoldCount, &created.CreatedAt, &created.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create product variant: %w", err)
	}

	logger.Info().
		Int("variant_id", created.ID).
		Int("product_id", created.ProductID).
		Msg("Product variant created successfully")

	return &created, nil
}

// BatchCreateProductVariants creates multiple product variants in a single transaction
func (s *Database) BatchCreateProductVariants(ctx context.Context, variants []*models.CreateProductVariantRequest) ([]*models.StorefrontProductVariant, error) {
	if len(variants) == 0 {
		return []*models.StorefrontProductVariant{}, nil
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck // Will be handled by commit

	created := make([]*models.StorefrontProductVariant, 0, len(variants))

	for _, variant := range variants {
		query := `
			INSERT INTO storefront_product_variants (
				product_id, sku, barcode, price, compare_at_price, cost_price,
				stock_quantity, stock_status, low_stock_threshold,
				variant_attributes, weight, dimensions,
				is_active, is_default
			) VALUES (
				$1, $2, $3, $4, $5, $6,
				$7, $8, $9,
				$10, $11, $12,
				$13, $14
			) RETURNING id, product_id, sku, barcode, price, compare_at_price, cost_price,
				stock_quantity, stock_status, low_stock_threshold,
				variant_attributes, weight, dimensions,
				is_active, is_default, view_count, sold_count, created_at, updated_at`

		var v models.StorefrontProductVariant
		err := tx.QueryRow(ctx, query,
			variant.ProductID, variant.SKU, variant.Barcode, variant.Price, variant.CompareAtPrice, variant.CostPrice,
			variant.StockQuantity, variant.StockStatus, variant.LowStockThreshold,
			variant.VariantAttributes, variant.Weight, variant.Dimensions,
			variant.IsActive, variant.IsDefault,
		).Scan(
			&v.ID, &v.ProductID, &v.SKU, &v.Barcode, &v.Price, &v.CompareAtPrice, &v.CostPrice,
			&v.StockQuantity, &v.StockStatus, &v.LowStockThreshold,
			&v.VariantAttributes, &v.Weight, &v.Dimensions,
			&v.IsActive, &v.IsDefault, &v.ViewCount, &v.SoldCount, &v.CreatedAt, &v.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create variant #%d: %w", len(created)+1, err)
		}

		created = append(created, &v)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	logger.Info().
		Int("count", len(created)).
		Int("product_id", created[0].ProductID).
		Msg("Batch created product variants successfully")

	return created, nil
}

// CreateProductVariantImage creates a new variant image
func (s *Database) CreateProductVariantImage(ctx context.Context, image *models.CreateProductVariantImageRequest) (*models.StorefrontProductVariantImage, error) {
	query := `
		INSERT INTO storefront_product_variant_images (
			variant_id, image_url, thumbnail_url, alt_text, display_order, is_main
		) VALUES (
			$1, $2, $3, $4, $5, $6
		) RETURNING id, variant_id, image_url, thumbnail_url, alt_text, display_order, is_main, created_at`

	var created models.StorefrontProductVariantImage
	err := s.pool.QueryRow(ctx, query,
		image.VariantID, image.ImageURL, image.ThumbnailURL, image.AltText, image.DisplayOrder, image.IsMain,
	).Scan(
		&created.ID, &created.VariantID, &created.ImageURL, &created.ThumbnailURL,
		&created.AltText, &created.DisplayOrder, &created.IsMain, &created.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create variant image: %w", err)
	}

	return &created, nil
}

// BatchCreateProductVariantImages creates multiple variant images in a single transaction
func (s *Database) BatchCreateProductVariantImages(ctx context.Context, images []*models.CreateProductVariantImageRequest) ([]*models.StorefrontProductVariantImage, error) {
	if len(images) == 0 {
		return []*models.StorefrontProductVariantImage{}, nil
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck // Will be handled by commit

	created := make([]*models.StorefrontProductVariantImage, 0, len(images))

	for _, image := range images {
		query := `
			INSERT INTO storefront_product_variant_images (
				variant_id, image_url, thumbnail_url, alt_text, display_order, is_main
			) VALUES (
				$1, $2, $3, $4, $5, $6
			) RETURNING id, variant_id, image_url, thumbnail_url, alt_text, display_order, is_main, created_at`

		var img models.StorefrontProductVariantImage
		err := tx.QueryRow(ctx, query,
			image.VariantID, image.ImageURL, image.ThumbnailURL, image.AltText, image.DisplayOrder, image.IsMain,
		).Scan(
			&img.ID, &img.VariantID, &img.ImageURL, &img.ThumbnailURL,
			&img.AltText, &img.DisplayOrder, &img.IsMain, &img.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create variant image #%d: %w", len(created)+1, err)
		}

		created = append(created, &img)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	logger.Info().
		Int("count", len(created)).
		Msg("Batch created variant images successfully")

	return created, nil
}

// GetProductVariants retrieves all variants for a product
func (s *Database) GetProductVariants(ctx context.Context, productID int) ([]*models.StorefrontProductVariant, error) {
	query := `
		SELECT
			id, product_id, sku, barcode, price, compare_at_price, cost_price,
			stock_quantity, stock_status, low_stock_threshold,
			variant_attributes, weight, dimensions,
			is_active, is_default, view_count, sold_count, created_at, updated_at
		FROM storefront_product_variants
		WHERE product_id = $1 AND is_active = true
		ORDER BY is_default DESC, id ASC`

	rows, err := s.pool.Query(ctx, query, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to query product variants: %w", err)
	}
	defer rows.Close()

	var variants []*models.StorefrontProductVariant
	for rows.Next() {
		var v models.StorefrontProductVariant
		err := rows.Scan(
			&v.ID, &v.ProductID, &v.SKU, &v.Barcode, &v.Price, &v.CompareAtPrice, &v.CostPrice,
			&v.StockQuantity, &v.StockStatus, &v.LowStockThreshold,
			&v.VariantAttributes, &v.Weight, &v.Dimensions,
			&v.IsActive, &v.IsDefault, &v.ViewCount, &v.SoldCount, &v.CreatedAt, &v.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan variant: %w", err)
		}
		variants = append(variants, &v)
	}

	return variants, nil
}

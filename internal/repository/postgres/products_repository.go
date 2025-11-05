package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/sveturs/listings/internal/domain"
)

// GetProductByID retrieves a single product by ID with optional storefront filter
func (r *Repository) GetProductByID(ctx context.Context, productID int64, storefrontID *int64) (*domain.Product, error) {
	r.logger.Debug().Int64("product_id", productID).Interface("storefront_id", storefrontID).Msg("getting product by ID")

	query := `
		SELECT
			p.id, p.storefront_id, p.name, p.description, p.price, p.currency,
			p.category_id, p.sku, p.barcode, p.stock_quantity, p.stock_status,
			p.is_active, p.attributes, p.view_count, p.sold_count,
			p.created_at, p.updated_at,
			p.has_individual_location, p.individual_address,
			p.individual_latitude, p.individual_longitude,
			p.location_privacy, p.show_on_map, p.has_variants
		FROM b2c_products p
		WHERE p.id = $1
		  AND ($2::bigint IS NULL OR p.storefront_id = $2)
		  AND p.deleted_at IS NULL
	`

	var product domain.Product
	var description sql.NullString
	var sku, barcode sql.NullString
	var individualAddress, locationPrivacy sql.NullString
	var individualLatitude, individualLongitude sql.NullFloat64
	var attributesJSON []byte

	err := r.db.QueryRowContext(ctx, query, productID, storefrontID).Scan(
		&product.ID,
		&product.StorefrontID,
		&product.Name,
		&description,
		&product.Price,
		&product.Currency,
		&product.CategoryID,
		&sku,
		&barcode,
		&product.StockQuantity,
		&product.StockStatus,
		&product.IsActive,
		&attributesJSON,
		&product.ViewCount,
		&product.SoldCount,
		&product.CreatedAt,
		&product.UpdatedAt,
		&product.HasIndividualLocation,
		&individualAddress,
		&individualLatitude,
		&individualLongitude,
		&locationPrivacy,
		&product.ShowOnMap,
		&product.HasVariants,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		r.logger.Error().Err(err).Int64("product_id", productID).Msg("failed to get product by ID")
		return nil, fmt.Errorf("failed to get product by ID: %w", err)
	}

	// Handle nullable fields
	if description.Valid {
		product.Description = description.String
	}
	if sku.Valid {
		product.SKU = &sku.String
	}
	if barcode.Valid {
		product.Barcode = &barcode.String
	}
	if individualAddress.Valid {
		product.IndividualAddress = &individualAddress.String
	}
	if individualLatitude.Valid {
		product.IndividualLatitude = &individualLatitude.Float64
	}
	if individualLongitude.Valid {
		product.IndividualLongitude = &individualLongitude.Float64
	}
	if locationPrivacy.Valid {
		product.LocationPrivacy = &locationPrivacy.String
	}

	// Parse JSONB attributes
	if len(attributesJSON) > 0 {
		if err := json.Unmarshal(attributesJSON, &product.Attributes); err != nil {
			r.logger.Error().Err(err).Msg("failed to unmarshal product attributes")
			return nil, fmt.Errorf("failed to unmarshal product attributes: %w", err)
		}
	}

	return &product, nil
}

// GetProductsBySKUs retrieves products by list of SKUs (для корзины)
func (r *Repository) GetProductsBySKUs(ctx context.Context, skus []string, storefrontID *int64) ([]*domain.Product, error) {
	r.logger.Debug().Int("sku_count", len(skus)).Interface("storefront_id", storefrontID).Msg("getting products by SKUs")

	if len(skus) == 0 {
		return []*domain.Product{}, nil
	}

	query := `
		SELECT
			p.id, p.storefront_id, p.name, p.description, p.price, p.currency,
			p.category_id, p.sku, p.barcode, p.stock_quantity, p.stock_status,
			p.is_active, p.attributes, p.view_count, p.sold_count,
			p.created_at, p.updated_at,
			p.has_individual_location, p.individual_address,
			p.individual_latitude, p.individual_longitude,
			p.location_privacy, p.show_on_map, p.has_variants
		FROM b2c_products p
		WHERE p.sku = ANY($1::text[])
		  AND ($2::bigint IS NULL OR p.storefront_id = $2)
		  AND p.is_active = true
		  AND p.deleted_at IS NULL
	`

	rows, err := r.db.QueryContext(ctx, query, pq.Array(skus), storefrontID)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to query products by SKUs")
		return nil, fmt.Errorf("failed to query products by SKUs: %w", err)
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		var product domain.Product
		var description sql.NullString
		var sku, barcode sql.NullString
		var individualAddress, locationPrivacy sql.NullString
		var individualLatitude, individualLongitude sql.NullFloat64
		var attributesJSON []byte

		err := rows.Scan(
			&product.ID,
			&product.StorefrontID,
			&product.Name,
			&description,
			&product.Price,
			&product.Currency,
			&product.CategoryID,
			&sku,
			&barcode,
			&product.StockQuantity,
			&product.StockStatus,
			&product.IsActive,
			&attributesJSON,
			&product.ViewCount,
			&product.SoldCount,
			&product.CreatedAt,
			&product.UpdatedAt,
			&product.HasIndividualLocation,
			&individualAddress,
			&individualLatitude,
			&individualLongitude,
			&locationPrivacy,
			&product.ShowOnMap,
			&product.HasVariants,
		)
		if err != nil {
			r.logger.Error().Err(err).Msg("failed to scan product")
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}

		// Handle nullable fields
		if description.Valid {
			product.Description = description.String
		}
		if sku.Valid {
			product.SKU = &sku.String
		}
		if barcode.Valid {
			product.Barcode = &barcode.String
		}
		if individualAddress.Valid {
			product.IndividualAddress = &individualAddress.String
		}
		if individualLatitude.Valid {
			product.IndividualLatitude = &individualLatitude.Float64
		}
		if individualLongitude.Valid {
			product.IndividualLongitude = &individualLongitude.Float64
		}
		if locationPrivacy.Valid {
			product.LocationPrivacy = &locationPrivacy.String
		}

		// Parse JSONB attributes
		if len(attributesJSON) > 0 {
			if err := json.Unmarshal(attributesJSON, &product.Attributes); err != nil {
				r.logger.Error().Err(err).Msg("failed to unmarshal product attributes")
				return nil, fmt.Errorf("failed to unmarshal product attributes: %w", err)
			}
		}

		products = append(products, &product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return products, nil
}

// GetProductsByIDs retrieves products by list of IDs (для корзины)
func (r *Repository) GetProductsByIDs(ctx context.Context, productIDs []int64, storefrontID *int64) ([]*domain.Product, error) {
	r.logger.Debug().Int("id_count", len(productIDs)).Interface("storefront_id", storefrontID).Msg("getting products by IDs")

	if len(productIDs) == 0 {
		return []*domain.Product{}, nil
	}

	query := `
		SELECT
			p.id, p.storefront_id, p.name, p.description, p.price, p.currency,
			p.category_id, p.sku, p.barcode, p.stock_quantity, p.stock_status,
			p.is_active, p.attributes, p.view_count, p.sold_count,
			p.created_at, p.updated_at,
			p.has_individual_location, p.individual_address,
			p.individual_latitude, p.individual_longitude,
			p.location_privacy, p.show_on_map, p.has_variants
		FROM b2c_products p
		WHERE p.id = ANY($1::bigint[])
		  AND ($2::bigint IS NULL OR p.storefront_id = $2)
		  AND p.is_active = true
		  AND p.deleted_at IS NULL
	`

	rows, err := r.db.QueryContext(ctx, query, pq.Array(productIDs), storefrontID)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to query products by IDs")
		return nil, fmt.Errorf("failed to query products by IDs: %w", err)
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		var product domain.Product
		var description sql.NullString
		var sku, barcode sql.NullString
		var individualAddress, locationPrivacy sql.NullString
		var individualLatitude, individualLongitude sql.NullFloat64
		var attributesJSON []byte

		err := rows.Scan(
			&product.ID,
			&product.StorefrontID,
			&product.Name,
			&description,
			&product.Price,
			&product.Currency,
			&product.CategoryID,
			&sku,
			&barcode,
			&product.StockQuantity,
			&product.StockStatus,
			&product.IsActive,
			&attributesJSON,
			&product.ViewCount,
			&product.SoldCount,
			&product.CreatedAt,
			&product.UpdatedAt,
			&product.HasIndividualLocation,
			&individualAddress,
			&individualLatitude,
			&individualLongitude,
			&locationPrivacy,
			&product.ShowOnMap,
			&product.HasVariants,
		)
		if err != nil {
			r.logger.Error().Err(err).Msg("failed to scan product")
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}

		// Handle nullable fields
		if description.Valid {
			product.Description = description.String
		}
		if sku.Valid {
			product.SKU = &sku.String
		}
		if barcode.Valid {
			product.Barcode = &barcode.String
		}
		if individualAddress.Valid {
			product.IndividualAddress = &individualAddress.String
		}
		if individualLatitude.Valid {
			product.IndividualLatitude = &individualLatitude.Float64
		}
		if individualLongitude.Valid {
			product.IndividualLongitude = &individualLongitude.Float64
		}
		if locationPrivacy.Valid {
			product.LocationPrivacy = &locationPrivacy.String
		}

		// Parse JSONB attributes
		if len(attributesJSON) > 0 {
			if err := json.Unmarshal(attributesJSON, &product.Attributes); err != nil {
				r.logger.Error().Err(err).Msg("failed to unmarshal product attributes")
				return nil, fmt.Errorf("failed to unmarshal product attributes: %w", err)
			}
		}

		products = append(products, &product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return products, nil
}

// ListProducts retrieves products with pagination
func (r *Repository) ListProducts(ctx context.Context, storefrontID int64, page, pageSize int, isActiveOnly bool) ([]*domain.Product, int, error) {
	r.logger.Debug().Int64("storefront_id", storefrontID).Int("page", page).Int("page_size", pageSize).Bool("is_active_only", isActiveOnly).Msg("listing products")

	offset := (page - 1) * pageSize

	// Get total count
	countQuery := `
		SELECT COUNT(*)
		FROM b2c_products p
		WHERE p.storefront_id = $1
		  AND ($2 = false OR p.is_active = true)
	`

	var total int
	err := r.db.QueryRowContext(ctx, countQuery, storefrontID, isActiveOnly).Scan(&total)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to count products")
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}

	// Get products with pagination
	query := `
		SELECT
			p.id, p.storefront_id, p.name, p.description, p.price, p.currency,
			p.category_id, p.sku, p.barcode, p.stock_quantity, p.stock_status,
			p.is_active, p.attributes, p.view_count, p.sold_count,
			p.created_at, p.updated_at,
			p.has_individual_location, p.individual_address,
			p.individual_latitude, p.individual_longitude,
			p.location_privacy, p.show_on_map, p.has_variants
		FROM b2c_products p
		WHERE p.storefront_id = $1
		  AND ($2 = false OR p.is_active = true)
		ORDER BY p.created_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.QueryContext(ctx, query, storefrontID, isActiveOnly, pageSize, offset)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to list products")
		return nil, 0, fmt.Errorf("failed to list products: %w", err)
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		var product domain.Product
		var description sql.NullString
		var sku, barcode sql.NullString
		var individualAddress, locationPrivacy sql.NullString
		var individualLatitude, individualLongitude sql.NullFloat64
		var attributesJSON []byte

		err := rows.Scan(
			&product.ID,
			&product.StorefrontID,
			&product.Name,
			&description,
			&product.Price,
			&product.Currency,
			&product.CategoryID,
			&sku,
			&barcode,
			&product.StockQuantity,
			&product.StockStatus,
			&product.IsActive,
			&attributesJSON,
			&product.ViewCount,
			&product.SoldCount,
			&product.CreatedAt,
			&product.UpdatedAt,
			&product.HasIndividualLocation,
			&individualAddress,
			&individualLatitude,
			&individualLongitude,
			&locationPrivacy,
			&product.ShowOnMap,
			&product.HasVariants,
		)
		if err != nil {
			r.logger.Error().Err(err).Msg("failed to scan product")
			return nil, 0, fmt.Errorf("failed to scan product: %w", err)
		}

		// Handle nullable fields
		if description.Valid {
			product.Description = description.String
		}
		if sku.Valid {
			product.SKU = &sku.String
		}
		if barcode.Valid {
			product.Barcode = &barcode.String
		}
		if individualAddress.Valid {
			product.IndividualAddress = &individualAddress.String
		}
		if individualLatitude.Valid {
			product.IndividualLatitude = &individualLatitude.Float64
		}
		if individualLongitude.Valid {
			product.IndividualLongitude = &individualLongitude.Float64
		}
		if locationPrivacy.Valid {
			product.LocationPrivacy = &locationPrivacy.String
		}

		// Parse JSONB attributes
		if len(attributesJSON) > 0 {
			if err := json.Unmarshal(attributesJSON, &product.Attributes); err != nil {
				r.logger.Error().Err(err).Msg("failed to unmarshal product attributes")
				return nil, 0, fmt.Errorf("failed to unmarshal product attributes: %w", err)
			}
		}

		products = append(products, &product)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows iteration error: %w", err)
	}

	return products, total, nil
}

// GetVariantByID retrieves a single variant by ID
func (r *Repository) GetVariantByID(ctx context.Context, variantID int64, productID *int64) (*domain.ProductVariant, error) {
	r.logger.Debug().Int64("variant_id", variantID).Interface("product_id", productID).Msg("getting variant by ID")

	query := `
		SELECT
			v.id, v.product_id, v.sku, v.barcode, v.price, v.compare_at_price,
			v.cost_price, v.stock_quantity, v.stock_status, v.low_stock_threshold,
			v.variant_attributes, v.weight, v.dimensions, v.is_active, v.is_default,
			v.view_count, v.sold_count, v.created_at, v.updated_at
		FROM b2c_product_variants v
		WHERE v.id = $1
		  AND ($2::bigint IS NULL OR v.product_id = $2)
	`

	var variant domain.ProductVariant
	var sku, barcode sql.NullString
	var price, compareAtPrice, costPrice, weight sql.NullFloat64
	var lowStockThreshold sql.NullInt32
	var variantAttributesJSON, dimensionsJSON []byte

	err := r.db.QueryRowContext(ctx, query, variantID, productID).Scan(
		&variant.ID,
		&variant.ProductID,
		&sku,
		&barcode,
		&price,
		&compareAtPrice,
		&costPrice,
		&variant.StockQuantity,
		&variant.StockStatus,
		&lowStockThreshold,
		&variantAttributesJSON,
		&weight,
		&dimensionsJSON,
		&variant.IsActive,
		&variant.IsDefault,
		&variant.ViewCount,
		&variant.SoldCount,
		&variant.CreatedAt,
		&variant.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("variant not found")
		}
		r.logger.Error().Err(err).Int64("variant_id", variantID).Msg("failed to get variant by ID")
		return nil, fmt.Errorf("failed to get variant by ID: %w", err)
	}

	// Handle nullable fields
	if sku.Valid {
		variant.SKU = &sku.String
	}
	if barcode.Valid {
		variant.Barcode = &barcode.String
	}
	if price.Valid {
		variant.Price = &price.Float64
	}
	if compareAtPrice.Valid {
		variant.CompareAtPrice = &compareAtPrice.Float64
	}
	if costPrice.Valid {
		variant.CostPrice = &costPrice.Float64
	}
	if weight.Valid {
		variant.Weight = &weight.Float64
	}
	if lowStockThreshold.Valid {
		variant.LowStockThreshold = &lowStockThreshold.Int32
	}

	// Parse JSONB fields
	if len(variantAttributesJSON) > 0 {
		if err := json.Unmarshal(variantAttributesJSON, &variant.VariantAttributes); err != nil {
			r.logger.Error().Err(err).Msg("failed to unmarshal variant attributes")
			return nil, fmt.Errorf("failed to unmarshal variant attributes: %w", err)
		}
	}
	if len(dimensionsJSON) > 0 {
		if err := json.Unmarshal(dimensionsJSON, &variant.Dimensions); err != nil {
			r.logger.Error().Err(err).Msg("failed to unmarshal variant dimensions")
			return nil, fmt.Errorf("failed to unmarshal variant dimensions: %w", err)
		}
	}

	return &variant, nil
}

// GetVariantsByProductID retrieves all variants for a product
func (r *Repository) GetVariantsByProductID(ctx context.Context, productID int64, isActiveOnly bool) ([]*domain.ProductVariant, error) {
	r.logger.Debug().Int64("product_id", productID).Bool("is_active_only", isActiveOnly).Msg("getting variants by product ID")

	query := `
		SELECT
			v.id, v.product_id, v.sku, v.barcode, v.price, v.compare_at_price,
			v.cost_price, v.stock_quantity, v.stock_status, v.low_stock_threshold,
			v.variant_attributes, v.weight, v.dimensions, v.is_active, v.is_default,
			v.view_count, v.sold_count, v.created_at, v.updated_at
		FROM b2c_product_variants v
		WHERE v.product_id = $1
		  AND ($2 = false OR v.is_active = true)
		ORDER BY v.is_default DESC, v.id ASC
	`

	rows, err := r.db.QueryContext(ctx, query, productID, isActiveOnly)
	if err != nil {
		r.logger.Error().Err(err).Int64("product_id", productID).Msg("failed to query variants")
		return nil, fmt.Errorf("failed to query variants: %w", err)
	}
	defer rows.Close()

	var variants []*domain.ProductVariant
	for rows.Next() {
		var variant domain.ProductVariant
		var sku, barcode sql.NullString
		var price, compareAtPrice, costPrice, weight sql.NullFloat64
		var lowStockThreshold sql.NullInt32
		var variantAttributesJSON, dimensionsJSON []byte

		err := rows.Scan(
			&variant.ID,
			&variant.ProductID,
			&sku,
			&barcode,
			&price,
			&compareAtPrice,
			&costPrice,
			&variant.StockQuantity,
			&variant.StockStatus,
			&lowStockThreshold,
			&variantAttributesJSON,
			&weight,
			&dimensionsJSON,
			&variant.IsActive,
			&variant.IsDefault,
			&variant.ViewCount,
			&variant.SoldCount,
			&variant.CreatedAt,
			&variant.UpdatedAt,
		)
		if err != nil {
			r.logger.Error().Err(err).Msg("failed to scan variant")
			return nil, fmt.Errorf("failed to scan variant: %w", err)
		}

		// Handle nullable fields
		if sku.Valid {
			variant.SKU = &sku.String
		}
		if barcode.Valid {
			variant.Barcode = &barcode.String
		}
		if price.Valid {
			variant.Price = &price.Float64
		}
		if compareAtPrice.Valid {
			variant.CompareAtPrice = &compareAtPrice.Float64
		}
		if costPrice.Valid {
			variant.CostPrice = &costPrice.Float64
		}
		if weight.Valid {
			variant.Weight = &weight.Float64
		}
		if lowStockThreshold.Valid {
			variant.LowStockThreshold = &lowStockThreshold.Int32
		}

		// Parse JSONB fields
		if len(variantAttributesJSON) > 0 {
			if err := json.Unmarshal(variantAttributesJSON, &variant.VariantAttributes); err != nil {
				r.logger.Error().Err(err).Msg("failed to unmarshal variant attributes")
				return nil, fmt.Errorf("failed to unmarshal variant attributes: %w", err)
			}
		}
		if len(dimensionsJSON) > 0 {
			if err := json.Unmarshal(dimensionsJSON, &variant.Dimensions); err != nil {
				r.logger.Error().Err(err).Msg("failed to unmarshal variant dimensions")
				return nil, fmt.Errorf("failed to unmarshal variant dimensions: %w", err)
			}
		}

		variants = append(variants, &variant)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return variants, nil
}

// CreateProduct creates a new product in the database
func (r *Repository) CreateProduct(ctx context.Context, input *domain.CreateProductInput) (*domain.Product, error) {
	r.logger.Debug().
		Int64("storefront_id", input.StorefrontID).
		Str("name", input.Name).
		Msg("creating product")

	// Validate input
	if input.StorefrontID <= 0 {
		return nil, fmt.Errorf("storefront_id must be greater than 0")
	}
	if input.Name == "" {
		return nil, fmt.Errorf("product name cannot be empty")
	}
	if input.Price < 0 {
		return nil, fmt.Errorf("price must be non-negative")
	}
	if input.StockQuantity < 0 {
		return nil, fmt.Errorf("stock_quantity must be non-negative")
	}

	// Check if storefront exists (basic validation)
	// Note: In production, this should be a foreign key constraint in DB
	// For now, we'll rely on DB constraint

	// Marshal attributes to JSON
	var attributesJSON []byte
	var err error
	if input.Attributes != nil {
		attributesJSON, err = json.Marshal(input.Attributes)
		if err != nil {
			r.logger.Error().Err(err).Msg("failed to marshal product attributes")
			return nil, fmt.Errorf("failed to marshal product attributes: %w", err)
		}
	} else {
		// Default to empty JSON object if nil
		attributesJSON = []byte("{}")
	}

	// Determine stock_status based on quantity
	stockStatus := domain.StockStatusOutOfStock
	if input.StockQuantity > 0 {
		stockStatus = domain.StockStatusInStock
	}

	// Set default values
	isActive := true
	showOnMap := input.ShowOnMap
	hasIndividualLocation := input.HasIndividualLocation

	// Insert product
	query := `
		INSERT INTO b2c_products (
			storefront_id, name, description, price, currency, category_id,
			sku, barcode, stock_quantity, stock_status, is_active,
			attributes, view_count, sold_count,
			has_individual_location, individual_address,
			individual_latitude, individual_longitude,
			location_privacy, show_on_map, has_variants
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11,
			$12, 0, 0,
			$13, $14, $15, $16, $17, $18, false
		)
		RETURNING
			id, storefront_id, name, description, price, currency,
			category_id, sku, barcode, stock_quantity, stock_status,
			is_active, attributes, view_count, sold_count,
			created_at, updated_at,
			has_individual_location, individual_address,
			individual_latitude, individual_longitude,
			location_privacy, show_on_map, has_variants
	`

	var product domain.Product
	var description sql.NullString
	var sku, barcode sql.NullString
	var individualAddress, locationPrivacy sql.NullString
	var individualLatitude, individualLongitude sql.NullFloat64
	var returnedAttributesJSON []byte

	err = r.db.QueryRowContext(
		ctx,
		query,
		input.StorefrontID,
		input.Name,
		input.Description,
		input.Price,
		input.Currency,
		input.CategoryID,
		input.SKU,
		input.Barcode,
		input.StockQuantity,
		stockStatus,
		isActive,
		attributesJSON,
		hasIndividualLocation,
		input.IndividualAddress,
		input.IndividualLatitude,
		input.IndividualLongitude,
		input.LocationPrivacy,
		showOnMap,
	).Scan(
		&product.ID,
		&product.StorefrontID,
		&product.Name,
		&description,
		&product.Price,
		&product.Currency,
		&product.CategoryID,
		&sku,
		&barcode,
		&product.StockQuantity,
		&product.StockStatus,
		&product.IsActive,
		&returnedAttributesJSON,
		&product.ViewCount,
		&product.SoldCount,
		&product.CreatedAt,
		&product.UpdatedAt,
		&product.HasIndividualLocation,
		&individualAddress,
		&individualLatitude,
		&individualLongitude,
		&locationPrivacy,
		&product.ShowOnMap,
		&product.HasVariants,
	)

	if err != nil {
		// Check for unique constraint violation (duplicate SKU)
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // unique_violation
				r.logger.Error().Err(err).Msg("duplicate SKU")
				return nil, fmt.Errorf("products.sku_duplicate")
			}
		}
		r.logger.Error().Err(err).Msg("failed to create product")
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	// Handle nullable fields
	if description.Valid {
		product.Description = description.String
	}
	if sku.Valid {
		product.SKU = &sku.String
	}
	if barcode.Valid {
		product.Barcode = &barcode.String
	}
	if individualAddress.Valid {
		product.IndividualAddress = &individualAddress.String
	}
	if individualLatitude.Valid {
		product.IndividualLatitude = &individualLatitude.Float64
	}
	if individualLongitude.Valid {
		product.IndividualLongitude = &individualLongitude.Float64
	}
	if locationPrivacy.Valid {
		product.LocationPrivacy = &locationPrivacy.String
	}

	// Parse returned JSONB attributes
	if len(returnedAttributesJSON) > 0 {
		if err := json.Unmarshal(returnedAttributesJSON, &product.Attributes); err != nil {
			r.logger.Error().Err(err).Msg("failed to unmarshal returned product attributes")
			return nil, fmt.Errorf("failed to unmarshal product attributes: %w", err)
		}
	}

	r.logger.Info().Int64("product_id", product.ID).Msg("product created successfully")
	return &product, nil
}

// UpdateProduct updates an existing product with ownership validation
func (r *Repository) UpdateProduct(ctx context.Context, productID int64, storefrontID int64, input *domain.UpdateProductInput) (*domain.Product, error) {
	r.logger.Debug().
		Int64("product_id", productID).
		Int64("storefront_id", storefrontID).
		Msg("updating product")

	// Validate input
	if productID <= 0 {
		return nil, fmt.Errorf("product_id must be greater than 0")
	}
	if storefrontID <= 0 {
		return nil, fmt.Errorf("storefront_id must be greater than 0")
	}

	// Build dynamic UPDATE query based on provided fields
	setClauses := []string{"updated_at = NOW()"}
	args := []interface{}{}
	argIndex := 1

	// Add fields to update
	if input.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, *input.Name)
		argIndex++
	}

	if input.Description != nil {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, *input.Description)
		argIndex++
	}

	if input.Price != nil {
		setClauses = append(setClauses, fmt.Sprintf("price = $%d", argIndex))
		args = append(args, *input.Price)
		argIndex++
	}

	if input.StockQuantity != nil {
		setClauses = append(setClauses, fmt.Sprintf("stock_quantity = $%d", argIndex))
		args = append(args, *input.StockQuantity)
		argIndex++

		// Auto-update stock_status based on quantity
		if *input.StockQuantity > 0 {
			setClauses = append(setClauses, "stock_status = 'in_stock'")
		} else {
			setClauses = append(setClauses, "stock_status = 'out_of_stock'")
		}
	}

	if input.StockStatus != nil {
		setClauses = append(setClauses, fmt.Sprintf("stock_status = $%d", argIndex))
		args = append(args, *input.StockStatus)
		argIndex++
	}

	if input.IsActive != nil {
		setClauses = append(setClauses, fmt.Sprintf("is_active = $%d", argIndex))
		args = append(args, *input.IsActive)
		argIndex++
	}

	if input.Attributes != nil {
		attributesJSON, err := json.Marshal(input.Attributes)
		if err != nil {
			r.logger.Error().Err(err).Msg("failed to marshal product attributes")
			return nil, fmt.Errorf("failed to marshal product attributes: %w", err)
		}
		setClauses = append(setClauses, fmt.Sprintf("attributes = $%d", argIndex))
		args = append(args, attributesJSON)
		argIndex++
	}

	if input.HasIndividualLocation != nil {
		setClauses = append(setClauses, fmt.Sprintf("has_individual_location = $%d", argIndex))
		args = append(args, *input.HasIndividualLocation)
		argIndex++
	}

	if input.IndividualAddress != nil {
		setClauses = append(setClauses, fmt.Sprintf("individual_address = $%d", argIndex))
		args = append(args, *input.IndividualAddress)
		argIndex++
	}

	if input.IndividualLatitude != nil {
		setClauses = append(setClauses, fmt.Sprintf("individual_latitude = $%d", argIndex))
		args = append(args, *input.IndividualLatitude)
		argIndex++
	}

	if input.IndividualLongitude != nil {
		setClauses = append(setClauses, fmt.Sprintf("individual_longitude = $%d", argIndex))
		args = append(args, *input.IndividualLongitude)
		argIndex++
	}

	if input.LocationPrivacy != nil {
		setClauses = append(setClauses, fmt.Sprintf("location_privacy = $%d", argIndex))
		args = append(args, *input.LocationPrivacy)
		argIndex++
	}

	if input.ShowOnMap != nil {
		setClauses = append(setClauses, fmt.Sprintf("show_on_map = $%d", argIndex))
		args = append(args, *input.ShowOnMap)
		argIndex++
	}

	// If no fields to update, return error
	if len(setClauses) == 1 { // Only updated_at
		return nil, fmt.Errorf("no fields to update")
	}

	// Add WHERE conditions
	args = append(args, productID, storefrontID)
	whereProductID := argIndex
	whereStorefrontID := argIndex + 1

	// Build final query
	query := fmt.Sprintf(`
		UPDATE b2c_products
		SET %s
		WHERE id = $%d AND storefront_id = $%d
		RETURNING
			id, storefront_id, name, description, price, currency,
			category_id, sku, barcode, stock_quantity, stock_status,
			is_active, attributes, view_count, sold_count,
			created_at, updated_at,
			has_individual_location, individual_address,
			individual_latitude, individual_longitude,
			location_privacy, show_on_map, has_variants
	`, strings.Join(setClauses, ", "), whereProductID, whereStorefrontID)

	var product domain.Product
	var description sql.NullString
	var sku, barcode sql.NullString
	var individualAddress, locationPrivacy sql.NullString
	var individualLatitude, individualLongitude sql.NullFloat64
	var returnedAttributesJSON []byte

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&product.ID,
		&product.StorefrontID,
		&product.Name,
		&description,
		&product.Price,
		&product.Currency,
		&product.CategoryID,
		&sku,
		&barcode,
		&product.StockQuantity,
		&product.StockStatus,
		&product.IsActive,
		&returnedAttributesJSON,
		&product.ViewCount,
		&product.SoldCount,
		&product.CreatedAt,
		&product.UpdatedAt,
		&product.HasIndividualLocation,
		&individualAddress,
		&individualLatitude,
		&individualLongitude,
		&locationPrivacy,
		&product.ShowOnMap,
		&product.HasVariants,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Error().Int64("product_id", productID).Int64("storefront_id", storefrontID).Msg("product not found or ownership validation failed")
			return nil, fmt.Errorf("products.not_found")
		}
		// Check for unique constraint violation (duplicate SKU)
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // unique_violation
				r.logger.Error().Err(err).Msg("duplicate SKU")
				return nil, fmt.Errorf("products.sku_duplicate")
			}
		}
		r.logger.Error().Err(err).Int64("product_id", productID).Msg("failed to update product")
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	// Handle nullable fields
	if description.Valid {
		product.Description = description.String
	}
	if sku.Valid {
		product.SKU = &sku.String
	}
	if barcode.Valid {
		product.Barcode = &barcode.String
	}
	if individualAddress.Valid {
		product.IndividualAddress = &individualAddress.String
	}
	if individualLatitude.Valid {
		product.IndividualLatitude = &individualLatitude.Float64
	}
	if individualLongitude.Valid {
		product.IndividualLongitude = &individualLongitude.Float64
	}
	if locationPrivacy.Valid {
		product.LocationPrivacy = &locationPrivacy.String
	}

	// Parse returned JSONB attributes
	if len(returnedAttributesJSON) > 0 {
		if err := json.Unmarshal(returnedAttributesJSON, &product.Attributes); err != nil {
			r.logger.Error().Err(err).Msg("failed to unmarshal returned product attributes")
			return nil, fmt.Errorf("failed to unmarshal product attributes: %w", err)
		}
	}

	r.logger.Info().Int64("product_id", product.ID).Msg("product updated successfully")
	return &product, nil
}

// DeleteProduct deletes a product (soft or hard delete) with ownership validation
func (r *Repository) DeleteProduct(ctx context.Context, productID, storefrontID int64, hardDelete bool) (int32, error) {
	r.logger.Debug().
		Int64("product_id", productID).
		Int64("storefront_id", storefrontID).
		Bool("hard_delete", hardDelete).
		Msg("deleting product")

	// Validate inputs
	if productID <= 0 {
		return 0, fmt.Errorf("product_id must be greater than 0")
	}
	if storefrontID <= 0 {
		return 0, fmt.Errorf("storefront_id must be greater than 0")
	}

	// Start transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to begin transaction")
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Step 1: Check ownership
	var exists bool
	ownershipQuery := `
		SELECT EXISTS(
			SELECT 1 FROM b2c_products
			WHERE id = $1 AND storefront_id = $2
		)
	`
	err = tx.QueryRowContext(ctx, ownershipQuery, productID, storefrontID).Scan(&exists)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to check product ownership")
		return 0, fmt.Errorf("failed to check product ownership: %w", err)
	}

	if !exists {
		return 0, fmt.Errorf("products.not_found")
	}

	// Step 2: Check for active orders (prevent deletion if product has pending orders)
	// TODO: Add check for active orders once orders table/microservice is available
	// For now, we skip this check

	// Step 3: Count variants before deletion
	var variantsCount int32
	countQuery := `
		SELECT COUNT(*) FROM b2c_product_variants
		WHERE product_id = $1
	`
	err = tx.QueryRowContext(ctx, countQuery, productID).Scan(&variantsCount)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to count variants")
		return 0, fmt.Errorf("failed to count variants: %w", err)
	}

	if hardDelete {
		// Hard delete: DELETE CASCADE will handle variants automatically
		deleteQuery := `
			DELETE FROM b2c_products
			WHERE id = $1 AND storefront_id = $2
		`
		result, err := tx.ExecContext(ctx, deleteQuery, productID, storefrontID)
		if err != nil {
			r.logger.Error().Err(err).Msg("failed to hard delete product")
			return 0, fmt.Errorf("failed to delete product: %w", err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			r.logger.Error().Err(err).Msg("failed to get rows affected")
			return 0, fmt.Errorf("failed to get rows affected: %w", err)
		}

		if rowsAffected == 0 {
			return 0, fmt.Errorf("products.not_found")
		}

		r.logger.Info().
			Int64("product_id", productID).
			Int32("variants_deleted", variantsCount).
			Msg("product hard deleted successfully")
	} else {
		// Soft delete: Set deleted_at timestamp
		// Note: If deleted_at column doesn't exist, we'll need to add it via migration
		softDeleteQuery := `
			UPDATE b2c_products
			SET deleted_at = NOW(), updated_at = NOW()
			WHERE id = $1 AND storefront_id = $2 AND deleted_at IS NULL
		`
		result, err := tx.ExecContext(ctx, softDeleteQuery, productID, storefrontID)
		if err != nil {
			// Check if column doesn't exist
			if pqErr, ok := err.(*pq.Error); ok {
				if pqErr.Code == "42703" { // undefined_column
					r.logger.Warn().Msg("deleted_at column not found, using is_active=false instead")
					// Fallback: use is_active = false
					fallbackQuery := `
						UPDATE b2c_products
						SET is_active = false, updated_at = NOW()
						WHERE id = $1 AND storefront_id = $2
					`
					result, err = tx.ExecContext(ctx, fallbackQuery, productID, storefrontID)
					if err != nil {
						r.logger.Error().Err(err).Msg("failed to deactivate product")
						return 0, fmt.Errorf("failed to deactivate product: %w", err)
					}
				} else {
					r.logger.Error().Err(err).Msg("failed to soft delete product")
					return 0, fmt.Errorf("failed to soft delete product: %w", err)
				}
			} else {
				r.logger.Error().Err(err).Msg("failed to soft delete product")
				return 0, fmt.Errorf("failed to soft delete product: %w", err)
			}
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			r.logger.Error().Err(err).Msg("failed to get rows affected")
			return 0, fmt.Errorf("failed to get rows affected: %w", err)
		}

		if rowsAffected == 0 {
			return 0, fmt.Errorf("products.not_found")
		}

		r.logger.Info().
			Int64("product_id", productID).
			Int32("variants_count", variantsCount).
			Msg("product soft deleted successfully")
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		r.logger.Error().Err(err).Msg("failed to commit transaction")
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return variantsCount, nil
}

// BulkCreateProducts creates multiple products in a single atomic transaction
func (r *Repository) BulkCreateProducts(ctx context.Context, storefrontID int64, inputs []*domain.CreateProductInput) ([]*domain.Product, []domain.BulkProductError, error) {
	r.logger.Debug().
		Int64("storefront_id", storefrontID).
		Int("product_count", len(inputs)).
		Msg("bulk creating products")

	// Validate inputs
	if storefrontID <= 0 {
		return nil, nil, fmt.Errorf("storefront_id must be greater than 0")
	}
	if len(inputs) == 0 {
		return nil, nil, fmt.Errorf("products.bulk_empty")
	}
	if len(inputs) > 1000 {
		return nil, nil, fmt.Errorf("products.bulk_too_large")
	}

	// Start transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to begin transaction")
		return nil, nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Collect all SKUs for uniqueness validation
	var skusToCheck []string
	for _, input := range inputs {
		if input.SKU != nil && *input.SKU != "" {
			skusToCheck = append(skusToCheck, *input.SKU)
		}
	}

	// Check for duplicate SKUs within storefront
	if len(skusToCheck) > 0 {
		duplicateQuery := `
			SELECT sku FROM b2c_products
			WHERE storefront_id = $1 AND sku = ANY($2::text[])
		`
		rows, err := tx.QueryContext(ctx, duplicateQuery, storefrontID, pq.Array(skusToCheck))
		if err != nil {
			r.logger.Error().Err(err).Msg("failed to check duplicate SKUs")
			return nil, nil, fmt.Errorf("failed to check duplicate SKUs: %w", err)
		}
		defer rows.Close()

		duplicateSKUs := make(map[string]bool)
		for rows.Next() {
			var sku string
			if err := rows.Scan(&sku); err != nil {
				return nil, nil, fmt.Errorf("failed to scan duplicate SKU: %w", err)
			}
			duplicateSKUs[sku] = true
		}

		// If duplicates found, return errors for all items with duplicate SKUs
		if len(duplicateSKUs) > 0 {
			var errors []domain.BulkProductError
			for idx, input := range inputs {
				if input.SKU != nil && duplicateSKUs[*input.SKU] {
					errors = append(errors, domain.BulkProductError{
						Index:        int32(idx),
						ErrorCode:    "products.sku_duplicate",
						ErrorMessage: fmt.Sprintf("SKU %s already exists in storefront", *input.SKU),
					})
				}
			}
			return nil, errors, fmt.Errorf("products.sku_duplicate")
		}
	}

	// Prepare batch insert using VALUES clause
	var createdProducts []*domain.Product
	var errors []domain.BulkProductError

	// Insert products one by one (can be optimized later with true batch INSERT)
	for idx, input := range inputs {
		// Validate individual input
		if input.Name == "" {
			errors = append(errors, domain.BulkProductError{
				Index:        int32(idx),
				ErrorCode:    "products.validation_failed",
				ErrorMessage: "product name cannot be empty",
			})
			continue
		}
		if input.Price < 0 {
			errors = append(errors, domain.BulkProductError{
				Index:        int32(idx),
				ErrorCode:    "products.validation_failed",
				ErrorMessage: "price must be non-negative",
			})
			continue
		}
		if input.StockQuantity < 0 {
			errors = append(errors, domain.BulkProductError{
				Index:        int32(idx),
				ErrorCode:    "products.validation_failed",
				ErrorMessage: "stock_quantity must be non-negative",
			})
			continue
		}

		// Marshal attributes
		var attributesJSON []byte
		if input.Attributes != nil {
			attributesJSON, err = json.Marshal(input.Attributes)
			if err != nil {
				errors = append(errors, domain.BulkProductError{
					Index:        int32(idx),
					ErrorCode:    "products.validation_failed",
					ErrorMessage: fmt.Sprintf("failed to marshal attributes: %v", err),
				})
				continue
			}
		} else {
			// Default to empty JSON object if nil
			attributesJSON = []byte("{}")
		}

		// Determine stock_status
		stockStatus := domain.StockStatusOutOfStock
		if input.StockQuantity > 0 {
			stockStatus = domain.StockStatusInStock
		}

		// Insert product
		query := `
			INSERT INTO b2c_products (
				storefront_id, name, description, price, currency, category_id,
				sku, barcode, stock_quantity, stock_status, is_active,
				attributes, view_count, sold_count,
				has_individual_location, individual_address,
				individual_latitude, individual_longitude,
				location_privacy, show_on_map, has_variants
			) VALUES (
				$1, $2, $3, $4, $5, $6,
				$7, $8, $9, $10, true,
				$11, 0, 0,
				$12, $13, $14, $15, $16, $17, false
			)
			RETURNING
				id, storefront_id, name, description, price, currency,
				category_id, sku, barcode, stock_quantity, stock_status,
				is_active, attributes, view_count, sold_count,
				created_at, updated_at,
				has_individual_location, individual_address,
				individual_latitude, individual_longitude,
				location_privacy, show_on_map, has_variants
		`

		var product domain.Product
		var description sql.NullString
		var sku, barcode sql.NullString
		var individualAddress, locationPrivacy sql.NullString
		var individualLatitude, individualLongitude sql.NullFloat64
		var returnedAttributesJSON []byte

		err = tx.QueryRowContext(
			ctx,
			query,
			storefrontID,
			input.Name,
			input.Description,
			input.Price,
			input.Currency,
			input.CategoryID,
			input.SKU,
			input.Barcode,
			input.StockQuantity,
			stockStatus,
			attributesJSON,
			input.HasIndividualLocation,
			input.IndividualAddress,
			input.IndividualLatitude,
			input.IndividualLongitude,
			input.LocationPrivacy,
			input.ShowOnMap,
		).Scan(
			&product.ID,
			&product.StorefrontID,
			&product.Name,
			&description,
			&product.Price,
			&product.Currency,
			&product.CategoryID,
			&sku,
			&barcode,
			&product.StockQuantity,
			&product.StockStatus,
			&product.IsActive,
			&returnedAttributesJSON,
			&product.ViewCount,
			&product.SoldCount,
			&product.CreatedAt,
			&product.UpdatedAt,
			&product.HasIndividualLocation,
			&individualAddress,
			&individualLatitude,
			&individualLongitude,
			&locationPrivacy,
			&product.ShowOnMap,
			&product.HasVariants,
		)

		if err != nil {
			// Check for unique constraint violation
			if pqErr, ok := err.(*pq.Error); ok {
				if pqErr.Code == "23505" { // unique_violation
					errors = append(errors, domain.BulkProductError{
						Index:        int32(idx),
						ErrorCode:    "products.sku_duplicate",
						ErrorMessage: fmt.Sprintf("SKU already exists"),
					})
					continue
				}
			}
			errors = append(errors, domain.BulkProductError{
				Index:        int32(idx),
				ErrorCode:    "products.bulk_create_failed",
				ErrorMessage: fmt.Sprintf("failed to create product: %v", err),
			})
			continue
		}

		// Handle nullable fields
		if description.Valid {
			product.Description = description.String
		}
		if sku.Valid {
			product.SKU = &sku.String
		}
		if barcode.Valid {
			product.Barcode = &barcode.String
		}
		if individualAddress.Valid {
			product.IndividualAddress = &individualAddress.String
		}
		if individualLatitude.Valid {
			product.IndividualLatitude = &individualLatitude.Float64
		}
		if individualLongitude.Valid {
			product.IndividualLongitude = &individualLongitude.Float64
		}
		if locationPrivacy.Valid {
			product.LocationPrivacy = &locationPrivacy.String
		}

		// Parse returned JSONB attributes
		if len(returnedAttributesJSON) > 0 {
			if err := json.Unmarshal(returnedAttributesJSON, &product.Attributes); err != nil {
				r.logger.Error().Err(err).Msg("failed to unmarshal returned product attributes")
				// Don't fail the entire operation for this
			}
		}

		createdProducts = append(createdProducts, &product)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		r.logger.Error().Err(err).Msg("failed to commit transaction")
		return nil, nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	r.logger.Info().
		Int("successful", len(createdProducts)).
		Int("failed", len(errors)).
		Msg("bulk product creation completed")

	return createdProducts, errors, nil
}

// BulkDeleteProducts deletes multiple products in a single atomic transaction
func (r *Repository) BulkDeleteProducts(ctx context.Context, storefrontID int64, productIDs []int64, hardDelete bool) (int32, int32, int32, map[int64]string, error) {
	r.logger.Debug().
		Int64("storefront_id", storefrontID).
		Int("product_count", len(productIDs)).
		Bool("hard_delete", hardDelete).
		Msg("bulk deleting products")

	// Validate inputs
	if storefrontID <= 0 {
		return 0, 0, 0, nil, fmt.Errorf("storefront_id must be greater than 0")
	}

	if len(productIDs) == 0 {
		return 0, 0, 0, nil, fmt.Errorf("product_ids list cannot be empty")
	}

	if len(productIDs) > 1000 {
		return 0, 0, 0, nil, fmt.Errorf("cannot delete more than 1000 products at once")
	}

	// Deduplicate product IDs
	uniqueIDs := make(map[int64]bool)
	deduplicatedIDs := make([]int64, 0, len(productIDs))
	for _, id := range productIDs {
		if id > 0 && !uniqueIDs[id] {
			uniqueIDs[id] = true
			deduplicatedIDs = append(deduplicatedIDs, id)
		}
	}

	if len(deduplicatedIDs) == 0 {
		return 0, 0, 0, nil, fmt.Errorf("no valid product IDs provided")
	}

	// Start transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to begin transaction")
		return 0, 0, 0, nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var successCount int32
	var failedCount int32
	var totalVariantsDeleted int32
	errors := make(map[int64]string)

	// Process in batches of 100 for better performance
	batchSize := 100
	for i := 0; i < len(deduplicatedIDs); i += batchSize {
		end := i + batchSize
		if end > len(deduplicatedIDs) {
			end = len(deduplicatedIDs)
		}
		batchIDs := deduplicatedIDs[i:end]

		// Step 1: Validate ownership for batch
		ownershipQuery := `
			SELECT id
			FROM b2c_products
			WHERE id = ANY($1::bigint[]) AND storefront_id = $2
		`
		rows, err := tx.QueryContext(ctx, ownershipQuery, pq.Array(batchIDs), storefrontID)
		if err != nil {
			r.logger.Error().Err(err).Msg("failed to check ownership")
			return 0, 0, 0, nil, fmt.Errorf("failed to check ownership: %w", err)
		}

		validIDs := make([]int64, 0, len(batchIDs))
		validIDsMap := make(map[int64]bool)
		for rows.Next() {
			var id int64
			if err := rows.Scan(&id); err != nil {
				rows.Close()
				r.logger.Error().Err(err).Msg("failed to scan product ID")
				return 0, 0, 0, nil, fmt.Errorf("failed to scan product ID: %w", err)
			}
			validIDs = append(validIDs, id)
			validIDsMap[id] = true
		}
		rows.Close()

		if err := rows.Err(); err != nil {
			r.logger.Error().Err(err).Msg("rows iteration error")
			return 0, 0, 0, nil, fmt.Errorf("rows iteration error: %w", err)
		}

		// Mark invalid IDs as failed (not found or not owned)
		for _, id := range batchIDs {
			if !validIDsMap[id] {
				errors[id] = "products.not_found"
				failedCount++
			}
		}

		// Skip if no valid IDs in this batch
		if len(validIDs) == 0 {
			continue
		}

		// Step 2: Count variants before deletion for valid IDs
		countQuery := `
			SELECT product_id, COUNT(*)
			FROM b2c_product_variants
			WHERE product_id = ANY($1::bigint[])
			GROUP BY product_id
		`
		variantRows, err := tx.QueryContext(ctx, countQuery, pq.Array(validIDs))
		if err != nil {
			r.logger.Error().Err(err).Msg("failed to count variants")
			return 0, 0, 0, nil, fmt.Errorf("failed to count variants: %w", err)
		}

		variantCounts := make(map[int64]int32)
		for variantRows.Next() {
			var productID int64
			var count int32
			if err := variantRows.Scan(&productID, &count); err != nil {
				variantRows.Close()
				r.logger.Error().Err(err).Msg("failed to scan variant count")
				return 0, 0, 0, nil, fmt.Errorf("failed to scan variant count: %w", err)
			}
			variantCounts[productID] = count
		}
		variantRows.Close()

		if err := variantRows.Err(); err != nil {
			r.logger.Error().Err(err).Msg("variant rows iteration error")
			return 0, 0, 0, nil, fmt.Errorf("variant rows iteration error: %w", err)
		}

		// Step 3: Perform deletion (soft or hard)
		if hardDelete {
			// Hard delete: DELETE CASCADE will handle variants automatically
			deleteQuery := `
				DELETE FROM b2c_products
				WHERE id = ANY($1::bigint[]) AND storefront_id = $2
			`
			result, err := tx.ExecContext(ctx, deleteQuery, pq.Array(validIDs), storefrontID)
			if err != nil {
				r.logger.Error().Err(err).Msg("failed to hard delete products")
				// Mark all valid IDs as failed
				for _, id := range validIDs {
					errors[id] = "products.delete_failed"
					failedCount++
				}
				continue
			}

			rowsAffected, err := result.RowsAffected()
			if err != nil {
				r.logger.Error().Err(err).Msg("failed to get rows affected")
				return 0, 0, 0, nil, fmt.Errorf("failed to get rows affected: %w", err)
			}

			successCount += int32(rowsAffected)

			// Sum up variants deleted
			for _, count := range variantCounts {
				totalVariantsDeleted += count
			}

			r.logger.Info().
				Int("batch_size", len(validIDs)).
				Int64("rows_affected", rowsAffected).
				Int32("variants_deleted", totalVariantsDeleted).
				Msg("products hard deleted in batch")
		} else {
			// Soft delete: Set deleted_at timestamp
			softDeleteQuery := `
				UPDATE b2c_products
				SET deleted_at = NOW(), updated_at = NOW()
				WHERE id = ANY($1::bigint[]) AND storefront_id = $2 AND deleted_at IS NULL
			`
			result, err := tx.ExecContext(ctx, softDeleteQuery, pq.Array(validIDs), storefrontID)
			if err != nil {
				// Check if column doesn't exist
				if pqErr, ok := err.(*pq.Error); ok {
					if pqErr.Code == "42703" { // undefined_column
						r.logger.Warn().Msg("deleted_at column not found, using is_active=false instead")
						// Fallback: use is_active = false
						fallbackQuery := `
							UPDATE b2c_products
							SET is_active = false, updated_at = NOW()
							WHERE id = ANY($1::bigint[]) AND storefront_id = $2
						`
						result, err = tx.ExecContext(ctx, fallbackQuery, pq.Array(validIDs), storefrontID)
						if err != nil {
							r.logger.Error().Err(err).Msg("failed to deactivate products")
							// Mark all valid IDs as failed
							for _, id := range validIDs {
								errors[id] = "products.delete_failed"
								failedCount++
							}
							continue
						}
					} else {
						r.logger.Error().Err(err).Msg("failed to soft delete products")
						// Mark all valid IDs as failed
						for _, id := range validIDs {
							errors[id] = "products.delete_failed"
							failedCount++
						}
						continue
					}
				} else {
					r.logger.Error().Err(err).Msg("failed to soft delete products")
					// Mark all valid IDs as failed
					for _, id := range validIDs {
						errors[id] = "products.delete_failed"
						failedCount++
					}
					continue
				}
			}

			rowsAffected, err := result.RowsAffected()
			if err != nil {
				r.logger.Error().Err(err).Msg("failed to get rows affected")
				return 0, 0, 0, nil, fmt.Errorf("failed to get rows affected: %w", err)
			}

			successCount += int32(rowsAffected)

			// Sum up variant counts (soft delete doesn't actually delete variants, but we track them)
			for _, count := range variantCounts {
				totalVariantsDeleted += count
			}

			r.logger.Info().
				Int("batch_size", len(validIDs)).
				Int64("rows_affected", rowsAffected).
				Int32("variants_count", totalVariantsDeleted).
				Msg("products soft deleted in batch")
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		r.logger.Error().Err(err).Msg("failed to commit transaction")
		return 0, 0, 0, nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	r.logger.Info().
		Int32("success_count", successCount).
		Int32("failed_count", failedCount).
		Int32("variants_deleted", totalVariantsDeleted).
		Bool("hard_delete", hardDelete).
		Msg("bulk delete products completed")

	return successCount, failedCount, totalVariantsDeleted, errors, nil
}


// UpdateProductInventory updates product stock with inventory movement tracking
// NOTE: Requires b2c_inventory_movements table (create via migration if not exists)
func (r *Repository) UpdateProductInventory(ctx context.Context, storefrontID, productID, variantID int64, movementType string, quantity int32, reason, notes string, userID int64) (int32, int32, error) {
	r.logger.Debug().
		Int64("storefront_id", storefrontID).
		Int64("product_id", productID).
		Int64("variant_id", variantID).
		Str("movement_type", movementType).
		Int32("quantity", quantity).
		Msg("updating product inventory")

	// Validate inputs
	if storefrontID <= 0 {
		return 0, 0, fmt.Errorf("storefront_id must be greater than 0")
	}
	if productID <= 0 {
		return 0, 0, fmt.Errorf("product_id must be greater than 0")
	}
	if movementType != "in" && movementType != "out" && movementType != "adjustment" {
		return 0, 0, fmt.Errorf("invalid movement_type: must be 'in', 'out', or 'adjustment'")
	}
	if quantity < 0 {
		return 0, 0, fmt.Errorf("quantity cannot be negative")
	}

	// Start transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to begin transaction")
		return 0, 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var currentQuantity, newQuantity int32
	var tableName, idColumn string

	// Determine target table and column
	if variantID > 0 {
		tableName = "b2c_product_variants"
		idColumn = "id"
	} else {
		tableName = "b2c_products"
		idColumn = "id"
	}

	// Get current stock quantity
	var query string
	var queryArgs []interface{}
	if variantID > 0 {
		query = fmt.Sprintf("SELECT stock_quantity FROM %s WHERE %s = $1 AND product_id = $2", tableName, idColumn)
		queryArgs = []interface{}{variantID, productID}
	} else {
		query = fmt.Sprintf("SELECT stock_quantity FROM %s WHERE %s = $1 AND storefront_id = $2", tableName, idColumn)
		queryArgs = []interface{}{productID, storefrontID}
	}

	err = tx.QueryRowContext(ctx, query, queryArgs...).Scan(&currentQuantity)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, 0, fmt.Errorf("products.not_found")
		}
		r.logger.Error().Err(err).Msg("failed to get current stock")
		return 0, 0, fmt.Errorf("failed to get current stock: %w", err)
	}

	// Calculate new quantity based on movement type
	switch movementType {
	case "in":
		newQuantity = currentQuantity + quantity
	case "out":
		newQuantity = currentQuantity - quantity
		if newQuantity < 0 {
			return currentQuantity, newQuantity, fmt.Errorf("products.insufficient_stock")
		}
	case "adjustment":
		newQuantity = quantity // Direct set
	default:
		return 0, 0, fmt.Errorf("invalid movement_type: %s", movementType)
	}

	// Update stock quantity
	var updateQuery string
	if variantID > 0 {
		updateQuery = fmt.Sprintf("UPDATE %s SET stock_quantity = $1, updated_at = NOW() WHERE %s = $2 AND product_id = $3", tableName, idColumn)
		_, err = tx.ExecContext(ctx, updateQuery, newQuantity, variantID, productID)
	} else {
		updateQuery = fmt.Sprintf("UPDATE %s SET stock_quantity = $1, updated_at = NOW() WHERE %s = $2 AND storefront_id = $3", tableName, idColumn)
		_, err = tx.ExecContext(ctx, updateQuery, newQuantity, productID, storefrontID)
	}

	if err != nil {
		r.logger.Error().Err(err).Msg("failed to update stock quantity")
		return currentQuantity, newQuantity, fmt.Errorf("failed to update stock quantity: %w", err)
	}

	// Record inventory movement in audit trail table
	// Note: This table might not exist yet, so we'll use INSERT with error handling
	movementQuery := `
		INSERT INTO b2c_inventory_movements (
			storefront_product_id, variant_id, type, quantity, reason, notes, user_id, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
	`
	var variantIDPtr *int64
	if variantID > 0 {
		variantIDPtr = &variantID
	}

	_, err = tx.ExecContext(ctx, movementQuery, productID, variantIDPtr, movementType, quantity, reason, notes, userID)
	if err != nil {
		// Check if table doesn't exist
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "42P01" { // undefined_table
				r.logger.Warn().Msg("b2c_inventory_movements table not found, skipping movement recording")
				// Continue without failing - movement tracking is optional for now
			} else {
				r.logger.Error().Err(err).Msg("failed to record inventory movement")
				return currentQuantity, newQuantity, fmt.Errorf("failed to record inventory movement: %w", err)
			}
		} else {
			r.logger.Error().Err(err).Msg("failed to record inventory movement")
			return currentQuantity, newQuantity, fmt.Errorf("failed to record inventory movement: %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		r.logger.Error().Err(err).Msg("failed to commit transaction")
		return currentQuantity, newQuantity, fmt.Errorf("failed to commit transaction: %w", err)
	}

	r.logger.Info().
		Int64("product_id", productID).
		Int64("variant_id", variantID).
		Int32("stock_before", currentQuantity).
		Int32("stock_after", newQuantity).
		Msg("inventory updated successfully")

	return currentQuantity, newQuantity, nil
}

// GetProductStats retrieves statistics for storefront products
func (r *Repository) GetProductStats(ctx context.Context, storefrontID int64) (*domain.ProductStats, error) {
	r.logger.Debug().Int64("storefront_id", storefrontID).Msg("getting product stats")

	if storefrontID <= 0 {
		return nil, fmt.Errorf("storefront_id must be greater than 0")
	}

	query := `
		SELECT
			COUNT(*) as total_products,
			COUNT(*) FILTER (WHERE is_active = true) as active_products,
			COUNT(*) FILTER (WHERE stock_status = 'out_of_stock') as out_of_stock,
			COUNT(*) FILTER (WHERE stock_status = 'low_stock') as low_stock,
			COALESCE(SUM(price * stock_quantity), 0) as total_value,
			COALESCE(SUM(sold_count), 0) as total_sold
		FROM b2c_products
		WHERE storefront_id = $1
	`

	var stats domain.ProductStats
	var totalValue sql.NullFloat64
	var totalSold sql.NullInt32

	err := r.db.QueryRowContext(ctx, query, storefrontID).Scan(
		&stats.TotalProducts,
		&stats.ActiveProducts,
		&stats.OutOfStock,
		&stats.LowStock,
		&totalValue,
		&totalSold,
	)

	if err != nil {
		r.logger.Error().Err(err).Msg("failed to get product stats")
		return nil, fmt.Errorf("failed to get product stats: %w", err)
	}

	// Handle nullable fields
	if totalValue.Valid {
		stats.TotalValue = totalValue.Float64
	}
	if totalSold.Valid {
		stats.TotalSold = totalSold.Int32
	}

	r.logger.Info().
		Int32("total_products", stats.TotalProducts).
		Int32("active_products", stats.ActiveProducts).
		Float64("total_value", stats.TotalValue).
		Msg("product stats retrieved")

	return &stats, nil
}

// IncrementProductViews increments the view counter for a product
func (r *Repository) IncrementProductViews(ctx context.Context, productID int64) error {
	r.logger.Debug().Int64("product_id", productID).Msg("incrementing product views")

	if productID <= 0 {
		return fmt.Errorf("product_id must be greater than 0")
	}

	query := `
		UPDATE b2c_products
		SET view_count = view_count + 1, updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, productID)
	if err != nil {
		r.logger.Error().Err(err).Int64("product_id", productID).Msg("failed to increment product views")
		return fmt.Errorf("failed to increment product views: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("products.not_found")
	}

	r.logger.Debug().Int64("product_id", productID).Msg("product views incremented")
	return nil
}

// BatchUpdateStock updates stock for multiple products/variants atomically
func (r *Repository) BatchUpdateStock(ctx context.Context, storefrontID int64, items []domain.StockUpdateItem, reason string, userID int64) (int32, int32, []domain.StockUpdateResult, error) {
	r.logger.Debug().
		Int64("storefront_id", storefrontID).
		Int("item_count", len(items)).
		Msg("batch updating stock")

	// Validate inputs
	if storefrontID <= 0 {
		return 0, 0, nil, fmt.Errorf("storefront_id must be greater than 0")
	}
	if len(items) == 0 {
		return 0, 0, nil, fmt.Errorf("items list cannot be empty")
	}
	if len(items) > 1000 {
		return 0, 0, nil, fmt.Errorf("cannot update more than 1000 items at once")
	}

	// Start transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to begin transaction")
		return 0, 0, nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var successCount, failedCount int32
	var results []domain.StockUpdateResult

	// Process each item
	for _, item := range items {
		result := domain.StockUpdateResult{
			ProductID: item.ProductID,
			VariantID: item.VariantID,
		}

		var currentQuantity int32
		var tableName, idColumn string
		var query string
		var queryArgs []interface{}

		// Determine target table
		if item.VariantID != nil && *item.VariantID > 0 {
			tableName = "b2c_product_variants"
			idColumn = "id"
			query = fmt.Sprintf("SELECT stock_quantity FROM %s WHERE %s = $1 AND product_id = $2", tableName, idColumn)
			queryArgs = []interface{}{*item.VariantID, item.ProductID}
			result.VariantID = item.VariantID
		} else {
			tableName = "b2c_products"
			idColumn = "id"
			query = fmt.Sprintf("SELECT stock_quantity FROM %s WHERE %s = $1 AND storefront_id = $2", tableName, idColumn)
			queryArgs = []interface{}{item.ProductID, storefrontID}
		}

		// Get current stock
		err := tx.QueryRowContext(ctx, query, queryArgs...).Scan(&currentQuantity)
		if err != nil {
			if err == sql.ErrNoRows {
				result.Success = false
				result.Error = strPtr("not found")
				failedCount++
				results = append(results, result)
				continue
			}
			r.logger.Error().Err(err).Msg("failed to get current stock")
			result.Success = false
			result.Error = strPtr(fmt.Sprintf("failed to get stock: %v", err))
			failedCount++
			results = append(results, result)
			continue
		}

		result.StockBefore = currentQuantity

		// Update stock quantity (absolute set)
		var updateQuery string
		if item.VariantID != nil && *item.VariantID > 0 {
			updateQuery = fmt.Sprintf("UPDATE %s SET stock_quantity = $1, updated_at = NOW() WHERE %s = $2 AND product_id = $3", tableName, idColumn)
			_, err = tx.ExecContext(ctx, updateQuery, item.Quantity, *item.VariantID, item.ProductID)
		} else {
			updateQuery = fmt.Sprintf("UPDATE %s SET stock_quantity = $1, updated_at = NOW() WHERE %s = $2 AND storefront_id = $3", tableName, idColumn)
			_, err = tx.ExecContext(ctx, updateQuery, item.Quantity, item.ProductID, storefrontID)
		}

		if err != nil {
			r.logger.Error().Err(err).Msg("failed to update stock")
			result.Success = false
			result.Error = strPtr(fmt.Sprintf("failed to update: %v", err))
			failedCount++
			results = append(results, result)
			continue
		}

		result.StockAfter = item.Quantity
		result.Success = true
		successCount++
		results = append(results, result)

		// Record movement (optional, don't fail if table doesn't exist)
		movementReason := reason
		if item.Reason != nil && *item.Reason != "" {
			movementReason = *item.Reason
		}

		movementQuery := `
			INSERT INTO b2c_inventory_movements (
				storefront_product_id, variant_id, type, quantity, reason, notes, user_id, created_at
			) VALUES ($1, $2, 'adjustment', $3, $4, '', $5, NOW())
		`
		_, movErr := tx.ExecContext(ctx, movementQuery, item.ProductID, item.VariantID, item.Quantity, movementReason, userID)
		if movErr != nil {
			// Log warning but don't fail - movement tracking is optional
			if pqErr, ok := movErr.(*pq.Error); ok {
				if pqErr.Code != "42P01" { // Only warn if not "table doesn't exist"
					r.logger.Warn().Err(movErr).Msg("failed to record inventory movement, continuing")
				}
			}
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		r.logger.Error().Err(err).Msg("failed to commit transaction")
		return 0, 0, nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	r.logger.Info().
		Int32("success_count", successCount).
		Int32("failed_count", failedCount).
		Msg("batch stock update completed")

	return successCount, failedCount, results, nil
}

// Helper function to create string pointer
func strPtr(s string) *string {
	return &s
}

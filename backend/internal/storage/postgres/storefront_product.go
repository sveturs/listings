package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"backend/internal/domain/models"
	"backend/internal/logger"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// ErrStorefrontProductNotFound возвращается когда товар витрины не найден
var ErrStorefrontProductNotFound = errors.New("storefront product not found")

// Sort field constants
const (
	sortFieldName          = "name"
	sortFieldPrice         = "price"
	sortFieldStockQuantity = "stock_quantity"
	sortFieldCreatedAt     = "created_at"
)

// GetBySlug retrieves a storefront by slug
func (s *Database) GetBySlug(ctx context.Context, slug string) (*models.Storefront, error) {
	return s.storefrontRepo.GetBySlug(ctx, slug)
}

// GetStorefrontProducts retrieves products for a storefront with filters
func (s *Database) GetStorefrontProducts(ctx context.Context, filter models.ProductFilter) ([]*models.StorefrontProduct, error) {
	query := `
		SELECT
			p.id, p.storefront_id, p.name, p.description, p.price, p.currency,
			p.category_id, p.sku, p.barcode, p.stock_quantity, p.stock_status,
			p.is_active, p.attributes, p.view_count, p.sold_count,
			p.created_at, p.updated_at,
			p.has_individual_location, p.individual_address, p.individual_latitude,
			p.individual_longitude, p.location_privacy, p.show_on_map, p.has_variants,
			c.id, c.name, c.slug, c.icon, c.parent_id
		FROM storefront_products p
		LEFT JOIN marketplace_categories c ON p.category_id = c.id
		WHERE p.storefront_id = $1`

	args := []interface{}{filter.StorefrontID}
	argIndex := 2

	// Apply filters
	// Фильтр по активности - применяем только если явно указан
	if filter.IsActive != nil {
		query += fmt.Sprintf(" AND p.is_active = $%d", argIndex)
		args = append(args, *filter.IsActive)
		argIndex++
	}

	if filter.CategoryID != nil {
		query += fmt.Sprintf(" AND p.category_id = $%d", argIndex)
		args = append(args, *filter.CategoryID)
		argIndex++
	}

	if filter.Search != nil && *filter.Search != "" {
		query += fmt.Sprintf(" AND (p.name ILIKE $%d OR p.description ILIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+*filter.Search+"%")
		argIndex++
	}

	if filter.MinPrice != nil {
		query += fmt.Sprintf(" AND p.price >= $%d", argIndex)
		args = append(args, *filter.MinPrice)
		argIndex++
	}

	if filter.MaxPrice != nil {
		query += fmt.Sprintf(" AND p.price <= $%d", argIndex)
		args = append(args, *filter.MaxPrice)
		argIndex++
	}

	if filter.StockStatus != nil {
		query += fmt.Sprintf(" AND p.stock_status = $%d", argIndex)
		args = append(args, *filter.StockStatus)
		argIndex++
	}

	if filter.SKU != nil {
		query += fmt.Sprintf(" AND p.sku = $%d", argIndex)
		args = append(args, *filter.SKU)
		argIndex++
	}

	if filter.Barcode != nil {
		query += fmt.Sprintf(" AND p.barcode = $%d", argIndex)
		args = append(args, *filter.Barcode)
		argIndex++
	}

	// Apply sorting
	sortBy := "p.created_at"
	if filter.SortBy != "" {
		switch filter.SortBy {
		case sortFieldName:
			sortBy = "p.name"
		case sortFieldPrice:
			sortBy = "p.price"
		case sortFieldStockQuantity:
			sortBy = "p.stock_quantity"
		case sortFieldCreatedAt:
			sortBy = "p.created_at"
		}
	}

	sortOrder := "DESC"
	if filter.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)

	// Apply pagination
	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, filter.Limit)
		argIndex++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, filter.Offset)
	}

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get storefront products: %w", err)
	}
	defer rows.Close()

	var products []*models.StorefrontProduct
	for rows.Next() {
		p := &models.StorefrontProduct{}
		c := &models.MarketplaceCategory{}
		var attributesJSON []byte
		var categoryID sql.NullInt64
		var categoryName, categorySlug, categoryIcon sql.NullString
		var categoryParentID sql.NullInt64

		err := rows.Scan(
			&p.ID, &p.StorefrontID, &p.Name, &p.Description, &p.Price, &p.Currency,
			&p.CategoryID, &p.SKU, &p.Barcode, &p.StockQuantity, &p.StockStatus,
			&p.IsActive, &attributesJSON, &p.ViewCount, &p.SoldCount,
			&p.CreatedAt, &p.UpdatedAt,
			&p.HasIndividualLocation, &p.IndividualAddress, &p.IndividualLatitude,
			&p.IndividualLongitude, &p.LocationPrivacy, &p.ShowOnMap, &p.HasVariants,
			&categoryID, &categoryName, &categorySlug, &categoryIcon, &categoryParentID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}

		// Only set category if it exists
		if categoryID.Valid {
			c.ID = int(categoryID.Int64)
			c.Name = categoryName.String
			c.Slug = categorySlug.String
			if categoryIcon.Valid {
				c.Icon = &categoryIcon.String
			}
			if categoryParentID.Valid {
				parentID := int(categoryParentID.Int64)
				c.ParentID = &parentID
			}
			p.Category = c
		}

		if attributesJSON != nil {
			if err := json.Unmarshal(attributesJSON, &p.Attributes); err != nil {
				return nil, fmt.Errorf("failed to unmarshal attributes: %w", err)
			}
		}

		p.Category = c
		products = append(products, p)
	}

	// Load images for all products
	if len(products) > 0 {
		productIDs := make([]int, len(products))
		productMap := make(map[int]*models.StorefrontProduct)
		for i, p := range products {
			productIDs[i] = p.ID
			productMap[p.ID] = p
			// Initialize empty slice for images
			p.Images = []models.StorefrontProductImage{}
		}

		images, err := s.getProductImages(ctx, productIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to get product images: %w", err)
		}

		for _, img := range images {
			if p, ok := productMap[img.StorefrontProductID]; ok {
				p.Images = append(p.Images, img)
			}
		}

		// Load address translations from geocoding_cache
		if err := s.loadAddressTranslations(ctx, products); err != nil {
			logger.Warn().Err(err).Msg("Failed to load address translations")
			// Не фатальная ошибка, продолжаем без переводов
		}
	}

	return products, nil
}

// GetStorefrontProduct retrieves a single product by ID
func (s *Database) GetStorefrontProduct(ctx context.Context, storefrontID, productID int) (*models.StorefrontProduct, error) {
	query := `
		SELECT
			p.id, p.storefront_id, p.name, p.description, p.price, p.currency,
			p.category_id, p.sku, p.barcode, p.stock_quantity, p.stock_status,
			p.is_active, p.attributes, p.view_count, p.sold_count,
			p.created_at, p.updated_at,
			p.has_individual_location, p.individual_address, p.individual_latitude,
			p.individual_longitude, p.location_privacy, p.show_on_map, p.has_variants,
			c.id, c.name, c.slug, c.icon, c.parent_id
		FROM storefront_products p
		LEFT JOIN marketplace_categories c ON p.category_id = c.id
		WHERE p.id = $1 AND p.storefront_id = $2`

	p := &models.StorefrontProduct{}
	c := &models.MarketplaceCategory{}
	var attributesJSON []byte
	var categoryID sql.NullInt64
	var categoryName, categorySlug, categoryIcon sql.NullString
	var categoryParentID sql.NullInt64

	err := s.pool.QueryRow(ctx, query, productID, storefrontID).Scan(
		&p.ID, &p.StorefrontID, &p.Name, &p.Description, &p.Price, &p.Currency,
		&p.CategoryID, &p.SKU, &p.Barcode, &p.StockQuantity, &p.StockStatus,
		&p.IsActive, &attributesJSON, &p.ViewCount, &p.SoldCount,
		&p.CreatedAt, &p.UpdatedAt,
		&p.HasIndividualLocation, &p.IndividualAddress, &p.IndividualLatitude,
		&p.IndividualLongitude, &p.LocationPrivacy, &p.ShowOnMap, &p.HasVariants,
		&categoryID, &categoryName, &categorySlug, &categoryIcon, &categoryParentID,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrStorefrontProductNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get storefront product: %w", err)
	}

	if attributesJSON != nil {
		if unmarshalErr := json.Unmarshal(attributesJSON, &p.Attributes); unmarshalErr != nil {
			return nil, fmt.Errorf("failed to unmarshal attributes: %w", unmarshalErr)
		}
	}

	// Only set category if it exists
	if categoryID.Valid {
		c.ID = int(categoryID.Int64)
		c.Name = categoryName.String
		c.Slug = categorySlug.String
		if categoryIcon.Valid {
			c.Icon = &categoryIcon.String
		}
		if categoryParentID.Valid {
			parentID := int(categoryParentID.Int64)
			c.ParentID = &parentID
		}
		p.Category = c
	}

	// Load images
	images, err := s.getProductImages(ctx, []int{p.ID})
	if err != nil {
		return nil, fmt.Errorf("failed to get product images: %w", err)
	}
	// Initialize empty slice if no images found
	if images == nil {
		p.Images = []models.StorefrontProductImage{}
	} else {
		p.Images = images
	}

	// Load variants
	variants, err := s.getProductVariants(ctx, p.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product variants: %w", err)
	}
	p.Variants = variants

	// Load translations
	translations, err := s.GetTranslationsForEntity(ctx, "storefront_product", p.ID)
	if err != nil {
		// Don't fail if translations load fails, just log it
		logger.Error().Err(err).Int("product_id", p.ID).Msg("Failed to load product translations")
	} else if len(translations) > 0 {
		// Convert []Translation to map[language]map[field]text
		p.Translations = make(map[string]map[string]string)
		for _, t := range translations {
			if p.Translations[t.Language] == nil {
				p.Translations[t.Language] = make(map[string]string)
			}
			p.Translations[t.Language][t.FieldName] = t.TranslatedText
		}
	}

	return p, nil
}

// GetStorefrontProductBySKU retrieves a single product by SKU and storefront_id
func (s *Database) GetStorefrontProductBySKU(ctx context.Context, storefrontID int, sku string) (*models.StorefrontProduct, error) {
	query := `
		SELECT
			p.id, p.storefront_id, p.name, p.description, p.price, p.currency,
			p.category_id, p.sku, p.barcode, p.stock_quantity, p.stock_status,
			p.is_active, p.attributes, p.view_count, p.sold_count,
			p.created_at, p.updated_at,
			p.has_individual_location, p.individual_address, p.individual_latitude,
			p.individual_longitude, p.location_privacy, p.show_on_map, p.has_variants,
			c.id, c.name, c.slug, c.icon, c.parent_id
		FROM storefront_products p
		LEFT JOIN marketplace_categories c ON p.category_id = c.id
		WHERE p.storefront_id = $1 AND p.sku = $2`

	p := &models.StorefrontProduct{}
	c := &models.MarketplaceCategory{}
	var attributesJSON []byte
	var categoryID sql.NullInt64
	var categoryName, categorySlug, categoryIcon sql.NullString
	var categoryParentID sql.NullInt64

	err := s.pool.QueryRow(ctx, query, storefrontID, sku).Scan(
		&p.ID, &p.StorefrontID, &p.Name, &p.Description, &p.Price, &p.Currency,
		&p.CategoryID, &p.SKU, &p.Barcode, &p.StockQuantity, &p.StockStatus,
		&p.IsActive, &attributesJSON, &p.ViewCount, &p.SoldCount,
		&p.CreatedAt, &p.UpdatedAt,
		&p.HasIndividualLocation, &p.IndividualAddress, &p.IndividualLatitude,
		&p.IndividualLongitude, &p.LocationPrivacy, &p.ShowOnMap, &p.HasVariants,
		&categoryID, &categoryName, &categorySlug, &categoryIcon, &categoryParentID,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrStorefrontProductNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get storefront product by SKU: %w", err)
	}

	if attributesJSON != nil {
		if unmarshalErr := json.Unmarshal(attributesJSON, &p.Attributes); unmarshalErr != nil {
			return nil, fmt.Errorf("failed to unmarshal attributes: %w", unmarshalErr)
		}
	}

	// Only set category if it exists
	if categoryID.Valid {
		c.ID = int(categoryID.Int64)
		c.Name = categoryName.String
		c.Slug = categorySlug.String
		if categoryIcon.Valid {
			c.Icon = &categoryIcon.String
		}
		if categoryParentID.Valid {
			parentID := int(categoryParentID.Int64)
			c.ParentID = &parentID
		}
		p.Category = c
	}

	// Load images
	images, err := s.getProductImages(ctx, []int{p.ID})
	if err != nil {
		return nil, fmt.Errorf("failed to get product images: %w", err)
	}
	// Initialize empty slice if no images found
	if images == nil {
		p.Images = []models.StorefrontProductImage{}
	} else {
		p.Images = images
	}

	// Load variants
	variants, err := s.getProductVariants(ctx, p.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product variants: %w", err)
	}
	p.Variants = variants

	// Load translations
	translations, err := s.GetTranslationsForEntity(ctx, "storefront_product", p.ID)
	if err != nil {
		// Don't fail if translations load fails, just log it
		logger.Error().Err(err).Int("product_id", p.ID).Msg("Failed to load product translations")
	} else if len(translations) > 0 {
		// Convert []Translation to map[language]map[field]text
		p.Translations = make(map[string]map[string]string)
		for _, t := range translations {
			if p.Translations[t.Language] == nil {
				p.Translations[t.Language] = make(map[string]string)
			}
			p.Translations[t.Language][t.FieldName] = t.TranslatedText
		}
	}

	// Load address translations from geocoding_cache
	if err := s.loadAddressTranslationsForProduct(ctx, p); err != nil {
		logger.Warn().Err(err).Int("product_id", p.ID).Msg("Failed to load address translations for product")
		// Не фатальная ошибка, продолжаем без переводов
	}

	return p, nil
}

// GetProductByID retrieves a product by ID (for ProductRepositoryInterface)
func (s *Database) GetProductByID(ctx context.Context, productID int64) (*models.StorefrontProduct, error) {
	query := `
		SELECT
			p.id, p.storefront_id, p.name, p.description, p.price, p.currency,
			p.category_id, p.sku, p.barcode, p.stock_quantity, p.stock_status,
			p.is_active, p.attributes, p.view_count, p.sold_count,
			p.created_at, p.updated_at,
			p.has_individual_location, p.individual_address, p.individual_latitude,
			p.individual_longitude, p.location_privacy, p.show_on_map, p.has_variants
		FROM storefront_products p
		WHERE p.id = $1`

	product := &models.StorefrontProduct{}
	var attributesJSON []byte

	err := s.pool.QueryRow(ctx, query, productID).Scan(
		&product.ID, &product.StorefrontID, &product.Name, &product.Description,
		&product.Price, &product.Currency, &product.CategoryID, &product.SKU,
		&product.Barcode, &product.StockQuantity, &product.StockStatus,
		&product.IsActive, &attributesJSON, &product.ViewCount, &product.SoldCount,
		&product.CreatedAt, &product.UpdatedAt,
		&product.HasIndividualLocation, &product.IndividualAddress,
		&product.IndividualLatitude, &product.IndividualLongitude,
		&product.LocationPrivacy, &product.ShowOnMap, &product.HasVariants,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("product not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	if attributesJSON != nil {
		if unmarshalErr := json.Unmarshal(attributesJSON, &product.Attributes); unmarshalErr != nil {
			return nil, fmt.Errorf("failed to unmarshal attributes: %w", unmarshalErr)
		}
	}

	// Load images
	images, err := s.getProductImages(ctx, []int{product.ID})
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get product images")
	}
	if images != nil {
		product.Images = images
	} else {
		product.Images = []models.StorefrontProductImage{}
	}

	// Load variants
	variants, err := s.getProductVariants(ctx, product.ID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get product variants")
	}
	product.Variants = variants

	return product, nil
}

// GetStorefrontProductByID retrieves a storefront product by ID without requiring storefront ID
func (s *Database) GetStorefrontProductByID(ctx context.Context, productID int) (*models.StorefrontProduct, error) {
	query := `
		SELECT
			p.id, p.storefront_id, p.name, p.description, p.price, p.currency,
			p.category_id, p.sku, p.barcode, p.stock_quantity, p.stock_status,
			p.is_active, p.attributes, p.view_count, p.sold_count,
			p.created_at, p.updated_at,
			p.has_individual_location, p.individual_address, p.individual_latitude,
			p.individual_longitude, p.location_privacy, p.show_on_map, p.has_variants,
			c.id, c.name, c.slug, c.icon, c.parent_id
		FROM storefront_products p
		LEFT JOIN marketplace_categories c ON p.category_id = c.id
		WHERE p.id = $1`

	p := &models.StorefrontProduct{}
	c := &models.MarketplaceCategory{}
	var attributesJSON []byte
	var categoryID sql.NullInt64
	var categoryName, categorySlug, categoryIcon sql.NullString
	var categoryParentID sql.NullInt64

	err := s.pool.QueryRow(ctx, query, productID).Scan(
		&p.ID, &p.StorefrontID, &p.Name, &p.Description, &p.Price, &p.Currency,
		&p.CategoryID, &p.SKU, &p.Barcode, &p.StockQuantity, &p.StockStatus,
		&p.IsActive, &attributesJSON, &p.ViewCount, &p.SoldCount,
		&p.CreatedAt, &p.UpdatedAt,
		&p.HasIndividualLocation, &p.IndividualAddress, &p.IndividualLatitude,
		&p.IndividualLongitude, &p.LocationPrivacy, &p.ShowOnMap, &p.HasVariants,
		&categoryID, &categoryName, &categorySlug, &categoryIcon, &categoryParentID,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrStorefrontProductNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get storefront product by ID: %w", err)
	}

	if attributesJSON != nil {
		if unmarshalErr := json.Unmarshal(attributesJSON, &p.Attributes); unmarshalErr != nil {
			return nil, fmt.Errorf("failed to unmarshal attributes: %w", unmarshalErr)
		}
	}

	// Only set category if it exists
	if categoryID.Valid {
		c.ID = int(categoryID.Int64)
		c.Name = categoryName.String
		c.Slug = categorySlug.String
		if categoryIcon.Valid {
			c.Icon = &categoryIcon.String
		}
		if categoryParentID.Valid {
			parentID := int(categoryParentID.Int64)
			c.ParentID = &parentID
		}
		p.Category = c
	}

	// Load images
	images, err := s.getProductImages(ctx, []int{p.ID})
	if err != nil {
		return nil, fmt.Errorf("failed to get product images: %w", err)
	}
	// Initialize empty slice if no images found
	if images == nil {
		p.Images = []models.StorefrontProductImage{}
	} else {
		p.Images = images
	}

	// Load variants
	variants, err := s.getProductVariants(ctx, p.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product variants: %w", err)
	}
	p.Variants = variants

	// Load translations
	translations, err := s.GetTranslationsForEntity(ctx, "storefront_product", p.ID)
	if err != nil {
		logger.Error().Err(err).Int("product_id", p.ID).Msg("Failed to load product translations")
	} else if len(translations) > 0 {
		// Convert []Translation to map[language]map[field]text
		p.Translations = make(map[string]map[string]string)
		for _, t := range translations {
			if p.Translations[t.Language] == nil {
				p.Translations[t.Language] = make(map[string]string)
			}
			p.Translations[t.Language][t.FieldName] = t.TranslatedText
		}
	}

	return p, nil
}

// GetProductVariantByID retrieves a product variant by ID
func (s *Database) GetProductVariantByID(ctx context.Context, variantID int64) (*models.StorefrontProductVariant, error) {
	query := `
		SELECT 
			id, storefront_product_id, sku, price, stock_quantity, 
			variant_attributes, is_active, created_at, updated_at
		FROM storefront_product_variants
		WHERE id = $1`

	variant := &models.StorefrontProductVariant{}
	var attributesJSON []byte

	err := s.pool.QueryRow(ctx, query, variantID).Scan(
		&variant.ID, &variant.StorefrontProductID, &variant.SKU,
		&variant.Price, &variant.StockQuantity, &attributesJSON,
		&variant.IsActive, &variant.CreatedAt, &variant.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("variant not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get variant: %w", err)
	}

	// Parse variant attributes to create name
	if attributesJSON != nil {
		if unmarshalErr := json.Unmarshal(attributesJSON, &variant.Attributes); unmarshalErr != nil {
			logger.Error().Err(unmarshalErr).Msg("Failed to unmarshal variant attributes")
		}
		// Create variant name from attributes (e.g., "Red - Large")
		variant.Name = s.createVariantNameFromAttributes(variant.Attributes)
	}

	return variant, nil
}

// createVariantNameFromAttributes creates a readable name from variant attributes
func (s *Database) createVariantNameFromAttributes(attributes map[string]interface{}) string {
	if len(attributes) == 0 {
		return ""
	}

	var parts []string
	// Common attribute order for display
	order := []string{"color", "size", "material", "style"}

	for _, key := range order {
		if val, ok := attributes[key]; ok {
			parts = append(parts, fmt.Sprintf("%v", val))
		}
	}

	// Add any remaining attributes not in the order list
	for key, val := range attributes {
		found := false
		for _, orderedKey := range order {
			if key == orderedKey {
				found = true
				break
			}
		}
		if !found {
			parts = append(parts, fmt.Sprintf("%v", val))
		}
	}

	return strings.Join(parts, " - ")
}

// UpdateProductStock updates the stock quantity for a product or variant
func (s *Database) UpdateProductStock(ctx context.Context, productID int64, variantID *int64, quantity int) error {
	if variantID != nil && *variantID > 0 {
		// Update variant stock
		query := `
			UPDATE storefront_product_variants 
			SET stock_quantity = stock_quantity - $1, 
			    updated_at = CURRENT_TIMESTAMP
			WHERE id = $2 AND stock_quantity >= $1`

		result, err := s.pool.Exec(ctx, query, quantity, *variantID)
		if err != nil {
			return fmt.Errorf("failed to update variant stock: %w", err)
		}

		if result.RowsAffected() == 0 {
			return fmt.Errorf("insufficient stock for variant %d", *variantID)
		}
	} else {
		// Update product stock
		query := `
			UPDATE storefront_products 
			SET stock_quantity = stock_quantity - $1,
			    updated_at = CURRENT_TIMESTAMP
			WHERE id = $2 AND stock_quantity >= $1`

		result, err := s.pool.Exec(ctx, query, quantity, productID)
		if err != nil {
			return fmt.Errorf("failed to update product stock: %w", err)
		}

		if result.RowsAffected() == 0 {
			return fmt.Errorf("insufficient stock for product %d", productID)
		}
	}

	return nil
}

// CreateStorefrontProduct creates a new product
func (s *Database) CreateStorefrontProduct(ctx context.Context, storefrontID int, req *models.CreateProductRequest) (*models.StorefrontProduct, error) {
	// Log incoming request
	logger.Info().
		Int("storefront_id", storefrontID).
		Str("name", req.Name).
		Float64("price", req.Price).
		Int("category_id", req.CategoryID).
		Interface("attributes", req.Attributes).
		Msg("Creating storefront product")

	var attributesJSON []byte
	if req.Attributes != nil {
		var err error
		attributesJSON, err = json.Marshal(req.Attributes)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to marshal attributes")
			return nil, fmt.Errorf("failed to marshal attributes: %w", err)
		}
	}

	query := `
		INSERT INTO storefront_products (
			storefront_id, name, description, price, currency, category_id,
			sku, barcode, stock_quantity, is_active, attributes,
			has_individual_location, individual_address, individual_latitude,
			individual_longitude, location_privacy, show_on_map, has_variants
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
		RETURNING id, stock_status, created_at, updated_at`

	// Log query parameters
	logger.Debug().
		Int("storefront_id", storefrontID).
		Str("name", req.Name).
		Str("description", req.Description).
		Float64("price", req.Price).
		Str("currency", req.Currency).
		Int("category_id", req.CategoryID).
		Interface("sku", req.SKU).
		Interface("barcode", req.Barcode).
		Int("stock_quantity", req.StockQuantity).
		Bool("is_active", req.IsActive).
		Str("attributes_json", string(attributesJSON)).
		Msg("Executing INSERT query")

	// Prepare location values
	hasIndividualLocation := req.HasIndividualLocation != nil && *req.HasIndividualLocation
	showOnMap := req.ShowOnMap == nil || *req.ShowOnMap // Default true

	// Handle empty strings as NULL for SKU and Barcode
	var sku, barcode interface{}
	if req.SKU != nil && *req.SKU != "" {
		sku = req.SKU
	}
	if req.Barcode != nil && *req.Barcode != "" {
		barcode = req.Barcode
	}

	var product models.StorefrontProduct
	err := s.pool.QueryRow(
		ctx, query,
		storefrontID, req.Name, req.Description, req.Price, req.Currency, req.CategoryID,
		sku, barcode, req.StockQuantity, req.IsActive, attributesJSON,
		hasIndividualLocation, req.IndividualAddress, req.IndividualLatitude,
		req.IndividualLongitude, req.LocationPrivacy, showOnMap, req.HasVariants,
	).Scan(&product.ID, &product.StockStatus, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		logger.Error().
			Err(err).
			Int("storefront_id", storefrontID).
			Str("name", req.Name).
			Msg("Failed to create storefront product")
		return nil, fmt.Errorf("failed to create storefront product: %w", err)
	}

	// Populate the product with request data
	product.StorefrontID = storefrontID
	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.Currency = req.Currency
	product.CategoryID = req.CategoryID
	product.SKU = req.SKU
	product.Barcode = req.Barcode
	product.StockQuantity = req.StockQuantity
	product.IsActive = req.IsActive
	product.Attributes = req.Attributes
	product.HasIndividualLocation = hasIndividualLocation
	product.IndividualAddress = req.IndividualAddress
	product.IndividualLatitude = req.IndividualLatitude
	product.IndividualLongitude = req.IndividualLongitude
	product.LocationPrivacy = req.LocationPrivacy
	product.ShowOnMap = showOnMap
	product.HasVariants = req.HasVariants

	return &product, nil
}

// CreateStorefrontProductTx creates a new product within a transaction
func (s *Database) CreateStorefrontProductTx(ctx context.Context, tx *sqlx.Tx, storefrontID int, req *models.CreateProductRequest) (*models.StorefrontProduct, error) {
	// Log incoming request
	logger.Info().
		Int("storefront_id", storefrontID).
		Str("name", req.Name).
		Float64("price", req.Price).
		Int("category_id", req.CategoryID).
		Interface("attributes", req.Attributes).
		Msg("Creating storefront product with transaction")

	var attributesJSON []byte
	if req.Attributes != nil {
		var err error
		attributesJSON, err = json.Marshal(req.Attributes)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to marshal attributes")
			return nil, fmt.Errorf("failed to marshal attributes: %w", err)
		}
	}

	query := `
		INSERT INTO storefront_products (
			storefront_id, name, description, price, currency, category_id,
			sku, barcode, stock_quantity, is_active, attributes,
			has_individual_location, individual_address, individual_latitude,
			individual_longitude, location_privacy, show_on_map, has_variants
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
		RETURNING id, stock_status, created_at, updated_at`

	// Log query parameters
	logger.Debug().
		Int("storefront_id", storefrontID).
		Str("name", req.Name).
		Str("description", req.Description).
		Float64("price", req.Price).
		Str("currency", req.Currency).
		Int("category_id", req.CategoryID).
		Interface("sku", req.SKU).
		Interface("barcode", req.Barcode).
		Int("stock_quantity", req.StockQuantity).
		Bool("is_active", req.IsActive).
		Str("attributes_json", string(attributesJSON)).
		Msg("Executing INSERT query with transaction")

	// Prepare location values
	hasIndividualLocation := req.HasIndividualLocation != nil && *req.HasIndividualLocation
	showOnMap := req.ShowOnMap == nil || *req.ShowOnMap // Default true

	// Handle empty strings as NULL for SKU and Barcode
	var sku, barcode interface{}
	if req.SKU != nil && *req.SKU != "" {
		sku = req.SKU
	}
	if req.Barcode != nil && *req.Barcode != "" {
		barcode = req.Barcode
	}

	var product models.StorefrontProduct
	err := tx.QueryRowContext(
		ctx, query,
		storefrontID, req.Name, req.Description, req.Price, req.Currency, req.CategoryID,
		sku, barcode, req.StockQuantity, req.IsActive, attributesJSON,
		hasIndividualLocation, req.IndividualAddress, req.IndividualLatitude,
		req.IndividualLongitude, req.LocationPrivacy, showOnMap, req.HasVariants,
	).Scan(&product.ID, &product.StockStatus, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		logger.Error().
			Err(err).
			Int("storefront_id", storefrontID).
			Str("name", req.Name).
			Msg("Failed to create storefront product with transaction")
		return nil, fmt.Errorf("failed to create storefront product: %w", err)
	}

	// Populate the product with request data
	product.StorefrontID = storefrontID
	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.Currency = req.Currency
	product.CategoryID = req.CategoryID
	product.SKU = req.SKU
	product.Barcode = req.Barcode
	product.StockQuantity = req.StockQuantity
	product.IsActive = req.IsActive
	product.Attributes = req.Attributes
	product.HasIndividualLocation = hasIndividualLocation
	product.IndividualAddress = req.IndividualAddress
	product.IndividualLatitude = req.IndividualLatitude
	product.IndividualLongitude = req.IndividualLongitude
	product.LocationPrivacy = req.LocationPrivacy
	product.ShowOnMap = showOnMap
	product.HasVariants = req.HasVariants

	return &product, nil
}

// UpdateStorefrontProduct updates an existing product
func (s *Database) UpdateStorefrontProduct(ctx context.Context, storefrontID, productID int, req *models.UpdateProductRequest) error {
	// Log incoming request
	logger.Info().
		Int("storefront_id", storefrontID).
		Int("product_id", productID).
		Interface("request", req).
		Msg("Updating storefront product")

	var setClauses []string
	var args []interface{}
	argIndex := 1

	if req.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, *req.Name)
		argIndex++
	}

	if req.Description != nil {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, *req.Description)
		argIndex++
	}

	if req.Price != nil {
		setClauses = append(setClauses, fmt.Sprintf("price = $%d", argIndex))
		args = append(args, *req.Price)
		argIndex++
	}

	if req.CategoryID != nil {
		setClauses = append(setClauses, fmt.Sprintf("category_id = $%d", argIndex))
		args = append(args, *req.CategoryID)
		argIndex++
	}

	if req.SKU != nil {
		setClauses = append(setClauses, fmt.Sprintf("sku = $%d", argIndex))
		// Treat empty string as NULL
		if *req.SKU == "" {
			args = append(args, nil)
		} else {
			args = append(args, *req.SKU)
		}
		argIndex++
	}

	if req.Barcode != nil {
		setClauses = append(setClauses, fmt.Sprintf("barcode = $%d", argIndex))
		// Treat empty string as NULL
		if *req.Barcode == "" {
			args = append(args, nil)
		} else {
			args = append(args, *req.Barcode)
		}
		argIndex++
	}

	if req.StockQuantity != nil {
		setClauses = append(setClauses, fmt.Sprintf("stock_quantity = $%d", argIndex))
		args = append(args, *req.StockQuantity)
		argIndex++
	}

	if req.IsActive != nil {
		setClauses = append(setClauses, fmt.Sprintf("is_active = $%d", argIndex))
		args = append(args, *req.IsActive)
		argIndex++
	}

	if req.Attributes != nil {
		attributesJSON, err := json.Marshal(req.Attributes)
		if err != nil {
			return fmt.Errorf("failed to marshal attributes: %w", err)
		}
		setClauses = append(setClauses, fmt.Sprintf("attributes = $%d", argIndex))
		args = append(args, attributesJSON)
		argIndex++
	}

	if req.HasIndividualLocation != nil {
		setClauses = append(setClauses, fmt.Sprintf("has_individual_location = $%d", argIndex))
		args = append(args, *req.HasIndividualLocation)
		argIndex++
	}

	if req.IndividualAddress != nil {
		setClauses = append(setClauses, fmt.Sprintf("individual_address = $%d", argIndex))
		args = append(args, req.IndividualAddress)
		argIndex++
	}

	if req.IndividualLatitude != nil {
		setClauses = append(setClauses, fmt.Sprintf("individual_latitude = $%d", argIndex))
		args = append(args, req.IndividualLatitude)
		argIndex++
	}

	if req.IndividualLongitude != nil {
		setClauses = append(setClauses, fmt.Sprintf("individual_longitude = $%d", argIndex))
		args = append(args, req.IndividualLongitude)
		argIndex++
	}

	if req.LocationPrivacy != nil {
		setClauses = append(setClauses, fmt.Sprintf("location_privacy = $%d::location_privacy_level", argIndex))
		args = append(args, *req.LocationPrivacy)
		argIndex++
	}

	if req.ShowOnMap != nil {
		setClauses = append(setClauses, fmt.Sprintf("show_on_map = $%d", argIndex))
		args = append(args, *req.ShowOnMap)
		argIndex++
	}

	if len(setClauses) == 0 {
		return nil // Nothing to update
	}

	// Add WHERE conditions
	args = append(args, productID, storefrontID)

	query := fmt.Sprintf(`
		UPDATE storefront_products
		SET %s, updated_at = CURRENT_TIMESTAMP
		WHERE id = $%d AND storefront_id = $%d`,
		strings.Join(setClauses, ", "), argIndex, argIndex+1)

	// Log the query for debugging
	logger.Debug().
		Str("query", query).
		Interface("args", args).
		Msg("Executing UPDATE query")

	result, err := s.pool.Exec(ctx, query, args...)
	if err != nil {
		logger.Error().
			Err(err).
			Str("query", query).
			Interface("args", args).
			Msg("Failed to execute UPDATE query")
		return fmt.Errorf("failed to update storefront product: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

// DeleteStorefrontProduct soft deletes a product by setting is_active to false
func (s *Database) DeleteStorefrontProduct(ctx context.Context, storefrontID, productID int) error {
	// Начинаем транзакцию для атомарности операции
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			// Игнорируем ошибку если транзакция уже была завершена
			_ = rollbackErr
		}
	}()

	// Деактивируем товар в storefront_products (если еще активен)
	query := `UPDATE storefront_products
		SET is_active = false, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND storefront_id = $2`
	result, err := tx.Exec(ctx, query, productID, storefrontID)
	if err != nil {
		return fmt.Errorf("failed to deactivate storefront product: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("product not found or already deleted")
	}

	// Обновляем статус в marketplace_listings на 'disabled'
	updateQuery := `UPDATE marketplace_listings 
		SET status = 'disabled', updated_at = CURRENT_TIMESTAMP 
		WHERE id = $1 AND storefront_id = $2`

	_, err = tx.Exec(ctx, updateQuery, productID, storefrontID)
	if err != nil {
		return fmt.Errorf("failed to update marketplace listing status: %w", err)
	}

	// Коммитим транзакцию
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// HardDeleteStorefrontProduct permanently deletes a product and all related data
func (s *Database) HardDeleteStorefrontProduct(ctx context.Context, storefrontID, productID int) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			_ = rollbackErr // Explicitly ignore error if transaction was already committed
		}
	}()

	// Удаляем связанные данные в правильном порядке (от зависимых к независимым)

	// 1. Удаляем изображения товара
	_, err = tx.Exec(ctx, `DELETE FROM storefront_product_images WHERE storefront_product_id = $1`, productID)
	if err != nil {
		return fmt.Errorf("failed to delete product images: %w", err)
	}

	// 2. Удаляем варианты товара (если есть)
	_, err = tx.Exec(ctx, `DELETE FROM storefront_product_variants WHERE product_id = $1`, productID)
	if err != nil {
		return fmt.Errorf("failed to delete product variants: %w", err)
	}

	// 3. Удаляем записи в избранном
	_, err = tx.Exec(ctx, `DELETE FROM marketplace_favorites WHERE listing_id = $1`, productID)
	if err != nil {
		return fmt.Errorf("failed to delete favorites: %w", err)
	}

	// 4. Удаляем отзывы и оценки
	_, err = tx.Exec(ctx, `DELETE FROM reviews WHERE entity_type = $1 AND entity_id = $2`, "storefront_product", productID)
	if err != nil {
		return fmt.Errorf("failed to delete reviews: %w", err)
	}

	// 5. Удаляем записи инвентаризации
	_, err = tx.Exec(ctx, `DELETE FROM storefront_inventory_movements WHERE storefront_product_id = $1`, productID)
	if err != nil {
		return fmt.Errorf("failed to delete inventory movements: %w", err)
	}

	// 6. Удаляем из marketplace_listings
	_, err = tx.Exec(ctx, `DELETE FROM marketplace_listings WHERE id = $1 AND storefront_id = $2`, productID, storefrontID)
	if err != nil {
		return fmt.Errorf("failed to delete marketplace listing: %w", err)
	}

	// 7. Наконец, удаляем сам товар из storefront_products
	result, err := tx.Exec(ctx, `DELETE FROM storefront_products WHERE id = $1 AND storefront_id = $2`, productID, storefrontID)
	if err != nil {
		return fmt.Errorf("failed to delete storefront product: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("product not found")
	}

	// Коммитим транзакцию
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// UpdateProductInventory updates product stock with tracking
func (s *Database) UpdateProductInventory(ctx context.Context, storefrontID, productID int, userID int, req *models.UpdateInventoryRequest) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			// Игнорируем ошибку если транзакция уже была завершена
			_ = rollbackErr // Explicitly ignore error
		}
	}()

	// Update stock quantity
	var newQuantity int
	if req.Type == "adjustment" {
		newQuantity = req.Quantity
	} else {
		// Get current quantity
		var currentQuantity int
		err = tx.QueryRow(ctx,
			`SELECT stock_quantity FROM storefront_products WHERE id = $1 AND storefront_id = $2`,
			productID, storefrontID,
		).Scan(&currentQuantity)
		if err != nil {
			return fmt.Errorf("failed to get current stock: %w", err)
		}

		if req.Type == "in" {
			newQuantity = currentQuantity + req.Quantity
		} else { // out
			newQuantity = currentQuantity - req.Quantity
			if newQuantity < 0 {
				return fmt.Errorf("insufficient stock")
			}
		}
	}

	// Update stock
	_, err = tx.Exec(ctx,
		`UPDATE storefront_products SET stock_quantity = $1 WHERE id = $2 AND storefront_id = $3`,
		newQuantity, productID, storefrontID,
	)
	if err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}

	// Record movement
	_, err = tx.Exec(ctx,
		`INSERT INTO storefront_inventory_movements (
			storefront_product_id, type, quantity, reason, notes, user_id
		) VALUES ($1, $2, $3, $4, $5, $6)`,
		productID, req.Type, req.Quantity, req.Reason, req.Notes, userID,
	)
	if err != nil {
		return fmt.Errorf("failed to record inventory movement: %w", err)
	}

	return tx.Commit(ctx)
}

// GetProductStats returns statistics for storefront products
func (s *Database) GetProductStats(ctx context.Context, storefrontID int) (*models.ProductStats, error) {
	query := `
		SELECT
			COUNT(*) as total_products,
			COUNT(*) FILTER (WHERE is_active = true) as active_products,
			COUNT(*) FILTER (WHERE stock_status = 'out_of_stock') as out_of_stock,
			COUNT(*) FILTER (WHERE stock_status = 'low_stock') as low_stock,
			SUM(price * stock_quantity) as total_value,
			SUM(sold_count) as total_sold
		FROM storefront_products
		WHERE storefront_id = $1`

	var stats models.ProductStats
	err := s.pool.QueryRow(ctx, query, storefrontID).Scan(
		&stats.TotalProducts,
		&stats.ActiveProducts,
		&stats.OutOfStock,
		&stats.LowStock,
		&stats.TotalValue,
		&stats.TotalSold,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get product stats: %w", err)
	}

	return &stats, nil
}

// Helper functions

func (s *Database) getProductImages(ctx context.Context, productIDs []int) ([]models.StorefrontProductImage, error) {
	if len(productIDs) == 0 {
		return []models.StorefrontProductImage{}, nil
	}

	// Get images from product images table
	query := `
		SELECT
			spi.id,
			spi.storefront_product_id,
			spi.image_url,
			spi.thumbnail_url,
			spi.display_order,
			spi.is_default,
			'' as file_path,
			'' as file_name,
			0 as file_size,
			'' as content_type,
			'minio' as storage_type,
			'storefronts' as storage_bucket,
			spi.image_url as public_url,
			spi.created_at
		FROM storefront_product_images spi
		WHERE spi.storefront_product_id = ANY($1)
		ORDER BY is_default DESC, display_order ASC`

	rows, err := s.pool.Query(ctx, query, pq.Array(productIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []models.StorefrontProductImage
	for rows.Next() {
		var img models.StorefrontProductImage
		err := rows.Scan(
			&img.ID, &img.StorefrontProductID, &img.ImageURL, &img.ThumbnailURL,
			&img.DisplayOrder, &img.IsDefault, &img.FilePath, &img.FileName,
			&img.FileSize, &img.ContentType, &img.StorageType, &img.StorageBucket,
			&img.PublicURL, &img.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		images = append(images, img)
	}

	return images, nil
}

func (s *Database) getProductVariants(ctx context.Context, productID int) ([]models.StorefrontProductVariant, error) {
	query := `
		SELECT
			id,
			product_id as storefront_product_id,
			COALESCE(sku, '') as name,
			sku,
			COALESCE(price, 0) as price,
			stock_quantity,
			variant_attributes as attributes,
			is_active,
			created_at,
			updated_at
		FROM storefront_product_variants
		WHERE product_id = $1 AND is_active = true
		ORDER BY is_default DESC, id ASC`

	rows, err := s.pool.Query(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var variants []models.StorefrontProductVariant
	for rows.Next() {
		var v models.StorefrontProductVariant
		var attributesJSON []byte

		err := rows.Scan(
			&v.ID, &v.StorefrontProductID, &v.Name, &v.SKU, &v.Price,
			&v.StockQuantity, &attributesJSON, &v.IsActive, &v.CreatedAt, &v.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if attributesJSON != nil {
			if err := json.Unmarshal(attributesJSON, &v.Attributes); err != nil {
				return nil, err
			}
		}

		variants = append(variants, v)
	}

	return variants, nil
}

// Bulk operation methods

// BulkCreateProducts creates multiple products in a single transaction
func (s *Database) BulkCreateProducts(ctx context.Context, storefrontID int, products []models.CreateProductRequest) ([]int, []error) {
	var createdIDs []int
	var errors []error

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to begin transaction: %w", err)}
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			// Игнорируем ошибку если транзакция уже была завершена
			_ = err // Explicitly ignore error
		}
	}()

	for i, req := range products {
		// Create product
		var productID int
		attributesJSON, _ := json.Marshal(req.Attributes)

		err := tx.QueryRow(ctx,
			`INSERT INTO storefront_products (
				storefront_id, name, description, price, currency, category_id,
				sku, barcode, stock_quantity, is_active, attributes
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`,
			storefrontID, req.Name, req.Description, req.Price, req.Currency,
			req.CategoryID, req.SKU, req.Barcode, req.StockQuantity, req.IsActive, attributesJSON,
		).Scan(&productID)
		if err != nil {
			errors = append(errors, fmt.Errorf("product %d: %w", i, err))
			continue
		}

		createdIDs = append(createdIDs, productID)
	}

	if len(createdIDs) > 0 {
		if err := tx.Commit(ctx); err != nil {
			return nil, []error{fmt.Errorf("failed to commit transaction: %w", err)}
		}
	}

	return createdIDs, errors
}

// BulkUpdateProducts updates multiple products in a single transaction
func (s *Database) BulkUpdateProducts(ctx context.Context, storefrontID int, updates []models.BulkUpdateItem) ([]int, []error) {
	var updatedIDs []int
	var errors []error

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to begin transaction: %w", err)}
	}
	defer func() {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			// Игнорируем ошибку если транзакция уже была завершена
			_ = rollbackErr // Explicitly ignore error
		}
	}()

	// Verify all products belong to the storefront
	productIDs := make([]int, len(updates))
	for i, update := range updates {
		productIDs[i] = update.ProductID
	}

	var validProductIDs []int
	rows, err := tx.Query(ctx,
		`SELECT id FROM storefront_products WHERE id = ANY($1) AND storefront_id = $2`,
		pq.Array(productIDs), storefrontID,
	)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to verify products: %w", err)}
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err == nil {
			validProductIDs = append(validProductIDs, id)
		}
	}

	// Create a map for quick lookup
	validProductMap := make(map[int]bool)
	for _, id := range validProductIDs {
		validProductMap[id] = true
	}

	// Update each product
	for _, update := range updates {
		if !validProductMap[update.ProductID] {
			errors = append(errors, fmt.Errorf("product %d: not found or doesn't belong to storefront", update.ProductID))
			continue
		}

		// Build dynamic update query
		setClauses := []string{}
		args := []interface{}{}
		argIndex := 1

		if update.Updates.Name != nil {
			setClauses = append(setClauses, fmt.Sprintf("name = $%d", argIndex))
			args = append(args, *update.Updates.Name)
			argIndex++
		}

		if update.Updates.Description != nil {
			setClauses = append(setClauses, fmt.Sprintf("description = $%d", argIndex))
			args = append(args, *update.Updates.Description)
			argIndex++
		}

		if update.Updates.Price != nil {
			setClauses = append(setClauses, fmt.Sprintf("price = $%d", argIndex))
			args = append(args, *update.Updates.Price)
			argIndex++
		}

		if update.Updates.CategoryID != nil {
			setClauses = append(setClauses, fmt.Sprintf("category_id = $%d", argIndex))
			args = append(args, *update.Updates.CategoryID)
			argIndex++
		}

		if update.Updates.StockQuantity != nil {
			setClauses = append(setClauses, fmt.Sprintf("stock_quantity = $%d", argIndex))
			args = append(args, *update.Updates.StockQuantity)
			argIndex++
		}

		if update.Updates.IsActive != nil {
			setClauses = append(setClauses, fmt.Sprintf("is_active = $%d", argIndex))
			args = append(args, *update.Updates.IsActive)
			argIndex++
		}

		if len(setClauses) == 0 {
			errors = append(errors, fmt.Errorf("product %d: no updates provided", update.ProductID))
			continue
		}

		// Add updated_at
		setClauses = append(setClauses, "updated_at = CURRENT_TIMESTAMP")
		// Don't increment argIndex as we're not adding a parameter

		// Add WHERE clause
		args = append(args, update.ProductID)
		query := fmt.Sprintf(
			"UPDATE storefront_products SET %s WHERE id = $%d",
			strings.Join(setClauses, ", "),
			argIndex,
		)

		_, err := tx.Exec(ctx, query, args...)
		if err != nil {
			errors = append(errors, fmt.Errorf("product %d: %w", update.ProductID, err))
			continue
		}

		updatedIDs = append(updatedIDs, update.ProductID)
	}

	if len(updatedIDs) > 0 {
		if err := tx.Commit(ctx); err != nil {
			return nil, []error{fmt.Errorf("failed to commit transaction: %w", err)}
		}
	}

	return updatedIDs, errors
}

// BulkDeleteProducts deletes multiple products in a single transaction
func (s *Database) BulkDeleteProducts(ctx context.Context, storefrontID int, productIDs []int) ([]int, []error) {
	var deletedIDs []int
	var errors []error

	if len(productIDs) == 0 {
		return deletedIDs, errors
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to begin transaction: %w", err)}
	}
	defer func() {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			// Игнорируем ошибку если транзакция уже была завершена
			_ = rollbackErr // Explicitly ignore error
		}
	}()

	// Soft delete products that belong to the storefront
	rows, err := tx.Query(ctx,
		`UPDATE storefront_products
		SET is_active = false, updated_at = CURRENT_TIMESTAMP
		WHERE id = ANY($1) AND storefront_id = $2 AND is_active = true
		RETURNING id`,
		pq.Array(productIDs), storefrontID,
	)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to deactivate products: %w", err)}
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err == nil {
			deletedIDs = append(deletedIDs, id)
		}
	}

	// Обновляем статус в marketplace_listings для удалённых товаров
	if len(deletedIDs) > 0 {
		_, err = tx.Exec(ctx,
			`UPDATE marketplace_listings 
			SET status = 'disabled', updated_at = CURRENT_TIMESTAMP 
			WHERE id = ANY($1) AND storefront_id = $2`,
			pq.Array(deletedIDs), storefrontID,
		)
		if err != nil {
			return nil, []error{fmt.Errorf("failed to update marketplace listings status: %w", err)}
		}
	}

	// Check which products were not deleted
	deletedMap := make(map[int]bool)
	for _, id := range deletedIDs {
		deletedMap[id] = true
	}

	for _, id := range productIDs {
		if !deletedMap[id] {
			errors = append(errors, fmt.Errorf("product %d: not found or doesn't belong to storefront", id))
		}
	}

	if len(deletedIDs) > 0 {
		if err := tx.Commit(ctx); err != nil {
			return nil, []error{fmt.Errorf("failed to commit transaction: %w", err)}
		}
	}

	return deletedIDs, errors
}

// IncrementProductViews increments the view count for a product
func (s *Database) IncrementProductViews(ctx context.Context, productID int) error {
	query := `
		UPDATE storefront_products
		SET view_count = view_count + 1
		WHERE id = $1`

	result, err := s.pool.Exec(ctx, query, productID)
	if err != nil {
		return fmt.Errorf("failed to increment product views: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

// BulkUpdateStatus updates the status of multiple products
func (s *Database) BulkUpdateStatus(ctx context.Context, storefrontID int, productIDs []int, isActive bool) ([]int, []error) {
	var updatedIDs []int
	var errors []error

	if len(productIDs) == 0 {
		return updatedIDs, errors
	}

	rows, err := s.pool.Query(ctx,
		`UPDATE storefront_products
		SET is_active = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = ANY($2) AND storefront_id = $3
		RETURNING id`,
		isActive, pq.Array(productIDs), storefrontID,
	)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to update status: %w", err)}
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err == nil {
			updatedIDs = append(updatedIDs, id)
		}
	}

	// Check which products were not updated
	updatedMap := make(map[int]bool)
	for _, id := range updatedIDs {
		updatedMap[id] = true
	}

	for _, id := range productIDs {
		if !updatedMap[id] {
			errors = append(errors, fmt.Errorf("product %d: not found or doesn't belong to storefront", id))
		}
	}

	return updatedIDs, errors
}

// loadAddressTranslations загружает переводы адресов из таблицы translations для списка продуктов
func (s *Database) loadAddressTranslations(ctx context.Context, products []*models.StorefrontProduct) error {
	if len(products) == 0 {
		return nil
	}

	// Собираем ID продуктов
	productIDs := make([]int, len(products))
	productMap := make(map[int]*models.StorefrontProduct)
	for i, p := range products {
		productIDs[i] = p.ID
		productMap[p.ID] = p
		// Инициализируем Translations если нужно
		if p.Translations == nil {
			p.Translations = make(map[string]map[string]string)
		}
	}

	// Запрашиваем переводы из таблицы translations
	query := `
		SELECT entity_id, language, field_name, translated_text
		FROM translations
		WHERE entity_type = 'storefront_product'
		  AND entity_id = ANY($1)
		  AND language IN ('en', 'ru', 'sr')
		ORDER BY entity_id, language, field_name`

	rows, err := s.pool.Query(ctx, query, pq.Array(productIDs))
	if err != nil {
		logger.Warn().Err(err).Msg("Failed to load translations")
		return nil // Не фатальная ошибка, продолжаем без переводов
	}
	defer rows.Close()

	// Группируем переводы по продукту и языку
	for rows.Next() {
		var entityID int
		var language, fieldName, translatedText string
		if err := rows.Scan(&entityID, &language, &fieldName, &translatedText); err != nil {
			logger.Warn().Err(err).Msg("Failed to scan translation")
			continue
		}

		if p, ok := productMap[entityID]; ok {
			if p.Translations[language] == nil {
				p.Translations[language] = make(map[string]string)
			}
			p.Translations[language][fieldName] = translatedText
		}
	}

	logger.Debug().
		Int("products_count", len(products)).
		Int("product_ids", len(productIDs)).
		Msg("Loaded product translations")

	return nil
}

// loadAddressTranslationsForProduct загружает переводы для одного продукта из таблицы translations
func (s *Database) loadAddressTranslationsForProduct(ctx context.Context, product *models.StorefrontProduct) error {
	// Инициализируем Translations если нужно
	if product.Translations == nil {
		product.Translations = make(map[string]map[string]string)
	}

	query := `
		SELECT language, field_name, translated_text
		FROM translations
		WHERE entity_type = 'storefront_product'
		  AND entity_id = $1
		  AND language IN ('en', 'ru', 'sr')
		ORDER BY language, field_name`

	rows, err := s.pool.Query(ctx, query, product.ID)
	if err != nil {
		logger.Warn().Err(err).Msg("Failed to load translations for product")
		return nil // Не фатальная ошибка
	}
	defer rows.Close()

	for rows.Next() {
		var language, fieldName, translatedText string
		if err := rows.Scan(&language, &fieldName, &translatedText); err != nil {
			logger.Warn().Err(err).Msg("Failed to scan translation")
			continue
		}

		if product.Translations[language] == nil {
			product.Translations[language] = make(map[string]string)
		}
		product.Translations[language][fieldName] = translatedText
	}

	if len(product.Translations) > 0 {
		logger.Debug().
			Int("product_id", product.ID).
			Int("languages_count", len(product.Translations)).
			Interface("translations", product.Translations).
			Msg("Loaded translations for product")
	}

	return nil
}

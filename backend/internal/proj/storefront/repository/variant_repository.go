package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"backend/internal/proj/storefront/types"

	"github.com/jmoiron/sqlx"
)

type VariantRepository struct {
	db *sqlx.DB
}

func NewVariantRepository(db *sqlx.DB) *VariantRepository {
	return &VariantRepository{db: db}
}

// GetVariantAttributes returns all variant attributes
func (r *VariantRepository) GetVariantAttributes(ctx context.Context) ([]types.ProductVariantAttribute, error) {
	query := `
		SELECT id, name, display_name, type, is_required, sort_order, created_at, updated_at
		FROM product_variant_attributes
		ORDER BY sort_order, name`

	var attributes []types.ProductVariantAttribute
	err := r.db.SelectContext(ctx, &attributes, query)
	return attributes, err
}

// GetVariantAttributeValues returns values for a specific attribute
func (r *VariantRepository) GetVariantAttributeValues(ctx context.Context, attributeID int) ([]types.ProductVariantAttributeValue, error) {
	query := `
		SELECT id, attribute_id, value, display_name, color_hex, image_url, sort_order, is_active, created_at, updated_at
		FROM product_variant_attribute_values
		WHERE attribute_id = $1 AND is_active = true
		ORDER BY sort_order, display_name`

	var values []types.ProductVariantAttributeValue
	err := r.db.SelectContext(ctx, &values, query, attributeID)
	return values, err
}

// CreateVariant creates a new product variant
func (r *VariantRepository) CreateVariant(ctx context.Context, req *types.CreateVariantRequest) (*types.ProductVariant, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			// Transaction was already committed or rolled back, ignore
		}
	}()

	// If this is set as default, unset other defaults for this product
	if req.IsDefault {
		_, err = tx.ExecContext(ctx,
			"UPDATE storefront_product_variants SET is_default = false WHERE product_id = $1",
			req.ProductID)
		if err != nil {
			return nil, err
		}
	}

	// Marshal JSON fields
	variantAttrsJSON, err := json.Marshal(req.VariantAttributes)
	if err != nil {
		return nil, err
	}

	var dimensionsJSON []byte
	if req.Dimensions != nil {
		dimensionsJSON, err = json.Marshal(req.Dimensions)
		if err != nil {
			return nil, err
		}
	}

	// Determine stock status
	stockStatus := "in_stock"
	if req.StockQuantity == 0 {
		stockStatus = "out_of_stock"
	} else if req.LowStockThreshold != nil && req.StockQuantity <= *req.LowStockThreshold {
		stockStatus = "low_stock"
	}

	// Insert variant
	query := `
		INSERT INTO storefront_product_variants (
			product_id, sku, barcode, price, compare_at_price, cost_price,
			stock_quantity, stock_status, low_stock_threshold, variant_attributes,
			weight, dimensions, is_default, is_active
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, true
		) RETURNING id, created_at, updated_at`

	var variant types.ProductVariant
	err = tx.QueryRowContext(ctx, query,
		req.ProductID, req.SKU, req.Barcode, req.Price, req.CompareAtPrice, req.CostPrice,
		req.StockQuantity, stockStatus, req.LowStockThreshold, variantAttrsJSON,
		req.Weight, dimensionsJSON, req.IsDefault,
	).Scan(&variant.ID, &variant.CreatedAt, &variant.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Set other fields
	variant.ProductID = req.ProductID
	variant.SKU = req.SKU
	variant.Barcode = req.Barcode
	variant.Price = req.Price
	variant.CompareAtPrice = req.CompareAtPrice
	variant.CostPrice = req.CostPrice
	variant.StockQuantity = req.StockQuantity
	variant.StockStatus = stockStatus
	variant.LowStockThreshold = req.LowStockThreshold
	variant.VariantAttributes = req.VariantAttributes
	variant.Weight = req.Weight
	variant.Dimensions = req.Dimensions
	variant.IsDefault = req.IsDefault
	variant.IsActive = true

	// Create variant images
	if len(req.Images) > 0 {
		for i, img := range req.Images {
			imageQuery := `
				INSERT INTO storefront_product_variant_images (
					variant_id, image_url, thumbnail_url, alt_text, display_order, is_main
				) VALUES ($1, $2, $3, $4, $5, $6)`

			_, err = tx.ExecContext(ctx, imageQuery,
				variant.ID, img.ImageURL, img.ThumbnailURL, img.AltText,
				img.DisplayOrder, img.IsMain && i == 0) // Only first image can be main if multiple are marked
			if err != nil {
				return nil, err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// Load complete variant with images
	return r.GetVariantByID(ctx, variant.ID)
}

// GetVariantByID returns a variant by ID with images
func (r *VariantRepository) GetVariantByID(ctx context.Context, id int) (*types.ProductVariant, error) {
	query := `
		SELECT id, product_id, sku, barcode, price, compare_at_price, cost_price,
			   stock_quantity, stock_status, low_stock_threshold, variant_attributes,
			   weight, dimensions, is_active, is_default, view_count, sold_count,
			   created_at, updated_at
		FROM storefront_product_variants
		WHERE id = $1`

	var variant types.ProductVariant
	var variantAttrsJSON, dimensionsJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&variant.ID, &variant.ProductID, &variant.SKU, &variant.Barcode,
		&variant.Price, &variant.CompareAtPrice, &variant.CostPrice,
		&variant.StockQuantity, &variant.StockStatus, &variant.LowStockThreshold,
		&variantAttrsJSON, &variant.Weight, &dimensionsJSON,
		&variant.IsActive, &variant.IsDefault, &variant.ViewCount, &variant.SoldCount,
		&variant.CreatedAt, &variant.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON fields
	if len(variantAttrsJSON) > 0 {
		err = json.Unmarshal(variantAttrsJSON, &variant.VariantAttributes)
		if err != nil {
			return nil, err
		}
	}

	if len(dimensionsJSON) > 0 {
		err = json.Unmarshal(dimensionsJSON, &variant.Dimensions)
		if err != nil {
			return nil, err
		}
	}

	// Load images
	variant.Images, err = r.getVariantImages(ctx, variant.ID)
	if err != nil {
		return nil, err
	}

	return &variant, nil
}

// GetVariantsByProductID returns all variants for a product
func (r *VariantRepository) GetVariantsByProductID(ctx context.Context, productID int) ([]types.ProductVariant, error) {
	query := `
		SELECT id, product_id, sku, barcode, price, compare_at_price, cost_price,
			   stock_quantity, stock_status, low_stock_threshold, variant_attributes,
			   weight, dimensions, is_active, is_default, view_count, sold_count,
			   created_at, updated_at
		FROM storefront_product_variants
		WHERE product_id = $1 AND is_active = true
		ORDER BY is_default DESC, created_at ASC`

	rows, err := r.db.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var variants []types.ProductVariant
	for rows.Next() {
		var variant types.ProductVariant
		var variantAttrsJSON, dimensionsJSON []byte

		err = rows.Scan(
			&variant.ID, &variant.ProductID, &variant.SKU, &variant.Barcode,
			&variant.Price, &variant.CompareAtPrice, &variant.CostPrice,
			&variant.StockQuantity, &variant.StockStatus, &variant.LowStockThreshold,
			&variantAttrsJSON, &variant.Weight, &dimensionsJSON,
			&variant.IsActive, &variant.IsDefault, &variant.ViewCount, &variant.SoldCount,
			&variant.CreatedAt, &variant.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Unmarshal JSON fields
		if len(variantAttrsJSON) > 0 {
			err = json.Unmarshal(variantAttrsJSON, &variant.VariantAttributes)
			if err != nil {
				return nil, err
			}
		}

		if len(dimensionsJSON) > 0 {
			err = json.Unmarshal(dimensionsJSON, &variant.Dimensions)
			if err != nil {
				return nil, err
			}
		}

		// Load images for each variant
		variant.Images, err = r.getVariantImages(ctx, variant.ID)
		if err != nil {
			return nil, err
		}

		variants = append(variants, variant)
	}

	return variants, rows.Err()
}

// GenerateVariants automatically generates variants based on attribute matrix
func (r *VariantRepository) GenerateVariants(ctx context.Context, req *types.GenerateVariantsRequest) ([]types.ProductVariant, error) {
	// Generate all possible combinations
	combinations := r.generateAttributeCombinations(req.AttributeMatrix)

	var variants []types.ProductVariant
	for _, combination := range combinations {
		// Create variant request
		variantReq := &types.CreateVariantRequest{
			ProductID:         req.ProductID,
			VariantAttributes: combination,
			StockQuantity:     r.getStockQuantity(combination, req.StockQuantities),
			Price:             r.calculatePrice(combination, req.PriceModifiers),
			IsDefault:         r.isDefaultCombination(combination, req.DefaultAttributes),
			Images:            r.getImagesForCombination(combination, req.ImageMappings),
		}

		// Generate SKU
		variantReq.SKU = r.generateSKU(req.ProductID, combination)

		// Create variant
		variant, err := r.CreateVariant(ctx, variantReq)
		if err != nil {
			return nil, fmt.Errorf("failed to create variant %v: %w", combination, err)
		}

		variants = append(variants, *variant)
	}

	return variants, nil
}

// generateAttributeCombinations creates all possible combinations from attribute matrix
func (r *VariantRepository) generateAttributeCombinations(matrix map[string][]string) []map[string]interface{} {
	if len(matrix) == 0 {
		return []map[string]interface{}{}
	}

	// Convert to slice for easier processing
	var attributes []string
	var values [][]string
	for attr, vals := range matrix {
		attributes = append(attributes, attr)
		values = append(values, vals)
	}

	// Generate combinations recursively
	var combinations []map[string]interface{}
	r.generateCombinationsRecursive(attributes, values, 0, make(map[string]interface{}), &combinations)

	return combinations
}

// generateCombinationsRecursive helper for recursive combination generation
func (r *VariantRepository) generateCombinationsRecursive(attributes []string, values [][]string, index int, current map[string]interface{}, result *[]map[string]interface{}) {
	if index == len(attributes) {
		// Make a copy of current combination
		combination := make(map[string]interface{})
		for k, v := range current {
			combination[k] = v
		}
		*result = append(*result, combination)
		return
	}

	// Try each value for current attribute
	for _, value := range values[index] {
		current[attributes[index]] = value
		r.generateCombinationsRecursive(attributes, values, index+1, current, result)
	}
}

// Helper functions for variant generation
func (r *VariantRepository) getStockQuantity(combination map[string]interface{}, stockQuantities map[string]int) int {
	// Create key from combination
	key := r.createCombinationKey(combination)
	if qty, exists := stockQuantities[key]; exists {
		return qty
	}
	return 10 // Default stock quantity
}

func (r *VariantRepository) calculatePrice(combination map[string]interface{}, priceModifiers map[string]float64) *float64 {
	var totalModifier float64
	for _, value := range combination {
		if modifier, exists := priceModifiers[value.(string)]; exists {
			totalModifier += modifier
		}
	}
	if totalModifier != 0 {
		return &totalModifier
	}
	return nil // Use parent product price
}

func (r *VariantRepository) isDefaultCombination(combination map[string]interface{}, defaultAttributes map[string]string) bool {
	if len(defaultAttributes) == 0 {
		return false
	}

	for attr, value := range defaultAttributes {
		if combination[attr] != value {
			return false
		}
	}
	return true
}

func (r *VariantRepository) getImagesForCombination(combination map[string]interface{}, imageMappings map[string][]types.CreateVariantImageRequest) []types.CreateVariantImageRequest {
	// Try to find images for any attribute value in the combination
	for _, value := range combination {
		if images, exists := imageMappings[value.(string)]; exists {
			return images
		}
	}
	return []types.CreateVariantImageRequest{}
}

func (r *VariantRepository) generateSKU(productID int, combination map[string]interface{}) *string {
	var parts []string
	parts = append(parts, fmt.Sprintf("PROD-%d", productID))

	for attr, value := range combination {
		parts = append(parts, fmt.Sprintf("%s-%s", strings.ToUpper(attr), strings.ToUpper(value.(string))))
	}

	sku := strings.Join(parts, "-")
	return &sku
}

func (r *VariantRepository) createCombinationKey(combination map[string]interface{}) string {
	var parts []string
	for attr, value := range combination {
		parts = append(parts, fmt.Sprintf("%s-%s", attr, value.(string)))
	}
	return strings.Join(parts, "-")
}

// SetupProductAttributes configures attributes for a seller's product
func (r *VariantRepository) SetupProductAttributes(ctx context.Context, req *types.SetupProductAttributesRequest) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			// Transaction was already committed or rolled back, ignore
		}
	}()

	// Delete existing attribute configurations for this product
	_, err = tx.ExecContext(ctx, "DELETE FROM storefront_product_attributes WHERE product_id = $1", req.ProductID)
	if err != nil {
		return err
	}

	// Insert new attribute configurations
	for _, attr := range req.Attributes {
		if !attr.IsEnabled {
			continue // Skip disabled attributes
		}

		// Combine custom values with selected global values
		allValues := attr.CustomValues

		// Add selected global values
		if len(attr.SelectedGlobalValues) > 0 {
			globalValues, err := r.GetVariantAttributeValues(ctx, attr.AttributeID)
			if err != nil {
				return err
			}

			for _, globalVal := range globalValues {
				for _, selectedVal := range attr.SelectedGlobalValues {
					if globalVal.Value == selectedVal {
						allValues = append(allValues, types.AttributeValue{
							Value:       globalVal.Value,
							DisplayName: globalVal.DisplayName,
							ColorHex:    globalVal.ColorHex,
							ImageURL:    globalVal.ImageURL,
							IsCustom:    false,
						})
						break
					}
				}
			}
		}

		// Marshal custom values to JSON
		customValuesJSON, err := json.Marshal(allValues)
		if err != nil {
			return err
		}

		// Insert attribute configuration
		_, err = tx.ExecContext(ctx, `
			INSERT INTO storefront_product_attributes (product_id, attribute_id, is_enabled, is_required, custom_values)
			VALUES ($1, $2, $3, $4, $5)`,
			req.ProductID, attr.AttributeID, attr.IsEnabled, attr.IsRequired, customValuesJSON)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetProductAttributes returns configured attributes for a product
func (r *VariantRepository) GetProductAttributes(ctx context.Context, productID int) ([]types.StorefrontProductAttribute, error) {
	query := `
		SELECT spa.id, spa.product_id, spa.attribute_id, spa.is_enabled, spa.is_required,
			   spa.custom_values, spa.created_at, spa.updated_at,
			   pva.name, pva.display_name, pva.type, pva.sort_order
		FROM storefront_product_attributes spa
		JOIN product_variant_attributes pva ON spa.attribute_id = pva.id
		WHERE spa.product_id = $1 AND spa.is_enabled = true
		ORDER BY pva.sort_order, pva.name`

	rows, err := r.db.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attributes []types.StorefrontProductAttribute
	for rows.Next() {
		var attr types.StorefrontProductAttribute
		var customValuesJSON []byte
		var attribute types.ProductVariantAttribute

		err = rows.Scan(
			&attr.ID, &attr.ProductID, &attr.AttributeID, &attr.IsEnabled, &attr.IsRequired,
			&customValuesJSON, &attr.CreatedAt, &attr.UpdatedAt,
			&attribute.Name, &attribute.DisplayName, &attribute.Type, &attribute.SortOrder,
		)
		if err != nil {
			return nil, err
		}

		// Unmarshal custom values
		if len(customValuesJSON) > 0 {
			err = json.Unmarshal(customValuesJSON, &attr.CustomValues)
			if err != nil {
				return nil, err
			}
		}

		attr.Attribute = &attribute
		attributes = append(attributes, attr)
	}

	return attributes, rows.Err()
}

// GetAvailableAttributesForCategory returns all attributes available for a category
func (r *VariantRepository) GetAvailableAttributesForCategory(ctx context.Context, categoryID int) ([]types.ProductVariantAttribute, error) {
	// For now, return all attributes. In the future, we can filter by category
	return r.GetVariantAttributes(ctx)
}

// getVariantImages loads images for a variant
func (r *VariantRepository) getVariantImages(ctx context.Context, variantID int) ([]types.ProductVariantImage, error) {
	query := `
		SELECT id, variant_id, image_url, thumbnail_url, alt_text, display_order, is_main, created_at
		FROM storefront_product_variant_images
		WHERE variant_id = $1
		ORDER BY is_main DESC, display_order ASC, created_at ASC`

	var images []types.ProductVariantImage
	err := r.db.SelectContext(ctx, &images, query, variantID)
	return images, err
}

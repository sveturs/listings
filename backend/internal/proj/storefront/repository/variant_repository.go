package repository

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"backend/internal/domain/models"
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
		SELECT id, name, display_name, type, is_required, sort_order, affects_stock, created_at, updated_at
		FROM product_variant_attributes
		ORDER BY sort_order, name`

	var attributes []types.ProductVariantAttribute
	err := r.db.SelectContext(ctx, &attributes, query)
	return attributes, err
}

// GetVariantAttributeValues returns values for a specific attribute
func (r *VariantRepository) GetVariantAttributeValues(ctx context.Context, attributeID int) ([]types.ProductVariantAttributeValue, error) {
	query := `
		SELECT id, attribute_id, value, display_name, color_hex, image_url, sort_order, is_active, 
		       is_popular, usage_count, metadata, created_at, updated_at
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
		if rollbackErr := tx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			// Transaction was already committed or rolled back, ignore
			_ = rollbackErr // Explicitly ignore the error
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

// CreateVariantTx creates a new product variant using the provided transaction
func (r *VariantRepository) CreateVariantTx(ctx context.Context, tx *sqlx.Tx, req *types.CreateVariantRequest) (*types.ProductVariant, error) {
	// If this is set as default, unset other defaults for this product
	if req.IsDefault {
		_, err := tx.ExecContext(ctx,
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

	// Note: Don't commit the transaction here, let the caller do it
	// Return the variant with all populated data
	return &variant, nil
}

// BulkCreateVariantsTx creates multiple variants in a single transaction
func (r *VariantRepository) BulkCreateVariantsTx(ctx context.Context, tx *sqlx.Tx, productID int, variants []types.CreateVariantRequest) ([]*types.ProductVariant, error) {
	if len(variants) == 0 {
		return []*types.ProductVariant{}, nil
	}

	createdVariants := make([]*types.ProductVariant, 0, len(variants))

	// Ensure productID is set for all variants
	for i := range variants {
		variants[i].ProductID = productID
	}

	// Check if any variant is set as default
	hasDefault := false
	for _, v := range variants {
		if v.IsDefault {
			hasDefault = true
			break
		}
	}

	// If we have a default variant, unset all existing defaults for this product
	if hasDefault {
		_, err := tx.ExecContext(ctx,
			"UPDATE storefront_product_variants SET is_default = false WHERE product_id = $1",
			productID)
		if err != nil {
			return nil, fmt.Errorf("failed to unset existing default variants: %w", err)
		}
	}

	// Create each variant
	for i, variantReq := range variants {
		variant, err := r.CreateVariantTx(ctx, tx, &variantReq)
		if err != nil {
			return nil, fmt.Errorf("failed to create variant %d (SKU: %s): %w", i, variantReq.SKU, err)
		}
		createdVariants = append(createdVariants, variant)
	}

	return createdVariants, nil
}

// GetVariantByID returns a variant by ID with images
func (r *VariantRepository) GetVariantByID(ctx context.Context, id int) (*types.ProductVariant, error) {
	query := `
		SELECT id, product_id, sku, barcode, price, compare_at_price, cost_price,
			   stock_quantity, reserved_quantity, available_quantity, stock_status, low_stock_threshold, variant_attributes,
			   weight, dimensions, is_active, is_default, view_count, sold_count,
			   created_at, updated_at
		FROM storefront_product_variants
		WHERE id = $1`

	var variant types.ProductVariant
	var variantAttrsJSON, dimensionsJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&variant.ID, &variant.ProductID, &variant.SKU, &variant.Barcode,
		&variant.Price, &variant.CompareAtPrice, &variant.CostPrice,
		&variant.StockQuantity, &variant.ReservedQuantity, &variant.AvailableQuantity, &variant.StockStatus, &variant.LowStockThreshold,
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
			   stock_quantity, reserved_quantity, available_quantity, stock_status, low_stock_threshold, variant_attributes,
			   weight, dimensions, is_active, is_default, view_count, sold_count,
			   created_at, updated_at
		FROM storefront_product_variants
		WHERE product_id = $1 AND is_active = true
		ORDER BY is_default DESC, created_at ASC`

	rows, err := r.db.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			// Логируем ошибку закрытия rows
			_ = closeErr // Explicitly ignore error
		}
	}()

	var variants []types.ProductVariant
	for rows.Next() {
		var variant types.ProductVariant
		var variantAttrsJSON, dimensionsJSON []byte

		err = rows.Scan(
			&variant.ID, &variant.ProductID, &variant.SKU, &variant.Barcode,
			&variant.Price, &variant.CompareAtPrice, &variant.CostPrice,
			&variant.StockQuantity, &variant.ReservedQuantity, &variant.AvailableQuantity, &variant.StockStatus, &variant.LowStockThreshold,
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
		if rollbackErr := tx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			// Transaction was already committed or rolled back, ignore
			_ = rollbackErr // Explicitly ignore the error
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
			globalValues, attrErr := r.GetVariantAttributeValues(ctx, attr.AttributeID)
			if attrErr != nil {
				return attrErr
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
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			// Логируем ошибку закрытия rows
			_ = closeErr // Explicitly ignore error
		}
	}()

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

// DeleteVariant soft deletes a variant by setting is_active to false
func (r *VariantRepository) DeleteVariant(ctx context.Context, variantID int) error {
	query := `UPDATE storefront_product_variants SET is_active = false WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, variantID)
	return err
}

// UpdateVariant updates a variant's information
func (r *VariantRepository) UpdateVariant(ctx context.Context, variantID int, req *types.UpdateVariantRequest) (*types.ProductVariant, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if rollbackErr := tx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			_ = rollbackErr
		}
	}()

	// Build dynamic query
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.SKU != nil {
		setParts = append(setParts, fmt.Sprintf("sku = $%d", argIndex))
		args = append(args, *req.SKU)
		argIndex++
	}

	if req.Barcode != nil {
		setParts = append(setParts, fmt.Sprintf("barcode = $%d", argIndex))
		args = append(args, *req.Barcode)
		argIndex++
	}

	if req.Price != nil {
		setParts = append(setParts, fmt.Sprintf("price = $%d", argIndex))
		args = append(args, *req.Price)
		argIndex++
	}

	if req.CompareAtPrice != nil {
		setParts = append(setParts, fmt.Sprintf("compare_at_price = $%d", argIndex))
		args = append(args, *req.CompareAtPrice)
		argIndex++
	}

	if req.CostPrice != nil {
		setParts = append(setParts, fmt.Sprintf("cost_price = $%d", argIndex))
		args = append(args, *req.CostPrice)
		argIndex++
	}

	if req.StockQuantity != nil {
		setParts = append(setParts, fmt.Sprintf("stock_quantity = $%d", argIndex))
		args = append(args, *req.StockQuantity)
		argIndex++
	}

	if req.LowStockThreshold != nil {
		setParts = append(setParts, fmt.Sprintf("low_stock_threshold = $%d", argIndex))
		args = append(args, *req.LowStockThreshold)
		argIndex++
	}

	if req.VariantAttributes != nil {
		variantAttrsJSON, err := json.Marshal(req.VariantAttributes)
		if err != nil {
			return nil, err
		}
		setParts = append(setParts, fmt.Sprintf("variant_attributes = $%d", argIndex))
		args = append(args, variantAttrsJSON)
		argIndex++
	}

	if req.Weight != nil {
		setParts = append(setParts, fmt.Sprintf("weight = $%d", argIndex))
		args = append(args, *req.Weight)
		argIndex++
	}

	if req.Dimensions != nil {
		dimensionsJSON, err := json.Marshal(req.Dimensions)
		if err != nil {
			return nil, err
		}
		setParts = append(setParts, fmt.Sprintf("dimensions = $%d", argIndex))
		args = append(args, dimensionsJSON)
		argIndex++
	}

	if req.IsActive != nil {
		setParts = append(setParts, fmt.Sprintf("is_active = $%d", argIndex))
		args = append(args, *req.IsActive)
		argIndex++
	}

	if req.IsDefault != nil && *req.IsDefault {
		// Unset other defaults for this product first
		_, err = tx.ExecContext(ctx, `
			UPDATE storefront_product_variants 
			SET is_default = false 
			WHERE product_id = (SELECT product_id FROM storefront_product_variants WHERE id = $1)`,
			variantID)
		if err != nil {
			return nil, err
		}

		setParts = append(setParts, fmt.Sprintf("is_default = $%d", argIndex))
		args = append(args, *req.IsDefault)
		argIndex++
	}

	if len(setParts) == 0 {
		return r.GetVariantByID(ctx, variantID) // No changes
	}

	// Add updated_at
	setParts = append(setParts, "updated_at = NOW()")

	// Add WHERE clause
	args = append(args, variantID)
	whereClause := fmt.Sprintf(" WHERE id = $%d", argIndex)

	query := fmt.Sprintf("UPDATE storefront_product_variants SET %s%s", strings.Join(setParts, ", "), whereClause)

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return r.GetVariantByID(ctx, variantID)
}

// GetVariantMatrix returns all possible variant combinations for a product
func (r *VariantRepository) GetVariantMatrix(ctx context.Context, productID int) (*types.VariantMatrixResponse, error) {
	// Get configured attributes for this product
	productAttrs, err := r.GetProductAttributes(ctx, productID)
	if err != nil {
		return nil, err
	}

	// Get existing variants
	existingVariants, err := r.GetVariantsByProductID(ctx, productID)
	if err != nil {
		return nil, err
	}

	// Build attributes map with available values
	attributes := make(map[string][]types.AttributeValue)
	possibleCombinations := 1

	for _, attr := range productAttrs {
		if attr.Attribute != nil {
			attributes[attr.Attribute.Name] = attr.CustomValues
			possibleCombinations *= len(attr.CustomValues)
		}
	}

	return &types.VariantMatrixResponse{
		Attributes:           attributes,
		ExistingVariants:     existingVariants,
		PossibleCombinations: possibleCombinations,
	}, nil
}

// BulkUpdateStock updates stock quantities for multiple variants
func (r *VariantRepository) BulkUpdateStock(ctx context.Context, productID int, req *types.BulkUpdateStockRequest) (int, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() {
		if rollbackErr := tx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			_ = rollbackErr
		}
	}()

	updatedCount := 0
	for _, update := range req.Updates {
		query := `
			UPDATE storefront_product_variants 
			SET stock_quantity = $1, updated_at = NOW()
			WHERE id = $2 AND product_id = $3`

		result, err := tx.ExecContext(ctx, query, update.StockQuantity, update.VariantID, productID)
		if err != nil {
			return 0, err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return 0, err
		}

		updatedCount += int(rowsAffected)
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return updatedCount, nil
}

// GetVariantAnalytics returns analytics data for product variants
func (r *VariantRepository) GetVariantAnalytics(ctx context.Context, productID int) (*types.VariantAnalyticsResponse, error) {
	// Basic analytics query
	query := `
		SELECT 
			COUNT(*) as total_variants,
			COALESCE(SUM(stock_quantity), 0) as total_stock,
			COALESCE(SUM(sold_count), 0) as total_sold
		FROM storefront_product_variants 
		WHERE product_id = $1 AND is_active = true`

	var analytics types.VariantAnalyticsResponse
	err := r.db.QueryRowContext(ctx, query, productID).Scan(
		&analytics.TotalVariants,
		&analytics.TotalStock,
		&analytics.TotalSold,
	)
	if err != nil {
		return nil, err
	}

	// Find best seller
	bestSellerQuery := `
		SELECT id, product_id, sku, barcode, price, compare_at_price, cost_price,
			   stock_quantity, reserved_quantity, available_quantity, stock_status, low_stock_threshold, variant_attributes,
			   weight, dimensions, is_active, is_default, view_count, sold_count,
			   created_at, updated_at
		FROM storefront_product_variants 
		WHERE product_id = $1 AND is_active = true 
		ORDER BY sold_count DESC 
		LIMIT 1`

	var bestSeller types.ProductVariant
	var variantAttrsJSON, dimensionsJSON []byte

	err = r.db.QueryRowContext(ctx, bestSellerQuery, productID).Scan(
		&bestSeller.ID, &bestSeller.ProductID, &bestSeller.SKU, &bestSeller.Barcode,
		&bestSeller.Price, &bestSeller.CompareAtPrice, &bestSeller.CostPrice,
		&bestSeller.StockQuantity, &bestSeller.ReservedQuantity, &bestSeller.AvailableQuantity, &bestSeller.StockStatus, &bestSeller.LowStockThreshold,
		&variantAttrsJSON, &bestSeller.Weight, &dimensionsJSON,
		&bestSeller.IsActive, &bestSeller.IsDefault, &bestSeller.ViewCount, &bestSeller.SoldCount,
		&bestSeller.CreatedAt, &bestSeller.UpdatedAt,
	)
	if err == nil {
		// Unmarshal JSON fields
		if len(variantAttrsJSON) > 0 {
			if err := json.Unmarshal(variantAttrsJSON, &bestSeller.VariantAttributes); err != nil {
				// Log error but don't fail the entire operation
				log.Printf("Failed to unmarshal variant attributes: %v", err)
			}
		}
		if len(dimensionsJSON) > 0 {
			if err := json.Unmarshal(dimensionsJSON, &bestSeller.Dimensions); err != nil {
				// Log error but don't fail the entire operation
				log.Printf("Failed to unmarshal dimensions: %v", err)
			}
		}
		analytics.BestSeller = &bestSeller
	}

	// Find low stock variants
	lowStockQuery := `
		SELECT id, product_id, sku, barcode, price, compare_at_price, cost_price,
			   stock_quantity, reserved_quantity, available_quantity, stock_status, low_stock_threshold, variant_attributes,
			   weight, dimensions, is_active, is_default, view_count, sold_count,
			   created_at, updated_at
		FROM storefront_product_variants 
		WHERE product_id = $1 AND is_active = true 
		AND available_quantity <= COALESCE(low_stock_threshold, 5)
		ORDER BY available_quantity ASC`

	rows, err := r.db.QueryContext(ctx, lowStockQuery, productID)
	if err == nil {
		defer rows.Close()
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("error iterating rows: %w", err)
		}
		for rows.Next() {
			var variant types.ProductVariant
			var variantAttrsJSON, dimensionsJSON []byte

			err = rows.Scan(
				&variant.ID, &variant.ProductID, &variant.SKU, &variant.Barcode,
				&variant.Price, &variant.CompareAtPrice, &variant.CostPrice,
				&variant.StockQuantity, &variant.ReservedQuantity, &variant.AvailableQuantity, &variant.StockStatus, &variant.LowStockThreshold,
				&variantAttrsJSON, &variant.Weight, &dimensionsJSON,
				&variant.IsActive, &variant.IsDefault, &variant.ViewCount, &variant.SoldCount,
				&variant.CreatedAt, &variant.UpdatedAt,
			)
			if err != nil {
				break
			}

			// Unmarshal JSON fields
			if len(variantAttrsJSON) > 0 {
				if err := json.Unmarshal(variantAttrsJSON, &variant.VariantAttributes); err != nil {
					// Log error but don't fail the entire operation
					log.Printf("Failed to unmarshal variant attributes: %v", err)
				}
			}
			if len(dimensionsJSON) > 0 {
				if err := json.Unmarshal(dimensionsJSON, &variant.Dimensions); err != nil {
					// Log error but don't fail the entire operation
					log.Printf("Failed to unmarshal dimensions: %v", err)
				}
			}

			analytics.LowStockVariants = append(analytics.LowStockVariants, variant)
		}
	}

	// Initialize maps for attribute analytics
	analytics.StockByAttribute = make(map[string]map[string]int)
	analytics.SalesByAttribute = make(map[string]map[string]int)

	return &analytics, nil
}

// ImportVariants imports variants from CSV data
func (r *VariantRepository) ImportVariants(ctx context.Context, productID int, csvData []byte) (int, error) {
	// Parse CSV data
	reader := csv.NewReader(bytes.NewReader(csvData))
	records, err := reader.ReadAll()
	if err != nil {
		return 0, fmt.Errorf("failed to parse CSV: %w", err)
	}

	if len(records) < 2 {
		return 0, fmt.Errorf("CSV must contain at least header and one data row")
	}

	// Validate header (first row)
	header := records[0]
	expectedHeader := map[string]int{
		"sku":                 -1,
		"name":                -1,
		"price":               -1,
		"compare_at_price":    -1,
		"cost_price":          -1,
		"stock_quantity":      -1,
		"reserved_quantity":   -1,
		"low_stock_threshold": -1,
		"weight":              -1,
		"dimensions":          -1,
		"is_active":           -1,
		"is_default":          -1,
		"attributes":          -1,
	}

	// Map header columns to indices
	for i, col := range header {
		if _, exists := expectedHeader[col]; exists {
			expectedHeader[col] = i
		}
	}

	// Validate required columns are present
	requiredCols := []string{"sku", "stock_quantity", "attributes"}
	for _, col := range requiredCols {
		if expectedHeader[col] == -1 {
			return 0, fmt.Errorf("required column '%s' not found in CSV", col)
		}
	}

	// Begin transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback() // ignore error - transaction might already be committed
	}()

	importedCount := 0

	// Process data rows
	for i, record := range records[1:] {
		if len(record) != len(header) {
			return 0, fmt.Errorf("row %d has incorrect number of columns", i+2)
		}

		// Create variant from CSV row
		variant := &types.ProductVariant{
			ProductID: productID,
		}

		// Parse CSV fields
		if expectedHeader["sku"] != -1 && record[expectedHeader["sku"]] != "" {
			variant.SKU = &record[expectedHeader["sku"]]
		}

		if expectedHeader["name"] != -1 && record[expectedHeader["name"]] != "" {
			// Name не является отдельным полем, оно должно быть сгенерировано из атрибутов
			// или сохранено в variant_attributes как дополнительное поле
			if variant.VariantAttributes == nil {
				variant.VariantAttributes = make(map[string]interface{})
			}
			variant.VariantAttributes["name"] = record[expectedHeader["name"]]
		}

		if expectedHeader["price"] != -1 && record[expectedHeader["price"]] != "" {
			if price, err := strconv.ParseFloat(record[expectedHeader["price"]], 64); err == nil {
				variant.Price = &price
			}
		}

		if expectedHeader["compare_at_price"] != -1 && record[expectedHeader["compare_at_price"]] != "" {
			if price, err := strconv.ParseFloat(record[expectedHeader["compare_at_price"]], 64); err == nil {
				variant.CompareAtPrice = &price
			}
		}

		if expectedHeader["cost_price"] != -1 && record[expectedHeader["cost_price"]] != "" {
			if price, err := strconv.ParseFloat(record[expectedHeader["cost_price"]], 64); err == nil {
				variant.CostPrice = &price
			}
		}

		if expectedHeader["stock_quantity"] != -1 {
			if stock, err := strconv.Atoi(record[expectedHeader["stock_quantity"]]); err == nil {
				variant.StockQuantity = stock
			}
		}

		if expectedHeader["reserved_quantity"] != -1 && record[expectedHeader["reserved_quantity"]] != "" {
			if reserved, err := strconv.Atoi(record[expectedHeader["reserved_quantity"]]); err == nil {
				variant.ReservedQuantity = reserved
			}
		}

		if expectedHeader["low_stock_threshold"] != -1 && record[expectedHeader["low_stock_threshold"]] != "" {
			if threshold, err := strconv.Atoi(record[expectedHeader["low_stock_threshold"]]); err == nil {
				variant.LowStockThreshold = &threshold
			} else {
				defaultThreshold := 5
				variant.LowStockThreshold = &defaultThreshold // default
			}
		} else {
			defaultThreshold := 5
			variant.LowStockThreshold = &defaultThreshold // default
		}

		if expectedHeader["weight"] != -1 && record[expectedHeader["weight"]] != "" {
			if weight, err := strconv.ParseFloat(record[expectedHeader["weight"]], 64); err == nil {
				variant.Weight = &weight
			}
		}

		if expectedHeader["dimensions"] != -1 && record[expectedHeader["dimensions"]] != "" {
			var dimensions map[string]interface{}
			if err := json.Unmarshal([]byte(record[expectedHeader["dimensions"]]), &dimensions); err == nil {
				variant.Dimensions = dimensions
			}
		}

		if expectedHeader["is_active"] != -1 {
			const trueValue = "true"
			variant.IsActive = record[expectedHeader["is_active"]] == trueValue
		} else {
			variant.IsActive = true // default
		}

		if expectedHeader["is_default"] != -1 {
			variant.IsDefault = record[expectedHeader["is_default"]] == "true"
		}

		if expectedHeader["attributes"] != -1 && record[expectedHeader["attributes"]] != "" {
			var attributes map[string]interface{}
			if err := json.Unmarshal([]byte(record[expectedHeader["attributes"]]), &attributes); err != nil {
				return 0, fmt.Errorf("invalid JSON in attributes column for row %d: %w", i+2, err)
			}
			variant.VariantAttributes = attributes
		}

		// Create the variant - конвертируем в правильный тип
		createReq := &types.CreateVariantRequest{
			ProductID:         variant.ProductID,
			SKU:               variant.SKU,
			Barcode:           variant.Barcode,
			Price:             variant.Price,
			CompareAtPrice:    variant.CompareAtPrice,
			CostPrice:         variant.CostPrice,
			StockQuantity:     variant.StockQuantity,
			LowStockThreshold: variant.LowStockThreshold,
			VariantAttributes: variant.VariantAttributes,
			Weight:            variant.Weight,
			Dimensions:        variant.Dimensions,
			IsDefault:         variant.IsDefault,
		}
		_, err = r.CreateVariant(ctx, createReq)
		if err != nil {
			return 0, fmt.Errorf("failed to create variant from row %d: %w", i+2, err)
		}

		importedCount++
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return importedCount, nil
}

// ExportVariants exports variants to CSV format
func (r *VariantRepository) ExportVariants(ctx context.Context, productID int) ([]byte, string, error) {
	// Get all variants for the product
	variants, err := r.GetVariantsByProductID(ctx, productID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get variants: %w", err)
	}

	// Create CSV buffer
	var csvBuffer bytes.Buffer
	writer := csv.NewWriter(&csvBuffer)

	// Write CSV header
	header := []string{
		"id",
		"sku",
		"name",
		"price",
		"compare_at_price",
		"cost_price",
		"stock_quantity",
		"reserved_quantity",
		"low_stock_threshold",
		"weight",
		"dimensions",
		"is_active",
		"is_default",
		"attributes",
	}
	if err := writer.Write(header); err != nil {
		return nil, "", fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write variant data
	for _, variant := range variants {
		// Convert JSON fields to strings
		attributesJSON, _ := json.Marshal(variant.VariantAttributes)
		dimensionsJSON, _ := json.Marshal(variant.Dimensions)

		// Извлекаем name из variant_attributes если он есть
		variantName := ""
		if nameValue, ok := variant.VariantAttributes["name"]; ok {
			if nameStr, ok := nameValue.(string); ok {
				variantName = nameStr
			}
		}

		record := []string{
			strconv.Itoa(variant.ID),
			stringOrEmpty(variant.SKU),
			variantName,
			floatToString(variant.Price),
			floatToString(variant.CompareAtPrice),
			floatToString(variant.CostPrice),
			strconv.Itoa(variant.StockQuantity),
			strconv.Itoa(variant.ReservedQuantity),
			intToString(variant.LowStockThreshold),
			floatToString(variant.Weight),
			string(dimensionsJSON),
			boolToString(variant.IsActive),
			boolToString(variant.IsDefault),
			string(attributesJSON),
		}

		if err := writer.Write(record); err != nil {
			return nil, "", fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, "", fmt.Errorf("CSV writer error: %w", err)
	}

	fileName := fmt.Sprintf("variants-product-%d-%s.csv", productID, time.Now().Format("2006-01-02"))
	return csvBuffer.Bytes(), fileName, nil
}

// Helper functions for CSV export
func stringOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func floatToString(f *float64) string {
	if f == nil {
		return ""
	}
	return strconv.FormatFloat(*f, 'f', 2, 64)
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func intToString(i *int) string {
	if i == nil {
		return ""
	}
	return strconv.Itoa(*i)
}

// =====================================================
// PUBLIC METHODS (for unauthenticated buyers)
// =====================================================

// GetVariantsByProductIDPublic returns variants for public viewing (without sensitive data)
func (r *VariantRepository) GetVariantsByProductIDPublic(ctx context.Context, productID int) ([]types.ProductVariant, error) {
	type VariantRow struct {
		ID                int              `db:"id"`
		ProductID         int              `db:"product_id"`
		SKU               *string          `db:"sku"`
		Price             *float64         `db:"price"`
		CompareAtPrice    *float64         `db:"compare_at_price"`
		StockQuantity     int              `db:"stock_quantity"`
		ReservedQuantity  int              `db:"reserved_quantity"`
		AvailableQuantity int              `db:"available_quantity"`
		StockStatus       string           `db:"stock_status"`
		VariantAttributes json.RawMessage  `db:"variant_attributes"`
		Weight            *float64         `db:"weight"`
		Dimensions        *json.RawMessage `db:"dimensions"`
		IsActive          bool             `db:"is_active"`
		IsDefault         bool             `db:"is_default"`
		ViewCount         int              `db:"view_count"`
		SoldCount         int              `db:"sold_count"`
		CreatedAt         time.Time        `db:"created_at"`
		UpdatedAt         time.Time        `db:"updated_at"`
	}

	query := `
		SELECT 
			v.id, v.product_id, v.sku, v.price, v.compare_at_price,
			v.stock_quantity, v.reserved_quantity, v.available_quantity, v.stock_status,
			v.variant_attributes, v.weight, v.dimensions, v.is_active, v.is_default,
			v.view_count, v.sold_count, v.created_at, v.updated_at
		FROM storefront_product_variants v
		WHERE v.product_id = $1 AND v.is_active = true
		ORDER BY v.is_default DESC, v.id`

	var rows []VariantRow
	err := r.db.SelectContext(ctx, &rows, query, productID)
	if err != nil {
		log.Printf("Failed to query variants: %v", err)
		return nil, err
	}

	// Convert to ProductVariant structs
	variants := make([]types.ProductVariant, len(rows))
	for i, row := range rows {
		variants[i] = types.ProductVariant{
			ID:                row.ID,
			ProductID:         row.ProductID,
			SKU:               row.SKU,
			Price:             row.Price,
			CompareAtPrice:    row.CompareAtPrice,
			StockQuantity:     row.StockQuantity,
			ReservedQuantity:  row.ReservedQuantity,
			AvailableQuantity: row.AvailableQuantity,
			StockStatus:       row.StockStatus,
			Weight:            row.Weight,
			IsActive:          row.IsActive,
			IsDefault:         row.IsDefault,
			ViewCount:         row.ViewCount,
			SoldCount:         row.SoldCount,
			CreatedAt:         row.CreatedAt,
			UpdatedAt:         row.UpdatedAt,
			Images:            []types.ProductVariantImage{},
		}

		// Parse JSON fields
		if len(row.VariantAttributes) > 0 {
			if err := json.Unmarshal(row.VariantAttributes, &variants[i].VariantAttributes); err != nil {
				log.Printf("Failed to parse variant_attributes for variant %d: %v", row.ID, err)
				variants[i].VariantAttributes = make(map[string]interface{})
			}
		} else {
			variants[i].VariantAttributes = make(map[string]interface{})
		}

		if row.Dimensions != nil && len(*row.Dimensions) > 0 {
			if err := json.Unmarshal(*row.Dimensions, &variants[i].Dimensions); err != nil {
				log.Printf("Failed to parse dimensions for variant %d: %v", row.ID, err)
				variants[i].Dimensions = make(map[string]interface{})
			}
		} else {
			variants[i].Dimensions = make(map[string]interface{})
		}
	}

	return variants, nil
}

// GetVariantAttributesPublic returns all variant attributes for public viewing
func (r *VariantRepository) GetVariantAttributesPublic(ctx context.Context) ([]types.ProductVariantAttribute, error) {
	query := `
		SELECT id, name, display_name, type, is_required, sort_order, created_at, updated_at
		FROM product_variant_attributes
		ORDER BY sort_order, name`

	var attributes []types.ProductVariantAttribute
	err := r.db.SelectContext(ctx, &attributes, query)
	return attributes, err
}

// GetVariantAttributeValuesPublic returns values for a specific attribute for public viewing
func (r *VariantRepository) GetVariantAttributeValuesPublic(ctx context.Context, attributeID int) ([]types.ProductVariantAttributeValue, error) {
	query := `
		SELECT id, attribute_id, value, display_name, color_hex, image_url, sort_order, is_active,
		       is_popular, usage_count, created_at, updated_at
		FROM product_variant_attribute_values
		WHERE attribute_id = $1 AND is_active = true
		ORDER BY is_popular DESC, sort_order, display_name`

	var values []types.ProductVariantAttributeValue
	err := r.db.SelectContext(ctx, &values, query, attributeID)
	return values, err
}

// GetProductPublic returns basic product information for public viewing
func (r *VariantRepository) GetProductPublic(ctx context.Context, slug string, productID int) (*models.StorefrontProduct, error) {
	query := `
		SELECT 
			p.id, p.storefront_id, p.title, p.description, p.price, p.compare_at_price,
			p.category_id, p.is_active, p.is_featured, p.view_count, p.created_at, p.updated_at,
			s.slug as storefront_slug
		FROM storefront_products p
		JOIN storefronts s ON p.storefront_id = s.id
		WHERE p.id = $1 AND s.slug = $2 AND p.is_active = true`

	var product models.StorefrontProduct
	err := r.db.GetContext(ctx, &product, query, productID, slug)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

// GetVariantByIDPublic returns a specific variant for public viewing
func (r *VariantRepository) GetVariantByIDPublic(ctx context.Context, variantID int) (*types.ProductVariant, error) {
	query := `
		SELECT 
			v.id, v.product_id, v.sku, v.price, v.compare_at_price,
			v.stock_quantity, v.reserved_quantity, v.available_quantity, v.stock_status,
			v.variant_attributes, v.weight, v.dimensions, v.is_active, v.is_default,
			v.view_count, v.sold_count, v.created_at, v.updated_at
		FROM storefront_product_variants v
		WHERE v.id = $1 AND v.is_active = true`

	var variant types.ProductVariant
	err := r.db.GetContext(ctx, &variant, query, variantID)
	if err != nil {
		return nil, err
	}

	// Get images for the variant
	images, err := r.getVariantImages(ctx, variant.ID)
	if err != nil {
		log.Printf("Failed to get images for variant %d: %v", variant.ID, err)
	} else {
		variant.Images = images
	}

	return &variant, nil
}

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

// BulkUpdateProducts updates multiple products atomically within a transaction
// Max 1000 items, processes in batches of 100 for performance
// Returns successfully updated products and errors for failed ones
func (r *Repository) BulkUpdateProducts(ctx context.Context, storefrontID int64, updates []*domain.BulkUpdateProductInput) (*domain.BulkUpdateProductsResult, error) {
	r.logger.Debug().
		Int64("storefront_id", storefrontID).
		Int("update_count", len(updates)).
		Msg("bulk updating products")

	// Validation
	if len(updates) == 0 {
		return &domain.BulkUpdateProductsResult{
			SuccessfulProducts: []*domain.Product{},
			FailedUpdates:      []domain.BulkUpdateError{},
		}, nil
	}

	if len(updates) > 1000 {
		return nil, fmt.Errorf("products.bulk_update_limit_exceeded")
	}

	// Start transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to begin transaction")
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Step 1: Validate ownership for all products
	productIDs := make([]int64, len(updates))
	for i, update := range updates {
		productIDs[i] = update.ProductID
	}

	ownershipQuery := `
		SELECT id
		FROM b2c_products
		WHERE id = ANY($1::bigint[]) AND storefront_id = $2
	`

	var ownedIDs []int64
	err = tx.SelectContext(ctx, &ownedIDs, ownershipQuery, pq.Array(productIDs), storefrontID)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to validate product ownership")
		return nil, fmt.Errorf("failed to validate product ownership: %w", err)
	}

	// Create map for quick lookup
	ownedMap := make(map[int64]bool)
	for _, id := range ownedIDs {
		ownedMap[id] = true
	}

	// Initialize result
	result := &domain.BulkUpdateProductsResult{
		SuccessfulProducts: make([]*domain.Product, 0),
		FailedUpdates:      make([]domain.BulkUpdateError, 0),
	}

	// Process updates in batches of 100
	batchSize := 100
	for batchStart := 0; batchStart < len(updates); batchStart += batchSize {
		batchEnd := batchStart + batchSize
		if batchEnd > len(updates) {
			batchEnd = len(updates)
		}

		batch := updates[batchStart:batchEnd]

		// Process each product in the batch
		for _, update := range batch {
			// Check ownership
			if !ownedMap[update.ProductID] {
				result.FailedUpdates = append(result.FailedUpdates, domain.BulkUpdateError{
					ProductID:    update.ProductID,
					ErrorCode:    "products.not_found",
					ErrorMessage: "Product not found or access denied",
				})
				continue
			}

			// Build dynamic UPDATE query based on update_mask or all non-nil fields
			setClauses := []string{"updated_at = NOW()"}
			args := []interface{}{}
			argIndex := 1

			// Helper to check if field is in update_mask
			shouldUpdate := func(fieldName string) bool {
				if len(update.UpdateMask) == 0 {
					return true // No mask = update all non-nil fields
				}
				for _, maskField := range update.UpdateMask {
					if maskField == fieldName {
						return true
					}
				}
				return false
			}

			// Add fields based on update_mask
			if shouldUpdate("name") && update.Name != nil {
				setClauses = append(setClauses, fmt.Sprintf("name = $%d", argIndex))
				args = append(args, *update.Name)
				argIndex++
			}

			if shouldUpdate("description") && update.Description != nil {
				setClauses = append(setClauses, fmt.Sprintf("description = $%d", argIndex))
				args = append(args, *update.Description)
				argIndex++
			}

			if shouldUpdate("price") && update.Price != nil {
				setClauses = append(setClauses, fmt.Sprintf("price = $%d", argIndex))
				args = append(args, *update.Price)
				argIndex++
			}

			if shouldUpdate("sku") && update.SKU != nil {
				setClauses = append(setClauses, fmt.Sprintf("sku = $%d", argIndex))
				args = append(args, *update.SKU)
				argIndex++
			}

			if shouldUpdate("barcode") && update.Barcode != nil {
				setClauses = append(setClauses, fmt.Sprintf("barcode = $%d", argIndex))
				args = append(args, *update.Barcode)
				argIndex++
			}

			if shouldUpdate("stock_quantity") && update.StockQuantity != nil {
				setClauses = append(setClauses, fmt.Sprintf("stock_quantity = $%d", argIndex))
				args = append(args, *update.StockQuantity)
				argIndex++

				// Auto-update stock_status based on quantity
				if *update.StockQuantity > 0 {
					setClauses = append(setClauses, "stock_status = 'in_stock'")
				} else {
					setClauses = append(setClauses, "stock_status = 'out_of_stock'")
				}
			}

			if shouldUpdate("stock_status") && update.StockStatus != nil {
				setClauses = append(setClauses, fmt.Sprintf("stock_status = $%d", argIndex))
				args = append(args, *update.StockStatus)
				argIndex++
			}

			if shouldUpdate("is_active") && update.IsActive != nil {
				setClauses = append(setClauses, fmt.Sprintf("is_active = $%d", argIndex))
				args = append(args, *update.IsActive)
				argIndex++
			}

			if shouldUpdate("attributes") && update.Attributes != nil {
				attributesJSON, err := json.Marshal(update.Attributes)
				if err != nil {
					result.FailedUpdates = append(result.FailedUpdates, domain.BulkUpdateError{
						ProductID:    update.ProductID,
						ErrorCode:    "products.invalid_attributes",
						ErrorMessage: "Failed to marshal attributes",
					})
					continue
				}
				setClauses = append(setClauses, fmt.Sprintf("attributes = $%d", argIndex))
				args = append(args, attributesJSON)
				argIndex++
			}

			if shouldUpdate("has_individual_location") && update.HasIndividualLocation != nil {
				setClauses = append(setClauses, fmt.Sprintf("has_individual_location = $%d", argIndex))
				args = append(args, *update.HasIndividualLocation)
				argIndex++
			}

			if shouldUpdate("individual_address") && update.IndividualAddress != nil {
				setClauses = append(setClauses, fmt.Sprintf("individual_address = $%d", argIndex))
				args = append(args, *update.IndividualAddress)
				argIndex++
			}

			if shouldUpdate("individual_latitude") && update.IndividualLatitude != nil {
				setClauses = append(setClauses, fmt.Sprintf("individual_latitude = $%d", argIndex))
				args = append(args, *update.IndividualLatitude)
				argIndex++
			}

			if shouldUpdate("individual_longitude") && update.IndividualLongitude != nil {
				setClauses = append(setClauses, fmt.Sprintf("individual_longitude = $%d", argIndex))
				args = append(args, *update.IndividualLongitude)
				argIndex++
			}

			if shouldUpdate("location_privacy") && update.LocationPrivacy != nil {
				setClauses = append(setClauses, fmt.Sprintf("location_privacy = $%d", argIndex))
				args = append(args, *update.LocationPrivacy)
				argIndex++
			}

			if shouldUpdate("show_on_map") && update.ShowOnMap != nil {
				setClauses = append(setClauses, fmt.Sprintf("show_on_map = $%d", argIndex))
				args = append(args, *update.ShowOnMap)
				argIndex++
			}

			// If no fields to update (only updated_at), skip
			if len(setClauses) == 1 {
				result.FailedUpdates = append(result.FailedUpdates, domain.BulkUpdateError{
					ProductID:    update.ProductID,
					ErrorCode:    "products.no_fields_to_update",
					ErrorMessage: "No fields provided for update",
				})
				continue
			}

			// Add WHERE clause
			args = append(args, update.ProductID, storefrontID)
			whereProductID := argIndex
			whereStorefrontID := argIndex + 1

			// Build final query with RETURNING
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

			// Execute update and scan result
			var product domain.Product
			var description sql.NullString
			var sku, barcode sql.NullString
			var individualAddress, locationPrivacy sql.NullString
			var individualLatitude, individualLongitude sql.NullFloat64
			var returnedAttributesJSON []byte

			err = tx.QueryRowContext(ctx, query, args...).Scan(
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
					result.FailedUpdates = append(result.FailedUpdates, domain.BulkUpdateError{
						ProductID:    update.ProductID,
						ErrorCode:    "products.not_found",
						ErrorMessage: "Product not found",
					})
					continue
				}

				// Check for unique constraint violation (duplicate SKU)
				if pqErr, ok := err.(*pq.Error); ok {
					if pqErr.Code == "23505" { // unique_violation
						result.FailedUpdates = append(result.FailedUpdates, domain.BulkUpdateError{
							ProductID:    update.ProductID,
							ErrorCode:    "products.sku_duplicate",
							ErrorMessage: "SKU already exists",
						})
						continue
					}
				}

				result.FailedUpdates = append(result.FailedUpdates, domain.BulkUpdateError{
					ProductID:    update.ProductID,
					ErrorCode:    "products.update_failed",
					ErrorMessage: fmt.Sprintf("Update failed: %v", err),
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
					r.logger.Error().Err(err).Int64("product_id", product.ID).Msg("failed to unmarshal attributes")
					// Don't fail the update, just log the error
				}
			}

			// Add to successful results
			result.SuccessfulProducts = append(result.SuccessfulProducts, &product)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		r.logger.Error().Err(err).Msg("failed to commit transaction")
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	r.logger.Info().
		Int("successful_count", len(result.SuccessfulProducts)).
		Int("failed_count", len(result.FailedUpdates)).
		Msg("bulk update products completed")

	return result, nil
}

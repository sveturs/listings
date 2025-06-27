package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"backend/internal/domain/models"
	"github.com/lib/pq"
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
			c.id, c.name, c.slug, c.icon, c.parent_id
		FROM storefront_products p
		LEFT JOIN marketplace_categories c ON p.category_id = c.id
		WHERE p.storefront_id = $1`

	args := []interface{}{filter.StorefrontID}
	argIndex := 2

	// Apply filters
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

	if filter.IsActive != nil {
		query += fmt.Sprintf(" AND p.is_active = $%d", argIndex)
		args = append(args, *filter.IsActive)
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
		case "name":
			sortBy = "p.name"
		case "price":
			sortBy = "p.price"
		case "stock_quantity":
			sortBy = "p.stock_quantity"
		case "created_at":
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
			c.Icon = categoryIcon.String
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
		&categoryID, &categoryName, &categorySlug, &categoryIcon, &categoryParentID,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get storefront product: %w", err)
	}

	if attributesJSON != nil {
		if err := json.Unmarshal(attributesJSON, &p.Attributes); err != nil {
			return nil, fmt.Errorf("failed to unmarshal attributes: %w", err)
		}
	}

	// Only set category if it exists
	if categoryID.Valid {
		c.ID = int(categoryID.Int64)
		c.Name = categoryName.String
		c.Slug = categorySlug.String
		c.Icon = categoryIcon.String
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
	p.Images = images

	// Load variants
	variants, err := s.getProductVariants(ctx, p.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product variants: %w", err)
	}
	p.Variants = variants

	return p, nil
}

// CreateStorefrontProduct creates a new product
func (s *Database) CreateStorefrontProduct(ctx context.Context, storefrontID int, req *models.CreateProductRequest) (*models.StorefrontProduct, error) {
	var attributesJSON []byte
	if req.Attributes != nil {
		var err error
		attributesJSON, err = json.Marshal(req.Attributes)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal attributes: %w", err)
		}
	}

	query := `
		INSERT INTO storefront_products (
			storefront_id, name, description, price, currency, category_id,
			sku, barcode, stock_quantity, is_active, attributes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, stock_status, created_at, updated_at`

	var product models.StorefrontProduct
	err := s.pool.QueryRow(
		ctx, query,
		storefrontID, req.Name, req.Description, req.Price, req.Currency, req.CategoryID,
		req.SKU, req.Barcode, req.StockQuantity, req.IsActive, attributesJSON,
	).Scan(&product.ID, &product.StockStatus, &product.CreatedAt, &product.UpdatedAt)

	if err != nil {
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

	return &product, nil
}

// UpdateStorefrontProduct updates an existing product
func (s *Database) UpdateStorefrontProduct(ctx context.Context, storefrontID, productID int, req *models.UpdateProductRequest) error {
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
		args = append(args, *req.SKU)
		argIndex++
	}

	if req.Barcode != nil {
		setClauses = append(setClauses, fmt.Sprintf("barcode = $%d", argIndex))
		args = append(args, *req.Barcode)
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

	result, err := s.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update storefront product: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

// DeleteStorefrontProduct deletes a product
func (s *Database) DeleteStorefrontProduct(ctx context.Context, storefrontID, productID int) error {
	query := `DELETE FROM storefront_products WHERE id = $1 AND storefront_id = $2`

	result, err := s.pool.Exec(ctx, query, productID, storefrontID)
	if err != nil {
		return fmt.Errorf("failed to delete storefront product: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

// UpdateProductInventory updates product stock with tracking
func (s *Database) UpdateProductInventory(ctx context.Context, storefrontID, productID int, userID int, req *models.UpdateInventoryRequest) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

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

	query := `
		SELECT id, storefront_product_id, image_url, thumbnail_url, display_order, is_default, created_at
		FROM storefront_product_images
		WHERE storefront_product_id = ANY($1)
		ORDER BY display_order, id`

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
			&img.DisplayOrder, &img.IsDefault, &img.CreatedAt,
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
		SELECT id, storefront_product_id, name, sku, price, stock_quantity, attributes, is_active, created_at, updated_at
		FROM storefront_product_variants
		WHERE storefront_product_id = $1
		ORDER BY id`

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
	defer tx.Rollback(ctx)

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
	defer tx.Rollback(ctx)

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
	defer tx.Rollback(ctx)

	// Delete products that belong to the storefront
	rows, err := tx.Query(ctx,
		`DELETE FROM storefront_products 
		WHERE id = ANY($1) AND storefront_id = $2 
		RETURNING id`,
		pq.Array(productIDs), storefrontID,
	)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to delete products: %w", err)}
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err == nil {
			deletedIDs = append(deletedIDs, id)
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

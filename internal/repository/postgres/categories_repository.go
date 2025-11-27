package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/gosimple/slug"
	"github.com/lib/pq"
	"github.com/vondi-global/listings/internal/domain"
)

// GetRootCategories retrieves all top-level categories (no parent)
func (r *Repository) GetRootCategories(ctx context.Context) ([]*domain.Category, error) {
	query := `
		SELECT
			id, name, slug, parent_id, icon, description,
			is_active, created_at, sort_order, level,
			has_custom_ui, custom_ui_component,
			COALESCE((SELECT COUNT(*) FROM listings WHERE category_id = mc.id AND status = 'active'), 0) as listing_count
		FROM categories mc
		WHERE is_active = true AND parent_id IS NULL
		ORDER BY sort_order ASC, name ASC
	`

	rows, err := r.db.QueryxContext(ctx, query)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to query root categories")
		return nil, fmt.Errorf("failed to query root categories: %w", err)
	}
	defer rows.Close()

	var categories []*domain.Category
	for rows.Next() {
		cat := &domain.Category{}
		var icon, description, customUIComponent sql.NullString
		var parentID sql.NullInt64

		err := rows.Scan(
			&cat.ID,
			&cat.Name,
			&cat.Slug,
			&parentID,
			&icon,
			&description,
			&cat.IsActive,
			&cat.CreatedAt,
			&cat.SortOrder,
			&cat.Level,
			&cat.HasCustomUI,
			&customUIComponent,
			&cat.ListingCount,
		)
		if err != nil {
			r.logger.Error().Err(err).Msg("failed to scan category")
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}

		// Handle nullable fields
		if parentID.Valid {
			cat.ParentID = &parentID.Int64
		}
		if icon.Valid {
			cat.Icon = &icon.String
		}
		if description.Valid {
			cat.Description = &description.String
		}
		if customUIComponent.Valid {
			cat.CustomUIComponent = &customUIComponent.String
		}

		categories = append(categories, cat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return categories, nil
}

// GetAllCategories retrieves all active categories
func (r *Repository) GetAllCategories(ctx context.Context) ([]*domain.Category, error) {
	query := `
		SELECT
			id, name, slug, parent_id, icon, description,
			is_active, created_at, sort_order, level,
			has_custom_ui, custom_ui_component,
			COALESCE((SELECT COUNT(*) FROM listings WHERE category_id = mc.id AND status = 'active'), 0) as listing_count
		FROM categories mc
		WHERE is_active = true
		ORDER BY sort_order ASC, name ASC
	`

	rows, err := r.db.QueryxContext(ctx, query)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to query categories")
		return nil, fmt.Errorf("failed to query categories: %w", err)
	}
	defer rows.Close()

	var categories []*domain.Category
	for rows.Next() {
		cat := &domain.Category{}
		var icon, description, customUIComponent sql.NullString
		var parentID sql.NullInt64

		err := rows.Scan(
			&cat.ID,
			&cat.Name,
			&cat.Slug,
			&parentID,
			&icon,
			&description,
			&cat.IsActive,
			&cat.CreatedAt,
			&cat.SortOrder,
			&cat.Level,
			&cat.HasCustomUI,
			&customUIComponent,
			&cat.ListingCount,
		)
		if err != nil {
			r.logger.Error().Err(err).Msg("failed to scan category")
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}

		// Handle nullable fields
		if parentID.Valid {
			cat.ParentID = &parentID.Int64
		}
		if icon.Valid {
			cat.Icon = &icon.String
		}
		if description.Valid {
			cat.Description = &description.String
		}
		if customUIComponent.Valid {
			cat.CustomUIComponent = &customUIComponent.String
		}

		categories = append(categories, cat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return categories, nil
}

// GetPopularCategories retrieves top N categories by listing count
func (r *Repository) GetPopularCategories(ctx context.Context, limit int) ([]*domain.Category, error) {
	if limit <= 0 {
		limit = 10 // default limit
	}

	query := `
		SELECT
			mc.id, mc.name, mc.slug, mc.parent_id, mc.icon, mc.description,
			mc.is_active, mc.created_at, mc.sort_order, mc.level,
			mc.has_custom_ui, mc.custom_ui_component,
			COUNT(l.id) as listing_count
		FROM categories mc
		LEFT JOIN listings l ON l.category_id = mc.id AND l.status = 'active' AND l.is_deleted = false
		WHERE mc.is_active = true
		GROUP BY mc.id, mc.name, mc.slug, mc.parent_id, mc.icon, mc.description,
				 mc.is_active, mc.created_at, mc.sort_order, mc.level,
				 mc.has_custom_ui, mc.custom_ui_component
		ORDER BY listing_count DESC, mc.name ASC
		LIMIT $1
	`

	rows, err := r.db.QueryxContext(ctx, query, limit)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to query popular categories")
		return nil, fmt.Errorf("failed to query popular categories: %w", err)
	}
	defer rows.Close()

	var categories []*domain.Category
	for rows.Next() {
		cat := &domain.Category{}
		var icon, description, customUIComponent sql.NullString
		var parentID sql.NullInt64
		var listingCount int64

		err := rows.Scan(
			&cat.ID,
			&cat.Name,
			&cat.Slug,
			&parentID,
			&icon,
			&description,
			&cat.IsActive,
			&cat.CreatedAt,
			&cat.SortOrder,
			&cat.Level,
			&cat.HasCustomUI,
			&customUIComponent,
			&listingCount,
		)
		if err != nil {
			r.logger.Error().Err(err).Msg("failed to scan popular category")
			return nil, fmt.Errorf("failed to scan popular category: %w", err)
		}

		cat.ListingCount = int32(listingCount)

		// Handle nullable fields
		if parentID.Valid {
			cat.ParentID = &parentID.Int64
		}
		if icon.Valid {
			cat.Icon = &icon.String
		}
		if description.Valid {
			cat.Description = &description.String
		}
		if customUIComponent.Valid {
			cat.CustomUIComponent = &customUIComponent.String
		}

		categories = append(categories, cat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return categories, nil
}

// GetCategoryByID retrieves a single category by ID
func (r *Repository) GetCategoryByID(ctx context.Context, categoryID int64) (*domain.Category, error) {
	query := `
		SELECT
			id, name, slug, parent_id, icon, description,
			is_active, created_at, sort_order, level,
			has_custom_ui, custom_ui_component,
			COALESCE((SELECT COUNT(*) FROM listings WHERE category_id = mc.id AND status = 'active'), 0) as listing_count
		FROM categories mc
		WHERE id = $1
	`

	cat := &domain.Category{}
	var icon, description, customUIComponent sql.NullString
	var parentID sql.NullInt64

	err := r.db.QueryRowxContext(ctx, query, categoryID).Scan(
		&cat.ID,
		&cat.Name,
		&cat.Slug,
		&parentID,
		&icon,
		&description,
		&cat.IsActive,
		&cat.CreatedAt,
		&cat.SortOrder,
		&cat.Level,
		&cat.HasCustomUI,
		&customUIComponent,
		&cat.ListingCount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("category not found")
		}
		r.logger.Error().Err(err).Int64("category_id", categoryID).Msg("failed to get category")
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	// Handle nullable fields
	if parentID.Valid {
		cat.ParentID = &parentID.Int64
	}
	if icon.Valid {
		cat.Icon = &icon.String
	}
	if description.Valid {
		cat.Description = &description.String
	}
	if customUIComponent.Valid {
		cat.CustomUIComponent = &customUIComponent.String
	}

	return cat, nil
}

// GetCategoryTree builds a hierarchical tree starting from a category
func (r *Repository) GetCategoryTree(ctx context.Context, categoryID int64) (*domain.CategoryTreeNode, error) {
	// Get all categories
	categories, err := r.GetAllCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	// Build map for fast lookup
	catMap := make(map[int64]*domain.CategoryTreeNode)
	var rootNode *domain.CategoryTreeNode

	// First pass: create nodes
	for _, cat := range categories {
		node := &domain.CategoryTreeNode{
			ID:           cat.ID,
			Name:         cat.Name,
			Slug:         cat.Slug,
			ParentID:     cat.ParentID,
			Level:        cat.Level,
			ListingCount: cat.ListingCount,
			Children:     []domain.CategoryTreeNode{},
			HasCustomUI:  cat.HasCustomUI,
			CreatedAt:    cat.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}

		// Handle nullable fields
		if cat.Icon != nil {
			node.Icon = cat.Icon
		}
		if cat.CustomUIComponent != nil {
			node.CustomUIComponent = cat.CustomUIComponent
		}

		// Build path (simple implementation - just slug for now)
		node.Path = cat.Slug

		catMap[cat.ID] = node

		// If this is the requested category, mark it as root
		if cat.ID == categoryID {
			rootNode = node
		}
	}

	// Second pass: build tree (link children to parents)
	for _, node := range catMap {
		if node.ParentID != nil {
			parent, exists := catMap[*node.ParentID]
			if exists {
				parent.Children = append(parent.Children, *node)
			}
		}
	}

	// Third pass: update ChildrenCount
	for _, node := range catMap {
		node.ChildrenCount = int32(len(node.Children))
	}

	if rootNode == nil {
		return nil, fmt.Errorf("category not found in tree")
	}

	return rootNode, nil
}

// GetCategoryBySlug retrieves a category by its slug
func (r *Repository) GetCategoryBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	query := `
		SELECT
			id, name, slug, parent_id, icon, description,
			is_active, created_at, sort_order, level,
			has_custom_ui, custom_ui_component,
			COALESCE((SELECT COUNT(*) FROM listings WHERE category_id = mc.id AND status = 'active'), 0) as listing_count
		FROM categories mc
		WHERE slug = $1
	`

	cat := &domain.Category{}
	var icon, description, customUIComponent sql.NullString
	var parentID sql.NullInt64

	err := r.db.QueryRowxContext(ctx, query, slug).Scan(
		&cat.ID,
		&cat.Name,
		&cat.Slug,
		&parentID,
		&icon,
		&description,
		&cat.IsActive,
		&cat.CreatedAt,
		&cat.SortOrder,
		&cat.Level,
		&cat.HasCustomUI,
		&customUIComponent,
		&cat.ListingCount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("category not found")
		}
		r.logger.Error().Err(err).Str("slug", slug).Msg("failed to get category by slug")
		return nil, fmt.Errorf("failed to get category by slug: %w", err)
	}

	// Handle nullable fields
	if parentID.Valid {
		cat.ParentID = &parentID.Int64
	}
	if icon.Valid {
		cat.Icon = &icon.String
	}
	if description.Valid {
		cat.Description = &description.String
	}
	if customUIComponent.Valid {
		cat.CustomUIComponent = &customUIComponent.String
	}

	return cat, nil
}

// GetCategoriesWithPagination returns paginated categories with optional filters
func (r *Repository) GetCategoriesWithPagination(ctx context.Context, parentID *int64, isActive *bool, limit, offset int32) ([]*domain.Category, int32, error) {
	// Build query with conditional filters
	query := `
		SELECT
			id, name, slug, parent_id, icon, description,
			is_active, created_at, sort_order, level,
			has_custom_ui, custom_ui_component,
			COALESCE((SELECT COUNT(*) FROM listings WHERE category_id = mc.id AND status = 'active'), 0) as listing_count
		FROM categories mc
		WHERE 1=1
	`

	countQuery := `SELECT COUNT(*) FROM categories WHERE 1=1`

	args := []interface{}{}
	argCount := 1

	// Add filters
	if parentID != nil {
		query += fmt.Sprintf(" AND parent_id = $%d", argCount)
		countQuery += fmt.Sprintf(" AND parent_id = $%d", argCount)
		args = append(args, *parentID)
		argCount++
	} else if parentID == nil {
		// Explicitly handle NULL parent_id
		query += " AND parent_id IS NULL"
		countQuery += " AND parent_id IS NULL"
	}

	if isActive != nil {
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		countQuery += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, *isActive)
		argCount++
	}

	// Get total count
	var totalCount int32
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to get categories count")
		return nil, 0, fmt.Errorf("failed to get categories count: %w", err)
	}

	// Add sorting and pagination
	query += ` ORDER BY sort_order ASC, name ASC`
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to query categories with pagination")
		return nil, 0, fmt.Errorf("failed to query categories: %w", err)
	}
	defer rows.Close()

	var categories []*domain.Category
	for rows.Next() {
		cat := &domain.Category{}
		var icon, description, customUIComponent sql.NullString
		var parentIDVal sql.NullInt64

		err := rows.Scan(
			&cat.ID,
			&cat.Name,
			&cat.Slug,
			&parentIDVal,
			&icon,
			&description,
			&cat.IsActive,
			&cat.CreatedAt,
			&cat.SortOrder,
			&cat.Level,
			&cat.HasCustomUI,
			&customUIComponent,
			&cat.ListingCount,
		)
		if err != nil {
			r.logger.Error().Err(err).Msg("failed to scan category")
			return nil, 0, fmt.Errorf("failed to scan category: %w", err)
		}

		// Handle nullable fields
		if parentIDVal.Valid {
			cat.ParentID = &parentIDVal.Int64
		}
		if icon.Valid {
			cat.Icon = &icon.String
		}
		if description.Valid {
			cat.Description = &description.String
		}
		if customUIComponent.Valid {
			cat.CustomUIComponent = &customUIComponent.String
		}

		categories = append(categories, cat)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows iteration error: %w", err)
	}

	return categories, totalCount, nil
}

// CreateCategory creates a new category
func (r *Repository) CreateCategory(ctx context.Context, cat *domain.Category) (*domain.Category, error) {
	// Generate slug if not provided
	if cat.Slug == "" {
		cat.Slug = slug.Make(cat.Name)
	}

	// Calculate level based on parent
	level := int32(0)
	if cat.ParentID != nil {
		parent, err := r.GetCategoryByID(ctx, *cat.ParentID)
		if err != nil {
			return nil, fmt.Errorf("invalid parent_id: %w", err)
		}
		level = parent.Level + 1
	}

	// Get max sort_order if not provided
	if cat.SortOrder == 0 {
		var maxSortOrder sql.NullInt32
		query := "SELECT MAX(sort_order) FROM categories"
		err := r.db.QueryRowContext(ctx, query).Scan(&maxSortOrder)
		if err != nil && err != sql.ErrNoRows {
			r.logger.Error().Err(err).Msg("failed to get max sort_order")
			return nil, fmt.Errorf("failed to get max sort_order: %w", err)
		}
		if maxSortOrder.Valid {
			cat.SortOrder = maxSortOrder.Int32 + 1
		} else {
			cat.SortOrder = 1
		}
	}

	cat.Level = level
	cat.CreatedAt = time.Now()

	query := `
		INSERT INTO categories (
			name, slug, parent_id, icon, description,
			is_active, sort_order, level, has_custom_ui, custom_ui_component,
			created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		)
		RETURNING id
	`

	var parentID sql.NullInt64
	if cat.ParentID != nil {
		parentID = sql.NullInt64{Int64: *cat.ParentID, Valid: true}
	}

	var icon, description, customUIComponent sql.NullString
	if cat.Icon != nil {
		icon = sql.NullString{String: *cat.Icon, Valid: true}
	}
	if cat.Description != nil {
		description = sql.NullString{String: *cat.Description, Valid: true}
	}
	if cat.CustomUIComponent != nil {
		customUIComponent = sql.NullString{String: *cat.CustomUIComponent, Valid: true}
	}

	err := r.db.QueryRowContext(ctx, query,
		cat.Name,
		cat.Slug,
		parentID,
		icon,
		description,
		cat.IsActive,
		cat.SortOrder,
		cat.Level,
		cat.HasCustomUI,
		customUIComponent,
		cat.CreatedAt,
	).Scan(&cat.ID)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // unique_violation
				return nil, fmt.Errorf("category with slug '%s' already exists", cat.Slug)
			}
		}
		r.logger.Error().Err(err).Msg("failed to create category")
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	r.logger.Info().Int64("category_id", cat.ID).Str("name", cat.Name).Msg("category created")
	return cat, nil
}

// UpdateCategory updates an existing category (partial updates supported)
func (r *Repository) UpdateCategory(ctx context.Context, cat *domain.Category) (*domain.Category, error) {
	// Validate category exists
	existing, err := r.GetCategoryByID(ctx, cat.ID)
	if err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}

	// Check for circular dependency if parent is being changed
	if cat.ParentID != nil {
		if *cat.ParentID == cat.ID {
			return nil, fmt.Errorf("category cannot be its own parent")
		}

		// Check if new parent is a descendant of this category
		isDescendant, err := r.isDescendant(ctx, cat.ID, *cat.ParentID)
		if err != nil {
			return nil, fmt.Errorf("failed to check circular dependency: %w", err)
		}
		if isDescendant {
			return nil, fmt.Errorf("circular dependency detected: cannot set descendant as parent")
		}

		// Recalculate level if parent changed
		if existing.ParentID == nil || *existing.ParentID != *cat.ParentID {
			parent, err := r.GetCategoryByID(ctx, *cat.ParentID)
			if err != nil {
				return nil, fmt.Errorf("invalid parent_id: %w", err)
			}
			cat.Level = parent.Level + 1

			// Update levels of all descendants
			err = r.updateDescendantLevels(ctx, cat.ID, cat.Level)
			if err != nil {
				return nil, fmt.Errorf("failed to update descendant levels: %w", err)
			}
		}
	} else if cat.ParentID == nil && existing.ParentID != nil {
		// Moving to root
		cat.Level = 0
		err = r.updateDescendantLevels(ctx, cat.ID, cat.Level)
		if err != nil {
			return nil, fmt.Errorf("failed to update descendant levels: %w", err)
		}
	}

	query := `
		UPDATE categories SET
			name = COALESCE($2, name),
			slug = COALESCE($3, slug),
			parent_id = $4,
			icon = $5,
			description = $6,
			is_active = COALESCE($7, is_active),
			sort_order = COALESCE($8, sort_order),
			level = COALESCE($9, level),
			has_custom_ui = COALESCE($10, has_custom_ui),
			custom_ui_component = $11
		WHERE id = $1
		RETURNING id, name, slug, parent_id, icon, description, is_active, created_at,
				  sort_order, level, has_custom_ui, custom_ui_component
	`

	var parentID sql.NullInt64
	if cat.ParentID != nil {
		parentID = sql.NullInt64{Int64: *cat.ParentID, Valid: true}
	}

	var icon, description, customUIComponent sql.NullString
	if cat.Icon != nil {
		icon = sql.NullString{String: *cat.Icon, Valid: true}
	}
	if cat.Description != nil {
		description = sql.NullString{String: *cat.Description, Valid: true}
	}
	if cat.CustomUIComponent != nil {
		customUIComponent = sql.NullString{String: *cat.CustomUIComponent, Valid: true}
	}

	// Prepare nullable values for COALESCE
	var name, slugVal *string
	var isActive *bool
	var sortOrder, level *int32
	var hasCustomUI *bool

	if cat.Name != "" {
		name = &cat.Name
	}
	if cat.Slug != "" {
		slugVal = &cat.Slug
	}
	isActive = &cat.IsActive
	if cat.SortOrder != 0 {
		sortOrder = &cat.SortOrder
	}
	if cat.Level != 0 {
		level = &cat.Level
	}
	hasCustomUI = &cat.HasCustomUI

	updatedCat := &domain.Category{}
	err = r.db.QueryRowContext(ctx, query,
		cat.ID,
		name,
		slugVal,
		parentID,
		icon,
		description,
		isActive,
		sortOrder,
		level,
		hasCustomUI,
		customUIComponent,
	).Scan(
		&updatedCat.ID,
		&updatedCat.Name,
		&updatedCat.Slug,
		&parentID,
		&icon,
		&description,
		&updatedCat.IsActive,
		&updatedCat.CreatedAt,
		&updatedCat.SortOrder,
		&updatedCat.Level,
		&updatedCat.HasCustomUI,
		&customUIComponent,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // unique_violation
				return nil, fmt.Errorf("category with slug '%s' already exists", cat.Slug)
			}
		}
		r.logger.Error().Err(err).Int64("category_id", cat.ID).Msg("failed to update category")
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	// Handle nullable fields
	if parentID.Valid {
		updatedCat.ParentID = &parentID.Int64
	}
	if icon.Valid {
		updatedCat.Icon = &icon.String
	}
	if description.Valid {
		updatedCat.Description = &description.String
	}
	if customUIComponent.Valid {
		updatedCat.CustomUIComponent = &customUIComponent.String
	}

	r.logger.Info().Int64("category_id", cat.ID).Msg("category updated")
	return updatedCat, nil
}

// DeleteCategory soft deletes a category (sets is_active = false)
func (r *Repository) DeleteCategory(ctx context.Context, categoryID int64) error {
	// Check if category has active listings
	var activeListingsCount int
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM listings
		WHERE category_id = $1 AND status = 'active' AND is_deleted = false
	`, categoryID).Scan(&activeListingsCount)
	if err != nil {
		r.logger.Error().Err(err).Int64("category_id", categoryID).Msg("failed to check active listings")
		return fmt.Errorf("failed to check active listings: %w", err)
	}

	if activeListingsCount > 0 {
		return fmt.Errorf("cannot delete category with %d active listings", activeListingsCount)
	}

	// Soft delete the category and all its children
	query := `
		WITH RECURSIVE category_tree AS (
			SELECT id FROM categories WHERE id = $1
			UNION ALL
			SELECT c.id FROM categories c
			INNER JOIN category_tree ct ON c.parent_id = ct.id
		)
		UPDATE categories SET is_active = false
		WHERE id IN (SELECT id FROM category_tree)
	`

	result, err := r.db.ExecContext(ctx, query, categoryID)
	if err != nil {
		r.logger.Error().Err(err).Int64("category_id", categoryID).Msg("failed to delete category")
		return fmt.Errorf("failed to delete category: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("category not found")
	}

	r.logger.Info().
		Int64("category_id", categoryID).
		Int64("categories_deactivated", rowsAffected).
		Msg("category and descendants deleted")

	return nil
}

// isDescendant checks if potentialDescendantID is a descendant of ancestorID
func (r *Repository) isDescendant(ctx context.Context, ancestorID, potentialDescendantID int64) (bool, error) {
	query := `
		WITH RECURSIVE category_tree AS (
			SELECT id, parent_id FROM categories WHERE id = $1
			UNION ALL
			SELECT c.id, c.parent_id FROM categories c
			INNER JOIN category_tree ct ON c.parent_id = ct.id
		)
		SELECT EXISTS(SELECT 1 FROM category_tree WHERE id = $2)
	`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, ancestorID, potentialDescendantID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check descendant: %w", err)
	}

	return exists, nil
}

// updateDescendantLevels updates the level of all descendants when a category's level changes
func (r *Repository) updateDescendantLevels(ctx context.Context, categoryID int64, newLevel int32) error {
	query := `
		WITH RECURSIVE category_tree AS (
			SELECT id, parent_id, $2::int + 1 as new_level FROM categories WHERE parent_id = $1
			UNION ALL
			SELECT c.id, c.parent_id, ct.new_level + 1
			FROM categories c
			INNER JOIN category_tree ct ON c.parent_id = ct.id
		)
		UPDATE categories c SET level = ct.new_level
		FROM category_tree ct
		WHERE c.id = ct.id
	`

	_, err := r.db.ExecContext(ctx, query, categoryID, newLevel)
	if err != nil {
		r.logger.Error().Err(err).Int64("category_id", categoryID).Msg("failed to update descendant levels")
		return fmt.Errorf("failed to update descendant levels: %w", err)
	}

	return nil
}

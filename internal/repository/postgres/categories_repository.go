package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sveturs/listings/internal/domain"
)

// GetRootCategories retrieves all top-level categories (no parent)
func (r *Repository) GetRootCategories(ctx context.Context) ([]*domain.Category, error) {
	query := `
		SELECT
			id, name, slug, parent_id, icon, description,
			is_active, created_at, sort_order, level,
			has_custom_ui, custom_ui_component,
			COALESCE((SELECT COUNT(*) FROM c2c_listings WHERE category_id = mc.id AND status = 'active'), 0) as listing_count
		FROM marketplace_categories mc
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
			COALESCE((SELECT COUNT(*) FROM c2c_listings WHERE category_id = mc.id AND status = 'active'), 0) as listing_count
		FROM marketplace_categories mc
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
		FROM marketplace_categories mc
		LEFT JOIN c2c_listings l ON l.category_id = mc.id AND l.status = 'active'
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
			COALESCE((SELECT COUNT(*) FROM c2c_listings WHERE category_id = mc.id AND status = 'active'), 0) as listing_count
		FROM marketplace_categories mc
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

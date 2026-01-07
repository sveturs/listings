package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/vondi-global/listings/internal/domain"
)

// GetByUUID retrieves a category by UUID
func (r *Repository) GetByUUID(ctx context.Context, id string) (*domain.CategoryV2, error) {
	query := `
		SELECT
			id, slug, parent_id, level, path, sort_order,
			name, description, meta_title, meta_description, meta_keywords,
			icon, image_url, is_active, created_at, updated_at
		FROM categories
		WHERE id = $1
	`

	cat := &domain.CategoryV2{
		Name:            make(map[string]string),
		Description:     make(map[string]string),
		MetaTitle:       make(map[string]string),
		MetaDescription: make(map[string]string),
		MetaKeywords:    make(map[string]string),
	}

	var (
		nameJSON, descJSON, metaTitleJSON, metaDescJSON, metaKeywordsJSON []byte
		parentIDStr                                                       sql.NullString
	)

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&cat.ID,
		&cat.Slug,
		&parentIDStr,
		&cat.Level,
		&cat.Path,
		&cat.SortOrder,
		&nameJSON,
		&descJSON,
		&metaTitleJSON,
		&metaDescJSON,
		&metaKeywordsJSON,
		&cat.Icon,
		&cat.ImageURL,
		&cat.IsActive,
		&cat.CreatedAt,
		&cat.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("category not found: %s", id)
	}
	if err != nil {
		r.logger.Error().Err(err).Str("id", id).Msg("failed to get category by UUID")
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	// Parse parent_id
	if parentIDStr.Valid {
		parentUUID, err := uuid.Parse(parentIDStr.String)
		if err == nil {
			cat.ParentID = &parentUUID
		}
	}

	// Parse JSONB fields
	if err := json.Unmarshal(nameJSON, &cat.Name); err != nil {
		r.logger.Warn().Err(err).Msg("failed to unmarshal name JSONB")
	}
	if len(descJSON) > 0 {
		if err := json.Unmarshal(descJSON, &cat.Description); err != nil {
			r.logger.Warn().Err(err).Msg("failed to unmarshal description JSONB")
		}
	}
	if len(metaTitleJSON) > 0 {
		if err := json.Unmarshal(metaTitleJSON, &cat.MetaTitle); err != nil {
			r.logger.Warn().Err(err).Msg("failed to unmarshal meta_title JSONB")
		}
	}
	if len(metaDescJSON) > 0 {
		if err := json.Unmarshal(metaDescJSON, &cat.MetaDescription); err != nil {
			r.logger.Warn().Err(err).Msg("failed to unmarshal meta_description JSONB")
		}
	}
	if len(metaKeywordsJSON) > 0 {
		if err := json.Unmarshal(metaKeywordsJSON, &cat.MetaKeywords); err != nil {
			r.logger.Warn().Err(err).Msg("failed to unmarshal meta_keywords JSONB")
		}
	}

	return cat, nil
}

// GetBySlugV2 retrieves a category by slug
func (r *Repository) GetBySlugV2(ctx context.Context, slug string) (*domain.CategoryV2, error) {
	query := `
		SELECT
			id, slug, parent_id, level, path, sort_order,
			name, description, meta_title, meta_description, meta_keywords,
			icon, image_url, is_active, created_at, updated_at
		FROM categories
		WHERE slug = $1
	`

	cat := &domain.CategoryV2{
		Name:            make(map[string]string),
		Description:     make(map[string]string),
		MetaTitle:       make(map[string]string),
		MetaDescription: make(map[string]string),
		MetaKeywords:    make(map[string]string),
	}

	var (
		nameJSON, descJSON, metaTitleJSON, metaDescJSON, metaKeywordsJSON []byte
		parentIDStr                                                       sql.NullString
	)

	err := r.db.QueryRowContext(ctx, query, slug).Scan(
		&cat.ID,
		&cat.Slug,
		&parentIDStr,
		&cat.Level,
		&cat.Path,
		&cat.SortOrder,
		&nameJSON,
		&descJSON,
		&metaTitleJSON,
		&metaDescJSON,
		&metaKeywordsJSON,
		&cat.Icon,
		&cat.ImageURL,
		&cat.IsActive,
		&cat.CreatedAt,
		&cat.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("category not found: %s", slug)
	}
	if err != nil {
		r.logger.Error().Err(err).Str("slug", slug).Msg("failed to get category by slug")
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	// Parse parent_id
	if parentIDStr.Valid {
		parentUUID, err := uuid.Parse(parentIDStr.String)
		if err == nil {
			cat.ParentID = &parentUUID
		}
	}

	// Parse JSONB fields
	if err := json.Unmarshal(nameJSON, &cat.Name); err != nil {
		r.logger.Warn().Err(err).Msg("failed to unmarshal name JSONB")
	}
	if len(descJSON) > 0 {
		if err := json.Unmarshal(descJSON, &cat.Description); err != nil {
			r.logger.Warn().Err(err).Msg("failed to unmarshal description JSONB")
		}
	}
	if len(metaTitleJSON) > 0 {
		if err := json.Unmarshal(metaTitleJSON, &cat.MetaTitle); err != nil {
			r.logger.Warn().Err(err).Msg("failed to unmarshal meta_title JSONB")
		}
	}
	if len(metaDescJSON) > 0 {
		if err := json.Unmarshal(metaDescJSON, &cat.MetaDescription); err != nil {
			r.logger.Warn().Err(err).Msg("failed to unmarshal meta_description JSONB")
		}
	}
	if len(metaKeywordsJSON) > 0 {
		if err := json.Unmarshal(metaKeywordsJSON, &cat.MetaKeywords); err != nil {
			r.logger.Warn().Err(err).Msg("failed to unmarshal meta_keywords JSONB")
		}
	}

	return cat, nil
}

// GetTreeV2 retrieves category tree with localization
func (r *Repository) GetTreeV2(ctx context.Context, filter *domain.GetCategoryTreeFilterV2) ([]*domain.CategoryTreeV2, error) {
	// Build WHERE clause
	var whereClauses []string
	var args []interface{}
	argPos := 1

	if filter.RootID != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("parent_id = $%d", argPos))
		args = append(args, filter.RootID.String())
	} else {
		whereClauses = append(whereClauses, "parent_id IS NULL")
	}

	if filter.ActiveOnly {
		whereClauses = append(whereClauses, "is_active = true")
	}

	whereClause := "WHERE " + strings.Join(whereClauses, " AND ")

	query := fmt.Sprintf(`
		SELECT
			id, slug, parent_id, level, path, sort_order,
			name, description, meta_title, meta_description, meta_keywords,
			icon, image_url, is_active, created_at, updated_at
		FROM categories
		%s
		ORDER BY sort_order ASC, slug ASC
	`, whereClause)

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to query category tree")
		return nil, fmt.Errorf("failed to query category tree: %w", err)
	}
	defer rows.Close()

	categories, err := scanCategoriesV2(rows, r)
	if err != nil {
		return nil, err
	}

	// Build tree structure with localization
	tree := make([]*domain.CategoryTreeV2, 0, len(categories))
	for _, cat := range categories {
		localized := cat.Localize(filter.Locale)
		node := &domain.CategoryTreeV2{
			Category:      localized,
			Subcategories: []*domain.CategoryTreeV2{},
		}

		// Recursively fetch children if needed
		if filter.MaxDepth == nil || cat.Level < *filter.MaxDepth {
			childFilter := &domain.GetCategoryTreeFilterV2{
				RootID:     &cat.ID,
				Locale:     filter.Locale,
				ActiveOnly: filter.ActiveOnly,
				MaxDepth:   filter.MaxDepth,
			}
			children, err := r.GetTreeV2(ctx, childFilter)
			if err != nil {
				r.logger.Warn().Err(err).Str("parent_id", cat.ID.String()).Msg("failed to fetch children")
			} else {
				node.Subcategories = children
			}
		}

		tree = append(tree, node)
	}

	return tree, nil
}

// GetBreadcrumb retrieves breadcrumb trail for a category
func (r *Repository) GetBreadcrumb(ctx context.Context, categoryID string, locale string) ([]*domain.CategoryBreadcrumb, error) {
	query := `
		WITH RECURSIVE category_path AS (
			-- Base case: the target category
			SELECT id, slug, parent_id, level, name, 1 as depth
			FROM categories
			WHERE id = $1

			UNION ALL

			-- Recursive case: parent categories
			SELECT c.id, c.slug, c.parent_id, c.level, c.name, cp.depth + 1
			FROM categories c
			INNER JOIN category_path cp ON c.id = cp.parent_id
		)
		SELECT id, slug, level, name
		FROM category_path
		ORDER BY level ASC
	`

	rows, err := r.db.QueryContext(ctx, query, categoryID)
	if err != nil {
		r.logger.Error().Err(err).Str("category_id", categoryID).Msg("failed to get breadcrumb")
		return nil, fmt.Errorf("failed to get breadcrumb: %w", err)
	}
	defer rows.Close()

	breadcrumbs := make([]*domain.CategoryBreadcrumb, 0)
	for rows.Next() {
		var (
			id       uuid.UUID
			slug     string
			level    int32
			nameJSON []byte
		)

		if err := rows.Scan(&id, &slug, &level, &nameJSON); err != nil {
			r.logger.Error().Err(err).Msg("failed to scan breadcrumb row")
			return nil, fmt.Errorf("failed to scan breadcrumb: %w", err)
		}

		// Parse JSONB name
		nameMap := make(map[string]string)
		if err := json.Unmarshal(nameJSON, &nameMap); err != nil {
			r.logger.Warn().Err(err).Msg("failed to unmarshal name JSONB")
			nameMap = map[string]string{"sr": slug} // Fallback to slug
		}

		// Extract localized name using same logic as domain.getLocalized
		localizedName := getLocalizedFromMap(nameMap, locale)

		breadcrumbs = append(breadcrumbs, &domain.CategoryBreadcrumb{
			ID:    id,
			Slug:  slug,
			Name:  localizedName,
			Level: level,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("breadcrumb rows iteration error: %w", err)
	}

	return breadcrumbs, nil
}

// ListV2 retrieves categories with pagination
func (r *Repository) ListV2(ctx context.Context, parentID *string, activeOnly bool, page, pageSize int32) ([]*domain.CategoryV2, int64, error) {
	// Build WHERE clause
	var whereClauses []string
	var args []interface{}
	argPos := 1

	if parentID != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("parent_id = $%d", argPos))
		args = append(args, *parentID)
		argPos++
	} else {
		whereClauses = append(whereClauses, "parent_id IS NULL")
	}

	if activeOnly {
		whereClauses = append(whereClauses, "is_active = true")
	}

	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM categories %s", whereClause)
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		r.logger.Error().Err(err).Msg("failed to count categories")
		return nil, 0, fmt.Errorf("failed to count categories: %w", err)
	}

	// Pagination
	offset := (page - 1) * pageSize
	args = append(args, pageSize, offset)

	query := fmt.Sprintf(`
		SELECT
			id, slug, parent_id, level, path, sort_order,
			name, description, meta_title, meta_description, meta_keywords,
			icon, image_url, is_active, created_at, updated_at
		FROM categories
		%s
		ORDER BY sort_order ASC, slug ASC
		LIMIT $%d OFFSET $%d
	`, whereClause, argPos, argPos+1)

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to list categories")
		return nil, 0, fmt.Errorf("failed to list categories: %w", err)
	}
	defer rows.Close()

	categories, err := scanCategoriesV2(rows, r)
	if err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

// scanCategoriesV2 is a helper to scan multiple CategoryV2 rows
func scanCategoriesV2(rows *sqlx.Rows, r *Repository) ([]*domain.CategoryV2, error) {
	categories := make([]*domain.CategoryV2, 0)

	for rows.Next() {
		cat := &domain.CategoryV2{
			Name:            make(map[string]string),
			Description:     make(map[string]string),
			MetaTitle:       make(map[string]string),
			MetaDescription: make(map[string]string),
			MetaKeywords:    make(map[string]string),
		}

		var (
			nameJSON, descJSON, metaTitleJSON, metaDescJSON, metaKeywordsJSON []byte
			parentIDStr                                                       sql.NullString
		)

		err := rows.Scan(
			&cat.ID,
			&cat.Slug,
			&parentIDStr,
			&cat.Level,
			&cat.Path,
			&cat.SortOrder,
			&nameJSON,
			&descJSON,
			&metaTitleJSON,
			&metaDescJSON,
			&metaKeywordsJSON,
			&cat.Icon,
			&cat.ImageURL,
			&cat.IsActive,
			&cat.CreatedAt,
			&cat.UpdatedAt,
		)

		if err != nil {
			r.logger.Error().Err(err).Msg("failed to scan category")
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}

		// Parse parent_id
		if parentIDStr.Valid {
			parentUUID, err := uuid.Parse(parentIDStr.String)
			if err == nil {
				cat.ParentID = &parentUUID
			}
		}

		// Parse JSONB fields
		if err := json.Unmarshal(nameJSON, &cat.Name); err != nil {
			r.logger.Warn().Err(err).Msg("failed to unmarshal name JSONB")
		}
		if len(descJSON) > 0 {
			json.Unmarshal(descJSON, &cat.Description)
		}
		if len(metaTitleJSON) > 0 {
			json.Unmarshal(metaTitleJSON, &cat.MetaTitle)
		}
		if len(metaDescJSON) > 0 {
			json.Unmarshal(metaDescJSON, &cat.MetaDescription)
		}
		if len(metaKeywordsJSON) > 0 {
			json.Unmarshal(metaKeywordsJSON, &cat.MetaKeywords)
		}

		categories = append(categories, cat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return categories, nil
}

// GetAllActiveV2 retrieves all active categories for category detection
func (r *Repository) GetAllActiveV2(ctx context.Context) ([]*domain.CategoryV2, error) {
	query := `
		SELECT
			id, slug, parent_id, level, path, sort_order,
			name, description, meta_title, meta_description, meta_keywords,
			icon, image_url, is_active, created_at, updated_at
		FROM categories
		WHERE is_active = true
		ORDER BY level ASC, sort_order ASC, slug ASC
	`

	rows, err := r.db.QueryxContext(ctx, query)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to query all active categories")
		return nil, fmt.Errorf("failed to query all active categories: %w", err)
	}
	defer rows.Close()

	return scanCategoriesV2(rows, r)
}

// getLocalizedFromMap extracts localized value with fallback logic
func getLocalizedFromMap(m map[string]string, locale string) string {
	if m == nil {
		return ""
	}

	// Try requested locale
	if val, ok := m[locale]; ok && val != "" {
		return val
	}

	// Fallback to Serbian
	if locale != "sr" {
		if val, ok := m["sr"]; ok && val != "" {
			return val
		}
	}

	// Fallback to English
	if locale != "en" {
		if val, ok := m["en"]; ok && val != "" {
			return val
		}
	}

	// Return first available
	for _, val := range m {
		if val != "" {
			return val
		}
	}

	return ""
}

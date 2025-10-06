package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"backend/internal/domain/models"
)

// CategoryMappingFilter represents filter options for listing category mappings
type CategoryMappingFilter struct {
	IsManual *bool
	Limit    int
	Offset   int
}

// CategoryMappingsRepositoryInterface defines the interface for category mappings repository
type CategoryMappingsRepositoryInterface interface {
	// Create creates a new category mapping (with UPSERT on conflict)
	Create(ctx context.Context, mapping *models.StorefrontCategoryMapping) error

	// GetByID retrieves a mapping by ID
	GetByID(ctx context.Context, id int) (*models.StorefrontCategoryMapping, error)

	// GetBySourcePath retrieves mapping by storefront ID and source category path
	GetBySourcePath(ctx context.Context, storefrontID int, sourcePath string) (*models.StorefrontCategoryMapping, error)

	// GetByStorefront retrieves all mappings for a storefront with optional filter
	GetByStorefront(ctx context.Context, storefrontID int, filter *CategoryMappingFilter) ([]*models.StorefrontCategoryMappingWithDetails, int, error)

	// Update updates an existing mapping
	Update(ctx context.Context, mapping *models.StorefrontCategoryMapping) error

	// Delete deletes a mapping by ID
	Delete(ctx context.Context, id int) error

	// DeleteBySourcePath deletes a mapping by storefront ID and source path
	DeleteBySourcePath(ctx context.Context, storefrontID int, sourcePath string) error
}

// categoryMappingsRepository implements the CategoryMappingsRepositoryInterface
type categoryMappingsRepository struct {
	pool *pgxpool.Pool
}

// NewCategoryMappingsRepository creates a new category mappings repository
func NewCategoryMappingsRepository(pool *pgxpool.Pool) CategoryMappingsRepositoryInterface {
	return &categoryMappingsRepository{pool: pool}
}

// Create creates a new category mapping with UPSERT semantics
func (r *categoryMappingsRepository) Create(ctx context.Context, mapping *models.StorefrontCategoryMapping) error {
	query := `
		INSERT INTO storefront_category_mappings (
			storefront_id, source_category_path, target_category_id,
			is_manual, confidence_score
		) VALUES (
			$1, $2, $3, $4, $5
		)
		ON CONFLICT (storefront_id, source_category_path)
		DO UPDATE SET
			target_category_id = EXCLUDED.target_category_id,
			is_manual = EXCLUDED.is_manual,
			confidence_score = EXCLUDED.confidence_score,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id, created_at, updated_at`

	err := r.pool.QueryRow(
		ctx, query,
		mapping.StorefrontID,
		mapping.SourceCategoryPath,
		mapping.TargetCategoryID,
		mapping.IsManual,
		mapping.ConfidenceScore,
	).Scan(&mapping.ID, &mapping.CreatedAt, &mapping.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create category mapping: %w", err)
	}

	return nil
}

// GetByID retrieves a mapping by ID
func (r *categoryMappingsRepository) GetByID(ctx context.Context, id int) (*models.StorefrontCategoryMapping, error) {
	query := `
		SELECT id, storefront_id, source_category_path, target_category_id,
		       is_manual, confidence_score, created_at, updated_at
		FROM storefront_category_mappings
		WHERE id = $1`

	var mapping models.StorefrontCategoryMapping
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&mapping.ID,
		&mapping.StorefrontID,
		&mapping.SourceCategoryPath,
		&mapping.TargetCategoryID,
		&mapping.IsManual,
		&mapping.ConfidenceScore,
		&mapping.CreatedAt,
		&mapping.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get category mapping: %w", err)
	}

	return &mapping, nil
}

// GetBySourcePath retrieves mapping by storefront ID and source category path
func (r *categoryMappingsRepository) GetBySourcePath(ctx context.Context, storefrontID int, sourcePath string) (*models.StorefrontCategoryMapping, error) {
	query := `
		SELECT id, storefront_id, source_category_path, target_category_id,
		       is_manual, confidence_score, created_at, updated_at
		FROM storefront_category_mappings
		WHERE storefront_id = $1 AND source_category_path = $2`

	var mapping models.StorefrontCategoryMapping
	err := r.pool.QueryRow(ctx, query, storefrontID, sourcePath).Scan(
		&mapping.ID,
		&mapping.StorefrontID,
		&mapping.SourceCategoryPath,
		&mapping.TargetCategoryID,
		&mapping.IsManual,
		&mapping.ConfidenceScore,
		&mapping.CreatedAt,
		&mapping.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get category mapping by source path: %w", err)
	}

	return &mapping, nil
}

// GetByStorefront retrieves all mappings for a storefront with optional filter
func (r *categoryMappingsRepository) GetByStorefront(ctx context.Context, storefrontID int, filter *CategoryMappingFilter) ([]*models.StorefrontCategoryMappingWithDetails, int, error) {
	// Build WHERE clause
	whereClause := "WHERE scm.storefront_id = $1"
	args := []interface{}{storefrontID}
	argIdx := 2

	if filter != nil && filter.IsManual != nil {
		whereClause += fmt.Sprintf(" AND scm.is_manual = $%d", argIdx)
		args = append(args, *filter.IsManual)
		argIdx++
	}

	// Count total
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM storefront_category_mappings scm
		%s`, whereClause)

	var total int
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count category mappings: %w", err)
	}

	// Build query with category details
	// TODO: Add recursive category path building in future version
	query := fmt.Sprintf(`
		SELECT
			scm.id, scm.storefront_id, scm.source_category_path,
			scm.target_category_id, scm.is_manual, scm.confidence_score,
			scm.created_at, scm.updated_at,
			mc.name as target_category_name,
			mc.slug as target_category_path
		FROM storefront_category_mappings scm
		JOIN marketplace_categories mc ON mc.id = scm.target_category_id
		%s
		ORDER BY scm.created_at DESC`, whereClause)

	// Add pagination
	if filter != nil {
		if filter.Limit > 0 {
			query += fmt.Sprintf(" LIMIT $%d", argIdx)
			args = append(args, filter.Limit)
			argIdx++
		}
		if filter.Offset > 0 {
			query += fmt.Sprintf(" OFFSET $%d", argIdx)
			args = append(args, filter.Offset)
		}
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get category mappings: %w", err)
	}
	defer rows.Close()

	var mappings []*models.StorefrontCategoryMappingWithDetails
	for rows.Next() {
		var m models.StorefrontCategoryMappingWithDetails
		err := rows.Scan(
			&m.ID,
			&m.StorefrontID,
			&m.SourceCategoryPath,
			&m.TargetCategoryID,
			&m.IsManual,
			&m.ConfidenceScore,
			&m.CreatedAt,
			&m.UpdatedAt,
			&m.TargetCategoryName,
			&m.TargetCategoryPath,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan category mapping: %w", err)
		}
		mappings = append(mappings, &m)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating category mappings: %w", err)
	}

	return mappings, total, nil
}

// Update updates an existing mapping
func (r *categoryMappingsRepository) Update(ctx context.Context, mapping *models.StorefrontCategoryMapping) error {
	query := `
		UPDATE storefront_category_mappings
		SET target_category_id = $1,
		    is_manual = $2,
		    confidence_score = $3,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
		RETURNING updated_at`

	err := r.pool.QueryRow(
		ctx, query,
		mapping.TargetCategoryID,
		mapping.IsManual,
		mapping.ConfidenceScore,
		mapping.ID,
	).Scan(&mapping.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return sql.ErrNoRows
		}
		return fmt.Errorf("failed to update category mapping: %w", err)
	}

	return nil
}

// Delete deletes a mapping by ID
func (r *categoryMappingsRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM storefront_category_mappings WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete category mapping: %w", err)
	}

	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// DeleteBySourcePath deletes a mapping by storefront ID and source path
func (r *categoryMappingsRepository) DeleteBySourcePath(ctx context.Context, storefrontID int, sourcePath string) error {
	query := `DELETE FROM storefront_category_mappings WHERE storefront_id = $1 AND source_category_path = $2`

	result, err := r.pool.Exec(ctx, query, storefrontID, sourcePath)
	if err != nil {
		return fmt.Errorf("failed to delete category mapping by source path: %w", err)
	}

	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

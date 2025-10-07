package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"backend/internal/domain/models"
)

// CategoryProposalFilter represents filter options for listing category proposals
type CategoryProposalFilter struct {
	Status       *models.CategoryProposalStatus
	StorefrontID *int
	Limit        int
	Offset       int
}

// CategoryProposalsRepositoryInterface defines the interface for category proposals repository
type CategoryProposalsRepositoryInterface interface {
	// Create creates a new category proposal
	Create(ctx context.Context, proposal *models.CategoryProposal) error

	// GetByID retrieves a proposal by ID
	GetByID(ctx context.Context, id int) (*models.CategoryProposal, error)

	// List retrieves proposals with filters and pagination
	List(ctx context.Context, filter *CategoryProposalFilter) ([]*models.CategoryProposal, int, error)

	// Update updates an existing proposal
	Update(ctx context.Context, proposal *models.CategoryProposal) error

	// Approve approves a proposal
	Approve(ctx context.Context, id int, reviewedByUserID int) error

	// Reject rejects a proposal
	Reject(ctx context.Context, id int, reviewedByUserID int, reason *string) error

	// Delete deletes a proposal by ID
	Delete(ctx context.Context, id int) error

	// GetPendingCount returns count of pending proposals
	GetPendingCount(ctx context.Context, storefrontID *int) (int, error)
}

// categoryProposalsRepository implements the CategoryProposalsRepositoryInterface
type categoryProposalsRepository struct {
	pool *pgxpool.Pool
}

// NewCategoryProposalsRepository creates a new category proposals repository
func NewCategoryProposalsRepository(pool *pgxpool.Pool) CategoryProposalsRepositoryInterface {
	return &categoryProposalsRepository{pool: pool}
}

// Create creates a new category proposal
func (r *categoryProposalsRepository) Create(ctx context.Context, proposal *models.CategoryProposal) error {
	query := `
		INSERT INTO category_proposals (
			proposed_by_user_id, storefront_id, name, name_translations,
			parent_category_id, description, reasoning, expected_products,
			external_category_source, similar_categories, tags, status
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
		)
		RETURNING id, created_at, updated_at`

	err := r.pool.QueryRow(
		ctx, query,
		proposal.ProposedByUserID,
		proposal.StorefrontID,
		proposal.Name,
		proposal.NameTranslations,
		proposal.ParentCategoryID,
		proposal.Description,
		proposal.Reasoning,
		proposal.ExpectedProducts,
		proposal.ExternalCategorySource,
		proposal.SimilarCategories,
		proposal.Tags,
		proposal.Status,
	).Scan(&proposal.ID, &proposal.CreatedAt, &proposal.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create category proposal: %w", err)
	}

	return nil
}

// GetByID retrieves a proposal by ID
func (r *categoryProposalsRepository) GetByID(ctx context.Context, id int) (*models.CategoryProposal, error) {
	query := `
		SELECT
			id, proposed_by_user_id, storefront_id, name, name_translations,
			parent_category_id, description, reasoning, expected_products,
			external_category_source, similar_categories, tags, status,
			reviewed_by_user_id, reviewed_at, created_at, updated_at
		FROM category_proposals
		WHERE id = $1`

	var proposal models.CategoryProposal
	var similarCategories []int
	var tags []string

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&proposal.ID,
		&proposal.ProposedByUserID,
		&proposal.StorefrontID,
		&proposal.Name,
		&proposal.NameTranslations,
		&proposal.ParentCategoryID,
		&proposal.Description,
		&proposal.Reasoning,
		&proposal.ExpectedProducts,
		&proposal.ExternalCategorySource,
		&similarCategories,
		&tags,
		&proposal.Status,
		&proposal.ReviewedByUserID,
		&proposal.ReviewedAt,
		&proposal.CreatedAt,
		&proposal.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("category proposal not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get category proposal: %w", err)
	}

	proposal.SimilarCategories = similarCategories
	proposal.Tags = tags

	return &proposal, nil
}

// List retrieves proposals with filters and pagination
func (r *categoryProposalsRepository) List(ctx context.Context, filter *CategoryProposalFilter) ([]*models.CategoryProposal, int, error) {
	// Build WHERE clause
	whereClauses := []string{}
	args := []interface{}{}
	argPos := 1

	if filter.Status != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("status = $%d", argPos))
		args = append(args, *filter.Status)
		argPos++
	}

	if filter.StorefrontID != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("storefront_id = $%d", argPos))
		args = append(args, *filter.StorefrontID)
		argPos++
	}

	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Count query
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM category_proposals
		%s`, whereClause)

	var total int
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count category proposals: %w", err)
	}

	// List query
	listQuery := fmt.Sprintf(`
		SELECT
			id, proposed_by_user_id, storefront_id, name, name_translations,
			parent_category_id, description, reasoning, expected_products,
			external_category_source, similar_categories, tags, status,
			reviewed_by_user_id, reviewed_at, created_at, updated_at
		FROM category_proposals
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, argPos, argPos+1)

	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.pool.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list category proposals: %w", err)
	}
	defer rows.Close()

	proposals := []*models.CategoryProposal{}
	for rows.Next() {
		var proposal models.CategoryProposal
		var similarCategories []int
		var tags []string

		err := rows.Scan(
			&proposal.ID,
			&proposal.ProposedByUserID,
			&proposal.StorefrontID,
			&proposal.Name,
			&proposal.NameTranslations,
			&proposal.ParentCategoryID,
			&proposal.Description,
			&proposal.Reasoning,
			&proposal.ExpectedProducts,
			&proposal.ExternalCategorySource,
			&similarCategories,
			&tags,
			&proposal.Status,
			&proposal.ReviewedByUserID,
			&proposal.ReviewedAt,
			&proposal.CreatedAt,
			&proposal.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan category proposal: %w", err)
		}

		proposal.SimilarCategories = similarCategories
		proposal.Tags = tags

		proposals = append(proposals, &proposal)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating category proposals: %w", err)
	}

	return proposals, total, nil
}

// Update updates an existing proposal
func (r *categoryProposalsRepository) Update(ctx context.Context, proposal *models.CategoryProposal) error {
	query := `
		UPDATE category_proposals
		SET name = $1, name_translations = $2, parent_category_id = $3,
		    description = $4, tags = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $6
		RETURNING updated_at`

	err := r.pool.QueryRow(
		ctx, query,
		proposal.Name,
		proposal.NameTranslations,
		proposal.ParentCategoryID,
		proposal.Description,
		proposal.Tags,
		proposal.ID,
	).Scan(&proposal.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("category proposal not found: %w", err)
		}
		return fmt.Errorf("failed to update category proposal: %w", err)
	}

	return nil
}

// Approve approves a proposal
func (r *categoryProposalsRepository) Approve(ctx context.Context, id int, reviewedByUserID int) error {
	query := `
		UPDATE category_proposals
		SET status = $1, reviewed_by_user_id = $2, reviewed_at = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4 AND status = 'pending'`

	result, err := r.pool.Exec(
		ctx, query,
		models.CategoryProposalStatusApproved,
		reviewedByUserID,
		time.Now(),
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to approve category proposal: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("category proposal not found or not pending")
	}

	return nil
}

// Reject rejects a proposal
func (r *categoryProposalsRepository) Reject(ctx context.Context, id int, reviewedByUserID int, reason *string) error {
	query := `
		UPDATE category_proposals
		SET status = $1, reviewed_by_user_id = $2, reviewed_at = $3,
		    description = COALESCE($4, description), updated_at = CURRENT_TIMESTAMP
		WHERE id = $5 AND status = 'pending'`

	result, err := r.pool.Exec(
		ctx, query,
		models.CategoryProposalStatusRejected,
		reviewedByUserID,
		time.Now(),
		reason,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to reject category proposal: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("category proposal not found or not pending")
	}

	return nil
}

// Delete deletes a proposal by ID
func (r *categoryProposalsRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM category_proposals WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete category proposal: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("category proposal not found")
	}

	return nil
}

// GetPendingCount returns count of pending proposals
func (r *categoryProposalsRepository) GetPendingCount(ctx context.Context, storefrontID *int) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM category_proposals
		WHERE status = 'pending'`

	args := []interface{}{}
	if storefrontID != nil {
		query += " AND storefront_id = $1"
		args = append(args, *storefrontID)
	}

	var count int
	err := r.pool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count pending proposals: %w", err)
	}

	return count, nil
}

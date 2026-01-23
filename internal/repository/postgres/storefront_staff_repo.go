package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/vondi-global/listings/internal/domain"
)

// StorefrontStaffRepository handles storefront staff database operations
type StorefrontStaffRepository struct {
	db *sql.DB
}

// NewStorefrontStaffRepository creates a new staff repository
func NewStorefrontStaffRepository(db *sql.DB) *StorefrontStaffRepository {
	return &StorefrontStaffRepository{db: db}
}

// Create adds a new staff member to a storefront
func (r *StorefrontStaffRepository) Create(ctx context.Context, staff *domain.StorefrontStaff) error {
	query := `
		INSERT INTO storefront_staff (
			storefront_id, user_id, role, permissions, actions_count
		) VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRowContext(ctx, query,
		staff.StorefrontID,
		staff.UserID,
		staff.Role,
		staff.Permissions,
		staff.ActionsCount,
	).Scan(&staff.ID, &staff.CreatedAt, &staff.UpdatedAt)
}

// GetByID retrieves a staff member by ID
func (r *StorefrontStaffRepository) GetByID(ctx context.Context, id int64) (*domain.StorefrontStaff, error) {
	query := `
		SELECT id, storefront_id, user_id, role, permissions,
		       last_active_at, actions_count, created_at, updated_at
		FROM storefront_staff
		WHERE id = $1`

	staff := &domain.StorefrontStaff{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&staff.ID,
		&staff.StorefrontID,
		&staff.UserID,
		&staff.Role,
		&staff.Permissions,
		&staff.LastActiveAt,
		&staff.ActionsCount,
		&staff.CreatedAt,
		&staff.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("staff member not found")
		}
		return nil, fmt.Errorf("failed to get staff member: %w", err)
	}

	return staff, nil
}

// GetByStorefrontAndUser retrieves a staff member by storefront and user
func (r *StorefrontStaffRepository) GetByStorefrontAndUser(ctx context.Context, storefrontID int64, userID int64) (*domain.StorefrontStaff, error) {
	query := `
		SELECT id, storefront_id, user_id, role, permissions,
		       last_active_at, actions_count, created_at, updated_at
		FROM storefront_staff
		WHERE storefront_id = $1 AND user_id = $2`

	staff := &domain.StorefrontStaff{}
	err := r.db.QueryRowContext(ctx, query, storefrontID, userID).Scan(
		&staff.ID,
		&staff.StorefrontID,
		&staff.UserID,
		&staff.Role,
		&staff.Permissions,
		&staff.LastActiveAt,
		&staff.ActionsCount,
		&staff.CreatedAt,
		&staff.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found - not an error
		}
		return nil, fmt.Errorf("failed to get staff member: %w", err)
	}

	return staff, nil
}

// ListByStorefront lists all staff members for a storefront
func (r *StorefrontStaffRepository) ListByStorefront(ctx context.Context, storefrontID int64) ([]*domain.StorefrontStaff, error) {
	query := `
		SELECT id, storefront_id, user_id, role, permissions,
		       last_active_at, actions_count, created_at, updated_at
		FROM storefront_staff
		WHERE storefront_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, storefrontID)
	if err != nil {
		return nil, fmt.Errorf("failed to list staff members: %w", err)
	}
	defer rows.Close()

	var staff []*domain.StorefrontStaff
	for rows.Next() {
		s := &domain.StorefrontStaff{}
		if err := rows.Scan(
			&s.ID,
			&s.StorefrontID,
			&s.UserID,
			&s.Role,
			&s.Permissions,
			&s.LastActiveAt,
			&s.ActionsCount,
			&s.CreatedAt,
			&s.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan staff member: %w", err)
		}
		staff = append(staff, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return staff, nil
}

// Delete removes a staff member
func (r *StorefrontStaffRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM storefront_staff WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete staff member: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("staff member not found")
	}

	return nil
}

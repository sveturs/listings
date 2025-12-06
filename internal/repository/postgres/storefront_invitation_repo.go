package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/vondi-global/listings/internal/domain"
)

// StorefrontInvitationRepository handles storefront invitation database operations
type StorefrontInvitationRepository struct {
	db *sql.DB
}

// NewStorefrontInvitationRepository creates a new invitation repository
func NewStorefrontInvitationRepository(db *sql.DB) *StorefrontInvitationRepository {
	return &StorefrontInvitationRepository{db: db}
}

// Create creates a new invitation
func (r *StorefrontInvitationRepository) Create(ctx context.Context, inv *domain.StorefrontInvitation) error {
	query := `
		INSERT INTO storefront_invitations (
			storefront_id, role, type,
			invited_email, invited_user_id,
			invite_code, expires_at, max_uses, current_uses,
			invited_by_id, status, comment
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRowContext(ctx, query,
		inv.StorefrontID,
		inv.Role,
		inv.Type,
		inv.InvitedEmail,
		inv.InvitedUserID,
		inv.InviteCode,
		inv.ExpiresAt,
		inv.MaxUses,
		inv.CurrentUses,
		inv.InvitedByID,
		inv.Status,
		inv.Comment,
	).Scan(&inv.ID, &inv.CreatedAt, &inv.UpdatedAt)
}

// GetByID returns invitation by ID
func (r *StorefrontInvitationRepository) GetByID(ctx context.Context, id int64) (*domain.StorefrontInvitation, error) {
	query := `
		SELECT
			id, storefront_id, role, type,
			invited_email, invited_user_id,
			invite_code, expires_at, max_uses, current_uses,
			invited_by_id, status, comment,
			created_at, updated_at, accepted_at, declined_at
		FROM storefront_invitations
		WHERE id = $1`

	inv := &domain.StorefrontInvitation{}
	err := r.scanInvitation(r.db.QueryRowContext(ctx, query, id), inv)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invitation not found")
		}
		return nil, fmt.Errorf("failed to get invitation: %w", err)
	}

	return inv, nil
}

// GetByCode returns invitation by invite code
func (r *StorefrontInvitationRepository) GetByCode(ctx context.Context, code string) (*domain.StorefrontInvitation, error) {
	query := `
		SELECT
			id, storefront_id, role, type,
			invited_email, invited_user_id,
			invite_code, expires_at, max_uses, current_uses,
			invited_by_id, status, comment,
			created_at, updated_at, accepted_at, declined_at
		FROM storefront_invitations
		WHERE invite_code = $1`

	inv := &domain.StorefrontInvitation{}
	err := r.scanInvitation(r.db.QueryRowContext(ctx, query, code), inv)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invitation not found")
		}
		return nil, fmt.Errorf("failed to get invitation by code: %w", err)
	}

	return inv, nil
}

// GetByEmailAndStorefront returns invitation by email and storefront (for checking duplicates)
func (r *StorefrontInvitationRepository) GetByEmailAndStorefront(ctx context.Context, email string, storefrontID int64) (*domain.StorefrontInvitation, error) {
	query := `
		SELECT
			id, storefront_id, role, type,
			invited_email, invited_user_id,
			invite_code, expires_at, max_uses, current_uses,
			invited_by_id, status, comment,
			created_at, updated_at, accepted_at, declined_at
		FROM storefront_invitations
		WHERE invited_email = $1 AND storefront_id = $2 AND status = 'pending'
		LIMIT 1`

	inv := &domain.StorefrontInvitation{}
	err := r.scanInvitation(r.db.QueryRowContext(ctx, query, email, storefrontID), inv)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No pending invitation found - not an error
		}
		return nil, fmt.Errorf("failed to get invitation by email: %w", err)
	}

	return inv, nil
}

// ListByStorefront returns invitations for a storefront with filtering and pagination
func (r *StorefrontInvitationRepository) ListByStorefront(ctx context.Context, storefrontID int64, filter *domain.ListInvitationsFilter) ([]*domain.StorefrontInvitation, int32, error) {
	whereConditions := []string{"storefront_id = $1"}
	args := []interface{}{storefrontID}
	argPos := 2

	if filter.Status != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("status = $%d", argPos))
		args = append(args, *filter.Status)
		argPos++
	}

	if filter.Type != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("type = $%d", argPos))
		args = append(args, *filter.Type)
		argPos++
	}

	if filter.InvitedByID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("invited_by_id = $%d", argPos))
		args = append(args, *filter.InvitedByID)
		argPos++
	}

	whereClause := strings.Join(whereConditions, " AND ")

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM storefront_invitations WHERE %s", whereClause)
	var total int32
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count invitations: %w", err)
	}

	// Get paginated results
	limit := int32(10)
	if filter.Limit > 0 {
		limit = filter.Limit
	}
	if limit > 100 {
		limit = 100
	}

	page := int32(1)
	if filter.Page > 0 {
		page = filter.Page
	}
	offset := (page - 1) * limit

	args = append(args, limit, offset)
	query := fmt.Sprintf(`
		SELECT
			id, storefront_id, role, type,
			invited_email, invited_user_id,
			invite_code, expires_at, max_uses, current_uses,
			invited_by_id, status, comment,
			created_at, updated_at, accepted_at, declined_at
		FROM storefront_invitations
		WHERE %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, argPos, argPos+1)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list invitations: %w", err)
	}
	defer rows.Close()

	var invitations []*domain.StorefrontInvitation
	for rows.Next() {
		inv := &domain.StorefrontInvitation{}
		if err := r.scanInvitationFromRows(rows, inv); err != nil {
			return nil, 0, fmt.Errorf("failed to scan invitation: %w", err)
		}
		invitations = append(invitations, inv)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows iteration error: %w", err)
	}

	return invitations, total, nil
}

// ListByUser returns pending invitations for a user (by email or user_id)
func (r *StorefrontInvitationRepository) ListByUser(ctx context.Context, userID int64, email string) ([]*domain.StorefrontInvitation, error) {
	query := `
		SELECT
			id, storefront_id, role, type,
			invited_email, invited_user_id,
			invite_code, expires_at, max_uses, current_uses,
			invited_by_id, status, comment,
			created_at, updated_at, accepted_at, declined_at
		FROM storefront_invitations
		WHERE (invited_user_id = $1 OR invited_email = $2)
			AND status = 'pending'
			AND (expires_at IS NULL OR expires_at > NOW())
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID, email)
	if err != nil {
		return nil, fmt.Errorf("failed to list user invitations: %w", err)
	}
	defer rows.Close()

	var invitations []*domain.StorefrontInvitation
	for rows.Next() {
		inv := &domain.StorefrontInvitation{}
		if err := r.scanInvitationFromRows(rows, inv); err != nil {
			return nil, fmt.Errorf("failed to scan invitation: %w", err)
		}
		invitations = append(invitations, inv)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return invitations, nil
}

// Update updates an invitation
func (r *StorefrontInvitationRepository) Update(ctx context.Context, inv *domain.StorefrontInvitation) error {
	query := `
		UPDATE storefront_invitations
		SET
			status = $2,
			current_uses = $3,
			accepted_at = $4,
			declined_at = $5,
			comment = $6,
			updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at`

	return r.db.QueryRowContext(ctx, query,
		inv.ID,
		inv.Status,
		inv.CurrentUses,
		inv.AcceptedAt,
		inv.DeclinedAt,
		inv.Comment,
	).Scan(&inv.UpdatedAt)
}

// Delete deletes an invitation (hard delete)
func (r *StorefrontInvitationRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM storefront_invitations WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete invitation: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("invitation not found")
	}

	return nil
}

// IncrementUses atomically increments the current_uses counter
func (r *StorefrontInvitationRepository) IncrementUses(ctx context.Context, id int64) error {
	query := `
		UPDATE storefront_invitations
		SET current_uses = current_uses + 1, updated_at = NOW()
		WHERE id = $1
			AND (max_uses IS NULL OR current_uses < max_uses)`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to increment uses: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("invitation not found or max uses reached")
	}

	return nil
}

// MarkAsExpired marks an invitation as expired
func (r *StorefrontInvitationRepository) MarkAsExpired(ctx context.Context, id int64) error {
	query := `
		UPDATE storefront_invitations
		SET status = 'expired', updated_at = NOW()
		WHERE id = $1 AND status = 'pending'`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to mark as expired: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("invitation not found or not pending")
	}

	return nil
}

// ExpirePendingInvitations marks all expired pending invitations as expired (batch operation)
func (r *StorefrontInvitationRepository) ExpirePendingInvitations(ctx context.Context) (int64, error) {
	query := `
		UPDATE storefront_invitations
		SET status = 'expired', updated_at = NOW()
		WHERE status = 'pending'
			AND expires_at IS NOT NULL
			AND expires_at < NOW()`

	result, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to expire pending invitations: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}

// MarkAsRevoked marks an invitation as revoked
func (r *StorefrontInvitationRepository) MarkAsRevoked(ctx context.Context, id int64) error {
	query := `
		UPDATE storefront_invitations
		SET status = 'revoked', updated_at = NOW()
		WHERE id = $1 AND status = 'pending'`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to mark as revoked: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("invitation not found or not pending")
	}

	return nil
}

// scanInvitation scans a single row into StorefrontInvitation
func (r *StorefrontInvitationRepository) scanInvitation(row *sql.Row, inv *domain.StorefrontInvitation) error {
	var invitedEmail, inviteCode, comment sql.NullString
	var invitedUserID sql.NullInt64
	var expiresAt, acceptedAt, declinedAt sql.NullTime
	var maxUses sql.NullInt32

	err := row.Scan(
		&inv.ID,
		&inv.StorefrontID,
		&inv.Role,
		&inv.Type,
		&invitedEmail,
		&invitedUserID,
		&inviteCode,
		&expiresAt,
		&maxUses,
		&inv.CurrentUses,
		&inv.InvitedByID,
		&inv.Status,
		&comment,
		&inv.CreatedAt,
		&inv.UpdatedAt,
		&acceptedAt,
		&declinedAt,
	)
	if err != nil {
		return err
	}

	// Handle nullable fields
	if invitedEmail.Valid {
		inv.InvitedEmail = &invitedEmail.String
	}
	if invitedUserID.Valid {
		inv.InvitedUserID = &invitedUserID.Int64
	}
	if inviteCode.Valid {
		inv.InviteCode = &inviteCode.String
	}
	if expiresAt.Valid {
		inv.ExpiresAt = &expiresAt.Time
	}
	if maxUses.Valid {
		inv.MaxUses = &maxUses.Int32
	}
	if comment.Valid {
		inv.Comment = comment.String
	}
	if acceptedAt.Valid {
		inv.AcceptedAt = &acceptedAt.Time
	}
	if declinedAt.Valid {
		inv.DeclinedAt = &declinedAt.Time
	}

	return nil
}

// scanInvitationFromRows scans rows into StorefrontInvitation
func (r *StorefrontInvitationRepository) scanInvitationFromRows(rows *sql.Rows, inv *domain.StorefrontInvitation) error {
	var invitedEmail, inviteCode, comment sql.NullString
	var invitedUserID sql.NullInt64
	var expiresAt, acceptedAt, declinedAt sql.NullTime
	var maxUses sql.NullInt32

	err := rows.Scan(
		&inv.ID,
		&inv.StorefrontID,
		&inv.Role,
		&inv.Type,
		&invitedEmail,
		&invitedUserID,
		&inviteCode,
		&expiresAt,
		&maxUses,
		&inv.CurrentUses,
		&inv.InvitedByID,
		&inv.Status,
		&comment,
		&inv.CreatedAt,
		&inv.UpdatedAt,
		&acceptedAt,
		&declinedAt,
	)
	if err != nil {
		return err
	}

	// Handle nullable fields
	if invitedEmail.Valid {
		inv.InvitedEmail = &invitedEmail.String
	}
	if invitedUserID.Valid {
		inv.InvitedUserID = &invitedUserID.Int64
	}
	if inviteCode.Valid {
		inv.InviteCode = &inviteCode.String
	}
	if expiresAt.Valid {
		inv.ExpiresAt = &expiresAt.Time
	}
	if maxUses.Valid {
		inv.MaxUses = &maxUses.Int32
	}
	if comment.Valid {
		inv.Comment = comment.String
	}
	if acceptedAt.Valid {
		inv.AcceptedAt = &acceptedAt.Time
	}
	if declinedAt.Valid {
		inv.DeclinedAt = &declinedAt.Time
	}

	return nil
}

// GetStatsByStorefront returns invitation statistics for a storefront
func (r *StorefrontInvitationRepository) GetStatsByStorefront(ctx context.Context, storefrontID int64) (map[string]int32, error) {
	query := `
		SELECT
			status,
			COUNT(*) as count
		FROM storefront_invitations
		WHERE storefront_id = $1
		GROUP BY status`

	rows, err := r.db.QueryContext(ctx, query, storefrontID)
	if err != nil {
		return nil, fmt.Errorf("failed to get invitation stats: %w", err)
	}
	defer rows.Close()

	stats := make(map[string]int32)
	for rows.Next() {
		var status string
		var count int32
		if err := rows.Scan(&status, &count); err != nil {
			return nil, fmt.Errorf("failed to scan stats: %w", err)
		}
		stats[status] = count
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return stats, nil
}

// CheckInviteCodeExists checks if an invite code already exists
func (r *StorefrontInvitationRepository) CheckInviteCodeExists(ctx context.Context, code string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM storefront_invitations WHERE invite_code = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, code).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check invite code: %w", err)
	}

	return exists, nil
}

// GetActiveLinkInvitations returns all active link invitations for a storefront
func (r *StorefrontInvitationRepository) GetActiveLinkInvitations(ctx context.Context, storefrontID int64) ([]*domain.StorefrontInvitation, error) {
	query := `
		SELECT
			id, storefront_id, role, type,
			invited_email, invited_user_id,
			invite_code, expires_at, max_uses, current_uses,
			invited_by_id, status, comment,
			created_at, updated_at, accepted_at, declined_at
		FROM storefront_invitations
		WHERE storefront_id = $1
			AND type = 'link'
			AND status = 'pending'
			AND (expires_at IS NULL OR expires_at > NOW())
			AND (max_uses IS NULL OR current_uses < max_uses)
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, storefrontID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active link invitations: %w", err)
	}
	defer rows.Close()

	var invitations []*domain.StorefrontInvitation
	for rows.Next() {
		inv := &domain.StorefrontInvitation{}
		if err := r.scanInvitationFromRows(rows, inv); err != nil {
			return nil, fmt.Errorf("failed to scan invitation: %w", err)
		}
		invitations = append(invitations, inv)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return invitations, nil
}

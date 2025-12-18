// Package postgres implements PostgreSQL repositories for the listings microservice.
// This file contains the StockReservationRepository implementation.
package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/vondi-global/listings/internal/domain"
)

// StockReservationRepository handles database operations for stock reservations
type StockReservationRepository struct {
	db     *sqlx.DB
	logger zerolog.Logger
}

// NewStockReservationRepository creates a new stock reservation repository
func NewStockReservationRepository(db *sqlx.DB, logger zerolog.Logger) *StockReservationRepository {
	return &StockReservationRepository{
		db:     db,
		logger: logger,
	}
}

// Create creates a new stock reservation within a transaction
func (r *StockReservationRepository) Create(ctx context.Context, tx *sqlx.Tx, reservation *domain.StockReservation) error {
	query := `
		INSERT INTO stock_reservations (
			id, variant_id, order_id, quantity, expires_at, status
		) VALUES (
			$1, $2, $3, $4, $5, $6
		)
		RETURNING created_at, updated_at
	`

	err := tx.QueryRowContext(ctx, query,
		reservation.ID,
		reservation.VariantID,
		reservation.OrderID,
		reservation.Quantity,
		reservation.ExpiresAt,
		reservation.Status,
	).Scan(&reservation.CreatedAt, &reservation.UpdatedAt)

	if err != nil {
		r.logger.Error().Err(err).
			Str("variant_id", reservation.VariantID.String()).
			Str("order_id", reservation.OrderID.String()).
			Msg("failed to create stock reservation")
		return fmt.Errorf("failed to create stock reservation: %w", err)
	}

	return nil
}

// GetByID retrieves a reservation by ID
func (r *StockReservationRepository) GetByID(ctx context.Context, id string) (*domain.StockReservation, error) {
	reservationID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid reservation ID: %w", err)
	}

	query := `
		SELECT
			id, variant_id, order_id, quantity, expires_at, status,
			created_at, updated_at
		FROM stock_reservations
		WHERE id = $1
	`

	reservation := &domain.StockReservation{}

	err = r.db.QueryRowContext(ctx, query, reservationID).Scan(
		&reservation.ID,
		&reservation.VariantID,
		&reservation.OrderID,
		&reservation.Quantity,
		&reservation.ExpiresAt,
		&reservation.Status,
		&reservation.CreatedAt,
		&reservation.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrStockReservationNotFound
	}
	if err != nil {
		r.logger.Error().Err(err).Str("id", id).Msg("failed to get reservation by ID")
		return nil, fmt.Errorf("failed to get reservation: %w", err)
	}

	return reservation, nil
}

// List retrieves reservations based on filters
func (r *StockReservationRepository) List(ctx context.Context, filter *domain.ListReservationsFilter) ([]*domain.StockReservation, error) {
	query := `
		SELECT
			id, variant_id, order_id, quantity, expires_at, status,
			created_at, updated_at
		FROM stock_reservations
		WHERE 1=1
	`

	args := make([]interface{}, 0)
	argIdx := 1

	if filter.VariantID != nil {
		query += fmt.Sprintf(" AND variant_id = $%d", argIdx)
		args = append(args, *filter.VariantID)
		argIdx++
	}

	if filter.OrderID != nil {
		query += fmt.Sprintf(" AND order_id = $%d", argIdx)
		args = append(args, *filter.OrderID)
		argIdx++
	}

	if filter.Status != nil {
		query += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, *filter.Status)
	}

	if filter.ActiveOnly {
		query += " AND status = 'active' AND expires_at > CURRENT_TIMESTAMP"
	}

	if filter.ExpiredOnly {
		query += " AND status = 'active' AND expires_at <= CURRENT_TIMESTAMP"
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.Error().Err(err).Interface("filter", filter).Msg("failed to list reservations")
		return nil, fmt.Errorf("failed to list reservations: %w", err)
	}
	defer rows.Close()

	reservations := make([]*domain.StockReservation, 0)

	for rows.Next() {
		reservation := &domain.StockReservation{}
		err = rows.Scan(
			&reservation.ID,
			&reservation.VariantID,
			&reservation.OrderID,
			&reservation.Quantity,
			&reservation.ExpiresAt,
			&reservation.Status,
			&reservation.CreatedAt,
			&reservation.UpdatedAt,
		)
		if err != nil {
			r.logger.Error().Err(err).Msg("failed to scan reservation row")
			continue
		}

		reservations = append(reservations, reservation)
	}

	return reservations, nil
}

// UpdateStatus updates the status of a reservation within a transaction
func (r *StockReservationRepository) UpdateStatus(ctx context.Context, tx *sqlx.Tx, id string, status string) error {
	reservationID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid reservation ID: %w", err)
	}

	query := `
		UPDATE stock_reservations
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	result, err := tx.ExecContext(ctx, query, status, reservationID)
	if err != nil {
		r.logger.Error().Err(err).Str("id", id).Str("status", status).Msg("failed to update reservation status")
		return fmt.Errorf("failed to update reservation status: %w", err)
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return domain.ErrStockReservationNotFound
	}

	return nil
}

// CleanupExpired marks expired reservations as 'expired' and returns count
func (r *StockReservationRepository) CleanupExpired(ctx context.Context) (int64, error) {
	query := `
		UPDATE stock_reservations
		SET status = 'expired', updated_at = CURRENT_TIMESTAMP
		WHERE status = 'active' AND expires_at < CURRENT_TIMESTAMP
	`

	result, err := r.db.ExecContext(ctx, query)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to cleanup expired reservations")
		return 0, fmt.Errorf("failed to cleanup expired reservations: %w", err)
	}

	count, _ := result.RowsAffected()

	if count > 0 {
		r.logger.Info().Int64("count", count).Msg("cleaned up expired stock reservations")
	}

	return count, nil
}

// GetActiveByVariant retrieves all active reservations for a variant
func (r *StockReservationRepository) GetActiveByVariant(ctx context.Context, variantID string) ([]*domain.StockReservation, error) {
	vid, err := uuid.Parse(variantID)
	if err != nil {
		return nil, fmt.Errorf("invalid variant ID: %w", err)
	}

	filter := &domain.ListReservationsFilter{
		VariantID:  &vid,
		ActiveOnly: true,
	}

	return r.List(ctx, filter)
}

// GetActiveByOrder retrieves all active reservations for an order
func (r *StockReservationRepository) GetActiveByOrder(ctx context.Context, orderID string) ([]*domain.StockReservation, error) {
	oid, err := uuid.Parse(orderID)
	if err != nil {
		return nil, fmt.Errorf("invalid order ID: %w", err)
	}

	filter := &domain.ListReservationsFilter{
		OrderID:    &oid,
		ActiveOnly: true,
	}

	return r.List(ctx, filter)
}

// Delete deletes a reservation (use with caution, prefer UpdateStatus)
func (r *StockReservationRepository) Delete(ctx context.Context, id string) error {
	reservationID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid reservation ID: %w", err)
	}

	query := `DELETE FROM stock_reservations WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, reservationID)
	if err != nil {
		r.logger.Error().Err(err).Str("id", id).Msg("failed to delete reservation")
		return fmt.Errorf("failed to delete reservation: %w", err)
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return domain.ErrStockReservationNotFound
	}

	return nil
}

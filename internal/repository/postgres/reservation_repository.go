// Package postgres implements PostgreSQL repository layer for listings microservice.
package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	"github.com/vondi-global/listings/internal/domain"
)

// ReservationRepository defines operations for inventory reservation management
type ReservationRepository interface {
	// CRUD operations
	Create(ctx context.Context, reservation *domain.InventoryReservation) error
	GetByID(ctx context.Context, reservationID int64) (*domain.InventoryReservation, error)
	GetByOrderID(ctx context.Context, orderID int64) ([]*domain.InventoryReservation, error)
	GetActiveByListing(ctx context.Context, listingID, variantID int64) ([]*domain.InventoryReservation, error)
	Update(ctx context.Context, reservation *domain.InventoryReservation) error
	Delete(ctx context.Context, reservationID int64) error

	// Batch operations
	CommitReservations(ctx context.Context, orderID int64) error
	ReleaseReservations(ctx context.Context, orderID int64) error

	// Cleanup
	ExpireStaleReservations(ctx context.Context) (int, error)

	// Transaction support
	WithTx(tx pgx.Tx) ReservationRepository
}

// reservationRepository implements ReservationRepository using PostgreSQL
type reservationRepository struct {
	db     dbOrTx
	logger zerolog.Logger
}

// NewReservationRepository creates a new reservation repository
func NewReservationRepository(pool *pgxpool.Pool, logger zerolog.Logger) ReservationRepository {
	return &reservationRepository{
		db:     pool,
		logger: logger.With().Str("component", "reservation_repository").Logger(),
	}
}

// WithTx returns a new repository instance using the provided transaction
func (r *reservationRepository) WithTx(tx pgx.Tx) ReservationRepository {
	return &reservationRepository{
		db:     tx,
		logger: r.logger,
	}
}

// Create creates a new inventory reservation
func (r *reservationRepository) Create(ctx context.Context, reservation *domain.InventoryReservation) error {
	if err := reservation.Validate(); err != nil {
		return fmt.Errorf("invalid reservation: %w", err)
	}

	query := `
		INSERT INTO inventory_reservations (
			listing_id, variant_id, order_id, quantity, status, expires_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		reservation.ListingID,
		reservation.VariantID,
		reservation.OrderID,
		reservation.Quantity,
		string(reservation.Status),
		reservation.ExpiresAt,
	).Scan(&reservation.ID, &reservation.CreatedAt, &reservation.UpdatedAt)

	if err != nil {
		r.logger.Error().Err(err).
			Int64("listing_id", reservation.ListingID).
			Int64("order_id", reservation.OrderID).
			Msg("failed to create reservation")
		return fmt.Errorf("failed to create reservation: %w", err)
	}

	r.logger.Info().
		Int64("reservation_id", reservation.ID).
		Int64("listing_id", reservation.ListingID).
		Int64("order_id", reservation.OrderID).
		Int32("quantity", reservation.Quantity).
		Msg("reservation created")
	return nil
}

// GetByID retrieves a reservation by its ID
func (r *reservationRepository) GetByID(ctx context.Context, reservationID int64) (*domain.InventoryReservation, error) {
	query := `
		SELECT id, listing_id, variant_id, order_id, quantity, status,
		       expires_at, created_at, updated_at, committed_at, released_at
		FROM inventory_reservations
		WHERE id = $1
	`

	var reservation domain.InventoryReservation
	var variantID sql.NullInt64
	var statusStr string
	var committedAt, releasedAt sql.NullTime

	err := r.db.QueryRow(ctx, query, reservationID).Scan(
		&reservation.ID,
		&reservation.ListingID,
		&variantID,
		&reservation.OrderID,
		&reservation.Quantity,
		&statusStr,
		&reservation.ExpiresAt,
		&reservation.CreatedAt,
		&reservation.UpdatedAt,
		&committedAt,
		&releasedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("reservation not found")
		}
		r.logger.Error().Err(err).Int64("reservation_id", reservationID).Msg("failed to get reservation by ID")
		return nil, fmt.Errorf("failed to get reservation by ID: %w", err)
	}

	// Handle nullable fields
	if variantID.Valid {
		reservation.VariantID = &variantID.Int64
	}
	reservation.Status = domain.ReservationStatus(statusStr)
	if committedAt.Valid {
		reservation.CommittedAt = &committedAt.Time
	}
	if releasedAt.Valid {
		reservation.ReleasedAt = &releasedAt.Time
	}

	return &reservation, nil
}

// GetByOrderID retrieves all reservations for an order
func (r *reservationRepository) GetByOrderID(ctx context.Context, orderID int64) ([]*domain.InventoryReservation, error) {
	query := `
		SELECT id, listing_id, variant_id, order_id, quantity, status,
		       expires_at, created_at, updated_at, committed_at, released_at
		FROM inventory_reservations
		WHERE order_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(ctx, query, orderID)
	if err != nil {
		r.logger.Error().Err(err).Int64("order_id", orderID).Msg("failed to get reservations by order ID")
		return nil, fmt.Errorf("failed to get reservations by order ID: %w", err)
	}
	defer rows.Close()

	var reservations []*domain.InventoryReservation
	for rows.Next() {
		var reservation domain.InventoryReservation
		var variantID sql.NullInt64
		var statusStr string
		var committedAt, releasedAt sql.NullTime

		err := rows.Scan(
			&reservation.ID,
			&reservation.ListingID,
			&variantID,
			&reservation.OrderID,
			&reservation.Quantity,
			&statusStr,
			&reservation.ExpiresAt,
			&reservation.CreatedAt,
			&reservation.UpdatedAt,
			&committedAt,
			&releasedAt,
		)
		if err != nil {
			r.logger.Error().Err(err).Msg("failed to scan reservation")
			return nil, fmt.Errorf("failed to scan reservation: %w", err)
		}

		// Handle nullable fields
		if variantID.Valid {
			reservation.VariantID = &variantID.Int64
		}
		reservation.Status = domain.ReservationStatus(statusStr)
		if committedAt.Valid {
			reservation.CommittedAt = &committedAt.Time
		}
		if releasedAt.Valid {
			reservation.ReleasedAt = &releasedAt.Time
		}

		reservations = append(reservations, &reservation)
	}

	if err = rows.Err(); err != nil {
		r.logger.Error().Err(err).Msg("error iterating reservation rows")
		return nil, fmt.Errorf("error iterating reservation rows: %w", err)
	}

	return reservations, nil
}

// GetActiveByListing retrieves all active reservations for a listing (and optional variant)
func (r *reservationRepository) GetActiveByListing(ctx context.Context, listingID, variantID int64) ([]*domain.InventoryReservation, error) {
	var query string
	var rows pgx.Rows
	var err error

	if variantID > 0 {
		query = `
			SELECT id, listing_id, variant_id, order_id, quantity, status,
			       expires_at, created_at, updated_at, committed_at, released_at
			FROM inventory_reservations
			WHERE listing_id = $1 AND variant_id = $2 AND status = 'active' AND expires_at > NOW()
			ORDER BY created_at ASC
		`
		rows, err = r.db.Query(ctx, query, listingID, variantID)
	} else {
		query = `
			SELECT id, listing_id, variant_id, order_id, quantity, status,
			       expires_at, created_at, updated_at, committed_at, released_at
			FROM inventory_reservations
			WHERE listing_id = $1 AND variant_id IS NULL AND status = 'active' AND expires_at > NOW()
			ORDER BY created_at ASC
		`
		rows, err = r.db.Query(ctx, query, listingID)
	}

	if err != nil {
		r.logger.Error().Err(err).Int64("listing_id", listingID).Msg("failed to get active reservations")
		return nil, fmt.Errorf("failed to get active reservations: %w", err)
	}
	defer rows.Close()

	var reservations []*domain.InventoryReservation
	for rows.Next() {
		var reservation domain.InventoryReservation
		var varID sql.NullInt64
		var statusStr string
		var committedAt, releasedAt sql.NullTime

		err := rows.Scan(
			&reservation.ID,
			&reservation.ListingID,
			&varID,
			&reservation.OrderID,
			&reservation.Quantity,
			&statusStr,
			&reservation.ExpiresAt,
			&reservation.CreatedAt,
			&reservation.UpdatedAt,
			&committedAt,
			&releasedAt,
		)
		if err != nil {
			r.logger.Error().Err(err).Msg("failed to scan reservation")
			return nil, fmt.Errorf("failed to scan reservation: %w", err)
		}

		// Handle nullable fields
		if varID.Valid {
			reservation.VariantID = &varID.Int64
		}
		reservation.Status = domain.ReservationStatus(statusStr)
		if committedAt.Valid {
			reservation.CommittedAt = &committedAt.Time
		}
		if releasedAt.Valid {
			reservation.ReleasedAt = &releasedAt.Time
		}

		reservations = append(reservations, &reservation)
	}

	if err = rows.Err(); err != nil {
		r.logger.Error().Err(err).Msg("error iterating reservation rows")
		return nil, fmt.Errorf("error iterating reservation rows: %w", err)
	}

	return reservations, nil
}

// Update updates an existing reservation
func (r *reservationRepository) Update(ctx context.Context, reservation *domain.InventoryReservation) error {
	if err := reservation.Validate(); err != nil {
		return fmt.Errorf("invalid reservation: %w", err)
	}

	query := `
		UPDATE inventory_reservations
		SET quantity = $1, status = $2, expires_at = $3, committed_at = $4, released_at = $5
		WHERE id = $6
		RETURNING updated_at
	`

	err := r.db.QueryRow(ctx, query,
		reservation.Quantity,
		string(reservation.Status),
		reservation.ExpiresAt,
		reservation.CommittedAt,
		reservation.ReleasedAt,
		reservation.ID,
	).Scan(&reservation.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("reservation not found")
		}
		r.logger.Error().Err(err).Int64("reservation_id", reservation.ID).Msg("failed to update reservation")
		return fmt.Errorf("failed to update reservation: %w", err)
	}

	r.logger.Info().Int64("reservation_id", reservation.ID).Str("status", string(reservation.Status)).Msg("reservation updated")
	return nil
}

// Delete deletes a reservation by ID
func (r *reservationRepository) Delete(ctx context.Context, reservationID int64) error {
	query := `DELETE FROM inventory_reservations WHERE id = $1`

	result, err := r.db.Exec(ctx, query, reservationID)
	if err != nil {
		r.logger.Error().Err(err).Int64("reservation_id", reservationID).Msg("failed to delete reservation")
		return fmt.Errorf("failed to delete reservation: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("reservation not found")
	}

	r.logger.Info().Int64("reservation_id", reservationID).Msg("reservation deleted")
	return nil
}

// CommitReservations commits all reservations for an order
func (r *reservationRepository) CommitReservations(ctx context.Context, orderID int64) error {
	query := `
		UPDATE inventory_reservations
		SET status = 'committed', committed_at = NOW(), updated_at = NOW()
		WHERE order_id = $1 AND status = 'active'
		RETURNING id
	`

	rows, err := r.db.Query(ctx, query, orderID)
	if err != nil {
		r.logger.Error().Err(err).Int64("order_id", orderID).Msg("failed to commit reservations")
		return fmt.Errorf("failed to commit reservations: %w", err)
	}
	defer rows.Close()

	committedCount := 0
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return fmt.Errorf("failed to scan committed reservation ID: %w", err)
		}
		committedCount++
	}

	if err = rows.Err(); err != nil {
		return fmt.Errorf("error iterating committed reservations: %w", err)
	}

	r.logger.Info().Int64("order_id", orderID).Int("committed_count", committedCount).Msg("reservations committed")
	return nil
}

// ReleaseReservations releases all reservations for an order
func (r *reservationRepository) ReleaseReservations(ctx context.Context, orderID int64) error {
	query := `
		UPDATE inventory_reservations
		SET status = 'released', released_at = NOW(), updated_at = NOW()
		WHERE order_id = $1 AND status = 'active'
		RETURNING id
	`

	rows, err := r.db.Query(ctx, query, orderID)
	if err != nil {
		r.logger.Error().Err(err).Int64("order_id", orderID).Msg("failed to release reservations")
		return fmt.Errorf("failed to release reservations: %w", err)
	}
	defer rows.Close()

	releasedCount := 0
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return fmt.Errorf("failed to scan released reservation ID: %w", err)
		}
		releasedCount++
	}

	if err = rows.Err(); err != nil {
		return fmt.Errorf("error iterating released reservations: %w", err)
	}

	r.logger.Info().Int64("order_id", orderID).Int("released_count", releasedCount).Msg("reservations released")
	return nil
}

// ExpireStaleReservations marks all expired active reservations as expired
func (r *reservationRepository) ExpireStaleReservations(ctx context.Context) (int, error) {
	query := `
		UPDATE inventory_reservations
		SET status = 'expired', updated_at = NOW()
		WHERE status = 'active' AND expires_at < NOW()
		RETURNING id
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to expire stale reservations")
		return 0, fmt.Errorf("failed to expire stale reservations: %w", err)
	}
	defer rows.Close()

	expiredCount := 0
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return expiredCount, fmt.Errorf("failed to scan expired reservation ID: %w", err)
		}
		expiredCount++
	}

	if err = rows.Err(); err != nil {
		return expiredCount, fmt.Errorf("error iterating expired reservations: %w", err)
	}

	if expiredCount > 0 {
		r.logger.Info().Int("expired_count", expiredCount).Msg("stale reservations expired")
	}

	return expiredCount, nil
}

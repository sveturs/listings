// Package service provides business logic layer for the listings microservice.
package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	"github.com/sveturs/listings/internal/domain"
	"github.com/sveturs/listings/internal/repository/postgres"
)

// InventoryService defines business logic operations for inventory and reservation management
type InventoryService interface {
	// Reservation management
	CreateReservation(ctx context.Context, req *CreateReservationRequest) (*domain.InventoryReservation, error)
	GetReservationsByOrder(ctx context.Context, orderID int64) ([]*domain.InventoryReservation, error)
	CommitReservation(ctx context.Context, reservationID int64) error
	ReleaseReservation(ctx context.Context, reservationID int64) error

	// Cleanup (called by cron job)
	CleanupExpiredReservations(ctx context.Context) (int, error)

	// Stock checks
	CheckStockAvailability(ctx context.Context, listingID int64, quantity int) (bool, error)
	GetAvailableStock(ctx context.Context, listingID int64) (int, error)
}

// CreateReservationRequest contains parameters for creating a reservation
type CreateReservationRequest struct {
	ListingID  int64
	VariantID  *int64
	OrderID    int64
	Quantity   int32
	TTLMinutes int // Optional, defaults to 30
}

// inventoryService implements InventoryService
type inventoryService struct {
	reservationRepo postgres.ReservationRepository
	productsRepo    *postgres.Repository
	orderRepo       postgres.OrderRepository
	pool            *pgxpool.Pool
	logger          zerolog.Logger
}

// NewInventoryService creates a new inventory service
func NewInventoryService(
	reservationRepo postgres.ReservationRepository,
	productsRepo *postgres.Repository,
	orderRepo postgres.OrderRepository,
	pool *pgxpool.Pool,
	logger zerolog.Logger,
) InventoryService {
	return &inventoryService{
		reservationRepo: reservationRepo,
		productsRepo:    productsRepo,
		orderRepo:       orderRepo,
		pool:            pool,
		logger:          logger.With().Str("component", "inventory_service").Logger(),
	}
}

// CreateReservation creates a new inventory reservation
func (s *inventoryService) CreateReservation(ctx context.Context, req *CreateReservationRequest) (*domain.InventoryReservation, error) {
	s.logger.Debug().
		Int64("listing_id", req.ListingID).
		Int64("order_id", req.OrderID).
		Int32("quantity", req.Quantity).
		Msg("creating reservation")

	// Validate quantity
	if req.Quantity <= 0 {
		return nil, fmt.Errorf("%w: quantity must be greater than 0", ErrInvalidInput)
	}

	// Get listing to validate stock
	listing, err := s.productsRepo.GetProductByID(ctx, req.ListingID, nil)
	if err != nil {
		return nil, &ErrListingNotFound{ListingID: req.ListingID}
	}

	// Check available stock (total stock - active reservations)
	availableStock, err := s.GetAvailableStock(ctx, req.ListingID)
	if err != nil {
		return nil, fmt.Errorf("failed to get available stock: %w", err)
	}

	if availableStock < int(req.Quantity) {
		return nil, &ErrStockNotAvailable{
			ListingID:      req.ListingID,
			RequestedQty:   req.Quantity,
			TotalStock:     listing.StockQuantity,
			ReservedStock:  listing.StockQuantity - int32(availableStock),
			AvailableStock: int32(availableStock),
		}
	}

	// Create reservation with TTL
	var reservation *domain.InventoryReservation
	if req.TTLMinutes > 0 {
		reservation = domain.NewInventoryReservationWithTTL(
			req.ListingID,
			req.VariantID,
			req.OrderID,
			req.Quantity,
			req.TTLMinutes,
		)
	} else {
		reservation = domain.NewInventoryReservation(
			req.ListingID,
			req.VariantID,
			req.OrderID,
			req.Quantity,
		)
	}

	if err := s.reservationRepo.Create(ctx, reservation); err != nil {
		s.logger.Error().Err(err).Msg("failed to create reservation")
		return nil, fmt.Errorf("failed to create reservation: %w", err)
	}

	s.logger.Info().Int64("reservation_id", reservation.ID).Msg("reservation created")
	return reservation, nil
}

// GetReservationsByOrder retrieves all reservations for an order
func (s *inventoryService) GetReservationsByOrder(ctx context.Context, orderID int64) ([]*domain.InventoryReservation, error) {
	reservations, err := s.reservationRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		s.logger.Error().Err(err).Int64("order_id", orderID).Msg("failed to get reservations")
		return nil, fmt.Errorf("failed to get reservations: %w", err)
	}
	return reservations, nil
}

// CommitReservation commits a reservation (marks it as committed, stock already deducted)
func (s *inventoryService) CommitReservation(ctx context.Context, reservationID int64) error {
	s.logger.Debug().Int64("reservation_id", reservationID).Msg("committing reservation")

	// Get reservation
	reservation, err := s.reservationRepo.GetByID(ctx, reservationID)
	if err != nil {
		if err.Error() == "reservation not found" {
			return ErrReservationNotFound
		}
		return fmt.Errorf("failed to get reservation: %w", err)
	}

	// Check if reservation can be committed
	if !reservation.CanCommit() {
		if reservation.IsExpired() {
			return &ErrReservationExpired{ReservationID: reservationID}
		}
		return &ErrReservationCannotCommit{
			ReservationID: reservationID,
			Reason:        fmt.Sprintf("reservation status is %s", reservation.Status),
		}
	}

	// Commit reservation using domain method + Update
	reservation.Commit()
	if err := s.reservationRepo.Update(ctx, reservation); err != nil {
		s.logger.Error().Err(err).Int64("reservation_id", reservationID).Msg("failed to commit reservation")
		return fmt.Errorf("failed to commit reservation: %w", err)
	}

	s.logger.Info().Int64("reservation_id", reservationID).Msg("reservation committed")
	return nil
}

// ReleaseReservation releases a reservation (restores stock)
func (s *inventoryService) ReleaseReservation(ctx context.Context, reservationID int64) error {
	s.logger.Debug().Int64("reservation_id", reservationID).Msg("releasing reservation")

	// Start transaction
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Get reservation
	reservation, err := s.reservationRepo.GetByID(ctx, reservationID)
	if err != nil {
		if err.Error() == "reservation not found" {
			return ErrReservationNotFound
		}
		return fmt.Errorf("failed to get reservation: %w", err)
	}

	// Check if reservation can be released
	if !reservation.CanRelease() {
		return &ErrReservationCannotRelease{
			ReservationID: reservationID,
			Reason:        fmt.Sprintf("reservation status is %s", reservation.Status),
		}
	}

	// Release reservation using domain method
	if err := reservation.Release(); err != nil {
		return fmt.Errorf("failed to release reservation: %w", err)
	}

	// Update reservation in database
	reservationRepoTx := s.reservationRepo.WithTx(tx)
	if err := reservationRepoTx.Update(ctx, reservation); err != nil {
		return fmt.Errorf("failed to update released reservation: %w", err)
	}

	// Restore stock
	if err := s.productsRepo.RestoreStockWithPgxTx(ctx, tx, reservation.ListingID, reservation.Quantity); err != nil {
		s.logger.Error().Err(err).Int64("listing_id", reservation.ListingID).Msg("failed to restore stock")
		return fmt.Errorf("failed to restore stock: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info().Int64("reservation_id", reservationID).Msg("reservation released")
	return nil
}

// CleanupExpiredReservations cleans up expired reservations (cron job, every 5 minutes)
func (s *inventoryService) CleanupExpiredReservations(ctx context.Context) (int, error) {
	s.logger.Info().Msg("cleaning up expired reservations")

	// Start transaction
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// 1. Use ExpireStaleReservations method from repository
	reservationRepoTx := s.reservationRepo.WithTx(tx)
	expiredCount, err := reservationRepoTx.ExpireStaleReservations(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to expire stale reservations")
		return 0, fmt.Errorf("failed to expire stale reservations: %w", err)
	}

	if expiredCount == 0 {
		s.logger.Debug().Msg("no expired reservations found")
		return 0, nil
	}

	s.logger.Info().Int("count", expiredCount).Msg("expired reservations marked")

	// 2. Restore stock for expired reservations
	// Get expired reservations to restore their stock
	if expiredCount > 0 {
		// Note: We need to get the expired reservations that were just marked
		// For now, we'll use a simplified approach - get all expired reservations by order
		// In production, ExpireStaleReservations should return the list of expired reservations
		// TODO: Refactor ExpireStaleReservations to return expired reservations list
		s.logger.Info().Msg("stock restoration for expired reservations needs reservation list from ExpireStaleReservations")
	}

	// 3. Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info().Int("count", expiredCount).Msg("expired reservations cleaned up successfully")
	return expiredCount, nil
}

// CheckStockAvailability checks if stock is available for a listing
func (s *inventoryService) CheckStockAvailability(ctx context.Context, listingID int64, quantity int) (bool, error) {
	availableStock, err := s.GetAvailableStock(ctx, listingID)
	if err != nil {
		return false, err
	}
	return availableStock >= quantity, nil
}

// GetAvailableStock calculates available stock (total stock - active reservations)
func (s *inventoryService) GetAvailableStock(ctx context.Context, listingID int64) (int, error) {
	// Get listing
	listing, err := s.productsRepo.GetProductByID(ctx, listingID, nil)
	if err != nil {
		return 0, &ErrListingNotFound{ListingID: listingID}
	}

	// Get active reservations for this listing (variantID = 0 for non-variant listings)
	reservations, err := s.reservationRepo.GetActiveByListing(ctx, listingID, 0)
	if err != nil {
		s.logger.Error().Err(err).Int64("listing_id", listingID).Msg("failed to get active reservations")
		return 0, fmt.Errorf("failed to get active reservations: %w", err)
	}

	// Calculate reserved quantity
	var reservedQty int32
	for _, res := range reservations {
		if res.Status == domain.ReservationStatusActive {
			reservedQty += res.Quantity
		}
	}

	// Available stock = total stock - reserved stock
	availableStock := int(listing.StockQuantity - reservedQty)
	if availableStock < 0 {
		availableStock = 0
	}

	return availableStock, nil
}

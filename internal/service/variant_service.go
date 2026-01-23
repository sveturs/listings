// Package service contains business logic for the listings microservice.
// This file implements VariantService with stock reservation operations.
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/vondi-global/listings/internal/domain"
	"github.com/vondi-global/listings/internal/repository/postgres"
)

// VariantService handles business logic for product variants and stock operations
type VariantService struct {
	variantRepo     *postgres.VariantRepository
	reservationRepo *postgres.StockReservationRepository
	skuGenerator    *SKUGenerator
	db              *sqlx.DB
	logger          zerolog.Logger
}

// NewVariantService creates a new variant service instance
func NewVariantService(
	variantRepo *postgres.VariantRepository,
	reservationRepo *postgres.StockReservationRepository,
	skuGenerator *SKUGenerator,
	db *sqlx.DB,
	logger zerolog.Logger,
) *VariantService {
	return &VariantService{
		variantRepo:     variantRepo,
		reservationRepo: reservationRepo,
		skuGenerator:    skuGenerator,
		db:              db,
		logger:          logger,
	}
}

// ReserveStockRequest represents a request to reserve stock
type ReserveStockRequest struct {
	VariantID  string
	OrderID    string
	Quantity   int32
	TTLMinutes int32
}

// ReserveStockResponse represents the response from stock reservation
type ReserveStockResponse struct {
	Success        bool
	ReservationID  string
	AvailableAfter int32
	ErrorMessage   string
}

// ReserveStock reserves stock for an order with transaction and row-level locking
// This is a CRITICAL operation that must use transactions to prevent race conditions
func (s *VariantService) ReserveStock(ctx context.Context, req *ReserveStockRequest) (*ReserveStockResponse, error) {
	// Validate request
	if req.Quantity <= 0 {
		return &ReserveStockResponse{
			Success:      false,
			ErrorMessage: "quantity must be positive",
		}, nil
	}

	if req.TTLMinutes <= 0 || req.TTLMinutes > domain.MaxReservationTTL {
		req.TTLMinutes = domain.DefaultReservationTTL
	}

	// 1. Begin transaction
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to begin transaction for stock reservation")
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 2. Lock variant row (SELECT FOR UPDATE)
	variant, err := s.variantRepo.GetForUpdate(ctx, tx, req.VariantID)
	if err != nil {
		s.logger.Error().Err(err).Str("variant_id", req.VariantID).Msg("failed to lock variant for update")
		return nil, fmt.Errorf("failed to lock variant: %w", err)
	}

	// 3. Check stock availability
	available := variant.StockQuantity - variant.ReservedQuantity

	if available < req.Quantity {
		s.logger.Warn().
			Str("variant_id", req.VariantID).
			Int32("requested", req.Quantity).
			Int32("available", available).
			Msg("insufficient stock for reservation")

		return &ReserveStockResponse{
			Success:        false,
			AvailableAfter: available,
			ErrorMessage:   fmt.Sprintf("insufficient stock: requested %d, available %d", req.Quantity, available),
		}, nil
	}

	// 4. Create stock reservation
	orderUUID, err := uuid.Parse(req.OrderID)
	if err != nil {
		return nil, fmt.Errorf("invalid order ID: %w", err)
	}

	reservation := &domain.StockReservation{
		ID:        uuid.New(),
		VariantID: variant.ID,
		OrderID:   orderUUID,
		Quantity:  req.Quantity,
		ExpiresAt: time.Now().Add(time.Duration(req.TTLMinutes) * time.Minute),
		Status:    domain.StockReservationStatusActive,
	}

	if err := s.reservationRepo.Create(ctx, tx, reservation); err != nil {
		s.logger.Error().Err(err).Msg("failed to create stock reservation")
		return nil, fmt.Errorf("failed to create reservation: %w", err)
	}

	// 5. Note: reserved_quantity is automatically updated by DB trigger
	// The trigger sync_variant_reserved_quantity() handles this

	// 6. Commit transaction
	if err := tx.Commit(); err != nil {
		s.logger.Error().Err(err).Msg("failed to commit stock reservation transaction")
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info().
		Str("reservation_id", reservation.ID.String()).
		Str("variant_id", req.VariantID).
		Str("order_id", req.OrderID).
		Int32("quantity", req.Quantity).
		Msg("stock reserved successfully")

	return &ReserveStockResponse{
		Success:        true,
		ReservationID:  reservation.ID.String(),
		AvailableAfter: available - req.Quantity,
	}, nil
}

// ReleaseStock cancels a stock reservation and returns stock to available pool
func (s *VariantService) ReleaseStock(ctx context.Context, reservationID string) error {
	// 1. Begin transaction
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to begin transaction for stock release")
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 2. Get reservation to validate
	reservation, err := s.reservationRepo.GetByID(ctx, reservationID)
	if err != nil {
		return fmt.Errorf("failed to get reservation: %w", err)
	}

	// 3. Check if reservation can be released
	if !reservation.CanBeReleased() {
		return fmt.Errorf("reservation cannot be released: status=%s, expired=%v",
			reservation.Status, reservation.IsExpired())
	}

	// 4. Update reservation status to 'cancelled'
	if err := s.reservationRepo.UpdateStatus(ctx, tx, reservationID, domain.StockReservationStatusCancelled); err != nil {
		s.logger.Error().Err(err).Str("reservation_id", reservationID).Msg("failed to update reservation status")
		return fmt.Errorf("failed to update reservation status: %w", err)
	}

	// 5. Note: reserved_quantity is automatically decreased by DB trigger

	// 6. Commit transaction
	if err := tx.Commit(); err != nil {
		s.logger.Error().Err(err).Msg("failed to commit stock release transaction")
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info().
		Str("reservation_id", reservationID).
		Str("variant_id", reservation.VariantID.String()).
		Int32("quantity", reservation.Quantity).
		Msg("stock released successfully")

	return nil
}

// ConfirmStockDeduction confirms a reservation and deducts stock from inventory
// This is called when order is confirmed/paid
func (s *VariantService) ConfirmStockDeduction(ctx context.Context, reservationID string) error {
	// 1. Begin transaction
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to begin transaction for stock deduction")
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 2. Get reservation
	reservation, err := s.reservationRepo.GetByID(ctx, reservationID)
	if err != nil {
		return fmt.Errorf("failed to get reservation: %w", err)
	}

	// 3. Check if reservation can be confirmed
	if !reservation.CanBeConfirmed() {
		return domain.ErrCannotConfirmReservation
	}

	// 4. Lock variant and deduct stock
	variant, err := s.variantRepo.GetForUpdate(ctx, tx, reservation.VariantID.String())
	if err != nil {
		return fmt.Errorf("failed to lock variant: %w", err)
	}

	// Verify stock availability (paranoid check)
	if variant.StockQuantity < reservation.Quantity {
		s.logger.Error().
			Str("variant_id", variant.ID.String()).
			Int32("stock", variant.StockQuantity).
			Int32("requested", reservation.Quantity).
			Msg("insufficient stock for deduction (data inconsistency)")
		return fmt.Errorf("insufficient stock for deduction")
	}

	// Deduct stock_quantity (direct SQL query for atomicity)
	newStock := variant.StockQuantity - reservation.Quantity

	// Update variant stock within transaction
	query := `
		UPDATE product_variants
		SET stock_quantity = stock_quantity - $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`
	_, err = tx.ExecContext(ctx, query, reservation.Quantity, variant.ID)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to deduct stock quantity")
		return fmt.Errorf("failed to deduct stock: %w", err)
	}

	// 5. Update reservation status to 'confirmed'
	if err := s.reservationRepo.UpdateStatus(ctx, tx, reservationID, domain.StockReservationStatusConfirmed); err != nil {
		s.logger.Error().Err(err).Str("reservation_id", reservationID).Msg("failed to confirm reservation")
		return fmt.Errorf("failed to confirm reservation: %w", err)
	}

	// 6. Note: reserved_quantity is automatically decreased by DB trigger

	// 7. Commit transaction
	if err := tx.Commit(); err != nil {
		s.logger.Error().Err(err).Msg("failed to commit stock deduction transaction")
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info().
		Str("reservation_id", reservationID).
		Str("variant_id", variant.ID.String()).
		Int32("quantity", reservation.Quantity).
		Int32("new_stock", newStock).
		Msg("stock deducted successfully")

	return nil
}

// CleanupExpiredReservations marks expired reservations as expired
// This should be called periodically (e.g., via cron job)
func (s *VariantService) CleanupExpiredReservations(ctx context.Context) (int64, error) {
	count, err := s.reservationRepo.CleanupExpired(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to cleanup expired reservations")
		return 0, err
	}

	if count > 0 {
		s.logger.Info().Int64("count", count).Msg("expired reservations cleaned up")
	}

	return count, nil
}

// GetVariantWithReservations retrieves a variant with its active reservations
func (s *VariantService) GetVariantWithReservations(ctx context.Context, variantID string) (*domain.ProductVariantV2, []*domain.StockReservation, error) {
	variant, err := s.variantRepo.GetByID(ctx, variantID)
	if err != nil {
		return nil, nil, err
	}

	reservations, err := s.reservationRepo.GetActiveByVariant(ctx, variantID)
	if err != nil {
		s.logger.Warn().Err(err).Msg("failed to get active reservations")
		reservations = []*domain.StockReservation{}
	}

	return variant, reservations, nil
}

// CreateVariant creates a new product variant with auto-generated SKU
func (s *VariantService) CreateVariant(ctx context.Context, input *domain.CreateVariantInputV2, categoryCode string) (*domain.ProductVariantV2, error) {
	// Auto-generate SKU if not provided or validate if provided
	if input.SKU == "" {
		// Auto-generate SKU
		attrs := make([]VariantAttributeForSKU, len(input.Attributes))
		for i, attr := range input.Attributes {
			attrs[i] = VariantAttributeForSKU{
				Code:       fmt.Sprintf("attr_%d", attr.AttributeID), // Simplified
				ValueLabel: s.extractValueLabel(&attr),
			}
		}

		input.SKU = s.skuGenerator.GenerateSKU(input.ProductID.String(), categoryCode, attrs)
	} else {
		// Validate provided SKU
		if err := s.skuGenerator.ValidateSKU(input.SKU); err != nil {
			return nil, fmt.Errorf("invalid SKU: %w", err)
		}
	}

	// Create variant via repository
	return s.variantRepo.Create(ctx, input)
}

// ListByProduct retrieves all variants for a product
// Accepts product_id as UUID string. Returns error if UUID is invalid.
func (s *VariantService) ListByProduct(ctx context.Context, productID string) ([]*domain.ProductVariantV2, error) {
	productUUID, err := uuid.Parse(productID)
	if err != nil {
		s.logger.Error().
			Str("product_id", productID).
			Err(err).
			Msg("invalid product UUID format")
		return nil, fmt.Errorf("invalid product ID: %w", err)
	}

	filter := &domain.ListVariantsFilter{
		ProductID:         productUUID,
		IncludeAttributes: true,
	}

	variants, err := s.variantRepo.ListByProduct(ctx, filter)
	if err != nil {
		s.logger.Error().Err(err).Str("product_id", productID).Msg("failed to list variants")
		return nil, fmt.Errorf("failed to list variants: %w", err)
	}

	return variants, nil
}

// Helper: extract value label from CreateVariantAttributeValue
func (s *VariantService) extractValueLabel(attr *domain.CreateVariantAttributeValue) string {
	if attr.ValueText != nil {
		return *attr.ValueText
	}
	if attr.ValueNumber != nil {
		return fmt.Sprintf("%.0f", *attr.ValueNumber)
	}
	if attr.ValueBoolean != nil {
		if *attr.ValueBoolean {
			return "YES"
		}
		return "NO"
	}
	return "UNK"
}

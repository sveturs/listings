// Package service contains business logic tests for the listings microservice.
// This file tests VariantService stock operations.
package service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vondi-global/listings/internal/domain"
	"github.com/vondi-global/listings/internal/repository/postgres"
)

// TestReserveStock_Success tests successful stock reservation
func TestReserveStock_Success(t *testing.T) {
	// Setup
	ctx := context.Background()
	db, cleanup := setupTestDB(t)
	defer cleanup()

	logger := zerolog.Nop()
	variantRepo := setupVariantRepo(db, logger)
	reservationRepo := setupReservationRepo(db, logger)
	skuGen := NewSKUGenerator()

	service := NewVariantService(variantRepo, reservationRepo, skuGen, db, logger)

	// Create test variant with stock=10, reserved=0
	productID := uuid.New()
	variant := createTestVariant(t, ctx, variantRepo, db, productID, 10, 0)

	// Test: Reserve 5 units
	req := &ReserveStockRequest{
		VariantID:  variant.ID.String(),
		OrderID:    uuid.New().String(),
		Quantity:   5,
		TTLMinutes: 30,
	}

	resp, err := service.ReserveStock(ctx, req)

	// Assert
	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.NotEmpty(t, resp.ReservationID)
	assert.Equal(t, int32(5), resp.AvailableAfter)

	// Verify reservation was created
	reservation, err := reservationRepo.GetByID(ctx, resp.ReservationID)
	require.NoError(t, err)
	assert.Equal(t, int32(5), reservation.Quantity)
	assert.Equal(t, domain.StockReservationStatusActive, reservation.Status)

	// Verify variant reserved_quantity was updated (by trigger)
	variantAfter, err := variantRepo.GetByID(ctx, variant.ID.String())
	require.NoError(t, err)
	assert.Equal(t, int32(5), variantAfter.ReservedQuantity)
}

// TestReserveStock_InsufficientStock tests reservation with insufficient stock
func TestReserveStock_InsufficientStock(t *testing.T) {
	ctx := context.Background()
	db, cleanup := setupTestDB(t)
	defer cleanup()

	logger := zerolog.Nop()
	variantRepo := setupVariantRepo(db, logger)
	reservationRepo := setupReservationRepo(db, logger)
	skuGen := NewSKUGenerator()

	service := NewVariantService(variantRepo, reservationRepo, skuGen, db, logger)

	// Create test variant with stock=10, reserved=5
	productID := uuid.New()
	variant := createTestVariant(t, ctx, variantRepo, db, productID, 10, 5)

	// Test: Try to reserve 6 units (only 5 available)
	req := &ReserveStockRequest{
		VariantID:  variant.ID.String(),
		OrderID:    uuid.New().String(),
		Quantity:   6,
		TTLMinutes: 30,
	}

	resp, err := service.ReserveStock(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "insufficient stock")
}

// TestReleaseStock tests releasing a stock reservation
func TestReleaseStock(t *testing.T) {
	ctx := context.Background()
	db, cleanup := setupTestDB(t)
	defer cleanup()

	logger := zerolog.Nop()
	variantRepo := setupVariantRepo(db, logger)
	reservationRepo := setupReservationRepo(db, logger)
	skuGen := NewSKUGenerator()

	service := NewVariantService(variantRepo, reservationRepo, skuGen, db, logger)

	// Create test variant and reservation
	productID := uuid.New()
	variant := createTestVariant(t, ctx, variantRepo, db, productID, 10, 5)

	tx, err := db.BeginTxx(ctx, nil)
	require.NoError(t, err)
	defer tx.Rollback()

	reservation := &domain.StockReservation{
		ID:        uuid.New(),
		VariantID: variant.ID,
		OrderID:   uuid.New(),
		Quantity:  5,
		ExpiresAt: time.Now().Add(30 * time.Minute),
		Status:    domain.StockReservationStatusActive,
	}
	err = reservationRepo.Create(ctx, tx, reservation)
	require.NoError(t, err)

	err = tx.Commit()
	require.NoError(t, err)

	// Test: Release stock
	err = service.ReleaseStock(ctx, reservation.ID.String())

	// Assert
	require.NoError(t, err)

	// Verify reservation status changed to cancelled
	reservationAfter, err := reservationRepo.GetByID(ctx, reservation.ID.String())
	require.NoError(t, err)
	assert.Equal(t, domain.StockReservationStatusCancelled, reservationAfter.Status)

	// Verify variant reserved_quantity was decreased (by trigger)
	variantAfter, err := variantRepo.GetByID(ctx, variant.ID.String())
	require.NoError(t, err)
	assert.Equal(t, int32(0), variantAfter.ReservedQuantity)
}

// TestConfirmStockDeduction tests confirming a reservation and deducting stock
func TestConfirmStockDeduction(t *testing.T) {
	ctx := context.Background()
	db, cleanup := setupTestDB(t)
	defer cleanup()

	logger := zerolog.Nop()
	variantRepo := setupVariantRepo(db, logger)
	reservationRepo := setupReservationRepo(db, logger)
	skuGen := NewSKUGenerator()

	service := NewVariantService(variantRepo, reservationRepo, skuGen, db, logger)

	// Create test variant with stock=10, reserved=5
	productID := uuid.New()
	variant := createTestVariant(t, ctx, variantRepo, db, productID, 10, 5)

	tx, err := db.BeginTxx(ctx, nil)
	require.NoError(t, err)
	defer tx.Rollback()

	reservation := &domain.StockReservation{
		ID:        uuid.New(),
		VariantID: variant.ID,
		OrderID:   uuid.New(),
		Quantity:  5,
		ExpiresAt: time.Now().Add(30 * time.Minute),
		Status:    domain.StockReservationStatusActive,
	}
	err = reservationRepo.Create(ctx, tx, reservation)
	require.NoError(t, err)

	err = tx.Commit()
	require.NoError(t, err)

	// Test: Confirm stock deduction
	err = service.ConfirmStockDeduction(ctx, reservation.ID.String())

	// Assert
	require.NoError(t, err)

	// Verify stock was deducted
	variantAfter, err := variantRepo.GetByID(ctx, variant.ID.String())
	require.NoError(t, err)
	assert.Equal(t, int32(5), variantAfter.StockQuantity) // 10 - 5 = 5

	// Verify reservation status changed to confirmed
	reservationAfter, err := reservationRepo.GetByID(ctx, reservation.ID.String())
	require.NoError(t, err)
	assert.Equal(t, domain.StockReservationStatusConfirmed, reservationAfter.Status)
}

// Helper functions

func setupTestDB(t *testing.T) (*sqlx.DB, func()) {
	// TODO: Setup dockertest PostgreSQL container
	// For now, return mock - implement when dockertest is available
	t.Skip("Test DB setup not yet implemented - requires dockertest")
	return nil, func() {}
}

func setupVariantRepo(db *sqlx.DB, logger zerolog.Logger) *postgres.VariantRepository {
	return postgres.NewVariantRepository(db, logger)
}

func setupReservationRepo(db *sqlx.DB, logger zerolog.Logger) *postgres.StockReservationRepository {
	return postgres.NewStockReservationRepository(db, logger)
}

func createTestVariant(t *testing.T, ctx context.Context, repo *postgres.VariantRepository, db *sqlx.DB, productID uuid.UUID, stock int32, reserved int32) *domain.ProductVariantV2 {
	input := &domain.CreateVariantInputV2{
		ProductID:     productID,
		SKU:          "TEST-SKU-" + uuid.New().String()[:8],
		StockQuantity: stock,
		LowStockAlert: 5,
		IsDefault:     true,
		Position:      0,
		Attributes: []domain.CreateVariantAttributeValue{
			{
				AttributeID: 1,
				ValueText:   ptr("Test"),
			},
		},
	}

	variant, err := repo.Create(ctx, input)
	require.NoError(t, err)

	// Manually set reserved_quantity if needed
	if reserved > 0 {
		_, err := db.ExecContext(ctx,
			"UPDATE product_variants SET reserved_quantity = $1 WHERE id = $2",
			reserved, variant.ID)
		require.NoError(t, err)
		variant.ReservedQuantity = reserved
	}

	return variant
}

func ptr[T any](v T) *T {
	return &v
}

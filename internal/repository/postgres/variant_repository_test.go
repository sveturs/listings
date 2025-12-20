// Package postgres contains repository tests for PostgreSQL.
// This file tests VariantRepository operations.
package postgres

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vondi-global/listings/internal/domain"
)

// TestVariantRepository_Create tests creating a variant with attributes
func TestVariantRepository_Create(t *testing.T) {
	// Setup
	db, cleanup := setupTestDB(t)
	defer cleanup()

	logger := zerolog.Nop()
	repo := NewVariantRepository(db, logger)
	ctx := context.Background()

	// Create test product
	productID := createTestProduct_helper(t, ctx, db)

	// Test: Create variant with attributes
	input := &domain.CreateVariantInputV2{
		ProductID:     productID,
		SKU:           "TEST-SKU-001",
		StockQuantity: 100,
		LowStockAlert: 10,
		IsDefault:     true,
		Position:      1,
		Attributes: []domain.CreateVariantAttributeValue{
			{
				AttributeID: 1, // Size
				ValueText:   ptr("M"),
			},
			{
				AttributeID: 2, // Color
				ValueText:   ptr("Black"),
			},
		},
	}

	variant, err := repo.Create(ctx, input)

	// Assert
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, variant.ID)
	assert.Equal(t, productID, variant.ProductID)
	assert.Equal(t, "TEST-SKU-001", variant.SKU)
	assert.Equal(t, int32(100), variant.StockQuantity)
	assert.Equal(t, int32(0), variant.ReservedQuantity)
	assert.Equal(t, int32(10), variant.LowStockAlert)
	assert.True(t, variant.IsDefault)
	assert.Equal(t, domain.VariantStatusActive, variant.Status)

	// Verify attributes were created
	variantWithAttrs, err := repo.GetByID(ctx, variant.ID.String())
	require.NoError(t, err)
	assert.Len(t, variantWithAttrs.Attributes, 2)
}

// TestVariantRepository_GetByID tests retrieving a variant by ID
func TestVariantRepository_GetByID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	logger := zerolog.Nop()
	repo := NewVariantRepository(db, logger)
	ctx := context.Background()

	// Create test variant
	productID := createTestProduct_helper(t, ctx, db)
	variant := createTestVariantHelper(t, ctx, repo, productID)

	// Test: Get by ID
	retrieved, err := repo.GetByID(ctx, variant.ID.String())

	// Assert
	require.NoError(t, err)
	assert.Equal(t, variant.ID, retrieved.ID)
	assert.Equal(t, variant.SKU, retrieved.SKU)
	assert.Equal(t, variant.StockQuantity, retrieved.StockQuantity)
}

// TestVariantRepository_FindByAttributes tests finding variant by attribute combination
func TestVariantRepository_FindByAttributes(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	logger := zerolog.Nop()
	repo := NewVariantRepository(db, logger)
	ctx := context.Background()

	// Create test product
	productID := createTestProduct_helper(t, ctx, db)

	// Create 3 variants: M-Black, M-White, L-Black
	createVariantWithAttrs(t, ctx, repo, productID, "M", "Black")
	createVariantWithAttrs(t, ctx, repo, productID, "M", "White")
	variantLBlack := createVariantWithAttrs(t, ctx, repo, productID, "L", "Black")

	// Test: Find L-Black by attributes
	attrs := map[int32]interface{}{
		1: "L",     // Size
		2: "Black", // Color
	}

	filter := &domain.FindVariantByAttributesFilter{
		ProductID:  productID,
		Attributes: attrs,
	}

	found, err := repo.FindByAttributes(ctx, filter)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, variantLBlack.ID, found.ID)
}

// TestVariantRepository_GetForUpdate tests SELECT FOR UPDATE locking
func TestVariantRepository_GetForUpdate(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	logger := zerolog.Nop()
	repo := NewVariantRepository(db, logger)
	ctx := context.Background()

	// Create test variant
	productID := createTestProduct_helper(t, ctx, db)
	variant := createTestVariantHelper(t, ctx, repo, productID)

	// Begin transaction
	tx, err := db.BeginTxx(ctx, nil)
	require.NoError(t, err)
	defer tx.Rollback()

	// Test: GetForUpdate should lock the row
	locked, err := repo.GetForUpdate(ctx, tx, variant.ID.String())

	// Assert
	require.NoError(t, err)
	assert.Equal(t, variant.ID, locked.ID)

	// Note: Testing actual locking behavior requires concurrent transactions
	// which is complex in unit tests. Integration tests should cover this.
}

// TestVariantRepository_Update tests updating a variant
func TestVariantRepository_Update(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	logger := zerolog.Nop()
	repo := NewVariantRepository(db, logger)
	ctx := context.Background()

	// Create test variant
	productID := createTestProduct_helper(t, ctx, db)
	variant := createTestVariantHelper(t, ctx, repo, productID)

	// Test: Update stock quantity and status
	newStock := int32(50)
	newStatus := domain.VariantStatusOutOfStock

	input := &domain.UpdateVariantInputV2{
		StockQuantity: &newStock,
		Status:        &newStatus,
	}

	updated, err := repo.Update(ctx, variant.ID.String(), input)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, newStock, updated.StockQuantity)
	assert.Equal(t, newStatus, updated.Status)
}

// TestVariantRepository_ListByProduct tests listing variants for a product
func TestVariantRepository_ListByProduct(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	logger := zerolog.Nop()
	repo := NewVariantRepository(db, logger)
	ctx := context.Background()

	// Create test product
	productID := createTestProduct_helper(t, ctx, db)

	// Create 3 variants
	createTestVariantHelper(t, ctx, repo, productID)
	createTestVariantHelper(t, ctx, repo, productID)
	createTestVariantHelper(t, ctx, repo, productID)

	// Test: List all variants for product
	filter := &domain.ListVariantsFilter{
		ProductID:         productID,
		IncludeAttributes: false,
	}

	variants, err := repo.ListByProduct(ctx, filter)

	// Assert
	require.NoError(t, err)
	assert.Len(t, variants, 3)
}

// TestVariantRepository_Delete tests soft-deleting a variant
func TestVariantRepository_Delete(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	logger := zerolog.Nop()
	repo := NewVariantRepository(db, logger)
	ctx := context.Background()

	// Create test variant
	productID := createTestProduct_helper(t, ctx, db)
	variant := createTestVariantHelper(t, ctx, repo, productID)

	// Test: Delete variant
	err := repo.Delete(ctx, variant.ID.String())

	// Assert
	require.NoError(t, err)

	// Verify variant status changed to discontinued
	deleted, err := repo.GetByID(ctx, variant.ID.String())
	require.NoError(t, err)
	assert.Equal(t, domain.VariantStatusDiscontinued, deleted.Status)
}

// Helper functions

func setupTestDB(t *testing.T) (*sqlx.DB, func()) {
	// TODO: Setup dockertest PostgreSQL container with migrations
	t.Skip("Test DB setup not yet implemented - requires dockertest")
	return nil, func() {}
}

func createTestProduct_helper(t *testing.T, ctx context.Context, db *sqlx.DB) uuid.UUID {
	productID := uuid.New()

	// TODO: Insert test product into products table
	// For now, just return UUID
	// require.NoError(t, err)

	return productID
}

func createTestVariantHelper(t *testing.T, ctx context.Context, repo *VariantRepository, productID uuid.UUID) *domain.ProductVariantV2 {
	input := &domain.CreateVariantInputV2{
		ProductID:     productID,
		SKU:           "TEST-SKU-" + uuid.New().String()[:8],
		StockQuantity: 100,
		LowStockAlert: 10,
		IsDefault:     false,
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

	return variant
}

func createVariantWithAttrs(t *testing.T, ctx context.Context, repo *VariantRepository, productID uuid.UUID, size, color string) *domain.ProductVariantV2 {
	input := &domain.CreateVariantInputV2{
		ProductID:     productID,
		SKU:           "TEST-" + size + "-" + color,
		StockQuantity: 100,
		LowStockAlert: 10,
		IsDefault:     false,
		Position:      0,
		Attributes: []domain.CreateVariantAttributeValue{
			{
				AttributeID: 1, // Size
				ValueText:   &size,
			},
			{
				AttributeID: 2, // Color
				ValueText:   &color,
			},
		},
	}

	variant, err := repo.Create(ctx, input)
	require.NoError(t, err)

	return variant
}

func ptr[T any](v T) *T {
	return &v
}

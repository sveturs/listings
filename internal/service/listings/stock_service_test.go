package listings

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vondi-global/listings/internal/service/listings/mocks"
)

// ============================================================================
// TEST SETUP HELPERS
// ============================================================================

// setupStockTest creates a service with mocked dependencies for stock testing
func setupStockTest(t *testing.T) (*Service, *mocks.MockRepository, sqlmock.Sqlmock, *sql.DB) {
	t.Helper()

	mockRepo := new(mocks.MockRepository)
	mockCache := new(mocks.MockCacheRepository)
	mockIndexer := new(mocks.MockIndexingService)

	// Create sqlmock for transaction testing
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	logger := zerolog.New(zerolog.NewTestWriter(t)).With().Timestamp().Logger()
	service := NewService(mockRepo, mockCache, mockIndexer, logger)

	return service, mockRepo, mock, db
}

// ============================================================================
// DECREMENT STOCK TESTS
// ============================================================================

func TestDecrementStock_Success_SingleProduct(t *testing.T) {
	svc, mockRepo, mock, db := setupStockTest(t)
	defer db.Close()

	ctx := context.Background()
	orderID := "ORDER-001"

	items := []StockItem{
		{ProductID: 100, VariantID: nil, Quantity: 5},
	}

	// Mock transaction
	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT quantity FROM listings WHERE id = \$1 AND source_type = 'b2c' FOR UPDATE`).
		WithArgs(int64(100)).
		WillReturnRows(sqlmock.NewRows([]string{"stock_quantity"}).AddRow(10))
	mock.ExpectExec(`UPDATE listings SET quantity = quantity - \$1`).
		WithArgs(int32(5), int64(100)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	// Mock BeginTx
	mockRepo.On("BeginTx", ctx).Return(db.Begin())

	results, err := svc.DecrementStock(ctx, items, &orderID)

	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.True(t, results[0].Success)
	assert.Equal(t, int64(100), results[0].ProductID)
	assert.Nil(t, results[0].VariantID)
	assert.Equal(t, int32(10), results[0].StockBefore)
	assert.Equal(t, int32(5), results[0].StockAfter)
	assert.Nil(t, results[0].Error)

	assert.NoError(t, mock.ExpectationsWereMet())
	mockRepo.AssertExpectations(t)
}

func TestDecrementStock_Success_MultipleProducts(t *testing.T) {
	svc, mockRepo, mock, db := setupStockTest(t)
	defer db.Close()

	ctx := context.Background()
	orderID := "ORDER-002"

	items := []StockItem{
		{ProductID: 100, VariantID: nil, Quantity: 3},
		{ProductID: 200, VariantID: nil, Quantity: 2},
	}

	// Mock transaction
	mock.ExpectBegin()

	// First product
	mock.ExpectQuery(`SELECT quantity FROM listings WHERE id = \$1 AND source_type = 'b2c' FOR UPDATE`).
		WithArgs(int64(100)).
		WillReturnRows(sqlmock.NewRows([]string{"stock_quantity"}).AddRow(20))
	mock.ExpectExec(`UPDATE listings SET quantity = quantity - \$1`).
		WithArgs(int32(3), int64(100)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Second product
	mock.ExpectQuery(`SELECT quantity FROM listings WHERE id = \$1 AND source_type = 'b2c' FOR UPDATE`).
		WithArgs(int64(200)).
		WillReturnRows(sqlmock.NewRows([]string{"stock_quantity"}).AddRow(15))
	mock.ExpectExec(`UPDATE listings SET quantity = quantity - \$1`).
		WithArgs(int32(2), int64(200)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	mockRepo.On("BeginTx", ctx).Return(db.Begin())

	results, err := svc.DecrementStock(ctx, items, &orderID)

	require.NoError(t, err)
	require.Len(t, results, 2)
	assert.True(t, results[0].Success)
	assert.True(t, results[1].Success)
	assert.Equal(t, int64(100), results[0].ProductID)
	assert.Equal(t, int64(200), results[1].ProductID)

	assert.NoError(t, mock.ExpectationsWereMet())
	mockRepo.AssertExpectations(t)
}

func TestDecrementStock_Success_SingleVariant(t *testing.T) {
	svc, mockRepo, mock, db := setupStockTest(t)
	defer db.Close()

	ctx := context.Background()
	orderID := "ORDER-003"
	variantID := int64(500)

	items := []StockItem{
		{ProductID: 100, VariantID: &variantID, Quantity: 4},
	}

	// Mock transaction
	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT stock_quantity FROM b2c_product_variants WHERE id = \$1 AND product_id = \$2 FOR UPDATE`).
		WithArgs(variantID, int64(100)).
		WillReturnRows(sqlmock.NewRows([]string{"stock_quantity"}).AddRow(12))
	mock.ExpectExec(`UPDATE b2c_product_variants SET stock_quantity = stock_quantity - \$1`).
		WithArgs(int32(4), variantID, int64(100)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	mockRepo.On("BeginTx", ctx).Return(db.Begin())

	results, err := svc.DecrementStock(ctx, items, &orderID)

	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.True(t, results[0].Success)
	assert.Equal(t, int64(100), results[0].ProductID)
	assert.NotNil(t, results[0].VariantID)
	assert.Equal(t, variantID, *results[0].VariantID)
	assert.Equal(t, int32(12), results[0].StockBefore)
	assert.Equal(t, int32(8), results[0].StockAfter)

	assert.NoError(t, mock.ExpectationsWereMet())
	mockRepo.AssertExpectations(t)
}

func TestDecrementStock_Success_MultipleVariants(t *testing.T) {
	svc, mockRepo, mock, db := setupStockTest(t)
	defer db.Close()

	ctx := context.Background()
	orderID := "ORDER-004"
	variantID1 := int64(501)
	variantID2 := int64(502)

	items := []StockItem{
		{ProductID: 100, VariantID: &variantID1, Quantity: 2},
		{ProductID: 100, VariantID: &variantID2, Quantity: 3},
	}

	// Mock transaction
	mock.ExpectBegin()

	// First variant
	mock.ExpectQuery(`SELECT stock_quantity FROM b2c_product_variants WHERE id = \$1 AND product_id = \$2 FOR UPDATE`).
		WithArgs(variantID1, int64(100)).
		WillReturnRows(sqlmock.NewRows([]string{"stock_quantity"}).AddRow(10))
	mock.ExpectExec(`UPDATE b2c_product_variants SET stock_quantity = stock_quantity - \$1`).
		WithArgs(int32(2), variantID1, int64(100)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Second variant
	mock.ExpectQuery(`SELECT stock_quantity FROM b2c_product_variants WHERE id = \$1 AND product_id = \$2 FOR UPDATE`).
		WithArgs(variantID2, int64(100)).
		WillReturnRows(sqlmock.NewRows([]string{"stock_quantity"}).AddRow(15))
	mock.ExpectExec(`UPDATE b2c_product_variants SET stock_quantity = stock_quantity - \$1`).
		WithArgs(int32(3), variantID2, int64(100)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	mockRepo.On("BeginTx", ctx).Return(db.Begin())

	results, err := svc.DecrementStock(ctx, items, &orderID)

	require.NoError(t, err)
	require.Len(t, results, 2)
	assert.True(t, results[0].Success)
	assert.True(t, results[1].Success)

	assert.NoError(t, mock.ExpectationsWereMet())
	mockRepo.AssertExpectations(t)
}

func TestDecrementStock_Success_MixedProductsAndVariants(t *testing.T) {
	svc, mockRepo, mock, db := setupStockTest(t)
	defer db.Close()

	ctx := context.Background()
	orderID := "ORDER-005"
	variantID := int64(500)

	items := []StockItem{
		{ProductID: 100, VariantID: nil, Quantity: 2},
		{ProductID: 200, VariantID: &variantID, Quantity: 3},
	}

	// Mock transaction
	mock.ExpectBegin()

	// Product
	mock.ExpectQuery(`SELECT quantity FROM listings WHERE id = \$1 AND source_type = 'b2c' FOR UPDATE`).
		WithArgs(int64(100)).
		WillReturnRows(sqlmock.NewRows([]string{"stock_quantity"}).AddRow(10))
	mock.ExpectExec(`UPDATE listings SET quantity = quantity - \$1`).
		WithArgs(int32(2), int64(100)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Variant
	mock.ExpectQuery(`SELECT stock_quantity FROM b2c_product_variants WHERE id = \$1 AND product_id = \$2 FOR UPDATE`).
		WithArgs(variantID, int64(200)).
		WillReturnRows(sqlmock.NewRows([]string{"stock_quantity"}).AddRow(12))
	mock.ExpectExec(`UPDATE b2c_product_variants SET stock_quantity = stock_quantity - \$1`).
		WithArgs(int32(3), variantID, int64(200)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	mockRepo.On("BeginTx", ctx).Return(db.Begin())

	results, err := svc.DecrementStock(ctx, items, &orderID)

	require.NoError(t, err)
	require.Len(t, results, 2)
	assert.True(t, results[0].Success)
	assert.True(t, results[1].Success)

	assert.NoError(t, mock.ExpectationsWereMet())
	mockRepo.AssertExpectations(t)
}

func TestDecrementStock_Error_InsufficientStock_Product(t *testing.T) {
	svc, mockRepo, mock, db := setupStockTest(t)
	defer db.Close()

	ctx := context.Background()
	orderID := "ORDER-006"

	items := []StockItem{
		{ProductID: 100, VariantID: nil, Quantity: 15},
	}

	// Mock transaction
	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT quantity FROM listings WHERE id = \$1 AND source_type = 'b2c' FOR UPDATE`).
		WithArgs(int64(100)).
		WillReturnRows(sqlmock.NewRows([]string{"stock_quantity"}).AddRow(10))
	mock.ExpectRollback()

	mockRepo.On("BeginTx", ctx).Return(db.Begin())

	results, err := svc.DecrementStock(ctx, items, &orderID)

	require.Error(t, err)
	require.Len(t, results, 1)
	assert.False(t, results[0].Success)
	assert.NotNil(t, results[0].Error)
	assert.Contains(t, *results[0].Error, "insufficient stock")

	assert.NoError(t, mock.ExpectationsWereMet())
	mockRepo.AssertExpectations(t)
}

func TestDecrementStock_Error_InsufficientStock_Variant(t *testing.T) {
	svc, mockRepo, mock, db := setupStockTest(t)
	defer db.Close()

	ctx := context.Background()
	orderID := "ORDER-007"
	variantID := int64(500)

	items := []StockItem{
		{ProductID: 100, VariantID: &variantID, Quantity: 20},
	}

	// Mock transaction
	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT stock_quantity FROM b2c_product_variants WHERE id = \$1 AND product_id = \$2 FOR UPDATE`).
		WithArgs(variantID, int64(100)).
		WillReturnRows(sqlmock.NewRows([]string{"stock_quantity"}).AddRow(8))
	mock.ExpectRollback()

	mockRepo.On("BeginTx", ctx).Return(db.Begin())

	results, err := svc.DecrementStock(ctx, items, &orderID)

	require.Error(t, err)
	require.Len(t, results, 1)
	assert.False(t, results[0].Success)
	assert.NotNil(t, results[0].Error)
	assert.Contains(t, *results[0].Error, "insufficient stock")

	assert.NoError(t, mock.ExpectationsWereMet())
	mockRepo.AssertExpectations(t)
}

func TestDecrementStock_Error_NotFound_Product(t *testing.T) {
	svc, mockRepo, mock, db := setupStockTest(t)
	defer db.Close()

	ctx := context.Background()
	orderID := "ORDER-008"

	items := []StockItem{
		{ProductID: 999, VariantID: nil, Quantity: 5},
	}

	// Mock transaction
	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT quantity FROM listings WHERE id = \$1 AND source_type = 'b2c' FOR UPDATE`).
		WithArgs(int64(999)).
		WillReturnError(sql.ErrNoRows)
	mock.ExpectRollback()

	mockRepo.On("BeginTx", ctx).Return(db.Begin())

	results, err := svc.DecrementStock(ctx, items, &orderID)

	require.Error(t, err)
	require.Len(t, results, 1)
	assert.False(t, results[0].Success)
	assert.NotNil(t, results[0].Error)
	assert.Contains(t, *results[0].Error, "not found")

	assert.NoError(t, mock.ExpectationsWereMet())
	mockRepo.AssertExpectations(t)
}

func TestDecrementStock_Error_NotFound_Variant(t *testing.T) {
	svc, mockRepo, mock, db := setupStockTest(t)
	defer db.Close()

	ctx := context.Background()
	orderID := "ORDER-009"
	variantID := int64(999)

	items := []StockItem{
		{ProductID: 100, VariantID: &variantID, Quantity: 5},
	}

	// Mock transaction
	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT stock_quantity FROM b2c_product_variants WHERE id = \$1 AND product_id = \$2 FOR UPDATE`).
		WithArgs(variantID, int64(100)).
		WillReturnError(sql.ErrNoRows)
	mock.ExpectRollback()

	mockRepo.On("BeginTx", ctx).Return(db.Begin())

	results, err := svc.DecrementStock(ctx, items, &orderID)

	require.Error(t, err)
	require.Len(t, results, 1)
	assert.False(t, results[0].Success)
	assert.NotNil(t, results[0].Error)
	assert.Contains(t, *results[0].Error, "not found")

	assert.NoError(t, mock.ExpectationsWereMet())
	mockRepo.AssertExpectations(t)
}

func TestDecrementStock_Error_TransactionRollback(t *testing.T) {
	svc, mockRepo, mock, db := setupStockTest(t)
	defer db.Close()

	ctx := context.Background()
	orderID := "ORDER-010"

	items := []StockItem{
		{ProductID: 100, VariantID: nil, Quantity: 5},
	}

	// Mock transaction that fails to commit
	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT quantity FROM listings WHERE id = \$1 AND source_type = 'b2c' FOR UPDATE`).
		WithArgs(int64(100)).
		WillReturnRows(sqlmock.NewRows([]string{"stock_quantity"}).AddRow(10))
	mock.ExpectExec(`UPDATE listings SET quantity = quantity - \$1`).
		WithArgs(int32(5), int64(100)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit().WillReturnError(errors.New("commit failed"))

	mockRepo.On("BeginTx", ctx).Return(db.Begin())

	results, err := svc.DecrementStock(ctx, items, &orderID)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to commit transaction")
	assert.Nil(t, results)

	assert.NoError(t, mock.ExpectationsWereMet())
	mockRepo.AssertExpectations(t)
}

func TestDecrementStock_Concurrency_StockChangedDuringTransaction_Product(t *testing.T) {
	svc, mockRepo, mock, db := setupStockTest(t)
	defer db.Close()

	ctx := context.Background()
	orderID := "ORDER-011"

	items := []StockItem{
		{ProductID: 100, VariantID: nil, Quantity: 5},
	}

	// Mock transaction - stock sufficient but update affects 0 rows (race condition)
	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT quantity FROM listings WHERE id = \$1 AND source_type = 'b2c' FOR UPDATE`).
		WithArgs(int64(100)).
		WillReturnRows(sqlmock.NewRows([]string{"stock_quantity"}).AddRow(10))
	mock.ExpectExec(`UPDATE listings SET quantity = quantity - \$1`).
		WithArgs(int32(5), int64(100)).
		WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected
	mock.ExpectRollback()

	mockRepo.On("BeginTx", ctx).Return(db.Begin())

	results, err := svc.DecrementStock(ctx, items, &orderID)

	require.Error(t, err)
	require.Len(t, results, 1)
	assert.False(t, results[0].Success)
	assert.NotNil(t, results[0].Error)
	assert.Contains(t, *results[0].Error, "stock changed during transaction")

	assert.NoError(t, mock.ExpectationsWereMet())
	mockRepo.AssertExpectations(t)
}

func TestDecrementStock_Concurrency_StockChangedDuringTransaction_Variant(t *testing.T) {
	svc, mockRepo, mock, db := setupStockTest(t)
	defer db.Close()

	ctx := context.Background()
	orderID := "ORDER-012"
	variantID := int64(500)

	items := []StockItem{
		{ProductID: 100, VariantID: &variantID, Quantity: 3},
	}

	// Mock transaction - stock sufficient but update affects 0 rows (race condition)
	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT stock_quantity FROM b2c_product_variants WHERE id = \$1 AND product_id = \$2 FOR UPDATE`).
		WithArgs(variantID, int64(100)).
		WillReturnRows(sqlmock.NewRows([]string{"stock_quantity"}).AddRow(10))
	mock.ExpectExec(`UPDATE b2c_product_variants SET stock_quantity = stock_quantity - \$1`).
		WithArgs(int32(3), variantID, int64(100)).
		WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected
	mock.ExpectRollback()

	mockRepo.On("BeginTx", ctx).Return(db.Begin())

	results, err := svc.DecrementStock(ctx, items, &orderID)

	require.Error(t, err)
	require.Len(t, results, 1)
	assert.False(t, results[0].Success)
	assert.NotNil(t, results[0].Error)
	assert.Contains(t, *results[0].Error, "stock changed during transaction")

	assert.NoError(t, mock.ExpectationsWereMet())
	mockRepo.AssertExpectations(t)
}

func TestDecrementStock_Concurrency_RaceCondition(t *testing.T) {
	svc, mockRepo, _, _ := setupStockTest(t)

	ctx := context.Background()

	// Mock BeginTx to return error simulating concurrent access
	mockRepo.On("BeginTx", ctx).Return(nil, errors.New("could not serialize access"))

	items := []StockItem{
		{ProductID: 100, VariantID: nil, Quantity: 5},
	}
	orderID := "ORDER-013"

	results, err := svc.DecrementStock(ctx, items, &orderID)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to begin transaction")
	assert.Nil(t, results)

	mockRepo.AssertExpectations(t)
}

// ============================================================================
// ROLLBACK STOCK TESTS
// ============================================================================

func TestRollbackStock_Success_SingleProduct(t *testing.T) {
	svc, mockRepo, mock, db := setupStockTest(t)
	defer db.Close()

	ctx := context.Background()
	orderID := "ORDER-101"

	items := []StockItem{
		{ProductID: 100, VariantID: nil, Quantity: 5},
	}

	// Mock transaction
	mock.ExpectBegin()

	// Check rollback doesn't exist
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM inventory_movements WHERE order_id = \$1 AND product_id = \$2 AND variant_id IS NULL AND movement_type = 'rollback'`).
		WithArgs(orderID, int64(100)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	// Get current stock
	mock.ExpectQuery(`SELECT quantity FROM listings WHERE id = \$1 AND source_type = 'b2c'`).
		WithArgs(int64(100)).
		WillReturnRows(sqlmock.NewRows([]string{"stock_quantity"}).AddRow(5))

	// Update stock
	mock.ExpectExec(`UPDATE listings SET quantity = quantity \+ \$1`).
		WithArgs(int32(5), int64(100)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Record rollback
	mock.ExpectExec(`INSERT INTO inventory_movements`).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	mockRepo.On("BeginTx", ctx).Return(db.Begin())

	results, err := svc.RollbackStock(ctx, items, &orderID)

	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.True(t, results[0].Success)
	assert.Equal(t, int64(100), results[0].ProductID)
	assert.Equal(t, int32(5), results[0].StockBefore)
	assert.Equal(t, int32(10), results[0].StockAfter)

	assert.NoError(t, mock.ExpectationsWereMet())
	mockRepo.AssertExpectations(t)
}

func TestRollbackStock_Success_MultipleProducts(t *testing.T) {
	svc, mockRepo, mock, db := setupStockTest(t)
	defer db.Close()

	ctx := context.Background()
	orderID := "ORDER-102"

	items := []StockItem{
		{ProductID: 100, VariantID: nil, Quantity: 3},
		{ProductID: 200, VariantID: nil, Quantity: 2},
	}

	// Mock transaction
	mock.ExpectBegin()

	// First product
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM inventory_movements WHERE order_id = \$1 AND product_id = \$2 AND variant_id IS NULL AND movement_type = 'rollback'`).
		WithArgs(orderID, int64(100)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectQuery(`SELECT quantity FROM listings WHERE id = \$1 AND source_type = 'b2c'`).
		WithArgs(int64(100)).
		WillReturnRows(sqlmock.NewRows([]string{"stock_quantity"}).AddRow(17))
	mock.ExpectExec(`UPDATE listings SET quantity = quantity \+ \$1`).
		WithArgs(int32(3), int64(100)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`INSERT INTO inventory_movements`).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Second product
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM inventory_movements WHERE order_id = \$1 AND product_id = \$2 AND variant_id IS NULL AND movement_type = 'rollback'`).
		WithArgs(orderID, int64(200)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectQuery(`SELECT quantity FROM listings WHERE id = \$1 AND source_type = 'b2c'`).
		WithArgs(int64(200)).
		WillReturnRows(sqlmock.NewRows([]string{"stock_quantity"}).AddRow(13))
	mock.ExpectExec(`UPDATE listings SET quantity = quantity \+ \$1`).
		WithArgs(int32(2), int64(200)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`INSERT INTO inventory_movements`).
		WillReturnResult(sqlmock.NewResult(2, 1))

	mock.ExpectCommit()

	mockRepo.On("BeginTx", ctx).Return(db.Begin())

	results, err := svc.RollbackStock(ctx, items, &orderID)

	require.NoError(t, err)
	require.Len(t, results, 2)
	assert.True(t, results[0].Success)
	assert.True(t, results[1].Success)

	assert.NoError(t, mock.ExpectationsWereMet())
	mockRepo.AssertExpectations(t)
}

func TestRollbackStock_Success_SingleVariant(t *testing.T) {
	svc, mockRepo, mock, db := setupStockTest(t)
	defer db.Close()

	ctx := context.Background()
	orderID := "ORDER-103"
	variantID := int64(500)

	items := []StockItem{
		{ProductID: 100, VariantID: &variantID, Quantity: 4},
	}

	// Mock transaction
	mock.ExpectBegin()

	// Check rollback doesn't exist
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM inventory_movements WHERE order_id = \$1 AND variant_id = \$2 AND movement_type = 'rollback'`).
		WithArgs(orderID, variantID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	// Get current stock
	mock.ExpectQuery(`SELECT stock_quantity FROM b2c_product_variants WHERE id = \$1 AND product_id = \$2`).
		WithArgs(variantID, int64(100)).
		WillReturnRows(sqlmock.NewRows([]string{"stock_quantity"}).AddRow(8))

	// Update stock
	mock.ExpectExec(`UPDATE b2c_product_variants SET stock_quantity = stock_quantity \+ \$1, updated_at = NOW\(\) WHERE id = \$2 AND product_id = \$3`).
		WithArgs(int32(4), variantID, int64(100)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Record rollback
	mock.ExpectExec(`INSERT INTO inventory_movements`).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	mockRepo.On("BeginTx", ctx).Return(db.Begin())

	results, err := svc.RollbackStock(ctx, items, &orderID)

	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.True(t, results[0].Success)
	assert.Equal(t, int64(100), results[0].ProductID)
	assert.NotNil(t, results[0].VariantID)
	assert.Equal(t, variantID, *results[0].VariantID)
	assert.Equal(t, int32(8), results[0].StockBefore)
	assert.Equal(t, int32(12), results[0].StockAfter)

	assert.NoError(t, mock.ExpectationsWereMet())
	mockRepo.AssertExpectations(t)
}

func TestRollbackStock_Success_MultipleVariants(t *testing.T) {
	svc, mockRepo, mock, db := setupStockTest(t)
	defer db.Close()

	ctx := context.Background()
	orderID := "ORDER-104"
	variantID1 := int64(501)
	variantID2 := int64(502)

	items := []StockItem{
		{ProductID: 100, VariantID: &variantID1, Quantity: 2},
		{ProductID: 100, VariantID: &variantID2, Quantity: 3},
	}

	// Mock transaction
	mock.ExpectBegin()

	// First variant
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM inventory_movements WHERE order_id = \$1 AND variant_id = \$2 AND movement_type = 'rollback'`).
		WithArgs(orderID, variantID1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectQuery(`SELECT stock_quantity FROM b2c_product_variants WHERE id = \$1 AND product_id = \$2`).
		WithArgs(variantID1, int64(100)).
		WillReturnRows(sqlmock.NewRows([]string{"stock_quantity"}).AddRow(8))
	mock.ExpectExec(`UPDATE b2c_product_variants SET stock_quantity = stock_quantity \+ \$1, updated_at = NOW\(\) WHERE id = \$2 AND product_id = \$3`).
		WithArgs(int32(2), variantID1, int64(100)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`INSERT INTO inventory_movements`).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Second variant
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM inventory_movements WHERE order_id = \$1 AND variant_id = \$2 AND movement_type = 'rollback'`).
		WithArgs(orderID, variantID2).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectQuery(`SELECT stock_quantity FROM b2c_product_variants WHERE id = \$1 AND product_id = \$2`).
		WithArgs(variantID2, int64(100)).
		WillReturnRows(sqlmock.NewRows([]string{"stock_quantity"}).AddRow(12))
	mock.ExpectExec(`UPDATE b2c_product_variants SET stock_quantity = stock_quantity \+ \$1, updated_at = NOW\(\) WHERE id = \$2 AND product_id = \$3`).
		WithArgs(int32(3), variantID2, int64(100)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`INSERT INTO inventory_movements`).
		WillReturnResult(sqlmock.NewResult(2, 1))

	mock.ExpectCommit()

	mockRepo.On("BeginTx", ctx).Return(db.Begin())

	results, err := svc.RollbackStock(ctx, items, &orderID)

	require.NoError(t, err)
	require.Len(t, results, 2)
	assert.True(t, results[0].Success)
	assert.True(t, results[1].Success)

	assert.NoError(t, mock.ExpectationsWereMet())
	mockRepo.AssertExpectations(t)
}

func TestRollbackStock_Idempotent_MultipleCallsSameOrder_Product(t *testing.T) {
	svc, mockRepo, mock, db := setupStockTest(t)
	defer db.Close()

	ctx := context.Background()
	orderID := "ORDER-105"

	items := []StockItem{
		{ProductID: 100, VariantID: nil, Quantity: 5},
	}

	// Mock transaction - rollback already exists
	mock.ExpectBegin()

	// Check rollback exists (idempotency)
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM inventory_movements WHERE order_id = \$1 AND product_id = \$2 AND variant_id IS NULL AND movement_type = 'rollback'`).
		WithArgs(orderID, int64(100)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	// Get current stock (no update)
	mock.ExpectQuery(`SELECT quantity FROM listings WHERE id = \$1 AND source_type = 'b2c'`).
		WithArgs(int64(100)).
		WillReturnRows(sqlmock.NewRows([]string{"stock_quantity"}).AddRow(10))

	mock.ExpectCommit()

	mockRepo.On("BeginTx", ctx).Return(db.Begin())

	results, err := svc.RollbackStock(ctx, items, &orderID)

	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.True(t, results[0].Success)
	assert.Equal(t, int32(10), results[0].StockBefore)
	assert.Equal(t, int32(10), results[0].StockAfter) // No change (idempotent)

	assert.NoError(t, mock.ExpectationsWereMet())
	mockRepo.AssertExpectations(t)
}

func TestRollbackStock_Idempotent_MultipleCallsSameOrder_Variant(t *testing.T) {
	svc, mockRepo, mock, db := setupStockTest(t)
	defer db.Close()

	ctx := context.Background()
	orderID := "ORDER-106"
	variantID := int64(500)

	items := []StockItem{
		{ProductID: 100, VariantID: &variantID, Quantity: 3},
	}

	// Mock transaction - rollback already exists
	mock.ExpectBegin()

	// Check rollback exists (idempotency)
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM inventory_movements WHERE order_id = \$1 AND variant_id = \$2 AND movement_type = 'rollback'`).
		WithArgs(orderID, variantID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	// Get current stock (no update)
	mock.ExpectQuery(`SELECT stock_quantity FROM b2c_product_variants WHERE id = \$1 AND product_id = \$2`).
		WithArgs(variantID, int64(100)).
		WillReturnRows(sqlmock.NewRows([]string{"stock_quantity"}).AddRow(12))

	mock.ExpectCommit()

	mockRepo.On("BeginTx", ctx).Return(db.Begin())

	results, err := svc.RollbackStock(ctx, items, &orderID)

	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.True(t, results[0].Success)
	assert.Equal(t, int32(12), results[0].StockBefore)
	assert.Equal(t, int32(12), results[0].StockAfter) // No change (idempotent)

	assert.NoError(t, mock.ExpectationsWereMet())
	mockRepo.AssertExpectations(t)
}

func TestRollbackStock_Error_MissingOrderID(t *testing.T) {
	svc, mockRepo, _, _ := setupStockTest(t)

	ctx := context.Background()

	items := []StockItem{
		{ProductID: 100, VariantID: nil, Quantity: 5},
	}

	results, err := svc.RollbackStock(ctx, items, nil)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "order_id is required")
	assert.Nil(t, results)

	mockRepo.AssertExpectations(t)
}

func TestRollbackStock_Error_EmptyOrderID(t *testing.T) {
	svc, mockRepo, _, _ := setupStockTest(t)

	ctx := context.Background()
	emptyOrderID := ""

	items := []StockItem{
		{ProductID: 100, VariantID: nil, Quantity: 5},
	}

	results, err := svc.RollbackStock(ctx, items, &emptyOrderID)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "order_id is required")
	assert.Nil(t, results)

	mockRepo.AssertExpectations(t)
}

func TestRollbackStock_Error_NotFound_Product(t *testing.T) {
	svc, mockRepo, mock, db := setupStockTest(t)
	defer db.Close()

	ctx := context.Background()
	orderID := "ORDER-107"

	items := []StockItem{
		{ProductID: 999, VariantID: nil, Quantity: 5},
	}

	// Mock transaction
	mock.ExpectBegin()

	// Check rollback doesn't exist
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM inventory_movements WHERE order_id = \$1 AND product_id = \$2 AND variant_id IS NULL AND movement_type = 'rollback'`).
		WithArgs(orderID, int64(999)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	// Product not found
	mock.ExpectQuery(`SELECT quantity FROM listings WHERE id = \$1 AND source_type = 'b2c'`).
		WithArgs(int64(999)).
		WillReturnError(sql.ErrNoRows)

	mock.ExpectCommit()

	mockRepo.On("BeginTx", ctx).Return(db.Begin())

	results, err := svc.RollbackStock(ctx, items, &orderID)

	require.NoError(t, err) // RollbackStock continues on errors
	require.Len(t, results, 1)
	assert.False(t, results[0].Success)
	assert.NotNil(t, results[0].Error)
	assert.Contains(t, *results[0].Error, "not found")

	assert.NoError(t, mock.ExpectationsWereMet())
	mockRepo.AssertExpectations(t)
}

func TestRollbackStock_Error_NotFound_Variant(t *testing.T) {
	svc, mockRepo, mock, db := setupStockTest(t)
	defer db.Close()

	ctx := context.Background()
	orderID := "ORDER-108"
	variantID := int64(999)

	items := []StockItem{
		{ProductID: 100, VariantID: &variantID, Quantity: 3},
	}

	// Mock transaction
	mock.ExpectBegin()

	// Check rollback doesn't exist
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM inventory_movements WHERE order_id = \$1 AND variant_id = \$2 AND movement_type = 'rollback'`).
		WithArgs(orderID, variantID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	// Variant not found
	mock.ExpectQuery(`SELECT stock_quantity FROM b2c_product_variants WHERE id = \$1 AND product_id = \$2`).
		WithArgs(variantID, int64(100)).
		WillReturnError(sql.ErrNoRows)

	mock.ExpectCommit()

	mockRepo.On("BeginTx", ctx).Return(db.Begin())

	results, err := svc.RollbackStock(ctx, items, &orderID)

	require.NoError(t, err) // RollbackStock continues on errors
	require.Len(t, results, 1)
	assert.False(t, results[0].Success)
	assert.NotNil(t, results[0].Error)
	assert.Contains(t, *results[0].Error, "not found")

	assert.NoError(t, mock.ExpectationsWereMet())
	mockRepo.AssertExpectations(t)
}

func TestRollbackStock_Audit_RecordCreated(t *testing.T) {
	svc, mockRepo, mock, db := setupStockTest(t)
	defer db.Close()

	ctx := context.Background()
	orderID := "ORDER-109"

	items := []StockItem{
		{ProductID: 100, VariantID: nil, Quantity: 5},
	}

	// Mock transaction
	mock.ExpectBegin()

	// Check rollback doesn't exist
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM inventory_movements WHERE order_id = \$1 AND product_id = \$2 AND variant_id IS NULL AND movement_type = 'rollback'`).
		WithArgs(orderID, int64(100)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	// Get current stock
	mock.ExpectQuery(`SELECT quantity FROM listings WHERE id = \$1 AND source_type = 'b2c'`).
		WithArgs(int64(100)).
		WillReturnRows(sqlmock.NewRows([]string{"stock_quantity"}).AddRow(5))

	// Update stock
	mock.ExpectExec(`UPDATE listings SET quantity = quantity \+ \$1`).
		WithArgs(int32(5), int64(100)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// CRITICAL: Record rollback in audit table
	mock.ExpectExec(`INSERT INTO inventory_movements`).
		WithArgs(int64(100), nil, "rollback", int32(5), "rollback", sqlmock.AnyArg(), int64(0), orderID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	mockRepo.On("BeginTx", ctx).Return(db.Begin())

	results, err := svc.RollbackStock(ctx, items, &orderID)

	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.True(t, results[0].Success)

	assert.NoError(t, mock.ExpectationsWereMet())
	mockRepo.AssertExpectations(t)
}

func TestRollbackStock_Audit_PreventsDuplicateRecord(t *testing.T) {
	svc, mockRepo, mock, db := setupStockTest(t)
	defer db.Close()

	ctx := context.Background()
	orderID := "ORDER-110"

	items := []StockItem{
		{ProductID: 100, VariantID: nil, Quantity: 5},
	}

	// Mock transaction
	mock.ExpectBegin()

	// Rollback already recorded (count > 0)
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM inventory_movements WHERE order_id = \$1 AND product_id = \$2 AND variant_id IS NULL AND movement_type = 'rollback'`).
		WithArgs(orderID, int64(100)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	// Should skip update and audit insert
	mock.ExpectQuery(`SELECT quantity FROM listings WHERE id = \$1 AND source_type = 'b2c'`).
		WithArgs(int64(100)).
		WillReturnRows(sqlmock.NewRows([]string{"stock_quantity"}).AddRow(10))

	// NO update or insert expected
	mock.ExpectCommit()

	mockRepo.On("BeginTx", ctx).Return(db.Begin())

	results, err := svc.RollbackStock(ctx, items, &orderID)

	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.True(t, results[0].Success)

	// Verify no duplicate insert attempted
	assert.NoError(t, mock.ExpectationsWereMet())
	mockRepo.AssertExpectations(t)
}

// ============================================================================
// CHECK STOCK AVAILABILITY TESTS
// ============================================================================

func TestCheckStockAvailability_Available_Product(t *testing.T) {
	svc, mockRepo, mock, db := setupStockTest(t)
	defer db.Close()

	ctx := context.Background()

	items := []StockItem{
		{ProductID: 100, VariantID: nil, Quantity: 5},
	}

	// Mock stock query
	mock.ExpectQuery(`SELECT quantity FROM listings WHERE id = \$1 AND source_type = 'b2c'`).
		WithArgs(int64(100)).
		WillReturnRows(sqlmock.NewRows([]string{"stock_quantity"}).AddRow(10))

	// Create sqlx.DB wrapper for GetDB mock
	sqlxDB := sqlx.NewDb(db, "postgres")
	mockRepo.On("GetDB").Return(sqlxDB)

	allAvailable, availabilities, err := svc.CheckStockAvailability(ctx, items)

	require.NoError(t, err)
	assert.True(t, allAvailable)
	require.Len(t, availabilities, 1)
	assert.True(t, availabilities[0].IsAvailable)
	assert.Equal(t, int32(10), availabilities[0].AvailableQuantity)
	assert.Equal(t, int32(5), availabilities[0].RequestedQuantity)

	assert.NoError(t, mock.ExpectationsWereMet())
	mockRepo.AssertExpectations(t)
}

func TestCheckStockAvailability_Available_Variant(t *testing.T) {
	svc, mockRepo, mock, db := setupStockTest(t)
	defer db.Close()

	ctx := context.Background()
	variantID := int64(500)

	items := []StockItem{
		{ProductID: 100, VariantID: &variantID, Quantity: 3},
	}

	// Mock variant stock query
	mock.ExpectQuery(`SELECT stock_quantity FROM b2c_product_variants WHERE id = \$1 AND product_id = \$2`).
		WithArgs(variantID, int64(100)).
		WillReturnRows(sqlmock.NewRows([]string{"stock_quantity"}).AddRow(12))

	sqlxDB := sqlx.NewDb(db, "postgres")
	mockRepo.On("GetDB").Return(sqlxDB)

	allAvailable, availabilities, err := svc.CheckStockAvailability(ctx, items)

	require.NoError(t, err)
	assert.True(t, allAvailable)
	require.Len(t, availabilities, 1)
	assert.True(t, availabilities[0].IsAvailable)
	assert.Equal(t, int32(12), availabilities[0].AvailableQuantity)
	assert.Equal(t, int32(3), availabilities[0].RequestedQuantity)

	assert.NoError(t, mock.ExpectationsWereMet())
	mockRepo.AssertExpectations(t)
}

func TestCheckStockAvailability_Unavailable_InsufficientStock(t *testing.T) {
	svc, mockRepo, mock, db := setupStockTest(t)
	defer db.Close()

	ctx := context.Background()

	items := []StockItem{
		{ProductID: 100, VariantID: nil, Quantity: 20},
	}

	// Mock insufficient stock
	mock.ExpectQuery(`SELECT quantity FROM listings WHERE id = \$1 AND source_type = 'b2c'`).
		WithArgs(int64(100)).
		WillReturnRows(sqlmock.NewRows([]string{"stock_quantity"}).AddRow(10))

	sqlxDB := sqlx.NewDb(db, "postgres")
	mockRepo.On("GetDB").Return(sqlxDB)

	allAvailable, availabilities, err := svc.CheckStockAvailability(ctx, items)

	require.NoError(t, err)
	assert.False(t, allAvailable)
	require.Len(t, availabilities, 1)
	assert.False(t, availabilities[0].IsAvailable)
	assert.Equal(t, int32(20), availabilities[0].RequestedQuantity)
	assert.Equal(t, int32(10), availabilities[0].AvailableQuantity)

	assert.NoError(t, mock.ExpectationsWereMet())
	mockRepo.AssertExpectations(t)
}

func TestCheckStockAvailability_Unavailable_NotFound_Product(t *testing.T) {
	svc, mockRepo, mock, db := setupStockTest(t)
	defer db.Close()

	ctx := context.Background()

	items := []StockItem{
		{ProductID: 999, VariantID: nil, Quantity: 5},
	}

	// Mock product not found
	mock.ExpectQuery(`SELECT quantity FROM listings WHERE id = \$1 AND source_type = 'b2c'`).
		WithArgs(int64(999)).
		WillReturnError(sql.ErrNoRows)

	sqlxDB := sqlx.NewDb(db, "postgres")
	mockRepo.On("GetDB").Return(sqlxDB)

	allAvailable, availabilities, err := svc.CheckStockAvailability(ctx, items)

	require.NoError(t, err)
	assert.False(t, allAvailable)
	require.Len(t, availabilities, 1)
	assert.False(t, availabilities[0].IsAvailable)
	assert.Equal(t, int32(0), availabilities[0].AvailableQuantity)

	assert.NoError(t, mock.ExpectationsWereMet())
	mockRepo.AssertExpectations(t)
}

func TestCheckStockAvailability_Unavailable_NotFound_Variant(t *testing.T) {
	svc, mockRepo, mock, db := setupStockTest(t)
	defer db.Close()

	ctx := context.Background()
	variantID := int64(999)

	items := []StockItem{
		{ProductID: 100, VariantID: &variantID, Quantity: 3},
	}

	// Mock variant not found
	mock.ExpectQuery(`SELECT stock_quantity FROM b2c_product_variants WHERE id = \$1 AND product_id = \$2`).
		WithArgs(variantID, int64(100)).
		WillReturnError(sql.ErrNoRows)

	sqlxDB := sqlx.NewDb(db, "postgres")
	mockRepo.On("GetDB").Return(sqlxDB)

	allAvailable, availabilities, err := svc.CheckStockAvailability(ctx, items)

	require.NoError(t, err)
	assert.False(t, allAvailable)
	require.Len(t, availabilities, 1)
	assert.False(t, availabilities[0].IsAvailable)
	assert.Equal(t, int32(0), availabilities[0].AvailableQuantity)

	assert.NoError(t, mock.ExpectationsWereMet())
	mockRepo.AssertExpectations(t)
}

func TestCheckStockAvailability_Mixed_SomeAvailableSomeNot(t *testing.T) {
	svc, mockRepo, mock, db := setupStockTest(t)
	defer db.Close()

	ctx := context.Background()

	items := []StockItem{
		{ProductID: 100, VariantID: nil, Quantity: 5},  // Available
		{ProductID: 200, VariantID: nil, Quantity: 20}, // Not available
	}

	// First product - available
	mock.ExpectQuery(`SELECT quantity FROM listings WHERE id = \$1 AND source_type = 'b2c'`).
		WithArgs(int64(100)).
		WillReturnRows(sqlmock.NewRows([]string{"stock_quantity"}).AddRow(10))

	// Second product - insufficient
	mock.ExpectQuery(`SELECT quantity FROM listings WHERE id = \$1 AND source_type = 'b2c'`).
		WithArgs(int64(200)).
		WillReturnRows(sqlmock.NewRows([]string{"stock_quantity"}).AddRow(15))

	sqlxDB := sqlx.NewDb(db, "postgres")
	mockRepo.On("GetDB").Return(sqlxDB)

	allAvailable, availabilities, err := svc.CheckStockAvailability(ctx, items)

	require.NoError(t, err)
	assert.False(t, allAvailable) // Not ALL available
	require.Len(t, availabilities, 2)

	// First item available
	assert.True(t, availabilities[0].IsAvailable)
	assert.Equal(t, int32(10), availabilities[0].AvailableQuantity)

	// Second item NOT available
	assert.False(t, availabilities[1].IsAvailable)
	assert.Equal(t, int32(15), availabilities[1].AvailableQuantity)
	assert.Equal(t, int32(20), availabilities[1].RequestedQuantity)

	assert.NoError(t, mock.ExpectationsWereMet())
	mockRepo.AssertExpectations(t)
}

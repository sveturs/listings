//go:build integration

package integration

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	pb "github.com/vondi-global/listings/api/proto/listings/v1"
	"github.com/vondi-global/listings/internal/repository/postgres"
	"github.com/vondi-global/listings/internal/service/listings"
	grpchandlers "github.com/vondi-global/listings/internal/transport/grpc"
	"github.com/vondi-global/listings/tests"
)

const (
	rollbackBufSize = 1024 * 1024

	// Test product IDs from fixtures
	testProductSingleRollback     = int64(8000)
	testProductBatchRollback1     = int64(8001)
	testProductBatchRollback2     = int64(8002)
	testProductBatchRollback3     = int64(8003)
	testProductIdempotency        = int64(8004)
	testProductPartialRollback    = int64(8005)
	testProductConcurrent         = int64(8006)
	testProductNoDecrementHistory = int64(8007)
	testProductE2EWorkflow        = int64(8008)
	testProductVariantRollback    = int64(8009)

	// Test variant IDs
	testVariantSizeS = int64(9000)
	testVariantSizeM = int64(9001)
	testVariantSizeL = int64(9002)
)

// setupRollbackTestServer creates a gRPC server with rollback test fixtures
func setupRollbackTestServer(t *testing.T) (pb.ListingsServiceClient, *tests.TestDB, func()) {
	t.Helper()

	tests.SkipIfNoDocker(t)

	// Setup test database
	testDB := tests.SetupTestPostgres(t)

	// Run migrations
	tests.RunMigrations(t, testDB.DB, "../../migrations")

	// Load rollback-specific fixtures
	tests.LoadRollbackStockFixtures(t, testDB.DB)

	// Create sqlx.DB wrapper
	db := sqlx.NewDb(testDB.DB, "postgres")

	// Create logger
	logger := zerolog.New(zerolog.NewTestWriter(t)).With().Timestamp().Logger()

	// Create repository
	repo := postgres.NewRepository(db, logger)

	// Create service (no Redis and no indexer for integration tests)
	service := listings.NewService(repo, nil, nil, logger)

	// Create gRPC server with metrics
	metricsInstance := getTestMetrics()
	server := grpchandlers.NewServer(
		service,
		nil, // storefrontService
		nil, // attrService
		nil, // categoryService
		nil, // orderService
		nil, // cartService
		nil, // chatService
		nil, // analyticsService
		nil, // storefrontAnalyticsService
		nil, // inventoryService
		nil, // invitationService
		nil, // minioClient
		metricsInstance,
		logger,
	)

	// Setup in-memory gRPC connection
	lis := bufconn.Listen(rollbackBufSize)

	grpcServer := grpc.NewServer()
	pb.RegisterListingsServiceServer(grpcServer, server)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			logger.Error().Err(err).Msg("gRPC server failed")
		}
	}()

	// Create client connection
	conn, err := grpc.DialContext(
		context.Background(),
		"bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)

	client := pb.NewListingsServiceClient(conn)

	cleanup := func() {
		conn.Close()
		grpcServer.Stop()
		lis.Close()
		testDB.TeardownTestPostgres(t)
	}

	return client, testDB, cleanup
}

// rollbackStringPtr is a helper to create string pointers for rollback tests
func rollbackStringPtr(s string) *string {
	return &s
}

// ============================================================================
// HAPPY PATH TESTS
// ============================================================================

// TestRollbackStock_SingleProduct_Success tests basic rollback of single product
func TestRollbackStock_SingleProduct_Success(t *testing.T) {
	client, testDB, cleanup := setupRollbackTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	// Verify initial state: product 8000 has 90 stock (was 100, decremented by 10)
	initialStock := tests.GetProductQuantity(t, testDB.DB, testProductSingleRollback)
	assert.Equal(t, int32(90), initialStock, "Initial stock should be 90")

	// Rollback the 10 units
	orderID := rollbackStringPtr("ORDER-001")
	req := &pb.RollbackStockRequest{
		OrderId: orderID,
		Items: []*pb.StockItem{
			{
				ProductId: testProductSingleRollback,
				Quantity:  10,
			},
		},
	}

	resp, err := client.RollbackStock(ctx, req)

	// Assertions
	require.NoError(t, err, "RollbackStock should succeed")
	require.NotNil(t, resp)
	assert.True(t, resp.Success, "Response should indicate success")
	assert.Len(t, resp.Results, 1, "Should have 1 result")

	// Verify result details
	result := resp.Results[0]
	assert.True(t, result.Success, "Result should be successful")
	assert.Equal(t, int32(90), result.StockBefore, "StockBefore should be 90")
	assert.Equal(t, int32(100), result.StockAfter, "StockAfter should be 100 (restored)")
	assert.Nil(t, result.Error, "Should have no error")

	// Verify final stock in database
	finalStock := tests.GetProductQuantity(t, testDB.DB, testProductSingleRollback)
	assert.Equal(t, int32(100), finalStock, "Stock should be restored to 100")
}

// TestRollbackStock_MultipleProducts_Success tests batch rollback for multiple products
func TestRollbackStock_MultipleProducts_Success(t *testing.T) {
	client, testDB, cleanup := setupRollbackTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	// Verify initial states
	stock1 := tests.GetProductQuantity(t, testDB.DB, testProductBatchRollback1)
	stock2 := tests.GetProductQuantity(t, testDB.DB, testProductBatchRollback2)
	stock3 := tests.GetProductQuantity(t, testDB.DB, testProductBatchRollback3)

	assert.Equal(t, int32(80), stock1, "Product 8001 should have 80 stock")
	assert.Equal(t, int32(85), stock2, "Product 8002 should have 85 stock")
	assert.Equal(t, int32(95), stock3, "Product 8003 should have 95 stock")

	// Rollback all three products from ORDER-002
	orderID := rollbackStringPtr("ORDER-002")
	req := &pb.RollbackStockRequest{
		OrderId: orderID,
		Items: []*pb.StockItem{
			{ProductId: testProductBatchRollback1, Quantity: 20}, // restore 20
			{ProductId: testProductBatchRollback2, Quantity: 15}, // restore 15
			{ProductId: testProductBatchRollback3, Quantity: 5},  // restore 5
		},
	}

	resp, err := client.RollbackStock(ctx, req)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.Len(t, resp.Results, 3, "Should have 3 results")

	// Verify all results are successful
	for i, result := range resp.Results {
		assert.True(t, result.Success, "Result %d should be successful", i)
		assert.Nil(t, result.Error, "Result %d should have no error", i)
	}

	// Verify final stocks in database (all should be restored to 100)
	finalStock1 := tests.GetProductQuantity(t, testDB.DB, testProductBatchRollback1)
	finalStock2 := tests.GetProductQuantity(t, testDB.DB, testProductBatchRollback2)
	finalStock3 := tests.GetProductQuantity(t, testDB.DB, testProductBatchRollback3)

	assert.Equal(t, int32(100), finalStock1, "Product 8001 should be restored to 100")
	assert.Equal(t, int32(100), finalStock2, "Product 8002 should be restored to 100")
	assert.Equal(t, int32(100), finalStock3, "Product 8003 should be restored to 100")
}

// TestRollbackStock_PartialOrder tests rollback of partial quantity (not full order)
func TestRollbackStock_PartialOrder(t *testing.T) {
	client, testDB, cleanup := setupRollbackTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	// Product 8005 has 50 stock (was 100, decremented by 50)
	initialStock := tests.GetProductQuantity(t, testDB.DB, testProductPartialRollback)
	assert.Equal(t, int32(50), initialStock)

	// Rollback only 25 units (half of what was decremented)
	orderID := rollbackStringPtr("ORDER-004-PARTIAL")
	req := &pb.RollbackStockRequest{
		OrderId: orderID,
		Items: []*pb.StockItem{
			{ProductId: testProductPartialRollback, Quantity: 25},
		},
	}

	resp, err := client.RollbackStock(ctx, req)

	// Assertions
	require.NoError(t, err)
	assert.True(t, resp.Success)

	// Verify partial restoration
	finalStock := tests.GetProductQuantity(t, testDB.DB, testProductPartialRollback)
	assert.Equal(t, int32(75), finalStock, "Stock should be partially restored to 75 (50 + 25)")
}

// ============================================================================
// IDEMPOTENCY TESTS (CRITICAL!)
// ============================================================================

// TestRollbackStock_DoubleRollback_SameOrderID tests that duplicate rollback is idempotent
// CRITICAL: This test will FAIL if there's no idempotency protection!
func TestRollbackStock_DoubleRollback_SameOrderID(t *testing.T) {
	client, testDB, cleanup := setupRollbackTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	// Product 8004 has 70 stock (was 100, decremented by 30)
	initialStock := tests.GetProductQuantity(t, testDB.DB, testProductIdempotency)
	assert.Equal(t, int32(70), initialStock)

	orderID := rollbackStringPtr("ORDER-003-IDEMPOTENCY")
	req := &pb.RollbackStockRequest{
		OrderId: orderID,
		Items: []*pb.StockItem{
			{ProductId: testProductIdempotency, Quantity: 30},
		},
	}

	// FIRST rollback - should succeed
	resp1, err1 := client.RollbackStock(ctx, req)
	require.NoError(t, err1, "First rollback should succeed")
	assert.True(t, resp1.Success)

	stockAfterFirst := tests.GetProductQuantity(t, testDB.DB, testProductIdempotency)
	assert.Equal(t, int32(100), stockAfterFirst, "Stock should be restored to 100")

	// SECOND rollback with SAME order_id - should be idempotent!
	resp2, err2 := client.RollbackStock(ctx, req)

	// CRITICAL ASSERTIONS for idempotency
	// Option 1: Should succeed but NOT change stock again
	// Option 2: Should return error indicating "already rolled back"

	// Current implementation doesn't have idempotency, so this will FAIL
	// Testing current behavior first, then we can fix the implementation

	// For now, let's document what SHOULD happen:
	// EXPECTED: err2 == nil && resp2.Success == true
	// EXPECTED: stock remains 100 (not 130!)

	if err2 != nil {
		t.Logf("Second rollback returned error (acceptable if it's 'already rolled back'): %v", err2)
	} else {
		assert.True(t, resp2.Success, "Second rollback should succeed (idempotent)")
	}

	stockAfterSecond := tests.GetProductQuantity(t, testDB.DB, testProductIdempotency)

	// THIS IS THE CRITICAL CHECK: Stock must NOT go above 100!
	assert.Equal(t, int32(100), stockAfterSecond,
		"CRITICAL: Stock should remain 100 after duplicate rollback (idempotency protection)")

	if stockAfterSecond != 100 {
		t.Error("❌ IDEMPOTENCY FAILURE: Double rollback incremented stock beyond original value!")
		t.Error("   This is a CRITICAL BUG in compensating transaction logic!")
		t.Error("   Stock went from 70 → 100 → ", stockAfterSecond)
	}
}

// TestRollbackStock_TripleRollback tests three identical rollback requests
func TestRollbackStock_TripleRollback(t *testing.T) {
	client, testDB, cleanup := setupRollbackTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	// Use product 8006 (60 stock, was 100, decremented by 40)
	initialStock := tests.GetProductQuantity(t, testDB.DB, testProductConcurrent)
	assert.Equal(t, int32(60), initialStock)

	orderID := rollbackStringPtr("ORDER-005-TRIPLE")
	req := &pb.RollbackStockRequest{
		OrderId: orderID,
		Items: []*pb.StockItem{
			{ProductId: testProductConcurrent, Quantity: 40},
		},
	}

	// Execute three times
	for i := 1; i <= 3; i++ {
		resp, err := client.RollbackStock(ctx, req)

		if err != nil {
			// Error after first is acceptable (idempotency)
			if i == 1 {
				require.NoError(t, err, "First rollback must succeed")
			} else {
				t.Logf("Rollback #%d returned error (acceptable): %v", i, err)
			}
		} else {
			assert.True(t, resp.Success, "Rollback #%d response should be successful", i)
		}

		// Check stock after each rollback
		currentStock := tests.GetProductQuantity(t, testDB.DB, testProductConcurrent)
		t.Logf("Stock after rollback #%d: %d", i, currentStock)

		// Stock must NEVER exceed 100
		assert.LessOrEqual(t, currentStock, int32(100),
			"Stock must not exceed 100 after rollback #%d", i)
	}

	// Final verification
	finalStock := tests.GetProductQuantity(t, testDB.DB, testProductConcurrent)
	assert.Equal(t, int32(100), finalStock,
		"Final stock should be exactly 100, not more (idempotency)")
}

// TestRollbackStock_ConcurrentRollbacks_SameOrder tests concurrent rollback calls
func TestRollbackStock_ConcurrentRollbacks_SameOrder(t *testing.T) {
	client, testDB, cleanup := setupRollbackTestServer(t)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Use fresh product with known state
	productID := testProductSingleRollback
	initialStock := tests.GetProductQuantity(t, testDB.DB, productID)

	orderID := rollbackStringPtr("ORDER-CONCURRENT")
	req := &pb.RollbackStockRequest{
		OrderId: orderID,
		Items: []*pb.StockItem{
			{ProductId: productID, Quantity: 10},
		},
	}

	// Launch 5 concurrent rollback requests
	concurrency := 5
	var wg sync.WaitGroup
	results := make(chan error, concurrency)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			_, err := client.RollbackStock(ctx, req)
			results <- err
		}(i)
	}

	wg.Wait()
	close(results)

	// Check results
	successCount := 0
	errorCount := 0
	for err := range results {
		if err == nil {
			successCount++
		} else {
			errorCount++
			t.Logf("Concurrent rollback error: %v", err)
		}
	}

	t.Logf("Concurrent rollbacks: %d succeeded, %d failed", successCount, errorCount)

	// At least ONE should succeed
	assert.GreaterOrEqual(t, successCount, 1, "At least one rollback should succeed")

	// CRITICAL: Final stock must be correct (initial + rollback quantity)
	finalStock := tests.GetProductQuantity(t, testDB.DB, productID)
	expectedStock := initialStock + 10

	assert.Equal(t, expectedStock, finalStock,
		"Stock should be incremented exactly once despite concurrent requests")

	if finalStock != expectedStock {
		t.Errorf("❌ RACE CONDITION: Expected stock %d, got %d", expectedStock, finalStock)
	}
}

// ============================================================================
// ERROR CASES TESTS
// ============================================================================

// TestRollbackStock_ProductNotFound tests rollback of non-existent product
func TestRollbackStock_ProductNotFound(t *testing.T) {
	client, _, cleanup := setupRollbackTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	orderID := rollbackStringPtr("ORDER-NOTFOUND")
	req := &pb.RollbackStockRequest{
		OrderId: orderID,
		Items: []*pb.StockItem{
			{ProductId: 99999, Quantity: 10}, // Non-existent product
		},
	}

	resp, err := client.RollbackStock(ctx, req)

	// Current implementation returns success=false with error in results
	// NOT a gRPC error
	require.NoError(t, err, "gRPC call should not error")
	assert.False(t, resp.Success, "Response should indicate failure")

	// Check result contains error
	require.Len(t, resp.Results, 1)
	assert.False(t, resp.Results[0].Success)
	assert.NotNil(t, resp.Results[0].Error)
	assert.Contains(t, *resp.Results[0].Error, "not found",
		"Error should mention product not found")
}

// TestRollbackStock_InvalidOrderID tests rollback with empty/invalid order ID
func TestRollbackStock_InvalidOrderID(t *testing.T) {
	client, _, cleanup := setupRollbackTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	// Test with nil order_id
	req := &pb.RollbackStockRequest{
		OrderId: nil, // No order ID
		Items: []*pb.StockItem{
			{ProductId: testProductSingleRollback, Quantity: 10},
		},
	}

	resp, err := client.RollbackStock(ctx, req)

	// Should still work (order_id is optional in current implementation)
	// But ideally should be required for audit trail
	require.NoError(t, err)

	// Log a warning if order_id is not enforced
	if resp.Success {
		t.Log("⚠️  WARNING: RollbackStock succeeded without order_id")
		t.Log("   Consider making order_id mandatory for audit trail")
	}
}

// TestRollbackStock_EmptyItems tests rollback with no items
func TestRollbackStock_EmptyItems(t *testing.T) {
	client, _, cleanup := setupRollbackTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	orderID := rollbackStringPtr("ORDER-EMPTY")
	req := &pb.RollbackStockRequest{
		OrderId: orderID,
		Items:   []*pb.StockItem{}, // Empty items
	}

	resp, err := client.RollbackStock(ctx, req)

	// Should return error
	require.NoError(t, err, "gRPC should not error")
	assert.False(t, resp.Success, "Response should indicate failure")
	assert.NotNil(t, resp.Error, "Should have error message")
	assert.Contains(t, *resp.Error, "no items", "Error should mention no items")
}

// TestRollbackStock_InvalidQuantity tests rollback with zero/negative quantity
func TestRollbackStock_InvalidQuantity(t *testing.T) {
	client, _, cleanup := setupRollbackTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	testCases := []struct {
		name     string
		quantity int32
	}{
		{"Zero quantity", 0},
		{"Negative quantity", -10},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := &pb.RollbackStockRequest{
				OrderId: rollbackStringPtr("ORDER-INVALID-QTY"),
				Items: []*pb.StockItem{
					{ProductId: testProductSingleRollback, Quantity: tc.quantity},
				},
			}

			resp, err := client.RollbackStock(ctx, req)

			require.NoError(t, err, "gRPC should not error")
			assert.False(t, resp.Success, "Should fail validation")
			assert.NotNil(t, resp.Error, "Should have error message")
		})
	}
}

// TODO: Continue with Saga Pattern tests, Integration tests, and E2E test
// This file is getting long, will continue in next section...

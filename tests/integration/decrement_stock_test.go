//go:build integration

package integration

import (
	"context"
	"fmt"
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

// setupDecrementStockTest creates a test environment with decrement stock fixtures
func setupDecrementStockTest(t *testing.T) (pb.ListingsServiceClient, *tests.TestDB, func()) {
	t.Helper()

	tests.SkipIfNoDocker(t)

	// Setup test database
	testDB := tests.SetupTestPostgres(t)

	// Run migrations
	tests.RunMigrations(t, testDB.DB, "../../migrations")

	// Load decrement stock fixtures
	tests.LoadTestFixtures(t, testDB.DB, "../fixtures/decrement_stock_fixtures.sql")

	// Create sqlx.DB wrapper
	db := sqlx.NewDb(testDB.DB, "postgres")

	// Create logger
	logger := zerolog.New(zerolog.NewTestWriter(t)).With().Timestamp().Logger()

	// Create repository
	repo := postgres.NewRepository(db, logger)

	// Create service (no Redis and no indexer for integration tests)
	service := listings.NewService(repo, nil, nil, logger)

	// Create gRPC server (with nil metrics for testing)
	server := grpchandlers.NewServer(
		service,
		nil, // storefrontService
		nil, // attrService
		nil, // categoryService
		nil, // categoryRepoV2 (Phase 1, not used in tests)
		nil, // categoryCache (Phase 1, not used in tests)
		nil, // orderService
		nil, // cartService
		nil, // chatService
		nil, // analyticsService
		nil, // storefrontAnalyticsService
		nil, // inventoryService
		nil, // invitationService
		nil, // minioClient
		nil, // metrics
		logger,
	)

	// Setup in-memory gRPC connection using bufconn
	const bufSize = 1024 * 1024
	lis := bufconn.Listen(bufSize)

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

// ============================================================================
// Happy Path Tests
// ============================================================================

// TestDecrementStock_SingleProduct_Success tests decrementing stock for a single product
func TestDecrementStock_SingleProduct_Success(t *testing.T) {
	client, testDB, cleanup := setupDecrementStockTest(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	productID := int64(8000) // Initial stock: 100
	initialStock := tests.GetProductQuantity(t, testDB.DB, productID)
	assert.Equal(t, int32(100), initialStock)

	req := &pb.DecrementStockRequest{
		Items: []*pb.StockItem{
			{
				ProductId: productID,
				Quantity:  25,
			},
		},
		OrderId: stringPtr("ORDER-TEST-001"),
	}

	resp, err := client.DecrementStock(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.Nil(t, resp.Error)
	assert.Len(t, resp.Results, 1)

	result := resp.Results[0]
	assert.Equal(t, productID, result.ProductId)
	assert.Nil(t, result.VariantId)
	assert.Equal(t, int32(100), result.StockBefore)
	assert.Equal(t, int32(75), result.StockAfter)
	assert.True(t, result.Success)
	assert.Nil(t, result.Error)

	// Verify database state
	finalStock := tests.GetProductQuantity(t, testDB.DB, productID)
	assert.Equal(t, int32(75), finalStock)
}

// TestDecrementStock_MultipleProducts_Success tests batch decrement for multiple products
func TestDecrementStock_MultipleProducts_Success(t *testing.T) {
	client, testDB, cleanup := setupDecrementStockTest(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	product1 := int64(8000) // Stock: 100
	product2 := int64(8001) // Stock: 50
	product3 := int64(8002) // Stock: 200

	req := &pb.DecrementStockRequest{
		Items: []*pb.StockItem{
			{ProductId: product1, Quantity: 10},
			{ProductId: product2, Quantity: 5},
			{ProductId: product3, Quantity: 20},
		},
		OrderId: stringPtr("ORDER-BATCH-001"),
	}

	resp, err := client.DecrementStock(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.Nil(t, resp.Error)
	assert.Len(t, resp.Results, 3)

	// Verify all results
	assert.Equal(t, int32(90), resp.Results[0].StockAfter)
	assert.Equal(t, int32(45), resp.Results[1].StockAfter)
	assert.Equal(t, int32(180), resp.Results[2].StockAfter)

	// Verify database state
	assert.Equal(t, int32(90), tests.GetProductQuantity(t, testDB.DB, product1))
	assert.Equal(t, int32(45), tests.GetProductQuantity(t, testDB.DB, product2))
	assert.Equal(t, int32(180), tests.GetProductQuantity(t, testDB.DB, product3))
}

// TestDecrementStock_ExactQuantity tests decrementing entire available stock
func TestDecrementStock_ExactQuantity(t *testing.T) {
	client, testDB, cleanup := setupDecrementStockTest(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	productID := int64(8003) // Initial stock: 10
	initialStock := tests.GetProductQuantity(t, testDB.DB, productID)
	assert.Equal(t, int32(10), initialStock)

	req := &pb.DecrementStockRequest{
		Items: []*pb.StockItem{
			{
				ProductId: productID,
				Quantity:  10, // Exact quantity
			},
		},
		OrderId: stringPtr("ORDER-EXACT-001"),
	}

	resp, err := client.DecrementStock(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.Len(t, resp.Results, 1)

	result := resp.Results[0]
	assert.Equal(t, int32(10), result.StockBefore)
	assert.Equal(t, int32(0), result.StockAfter)
	assert.True(t, result.Success)

	// Verify stock is now zero
	finalStock := tests.GetProductQuantity(t, testDB.DB, productID)
	assert.Equal(t, int32(0), finalStock)
}

// TestDecrementStock_VariantLevel tests decrementing variant stock
func TestDecrementStock_VariantLevel(t *testing.T) {
	client, testDB, cleanup := setupDecrementStockTest(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	productID := int64(8004) // Product with variants
	variantID := int64(9000) // Size S, stock: 50
	initialStock := tests.GetVariantQuantity(t, testDB.DB, variantID)
	assert.Equal(t, int32(50), initialStock)

	req := &pb.DecrementStockRequest{
		Items: []*pb.StockItem{
			{
				ProductId: productID,
				VariantId: &variantID,
				Quantity:  15,
			},
		},
		OrderId: stringPtr("ORDER-VARIANT-001"),
	}

	resp, err := client.DecrementStock(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.Len(t, resp.Results, 1)

	result := resp.Results[0]
	assert.Equal(t, productID, result.ProductId)
	assert.NotNil(t, result.VariantId)
	assert.Equal(t, variantID, *result.VariantId)
	assert.Equal(t, int32(50), result.StockBefore)
	assert.Equal(t, int32(35), result.StockAfter)
	assert.True(t, result.Success)

	// Verify variant stock
	finalStock := tests.GetVariantQuantity(t, testDB.DB, variantID)
	assert.Equal(t, int32(35), finalStock)
}

// ============================================================================
// Error Cases Tests
// ============================================================================

// TestDecrementStock_InsufficientStock tests attempting to decrement more than available
func TestDecrementStock_InsufficientStock(t *testing.T) {
	client, testDB, cleanup := setupDecrementStockTest(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	productID := int64(8003) // Stock: 10
	initialStock := tests.GetProductQuantity(t, testDB.DB, productID)
	assert.Equal(t, int32(10), initialStock)

	req := &pb.DecrementStockRequest{
		Items: []*pb.StockItem{
			{
				ProductId: productID,
				Quantity:  50, // More than available
			},
		},
		OrderId: stringPtr("ORDER-INSUFFICIENT-001"),
	}

	resp, err := client.DecrementStock(ctx, req)

	// Service returns response with error, not gRPC error
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.False(t, resp.Success)
	assert.NotNil(t, resp.Error)
	assert.Contains(t, *resp.Error, "insufficient stock")

	// Verify stock NOT changed
	finalStock := tests.GetProductQuantity(t, testDB.DB, productID)
	assert.Equal(t, initialStock, finalStock, "Stock should not change on failure")
}

// TestDecrementStock_ProductNotFound tests attempting to decrement non-existent product
func TestDecrementStock_ProductNotFound(t *testing.T) {
	client, _, cleanup := setupDecrementStockTest(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	req := &pb.DecrementStockRequest{
		Items: []*pb.StockItem{
			{
				ProductId: 99999, // Non-existent
				Quantity:  10,
			},
		},
		OrderId: stringPtr("ORDER-NOTFOUND-001"),
	}

	resp, err := client.DecrementStock(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.False(t, resp.Success)
	assert.NotNil(t, resp.Error)
	assert.Contains(t, *resp.Error, "not found")
}

// TestDecrementStock_VariantNotFound tests attempting to decrement non-existent variant
func TestDecrementStock_VariantNotFound(t *testing.T) {
	client, _, cleanup := setupDecrementStockTest(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	productID := int64(8000)
	variantID := int64(99999) // Non-existent variant

	req := &pb.DecrementStockRequest{
		Items: []*pb.StockItem{
			{
				ProductId: productID,
				VariantId: &variantID,
				Quantity:  10,
			},
		},
		OrderId: stringPtr("ORDER-VARIANT-NOTFOUND-001"),
	}

	resp, err := client.DecrementStock(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.False(t, resp.Success)
	assert.NotNil(t, resp.Error)
	assert.Contains(t, *resp.Error, "not found")
}

// TestDecrementStock_InvalidQuantity tests invalid quantity values
func TestDecrementStock_InvalidQuantity(t *testing.T) {
	client, _, cleanup := setupDecrementStockTest(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	testCases := []struct {
		name     string
		quantity int32
		wantErr  string
	}{
		{
			name:     "Zero quantity",
			quantity: 0,
			wantErr:  "quantity must be positive",
		},
		{
			name:     "Negative quantity",
			quantity: -10,
			wantErr:  "quantity must be positive",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := &pb.DecrementStockRequest{
				Items: []*pb.StockItem{
					{
						ProductId: 8000,
						Quantity:  tc.quantity,
					},
				},
			}

			resp, err := client.DecrementStock(ctx, req)

			require.NoError(t, err)
			require.NotNil(t, resp)
			assert.False(t, resp.Success)
			assert.NotNil(t, resp.Error)
			assert.Contains(t, *resp.Error, tc.wantErr)
		})
	}
}

// TestDecrementStock_EmptyItems tests request with no items
func TestDecrementStock_EmptyItems(t *testing.T) {
	client, _, cleanup := setupDecrementStockTest(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	req := &pb.DecrementStockRequest{
		Items:   []*pb.StockItem{}, // Empty
		OrderId: stringPtr("ORDER-EMPTY-001"),
	}

	resp, err := client.DecrementStock(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.False(t, resp.Success)
	assert.NotNil(t, resp.Error)
	assert.Contains(t, *resp.Error, "no items")
}

// ============================================================================
// Concurrency Tests (CRITICAL!)
// ============================================================================

// TestDecrementStock_ConcurrentDecrements_SameProduct tests concurrent decrements of same product
// This is CRITICAL to verify no overselling due to race conditions
func TestDecrementStock_ConcurrentDecrements_SameProduct(t *testing.T) {
	client, testDB, cleanup := setupDecrementStockTest(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	productID := int64(8005) // Stock: 100 (for concurrency tests)
	initialStock := tests.GetProductQuantity(t, testDB.DB, productID)
	assert.Equal(t, int32(100), initialStock)

	concurrentRequests := 10
	quantityPerRequest := int32(5)
	expectedFinalStock := initialStock - (int32(concurrentRequests) * quantityPerRequest)

	var wg sync.WaitGroup
	results := make(chan *pb.DecrementStockResponse, concurrentRequests)
	errors := make(chan error, concurrentRequests)

	// Launch concurrent decrements
	for i := 0; i < concurrentRequests; i++ {
		wg.Add(1)
		go func(requestNum int) {
			defer wg.Done()

			req := &pb.DecrementStockRequest{
				Items: []*pb.StockItem{
					{
						ProductId: productID,
						Quantity:  quantityPerRequest,
					},
				},
				OrderId: stringPtr(fmt.Sprintf("ORDER-CONCURRENT-%d", requestNum)),
			}

			resp, err := client.DecrementStock(ctx, req)
			if err != nil {
				errors <- err
				return
			}
			results <- resp
		}(i)
	}

	wg.Wait()
	close(results)
	close(errors)

	// Check for errors
	for err := range errors {
		require.NoError(t, err, "No concurrent requests should fail")
	}

	// Verify all responses are successful
	successCount := 0
	for resp := range results {
		require.NotNil(t, resp)
		assert.True(t, resp.Success, "All concurrent decrements should succeed")
		assert.Len(t, resp.Results, 1)
		successCount++
	}

	assert.Equal(t, concurrentRequests, successCount, "All requests should succeed")

	// CRITICAL: Verify final stock is exactly as expected (no overselling!)
	finalStock := tests.GetProductQuantity(t, testDB.DB, productID)
	assert.Equal(t, expectedFinalStock, finalStock,
		"Final stock must be exactly initial - (concurrent_requests * quantity). NO overselling allowed!")
}

// TestDecrementStock_ConcurrentDecrements_DifferentProducts tests concurrent decrements of different products
func TestDecrementStock_ConcurrentDecrements_DifferentProducts(t *testing.T) {
	client, testDB, cleanup := setupDecrementStockTest(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	products := []struct {
		id            int64
		initialStock  int32
		decrement     int32
		expectedFinal int32
	}{
		{8000, 100, 10, 90},
		{8001, 50, 5, 45},
		{8002, 200, 20, 180},
	}

	var wg sync.WaitGroup
	errors := make(chan error, len(products))

	// Launch concurrent decrements for different products
	for _, p := range products {
		wg.Add(1)
		go func(productID int64, quantity int32) {
			defer wg.Done()

			req := &pb.DecrementStockRequest{
				Items: []*pb.StockItem{
					{
						ProductId: productID,
						Quantity:  quantity,
					},
				},
				OrderId: stringPtr(fmt.Sprintf("ORDER-DIFF-PROD-%d", productID)),
			}

			resp, err := client.DecrementStock(ctx, req)
			if err != nil {
				errors <- err
				return
			}

			if !resp.Success {
				errors <- fmt.Errorf("decrement failed for product %d", productID)
			}
		}(p.id, p.decrement)
	}

	wg.Wait()
	close(errors)

	// Check for errors
	for err := range errors {
		require.NoError(t, err)
	}

	// Verify each product's final stock
	for _, p := range products {
		finalStock := tests.GetProductQuantity(t, testDB.DB, p.id)
		assert.Equal(t, p.expectedFinal, finalStock,
			"Product %d should have correct final stock", p.id)
	}
}

// TestDecrementStock_ConcurrentDecrement_Overselling tests prevention of overselling
// Two concurrent requests trying to buy more than available - one MUST fail
func TestDecrementStock_ConcurrentDecrement_Overselling(t *testing.T) {
	client, testDB, cleanup := setupDecrementStockTest(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	productID := int64(8006) // Stock: 20
	initialStock := tests.GetProductQuantity(t, testDB.DB, productID)
	assert.Equal(t, int32(20), initialStock)

	var wg sync.WaitGroup
	results := make(chan *pb.DecrementStockResponse, 2)

	// Launch 2 concurrent requests: each wants 15 units
	// Total: 30 units, but only 20 available
	// Expected: ONE succeeds, ONE fails
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(requestNum int) {
			defer wg.Done()

			req := &pb.DecrementStockRequest{
				Items: []*pb.StockItem{
					{
						ProductId: productID,
						Quantity:  15, // Each wants 15 (total 30 > 20 available)
					},
				},
				OrderId: stringPtr(fmt.Sprintf("ORDER-OVERSELL-%d", requestNum)),
			}

			// Add small random delay to increase race condition chance
			time.Sleep(time.Millisecond * time.Duration(requestNum*2))

			resp, _ := client.DecrementStock(ctx, req)
			results <- resp
		}(i)
	}

	wg.Wait()
	close(results)

	// Count successes and failures
	successCount := 0
	failureCount := 0

	for resp := range results {
		require.NotNil(t, resp)
		if resp.Success {
			successCount++
		} else {
			failureCount++
			assert.NotNil(t, resp.Error)
			assert.Contains(t, *resp.Error, "insufficient stock")
		}
	}

	// CRITICAL: Exactly ONE should succeed, ONE should fail
	assert.Equal(t, 1, successCount, "Exactly ONE request should succeed")
	assert.Equal(t, 1, failureCount, "Exactly ONE request should fail")

	// CRITICAL: Final stock should be 5 (20 - 15), NOT negative or zero
	finalStock := tests.GetProductQuantity(t, testDB.DB, productID)
	assert.Equal(t, int32(5), finalStock,
		"Final stock must be 5 (20-15). NO overselling allowed!")
}

// ============================================================================
// Transaction Tests
// ============================================================================

// TestDecrementStock_BatchPartialFailure tests atomic batch operations
// If ONE item fails, ENTIRE batch should rollback
func TestDecrementStock_BatchPartialFailure(t *testing.T) {
	client, testDB, cleanup := setupDecrementStockTest(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	product1 := int64(8000) // Stock: 100
	product2 := int64(8003) // Stock: 10
	product3 := int64(8001) // Stock: 50

	initial1 := tests.GetProductQuantity(t, testDB.DB, product1)
	initial2 := tests.GetProductQuantity(t, testDB.DB, product2)
	initial3 := tests.GetProductQuantity(t, testDB.DB, product3)

	req := &pb.DecrementStockRequest{
		Items: []*pb.StockItem{
			{ProductId: product1, Quantity: 10}, // OK
			{ProductId: product2, Quantity: 50}, // FAIL: insufficient (only 10 available)
			{ProductId: product3, Quantity: 5},  // OK if not for batch failure
		},
		OrderId: stringPtr("ORDER-PARTIAL-FAIL-001"),
	}

	resp, err := client.DecrementStock(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.False(t, resp.Success, "Batch should fail if any item fails")
	assert.NotNil(t, resp.Error)

	// CRITICAL: Verify NO stock was changed (atomic rollback)
	final1 := tests.GetProductQuantity(t, testDB.DB, product1)
	final2 := tests.GetProductQuantity(t, testDB.DB, product2)
	final3 := tests.GetProductQuantity(t, testDB.DB, product3)

	assert.Equal(t, initial1, final1, "Product 1 stock should NOT change on batch failure")
	assert.Equal(t, initial2, final2, "Product 2 stock should NOT change")
	assert.Equal(t, initial3, final3, "Product 3 stock should NOT change on batch failure")
}

// TestDecrementStock_TransactionIsolation tests that DecrementStock is atomic
// We verify that failed transactions don't change stock
func TestDecrementStock_TransactionIsolation(t *testing.T) {
	client, testDB, cleanup := setupDecrementStockTest(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	productID := int64(8007) // Stock: 50
	initialStock := tests.GetProductQuantity(t, testDB.DB, productID)
	assert.Equal(t, int32(50), initialStock)

	// Try to decrement with insufficient stock (should fail and rollback)
	req := &pb.DecrementStockRequest{
		Items: []*pb.StockItem{
			{
				ProductId: productID,
				Quantity:  100, // More than available
			},
		},
		OrderId: stringPtr("ORDER-ISOLATION-FAIL"),
	}

	resp, err := client.DecrementStock(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.False(t, resp.Success, "Request should fail")

	// Verify stock NOT changed (transaction was rolled back)
	stockAfterFail := tests.GetProductQuantity(t, testDB.DB, productID)
	assert.Equal(t, initialStock, stockAfterFail, "Stock should NOT change on failed transaction")

	// Now do a successful decrement
	req2 := &pb.DecrementStockRequest{
		Items: []*pb.StockItem{
			{
				ProductId: productID,
				Quantity:  5,
			},
		},
		OrderId: stringPtr("ORDER-ISOLATION-SUCCESS"),
	}

	resp2, err := client.DecrementStock(ctx, req2)
	require.NoError(t, err)
	require.NotNil(t, resp2)
	assert.True(t, resp2.Success)

	// Verify stock DID change (transaction was committed)
	finalStock := tests.GetProductQuantity(t, testDB.DB, productID)
	assert.Equal(t, int32(45), finalStock, "Stock should change on successful transaction")
}

// ============================================================================
// Performance Tests
// ============================================================================

// TestDecrementStock_LargeBatch tests performance with large batch
func TestDecrementStock_LargeBatch(t *testing.T) {
	client, testDB, cleanup := setupDecrementStockTest(t)
	defer cleanup()
	_ = testDB // Used for potential verification

	ctx := tests.TestContext(t)

	// Create batch of 50 items (using products 8010-8059)
	batchSize := 50
	items := make([]*pb.StockItem, batchSize)
	for i := 0; i < batchSize; i++ {
		items[i] = &pb.StockItem{
			ProductId: int64(8010 + i),
			Quantity:  1,
		}
	}

	req := &pb.DecrementStockRequest{
		Items:   items,
		OrderId: stringPtr("ORDER-LARGE-BATCH-001"),
	}

	start := time.Now()
	resp, err := client.DecrementStock(ctx, req)
	elapsed := time.Since(start)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.Len(t, resp.Results, batchSize)

	// Performance requirement: < 500ms for 50 items
	assert.Less(t, elapsed, 500*time.Millisecond,
		"Large batch (50 items) should complete in < 500ms, got: %v", elapsed)

	t.Logf("Large batch (%d items) completed in: %v", batchSize, elapsed)
}

// TestDecrementStock_HighFrequency tests sequential high-frequency decrements
func TestDecrementStock_HighFrequency(t *testing.T) {
	client, testDB, cleanup := setupDecrementStockTest(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	productID := int64(8008) // Stock: 1000 (large stock for high frequency)
	initialStock := tests.GetProductQuantity(t, testDB.DB, productID)
	assert.Equal(t, int32(1000), initialStock)

	iterations := 100
	quantityPerIteration := int32(1)

	start := time.Now()
	for i := 0; i < iterations; i++ {
		req := &pb.DecrementStockRequest{
			Items: []*pb.StockItem{
				{
					ProductId: productID,
					Quantity:  quantityPerIteration,
				},
			},
			OrderId: stringPtr(fmt.Sprintf("ORDER-HIGHFREQ-%d", i)),
		}

		resp, err := client.DecrementStock(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.True(t, resp.Success, "Iteration %d should succeed", i)
	}
	elapsed := time.Since(start)

	// Verify final stock
	expectedFinalStock := initialStock - (int32(iterations) * quantityPerIteration)
	finalStock := tests.GetProductQuantity(t, testDB.DB, productID)
	assert.Equal(t, expectedFinalStock, finalStock,
		"All %d high-frequency decrements should be applied correctly", iterations)

	// Performance: Should complete 100 sequential operations reasonably fast
	t.Logf("High frequency test (%d iterations) completed in: %v (avg: %v per op)",
		iterations, elapsed, elapsed/time.Duration(iterations))

	// No deadlocks should occur
	assert.NotPanics(t, func() {
		// Additional verification: can still operate
		req := &pb.DecrementStockRequest{
			Items: []*pb.StockItem{
				{ProductId: productID, Quantity: 1},
			},
		}
		_, _ = client.DecrementStock(ctx, req)
	})
}

// TestDecrementStock_Performance_SingleItem tests single item performance
func TestDecrementStock_Performance_SingleItem(t *testing.T) {
	client, testDB, cleanup := setupDecrementStockTest(t)
	defer cleanup()
	_ = testDB // Used for potential verification

	ctx := tests.TestContext(t)

	productID := int64(8000) // Stock: 100

	req := &pb.DecrementStockRequest{
		Items: []*pb.StockItem{
			{
				ProductId: productID,
				Quantity:  1,
			},
		},
		OrderId: stringPtr("ORDER-PERF-SINGLE-001"),
	}

	start := time.Now()
	resp, err := client.DecrementStock(ctx, req)
	elapsed := time.Since(start)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.True(t, resp.Success)

	// Performance requirement: single decrement < 50ms
	assert.Less(t, elapsed, 50*time.Millisecond,
		"Single decrement should complete in < 50ms, got: %v", elapsed)

	t.Logf("Single decrement completed in: %v", elapsed)
}

// ============================================================================
// Helper Functions
// ============================================================================

//go:build integration
// +build integration

package integration

import (
	"context"
	"net"
	"sync"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/sveturs/listings/api/proto/listings/v1"
	"github.com/sveturs/listings/internal/metrics"
	"github.com/sveturs/listings/internal/repository/postgres"
	"github.com/sveturs/listings/internal/service/listings"
	grpchandlers "github.com/sveturs/listings/internal/transport/grpc"
	"github.com/sveturs/listings/tests"
)

const bufSize = 1024 * 1024

var (
	// testMetrics is a singleton instance of metrics for all integration tests
	// to avoid duplicate Prometheus registration errors
	testMetrics     *metrics.Metrics
	testMetricsOnce sync.Once
)

// getTestMetrics returns a singleton metrics instance for testing
func getTestMetrics() *metrics.Metrics {
	testMetricsOnce.Do(func() {
		testMetrics = metrics.NewMetrics("listings_test")
	})
	return testMetrics
}

// setupGRPCTestServer creates a gRPC server with real database
func setupGRPCTestServer(t *testing.T) (pb.ListingsServiceClient, *tests.TestDB, func()) {
	t.Helper()

	tests.SkipIfNoDocker(t)

	// Setup test database
	testDB := tests.SetupTestPostgres(t)

	// Run migrations
	tests.RunMigrations(t, testDB.DB, "../../migrations")

	// Load inventory fixtures
	tests.LoadInventoryFixtures(t, testDB.DB)

	// Create sqlx.DB wrapper
	db := sqlx.NewDb(testDB.DB, "postgres")

	// Create logger
	logger := zerolog.New(zerolog.NewTestWriter(t)).With().Timestamp().Logger()

	// Create repository
	repo := postgres.NewRepository(db, logger)

	// Create service (no Redis and no indexer for integration tests)
	service := listings.NewService(repo, nil, nil, logger)

	// Get singleton metrics instance
	m := getTestMetrics()

	// Create gRPC server
	server := grpchandlers.NewServer(service, m, logger)

	// Setup in-memory gRPC connection
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

// TestGRPCRecordInventoryMovement_FullCycle tests full gRPC request/response cycle
func TestGRPCRecordInventoryMovement_FullCycle(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	productID := int64(5000) // Initial quantity: 100

	testCases := []struct {
		name         string
		request      *pb.RecordInventoryMovementRequest
		wantErr      bool
		wantCode     codes.Code
		wantBefore   int32
		wantAfter    int32
	}{
		{
			name: "Valid stock IN movement",
			request: &pb.RecordInventoryMovementRequest{
				StorefrontId: 1000,
				ProductId:    productID,
				MovementType: "in",
				Quantity:     50,
				Reason:       stringPtr("restock"),
				Notes:        stringPtr("Integration test restock"),
				UserId:       1000,
			},
			wantErr:    false,
			wantBefore: 100,
			wantAfter:  150,
		},
		{
			name: "Valid stock OUT movement",
			request: &pb.RecordInventoryMovementRequest{
				StorefrontId: 1000,
				ProductId:    productID,
				MovementType: "out",
				Quantity:     30,
				Reason:       stringPtr("sale"),
				Notes:        stringPtr("Customer order"),
				UserId:       1000,
			},
			wantErr:    false,
			wantBefore: 150,
			wantAfter:  120,
		},
		{
			name: "Invalid storefront ID",
			request: &pb.RecordInventoryMovementRequest{
				StorefrontId: 0,
				ProductId:    productID,
				MovementType: "in",
				Quantity:     10,
				UserId:       1000,
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "Invalid movement type",
			request: &pb.RecordInventoryMovementRequest{
				StorefrontId: 1000,
				ProductId:    productID,
				MovementType: "invalid",
				Quantity:     10,
				UserId:       1000,
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "Non-existent product",
			request: &pb.RecordInventoryMovementRequest{
				StorefrontId: 1000,
				ProductId:    99999,
				MovementType: "in",
				Quantity:     10,
				UserId:       1000,
			},
			wantErr:  true,
			wantCode: codes.NotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := client.RecordInventoryMovement(ctx, tc.request)

			if tc.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, tc.wantCode, st.Code())
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			assert.True(t, resp.Success)
			assert.Equal(t, tc.wantBefore, resp.StockBefore)
			assert.Equal(t, tc.wantAfter, resp.StockAfter)

			// Verify database state
			actualQty := tests.GetProductQuantity(t, testDB.DB, tc.request.ProductId)
			assert.Equal(t, tc.wantAfter, actualQty)
		})
	}
}

// TestGRPCRecordInventoryMovement_VariantLevel tests variant-level inventory
func TestGRPCRecordInventoryMovement_VariantLevel(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	variantID := int64(6000) // Size S variant, qty 50

	req := &pb.RecordInventoryMovementRequest{
		StorefrontId: 1000,
		ProductId:    5000,
		VariantId:    &variantID,
		MovementType: "in",
		Quantity:     25,
		Reason:       stringPtr("variant_restock"),
		UserId:       1000,
	}

	resp, err := client.RecordInventoryMovement(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.Equal(t, int32(50), resp.StockBefore)
	assert.Equal(t, int32(75), resp.StockAfter)

	// Verify variant quantity in database
	actualQty := tests.GetVariantQuantity(t, testDB.DB, variantID)
	assert.Equal(t, int32(75), actualQty)
}

// TestGRPCBatchUpdateStock_FullCycle tests batch stock update via gRPC
func TestGRPCBatchUpdateStock_FullCycle(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	testCases := []struct {
		name            string
		request         *pb.BatchUpdateStockRequest
		wantErr         bool
		wantCode        codes.Code
		wantSuccessful  int32
		wantFailed      int32
	}{
		{
			name: "Successful batch update",
			request: &pb.BatchUpdateStockRequest{
				StorefrontId: 1000,
				Items: []*pb.StockUpdateItem{
					{
						ProductId: 5003,
						Quantity:  100,
						Reason:    stringPtr("restock"),
					},
					{
						ProductId: 5004,
						Quantity:  80,
						Reason:    stringPtr("adjustment"),
					},
				},
				Reason: stringPtr("bulk_inventory_update"),
				UserId: 1000,
			},
			wantErr:        false,
			wantSuccessful: 2,
			wantFailed:     0,
		},
		{
			name: "Empty items list",
			request: &pb.BatchUpdateStockRequest{
				StorefrontId: 1000,
				Items:        []*pb.StockUpdateItem{},
				UserId:       1000,
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "Invalid storefront ID",
			request: &pb.BatchUpdateStockRequest{
				StorefrontId: 0,
				Items: []*pb.StockUpdateItem{
					{ProductId: 5003, Quantity: 50},
				},
				UserId: 1000,
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "Partial success (some products not found)",
			request: &pb.BatchUpdateStockRequest{
				StorefrontId: 1000,
				Items: []*pb.StockUpdateItem{
					{ProductId: 5003, Quantity: 60}, // Valid
					{ProductId: 99999, Quantity: 50}, // Invalid
				},
				UserId: 1000,
			},
			wantErr:        false,
			wantSuccessful: 1,
			wantFailed:     1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := client.BatchUpdateStock(ctx, tc.request)

			if tc.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, tc.wantCode, st.Code())
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			assert.Equal(t, tc.wantSuccessful, resp.SuccessfulCount)
			assert.Equal(t, tc.wantFailed, resp.FailedCount)
			assert.Len(t, resp.Results, int(tc.wantSuccessful+tc.wantFailed))

			// Verify successful items in database
			for _, item := range tc.request.Items {
				if item.ProductId < 90000 && tests.ProductExists(t, testDB.DB, item.ProductId) {
					actualQty := tests.GetProductQuantity(t, testDB.DB, item.ProductId)
					// Quantity should be updated
					assert.GreaterOrEqual(t, actualQty, int32(0))
				}
			}
		})
	}
}

// TestGRPCGetProductStats_FullCycle tests stats retrieval via gRPC
func TestGRPCGetProductStats_FullCycle(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	testCases := []struct {
		name         string
		storefrontID int64
		wantErr      bool
		wantCode     codes.Code
	}{
		{
			name:         "Valid storefront with products",
			storefrontID: 1000,
			wantErr:      false,
		},
		{
			name:         "Invalid storefront ID (zero)",
			storefrontID: 0,
			wantErr:      true,
			wantCode:     codes.InvalidArgument,
		},
		{
			name:         "Negative storefront ID",
			storefrontID: -1,
			wantErr:      true,
			wantCode:     codes.InvalidArgument,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := &pb.GetProductStatsRequest{
				StorefrontId: tc.storefrontID,
			}

			resp, err := client.GetProductStats(ctx, req)

			if tc.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, tc.wantCode, st.Code())
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			require.NotNil(t, resp.Stats)

			// Verify stats accuracy against database
			dbTotal := tests.CountProductsByStorefront(t, testDB.DB, tc.storefrontID)
			assert.Equal(t, dbTotal, resp.Stats.TotalProducts)

			dbActive := tests.CountActiveProductsByStorefront(t, testDB.DB, tc.storefrontID)
			assert.Equal(t, dbActive, resp.Stats.ActiveProducts)

			dbOutOfStock := tests.CountOutOfStockProducts(t, testDB.DB, tc.storefrontID)
			assert.Equal(t, dbOutOfStock, resp.Stats.OutOfStock)

			dbLowStock := tests.CountLowStockProducts(t, testDB.DB, tc.storefrontID)
			assert.Equal(t, dbLowStock, resp.Stats.LowStock)

			dbValue := tests.GetTotalInventoryValue(t, testDB.DB, tc.storefrontID)
			assert.InDelta(t, dbValue, resp.Stats.TotalValue, 0.01)
		})
	}
}

// TestGRPCIncrementProductViews_FullCycle tests view increment via gRPC
func TestGRPCIncrementProductViews_FullCycle(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	productID := int64(5000)

	testCases := []struct {
		name      string
		productID int64
		wantErr   bool
		wantCode  codes.Code
	}{
		{
			name:      "Valid product",
			productID: productID,
			wantErr:   false,
		},
		{
			name:      "Zero product ID",
			productID: 0,
			wantErr:   true,
			wantCode:  codes.InvalidArgument,
		},
		{
			name:      "Negative product ID",
			productID: -1,
			wantErr:   true,
			wantCode:  codes.InvalidArgument,
		},
		{
			name:      "Non-existent product",
			productID: 99999,
			wantErr:   true,
			wantCode:  codes.Internal, // Or NotFound depending on implementation
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var initialViews int32
			if tc.productID == productID {
				initialViews = tests.GetProductViewCount(t, testDB.DB, tc.productID)
			}

			req := &pb.IncrementProductViewsRequest{
				ProductId: tc.productID,
			}

			resp, err := client.IncrementProductViews(ctx, req)

			if tc.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, tc.wantCode, st.Code())
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			assert.IsType(t, &emptypb.Empty{}, resp)

			// Verify view count incremented in database
			newViews := tests.GetProductViewCount(t, testDB.DB, tc.productID)
			assert.Equal(t, initialViews+1, newViews)
		})
	}
}

// TestGRPCInventoryWorkflow_CompleteScenario tests complete inventory workflow via gRPC
func TestGRPCInventoryWorkflow_CompleteScenario(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	storefrontID := int64(1000)
	productID := int64(5000)
	userID := int64(1000)

	// Step 1: Get initial stats
	statsResp, err := client.GetProductStats(ctx, &pb.GetProductStatsRequest{
		StorefrontId: storefrontID,
	})
	require.NoError(t, err)
	initialValue := statsResp.Stats.TotalValue

	// Step 2: Record stock IN
	movementResp, err := client.RecordInventoryMovement(ctx, &pb.RecordInventoryMovementRequest{
		StorefrontId: storefrontID,
		ProductId:    productID,
		MovementType: "in",
		Quantity:     100,
		Reason:       stringPtr("bulk_restock"),
		UserId:       userID,
	})
	require.NoError(t, err)
	assert.True(t, movementResp.Success)

	// Step 3: Batch update multiple products
	batchResp, err := client.BatchUpdateStock(ctx, &pb.BatchUpdateStockRequest{
		StorefrontId: storefrontID,
		Items: []*pb.StockUpdateItem{
			{ProductId: 5003, Quantity: 150},
			{ProductId: 5004, Quantity: 120},
		},
		Reason: stringPtr("quarterly_audit"),
		UserId: userID,
	})
	require.NoError(t, err)
	assert.Equal(t, int32(2), batchResp.SuccessfulCount)

	// Step 4: Increment product views
	_, err = client.IncrementProductViews(ctx, &pb.IncrementProductViewsRequest{
		ProductId: productID,
	})
	require.NoError(t, err)

	// Step 5: Get updated stats
	finalStatsResp, err := client.GetProductStats(ctx, &pb.GetProductStatsRequest{
		StorefrontId: storefrontID,
	})
	require.NoError(t, err)
	assert.Greater(t, finalStatsResp.Stats.TotalValue, initialValue,
		"Total inventory value should increase after restocking")

	// Step 6: Verify all changes in database
	finalQty := tests.GetProductQuantity(t, testDB.DB, productID)
	assert.Equal(t, int32(200), finalQty, "Product quantity should be 200 after +100")

	finalViews := tests.GetProductViewCount(t, testDB.DB, productID)
	assert.Greater(t, finalViews, int32(0), "Product should have views")

	movementCount := tests.GetInventoryMovementCount(t, testDB.DB, productID)
	assert.GreaterOrEqual(t, movementCount, 1, "Should have inventory movements recorded")
}

// TestGRPCConcurrentRequests tests concurrent gRPC requests
func TestGRPCConcurrentRequests(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	productID := int64(5000)
	concurrentCalls := 10

	t.Run("Concurrent view increments", func(t *testing.T) {
		initialViews := tests.GetProductViewCount(t, testDB.DB, productID)

		done := make(chan error, concurrentCalls)
		for i := 0; i < concurrentCalls; i++ {
			go func() {
				req := &pb.IncrementProductViewsRequest{ProductId: productID}
				_, err := client.IncrementProductViews(ctx, req)
				done <- err
			}()
		}

		// Wait for all
		for i := 0; i < concurrentCalls; i++ {
			err := <-done
			assert.NoError(t, err)
		}

		finalViews := tests.GetProductViewCount(t, testDB.DB, productID)
		assert.Equal(t, initialViews+int32(concurrentCalls), finalViews,
			"All concurrent increments should be recorded")
	})

	t.Run("Concurrent stats reads", func(t *testing.T) {
		done := make(chan error, concurrentCalls)
		for i := 0; i < concurrentCalls; i++ {
			go func() {
				req := &pb.GetProductStatsRequest{StorefrontId: 1000}
				_, err := client.GetProductStats(ctx, req)
				done <- err
			}()
		}

		// All should succeed
		for i := 0; i < concurrentCalls; i++ {
			err := <-done
			assert.NoError(t, err)
		}
	})
}

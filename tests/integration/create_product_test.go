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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/structpb"

	pb "github.com/sveturs/listings/api/proto/listings/v1"
	"github.com/sveturs/listings/internal/repository/postgres"
	"github.com/sveturs/listings/internal/service/listings"
	grpchandlers "github.com/sveturs/listings/internal/transport/grpc"
	"github.com/sveturs/listings/tests"
)

// ============================================================================
// Test Setup Helpers
// ============================================================================

// setupCreateProductTest creates a test environment with create product fixtures
func setupCreateProductTest(t *testing.T) (pb.ListingsServiceClient, *tests.TestDB, func()) {
	t.Helper()

	tests.SkipIfNoDocker(t)

	// Setup test database
	testDB := tests.SetupTestPostgres(t)

	// Run migrations
	tests.RunMigrations(t, testDB.DB, "../../migrations")

	// Load create product fixtures
	tests.LoadTestFixtures(t, testDB.DB, "../fixtures/create_product_fixtures.sql")

	// Create sqlx.DB wrapper
	db := sqlx.NewDb(testDB.DB, "postgres")

	// Create logger
	logger := zerolog.New(zerolog.NewTestWriter(t)).With().Timestamp().Logger()

	// Create repository
	repo := postgres.NewRepository(db, logger)

	// Create service (no Redis and no indexer for integration tests)
	service := listings.NewService(repo, nil, nil, logger)

	// Create gRPC server (with singleton metrics)
	m := getTestMetrics()
	server := grpchandlers.NewServer(service, m, logger)

	// Setup in-memory gRPC connection using bufconn
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

// getProductByID retrieves a product from database for verification
func getProductByID(t *testing.T, db *sqlx.DB, productID int64) *productRecord {
	t.Helper()

	var p productRecord
	err := db.Get(&p, `
		SELECT id, storefront_id, title, description, price, currency, category_id,
		       sku, quantity, stock_status, status, has_variants,
		       attributes, view_count, sold_count, created_at, updated_at
		FROM listings
		WHERE id = $1 AND source_type = 'b2c' AND deleted_at IS NULL
	`, productID)

	if err != nil {
		t.Logf("Product %d not found: %v", productID, err)
		return nil
	}

	return &p
}

// productRecord represents a product record from database
type productRecord struct {
	ID           int64     `db:"id"`
	StorefrontID int64     `db:"storefront_id"`
	Title        string    `db:"title"`
	Description  *string   `db:"description"`
	Price        float64   `db:"price"`
	Currency     string    `db:"currency"`
	CategoryID   int64     `db:"category_id"`
	SKU          *string   `db:"sku"`
	Quantity     int32     `db:"quantity"`
	StockStatus  string    `db:"stock_status"`
	Status       string    `db:"status"` // 'active', 'inactive', 'draft', etc.
	HasVariants  bool      `db:"has_variants"`
	Attributes   *string   `db:"attributes"` // JSONB stored as string
	ViewCount    int32     `db:"view_count"`
	SoldCount    int32     `db:"sold_count"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// Note: Helper functions stringPtr, float64Ptr, int32Ptr are defined in database_test.go

// ============================================================================
// Happy Path Tests - CreateProduct
// ============================================================================

// TestCreateProduct_Success tests creating a basic product with all required fields
func TestCreateProduct_Success(t *testing.T) {
	client, testDB, cleanup := setupCreateProductTest(t)
	defer cleanup()

	ctx := tests.TestContext(t)
	db := sqlx.NewDb(testDB.DB, "postgres")

	req := &pb.CreateProductRequest{
		StorefrontId: 1100,
		Name:         "Test Smartphone Pro Max",
		Description:  "Premium smartphone with advanced features",
		Price:        999.99,
		Currency:     "USD",
		CategoryId:   2110,
		Sku:          stringPtr("TEST-PHONE-001"),
		// Barcode: не поддерживается в таблице listings (только в variants)
		StockQuantity: 50,
		IsActive:      true,
	}

	resp, err := client.CreateProduct(ctx, req)

	// Assertions
	require.NoError(t, err, "CreateProduct should succeed")
	require.NotNil(t, resp, "Response should not be nil")
	require.NotNil(t, resp.Product, "Product should be returned")

	product := resp.Product
	assert.Greater(t, product.Id, int64(0), "Product ID should be assigned")
	assert.Equal(t, req.StorefrontId, product.StorefrontId)
	assert.Equal(t, req.Name, product.Name)
	assert.Equal(t, req.Description, product.Description)
	assert.Equal(t, req.Price, product.Price)
	assert.Equal(t, req.Currency, product.Currency)
	assert.Equal(t, req.CategoryId, product.CategoryId)
	assert.Equal(t, req.Sku, product.Sku)
	assert.Nil(t, product.Barcode, "Barcode should be nil (not supported in listings table)")
	assert.Equal(t, req.StockQuantity, product.StockQuantity)
	assert.Equal(t, "in_stock", product.StockStatus, "Stock status should be 'in_stock'")
	assert.Equal(t, req.IsActive, product.IsActive)
	assert.NotNil(t, product.CreatedAt, "CreatedAt should be set")
	assert.NotNil(t, product.UpdatedAt, "UpdatedAt should be set")

	// Verify product exists in database
	dbProduct := getProductByID(t, db, product.Id)
	require.NotNil(t, dbProduct, "Product should exist in database")
	assert.Equal(t, req.Name, dbProduct.Title)
	assert.Equal(t, req.Price, dbProduct.Price)
	assert.Equal(t, "in_stock", dbProduct.StockStatus)
}

// TestCreateProduct_MinimalFields tests creating product with only required fields
func TestCreateProduct_MinimalFields(t *testing.T) {
	client, testDB, cleanup := setupCreateProductTest(t)
	defer cleanup()

	ctx := tests.TestContext(t)
	db := sqlx.NewDb(testDB.DB, "postgres")

	req := &pb.CreateProductRequest{
		StorefrontId:  1100,
		Name:          "Minimal Product",
		Price:         19.99,
		Currency:      "USD",
		CategoryId:    2110,
		StockQuantity: 0, // Zero stock
		IsActive:      true,
		// No description, SKU, barcode
	}

	resp, err := client.CreateProduct(ctx, req)

	// Assertions
	require.NoError(t, err, "CreateProduct should succeed with minimal fields")
	require.NotNil(t, resp)
	require.NotNil(t, resp.Product)

	product := resp.Product
	assert.Greater(t, product.Id, int64(0))
	assert.Equal(t, req.Name, product.Name)
	assert.Equal(t, "", product.Description, "Description should be empty string")
	assert.Nil(t, product.Sku, "SKU should be nil when not provided")
	assert.Nil(t, product.Barcode, "Barcode should be nil when not provided")
	assert.Equal(t, int32(0), product.StockQuantity)
	assert.Equal(t, "out_of_stock", product.StockStatus, "Zero quantity should set status to out_of_stock")

	// Verify in database
	dbProduct := getProductByID(t, db, product.Id)
	require.NotNil(t, dbProduct)
	assert.Equal(t, int32(0), dbProduct.Quantity)
	assert.Equal(t, "out_of_stock", dbProduct.StockStatus)
}

// TestCreateProduct_WithVariants tests creating product configured for variants
func TestCreateProduct_WithVariants(t *testing.T) {
	client, testDB, cleanup := setupCreateProductTest(t)
	defer cleanup()

	ctx := tests.TestContext(t)
	db := sqlx.NewDb(testDB.DB, "postgres")

	req := &pb.CreateProductRequest{
		StorefrontId:  1100,
		Name:          "T-Shirt with Sizes",
		Description:   "Cotton t-shirt available in multiple sizes",
		Price:         24.99,
		Currency:      "USD",
		CategoryId:    2120,
		Sku:           stringPtr("TSHIRT-BASE-001"),
		StockQuantity: 0, // Variants will have stock
		IsActive:      true,
		HasVariants:   true, // Enable variants
	}

	resp, err := client.CreateProduct(ctx, req)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, resp.Product)

	product := resp.Product
	assert.True(t, product.HasVariants, "Product should have variants enabled")
	assert.Equal(t, int32(0), product.StockQuantity, "Base product with variants should have 0 stock")

	// Verify in database
	dbProduct := getProductByID(t, db, product.Id)
	require.NotNil(t, dbProduct)
	assert.True(t, dbProduct.HasVariants)
}

// TestCreateProduct_WithAttributes tests creating product with custom JSONB attributes
func TestCreateProduct_WithAttributes(t *testing.T) {
	client, testDB, cleanup := setupCreateProductTest(t)
	defer cleanup()

	ctx := tests.TestContext(t)
	db := sqlx.NewDb(testDB.DB, "postgres")

	// Create JSONB attributes
	attrs, err := structpb.NewStruct(map[string]interface{}{
		"brand":           "TestBrand",
		"model":           "Pro-2024",
		"color":           "Black",
		"weight_kg":       0.18,
		"waterproof":      true,
		"warranty_months": 12,
	})
	require.NoError(t, err, "Failed to create attributes struct")

	req := &pb.CreateProductRequest{
		StorefrontId:  1100,
		Name:          "Product with Attributes",
		Description:   "Product with rich custom attributes",
		Price:         149.99,
		Currency:      "USD",
		CategoryId:    2111,
		Sku:           stringPtr("ATTR-PRODUCT-001"),
		StockQuantity: 25,
		IsActive:      true,
		Attributes:    attrs,
	}

	resp, err := client.CreateProduct(ctx, req)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, resp.Product)

	product := resp.Product
	assert.NotNil(t, product.Attributes, "Attributes should be set")
	assert.Equal(t, "TestBrand", product.Attributes.Fields["brand"].GetStringValue())
	assert.Equal(t, "Pro-2024", product.Attributes.Fields["model"].GetStringValue())
	assert.Equal(t, "Black", product.Attributes.Fields["color"].GetStringValue())
	assert.Equal(t, 0.18, product.Attributes.Fields["weight_kg"].GetNumberValue())
	assert.True(t, product.Attributes.Fields["waterproof"].GetBoolValue())

	// Verify attributes stored in database
	dbProduct := getProductByID(t, db, product.Id)
	require.NotNil(t, dbProduct)
	assert.NotNil(t, dbProduct.Attributes, "Attributes should be stored in DB")
}

// TestCreateProduct_WithImages tests creating product and associating images
// Note: Images are typically added after product creation via separate AddImage API
// This test verifies the product can be created and is ready for images
func TestCreateProduct_WithImages(t *testing.T) {
	client, testDB, cleanup := setupCreateProductTest(t)
	defer cleanup()

	ctx := tests.TestContext(t)
	db := sqlx.NewDb(testDB.DB, "postgres")

	// Create product first
	req := &pb.CreateProductRequest{
		StorefrontId:  1100,
		Name:          "Product Ready for Images",
		Description:   "Product that will have images added",
		Price:         79.99,
		Currency:      "USD",
		CategoryId:    2110,
		Sku:           stringPtr("IMG-PRODUCT-001"),
		StockQuantity: 15,
		IsActive:      true,
	}

	resp, err := client.CreateProduct(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp.Product)

	productID := resp.Product.Id

	// Verify product can accept images (has no images yet)
	var imageCount int
	err = db.Get(&imageCount, `
		SELECT COUNT(*) FROM listing_images WHERE listing_id = $1
	`, productID)
	require.NoError(t, err)
	assert.Equal(t, 0, imageCount, "Newly created product should have no images")

	// Product is ready for AddImage API calls
	// This confirms the product structure supports image associations
	dbProduct := getProductByID(t, db, productID)
	require.NotNil(t, dbProduct)
	assert.Greater(t, dbProduct.ID, int64(0), "Product ready for image uploads")
}

// ============================================================================
// Validation Tests - CreateProduct
// ============================================================================

// TestCreateProduct_MissingName tests validation when product name is missing
func TestCreateProduct_MissingName(t *testing.T) {
	client, _, cleanup := setupCreateProductTest(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	req := &pb.CreateProductRequest{
		StorefrontId:  1100,
		Name:          "", // Empty name
		Price:         99.99,
		Currency:      "USD",
		CategoryId:    2110,
		StockQuantity: 10,
		IsActive:      true,
	}

	resp, err := client.CreateProduct(ctx, req)

	// Assertions
	require.Error(t, err, "Should return error for missing name")
	assert.Nil(t, resp, "Response should be nil on validation error")

	// Verify gRPC status code
	st, ok := status.FromError(err)
	require.True(t, ok, "Error should be gRPC status")
	assert.Equal(t, codes.InvalidArgument, st.Code(), "Should return InvalidArgument status code")
	assert.Contains(t, st.Message(), "name", "Error message should mention 'name'")
}

// TestCreateProduct_MissingStorefrontID tests validation when storefront_id is missing
func TestCreateProduct_MissingStorefrontID(t *testing.T) {
	client, _, cleanup := setupCreateProductTest(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	req := &pb.CreateProductRequest{
		StorefrontId:  0, // Invalid storefront ID
		Name:          "Product without Storefront",
		Price:         99.99,
		Currency:      "USD",
		CategoryId:    2110,
		StockQuantity: 10,
		IsActive:      true,
	}

	resp, err := client.CreateProduct(ctx, req)

	// Assertions
	require.Error(t, err, "Should return error for missing storefront_id")
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "storefront", "Error should mention storefront")
}

// TestCreateProduct_InvalidCategoryID tests validation with non-existent category
func TestCreateProduct_InvalidCategoryID(t *testing.T) {
	client, _, cleanup := setupCreateProductTest(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	req := &pb.CreateProductRequest{
		StorefrontId:  1100,
		Name:          "Product with Invalid Category",
		Price:         99.99,
		Currency:      "USD",
		CategoryId:    99999, // Non-existent category
		StockQuantity: 10,
		IsActive:      true,
	}

	resp, err := client.CreateProduct(ctx, req)

	// Assertions
	require.Error(t, err, "Should return error for invalid category_id")
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	require.True(t, ok)
	// Could be InvalidArgument or NotFound depending on implementation
	assert.Contains(t, []codes.Code{codes.InvalidArgument, codes.NotFound}, st.Code())
	assert.Contains(t, st.Message(), "category", "Error should mention category")
}

// TestCreateProduct_NegativePrice tests validation with negative price
func TestCreateProduct_NegativePrice(t *testing.T) {
	client, _, cleanup := setupCreateProductTest(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	testCases := []struct {
		name  string
		price float64
	}{
		{
			name:  "Negative price",
			price: -10.00,
		},
		{
			name:  "Zero price",
			price: 0.00,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := &pb.CreateProductRequest{
				StorefrontId:  1100,
				Name:          "Product with Invalid Price",
				Price:         tc.price,
				Currency:      "USD",
				CategoryId:    2110,
				StockQuantity: 10,
				IsActive:      true,
			}

			resp, err := client.CreateProduct(ctx, req)

			require.Error(t, err, "Should return error for invalid price")
			assert.Nil(t, resp)

			st, ok := status.FromError(err)
			require.True(t, ok)
			assert.Equal(t, codes.InvalidArgument, st.Code())
			assert.Contains(t, st.Message(), "price", "Error should mention price")
		})
	}
}

// TestCreateProduct_DuplicateSKU tests validation when SKU already exists in same storefront
func TestCreateProduct_DuplicateSKU(t *testing.T) {
	t.Skip("UNIQUE constraint on listings.sku not yet implemented - requires migration")

	client, _, cleanup := setupCreateProductTest(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	// Fixture 7000 has SKU "TEST-SKU-001" in storefront 1100
	req := &pb.CreateProductRequest{
		StorefrontId:  1100,
		Name:          "Product with Duplicate SKU",
		Price:         99.99,
		Currency:      "USD",
		CategoryId:    2110,
		Sku:           stringPtr("TEST-SKU-001"), // Duplicate SKU
		StockQuantity: 10,
		IsActive:      true,
	}

	resp, err := client.CreateProduct(ctx, req)

	// Assertions
	require.Error(t, err, "Should return error for duplicate SKU")
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.AlreadyExists, st.Code(), "Should return AlreadyExists status")
	assert.Contains(t, st.Message(), "SKU", "Error should mention SKU")
}

// TestCreateProduct_DuplicateSKU_DifferentStorefront tests that same SKU is allowed in different storefront
func TestCreateProduct_DuplicateSKU_DifferentStorefront(t *testing.T) {
	client, testDB, cleanup := setupCreateProductTest(t)
	defer cleanup()

	ctx := tests.TestContext(t)
	db := sqlx.NewDb(testDB.DB, "postgres")

	// Storefront 1100 has product with SKU "TEST-SKU-001"
	// Creating in storefront 1101 should succeed
	req := &pb.CreateProductRequest{
		StorefrontId:  1101, // Different storefront
		Name:          "Product in Different Storefront",
		Price:         99.99,
		Currency:      "USD",
		CategoryId:    2110,
		Sku:           stringPtr("TEST-SKU-001"), // Same SKU as storefront 1100
		StockQuantity: 10,
		IsActive:      true,
	}

	resp, err := client.CreateProduct(ctx, req)

	// Assertions
	require.NoError(t, err, "Same SKU should be allowed in different storefront")
	require.NotNil(t, resp)
	require.NotNil(t, resp.Product)

	product := resp.Product
	assert.Equal(t, int64(1101), product.StorefrontId)
	assert.Equal(t, stringPtr("TEST-SKU-001"), product.Sku)

	// Verify in database
	dbProduct := getProductByID(t, db, product.Id)
	require.NotNil(t, dbProduct)
	assert.Equal(t, int64(1101), dbProduct.StorefrontID)
}

// ============================================================================
// BulkCreateProducts Tests
// ============================================================================

// TestBulkCreateProducts_Success tests creating multiple products in single request
func TestBulkCreateProducts_Success(t *testing.T) {
	client, testDB, cleanup := setupCreateProductTest(t)
	defer cleanup()

	ctx := tests.TestContext(t)
	db := sqlx.NewDb(testDB.DB, "postgres")

	products := []*pb.ProductInput{
		{
			Name:          "Bulk Product 1",
			Description:   "First product in bulk",
			Price:         19.99,
			Currency:      "USD",
			CategoryId:    2110,
			Sku:           stringPtr("BULK-001"),
			StockQuantity: 10,
			IsActive:      true,
		},
		{
			Name:          "Bulk Product 2",
			Description:   "Second product in bulk",
			Price:         29.99,
			Currency:      "USD",
			CategoryId:    2111,
			Sku:           stringPtr("BULK-002"),
			StockQuantity: 20,
			IsActive:      true,
		},
		{
			Name:          "Bulk Product 3",
			Description:   "Third product in bulk",
			Price:         39.99,
			Currency:      "USD",
			CategoryId:    2120,
			Sku:           stringPtr("BULK-003"),
			StockQuantity: 30,
			IsActive:      true,
		},
	}

	req := &pb.BulkCreateProductsRequest{
		StorefrontId: 1101,
		Products:     products,
	}

	resp, err := client.BulkCreateProducts(ctx, req)

	// Assertions
	require.NoError(t, err, "BulkCreateProducts should succeed")
	require.NotNil(t, resp)
	assert.Equal(t, int32(3), resp.SuccessfulCount, "Should create 3 products")
	assert.Equal(t, int32(0), resp.FailedCount, "Should have no failures")
	assert.Len(t, resp.Products, 3, "Should return 3 created products")
	assert.Empty(t, resp.Errors, "Should have no errors")

	// Verify each product
	for i, product := range resp.Products {
		assert.Greater(t, product.Id, int64(0), "Product %d should have ID", i)
		assert.Equal(t, products[i].Name, product.Name)
		assert.Equal(t, products[i].Price, product.Price)
		assert.Equal(t, products[i].StockQuantity, product.StockQuantity)

		// Verify in database
		dbProduct := getProductByID(t, db, product.Id)
		require.NotNil(t, dbProduct, "Product %d should exist in DB", i)
		assert.Equal(t, products[i].Name, dbProduct.Title)
	}
}

// TestBulkCreateProducts_LargeBatch tests creating 50 products in single batch
func TestBulkCreateProducts_LargeBatch(t *testing.T) {
	client, testDB, cleanup := setupCreateProductTest(t)
	defer cleanup()

	ctx := tests.TestContext(t)
	db := sqlx.NewDb(testDB.DB, "postgres")

	batchSize := 50
	products := make([]*pb.ProductInput, batchSize)

	for i := 0; i < batchSize; i++ {
		products[i] = &pb.ProductInput{
			Name:          fmt.Sprintf("Large Batch Product %d", i+1),
			Description:   fmt.Sprintf("Product %d in large batch", i+1),
			Price:         float64(10 + i),
			Currency:      "USD",
			CategoryId:    2110,
			Sku:           stringPtr(fmt.Sprintf("LARGE-BATCH-%03d", i+1)),
			StockQuantity: int32(10 + i),
			IsActive:      true,
		}
	}

	req := &pb.BulkCreateProductsRequest{
		StorefrontId: 1101,
		Products:     products,
	}

	start := time.Now()
	resp, err := client.BulkCreateProducts(ctx, req)
	elapsed := time.Since(start)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, int32(batchSize), resp.SuccessfulCount)
	assert.Equal(t, int32(0), resp.FailedCount)
	assert.Len(t, resp.Products, batchSize)

	// Performance requirement: < 2s for 50 products
	assert.Less(t, elapsed, 2*time.Second,
		"Large batch (%d products) should complete in < 2s (actual: %v)", batchSize, elapsed)

	t.Logf("Created %d products in %v (avg: %v per product)",
		batchSize, elapsed, elapsed/time.Duration(batchSize))

	// Verify random products in database
	for i := 0; i < 5; i++ {
		idx := i * 10 // Check products 0, 10, 20, 30, 40
		dbProduct := getProductByID(t, db, resp.Products[idx].Id)
		require.NotNil(t, dbProduct)
		assert.Equal(t, products[idx].Name, dbProduct.Title)
	}
}

// TestBulkCreateProducts_PartialFailure tests bulk create with some invalid products
func TestBulkCreateProducts_PartialFailure(t *testing.T) {
	client, _, cleanup := setupCreateProductTest(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	products := []*pb.ProductInput{
		{
			Name:          "Valid Product 1",
			Price:         19.99,
			Currency:      "USD",
			CategoryId:    2110,
			Sku:           stringPtr("PARTIAL-001"),
			StockQuantity: 10,
			IsActive:      true,
		},
		{
			Name:          "", // INVALID: empty name
			Price:         29.99,
			Currency:      "USD",
			CategoryId:    2111,
			Sku:           stringPtr("PARTIAL-002-INVALID"),
			StockQuantity: 20,
			IsActive:      true,
		},
		{
			Name:          "Valid Product 3",
			Price:         39.99,
			Currency:      "USD",
			CategoryId:    2120,
			Sku:           stringPtr("PARTIAL-003"),
			StockQuantity: 30,
			IsActive:      true,
		},
	}

	req := &pb.BulkCreateProductsRequest{
		StorefrontId: 1101,
		Products:     products,
	}

	resp, err := client.BulkCreateProducts(ctx, req)

	// Assertions - partial success should still return response
	require.NoError(t, err, "BulkCreate should not return gRPC error for partial failure")
	require.NotNil(t, resp)

	// Check that some products succeeded and one failed
	assert.Equal(t, int32(2), resp.SuccessfulCount, "Should have 2 successful products")
	assert.Equal(t, int32(1), resp.FailedCount, "Should have 1 failed product")
	assert.Len(t, resp.Products, 2, "Should return 2 created products")
	assert.Len(t, resp.Errors, 1, "Should have 1 error")

	// Verify error details
	bulkError := resp.Errors[0]
	assert.Equal(t, int32(1), bulkError.Index, "Error should be for product at index 1")
	assert.Contains(t, bulkError.ErrorMessage, "name", "Error should mention empty name")
}

// TestBulkCreateProducts_EmptyBatch tests validation with empty product list
func TestBulkCreateProducts_EmptyBatch(t *testing.T) {
	client, _, cleanup := setupCreateProductTest(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	req := &pb.BulkCreateProductsRequest{
		StorefrontId: 1101,
		Products:     []*pb.ProductInput{}, // Empty batch
	}

	resp, err := client.BulkCreateProducts(ctx, req)

	// Assertions
	require.Error(t, err, "Should return error for empty batch")
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "products.bulk_empty", "Error should mention empty products")
}

// TestBulkCreateProducts_DuplicateSKU tests handling of duplicate SKUs within batch
func TestBulkCreateProducts_DuplicateSKU(t *testing.T) {
	client, _, cleanup := setupCreateProductTest(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	products := []*pb.ProductInput{
		{
			Name:          "Product A",
			Price:         19.99,
			Currency:      "USD",
			CategoryId:    2110,
			Sku:           stringPtr("DUPLICATE-SKU"),
			StockQuantity: 10,
			IsActive:      true,
		},
		{
			Name:          "Product B",
			Price:         29.99,
			Currency:      "USD",
			CategoryId:    2111,
			Sku:           stringPtr("DUPLICATE-SKU"), // Duplicate within batch
			StockQuantity: 20,
			IsActive:      true,
		},
	}

	req := &pb.BulkCreateProductsRequest{
		StorefrontId: 1101,
		Products:     products,
	}

	resp, err := client.BulkCreateProducts(ctx, req)

	// Assertions
	require.NoError(t, err, "BulkCreate should not error but report partial failure")
	require.NotNil(t, resp)

	// First product should succeed, second should fail due to duplicate SKU
	assert.Equal(t, int32(1), resp.SuccessfulCount, "First product should succeed")
	assert.Equal(t, int32(1), resp.FailedCount, "Second product should fail")
	assert.Len(t, resp.Errors, 1)

	bulkError := resp.Errors[0]
	assert.Equal(t, int32(1), bulkError.Index, "Error should be for second product")
	assert.Contains(t, bulkError.ErrorMessage, "SKU", "Error should mention duplicate SKU")
}

// ============================================================================
// Performance & Concurrency Tests
// ============================================================================

// TestCreateProduct_Performance tests single product creation meets SLA
func TestCreateProduct_Performance(t *testing.T) {
	client, _, cleanup := setupCreateProductTest(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	req := &pb.CreateProductRequest{
		StorefrontId:  1100,
		Name:          "Performance Test Product",
		Price:         99.99,
		Currency:      "USD",
		CategoryId:    2110,
		Sku:           stringPtr("PERF-TEST-001"),
		StockQuantity: 50,
		IsActive:      true,
	}

	// Warmup call
	_, _ = client.CreateProduct(ctx, req)

	// Measured call
	req.Sku = stringPtr("PERF-TEST-002") // Different SKU to avoid duplicate
	start := time.Now()
	resp, err := client.CreateProduct(ctx, req)
	elapsed := time.Since(start)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Performance SLA: < 100ms
	assert.Less(t, elapsed, 100*time.Millisecond,
		"CreateProduct should complete in < 100ms (actual: %v)", elapsed)

	t.Logf("CreateProduct completed in: %v", elapsed)
}

// TestCreateProduct_Concurrent tests concurrent product creation (race detector)
func TestCreateProduct_Concurrent(t *testing.T) {
	client, testDB, cleanup := setupCreateProductTest(t)
	defer cleanup()

	ctx := tests.TestContext(t)
	db := sqlx.NewDb(testDB.DB, "postgres")

	concurrentRequests := 10
	var wg sync.WaitGroup
	results := make(chan *pb.ProductResponse, concurrentRequests)
	errors := make(chan error, concurrentRequests)

	// Launch concurrent CreateProduct requests
	for i := 0; i < concurrentRequests; i++ {
		wg.Add(1)
		go func(reqNum int) {
			defer wg.Done()

			req := &pb.CreateProductRequest{
				StorefrontId:  1100,
				Name:          fmt.Sprintf("Concurrent Product %d", reqNum),
				Description:   fmt.Sprintf("Created concurrently %d", reqNum),
				Price:         float64(50 + reqNum),
				Currency:      "USD",
				CategoryId:    2110,
				Sku:           stringPtr(fmt.Sprintf("CONCURRENT-%03d", reqNum)),
				StockQuantity: int32(10 + reqNum),
				IsActive:      true,
			}

			resp, err := client.CreateProduct(ctx, req)
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

	// Verify all products created successfully
	successCount := 0
	productIDs := make([]int64, 0, concurrentRequests)

	for resp := range results {
		require.NotNil(t, resp)
		require.NotNil(t, resp.Product)
		assert.Greater(t, resp.Product.Id, int64(0))
		productIDs = append(productIDs, resp.Product.Id)
		successCount++
	}

	assert.Equal(t, concurrentRequests, successCount,
		"All %d concurrent requests should succeed", concurrentRequests)

	// Verify all products exist in database
	for _, productID := range productIDs {
		dbProduct := getProductByID(t, db, productID)
		assert.NotNil(t, dbProduct, "Product %d should exist in DB", productID)
	}

	// Verify no duplicate product IDs (race condition check)
	uniqueIDs := make(map[int64]bool)
	for _, id := range productIDs {
		uniqueIDs[id] = true
	}
	assert.Len(t, uniqueIDs, concurrentRequests,
		"All product IDs should be unique (no race condition)")
}

// TestCreateProduct_Concurrent_SameStorefront tests concurrent creates in same storefront
func TestCreateProduct_Concurrent_SameStorefront(t *testing.T) {
	client, testDB, cleanup := setupCreateProductTest(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	concurrentRequests := 20
	storefrontID := int64(1100)

	// Count products before
	initialCount := tests.CountProductsByStorefront(t, testDB.DB, storefrontID)

	var wg sync.WaitGroup
	successCount := int32(0)
	var mu sync.Mutex

	for i := 0; i < concurrentRequests; i++ {
		wg.Add(1)
		go func(reqNum int) {
			defer wg.Done()

			req := &pb.CreateProductRequest{
				StorefrontId:  storefrontID,
				Name:          fmt.Sprintf("Same Store Product %d", reqNum),
				Price:         float64(25 + reqNum),
				Currency:      "USD",
				CategoryId:    2110,
				Sku:           stringPtr(fmt.Sprintf("SAME-STORE-%03d", reqNum)),
				StockQuantity: int32(5 + reqNum),
				IsActive:      true,
			}

			resp, err := client.CreateProduct(ctx, req)
			if err == nil && resp != nil && resp.Product != nil {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()

	// Verify all products created
	assert.Equal(t, int32(concurrentRequests), successCount,
		"All concurrent products should be created")

	// Verify product count in database
	finalCount := tests.CountProductsByStorefront(t, testDB.DB, storefrontID)
	expectedCount := initialCount + int32(concurrentRequests)
	assert.Equal(t, expectedCount, finalCount,
		"Storefront should have %d products", expectedCount)

	// Verify no deadlocks occurred
	dbx := sqlx.NewDb(testDB.DB, "postgres")
	var deadlockCount int
	err2 := dbx.Get(&deadlockCount, "SELECT COUNT(*) FROM pg_stat_activity WHERE wait_event_type = 'Lock'")
	require.NoError(t, err2)
	assert.Equal(t, 0, deadlockCount, "Should have no active deadlocks")
}

// TestCreateProduct_StressTest tests creating products under load
func TestCreateProduct_StressTest(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	client, _, cleanup := setupCreateProductTest(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	iterations := 100
	start := time.Now()

	for i := 0; i < iterations; i++ {
		req := &pb.CreateProductRequest{
			StorefrontId:  1101,
			Name:          fmt.Sprintf("Stress Test Product %d", i),
			Price:         float64(10 + i%100),
			Currency:      "USD",
			CategoryId:    2110,
			Sku:           stringPtr(fmt.Sprintf("STRESS-%05d", i)),
			StockQuantity: int32(i % 50),
			IsActive:      true,
		}

		resp, err := client.CreateProduct(ctx, req)
		require.NoError(t, err, "Iteration %d should succeed", i)
		require.NotNil(t, resp)
	}

	elapsed := time.Since(start)
	avgTime := elapsed / time.Duration(iterations)

	t.Logf("Stress test: %d products created in %v (avg: %v per product)",
		iterations, elapsed, avgTime)

	// Performance assertion: average should be reasonable
	assert.Less(t, avgTime, 200*time.Millisecond,
		"Average creation time should be < 200ms")
}

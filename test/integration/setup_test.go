package integration

import (
	"context"
	"net"
	"sync"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	pb "github.com/vondi-global/listings/api/proto/listings/v1"
	"github.com/vondi-global/listings/internal/domain"
	"github.com/vondi-global/listings/internal/metrics"
	miniorepo "github.com/vondi-global/listings/internal/repository/minio"
	"github.com/vondi-global/listings/internal/repository/postgres"
	"github.com/vondi-global/listings/internal/service"
	"github.com/vondi-global/listings/internal/service/listings"
	testutils "github.com/vondi-global/listings/internal/testing"
	grpchandlers "github.com/vondi-global/listings/internal/transport/grpc"
)

// =============================================================================
// Constants
// =============================================================================

const (
	// bufSize is the size of the in-memory buffer for bufconn
	bufSize = 1024 * 1024 // 1MB

	// defaultTimeout is the default timeout for test contexts
)

// =============================================================================
// Mock OrderService for Integration Tests
// =============================================================================

// mockOrderService is a minimal implementation of OrderService for testing
type mockOrderService struct{}

func (m *mockOrderService) CreateOrder(ctx context.Context, req *service.CreateOrderRequest) (*domain.Order, error) {
	return nil, nil
}

func (m *mockOrderService) GetOrder(ctx context.Context, orderID int64) (*domain.Order, error) {
	return nil, nil
}

func (m *mockOrderService) GetOrderByNumber(ctx context.Context, orderNumber string) (*domain.Order, error) {
	return nil, nil
}

func (m *mockOrderService) ListOrders(ctx context.Context, req *service.ListOrdersRequest) ([]*domain.Order, int64, error) {
	return nil, 0, nil
}

func (m *mockOrderService) CancelOrder(ctx context.Context, orderID int64, userID int64, reason string) (*domain.Order, error) {
	return nil, nil
}

func (m *mockOrderService) UpdateOrderStatus(ctx context.Context, orderID int64, status domain.OrderStatus) (*domain.Order, error) {
	return nil, nil
}

func (m *mockOrderService) GetOrderStats(ctx context.Context, userID *int64, storefrontID *int64) (*service.OrderStats, error) {
	return nil, nil
}

func (m *mockOrderService) ConfirmOrderPayment(ctx context.Context, orderID int64, transactionID string) error {
	return nil
}

func (m *mockOrderService) ProcessRefund(ctx context.Context, orderID int64) error {
	return nil
}

// =============================================================================
// Mock Cart Service
// =============================================================================

type mockCartService struct{}

func (m *mockCartService) AddToCart(ctx context.Context, req *service.AddToCartRequest) (*domain.Cart, error) {
	return nil, nil
}

func (m *mockCartService) UpdateCartItem(ctx context.Context, req *service.UpdateCartItemRequest) (*domain.Cart, error) {
	return nil, nil
}

func (m *mockCartService) UpdateCartItemByItemID(ctx context.Context, cartItemID int64, quantity int32, userID *int64, sessionID *string) (*domain.Cart, error) {
	return nil, nil
}

func (m *mockCartService) RemoveFromCart(ctx context.Context, cartID, itemID int64) error {
	return nil
}

func (m *mockCartService) RemoveFromCartByItemID(ctx context.Context, cartItemID int64, userID *int64, sessionID *string) error {
	return nil
}

func (m *mockCartService) GetCart(ctx context.Context, userID *int64, sessionID *string, storefrontID int64) (*domain.Cart, error) {
	return nil, nil
}

func (m *mockCartService) ClearCart(ctx context.Context, cartID int64) error {
	return nil
}

func (m *mockCartService) GetUserCarts(ctx context.Context, userID int64) ([]*domain.Cart, error) {
	return nil, nil
}

func (m *mockCartService) MergeSessionCartToUser(ctx context.Context, sessionID string, userID int64) error {
	return nil
}

func (m *mockCartService) RecalculateCart(ctx context.Context, cartID int64) (*domain.Cart, error) {
	return nil, nil
}

func (m *mockCartService) ValidateCartItems(ctx context.Context, cartID int64) ([]service.PriceChangeItem, error) {
	return nil, nil
}

// =============================================================================
// Global Metrics Singleton
// =============================================================================

var (
	// testMetrics is a singleton instance of metrics for all integration tests.
	// This prevents duplicate Prometheus registration errors across tests.
	testMetrics     *metrics.Metrics
	testMetricsOnce sync.Once
)

// getTestMetrics returns the singleton metrics instance for testing.
// This is thread-safe and ensures only one metrics instance exists.
func getTestMetrics() *metrics.Metrics {
	testMetricsOnce.Do(func() {
		testMetrics = metrics.NewMetrics("listings_integration_test")
	})
	return testMetrics
}

// =============================================================================
// Test Server Setup
// =============================================================================

// TestServerConfig holds configuration for creating a test server
type TestServerConfig struct {
	// UseDocker determines whether to use Docker for database (true) or existing DB (false)
	UseDocker bool

	// DSN is the connection string (used when UseDocker is false)
	DSN string

	// MigrationsPath is the path to migration files
	MigrationsPath string

	// FixturesPath is the path to fixture files (optional)
	FixturesPath string

	// UseTransactions enables transactional test isolation
	UseTransactions bool

	// Logger is the logger to use (if nil, a test logger is created)
	Logger *zerolog.Logger
}

// DefaultTestServerConfig returns a configuration with sensible defaults
func DefaultTestServerConfig() TestServerConfig {
	return TestServerConfig{
		UseDocker:       true,
		MigrationsPath:  "../../migrations",
		UseTransactions: false,
		Logger:          nil,
	}
}

// TestServer wraps a gRPC test server with all necessary components
type TestServer struct {
	Client   pb.ListingsServiceClient
	DB       *testutils.TestDatabase
	Conn     *grpc.ClientConn
	Server   *grpc.Server
	Listener *bufconn.Listener
	Logger   zerolog.Logger

	// cleanup functions to run when server is torn down
	cleanupFuncs []func() error
}

// SetupTestServer creates a complete test server with database, gRPC server, and client.
// This is the recommended way to set up integration tests.
//
// Example:
//
//	func TestMyFeature(t *testing.T) {
//	    testutils.SkipIfShort(t)
//	    testutils.SkipIfNoDocker(t)
//
//	    config := DefaultTestServerConfig()
//	    server := SetupTestServer(t, config)
//	    defer server.Teardown(t)
//
//	    ctx := testutils.TestContext(t)
//	    resp, err := server.Client.GetListing(ctx, &pb.GetListingRequest{Id: 1})
//	    require.NoError(t, err)
//	}
func SetupTestServer(t *testing.T, config TestServerConfig) *TestServer {
	t.Helper()

	// Setup test database
	dbConfig := testutils.TestDatabaseConfig{
		UseDocker:       config.UseDocker,
		DSN:             config.DSN,
		MigrationsPath:  config.MigrationsPath,
		FixturesPath:    config.FixturesPath,
		UseTransactions: config.UseTransactions,
	}
	testDB := testutils.SetupTestDatabase(t, dbConfig)

	// Setup logger
	var logger zerolog.Logger
	if config.Logger != nil {
		logger = *config.Logger
	} else {
		logger = zerolog.New(zerolog.NewTestWriter(t)).With().Timestamp().Logger()
	}

	// Create sqlx.DB wrapper (if not already wrapped)
	db := testDB.GetDBx()

	// Create repository
	repo := postgres.NewRepository(db, logger)

	// Create service (no Redis and no indexer for integration tests)
	listingsService := listings.NewService(repo, nil, nil, logger)

	// Create storefront service
	storefrontService := listings.NewStorefrontService(repo, &logger)

	// Create attribute service (uses dedicated AttributeRepository)
	attrRepo := postgres.NewAttributeRepository(db, logger)
	attrService := service.NewAttributeService(attrRepo, nil, logger)

	// Create category service (uses main Repository which implements CategoryRepository interface)
	categoryService := service.NewCategoryService(repo, nil, logger)

	// Create mock order service for integration tests
	orderService := &mockOrderService{}

	// Create mock cart service for integration tests
	cartService := &mockCartService{}

	// Mock analytics service (nil is OK for integration tests that don't need analytics)
	var analyticsService service.AnalyticsService = nil

	// Mock minio client (nil is OK for integration tests that don't need image operations)
	var minioClient *miniorepo.Client = nil

	// Get singleton metrics instance
	m := getTestMetrics()

	// Create gRPC server
	grpcServer := grpchandlers.NewServer(
		listingsService,
		storefrontService,
		attrService,
		categoryService,
		orderService,
		cartService,
		analyticsService,
		minioClient,
		m,
		logger,
	)

	// Setup in-memory gRPC connection using bufconn
	listener := bufconn.Listen(bufSize)

	server := grpc.NewServer()
	pb.RegisterListingsServiceServer(server, grpcServer)

	// Start gRPC server in background
	go func() {
		if err := server.Serve(listener); err != nil {
			logger.Error().Err(err).Msg("gRPC server failed")
		}
	}()

	// Create client connection
	conn, err := grpc.DialContext(
		context.Background(),
		"bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err, "Failed to create gRPC client connection")

	client := pb.NewListingsServiceClient(conn)

	return &TestServer{
		Client:       client,
		DB:           testDB,
		Conn:         conn,
		Server:       server,
		Listener:     listener,
		Logger:       logger,
		cleanupFuncs: []func() error{},
	}
}

// AddCleanupFunc registers a cleanup function to be called during Teardown.
// This is useful for custom cleanup logic.
func (ts *TestServer) AddCleanupFunc(fn func() error) {
	ts.cleanupFuncs = append(ts.cleanupFuncs, fn)
}

// Teardown cleans up the test server.
// This should be called in a defer statement after creating the test server.
func (ts *TestServer) Teardown(t *testing.T) {
	t.Helper()

	// Run custom cleanup functions
	for i, fn := range ts.cleanupFuncs {
		if err := fn(); err != nil {
			t.Logf("Warning: cleanup function %d failed: %v", i, err)
		}
	}

	// Close client connection
	if ts.Conn != nil {
		if err := ts.Conn.Close(); err != nil {
			t.Logf("Warning: failed to close client connection: %v", err)
		}
	}

	// Stop gRPC server
	if ts.Server != nil {
		ts.Server.GracefulStop()
	}

	// Close listener
	if ts.Listener != nil {
		if err := ts.Listener.Close(); err != nil {
			t.Logf("Warning: failed to close listener: %v", err)
		}
	}

	// Teardown database
	if ts.DB != nil {
		ts.DB.Teardown(t)
	}
}

// =============================================================================
// Parallel Test Server Pool
// =============================================================================

// TestServerPool manages multiple test servers for parallel testing.
// This is useful for load testing or simulating multiple concurrent clients.
type TestServerPool struct {
	servers []*TestServer
	size    int
	logger  zerolog.Logger
}

// NewTestServerPool creates a pool of test servers for parallel testing.
//
// Example:
//
//	pool := NewTestServerPool(t, 5, DefaultTestServerConfig())
//	defer pool.TeardownAll(t)
//
//	// Run tests in parallel
//	t.Run("Parallel", func(t *testing.T) {
//	    for i := 0; i < pool.Size(); i++ {
//	        serverIndex := i
//	        t.Run(fmt.Sprintf("Server%d", serverIndex), func(t *testing.T) {
//	            t.Parallel()
//	            server := pool.Get(serverIndex)
//	            // Use server for testing
//	        })
//	    }
//	})
func NewTestServerPool(t *testing.T, size int, config TestServerConfig) *TestServerPool {
	t.Helper()

	if size <= 0 {
		t.Fatalf("pool size must be positive, got %d", size)
	}

	var logger zerolog.Logger
	if config.Logger != nil {
		logger = *config.Logger
	} else {
		logger = zerolog.New(zerolog.NewTestWriter(t)).With().Timestamp().Logger()
	}

	servers := make([]*TestServer, size)
	for i := 0; i < size; i++ {
		servers[i] = SetupTestServer(t, config)
	}

	return &TestServerPool{
		servers: servers,
		size:    size,
		logger:  logger,
	}
}

// Get returns the test server at the specified index
func (p *TestServerPool) Get(index int) *TestServer {
	if index < 0 || index >= p.size {
		p.logger.Warn().Int("index", index).Int("size", p.size).Msg("index out of range")
		return nil
	}
	return p.servers[index]
}

// GetAll returns all test servers in the pool
func (p *TestServerPool) GetAll() []*TestServer {
	return p.servers
}

// Size returns the number of test servers in the pool
func (p *TestServerPool) Size() int {
	return p.size
}

// TeardownAll tears down all test servers in the pool
func (p *TestServerPool) TeardownAll(t *testing.T) {
	t.Helper()

	for i, server := range p.servers {
		t.Logf("Tearing down test server %d/%d", i+1, p.size)
		server.Teardown(t)
	}
}

// =============================================================================
// Helper Functions
// =============================================================================

// ExecuteSQL is a helper to execute SQL directly on the test database.
// This is useful for complex test setup that can't be done through the API.
//
// Example:
//
//	ExecuteSQL(t, server, "INSERT INTO categories (name, slug) VALUES ($1, $2)", "Electronics", "electronics")
func ExecuteSQL(t *testing.T, server *TestServer, query string, args ...interface{}) {
	t.Helper()
	server.DB.ExecuteSQL(t, query, args...)
}

// CountRows is a helper to count rows in a table matching a condition.
//
// Example:
//
//	count := CountRows(t, server, "listings", "status = $1", "active")
//	require.Equal(t, 5, count)
func CountRows(t *testing.T, server *TestServer, table string, condition string, args ...interface{}) int {
	t.Helper()
	return server.DB.CountRows(t, table, condition, args...)
}

// RowExists is a helper to check if a row exists in a table.
//
// Example:
//
//	exists := RowExists(t, server, "categories", "id = $1", 1)
//	require.True(t, exists)
func RowExists(t *testing.T, server *TestServer, table string, condition string, args ...interface{}) bool {
	t.Helper()
	return server.DB.RowExists(t, table, condition, args...)
}

// TruncateTables is a helper to truncate tables for test isolation.
//
// Example:
//
//	TruncateTables(t, server, "listings", "categories", "images")
func TruncateTables(t *testing.T, server *TestServer, tables ...string) {
	t.Helper()
	server.DB.TruncateTables(t, tables...)
}

// CleanupTestData is a helper to remove test data based on ID ranges.
//
// Example:
//
//	CleanupTestData(t, server, "listings", 1000, 2000)
func CleanupTestData(t *testing.T, server *TestServer, table string, minID, maxID int64) {
	t.Helper()
	server.DB.CleanupTestData(t, table, minID, maxID)
}

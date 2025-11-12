//go:build integration

package integration

import (
	"context"
	"net"
	"sync"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	pb "github.com/sveturs/listings/api/proto/listings/v1"
	"github.com/sveturs/listings/internal/metrics"
	"github.com/sveturs/listings/internal/repository/postgres"
	"github.com/sveturs/listings/internal/service/listings"
	grpchandlers "github.com/sveturs/listings/internal/transport/grpc"
	"github.com/sveturs/listings/tests"
)

// ============================================================================
// Shared Test Constants
// ============================================================================

const bufSize = 1024 * 1024

// ============================================================================
// Shared Test Helpers
// ============================================================================

// stringPtr returns a pointer to a string
func stringPtr(s string) *string {
	return &s
}

// int32Ptr returns a pointer to an int32
func int32Ptr(i int32) *int32 {
	return &i
}

// int64Ptr returns a pointer to an int64

// float64Ptr returns a pointer to a float64
func float64Ptr(f float64) *float64 {
	return &f
}

// boolPtr returns a pointer to a bool
func boolPtr(b bool) *bool {
	return &b
}

// stringRepeat repeats a string n times
func stringRepeat(s string, count int) string {
	if count <= 0 {
		return ""
	}
	result := make([]byte, len(s)*count)
	bp := copy(result, s)
	for bp < len(result) {
		copy(result[bp:], result[:bp])
		bp *= 2
	}
	return string(result)
}

// ============================================================================
// Metrics Singleton
// ============================================================================

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

// ============================================================================
// Test Server Setup
// ============================================================================

// setupGRPCTestServer creates a gRPC server with real database and inventory fixtures
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

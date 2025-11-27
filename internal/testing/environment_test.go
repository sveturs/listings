package testing_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	testingpkg "github.com/vondi-global/listings/internal/testing"
)

// TestNewTestEnvironment_Docker tests environment creation with Docker containers
func TestNewTestEnvironment_Docker(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Docker test in short mode")
	}

	// Create test environment
	env := testingpkg.NewTestEnvironment(t)
	defer env.Cleanup()

	// Verify PostgreSQL is connected
	assert.NotNil(t, env.DB)
	err := env.DB.Ping()
	assert.NoError(t, err, "PostgreSQL should be connected")

	// Verify Redis is connected
	assert.NotNil(t, env.Redis)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = env.Redis.Ping(ctx).Err()
	assert.NoError(t, err, "Redis should be connected")

	// Verify repository is initialized
	assert.NotNil(t, env.Repo)

	// Verify services are initialized
	assert.NotNil(t, env.CartService)
	assert.NotNil(t, env.OrderService)
	assert.NotNil(t, env.InventoryService)

	// Verify logger is initialized
	assert.NotNil(t, env.Logger)
}

// TestNewTestEnvironment_CustomConfig tests environment with custom configuration
func TestNewTestEnvironment_CustomConfig(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Docker test in short mode")
	}

	config := testingpkg.DefaultTestEnvironmentConfig()
	config.LogLevel = "debug"
	config.MaxWaitTime = 90 * time.Second

	env := testingpkg.NewTestEnvironmentWithConfig(t, config)
	defer env.Cleanup()

	assert.NotNil(t, env.DB)
	assert.NotNil(t, env.Redis)
	assert.Equal(t, "debug", env.Config.LogLevel)
}

// TestSeedTestData tests seeding test data
func TestSeedTestData(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Docker test in short mode")
	}

	env := testingpkg.NewTestEnvironment(t)
	defer env.Cleanup()

	// Seed test data
	env.SeedTestData(t)

	// Verify categories exist
	var count int
	err := env.DB.Get(&count, "SELECT COUNT(*) FROM categories WHERE id IN (1000, 1001)")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, 2, "Should have at least 2 test categories")

	// Verify storefronts exist
	err = env.DB.Get(&count, "SELECT COUNT(*) FROM storefronts WHERE id IN (1, 2)")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, 2, "Should have at least 2 test storefronts")

	// Verify products exist
	err = env.DB.Get(&count, "SELECT COUNT(*) FROM listings WHERE id IN (100, 101, 200)")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, 3, "Should have at least 3 test products")

	// Verify inventory exists
	err = env.DB.Get(&count, "SELECT COUNT(*) FROM inventory WHERE listing_id IN (100, 101, 200)")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, 3, "Should have at least 3 inventory records")
}

// TestTruncateTables tests table truncation
func TestTruncateTables(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Docker test in short mode")
	}

	env := testingpkg.NewTestEnvironment(t)
	defer env.Cleanup()

	// Seed data
	env.SeedTestData(t)

	// Verify data exists
	var count int
	err := env.DB.Get(&count, "SELECT COUNT(*) FROM listings WHERE id IN (100, 101, 200)")
	require.NoError(t, err)
	assert.Greater(t, count, 0)

	// Truncate
	env.TruncateTables(t, "listings")

	// Verify data is gone
	err = env.DB.Get(&count, "SELECT COUNT(*) FROM listings WHERE id IN (100, 101, 200)")
	require.NoError(t, err)
	assert.Equal(t, 0, count, "Products should be truncated")
}

// TestFlushRedis tests Redis flushing
func TestFlushRedis(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Docker test in short mode")
	}

	env := testingpkg.NewTestEnvironment(t)
	defer env.Cleanup()

	ctx := context.Background()

	// Set some data
	err := env.Redis.Set(ctx, "test-key", "test-value", 0).Err()
	require.NoError(t, err)

	// Verify it exists
	val, err := env.Redis.Get(ctx, "test-key").Result()
	require.NoError(t, err)
	assert.Equal(t, "test-value", val)

	// Flush Redis
	env.FlushRedis(t)

	// Verify data is gone
	_, err = env.Redis.Get(ctx, "test-key").Result()
	assert.Error(t, err, "Key should not exist after flush")
}

// TestCreateTestProduct tests product creation helper
func TestCreateTestProduct(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Docker test in short mode")
	}

	env := testingpkg.NewTestEnvironment(t)
	defer env.Cleanup()

	// Seed base data
	env.SeedTestData(t)

	// Create test product
	productID := env.CreateTestProduct(t, 1, 299.99)

	// Verify product was created
	assert.NotZero(t, productID)

	// Verify in database
	var count int
	err := env.DB.Get(&count, "SELECT COUNT(*) FROM listings WHERE id = $1", productID)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

// TestAddCleanupFunc tests custom cleanup function registration
func TestAddCleanupFunc(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Docker test in short mode")
	}

	cleanupCalled := false

	env := testingpkg.NewTestEnvironment(t)

	// Add custom cleanup
	env.AddCleanupFunc(func() error {
		cleanupCalled = true
		return nil
	})

	// Cleanup
	env.Cleanup()

	// Verify cleanup was called
	assert.True(t, cleanupCalled, "Custom cleanup function should be called")
}

// TestEnvironmentIsolation tests that multiple test environments are isolated
func TestEnvironmentIsolation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Docker test in short mode")
	}

	// Create first environment
	env1 := testingpkg.NewTestEnvironment(t)
	defer env1.Cleanup()

	env1.SeedTestData(t)

	// Insert data in env1
	_, err := env1.DB.Exec("INSERT INTO listings (id, user_id, storefront_id, title, slug, price, currency, category_id, status) VALUES (9001, 1, 1, 'Env1 Product', 'env1-product', 99.99, 'EUR', 1001, 'active')")
	require.NoError(t, err)

	// Create second environment (different container)
	env2 := testingpkg.NewTestEnvironment(t)
	defer env2.Cleanup()

	env2.SeedTestData(t)

	// Verify env2 doesn't have env1's data
	var count int
	err = env2.DB.Get(&count, "SELECT COUNT(*) FROM listings WHERE id = 9001")
	require.NoError(t, err)
	assert.Equal(t, 0, count, "Second environment should not have first environment's data")

	// Verify env1 still has its data
	err = env1.DB.Get(&count, "SELECT COUNT(*) FROM listings WHERE id = 9001")
	require.NoError(t, err)
	assert.Equal(t, 1, count, "First environment should still have its data")
}

// TestMigrationsApplied tests that migrations are properly applied
func TestMigrationsApplied(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Docker test in short mode")
	}

	env := testingpkg.NewTestEnvironment(t)
	defer env.Cleanup()

	// Check that expected tables exist
	expectedTables := []string{
		"listings",
		"inventory_reservations",
		"shopping_carts",
		"cart_items",
		"orders",
		"order_items",
		"categories",
		"storefronts",
	}

	for _, table := range expectedTables {
		var exists bool
		err := env.DB.Get(&exists, `
			SELECT EXISTS (
				SELECT FROM information_schema.tables
				WHERE table_schema = 'public'
				AND table_name = $1
			)
		`, table)
		require.NoError(t, err)
		assert.True(t, exists, "Table %s should exist after migrations", table)
	}
}

// NOTE: BenchmarkEnvironmentCreation was removed because it benchmarks Docker container
// startup/shutdown time (30s+ per iteration), not code performance. This would cause
// CI/CD timeouts and produce meaningless metrics.
//
// If you need to measure environment setup overhead, use manual timing with a single
// iteration, or benchmark WITHOUT Docker using in-memory databases.

package testing

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/vondi-global/listings/internal/repository/postgres"
	"github.com/vondi-global/listings/internal/service"
)

// TestEnvironment provides a complete testing environment with all dependencies.
// It includes PostgreSQL, Redis, and all services initialized with test data.
type TestEnvironment struct {
	// Infrastructure
	DB       *sqlx.DB
	PgPool   *pgxpool.Pool // pgx connection pool for new repositories
	Redis    *redis.Client
	Logger   zerolog.Logger
	Pool     *dockertest.Pool
	DBRes    *dockertest.Resource
	RedisRes *dockertest.Resource

	// Configuration
	Config TestEnvironmentConfig

	// Repositories
	Repo            *postgres.Repository           // main repository (sqlx-based)
	OrderRepo       postgres.OrderRepository       // order repository (pgxpool-based)
	ReservationRepo postgres.ReservationRepository // reservation repository (pgxpool-based)
	CartRepo        postgres.CartRepository        // cart repository (sqlx-based)

	// Services
	CartService      service.CartService
	OrderService     service.OrderService
	InventoryService service.InventoryService

	// Cleanup
	cleanupFuncs []func() error
}

// TestEnvironmentConfig holds configuration for test environment
type TestEnvironmentConfig struct {
	// UseDocker determines whether to use Docker containers
	UseDocker bool

	// PostgreSQL settings
	PostgresDSN     string
	PostgresVersion string

	// Redis settings
	RedisAddr    string
	RedisVersion string

	// Test data
	LoadFixtures    bool
	MigrationsPath  string
	FixturesPath    string
	UseTransactions bool

	// Timeouts
	MaxWaitTime time.Duration

	// Logging
	LogLevel string
}

// DefaultTestEnvironmentConfig returns sensible defaults
func DefaultTestEnvironmentConfig() TestEnvironmentConfig {
	return TestEnvironmentConfig{
		UseDocker:       true,
		PostgresVersion: "15-alpine",
		RedisVersion:    "7-alpine",
		LoadFixtures:    true,
		MigrationsPath:  "../../migrations",
		FixturesPath:    "",
		UseTransactions: false,
		MaxWaitTime:     60 * time.Second,
		LogLevel:        "error", // Quiet logs in tests
	}
}

// NewTestEnvironment creates a complete test environment with all dependencies.
// It automatically sets up PostgreSQL, Redis, runs migrations, and initializes all services.
//
// Example usage:
//
//	env := testing.NewTestEnvironment(t)
//	defer env.Cleanup()
//
//	// Use env.CartService, env.OrderService, etc.
func NewTestEnvironment(tb testing.TB) *TestEnvironment {
	return NewTestEnvironmentWithConfig(tb, DefaultTestEnvironmentConfig())
}

// NewTestEnvironmentWithConfig creates a test environment with custom configuration
func NewTestEnvironmentWithConfig(tb testing.TB, config TestEnvironmentConfig) *TestEnvironment {
	tb.Helper()

	env := &TestEnvironment{
		Config:       config,
		cleanupFuncs: []func() error{},
	}

	// Setup logger
	env.setupLogger(tb)

	// Setup Docker pool
	if config.UseDocker {
		env.setupDockerPool(tb)
	}

	// Setup PostgreSQL
	env.setupPostgreSQL(tb)

	// Setup Redis
	env.setupRedis(tb)

	// Run migrations
	if config.MigrationsPath != "" {
		env.runMigrations(tb)
	}

	// Load fixtures
	if config.LoadFixtures && config.FixturesPath != "" {
		env.loadFixtures(tb)
	}

	// Initialize repositories
	env.setupRepositories(tb)

	// Initialize services
	env.setupServices(tb)

	return env
}

// setupLogger initializes the logger for tests
func (env *TestEnvironment) setupLogger(tb testing.TB) {
	tb.Helper()

	// Parse log level
	logLevel, err := zerolog.ParseLevel(env.Config.LogLevel)
	if err != nil {
		logLevel = zerolog.ErrorLevel
	}

	// Create logger (output to test log, not stdout)
	env.Logger = zerolog.New(zerolog.NewTestWriter(tb)).
		Level(logLevel).
		With().
		Timestamp().
		Logger()

	tb.Log("Logger initialized")
}

// setupDockerPool creates Docker pool for running containers
func (env *TestEnvironment) setupDockerPool(tb testing.TB) {
	tb.Helper()

	pool, err := dockertest.NewPool("")
	require.NoError(tb, err, "Failed to connect to Docker")

	err = pool.Client.Ping()
	require.NoError(tb, err, "Failed to ping Docker")

	env.Pool = pool
	tb.Log("Docker pool initialized")
}

// setupPostgreSQL starts PostgreSQL container and connects to it
func (env *TestEnvironment) setupPostgreSQL(tb testing.TB) {
	tb.Helper()

	if env.Config.UseDocker {
		env.setupPostgreSQLDocker(tb)
	} else {
		env.setupPostgreSQLExisting(tb)
	}

	// Verify connection
	err := env.DB.Ping()
	require.NoError(tb, err, "Failed to ping PostgreSQL")

	tb.Log("PostgreSQL connected")
}

// setupPostgreSQLDocker starts PostgreSQL in Docker container
func (env *TestEnvironment) setupPostgreSQLDocker(tb testing.TB) {
	tb.Helper()

	dbName := fmt.Sprintf("test_listings_%d", time.Now().UnixNano())

	resource, err := env.Pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        env.Config.PostgresVersion,
		Env: []string{
			"POSTGRES_USER=test_user",
			"POSTGRES_PASSWORD=test_password",
			fmt.Sprintf("POSTGRES_DB=%s", dbName),
			"listen_addresses=*",
		},
	}, func(hostConfig *docker.HostConfig) {
		hostConfig.AutoRemove = true
		hostConfig.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	require.NoError(tb, err, "Failed to start PostgreSQL container")

	env.DBRes = resource

	// Set expiry to clean up in case of test panic
	err = resource.Expire(180) // 3 minutes
	require.NoError(tb, err, "Failed to set container expiry")

	hostAndPort := resource.GetHostPort("5432/tcp")
	dsn := fmt.Sprintf("postgres://test_user:test_password@%s/%s?sslmode=disable", hostAndPort, dbName)

	// Wait for database to be ready
	env.Pool.MaxWait = env.Config.MaxWaitTime
	err = env.Pool.Retry(func() error {
		db, connErr := sqlx.Connect("postgres", dsn)
		if connErr != nil {
			return connErr
		}
		env.DB = db
		return db.Ping()
	})
	require.NoError(tb, err, "Failed to connect to PostgreSQL container")

	// Also create pgxpool connection for new repositories
	pgPool, err := pgxpool.New(context.Background(), dsn)
	require.NoError(tb, err, "Failed to create pgxpool connection")
	env.PgPool = pgPool

	tb.Logf("PostgreSQL container started: %s", hostAndPort)
}

// setupPostgreSQLExisting connects to existing PostgreSQL instance
func (env *TestEnvironment) setupPostgreSQLExisting(tb testing.TB) {
	tb.Helper()

	require.NotEmpty(tb, env.Config.PostgresDSN, "PostgresDSN is required when UseDocker is false")

	db, err := sqlx.Connect("postgres", env.Config.PostgresDSN)
	require.NoError(tb, err, "Failed to connect to existing PostgreSQL")
	env.DB = db

	// Also create pgxpool connection
	pgPool, err := pgxpool.New(context.Background(), env.Config.PostgresDSN)
	require.NoError(tb, err, "Failed to create pgxpool connection")
	env.PgPool = pgPool

	tb.Log("Connected to existing PostgreSQL")
}

// setupRedis starts Redis container and connects to it
func (env *TestEnvironment) setupRedis(tb testing.TB) {
	tb.Helper()

	if env.Config.UseDocker {
		env.setupRedisDocker(tb)
	} else {
		env.setupRedisExisting(tb)
	}

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := env.Redis.Ping(ctx).Err()
	require.NoError(tb, err, "Failed to ping Redis")

	tb.Log("Redis connected")
}

// setupRedisDocker starts Redis in Docker container
func (env *TestEnvironment) setupRedisDocker(tb testing.TB) {
	tb.Helper()

	resource, err := env.Pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "redis",
		Tag:        env.Config.RedisVersion,
	}, func(hostConfig *docker.HostConfig) {
		hostConfig.AutoRemove = true
		hostConfig.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	require.NoError(tb, err, "Failed to start Redis container")

	env.RedisRes = resource

	// Set expiry
	err = resource.Expire(180)
	require.NoError(tb, err, "Failed to set Redis container expiry")

	hostAndPort := resource.GetHostPort("6379/tcp")

	// Wait for Redis to be ready
	env.Pool.MaxWait = env.Config.MaxWaitTime
	err = env.Pool.Retry(func() error {
		client := redis.NewClient(&redis.Options{
			Addr: hostAndPort,
		})
		env.Redis = client
		return client.Ping(context.Background()).Err()
	})
	require.NoError(tb, err, "Failed to connect to Redis container")

	tb.Logf("Redis container started: %s", hostAndPort)
}

// setupRedisExisting connects to existing Redis instance
func (env *TestEnvironment) setupRedisExisting(tb testing.TB) {
	tb.Helper()

	require.NotEmpty(tb, env.Config.RedisAddr, "RedisAddr is required when UseDocker is false")

	client := redis.NewClient(&redis.Options{
		Addr: env.Config.RedisAddr,
	})

	env.Redis = client
	tb.Log("Connected to existing Redis")
}

// runMigrations runs database migrations
func (env *TestEnvironment) runMigrations(tb testing.TB) {
	tb.Helper()

	files, err := os.ReadDir(env.Config.MigrationsPath)
	require.NoError(tb, err, "Failed to read migrations directory")

	// Filter .up.sql files
	var upFiles []string
	for _, f := range files {
		if !f.IsDir() && len(f.Name()) > 7 && f.Name()[len(f.Name())-7:] == ".up.sql" {
			upFiles = append(upFiles, f.Name())
		}
	}

	tb.Logf("Running %d migrations", len(upFiles))

	for _, fileName := range upFiles {
		path := fmt.Sprintf("%s/%s", env.Config.MigrationsPath, fileName)
		data, err := os.ReadFile(path)
		require.NoError(tb, err, "Failed to read migration: %s", fileName)

		_, err = env.DB.Exec(string(data))
		require.NoError(tb, err, "Failed to run migration: %s", fileName)

		tb.Logf("Applied migration: %s", fileName)
	}
}

// loadFixtures loads test fixtures
func (env *TestEnvironment) loadFixtures(tb testing.TB) {
	tb.Helper()

	if env.Config.FixturesPath == "" {
		return
	}

	data, err := os.ReadFile(env.Config.FixturesPath)
	require.NoError(tb, err, "Failed to read fixtures")

	_, err = env.DB.Exec(string(data))
	require.NoError(tb, err, "Failed to load fixtures")

	tb.Log("Fixtures loaded")
}

// setupRepositories initializes all repositories
func (env *TestEnvironment) setupRepositories(tb testing.TB) {
	tb.Helper()

	// Main repository (sqlx-based)
	env.Repo = postgres.NewRepository(env.DB, env.Logger)

	// Specialized repositories
	env.CartRepo = postgres.NewCartRepository(env.PgPool, env.Logger)
	env.OrderRepo = postgres.NewOrderRepository(env.PgPool, env.Logger)
	env.ReservationRepo = postgres.NewReservationRepository(env.PgPool, env.Logger)

	tb.Log("Repositories initialized")
}

// setupServices initializes all services
func (env *TestEnvironment) setupServices(tb testing.TB) {
	tb.Helper()

	// Inventory service
	env.InventoryService = service.NewInventoryService(
		env.ReservationRepo, // reservationRepo
		env.Repo,            // productsRepo
		env.OrderRepo,       // orderRepo
		env.PgPool,          // pool
		env.Logger,
	)

	// Cart service
	env.CartService = service.NewCartService(
		env.CartRepo, // cartRepo
		env.Repo,     // productsRepo
		env.Repo,     // variantsRepo
		env.Repo,     // inventoryRepo
		env.Logger,
	)

	// Order service
	env.OrderService = service.NewOrderService(
		env.OrderRepo,       // orderRepo
		env.CartRepo,        // cartRepo
		env.ReservationRepo, // reservationRepo
		env.Repo,            // productsRepo
		env.PgPool,          // pool
		nil,                 // config (uses default)
		env.Logger,
	)

	tb.Log("Services initialized")
}

// SeedTestData seeds minimal test data (products, inventory, etc.)
func (env *TestEnvironment) SeedTestData(tb testing.TB) {
	tb.Helper()

	ctx := context.Background()

	// Create test categories
	_, err := env.DB.ExecContext(ctx, `
		INSERT INTO categories (id, name, slug, parent_id, level, created_at)
		VALUES
			(1000, 'Test Electronics', 'test-electronics', NULL, 0, NOW()),
			(1001, 'Test Phones', 'test-phones', 1000, 1, NOW())
		ON CONFLICT (id) DO NOTHING
	`)
	require.NoError(tb, err, "Failed to seed categories")

	// Create test storefronts
	_, err = env.DB.ExecContext(ctx, `
		INSERT INTO storefronts (id, name, slug, user_id, created_at, updated_at)
		VALUES
			(1, 'Test Store 1', 'test-store-1', 1, NOW(), NOW()),
			(2, 'Test Store 2', 'test-store-2', 2, NOW(), NOW())
		ON CONFLICT (id) DO NOTHING
	`)
	require.NoError(tb, err, "Failed to seed storefronts")

	// Create test products (listings)
	_, err = env.DB.ExecContext(ctx, `
		INSERT INTO listings (id, user_id, storefront_id, title, slug, description, price, currency, category_id, status, created_at, updated_at)
		VALUES
			(100, 1, 1, 'Test Product 1', 'test-product-1', 'Description 1', 99.99, 'EUR', 1001, 'active', NOW(), NOW()),
			(101, 1, 1, 'Test Product 2', 'test-product-2', 'Description 2', 149.99, 'EUR', 1001, 'active', NOW(), NOW()),
			(200, 2, 2, 'Test Product 3', 'test-product-3', 'Description 3', 199.99, 'EUR', 1001, 'active', NOW(), NOW())
		ON CONFLICT (id) DO NOTHING
	`)
	require.NoError(tb, err, "Failed to seed products")

	tb.Log("Test data seeded")
}

// TruncateTables truncates specified tables for cleanup
func (env *TestEnvironment) TruncateTables(tb testing.TB, tables ...string) {
	tb.Helper()

	for _, table := range tables {
		_, err := env.DB.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		require.NoError(tb, err, "Failed to truncate table: %s", table)
	}
}

// FlushRedis flushes all Redis data
func (env *TestEnvironment) FlushRedis(tb testing.TB) {
	tb.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := env.Redis.FlushAll(ctx).Err()
	require.NoError(tb, err, "Failed to flush Redis")

	tb.Log("Redis flushed")
}

// CreateTestProduct creates a test product and returns its ID
func (env *TestEnvironment) CreateTestProduct(tb testing.TB, storefrontID int64, price float64) int64 {
	tb.Helper()

	ctx := context.Background()

	var productID int64
	err := env.DB.QueryRowContext(ctx, `
		INSERT INTO listings
		(user_id, storefront_id, title, slug, description, price, currency, category_id, status, created_at, updated_at)
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
		RETURNING id
	`,
		1, // user_id
		storefrontID,
		fmt.Sprintf("Test Product %d", time.Now().UnixNano()),
		fmt.Sprintf("test-product-%d", time.Now().UnixNano()),
		"Test description",
		price,
		"EUR",
		1001,
		"active",
	).Scan(&productID)
	require.NoError(tb, err, "Failed to create test product")

	return productID
}

// AddCleanupFunc adds a custom cleanup function
func (env *TestEnvironment) AddCleanupFunc(fn func() error) {
	env.cleanupFuncs = append(env.cleanupFuncs, fn)
}

// Cleanup tears down the test environment
func (env *TestEnvironment) Cleanup() {
	// Run custom cleanup functions
	for _, fn := range env.cleanupFuncs {
		_ = fn()
	}

	// Close Redis
	if env.Redis != nil {
		_ = env.Redis.Close()
	}

	// Close pgxpool
	if env.PgPool != nil {
		env.PgPool.Close()
	}

	// Close DB
	if env.DB != nil {
		_ = env.DB.Close()
	}

	// Purge Docker containers
	if env.Pool != nil {
		if env.DBRes != nil {
			_ = env.Pool.Purge(env.DBRes)
		}
		if env.RedisRes != nil {
			_ = env.Pool.Purge(env.RedisRes)
		}
	}
}

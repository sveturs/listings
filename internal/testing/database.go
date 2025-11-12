package testing

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// TestDatabase represents a test database with transaction support.
// It provides automatic cleanup and test isolation.
type TestDatabase struct {
	DB       *sqlx.DB
	Pool     *dockertest.Pool
	Resource *dockertest.Resource
	DSN      string

	// tx holds the current transaction (if using transactional isolation)
	tx *sqlx.Tx

	// cleanup functions to run when database is torn down
	cleanupFuncs []func() error
}

// TestDatabaseConfig holds configuration for creating a test database
type TestDatabaseConfig struct {
	// UseDocker determines whether to use dockertest (true) or existing DB (false)
	UseDocker bool

	// DSN is the connection string (used when UseDocker is false)
	DSN string

	// MigrationsPath is the path to migration files
	MigrationsPath string

	// FixturesPath is the path to fixture files
	FixturesPath string

	// UseTransactions enables transactional test isolation
	// When true, each test runs in a transaction that is rolled back
	UseTransactions bool

	// MaxWaitTime is the maximum time to wait for Docker container to be ready
	MaxWaitTime time.Duration

	// PostgresVersion is the PostgreSQL Docker image tag
	PostgresVersion string
}

// DefaultTestDatabaseConfig returns a configuration with sensible defaults
func DefaultTestDatabaseConfig() TestDatabaseConfig {
	return TestDatabaseConfig{
		UseDocker:       true,
		MigrationsPath:  "../../migrations",
		MaxWaitTime:     30 * time.Second,
		PostgresVersion: "15-alpine",
		UseTransactions: false,
	}
}

// SetupTestDatabase creates a test database with optional Docker containerization.
// It automatically runs migrations and can load fixtures.
//
// Example with Docker:
//
//	testDB := testing.SetupTestDatabase(t, testing.DefaultTestDatabaseConfig())
//	defer testDB.Teardown(t)
//
// Example with existing database (for faster tests):
//
//	config := testing.DefaultTestDatabaseConfig()
//	config.UseDocker = false
//	config.DSN = "postgres://test:test@localhost:5432/testdb?sslmode=disable"
//	testDB := testing.SetupTestDatabase(t, config)
//	defer testDB.Teardown(t)
func SetupTestDatabase(t *testing.T, config TestDatabaseConfig) *TestDatabase {
	t.Helper()

	var db *sqlx.DB
	var pool *dockertest.Pool
	var resource *dockertest.Resource
	var dsn string

	if config.UseDocker {
		// Setup with Docker
		db, pool, resource, dsn = setupDockerDatabase(t, config)
	} else {
		// Use existing database
		if config.DSN == "" {
			t.Fatal("DSN is required when UseDocker is false")
		}
		var err error
		db, err = sqlx.Connect("postgres", config.DSN)
		require.NoError(t, err, "Failed to connect to database")
		dsn = config.DSN
	}

	testDB := &TestDatabase{
		DB:           db,
		Pool:         pool,
		Resource:     resource,
		DSN:          dsn,
		cleanupFuncs: []func() error{},
	}

	// Run migrations if path is provided
	if config.MigrationsPath != "" {
		testDB.RunMigrations(t, config.MigrationsPath)
	}

	// Load fixtures if path is provided
	if config.FixturesPath != "" {
		testDB.LoadFixtures(t, config.FixturesPath)
	}

	// Start transaction if using transactional isolation
	if config.UseTransactions {
		testDB.BeginTransaction(t)
	}

	return testDB
}

// setupDockerDatabase creates a PostgreSQL container for testing
func setupDockerDatabase(t *testing.T, config TestDatabaseConfig) (*sqlx.DB, *dockertest.Pool, *dockertest.Resource, string) {
	t.Helper()

	pool, err := dockertest.NewPool("")
	require.NoError(t, err, "Could not connect to docker")

	err = pool.Client.Ping()
	require.NoError(t, err, "Could not ping docker")

	// Generate unique database name to avoid conflicts
	dbName := fmt.Sprintf("test_db_%d", time.Now().UnixNano())

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        config.PostgresVersion,
		Env: []string{
			"POSTGRES_USER=test_user",
			"POSTGRES_PASSWORD=test_password",
			fmt.Sprintf("POSTGRES_DB=%s", dbName),
			"listen_addresses = '*'",
		},
	}, func(hostConfig *docker.HostConfig) {
		hostConfig.AutoRemove = true
		hostConfig.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	require.NoError(t, err, "Could not start PostgreSQL container")

	// Set expiry to clean up in case of test panic
	err = resource.Expire(120)
	require.NoError(t, err, "Could not set container expiry")

	var db *sqlx.DB
	hostAndPort := resource.GetHostPort("5432/tcp")
	dsn := fmt.Sprintf("postgres://test_user:test_password@%s/%s?sslmode=disable", hostAndPort, dbName)

	// Wait for database to be ready
	pool.MaxWait = config.MaxWaitTime
	err = pool.Retry(func() error {
		var connErr error
		db, connErr = sqlx.Connect("postgres", dsn)
		if connErr != nil {
			return connErr
		}
		return db.Ping()
	})
	require.NoError(t, err, "Could not connect to PostgreSQL container")

	return db, pool, resource, dsn
}

// BeginTransaction starts a transaction for test isolation.
// The transaction will be rolled back on Teardown.
func (tdb *TestDatabase) BeginTransaction(t *testing.T) {
	t.Helper()

	if tdb.tx != nil {
		t.Fatal("Transaction already started")
	}

	tx, err := tdb.DB.Beginx()
	require.NoError(t, err, "Could not begin transaction")

	tdb.tx = tx
}

// GetDB returns the database connection.
// If a transaction is active, it returns the transaction.
func (tdb *TestDatabase) GetDB() sqlx.Ext {
	if tdb.tx != nil {
		return tdb.tx
	}
	return tdb.DB
}

// GetDBx returns the database connection as *sqlx.DB
func (tdb *TestDatabase) GetDBx() *sqlx.DB {
	return tdb.DB
}

// GetTx returns the current transaction (if any)
func (tdb *TestDatabase) GetTx() *sqlx.Tx {
	return tdb.tx
}

// Commit commits the current transaction (if any).
// This is useful for testing transaction-dependent code.
func (tdb *TestDatabase) Commit(t *testing.T) {
	t.Helper()

	if tdb.tx == nil {
		t.Fatal("No active transaction to commit")
	}

	err := tdb.tx.Commit()
	require.NoError(t, err, "Could not commit transaction")

	tdb.tx = nil
}

// Rollback rolls back the current transaction (if any)
func (tdb *TestDatabase) Rollback(t *testing.T) {
	t.Helper()

	if tdb.tx == nil {
		return
	}

	err := tdb.tx.Rollback()
	require.NoError(t, err, "Could not rollback transaction")

	tdb.tx = nil
}

// RunMigrations runs all migration files from the specified directory.
// Migrations are executed in alphanumeric order (e.g., 001_init.up.sql, 002_users.up.sql).
func (tdb *TestDatabase) RunMigrations(t *testing.T, migrationsDir string) {
	t.Helper()

	files, err := os.ReadDir(migrationsDir)
	require.NoError(t, err, "Could not read migrations directory: %s", migrationsDir)

	// Filter and sort .up.sql files
	var migrationFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".up.sql") {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}
	sort.Strings(migrationFiles)

	t.Logf("Running %d migrations from %s", len(migrationFiles), migrationsDir)

	for _, fileName := range migrationFiles {
		migrationPath := filepath.Join(migrationsDir, fileName)
		data, err := os.ReadFile(migrationPath)
		require.NoError(t, err, "Could not read migration file: %s", migrationPath)

		_, err = tdb.DB.Exec(string(data))
		require.NoError(t, err, "Could not run migration: %s", fileName)

		t.Logf("Applied migration: %s", fileName)
	}
}

// LoadFixtures loads SQL fixtures from a file or directory.
// If path is a directory, all .sql files are loaded in alphanumeric order.
func (tdb *TestDatabase) LoadFixtures(t *testing.T, path string) {
	t.Helper()

	info, err := os.Stat(path)
	require.NoError(t, err, "Could not stat fixtures path: %s", path)

	if info.IsDir() {
		tdb.loadFixturesFromDir(t, path)
	} else {
		tdb.loadFixtureFile(t, path)
	}
}

// loadFixturesFromDir loads all .sql files from a directory
func (tdb *TestDatabase) loadFixturesFromDir(t *testing.T, dir string) {
	t.Helper()

	files, err := os.ReadDir(dir)
	require.NoError(t, err, "Could not read fixtures directory: %s", dir)

	var fixtureFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			fixtureFiles = append(fixtureFiles, file.Name())
		}
	}
	sort.Strings(fixtureFiles)

	t.Logf("Loading %d fixture files from %s", len(fixtureFiles), dir)

	for _, fileName := range fixtureFiles {
		fixturePath := filepath.Join(dir, fileName)
		tdb.loadFixtureFile(t, fixturePath)
	}
}

// loadFixtureFile loads a single SQL fixture file
func (tdb *TestDatabase) loadFixtureFile(t *testing.T, filePath string) {
	t.Helper()

	data, err := os.ReadFile(filePath)
	require.NoError(t, err, "Could not read fixture file: %s", filePath)

	_, err = tdb.DB.Exec(string(data))
	require.NoError(t, err, "Could not load fixtures from: %s", filePath)

	t.Logf("Loaded fixture: %s", filepath.Base(filePath))
}

// TruncateTables truncates the specified tables in the correct order (respecting foreign keys).
// Use CASCADE to also truncate dependent tables.
func (tdb *TestDatabase) TruncateTables(t *testing.T, tables ...string) {
	t.Helper()

	for _, table := range tables {
		_, err := tdb.DB.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		require.NoError(t, err, "Could not truncate table: %s", table)
	}
}

// CleanupTestData removes test data based on ID ranges.
// This is useful for cleaning up test data without truncating entire tables.
func (tdb *TestDatabase) CleanupTestData(t *testing.T, table string, minID, maxID int64) {
	t.Helper()

	query := fmt.Sprintf("DELETE FROM %s WHERE id >= $1 AND id < $2", table)
	_, err := tdb.DB.Exec(query, minID, maxID)
	require.NoError(t, err, "Could not cleanup test data from table: %s", table)
}

// ExecuteSQL executes arbitrary SQL (useful for complex test setup)
func (tdb *TestDatabase) ExecuteSQL(t *testing.T, sql string, args ...interface{}) {
	t.Helper()

	_, err := tdb.DB.Exec(sql, args...)
	require.NoError(t, err, "Could not execute SQL: %s", sql)
}

// QueryOne executes a query and returns a single row
func (tdb *TestDatabase) QueryOne(t *testing.T, dest interface{}, query string, args ...interface{}) {
	t.Helper()

	err := tdb.DB.Get(dest, query, args...)
	require.NoError(t, err, "Query failed: %s", query)
}

// QueryMany executes a query and returns multiple rows
func (tdb *TestDatabase) QueryMany(t *testing.T, dest interface{}, query string, args ...interface{}) {
	t.Helper()

	err := tdb.DB.Select(dest, query, args...)
	require.NoError(t, err, "Query failed: %s", query)
}

// CountRows returns the number of rows in a table matching the condition
func (tdb *TestDatabase) CountRows(t *testing.T, table string, condition string, args ...interface{}) int {
	t.Helper()

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", table, condition)
	var count int
	err := tdb.DB.Get(&count, query, args...)
	require.NoError(t, err, "Could not count rows in table: %s", table)

	return count
}

// RowExists checks if at least one row exists matching the condition
func (tdb *TestDatabase) RowExists(t *testing.T, table string, condition string, args ...interface{}) bool {
	t.Helper()

	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE %s)", table, condition)
	var exists bool
	err := tdb.DB.Get(&exists, query, args...)
	require.NoError(t, err, "Could not check row existence in table: %s", table)

	return exists
}

// AddCleanupFunc registers a cleanup function to be called during Teardown.
// This is useful for custom cleanup logic.
func (tdb *TestDatabase) AddCleanupFunc(fn func() error) {
	tdb.cleanupFuncs = append(tdb.cleanupFuncs, fn)
}

// Teardown cleans up the test database.
// This should be called in a defer statement after creating the test database.
func (tdb *TestDatabase) Teardown(t *testing.T) {
	t.Helper()

	// Rollback transaction if active
	if tdb.tx != nil {
		tdb.Rollback(t)
	}

	// Run custom cleanup functions
	for i, fn := range tdb.cleanupFuncs {
		if err := fn(); err != nil {
			t.Logf("Warning: cleanup function %d failed: %v", i, err)
		}
	}

	// Close database connection
	if tdb.DB != nil {
		_ = tdb.DB.Close()
	}

	// Purge Docker container if using Docker
	if tdb.Pool != nil && tdb.Resource != nil {
		err := tdb.Pool.Purge(tdb.Resource)
		require.NoError(t, err, "Could not purge PostgreSQL container")
	}
}

// =============================================================================
// Helper Functions for Test Isolation
// =============================================================================

// WithTransaction runs a test function inside a transaction that is automatically rolled back.
// This provides test isolation without the overhead of Docker containers.
//
// Example:
//
//	testing.WithTransaction(t, testDB, func(t *testing.T, tx *sqlx.Tx) {
//	    // Test code here - changes will be rolled back automatically
//	    _, err := tx.Exec("INSERT INTO listings (...) VALUES (...)")
//	    require.NoError(t, err)
//	})
func WithTransaction(t *testing.T, db *sqlx.DB, fn func(t *testing.T, tx *sqlx.Tx)) {
	t.Helper()

	tx, err := db.Beginx()
	require.NoError(t, err, "Could not begin transaction")

	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			t.Logf("Warning: failed to rollback transaction: %v", err)
		}
	}()

	fn(t, tx)
}

// WithIsolation runs a test function with complete database isolation.
// It truncates specified tables before and after the test.
//
// Example:
//
//	testing.WithIsolation(t, testDB, []string{"listings", "categories"}, func(t *testing.T) {
//	    // Test code here - tables are clean before and after
//	})
func WithIsolation(t *testing.T, tdb *TestDatabase, tables []string, fn func(t *testing.T)) {
	t.Helper()

	// Truncate before test
	tdb.TruncateTables(t, tables...)

	// Run test
	fn(t)

	// Truncate after test
	tdb.TruncateTables(t, tables...)
}

// =============================================================================
// Context Helpers
// =============================================================================

// TestContext creates a test context with a reasonable timeout
func TestContext(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(cancel)
	return ctx
}

// TestContextWithTimeout creates a test context with a custom timeout
func TestContextWithTimeout(t *testing.T, timeout time.Duration) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	t.Cleanup(cancel)
	return ctx
}

// =============================================================================
// Skip Helpers
// =============================================================================

// SkipIfShort skips the test if running in short mode
func SkipIfShort(t *testing.T) {
	t.Helper()
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}
}

// SkipIfNoDocker skips the test if Docker is not available
func SkipIfNoDocker(t *testing.T) {
	t.Helper()
	pool, err := dockertest.NewPool("")
	if err != nil || pool.Client.Ping() != nil {
		t.Skip("Docker not available, skipping test")
	}
}

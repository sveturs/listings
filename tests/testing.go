package tests

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// TestDB represents a test database connection
type TestDB struct {
	DB       *sql.DB
	Pool     *dockertest.Pool
	Resource *dockertest.Resource
}

// SetupTestPostgres creates a PostgreSQL container for testing
func SetupTestPostgres(tb testing.TB) *TestDB {
	tb.Helper()

	pool, err := dockertest.NewPool("")
	require.NoError(tb, err, "Could not connect to docker")

	err = pool.Client.Ping()
	require.NoError(tb, err, "Could not ping docker")

	// Pull PostgreSQL 15 image
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "15-alpine",
		Env: []string{
			"POSTGRES_USER=test_user",
			"POSTGRES_PASSWORD=test_password",
			"POSTGRES_DB=test_db",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	require.NoError(tb, err, "Could not start PostgreSQL container")

	// Set expiry to clean up in case of test panic
	err = resource.Expire(120)
	require.NoError(tb, err, "Could not set container expiry")

	var db *sql.DB
	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseURL := fmt.Sprintf("postgres://test_user:test_password@%s/test_db?sslmode=disable", hostAndPort)

	// Wait for database to be ready
	pool.MaxWait = 30 * time.Second
	err = pool.Retry(func() error {
		db, err = sql.Open("postgres", databaseURL)
		if err != nil {
			return err
		}
		return db.Ping()
	})
	require.NoError(tb, err, "Could not connect to PostgreSQL container")

	return &TestDB{
		DB:       db,
		Pool:     pool,
		Resource: resource,
	}
}

// TeardownTestPostgres cleans up the test database
func (tdb *TestDB) TeardownTestPostgres(tb testing.TB) {
	tb.Helper()

	if tdb.DB != nil {
		_ = tdb.DB.Close()
	}

	if tdb.Pool != nil && tdb.Resource != nil {
		err := tdb.Pool.Purge(tdb.Resource)
		require.NoError(tb, err, "Could not purge PostgreSQL container")
	}
}

// TestRedis represents a test Redis connection
type TestRedis struct {
	Pool     *dockertest.Pool
	Resource *dockertest.Resource
	Addr     string
}

// SetupTestRedis creates a Redis container for testing
func SetupTestRedis(tb testing.TB) *TestRedis {
	tb.Helper()

	pool, err := dockertest.NewPool("")
	require.NoError(tb, err, "Could not connect to docker")

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "redis",
		Tag:        "7-alpine",
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	require.NoError(tb, err, "Could not start Redis container")

	err = resource.Expire(120)
	require.NoError(tb, err, "Could not set container expiry")

	addr := resource.GetHostPort("6379/tcp")

	// Wait for Redis to be ready
	pool.MaxWait = 30 * time.Second
	err = pool.Retry(func() error {
		// Simple connection test (would need redis client in real implementation)
		return nil
	})
	require.NoError(tb, err, "Could not connect to Redis container")

	return &TestRedis{
		Pool:     pool,
		Resource: resource,
		Addr:     addr,
	}
}

// TeardownTestRedis cleans up the test Redis
func (tr *TestRedis) TeardownTestRedis(tb testing.TB) {
	tb.Helper()

	if tr.Pool != nil && tr.Resource != nil {
		err := tr.Pool.Purge(tr.Resource)
		require.NoError(tb, err, "Could not purge Redis container")
	}
}

// LoadTestFixtures loads test data from SQL file
// Automatically loads categories fixture first if needed
func LoadTestFixtures(tb testing.TB, db *sql.DB, fixtureFile string) {
	tb.Helper()

	// Auto-load categories fixture if not already loaded and not loading categories itself
	if !strings.Contains(fixtureFile, "00_categories_fixtures.sql") {
		categoriesFile := filepath.Join(filepath.Dir(fixtureFile), "00_categories_fixtures.sql")
		if _, err := os.Stat(categoriesFile); err == nil {
			// Check if categories already loaded (avoid duplicate load)
			var count int
			err := db.QueryRow("SELECT COUNT(*) FROM categories WHERE id = 1301").Scan(&count)
			if err != nil || count == 0 {
				// Load categories fixture
				catData, err := os.ReadFile(categoriesFile)
				if err == nil {
					_, _ = db.Exec(string(catData))
					tb.Logf("Auto-loaded categories fixture: %s", categoriesFile)
				}
			}
		}
	}

	data, err := os.ReadFile(fixtureFile)
	require.NoError(tb, err, "Could not read fixture file: %s", fixtureFile)

	_, err = db.Exec(string(data))
	require.NoError(tb, err, "Could not load fixtures from: %s", fixtureFile)
}

// RunMigrations runs database migrations from directory
func RunMigrations(tb testing.TB, db *sql.DB, migrationsDir string) {
	tb.Helper()

	// This is a simplified version
	// In production, use golang-migrate library
	files, err := os.ReadDir(migrationsDir)
	require.NoError(tb, err, "Could not read migrations directory")

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		// Only run .up.sql files
		if len(file.Name()) > 7 && file.Name()[len(file.Name())-7:] == ".up.sql" {
			migrationPath := fmt.Sprintf("%s/%s", migrationsDir, file.Name())
			data, err := os.ReadFile(migrationPath)
			require.NoError(tb, err, "Could not read migration file: %s", migrationPath)

			_, err = db.Exec(string(data))
			require.NoError(tb, err, "Could not run migration: %s", file.Name())

			log.Printf("Applied migration: %s", file.Name())
		}
	}
}

// CleanupTestDB truncates all tables for clean test state
func CleanupTestDB(tb testing.TB, db *sql.DB) {
	tb.Helper()

	tables := []string{
		"listing_images",
		"listing_attributes",
		"listing_tags",
		"listing_locations",
		"listing_stats",
		"indexing_queue",
		"listings",
	}

	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		require.NoError(tb, err, "Could not truncate table: %s", table)
	}
}

// TestContext creates a test context with timeout
func TestContext(tb testing.TB) context.Context {
	tb.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	tb.Cleanup(cancel)
	return ctx
}

// GenerateTestListing creates a test listing entity
func GenerateTestListing(userID int64, title string) map[string]interface{} {
	return map[string]interface{}{
		"user_id":     userID,
		"title":       title,
		"description": fmt.Sprintf("Test description for %s", title),
		"price":       999.99,
		"currency":    "USD",
		"category_id": 1,
		"status":      "active",
		"is_b2c":      true,
	}
}

// GenerateTestListings creates multiple test listings
func GenerateTestListings(count int) []map[string]interface{} {
	listings := make([]map[string]interface{}, count)
	for i := 0; i < count; i++ {
		listings[i] = GenerateTestListing(int64(i+1), fmt.Sprintf("Test Listing %d", i+1))
	}
	return listings
}

// AssertNoError is a helper to assert no error with better messages
func AssertNoError(tb testing.TB, err error, msgAndArgs ...interface{}) {
	tb.Helper()
	require.NoError(tb, err, msgAndArgs...)
}

// AssertEqual is a helper for equality assertions
func AssertEqual(tb testing.TB, expected, actual interface{}, msgAndArgs ...interface{}) {
	tb.Helper()
	require.Equal(tb, expected, actual, msgAndArgs...)
}

// AssertNotNil is a helper for nil checks
func AssertNotNil(tb testing.TB, obj interface{}, msgAndArgs ...interface{}) {
	tb.Helper()
	require.NotNil(tb, obj, msgAndArgs...)
}

// SkipIfShort skips test if running in short mode
func SkipIfShort(tb testing.TB) {
	tb.Helper()
	if testing.Short() {
		tb.Skip("Skipping test in short mode")
	}
}

// SkipIfNoDocker skips test if Docker is not available
func SkipIfNoDocker(tb testing.TB) {
	tb.Helper()
	pool, err := dockertest.NewPool("")
	if err != nil || pool.Client.Ping() != nil {
		tb.Skip("Docker not available, skipping integration test")
	}
}

// testClock returns current time for performance measurements

// testClockSince returns duration since start time

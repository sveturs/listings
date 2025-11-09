package health

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestDB creates a mock database for testing
func setupTestDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	require.NoError(t, err)
	return db, mock
}

// setupTestRedis creates a mock Redis client for testing
func setupTestRedis(t *testing.T) *redis.Client {
	// Use a test Redis instance if available, otherwise mock
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:36380",
		DB:   15, // Use a separate DB for tests
	})

	// Clean up test keys
	t.Cleanup(func() {
		client.FlushDB(context.Background())
		client.Close()
	})

	return client
}

func TestNewService(t *testing.T) {
	logger := zerolog.Nop()
	db, _ := setupTestDB(t)
	defer db.Close()

	config := DefaultConfig()

	service := NewService(db, nil, nil, nil, config, logger)

	assert.NotNil(t, service)
	assert.Equal(t, config, service.config)
	assert.False(t, service.isStarted) // Initially not started
}

func TestCheckDatabase_Success(t *testing.T) {
	logger := zerolog.Nop()
	db, mock := setupTestDB(t)
	defer db.Close()

	// Mock successful ping
	mock.ExpectPing()

	// Mock successful query
	rows := sqlmock.NewRows([]string{"count"}).AddRow(5)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM listings").WillReturnRows(rows)

	config := DefaultConfig()
	service := NewService(db, nil, nil, nil, config, logger)

	ctx := context.Background()
	result := service.CheckDatabase(ctx)

	assert.Equal(t, CheckStatusHealthy, result.Status)
	assert.Contains(t, result.Details, "connections active")
	assert.NoError(t, result.Error)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCheckDatabase_PingFailed(t *testing.T) {
	logger := zerolog.Nop()
	db, mock := setupTestDB(t)
	defer db.Close()

	// Mock failed ping
	mock.ExpectPing().WillReturnError(sql.ErrConnDone)

	config := DefaultConfig()
	service := NewService(db, nil, nil, nil, config, logger)

	ctx := context.Background()
	result := service.CheckDatabase(ctx)

	assert.Equal(t, CheckStatusUnhealthy, result.Status)
	assert.Equal(t, "ping failed", result.Details)
	assert.Error(t, result.Error)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCheckDatabase_QueryFailed(t *testing.T) {
	logger := zerolog.Nop()
	db, mock := setupTestDB(t)
	defer db.Close()

	// Mock successful ping
	mock.ExpectPing()

	// Mock failed query
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM listings").WillReturnError(sql.ErrNoRows)

	config := DefaultConfig()
	service := NewService(db, nil, nil, nil, config, logger)

	ctx := context.Background()
	result := service.CheckDatabase(ctx)

	assert.Equal(t, CheckStatusUnhealthy, result.Status)
	assert.Equal(t, "query failed", result.Details)
	assert.Error(t, result.Error)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCheckDatabase_Caching(t *testing.T) {
	logger := zerolog.Nop()
	db, mock := setupTestDB(t)
	defer db.Close()

	// Mock successful ping and query (only once)
	mock.ExpectPing()
	rows := sqlmock.NewRows([]string{"count"}).AddRow(5)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM listings").WillReturnRows(rows)

	config := DefaultConfig()
	config.CacheDuration = 1 * time.Second
	service := NewService(db, nil, nil, nil, config, logger)

	ctx := context.Background()

	// First call - should hit DB
	result1 := service.CheckDatabase(ctx)
	assert.Equal(t, CheckStatusHealthy, result1.Status)

	// Second call - should use cache (no additional mock expectations)
	result2 := service.CheckDatabase(ctx)
	assert.Equal(t, CheckStatusHealthy, result2.Status)

	// Verify only one set of expectations was met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCheckRedis_NilClient(t *testing.T) {
	logger := zerolog.Nop()
	config := DefaultConfig()
	service := NewService(nil, nil, nil, nil, config, logger)

	ctx := context.Background()
	result := service.CheckRedis(ctx)

	assert.Equal(t, CheckStatusUnhealthy, result.Status)
	assert.Contains(t, result.Details, "not initialized")
	assert.Error(t, result.Error)
}

func TestCheckOpenSearch_NilClient(t *testing.T) {
	logger := zerolog.Nop()
	config := DefaultConfig()
	service := NewService(nil, nil, nil, nil, config, logger)

	ctx := context.Background()
	result := service.CheckOpenSearch(ctx)

	// Should be healthy (optional dependency)
	assert.Equal(t, CheckStatusHealthy, result.Status)
	assert.Contains(t, result.Details, "not configured")
	assert.NoError(t, result.Error)
}

func TestCheckMinIO_NilClient(t *testing.T) {
	logger := zerolog.Nop()
	config := DefaultConfig()
	service := NewService(nil, nil, nil, nil, config, logger)

	ctx := context.Background()
	result := service.CheckMinIO(ctx)

	// Should be healthy (optional dependency)
	assert.Equal(t, CheckStatusHealthy, result.Status)
	assert.Contains(t, result.Details, "not configured")
	assert.NoError(t, result.Error)
}

func TestCheckAll(t *testing.T) {
	logger := zerolog.Nop()
	db, mock := setupTestDB(t)
	defer db.Close()

	// Mock successful DB check
	mock.ExpectPing()
	rows := sqlmock.NewRows([]string{"count"}).AddRow(5)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM listings").WillReturnRows(rows)

	config := DefaultConfig()
	service := NewService(db, nil, nil, nil, config, logger)

	ctx := context.Background()
	response := service.CheckAll(ctx)

	assert.Equal(t, HealthStatusUnhealthy, response.Status) // DB healthy, Redis unhealthy (critical)
	assert.NotEmpty(t, response.Version)
	assert.NotEmpty(t, response.Uptime)
	assert.Len(t, response.Checks, 4) // database, redis, opensearch, minio
	assert.NotZero(t, response.Timestamp)

	// Check individual results
	assert.Equal(t, CheckStatusHealthy, response.Checks["database"].Status)
	assert.Equal(t, CheckStatusUnhealthy, response.Checks["redis"].Status)
	assert.Equal(t, CheckStatusHealthy, response.Checks["opensearch"].Status) // Optional
	assert.Equal(t, CheckStatusHealthy, response.Checks["minio"].Status)      // Optional

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCheckLiveness(t *testing.T) {
	logger := zerolog.Nop()
	config := DefaultConfig()
	service := NewService(nil, nil, nil, nil, config, logger)

	ctx := context.Background()
	isAlive := service.CheckLiveness(ctx)

	// Liveness should always be true (service is running)
	assert.True(t, isAlive)
}

func TestCheckReadiness(t *testing.T) {
	logger := zerolog.Nop()
	db, mock := setupTestDB(t)
	defer db.Close()

	// Mock successful DB check
	mock.ExpectPing()
	rows := sqlmock.NewRows([]string{"count"}).AddRow(5)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM listings").WillReturnRows(rows)

	config := DefaultConfig()
	service := NewService(db, nil, nil, nil, config, logger)

	ctx := context.Background()
	isReady := service.CheckReadiness(ctx)

	// Should be false because Redis is nil
	assert.False(t, isReady)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCheckStartup(t *testing.T) {
	logger := zerolog.Nop()
	db, mock := setupTestDB(t)
	defer db.Close()

	// Mock successful DB check
	mock.ExpectPing()
	rows := sqlmock.NewRows([]string{"count"}).AddRow(5)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM listings").WillReturnRows(rows)

	config := DefaultConfig()
	service := NewService(db, nil, nil, nil, config, logger)

	ctx := context.Background()

	// Should be false initially (grace period not elapsed)
	isStarted := service.CheckStartup(ctx)
	assert.False(t, isStarted)

	// Manually mark as started for testing
	service.startupMu.Lock()
	service.isStarted = true
	service.startupMu.Unlock()

	// Should still be false because Redis is nil
	isStarted = service.CheckStartup(ctx)
	assert.False(t, isStarted)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDetermineOverallStatus(t *testing.T) {
	logger := zerolog.Nop()
	config := DefaultConfig()
	service := NewService(nil, nil, nil, nil, config, logger)

	tests := []struct {
		name     string
		checks   map[string]CheckResult
		expected HealthStatus
	}{
		{
			name: "all healthy",
			checks: map[string]CheckResult{
				"database":   {Status: CheckStatusHealthy},
				"redis":      {Status: CheckStatusHealthy},
				"opensearch": {Status: CheckStatusHealthy},
				"minio":      {Status: CheckStatusHealthy},
			},
			expected: HealthStatusHealthy,
		},
		{
			name: "critical unhealthy",
			checks: map[string]CheckResult{
				"database":   {Status: CheckStatusUnhealthy},
				"redis":      {Status: CheckStatusHealthy},
				"opensearch": {Status: CheckStatusHealthy},
				"minio":      {Status: CheckStatusHealthy},
			},
			expected: HealthStatusUnhealthy,
		},
		{
			name: "optional unhealthy",
			checks: map[string]CheckResult{
				"database":   {Status: CheckStatusHealthy},
				"redis":      {Status: CheckStatusHealthy},
				"opensearch": {Status: CheckStatusUnhealthy},
				"minio":      {Status: CheckStatusHealthy},
			},
			expected: HealthStatusDegraded,
		},
		{
			name: "critical and optional unhealthy",
			checks: map[string]CheckResult{
				"database":   {Status: CheckStatusUnhealthy},
				"redis":      {Status: CheckStatusHealthy},
				"opensearch": {Status: CheckStatusUnhealthy},
				"minio":      {Status: CheckStatusHealthy},
			},
			expected: HealthStatusUnhealthy, // Critical takes precedence
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.determineOverallStatus(tt.checks)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRecordError(t *testing.T) {
	logger := zerolog.Nop()
	config := DefaultConfig()
	service := NewService(nil, nil, nil, nil, config, logger)

	// Record errors
	for i := 0; i < 15; i++ {
		service.recordError("test error")
	}

	// Should keep only last 10 errors
	service.errorsMu.RLock()
	errorCount := len(service.recentErrors)
	service.errorsMu.RUnlock()

	assert.Equal(t, 10, errorCount)
}

func TestGetUptime(t *testing.T) {
	logger := zerolog.Nop()
	config := DefaultConfig()
	service := NewService(nil, nil, nil, nil, config, logger)

	// Set startup time to 2 hours ago
	service.startupTime = time.Now().Add(-2 * time.Hour)

	uptime := service.getUptime()
	assert.Contains(t, uptime, "2h")
}

func TestCheckDeep(t *testing.T) {
	logger := zerolog.Nop()
	db, mock := setupTestDB(t)
	defer db.Close()

	// Mock successful DB check
	mock.ExpectPing()
	rows := sqlmock.NewRows([]string{"count"}).AddRow(5)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM listings").WillReturnRows(rows)

	config := DefaultConfig()
	config.EnableDeepChecks = true
	service := NewService(db, nil, nil, nil, config, logger)

	// Record some errors
	service.recordError("test error 1")
	service.recordError("test error 2")

	ctx := context.Background()
	response := service.CheckDeep(ctx)

	assert.Equal(t, HealthStatusUnhealthy, response.Status) // Redis unhealthy (critical)
	assert.NotZero(t, response.Diagnostics.Goroutines)
	assert.GreaterOrEqual(t, response.Diagnostics.MemoryAllocMB, uint64(0)) // Can be zero
	assert.NotZero(t, response.Diagnostics.MemorySysMB)
	assert.Len(t, response.Diagnostics.RecentErrors, 2)
	assert.NotEmpty(t, response.Diagnostics.ConnectionPools)

	// Check database pool metrics
	dbMetrics, exists := response.Diagnostics.ConnectionPools["database"]
	assert.True(t, exists)
	assert.GreaterOrEqual(t, dbMetrics.MaxOpen, 0)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCheckDeep_Disabled(t *testing.T) {
	logger := zerolog.Nop()
	db, _ := setupTestDB(t)
	defer db.Close()

	config := DefaultConfig()
	config.EnableDeepChecks = false
	service := NewService(db, nil, nil, nil, config, logger)

	ctx := context.Background()
	response := service.CheckDeep(ctx)

	// Should return empty diagnostics
	assert.Equal(t, 0, response.Diagnostics.Goroutines)
	assert.Equal(t, uint64(0), response.Diagnostics.MemoryAllocMB)
	assert.Nil(t, response.Diagnostics.RecentErrors)
}

func TestCacheExpiry(t *testing.T) {
	logger := zerolog.Nop()
	db, mock := setupTestDB(t)
	defer db.Close()

	// Mock successful ping and query (twice - cache will expire)
	mock.ExpectPing()
	rows1 := sqlmock.NewRows([]string{"count"}).AddRow(5)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM listings").WillReturnRows(rows1)

	mock.ExpectPing()
	rows2 := sqlmock.NewRows([]string{"count"}).AddRow(10)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM listings").WillReturnRows(rows2)

	config := DefaultConfig()
	config.CacheDuration = 100 * time.Millisecond
	service := NewService(db, nil, nil, nil, config, logger)

	ctx := context.Background()

	// First call - should hit DB
	result1 := service.CheckDatabase(ctx)
	assert.Equal(t, CheckStatusHealthy, result1.Status)

	// Wait for cache to expire
	time.Sleep(150 * time.Millisecond)

	// Second call - cache expired, should hit DB again
	result2 := service.CheckDatabase(ctx)
	assert.Equal(t, CheckStatusHealthy, result2.Status)

	// Verify both sets of expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

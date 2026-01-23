package health

import (
	"context"
	"database/sql"
	"fmt"
	"runtime"
	"sync"
	"time"

	opensearchapi "github.com/opensearch-project/opensearch-go/v2/opensearchapi"
	"github.com/rs/zerolog"

	"github.com/vondi-global/listings/internal/cache"
	"github.com/vondi-global/listings/internal/repository/minio"
	"github.com/vondi-global/listings/internal/repository/opensearch"
)

// Service implements health checking for all dependencies
type Service struct {
	db           *sql.DB
	redisCache   *cache.RedisCache
	searchClient *opensearch.Client
	minioClient  *minio.Client
	config       *Config
	logger       zerolog.Logger

	// Startup tracking
	startupTime time.Time
	isStarted   bool
	startupMu   sync.RWMutex

	// Error tracking
	recentErrors []string
	errorsMu     sync.RWMutex

	// Cache for check results
	cachedResults map[string]*cachedCheckResult
	cacheMu       sync.RWMutex
}

type cachedCheckResult struct {
	result    CheckResult
	timestamp time.Time
}

// NewService creates a new health check service
func NewService(
	db *sql.DB,
	redisCache *cache.RedisCache,
	searchClient *opensearch.Client,
	minioClient *minio.Client,
	config *Config,
	logger zerolog.Logger,
) *Service {
	if config == nil {
		config = DefaultConfig()
	}

	s := &Service{
		db:            db,
		redisCache:    redisCache,
		searchClient:  searchClient,
		minioClient:   minioClient,
		config:        config,
		logger:        logger.With().Str("component", "health").Logger(),
		startupTime:   time.Now(),
		isStarted:     false,
		recentErrors:  make([]string, 0, 10),
		cachedResults: make(map[string]*cachedCheckResult),
	}

	// Mark as started after a grace period
	go func() {
		time.Sleep(5 * time.Second)
		s.startupMu.Lock()
		s.isStarted = true
		s.startupMu.Unlock()
		s.logger.Info().Msg("health check service marked as started")
	}()

	return s
}

// CheckDatabase verifies PostgreSQL connectivity
func (s *Service) CheckDatabase(ctx context.Context) CheckResult {
	start := time.Now()

	// Check cache first
	if cached := s.getCachedResult("database"); cached != nil {
		return *cached
	}

	ctx, cancel := context.WithTimeout(ctx, s.config.CheckTimeout)
	defer cancel()

	// Test ping
	if err := s.db.PingContext(ctx); err != nil {
		s.recordError(fmt.Sprintf("database ping failed: %v", err))
		result := CheckResult{
			Status:       CheckStatusUnhealthy,
			ResponseTime: time.Since(start),
			Details:      "ping failed",
			Error:        err,
		}
		s.cacheResult("database", result)
		return result
	}

	// Test simple query
	var count int
	err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM listings").Scan(&count)
	if err != nil {
		s.recordError(fmt.Sprintf("database query failed: %v", err))
		result := CheckResult{
			Status:       CheckStatusUnhealthy,
			ResponseTime: time.Since(start),
			Details:      "query failed",
			Error:        err,
		}
		s.cacheResult("database", result)
		return result
	}

	// Get connection stats
	stats := s.db.Stats()
	details := fmt.Sprintf("%d connections active, %d idle, %d in use",
		stats.OpenConnections, stats.Idle, stats.InUse)

	result := CheckResult{
		Status:       CheckStatusHealthy,
		ResponseTime: time.Since(start),
		Details:      details,
		Error:        nil,
	}

	s.cacheResult("database", result)
	return result
}

// CheckRedis verifies Redis connectivity
func (s *Service) CheckRedis(ctx context.Context) CheckResult {
	start := time.Now()

	// Check cache first
	if cached := s.getCachedResult("redis"); cached != nil {
		return *cached
	}

	if s.redisCache == nil {
		result := CheckResult{
			Status:       CheckStatusUnhealthy,
			ResponseTime: time.Since(start),
			Details:      "redis client not initialized",
			Error:        fmt.Errorf("redis client is nil"),
		}
		s.cacheResult("redis", result)
		return result
	}

	ctx, cancel := context.WithTimeout(ctx, s.config.CheckTimeout)
	defer cancel()

	// Test ping
	if err := s.redisCache.HealthCheck(ctx); err != nil {
		s.recordError(fmt.Sprintf("redis ping failed: %v", err))
		result := CheckResult{
			Status:       CheckStatusUnhealthy,
			ResponseTime: time.Since(start),
			Details:      "ping failed",
			Error:        err,
		}
		s.cacheResult("redis", result)
		return result
	}

	// Test set/get
	testKey := "health:check:test"
	testValue := time.Now().Unix()
	client := s.redisCache.GetClient()

	if err := client.Set(ctx, testKey, testValue, 10*time.Second).Err(); err != nil {
		s.recordError(fmt.Sprintf("redis set failed: %v", err))
		result := CheckResult{
			Status:       CheckStatusUnhealthy,
			ResponseTime: time.Since(start),
			Details:      "set operation failed",
			Error:        err,
		}
		s.cacheResult("redis", result)
		return result
	}

	val, err := client.Get(ctx, testKey).Int64()
	if err != nil || val != testValue {
		s.recordError(fmt.Sprintf("redis get failed: %v", err))
		result := CheckResult{
			Status:       CheckStatusUnhealthy,
			ResponseTime: time.Since(start),
			Details:      "get operation failed",
			Error:        err,
		}
		s.cacheResult("redis", result)
		return result
	}

	// Clean up test key
	_ = client.Del(ctx, testKey)

	// Get pool stats
	poolStats := s.redisCache.GetPoolStats()
	details := fmt.Sprintf("%d hits, %d misses, %d total conns",
		poolStats.Hits, poolStats.Misses, poolStats.TotalConns)

	result := CheckResult{
		Status:       CheckStatusHealthy,
		ResponseTime: time.Since(start),
		Details:      details,
		Error:        nil,
	}

	s.cacheResult("redis", result)
	return result
}

// CheckOpenSearch verifies OpenSearch connectivity
func (s *Service) CheckOpenSearch(ctx context.Context) CheckResult {
	start := time.Now()

	// Check cache first
	if cached := s.getCachedResult("opensearch"); cached != nil {
		return *cached
	}

	if s.searchClient == nil {
		result := CheckResult{
			Status:       CheckStatusHealthy, // Degraded, not unhealthy (optional dependency)
			ResponseTime: time.Since(start),
			Details:      "opensearch not configured (optional)",
			Error:        nil,
		}
		s.cacheResult("opensearch", result)
		return result
	}

	ctx, cancel := context.WithTimeout(ctx, s.config.CheckTimeout)
	defer cancel()

	// Check cluster health
	req := opensearchapi.ClusterHealthRequest{}
	res, err := req.Do(ctx, s.searchClient.GetClient())
	if err != nil {
		s.recordError(fmt.Sprintf("opensearch cluster health failed: %v", err))
		result := CheckResult{
			Status:       CheckStatusUnhealthy,
			ResponseTime: time.Since(start),
			Details:      "cluster health check failed",
			Error:        err,
		}
		s.cacheResult("opensearch", result)
		return result
	}
	defer res.Body.Close()

	if res.IsError() {
		err := fmt.Errorf("cluster health returned error: %s", res.Status())
		s.recordError(fmt.Sprintf("opensearch: %v", err))
		result := CheckResult{
			Status:       CheckStatusUnhealthy,
			ResponseTime: time.Since(start),
			Details:      res.Status(),
			Error:        err,
		}
		s.cacheResult("opensearch", result)
		return result
	}

	details := fmt.Sprintf("cluster status: %s", res.Status())

	result := CheckResult{
		Status:       CheckStatusHealthy,
		ResponseTime: time.Since(start),
		Details:      details,
		Error:        nil,
	}

	s.cacheResult("opensearch", result)
	return result
}

// CheckMinIO verifies MinIO connectivity
func (s *Service) CheckMinIO(ctx context.Context) CheckResult {
	start := time.Now()

	// Check cache first
	if cached := s.getCachedResult("minio"); cached != nil {
		return *cached
	}

	if s.minioClient == nil {
		result := CheckResult{
			Status:       CheckStatusHealthy, // Degraded, not unhealthy (optional dependency)
			ResponseTime: time.Since(start),
			Details:      "minio not configured (optional)",
			Error:        nil,
		}
		s.cacheResult("minio", result)
		return result
	}

	ctx, cancel := context.WithTimeout(ctx, s.config.CheckTimeout)
	defer cancel()

	// Check bucket exists
	exists, err := s.minioClient.BucketExists(ctx)
	if err != nil {
		s.recordError(fmt.Sprintf("minio bucket check failed: %v", err))
		result := CheckResult{
			Status:       CheckStatusUnhealthy,
			ResponseTime: time.Since(start),
			Details:      "bucket check failed",
			Error:        err,
		}
		s.cacheResult("minio", result)
		return result
	}

	if !exists {
		err := fmt.Errorf("bucket does not exist")
		s.recordError(fmt.Sprintf("minio: %v", err))
		result := CheckResult{
			Status:       CheckStatusUnhealthy,
			ResponseTime: time.Since(start),
			Details:      "bucket not found",
			Error:        err,
		}
		s.cacheResult("minio", result)
		return result
	}

	details := "bucket accessible"

	result := CheckResult{
		Status:       CheckStatusHealthy,
		ResponseTime: time.Since(start),
		Details:      details,
		Error:        nil,
	}

	s.cacheResult("minio", result)
	return result
}

// CheckAll performs all health checks in parallel
func (s *Service) CheckAll(ctx context.Context) HealthResponse {
	type checkResult struct {
		name   string
		result CheckResult
	}

	results := make(chan checkResult, 4)
	var wg sync.WaitGroup

	// Run all checks in parallel
	checks := []struct {
		name string
		fn   func(context.Context) CheckResult
	}{
		{"database", s.CheckDatabase},
		{"redis", s.CheckRedis},
		{"opensearch", s.CheckOpenSearch},
		{"minio", s.CheckMinIO},
	}

	for _, check := range checks {
		wg.Add(1)
		go func(name string, fn func(context.Context) CheckResult) {
			defer wg.Done()
			results <- checkResult{
				name:   name,
				result: fn(ctx),
			}
		}(check.name, check.fn)
	}

	// Close results channel when all checks complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	checksMap := make(map[string]CheckResult)
	for result := range results {
		checksMap[result.name] = result.result
	}

	// Determine overall status
	overallStatus := s.determineOverallStatus(checksMap)

	return HealthResponse{
		Status:    overallStatus,
		Version:   "0.1.0", // TODO: get from build info
		Uptime:    s.getUptime(),
		Checks:    checksMap,
		Timestamp: time.Now(),
	}
}

// CheckDeep performs extended diagnostics
func (s *Service) CheckDeep(ctx context.Context) DeepHealthResponse {
	if !s.config.EnableDeepChecks {
		return DeepHealthResponse{
			HealthResponse: s.CheckAll(ctx),
			Diagnostics:    DiagnosticsInfo{},
		}
	}

	healthResp := s.CheckAll(ctx)

	// Get runtime stats
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Get connection pool metrics
	poolMetrics := make(map[string]PoolMetrics)

	// Database pool metrics
	if s.db != nil {
		stats := s.db.Stats()
		poolMetrics["database"] = PoolMetrics{
			Active:   stats.InUse,
			Idle:     stats.Idle,
			MaxOpen:  stats.MaxOpenConnections,
			WaitTime: int(stats.WaitDuration.Milliseconds()),
		}
	}

	// Redis pool metrics
	if s.redisCache != nil {
		poolStats := s.redisCache.GetPoolStats()
		poolMetrics["redis"] = PoolMetrics{
			Active:   int(poolStats.TotalConns) - int(poolStats.IdleConns),
			Idle:     int(poolStats.IdleConns),
			MaxOpen:  10, // From config, hardcoded here for simplicity
			WaitTime: 0,
		}
	}

	// Get recent errors
	s.errorsMu.RLock()
	recentErrors := make([]string, len(s.recentErrors))
	copy(recentErrors, s.recentErrors)
	s.errorsMu.RUnlock()

	diagnostics := DiagnosticsInfo{
		Goroutines:      runtime.NumGoroutine(),
		MemoryAllocMB:   m.Alloc / 1024 / 1024,
		MemorySysMB:     m.Sys / 1024 / 1024,
		NumGC:           m.NumGC,
		RecentErrors:    recentErrors,
		ConnectionPools: poolMetrics,
	}

	return DeepHealthResponse{
		HealthResponse: healthResp,
		Diagnostics:    diagnostics,
	}
}

// CheckLiveness performs minimal liveness check
func (s *Service) CheckLiveness(ctx context.Context) bool {
	// Just check if the service is running
	// No dependency checks required for liveness
	return true
}

// CheckReadiness performs readiness check
func (s *Service) CheckReadiness(ctx context.Context) bool {
	// Check critical dependencies (DB and Redis)
	dbResult := s.CheckDatabase(ctx)
	redisResult := s.CheckRedis(ctx)

	return dbResult.Status == CheckStatusHealthy &&
		redisResult.Status == CheckStatusHealthy
}

// CheckStartup performs startup probe check
func (s *Service) CheckStartup(ctx context.Context) bool {
	s.startupMu.RLock()
	started := s.isStarted
	s.startupMu.RUnlock()

	if !started {
		return false
	}

	// Check all critical dependencies with longer timeout
	ctx, cancel := context.WithTimeout(ctx, s.config.StartupTimeout)
	defer cancel()

	dbResult := s.CheckDatabase(ctx)
	redisResult := s.CheckRedis(ctx)

	return dbResult.Status == CheckStatusHealthy &&
		redisResult.Status == CheckStatusHealthy
}

// Helper methods

func (s *Service) determineOverallStatus(checks map[string]CheckResult) HealthStatus {
	hasUnhealthy := false
	hasDegraded := false

	for name, result := range checks {
		// Critical dependencies
		if name == "database" || name == "redis" {
			if result.Status == CheckStatusUnhealthy {
				hasUnhealthy = true
			}
		} else {
			// Optional dependencies
			if result.Status == CheckStatusUnhealthy {
				hasDegraded = true
			}
		}
	}

	if hasUnhealthy {
		return HealthStatusUnhealthy
	}
	if hasDegraded {
		return HealthStatusDegraded
	}
	return HealthStatusHealthy
}

func (s *Service) getUptime() string {
	uptime := time.Since(s.startupTime)
	hours := int(uptime.Hours())
	minutes := int(uptime.Minutes()) % 60
	return fmt.Sprintf("%dh%dm", hours, minutes)
}

func (s *Service) recordError(errMsg string) {
	s.errorsMu.Lock()
	defer s.errorsMu.Unlock()

	// Keep last 10 errors
	if len(s.recentErrors) >= 10 {
		s.recentErrors = s.recentErrors[1:]
	}
	s.recentErrors = append(s.recentErrors, fmt.Sprintf("%s: %s", time.Now().Format(time.RFC3339), errMsg))

	s.logger.Warn().Str("error", errMsg).Msg("health check error recorded")
}

func (s *Service) getCachedResult(name string) *CheckResult {
	s.cacheMu.RLock()
	defer s.cacheMu.RUnlock()

	cached, exists := s.cachedResults[name]
	if !exists {
		return nil
	}

	// Check if cache is still valid
	if time.Since(cached.timestamp) > s.config.CacheDuration {
		return nil
	}

	return &cached.result
}

func (s *Service) cacheResult(name string, result CheckResult) {
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()

	s.cachedResults[name] = &cachedCheckResult{
		result:    result,
		timestamp: time.Now(),
	}
}

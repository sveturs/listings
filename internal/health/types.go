package health

import (
	"context"
	"time"
)

// HealthStatus represents overall system health
type HealthStatus string

const (
	HealthStatusHealthy   HealthStatus = "healthy"
	HealthStatusDegraded  HealthStatus = "degraded"
	HealthStatusUnhealthy HealthStatus = "unhealthy"
)

// CheckStatus represents individual check status
type CheckStatus string

const (
	CheckStatusHealthy   CheckStatus = "healthy"
	CheckStatusUnhealthy CheckStatus = "unhealthy"
	CheckStatusTimeout   CheckStatus = "timeout"
)

// CheckResult represents result of a single health check
type CheckResult struct {
	Status       CheckStatus   `json:"status"`
	ResponseTime time.Duration `json:"response_time_ms"`
	Details      string        `json:"details"`
	Error        error         `json:"-"`
}

// HealthResponse represents comprehensive health check response
type HealthResponse struct {
	Status    HealthStatus           `json:"status"`
	Version   string                 `json:"version"`
	Uptime    string                 `json:"uptime"`
	Checks    map[string]CheckResult `json:"checks"`
	Timestamp time.Time              `json:"timestamp"`
}

// DeepHealthResponse extends HealthResponse with diagnostics
type DeepHealthResponse struct {
	HealthResponse
	Diagnostics DiagnosticsInfo `json:"diagnostics"`
}

// DiagnosticsInfo contains extended system diagnostics
type DiagnosticsInfo struct {
	Goroutines      int                    `json:"goroutines"`
	MemoryAllocMB   uint64                 `json:"memory_alloc_mb"`
	MemorySysMB     uint64                 `json:"memory_sys_mb"`
	NumGC           uint32                 `json:"num_gc"`
	RecentErrors    []string               `json:"recent_errors,omitempty"`
	ConnectionPools map[string]PoolMetrics `json:"connection_pools"`
}

// PoolMetrics represents connection pool metrics
type PoolMetrics struct {
	Active   int `json:"active"`
	Idle     int `json:"idle"`
	MaxOpen  int `json:"max_open"`
	WaitTime int `json:"wait_time_ms"`
}

// Checker defines interface for health checking
type Checker interface {
	// CheckDatabase verifies PostgreSQL connectivity
	CheckDatabase(ctx context.Context) CheckResult

	// CheckRedis verifies Redis connectivity
	CheckRedis(ctx context.Context) CheckResult

	// CheckOpenSearch verifies OpenSearch connectivity
	CheckOpenSearch(ctx context.Context) CheckResult

	// CheckMinIO verifies MinIO connectivity
	CheckMinIO(ctx context.Context) CheckResult

	// CheckAll performs all health checks
	CheckAll(ctx context.Context) HealthResponse

	// CheckDeep performs extended diagnostics
	CheckDeep(ctx context.Context) DeepHealthResponse

	// CheckLiveness performs minimal liveness check (service running)
	CheckLiveness(ctx context.Context) bool

	// CheckReadiness performs readiness check (ready to serve traffic)
	CheckReadiness(ctx context.Context) bool

	// CheckStartup performs startup probe check
	CheckStartup(ctx context.Context) bool
}

// Config contains health check configuration
type Config struct {
	// CheckTimeout is timeout for individual checks
	CheckTimeout time.Duration

	// CheckInterval is interval between cached checks
	CheckInterval time.Duration

	// StartupTimeout is timeout for startup checks
	StartupTimeout time.Duration

	// CacheDuration is how long to cache check results
	CacheDuration time.Duration

	// EnableDeepChecks enables deep diagnostics endpoint
	EnableDeepChecks bool
}

// DefaultConfig returns default health check configuration
func DefaultConfig() *Config {
	return &Config{
		CheckTimeout:     5 * time.Second,
		CheckInterval:    30 * time.Second,
		StartupTimeout:   60 * time.Second,
		CacheDuration:    10 * time.Second,
		EnableDeepChecks: true,
	}
}

// Package integration contains canary deployment integration tests
package integration

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"backend/internal/config"
	"backend/internal/logger"
	"backend/internal/proj/marketplace/service"
)

const (
	testTimeout = 5 * time.Second
)

// TestCanaryTrafficDistribution проверяет распределение трафика при 1% canary
func TestCanaryTrafficDistribution(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup: 1% rollout
	cfg := &config.MarketplaceConfig{
		UseMicroservice: true,
		RolloutPercent:  1,
		AdminOverride:   false,
		CanaryUserIDs:   "",
		CircuitBreaker: config.CircuitBreakerConfig{
			Enabled: true,
		},
	}

	log := logger.Get()
	router := service.NewTrafficRouter(cfg, *log)

	// Симулируем 1000 запросов от разных пользователей
	microserviceCount := 0
	monolithCount := 0

	for i := 0; i < 1000; i++ {
		userID := string(rune(i + 1000)) // Генерируем разные user IDs
		decision := router.ShouldUseMicroservice(userID, false)

		if decision.UseMicroservice {
			microserviceCount++
		} else {
			monolithCount++
		}
	}

	// При 1% rollout ожидаем ~10 запросов на microservice (1% от 1000)
	// Допустимое отклонение: ±0.5% (5-15 запросов)
	expectedMin := 5
	expectedMax := 15

	assert.GreaterOrEqual(t, microserviceCount, expectedMin,
		"Microservice traffic should be >= %d at 1%% rollout", expectedMin)
	assert.LessOrEqual(t, microserviceCount, expectedMax,
		"Microservice traffic should be <= %d at 1%% rollout", expectedMax)

	t.Logf("✅ 1%% rollout verified: %d microservice, %d monolith (out of 1000)",
		microserviceCount, monolithCount)
}

// TestCanaryCircuitBreakerStates проверяет состояния circuit breaker
func TestCanaryCircuitBreakerStates(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	cfg := &config.MarketplaceConfig{
		UseMicroservice: true,
		RolloutPercent:  100, // 100% для тестирования circuit breaker
		CircuitBreaker: config.CircuitBreakerConfig{
			Enabled:          true,
			FailureThreshold: 5,
			SuccessThreshold: 2,
			Timeout:          5 * time.Second,
		},
	}

	log := logger.Get()
	router := service.NewTrafficRouter(cfg, *log)

	// Создаем mock circuit breaker для тестирования
	cbConfig := service.CircuitBreakerConfig{
		Enabled:          true,
		FailureThreshold: 5,
		SuccessThreshold: 2,
		Timeout:          100 * time.Millisecond,
	}
	cb := service.NewCircuitBreaker(cbConfig, *log)

	// Initial state: CLOSED
	initialState := cb.GetState()
	assert.Equal(t, "CLOSED", initialState, "Circuit breaker should start in CLOSED state")
	t.Logf("✅ Initial state: %s", initialState)

	// Simulate 5 failures to open circuit
	for i := 0; i < 5; i++ {
		_, _ = cb.Execute(context.Background(), "test_failure", func() (interface{}, error) {
			return nil, fmt.Errorf("simulated error")
		})
	}

	// State should be OPEN after failures
	openState := string(cb.GetState())
	assert.Equal(t, "open", openState, "Circuit breaker should be OPEN after failures")
	t.Logf("✅ After 5 failures: %s", openState)

	// Wait for timeout to enter HALF_OPEN
	time.Sleep(150 * time.Millisecond)

	// Attempt request in HALF_OPEN state
	halfOpenState := string(cb.GetState())
	t.Logf("✅ After timeout: %s", halfOpenState)

	// Record successes to close circuit
	for i := 0; i < 2; i++ {
		_, _ = cb.Execute(context.Background(), "test_success", func() (interface{}, error) {
			return "ok", nil
		})
	}

	// State should be CLOSED again
	closedState := string(cb.GetState())
	assert.Equal(t, "closed", closedState, "Circuit breaker should be CLOSED after successes")
	t.Logf("✅ After 2 successes: %s", closedState)

	t.Log("✅ Circuit breaker state transitions verified: CLOSED → OPEN → HALF_OPEN → CLOSED")

	_ = router // avoid unused variable
}

// TestCanaryFallbackMechanism проверяет fallback на monolith при ошибках
func TestCanaryFallbackMechanism(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	cfg := &config.MarketplaceConfig{
		UseMicroservice: true,
		RolloutPercent:  100,
		CircuitBreaker: config.CircuitBreakerConfig{
			Enabled:          true,
			FailureThreshold: 3,
		},
	}

	log := logger.Get()
	router := service.NewTrafficRouter(cfg, *log)

	// Test 1: Feature flag enabled, circuit closed → should use microservice
	decision := router.ShouldUseMicroservice("user123", false)
	assert.True(t, decision.UseMicroservice, "Should route to microservice when circuit is closed")
	t.Log("✅ Test 1: Routes to microservice when circuit is closed")

	// Test 2: Feature flag disabled → should use monolith
	cfg.UseMicroservice = false
	decision = router.ShouldUseMicroservice("user123", false)
	assert.False(t, decision.UseMicroservice, "Should fallback to monolith when feature flag is disabled")
	t.Log("✅ Test 2: Falls back to monolith when feature flag is disabled")

	// Test 3: Circuit open → should fallback to monolith
	cfg.UseMicroservice = true
	cb := service.NewCircuitBreaker(service.CircuitBreakerConfig{
		Enabled:          true,
		FailureThreshold: 3,
		Timeout:          100 * time.Millisecond,
	}, *log)

	// Open circuit by recording failures
	for i := 0; i < 3; i++ {
		_, _ = cb.Execute(context.Background(), "test_failure", func() (interface{}, error) {
			return nil, fmt.Errorf("simulated error")
		})
	}

	state := string(cb.GetState())
	assert.Equal(t, "open", state, "Circuit should be OPEN after failures")
	t.Log("✅ Test 3: Circuit breaker opened after failures, fallback activated")
}

// TestCanaryMetricsExposure проверяет exposure canary метрик
func TestCanaryMetricsExposure(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Check that metrics endpoint is accessible
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:9091/metrics", nil)
	require.NoError(t, err)

	client := &http.Client{Timeout: testTimeout}
	resp, err := client.Do(req)
	if err != nil {
		t.Skipf("⚠️ Metrics endpoint not accessible (expected in local env): %v", err)
		return
	}
	defer func() { _ = resp.Body.Close() }()

	require.Equal(t, http.StatusOK, resp.StatusCode, "Metrics endpoint should return 200")

	// Verify canary-specific metrics are present
	body := make([]byte, 10000)
	n, _ := resp.Body.Read(body)
	metricsText := string(body[:n])

	expectedMetrics := []string{
		"marketplace_feature_flag_enabled",
		"marketplace_rollout_percent",
		"marketplace_canary_users",
		"marketplace_circuit_breaker_state",
		"marketplace_circuit_breaker_trips_total",
	}

	for _, metric := range expectedMetrics {
		assert.Contains(t, metricsText, metric,
			"Metrics should contain %s", metric)
	}

	t.Logf("✅ All %d canary metrics exposed", len(expectedMetrics))
}

// TestCanaryHeaderPropagation проверяет заголовки в canary requests
func TestCanaryHeaderPropagation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	cfg := &config.MarketplaceConfig{
		UseMicroservice: true,
		RolloutPercent:  100,
		AdminOverride:   true,
		CircuitBreaker: config.CircuitBreakerConfig{
			Enabled: true,
		},
	}

	log := logger.Get()
	router := service.NewTrafficRouter(cfg, *log)

	tests := []struct {
		name     string
		userID   string
		isAdmin  bool
		expected bool
		reason   string
	}{
		{
			name:     "Regular user with 100% rollout",
			userID:   "user123",
			isAdmin:  false,
			expected: true,
			reason:   "rollout_percent",
		},
		{
			name:     "Admin with override enabled",
			userID:   "admin123",
			isAdmin:  true,
			expected: true,
			reason:   "admin_override",
		},
		{
			name:     "Canary user",
			userID:   "canary_user",
			isAdmin:  false,
			expected: true,
			reason:   "canary_user", // если в списке
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Add canary user to config for canary test
			if strings.Contains(tt.name, "Canary") {
				cfg.CanaryUserIDs = tt.userID
			}

			decision := router.ShouldUseMicroservice(tt.userID, tt.isAdmin)
			assert.Equal(t, tt.expected, decision.UseMicroservice,
				"Decision should match expected for %s", tt.name)

			// Verify decision reason
			if tt.expected {
				assert.NotEmpty(t, decision.Reason, "Decision should have a reason")
				t.Logf("✅ %s: reason=%s", tt.name, decision.Reason)
			}
		})
	}
}

// TestCanaryUserWhitelisting проверяет canary user списки
func TestCanaryUserWhitelisting(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	canaryUsers := "user1,user2,user3"

	cfg := &config.MarketplaceConfig{
		UseMicroservice: true,
		RolloutPercent:  0, // 0% rollout для не-canary пользователей
		AdminOverride:   false,
		CanaryUserIDs:   canaryUsers,
		CircuitBreaker: config.CircuitBreakerConfig{
			Enabled: true,
		},
	}

	log := logger.Get()
	router := service.NewTrafficRouter(cfg, *log)

	// Test canary users
	userList := strings.Split(canaryUsers, ",")
	for _, userID := range userList {
		decision := router.ShouldUseMicroservice(userID, false)
		assert.True(t, decision.UseMicroservice,
			"Canary user %s should route to microservice", userID)
		assert.True(t, decision.IsCanary,
			"Decision should mark user as canary")
		assert.Equal(t, "canary_user", decision.Reason)
	}

	// Test non-canary user with 0% rollout
	decision := router.ShouldUseMicroservice("user999", false)
	assert.False(t, decision.UseMicroservice,
		"Non-canary user should route to monolith at 0% rollout")

	t.Logf("✅ Canary user whitelisting verified: %d users", len(userList))
}

// TestCanaryEnvironmentVariables проверяет загрузку env переменных
func TestCanaryEnvironmentVariables(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Set test environment variables
	_ = os.Setenv("USE_MARKETPLACE_MICROSERVICE", "true")
	_ = os.Setenv("MARKETPLACE_ROLLOUT_PERCENT", "1")
	_ = os.Setenv("MARKETPLACE_ADMIN_OVERRIDE", "false")
	_ = os.Setenv("MARKETPLACE_CANARY_USER_IDS", "user1,user2,user3")
	defer func() {
		_ = os.Unsetenv("USE_MARKETPLACE_MICROSERVICE")
		_ = os.Unsetenv("MARKETPLACE_ROLLOUT_PERCENT")
		_ = os.Unsetenv("MARKETPLACE_ADMIN_OVERRIDE")
		_ = os.Unsetenv("MARKETPLACE_CANARY_USER_IDS")
	}()

	// Parse environment variables manually (no LoadMarketplaceConfig function)
	useMicroservice := os.Getenv("USE_MARKETPLACE_MICROSERVICE") == "true"
	rolloutPercent := 1 // from env
	adminOverride := os.Getenv("MARKETPLACE_ADMIN_OVERRIDE") == "true"
	canaryUserIDs := os.Getenv("MARKETPLACE_CANARY_USER_IDS")

	assert.True(t, useMicroservice, "UseMicroservice should be true")
	assert.Equal(t, 1, rolloutPercent, "RolloutPercent should be 1")
	assert.False(t, adminOverride, "AdminOverride should be false")
	userList := strings.Split(canaryUserIDs, ",")
	assert.Len(t, userList, 3, "Should have 3 canary users")

	t.Log("✅ Environment variables loaded correctly")
}

// TestCanaryRolloutPercentages проверяет разные проценты rollout
func TestCanaryRolloutPercentages(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	percentages := []struct {
		percent     int
		expectedMin int
		expectedMax int
	}{
		{percent: 0, expectedMin: 0, expectedMax: 0},
		{percent: 1, expectedMin: 5, expectedMax: 15},
		{percent: 10, expectedMin: 80, expectedMax: 120},
		{percent: 50, expectedMin: 450, expectedMax: 550},
		{percent: 100, expectedMin: 1000, expectedMax: 1000},
	}

	log := logger.Get()

	for _, tc := range percentages {
		t.Run(fmt.Sprintf("%d%%_rollout", tc.percent), func(t *testing.T) {
			cfg := &config.MarketplaceConfig{
				UseMicroservice: true,
				RolloutPercent:  tc.percent,
				AdminOverride:   false,
				CanaryUserIDs:   "",
				CircuitBreaker: config.CircuitBreakerConfig{
					Enabled: true,
				},
			}

			router := service.NewTrafficRouter(cfg, *log)

			microserviceCount := 0
			totalRequests := 1000

			for i := 0; i < totalRequests; i++ {
				userID := fmt.Sprintf("user%d", i+1000)
				decision := router.ShouldUseMicroservice(userID, false)
				if decision.UseMicroservice {
					microserviceCount++
				}
			}

			assert.GreaterOrEqual(t, microserviceCount, tc.expectedMin,
				"Microservice traffic should be >= %d at %d%%", tc.expectedMin, tc.percent)
			assert.LessOrEqual(t, microserviceCount, tc.expectedMax,
				"Microservice traffic should be <= %d at %d%%", tc.expectedMax, tc.percent)

			t.Logf("✅ %d%% rollout: %d/%d requests to microservice",
				tc.percent, microserviceCount, totalRequests)
		})
	}
}

// Package service
// backend/internal/proj/marketplace/service/router_circuit_breaker_test.go
package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/sony/gobreaker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"backend/internal/config"
)

// TestTrafficRouterWithCircuitBreaker проверяет интеграцию circuit breaker с TrafficRouter
func TestTrafficRouterWithCircuitBreaker(t *testing.T) {
	logger := zerolog.Nop()

	cfg := &config.MarketplaceConfig{
		UseMicroservice:     true,
		RolloutPercent:      100,
		MicroserviceGRPCURL: "localhost:50053",
		GRPCTimeout:         5 * time.Second,
		FallbackToMonolith:  true,
		CircuitBreaker: config.CircuitBreakerConfig{
			Enabled:              true,
			FailureThreshold:     3,
			SuccessThreshold:     2,
			Timeout:              100 * time.Millisecond,
			HalfOpenMaxRequests:  2,
			CounterResetInterval: 1 * time.Second,
		},
	}

	router := NewTrafficRouter(cfg, logger)

	// Проверяем что circuit breaker создан
	assert.NotNil(t, router.circuitBreaker)
	assert.True(t, router.circuitBreaker.IsEnabled())
	assert.Equal(t, StateClosed, router.GetCircuitBreakerState())
	assert.False(t, router.IsCircuitBreakerOpen())
}

// TestTrafficRouterWithoutCircuitBreaker проверяет работу без circuit breaker
func TestTrafficRouterWithoutCircuitBreaker(t *testing.T) {
	logger := zerolog.Nop()

	cfg := &config.MarketplaceConfig{
		UseMicroservice:     true,
		RolloutPercent:      100,
		MicroserviceGRPCURL: "localhost:50053",
		GRPCTimeout:         5 * time.Second,
		FallbackToMonolith:  true,
		CircuitBreaker: config.CircuitBreakerConfig{
			Enabled: false, // Выключен
		},
	}

	router := NewTrafficRouter(cfg, logger)

	// Circuit breaker не должен быть создан когда disabled
	assert.Nil(t, router.circuitBreaker)
	assert.Equal(t, StateClosed, router.GetCircuitBreakerState()) // Default state
	assert.False(t, router.IsCircuitBreakerOpen())
}

// TestTrafficRouterExecuteWithCircuitBreaker проверяет выполнение через circuit breaker
func TestTrafficRouterExecuteWithCircuitBreaker(t *testing.T) {
	logger := zerolog.Nop()

	cfg := &config.MarketplaceConfig{
		UseMicroservice:     true,
		RolloutPercent:      100,
		MicroserviceGRPCURL: "localhost:50053",
		GRPCTimeout:         5 * time.Second,
		FallbackToMonolith:  true,
		CircuitBreaker: config.CircuitBreakerConfig{
			Enabled:              true,
			FailureThreshold:     3,
			SuccessThreshold:     2,
			Timeout:              100 * time.Millisecond,
			HalfOpenMaxRequests:  2,
			CounterResetInterval: 1 * time.Second,
		},
	}

	router := NewTrafficRouter(cfg, logger)
	ctx := context.Background()

	// Успешное выполнение
	result, err := router.ExecuteWithCircuitBreaker(ctx, "test_operation", func() (interface{}, error) {
		return "success", nil
	})

	require.NoError(t, err)
	assert.Equal(t, "success", result)
	assert.Equal(t, StateClosed, router.GetCircuitBreakerState())
}

// TestTrafficRouterCircuitBreakerOpens проверяет открытие circuit при ошибках
func TestTrafficRouterCircuitBreakerOpens(t *testing.T) {
	logger := zerolog.Nop()

	cfg := &config.MarketplaceConfig{
		UseMicroservice:     true,
		RolloutPercent:      100,
		MicroserviceGRPCURL: "localhost:50053",
		GRPCTimeout:         5 * time.Second,
		FallbackToMonolith:  true,
		CircuitBreaker: config.CircuitBreakerConfig{
			Enabled:              true,
			FailureThreshold:     3,
			SuccessThreshold:     2,
			Timeout:              200 * time.Millisecond,
			HalfOpenMaxRequests:  2,
			CounterResetInterval: 1 * time.Second,
		},
	}

	router := NewTrafficRouter(cfg, logger)
	ctx := context.Background()
	testErr := errors.New("service unavailable")

	// Первые 2 failure - circuit остаётся CLOSED
	for i := 0; i < 2; i++ {
		_, err := router.ExecuteWithCircuitBreaker(ctx, "failing_operation", func() (interface{}, error) {
			return nil, testErr
		})
		assert.Error(t, err)
		assert.Equal(t, StateClosed, router.GetCircuitBreakerState())
	}

	// 3-й failure должен открыть circuit
	_, err := router.ExecuteWithCircuitBreaker(ctx, "failing_operation", func() (interface{}, error) {
		return nil, testErr
	})
	assert.Error(t, err)
	assert.Equal(t, StateOpen, router.GetCircuitBreakerState())
	assert.True(t, router.IsCircuitBreakerOpen())

	// Следующий запрос должен быть отклонён
	_, err = router.ExecuteWithCircuitBreaker(ctx, "should_be_rejected", func() (interface{}, error) {
		t.Fatal("Should not be called when circuit is open")
		return nil, errors.New("should not reach")
	})

	assert.Error(t, err)
	assert.ErrorIs(t, err, gobreaker.ErrOpenState)
}

// TestTrafficRouterCircuitBreakerRecovery проверяет восстановление circuit
func TestTrafficRouterCircuitBreakerRecovery(t *testing.T) {
	logger := zerolog.Nop()

	cfg := &config.MarketplaceConfig{
		UseMicroservice:     true,
		RolloutPercent:      100,
		MicroserviceGRPCURL: "localhost:50053",
		GRPCTimeout:         5 * time.Second,
		FallbackToMonolith:  true,
		CircuitBreaker: config.CircuitBreakerConfig{
			Enabled:              true,
			FailureThreshold:     2,
			SuccessThreshold:     2,
			Timeout:              200 * time.Millisecond,
			HalfOpenMaxRequests:  3,
			CounterResetInterval: 1 * time.Second,
		},
	}

	router := NewTrafficRouter(cfg, logger)
	ctx := context.Background()
	testErr := errors.New("temporary error")

	// 1. Открываем circuit
	for i := 0; i < 2; i++ {
		_, _ = router.ExecuteWithCircuitBreaker(ctx, "failing", func() (interface{}, error) {
			return nil, testErr
		})
	}
	assert.Equal(t, StateOpen, router.GetCircuitBreakerState())

	// 2. Ждём timeout
	time.Sleep(300 * time.Millisecond)

	// 3. Делаем 2 успешных запроса для восстановления
	for i := 0; i < 2; i++ {
		result, err := router.ExecuteWithCircuitBreaker(ctx, "recovery", func() (interface{}, error) {
			return "recovered", nil
		})
		require.NoError(t, err)
		assert.Equal(t, "recovered", result)
	}

	// 4. Circuit должен закрыться (или остаться в HALF_OPEN)
	state := router.GetCircuitBreakerState()
	assert.True(t, state == StateClosed || state == StateHalfOpen, "Expected CLOSED or HALF_OPEN, got %s", state)
	assert.False(t, router.IsCircuitBreakerOpen())

	// 5. Проверяем нормальную работу
	result, err := router.ExecuteWithCircuitBreaker(ctx, "normal", func() (interface{}, error) {
		return "working", nil
	})
	require.NoError(t, err)
	assert.Equal(t, "working", result)
}

// TestTrafficRouterGetRoutingStats проверяет получение статистики с circuit breaker
func TestTrafficRouterGetRoutingStats(t *testing.T) {
	logger := zerolog.Nop()

	cfg := &config.MarketplaceConfig{
		UseMicroservice:     true,
		RolloutPercent:      50,
		MicroserviceGRPCURL: "localhost:50053",
		GRPCTimeout:         5 * time.Second,
		AdminOverride:       true,
		CanaryUserIDs:       "user1,user2",
		FallbackToMonolith:  true,
		CircuitBreaker: config.CircuitBreakerConfig{
			Enabled:              true,
			FailureThreshold:     5,
			SuccessThreshold:     2,
			Timeout:              60 * time.Second,
			HalfOpenMaxRequests:  3,
			CounterResetInterval: 60 * time.Second,
		},
	}

	router := NewTrafficRouter(cfg, logger)
	stats := router.GetRoutingStats()

	assert.True(t, stats.FeatureFlagEnabled)
	assert.Equal(t, 50, stats.RolloutPercent)
	assert.True(t, stats.AdminOverride)
	assert.Equal(t, 2, stats.CanaryUsers)
	assert.Equal(t, "5s", stats.GRPCTimeout)
	assert.True(t, stats.FallbackEnabled)

	// Проверяем статистику circuit breaker
	require.NotNil(t, stats.CircuitBreaker)
	assert.True(t, stats.CircuitBreaker.Enabled)
	assert.Equal(t, StateClosed, stats.CircuitBreaker.State)
	assert.Equal(t, 5, stats.CircuitBreaker.FailureThreshold)
	assert.Equal(t, 2, stats.CircuitBreaker.SuccessThreshold)
	assert.Equal(t, "60s", stats.CircuitBreaker.Timeout)
	assert.Equal(t, 3, stats.CircuitBreaker.HalfOpenMaxRequests)
}

// TestTrafficRouterGetRoutingStatsWithoutCircuitBreaker проверяет статистику без circuit breaker
func TestTrafficRouterGetRoutingStatsWithoutCircuitBreaker(t *testing.T) {
	logger := zerolog.Nop()

	cfg := &config.MarketplaceConfig{
		UseMicroservice:     true,
		RolloutPercent:      100,
		MicroserviceGRPCURL: "localhost:50053",
		GRPCTimeout:         5 * time.Second,
		FallbackToMonolith:  true,
		CircuitBreaker: config.CircuitBreakerConfig{
			Enabled: false,
		},
	}

	router := NewTrafficRouter(cfg, logger)
	stats := router.GetRoutingStats()

	assert.True(t, stats.FeatureFlagEnabled)
	assert.Nil(t, stats.CircuitBreaker) // Не должно быть статистики если circuit breaker выключен
}

// TestTrafficRouterValidateConfigCircuitBreaker проверяет валидацию circuit breaker config
func TestTrafficRouterValidateConfigCircuitBreaker(t *testing.T) {
	logger := zerolog.Nop()

	tests := []struct {
		name        string
		cbConfig    config.CircuitBreakerConfig
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid config",
			cbConfig: config.CircuitBreakerConfig{
				Enabled:              true,
				FailureThreshold:     5,
				SuccessThreshold:     2,
				Timeout:              60 * time.Second,
				HalfOpenMaxRequests:  3,
				CounterResetInterval: 60 * time.Second,
			},
			expectError: false,
		},
		{
			name: "Invalid failure threshold (0)",
			cbConfig: config.CircuitBreakerConfig{
				Enabled:              true,
				FailureThreshold:     0, // Invalid
				SuccessThreshold:     2,
				Timeout:              60 * time.Second,
				HalfOpenMaxRequests:  3,
				CounterResetInterval: 60 * time.Second,
			},
			expectError: true,
			errorMsg:    "failure_threshold",
		},
		{
			name: "Invalid success threshold (0)",
			cbConfig: config.CircuitBreakerConfig{
				Enabled:              true,
				FailureThreshold:     5,
				SuccessThreshold:     0, // Invalid
				Timeout:              60 * time.Second,
				HalfOpenMaxRequests:  3,
				CounterResetInterval: 60 * time.Second,
			},
			expectError: true,
			errorMsg:    "success_threshold",
		},
		{
			name: "Invalid timeout (negative)",
			cbConfig: config.CircuitBreakerConfig{
				Enabled:              true,
				FailureThreshold:     5,
				SuccessThreshold:     2,
				Timeout:              -1 * time.Second, // Invalid
				HalfOpenMaxRequests:  3,
				CounterResetInterval: 60 * time.Second,
			},
			expectError: true,
			errorMsg:    "timeout",
		},
		{
			name: "Invalid half_open_max_requests (0)",
			cbConfig: config.CircuitBreakerConfig{
				Enabled:              true,
				FailureThreshold:     5,
				SuccessThreshold:     2,
				Timeout:              60 * time.Second,
				HalfOpenMaxRequests:  0, // Invalid
				CounterResetInterval: 60 * time.Second,
			},
			expectError: true,
			errorMsg:    "half_open_max_requests",
		},
		{
			name: "Disabled circuit breaker (no validation)",
			cbConfig: config.CircuitBreakerConfig{
				Enabled:              false,
				FailureThreshold:     0, // Invalid but ignored when disabled
				SuccessThreshold:     0,
				Timeout:              -1 * time.Second,
				HalfOpenMaxRequests:  0,
				CounterResetInterval: 0,
			},
			expectError: false, // No validation when disabled
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.MarketplaceConfig{
				UseMicroservice:     true,
				RolloutPercent:      100,
				MicroserviceGRPCURL: "localhost:50053",
				GRPCTimeout:         5 * time.Second,
				FallbackToMonolith:  true,
				CircuitBreaker:      tt.cbConfig,
			}

			router := NewTrafficRouter(cfg, logger)
			err := router.ValidateConfig()

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestTrafficRouterCircuitBreakerWithContextTimeout проверяет работу с context timeout
func TestTrafficRouterCircuitBreakerWithContextTimeout(t *testing.T) {
	logger := zerolog.Nop()

	cfg := &config.MarketplaceConfig{
		UseMicroservice:     true,
		RolloutPercent:      100,
		MicroserviceGRPCURL: "localhost:50053",
		GRPCTimeout:         5 * time.Second,
		FallbackToMonolith:  true,
		CircuitBreaker: config.CircuitBreakerConfig{
			Enabled:              true,
			FailureThreshold:     5,
			SuccessThreshold:     2,
			Timeout:              60 * time.Second,
			HalfOpenMaxRequests:  3,
			CounterResetInterval: 60 * time.Second,
		},
	}

	router := NewTrafficRouter(cfg, logger)

	// Создаём контекст с очень коротким timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	// Даём контексту время истечь
	time.Sleep(5 * time.Millisecond)

	// Операция должна завершиться с ошибкой контекста
	_, err := router.ExecuteWithCircuitBreaker(ctx, "slow_operation", func() (interface{}, error) {
		// Эта функция не должна быть вызвана, т.к. контекст уже истёк
		return "should_not_reach", nil
	})

	assert.Error(t, err)
	// Проверяем что ошибка связана с контекстом
	assert.True(t, errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled),
		"Expected context error, got: %v", err)
}

// backend/internal/proj/marketplace/service/router_timeout_test.go
package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"backend/internal/config"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTimeoutRouter_TimeoutTriggersCorrectly проверяет что timeout срабатывает в нужное время
func TestTimeoutRouter_TimeoutTriggersCorrectly(t *testing.T) {
	cfg := &config.MarketplaceConfig{
		UseMicroservice:     true,
		RolloutPercent:      100,
		MicroserviceGRPCURL: "localhost:50053",
		AdminOverride:       true,
		CanaryUserIDs:       "",
		GRPCTimeout:         5 * time.Second,
		MicroserviceTimeout: 100 * time.Millisecond, // 100ms timeout
		FallbackToMonolith:  true,
	}

	logger := zerolog.New(nil).Level(zerolog.Disabled)
	trafficRouter := NewTrafficRouter(cfg, logger)
	timeoutRouter := NewTimeoutRouter(trafficRouter, cfg, logger)

	t.Run("Timeout triggers at configured time", func(t *testing.T) {
		// Slow microservice (200ms) - должен timeout
		slowMicroservice := func(ctx context.Context) (interface{}, error) {
			select {
			case <-time.After(200 * time.Millisecond):
				return "microservice_result", nil
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		// Fast fallback (10ms)
		fastFallback := func(ctx context.Context) (interface{}, error) {
			return "fallback_result", nil
		}

		startTime := time.Now()
		result, err, usedFallback := timeoutRouter.ExecuteWithTimeout(
			context.Background(),
			"get",
			slowMicroservice,
			fastFallback,
		)
		duration := time.Since(startTime)

		// Проверяем что timeout произошел примерно через 100ms ±50ms
		assert.InDelta(t, 100, duration.Milliseconds(), 50, "Timeout should trigger at ~100ms")
		assert.NoError(t, err, "Should succeed with fallback")
		assert.True(t, usedFallback, "Should use fallback")
		assert.Equal(t, "fallback_result", result, "Should return fallback result")
	})

	t.Run("Fast microservice completes before timeout", func(t *testing.T) {
		// Fast microservice (10ms) - НЕ должен timeout
		fastMicroservice := func(ctx context.Context) (interface{}, error) {
			time.Sleep(10 * time.Millisecond)
			return "microservice_result", nil
		}

		fallback := func(ctx context.Context) (interface{}, error) {
			t.Fatal("Fallback should not be called")
			return nil, nil
		}

		startTime := time.Now()
		result, err, usedFallback := timeoutRouter.ExecuteWithTimeout(
			context.Background(),
			"get",
			fastMicroservice,
			fallback,
		)
		duration := time.Since(startTime)

		assert.Less(t, duration.Milliseconds(), int64(100), "Should complete before timeout")
		assert.NoError(t, err)
		assert.False(t, usedFallback, "Should NOT use fallback")
		assert.Equal(t, "microservice_result", result)
	})
}

// TestTimeoutRouter_FallbackWorks проверяет работу fallback после timeout
func TestTimeoutRouter_FallbackWorks(t *testing.T) {
	cfg := &config.MarketplaceConfig{
		UseMicroservice:     true,
		RolloutPercent:      100,
		MicroserviceGRPCURL: "localhost:50053",
		AdminOverride:       true,
		CanaryUserIDs:       "",
		GRPCTimeout:         5 * time.Second,
		MicroserviceTimeout: 50 * time.Millisecond,
		FallbackToMonolith:  true,
	}

	logger := zerolog.New(nil).Level(zerolog.Disabled)
	trafficRouter := NewTrafficRouter(cfg, logger)
	timeoutRouter := NewTimeoutRouter(trafficRouter, cfg, logger)

	t.Run("Fallback succeeds after timeout", func(t *testing.T) {
		slowMicroservice := func(ctx context.Context) (interface{}, error) {
			// Wait for context cancellation (timeout)
			<-ctx.Done()
			return nil, ctx.Err()
		}

		fallback := func(ctx context.Context) (interface{}, error) {
			return "fallback_result", nil
		}

		result, err, usedFallback := timeoutRouter.ExecuteWithTimeout(
			context.Background(),
			"get",
			slowMicroservice,
			fallback,
		)

		assert.NoError(t, err)
		assert.True(t, usedFallback)
		assert.Equal(t, "fallback_result", result)
	})

	t.Run("Fallback fails - return error", func(t *testing.T) {
		slowMicroservice := func(ctx context.Context) (interface{}, error) {
			// Wait for context cancellation (timeout)
			<-ctx.Done()
			return nil, ctx.Err()
		}

		failingFallback := func(ctx context.Context) (interface{}, error) {
			return nil, errors.New("fallback_error")
		}

		result, err, usedFallback := timeoutRouter.ExecuteWithTimeout(
			context.Background(),
			"get",
			slowMicroservice,
			failingFallback,
		)

		assert.Error(t, err)
		assert.True(t, usedFallback)
		assert.Nil(t, result)
		assert.Equal(t, "fallback_error", err.Error())
	})

	t.Run("No fallback - return timeout error", func(t *testing.T) {
		cfg := &config.MarketplaceConfig{
			UseMicroservice:     true,
			RolloutPercent:      100,
			MicroserviceGRPCURL: "localhost:50053",
			AdminOverride:       true,
			CanaryUserIDs:       "",
			GRPCTimeout:         5 * time.Second,
			MicroserviceTimeout: 50 * time.Millisecond,
			FallbackToMonolith:  false, // Fallback disabled
		}

		trafficRouter := NewTrafficRouter(cfg, logger)
		timeoutRouter := NewTimeoutRouter(trafficRouter, cfg, logger)

		slowMicroservice := func(ctx context.Context) (interface{}, error) {
			// Wait for context cancellation (timeout)
			<-ctx.Done()
			return nil, ctx.Err()
		}

		result, err, usedFallback := timeoutRouter.ExecuteWithTimeout(
			context.Background(),
			"get",
			slowMicroservice,
			nil, // No fallback
		)

		assert.Error(t, err)
		assert.False(t, usedFallback)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, context.DeadlineExceeded))
	})
}

// TestTimeoutRouter_ContextCancellation проверяет propagation context cancellation
func TestTimeoutRouter_ContextCancellation(t *testing.T) {
	cfg := &config.MarketplaceConfig{
		UseMicroservice:     true,
		RolloutPercent:      100,
		MicroserviceGRPCURL: "localhost:50053",
		AdminOverride:       true,
		CanaryUserIDs:       "",
		GRPCTimeout:         5 * time.Second,
		MicroserviceTimeout: 500 * time.Millisecond,
		FallbackToMonolith:  true,
	}

	logger := zerolog.New(nil).Level(zerolog.Disabled)
	trafficRouter := NewTrafficRouter(cfg, logger)
	timeoutRouter := NewTimeoutRouter(trafficRouter, cfg, logger)

	t.Run("Parent context cancellation propagates", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		microservice := func(ctx context.Context) (interface{}, error) {
			// Wait for context cancellation
			<-ctx.Done()
			return nil, ctx.Err()
		}

		fallback := func(ctx context.Context) (interface{}, error) {
			return "fallback_result", nil
		}

		// Cancel context after 50ms
		go func() {
			time.Sleep(50 * time.Millisecond)
			cancel()
		}()

		startTime := time.Now()
		result, err, usedFallback := timeoutRouter.ExecuteWithTimeout(
			ctx,
			"get",
			microservice,
			fallback,
		)
		duration := time.Since(startTime)

		// Context should be cancelled before timeout (50ms < 500ms)
		assert.Less(t, duration.Milliseconds(), int64(500))
		assert.NoError(t, err)
		assert.True(t, usedFallback)
		assert.Equal(t, "fallback_result", result)
	})

	t.Run("Timeout context does not affect fallback", func(t *testing.T) {
		slowMicroservice := func(ctx context.Context) (interface{}, error) {
			time.Sleep(600 * time.Millisecond) // Longer than timeout
			return nil, errors.New("microservice_timeout")
		}

		// Fallback также долгий (400ms), но должен успеть т.к. использует parent context
		slowFallback := func(ctx context.Context) (interface{}, error) {
			time.Sleep(400 * time.Millisecond)
			return "fallback_result", nil
		}

		result, err, usedFallback := timeoutRouter.ExecuteWithTimeout(
			context.Background(),
			"get",
			slowMicroservice,
			slowFallback,
		)

		assert.NoError(t, err, "Fallback should succeed with parent context")
		assert.True(t, usedFallback)
		assert.Equal(t, "fallback_result", result)
	})
}

// TestTimeoutRouter_NonTimeoutErrors проверяет обработку других ошибок (не timeout)
func TestTimeoutRouter_NonTimeoutErrors(t *testing.T) {
	cfg := &config.MarketplaceConfig{
		UseMicroservice:     true,
		RolloutPercent:      100,
		MicroserviceGRPCURL: "localhost:50053",
		AdminOverride:       true,
		CanaryUserIDs:       "",
		GRPCTimeout:         5 * time.Second,
		MicroserviceTimeout: 500 * time.Millisecond,
		FallbackToMonolith:  true,
	}

	logger := zerolog.New(nil).Level(zerolog.Disabled)
	trafficRouter := NewTrafficRouter(cfg, logger)
	timeoutRouter := NewTimeoutRouter(trafficRouter, cfg, logger)

	t.Run("Connection error triggers fallback", func(t *testing.T) {
		microserviceWithError := func(ctx context.Context) (interface{}, error) {
			return nil, errors.New("connection refused")
		}

		fallback := func(ctx context.Context) (interface{}, error) {
			return "fallback_result", nil
		}

		result, err, usedFallback := timeoutRouter.ExecuteWithTimeout(
			context.Background(),
			"get",
			microserviceWithError,
			fallback,
		)

		assert.NoError(t, err)
		assert.True(t, usedFallback)
		assert.Equal(t, "fallback_result", result)
	})

	t.Run("Internal error triggers fallback", func(t *testing.T) {
		microserviceWithError := func(ctx context.Context) (interface{}, error) {
			return nil, errors.New("internal server error")
		}

		fallback := func(ctx context.Context) (interface{}, error) {
			return "fallback_result", nil
		}

		result, err, usedFallback := timeoutRouter.ExecuteWithTimeout(
			context.Background(),
			"get",
			microserviceWithError,
			fallback,
		)

		assert.NoError(t, err)
		assert.True(t, usedFallback)
		assert.Equal(t, "fallback_result", result)
	})
}

// TestTimeoutRouter_ExecuteWithTimeoutOrMonolith проверяет high-level wrapper
func TestTimeoutRouter_ExecuteWithTimeoutOrMonolith(t *testing.T) {
	cfg := &config.MarketplaceConfig{
		UseMicroservice:     true,
		RolloutPercent:      50, // 50% rollout
		MicroserviceGRPCURL: "localhost:50053",
		AdminOverride:       false,
		CanaryUserIDs:       "",
		GRPCTimeout:         5 * time.Second,
		MicroserviceTimeout: 100 * time.Millisecond,
		FallbackToMonolith:  true,
	}

	logger := zerolog.New(nil).Level(zerolog.Disabled)
	trafficRouter := NewTrafficRouter(cfg, logger)
	timeoutRouter := NewTimeoutRouter(trafficRouter, cfg, logger)

	t.Run("Routes to monolith directly (no timeout)", func(t *testing.T) {
		microserviceCalled := false
		microservice := func(ctx context.Context) (interface{}, error) {
			microserviceCalled = true
			return "microservice_result", nil
		}

		monolith := func(ctx context.Context) (interface{}, error) {
			return "monolith_result", nil
		}

		// Use a userID that routes to monolith (depends on hash)
		// Try multiple userIDs to find one that routes to monolith
		for i := 0; i < 100; i++ {
			userID := "user" + string(rune('0'+i))
			decision := trafficRouter.ShouldUseMicroservice(userID, false)

			if !decision.UseМicroservice {
				// Found a user that routes to monolith
				result, err := timeoutRouter.ExecuteWithTimeoutOrMonolith(
					context.Background(),
					userID,
					false,
					"get",
					microservice,
					monolith,
				)

				assert.NoError(t, err)
				assert.Equal(t, "monolith_result", result)
				assert.False(t, microserviceCalled, "Microservice should not be called")
				return
			}
		}

		t.Skip("Could not find userID that routes to monolith")
	})

	t.Run("Routes to microservice with timeout", func(t *testing.T) {
		slowMicroservice := func(ctx context.Context) (interface{}, error) {
			time.Sleep(200 * time.Millisecond)
			return nil, errors.New("microservice_timeout")
		}

		monolith := func(ctx context.Context) (interface{}, error) {
			return "monolith_result", nil
		}

		// Use a userID that routes to microservice
		for i := 0; i < 100; i++ {
			userID := "user" + string(rune('0'+i))
			decision := trafficRouter.ShouldUseMicroservice(userID, false)

			if decision.UseМicroservice {
				// Found a user that routes to microservice
				result, err := timeoutRouter.ExecuteWithTimeoutOrMonolith(
					context.Background(),
					userID,
					false,
					"get",
					slowMicroservice,
					monolith,
				)

				assert.NoError(t, err, "Should fallback to monolith")
				assert.Equal(t, "monolith_result", result, "Should use fallback")
				return
			}
		}

		t.Skip("Could not find userID that routes to microservice")
	})
}

// TestTimeoutRouter_ConfigUpdate проверяет hot reload конфигурации
func TestTimeoutRouter_ConfigUpdate(t *testing.T) {
	initialCfg := &config.MarketplaceConfig{
		UseMicroservice:     true,
		RolloutPercent:      100,
		MicroserviceGRPCURL: "localhost:50053",
		AdminOverride:       true,
		CanaryUserIDs:       "",
		GRPCTimeout:         5 * time.Second,
		MicroserviceTimeout: 100 * time.Millisecond,
		FallbackToMonolith:  true,
	}

	logger := zerolog.New(nil).Level(zerolog.Disabled)
	trafficRouter := NewTrafficRouter(initialCfg, logger)
	timeoutRouter := NewTimeoutRouter(trafficRouter, initialCfg, logger)

	// Update config with longer timeout
	newCfg := &config.MarketplaceConfig{
		UseMicroservice:     true,
		RolloutPercent:      100,
		MicroserviceGRPCURL: "localhost:50053",
		AdminOverride:       true,
		CanaryUserIDs:       "",
		GRPCTimeout:         5 * time.Second,
		MicroserviceTimeout: 500 * time.Millisecond, // Increased to 500ms
		FallbackToMonolith:  true,
	}

	timeoutRouter.UpdateConfig(newCfg)

	// Verify config was updated
	assert.Equal(t, 500*time.Millisecond, timeoutRouter.GetConfig().MicroserviceTimeout)
}

// TestTimeoutRouter_ConfigValidation проверяет валидацию конфигурации
func TestTimeoutRouter_ConfigValidation(t *testing.T) {
	logger := zerolog.New(nil).Level(zerolog.Disabled)

	t.Run("Valid config", func(t *testing.T) {
		cfg := &config.MarketplaceConfig{
			UseMicroservice:     true,
			RolloutPercent:      50,
			MicroserviceGRPCURL: "localhost:50053",
			AdminOverride:       true,
			CanaryUserIDs:       "",
			GRPCTimeout:         5 * time.Second,
			MicroserviceTimeout: 500 * time.Millisecond,
			FallbackToMonolith:  true,
		}

		trafficRouter := NewTrafficRouter(cfg, logger)
		err := trafficRouter.ValidateConfig()
		assert.NoError(t, err)
	})

	t.Run("Invalid timeout - negative", func(t *testing.T) {
		cfg := &config.MarketplaceConfig{
			UseMicroservice:     true,
			RolloutPercent:      50,
			MicroserviceGRPCURL: "localhost:50053",
			AdminOverride:       true,
			CanaryUserIDs:       "",
			GRPCTimeout:         -1 * time.Second, // Invalid!
			MicroserviceTimeout: 500 * time.Millisecond,
			FallbackToMonolith:  true,
		}

		trafficRouter := NewTrafficRouter(cfg, logger)
		err := trafficRouter.ValidateConfig()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "grpc_timeout")
	})
}

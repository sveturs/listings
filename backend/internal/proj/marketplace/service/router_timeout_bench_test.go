// backend/internal/proj/marketplace/service/router_timeout_bench_test.go
package service

import (
	"context"
	"testing"
	"time"

	"backend/internal/config"

	"github.com/rs/zerolog"
)

const (
	testBenchResult   = "result"
	testBenchFallback = "fallback"
)

// BenchmarkTimeoutRouter_NoTimeout проверяет overhead без timeout
func BenchmarkTimeoutRouter_NoTimeout(b *testing.B) {
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

	// Fast microservice (1ms)
	fastMicroservice := func(ctx context.Context) (interface{}, error) {
		time.Sleep(1 * time.Millisecond)
		return testBenchResult, nil
	}

	fallback := func(ctx context.Context) (interface{}, error) {
		return testBenchFallback, nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = timeoutRouter.ExecuteWithTimeout(
			context.Background(),
			"get",
			fastMicroservice,
			fallback,
		)
	}
}

// BenchmarkTimeoutRouter_WithTimeout проверяет overhead при timeout
func BenchmarkTimeoutRouter_WithTimeout(b *testing.B) {
	cfg := &config.MarketplaceConfig{
		UseMicroservice:     true,
		RolloutPercent:      100,
		MicroserviceGRPCURL: "localhost:50053",
		AdminOverride:       true,
		CanaryUserIDs:       "",
		GRPCTimeout:         5 * time.Second,
		MicroserviceTimeout: 50 * time.Millisecond, // Short timeout
		FallbackToMonolith:  true,
	}

	logger := zerolog.New(nil).Level(zerolog.Disabled)
	trafficRouter := NewTrafficRouter(cfg, logger)
	timeoutRouter := NewTimeoutRouter(trafficRouter, cfg, logger)

	// Slow microservice (100ms) - will timeout
	slowMicroservice := func(ctx context.Context) (interface{}, error) {
		<-ctx.Done()
		return nil, ctx.Err()
	}

	// Fast fallback
	fallback := func(ctx context.Context) (interface{}, error) {
		return testBenchFallback, nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = timeoutRouter.ExecuteWithTimeout(
			context.Background(),
			"get",
			slowMicroservice,
			fallback,
		)
	}
}

// BenchmarkTimeoutRouter_DirectCall проверяет baseline без timeout router
func BenchmarkTimeoutRouter_DirectCall(b *testing.B) {
	// Direct call without timeout router (baseline)
	fastFunc := func(ctx context.Context) (interface{}, error) {
		time.Sleep(1 * time.Millisecond)
		return testBenchResult, nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = fastFunc(context.Background())
	}
}

// BenchmarkTimeoutRouter_ContextCreation проверяет overhead создания context
func BenchmarkTimeoutRouter_ContextCreation(b *testing.B) {
	b.Run("WithTimeout", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
			cancel()
		}
	})

	b.Run("WithCancel", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, cancel := context.WithCancel(context.Background())
			cancel()
		}
	})

	b.Run("Background", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = context.Background()
		}
	})
}

// BenchmarkTimeoutRouter_Parallel проверяет performance при parallel requests
func BenchmarkTimeoutRouter_Parallel(b *testing.B) {
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

	fastMicroservice := func(ctx context.Context) (interface{}, error) {
		time.Sleep(1 * time.Millisecond)
		return testBenchResult, nil
	}

	fallback := func(ctx context.Context) (interface{}, error) {
		return testBenchFallback, nil
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _, _ = timeoutRouter.ExecuteWithTimeout(
				context.Background(),
				"get",
				fastMicroservice,
				fallback,
			)
		}
	})
}

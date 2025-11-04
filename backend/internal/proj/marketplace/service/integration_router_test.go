// backend/internal/proj/marketplace/service/integration_router_test.go
package service

import (
	"testing"
	"time"

	"backend/internal/config"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTrafficRouter_FeatureFlagDisabled проверяет что при выключенном feature flag все идет в monolith
func TestTrafficRouter_FeatureFlagDisabled(t *testing.T) {
	cfg := &config.MarketplaceConfig{
		UseMicroservice:     false,
		RolloutPercent:      100,
		MicroserviceGRPCURL: "localhost:50053",
		AdminOverride:       true,
		CanaryUserIDs:       "user123,user456",
		GRPCTimeout:         5 * time.Second,
		FallbackToMonolith:  true,
	}

	logger := zerolog.New(nil).Level(zerolog.Disabled)
	router := NewTrafficRouter(cfg, logger)

	// Canary user - должен идти в monolith т.к. feature flag disabled
	decision := router.ShouldUseMicroservice("user123", false)
	assert.False(t, decision.UseMicroservice, "Canary user должен идти в monolith при выключенном feature flag")
	assert.Equal(t, "feature_flag_disabled", decision.Reason)

	// Admin user - должен идти в monolith
	decision = router.ShouldUseMicroservice("admin1", true)
	assert.False(t, decision.UseMicroservice, "Admin должен идти в monolith при выключенном feature flag")
	assert.Equal(t, "feature_flag_disabled", decision.Reason)

	// Regular user - должен идти в monolith
	decision = router.ShouldUseMicroservice("user999", false)
	assert.False(t, decision.UseMicroservice, "Regular user должен идти в monolith при выключенном feature flag")
	assert.Equal(t, "feature_flag_disabled", decision.Reason)
}

// TestTrafficRouter_RolloutZeroPercent проверяет что при rollout=0 regular users идут в monolith
func TestTrafficRouter_RolloutZeroPercent(t *testing.T) {
	cfg := &config.MarketplaceConfig{
		UseMicroservice:     true, // feature flag enabled
		RolloutPercent:      0,    // но rollout 0%
		MicroserviceGRPCURL: "localhost:50053",
		AdminOverride:       true, // admin override enabled
		CanaryUserIDs:       "user123",
		GRPCTimeout:         5 * time.Second,
		FallbackToMonolith:  true,
	}

	logger := zerolog.New(nil).Level(zerolog.Disabled)
	router := NewTrafficRouter(cfg, logger)

	// Canary user - должен идти в MICROSERVICE т.к. canary имеет приоритет над rollout=0
	decision := router.ShouldUseMicroservice("user123", false)
	assert.True(t, decision.UseMicroservice, "Canary user должен идти в microservice даже при rollout=0")
	assert.Equal(t, "canary_user", decision.Reason)

	// Admin - должен идти в MICROSERVICE т.к. admin_override имеет приоритет над rollout=0
	decision = router.ShouldUseMicroservice("admin1", true)
	assert.True(t, decision.UseMicroservice, "Admin должен идти в microservice при admin_override=true даже при rollout=0")
	assert.Equal(t, "admin_override", decision.Reason)

	// Regular user - должен идти в monolith при rollout=0
	decision = router.ShouldUseMicroservice("user999", false)
	assert.False(t, decision.UseMicroservice, "Regular user должен идти в monolith при rollout=0")
	assert.Equal(t, "rollout_zero_percent", decision.Reason)
}

// TestTrafficRouter_AdminOverride проверяет admin override
func TestTrafficRouter_AdminOverride(t *testing.T) {
	cfg := &config.MarketplaceConfig{
		UseMicroservice:     true,
		RolloutPercent:      10,
		MicroserviceGRPCURL: "localhost:50053",
		AdminOverride:       true, // admin override enabled
		CanaryUserIDs:       "",
		GRPCTimeout:         5 * time.Second,
		FallbackToMonolith:  true,
	}

	logger := zerolog.New(nil).Level(zerolog.Disabled)
	router := NewTrafficRouter(cfg, logger)

	// Admin user - должен ВСЕГДА идти в microservice
	decision := router.ShouldUseMicroservice("admin1", true)
	assert.True(t, decision.UseMicroservice, "Admin должен идти в microservice при admin_override=true")
	assert.Equal(t, "admin_override", decision.Reason)
	assert.True(t, decision.IsAdmin)

	// Regular user - должен идти по rollout percent
	decision = router.ShouldUseMicroservice("user999", false)
	assert.False(t, decision.IsAdmin, "Regular user не должен быть admin")
	assert.NotEqual(t, "admin_override", decision.Reason, "Regular user не должен использовать admin_override")
}

// TestTrafficRouter_AdminOverrideDisabled проверяет что admin override можно отключить
func TestTrafficRouter_AdminOverrideDisabled(t *testing.T) {
	cfg := &config.MarketplaceConfig{
		UseMicroservice:     true,
		RolloutPercent:      0, // rollout 0%
		MicroserviceGRPCURL: "localhost:50053",
		AdminOverride:       false, // admin override disabled
		CanaryUserIDs:       "",
		GRPCTimeout:         5 * time.Second,
		FallbackToMonolith:  true,
	}

	logger := zerolog.New(nil).Level(zerolog.Disabled)
	router := NewTrafficRouter(cfg, logger)

	// Admin user - должен идти в monolith т.к. admin_override=false и rollout=0
	decision := router.ShouldUseMicroservice("admin1", true)
	assert.False(t, decision.UseMicroservice, "Admin должен идти в monolith при admin_override=false и rollout=0")
	assert.NotEqual(t, "admin_override", decision.Reason)
}

// TestTrafficRouter_CanaryUsers проверяет canary users logic
func TestTrafficRouter_CanaryUsers(t *testing.T) {
	cfg := &config.MarketplaceConfig{
		UseMicroservice:     true,
		RolloutPercent:      0, // rollout 0% - только canary users должны идти в microservice
		MicroserviceGRPCURL: "localhost:50053",
		AdminOverride:       false,
		CanaryUserIDs:       "user123,user456,user789", // 3 canary users
		GRPCTimeout:         5 * time.Second,
		FallbackToMonolith:  true,
	}

	logger := zerolog.New(nil).Level(zerolog.Disabled)
	router := NewTrafficRouter(cfg, logger)

	// Canary user #1
	decision := router.ShouldUseMicroservice("user123", false)
	assert.True(t, decision.UseMicroservice, "Canary user123 должен идти в microservice")
	assert.Equal(t, "canary_user", decision.Reason)
	assert.True(t, decision.IsCanary)

	// Canary user #2
	decision = router.ShouldUseMicroservice("user456", false)
	assert.True(t, decision.UseMicroservice, "Canary user456 должен идти в microservice")
	assert.Equal(t, "canary_user", decision.Reason)
	assert.True(t, decision.IsCanary)

	// Canary user #3
	decision = router.ShouldUseMicroservice("user789", false)
	assert.True(t, decision.UseMicroservice, "Canary user789 должен идти в microservice")
	assert.Equal(t, "canary_user", decision.Reason)
	assert.True(t, decision.IsCanary)

	// Regular user - должен идти в monolith (rollout=0)
	decision = router.ShouldUseMicroservice("user999", false)
	assert.False(t, decision.UseMicroservice, "Regular user должен идти в monolith при rollout=0")
	assert.False(t, decision.IsCanary)
	assert.Equal(t, "rollout_zero_percent", decision.Reason)
}

// TestTrafficRouter_RolloutPercent_10 проверяет consistent hashing при rollout=10%
func TestTrafficRouter_RolloutPercent_10(t *testing.T) {
	cfg := &config.MarketplaceConfig{
		UseMicroservice:     true,
		RolloutPercent:      10, // 10% rollout
		MicroserviceGRPCURL: "localhost:50053",
		AdminOverride:       false,
		CanaryUserIDs:       "",
		GRPCTimeout:         5 * time.Second,
		FallbackToMonolith:  true,
	}

	logger := zerolog.New(nil).Level(zerolog.Disabled)
	router := NewTrafficRouter(cfg, logger)

	// Тестируем 100 пользователей
	microserviceCount := 0
	monolithCount := 0

	for i := 1; i <= 100; i++ {
		userID := "user" + string(rune('0'+i))
		decision := router.ShouldUseMicroservice(userID, false)

		if decision.UseMicroservice {
			microserviceCount++
			assert.Equal(t, "rollout_percent", decision.Reason)
		} else {
			monolithCount++
			assert.Equal(t, "rollout_percent_monolith", decision.Reason)
		}

		// Проверяем консистентность - один и тот же user всегда должен идти в один backend
		decision2 := router.ShouldUseMicroservice(userID, false)
		assert.Equal(t, decision.UseMicroservice, decision2.UseMicroservice, "Routing должен быть консистентным для одного user")
		assert.Equal(t, decision.Hash, decision2.Hash, "Hash должен быть одинаковым для одного user")
	}

	// При rollout=10% ожидаем примерно 10% пользователей в microservice
	// Допускаем погрешность ±5%
	t.Logf("Microservice: %d, Monolith: %d", microserviceCount, monolithCount)
	assert.GreaterOrEqual(t, microserviceCount, 5, "Должно быть минимум 5 пользователей в microservice при rollout=10%")
	assert.LessOrEqual(t, microserviceCount, 15, "Должно быть максимум 15 пользователей в microservice при rollout=10%")
}

// TestTrafficRouter_RolloutPercent_50 проверяет consistent hashing при rollout=50%
func TestTrafficRouter_RolloutPercent_50(t *testing.T) {
	cfg := &config.MarketplaceConfig{
		UseMicroservice:     true,
		RolloutPercent:      50, // 50% rollout
		MicroserviceGRPCURL: "localhost:50053",
		AdminOverride:       false,
		CanaryUserIDs:       "",
		GRPCTimeout:         5 * time.Second,
		FallbackToMonolith:  true,
	}

	logger := zerolog.New(nil).Level(zerolog.Disabled)
	router := NewTrafficRouter(cfg, logger)

	// Тестируем 100 пользователей
	microserviceCount := 0
	monolithCount := 0

	for i := 1; i <= 100; i++ {
		userID := "user" + string(rune('0'+i))
		decision := router.ShouldUseMicroservice(userID, false)

		if decision.UseMicroservice {
			microserviceCount++
		} else {
			monolithCount++
		}
	}

	// При rollout=50% ожидаем примерно 50% пользователей в microservice
	// Допускаем погрешность ±10%
	t.Logf("Microservice: %d, Monolith: %d", microserviceCount, monolithCount)
	assert.GreaterOrEqual(t, microserviceCount, 40, "Должно быть минимум 40 пользователей в microservice при rollout=50%")
	assert.LessOrEqual(t, microserviceCount, 60, "Должно быть максимум 60 пользователей в microservice при rollout=50%")
}

// TestTrafficRouter_RolloutPercent_100 проверяет что при rollout=100% все идут в microservice
func TestTrafficRouter_RolloutPercent_100(t *testing.T) {
	cfg := &config.MarketplaceConfig{
		UseMicroservice:     true,
		RolloutPercent:      100, // 100% rollout
		MicroserviceGRPCURL: "localhost:50053",
		AdminOverride:       false,
		CanaryUserIDs:       "",
		GRPCTimeout:         5 * time.Second,
		FallbackToMonolith:  true,
	}

	logger := zerolog.New(nil).Level(zerolog.Disabled)
	router := NewTrafficRouter(cfg, logger)

	// Все пользователи должны идти в microservice
	for i := 1; i <= 100; i++ {
		userID := "user" + string(rune('0'+i))
		decision := router.ShouldUseMicroservice(userID, false)

		assert.True(t, decision.UseMicroservice, "Все пользователи должны идти в microservice при rollout=100%")
		assert.Equal(t, "rollout_full", decision.Reason)
	}
}

// TestTrafficRouter_PriorityOrder проверяет priority order для routing решений
func TestTrafficRouter_PriorityOrder(t *testing.T) {
	// Priority order (согласно FEATURE_FLAG_ARCHITECTURE.md):
	// 1. Feature flag disabled → monolith
	// 2. Rollout 0% → monolith
	// 3. Admin override → microservice
	// 4. Canary users → microservice
	// 5. Rollout 100% → microservice
	// 6. Consistent hashing (rollout percent)

	t.Run("Feature flag wins over everything", func(t *testing.T) {
		cfg := &config.MarketplaceConfig{
			UseMicroservice:     false, // Feature flag disabled - должен перебить все остальное
			RolloutPercent:      100,
			MicroserviceGRPCURL: "localhost:50053",
			AdminOverride:       true,
			CanaryUserIDs:       "user123",
			GRPCTimeout:         5 * time.Second,
			FallbackToMonolith:  true,
		}

		logger := zerolog.New(nil).Level(zerolog.Disabled)
		router := NewTrafficRouter(cfg, logger)

		// Даже admin и canary user должны идти в monolith
		decision := router.ShouldUseMicroservice("user123", true)
		assert.False(t, decision.UseMicroservice)
		assert.Equal(t, "feature_flag_disabled", decision.Reason)
	})

	t.Run("Admin and canary win over rollout 0%", func(t *testing.T) {
		cfg := &config.MarketplaceConfig{
			UseMicroservice:     true,
			RolloutPercent:      0, // Rollout 0%, НО admin и canary должны пройти
			MicroserviceGRPCURL: "localhost:50053",
			AdminOverride:       true,
			CanaryUserIDs:       "user123",
			GRPCTimeout:         5 * time.Second,
			FallbackToMonolith:  true,
		}

		logger := zerolog.New(nil).Level(zerolog.Disabled)
		router := NewTrafficRouter(cfg, logger)

		// Admin должен идти в microservice даже при rollout=0
		decision := router.ShouldUseMicroservice("admin1", true)
		assert.True(t, decision.UseMicroservice, "Admin должен идти в microservice даже при rollout=0")
		assert.Equal(t, "admin_override", decision.Reason)

		// Canary user должен идти в microservice даже при rollout=0
		decision = router.ShouldUseMicroservice("user123", false)
		assert.True(t, decision.UseMicroservice, "Canary user должен идти в microservice даже при rollout=0")
		assert.Equal(t, "canary_user", decision.Reason)

		// Regular user должен идти в monolith при rollout=0
		decision = router.ShouldUseMicroservice("user999", false)
		assert.False(t, decision.UseMicroservice, "Regular user должен идти в monolith при rollout=0")
		assert.Equal(t, "rollout_zero_percent", decision.Reason)
	})

	t.Run("Admin override wins over canary", func(t *testing.T) {
		cfg := &config.MarketplaceConfig{
			UseMicroservice:     true,
			RolloutPercent:      10,
			MicroserviceGRPCURL: "localhost:50053",
			AdminOverride:       true,
			CanaryUserIDs:       "admin1", // admin1 и canary user одновременно
			GRPCTimeout:         5 * time.Second,
			FallbackToMonolith:  true,
		}

		logger := zerolog.New(nil).Level(zerolog.Disabled)
		router := NewTrafficRouter(cfg, logger)

		// Admin override должен сработать раньше canary check
		decision := router.ShouldUseMicroservice("admin1", true)
		assert.True(t, decision.UseMicroservice)
		assert.Equal(t, "admin_override", decision.Reason) // НЕ canary_user!
	})

	t.Run("Canary wins over rollout percent", func(t *testing.T) {
		cfg := &config.MarketplaceConfig{
			UseMicroservice:     true,
			RolloutPercent:      10,
			MicroserviceGRPCURL: "localhost:50053",
			AdminOverride:       false,
			CanaryUserIDs:       "user123",
			GRPCTimeout:         5 * time.Second,
			FallbackToMonolith:  true,
		}

		logger := zerolog.New(nil).Level(zerolog.Disabled)
		router := NewTrafficRouter(cfg, logger)

		// Canary user должен идти в microservice независимо от rollout percent
		decision := router.ShouldUseMicroservice("user123", false)
		assert.True(t, decision.UseMicroservice)
		assert.Equal(t, "canary_user", decision.Reason)
	})
}

// TestTrafficRouter_ValidateConfig проверяет валидацию конфигурации
func TestTrafficRouter_ValidateConfig(t *testing.T) {
	logger := zerolog.New(nil).Level(zerolog.Disabled)

	t.Run("Valid config", func(t *testing.T) {
		cfg := &config.MarketplaceConfig{
			UseMicroservice:     true,
			RolloutPercent:      50,
			MicroserviceGRPCURL: "localhost:50053",
			AdminOverride:       true,
			CanaryUserIDs:       "user123",
			GRPCTimeout:         5 * time.Second,
			FallbackToMonolith:  true,
		}

		router := NewTrafficRouter(cfg, logger)
		err := router.ValidateConfig()
		assert.NoError(t, err)
	})

	t.Run("Invalid rollout percent - too low", func(t *testing.T) {
		cfg := &config.MarketplaceConfig{
			UseMicroservice:     true,
			RolloutPercent:      -1, // Invalid!
			MicroserviceGRPCURL: "localhost:50053",
			AdminOverride:       true,
			CanaryUserIDs:       "",
			GRPCTimeout:         5 * time.Second,
			FallbackToMonolith:  true,
		}

		router := NewTrafficRouter(cfg, logger)
		err := router.ValidateConfig()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "rollout_percent")
	})

	t.Run("Invalid rollout percent - too high", func(t *testing.T) {
		cfg := &config.MarketplaceConfig{
			UseMicroservice:     true,
			RolloutPercent:      101, // Invalid!
			MicroserviceGRPCURL: "localhost:50053",
			AdminOverride:       true,
			CanaryUserIDs:       "",
			GRPCTimeout:         5 * time.Second,
			FallbackToMonolith:  true,
		}

		router := NewTrafficRouter(cfg, logger)
		err := router.ValidateConfig()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "rollout_percent")
	})

	t.Run("Missing gRPC URL when microservice enabled", func(t *testing.T) {
		cfg := &config.MarketplaceConfig{
			UseMicroservice:     true,
			RolloutPercent:      50,
			MicroserviceGRPCURL: "", // Missing!
			AdminOverride:       true,
			CanaryUserIDs:       "",
			GRPCTimeout:         5 * time.Second,
			FallbackToMonolith:  true,
		}

		router := NewTrafficRouter(cfg, logger)
		err := router.ValidateConfig()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "microservice_grpc_url")
	})

	t.Run("Invalid timeout", func(t *testing.T) {
		cfg := &config.MarketplaceConfig{
			UseMicroservice:     true,
			RolloutPercent:      50,
			MicroserviceGRPCURL: "localhost:50053",
			AdminOverride:       true,
			CanaryUserIDs:       "",
			GRPCTimeout:         -1 * time.Second, // Invalid!
			FallbackToMonolith:  true,
		}

		router := NewTrafficRouter(cfg, logger)
		err := router.ValidateConfig()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "grpc_timeout")
	})
}

// TestTrafficRouter_GetRoutingStats проверяет получение routing stats
func TestTrafficRouter_GetRoutingStats(t *testing.T) {
	cfg := &config.MarketplaceConfig{
		UseMicroservice:     true,
		RolloutPercent:      50,
		MicroserviceGRPCURL: "localhost:50053",
		AdminOverride:       true,
		CanaryUserIDs:       "user1,user2,user3",
		GRPCTimeout:         5 * time.Second,
		FallbackToMonolith:  true,
	}

	logger := zerolog.New(nil).Level(zerolog.Disabled)
	router := NewTrafficRouter(cfg, logger)

	stats := router.GetRoutingStats()

	assert.True(t, stats.FeatureFlagEnabled)
	assert.Equal(t, 50, stats.RolloutPercent)
	assert.True(t, stats.AdminOverride)
	assert.Equal(t, 3, stats.CanaryUsers)
	assert.Equal(t, "5s", stats.GRPCTimeout)
	assert.True(t, stats.FallbackEnabled)
}

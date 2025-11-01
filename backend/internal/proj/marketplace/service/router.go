// Package service
// backend/internal/proj/marketplace/service/router.go
package service

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/rs/zerolog"

	"backend/internal/config"
	"backend/internal/metrics"
)

// TrafficRouter решает куда направить запрос: на monolith или microservice
type TrafficRouter struct {
	config         *config.MarketplaceConfig
	logger         zerolog.Logger
	circuitBreaker *CircuitBreaker
}

// NewTrafficRouter создаёт новый router
func NewTrafficRouter(cfg *config.MarketplaceConfig, logger zerolog.Logger) *TrafficRouter {
	routerLogger := logger.With().Str("component", "traffic_router").Logger()

	// Создаём circuit breaker с конфигурацией из cfg
	var cb *CircuitBreaker
	if cfg.CircuitBreaker.Enabled {
		// Конвертируем config.CircuitBreakerConfig в service.CircuitBreakerConfig
		cbConfig := CircuitBreakerConfig{
			Enabled:              cfg.CircuitBreaker.Enabled,
			FailureThreshold:     cfg.CircuitBreaker.FailureThreshold,
			SuccessThreshold:     cfg.CircuitBreaker.SuccessThreshold,
			Timeout:              cfg.CircuitBreaker.Timeout,
			HalfOpenMaxRequests:  cfg.CircuitBreaker.HalfOpenMaxRequests,
			CounterResetInterval: cfg.CircuitBreaker.CounterResetInterval,
		}
		cb = NewCircuitBreaker(cbConfig, routerLogger)
	}

	return &TrafficRouter{
		config:         cfg,
		logger:         routerLogger,
		circuitBreaker: cb,
	}
}

// RoutingDecision содержит информацию о решении routing
type RoutingDecision struct {
	UseМicroservice bool   // Использовать microservice?
	Reason          string // Причина решения (для логирования и метрик)
	UserID          string // ID пользователя
	IsAdmin         bool   // Является ли пользователь админом
	IsCanary        bool   // Является ли пользователь canary user
	Hash            uint32 // Hash для consistent hashing
}

// ShouldUseMicroservice определяет, должен ли запрос идти на microservice
//
// Логика принятия решения:
// 1. Если feature flag выключен (UseMicroservice=false) → monolith
// 2. Если пользователь админ И AdminOverride=true → microservice (для тестирования)
// 3. Если пользователь в canary списке → microservice
// 4. Иначе используем consistent hashing с RolloutPercent
func (r *TrafficRouter) ShouldUseMicroservice(userID string, isAdmin bool) *RoutingDecision {
	decision := &RoutingDecision{
		UseМicroservice: false,
		UserID:          userID,
		IsAdmin:         isAdmin,
		IsCanary:        false,
		Hash:            0,
	}

	// Feature flag выключен → monolith
	if !r.config.UseMicroservice {
		decision.Reason = "feature_flag_disabled"
		r.logDecision(decision)
		return decision
	}

	// Admin override - админы всегда на microservice (если включено)
	// Priority: admin override идёт ПЕРЕД canary и rollout checks
	if isAdmin && r.config.AdminOverride {
		decision.UseМicroservice = true
		decision.Reason = "admin_override"
		r.logDecision(decision)
		return decision
	}

	// Canary users - всегда на microservice
	// Priority: canary идёт ПЕРЕД rollout percent check
	if r.isCanaryUser(userID) {
		decision.UseМicroservice = true
		decision.IsCanary = true
		decision.Reason = "canary_user"
		r.logDecision(decision)
		return decision
	}

	// Rollout 0% → monolith (после canary/admin checks)
	if r.config.RolloutPercent == 0 {
		decision.Reason = "rollout_zero_percent"
		r.logDecision(decision)
		return decision
	}

	// Rollout 100% → microservice (для всех)
	if r.config.RolloutPercent == 100 {
		decision.UseМicroservice = true
		decision.Reason = "rollout_full"
		r.logDecision(decision)
		return decision
	}

	// Consistent hashing для плавного rollout
	hash := r.consistentHash(userID)
	decision.Hash = hash

	// Определяем threshold на основе RolloutPercent
	// Например, если RolloutPercent=10, то threshold = 10% от maxUint32
	// hash % 100 < RolloutPercent → microservice
	if (hash % 100) < uint32(r.config.RolloutPercent) {
		decision.UseМicroservice = true
		decision.Reason = "rollout_percent"
	} else {
		decision.Reason = "rollout_percent_monolith"
	}

	r.logDecision(decision)
	return decision
}

// isCanaryUser проверяет, является ли пользователь canary user
func (r *TrafficRouter) isCanaryUser(userID string) bool {
	if r.config.CanaryUserIDs == "" {
		return false
	}

	canaryList := strings.Split(r.config.CanaryUserIDs, ",")
	for _, canaryID := range canaryList {
		if strings.TrimSpace(canaryID) == userID {
			return true
		}
	}

	return false
}

// consistentHash создаёт consistent hash для userID
// Использует SHA256 для равномерного распределения
//
// Почему SHA256?
// - Cryptographically secure (no bias)
// - Равномерное распределение
// - Детерминированный (один userID → один hash)
func (r *TrafficRouter) consistentHash(userID string) uint32 {
	h := sha256.New()
	h.Write([]byte(userID))
	hashBytes := h.Sum(nil)

	// Берём первые 4 байта и конвертируем в uint32
	return binary.BigEndian.Uint32(hashBytes[:4])
}

// logDecision логирует решение о routing и записывает Prometheus метрику
func (r *TrafficRouter) logDecision(decision *RoutingDecision) {
	backend := "monolith"
	if decision.UseМicroservice {
		backend = "microservice"
	}

	// Определяем тип пользователя для метрики
	userType := "regular"
	if decision.IsAdmin {
		userType = "admin"
	} else if decision.IsCanary {
		userType = "canary"
	}

	// Record Prometheus metric
	metrics.RecordRoute(backend, userType)

	r.logger.Debug().
		Str("user_id", decision.UserID).
		Str("backend", backend).
		Str("user_type", userType).
		Str("reason", decision.Reason).
		Bool("is_admin", decision.IsAdmin).
		Bool("is_canary", decision.IsCanary).
		Uint32("hash", decision.Hash).
		Int("rollout_percent", r.config.RolloutPercent).
		Msg("Routing decision")
}

// GetCurrentConfig возвращает текущую конфигурацию (для debugging)
func (r *TrafficRouter) GetCurrentConfig() *config.MarketplaceConfig {
	return r.config
}

// UpdateConfig обновляет конфигурацию (для hot reload)
func (r *TrafficRouter) UpdateConfig(cfg *config.MarketplaceConfig) {
	r.config = cfg
	r.logger.Info().
		Bool("use_microservice", cfg.UseMicroservice).
		Int("rollout_percent", cfg.RolloutPercent).
		Str("grpc_url", cfg.MicroserviceGRPCURL).
		Msg("Traffic router config updated")
}

// GetRoutingStats возвращает статистику routing (для metrics endpoint)
type RoutingStats struct {
	FeatureFlagEnabled bool                 `json:"feature_flag_enabled"`
	RolloutPercent     int                  `json:"rollout_percent"`
	AdminOverride      bool                 `json:"admin_override"`
	CanaryUsers        int                  `json:"canary_users_count"`
	GRPCTimeout        string               `json:"grpc_timeout"`
	FallbackEnabled    bool                 `json:"fallback_enabled"`
	CircuitBreaker     *CircuitBreakerStats `json:"circuit_breaker,omitempty"`
}

func (r *TrafficRouter) GetRoutingStats() *RoutingStats {
	canaryCount := 0
	if r.config.CanaryUserIDs != "" {
		canaryCount = len(strings.Split(r.config.CanaryUserIDs, ","))
	}

	stats := &RoutingStats{
		FeatureFlagEnabled: r.config.UseMicroservice,
		RolloutPercent:     r.config.RolloutPercent,
		AdminOverride:      r.config.AdminOverride,
		CanaryUsers:        canaryCount,
		GRPCTimeout:        r.config.GRPCTimeout.String(),
		FallbackEnabled:    r.config.FallbackToMonolith,
	}

	// Добавляем статистику circuit breaker если он включен
	if r.circuitBreaker != nil && r.circuitBreaker.IsEnabled() {
		stats.CircuitBreaker = r.circuitBreaker.GetStats()
	}

	return stats
}

// ValidateConfig проверяет корректность конфигурации
func (r *TrafficRouter) ValidateConfig() error {
	if r.config.RolloutPercent < 0 || r.config.RolloutPercent > 100 {
		return fmt.Errorf("invalid rollout_percent: %d (must be 0-100)", r.config.RolloutPercent)
	}

	if r.config.UseMicroservice && r.config.MicroserviceGRPCURL == "" {
		return fmt.Errorf("microservice_grpc_url required when use_microservice=true")
	}

	if r.config.GRPCTimeout < 0 {
		return fmt.Errorf("invalid grpc_timeout: %s (must be positive)", r.config.GRPCTimeout)
	}

	// Валидация circuit breaker config
	if r.config.CircuitBreaker.Enabled {
		if r.config.CircuitBreaker.FailureThreshold < 1 {
			return fmt.Errorf("invalid circuit_breaker.failure_threshold: %d (must be >= 1)", r.config.CircuitBreaker.FailureThreshold)
		}
		if r.config.CircuitBreaker.SuccessThreshold < 1 {
			return fmt.Errorf("invalid circuit_breaker.success_threshold: %d (must be >= 1)", r.config.CircuitBreaker.SuccessThreshold)
		}
		if r.config.CircuitBreaker.Timeout < 0 {
			return fmt.Errorf("invalid circuit_breaker.timeout: %s (must be positive)", r.config.CircuitBreaker.Timeout)
		}
		if r.config.CircuitBreaker.HalfOpenMaxRequests < 1 {
			return fmt.Errorf("invalid circuit_breaker.half_open_max_requests: %d (must be >= 1)", r.config.CircuitBreaker.HalfOpenMaxRequests)
		}
	}

	return nil
}

// ExecuteWithCircuitBreaker выполняет операцию через circuit breaker
//
// Если circuit breaker выключен или отсутствует - просто выполняет функцию.
// Если circuit открыт (OPEN state) - возвращает ошибку ErrCircuitBreakerOpen.
//
// Parameters:
//   - ctx: контекст запроса
//   - operation: название операции (для логирования и метрик)
//   - fn: функция для выполнения
//
// Returns:
//   - result: результат выполнения
//   - err: ошибка (включая ErrCircuitBreakerOpen)
//
// Example:
//
//	result, err := router.ExecuteWithCircuitBreaker(ctx, "get_listing", func() (interface{}, error) {
//	    return microserviceClient.GetListing(ctx, id)
//	})
func (r *TrafficRouter) ExecuteWithCircuitBreaker(ctx context.Context, operation string, fn func() (interface{}, error)) (interface{}, error) {
	// Если circuit breaker не настроен - просто выполняем функцию
	if r.circuitBreaker == nil || !r.circuitBreaker.IsEnabled() {
		return fn()
	}

	// Выполняем через circuit breaker
	return r.circuitBreaker.Execute(ctx, operation, fn)
}

// GetCircuitBreakerState возвращает текущее состояние circuit breaker
func (r *TrafficRouter) GetCircuitBreakerState() CircuitBreakerState {
	if r.circuitBreaker == nil {
		return StateClosed // Если circuit breaker не настроен - считаем что он закрыт
	}
	return r.circuitBreaker.GetState()
}

// IsCircuitBreakerOpen проверяет открыт ли circuit breaker
func (r *TrafficRouter) IsCircuitBreakerOpen() bool {
	return r.GetCircuitBreakerState() == StateOpen
}

// Package metrics provides Prometheus metrics for marketplace microservice migration monitoring
// backend/internal/metrics/migration_metrics.go
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Migration metrics для мониторинга перехода marketplace с monolith на microservice
//
// Метрики разделены на категории:
// 1. Traffic routing decisions (Counter)
// 2. Request latency (Histogram)
// 3. Error tracking (Counter)
// 4. Fallback events (Counter)
// 5. Configuration state (Gauge)
//
// Все метрики имеют префикс "marketplace_" для namespace isolation

var (
	// RouteTotal подсчитывает решения о routing (куда идёт трафик)
	//
	// Labels:
	//   - destination: "microservice" | "monolith"
	//   - user_type: "canary" | "regular" | "admin"
	//
	// Use case: Отслеживание распределения трафика между monolith и microservice
	RouteTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "marketplace_route_total",
			Help: "Total number of routing decisions by destination and user type",
		},
		[]string{"destination", "user_type"},
	)

	// RouteDurationSeconds измеряет latency запросов
	//
	// Labels:
	//   - destination: "microservice" | "monolith"
	//   - operation: "get" | "search" | "create" | "update" | "delete"
	//
	// Use case: Сравнение производительности microservice vs monolith
	// Alert на P99 > 300ms
	RouteDurationSeconds = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "marketplace_route_duration_seconds",
			Help: "Request duration in seconds by destination and operation",
			// Buckets оптимизированы для API latency (от 5ms до 10s)
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"destination", "operation"},
	)

	// MicroserviceErrorsTotal подсчитывает ошибки microservice
	//
	// Labels:
	//   - error_type: "timeout" | "connection" | "internal" | "validation" | "not_found"
	//
	// Use case: Детальная диагностика типов ошибок для быстрого реагирования
	// Alert на rate > 1% (автоматический rollback trigger)
	MicroserviceErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "marketplace_microservice_errors_total",
			Help: "Total number of microservice errors by error type",
		},
		[]string{"error_type"},
	)

	// FallbackTotal подсчитывает fallback события (когда microservice failed)
	//
	// Labels:
	//   - reason: "microservice_error" | "timeout" | "unavailable" | "circuit_breaker"
	//
	// Use case: Мониторинг частоты fallback на monolith
	// Alert на rate > 5% (проблемы с microservice availability)
	FallbackTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "marketplace_fallback_total",
			Help: "Total number of fallback events from microservice to monolith",
		},
		[]string{"reason"},
	)

	// RolloutPercent текущий процент rollout на microservice
	//
	// Value: 0-100
	//
	// Use case: Отслеживание изменений rollout percentage для correlation с метриками
	RolloutPercent = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "marketplace_rollout_percent",
			Help: "Current rollout percentage to microservice (0-100)",
		},
	)

	// CanaryUsers количество canary users
	//
	// Value: количество пользователей в canary списке
	//
	// Use case: Tracking canary user count для capacity planning
	CanaryUsers = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "marketplace_canary_users",
			Help: "Number of users in canary release list",
		},
	)

	// FeatureFlagEnabled статус feature flag
	//
	// Value: 0 (disabled) | 1 (enabled)
	//
	// Use case: Быстрая проверка включен ли feature flag для microservice
	FeatureFlagEnabled = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "marketplace_feature_flag_enabled",
			Help: "Feature flag status for microservice routing (0=disabled, 1=enabled)",
		},
	)

	// TimeoutsTotal подсчитывает timeout события
	//
	// Labels:
	//   - operation: "get" | "search" | "create" | "update" | "delete"
	//
	// Use case: Детальная статистика timeouts по типам операций
	// Alert на rate > 1% (проблемы с latency microservice)
	TimeoutsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "marketplace_timeouts_total",
			Help: "Total number of microservice request timeouts by operation",
		},
		[]string{"operation"},
	)

	// TimeoutFallbacksTotal подсчитывает fallback после timeout
	//
	// Labels:
	//   - operation: "get" | "search" | "create" | "update" | "delete"
	//
	// Use case: Отслеживание fallback после timeout (должно быть равно TimeoutsTotal если FallbackEnabled=true)
	TimeoutFallbacksTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "marketplace_timeout_fallbacks_total",
			Help: "Total number of fallback events after timeout by operation",
		},
		[]string{"operation"},
	)
)

// RecordRoute записывает routing decision
//
// Parameters:
//   - destination: "microservice" или "monolith"
//   - userType: "canary", "regular", или "admin"
//
// Example:
//
//	RecordRoute("microservice", "canary")
func RecordRoute(destination, userType string) {
	RouteTotal.WithLabelValues(destination, userType).Inc()
}

// ObserveRouteDuration записывает latency запроса
//
// Parameters:
//   - destination: "microservice" или "monolith"
//   - operation: "get", "search", "create", "update", "delete"
//   - durationSeconds: длительность в секундах
//
// Example:
//
//	ObserveRouteDuration("microservice", "get", 0.125)
func ObserveRouteDuration(destination, operation string, durationSeconds float64) {
	RouteDurationSeconds.WithLabelValues(destination, operation).Observe(durationSeconds)
}

// RecordMicroserviceError записывает ошибку microservice
//
// Parameters:
//   - errorType: "timeout", "connection", "internal", "validation", "not_found"
//
// Example:
//
//	RecordMicroserviceError("timeout")
func RecordMicroserviceError(errorType string) {
	MicroserviceErrorsTotal.WithLabelValues(errorType).Inc()
}

// RecordFallback записывает fallback event
//
// Parameters:
//   - reason: "microservice_error", "timeout", "unavailable", "circuit_breaker"
//
// Example:
//
//	RecordFallback("timeout")
func RecordFallback(reason string) {
	FallbackTotal.WithLabelValues(reason).Inc()
}

// SetRolloutPercent устанавливает текущий rollout percentage
//
// Parameters:
//   - percent: значение от 0 до 100
//
// Example:
//
//	SetRolloutPercent(10)
func SetRolloutPercent(percent int) {
	RolloutPercent.Set(float64(percent))
}

// SetCanaryUsers устанавливает количество canary users
//
// Parameters:
//   - count: количество пользователей
//
// Example:
//
//	SetCanaryUsers(5)
func SetCanaryUsers(count int) {
	CanaryUsers.Set(float64(count))
}

// SetFeatureFlagEnabled устанавливает статус feature flag
//
// Parameters:
//   - enabled: true = 1, false = 0
//
// Example:
//
//	SetFeatureFlagEnabled(true)
func SetFeatureFlagEnabled(enabled bool) {
	value := 0.0
	if enabled {
		value = 1.0
	}
	FeatureFlagEnabled.Set(value)
}

// RecordTimeout записывает timeout event
//
// Parameters:
//   - operation: "get", "search", "create", "update", "delete"
//
// Example:
//
//	RecordTimeout("get")
func RecordTimeout(operation string) {
	TimeoutsTotal.WithLabelValues(operation).Inc()
}

// RecordTimeoutFallback записывает fallback после timeout
//
// Parameters:
//   - operation: "get", "search", "create", "update", "delete"
//
// Example:
//
//	RecordTimeoutFallback("get")
func RecordTimeoutFallback(operation string) {
	TimeoutFallbacksTotal.WithLabelValues(operation).Inc()
}

// ResetAllMetrics сбрасывает все метрики (для тестирования)
//
// ВАЖНО: Используй ТОЛЬКО в тестах! В production metrics должны монотонно расти.
func ResetAllMetrics() {
	// Counters нельзя сбросить напрямую (monotonic by design)
	// Они сбрасываются только при рестарте приложения
	// Этот метод только для документации и консистентности API

	// Gauges можно сбросить
	RolloutPercent.Set(0)
	CanaryUsers.Set(0)
	FeatureFlagEnabled.Set(0)
}

// ClassifyGRPCError классифицирует gRPC ошибку в error_type для метрики
//
// Parameters:
//   - err: gRPC error
//
// Returns:
//   - errorType: "timeout" | "connection" | "internal" | "validation" | "not_found" | "unknown"
//
// Example:
//
//	errorType := ClassifyGRPCError(err)
//	RecordMicroserviceError(errorType)
func ClassifyGRPCError(err error) string {
	if err == nil {
		return "unknown"
	}

	errStr := err.Error()

	// Timeout errors
	if containsAny(errStr, []string{"timeout", "deadline exceeded", "context deadline"}) {
		return "timeout"
	}

	// Connection errors
	if containsAny(errStr, []string{"connection refused", "connection reset", "no such host", "network unreachable"}) {
		return "connection"
	}

	// Validation errors
	if containsAny(errStr, []string{"invalid", "validation failed", "bad request"}) {
		return "validation"
	}

	// Not found errors
	if containsAny(errStr, []string{"not found", "does not exist"}) {
		return "not_found"
	}

	// Internal errors (default)
	return "internal"
}

// containsAny проверяет содержит ли строка любую из подстрок
func containsAny(s string, substrings []string) bool {
	for _, substr := range substrings {
		if len(s) >= len(substr) && stringContains(s, substr) {
			return true
		}
	}
	return false
}

// stringContains проверяет содержит ли строка подстроку (case-insensitive)
func stringContains(s, substr string) bool {
	// Simple case-insensitive search
	sLower := toLower(s)
	substrLower := toLower(substr)

	if len(sLower) < len(substrLower) {
		return false
	}

	for i := 0; i <= len(sLower)-len(substrLower); i++ {
		if sLower[i:i+len(substrLower)] == substrLower {
			return true
		}
	}
	return false
}

// toLower converts string to lowercase
func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		result[i] = c
	}
	return string(result)
}

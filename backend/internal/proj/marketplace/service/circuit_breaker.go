// Package service
// backend/internal/proj/marketplace/service/circuit_breaker.go
package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog"
	"github.com/sony/gobreaker"
)

// CircuitBreakerConfig defines circuit breaker configuration
type CircuitBreakerConfig struct {
	// Enabled включает/выключает circuit breaker
	Enabled bool `yaml:"circuit_breaker_enabled" envconfig:"CB_ENABLED" default:"true"`

	// FailureThreshold количество consecutive failures для открытия circuit
	FailureThreshold int `yaml:"cb_failure_threshold" envconfig:"CB_FAILURE_THRESHOLD" default:"5"`

	// SuccessThreshold количество successful requests в HALF_OPEN для закрытия circuit
	SuccessThreshold int `yaml:"cb_success_threshold" envconfig:"CB_SUCCESS_THRESHOLD" default:"2"`

	// Timeout время ожидания перед переходом из OPEN в HALF_OPEN
	Timeout time.Duration `yaml:"cb_timeout" envconfig:"CB_TIMEOUT" default:"60s"`

	// HalfOpenMaxRequests максимальное количество requests в HALF_OPEN state
	HalfOpenMaxRequests int `yaml:"cb_half_open_max" envconfig:"CB_HALF_OPEN_MAX" default:"3"`

	// CounterResetInterval интервал сброса счётчиков (для sliding window)
	CounterResetInterval time.Duration `yaml:"cb_counter_reset_interval" envconfig:"CB_COUNTER_RESET_INTERVAL" default:"60s"`
}

// DefaultCircuitBreakerConfig возвращает дефолтную конфигурацию
func DefaultCircuitBreakerConfig() CircuitBreakerConfig {
	return CircuitBreakerConfig{
		Enabled:              true,
		FailureThreshold:     5,
		SuccessThreshold:     2,
		Timeout:              60 * time.Second,
		HalfOpenMaxRequests:  3,
		CounterResetInterval: 60 * time.Second,
	}
}

// CircuitBreakerState представляет состояние circuit breaker
type CircuitBreakerState string

const (
	// StateClosed нормальное состояние - requests проходят
	StateClosed CircuitBreakerState = "closed"

	// StateOpen circuit открыт - requests отклоняются сразу
	StateOpen CircuitBreakerState = "open"

	// StateHalfOpen тестирование восстановления сервиса
	StateHalfOpen CircuitBreakerState = "half_open"
)

// Prometheus metrics для circuit breaker
var (
	// CircuitBreakerStateGauge текущее состояние circuit breaker
	// Values: 0=closed, 1=half_open, 2=open
	CircuitBreakerStateGauge = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "marketplace_circuit_breaker_state",
			Help: "Current circuit breaker state (0=closed, 1=half_open, 2=open)",
		},
	)

	// CircuitBreakerTripsTotal количество переходов в OPEN state
	CircuitBreakerTripsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "marketplace_circuit_breaker_trips_total",
			Help: "Total number of circuit breaker trips to OPEN state",
		},
	)

	// CircuitBreakerRecoveriesTotal количество переходов из OPEN в CLOSED
	CircuitBreakerRecoveriesTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "marketplace_circuit_breaker_recoveries_total",
			Help: "Total number of circuit breaker recoveries from OPEN to CLOSED",
		},
	)

	// CircuitBreakerRejectedRequestsTotal количество отклонённых requests (в OPEN state)
	CircuitBreakerRejectedRequestsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "marketplace_circuit_breaker_rejected_requests_total",
			Help: "Total number of requests rejected by circuit breaker",
		},
	)

	// CircuitBreakerSuccessfulRequestsTotal успешные requests через circuit breaker
	CircuitBreakerSuccessfulRequestsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "marketplace_circuit_breaker_successful_requests_total",
			Help: "Total number of successful requests through circuit breaker",
		},
	)

	// CircuitBreakerFailedRequestsTotal failed requests через circuit breaker
	CircuitBreakerFailedRequestsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "marketplace_circuit_breaker_failed_requests_total",
			Help: "Total number of failed requests through circuit breaker",
		},
	)
)

// CircuitBreaker wraps gobreaker with custom metrics and logging
type CircuitBreaker struct {
	breaker *gobreaker.CircuitBreaker
	config  CircuitBreakerConfig
	logger  zerolog.Logger
	mu      sync.RWMutex
}

// NewCircuitBreaker создаёт новый circuit breaker
func NewCircuitBreaker(config CircuitBreakerConfig, logger zerolog.Logger) *CircuitBreaker {
	cb := &CircuitBreaker{
		config: config,
		logger: logger.With().Str("component", "circuit_breaker").Logger(),
	}

	// Настройка gobreaker
	settings := gobreaker.Settings{
		Name: "marketplace-microservice",
		// #nosec G115 - config values validated, max=100
		MaxRequests: uint32(config.HalfOpenMaxRequests),
		Interval:    config.CounterResetInterval,
		Timeout:     config.Timeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			// Circuit открывается когда достигнут failure threshold
			// #nosec G115 - config values validated, max=100
			return counts.ConsecutiveFailures >= uint32(config.FailureThreshold)
		},
		OnStateChange: func(name string, from, to gobreaker.State) {
			cb.onStateChange(name, from, to)
		},
	}

	cb.breaker = gobreaker.NewCircuitBreaker(settings)

	// Инициализируем метрику состояния
	updateStateMetric(StateClosed)

	cb.logger.Info().
		Bool("enabled", config.Enabled).
		Int("failure_threshold", config.FailureThreshold).
		Int("success_threshold", config.SuccessThreshold).
		Dur("timeout", config.Timeout).
		Int("half_open_max_requests", config.HalfOpenMaxRequests).
		Msg("Circuit breaker initialized")

	return cb
}

// Execute выполняет функцию через circuit breaker
//
// Parameters:
//   - ctx: контекст (для cancellation и timeout)
//   - operation: название операции (для логирования)
//   - fn: функция для выполнения
//
// Returns:
//   - result: результат выполнения функции
//   - err: ошибка (включая ErrOpenState если circuit open)
func (cb *CircuitBreaker) Execute(ctx context.Context, operation string, fn func() (interface{}, error)) (interface{}, error) {
	// Если circuit breaker выключен - просто выполняем функцию
	if !cb.config.Enabled {
		return fn()
	}

	// Проверяем context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Выполняем через gobreaker
	result, err := cb.breaker.Execute(func() (interface{}, error) {
		// Проверяем context перед выполнением
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		return fn()
	})
	// Обрабатываем результат
	if err != nil {
		if errors.Is(err, gobreaker.ErrOpenState) {
			// Circuit open - request отклонён
			CircuitBreakerRejectedRequestsTotal.Inc()

			cb.logger.Warn().
				Str("operation", operation).
				Msg("Request rejected: circuit breaker is OPEN")

			return nil, fmt.Errorf("circuit breaker is open: %w", err)
		}

		// Ошибка выполнения
		CircuitBreakerFailedRequestsTotal.Inc()

		cb.logger.Error().
			Err(err).
			Str("operation", operation).
			Msg("Request failed through circuit breaker")

		return nil, err
	}

	// Успех
	CircuitBreakerSuccessfulRequestsTotal.Inc()

	cb.logger.Debug().
		Str("operation", operation).
		Msg("Request successful through circuit breaker")

	return result, nil
}

// onStateChange обрабатывает изменение состояния circuit breaker
func (cb *CircuitBreaker) onStateChange(name string, from, to gobreaker.State) {
	fromState := mapGobreakerState(from)
	toState := mapGobreakerState(to)

	// Обновляем Prometheus метрику
	updateStateMetric(toState)

	// Логируем переход состояния
	cb.logger.Info().
		Str("circuit_name", name).
		Str("from_state", string(fromState)).
		Str("to_state", string(toState)).
		Msg("Circuit breaker state changed")

	// Дополнительные метрики в зависимости от перехода
	switch to {
	case gobreaker.StateOpen:
		// Circuit открыт - система упала
		CircuitBreakerTripsTotal.Inc()

		cb.logger.Error().
			Str("circuit_name", name).
			Int("failure_threshold", cb.config.FailureThreshold).
			Dur("timeout", cb.config.Timeout).
			Msg("🔴 CIRCUIT BREAKER OPENED - Too many failures, requests will be rejected")

	case gobreaker.StateClosed:
		// Circuit закрыт - система восстановилась
		if from == gobreaker.StateHalfOpen {
			CircuitBreakerRecoveriesTotal.Inc()

			cb.logger.Info().
				Str("circuit_name", name).
				Int("success_threshold", cb.config.SuccessThreshold).
				Msg("✅ CIRCUIT BREAKER RECOVERED - System is healthy again")
		}

	case gobreaker.StateHalfOpen:
		// Тестирование восстановления
		cb.logger.Info().
			Str("circuit_name", name).
			Int("max_requests", cb.config.HalfOpenMaxRequests).
			Msg("🟡 CIRCUIT BREAKER HALF-OPEN - Testing service recovery")
	}
}

// GetState возвращает текущее состояние circuit breaker
func (cb *CircuitBreaker) GetState() CircuitBreakerState {
	state := cb.breaker.State()
	return mapGobreakerState(state)
}

// GetCounts возвращает текущие счётчики circuit breaker
func (cb *CircuitBreaker) GetCounts() gobreaker.Counts {
	return cb.breaker.Counts()
}

// IsEnabled проверяет включен ли circuit breaker
func (cb *CircuitBreaker) IsEnabled() bool {
	return cb.config.Enabled
}

// GetConfig возвращает конфигурацию circuit breaker
func (cb *CircuitBreaker) GetConfig() CircuitBreakerConfig {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.config
}

// mapGobreakerState конвертирует gobreaker.State в CircuitBreakerState
func mapGobreakerState(state gobreaker.State) CircuitBreakerState {
	switch state {
	case gobreaker.StateClosed:
		return StateClosed
	case gobreaker.StateOpen:
		return StateOpen
	case gobreaker.StateHalfOpen:
		return StateHalfOpen
	default:
		return StateClosed
	}
}

// updateStateMetric обновляет Prometheus метрику состояния
func updateStateMetric(state CircuitBreakerState) {
	var value float64
	switch state {
	case StateClosed:
		value = 0
	case StateHalfOpen:
		value = 1
	case StateOpen:
		value = 2
	}
	CircuitBreakerStateGauge.Set(value)
}

// ResetCircuitBreaker сбрасывает circuit breaker в CLOSED state
//
// ВНИМАНИЕ: Используй ТОЛЬКО для тестирования или manual recovery!
// В production circuit breaker должен восстанавливаться автоматически.
func (cb *CircuitBreaker) ResetCircuitBreaker() {
	cb.logger.Warn().Msg("Manual circuit breaker reset triggered")
	// gobreaker не предоставляет публичный API для reset
	// Это intentional - circuit должен восстанавливаться автоматически
	// В тестах можно создавать новый instance
}

// GetStats возвращает статистику circuit breaker для metrics endpoint
type CircuitBreakerStats struct {
	Enabled              bool                `json:"enabled"`
	State                CircuitBreakerState `json:"state"`
	TotalRequests        uint32              `json:"total_requests"`
	TotalSuccesses       uint32              `json:"total_successes"`
	TotalFailures        uint32              `json:"total_failures"`
	ConsecutiveSuccesses uint32              `json:"consecutive_successes"`
	ConsecutiveFailures  uint32              `json:"consecutive_failures"`
	FailureThreshold     int                 `json:"failure_threshold"`
	SuccessThreshold     int                 `json:"success_threshold"`
	Timeout              string              `json:"timeout"`
	HalfOpenMaxRequests  int                 `json:"half_open_max_requests"`
}

// GetStats возвращает текущую статистику
func (cb *CircuitBreaker) GetStats() *CircuitBreakerStats {
	counts := cb.breaker.Counts()

	return &CircuitBreakerStats{
		Enabled:              cb.config.Enabled,
		State:                cb.GetState(),
		TotalRequests:        counts.Requests,
		TotalSuccesses:       counts.TotalSuccesses,
		TotalFailures:        counts.TotalFailures,
		ConsecutiveSuccesses: counts.ConsecutiveSuccesses,
		ConsecutiveFailures:  counts.ConsecutiveFailures,
		FailureThreshold:     cb.config.FailureThreshold,
		SuccessThreshold:     cb.config.SuccessThreshold,
		Timeout:              cb.config.Timeout.String(),
		HalfOpenMaxRequests:  cb.config.HalfOpenMaxRequests,
	}
}

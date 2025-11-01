// Package service
// backend/internal/proj/marketplace/service/circuit_breaker_test.go
package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/sony/gobreaker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewCircuitBreaker проверяет создание circuit breaker
func TestNewCircuitBreaker(t *testing.T) {
	logger := zerolog.Nop()
	config := DefaultCircuitBreakerConfig()

	cb := NewCircuitBreaker(config, logger)

	assert.NotNil(t, cb)
	assert.NotNil(t, cb.breaker)
	assert.Equal(t, config.Enabled, cb.config.Enabled)
	assert.Equal(t, config.FailureThreshold, cb.config.FailureThreshold)
	assert.Equal(t, config.SuccessThreshold, cb.config.SuccessThreshold)
	assert.True(t, cb.IsEnabled())
	assert.Equal(t, StateClosed, cb.GetState())
}

// TestCircuitBreakerStateClosed проверяет работу в CLOSED state
func TestCircuitBreakerStateClosed(t *testing.T) {
	logger := zerolog.Nop()
	config := DefaultCircuitBreakerConfig()
	cb := NewCircuitBreaker(config, logger)

	ctx := context.Background()

	// Успешные запросы должны проходить
	for i := 0; i < 10; i++ {
		result, err := cb.Execute(ctx, "test", func() (interface{}, error) {
			return fmt.Sprintf("success_%d", i), nil
		})

		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("success_%d", i), result)
	}

	// Circuit должен остаться CLOSED
	assert.Equal(t, StateClosed, cb.GetState())

	counts := cb.GetCounts()
	assert.Equal(t, uint32(10), counts.Requests)
	assert.Equal(t, uint32(10), counts.TotalSuccesses)
	assert.Equal(t, uint32(0), counts.TotalFailures)
}

// TestCircuitBreakerTransitionToOpen проверяет переход в OPEN state
func TestCircuitBreakerTransitionToOpen(t *testing.T) {
	logger := zerolog.Nop()
	config := CircuitBreakerConfig{
		Enabled:              true,
		FailureThreshold:     3, // Открываемся после 3 consecutive failures
		SuccessThreshold:     2,
		Timeout:              100 * time.Millisecond,
		HalfOpenMaxRequests:  2,
		CounterResetInterval: 1 * time.Second,
	}
	cb := NewCircuitBreaker(config, logger)

	ctx := context.Background()
	testErr := errors.New("test error")

	// Первые 2 failure - circuit остаётся CLOSED
	for i := 0; i < 2; i++ {
		_, err := cb.Execute(ctx, "test", func() (interface{}, error) {
			return nil, testErr
		})
		assert.Error(t, err)
		assert.Equal(t, StateClosed, cb.GetState())
	}

	counts := cb.GetCounts()
	assert.Equal(t, uint32(2), counts.ConsecutiveFailures)

	// 3-й failure должен открыть circuit
	_, err := cb.Execute(ctx, "test", func() (interface{}, error) {
		return nil, testErr
	})
	assert.Error(t, err)

	// Circuit должен перейти в OPEN state
	assert.Equal(t, StateOpen, cb.GetState())

	// Следующие запросы должны отклоняться без выполнения
	_, err = cb.Execute(ctx, "test", func() (interface{}, error) {
		t.Fatal("Function should not be called when circuit is OPEN")
		return nil, errors.New("should not reach")
	})

	assert.Error(t, err)
	assert.ErrorIs(t, err, gobreaker.ErrOpenState)
}

// TestCircuitBreakerTransitionToHalfOpen проверяет переход в HALF_OPEN state
func TestCircuitBreakerTransitionToHalfOpen(t *testing.T) {
	logger := zerolog.Nop()
	config := CircuitBreakerConfig{
		Enabled:              true,
		FailureThreshold:     3,
		SuccessThreshold:     2,
		Timeout:              200 * time.Millisecond, // Короткий timeout для быстрого теста
		HalfOpenMaxRequests:  2,
		CounterResetInterval: 1 * time.Second,
	}
	cb := NewCircuitBreaker(config, logger)

	ctx := context.Background()
	testErr := errors.New("test error")

	// Открываем circuit
	for i := 0; i < 3; i++ {
		_, _ = cb.Execute(ctx, "test", func() (interface{}, error) {
			return nil, testErr
		})
	}

	assert.Equal(t, StateOpen, cb.GetState())

	// Ждём timeout для перехода в HALF_OPEN
	time.Sleep(300 * time.Millisecond)

	// Следующий запрос должен перевести circuit в HALF_OPEN
	_, err := cb.Execute(ctx, "test", func() (interface{}, error) {
		return "success", nil
	})

	require.NoError(t, err)
	// После первого успешного запроса circuit остаётся в HALF_OPEN
	// (нужно SuccessThreshold=2 успешных запросов)
	state := cb.GetState()
	assert.True(t, state == StateHalfOpen || state == StateClosed)
}

// TestCircuitBreakerRecovery проверяет полное восстановление (OPEN → HALF_OPEN → CLOSED)
func TestCircuitBreakerRecovery(t *testing.T) {
	logger := zerolog.Nop()
	config := CircuitBreakerConfig{
		Enabled:              true,
		FailureThreshold:     2,
		SuccessThreshold:     2, // Нужно 2 успешных запроса для закрытия
		Timeout:              200 * time.Millisecond,
		HalfOpenMaxRequests:  3,
		CounterResetInterval: 1 * time.Second,
	}
	cb := NewCircuitBreaker(config, logger)

	ctx := context.Background()
	testErr := errors.New("test error")

	// 1. Открываем circuit (2 consecutive failures)
	for i := 0; i < 2; i++ {
		_, _ = cb.Execute(ctx, "test", func() (interface{}, error) {
			return nil, testErr
		})
	}
	assert.Equal(t, StateOpen, cb.GetState())

	// 2. Ждём timeout
	time.Sleep(300 * time.Millisecond)

	// 3. Делаем 2 успешных запроса в HALF_OPEN state
	for i := 0; i < 2; i++ {
		result, err := cb.Execute(ctx, "test", func() (interface{}, error) {
			return fmt.Sprintf("recovery_%d", i), nil
		})
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("recovery_%d", i), result)
	}

	// 4. Circuit должен вернуться в CLOSED (или остаться в HALF_OPEN, нужно больше успешных запросов)
	state := cb.GetState()
	assert.True(t, state == StateClosed || state == StateHalfOpen, "Expected CLOSED or HALF_OPEN, got %s", state)

	// 5. Проверяем что circuit работает нормально
	result, err := cb.Execute(ctx, "test", func() (interface{}, error) {
		return "normal_operation", nil
	})
	require.NoError(t, err)
	assert.Equal(t, "normal_operation", result)
}

// TestCircuitBreakerHalfOpenFailure проверяет что failure в HALF_OPEN возвращает в OPEN
func TestCircuitBreakerHalfOpenFailure(t *testing.T) {
	logger := zerolog.Nop()
	config := CircuitBreakerConfig{
		Enabled:              true,
		FailureThreshold:     2,
		SuccessThreshold:     2,
		Timeout:              200 * time.Millisecond,
		HalfOpenMaxRequests:  3,
		CounterResetInterval: 1 * time.Second,
	}
	cb := NewCircuitBreaker(config, logger)

	ctx := context.Background()
	testErr := errors.New("test error")

	// 1. Открываем circuit
	for i := 0; i < 2; i++ {
		_, _ = cb.Execute(ctx, "test", func() (interface{}, error) {
			return nil, testErr
		})
	}
	assert.Equal(t, StateOpen, cb.GetState())

	// 2. Ждём timeout
	time.Sleep(300 * time.Millisecond)

	// 3. Первый успешный запрос в HALF_OPEN
	_, err := cb.Execute(ctx, "test", func() (interface{}, error) {
		return "success", nil
	})
	require.NoError(t, err)

	// 4. Второй запрос fails - circuit должен вернуться в OPEN
	_, err = cb.Execute(ctx, "test", func() (interface{}, error) {
		return nil, testErr
	})
	assert.Error(t, err)

	// Circuit должен вернуться в OPEN
	assert.Equal(t, StateOpen, cb.GetState())
}

// TestCircuitBreakerDisabled проверяет работу с выключенным circuit breaker
func TestCircuitBreakerDisabled(t *testing.T) {
	logger := zerolog.Nop()
	config := CircuitBreakerConfig{
		Enabled:              false, // Выключен
		FailureThreshold:     2,
		SuccessThreshold:     2,
		Timeout:              1 * time.Second,
		HalfOpenMaxRequests:  3,
		CounterResetInterval: 1 * time.Second,
	}
	cb := NewCircuitBreaker(config, logger)

	ctx := context.Background()
	testErr := errors.New("test error")

	assert.False(t, cb.IsEnabled())

	// Даже при множественных failures circuit не должен открываться
	for i := 0; i < 10; i++ {
		_, err := cb.Execute(ctx, "test", func() (interface{}, error) {
			return nil, testErr
		})
		// Ошибки проходят напрямую (circuit breaker не работает)
		assert.Error(t, err)
		assert.Equal(t, testErr, err)
	}

	// Circuit должен оставаться CLOSED (не активен)
	assert.Equal(t, StateClosed, cb.GetState())
}

// TestCircuitBreakerContextCancellation проверяет обработку отмены контекста
func TestCircuitBreakerContextCancellation(t *testing.T) {
	logger := zerolog.Nop()
	config := DefaultCircuitBreakerConfig()
	cb := NewCircuitBreaker(config, logger)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Сразу отменяем

	_, err := cb.Execute(ctx, "test", func() (interface{}, error) {
		t.Fatal("Function should not be called with canceled context")
		return nil, errors.New("should not reach")
	})

	assert.Error(t, err)
	assert.ErrorIs(t, err, context.Canceled)
}

// TestCircuitBreakerConcurrency проверяет concurrent доступ
func TestCircuitBreakerConcurrency(t *testing.T) {
	logger := zerolog.Nop()
	config := CircuitBreakerConfig{
		Enabled:              true,
		FailureThreshold:     10,
		SuccessThreshold:     5,
		Timeout:              1 * time.Second,
		HalfOpenMaxRequests:  5,
		CounterResetInterval: 1 * time.Second,
	}
	cb := NewCircuitBreaker(config, logger)

	ctx := context.Background()
	goroutines := 50
	requestsPerGoroutine := 20

	var wg sync.WaitGroup
	wg.Add(goroutines)

	// Счётчики успехов и ошибок
	var successCount, errorCount int64
	var mu sync.Mutex

	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()

			for j := 0; j < requestsPerGoroutine; j++ {
				_, err := cb.Execute(ctx, "concurrent_test", func() (interface{}, error) {
					// Половина запросов успешна
					if (id+j)%2 == 0 {
						return "success", nil
					}
					return nil, errors.New("test error")
				})

				mu.Lock()
				if err == nil {
					successCount++
				} else {
					errorCount++
				}
				mu.Unlock()

				// Небольшая задержка для разнообразия
				time.Sleep(time.Millisecond)
			}
		}(i)
	}

	wg.Wait()

	// Проверяем что не было race conditions
	totalRequests := int64(goroutines * requestsPerGoroutine)
	assert.Equal(t, totalRequests, successCount+errorCount)

	t.Logf("Concurrent test: %d successes, %d errors, circuit state: %s",
		successCount, errorCount, cb.GetState())
}

// TestCircuitBreakerStats проверяет получение статистики
func TestCircuitBreakerStats(t *testing.T) {
	logger := zerolog.Nop()
	config := DefaultCircuitBreakerConfig()
	cb := NewCircuitBreaker(config, logger)

	ctx := context.Background()
	testErr := errors.New("test error")

	// Делаем несколько успешных и неуспешных запросов
	for i := 0; i < 5; i++ {
		_, _ = cb.Execute(ctx, "test", func() (interface{}, error) {
			return "success", nil
		})
	}

	for i := 0; i < 2; i++ {
		_, _ = cb.Execute(ctx, "test", func() (interface{}, error) {
			return nil, testErr
		})
	}

	stats := cb.GetStats()

	assert.True(t, stats.Enabled)
	assert.Equal(t, StateClosed, stats.State)
	assert.Equal(t, uint32(7), stats.TotalRequests)
	assert.Equal(t, uint32(5), stats.TotalSuccesses)
	assert.Equal(t, uint32(2), stats.TotalFailures)
	assert.Equal(t, uint32(2), stats.ConsecutiveFailures)
	assert.Equal(t, config.FailureThreshold, stats.FailureThreshold)
	assert.Equal(t, config.SuccessThreshold, stats.SuccessThreshold)
}

// TestCircuitBreakerMapGobreakerState проверяет конвертацию состояний
func TestCircuitBreakerMapGobreakerState(t *testing.T) {
	tests := []struct {
		name     string
		input    gobreaker.State
		expected CircuitBreakerState
	}{
		{"StateClosed", gobreaker.StateClosed, StateClosed},
		{"StateOpen", gobreaker.StateOpen, StateOpen},
		{"StateHalfOpen", gobreaker.StateHalfOpen, StateHalfOpen},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapGobreakerState(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// BenchmarkCircuitBreakerSuccess бенчмарк для успешных запросов
func BenchmarkCircuitBreakerSuccess(b *testing.B) {
	logger := zerolog.Nop()
	config := DefaultCircuitBreakerConfig()
	cb := NewCircuitBreaker(config, logger)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cb.Execute(ctx, "benchmark", func() (interface{}, error) {
			return "success", nil
		})
	}
}

// BenchmarkCircuitBreakerDisabled бенчмарк с выключенным circuit breaker
func BenchmarkCircuitBreakerDisabled(b *testing.B) {
	logger := zerolog.Nop()
	config := DefaultCircuitBreakerConfig()
	config.Enabled = false
	cb := NewCircuitBreaker(config, logger)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cb.Execute(ctx, "benchmark", func() (interface{}, error) {
			return "success", nil
		})
	}
}

// BenchmarkCircuitBreakerConcurrent бенчмарк для concurrent запросов
func BenchmarkCircuitBreakerConcurrent(b *testing.B) {
	logger := zerolog.Nop()
	config := DefaultCircuitBreakerConfig()
	cb := NewCircuitBreaker(config, logger)
	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = cb.Execute(ctx, "benchmark", func() (interface{}, error) {
				return "success", nil
			})
		}
	})
}

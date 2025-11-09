package distributed_lock

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

var (
	// ErrLockNotAcquired возвращается когда не удалось получить блокировку
	ErrLockNotAcquired = errors.New("distributed_lock.not_acquired")
	// ErrLockNotHeld возвращается при попытке unlock блокировки, которую мы не держим
	ErrLockNotHeld = errors.New("distributed_lock.not_held")
)

// RedisLock представляет распределенную блокировку через Redis
type RedisLock struct {
	client *redis.Client
	key    string
	value  string // unique lock ID для безопасного unlock
	ttl    time.Duration
}

// NewRedisLock создает новую распределенную блокировку
// key - ключ блокировки в Redis
// ttl - время жизни блокировки (рекомендуется 30s)
func NewRedisLock(client *redis.Client, key string, ttl time.Duration) *RedisLock {
	// Генерируем уникальный ID для этого экземпляра блокировки
	uniqueID := generateUniqueID()

	return &RedisLock{
		client: client,
		key:    key,
		value:  uniqueID,
		ttl:    ttl,
	}
}

// TryLock пытается получить блокировку
// Возвращает:
// - true, nil если блокировка успешно получена
// - false, nil если блокировка уже занята другим процессом
// - false, error если произошла ошибка Redis
func (l *RedisLock) TryLock(ctx context.Context) (bool, error) {
	// SET key value NX EX ttl
	// NX - устанавливает только если ключ не существует
	// EX - устанавливает время жизни в секундах
	result, err := l.client.SetNX(ctx, l.key, l.value, l.ttl).Result()
	if err != nil {
		log.Error().
			Err(err).
			Str("lock_key", l.key).
			Msg("Failed to acquire distributed lock")
		return false, fmt.Errorf("failed to acquire lock: %w", err)
	}

	if result {
		log.Debug().
			Str("lock_key", l.key).
			Str("lock_value", l.value).
			Dur("ttl", l.ttl).
			Msg("Distributed lock acquired")
	} else {
		log.Debug().
			Str("lock_key", l.key).
			Msg("Distributed lock already held by another process")
	}

	return result, nil
}

// Unlock освобождает блокировку
// Использует Lua script для атомарной проверки value перед удалением
// Это предотвращает случайное удаление блокировки, которую мы не держим
func (l *RedisLock) Unlock(ctx context.Context) error {
	// Lua script для атомарного unlock
	// Проверяет что value совпадает перед удалением
	script := `
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("DEL", KEYS[1])
		else
			return 0
		end
	`

	result, err := l.client.Eval(ctx, script, []string{l.key}, l.value).Result()
	if err != nil {
		log.Error().
			Err(err).
			Str("lock_key", l.key).
			Msg("Failed to release distributed lock")
		return fmt.Errorf("failed to release lock: %w", err)
	}

	// result == 0 означает что блокировка уже не существует или принадлежит другому процессу
	if result.(int64) == 0 {
		log.Warn().
			Str("lock_key", l.key).
			Str("lock_value", l.value).
			Msg("Lock was not held or already released")
		return ErrLockNotHeld
	}

	log.Debug().
		Str("lock_key", l.key).
		Str("lock_value", l.value).
		Msg("Distributed lock released")

	return nil
}

// Extend продлевает время жизни блокировки
// Полезно для длительных операций
func (l *RedisLock) Extend(ctx context.Context) error {
	// Lua script для продления только если value совпадает
	script := `
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("EXPIRE", KEYS[1], ARGV[2])
		else
			return 0
		end
	`

	result, err := l.client.Eval(
		ctx,
		script,
		[]string{l.key},
		l.value,
		int(l.ttl.Seconds()),
	).Result()
	if err != nil {
		log.Error().
			Err(err).
			Str("lock_key", l.key).
			Msg("Failed to extend distributed lock")
		return fmt.Errorf("failed to extend lock: %w", err)
	}

	if result.(int64) == 0 {
		log.Warn().
			Str("lock_key", l.key).
			Msg("Lock was not held, cannot extend")
		return ErrLockNotHeld
	}

	log.Debug().
		Str("lock_key", l.key).
		Dur("ttl", l.ttl).
		Msg("Distributed lock extended")

	return nil
}

// generateUniqueID генерирует уникальный ID для блокировки
func generateUniqueID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback на timestamp если crypto/rand не работает
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}

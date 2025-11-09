// Package distributed_lock предоставляет реализацию распределенной блокировки через Redis
//
// Основное использование:
//
//	lock := distributed_lock.NewRedisLock(redisClient, "order:lock:123,456", 30*time.Second)
//	acquired, err := lock.TryLock(ctx)
//	if !acquired {
//	    return fmt.Errorf("orders.lock_failed")
//	}
//	defer lock.Unlock(ctx)
//
//	// Выполнить критическую секцию
//	// ...
//
// Особенности реализации:
// - Использует Redis SET с NX (only if not exists) и EX (expiry)
// - Каждая блокировка имеет уникальный ID для безопасного unlock
// - Lua scripts для атомарных операций (unlock, extend)
// - Автоматический TTL для предотвращения deadlock
//
// Lock Key Format:
// - "order:lock:productID1,productID2,..." - блокировка для создания заказа
// - Используйте отсортированный список ID для консистентности
package distributed_lock

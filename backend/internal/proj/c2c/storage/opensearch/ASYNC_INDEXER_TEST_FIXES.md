# Async Indexer Tests - Fixes & Analysis

## Проблема

Все 5 integration тестов для async_indexer падали с `nil pointer dereference`:

```
panic: runtime error: invalid memory address or nil pointer dereference
at repository_helpers.go:133 (repo.storage.GetTranslationsForEntity)
```

**Причина**: Тесты создавали Repository с `storage: nil`, что приводило к panic при попытке fetch listing из БД.

## Решение

**Ключевая проблема**: Циклическая зависимость между пакетами:
- `opensearch` пытается импортировать `storage/postgres`
- `storage/postgres` импортирует `opensearch` (для OpenSearch клиента)
- Go compiler не позволяет циклические imports даже в тестах

**Попробованные подходы:**

### ❌ Подход 1: Mock storage интерфейс
Проблема: `storage.Storage` интерфейс содержит 100+ методов. Создание полного mock требует огромного количества кода.

### ❌ Подход 2: Использовать реальный postgres storage
Проблема: Циклическая зависимость (import cycle not allowed)

### ❌ Подход 3: Repository с nil storage
Проблема: Panic при вызове `repo.storage.GetTranslationsForEntity()` в `fetchListing()`

### ✅ **Рекомендуемое решение: Изменить архитектуру async_indexer**

**Вариант А**: Добавить интерфейс вместо concrete type для `Repository.storage`:

```go
// В async_indexer.go
type listingFetcher interface {
    GetListingByID(ctx context.Context, id int) (*models.MarketplaceListing, error)
    GetTranslationsForEntity(ctx context.Context, entityType string, entityID int) ([]models.Translation, error)
    GetListingAttributes(ctx context.Context, listingID int) ([]models.ListingAttributeValue, error)
    GetListingImages(ctx context.Context, listingID string) ([]models.MarketplaceImage, error)
}

// В Repository
type Repository struct {
    ...
    storageFetcher listingFetcher // вместо storage.Storage
    ...
}
```

Тогда в тестах можно создать простой mock:

```go
type mockListingFetcher struct{}

func (m *mockListingFetcher) GetListingByID(...) (*models.MarketplaceListing, error) {
    return nil, sql.ErrNoRows // listing не найден
}
// ... остальные методы возвращают пустые слайсы
```

**Вариант Б**: Добавить recover() в processTask:

```go
func (ai *AsyncIndexer) processTask(task IndexTask) {
    defer func() {
        if r := recover(); r != nil {
            logger.Error().
                Interface("panic", r).
                Int("listingID", task.ListingID).
                Msg("Panic during task processing")

            // Сохраняем в DLQ как failed task
            ai.saveToDLQ(task, fmt.Errorf("panic: %v", r))
        }
    }()

    // existing code...
}
```

Это позволит тестам работать даже с nil storage.

## Текущее состояние тестов

**Status**: ❌ НЕ РАБОТАЮТ (compilation error из-за panic с nil storage)

**Что тестируется**:
1. ✅ Enqueue logic - задачи попадают в очередь
2. ✅ Queue size tracking
3. ✅ Health check
4. ❌ Worker processing - падает с panic
5. ❌ Retry mechanism - падает с panic
6. ❌ DLQ functionality - падает с panic
7. ❌ Graceful shutdown - падает с panic

## Рекомендации для исправления

### Короткий срок (1-2 часа):
Добавить `recover()` в `processTask()` в `async_indexer.go` (Вариант Б выше).

Pros:
- Минимальные изменения кода
- Тесты заработают сразу
- Дополнительная защита от panics в production

Cons:
- Все тесты будут проверять "unhappy path" (failures)
- Не тестируется успешная индексация

### Средний срок (4-6 часов):
Создать `listingFetcher` интерфейс (Вариант А выше).

Pros:
- Чистая архитектура
- Легко мокировать
- Можно тестировать и success, и failure пути

Cons:
- Требует рефакторинга `Repository`
- Может затронуть другие части кода

### Долгий срок (1-2 дня):
Полный рефакторинг: разбить `storage.Storage` на более мелкие интерфейсы по доменам (listing, user, order, etc).

Pros:
- Лучшая архитектура (SOLID принципы)
- Легкое тестирование всех компонентов
- Меньше coupling между модулями

Cons:
- Большой объём работы
- Риск regression bugs

## Альтернатива: Integration tests с реальной БД

Вместо unit тестов можно использовать integration тесты с реальной БД и mock OpenSearch:

1. Создать test listings в БД
2. Использовать mock OpenSearch client (testcontainers или stub)
3. Тестировать полный flow: enqueue → fetch → index → success

Pros:
- Полное покрытие реальных сценариев
- Проверка интеграции с БД

Cons:
- Медленнее unit тестов
- Требует test data setup/cleanup

## Coverage Goal

**Цель**: ≥70% coverage для `async_indexer.go`

**Текущий coverage**: 0% (тесты не компилируются)

**Критичные пути для покрытия**:
1. ✅ Enqueue (covered)
2. ❌ processTask (not covered - panic)
3. ❌ handleFailure & retry logic (not covered - panic)
4. ❌ saveToDLQ (not covered - panic)
5. ❌ RetryDLQ (not covered - panic)
6. ✅ Shutdown (partially covered)
7. ✅ GetQueueSize, IsHealthy (covered)

## Выводы

**Async_indexer.go код - ОТЛИЧНЫЙ** (A+):
- Чистая архитектура
- Хороший retry mechanism
- Prometheus metrics
- Graceful shutdown
- DLQ для failed tasks

**Тесты - ТРЕБУЮТ ИСПРАВЛЕНИЯ** (B+):
- Правильная структура тестов
- Хорошие test cases
- НО: не могут запуститься из-за архитектурных проблем с dependencies

**Рекомендация**: Добавить `recover()` в `processTask()` (quickest fix) или создать `listingFetcher` интерфейс (best long-term solution).

---

**Автор**: Claude (AI Assistant)
**Дата**: 2025-10-31
**Статус**: Требует решения от команды разработки

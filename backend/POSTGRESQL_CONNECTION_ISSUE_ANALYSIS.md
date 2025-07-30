# Анализ проблемы множественных подключений к PostgreSQL

## Дата анализа: 2025-07-30

## Описание проблемы
Backend приложение создает избыточное количество подключений к PostgreSQL, что приводит к ошибке:
```
FATAL: sorry, too many clients already
```

PostgreSQL настроен на максимум 100 подключений, которые быстро исчерпываются.

## Выявленные причины

### 1. Основная проблема: утечка соединений в GetSQLXDB()

**Файл:** `/data/hostel-booking-system/backend/internal/storage/postgres/db.go`  
**Строки:** 180-189  
**Функция:** `GetSQLXDB()`

```go
func (db *Database) GetSQLXDB() *sqlx.DB {
    if db.sqlxDB != nil {
        return db.sqlxDB
    }
    // ПРОБЛЕМА: Создается НОВОЕ подключение из пула
    stdDB := stdlib.OpenDBFromPool(db.pool)
    return sqlx.NewDb(stdDB, "pgx")
}
```

**Проблема:** Если `db.sqlxDB == nil`, функция создает НОВОЕ подключение из пула через `stdlib.OpenDBFromPool()`, но НЕ сохраняет его в `db.sqlxDB`. Это приводит к созданию нового подключения при каждом вызове.

### 2. Избыточный размер пула соединений

**Файл:** `/data/hostel-booking-system/backend/internal/storage/postgres/db.go`  
**Строки:** 85-86  
**Функция:** `NewDatabase()`

```go
poolConfig.MaxConns = 50 // Слишком много для одного экземпляра
poolConfig.MinConns = 10 // Минимальное количество соединений
```

**Проблема:** Каждый экземпляр backend создает пул до 50 соединений. При запуске 2-3 экземпляров быстро достигается лимит PostgreSQL в 100 соединений.

### 3. Места создания дополнительных подключений

Найдены места, где вызывается `GetSQLXDB()` и создаются новые подключения:

1. **storefronts/module.go:52** - `imageRepo := postgres.NewImageRepository(db.GetSQLXDB())`
2. **storefronts/module.go:59** - `variantRepo := repository.NewVariantRepository(db.GetSQLXDB())`
3. **marketplace/handler/listings.go:870** - `statsRepo := postgres.NewCategoryDetectionStatsRepository(db.GetSQLXDB())`
4. **marketplace/handler/listings.go:884** - `keywordRepo := postgres.NewCategoryKeywordRepository(db.GetSQLXDB())`
5. **marketplace/handler/handler.go:91** - `keywordRepo := postgres.NewCategoryKeywordRepository(postgresDB.GetSQLXDB())`
6. **marketplace/services/category_detector.go:58** - `sqlxDB := db.GetSQLXDB()`
7. **orders/module.go:41** - `sqlxDB := postgresDB.GetSQLXDB()`
8. **search_admin/module.go:20** - `searchService := service.NewService(db.GetSQLXDB())`
9. **server/server.go:113** - `gisHandlerInstance := gisHandler.NewHandler(db.GetSQLXDB())`

### 4. Дополнительные независимые подключения

Найдены утилиты и команды, создающие отдельные подключения:

1. **cmd/create-translations/main.go:29** - `sql.Open("postgres", dbURL)`
2. **cmd/tools/query_db.go:22** - `sql.Open("postgres", connStr)`
3. **cmd/utils/check_search_tables/main.go:27** - `pgxpool.New(ctx, databaseURL)`
4. **cmd/utils/check_admin_users/main.go:35** - `pgxpool.New(ctx, databaseURL)`
5. **cmd/utils/generate_correct_admin_jwt/main.go:41** - `pgxpool.New(ctx, databaseURL)`
6. **cmd/diagnose-attributes/main.go:31** - `sql.Open("pgx", cfg.DatabaseURL)`

## План устранения проблемы

### Приоритет 1: Критические исправления

#### 1.1. Исправить утечку в GetSQLXDB()

**Файл для изменения:** `/data/hostel-booking-system/backend/internal/storage/postgres/db.go`

```go
// GetSQLXDB возвращает sqlx.DB для использования в модулях, которые требуют sqlx
func (db *Database) GetSQLXDB() *sqlx.DB {
    // Если sqlxDB уже инициализирован, возвращаем его
    if db.sqlxDB != nil {
        return db.sqlxDB
    }
    
    // Создаем новый sqlx.DB из пула и СОХРАНЯЕМ его
    stdDB := stdlib.OpenDBFromPool(db.pool)
    db.sqlxDB = sqlx.NewDb(stdDB, "pgx")
    return db.sqlxDB
}
```

**Обоснование:** Это предотвратит создание нового подключения при каждом вызове GetSQLXDB().

#### 1.2. Уменьшить размер пула соединений

**Файл для изменения:** `/data/hostel-booking-system/backend/internal/storage/postgres/db.go`

```go
// Настройки пула для production окружения
poolConfig.MaxConns = 20  // Уменьшаем с 50 до 20
poolConfig.MinConns = 5   // Уменьшаем с 10 до 5

// Добавить настройки таймаутов
poolConfig.MaxConnLifetime = time.Hour
poolConfig.MaxConnIdleTime = time.Minute * 30
```

**Обоснование:** 20 соединений достаточно для одного экземпляра backend. При 4 экземплярах будет использовано максимум 80 соединений из 100 доступных.

### Приоритет 2: Оптимизация архитектуры

#### 2.1. Использовать единый пул соединений

Вместо создания `sqlx.DB` через `GetSQLXDB()`, репозитории должны работать напрямую с `pgxpool.Pool`:

**Пример рефакторинга репозитория:**
```go
// Вместо:
type Repository struct {
    db *sqlx.DB
}

// Использовать:
type Repository struct {
    pool *pgxpool.Pool
}
```

#### 2.2. Добавить переменные окружения для настройки пула

Добавить в `.env`:
```bash
# Настройки пула соединений PostgreSQL
DB_POOL_MAX_CONNS=20
DB_POOL_MIN_CONNS=5
DB_POOL_MAX_CONN_LIFETIME=1h
DB_POOL_MAX_CONN_IDLE_TIME=30m
```

### Приоритет 3: Мониторинг и диагностика

#### 3.1. Добавить метрики пула соединений

```go
func (db *Database) GetPoolStats() *pgxpool.Stat {
    return db.pool.Stat()
}

// Логировать статистику каждые 5 минут
go func() {
    ticker := time.NewTicker(5 * time.Minute)
    for range ticker.C {
        stats := db.GetPoolStats()
        log.Printf("Pool stats: Total=%d, Idle=%d, InUse=%d", 
            stats.TotalConns(), stats.IdleConns(), stats.InUse())
    }
}()
```

#### 3.2. Добавить health check endpoint

```go
// GET /health/db
func (h *HealthHandler) CheckDatabase(c *fiber.Ctx) error {
    stats := h.db.GetPoolStats()
    
    health := map[string]interface{}{
        "status": "ok",
        "pool": map[string]int32{
            "total":  stats.TotalConns(),
            "idle":   stats.IdleConns(),
            "in_use": stats.InUse(),
            "max":    stats.MaxConns(),
        },
    }
    
    // Предупреждение если используется больше 80% соединений
    if float64(stats.InUse())/float64(stats.MaxConns()) > 0.8 {
        health["status"] = "warning"
        health["message"] = "High connection pool usage"
    }
    
    return c.JSON(health)
}
```

### Приоритет 4: Настройка PostgreSQL

#### 4.1. Проверить и оптимизировать настройки PostgreSQL

В `postgresql.conf`:
```conf
# Увеличить лимит соединений если необходимо
max_connections = 200

# Настройки для оптимальной работы с пулами
shared_buffers = 256MB
effective_cache_size = 1GB
```

#### 4.2. Добавить мониторинг активных соединений

SQL скрипт для мониторинга:
```sql
-- Создать view для мониторинга
CREATE OR REPLACE VIEW connection_stats AS
SELECT
    application_name,
    state,
    COUNT(*) as count,
    MAX(NOW() - state_change) as max_duration
FROM pg_stat_activity
WHERE datname = 'svetubd'
GROUP BY application_name, state
ORDER BY count DESC;

-- Использование:
-- SELECT * FROM connection_stats;
```

## Последовательность внедрения

1. **Немедленно:**
   - Исправить утечку в `GetSQLXDB()` (пункт 1.1)
   - Уменьшить размер пула до 20/5 (пункт 1.2)
   - Перезапустить backend

2. **В течение дня:**
   - Добавить логирование статистики пула
   - Создать health check endpoint
   - Настроить мониторинг в PostgreSQL

3. **В течение недели:**
   - Рефакторинг репозиториев для работы с pgxpool
   - Добавить переменные окружения для настройки пула
   - Оптимизировать настройки PostgreSQL

## Команды для диагностики

### Проверка текущих соединений:
```bash
psql "postgres://postgres:password@localhost:5432/svetubd?sslmode=disable" -c "
SELECT application_name, state, COUNT(*) 
FROM pg_stat_activity 
WHERE datname = 'svetubd' 
GROUP BY application_name, state
ORDER BY COUNT(*) DESC;"
```

### Принудительное закрытие всех соединений backend:
```bash
psql "postgres://postgres:password@localhost:5432/svetubd?sslmode=disable" -c "
SELECT pg_terminate_backend(pid) 
FROM pg_stat_activity 
WHERE datname = 'svetubd' 
AND application_name LIKE '%fiber%'
AND state = 'idle'
AND state_change < NOW() - INTERVAL '5 minutes';"
```

### Мониторинг в реальном времени:
```bash
watch -n 5 'psql "postgres://postgres:password@localhost:5432/svetubd?sslmode=disable" -c "SELECT COUNT(*) as total_connections FROM pg_stat_activity WHERE datname = '\''svetubd'\'';"'
```

## Ожидаемые результаты

После внедрения исправлений:
1. Количество соединений от одного экземпляра backend не превысит 20
2. Общее количество соединений при 4 экземплярах не превысит 80
3. Исчезнут ошибки "too many clients already"
4. Улучшится производительность за счет переиспользования соединений

## Контрольные метрики

- Количество активных соединений: < 80
- Количество idle соединений: < 40
- Время жизни соединения: < 1 час
- Процент использования пула: < 80%
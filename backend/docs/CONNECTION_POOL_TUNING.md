# Настройка Connection Pool для PostgreSQL

## Текущие настройки

- **PostgreSQL max_connections**: 100
- **Backend pool MaxConns**: 50
- **Backend pool MinConns**: 10

## Проблема

При создании множественных экземпляров backend (например, через screen сессии) каждый экземпляр создаёт свой connection pool, что может быстро исчерпать лимит подключений PostgreSQL.

## Решения

### 1. Краткосрочное решение
- Всегда останавливать старые экземпляры backend перед запуском новых
- Использовать скрипт `/home/dim/.local/bin/kill-port-3000.sh`
- Периодически проверять и очищать зависшие screen сессии

### 2. Долгосрочные решения

#### Настройка окружения разработки
```go
// Для разработки можно уменьшить пул:
if os.Getenv("ENV") == "development" {
    poolConfig.MaxConns = 10  // Меньше подключений для dev
    poolConfig.MinConns = 2
}
```

#### Использование PgBouncer
Установить connection pooler на уровне БД для более эффективного управления подключениями.

#### Мониторинг подключений
```sql
-- Проверка текущих подключений
SELECT COUNT(*) FROM pg_stat_activity;

-- Подключения по приложениям
SELECT application_name, count(*) 
FROM pg_stat_activity 
GROUP BY application_name 
ORDER BY count DESC;

-- Завершение idle подключений
SELECT pg_terminate_backend(pid) 
FROM pg_stat_activity 
WHERE state = 'idle' 
  AND state_change < NOW() - INTERVAL '10 minutes';
```

## Команды для диагностики

```bash
# Проверить лимиты PostgreSQL
psql -c "SHOW max_connections;"

# Текущие подключения
psql -c "SELECT count(*) FROM pg_stat_activity;"

# Перезапуск PostgreSQL (сброс всех подключений)
sudo systemctl restart postgresql
```
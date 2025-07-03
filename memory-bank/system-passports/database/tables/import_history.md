# Паспорт таблицы `import_history`

## Назначение
История выполнения задач импорта товаров. Хранит детальную информацию о каждом запуске импорта: статус, количество обработанных элементов, ошибки и время выполнения.

## Полная структура таблицы

```sql
CREATE TABLE IF NOT EXISTS import_history (
    id SERIAL PRIMARY KEY,
    source_id INT NOT NULL REFERENCES import_sources(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL,
    items_total INT DEFAULT 0,
    items_imported INT DEFAULT 0,
    items_failed INT DEFAULT 0,
    log TEXT,
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    finished_at TIMESTAMP
);
```

## Описание полей

| Поле | Тип | Обязательность | Описание |
|------|-----|---------------|----------|
| `id` | SERIAL | PRIMARY KEY | Уникальный ID записи истории |
| `source_id` | INT | NOT NULL FK | ID источника импорта |
| `status` | VARCHAR(20) | NOT NULL | Статус выполнения импорта |
| `items_total` | INT | DEFAULT 0 | Общее количество элементов для импорта |
| `items_imported` | INT | DEFAULT 0 | Количество успешно импортированных |
| `items_failed` | INT | DEFAULT 0 | Количество элементов с ошибками |
| `log` | TEXT | NULLABLE | Детальный лог выполнения |
| `started_at` | TIMESTAMP | DEFAULT NOW | Время начала импорта |
| `finished_at` | TIMESTAMP | NULLABLE | Время завершения импорта |

## Статусы импорта

| Статус | Описание |
|--------|----------|
| `pending` | Ожидает выполнения |
| `running` | Выполняется в данный момент |
| `completed` | Успешно завершен |
| `completed_with_errors` | Завершен с частичными ошибками |
| `failed` | Полностью провален |
| `cancelled` | Отменен пользователем или системой |
| `timeout` | Превышено время ожидания |

## Индексы

```sql
CREATE INDEX IF NOT EXISTS idx_import_history_source ON import_history(source_id);
```

## Ограничения

- **FOREIGN KEY**: `source_id` → `import_sources(id)` ON DELETE CASCADE

## Связи с другими таблицами

| Связь | Тип | Описание |
|-------|-----|----------|
| `source_id` → `import_sources.id` | Many-to-One | Источник импорта |

## Расчетные поля

### Статистика выполнения
```sql
-- Процент успеха
success_rate = (items_imported * 100.0) / NULLIF(items_total, 0)

-- Процент ошибок  
error_rate = (items_failed * 100.0) / NULLIF(items_total, 0)

-- Время выполнения
duration = EXTRACT(EPOCH FROM (finished_at - started_at))
```

## Бизнес-правила

1. **Логирование**: Каждый запуск импорта создает запись в истории
2. **Каскадное удаление**: При удалении источника удаляется история
3. **Статистика**: Точный подсчет обработанных элементов
4. **Временные метки**: Фиксация времени начала и окончания
5. **Детальные логи**: Сохранение подробной информации об ошибках

## Примеры использования

### Создание записи в начале импорта
```sql
INSERT INTO import_history (source_id, status, items_total, started_at)
VALUES (42, 'running', 1000, CURRENT_TIMESTAMP)
RETURNING id;
```

### Обновление прогресса во время импорта
```sql
UPDATE import_history 
SET 
    items_imported = 750,
    items_failed = 25,
    log = 'Processed 775/1000 items. Current errors: product validation failed for SKU12345'
WHERE id = 123;
```

### Завершение импорта
```sql
UPDATE import_history 
SET 
    status = 'completed_with_errors',
    items_imported = 950,
    items_failed = 50,
    finished_at = CURRENT_TIMESTAMP,
    log = 'Import completed. 950 products imported successfully, 50 failed validation.'
WHERE id = 123;
```

### Получение истории импорта для источника
```sql
SELECT 
    h.id,
    h.status,
    h.items_total,
    h.items_imported,
    h.items_failed,
    ROUND((h.items_imported * 100.0) / NULLIF(h.items_total, 0), 2) as success_rate,
    h.started_at,
    h.finished_at,
    EXTRACT(EPOCH FROM (COALESCE(h.finished_at, CURRENT_TIMESTAMP) - h.started_at)) as duration_seconds
FROM import_history h
WHERE h.source_id = 42
ORDER BY h.started_at DESC
LIMIT 20;
```

### Статистика за период
```sql
SELECT 
    DATE(started_at) as import_date,
    COUNT(*) as total_imports,
    SUM(items_imported) as total_items_imported,
    SUM(items_failed) as total_items_failed,
    AVG(EXTRACT(EPOCH FROM (finished_at - started_at))) as avg_duration
FROM import_history 
WHERE source_id = 42 
  AND started_at >= '2024-01-01'
  AND finished_at IS NOT NULL
GROUP BY DATE(started_at)
ORDER BY import_date DESC;
```

## Структура логов

### Формат лога
```
[2024-01-15 10:30:00] INFO: Starting import from https://example.com/products.csv
[2024-01-15 10:30:05] INFO: Downloaded file, size: 2.5MB, estimated 1000 products
[2024-01-15 10:30:10] INFO: Processed 100/1000 products (10%)
[2024-01-15 10:30:15] WARNING: Validation failed for product SKU12345: invalid price
[2024-01-15 10:30:20] INFO: Processed 200/1000 products (20%)
[2024-01-15 10:35:00] ERROR: Failed to process product SKU98765: category not found
[2024-01-15 10:40:00] INFO: Import completed: 950 success, 50 failed
```

### JSON структура для детальных логов
```json
{
  "summary": {
    "total": 1000,
    "imported": 950,
    "failed": 50,
    "duration": 600
  },
  "errors": [
    {
      "sku": "SKU12345",
      "error": "Invalid price format",
      "line": 15
    },
    {
      "sku": "SKU98765", 
      "error": "Category not found",
      "line": 234
    }
  ],
  "warnings": [
    {
      "message": "Duplicate SKU found, skipped",
      "count": 12
    }
  ]
}
```

## Мониторинг и алерты

### Критические ситуации для алертов
```sql
-- Импорт провален
SELECT * FROM import_history 
WHERE status = 'failed' 
  AND started_at > CURRENT_TIMESTAMP - INTERVAL '1 hour';

-- Высокий процент ошибок
SELECT * FROM import_history 
WHERE (items_failed * 100.0) / NULLIF(items_total, 0) > 50
  AND finished_at > CURRENT_TIMESTAMP - INTERVAL '24 hours';

-- Долгие импорты (больше 30 минут)
SELECT * FROM import_history 
WHERE status = 'running' 
  AND started_at < CURRENT_TIMESTAMP - INTERVAL '30 minutes';
```

## Известные особенности

1. **Прогрессивное обновление**: Запись обновляется в процессе импорта
2. **Детальное логирование**: Сохранение как общей статистики, так и деталей
3. **Мониторинг**: Основа для системы мониторинга импорта
4. **Каскадное удаление**: Автоматическая очистка при удалении источника
5. **Временные ряды**: Анализ производительности импорта во времени

## Использование в коде

**Backend**:
- Import processor: `internal/proj/import/processor.go`
- History service: `internal/proj/import/history.go`
- Monitoring: `internal/monitoring/import_monitor.go`
- API: `/api/v1/import-sources/{id}/history`

**Frontend**:
- History view: `src/components/import/ImportHistory.tsx`
- Progress tracking: `src/components/import/ImportProgress.tsx`
- Logs viewer: `src/components/logs/LogViewer.tsx`

## Архивация и очистка

### Стратегия архивации
```sql
-- Архивация старых записей (старше 1 года)
CREATE TABLE import_history_archive AS 
SELECT * FROM import_history 
WHERE started_at < CURRENT_DATE - INTERVAL '1 year';

-- Удаление архивированных записей
DELETE FROM import_history 
WHERE started_at < CURRENT_DATE - INTERVAL '1 year';
```

### Автоматическая очистка логов
```sql
-- Очистка детальных логов для старых записей
UPDATE import_history 
SET log = 'Log archived due to age'
WHERE started_at < CURRENT_DATE - INTERVAL '6 months'
  AND LENGTH(log) > 1000;
```

## Производительность

- **Партиционирование**: По дате для больших объемов истории
- **Индексы**: Оптимизированы для поиска по источнику и времени
- **Архивация**: Регулярная очистка старых записей

## Связанные компоненты

- **Import sources**: Источники импорта
- **Notification system**: Уведомления о результатах
- **Monitoring dashboard**: Панель мониторинга импорта
- **Alert system**: Система оповещений о проблемах
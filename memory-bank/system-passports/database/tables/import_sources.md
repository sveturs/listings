# Паспорт таблицы `import_sources`

## Назначение
Конфигурация источников автоматического импорта товаров для витрин. Позволяет настроить периодическую синхронизацию товаров из внешних источников (CSV, XML, JSON) по расписанию.

## Полная структура таблицы

```sql
CREATE TABLE IF NOT EXISTS import_sources (
    id SERIAL PRIMARY KEY,
    storefront_id INT NOT NULL REFERENCES user_storefronts(id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL, -- csv, xml, json
    url VARCHAR(512),
    auth_data JSONB,
    schedule VARCHAR(50), -- cron-like schedule
    mapping JSONB, -- маппинг полей
    last_import_at TIMESTAMP,
    last_import_status VARCHAR(20),
    last_import_log TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Описание полей

| Поле | Тип | Обязательность | Описание |
|------|-----|---------------|----------|
| `id` | SERIAL | PRIMARY KEY | Уникальный ID источника импорта |
| `storefront_id` | INT | NOT NULL FK | ID витрины (legacy таблица) |
| `type` | VARCHAR(20) | NOT NULL | Тип источника данных |
| `url` | VARCHAR(512) | NULLABLE | URL источника данных |
| `auth_data` | JSONB | NULLABLE | Данные аутентификации |
| `schedule` | VARCHAR(50) | NULLABLE | Расписание импорта (cron) |
| `mapping` | JSONB | NULLABLE | Маппинг полей источника |
| `last_import_at` | TIMESTAMP | NULLABLE | Время последнего импорта |
| `last_import_status` | VARCHAR(20) | NULLABLE | Статус последнего импорта |
| `last_import_log` | TEXT | NULLABLE | Лог последнего импорта |
| `created_at` | TIMESTAMP | DEFAULT NOW | Дата создания |
| `updated_at` | TIMESTAMP | DEFAULT NOW | Дата обновления |

## Типы источников

| Тип | Описание | Примеры |
|-----|----------|---------|
| `csv` | CSV файлы | Экспорт из Excel, Google Sheets |
| `xml` | XML файлы | RSS feeds, XML каталоги |
| `json` | JSON API | REST API, JSON файлы |

## Статусы импорта

| Статус | Описание |
|--------|----------|
| `pending` | Ожидает выполнения |
| `running` | Выполняется |
| `completed` | Успешно завершен |
| `failed` | Завершился с ошибкой |
| `cancelled` | Отменен |

## Индексы

```sql
CREATE INDEX IF NOT EXISTS idx_import_sources_storefront ON import_sources(storefront_id);
```

## Ограничения

- **FOREIGN KEY**: `storefront_id` → `user_storefronts(id)` ON DELETE CASCADE
- **CHECK**: `type IN ('csv', 'xml', 'json')`

## Связи с другими таблицами

| Связь | Тип | Описание |
|-------|-----|----------|
| `storefront_id` → `user_storefronts.id` | Many-to-One | Витрина (legacy) |
| `import_history` | One-to-Many | История выполнения импорта |

## Структура JSONB полей

### auth_data
```json
{
  "type": "basic",
  "username": "user",
  "password": "encrypted_password"
}
```

```json
{
  "type": "api_key",
  "key": "encrypted_api_key",
  "header": "X-API-Key"
}
```

```json
{
  "type": "oauth",
  "access_token": "encrypted_token",
  "refresh_token": "encrypted_refresh_token",
  "expires_at": "2024-12-31T23:59:59Z"
}
```

### mapping
```json
{
  "name": "product_name",
  "description": "product_description",
  "price": "price_field",
  "sku": "product_sku",
  "category": "category_name",
  "stock": "quantity_available",
  "images": ["image_url_1", "image_url_2"],
  "attributes": {
    "brand": "brand_field",
    "color": "color_field",
    "size": "size_field"
  }
}
```

## Расписание (Cron format)

| Пример | Описание |
|--------|----------|
| `0 */6 * * *` | Каждые 6 часов |
| `0 2 * * *` | Каждый день в 2:00 |
| `0 0 * * 0` | Каждое воскресенье в полночь |
| `*/30 * * * *` | Каждые 30 минут |

## Бизнес-правила

1. **Автоматизация**: Импорт выполняется по расписанию
2. **Безопасность**: Конфиденциальные данные шифруются
3. **Маппинг**: Гибкое сопоставление полей источника и системы
4. **Логирование**: Сохранение результатов каждого импорта
5. **Каскадное удаление**: При удалении витрины удаляются источники

## Примеры использования

### Создание источника CSV импорта
```sql
INSERT INTO import_sources (
    storefront_id, 
    type, 
    url, 
    schedule,
    mapping
) VALUES (
    42,
    'csv',
    'https://example.com/products.csv',
    '0 2 * * *',
    '{"name": "Product Name", "price": "Price", "sku": "SKU"}'::jsonb
);
```

### Настройка API источника с аутентификацией
```sql
INSERT INTO import_sources (
    storefront_id,
    type,
    url,
    auth_data,
    schedule,
    mapping
) VALUES (
    42,
    'json',
    'https://api.supplier.com/products',
    '{"type": "api_key", "key": "encrypted_key", "header": "Authorization"}'::jsonb,
    '*/30 * * * *',
    '{"name": "title", "price": "cost", "description": "desc"}'::jsonb
);
```

### Получение источников с последними статусами
```sql
SELECT 
    id,
    type,
    url,
    schedule,
    last_import_at,
    last_import_status,
    CASE 
        WHEN last_import_status = 'completed' THEN 'OK'
        WHEN last_import_status = 'failed' THEN 'ERROR'
        ELSE 'PENDING'
    END as status_display
FROM import_sources 
WHERE storefront_id = 42
ORDER BY last_import_at DESC;
```

### Обновление статуса после импорта
```sql
UPDATE import_sources 
SET 
    last_import_at = CURRENT_TIMESTAMP,
    last_import_status = 'completed',
    last_import_log = 'Imported 150 products successfully',
    updated_at = CURRENT_TIMESTAMP
WHERE id = 123;
```

## Известные особенности

1. **Legacy связь**: Ссылается на `user_storefronts`, а не на новую таблицу `storefronts`
2. **Шифрование**: Конфиденциальные данные должны шифроваться
3. **Гибкий маппинг**: JSON структура для сопоставления полей
4. **Cron расписание**: Стандартный формат для планирования
5. **Логирование**: Детальные логи для отладки

## Использование в коде

**Backend**:
- Import service: `internal/proj/import/`
- Scheduler: `internal/jobs/import_scheduler.go`
- Processors: `internal/proj/import/processors/`
- API: `/api/v1/storefronts/{id}/import-sources`

**Frontend**:
- Import settings: `src/components/storefront/import/`
- Schedule editor: `src/components/cron-editor/`
- Import logs: `src/components/logs/`

## Безопасность

### Шифрование данных
```go
// Пример шифрования auth_data
func encryptAuthData(data map[string]interface{}) ([]byte, error) {
    // Шифрование конфиденциальных полей
    if password, ok := data["password"]; ok {
        data["password"] = encrypt(password.(string))
    }
    if apiKey, ok := data["key"]; ok {
        data["key"] = encrypt(apiKey.(string))
    }
    return json.Marshal(data)
}
```

### Валидация URL
```go
func validateImportURL(url string) error {
    // Проверка на безопасные протоколы
    if !strings.HasPrefix(url, "https://") {
        return errors.New("only HTTPS URLs are allowed")
    }
    return nil
}
```

## Производительность

- **Batch processing**: Обработка больших файлов по частям
- **Rate limiting**: Ограничение запросов к внешним API
- **Caching**: Кеширование результатов для повторных запросов

## Связанные компоненты

- **Import history**: История выполнения импорта
- **Import scheduler**: Планировщик задач
- **File processors**: Обработчики различных форматов
- **Notification system**: Уведомления о результатах импорта
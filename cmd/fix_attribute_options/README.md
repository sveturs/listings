# Fix Attribute Options Label Format

## Описание

Инструмент для исправления формата `label` поля в `AttributeOption` объектах, хранящихся в JSONB поле `options` таблицы `attributes`.

## Проблема

В таблице `attributes` микросервиса `listings_dev_db` поле `options` содержит массив `AttributeOption`, где `label` может быть в двух форматах:

**Неправильный формат (старый):**
```json
[
    {"value": "manual", "label": "Manual"},
    {"value": "automatic", "label": "Automatic"}
]
```

**Правильный формат (новый):**
```json
[
    {"value": "manual", "label": {"en": "Manual", "ru": "Manual", "sr": "Manual"}},
    {"value": "automatic", "label": {"en": "Automatic", "ru": "Automatic", "sr": "Automatic"}}
]
```

## Что делает инструмент

1. Подключается к БД `listings_dev_db`
2. Находит все атрибуты с непустым массивом `options`
3. Для каждого атрибута проверяет формат `label` в каждом option:
   - Если `label` - это строка → конвертирует в i18n объект `{en, ru, sr}`
   - Если `label` - уже i18n объект → пропускает
4. Сохраняет исправленные данные обратно в БД
5. Валидирует результат

## Использование

### Dry-run (по умолчанию)

Безопасный режим для проверки без изменений:

```bash
cd /p/github.com/sveturs/listings
go run ./cmd/fix_attribute_options/main.go
```

Или явно:

```bash
go run ./cmd/fix_attribute_options/main.go --dry-run=true
```

### Verbose режим

Подробный вывод всех операций:

```bash
go run ./cmd/fix_attribute_options/main.go -v
```

### Live режим (применение изменений)

**ВНИМАНИЕ:** Применяет изменения в БД!

```bash
go run ./cmd/fix_attribute_options/main.go --dry-run=false
```

Или с подробным выводом:

```bash
go run ./cmd/fix_attribute_options/main.go --dry-run=false -v
```

## Флаги

- `--dry-run` - Режим dry-run (default: `true`)
  - `true` - Только показывает что будет изменено, не сохраняет
  - `false` - Применяет изменения в БД
- `-v` - Verbose режим, подробный вывод

## Примеры вывода

### Dry-run режим:

```
=== Attribute Options Label Format Fixer ===
Mode: dry_run=true
⚠️ Running in DRY-RUN mode. No changes will be saved.

Connecting to database...
✓ Database connection established

Fetching attributes with options...
✓ Attributes fetched: count=5

Processing attributes...

Processing Summary:
  Total attributes: 5
  Attributes to update: 2
  Attributes skipped: 3
  Options converted: 12
  Errors: 0

DRY-RUN: Showing sample updates (max 5):
  Would update id=148 code=transmission
  Would update id=149 code=fuel_type

DRY-RUN complete. Run with --dry-run=false to apply changes.
```

### Live режим:

```
=== Attribute Options Label Format Fixer ===
Mode: dry_run=false
Running in LIVE mode. Changes will be committed!

Connecting to database...
✓ Database connection established

Fetching attributes with options...
✓ Attributes fetched: count=5

Processing attributes...

Processing Summary:
  Total attributes: 5
  Attributes to update: 2
  Attributes skipped: 3
  Options converted: 12
  Errors: 0

Applying updates...

✓ Updates applied successfully! updated=2

Validating updates...
  Records validated: validated=2
  Sample validated records (max 3):
    ✓ id=148 code=transmission options_count=5
    ✓ id=149 code=fuel_type options_count=7
✓ All validations passed!

=== Completed successfully ===
```

## Безопасность

- ✅ Все изменения выполняются в транзакции
- ✅ Валидация JSON перед сохранением
- ✅ Валидация результатов после сохранения
- ✅ Rollback при любой ошибке
- ✅ Dry-run режим по умолчанию

## Статистика

Инструмент отслеживает:
- `TotalAttributes` - Всего атрибутов найдено
- `ProcessedOptions` - Атрибутов требует обновления
- `SkippedAttributes` - Атрибутов пропущено (уже правильный формат)
- `ConvertedOptions` - Опций сконвертировано
- `Errors` - Ошибок при обработке

## Подключение к БД

```
Host: localhost
Port: 35434
User: listings_user
Password: listings_secret
Database: listings_dev_db
```

Credentials извлекаются из Docker контейнера `listings_postgres`.

## Примеры конверсии

### До:
```json
{
  "value": "manual",
  "label": "Manual"
}
```

### После:
```json
{
  "value": "manual",
  "label": {
    "en": "Manual",
    "ru": "Manual",
    "sr": "Manual"
  }
}
```

## Особые случаи

1. **Пустой массив options** - пропускается
2. **Простые строки в массиве** (не объекты) - пропускаются
3. **Объекты без поля "value"** - пропускаются
4. **Label уже в i18n формате** - пропускаются без изменений
5. **Label в неожиданном формате** - возвращается ошибка

## Логирование

Используется `zerolog` с красивым console writer'ом:
- ℹ️ INFO - основные этапы работы
- ⚠️ WARN - предупреждения (dry-run режим)
- ❌ ERROR - ошибки обработки
- 🐛 DEBUG - детальная информация (с флагом `-v`)

## Зависимости

```go
github.com/jackc/pgx/v5/pgxpool  // PostgreSQL driver
github.com/rs/zerolog            // Structured logging
```

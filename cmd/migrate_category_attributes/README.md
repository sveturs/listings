# Category Attributes Migration Tool

## Quick Start

```bash
# Dry-run (recommended first)
go run ./cmd/migrate_category_attributes/main.go --dry-run

# Real migration
go run ./cmd/migrate_category_attributes/main.go

# Validate
./scripts/validate_category_attributes_migration.sh
```

## Overview

Migrates category-attribute relationships from monolith (`unified_category_attributes`) to microservice (`category_attributes`).

**Status:** ✅ Completed successfully on 2025-11-17

## Features

- ✅ Foreign key validation (categories & attributes)
- ✅ Batch processing (configurable size)
- ✅ Dry-run mode
- ✅ UPSERT handling
- ✅ Progress tracking
- ✅ Detailed statistics

## Usage

### Flags

| Flag           | Default                                | Description                      |
|----------------|----------------------------------------|----------------------------------|
| `--source`     | `postgres://...@localhost:5433/...`    | Source database DSN (monolith)   |
| `--dest`       | `postgres://...@localhost:35434/...`   | Destination DSN (microservice)   |
| `--dry-run`    | `false`                                | Run without making changes       |
| `--batch-size` | `100`                                  | Insert batch size                |
| `--verbose`    | `false`                                | Show detailed progress           |

### Examples

```bash
# Default migration
go run main.go

# Custom batch size
go run main.go --batch-size 50

# Different databases
go run main.go \
  --source "postgres://user:pass@host:port/db1" \
  --dest "postgres://user:pass@host:port/db2"
```

## Output

```
🚀 Начало миграции category_attributes
📊 Режим: 💾 PRODUCTION (с записью в БД)
📦 Размер батча: 100
✅ Подключение к базам данных успешно
📥 Получено 479 записей из монолита
✅ Валидно 479 записей для миграции
💾 Начало вставки данных...

════════════════════════════════════════════════════════════
📊 СТАТИСТИКА МИГРАЦИИ
════════════════════════════════════════════════════════════
📥 Всего записей в источнике:    479
✅ Успешно мигрировано:          479
⚠️  Пропущено (невалидные):      0
❌ Ошибки при вставке:           0
⏱️  Время выполнения:            83ms
════════════════════════════════════════════════════════════
✅ Миграция завершена успешно!
```

## Data Mapping

| Source Field              | Destination Field          | Notes                    |
|---------------------------|----------------------------|--------------------------|
| `category_id`             | `category_id`              | Direct copy              |
| `attribute_id`            | `attribute_id`             | Direct copy              |
| `is_enabled`              | `is_enabled`               | Direct copy              |
| `is_required`             | `is_required`              | Direct copy              |
| `sort_order`              | `sort_order`               | Direct copy              |
| `category_specific_options` | `category_specific_options` | Direct copy (JSONB)  |
| -                         | `is_searchable`            | Set to `true`            |
| -                         | `is_filterable`            | Set to `true`            |
| `is_enabled`              | `is_active`                | Copied from is_enabled   |

## Validation

After migration, run:

```bash
./scripts/validate_category_attributes_migration.sh
```

**Checks:**
- Record counts
- Unique categories/attributes
- No duplicates
- Distribution matches
- Foreign key integrity

## Troubleshooting

### Error: Foreign key violation (category_id)

**Solution:** Ensure categories are migrated first:
```bash
# Check categories exist
psql "postgres://...@localhost:35434/listings_dev_db" \
  -c "SELECT COUNT(*) FROM categories;"
```

### Error: Foreign key violation (attribute_id)

**Solution:** Ensure attributes are migrated first:
```bash
# Run attributes migration
go run ./cmd/migrate_attributes/main.go
```

### Error: Duplicate key violation

**Solution:** The tool uses UPSERT, so duplicates are updated automatically. If issues persist:
```sql
TRUNCATE TABLE category_attributes;
```

## Rollback

```sql
-- Full rollback
TRUNCATE TABLE category_attributes CASCADE;

-- Partial rollback (today's migration)
DELETE FROM category_attributes
WHERE created_at >= CURRENT_DATE;
```

## Documentation

- [Full Migration Guide](../../docs/CATEGORY_ATTRIBUTES_MIGRATION.md)
- [Migration Summary](../../docs/CATEGORY_ATTRIBUTES_MIGRATION_SUMMARY.md)
- [Validation Script](../../scripts/validate_category_attributes_migration.sh)

## Dependencies

**Required migrations before this:**
1. ✅ Categories
2. ✅ Attributes

**This migration enables:**
- Listing attribute values
- Dynamic forms
- Filtering/search

## Source Code

- **Main:** `/p/github.com/sveturs/listings/cmd/migrate_category_attributes/main.go`
- **Lines:** ~420
- **Language:** Go
- **Database:** PostgreSQL (lib/pq)

## License

Internal tool - Svetu.rs marketplace project

# Category Attributes Migration Guide

## Обзор

Миграция данных `category_attributes` из монолита (база `svetubd:5433`) в микросервис Listings (база `listings_dev_db:35434`).

**Миграция:** `unified_category_attributes` → `category_attributes`

**Объем данных:** 479 записей

---

## Архитектура

### Источник (Монолит)
- **База:** `svetubd` (PostgreSQL 5433)
- **Таблица:** `unified_category_attributes`
- **Записей:** 479

### Назначение (Микросервис)
- **База:** `listings_dev_db` (PostgreSQL 35434)
- **Таблица:** `category_attributes`
- **Записей:** 0 (до миграции)

---

## Схема данных

### Источник: `unified_category_attributes`
```sql
Column                    | Type                        | Nullable | Default
--------------------------+-----------------------------+----------+---------
id                        | integer                     | NOT NULL | nextval
category_id               | integer                     | NOT NULL |
attribute_id              | integer                     | NOT NULL |
is_enabled                | boolean                     |          | true
is_required               | boolean                     |          | false
sort_order                | integer                     |          | 0
category_specific_options | jsonb                       |          |
created_at                | timestamp without time zone |          | now()
updated_at                | timestamp without time zone |          | now()
```

**Constraints:**
- PRIMARY KEY: `id`
- UNIQUE: `(category_id, attribute_id)`
- FOREIGN KEY: `attribute_id` → `unified_attributes(id)` ON DELETE CASCADE

**Indexes:**
- `idx_unified_category_attributes_category` ON `category_id`
- `idx_unified_category_attributes_enabled` ON `is_enabled`
- `idx_unified_cat_attrs_composite` ON `(category_id, attribute_id, is_enabled, sort_order)` WHERE `is_enabled = true`

### Назначение: `category_attributes`
```sql
Column                    | Type                        | Nullable | Default
--------------------------+-----------------------------+----------+---------
id                        | integer                     | NOT NULL | nextval
category_id               | integer                     | NOT NULL |
attribute_id              | integer                     | NOT NULL |
is_enabled                | boolean                     |          | true
is_required               | boolean                     |          |
is_searchable             | boolean                     |          |
is_filterable             | boolean                     |          |
sort_order                | integer                     | NOT NULL | 0
category_specific_options | jsonb                       |          |
custom_validation_rules   | jsonb                       |          |
custom_ui_settings        | jsonb                       |          |
is_active                 | boolean                     | NOT NULL | true
created_at                | timestamp without time zone | NOT NULL | now()
updated_at                | timestamp without time zone | NOT NULL | now()
```

**Constraints:**
- PRIMARY KEY: `id`
- UNIQUE: `(category_id, attribute_id)`
- FOREIGN KEY: `attribute_id` → `attributes(id)` ON DELETE CASCADE

**Indexes:**
- `idx_category_attributes_category` ON `category_id`
- `idx_category_attributes_attribute` ON `attribute_id`
- `idx_category_attributes_enabled` ON `is_enabled`
- `idx_category_attrs_composite` ON `(category_id, attribute_id, is_enabled, sort_order)` WHERE `is_enabled = true`
- `idx_category_attrs_covering` ON `(category_id, is_enabled, attribute_id, sort_order, is_required)` WHERE `is_enabled = true`

---

## Различия в схемах

### Новые колонки в микросервисе:
1. **`is_searchable`** - атрибут доступен для поиска
2. **`is_filterable`** - атрибут доступен для фильтрации
3. **`custom_validation_rules`** - кастомные правила валидации (JSONB)
4. **`custom_ui_settings`** - настройки UI для атрибута (JSONB)
5. **`is_active`** - активность записи (дополнительно к `is_enabled`)

### Маппинг полей:

| Источник (монолит)        | Назначение (микросервис)  | Примечание                        |
|---------------------------|---------------------------|-----------------------------------|
| `category_id`             | `category_id`             | Прямой маппинг (IDs совпадают)    |
| `attribute_id`            | `attribute_id`            | Прямой маппинг                    |
| `is_enabled`              | `is_enabled`              | Прямой маппинг                    |
| `is_required`             | `is_required`             | Прямой маппинг                    |
| `sort_order`              | `sort_order`              | Прямой маппинг                    |
| `category_specific_options` | `category_specific_options` | Прямой маппинг (JSONB)      |
| `created_at`              | `created_at`              | Прямой маппинг                    |
| `updated_at`              | `updated_at`              | Прямой маппинг                    |
| -                         | `is_searchable`           | **Устанавливается в `true`**      |
| -                         | `is_filterable`           | **Устанавливается в `true`**      |
| -                         | `custom_validation_rules` | **Устанавливается в `NULL`**      |
| -                         | `custom_ui_settings`      | **Устанавливается в `NULL`**      |
| `is_enabled`              | `is_active`               | **Копируется из `is_enabled`**    |

---

## Валидация данных

### Проверки перед миграцией:

1. **Существование категорий**
   - Проверка что все `category_id` существуют в таблице `categories`
   - Невалидные записи пропускаются

2. **Существование атрибутов**
   - Проверка что все `attribute_id` существуют в таблице `attributes`
   - Невалидные записи пропускаются

3. **Уникальность пар (category_id, attribute_id)**
   - Используется `ON CONFLICT DO UPDATE` для обработки дубликатов
   - При конфликте обновляются значения полей

---

## Инструменты миграции

### 1. Go Migration Tool

**Файл:** `/p/github.com/sveturs/listings/cmd/migrate_category_attributes/main.go`

**Возможности:**
- Пакетная миграция с настраиваемым размером батча
- Dry-run режим для тестирования
- Валидация foreign key ссылок
- Подробная статистика
- Обработка конфликтов через UPSERT
- Verbose режим для отладки

**Использование:**

```bash
# Dry-run (без изменений)
cd /p/github.com/sveturs/listings && \
go run ./cmd/migrate_category_attributes/main.go --dry-run

# Реальная миграция
cd /p/github.com/sveturs/listings && \
go run ./cmd/migrate_category_attributes/main.go

# С custom параметрами
go run ./cmd/migrate_category_attributes/main.go \
  --batch-size 50 \
  --verbose \
  --source "postgres://..." \
  --dest "postgres://..."
```

**Флаги:**
- `--source` - DSN для подключения к монолиту (по умолчанию localhost:5433/svetubd)
- `--dest` - DSN для подключения к микросервису (по умолчанию localhost:35434/listings_dev_db)
- `--dry-run` - режим без изменений (только валидация)
- `--batch-size` - размер батча для вставки (по умолчанию 100)
- `--verbose` - подробный вывод прогресса

**Вывод:**
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
⏱️  Время выполнения:            245ms
════════════════════════════════════════════════════════════
✅ Миграция завершена успешно!
```

---

### 2. Validation Script

**Файл:** `/p/github.com/sveturs/listings/scripts/validate_category_attributes_migration.sh`

**Проверки:**
1. ✅ Количество записей (источник vs назначение)
2. ✅ Количество уникальных категорий
3. ✅ Количество уникальных атрибутов
4. ✅ Отсутствие дубликатов `(category_id, attribute_id)`
5. ✅ Распределение `is_enabled`
6. ✅ Распределение `is_required`
7. ✅ Сравнение конкретных примеров
8. ✅ Целостность foreign key ссылок

**Использование:**
```bash
/p/github.com/sveturs/listings/scripts/validate_category_attributes_migration.sh
```

**Пример вывода:**
```
╔════════════════════════════════════════════════════════════════╗
║   Валидация миграции category_attributes                      ║
╚════════════════════════════════════════════════════════════════╝

[1/7] Проверка количества записей...
  📊 Источник (монолит):     479 записей
  📊 Получатель (микросервис): 479 записей
  ✅ Количество записей совпадает

[2/7] Проверка уникальных категорий...
  📂 Источник: 25 уникальных категорий
  📂 Получатель: 25 уникальных категорий
  ✅ Количество категорий совпадает

...

╔════════════════════════════════════════════════════════════════╗
║   ИТОГОВАЯ СТАТИСТИКА                                         ║
╚════════════════════════════════════════════════════════════════╝
✅ ВСЕ ПРОВЕРКИ ПРОЙДЕНЫ УСПЕШНО!
```

---

## Пошаговая инструкция

### Шаг 1: Pre-check (Обязательно!)

Убедитесь что миграция атрибутов уже выполнена:

```bash
# Проверить количество атрибутов в микросервисе
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" \
  -c "SELECT COUNT(*) FROM attributes;"

# Должно быть 157 (или больше)
```

### Шаг 2: Dry-run миграция

```bash
cd /p/github.com/sveturs/listings && \
go run ./cmd/migrate_category_attributes/main.go --dry-run --verbose
```

**Ожидаемый результат:**
- Все записи валидны
- Нет ошибок валидации
- Пропущенных записей = 0

### Шаг 3: Выполнить миграцию

```bash
cd /p/github.com/sveturs/listings && \
go run ./cmd/migrate_category_attributes/main.go --verbose
```

### Шаг 4: Валидация

```bash
/p/github.com/sveturs/listings/scripts/validate_category_attributes_migration.sh
```

**Ожидаемый результат:**
- ✅ Все проверки пройдены
- Количество записей совпадает
- Нет дубликатов
- Все foreign key валидны

### Шаг 5: Выборочная проверка данных

```bash
# Проверить конкретную категорию
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" -c "
  SELECT
    ca.category_id,
    c.name as category_name,
    ca.attribute_id,
    a.name as attribute_name,
    ca.is_enabled,
    ca.is_required,
    ca.sort_order
  FROM category_attributes ca
  JOIN categories c ON ca.category_id = c.id
  JOIN attributes a ON ca.attribute_id = a.id
  WHERE ca.category_id = 1001
  ORDER BY ca.sort_order
  LIMIT 10;
"
```

---

## Распределение данных

### По категориям (TOP 10):

```sql
SELECT
    category_id,
    COUNT(*) as attributes_count
FROM unified_category_attributes
GROUP BY category_id
ORDER BY attributes_count DESC
LIMIT 10;
```

**Результат:**
```
 category_id | attributes_count
-------------+-----------------
        1301 |              34  -- Lični automobili
        1103 |              27
        1003 |              18  -- Automobili
        1401 |              17
        1101 |              17
        1102 |              16
        1104 |              16
        1402 |              15
        1302 |              13
        1202 |              11
```

### По атрибутам:

```sql
SELECT COUNT(DISTINCT attribute_id) FROM unified_category_attributes;
-- Результат: ~150+ уникальных атрибутов
```

---

## Troubleshooting

### Проблема: Foreign key violation на category_id

**Причина:** Категория не существует в микросервисе

**Решение:**
1. Проверить какие категории отсутствуют:
```bash
psql "postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5433/svetubd?sslmode=disable" -c "
  SELECT DISTINCT ca.category_id
  FROM unified_category_attributes ca
  WHERE NOT EXISTS (
    SELECT 1 FROM categories c WHERE c.id = ca.category_id
  );
"
```

2. Сначала мигрировать недостающие категории
3. Повторить миграцию category_attributes

### Проблема: Foreign key violation на attribute_id

**Причина:** Атрибут не был мигрирован

**Решение:**
1. Проверить миграцию атрибутов:
```bash
/p/github.com/sveturs/listings/scripts/validate_attributes_migration.sh
```

2. При необходимости повторить миграцию атрибутов
3. Повторить миграцию category_attributes

### Проблема: Duplicate key violation

**Причина:** Запись с такой парой (category_id, attribute_id) уже существует

**Решение:**
- Миграция использует `ON CONFLICT DO UPDATE`, поэтому дубликаты автоматически обновляются
- Если проблема сохраняется, очистить таблицу и повторить миграцию:
```sql
TRUNCATE TABLE category_attributes CASCADE;
```

---

## Rollback

### Полный откат миграции:

```bash
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" -c "
  TRUNCATE TABLE category_attributes;
"
```

### Частичный откат (удалить только мигрированные записи):

```bash
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" -c "
  DELETE FROM category_attributes
  WHERE created_at >= '2025-11-17';
"
```

---

## Зависимости

### Должны быть выполнены ДО этой миграции:
1. ✅ Миграция категорий (categories)
2. ✅ Миграция атрибутов (attributes)

### Эта миграция требуется ДЛЯ:
1. Listing values (связь атрибутов с конкретными объявлениями)
2. Фильтрация и поиск по атрибутам
3. Валидация атрибутов при создании/редактировании объявлений

---

## SQL запросы для анализа

### Проверить распределение по is_enabled:
```sql
SELECT is_enabled, COUNT(*)
FROM category_attributes
GROUP BY is_enabled;
```

### Проверить распределение по is_required:
```sql
SELECT is_required, COUNT(*)
FROM category_attributes
GROUP BY is_required;
```

### Найти категории с наибольшим количеством атрибутов:
```sql
SELECT
    c.id,
    c.name,
    COUNT(ca.id) as attributes_count
FROM categories c
LEFT JOIN category_attributes ca ON c.id = ca.category_id
GROUP BY c.id, c.name
ORDER BY attributes_count DESC
LIMIT 10;
```

### Проверить category_specific_options:
```sql
SELECT
    category_id,
    attribute_id,
    category_specific_options
FROM category_attributes
WHERE category_specific_options IS NOT NULL
LIMIT 10;
```

---

## Changelog

### 2025-11-17
- ✅ Создан Go migration tool
- ✅ Создан validation script
- ✅ Написана документация
- ⏳ Dry-run тестирование
- ⏳ Реальная миграция

---

## Полезные ссылки

- [Attributes Migration Guide](./ATTRIBUTES_MIGRATION.md)
- [Categories Migration Guide](./CATEGORIES_MIGRATION.md) (если существует)
- [Listings Microservice README](../README.md)

---

**Автор:** Automated Migration Tool
**Дата создания:** 2025-11-17
**Последнее обновление:** 2025-11-17

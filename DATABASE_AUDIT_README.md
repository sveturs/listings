# Listings Database Audit - Навигация по документации

**Дата аудита:** 2025-12-16
**База данных:** listings_dev_db (PostgreSQL 15, порт 35434)
**Статус:** ⚠️ Cleanup ожидает выполнения

---

## 📚 Документация

### 1. AUDIT_SUMMARY.txt
**Краткая сводка с визуализацией**

Быстрый обзор результатов аудита в ASCII-формате:
- Общая статистика БД
- Критические проблемы (пустые таблицы, дубликаты, NULL-колонки)
- План действий
- Команды для выполнения cleanup

```bash
cat /p/github.com/vondi-global/listings/AUDIT_SUMMARY.txt
```

---

### 2. LISTINGS_DB_AUDIT_REPORT.md
**Полный детальный отчёт**

Подробный анализ всех аспектов БД:
- Список всех пустых таблиц с рекомендациями
- Детали по каждому дублирующемуся индексу
- Анализ колонок-всегда-NULL
- Проверка orphan records, PK, FK
- Статистика использования индексов
- Materialized views
- Полный план действий по фазам

```bash
# Открыть в редакторе
code /p/github.com/vondi-global/listings/LISTINGS_DB_AUDIT_REPORT.md

# Или просмотреть в терминале
less /p/github.com/vondi-global/listings/LISTINGS_DB_AUDIT_REPORT.md
```

---

### 3. cleanup_rudiments.sql
**Исполняемый SQL скрипт**

Ready-to-run SQL для удаления всех рудиментов:
- ФАЗА 1: Удаление 16 дублирующихся индексов
- ФАЗА 2: Удаление 9 колонок-рудиментов
- ФАЗА 3: Удаление 7 пустых таблиц
- ФАЗА 4: Оптимизация (ANALYZE)
- Проверка результатов
- Rollback план

**Безопасность:**
- Всё в одной транзакции (BEGIN...COMMIT)
- Можно откатить через ROLLBACK
- Использует IF EXISTS

```sql
-- Предварительный просмотр
less /p/github.com/vondi-global/listings/cleanup_rudiments.sql

-- Выполнение (ТОЛЬКО после backup!)
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db" \
  -f cleanup_rudiments.sql
```

---

### 4. CLEANUP_CHECKLIST.md
**Пошаговый чек-лист выполнения**

Детальная инструкция для безопасного выполнения cleanup:
- ✅ Pre-Flight Checklist
- Шаг 1: Backup БД
- Шаг 2: Остановка сервисов
- Шаг 3: Проверка состояния (до)
- Шаг 4: Выполнение cleanup
- Шаг 5: Проверка состояния (после)
- Шаг 6: Запуск сервисов
- Шаг 7: Функциональное тестирование
- Шаг 8: Оптимизация
- Шаг 9: Документация изменений
- Шаг 10: Финальная проверка
- 🚨 Rollback Plan

```bash
# Открыть чек-лист
code /p/github.com/vondi-global/listings/CLEANUP_CHECKLIST.md

# Или работать построчно
cat /p/github.com/vondi-global/listings/CLEANUP_CHECKLIST.md
```

---

## 🚀 Quick Start

### Вариант 1: Быстрый cleanup (для опытных)

```bash
# 1. Backup
mkdir -p /p/github.com/vondi-global/listings/backups
cd /p/github.com/vondi-global/listings/backups
pg_dump -h localhost -p 35434 -U listings_user \
  -d listings_dev_db \
  -f "listings_dev_db_backup_$(date +%Y%m%d_%H%M%S).sql"

# 2. Остановить сервисы
/home/dim/.local/bin/stop-listings-microservice.sh

# 3. Выполнить cleanup
cd /p/github.com/vondi-global/listings
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db" \
  -f cleanup_rudiments.sql

# 4. Запустить сервисы
/home/dim/.local/bin/start-listings-microservice.sh

# 5. Проверить
curl http://localhost:8086/health
```

### Вариант 2: Пошаговый cleanup (рекомендуется)

Следовать инструкциям в **CLEANUP_CHECKLIST.md** - отмечать каждый шаг.

---

## 📊 Ожидаемые результаты

| Метрика | До | После | Экономия |
|---------|-----|-------|----------|
| **Таблиц** | 38 | ~31 | -7 (18%) |
| **Индексов** | 281 | ~200-220 | -60-80 (25-30%) |
| **Размер БД** | 16 MB | ~14-15 MB | ~1-2 MB (10%) |
| **NULL-колонок** | 33 | 24 | -9 |

---

## 🎯 Что будет удалено

### Таблицы (7 шт):
- `analytics_events` - заменено на `listing_stats`
- `attribute_options` - устарела
- `attribute_search_cache` - не используется
- `b2c_product_variants` - не нужна
- `category_variant_attributes` - рудимент
- `listing_attribute_values` - дубликат
- `variant_attribute_values` - не используется

### Индексы (16 шт):
Дублирующие индексы на `listings`, `categories`, `attributes`, `storefronts`, `orders`, `shopping_carts`, `c2c_chats` и других таблицах.

### Колонки (9 шт):
- `attributes.legacy_product_variant_attribute_id`
- `categories.external_id`
- `category_attributes.*` (5 колонок)
- `orders.notes`, `orders.shipping_method`

---

## 🔍 Проверка текущего состояния

```bash
# Размер БД
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db" \
  -c "SELECT pg_size_pretty(pg_database_size('listings_dev_db'));"

# Количество таблиц
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db" \
  -c "SELECT COUNT(*) FROM pg_tables WHERE schemaname = 'public';"

# Количество индексов
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db" \
  -c "SELECT COUNT(*) FROM pg_indexes WHERE schemaname = 'public';"

# Проверить рудиментные таблицы
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db" \
  -c "SELECT tablename, pg_size_pretty(pg_total_relation_size('public.'||tablename))
      FROM pg_tables
      WHERE schemaname = 'public'
      AND tablename IN (
        'analytics_events',
        'attribute_options',
        'attribute_search_cache',
        'b2c_product_variants',
        'category_variant_attributes',
        'listing_attribute_values',
        'variant_attribute_values'
      );"
```

---

## ✅ Что НЕ будет затронуто

Все критические таблицы остаются нетронутыми:
- ✅ `listings` - основная таблица объявлений
- ✅ `categories` - категории
- ✅ `storefronts` - витрины
- ✅ `orders`, `order_items` - заказы
- ✅ `chats`, `messages` - чаты
- ✅ `listing_images` - изображения
- ✅ `listing_favorites` - избранное
- ✅ `shopping_carts`, `cart_items` - корзина
- ✅ `listing_attributes` - атрибуты
- ✅ `inventory_movements`, `inventory_reservations` - инвентарь

---

## 🚨 Безопасность

### Backup стратегия:
1. **Full dump** перед cleanup (обязательно!)
2. **Compressed dump** для архива (рекомендуется)
3. **Test restore** на копии БД (опционально)

### Rollback план:
Если что-то пошло не так:
```bash
# Остановить сервисы
/home/dim/.local/bin/stop-listings-microservice.sh

# Восстановить из backup
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db" \
  -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"

psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db" \
  -f backups/listings_dev_db_backup_YYYYMMDD_HHMMSS.sql

# Запустить сервисы
/home/dim/.local/bin/start-listings-microservice.sh
```

---

## 📞 Support

При возникновении проблем:
1. Проверить логи: `tail -f /tmp/listings-microservice.log`
2. Проверить подключение к БД: `psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db"`
3. Использовать rollback план из `CLEANUP_CHECKLIST.md`

---

## 📝 История изменений

| Дата | Версия | Изменения |
|------|--------|-----------|
| 2025-12-16 | 1.0 | Первый аудит БД, создание cleanup плана |

---

**Создано:** Claude Code
**Последнее обновление:** 2025-12-16

# Listings Database Cleanup - Чек-лист выполнения

**Дата:** 2025-12-16
**Статус:** ⚠️ ОЖИДАЕТ ВЫПОЛНЕНИЯ
**Ответственный:** DevOps / Database Admin

---

## ✅ Pre-Flight Checklist

Перед выполнением cleanup убедись:

- [ ] **Backup БД создан** (обязательно!)
- [ ] **Сервис listings остановлен**
- [ ] **Монолит остановлен** (на всякий случай)
- [ ] **Проверено использование таблиц в коде** (grep показал отсутствие)
- [ ] **Production БД не затронута** (работаем только с Dev)

---

## 📋 Шаг 1: Backup

```bash
# 1. Создать директорию для бэкапов
mkdir -p /p/github.com/vondi-global/listings/backups
cd /p/github.com/vondi-global/listings/backups

# 2. Создать backup БД
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
pg_dump -h localhost -p 35434 -U listings_user \
  -d listings_dev_db \
  --no-owner --no-acl \
  -f "listings_dev_db_before_cleanup_${TIMESTAMP}.sql"

# 3. Проверить размер backup
ls -lh "listings_dev_db_before_cleanup_${TIMESTAMP}.sql"

# 4. Создать compressed backup (опционально)
gzip -k "listings_dev_db_before_cleanup_${TIMESTAMP}.sql"
```

**Результат:**
- [ ] Backup создан
- [ ] Размер backup: ___________ MB
- [ ] Файл: `listings_dev_db_before_cleanup_YYYYMMDD_HHMMSS.sql`

---

## 📋 Шаг 2: Остановка сервисов

```bash
# 1. Остановить Listings Microservice
/home/dim/.local/bin/stop-listings-microservice.sh

# 2. Проверить, что процессы остановлены
netstat -tlnp | grep ":50053"  # Должно быть пусто
screen -ls | grep listings      # Должно быть пусто

# 3. Остановить монолит (опционально, на всякий случай)
/home/dim/.local/bin/kill-port-3000.sh
```

**Результат:**
- [ ] Listings Microservice остановлен
- [ ] Порт 50053 свободен
- [ ] Screen сессии закрыты

---

## 📋 Шаг 3: Проверка состояния БД (до cleanup)

```bash
# 1. Подключиться к БД
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db"

# 2. Запустить проверки
\dt  -- Список таблиц (должно быть 38)

SELECT pg_size_pretty(pg_database_size('listings_dev_db'));  -- Размер БД

SELECT count(*) FROM pg_indexes WHERE schemaname = 'public';  -- Количество индексов

\q
```

**Результат (до cleanup):**
- [ ] Таблиц: 38
- [ ] Размер БД: ___________ MB (ожидается ~16 MB)
- [ ] Индексов: ___________ (ожидается ~281)

---

## 📋 Шаг 4: Выполнение Cleanup

```bash
# 1. Перейти в директорию listings
cd /p/github.com/vondi-global/listings

# 2. Применить cleanup скрипт
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db" \
  -f cleanup_rudiments.sql

# 3. Проверить вывод
# Должно быть: BEGIN -> ... -> COMMIT
# НЕ должно быть: ERROR или ROLLBACK
```

**Результат:**
- [ ] Скрипт выполнен успешно (COMMIT)
- [ ] Нет ошибок в выводе
- [ ] Транзакция завершена

---

## 📋 Шаг 5: Проверка состояния БД (после cleanup)

```bash
# 1. Подключиться к БД
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db"

# 2. Запустить проверки
\dt  -- Список таблиц

SELECT pg_size_pretty(pg_database_size('listings_dev_db'));  -- Размер БД

SELECT count(*) FROM pg_indexes WHERE schemaname = 'public';  -- Количество индексов

-- Проверить, что критические таблицы существуют
SELECT COUNT(*) FROM listings;
SELECT COUNT(*) FROM categories;
SELECT COUNT(*) FROM storefronts;
SELECT COUNT(*) FROM orders;

-- Проверить, что рудименты удалены
SELECT tablename FROM pg_tables
WHERE schemaname = 'public'
AND tablename IN (
  'analytics_events',
  'attribute_options',
  'attribute_search_cache',
  'b2c_product_variants',
  'category_variant_attributes',
  'listing_attribute_values',
  'variant_attribute_values'
);
-- Должно быть: 0 rows

\q
```

**Результат (после cleanup):**
- [ ] Таблиц: _______ (ожидается ~31, было 38)
- [ ] Размер БД: _______ MB (ожидается ~14-15 MB, было 16 MB)
- [ ] Индексов: _______ (ожидается ~200-220, было 281)
- [ ] Рудиментные таблицы удалены: ✅
- [ ] Критические таблицы на месте: ✅

---

## 📋 Шаг 6: Запуск сервисов

```bash
# 1. Запустить Listings Microservice
/home/dim/.local/bin/start-listings-microservice.sh

# 2. Проверить логи
tail -f /tmp/listings-microservice.log
# Ждём: "Server listening on [::]:50053"
# НЕ должно быть: "relation ... does not exist"

# 3. Проверить health check
curl http://localhost:8086/health
# Ожидается: {"status":"healthy"}
```

**Результат:**
- [ ] Listings Microservice запущен
- [ ] Health check проходит
- [ ] Нет ошибок в логах

---

## 📋 Шаг 7: Функциональное тестирование

```bash
# 1. Получить JWT токен
TOKEN=$(cat /tmp/token)

# 2. Тест: получить список объявлений
curl -s -H "Authorization: Bearer $TOKEN" \
  http://localhost:3000/api/v1/marketplace/listings?limit=5 | jq '.data | length'
# Ожидается: число > 0

# 3. Тест: получить избранное
curl -s -H "Authorization: Bearer $TOKEN" \
  http://localhost:3000/api/v1/marketplace/favorites | jq '.data'
# Ожидается: массив (может быть пустым)

# 4. Тест: получить категории
curl -s http://localhost:3000/api/v1/marketplace/categories | jq '.data | length'
# Ожидается: число > 0

# 5. Тест: получить заказы (если есть)
curl -s -H "Authorization: Bearer $TOKEN" \
  http://localhost:3000/api/v1/marketplace/orders | jq '.data'
# Ожидается: массив
```

**Результат:**
- [ ] Listings API работает
- [ ] Favorites API работает
- [ ] Categories API работает
- [ ] Orders API работает
- [ ] Нет ошибок 500

---

## 📋 Шаг 8: Оптимизация после cleanup

```bash
# 1. Обновить статистику планировщика
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db" \
  -c "ANALYZE;"

# 2. Освободить место (опционально, может быть долгим)
# psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db" \
#   -c "VACUUM FULL ANALYZE;"
```

**Результат:**
- [ ] ANALYZE выполнен
- [ ] VACUUM выполнен (если нужно)

---

## 📋 Шаг 9: Документация изменений

```bash
# 1. Обновить CHANGELOG
cat >> /p/github.com/vondi-global/listings/CHANGELOG.md << EOF

## [Unreleased] - $(date +%Y-%m-%d)

### Removed
- Удалены пустые рудиментные таблицы (7 шт):
  - analytics_events
  - attribute_options
  - attribute_search_cache
  - b2c_product_variants
  - category_variant_attributes
  - listing_attribute_values
  - variant_attribute_values
- Удалены дублирующиеся индексы (16 шт)
- Удалены неиспользуемые колонки (9 шт):
  - attributes.legacy_product_variant_attribute_id
  - categories.external_id
  - category_attributes (5 колонок)
  - orders.notes, orders.shipping_method

### Optimized
- Экономия дискового пространства: ~1.5-2 MB
- Сокращение количества индексов: ~60-80 шт
- Упрощение схемы БД на ~18%
EOF
```

**Результат:**
- [ ] CHANGELOG обновлён
- [ ] Изменения задокументированы

---

## 📋 Шаг 10: Финальная проверка

```bash
# 1. Создать финальный отчёт
cat > /p/github.com/vondi-global/listings/CLEANUP_REPORT.txt << EOF
Listings Database Cleanup Report
=================================
Date: $(date +"%Y-%m-%d %H:%M:%S")
Performed by: ${USER}

BEFORE CLEANUP:
- Tables: 38
- Indexes: 281
- Database size: 16 MB

AFTER CLEANUP:
- Tables: $(psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db" -t -c "SELECT COUNT(*) FROM pg_tables WHERE schemaname = 'public';")
- Indexes: $(psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db" -t -c "SELECT COUNT(*) FROM pg_indexes WHERE schemaname = 'public';")
- Database size: $(psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db" -t -c "SELECT pg_size_pretty(pg_database_size('listings_dev_db'));")

REMOVED:
- Tables: 7 (rudiments)
- Indexes: 16+ (duplicates)
- Columns: 9 (always NULL)

SPACE SAVED: ~1.5-2 MB

STATUS: ✅ SUCCESS
EOF

cat /p/github.com/vondi-global/listings/CLEANUP_REPORT.txt
```

**Результат:**
- [ ] Отчёт создан
- [ ] Все метрики собраны
- [ ] Cleanup успешно завершён

---

## 🚨 Rollback Plan (на случай проблем)

Если после cleanup возникли проблемы:

```bash
# 1. Остановить сервисы
/home/dim/.local/bin/stop-listings-microservice.sh
/home/dim/.local/bin/kill-port-3000.sh

# 2. Найти backup
cd /p/github.com/vondi-global/listings/backups
ls -lht

# 3. Восстановить из backup
BACKUP_FILE="listings_dev_db_before_cleanup_YYYYMMDD_HHMMSS.sql"
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db" \
  -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"

psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db" \
  -f "$BACKUP_FILE"

# 4. Перезапустить сервисы
/home/dim/.local/bin/start-listings-microservice.sh
```

---

## 📊 Ожидаемые результаты

| Метрика | До cleanup | После cleanup | Изменение |
|---------|------------|---------------|-----------|
| Таблиц | 38 | ~31 | -7 (18%) |
| Индексов | 281 | ~200-220 | -60-80 (25-30%) |
| Размер БД | 16 MB | ~14-15 MB | -1-2 MB (10%) |
| Колонок-NULL | 33 | 24 | -9 |

---

## ✅ Finalize

После успешного выполнения:

- [ ] Все чек-боксы отмечены
- [ ] Сервисы работают
- [ ] Тесты проходят
- [ ] Backup сохранён
- [ ] Документация обновлена
- [ ] Отчёт создан

**Статус cleanup:** ✅ ЗАВЕРШЁН / ⚠️ ПРОБЛЕМЫ / ❌ ROLLBACK

---

**Примечания:**

_Добавь сюда любые заметки или проблемы, возникшие в процессе выполнения..._

---

**Подготовил:** Claude Code
**Версия:** 1.0
**Дата:** 2025-12-16

# Миграция C2C/B2C - Завершение

**Дата:** 2025-10-09
**Ветка:** `feature/c2c-b2c-migration`
**Статус:** ✅ Завершено

## 📋 Обзор

Успешно выполнена миграция терминологии проекта:
- `marketplace` → `c2c` (Customer-to-Customer)
- `storefronts` → `b2c` (Business-to-Customer)

Миграция включала переименование модулей, обновление базы данных, OpenSearch индексов и всего frontend/backend кода.

---

## ✅ Выполненные фазы

### Фаза 0-4: Подготовка и миграция кода (выполнено ранее)
- ✅ Создание ветки feature/c2c-b2c-migration
- ✅ Создание backup перед миграцией
- ✅ Создание naming-map.json для отслеживания переименований
- ✅ Миграции базы данных (000172-000173):
  - Таблицы: `marketplace_*` → `c2c_*`
  - Таблицы: `storefront_*` → `b2c_*`
- ✅ Переименование backend модулей:
  - `internal/proj/marketplace` → `internal/proj/c2c`
  - `internal/proj/storefronts` → `internal/proj/b2c`
- ✅ Переименование frontend компонентов, routes, i18n

### Фаза 5: OpenSearch миграция ✅
**Дата выполнения:** 2025-10-09

#### Создание новых индексов
- ✅ Создан скрипт миграции: `backend/migrate_opensearch_indexes.py`
- ✅ Созданы индексы:
  - `c2c_listings` (вместо marketplace_listings)
  - `b2c_products` (вместо storefront_products)
- ✅ Перенесено 7 документов из marketplace_listings → c2c_listings

#### Обновление backend конфигурации
**Файл:** `backend/internal/config/config.go`
```go
type OpenSearchConfig struct {
    URL              string `yaml:"url"`
    Username         string `yaml:"username"`
    Password         string `yaml:"password"`
    MarketplaceIndex string `yaml:"marketplace_index"` // Deprecated
    C2CIndex         string `yaml:"c2c_index"`          // NEW
    B2CIndex         string `yaml:"b2c_index"`          // NEW
}
```

**Файл:** `backend/internal/server/server.go:147`
```go
// Старый код:
// db, err := postgres.NewDatabase(ctx, cfg.DatabaseURL, osClient, cfg.OpenSearch.MarketplaceIndex, fileStorage, cfg.SearchWeights)

// Новый код:
db, err := postgres.NewDatabase(ctx, cfg.DatabaseURL, osClient, cfg.OpenSearch.C2CIndex, fileStorage, cfg.SearchWeights)
```

### Фаза 6: MinIO/S3 миграция ✅
**Решение:** Пропущено - существующие buckets можно продолжать использовать без переименования.

### Фаза 7: Тестирование ✅
- ✅ Backend компиляция успешна (87MB binary)
- ✅ Frontend сборка успешна (64.20s)
- ✅ Backend запуск успешен, использует новый индекс `c2c_listings`
- ✅ API endpoints работают корректно

### Фаза 8: Pre-commit проверка ✅
**Дата выполнения:** 2025-10-09

#### Backend ✅
1. **Format:** `make format` - успешно
2. **Lint:** `make lint` - успешно (исправлены 2 проблемы):
   - Переписан if-else на switch в `opensearch/repository.go:1447`
   - Исправлен формат deprecated коммента в `chat.go:834`

#### B2C Naming Convention ✅
Переименован метод `GetStorefrontProductImages` → `GetB2CProductImages` в 8 файлах:
- `internal/storage/storage.go:78` (интерфейс)
- `internal/storage/postgres/db.go:789-790` (делегирование)
- `internal/proj/c2c/storage/postgres/marketplace.go:3614` (реализация)
- `internal/proj/c2c/storage/postgres/marketplace.go:3121, 3405` (вызовы)
- `internal/proj/c2c/storage/opensearch/repository.go:1469` (использование)
- `internal/proj/c2c/service/category_test.go:1388` (mock)
- `internal/proj/c2c/service/integration_test.go:292` (mock)

#### Frontend ✅
1. **Format:** `yarn format` - успешно (16.16s)
2. **Lint:** `yarn lint` - успешно (исправлена 1 проблема):
   - Заменен `<a href="/b2c">` на `<Link href="/b2c">` в `ideal-homepage/page.tsx:163`

---

## 📊 Статистика миграции

### OpenSearch
- **Индексы созданы:** 2 (c2c_listings, b2c_products)
- **Документов перенесено:** 7
- **Старые индексы:** сохранены для возможного rollback

### Backend
- **Модулей переименовано:** 2 (marketplace→c2c, storefronts→b2c)
- **Файлов изменено:** 15+
- **Методов переименовано:** 8 (GetStorefrontProductImages)
- **Lint issues исправлено:** 2

### Frontend
- **Компонентов обновлено:** множество
- **Routes изменено:** /marketplace→/c2c, /storefronts→/b2c
- **i18n keys обновлено:** множество
- **Lint issues исправлено:** 1

---

## 🔧 Изменения в конфигурации

### Backend Config (YAML)
```yaml
opensearch:
  url: "http://localhost:9200"
  c2c_index: "c2c_listings"    # NEW
  b2c_index: "b2c_products"    # NEW
  marketplace_index: "..."      # DEPRECATED
```

### Environment Variables
Не требуется изменений - все настройки в YAML config.

---

## 🚀 Следующие шаги

### После верификации
1. **Удалить старые OpenSearch индексы (опционально):**
   ```bash
   curl -X DELETE http://localhost:9200/marketplace_listings
   curl -X DELETE http://localhost:9200/storefront_products
   ```

2. **Удалить deprecated код:**
   - Удалить поле `MarketplaceIndex` из `OpenSearchConfig`
   - Обновить все ссылки на старую терминологию в комментариях

3. **Обновить документацию:**
   - API документация (Swagger)
   - README файлы
   - Архитектурные диаграммы

### Деплой на production
1. **Backup production БД и OpenSearch**
2. **Выполнить миграции:**
   ```bash
   # На production сервере
   cd backend && ./migrator up
   ```
3. **Запустить скрипт миграции OpenSearch:**
   ```bash
   python3 migrate_opensearch_indexes.py
   ```
4. **Обновить конфигурацию:**
   ```yaml
   opensearch:
     c2c_index: "c2c_listings"
     b2c_index: "b2c_products"
   ```
5. **Перезапустить сервисы**
6. **Мониторинг логов и метрик**

---

## 📝 Важные заметки

### Backward Compatibility
- ✅ Старые индексы OpenSearch НЕ удалены (можно откатиться)
- ✅ База данных поддерживает старые и новые названия таблиц через views
- ⚠️ Frontend routes изменились: `/marketplace` → `/c2c`, `/storefronts` → `/b2c`

### Breaking Changes
- ⚠️ API endpoints изменились (если используются внешними клиентами)
- ⚠️ OpenSearch query paths изменились

### Rollback Plan
В случае проблем:
1. Откатить миграции БД: `./migrator down`
2. Вернуться к старой ветке: `git checkout main`
3. Восстановить старые индексы OpenSearch из backup

---

## ✅ Checklist финальной проверки

- [x] Backend компилируется без ошибок
- [x] Frontend собирается без ошибок
- [x] Backend lint проходит без ошибок
- [x] Frontend lint проходит без ошибок
- [x] OpenSearch индексы созданы и заполнены
- [x] Backend использует новые индексы
- [x] API endpoints работают
- [x] Все тесты проходят
- [x] Pre-commit hooks настроены

---

## 📚 Связанные документы

- [Детальный план миграции](C2C_B2C_MIGRATION_PLAN_DETAILED.md)
- [Naming Map](naming-map.json)
- [Скрипт миграции OpenSearch](../backend/migrate_opensearch_indexes.py)
- [Миграции БД](../backend/migrations/000172_rename_marketplace_to_c2c.up.sql)

---

**Миграция завершена успешно!** 🎉

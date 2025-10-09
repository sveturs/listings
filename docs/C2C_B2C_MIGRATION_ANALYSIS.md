# 📊 Комплексная оценка миграции marketplace → c2c, storefronts → b2c

**Дата анализа**: 2025-10-08
**Автор**: Claude Code Analysis
**Статус**: Предварительная оценка

⚡ **ВАЖНО**: Проект НЕ в продакшене - обратная совместимость НЕ требуется!

---

## 🔍 Анализ масштаба изменений

### Статистика по коду:

**Backend:**
- 📁 **13,021 упоминаний** в **297 файлах**
- 🗂️ **3 модуля проекта**: `internal/proj/marketplace`, `internal/proj/storefront`, `internal/proj/storefronts`
- 📄 **~47 миграций** содержат упоминания

**Frontend:**
- 📁 **5,053 упоминаний** в **381 файлах**
- 🗂️ **7 директорий**: app/marketplace, app/storefronts, components/marketplace, components/storefronts, и др.
- 📝 **Переводы**: 3 языка × множество файлов (marketplace.json, storefronts.json и т.д.)

**База данных:**
- 🗄️ **23 таблицы** (общий размер ~4.5 MB):
  - `marketplace_*`: 8 таблиц (2.6 MB)
  - `storefront_*`: 15 таблиц (1.9 MB)

**OpenSearch:**
- 🔍 **3 индекса**:
  - `marketplace_listings` (12 документов, 1.3 MB)
  - `storefront_products` (0 документов, 3.3 KB)
  - `storefronts` (1 документ, 12.5 KB)

**S3/MinIO:**
- 📦 Структура: `marketplace-images/`, `storefront-images/`

---

## 💡 Моя оценка сложности

### ✅ **РЕАЛИЗУЕМО И УПРОЩЕНО**

**Общая оценка**: 🟡 **7/10** по сложности

**Время выполнения**: 📅 **2-4 недели** (при полной занятости)

**Риски**: ⚠️ **СРЕДНИЕ** - проект не в продакшене, можно мигрировать без оглядки на пользователей

**Преимущества**: 🚀 Не нужна обратная совместимость - миграция в 1.5-2 раза быстрее!

---

## 🎯 Почему это имеет смысл?

### **Преимущества:**

1. ✅ **Семантическая ясность**
   - `c2c` (customer-to-customer) — объявления от пользователей
   - `b2c` (business-to-customer) — товары от витрин/магазинов

2. ✅ **Меньше путаницы в коде**
   - Нет двусмысленности: "marketplace" = рынок или объявление?
   - "storefront" = витрина или продукт витрины?

3. ✅ **Лучшая масштабируемость**
   - Легче добавить в будущем b2b (business-to-business)
   - Логическое разделение бизнес-моделей

4. ✅ **Соответствие индустрии**
   - C2C/B2C/B2B — стандартная терминология e-commerce

### **Текущие проблемы:**
- `marketplace_listings` — непонятно, что это C2C
- `storefront_products` — длинное название
- Смешение понятий в API endpoints

---

## 📋 Детальный план миграции

### **Фаза 1: Подготовка (3-5 дней)**

#### 1.1 Создание маппинга имен
```bash
# Создать файл migration_map.json
marketplace_listings → c2c_listings
marketplace_categories → c2c_categories
marketplace_images → c2c_images
marketplace_chats → c2c_chats
marketplace_messages → c2c_messages
marketplace_favorites → c2c_favorites
marketplace_orders → c2c_orders
marketplace_listing_variants → c2c_listing_variants

storefronts → b2c_stores
storefront_products → b2c_products
storefront_product_images → b2c_product_images
storefront_product_variants → b2c_product_variants
storefront_product_attributes → b2c_product_attributes
storefront_orders → b2c_orders
storefront_order_items → b2c_order_items
storefront_favorites → b2c_favorites
storefront_hours → b2c_store_hours
storefront_staff → b2c_store_staff
storefront_payment_methods → b2c_payment_methods
storefront_delivery_options → b2c_delivery_options
storefront_inventory_movements → b2c_inventory_movements
user_storefronts → user_b2c_stores
storefront_product_variant_images → b2c_product_variant_images
```

#### 1.2 Резервное копирование
```bash
# Полный дамп БД
pg_dump -Fc svetubd > backup_pre_migration.dump

# Резервная копия OpenSearch
curl -X PUT "localhost:9200/_snapshot/my_backup/snapshot_1?wait_for_completion=true"

# Бэкап MinIO
mc mirror local/marketplace-images local/marketplace-images-backup
mc mirror local/storefront-images local/storefront-images-backup
```

#### 1.3 Создание тестовой среды
```bash
# Копия БД для тестов
createdb svetubd_migration_test
pg_restore -d svetubd_migration_test backup_pre_migration.dump
```

---

### **Фаза 2: Миграция базы данных (5-7 дней)**

#### 2.1 Создание новых таблиц с новыми именами
```sql
-- migration 000172_create_c2c_b2c_tables.up.sql

-- C2C таблицы
CREATE TABLE c2c_listings (LIKE marketplace_listings INCLUDING ALL);
CREATE TABLE c2c_categories (LIKE marketplace_categories INCLUDING ALL);
CREATE TABLE c2c_images (LIKE marketplace_images INCLUDING ALL);
-- ... и так далее для всех 8 таблиц

-- B2C таблицы
CREATE TABLE b2c_stores (LIKE storefronts INCLUDING ALL);
CREATE TABLE b2c_products (LIKE storefront_products INCLUDING ALL);
CREATE TABLE b2c_product_images (LIKE storefront_product_images INCLUDING ALL);
-- ... и так далее для всех 15 таблиц
```

#### 2.2 Копирование данных
```sql
-- Копирование данных (с сохранением ID)
INSERT INTO c2c_listings SELECT * FROM marketplace_listings;
INSERT INTO c2c_categories SELECT * FROM marketplace_categories;
-- ... и т.д.

INSERT INTO b2c_stores SELECT * FROM storefronts;
INSERT INTO b2c_products SELECT * FROM storefront_products;
-- ... и т.д.
```

#### 2.3 Пересоздание индексов и ограничений
```sql
-- Автоматически при LIKE INCLUDING ALL
-- Но нужно проверить и обновить внешние ключи с новыми именами
ALTER TABLE c2c_listings
  DROP CONSTRAINT IF EXISTS marketplace_listings_category_id_fkey,
  ADD CONSTRAINT c2c_listings_category_id_fkey
    FOREIGN KEY (category_id) REFERENCES c2c_categories(id);
```

#### 2.4 Обновление триггеров и функций
```sql
-- Найти все триггеры на marketplace/storefront таблицы
SELECT * FROM pg_trigger
WHERE tgrelid::regclass::text LIKE '%marketplace%'
   OR tgrelid::regclass::text LIKE '%storefront%';

-- Пересоздать их для новых таблиц
-- Например: trigger_update_marketplace_listings_updated_at → trigger_update_c2c_listings_updated_at
```

---

### **Фаза 3: Backend миграция (7-10 дней)**

#### 3.1 Переименование модулей проекта
```bash
# Переименовать директории
mv internal/proj/marketplace internal/proj/c2c
mv internal/proj/storefronts internal/proj/b2c

# Обновить импорты во ВСЕХ файлах
find . -name "*.go" -exec sed -i 's|internal/proj/marketplace|internal/proj/c2c|g' {} +
find . -name "*.go" -exec sed -i 's|internal/proj/storefronts|internal/proj/b2c|g' {} +
```

#### 3.2 Обновление моделей (domain/models)
```bash
# Переименовать файлы
mv internal/domain/models/marketplace_listing.go → c2c_listing.go
mv internal/domain/models/storefront_product.go → b2c_product.go

# Обновить структуры
type MarketplaceListing → type C2CListing
type StorefrontProduct → type B2CProduct
```

#### 3.3 Обновление SQL запросов
```bash
# Найти все упоминания таблиц в Go коде
grep -r "marketplace_listings" internal/
grep -r "storefront_products" internal/

# Заменить через sed или вручную
sed -i 's/marketplace_listings/c2c_listings/g' internal/proj/c2c/**/*.go
sed -i 's/storefront_products/b2c_products/g' internal/proj/b2c/**/*.go
```

#### 3.4 Обновление API endpoints
```go
// Было:
app.Get("/api/v1/marketplace/listings", handler.GetListings)
app.Get("/api/v1/storefronts/:id", handler.GetStorefront)

// Стало:
app.Get("/api/v1/c2c/listings", handler.GetListings)
app.Get("/api/v1/b2c/stores/:id", handler.GetStore)
```

**✅ Проект не в продакшене** - старые endpoints удаляются полностью!

---

### **Фаза 4: Frontend миграция (5-7 дней)**

#### 4.1 Переименование директорий
```bash
mv src/app/[locale]/marketplace → src/app/[locale]/c2c
mv src/app/[locale]/storefronts → src/app/[locale]/b2c
mv src/components/marketplace → src/components/c2c
mv src/components/storefronts → src/components/b2c
```

#### 4.2 Обновление типов
```typescript
// src/types/marketplace.ts → src/types/c2c.ts
export interface MarketplaceListing → export interface C2CListing

// src/types/storefront.ts → src/types/b2c.ts
export interface StorefrontProduct → export interface B2CProduct
```

#### 4.3 Обновление API клиентов
```typescript
// src/services/marketplaceApi.ts → src/services/c2cApi.ts
const response = await apiClient.get('/marketplace/listings')
// →
const response = await apiClient.get('/c2c/listings')
```

#### 4.4 Обновление переводов
```bash
# Переименовать файлы
mv src/messages/en/marketplace.json → c2c.json
mv src/messages/en/storefronts.json → b2c.json

# Обновить ключи переводов
"marketplace.title" → "c2c.title"
"storefronts.product" → "b2c.product"
```

#### 4.5 Обновление роутинга
```typescript
// next.config.ts или app router
// /marketplace/[id] → /c2c/[id]
// /storefronts/[slug] → /b2c/[slug]

// ✅ Проект не в продакшене - редиректы не нужны!
// Просто удаляем старые роуты и используем новые
```

---

### **Фаза 5: OpenSearch миграция (2-3 дня)**

#### 5.1 Создание новых индексов
```bash
# Создать c2c_listings индекс
curl -X PUT "localhost:9200/c2c_listings" \
  -H 'Content-Type: application/json' \
  -d @opensearch/c2c_mapping.json

# Создать b2c_products индекс
curl -X PUT "localhost:9200/b2c_products" \
  -H 'Content-Type: application/json' \
  -d @opensearch/b2c_products_mapping.json

# Создать b2c_stores индекс
curl -X PUT "localhost:9200/b2c_stores" \
  -H 'Content-Type: application/json' \
  -d @opensearch/b2c_stores_mapping.json
```

#### 5.2 Переиндексация данных
```bash
# Reindex из старых индексов в новые
curl -X POST "localhost:9200/_reindex" -H 'Content-Type: application/json' -d'
{
  "source": { "index": "marketplace_listings" },
  "dest": { "index": "c2c_listings" }
}'

curl -X POST "localhost:9200/_reindex" -H 'Content-Type: application/json' -d'
{
  "source": { "index": "storefront_products" },
  "dest": { "index": "b2c_products" }
}'
```

#### 5.3 Обновление кода поиска
```go
// internal/proj/c2c/storage/opensearch/repository.go
const indexName = "c2c_listings" // было "marketplace_listings"

// internal/proj/b2c/storage/opensearch/product_repository.go
const indexName = "b2c_products" // было "storefront_products"
```

---

### **Фаза 6: MinIO/S3 миграция (1-2 дня)**

#### 6.1 Переименование bucket'ов
```bash
# Вариант 1: Создать новые bucket'ы и скопировать
mc mb local/c2c-images
mc mb local/b2c-images

mc mirror local/marketplace-images local/c2c-images
mc mirror local/storefront-images local/b2c-images

# Вариант 2: Использовать алиасы (через код)
```

#### 6.2 Обновление кода загрузки файлов
```go
// internal/storage/minio/client.go
const c2cBucket = "c2c-images" // было "marketplace-images"
const b2cBucket = "b2c-images" // было "storefront-images"
```

---

### **Фаза 7: Тестирование (5-7 дней)**

#### 7.1 Unit тесты
```bash
# Backend
cd backend && go test ./... -v

# Frontend
cd frontend/svetu && yarn test
```

#### 7.2 Интеграционные тесты
```bash
# Проверить все API endpoints
curl http://localhost:3000/api/v1/c2c/listings
curl http://localhost:3000/api/v1/b2c/stores

# Проверить поиск
curl "http://localhost:9200/c2c_listings/_search?q=*"
```

#### 7.3 E2E тестирование
```bash
# Playwright/Cypress тесты
- Создание C2C объявления
- Создание B2C продукта
- Поиск и фильтрация
- Корзина и заказ
```

---

### **Фаза 8: Деплой и мониторинг (3-5 дней)**

#### 8.1 Staging деплой
```bash
# Развернуть на dev.svetu.rs
./deploy-to-dev.sh

# Smoke testing
curl https://devapi.svetu.rs/api/v1/c2c/listings
```

#### 8.2 Production деплой
```bash
# ✅ Проект не в продакшене - можно мигрировать сразу!
# 1. Обновить БД (создать новые таблицы, скопировать данные)
# 2. Обновить OpenSearch индексы
# 3. Развернуть backend с новыми endpoint'ами
# 4. Развернуть frontend с новыми путями
# 5. Мониторинг логов и метрик
```

#### 8.3 Удаление старых сущностей
```sql
-- ✅ Проект не в продакшене - можно удалять сразу после миграции!
-- После базовой проверки работоспособности (1-2 дня)
DROP TABLE marketplace_listings CASCADE;
DROP TABLE storefront_products CASCADE;
-- ... и т.д.

-- OpenSearch
DELETE /marketplace_listings
DELETE /storefront_products
```

---

## ⚠️ Критические риски

### 1. **Пропущенные упоминания**
- С 18,000+ упоминаний легко что-то пропустить
- **Решение**: Автоматизация через скрипты + тщательный code review

### 2. **Сломанные зависимости**
- Сторонние библиотеки могут ссылаться на старые имена
- **Решение**: Полная перепись кода, тщательное тестирование

### 3. **Данные пользователей**
- Избранное, корзина, история могут сломаться
- **Решение**: Тщательное тестирование + миграция данных с проверкой FK

### 4. **Downtime**
- Миграция БД может занять время (но проект не в продакшене!)
- **Решение**: Можно мигрировать в любое время, без ограничений

---

## 💰 Оценка трудозатрат

**✅ УПРОЩЕНО**: Проект не в продакшене - нет обратной совместимости!

| Фаза | Дни | Описание |
|------|-----|----------|
| Подготовка | 2-3 | Планирование, бэкапы, маппинг |
| БД миграция | 4-5 | Создание таблиц, копирование данных, FK |
| Backend | 5-7 | Переименование модулей, обновление кода |
| Frontend | 4-5 | Компоненты, типы, роутинг, переводы |
| OpenSearch | 2-3 | Новые индексы, reindex |
| MinIO/S3 | 1-2 | Копирование bucket'ов |
| Тестирование | 3-5 | Unit, integration, E2E тесты |
| Деплой | 2-3 | Staging, production, мониторинг |
| **ИТОГО** | **23-33 дня** | **~4-7 недель** |

**Экономия**: ~8-13 дней за счёт отсутствия обратной совместимости!

---

## 🎯 Рекомендации

### ✅ **ДА, стоит делать!**

**Преимущество**: Проект не в продакшене - можем мигрировать смело и быстро!

1. **Разбить на этапы**
   - Не делать всё сразу
   - Начать с backend + БД
   - Потом frontend
   - Постепенный rollout по модулям

2. **Автоматизировать**
   - Написать скрипты для массового rename
   - CI/CD проверки на пропущенные упоминания
   - Автоматические тесты для проверки миграции

3. **Когда лучше делать**
   - ✅ Сейчас самое время - нет пользователей!
   - Пока нет срочных фич
   - Можно спокойно тестировать и исправлять ошибки

---

## 📝 Альтернативные подходы

### **Вариант 1: Только внутренние имена (умеренная сложность)**
- Переименовать только в коде (модели, переменные)
- Оставить БД таблицы и API как есть
- **Плюс**: Меньше работы, нет breaking changes
- **Минус**: Не решает проблему полностью

### **Вариант 2: Новый микросервис (максимальная сложность)**
- Создать новые C2C/B2C сервисы с нуля
- Постепенно мигрировать функционал
- **Плюс**: Чистая архитектура
- **Минус**: Очень долго и дорого

### **Вариант 3: Полная миграция (текущий план)**
- Переименовать везде: БД, код, API, UI
- **Плюс**: Полное решение проблемы
- **Минус**: Высокая сложность

---

## 📊 Статистика анализа

**Файлы проанализированы:**
- Backend: 297 файлов Go
- Frontend: 381 файл TypeScript/React
- База данных: 23 таблицы
- OpenSearch: 3 индекса
- Миграции: 47 SQL файлов

**Общее количество упоминаний:**
- Backend: 13,021
- Frontend: 5,053
- **TOTAL: 18,074 упоминаний**

---

## 🚀 Следующие шаги

Если принято решение о миграции:

1. ✅ Создать отдельную feature ветку `feature/c2c-b2c-migration`
2. ✅ Написать автоматические скрипты для rename
3. ✅ Создать первую миграцию БД (создание новых таблиц)
4. ✅ Обновить backend модели и репозитории
5. ✅ Обновить API endpoints (без алиасов - проект не в продакшене!)
6. ✅ Обновить frontend (по модулям)
7. ✅ Тестирование на каждом этапе
8. ✅ Деплой через staging → production

**Статус**: ⏸️ Ожидает решения

**Преимущество**: 🚀 Проект не в продакшене - можно мигрировать быстро и смело!

---

**Последнее обновление**: 2025-10-08 (убраны требования обратной совместимости)

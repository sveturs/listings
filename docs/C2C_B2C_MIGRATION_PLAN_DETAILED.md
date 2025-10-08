# 🚀 PRODUCTION-READY ПЛАН МИГРАЦИИ C2C→B2C (Детальный)

**Статус документа**: 🔴 ЖИВОЙ ДОКУМЕНТ - ПОСТОЯННОЕ ОБНОВЛЕНИЕ ОБЯЗАТЕЛЬНО!
**Дата создания**: 2025-10-09
**Версия плана**: 1.0
**Цель**: Выход в продакшн с чистой архитектурой

---

## ⚠️ КРИТИЧЕСКИЕ ПРИНЦИПЫ

### 🎯 **ZERO TECHNICAL DEBT POLICY**
- ✅ Каждая фаза проверяется на 100% завершенность
- ✅ НЕТ "TODO" комментариев в финальном коде
- ✅ НЕТ временных workaround'ов
- ✅ НЕТ копипасты - только переиспользуемый код
- ✅ Полное покрытие тестами (unit + integration)

### 📋 **ОБЯЗАТЕЛЬНОЕ ОБНОВЛЕНИЕ ПЛАНА**
После каждой фазы:
1. ✅ Обновить статус выполнения (✅/🚧/❌)
3. ✅ Задокументировать проблемы и решения
4. ✅ Обновить риски (новые/закрытые)
5. ✅ Скорректировать оценки следующих фаз

### 🎖️ **PRODUCTION-GRADE QUALITY**
- ✅ Код проходит pre-commit hooks
- ✅ Линтеры: 0 warnings, 0 errors
- ✅ Тесты: 100% критических путей покрыто
- ✅ Документация обновлена (Swagger, CLAUDE.md)


---

## 📊 ГЛОБАЛЬНАЯ МЕТРИКА ПРОГРЕССА

| Фаза | Статус | Прогресс | Качество | Дата начала | Дата окончания |
|------|--------|----------|----------|-------------|----------------|
| 0. Инициализация | ⏸️ Pending | 0% | N/A | - | - |
| 1. Подготовка | ⏸️ Pending | 0% | N/A | - | - |
| 2. БД миграция | ⏸️ Pending | 0% | N/A | - | - |
| 3. Backend | ⏸️ Pending | 0% | N/A | - | - |
| 4. Frontend | ⏸️ Pending | 0% | N/A | - | - |
| 5. OpenSearch | ⏸️ Pending | 0% | N/A | - | - |
| 6. MinIO/S3 | ⏸️ Pending | 0% | N/A | - | - |
| 7. Тестирование | ⏸️ Pending | 0% | N/A | - | - |
| 8. Production Деплой | ⏸️ Pending | 0% | N/A | - | - |

**Общий прогресс**: 0% (0/8 фаз)

---

## 🎬 ФАЗА 0: ИНИЦИАЛИЗАЦИЯ (День 0)

### Цель
Подготовить инфраструктуру для миграции, создать точки отката.

### Задачи

#### 0.1 Git Workflow
```bash
# Создать feature ветку
git checkout -b feature/c2c-b2c-migration

# Создать backup tag ПЕРЕД началом
git tag -a migration-pre-start -m "Snapshot before C2C→B2C migration"
git push origin migration-pre-start

# Защита основной ветки (на GitHub/GitLab)
# - Требовать PR для всех изменений
# - Минимум 1 reviewer
# - Проходить CI/CD checks
```

#### 0.2 Резервные копии (КРИТИЧНО!)
```bash
# 1. PostgreSQL - полный дамп (SQL формат для читаемости и отладки)
PGPASSWORD=mX3g1XGhMRUZEX3l pg_dump \
  -h localhost \
  -U postgres \
  -d svetubd \
  --no-owner \
  --no-acl \
  --column-inserts \
  --inserts \
  -f /tmp/svetubd_migration_backup_$(date +%Y%m%d_%H%M%S).sql

# 2. OpenSearch snapshot
curl -X PUT "localhost:9200/_snapshot/migration_backup" -H 'Content-Type: application/json' -d'
{
  "type": "fs",
  "settings": {
    "location": "/var/opensearch/backups/migration_$(date +%Y%m%d)"
  }
}'

curl -X PUT "localhost:9200/_snapshot/migration_backup/snapshot_pre_migration?wait_for_completion=true"

# 3. MinIO/S3 - зеркалирование
mc mirror --preserve local/marketplace-images /tmp/backup/marketplace-images
mc mirror --preserve local/storefront-images /tmp/backup/storefront-images

# 4. Файлы конфигурации
tar -czf /tmp/config_backup_$(date +%Y%m%d).tar.gz \
  backend/config.yaml \
  frontend/svetu/.env.local \
  docker-compose.yml
```

#### 0.3 Тестовая среда
```bash
# Создать изолированную БД для тестов миграции
createdb svetubd_migration_test

# Восстановить дамп (SQL формат)
PGPASSWORD=mX3g1XGhMRUZEX3l psql \
  -h localhost \
  -U postgres \
  -d svetubd_migration_test \
  -f /tmp/svetubd_migration_backup_*.sql

# Обновить .env для тестов
DATABASE_URL_TEST=postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5432/svetubd_migration_test?sslmode=disable
```

#### 0.4 Скрипты миграции
```bash
# Создать директорию для утилит
mkdir -p /data/hostel-booking-system/migration-tools

# Скрипт поиска упоминаний (для контроля прогресса)
cat > /data/hostel-booking-system/migration-tools/check-references.sh << 'EOF'
#!/bin/bash
echo "=== Checking old references ==="
echo "Backend marketplace: $(grep -r "marketplace" backend/internal/proj --include="*.go" | wc -l)"
echo "Backend storefronts: $(grep -r "storefronts" backend/internal/proj --include="*.go" | wc -l)"
echo "Frontend marketplace: $(grep -r "marketplace" frontend/svetu/src --include="*.ts*" | wc -l)"
echo "Frontend storefronts: $(grep -r "storefronts" frontend/svetu/src --include="*.ts*" | wc -l)"
EOF
chmod +x /data/hostel-booking-system/migration-tools/check-references.sh
```

### Критерии приёмки (Definition of Done)
- [ ] ✅ Git branch создана и защищена
- [ ] ✅ 3 резервные копии (БД, OpenSearch, MinIO) проверены
- [ ] ✅ Тестовая БД работает и доступна
- [ ] ✅ Скрипты проверки созданы и протестированы
- [ ] ✅ План обновлён: дата начала, статус → 🚧 In Progress

### Риски и митигация
| Риск | Вероятность | Решение |
|------|-------------|---------|
| Бэкапы повреждены | Низкая | Проверить восстановление до начала миграции |
| Нехватка места | Средняя | Освободить минимум 10GB перед стартом |

### Фактические данные (заполнить после выполнения)
- **Время**: ___ часов (план: 4-6 часов)
- **Проблемы**: ___
- **Решения**: ___

---

## 📂 ФАЗА 1: ПОДГОТОВКА И МАППИНГ (Дни 1-2)

### Цель
Создать полный маппинг сущностей, автоматизировать проверки, подготовить миграционные скрипты.

### Задачи

#### 1.1 Создание маппинга имён (КРИТИЧНО!)

**Файл**: `/data/hostel-booking-system/migration-tools/naming-map.json`

```json
{
  "version": "1.0",
  "database_tables": {
    "marketplace_listings": "c2c_listings",
    "marketplace_categories": "c2c_categories",
    "marketplace_images": "c2c_images",
    "marketplace_chats": "c2c_chats",
    "marketplace_messages": "c2c_messages",
    "marketplace_favorites": "c2c_favorites",
    "marketplace_orders": "c2c_orders",
    "marketplace_listing_variants": "c2c_listing_variants",
    "storefronts": "b2c_stores",
    "storefront_products": "b2c_products",
    "storefront_product_images": "b2c_product_images",
    "storefront_product_variants": "b2c_product_variants",
    "storefront_product_attributes": "b2c_product_attributes",
    "storefront_orders": "b2c_orders",
    "storefront_order_items": "b2c_order_items",
    "storefront_favorites": "b2c_favorites",
    "storefront_hours": "b2c_store_hours",
    "storefront_staff": "b2c_store_staff",
    "storefront_payment_methods": "b2c_payment_methods",
    "storefront_delivery_options": "b2c_delivery_options",
    "storefront_inventory_movements": "b2c_inventory_movements",
    "user_storefronts": "user_b2c_stores",
    "storefront_product_variant_images": "b2c_product_variant_images"
  },
  "api_endpoints": {
    "/api/v1/marketplace": "/api/v1/c2c",
    "/api/v1/storefronts": "/api/v1/b2c/stores",
    "/api/v1/admin/marketplace": "/api/v1/admin/c2c",
    "/api/v1/admin/storefronts": "/api/v1/admin/b2c"
  },
  "go_modules": {
    "internal/proj/marketplace": "internal/proj/c2c",
    "internal/proj/storefronts": "internal/proj/b2c"
  },
  "go_types": {
    "MarketplaceListing": "C2CListing",
    "MarketplaceCategory": "C2CCategory",
    "MarketplaceImage": "C2CImage",
    "MarketplaceChat": "C2CChat",
    "MarketplaceMessage": "C2CMessage",
    "MarketplaceFavorite": "C2CFavorite",
    "MarketplaceOrder": "C2COrder",
    "Storefront": "B2CStore",
    "StorefrontProduct": "B2CProduct",
    "StorefrontProductImage": "B2CProductImage",
    "StorefrontProductVariant": "B2CProductVariant"
  },
  "frontend_routes": {
    "/marketplace": "/c2c",
    "/storefronts": "/b2c"
  },
  "frontend_components": {
    "components/marketplace": "components/c2c",
    "components/storefronts": "components/b2c"
  },
  "i18n_keys": {
    "marketplace": "c2c",
    "storefronts": "b2c"
  },
  "opensearch_indices": {
    "marketplace_listings": "c2c_listings",
    "storefront_products": "b2c_products",
    "storefronts": "b2c_stores"
  },
  "minio_buckets": {
    "marketplace-images": "c2c-images",
    "storefront-images": "b2c-images"
  }
}
```

#### 1.2 Валидационный скрипт

**Файл**: `/data/hostel-booking-system/migration-tools/validate-mapping.py`

```python
#!/usr/bin/env python3
"""
Валидирует полноту маппинга и находит потенциальные пропуски.
"""
import json
import os
import re
import subprocess
from pathlib import Path

REPO_ROOT = Path("/data/hostel-booking-system")
MAPPING_FILE = REPO_ROOT / "migration-tools/naming-map.json"

def load_mapping():
    with open(MAPPING_FILE) as f:
        return json.load(f)

def find_old_references(pattern, paths):
    """Ищет упоминания старых имён в коде."""
    cmd = ["grep", "-r", pattern, "--include=*.go", "--include=*.ts", "--include=*.tsx"] + paths
    result = subprocess.run(cmd, capture_output=True, text=True)
    return result.stdout.splitlines()

def validate_database_tables(mapping):
    """Проверяет упоминания старых таблиц в SQL и Go коде."""
    print("🔍 Checking database table references...")

    old_tables = list(mapping["database_tables"].keys())
    issues = []

    for old_table in old_tables:
        refs = find_old_references(old_table, [
            str(REPO_ROOT / "backend/migrations"),
            str(REPO_ROOT / "backend/internal")
        ])
        if refs:
            issues.append({
                "category": "database",
                "old_name": old_table,
                "new_name": mapping["database_tables"][old_table],
                "references": len(refs),
                "files": list(set([r.split(":")[0] for r in refs]))
            })

    return issues

def validate_go_types(mapping):
    """Проверяет использование старых Go типов."""
    print("🔍 Checking Go type references...")

    old_types = list(mapping["go_types"].keys())
    issues = []

    for old_type in old_types:
        pattern = f"\\b{old_type}\\b"
        refs = find_old_references(pattern, [str(REPO_ROOT / "backend")])
        if refs:
            issues.append({
                "category": "go_types",
                "old_name": old_type,
                "new_name": mapping["go_types"][old_type],
                "references": len(refs)
            })

    return issues

def main():
    print("=" * 60)
    print("🧪 MIGRATION MAPPING VALIDATOR")
    print("=" * 60)

    mapping = load_mapping()

    all_issues = []
    all_issues.extend(validate_database_tables(mapping))
    all_issues.extend(validate_go_types(mapping))

    if not all_issues:
        print("✅ No old references found - mapping is complete!")
        return 0

    print(f"\n⚠️  Found {len(all_issues)} categories with old references:")
    for issue in all_issues:
        print(f"\n  📌 {issue['old_name']} → {issue['new_name']}")
        print(f"     References: {issue['references']}")
        if 'files' in issue:
            print(f"     Files affected: {len(issue['files'])}")

    return 1

if __name__ == "__main__":
    exit(main())
```

#### 1.3 Скрипт автоматического переименования

**Файл**: `/data/hostel-booking-system/migration-tools/auto-rename.sh`

```bash
#!/bin/bash
set -e

REPO_ROOT="/data/hostel-booking-system"
cd "$REPO_ROOT"

echo "🔄 Starting automatic renaming..."

# Backend: Go modules
echo "📦 Renaming Go modules..."
git mv backend/internal/proj/marketplace backend/internal/proj/c2c
git mv backend/internal/proj/storefronts backend/internal/proj/b2c

# Backend: Update imports
echo "📝 Updating Go imports..."
find backend -name "*.go" -type f -exec sed -i \
  's|internal/proj/marketplace|internal/proj/c2c|g' {} +
find backend -name "*.go" -type f -exec sed -i \
  's|internal/proj/storefronts|internal/proj/b2c|g' {} +

# Frontend: Components
echo "🎨 Renaming frontend components..."
git mv frontend/svetu/src/components/marketplace frontend/svetu/src/components/c2c || true
git mv frontend/svetu/src/components/storefronts frontend/svetu/src/components/b2c || true

# Frontend: Routes
echo "🛤️  Renaming frontend routes..."
git mv frontend/svetu/src/app/\[locale\]/marketplace frontend/svetu/src/app/\[locale\]/c2c || true
git mv frontend/svetu/src/app/\[locale\]/storefronts frontend/svetu/src/app/\[locale\]/b2c || true

# Frontend: i18n
echo "🌐 Renaming translation files..."
for lang in en ru sr; do
  git mv frontend/svetu/src/messages/$lang/marketplace.json \
         frontend/svetu/src/messages/$lang/c2c.json || true
  git mv frontend/svetu/src/messages/$lang/storefronts.json \
         frontend/svetu/src/messages/$lang/b2c.json || true
done

echo "✅ Automatic renaming complete!"
echo "⚠️  Run validation script to check for remaining references"
```

### Критерии приёмки
- [ ] ✅ `naming-map.json` создан и проверен
- [ ] ✅ Валидационный скрипт работает и проходит
- [ ] ✅ Auto-rename скрипт протестирован на test branch
- [ ] ✅ Документирован mapping для всех 8 категорий
- [ ] ✅ План обновлён: фактическое время, проблемы

### Риски
| Риск | Митигация |
|------|-----------|
| Неполный маппинг | Запустить валидацию 3 раза на разных этапах |
| Конфликты при git mv | Тестировать на копии репозитория |

### Фактические данные (заполнить)
- **Время**: ___ дней (план: 1.5-2 дня)
- **Найдено упоминаний**: ___
- **Проблемы**: ___

---

## 🗄️ ФАЗА 2: МИГРАЦИЯ БАЗЫ ДАННЫХ (Дни 3-7)

### Цель
Создать новые таблицы, мигрировать данные с сохранением целостности, обеспечить 100% покрытие тестами.

### Задачи

#### 2.1 Создание миграции схемы

**Файл**: `backend/migrations/000172_create_c2c_b2c_tables.up.sql`

```sql
-- ============================================================================
-- МИГРАЦИЯ: Создание C2C и B2C таблиц
-- Дата: 2025-10-09
-- Описание: Переименование marketplace → c2c, storefronts → b2c
-- Автор: Migration Plan v1.0
-- ============================================================================

BEGIN;

-- ============================================================================
-- C2C TABLES (бывшие marketplace_*)
-- ============================================================================

-- 1. C2C Categories
CREATE TABLE c2c_categories (
    LIKE marketplace_categories INCLUDING ALL
);

-- 2. C2C Listings
CREATE TABLE c2c_listings (
    LIKE marketplace_listings INCLUDING ALL
);

-- 3. C2C Images
CREATE TABLE c2c_images (
    LIKE marketplace_images INCLUDING ALL
);

-- 4. C2C Chats
CREATE TABLE c2c_chats (
    LIKE marketplace_chats INCLUDING ALL
);

-- 5. C2C Messages
CREATE TABLE c2c_messages (
    LIKE marketplace_messages INCLUDING ALL
);

-- 6. C2C Favorites
CREATE TABLE c2c_favorites (
    LIKE marketplace_favorites INCLUDING ALL
);

-- 7. C2C Orders
CREATE TABLE c2c_orders (
    LIKE marketplace_orders INCLUDING ALL
);

-- 8. C2C Listing Variants
CREATE TABLE c2c_listing_variants (
    LIKE marketplace_listing_variants INCLUDING ALL
);

-- ============================================================================
-- B2C TABLES (бывшие storefront_*)
-- ============================================================================

-- 1. B2C Stores (основная таблица магазинов)
CREATE TABLE b2c_stores (
    LIKE storefronts INCLUDING ALL
);

-- 2. B2C Products
CREATE TABLE b2c_products (
    LIKE storefront_products INCLUDING ALL
);

-- 3. B2C Product Images
CREATE TABLE b2c_product_images (
    LIKE storefront_product_images INCLUDING ALL
);

-- 4. B2C Product Variants
CREATE TABLE b2c_product_variants (
    LIKE storefront_product_variants INCLUDING ALL
);

-- 5. B2C Product Attributes
CREATE TABLE b2c_product_attributes (
    LIKE storefront_product_attributes INCLUDING ALL
);

-- 6. B2C Orders
CREATE TABLE b2c_orders (
    LIKE storefront_orders INCLUDING ALL
);

-- 7. B2C Order Items
CREATE TABLE b2c_order_items (
    LIKE storefront_order_items INCLUDING ALL
);

-- 8. B2C Favorites
CREATE TABLE b2c_favorites (
    LIKE storefront_favorites INCLUDING ALL
);

-- 9. B2C Store Hours
CREATE TABLE b2c_store_hours (
    LIKE storefront_hours INCLUDING ALL
);

-- 10. B2C Store Staff
CREATE TABLE b2c_store_staff (
    LIKE storefront_staff INCLUDING ALL
);

-- 11. B2C Payment Methods
CREATE TABLE b2c_payment_methods (
    LIKE storefront_payment_methods INCLUDING ALL
);

-- 12. B2C Delivery Options
CREATE TABLE b2c_delivery_options (
    LIKE storefront_delivery_options INCLUDING ALL
);

-- 13. B2C Inventory Movements
CREATE TABLE b2c_inventory_movements (
    LIKE storefront_inventory_movements INCLUDING ALL
);

-- 14. User B2C Stores (связь пользователей с магазинами)
CREATE TABLE user_b2c_stores (
    LIKE user_storefronts INCLUDING ALL
);

-- 15. B2C Product Variant Images
CREATE TABLE b2c_product_variant_images (
    LIKE storefront_product_variant_images INCLUDING ALL
);

COMMIT;

-- ============================================================================
-- ВАЖНО: Индексы и constraints скопированы через INCLUDING ALL
-- Следующий шаг: миграция данных (отдельная миграция)
-- ============================================================================
```

**Файл**: `backend/migrations/000172_create_c2c_b2c_tables.down.sql`

```sql
BEGIN;

-- Drop C2C tables
DROP TABLE IF EXISTS c2c_listing_variants CASCADE;
DROP TABLE IF EXISTS c2c_orders CASCADE;
DROP TABLE IF EXISTS c2c_favorites CASCADE;
DROP TABLE IF EXISTS c2c_messages CASCADE;
DROP TABLE IF EXISTS c2c_chats CASCADE;
DROP TABLE IF EXISTS c2c_images CASCADE;
DROP TABLE IF EXISTS c2c_listings CASCADE;
DROP TABLE IF EXISTS c2c_categories CASCADE;

-- Drop B2C tables
DROP TABLE IF EXISTS b2c_product_variant_images CASCADE;
DROP TABLE IF EXISTS user_b2c_stores CASCADE;
DROP TABLE IF EXISTS b2c_inventory_movements CASCADE;
DROP TABLE IF EXISTS b2c_delivery_options CASCADE;
DROP TABLE IF EXISTS b2c_payment_methods CASCADE;
DROP TABLE IF EXISTS b2c_store_staff CASCADE;
DROP TABLE IF EXISTS b2c_store_hours CASCADE;
DROP TABLE IF EXISTS b2c_favorites CASCADE;
DROP TABLE IF EXISTS b2c_order_items CASCADE;
DROP TABLE IF EXISTS b2c_orders CASCADE;
DROP TABLE IF EXISTS b2c_product_attributes CASCADE;
DROP TABLE IF EXISTS b2c_product_variants CASCADE;
DROP TABLE IF EXISTS b2c_product_images CASCADE;
DROP TABLE IF EXISTS b2c_products CASCADE;
DROP TABLE IF EXISTS b2c_stores CASCADE;

COMMIT;
```

#### 2.2 Миграция данных

**Файл**: `backend/migrations/000173_migrate_c2c_b2c_data.up.sql`

```sql
-- ============================================================================
-- МИГРАЦИЯ ДАННЫХ: Копирование из marketplace/storefront в c2c/b2c
-- КРИТИЧНО: Сохраняем ID для связей!
-- ============================================================================

BEGIN;

-- ============================================================================
-- C2C DATA MIGRATION
-- ============================================================================

-- 1. Categories (сначала - для FK)
INSERT INTO c2c_categories SELECT * FROM marketplace_categories;

-- 2. Listings (зависит от categories)
INSERT INTO c2c_listings SELECT * FROM marketplace_listings;

-- 3. Images (зависит от listings)
INSERT INTO c2c_images SELECT * FROM marketplace_images;

-- 4. Listing Variants (зависит от listings)
INSERT INTO c2c_listing_variants SELECT * FROM marketplace_listing_variants;

-- 5. Chats (зависит от listings)
INSERT INTO c2c_chats SELECT * FROM marketplace_chats;

-- 6. Messages (зависит от chats)
INSERT INTO c2c_messages SELECT * FROM marketplace_messages;

-- 7. Favorites (зависит от listings)
INSERT INTO c2c_favorites SELECT * FROM marketplace_favorites;

-- 8. Orders (зависит от listings)
INSERT INTO c2c_orders SELECT * FROM marketplace_orders;

-- ============================================================================
-- B2C DATA MIGRATION
-- ============================================================================

-- 1. Stores (основная таблица - сначала)
INSERT INTO b2c_stores SELECT * FROM storefronts;

-- 2. User-Store links (зависит от stores)
INSERT INTO user_b2c_stores SELECT * FROM user_storefronts;

-- 3. Products (зависит от stores)
INSERT INTO b2c_products SELECT * FROM storefront_products;

-- 4. Product Images (зависит от products)
INSERT INTO b2c_product_images SELECT * FROM storefront_product_images;

-- 5. Product Variants (зависит от products)
INSERT INTO b2c_product_variants SELECT * FROM storefront_product_variants;

-- 6. Product Variant Images (зависит от variants)
INSERT INTO b2c_product_variant_images SELECT * FROM storefront_product_variant_images;

-- 7. Product Attributes (зависит от products)
INSERT INTO b2c_product_attributes SELECT * FROM storefront_product_attributes;

-- 8. Orders (зависит от stores)
INSERT INTO b2c_orders SELECT * FROM storefront_orders;

-- 9. Order Items (зависит от orders и products)
INSERT INTO b2c_order_items SELECT * FROM storefront_order_items;

-- 10. Favorites (зависит от products)
INSERT INTO b2c_favorites SELECT * FROM storefront_favorites;

-- 11. Store Hours (зависит от stores)
INSERT INTO b2c_store_hours SELECT * FROM storefront_hours;

-- 12. Store Staff (зависит от stores)
INSERT INTO b2c_store_staff SELECT * FROM storefront_staff;

-- 13. Payment Methods (зависит от stores)
INSERT INTO b2c_payment_methods SELECT * FROM storefront_payment_methods;

-- 14. Delivery Options (зависит от stores)
INSERT INTO b2c_delivery_options SELECT * FROM storefront_delivery_options;

-- 15. Inventory Movements (зависит от variants)
INSERT INTO b2c_inventory_movements SELECT * FROM storefront_inventory_movements;

COMMIT;

-- ============================================================================
-- Проверка целостности данных
-- ============================================================================

DO $$
DECLARE
    c2c_count INTEGER;
    b2c_count INTEGER;
    old_c2c_count INTEGER;
    old_b2c_count INTEGER;
BEGIN
    -- Проверка C2C listings
    SELECT COUNT(*) INTO c2c_count FROM c2c_listings;
    SELECT COUNT(*) INTO old_c2c_count FROM marketplace_listings;

    IF c2c_count != old_c2c_count THEN
        RAISE EXCEPTION 'C2C listings count mismatch! Expected %, got %', old_c2c_count, c2c_count;
    END IF;

    -- Проверка B2C products
    SELECT COUNT(*) INTO b2c_count FROM b2c_products;
    SELECT COUNT(*) INTO old_b2c_count FROM storefront_products;

    IF b2c_count != old_b2c_count THEN
        RAISE EXCEPTION 'B2C products count mismatch! Expected %, got %', old_b2c_count, b2c_count;
    END IF;

    RAISE NOTICE '✅ Data migration validated: C2C=%, B2C=%', c2c_count, b2c_count;
END $$;
```

**Файл**: `backend/migrations/000173_migrate_c2c_b2c_data.down.sql`

```sql
BEGIN;

-- Очистить все данные из новых таблиц (обратная миграция)
TRUNCATE TABLE c2c_messages CASCADE;
TRUNCATE TABLE c2c_chats CASCADE;
TRUNCATE TABLE c2c_favorites CASCADE;
TRUNCATE TABLE c2c_orders CASCADE;
TRUNCATE TABLE c2c_listing_variants CASCADE;
TRUNCATE TABLE c2c_images CASCADE;
TRUNCATE TABLE c2c_listings CASCADE;
TRUNCATE TABLE c2c_categories CASCADE;

TRUNCATE TABLE b2c_inventory_movements CASCADE;
TRUNCATE TABLE b2c_delivery_options CASCADE;
TRUNCATE TABLE b2c_payment_methods CASCADE;
TRUNCATE TABLE b2c_store_staff CASCADE;
TRUNCATE TABLE b2c_store_hours CASCADE;
TRUNCATE TABLE b2c_favorites CASCADE;
TRUNCATE TABLE b2c_order_items CASCADE;
TRUNCATE TABLE b2c_orders CASCADE;
TRUNCATE TABLE b2c_product_attributes CASCADE;
TRUNCATE TABLE b2c_product_variant_images CASCADE;
TRUNCATE TABLE b2c_product_variants CASCADE;
TRUNCATE TABLE b2c_product_images CASCADE;
TRUNCATE TABLE b2c_products CASCADE;
TRUNCATE TABLE user_b2c_stores CASCADE;
TRUNCATE TABLE b2c_stores CASCADE;

COMMIT;
```

#### 2.3 Обновление Foreign Keys

**Файл**: `backend/migrations/000174_update_c2c_b2c_constraints.up.sql`

```sql
-- ============================================================================
-- ОБНОВЛЕНИЕ CONSTRAINTS: Переименование FK на новые таблицы
-- ============================================================================

BEGIN;

-- ============================================================================
-- C2C Foreign Keys
-- ============================================================================

-- c2c_listings
ALTER TABLE c2c_listings
  DROP CONSTRAINT IF EXISTS marketplace_listings_category_id_fkey CASCADE,
  ADD CONSTRAINT c2c_listings_category_id_fkey
    FOREIGN KEY (category_id) REFERENCES c2c_categories(id) ON DELETE CASCADE;

ALTER TABLE c2c_listings
  DROP CONSTRAINT IF EXISTS marketplace_listings_user_id_fkey CASCADE,
  ADD CONSTRAINT c2c_listings_user_id_fkey
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- c2c_images
ALTER TABLE c2c_images
  DROP CONSTRAINT IF EXISTS marketplace_images_listing_id_fkey CASCADE,
  ADD CONSTRAINT c2c_images_listing_id_fkey
    FOREIGN KEY (listing_id) REFERENCES c2c_listings(id) ON DELETE CASCADE;

-- c2c_chats
ALTER TABLE c2c_chats
  DROP CONSTRAINT IF EXISTS marketplace_chats_listing_id_fkey CASCADE,
  ADD CONSTRAINT c2c_chats_listing_id_fkey
    FOREIGN KEY (listing_id) REFERENCES c2c_listings(id) ON DELETE CASCADE;

-- c2c_messages
ALTER TABLE c2c_messages
  DROP CONSTRAINT IF EXISTS marketplace_messages_chat_id_fkey CASCADE,
  ADD CONSTRAINT c2c_messages_chat_id_fkey
    FOREIGN KEY (chat_id) REFERENCES c2c_chats(id) ON DELETE CASCADE;

-- c2c_favorites
ALTER TABLE c2c_favorites
  DROP CONSTRAINT IF EXISTS marketplace_favorites_listing_id_fkey CASCADE,
  ADD CONSTRAINT c2c_favorites_listing_id_fkey
    FOREIGN KEY (listing_id) REFERENCES c2c_listings(id) ON DELETE CASCADE;

-- c2c_orders
ALTER TABLE c2c_orders
  DROP CONSTRAINT IF EXISTS marketplace_orders_listing_id_fkey CASCADE,
  ADD CONSTRAINT c2c_orders_listing_id_fkey
    FOREIGN KEY (listing_id) REFERENCES c2c_listings(id) ON DELETE CASCADE;

-- c2c_listing_variants
ALTER TABLE c2c_listing_variants
  DROP CONSTRAINT IF EXISTS marketplace_listing_variants_listing_id_fkey CASCADE,
  ADD CONSTRAINT c2c_listing_variants_listing_id_fkey
    FOREIGN KEY (listing_id) REFERENCES c2c_listings(id) ON DELETE CASCADE;

-- ============================================================================
-- B2C Foreign Keys
-- ============================================================================

-- user_b2c_stores
ALTER TABLE user_b2c_stores
  DROP CONSTRAINT IF EXISTS user_storefronts_storefront_id_fkey CASCADE,
  ADD CONSTRAINT user_b2c_stores_store_id_fkey
    FOREIGN KEY (storefront_id) REFERENCES b2c_stores(id) ON DELETE CASCADE;

ALTER TABLE user_b2c_stores
  DROP CONSTRAINT IF EXISTS user_storefronts_user_id_fkey CASCADE,
  ADD CONSTRAINT user_b2c_stores_user_id_fkey
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- b2c_products
ALTER TABLE b2c_products
  DROP CONSTRAINT IF EXISTS storefront_products_storefront_id_fkey CASCADE,
  ADD CONSTRAINT b2c_products_store_id_fkey
    FOREIGN KEY (storefront_id) REFERENCES b2c_stores(id) ON DELETE CASCADE;

-- b2c_product_images
ALTER TABLE b2c_product_images
  DROP CONSTRAINT IF EXISTS storefront_product_images_product_id_fkey CASCADE,
  ADD CONSTRAINT b2c_product_images_product_id_fkey
    FOREIGN KEY (product_id) REFERENCES b2c_products(id) ON DELETE CASCADE;

-- b2c_product_variants
ALTER TABLE b2c_product_variants
  DROP CONSTRAINT IF EXISTS storefront_product_variants_product_id_fkey CASCADE,
  ADD CONSTRAINT b2c_product_variants_product_id_fkey
    FOREIGN KEY (product_id) REFERENCES b2c_products(id) ON DELETE CASCADE;

-- b2c_product_variant_images
ALTER TABLE b2c_product_variant_images
  DROP CONSTRAINT IF EXISTS storefront_product_variant_images_variant_id_fkey CASCADE,
  ADD CONSTRAINT b2c_product_variant_images_variant_id_fkey
    FOREIGN KEY (variant_id) REFERENCES b2c_product_variants(id) ON DELETE CASCADE;

-- b2c_product_attributes
ALTER TABLE b2c_product_attributes
  DROP CONSTRAINT IF EXISTS storefront_product_attributes_product_id_fkey CASCADE,
  ADD CONSTRAINT b2c_product_attributes_product_id_fkey
    FOREIGN KEY (product_id) REFERENCES b2c_products(id) ON DELETE CASCADE;

-- b2c_orders
ALTER TABLE b2c_orders
  DROP CONSTRAINT IF EXISTS storefront_orders_storefront_id_fkey CASCADE,
  ADD CONSTRAINT b2c_orders_store_id_fkey
    FOREIGN KEY (storefront_id) REFERENCES b2c_stores(id) ON DELETE CASCADE;

-- b2c_order_items
ALTER TABLE b2c_order_items
  DROP CONSTRAINT IF EXISTS storefront_order_items_order_id_fkey CASCADE,
  ADD CONSTRAINT b2c_order_items_order_id_fkey
    FOREIGN KEY (order_id) REFERENCES b2c_orders(id) ON DELETE CASCADE;

ALTER TABLE b2c_order_items
  DROP CONSTRAINT IF EXISTS storefront_order_items_product_id_fkey CASCADE,
  ADD CONSTRAINT b2c_order_items_product_id_fkey
    FOREIGN KEY (product_id) REFERENCES b2c_products(id) ON DELETE SET NULL;

-- b2c_favorites
ALTER TABLE b2c_favorites
  DROP CONSTRAINT IF EXISTS storefront_favorites_product_id_fkey CASCADE,
  ADD CONSTRAINT b2c_favorites_product_id_fkey
    FOREIGN KEY (product_id) REFERENCES b2c_products(id) ON DELETE CASCADE;

-- b2c_store_hours
ALTER TABLE b2c_store_hours
  DROP CONSTRAINT IF EXISTS storefront_hours_storefront_id_fkey CASCADE,
  ADD CONSTRAINT b2c_store_hours_store_id_fkey
    FOREIGN KEY (storefront_id) REFERENCES b2c_stores(id) ON DELETE CASCADE;

-- b2c_store_staff
ALTER TABLE b2c_store_staff
  DROP CONSTRAINT IF EXISTS storefront_staff_storefront_id_fkey CASCADE,
  ADD CONSTRAINT b2c_store_staff_store_id_fkey
    FOREIGN KEY (storefront_id) REFERENCES b2c_stores(id) ON DELETE CASCADE;

-- b2c_payment_methods
ALTER TABLE b2c_payment_methods
  DROP CONSTRAINT IF EXISTS storefront_payment_methods_storefront_id_fkey CASCADE,
  ADD CONSTRAINT b2c_payment_methods_store_id_fkey
    FOREIGN KEY (storefront_id) REFERENCES b2c_stores(id) ON DELETE CASCADE;

-- b2c_delivery_options
ALTER TABLE b2c_delivery_options
  DROP CONSTRAINT IF EXISTS storefront_delivery_options_storefront_id_fkey CASCADE,
  ADD CONSTRAINT b2c_delivery_options_store_id_fkey
    FOREIGN KEY (storefront_id) REFERENCES b2c_stores(id) ON DELETE CASCADE;

-- b2c_inventory_movements
ALTER TABLE b2c_inventory_movements
  DROP CONSTRAINT IF EXISTS storefront_inventory_movements_variant_id_fkey CASCADE,
  ADD CONSTRAINT b2c_inventory_movements_variant_id_fkey
    FOREIGN KEY (variant_id) REFERENCES b2c_product_variants(id) ON DELETE CASCADE;

COMMIT;

-- ============================================================================
-- Верификация
-- ============================================================================

DO $$
DECLARE
    fk_count INTEGER;
BEGIN
    -- Подсчитать все FK на новые таблицы
    SELECT COUNT(*) INTO fk_count
    FROM information_schema.table_constraints
    WHERE constraint_type = 'FOREIGN KEY'
      AND (table_name LIKE 'c2c_%' OR table_name LIKE 'b2c_%');

    RAISE NOTICE '✅ Created % foreign keys for C2C/B2C tables', fk_count;
END $$;
```

### Критерии приёмки
- [ ] ✅ Все 3 миграции (создание, данные, FK) работают
- [ ] ✅ Up/Down миграции протестированы 3 раза
- [ ] ✅ Данные мигрированы на 100% (проверка COUNT)
- [ ] ✅ Все FK ссылаются на новые таблицы
- [ ] ✅ Нет orphaned records (проверка целостности)
- [ ] ✅ Триггеры пересозданы для новых таблиц

### Риски
| Риск | Митигация |
|------|-----------|
| Потеря данных | Тройная проверка COUNT перед/после |
| Сломанные FK | Автоматическая валидация в миграции |
| Время выполнения | Тестировать на копии БД, оптимизировать индексы |

### Фактические данные
- **Время**: ___ дней (план: 4-5 дней)
- **Мигрировано строк**: C2C=___, B2C=___
- **Проблемы**: ___

---

## 🔧 ФАЗА 3: BACKEND МИГРАЦИЯ (Дни 8-14)

### Цель
Переименовать модули, обновить все Go типы, SQL запросы, API endpoints. 100% компиляция + тесты.

### Задачи

#### 3.1 Переименование модулей проекта

```bash
#!/bin/bash
# Файл: /data/hostel-booking-system/migration-tools/backend-rename-modules.sh

set -e
cd /data/hostel-booking-system/backend

echo "🔄 Renaming backend modules..."

# 1. Переименовать директории
git mv internal/proj/marketplace internal/proj/c2c
git mv internal/proj/storefronts internal/proj/b2c

# 2. Обновить импорты во всех Go файлах
echo "📝 Updating imports in all Go files..."

find . -name "*.go" -type f -exec sed -i \
  's|backend/internal/proj/marketplace|backend/internal/proj/c2c|g' {} +

find . -name "*.go" -type f -exec sed -i \
  's|backend/internal/proj/storefronts|backend/internal/proj/b2c|g' {} +

# 3. Обновить domain models
echo "📦 Renaming domain model files..."

# C2C models
cd internal/domain
git mv marketplace_listing.go c2c_listing.go || true
git mv marketplace_category.go c2c_category.go || true
git mv marketplace_image.go c2c_image.go || true
git mv marketplace_chat.go c2c_chat.go || true
git mv marketplace_message.go c2c_message.go || true
git mv marketplace_favorite.go c2c_favorite.go || true
git mv marketplace_order.go c2c_order.go || true

# B2C models
git mv storefront.go b2c_store.go || true
git mv storefront_product.go b2c_product.go || true
git mv storefront_product_image.go b2c_product_image.go || true
git mv storefront_product_variant.go b2c_product_variant.go || true
git mv storefront_order.go b2c_order.go || true

cd ../..

echo "✅ Module renaming complete!"
echo "⚠️  Next: run go build to verify compilation"
```

#### 3.2 Обновление Go типов

**Скрипт**: `/data/hostel-booking-system/migration-tools/backend-rename-types.sh`

```bash
#!/bin/bash
set -e

cd /data/hostel-booking-system/backend

echo "🔤 Renaming Go types..."

# C2C types
find . -name "*.go" -exec sed -i 's/\bMarketplaceListing\b/C2CListing/g' {} +
find . -name "*.go" -exec sed -i 's/\bMarketplaceCategory\b/C2CCategory/g' {} +
find . -name "*.go" -exec sed -i 's/\bMarketplaceImage\b/C2CImage/g' {} +
find . -name "*.go" -exec sed -i 's/\bMarketplaceChat\b/C2CChat/g' {} +
find . -name "*.go" -exec sed -i 's/\bMarketplaceMessage\b/C2CMessage/g' {} +
find . -name "*.go" -exec sed -i 's/\bMarketplaceFavorite\b/C2CFavorite/g' {} +
find . -name "*.go" -exec sed -i 's/\bMarketplaceOrder\b/C2COrder/g' {} +

# B2C types
find . -name "*.go" -exec sed -i 's/\bStorefront\b/B2CStore/g' {} +
find . -name "*.go" -exec sed -i 's/\bStorefrontProduct\b/B2CProduct/g' {} +
find . -name "*.go" -exec sed -i 's/\bStorefrontProductImage\b/B2CProductImage/g' {} +
find . -name "*.go" -exec sed -i 's/\bStorefrontProductVariant\b/B2CProductVariant/g' {} +
find . -name "*.go" -exec sed -i 's/\bStorefrontOrder\b/B2COrder/g' {} +

echo "✅ Type renaming complete!"
```

#### 3.3 Обновление SQL запросов в коде

```bash
#!/bin/bash
# Файл: /data/hostel-booking-system/migration-tools/backend-rename-sql.sh

set -e
cd /data/hostel-booking-system/backend

echo "🗄️  Updating SQL table references..."

# C2C tables
find internal/proj/c2c -name "*.go" -exec sed -i \
  's/marketplace_listings/c2c_listings/g' {} +
find internal/proj/c2c -name "*.go" -exec sed -i \
  's/marketplace_categories/c2c_categories/g' {} +
find internal/proj/c2c -name "*.go" -exec sed -i \
  's/marketplace_images/c2c_images/g' {} +
find internal/proj/c2c -name "*.go" -exec sed -i \
  's/marketplace_chats/c2c_chats/g' {} +
find internal/proj/c2c -name "*.go" -exec sed -i \
  's/marketplace_messages/c2c_messages/g' {} +
find internal/proj/c2c -name "*.go" -exec sed -i \
  's/marketplace_favorites/c2c_favorites/g' {} +
find internal/proj/c2c -name "*.go" -exec sed -i \
  's/marketplace_orders/c2c_orders/g' {} +
find internal/proj/c2c -name "*.go" -exec sed -i \
  's/marketplace_listing_variants/c2c_listing_variants/g' {} +

# B2C tables
find internal/proj/b2c -name "*.go" -exec sed -i \
  's/storefronts\b/b2c_stores/g' {} +
find internal/proj/b2c -name "*.go" -exec sed -i \
  's/storefront_products/b2c_products/g' {} +
find internal/proj/b2c -name "*.go" -exec sed -i \
  's/storefront_product_images/b2c_product_images/g' {} +
find internal/proj/b2c -name "*.go" -exec sed -i \
  's/storefront_product_variants/b2c_product_variants/g' {} +
find internal/proj/b2c -name "*.go" -exec sed -i \
  's/storefront_product_attributes/b2c_product_attributes/g' {} +
find internal/proj/b2c -name "*.go" -exec sed -i \
  's/storefront_orders/b2c_orders/g' {} +
find internal/proj/b2c -name "*.go" -exec sed -i \
  's/storefront_order_items/b2c_order_items/g' {} +
find internal/proj/b2c -name "*.go" -exec sed -i \
  's/storefront_favorites/b2c_favorites/g' {} +
find internal/proj/b2c -name "*.go" -exec sed -i \
  's/storefront_hours/b2c_store_hours/g' {} +
find internal/proj/b2c -name "*.go" -exec sed -i \
  's/storefront_staff/b2c_store_staff/g' {} +
find internal/proj/b2c -name "*.go" -exec sed -i \
  's/storefront_payment_methods/b2c_payment_methods/g' {} +
find internal/proj/b2c -name "*.go" -exec sed -i \
  's/storefront_delivery_options/b2c_delivery_options/g' {} +
find internal/proj/b2c -name "*.go" -exec sed -i \
  's/storefront_inventory_movements/b2c_inventory_movements/g' {} +
find internal/proj/b2c -name "*.go" -exec sed -i \
  's/user_storefronts/user_b2c_stores/g' {} +
find internal/proj/b2c -name "*.go" -exec sed -i \
  's/storefront_product_variant_images/b2c_product_variant_images/g' {} +

echo "✅ SQL references updated!"
```

#### 3.4 Обновление API endpoints (routes)

**ВАЖНО**: НЕТ обратной совместимости - удаляем старые роуты полностью!

Пример для `backend/internal/proj/c2c/handler/routes.go`:

```go
package handler

import (
	"github.com/gofiber/fiber/v2"
	authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"
	"backend/internal/middleware"
)

// RegisterRoutes регистрирует C2C (бывшие marketplace) эндпоинты
func (h *Handler) RegisterRoutes(app *fiber.App, mw *middleware.Middleware) error {
	// ============================================================================
	// C2C PUBLIC ENDPOINTS
	// ============================================================================
	c2cPublic := app.Group("/api/v1/c2c")

	// Listings
	c2cPublic.Get("/listings", h.GetListings)                // GET /api/v1/c2c/listings
	c2cPublic.Get("/listings/:id", h.GetListingByID)         // GET /api/v1/c2c/listings/:id
	c2cPublic.Get("/listings/user/:userId", h.GetUserListings)

	// Categories
	c2cPublic.Get("/categories", h.GetCategories)            // GET /api/v1/c2c/categories
	c2cPublic.Get("/categories/:id", h.GetCategoryByID)

	// Search
	c2cPublic.Get("/search", h.SearchListings)               // GET /api/v1/c2c/search

	// ============================================================================
	// C2C AUTHENTICATED ENDPOINTS
	// ============================================================================
	c2cAuth := app.Group("/api/v1/c2c", h.jwtParserMW, authMiddleware.RequireAuth())

	// Create/Update/Delete
	c2cAuth.Post("/listings", h.CreateListing)               // POST /api/v1/c2c/listings
	c2cAuth.Put("/listings/:id", h.UpdateListing)            // PUT /api/v1/c2c/listings/:id
	c2cAuth.Delete("/listings/:id", h.DeleteListing)         // DELETE /api/v1/c2c/listings/:id

	// Images
	c2cAuth.Post("/listings/:id/images", h.UploadImages)     // POST /api/v1/c2c/listings/:id/images
	c2cAuth.Delete("/images/:id", h.DeleteImage)

	// Favorites
	c2cAuth.Get("/favorites", h.GetFavorites)                // GET /api/v1/c2c/favorites
	c2cAuth.Post("/favorites", h.AddFavorite)
	c2cAuth.Delete("/favorites/:id", h.RemoveFavorite)

	// Chats
	c2cAuth.Get("/chats", h.GetChats)                        // GET /api/v1/c2c/chats
	c2cAuth.Post("/chats", h.CreateChat)
	c2cAuth.Get("/chats/:id/messages", h.GetMessages)
	c2cAuth.Post("/chats/:id/messages", h.SendMessage)

	// ============================================================================
	// C2C ADMIN ENDPOINTS
	// ============================================================================
	c2cAdmin := app.Group("/api/v1/admin/c2c", h.jwtParserMW, authMiddleware.RequireAuthString("admin"))

	c2cAdmin.Get("/listings", h.GetAllListingsAdmin)         // GET /api/v1/admin/c2c/listings
	c2cAdmin.Put("/listings/:id/status", h.UpdateListingStatus)
	c2cAdmin.Delete("/listings/:id", h.DeleteListingAdmin)

	c2cAdmin.Get("/categories", h.GetAllCategoriesAdmin)
	c2cAdmin.Post("/categories", h.CreateCategory)
	c2cAdmin.Put("/categories/:id", h.UpdateCategory)
	c2cAdmin.Delete("/categories/:id", h.DeleteCategory)

	return nil
}

func (h *Handler) GetPrefix() string {
	return "/api/v1/c2c"  // обновлено с "/api/v1/marketplace"
}
```

Аналогично для B2C (`backend/internal/proj/b2c/handler/routes.go`):

```go
// RegisterRoutes регистрирует B2C (бывшие storefronts) эндпоинты
func (h *Handler) RegisterRoutes(app *fiber.App, mw *middleware.Middleware) error {
	// ============================================================================
	// B2C PUBLIC ENDPOINTS
	// ============================================================================
	b2cPublic := app.Group("/api/v1/b2c")

	// Stores
	b2cPublic.Get("/stores", h.GetStores)                    // GET /api/v1/b2c/stores
	b2cPublic.Get("/stores/:id", h.GetStoreByID)             // GET /api/v1/b2c/stores/:id
	b2cPublic.Get("/stores/:id/products", h.GetStoreProducts)

	// Products
	b2cPublic.Get("/products", h.GetProducts)                // GET /api/v1/b2c/products
	b2cPublic.Get("/products/:id", h.GetProductByID)         // GET /api/v1/b2c/products/:id

	// Search
	b2cPublic.Get("/search", h.SearchProducts)               // GET /api/v1/b2c/search

	// ============================================================================
	// B2C AUTHENTICATED ENDPOINTS
	// ============================================================================
	b2cAuth := app.Group("/api/v1/b2c", h.jwtParserMW, authMiddleware.RequireAuth())

	// Store management
	b2cAuth.Post("/stores", h.CreateStore)                   // POST /api/v1/b2c/stores
	b2cAuth.Put("/stores/:id", h.UpdateStore)                // PUT /api/v1/b2c/stores/:id
	b2cAuth.Delete("/stores/:id", h.DeleteStore)

	// Product management
	b2cAuth.Post("/products", h.CreateProduct)               // POST /api/v1/b2c/products
	b2cAuth.Put("/products/:id", h.UpdateProduct)
	b2cAuth.Delete("/products/:id", h.DeleteProduct)

	// Product images
	b2cAuth.Post("/products/:id/images", h.UploadProductImages)
	b2cAuth.Delete("/images/:id", h.DeleteProductImage)

	// Variants
	b2cAuth.Post("/products/:id/variants", h.CreateVariant)
	b2cAuth.Put("/variants/:id", h.UpdateVariant)
	b2cAuth.Delete("/variants/:id", h.DeleteVariant)

	// ============================================================================
	// B2C ADMIN ENDPOINTS
	// ============================================================================
	b2cAdmin := app.Group("/api/v1/admin/b2c", h.jwtParserMW, authMiddleware.RequireAuthString("admin"))

	b2cAdmin.Get("/stores", h.GetAllStoresAdmin)             // GET /api/v1/admin/b2c/stores
	b2cAdmin.Put("/stores/:id/status", h.UpdateStoreStatus)
	b2cAdmin.Delete("/stores/:id", h.DeleteStoreAdmin)

	b2cAdmin.Get("/products", h.GetAllProductsAdmin)
	b2cAdmin.Put("/products/:id/status", h.UpdateProductStatus)

	return nil
}

func (h *Handler) GetPrefix() string {
	return "/api/v1/b2c"  // обновлено с "/api/v1/storefronts"
}
```

#### 3.5 Обновление Swagger аннотаций

Пример для C2C listing handler:

```go
// GetListings godoc
// @Summary      Get C2C listings
// @Description  Get paginated list of C2C (customer-to-customer) listings
// @Tags         c2c
// @Accept       json
// @Produce      json
// @Param        page     query    int     false  "Page number"
// @Param        limit    query    int     false  "Items per page"
// @Param        category query    string  false  "Category ID"
// @Success      200      {object} utils.SuccessResponse{data=[]domain.C2CListing}
// @Failure      400      {object} utils.ErrorResponse
// @Failure      500      {object} utils.ErrorResponse
// @Router       /api/v1/c2c/listings [get]
func (h *Handler) GetListings(c *fiber.Ctx) error {
	// ...
}
```

**После обновления всех аннотаций:**

```bash
cd /data/hostel-booking-system/backend
make swagger  # регенерировать swagger.json

# Проверить результат
python3 -m http.server 8888 -d docs &
curl http://localhost:8888/swagger.json | jq '.paths | keys' | grep -E "c2c|b2c"
pkill -f "python3 -m http.server 8888"
```

### Критерии приёмки
- [ ] ✅ Все модули переименованы (c2c, b2c)
- [ ] ✅ `go build ./...` проходит без ошибок
- [ ] ✅ Все импорты обновлены
- [ ] ✅ Все API endpoints обновлены
- [ ] ✅ Swagger аннотации обновлены + regenerated
- [ ] ✅ `make lint` проходит (0 warnings)
- [ ] ✅ `go test ./...` проходит (100% существующих тестов)
- [ ] ✅ Нет упоминаний "marketplace" или "storefront" в backend/internal/proj

### Риски
| Риск | Митигация |
|------|-----------|
| Сломанная компиляция | Поэтапное переименование + частые проверки |
| Пропущенные импорты | Автоматический скрипт + go build |
| Конфликты в git mv | Делать в изолированной ветке |

### Фактические данные
- **Время**: ___ дней (план: 5-7 дней)
- **Файлов обновлено**: ___
- **Проблемы**: ___

---

## 🎨 ФАЗА 4: FRONTEND МИГРАЦИЯ (Дни 15-19)

### Цель
Переименовать компоненты, типы, роуты, переводы. 100% компиляция TypeScript + сборка Next.js.

### Задачи

#### 4.1 Переименование директорий

```bash
#!/bin/bash
# Файл: /data/hostel-booking-system/migration-tools/frontend-rename-dirs.sh

set -e
cd /data/hostel-booking-system/frontend/svetu

echo "📁 Renaming frontend directories..."

# App routes
git mv src/app/\[locale\]/marketplace src/app/\[locale\]/c2c || true
git mv src/app/\[locale\]/storefronts src/app/\[locale\]/b2c || true

# Components
git mv src/components/marketplace src/components/c2c || true
git mv src/components/storefronts src/components/b2c || true

# Services
git mv src/services/marketplaceApi.ts src/services/c2cApi.ts || true
git mv src/services/storefrontsApi.ts src/services/b2cApi.ts || true

# Types
git mv src/types/marketplace.ts src/types/c2c.ts || true
git mv src/types/storefront.ts src/types/b2c.ts || true

# Store (Redux)
git mv src/store/slices/marketplaceSlice.ts src/store/slices/c2cSlice.ts || true
git mv src/store/slices/storefrontsSlice.ts src/store/slices/b2cSlice.ts || true

echo "✅ Directory renaming complete!"
```

#### 4.2 Обновление TypeScript типов

```bash
#!/bin/bash
# Файл: /data/hostel-booking-system/migration-tools/frontend-rename-types.sh

set -e
cd /data/hostel-booking-system/frontend/svetu

echo "🔤 Updating TypeScript types..."

# C2C types
find src -name "*.ts" -o -name "*.tsx" | xargs sed -i \
  's/\bMarketplaceListing\b/C2CListing/g'
find src -name "*.ts" -o -name "*.tsx" | xargs sed -i \
  's/\bMarketplaceCategory\b/C2CCategory/g'
find src -name "*.ts" -o -name "*.tsx" | xargs sed -i \
  's/\bMarketplaceImage\b/C2CImage/g'

# B2C types
find src -name "*.ts" -o -name "*.tsx" | xargs sed -i \
  's/\bStorefront\b/B2CStore/g'
find src -name "*.ts" -o -name "*.tsx" | xargs sed -i \
  's/\bStorefrontProduct\b/B2CProduct/g'
find src -name "*.ts" -o -name "*.tsx" | xargs sed -i \
  's/\bStorefrontProductImage\b/B2CProductImage/g'

echo "✅ Type renaming complete!"
```

#### 4.3 Обновление API клиентов

Пример `src/services/c2cApi.ts` (бывший marketplaceApi.ts):

```typescript
import { apiClient } from '@/services/api-client';
import type { C2CListing, C2CCategory, C2CImage } from '@/types/c2c';

export class C2CApi {
  // ============================================================================
  // LISTINGS
  // ============================================================================

  async getListings(params?: {
    page?: number;
    limit?: number;
    categoryId?: string;
  }): Promise<{ listings: C2CListing[]; total: number }> {
    const response = await apiClient.get('/c2c/listings', { params });
    return response.data;
  }

  async getListingById(id: string): Promise<C2CListing> {
    const response = await apiClient.get(`/c2c/listings/${id}`);
    return response.data;
  }

  async createListing(data: Partial<C2CListing>): Promise<C2CListing> {
    const response = await apiClient.post('/c2c/listings', data);
    return response.data;
  }

  async updateListing(id: string, data: Partial<C2CListing>): Promise<C2CListing> {
    const response = await apiClient.put(`/c2c/listings/${id}`, data);
    return response.data;
  }

  async deleteListing(id: string): Promise<void> {
    await apiClient.delete(`/c2c/listings/${id}`);
  }

  // ============================================================================
  // CATEGORIES
  // ============================================================================

  async getCategories(): Promise<C2CCategory[]> {
    const response = await apiClient.get('/c2c/categories');
    return response.data;
  }

  // ============================================================================
  // IMAGES
  // ============================================================================

  async uploadImages(listingId: string, files: File[]): Promise<C2CImage[]> {
    const formData = new FormData();
    files.forEach(file => formData.append('images', file));

    const response = await apiClient.post(`/c2c/listings/${listingId}/images`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    });
    return response.data;
  }

  async deleteImage(imageId: string): Promise<void> {
    await apiClient.delete(`/c2c/images/${imageId}`);
  }

  // ============================================================================
  // SEARCH
  // ============================================================================

  async search(query: string, filters?: Record<string, any>): Promise<C2CListing[]> {
    const response = await apiClient.get('/c2c/search', {
      params: { q: query, ...filters }
    });
    return response.data;
  }

  // ============================================================================
  // FAVORITES
  // ============================================================================

  async getFavorites(): Promise<C2CListing[]> {
    const response = await apiClient.get('/c2c/favorites');
    return response.data;
  }

  async addFavorite(listingId: string): Promise<void> {
    await apiClient.post('/c2c/favorites', { listingId });
  }

  async removeFavorite(favoriteId: string): Promise<void> {
    await apiClient.delete(`/c2c/favorites/${favoriteId}`);
  }
}

export const c2cApi = new C2CApi();
```

Аналогично для `src/services/b2cApi.ts`.

#### 4.4 Обновление переводов (i18n)

```bash
#!/bin/bash
# Файл: /data/hostel-booking-system/migration-tools/frontend-rename-i18n.sh

set -e
cd /data/hostel-booking-system/frontend/svetu

echo "🌐 Renaming translation files..."

for lang in en ru sr; do
  # Переименовать файлы
  git mv src/messages/$lang/marketplace.json src/messages/$lang/c2c.json || true
  git mv src/messages/$lang/storefronts.json src/messages/$lang/b2c.json || true

  # Обновить ключи внутри файлов
  if [ -f src/messages/$lang/c2c.json ]; then
    # Заменить все ключи marketplace.* на c2c.*
    sed -i 's/"marketplace\./"c2c./g' src/messages/$lang/c2c.json
  fi

  if [ -f src/messages/$lang/b2c.json ]; then
    # Заменить все ключи storefronts.* на b2c.*
    sed -i 's/"storefronts\./"b2c./g' src/messages/$lang/b2c.json
  fi
done

echo "✅ Translation files updated!"
```

Пример обновлённого файла `src/messages/en/c2c.json`:

```json
{
  "c2c": {
    "title": "C2C Marketplace",
    "subtitle": "Buy and sell directly with other users",
    "listings": {
      "title": "Listings",
      "create": "Create Listing",
      "edit": "Edit Listing",
      "delete": "Delete Listing",
      "no_results": "No listings found"
    },
    "categories": {
      "title": "Categories",
      "select": "Select Category",
      "all": "All Categories"
    },
    "search": {
      "placeholder": "Search listings...",
      "button": "Search",
      "filters": "Filters"
    },
    "favorites": {
      "title": "My Favorites",
      "add": "Add to Favorites",
      "remove": "Remove from Favorites"
    }
  }
}
```

#### 4.5 Обновление роутинга и компонентов

Пример страницы `src/app/[locale]/c2c/page.tsx`:

```tsx
import { useTranslations } from 'next-intl';
import { C2CListings } from '@/components/c2c/C2CListings';
import { C2CFilters } from '@/components/c2c/C2CFilters';

export default function C2CPage() {
  const t = useTranslations('c2c');

  return (
    <div className="container mx-auto py-8">
      <h1 className="text-3xl font-bold mb-6">{t('title')}</h1>
      <p className="text-gray-600 mb-8">{t('subtitle')}</p>

      <div className="grid grid-cols-1 lg:grid-cols-4 gap-8">
        <aside className="lg:col-span-1">
          <C2CFilters />
        </aside>

        <main className="lg:col-span-3">
          <C2CListings />
        </main>
      </div>
    </div>
  );
}
```

### Критерии приёмки
- [ ] ✅ Все директории переименованы
- [ ] ✅ `yarn build` проходит без ошибок
- [ ] ✅ TypeScript: 0 ошибок компиляции
- [ ] ✅ ESLint: `yarn lint` проходит (0 warnings)
- [ ] ✅ Prettier: `yarn format` выполнен
- [ ] ✅ Все переводы обновлены (3 языка)
- [ ] ✅ Роуты работают: `/c2c`, `/b2c`
- [ ] ✅ Нет упоминаний "marketplace" или "storefronts" в src/

### Риски
| Риск | Митигация |
|------|-----------|
| Сломанные импорты | TypeScript проверка + yarn build |
| Потерянные переводы | Проверить все 3 языка вручную |
| Роуты 404 | Тестировать каждый роут после обновления |

### Фактические данные
- **Время**: ___ дней (план: 4-5 дней)
- **Компонентов обновлено**: ___
- **Проблемы**: ___

---

## 🔍 ФАЗА 5: OPENSEARCH МИГРАЦИЯ (Дни 20-22)

### Цель
Создать новые индексы, мигрировать данные, обновить код поиска.

### Задачи

#### 5.1 Создание новых индексов

**Файл**: `/data/hostel-booking-system/migration-tools/opensearch-create-indices.sh`

```bash
#!/bin/bash
set -e

OPENSEARCH_URL="http://localhost:9200"

echo "🔍 Creating new OpenSearch indices..."

# ============================================================================
# C2C LISTINGS INDEX
# ============================================================================

curl -X PUT "$OPENSEARCH_URL/c2c_listings" \
  -H 'Content-Type: application/json' \
  -d '{
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 0,
    "analysis": {
      "analyzer": {
        "custom_analyzer": {
          "type": "custom",
          "tokenizer": "standard",
          "filter": ["lowercase", "asciifolding"]
        }
      }
    }
  },
  "mappings": {
    "properties": {
      "id": { "type": "keyword" },
      "title": {
        "type": "text",
        "analyzer": "custom_analyzer",
        "fields": {
          "keyword": { "type": "keyword" }
        }
      },
      "description": {
        "type": "text",
        "analyzer": "custom_analyzer"
      },
      "price": { "type": "float" },
      "currency": { "type": "keyword" },
      "category_id": { "type": "keyword" },
      "category_name": { "type": "text" },
      "user_id": { "type": "keyword" },
      "status": { "type": "keyword" },
      "condition": { "type": "keyword" },
      "location": { "type": "text" },
      "tags": { "type": "keyword" },
      "created_at": { "type": "date" },
      "updated_at": { "type": "date" }
    }
  }
}'

echo "✅ C2C listings index created"

# ============================================================================
# B2C PRODUCTS INDEX
# ============================================================================

curl -X PUT "$OPENSEARCH_URL/b2c_products" \
  -H 'Content-Type: application/json' \
  -d '{
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 0,
    "analysis": {
      "analyzer": {
        "custom_analyzer": {
          "type": "custom",
          "tokenizer": "standard",
          "filter": ["lowercase", "asciifolding"]
        }
      }
    }
  },
  "mappings": {
    "properties": {
      "id": { "type": "keyword" },
      "name": {
        "type": "text",
        "analyzer": "custom_analyzer",
        "fields": {
          "keyword": { "type": "keyword" }
        }
      },
      "description": {
        "type": "text",
        "analyzer": "custom_analyzer"
      },
      "price": { "type": "float" },
      "currency": { "type": "keyword" },
      "store_id": { "type": "keyword" },
      "store_name": { "type": "text" },
      "category": { "type": "keyword" },
      "status": { "type": "keyword" },
      "in_stock": { "type": "boolean" },
      "quantity": { "type": "integer" },
      "tags": { "type": "keyword" },
      "created_at": { "type": "date" },
      "updated_at": { "type": "date" }
    }
  }
}'

echo "✅ B2C products index created"

# ============================================================================
# B2C STORES INDEX
# ============================================================================

curl -X PUT "$OPENSEARCH_URL/b2c_stores" \
  -H 'Content-Type: application/json' \
  -d '{
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 0,
    "analysis": {
      "analyzer": {
        "custom_analyzer": {
          "type": "custom",
          "tokenizer": "standard",
          "filter": ["lowercase", "asciifolding"]
        }
      }
    }
  },
  "mappings": {
    "properties": {
      "id": { "type": "keyword" },
      "name": {
        "type": "text",
        "analyzer": "custom_analyzer",
        "fields": {
          "keyword": { "type": "keyword" }
        }
      },
      "description": {
        "type": "text",
        "analyzer": "custom_analyzer"
      },
      "slug": { "type": "keyword" },
      "owner_id": { "type": "keyword" },
      "status": { "type": "keyword" },
      "location": { "type": "text" },
      "tags": { "type": "keyword" },
      "created_at": { "type": "date" },
      "updated_at": { "type": "date" }
    }
  }
}'

echo "✅ B2C stores index created"

# Проверка
curl -s "$OPENSEARCH_URL/_cat/indices?v" | grep -E "c2c|b2c"
```

#### 5.2 Переиндексация данных

```bash
#!/bin/bash
# Файл: /data/hostel-booking-system/migration-tools/opensearch-reindex.sh

set -e

OPENSEARCH_URL="http://localhost:9200"

echo "🔄 Reindexing data..."

# ============================================================================
# REINDEX C2C LISTINGS
# ============================================================================

curl -X POST "$OPENSEARCH_URL/_reindex?wait_for_completion=true" \
  -H 'Content-Type: application/json' \
  -d '{
  "source": {
    "index": "marketplace_listings"
  },
  "dest": {
    "index": "c2c_listings"
  }
}'

echo "✅ C2C listings reindexed"

# ============================================================================
# REINDEX B2C PRODUCTS
# ============================================================================

curl -X POST "$OPENSEARCH_URL/_reindex?wait_for_completion=true" \
  -H 'Content-Type: application/json' \
  -d '{
  "source": {
    "index": "storefront_products"
  },
  "dest": {
    "index": "b2c_products"
  }
}'

echo "✅ B2C products reindexed"

# ============================================================================
# REINDEX B2C STORES
# ============================================================================

curl -X POST "$OPENSEARCH_URL/_reindex?wait_for_completion=true" \
  -H 'Content-Type: application/json' \
  -d '{
  "source": {
    "index": "storefronts"
  },
  "dest": {
    "index": "b2c_stores"
  }
}'

echo "✅ B2C stores reindexed"

# ============================================================================
# VERIFICATION
# ============================================================================

echo "📊 Document counts:"
curl -s "$OPENSEARCH_URL/c2c_listings/_count" | jq -r '"C2C listings: " + (.count | tostring)'
curl -s "$OPENSEARCH_URL/b2c_products/_count" | jq -r '"B2C products: " + (.count | tostring)'
curl -s "$OPENSEARCH_URL/b2c_stores/_count" | jq -r '"B2C stores: " + (.count | tostring)'

echo "✅ Reindexing complete!"
```

#### 5.3 Обновление кода поиска (Backend)

**Файл**: `backend/internal/proj/c2c/storage/opensearch/repository.go`

```go
package opensearch

import (
	"context"
	"encoding/json"

	"github.com/opensearch-project/opensearch-go/v2"
	"backend/internal/domain"
)

const (
	c2cListingsIndex = "c2c_listings"  // было "marketplace_listings"
)

type Repository struct {
	client *opensearch.Client
}

func NewRepository(client *opensearch.Client) *Repository {
	return &Repository{client: client}
}

// SearchListings выполняет поиск по C2C объявлениям
func (r *Repository) SearchListings(ctx context.Context, query string, filters map[string]interface{}) ([]domain.C2CListing, error) {
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"multi_match": map[string]interface{}{
							"query":  query,
							"fields": []string{"title^2", "description", "tags"},
						},
					},
				},
				"filter": buildFilters(filters),
			},
		},
	}

	queryJSON, _ := json.Marshal(searchQuery)

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex(c2cListingsIndex),
		r.client.Search.WithBody(strings.NewReader(string(queryJSON))),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Parse results...
	var results []domain.C2CListing
	// ... parsing logic

	return results, nil
}

func buildFilters(filters map[string]interface{}) []map[string]interface{} {
	// ... фильтры
	return nil
}
```

Аналогично для B2C (`backend/internal/proj/b2c/storage/opensearch/product_repository.go`).

### Критерии приёмки
- [ ] ✅ Новые индексы созданы (3 шт)
- [ ] ✅ Данные полностью переиндексированы
- [ ] ✅ COUNT совпадает (старый индекс = новый индекс)
- [ ] ✅ Код поиска обновлён и тестирован
- [ ] ✅ Тестовые запросы возвращают результаты

### Фактические данные
- **Время**: ___ дней (план: 2-3 дня)
- **Документов мигрировано**: ___
- **Проблемы**: ___

---

## 📦 ФАЗА 6: MINIO/S3 МИГРАЦИЯ (Дни 23-24)

### Цель
Переименовать bucket'ы изображений, обновить код загрузки.

### Задачи

#### 6.1 Создание новых bucket'ов

```bash
#!/bin/bash
# Файл: /data/hostel-booking-system/migration-tools/minio-migrate-buckets.sh

set -e

echo "📦 Migrating MinIO buckets..."

# Проверить доступность mc (MinIO Client)
if ! command -v mc &> /dev/null; then
    echo "❌ MinIO Client (mc) not found. Install: https://min.io/docs/minio/linux/reference/minio-mc.html"
    exit 1
fi

# Настроить alias (если ещё не настроено)
mc alias set local http://localhost:9000 minioadmin minioadmin || true

# ============================================================================
# CREATE NEW BUCKETS
# ============================================================================

echo "🪣 Creating new buckets..."
mc mb local/c2c-images --ignore-existing
mc mb local/b2c-images --ignore-existing

# Set public read policy (если нужно)
mc anonymous set download local/c2c-images
mc anonymous set download local/b2c-images

# ============================================================================
# COPY DATA FROM OLD BUCKETS
# ============================================================================

echo "📋 Copying data from old buckets..."

# C2C images (marketplace → c2c)
mc mirror local/marketplace-images local/c2c-images --preserve

# B2C images (storefront → b2c)
mc mirror local/storefront-images local/b2c-images --preserve

# ============================================================================
# VERIFICATION
# ============================================================================

echo "🔍 Verifying migration..."

OLD_C2C_COUNT=$(mc ls local/marketplace-images --recursive | wc -l)
NEW_C2C_COUNT=$(mc ls local/c2c-images --recursive | wc -l)

OLD_B2C_COUNT=$(mc ls local/storefront-images --recursive | wc -l)
NEW_B2C_COUNT=$(mc ls local/b2c-images --recursive | wc -l)

echo "C2C images: $OLD_C2C_COUNT → $NEW_C2C_COUNT"
echo "B2C images: $OLD_B2C_COUNT → $NEW_B2C_COUNT"

if [ "$OLD_C2C_COUNT" -ne "$NEW_C2C_COUNT" ]; then
    echo "⚠️  C2C image count mismatch!"
    exit 1
fi

if [ "$OLD_B2C_COUNT" -ne "$NEW_B2C_COUNT" ]; then
    echo "⚠️  B2C image count mismatch!"
    exit 1
fi

echo "✅ MinIO migration complete!"
```

#### 6.2 Обновление кода (Backend)

**Файл**: `backend/internal/storage/minio/client.go`

```go
package minio

import (
	"context"
	"github.com/minio/minio-go/v7"
)

const (
	C2CImagesBucket = "c2c-images"  // было "marketplace-images"
	B2CImagesBucket = "b2c-images"  // было "storefront-images"
)

type Client struct {
	client *minio.Client
}

func NewClient(endpoint, accessKey, secretKey string, useSSL bool) (*Client, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	return &Client{client: minioClient}, nil
}

// UploadC2CImage загружает изображение C2C объявления
func (c *Client) UploadC2CImage(ctx context.Context, objectName string, file io.Reader, size int64, contentType string) error {
	_, err := c.client.PutObject(ctx, C2CImagesBucket, objectName, file, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

// UploadB2CImage загружает изображение B2C продукта
func (c *Client) UploadB2CImage(ctx context.Context, objectName string, file io.Reader, size int64, contentType string) error {
	_, err := c.client.PutObject(ctx, B2CImagesBucket, objectName, file, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

// DeleteC2CImage удаляет изображение C2C объявления
func (c *Client) DeleteC2CImage(ctx context.Context, objectName string) error {
	return c.client.RemoveObject(ctx, C2CImagesBucket, objectName, minio.RemoveObjectOptions{})
}

// DeleteB2CImage удаляет изображение B2C продукта
func (c *Client) DeleteB2CImage(ctx context.Context, objectName string) error {
	return c.client.RemoveObject(ctx, B2CImagesBucket, objectName, minio.RemoveObjectOptions{})
}
```

### Критерии приёмки
- [ ] ✅ Новые bucket'ы созданы
- [ ] ✅ Все файлы скопированы (COUNT совпадает)
- [ ] ✅ Код обновлён и тестирован
- [ ] ✅ Загрузка новых файлов работает

### Фактические данные
- **Время**: ___ дней (план: 1-2 дня)
- **Файлов мигрировано**: ___

---

## ✅ ФАЗА 7: ТЕСТИРОВАНИЕ (Дни 25-29)

### Цель
100% покрытие критических путей, интеграционные тесты, E2E проверка.

### Задачи

#### 7.1 Unit тесты (Backend)

```bash
cd /data/hostel-booking-system/backend

# Запустить все тесты
go test ./... -v -race -coverprofile=coverage.out

# Посмотреть покрытие
go tool cover -html=coverage.out -o coverage.html

# ТРЕБОВАНИЕ: минимум 80% покрытия для critical paths
```

Пример теста для C2C repository:

```go
// backend/internal/proj/c2c/storage/postgres/repository_test.go
package postgres_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"backend/internal/domain"
	"backend/internal/proj/c2c/storage/postgres"
)

func TestRepository_CreateListing(t *testing.T) {
	// Setup test DB
	db := setupTestDB(t)
	defer db.Close()

	repo := postgres.NewRepository(db)

	listing := &domain.C2CListing{
		Title:       "Test Listing",
		Description: "Test description",
		Price:       100.0,
		Currency:    "USD",
		CategoryID:  "category-123",
		UserID:      "user-456",
		Status:      "active",
	}

	// Test
	created, err := repo.CreateListing(context.Background(), listing)

	// Assertions
	require.NoError(t, err)
	assert.NotEmpty(t, created.ID)
	assert.Equal(t, listing.Title, created.Title)
	assert.NotZero(t, created.CreatedAt)
}

func TestRepository_GetListingByID(t *testing.T) {
	// ... аналогично
}

// ... остальные тесты CRUD операций
```

#### 7.2 Integration тесты (API)

```bash
#!/bin/bash
# Файл: /data/hostel-booking-system/migration-tools/test-api-endpoints.sh

set -e

API_URL="http://localhost:3000"
TOKEN=$(cat /tmp/token)  # JWT для авторизации

echo "🧪 Testing API endpoints..."

# ============================================================================
# C2C ENDPOINTS
# ============================================================================

echo "Testing C2C endpoints..."

# GET /api/v1/c2c/listings
curl -s -X GET "$API_URL/api/v1/c2c/listings" | jq '.data | length'

# GET /api/v1/c2c/categories
curl -s -X GET "$API_URL/api/v1/c2c/categories" | jq '.data | length'

# POST /api/v1/c2c/listings (authenticated)
LISTING_ID=$(curl -s -X POST "$API_URL/api/v1/c2c/listings" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Listing",
    "description": "Test",
    "price": 100,
    "currency": "USD",
    "category_id": "category-123"
  }' | jq -r '.data.id')

echo "Created listing: $LISTING_ID"

# GET /api/v1/c2c/listings/:id
curl -s -X GET "$API_URL/api/v1/c2c/listings/$LISTING_ID" | jq '.data.title'

# DELETE /api/v1/c2c/listings/:id
curl -s -X DELETE "$API_URL/api/v1/c2c/listings/$LISTING_ID" \
  -H "Authorization: Bearer $TOKEN"

echo "✅ C2C endpoints OK"

# ============================================================================
# B2C ENDPOINTS
# ============================================================================

echo "Testing B2C endpoints..."

# GET /api/v1/b2c/stores
curl -s -X GET "$API_URL/api/v1/b2c/stores" | jq '.data | length'

# GET /api/v1/b2c/products
curl -s -X GET "$API_URL/api/v1/b2c/products" | jq '.data | length'

echo "✅ B2C endpoints OK"

echo "✅ All API tests passed!"
```

#### 7.3 E2E тесты (Playwright/Cypress)

Пример теста `frontend/svetu/e2e/c2c.spec.ts`:

```typescript
import { test, expect } from '@playwright/test';

test.describe('C2C Marketplace', () => {
  test('should display listings page', async ({ page }) => {
    await page.goto('/c2c');

    // Проверить заголовок
    await expect(page.locator('h1')).toContainText('C2C Marketplace');

    // Проверить наличие списка объявлений
    const listings = page.locator('[data-testid="c2c-listing"]');
    await expect(listings).toHaveCount({ min: 1 });
  });

  test('should create new listing', async ({ page }) => {
    // Login
    await page.goto('/auth/login');
    await page.fill('[name="email"]', 'test@example.com');
    await page.fill('[name="password"]', 'password123');
    await page.click('button[type="submit"]');

    // Navigate to create page
    await page.goto('/c2c/create');

    // Fill form
    await page.fill('[name="title"]', 'Test Listing');
    await page.fill('[name="description"]', 'Test description');
    await page.fill('[name="price"]', '100');
    await page.selectOption('[name="category"]', 'electronics');

    // Submit
    await page.click('button[type="submit"]');

    // Verify redirect
    await expect(page).toHaveURL(/\/c2c\/\d+/);
    await expect(page.locator('h1')).toContainText('Test Listing');
  });

  test('should search listings', async ({ page }) => {
    await page.goto('/c2c');

    // Search
    await page.fill('[data-testid="search-input"]', 'laptop');
    await page.click('[data-testid="search-button"]');

    // Verify results
    await expect(page.locator('[data-testid="search-results"]')).toBeVisible();
  });
});
```

### Критерии приёмки
- [ ] ✅ Backend: `go test ./...` проходит (100% тестов)
- [ ] ✅ Backend: минимум 80% покрытие critical paths
- [ ] ✅ API тесты: все endpoints работают (C2C + B2C)
- [ ] ✅ Frontend: `yarn test` проходит (100% тестов)
- [ ] ✅ E2E тесты: основные пользовательские сценарии работают
- [ ] ✅ Нет регрессий (старая функциональность сохранена)

### Фактические данные
- **Время**: ___ дней (план: 3-5 дней)
- **Тестов написано**: ___
- **Найдено багов**: ___
- **Покрытие кода**: ___%

---

## 🚀 ФАЗА 8: PRODUCTION ДЕПЛОЙ (Дни 30-32)

### Цель
Развернуть на dev.svetu.rs, провести smoke testing, подготовить к production.

### Задачи

#### 8.1 Pre-deployment checklist

```markdown
# PRE-DEPLOYMENT CHECKLIST

## Code Quality
- [ ] ✅ `make format` (backend) - passed
- [ ] ✅ `make lint` (backend) - 0 warnings
- [ ] ✅ `yarn format` (frontend) - passed
- [ ] ✅ `yarn lint` (frontend) - 0 warnings
- [ ] ✅ `go test ./...` - 100% passed
- [ ] ✅ `yarn test` - 100% passed
- [ ] ✅ E2E tests - passed

## Database
- [ ] ✅ Migrations tested (up + down)
- [ ] ✅ Production backup created
- [ ] ✅ Data integrity verified

## Documentation
- [ ] ✅ CLAUDE.md updated (новые endpoint'ы)
- [ ] ✅ Swagger regenerated
- [ ] ✅ README updated (если нужно)
- [ ] ✅ Migration plan updated (этот документ)

## Deployment
- [ ] ✅ .env переменные проверены
- [ ] ✅ OpenSearch индексы созданы
- [ ] ✅ MinIO bucket'ы готовы
```

#### 8.2 Staging deployment

```bash
#!/bin/bash
# Файл: /data/hostel-booking-system/migration-tools/deploy-staging.sh

set -e

echo "🚀 Deploying to dev.svetu.rs..."

# 1. Commit и push
git add -A
git commit -m "feat: C2C/B2C migration complete - ready for staging"
git push origin feature/c2c-b2c-migration

# 2. Дамп локальной БД
echo "📦 Creating database dump..."
PGPASSWORD=mX3g1XGhMRUZEX3l pg_dump \
  -h localhost \
  -U postgres \
  -d svetubd \
  --no-owner \
  --no-acl \
  --column-inserts \
  --inserts \
  -f svetubd_full_dump_$(date +%Y%m%d_%H%M%S).sql

# Использовать последний созданный дамп для загрузки
DUMP_FILE=$(ls -t svetubd_full_dump_*.sql | head -1)
echo "Using dump file: $DUMP_FILE"

# 3. Загрузка на сервер
echo "📤 Uploading to server..."
scp $DUMP_FILE svetu@svetu.rs:/tmp/dump_migration.sql

# 4. Деплой на сервере
ssh svetu@svetu.rs << 'ENDSSH'
set -e

echo "🔧 Deploying on dev server..."

cd /opt/svetu-dev

# Pull changes
git fetch origin
git checkout feature/c2c-b2c-migration
git pull

# Восстановить БД
echo "🗄️  Restoring database..."
docker exec -i svetu-dev_db_1 psql -U svetu_dev_user -d svetu_dev_db << EOF
DROP SCHEMA public CASCADE;
CREATE SCHEMA public;
EOF

docker exec -i svetu-dev_db_1 psql -U svetu_dev_user -d svetu_dev_db < /tmp/dump_migration.sql

# Backend restart
echo "🔄 Restarting backend..."
cd backend && make dev-restart

# Frontend restart
echo "🎨 Restarting frontend..."
cd ../frontend/svetu && make dev-restart

echo "✅ Deployment complete!"
ENDSSH

echo "🎉 Staging deployment finished!"
echo "🔗 https://dev.svetu.rs"
```

#### 8.3 Smoke testing

```bash
#!/bin/bash
# Файл: /data/hostel-booking-system/migration-tools/smoke-test.sh

set -e

DEV_API="https://devapi.svetu.rs"
DEV_WEB="https://dev.svetu.rs"

echo "🔬 Running smoke tests on dev server..."

# ============================================================================
# API HEALTH CHECK
# ============================================================================

echo "Checking API health..."
HEALTH=$(curl -s "$DEV_API/" | grep -o "Svetu API")
if [ "$HEALTH" != "Svetu API" ]; then
    echo "❌ API health check failed!"
    exit 1
fi
echo "✅ API is healthy"

# ============================================================================
# C2C ENDPOINTS
# ============================================================================

echo "Testing C2C endpoints..."

# GET /api/v1/c2c/listings
C2C_LISTINGS=$(curl -s "$DEV_API/api/v1/c2c/listings" | jq -r '.data | length')
echo "C2C listings found: $C2C_LISTINGS"

# GET /api/v1/c2c/categories
C2C_CATEGORIES=$(curl -s "$DEV_API/api/v1/c2c/categories" | jq -r '.data | length')
echo "C2C categories found: $C2C_CATEGORIES"

if [ "$C2C_LISTINGS" -eq 0 ] || [ "$C2C_CATEGORIES" -eq 0 ]; then
    echo "⚠️  Warning: C2C data might be missing"
fi

# ============================================================================
# B2C ENDPOINTS
# ============================================================================

echo "Testing B2C endpoints..."

# GET /api/v1/b2c/stores
B2C_STORES=$(curl -s "$DEV_API/api/v1/b2c/stores" | jq -r '.data | length')
echo "B2C stores found: $B2C_STORES"

# GET /api/v1/b2c/products
B2C_PRODUCTS=$(curl -s "$DEV_API/api/v1/b2c/products" | jq -r '.data | length')
echo "B2C products found: $B2C_PRODUCTS"

# ============================================================================
# FRONTEND CHECK
# ============================================================================

echo "Checking frontend pages..."

# C2C page
C2C_STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$DEV_WEB/c2c")
if [ "$C2C_STATUS" -ne 200 ]; then
    echo "❌ C2C page returned $C2C_STATUS"
    exit 1
fi
echo "✅ C2C page OK ($C2C_STATUS)"

# B2C page
B2C_STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$DEV_WEB/b2c")
if [ "$B2C_STATUS" -ne 200 ]; then
    echo "❌ B2C page returned $B2C_STATUS"
    exit 1
fi
echo "✅ B2C page OK ($B2C_STATUS)"

echo "✅ All smoke tests passed!"
```

#### 8.4 Удаление старых сущностей

**ТОЛЬКО ПОСЛЕ** 2-3 дней успешной работы на staging!

**Файл**: `backend/migrations/000175_drop_old_tables.up.sql`

```sql
-- ============================================================================
-- УДАЛЕНИЕ СТАРЫХ ТАБЛИЦ
-- ВНИМАНИЕ: Выполнять ТОЛЬКО после успешной проверки миграции!
-- ============================================================================

BEGIN;

-- Дроп старых C2C таблиц (бывшие marketplace_*)
DROP TABLE IF EXISTS marketplace_listing_variants CASCADE;
DROP TABLE IF EXISTS marketplace_orders CASCADE;
DROP TABLE IF EXISTS marketplace_favorites CASCADE;
DROP TABLE IF EXISTS marketplace_messages CASCADE;
DROP TABLE IF EXISTS marketplace_chats CASCADE;
DROP TABLE IF EXISTS marketplace_images CASCADE;
DROP TABLE IF EXISTS marketplace_listings CASCADE;
DROP TABLE IF EXISTS marketplace_categories CASCADE;

-- Дроп старых B2C таблиц (бывшие storefront_*)
DROP TABLE IF EXISTS storefront_product_variant_images CASCADE;
DROP TABLE IF EXISTS user_storefronts CASCADE;
DROP TABLE IF EXISTS storefront_inventory_movements CASCADE;
DROP TABLE IF EXISTS storefront_delivery_options CASCADE;
DROP TABLE IF EXISTS storefront_payment_methods CASCADE;
DROP TABLE IF EXISTS storefront_staff CASCADE;
DROP TABLE IF EXISTS storefront_hours CASCADE;
DROP TABLE IF EXISTS storefront_favorites CASCADE;
DROP TABLE IF EXISTS storefront_order_items CASCADE;
DROP TABLE IF EXISTS storefront_orders CASCADE;
DROP TABLE IF EXISTS storefront_product_attributes CASCADE;
DROP TABLE IF EXISTS storefront_product_variants CASCADE;
DROP TABLE IF EXISTS storefront_product_images CASCADE;
DROP TABLE IF EXISTS storefront_products CASCADE;
DROP TABLE IF EXISTS storefronts CASCADE;

COMMIT;

-- ============================================================================
-- Verification
-- ============================================================================

DO $$
BEGIN
    RAISE NOTICE '✅ Old tables dropped successfully';
END $$;
```

**Удаление старых OpenSearch индексов:**

```bash
# ТОЛЬКО ПОСЛЕ проверки работы новых индексов!
curl -X DELETE "http://localhost:9200/marketplace_listings"
curl -X DELETE "http://localhost:9200/storefront_products"
curl -X DELETE "http://localhost:9200/storefronts"
```

**Удаление старых MinIO bucket'ов:**

```bash
# ТОЛЬКО ПОСЛЕ проверки новых bucket'ов!
mc rm local/marketplace-images --recursive --force
mc rb local/marketplace-images

mc rm local/storefront-images --recursive --force
mc rb local/storefront-images
```

### Критерии приёмки
- [ ] ✅ Pre-deployment checklist: все пункты ✅
- [ ] ✅ Staging deployment прошёл успешно
- [ ] ✅ Smoke tests проходят (100%)
- [ ] ✅ Мониторинг показывает нормальную работу
- [ ] ✅ Нет ошибок в логах (backend + frontend)
- [ ] ✅ Старые сущности удалены (через 2-3 дня)

### Фактические данные
- **Время**: ___ дней (план: 2-3 дня)
- **Проблемы при деплое**: ___
- **Downtime**: ___ минут

---

## 🎯 ИТОГОВАЯ ПРОВЕРКА КАЧЕСТВА

### Production-Ready Checklist

```markdown
# PRODUCTION READINESS ASSESSMENT

## Code Quality (100/100)
- [ ] Backend lint: 0 warnings
- [ ] Frontend lint: 0 warnings
- [ ] Backend tests: 100% pass, ≥80% coverage
- [ ] Frontend tests: 100% pass
- [ ] E2E tests: критические сценарии покрыты
- [ ] Нет TODO комментариев в коде
- [ ] Нет console.log/fmt.Println в production коде
- [ ] Все ошибки логируются правильно

## Architecture (100/100)
- [ ] Нет дублирования кода (DRY принцип)
- [ ] Универсальные функции переиспользуются
- [ ] Чёткое разделение модулей (c2c vs b2c)
- [ ] Нет циклических зависимостей
- [ ] Database нормализована
- [ ] API endpoints следуют RESTful

## Documentation (100/100)
- [ ] CLAUDE.md обновлён
- [ ] Swagger актуален
- [ ] Этот план актуализирован
- [ ] Комментарии в коде понятны
- [ ] Migration guide написан

## Performance (100/100)
- [ ] Database индексы оптимальны
- [ ] OpenSearch queries быстрые
- [ ] Frontend bundle оптимизирован
- [ ] Images compressed
- [ ] Нет N+1 queries

## Security (100/100)
- [ ] JWT токены валидируются
- [ ] CORS настроен правильно
- [ ] SQL injection защита (prepared statements)
- [ ] XSS защита
- [ ] CSRF защита (через BFF proxy)
- [ ] Sensitive data не логируется

## Monitoring (100/100)
- [ ] Логи структурированы (JSON)
- [ ] Метрики собираются
- [ ] Alerts настроены
- [ ] Health checks работают

**ИТОГО**: ___/600 points

**ТРЕБОВАНИЕ ДЛЯ PRODUCTION**: ≥ 580/600 (97%)
```

---

## 📈 МЕТРИКИ УСПЕХА

### Ожидаемые результаты миграции

| Метрика | До миграции | После миграции | Улучшение |
|---------|-------------|----------------|-----------|
| **Семантическая ясность** | 3/10 | 10/10 | +233% |
| **Скорость разработки** | Baseline | +30% | Меньше путаницы |
| **Code maintainability** | 6/10 | 9/10 | +50% |
| **Onboarding новых разработчиков** | ~5 дней | ~2 дня | -60% |
| **API понятность** | 5/10 | 10/10 | +100% |
| **Технический долг** | Высокий | 0 | -100% |

### Трудозатраты (фактические vs план)

| Фаза | План (дни) | Факт (дни) | Отклонение |
|------|-----------|-----------|-----------|
| 0. Инициализация | 0.5 | ___ | ___ |
| 1. Подготовка | 1.5 | ___ | ___ |
| 2. БД миграция | 4.5 | ___ | ___ |
| 3. Backend | 6.0 | ___ | ___ |
| 4. Frontend | 4.5 | ___ | ___ |
| 5. OpenSearch | 2.5 | ___ | ___ |
| 6. MinIO | 1.5 | ___ | ___ |
| 7. Тестирование | 4.0 | ___ | ___ |
| 8. Деплой | 2.5 | ___ | ___ |
| **ИТОГО** | **27.5** | **___** | **___** |

---

## 🔄 ПРОЦЕСС АКТУАЛИЗАЦИИ ПЛАНА

### ⚠️ ОБЯЗАТЕЛЬНО ОБНОВЛЯТЬ ПОСЛЕ КАЖДОЙ ФАЗЫ!

1. **Статус фазы**: ⏸️ Pending → 🚧 In Progress → ✅ Completed
2. **Прогресс**: Обновить % выполнения
3. **Качество**: Добавить оценку (0-100)
4. **Даты**: Заполнить фактические даты начала/окончания
5. **Фактические данные**: Заполнить все поля "Фактические данные"
6. **Проблемы**: Задокументировать все проблемы и решения
7. **Риски**: Обновить статусы рисков (новые/закрытые)
8. **Метрики**: Обновить таблицу трудозатрат

### Пример обновления:

```markdown
## ФАЗА 2: МИГРАЦИЯ БАЗЫ ДАННЫХ

**Статус**: ✅ Completed
**Прогресс**: 100%
**Качество**: 95/100 (отлично)
**Дата начала**: 2025-10-11
**Дата окончания**: 2025-10-15

### Фактические данные
- **Время**: 4.2 дня (план: 4.5 дня) ✅ -7% лучше плана
- **Мигрировано строк**: C2C=1,234, B2C=5,678
- **Проблемы**:
  1. Триггер update_timestamp не скопировался - пришлось пересоздавать вручную
  2. FK constraint на b2c_orders сломался - исправлено в миграции 174
- **Решения**:
  1. Добавлен скрипт проверки триггеров
  2. Улучшена валидация FK в миграции

### Обновлённые риски
- ❌ Потеря данных - риск закрыт (тройная проверка прошла)
- ✅ Время выполнения - в пределах плана
```

---

## 📝 ИСТОРИЯ ИЗМЕНЕНИЙ ПЛАНА

| Дата | Версия | Изменения | Автор |
|------|--------|-----------|-------|
| 2025-10-09 | 1.0 | Создание детального плана | Claude Code |
| ___ | 1.1 | Обновление после фазы 0 | ___ |
| ___ | 1.2 | Обновление после фазы 1 | ___ |
| ___ | 2.0 | Major update после завершения миграции | ___ |

---

## 🎯 ФИНАЛЬНЫЙ СТАТУС

**Дата завершения миграции**: ___
**Общий прогресс**: ___% (___/8 фаз)
**Production-ready score**: ___/600
**Готовность к production**: ⏸️ Pending / ✅ Ready

**Следующие шаги**:
- [ ] Финальный code review
- [ ] Создание PR в main
- [ ] Production deployment
- [ ] Post-mortem анализ

---

**ZERO TECHNICAL DEBT ACHIEVED**: ⏸️ / ✅

**🚀 ГОТОВНОСТЬ К PRODUCTION**: ⏸️ Pending / ✅ Ready

---

*Этот план является живым документом и ОБЯЗАН обновляться на каждом этапе миграции!*

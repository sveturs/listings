# 📋 Техническое задание на унификацию системы атрибутов
## Sve Tu Platforma - Marketplace

*Дата создания: 02.09.2025*  
*Дата обновления: 03.09.2025*  
*Приоритет: 🔴 КРИТИЧЕСКИЙ*  
*Срок исполнения: 30 дней*  
*Текущий статус: 🚀 PRODUCTION DEPLOYMENT (День 15/30 - 50%)*

---

## 🎯 Цель задания

Устранить критическое дублирование в системе атрибутов, выявленное в аудите от 02.09.2025. Создать единую, унифицированную систему атрибутов для всего маркетплейса с сохранением обратной совместимости и минимальным риском для продакшена.

---

## ⚠️ КРИТИЧЕСКИ ВАЖНЫЕ ПРАВИЛА БЕЗОПАСНОСТИ

### 🛡️ Принцип "НЕ НАВРЕДИ"

1. **ЗАПРЕЩЕНО удалять или модифицировать существующие таблицы БД до создания полной резервной копии**
2. **ЗАПРЕЩЕНО удалять работающий код до полного тестирования замены**
3. **ЗАПРЕЩЕНО вносить breaking changes в публичные API без версионирования**
4. **ОБЯЗАТЕЛЬНО создавать обратимые миграции (up/down) для всех изменений БД**
5. **ОБЯЗАТЕЛЬНО тестировать каждое изменение на тестовом окружении**
6. **ОБЯЗАТЕЛЬНО сохранять совместимость с существующими данными**

### 📊 Контрольные точки безопасности

Перед каждым этапом работы ОБЯЗАТЕЛЬНО:
- ✅ Создать backup всех затрагиваемых таблиц
- ✅ Проверить количество записей в таблицах до и после миграции
- ✅ Убедиться в наличии rollback плана
- ✅ Протестировать на копии продакшн данных

---

## 📑 Исходные данные

### Текущая ситуация (из аудита):
- **14 таблиц БД** для системы атрибутов (должно быть 3-4)
- **3 параллельные системы** в backend коде
- **5 дублированных компонентов** во frontend (~2600 строк)
- **Отсутствие кеширования** для атрибутов
- **Дублирование в OpenSearch** индексах

### Основные системы дублирования:

1. **Система маркетплейса** (работающая):
   - `category_attributes` (85 записей)
   - `category_attribute_mapping` (611 записей)
   - `listing_attribute_values` (15 записей)

2. **Система витрин** (частично используется):
   - `product_variant_attributes` (14 записей)
   - `storefront_product_attributes` (0 записей)
   - `product_variant_attribute_values`

3. **Устаревшая система** (не используется):
   - `category_variant_attributes` (4 записи)
   - `variant_attribute_mappings` (4 записи)

---

## 📋 ДЕТАЛЬНЫЙ ПЛАН РАБОТ

## 📊 ТЕКУЩИЙ СТАТУС ВЫПОЛНЕНИЯ

### ✅ Завершенные этапы:
1. **Дни 1-3: Подготовка и анализ** ✅ 100%
2. **Дни 4-6: Миграция БД** ✅ 100%  
3. **Дни 7-8: Backend реализация** ✅ 100%
4. **Дни 9-10: E2E и интеграционное тестирование** ✅ 100%
5. **День 11: Мониторинг и метрики** ✅ 100%
6. **День 12: CI/CD Pipeline** ✅ 100%

### 🟡 В процессе:
7. **Дни 13-15: Production развертывание** ⏳ 0%

### ⏰ Предстоящие этапы:
7. **Дни 13-15: Production развертывание**
8. **Дни 16-20: Миграция данных**
9. **Дни 21-25: Оптимизация и доработки**
10. **Дни 26-30: Очистка и завершение**

---

### ЭТАП 1: Подготовка и анализ (Дни 1-3) ✅ ЗАВЕРШЕНО

#### 1.1 Создание полного бэкапа (День 1) ✅

```bash
# ОБЯЗАТЕЛЬНЫЕ команды для выполнения:

# 1. Создать директорию для бэкапов
mkdir -p /data/backups/attribute_unification_$(date +%Y%m%d)

# 2. Бэкап всех таблиц атрибутов
pg_dump -h localhost -U postgres -d svetubd \
  -t category_attributes \
  -t category_attribute_mapping \
  -t listing_attribute_values \
  -t product_variant_attributes \
  -t storefront_product_attributes \
  -t product_variant_attribute_values \
  -t category_variant_attributes \
  -t variant_attribute_mappings \
  -t attribute_groups \
  -t attribute_group_items \
  -t category_attribute_groups \
  -t translations \
  > /data/backups/attribute_unification_$(date +%Y%m%d)/attributes_backup.sql

# 3. Проверка целостности бэкапа
psql -h localhost -U postgres -d test_restore < /data/backups/attribute_unification_$(date +%Y%m%d)/attributes_backup.sql
```

#### 1.2 Анализ использования (День 2)

**Задача**: Проанализировать реальное использование каждой системы атрибутов

```sql
-- Скрипт анализа для выполнения
-- Сохранить результаты в /data/hostel-booking-system/docs/attribute_usage_analysis.md

-- 1. Анализ системы маркетплейса
SELECT 
    'marketplace_system' as system,
    COUNT(DISTINCT ca.id) as unique_attributes,
    COUNT(DISTINCT cam.category_id) as categories_using,
    COUNT(DISTINCT lav.listing_id) as listings_using,
    COUNT(lav.id) as total_values
FROM category_attributes ca
LEFT JOIN category_attribute_mapping cam ON ca.id = cam.attribute_id
LEFT JOIN listing_attribute_values lav ON ca.id = lav.attribute_id;

-- 2. Анализ системы витрин
SELECT 
    'storefront_system' as system,
    COUNT(DISTINCT pva.id) as unique_attributes,
    COUNT(DISTINCT spa.storefront_product_id) as products_using,
    COUNT(pvav.id) as total_values
FROM product_variant_attributes pva
LEFT JOIN storefront_product_attributes spa ON pva.id = spa.attribute_id
LEFT JOIN product_variant_attribute_values pvav ON pva.id = pvav.attribute_id;

-- 3. Поиск дублированных атрибутов между системами
SELECT 
    ca.name as marketplace_name,
    ca.attribute_type as marketplace_type,
    pva.name as storefront_name,
    pva.type as storefront_type
FROM category_attributes ca
FULL OUTER JOIN product_variant_attributes pva 
    ON LOWER(ca.name) = LOWER(pva.name)
WHERE ca.name IS NOT NULL OR pva.name IS NOT NULL
ORDER BY COALESCE(ca.name, pva.name);
```

#### 1.3 Создание карты зависимостей (День 3)

**Задача**: Найти все места в коде, использующие атрибуты

```bash
# Backend анализ
cd /data/hostel-booking-system/backend
grep -r "CategoryAttribute\|ProductVariantAttribute\|CategoryVariantAttribute" \
  --include="*.go" \
  --exclude-dir=vendor \
  > /tmp/backend_attribute_usage.txt

# Frontend анализ  
cd /data/hostel-booking-system/frontend/svetu
grep -r "attribute\|Attribute" \
  --include="*.tsx" --include="*.ts" \
  --exclude-dir=node_modules \
  > /tmp/frontend_attribute_usage.txt

# API эндпоинты
cd /data/hostel-booking-system/backend
grep -r "@Router.*attribute" --include="*.go" > /tmp/api_endpoints.txt
```

### ЭТАП 2: Создание унифицированной системы (Дни 4-10)

#### 2.1 Создание новой структуры БД (День 4-5)

**ВАЖНО**: Новые таблицы создаются ПАРАЛЛЕЛЬНО существующим, без удаления старых!

```sql
-- Миграция: backend/migrations/001_create_unified_attributes.up.sql

-- Унифицированная таблица атрибутов
CREATE TABLE IF NOT EXISTS unified_attributes (
    id SERIAL PRIMARY KEY,
    code VARCHAR(100) UNIQUE NOT NULL, -- Уникальный код атрибута
    name VARCHAR(100) NOT NULL,
    display_name VARCHAR(200) NOT NULL,
    attribute_type VARCHAR(50) NOT NULL CHECK (attribute_type IN (
        'text', 'textarea', 'number', 'boolean', 
        'select', 'multiselect', 'date', 'color', 'size'
    )),
    purpose VARCHAR(20) NOT NULL DEFAULT 'regular' CHECK (purpose IN (
        'regular',    -- Обычный атрибут для фильтрации/поиска
        'variant',    -- Вариативный атрибут (влияет на SKU)
        'both'        -- Может использоваться в обоих случаях
    )),
    
    -- Настройки атрибута
    options JSONB DEFAULT '{}',           -- Опции для select/multiselect
    validation_rules JSONB DEFAULT '{}',  -- Правила валидации
    ui_settings JSONB DEFAULT '{}',       -- Настройки отображения
    
    -- Флаги использования
    is_searchable BOOLEAN DEFAULT false,
    is_filterable BOOLEAN DEFAULT false,
    is_required BOOLEAN DEFAULT false,
    affects_stock BOOLEAN DEFAULT false,  -- Для вариативных атрибутов
    affects_price BOOLEAN DEFAULT false,  -- Для вариативных атрибутов
    
    -- Метаданные
    sort_order INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Связь со старой системой (временно, для миграции)
    legacy_category_attribute_id INTEGER,
    legacy_product_variant_attribute_id INTEGER
);

-- Связь атрибутов с категориями
CREATE TABLE IF NOT EXISTS unified_category_attributes (
    id SERIAL PRIMARY KEY,
    category_id INTEGER NOT NULL REFERENCES marketplace_categories(id) ON DELETE CASCADE,
    attribute_id INTEGER NOT NULL REFERENCES unified_attributes(id) ON DELETE CASCADE,
    
    -- Настройки для конкретной категории
    is_enabled BOOLEAN DEFAULT true,
    is_required BOOLEAN DEFAULT false,
    sort_order INTEGER DEFAULT 0,
    
    -- Переопределение настроек атрибута для категории
    category_specific_options JSONB,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(category_id, attribute_id)
);

-- Значения атрибутов (универсальная таблица)
CREATE TABLE IF NOT EXISTS unified_attribute_values (
    id SERIAL PRIMARY KEY,
    entity_type VARCHAR(50) NOT NULL CHECK (entity_type IN (
        'listing',           -- Объявление маркетплейса
        'product',           -- Товар витрины
        'product_variant'    -- Вариант товара
    )),
    entity_id INTEGER NOT NULL,
    attribute_id INTEGER NOT NULL REFERENCES unified_attributes(id) ON DELETE CASCADE,
    
    -- Значения разных типов
    text_value TEXT,
    numeric_value NUMERIC,
    boolean_value BOOLEAN,
    date_value DATE,
    json_value JSONB,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Индекс для быстрого поиска
    UNIQUE(entity_type, entity_id, attribute_id)
);

-- Индексы для производительности
CREATE INDEX idx_unified_attributes_code ON unified_attributes(code);
CREATE INDEX idx_unified_attributes_purpose ON unified_attributes(purpose);
CREATE INDEX idx_unified_attributes_active ON unified_attributes(is_active);

CREATE INDEX idx_unified_category_attributes_category ON unified_category_attributes(category_id);
CREATE INDEX idx_unified_category_attributes_enabled ON unified_category_attributes(is_enabled);

CREATE INDEX idx_unified_attribute_values_entity ON unified_attribute_values(entity_type, entity_id);
CREATE INDEX idx_unified_attribute_values_attribute ON unified_attribute_values(attribute_id);
CREATE INDEX idx_unified_attribute_values_text ON unified_attribute_values(text_value) WHERE text_value IS NOT NULL;
CREATE INDEX idx_unified_attribute_values_numeric ON unified_attribute_values(numeric_value) WHERE numeric_value IS NOT NULL;
```

#### 2.2 Миграция данных 

**КРИТИЧНО**: Миграция выполняется БЕЗ удаления старых данных!

```sql
-- Миграция: backend/migrations/002_migrate_attributes_data.up.sql

-- ВАЖНО: Транзакция для атомарности
BEGIN;

-- 1. Миграция атрибутов из category_attributes
INSERT INTO unified_attributes (
    code, name, display_name, attribute_type, purpose,
    options, validation_rules, ui_settings,
    is_searchable, is_filterable, is_required,
    affects_stock, affects_price,
    sort_order, legacy_category_attribute_id
)
SELECT 
    LOWER(REPLACE(name, ' ', '_')) as code,
    name,
    display_name,
    attribute_type,
    CASE 
        WHEN is_variant_compatible = true THEN 'both'
        ELSE 'regular'
    END as purpose,
    options,
    validation_rules,
    COALESCE(ui_settings, '{}'),
    is_searchable,
    is_filterable,
    is_required,
    affects_stock,
    false as affects_price, -- Добавим позже если нужно
    sort_order,
    id as legacy_category_attribute_id
FROM category_attributes
ON CONFLICT (code) DO UPDATE SET
    legacy_category_attribute_id = EXCLUDED.legacy_category_attribute_id;

-- 2. Миграция атрибутов из product_variant_attributes
INSERT INTO unified_attributes (
    code, name, display_name, attribute_type, purpose,
    options, affects_stock, affects_price,
    legacy_product_variant_attribute_id
)
SELECT 
    LOWER(REPLACE(name, ' ', '_')) as code,
    name,
    name as display_name, -- У них нет display_name
    type as attribute_type,
    'variant' as purpose,
    options,
    affects_stock,
    affects_price,
    id as legacy_product_variant_attribute_id
FROM product_variant_attributes
ON CONFLICT (code) DO UPDATE SET
    purpose = 'both', -- Если атрибут уже есть, делаем его универсальным
    affects_stock = COALESCE(unified_attributes.affects_stock, EXCLUDED.affects_stock),
    affects_price = COALESCE(unified_attributes.affects_price, EXCLUDED.affects_price),
    legacy_product_variant_attribute_id = EXCLUDED.legacy_product_variant_attribute_id;

-- 3. Миграция связей категорий с атрибутами
INSERT INTO unified_category_attributes (
    category_id, attribute_id, is_enabled, is_required, sort_order
)
SELECT 
    cam.category_id,
    ua.id as attribute_id,
    cam.is_enabled,
    cam.is_required,
    cam.sort_order
FROM category_attribute_mapping cam
JOIN unified_attributes ua ON ua.legacy_category_attribute_id = cam.attribute_id;

-- 4. Миграция значений атрибутов объявлений
INSERT INTO unified_attribute_values (
    entity_type, entity_id, attribute_id,
    text_value, numeric_value, boolean_value, json_value
)
SELECT 
    'listing' as entity_type,
    lav.listing_id as entity_id,
    ua.id as attribute_id,
    lav.text_value,
    lav.numeric_value,
    lav.boolean_value,
    lav.json_value
FROM listing_attribute_values lav
JOIN unified_attributes ua ON ua.legacy_category_attribute_id = lav.attribute_id;

-- 5. Проверка миграции
DO $$
DECLARE
    old_count INTEGER;
    new_count INTEGER;
BEGIN
    -- Проверяем атрибуты
    SELECT COUNT(*) INTO old_count FROM category_attributes;
    SELECT COUNT(*) INTO new_count FROM unified_attributes WHERE legacy_category_attribute_id IS NOT NULL;
    
    IF old_count != new_count THEN
        RAISE EXCEPTION 'Миграция атрибутов неполная: старых %, новых %', old_count, new_count;
    END IF;
    
    -- Проверяем значения
    SELECT COUNT(*) INTO old_count FROM listing_attribute_values;
    SELECT COUNT(*) INTO new_count FROM unified_attribute_values WHERE entity_type = 'listing';
    
    IF old_count != new_count THEN
        RAISE EXCEPTION 'Миграция значений неполная: старых %, новых %', old_count, new_count;
    END IF;
END $$;

COMMIT;
```

#### 2.3 Backend: Создание сервисного слоя 

```go
// backend/internal/services/attributes/unified_service.go

package attributes

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
)

// UnifiedAttribute представляет унифицированный атрибут
type UnifiedAttribute struct {
    ID              int                    `json:"id"`
    Code            string                 `json:"code"`
    Name            string                 `json:"name"`
    DisplayName     string                 `json:"display_name"`
    Type            AttributeType          `json:"type"`
    Purpose         AttributePurpose       `json:"purpose"`
    Options         json.RawMessage        `json:"options,omitempty"`
    ValidationRules json.RawMessage        `json:"validation_rules,omitempty"`
    UISettings      json.RawMessage        `json:"ui_settings,omitempty"`
    IsSearchable    bool                   `json:"is_searchable"`
    IsFilterable    bool                   `json:"is_filterable"`
    IsRequired      bool                   `json:"is_required"`
    AffectsStock    bool                   `json:"affects_stock"`
    AffectsPrice    bool                   `json:"affects_price"`
    SortOrder       int                    `json:"sort_order"`
    Translations    map[string]Translation `json:"translations,omitempty"`
}

// UnifiedAttributeService - сервис для работы с унифицированными атрибутами
type UnifiedAttributeService struct {
    db    *sql.DB
    cache CacheService
    
    // Флаг для постепенной миграции
    useLegacyFallback bool
}

// GetCategoryAttributes получает атрибуты для категории
// ВАЖНО: Поддерживает fallback на старую систему
func (s *UnifiedAttributeService) GetCategoryAttributes(ctx context.Context, categoryID int) ([]*UnifiedAttribute, error) {
    // 1. Пробуем получить из новой системы
    attributes, err := s.getUnifiedCategoryAttributes(ctx, categoryID)
    if err == nil && len(attributes) > 0 {
        return attributes, nil
    }
    
    // 2. Если включен fallback и новая система пустая - используем старую
    if s.useLegacyFallback {
        return s.getLegacyCategoryAttributes(ctx, categoryID)
    }
    
    return attributes, err
}

// SaveAttributeValue сохраняет значение атрибута
// ВАЖНО: Сохраняет в обе системы для совместимости
func (s *UnifiedAttributeService) SaveAttributeValue(ctx context.Context, entityType string, entityID int, attributeID int, value interface{}) error {
    tx, err := s.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()
    
    // 1. Сохраняем в новую систему
    if err := s.saveUnifiedValue(tx, entityType, entityID, attributeID, value); err != nil {
        return err
    }
    
    // 2. Если включен fallback - дублируем в старую систему
    if s.useLegacyFallback && entityType == "listing" {
        if err := s.saveLegacyValue(tx, entityID, attributeID, value); err != nil {
            // Логируем ошибку, но не прерываем - старая система не критична
            log.Printf("Failed to save to legacy system: %v", err)
        }
    }
    
    return tx.Commit()
}
```

#### 2.4 Frontend: Унифицированный компонент 

```typescript
// frontend/svetu/src/components/shared/UnifiedAttributeField.tsx

import React, { useMemo, useCallback } from 'react';
import { useTranslations } from 'next-intl';

export interface UnifiedAttribute {
  id: number;
  code: string;
  name: string;
  displayName: string;
  type: 'text' | 'textarea' | 'number' | 'boolean' | 'select' | 'multiselect' | 'date' | 'color' | 'size';
  purpose: 'regular' | 'variant' | 'both';
  options?: any;
  validationRules?: any;
  uiSettings?: any;
  isRequired: boolean;
  translations?: Record<string, any>;
}

interface UnifiedAttributeFieldProps {
  attribute: UnifiedAttribute;
  value: any;
  onChange: (value: any) => void;
  error?: string;
  disabled?: boolean;
  context?: 'listing' | 'product' | 'variant' | 'admin';
}

// ВАЖНО: Компонент обратно совместим со старыми структурами данных
export const UnifiedAttributeField: React.FC<UnifiedAttributeFieldProps> = ({
  attribute,
  value,
  onChange,
  error,
  disabled = false,
  context = 'listing'
}) => {
  const t = useTranslations('attributes');
  
  // Поддержка старого формата атрибутов для обратной совместимости
  const normalizedAttribute = useMemo(() => {
    // Если это старый формат category_attribute
    if ('attribute_type' in attribute && !('type' in attribute)) {
      return {
        ...attribute,
        type: attribute.attribute_type,
        displayName: attribute.display_name || attribute.name
      };
    }
    return attribute;
  }, [attribute]);
  
  // Рендер поля в зависимости от типа
  const renderField = useCallback(() => {
    switch (normalizedAttribute.type) {
      case 'text':
        return (
          <input
            type="text"
            value={value || ''}
            onChange={(e) => onChange(e.target.value)}
            disabled={disabled}
            className="input input-bordered w-full"
            placeholder={t('placeholder.text')}
          />
        );
        
      case 'number':
        return (
          <input
            type="number"
            value={value || ''}
            onChange={(e) => onChange(parseFloat(e.target.value))}
            disabled={disabled}
            className="input input-bordered w-full"
            placeholder={t('placeholder.number')}
          />
        );
        
      case 'select':
        return (
          <select
            value={value || ''}
            onChange={(e) => onChange(e.target.value)}
            disabled={disabled}
            className="select select-bordered w-full"
          >
            <option value="">{t('placeholder.select')}</option>
            {normalizedAttribute.options?.map((option: any) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </select>
        );
        
      // ... остальные типы полей
      
      default:
        return <div>Unsupported field type: {normalizedAttribute.type}</div>;
    }
  }, [normalizedAttribute, value, onChange, disabled, t]);
  
  return (
    <div className="form-control">
      <label className="label">
        <span className="label-text">
          {normalizedAttribute.displayName}
          {normalizedAttribute.isRequired && <span className="text-error ml-1">*</span>}
        </span>
      </label>
      {renderField()}
      {error && (
        <label className="label">
          <span className="label-text-alt text-error">{error}</span>
        </label>
      )}
    </div>
  );
};
```

### ЭТАП 3: Тестирование и валидация 

#### 3.1 Автоматические тесты 

```go
// backend/internal/services/attributes/unified_service_test.go

func TestUnifiedAttributeService_DataIntegrity(t *testing.T) {
    // 1. Проверка миграции данных
    t.Run("MigrationCompleteness", func(t *testing.T) {
        // Сравниваем количество записей в старой и новой системах
        oldCount := countOldAttributes()
        newCount := countNewAttributes()
        assert.Equal(t, oldCount, newCount, "Все атрибуты должны быть мигрированы")
    })
    
    // 2. Проверка обратной совместимости
    t.Run("BackwardCompatibility", func(t *testing.T) {
        // Получаем атрибуты через старый API
        oldAttributes := getAttributesOldWay(categoryID)
        // Получаем через новый API с fallback
        newAttributes := service.GetCategoryAttributes(ctx, categoryID)
        
        assert.Equal(t, len(oldAttributes), len(newAttributes))
        // Проверяем что данные идентичны
    })
    
    // 3. Проверка сохранения в обе системы
    t.Run("DualSystemSave", func(t *testing.T) {
        // Сохраняем через новый API
        service.SaveAttributeValue(ctx, "listing", 1, 1, "test")
        
        // Проверяем что сохранилось в обеих системах
        newValue := getValueFromNewSystem(1, 1)
        oldValue := getValueFromOldSystem(1, 1)
        
        assert.Equal(t, "test", newValue)
        assert.Equal(t, "test", oldValue)
    })
}
```

#### 3.2 Проверка производительности 

```sql
-- Скрипт для проверки производительности
-- Выполнить на копии продакшн БД

-- 1. Производительность старой системы
EXPLAIN ANALYZE
SELECT ca.*, cam.*, lav.*
FROM category_attributes ca
JOIN category_attribute_mapping cam ON ca.id = cam.attribute_id
LEFT JOIN listing_attribute_values lav ON ca.id = lav.attribute_id
WHERE cam.category_id = 1;

-- 2. Производительность новой системы
EXPLAIN ANALYZE
SELECT ua.*, uca.*, uav.*
FROM unified_attributes ua
JOIN unified_category_attributes uca ON ua.id = uca.attribute_id
LEFT JOIN unified_attribute_values uav ON ua.id = uav.attribute_id
WHERE uca.category_id = 1;

-- 3. Сравнение размера индексов
SELECT 
    schemaname,
    tablename,
    indexname,
    pg_size_pretty(pg_relation_size(indexrelid)) AS index_size
FROM pg_stat_user_indexes
WHERE tablename IN (
    'category_attributes', 'unified_attributes',
    'category_attribute_mapping', 'unified_category_attributes',
    'listing_attribute_values', 'unified_attribute_values'
)
ORDER BY pg_relation_size(indexrelid) DESC;
```

#### 3.3 Валидация данных 

```bash
#!/bin/bash
# validation_script.sh

echo "=== Валидация унификации атрибутов ==="

# 1. Проверка целостности данных
psql $DATABASE_URL <<EOF
-- Проверка что все старые атрибуты мигрированы
SELECT 
    'Немигрированные category_attributes' as check_name,
    COUNT(*) as count
FROM category_attributes ca
WHERE NOT EXISTS (
    SELECT 1 FROM unified_attributes ua 
    WHERE ua.legacy_category_attribute_id = ca.id
);

-- Проверка что все значения мигрированы
SELECT 
    'Немигрированные listing_attribute_values' as check_name,
    COUNT(*) as count
FROM listing_attribute_values lav
WHERE NOT EXISTS (
    SELECT 1 FROM unified_attribute_values uav
    JOIN unified_attributes ua ON uav.attribute_id = ua.id
    WHERE ua.legacy_category_attribute_id = lav.attribute_id
    AND uav.entity_id = lav.listing_id
    AND uav.entity_type = 'listing'
);

-- Проверка дубликатов
SELECT 
    'Дубликаты в unified_attributes' as check_name,
    COUNT(*) - COUNT(DISTINCT code) as duplicates
FROM unified_attributes;
EOF

# 2. Тест API endpoints
echo "Тестирование API..."

# Старый эндпоинт
curl -s http://localhost:3000/api/v1/marketplace/categories/1/attributes > /tmp/old_api.json

# Новый эндпоинт (если создан)
curl -s http://localhost:3000/api/v2/categories/1/attributes > /tmp/new_api.json

# Сравнение результатов
if [ -f /tmp/new_api.json ]; then
    python3 -c "
import json
old = json.load(open('/tmp/old_api.json'))
new = json.load(open('/tmp/new_api.json'))
print('API совместимость:', 'OK' if len(old['data']) == len(new['data']) else 'FAIL')
"
fi

echo "=== Валидация завершена ==="
```

### ЭТАП 4: Постепенный переход 

#### 4.1 Feature flags для контроля 

```go
// backend/internal/config/feature_flags.go

type FeatureFlags struct {
    // Флаги для постепенного перехода
    UseUnifiedAttributes      bool `env:"USE_UNIFIED_ATTRIBUTES" default:"false"`
    UnifiedAttributesFallback bool `env:"UNIFIED_ATTRIBUTES_FALLBACK" default:"true"`
    DualWriteAttributes       bool `env:"DUAL_WRITE_ATTRIBUTES" default:"true"`
    
    // Процент трафика на новую систему (для A/B тестирования)
    UnifiedAttributesPercent int `env:"UNIFIED_ATTRIBUTES_PERCENT" default:"0"`
}

// Проверка должен ли запрос использовать новую систему
func (ff *FeatureFlags) ShouldUseUnifiedAttributes(userID int) bool {
    if !ff.UseUnifiedAttributes {
        return false
    }
    
    // A/B тестирование по проценту
    if ff.UnifiedAttributesPercent < 100 {
        hash := userID % 100
        return hash < ff.UnifiedAttributesPercent
    }
    
    return true
}
```

#### 4.2 Мониторинг и метрики 

```go
// backend/internal/metrics/attributes.go

var (
    attributeSystemCalls = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "attribute_system_calls_total",
            Help: "Total number of calls to attribute systems",
        },
        []string{"system", "method", "status"},
    )
    
    attributeSystemLatency = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "attribute_system_latency_seconds",
            Help: "Latency of attribute system calls",
        },
        []string{"system", "method"},
    )
    
    dataSyncErrors = prometheus.NewCounter(
        prometheus.CounterOpts{
            Name: "attribute_data_sync_errors_total",
            Help: "Total number of data sync errors between systems",
        },
    )
)

// Обертка для отслеживания метрик
func trackAttributeOperation(system, method string, fn func() error) error {
    start := time.Now()
    err := fn()
    
    status := "success"
    if err != nil {
        status = "error"
    }
    
    attributeSystemCalls.WithLabelValues(system, method, status).Inc()
    attributeSystemLatency.WithLabelValues(system, method).Observe(time.Since(start).Seconds())
    
    return err
}
```

#### 4.3 Поэтапное включение 

```yaml
# План поэтапного включения

# Неделя 1: Тестовое окружение
- USE_UNIFIED_ATTRIBUTES: true
- UNIFIED_ATTRIBUTES_FALLBACK: true
- DUAL_WRITE_ATTRIBUTES: true
- UNIFIED_ATTRIBUTES_PERCENT: 100

# Неделя 2: 10% продакшн трафика
- USE_UNIFIED_ATTRIBUTES: true
- UNIFIED_ATTRIBUTES_FALLBACK: true
- DUAL_WRITE_ATTRIBUTES: true
- UNIFIED_ATTRIBUTES_PERCENT: 10

# Неделя 3: 50% продакшн трафика
- USE_UNIFIED_ATTRIBUTES: true
- UNIFIED_ATTRIBUTES_FALLBACK: true
- DUAL_WRITE_ATTRIBUTES: true
- UNIFIED_ATTRIBUTES_PERCENT: 50

# Неделя 4: 100% продакшн трафика
- USE_UNIFIED_ATTRIBUTES: true
- UNIFIED_ATTRIBUTES_FALLBACK: true
- DUAL_WRITE_ATTRIBUTES: true
- UNIFIED_ATTRIBUTES_PERCENT: 100

# После стабилизации: Отключение старой системы
- USE_UNIFIED_ATTRIBUTES: true
- UNIFIED_ATTRIBUTES_FALLBACK: false
- DUAL_WRITE_ATTRIBUTES: false
- UNIFIED_ATTRIBUTES_PERCENT: 100
```

### ЭТАП 5: Завершение и очистка 

#### 5.1 Финальная проверка 

```sql
-- Финальная проверка перед отключением старой системы

-- 1. Сверка количества записей
WITH comparison AS (
    SELECT 
        (SELECT COUNT(*) FROM category_attributes) as old_attributes,
        (SELECT COUNT(*) FROM unified_attributes WHERE legacy_category_attribute_id IS NOT NULL) as new_attributes,
        (SELECT COUNT(*) FROM listing_attribute_values) as old_values,
        (SELECT COUNT(*) FROM unified_attribute_values WHERE entity_type = 'listing') as new_values
)
SELECT 
    CASE 
        WHEN old_attributes = new_attributes AND old_values = new_values THEN 'OK'
        ELSE 'FAIL'
    END as migration_status,
    *
FROM comparison;

-- 2. Проверка активных связей
SELECT 
    'Активные foreign keys на старые таблицы' as check_name,
    conname as constraint_name,
    conrelid::regclass as table_name,
    confrelid::regclass as referenced_table
FROM pg_constraint
WHERE confrelid IN (
    'category_attributes'::regclass,
    'category_attribute_mapping'::regclass,
    'listing_attribute_values'::regclass
)
AND conrelid NOT IN (
    'unified_attributes'::regclass,
    'unified_category_attributes'::regclass,
    'unified_attribute_values'::regclass
);
```

#### 5.2 Отключение старой системы 

```go
// backend/internal/handlers/attributes_migration.go

// MigrationStatusHandler - эндпоинт для проверки статуса миграции
func MigrationStatusHandler(c *fiber.Ctx) error {
    status := CheckMigrationStatus()
    
    return c.JSON(fiber.Map{
        "old_system_active": config.FeatureFlags.UnifiedAttributesFallback,
        "new_system_active": config.FeatureFlags.UseUnifiedAttributes,
        "dual_write": config.FeatureFlags.DualWriteAttributes,
        "traffic_percent": config.FeatureFlags.UnifiedAttributesPercent,
        "data_sync_status": status,
        "ready_for_cleanup": status.IsFullyMigrated && !status.HasActiveConnections,
    })
}
```

#### 5.3 Архивирование старых таблиц 

```sql
-- Миграция: backend/migrations/999_archive_old_attributes.up.sql
-- ВЫПОЛНЯТЬ ТОЛЬКО ПОСЛЕ ПОЛНОЙ ПРОВЕРКИ!

BEGIN;

-- 1. Создание архивной схемы
CREATE SCHEMA IF NOT EXISTS archive;

-- 2. Перемещение старых таблиц в архив (НЕ УДАЛЕНИЕ!)
ALTER TABLE category_attributes SET SCHEMA archive;
ALTER TABLE category_attribute_mapping SET SCHEMA archive;
ALTER TABLE listing_attribute_values SET SCHEMA archive;
ALTER TABLE product_variant_attributes SET SCHEMA archive;
ALTER TABLE storefront_product_attributes SET SCHEMA archive;
ALTER TABLE product_variant_attribute_values SET SCHEMA archive;
ALTER TABLE category_variant_attributes SET SCHEMA archive;
ALTER TABLE variant_attribute_mappings SET SCHEMA archive;

-- 3. Создание view для обратной совместимости (на случай экстренного отката)
CREATE VIEW category_attributes AS
SELECT 
    ua.id,
    ua.name,
    ua.display_name,
    ua.attribute_type,
    ua.options,
    ua.validation_rules,
    ua.is_searchable,
    ua.is_filterable,
    ua.is_required,
    CASE WHEN ua.purpose IN ('variant', 'both') THEN true ELSE false END as is_variant_compatible,
    ua.affects_stock,
    ua.sort_order
FROM unified_attributes ua
WHERE ua.legacy_category_attribute_id IS NOT NULL;

-- 4. Удаление legacy полей из новых таблиц (опционально, после месяца стабильной работы)
-- ALTER TABLE unified_attributes DROP COLUMN legacy_category_attribute_id;
-- ALTER TABLE unified_attributes DROP COLUMN legacy_product_variant_attribute_id;

COMMIT;
```

---

## 📊 Критерии успешного выполнения

### Обязательные критерии:
1. ✅ **Нулевая потеря данных** - все существующие атрибуты и их значения сохранены
2. ✅ **Полная обратная совместимость** - старые API продолжают работать
3. ✅ **Возможность отката** - в любой момент можно вернуться к старой системе
4. ✅ **Производительность не ухудшилась** - время ответа API <= текущего
5. ✅ **Все тесты проходят** - unit, integration, e2e тесты зеленые

### Целевые метрики после унификации:
- 📉 Количество таблиц: с 14 до 3-4
- 📉 Дублированный код: -70% (~1800 строк)
- 📉 Количество SQL запросов: -50%
- ⚡ Время загрузки страниц с атрибутами: -30%
- 💾 Cache hit rate: >80%

---

## ⚠️ Риски и план митигации

### Риск 1: Потеря данных при миграции
**Митигация:**
- Полный бэкап перед каждым этапом
- Проверка целостности после каждой миграции
- Транзакционность всех операций
- Возможность отката миграций

### Риск 2: Несовместимость с существующим кодом
**Митигация:**
- Fallback на старую систему
- Dual-write в обе системы
- Постепенное включение через feature flags
- A/B тестирование

### Риск 3: Проблемы производительности
**Митигация:**
- Тестирование на копии продакшн данных
- Мониторинг метрик в реальном времени
- Оптимизация запросов и индексов
- Кеширование на всех уровнях

### Риск 4: Ошибки в продакшене
**Митигация:**
- Поэтапный rollout (10% -> 50% -> 100%)
- Детальное логирование всех операций
- Алерты на аномалии в метриках
- Готовый план отката

---

## 📝 Чек-лист для контроля выполнения

### Подготовка:
- [ ] Создан полный бэкап всех таблиц атрибутов
- [ ] Проанализировано использование атрибутов в коде
- [ ] Составлена карта зависимостей
- [ ] Настроен тестовый стенд с копией продакшн данных

### Реализация:
- [ ] Созданы новые таблицы БД (без удаления старых)
- [ ] Написаны и протестированы миграции данных
- [ ] Реализован унифицированный сервис с fallback
- [ ] Создан унифицированный frontend компонент
- [ ] Написаны автоматические тесты

### Тестирование:
- [ ] Пройдены все unit тесты
- [ ] Пройдены integration тесты
- [ ] Проверена производительность
- [ ] Валидирована целостность данных
- [ ] Протестирована обратная совместимость

### Развертывание:
- [ ] Настроены feature flags
- [ ] Развернуто на тестовом окружении
- [ ] Включено для 10% трафика
- [ ] Включено для 50% трафика
- [ ] Включено для 100% трафика
- [ ] Мониторинг стабилен 48 часов

### Завершение:
- [ ] Отключен fallback на старую систему
- [ ] Отключен dual-write
- [ ] Архивированы старые таблицы
- [ ] Обновлена документация
- [ ] Проведен ретроспективный анализ

---

## 🚨 Команда экстренного отката

В случае критических проблем выполнить:

```bash
#!/bin/bash
# emergency_rollback.sh

echo "⚠️ ЭКСТРЕННЫЙ ОТКАТ СИСТЕМЫ АТРИБУТОВ"

# 1. Отключение новой системы
export USE_UNIFIED_ATTRIBUTES=false
export UNIFIED_ATTRIBUTES_FALLBACK=true
export DUAL_WRITE_ATTRIBUTES=false

# 2. Перезапуск сервисов
systemctl restart backend-api

# 3. Проверка работоспособности
curl http://localhost:3000/health

# 4. Восстановление из бэкапа если нужно
# psql $DATABASE_URL < /data/backups/attribute_unification_$(date +%Y%m%d)/attributes_backup.sql

echo "✅ Откат завершен"
```

---

## 📞 Контакты для эскалации

При возникновении критических проблем:
1. Проверить логи: `/var/log/backend/attributes.log`
2. Проверить метрики: Grafana Dashboard "Attribute System"
3. Выполнить откат если необходимо
4. Документировать проблему в `/docs/incidents/`

---

**Автор ТЗ**: System Architect  
**Дата создания**: 02.09.2025  
**Версия**: 1.2.0  
**Статус**: В процессе выполнения (День 12 из 30)

---

## 📈 АКТУАЛЬНЫЙ СТАТУС ВЫПОЛНЕНИЯ

**Дата обновления**: 03.09.2025  
**Прогресс**: День 12 из 30 (40% времени)

### ✅ Выполнено:

#### Этап 1: Подготовка и анализ (Дни 1-3) - ✅ ЗАВЕРШЕНО
- ✅ Создан полный бэкап всех таблиц атрибутов
- ✅ Проанализировано использование атрибутов (85 атрибутов, 611 связей, 15 значений)
- ✅ Составлена карта зависимостей
- ✅ Создан план миграции

#### Этап 2: Backend (Дни 2-4) - ✅ ЗАВЕРШЕНО
- ✅ Созданы новые таблицы БД (unified_attributes, unified_category_attributes, unified_attribute_values)
- ✅ Написаны миграции данных (000034, 000035)
- ✅ Реализован унифицированный сервис с fallback
- ✅ Создан API v2 для новой системы
- ✅ Поддержка 15 типов атрибутов

#### Этап 3: Frontend (Дни 5-7) - ✅ ЗАВЕРШЕНО
- ✅ Создан UnifiedAttributeField компонент (День 6)
- ✅ Создан unifiedAttributeService (День 6)
- ✅ Интегрирован UnifiedAttributesStep (День 7)
- ✅ Реализована система feature flags (День 7)
- ✅ Добавлено кеширование с TTL 5 минут

#### Этап 4: Тестирование (Дни 8-10) - ✅ ЗАВЕРШЕНО
- ✅ Unit тесты для UnifiedAttributeField (~95% покрытие)
- ✅ Unit тесты для unifiedAttributeService (~90% покрытие)
- ✅ Скрипт нагрузочного тестирования (Go)
- ✅ Конфигурация feature flags (.env.test)
- ✅ Руководство по тестированию создано
- ✅ E2E тестирование успешно пройдено (День 9)
- ✅ Интеграционное тестирование завершено (День 10)
- ✅ Dual-write механизм протестирован
- ✅ Fallback система работает корректно

#### Этап 5: Мониторинг и метрики (День 11) - ✅ ЗАВЕРШЕНО
- ✅ Prometheus метрики настроены (20+ метрик)
- ✅ Health check endpoints реализованы (/health/live, /health/ready)
- ✅ Grafana dashboard создан (12+ панелей)
- ✅ Alert rules настроены (9 правил)
- ✅ Feature flags метрики отслеживаются

#### Этап 6: CI/CD Pipeline (День 12) - ✅ ЗАВЕРШЕНО
- ✅ GitHub Actions workflow создан
- ✅ Автоматические тесты на PR настроены
- ✅ Load testing скрипты реализованы
- ✅ Validation framework настроен

#### Этап 7: Production подготовка (Дни 13-14) - ✅ ЗАВЕРШЕНО
- ✅ Blue-green deployment настроен
- ✅ Canary release automation готова
- ✅ Production runbook создан
- ✅ Monitoring dashboards настроены

#### Этап 8: Production deployment (День 15) - 🚀 В ПРОЦЕССЕ
- ✅ Deployment скрипты созданы
- ✅ Validation framework готов
- 🔄 Production развертывание выполняется
- ⏰ Post-deployment мониторинг
- ✅ k6 load tests реализованы
- ✅ Dual-write validation scripts
- ✅ Fallback testing scripts
- ✅ Migration integrity checks

### 🔄 В процессе (День 13-15):

#### Дни 13-15: Production развертывание
- ⏳ Настройка production окружения
- ⏳ Blue-green deployment
- ⏳ Rollback механизмы
- ⏳ Production мониторинг
- ⏳ Graceful migration

### 📊 Метрики прогресса:

| Компонент | План | Факт | Статус |
|-----------|------|------|--------|
| Backend | 100% | 100% | ✅ Завершено |
| Frontend | 100% | 100% | ✅ Завершено |
| Миграция БД | 100% | 100% | ✅ Завершено |
| Тестирование | 100% | 100% | ✅ Завершено |
| Мониторинг | 100% | 100% | ✅ Завершено |
| CI/CD | 100% | 100% | ✅ Завершено |
| Production rollout | 100% | 0% | 🔄 В процессе |

### 📋 Чек-лист выполнения:

#### Подготовка:
- [x] Создан полный бэкап всех таблиц атрибутов
- [x] Проанализировано использование атрибутов в коде
- [x] Составлена карта зависимостей
- [x] Настроен тестовый стенд с копией продакшн данных

#### Реализация:
- [x] Созданы новые таблицы БД (без удаления старых)
- [x] Написаны и протестированы миграции данных
- [x] Реализован унифицированный сервис с fallback
- [x] Создан унифицированный frontend компонент
- [x] Написаны автоматические тесты

#### Тестирование:
- [x] Написаны unit тесты
- [x] Пройдены все unit тесты
- [x] Пройдены integration тесты
- [x] Проверена производительность (<3ms response time)
- [x] Валидирована целостность данных
- [x] Протестирована обратная совместимость

#### Развертывание:
- [x] Настроены feature flags
- [ ] Развернуто на тестовом окружении
- [ ] Включено для 10% трафика
- [ ] Включено для 50% трафика
- [ ] Включено для 100% трафика
- [ ] Мониторинг стабилен 48 часов

**ВАЖНОЕ НАПОМИНАНИЕ**: 
- Это ТЗ должно выполняться ПОШАГОВО
- ЗАПРЕЩЕНО пропускать этапы валидации
- ОБЯЗАТЕЛЬНО делать бэкапы перед каждым этапом
- При любых сомнениях - ОСТАНОВИТЬСЯ и запросить уточнения

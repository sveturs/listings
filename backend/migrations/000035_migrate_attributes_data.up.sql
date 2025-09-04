-- Миграция данных из старых таблиц атрибутов в унифицированную систему
-- ВАЖНО: Транзакция для атомарности
-- Дата: 02.09.2025

BEGIN;

-- =====================================================
-- 1. МИГРАЦИЯ АТРИБУТОВ ИЗ category_attributes
-- =====================================================
INSERT INTO unified_attributes (
    code, name, display_name, attribute_type, purpose,
    options, validation_rules, ui_settings,
    is_searchable, is_filterable, is_required,
    affects_stock, affects_price,
    sort_order, legacy_category_attribute_id
)
SELECT 
    LOWER(REPLACE(REPLACE(name, ' ', '_'), '-', '_')) as code,
    name,
    COALESCE(display_name, name),
    attribute_type,
    CASE 
        WHEN is_variant_compatible = true THEN 'both'
        ELSE 'regular'
    END as purpose,
    COALESCE(options, '{}'),
    COALESCE(validation_rules, '{}'),
    '{}' as ui_settings, -- Поле ui_settings отсутствует в старой таблице
    COALESCE(is_searchable, false),
    COALESCE(is_filterable, false),
    COALESCE(is_required, false),
    COALESCE(affects_stock, false),
    false as affects_price, -- Добавим позже если нужно
    COALESCE(sort_order, 0),
    id as legacy_category_attribute_id
FROM category_attributes
ON CONFLICT (code) DO UPDATE SET
    legacy_category_attribute_id = EXCLUDED.legacy_category_attribute_id,
    purpose = CASE 
        WHEN unified_attributes.purpose = 'variant' THEN 'both'
        WHEN EXCLUDED.purpose = 'variant' THEN 'both'
        ELSE unified_attributes.purpose
    END;

-- =====================================================
-- 2. МИГРАЦИЯ АТРИБУТОВ ИЗ product_variant_attributes
-- =====================================================
INSERT INTO unified_attributes (
    code, name, display_name, attribute_type, purpose,
    options, affects_stock, affects_price,
    is_filterable, is_searchable,
    legacy_product_variant_attribute_id
)
SELECT 
    LOWER(REPLACE(REPLACE(name, ' ', '_'), '-', '_')) as code,
    name,
    COALESCE(display_name, name) as display_name,
    type as attribute_type,
    'variant' as purpose,
    '{}' as options, -- Поле options отсутствует в таблице
    COALESCE(affects_stock, false),
    false as affects_price, -- Поле affects_price отсутствует в таблице
    true as is_filterable, -- Вариативные атрибуты обычно фильтруемые
    true as is_searchable,
    id as legacy_product_variant_attribute_id
FROM product_variant_attributes
ON CONFLICT (code) DO UPDATE SET
    purpose = 'both', -- Если атрибут уже есть, делаем его универсальным
    affects_stock = COALESCE(unified_attributes.affects_stock, EXCLUDED.affects_stock),
    affects_price = COALESCE(unified_attributes.affects_price, EXCLUDED.affects_price),
    legacy_product_variant_attribute_id = EXCLUDED.legacy_product_variant_attribute_id;

-- =====================================================
-- 3. МИГРАЦИЯ СВЯЗЕЙ КАТЕГОРИЙ С АТРИБУТАМИ
-- =====================================================
INSERT INTO unified_category_attributes (
    category_id, attribute_id, is_enabled, is_required, sort_order
)
SELECT 
    cam.category_id,
    ua.id as attribute_id,
    COALESCE(cam.is_enabled, true),
    COALESCE(cam.is_required, false),
    COALESCE(cam.sort_order, 0)
FROM category_attribute_mapping cam
JOIN unified_attributes ua ON ua.legacy_category_attribute_id = cam.attribute_id
ON CONFLICT (category_id, attribute_id) DO NOTHING;

-- =====================================================
-- 4. МИГРАЦИЯ ЗНАЧЕНИЙ АТРИБУТОВ ОБЪЯВЛЕНИЙ
-- =====================================================
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
JOIN unified_attributes ua ON ua.legacy_category_attribute_id = lav.attribute_id
ON CONFLICT (entity_type, entity_id, attribute_id) DO NOTHING;

-- =====================================================
-- 5. МИГРАЦИЯ ЗНАЧЕНИЙ ВАРИАТИВНЫХ АТРИБУТОВ
-- =====================================================
INSERT INTO unified_attribute_values (
    entity_type, entity_id, attribute_id,
    text_value, numeric_value, boolean_value, json_value
)
SELECT 
    'product_variant' as entity_type,
    pvav.variant_id as entity_id,
    ua.id as attribute_id,
    pvav.text_value,
    pvav.numeric_value,
    pvav.boolean_value,
    pvav.json_value
FROM product_variant_attribute_values pvav
JOIN unified_attributes ua ON ua.legacy_product_variant_attribute_id = pvav.attribute_id
ON CONFLICT (entity_type, entity_id, attribute_id) DO NOTHING;

-- =====================================================
-- 6. МИГРАЦИЯ АТРИБУТОВ ИЗ category_variant_attributes (устаревшая система)
-- =====================================================
-- ПРОПУСКАЕМ category_variant_attributes так как структура таблицы отличается
-- и в ней нет необходимых полей (name, type, options)
ON CONFLICT (code) DO UPDATE SET
    purpose = 'both'; -- Если атрибут уже есть, делаем его универсальным

-- =====================================================
-- 7. ПРОВЕРКА МИГРАЦИИ
-- =====================================================
DO $$
DECLARE
    old_count INTEGER;
    new_count INTEGER;
    missing_count INTEGER;
BEGIN
    -- Проверяем атрибуты из category_attributes
    SELECT COUNT(*) INTO old_count FROM category_attributes;
    SELECT COUNT(*) INTO new_count FROM unified_attributes WHERE legacy_category_attribute_id IS NOT NULL;
    
    IF old_count != new_count THEN
        RAISE WARNING 'Миграция category_attributes: старых %, новых %', old_count, new_count;
    ELSE
        RAISE NOTICE 'category_attributes мигрированы успешно: %', old_count;
    END IF;
    
    -- Проверяем значения
    SELECT COUNT(*) INTO old_count FROM listing_attribute_values;
    SELECT COUNT(*) INTO new_count FROM unified_attribute_values WHERE entity_type = 'listing';
    
    IF old_count != new_count THEN
        RAISE WARNING 'Миграция listing_attribute_values: старых %, новых %', old_count, new_count;
    ELSE
        RAISE NOTICE 'listing_attribute_values мигрированы успешно: %', old_count;
    END IF;
    
    -- Проверяем маппинги категорий
    SELECT COUNT(*) INTO old_count FROM category_attribute_mapping;
    SELECT COUNT(*) INTO new_count FROM unified_category_attributes;
    
    RAISE NOTICE 'category_attribute_mapping: старых %, новых %', old_count, new_count;
END $$;

COMMIT;

-- =====================================================
-- 8. СТАТИСТИКА ПОСЛЕ МИГРАЦИИ
-- =====================================================
SELECT 
    'Статистика миграции' as info,
    (SELECT COUNT(*) FROM unified_attributes) as total_attributes,
    (SELECT COUNT(*) FROM unified_attributes WHERE purpose = 'regular') as regular_attributes,
    (SELECT COUNT(*) FROM unified_attributes WHERE purpose = 'variant') as variant_attributes,
    (SELECT COUNT(*) FROM unified_attributes WHERE purpose = 'both') as universal_attributes,
    (SELECT COUNT(*) FROM unified_category_attributes) as category_mappings,
    (SELECT COUNT(*) FROM unified_attribute_values) as total_values;
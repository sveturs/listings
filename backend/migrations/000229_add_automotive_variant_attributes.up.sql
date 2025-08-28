-- Добавление вариативных атрибутов для автомобилей в таблицу product_variant_attributes
-- Сначала проверяем и добавляем атрибуты, если их нет
INSERT INTO product_variant_attributes (name, display_name, affects_stock, created_at, updated_at)
SELECT 'engine_type', 'Engine Type', false, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM product_variant_attributes WHERE name = 'engine_type');

INSERT INTO product_variant_attributes (name, display_name, affects_stock, created_at, updated_at)
SELECT 'trim_level', 'Trim Level', false, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM product_variant_attributes WHERE name = 'trim_level');

INSERT INTO product_variant_attributes (name, display_name, affects_stock, created_at, updated_at)
SELECT 'transmission', 'Transmission Type', false, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM product_variant_attributes WHERE name = 'transmission');

INSERT INTO product_variant_attributes (name, display_name, affects_stock, created_at, updated_at)
SELECT 'fuel_type', 'Fuel Type', false, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM product_variant_attributes WHERE name = 'fuel_type');

-- Переводы будут добавлены через product_variant_attributes, так как translations требует integer entity_id
-- Пропускаем вставку переводов в таблицу translations

-- Добавление атрибутов в категорию автомобилей как доступные для вариантов
-- Сначала найдем ID категории Cars
DO $$
DECLARE
    cars_category_id INTEGER;
BEGIN
    SELECT id INTO cars_category_id FROM marketplace_categories WHERE slug = 'cars' LIMIT 1;
    
    IF cars_category_id IS NOT NULL THEN
        -- Добавляем вариативные атрибуты для категории Cars
        INSERT INTO category_variant_attributes (category_id, variant_attribute_name, sort_order, is_required, created_at)
        VALUES
            (cars_category_id, 'engine_type', 1, false, NOW()),
            (cars_category_id, 'trim_level', 2, false, NOW()),
            (cars_category_id, 'transmission', 3, false, NOW()),
            (cars_category_id, 'fuel_type', 4, false, NOW())
        ON CONFLICT DO NOTHING;
    END IF;
END$$;

-- Создание атрибутов категории для автомобильных вариантов (если их еще нет)
INSERT INTO category_attributes (name, display_name, attribute_type, is_required, is_filterable, is_variant_compatible, data_source, sort_order, created_at)
SELECT 'engine_type', 'Engine Type', 'text', false, true, true, 'manual', 200, NOW()
WHERE NOT EXISTS (SELECT 1 FROM category_attributes WHERE name = 'engine_type');

INSERT INTO category_attributes (name, display_name, attribute_type, is_required, is_filterable, is_variant_compatible, data_source, sort_order, created_at)
SELECT 'trim_level', 'Trim Level', 'text', false, true, true, 'manual', 201, NOW()
WHERE NOT EXISTS (SELECT 1 FROM category_attributes WHERE name = 'trim_level');

INSERT INTO category_attributes (name, display_name, attribute_type, is_required, is_filterable, is_variant_compatible, data_source, sort_order, created_at)
SELECT 'transmission', 'Transmission', 'select', false, true, true, 'manual', 202, NOW()
WHERE NOT EXISTS (SELECT 1 FROM category_attributes WHERE name = 'transmission');

INSERT INTO category_attributes (name, display_name, attribute_type, is_required, is_filterable, is_variant_compatible, data_source, sort_order, created_at)
SELECT 'fuel_type', 'Fuel Type', 'select', false, true, true, 'manual', 203, NOW()
WHERE NOT EXISTS (SELECT 1 FROM category_attributes WHERE name = 'fuel_type');

-- Обновляем существующие атрибуты, если они уже есть
UPDATE category_attributes SET is_variant_compatible = true, is_filterable = true 
WHERE name IN ('engine_type', 'trim_level', 'transmission', 'fuel_type');

-- Добавление предопределенных значений для transmission
UPDATE category_attributes
SET options = jsonb_build_array(
    jsonb_build_object('value', 'manual', 'label', 'Manual'),
    jsonb_build_object('value', 'automatic', 'label', 'Automatic'),
    jsonb_build_object('value', 'semi-automatic', 'label', 'Semi-Automatic'),
    jsonb_build_object('value', 'cvt', 'label', 'CVT'),
    jsonb_build_object('value', 'dsg', 'label', 'DSG')
)
WHERE name = 'transmission';

-- Добавление предопределенных значений для fuel_type
UPDATE category_attributes
SET options = jsonb_build_array(
    jsonb_build_object('value', 'petrol', 'label', 'Petrol'),
    jsonb_build_object('value', 'diesel', 'label', 'Diesel'),
    jsonb_build_object('value', 'electric', 'label', 'Electric'),
    jsonb_build_object('value', 'hybrid', 'label', 'Hybrid'),
    jsonb_build_object('value', 'plug-in-hybrid', 'label', 'Plug-in Hybrid'),
    jsonb_build_object('value', 'lpg', 'label', 'LPG'),
    jsonb_build_object('value', 'cng', 'label', 'CNG')
)
WHERE name = 'fuel_type';
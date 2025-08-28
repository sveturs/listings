-- Добавляем атрибуты для хранения ID марки и модели автомобиля
INSERT INTO category_attributes (name, display_name, attribute_type, is_required, is_searchable, sort_order, options)
SELECT 'car_make_id', 'Car Make ID', 'number', false, true, 99, '{}'
WHERE NOT EXISTS (SELECT 1 FROM category_attributes WHERE name = 'car_make_id');

INSERT INTO category_attributes (name, display_name, attribute_type, is_required, is_searchable, sort_order, options)
SELECT 'car_model_id', 'Car Model ID', 'number', false, true, 100, '{}'
WHERE NOT EXISTS (SELECT 1 FROM category_attributes WHERE name = 'car_model_id');

-- Связываем новые атрибуты с автомобильными категориями
INSERT INTO category_attribute_mapping (category_id, attribute_id, is_enabled, is_required, sort_order)
SELECT 
    c.id as category_id,
    a.id as attribute_id,
    true as is_enabled,
    false as is_required,
    CASE 
        WHEN a.name = 'car_make_id' THEN 99
        WHEN a.name = 'car_model_id' THEN 100
    END as sort_order
FROM marketplace_categories c
CROSS JOIN category_attributes a
WHERE 
    -- Применяем к автомобильным категориям
    c.id IN (
        2000, 2100, 2200, 2300, 2400, 2500, 2600, 2700, 2800,
        10100, 10101, 10102, 10103, 10104, 10110, 10111, 10112, 10113,
        10170, 10171, 10172, 10173, 10174, 10175, 10176, 10177, 10178,
        10179, 10180, 10181, 10182, 10183, 10184, 10185, 10186, 10187
    )
    AND a.name IN ('car_make_id', 'car_model_id')
ON CONFLICT (category_id, attribute_id) DO NOTHING;

-- Добавляем переводы для новых атрибутов
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text)
VALUES 
    ('category_attribute', (SELECT id FROM category_attributes WHERE name = 'car_make_id'), 'en', 'display_name', 'Car Make ID'),
    ('category_attribute', (SELECT id FROM category_attributes WHERE name = 'car_make_id'), 'ru', 'display_name', 'ID марки автомобиля'),
    ('category_attribute', (SELECT id FROM category_attributes WHERE name = 'car_make_id'), 'sr', 'display_name', 'ID marke automobila'),
    ('category_attribute', (SELECT id FROM category_attributes WHERE name = 'car_model_id'), 'en', 'display_name', 'Car Model ID'),
    ('category_attribute', (SELECT id FROM category_attributes WHERE name = 'car_model_id'), 'ru', 'display_name', 'ID модели автомобиля'),
    ('category_attribute', (SELECT id FROM category_attributes WHERE name = 'car_model_id'), 'sr', 'display_name', 'ID modela automobila')
ON CONFLICT (entity_type, entity_id, language, field_name) DO NOTHING;
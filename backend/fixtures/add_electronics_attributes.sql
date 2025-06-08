-- Создаем атрибуты для категории "Электроника"

-- 1. Бренд (текстовый)
INSERT INTO category_attributes (name, display_name, attribute_type, options, is_searchable, is_filterable, is_required, sort_order)
VALUES ('brand', 'Бренд', 'text', '{}', true, true, true, 1)
ON CONFLICT (name) DO NOTHING;

-- 2. Модель (текстовый)
INSERT INTO category_attributes (name, display_name, attribute_type, options, is_searchable, is_filterable, is_required, sort_order)
VALUES ('model', 'Модель', 'text', '{}', true, true, false, 2)
ON CONFLICT (name) DO NOTHING;

-- 3. Состояние (select)
INSERT INTO category_attributes (name, display_name, attribute_type, options, is_searchable, is_filterable, is_required, sort_order)
VALUES ('device_condition', 'Состояние устройства', 'select', 
    '{"values": ["Новое", "Как новое", "Отличное", "Хорошее", "Удовлетворительное"]}', 
    false, true, true, 3)
ON CONFLICT (name) DO NOTHING;

-- 4. Память/Хранилище (select)
INSERT INTO category_attributes (name, display_name, attribute_type, options, is_searchable, is_filterable, is_required, sort_order)
VALUES ('storage', 'Память', 'select', 
    '{"values": ["16GB", "32GB", "64GB", "128GB", "256GB", "512GB", "1TB", "2TB", "Другое"]}', 
    false, true, false, 4)
ON CONFLICT (name) DO NOTHING;

-- 5. Цвет (select)
INSERT INTO category_attributes (name, display_name, attribute_type, options, is_searchable, is_filterable, is_required, sort_order)
VALUES ('color', 'Цвет', 'select', 
    '{"values": ["Черный", "Белый", "Серый", "Синий", "Красный", "Зеленый", "Золотой", "Серебристый", "Другой"]}', 
    false, true, false, 5)
ON CONFLICT (name) DO NOTHING;

-- 6. Гарантия (boolean)
INSERT INTO category_attributes (name, display_name, attribute_type, options, is_searchable, is_filterable, is_required, sort_order)
VALUES ('warranty', 'Гарантия', 'boolean', '{}', false, true, false, 6)
ON CONFLICT (name) DO NOTHING;

-- 7. Комплектация (multiselect)
INSERT INTO category_attributes (name, display_name, attribute_type, options, is_searchable, is_filterable, is_required, sort_order)
VALUES ('accessories', 'Комплектация', 'multiselect', 
    '{"values": ["Коробка", "Зарядное устройство", "Кабель", "Наушники", "Чехол", "Защитное стекло", "Документы"], "multiselect": true}', 
    false, false, false, 7)
ON CONFLICT (name) DO NOTHING;

-- 8. Размер экрана (number)
INSERT INTO category_attributes (name, display_name, attribute_type, options, is_searchable, is_filterable, is_required, sort_order)
VALUES ('screen_size', 'Размер экрана (дюймы)', 'number', 
    '{"min": 1, "max": 100, "step": 0.1}', 
    false, true, false, 8)
ON CONFLICT (name) DO NOTHING;

-- Теперь привязываем атрибуты к категории "Электроника" (ID = 1)
INSERT INTO category_attribute_mapping (category_id, attribute_id, is_enabled, is_required, sort_order)
SELECT 1, ca.id, true, 
    CASE 
        WHEN ca.name IN ('brand', 'device_condition') THEN true 
        ELSE false 
    END,
    ca.sort_order
FROM category_attributes ca
WHERE ca.name IN ('brand', 'model', 'device_condition', 'storage', 'color', 'warranty', 'accessories', 'screen_size')
ON CONFLICT (category_id, attribute_id) DO UPDATE
SET is_enabled = EXCLUDED.is_enabled,
    is_required = EXCLUDED.is_required,
    sort_order = EXCLUDED.sort_order;

-- Добавляем переводы для атрибутов
INSERT INTO attribute_translations (attribute_id, language_code, display_name)
SELECT ca.id, 'ru', 
    CASE ca.name
        WHEN 'brand' THEN 'Бренд'
        WHEN 'model' THEN 'Модель'
        WHEN 'device_condition' THEN 'Состояние устройства'
        WHEN 'storage' THEN 'Память'
        WHEN 'color' THEN 'Цвет'
        WHEN 'warranty' THEN 'Гарантия'
        WHEN 'accessories' THEN 'Комплектация'
        WHEN 'screen_size' THEN 'Размер экрана (дюймы)'
    END
FROM category_attributes ca
WHERE ca.name IN ('brand', 'model', 'device_condition', 'storage', 'color', 'warranty', 'accessories', 'screen_size')
ON CONFLICT (attribute_id, language_code) DO UPDATE
SET display_name = EXCLUDED.display_name;

-- Английские переводы
INSERT INTO attribute_translations (attribute_id, language_code, display_name)
SELECT ca.id, 'en', 
    CASE ca.name
        WHEN 'brand' THEN 'Brand'
        WHEN 'model' THEN 'Model'
        WHEN 'device_condition' THEN 'Device Condition'
        WHEN 'storage' THEN 'Storage'
        WHEN 'color' THEN 'Color'
        WHEN 'warranty' THEN 'Warranty'
        WHEN 'accessories' THEN 'Accessories'
        WHEN 'screen_size' THEN 'Screen Size (inches)'
    END
FROM category_attributes ca
WHERE ca.name IN ('brand', 'model', 'device_condition', 'storage', 'color', 'warranty', 'accessories', 'screen_size')
ON CONFLICT (attribute_id, language_code) DO UPDATE
SET display_name = EXCLUDED.display_name;

-- Добавляем переводы для значений опций
INSERT INTO attribute_option_translations (attribute_id, option_value, language_code, translated_value)
SELECT ca.id, v.value, 'ru', v.value
FROM category_attributes ca,
LATERAL (
    SELECT unnest(ARRAY['Новое', 'Как новое', 'Отличное', 'Хорошее', 'Удовлетворительное']) as value
    WHERE ca.name = 'device_condition'
    UNION ALL
    SELECT unnest(ARRAY['16GB', '32GB', '64GB', '128GB', '256GB', '512GB', '1TB', '2TB', 'Другое']) as value
    WHERE ca.name = 'storage'
    UNION ALL
    SELECT unnest(ARRAY['Черный', 'Белый', 'Серый', 'Синий', 'Красный', 'Зеленый', 'Золотой', 'Серебристый', 'Другой']) as value
    WHERE ca.name = 'color'
    UNION ALL
    SELECT unnest(ARRAY['Коробка', 'Зарядное устройство', 'Кабель', 'Наушники', 'Чехол', 'Защитное стекло', 'Документы']) as value
    WHERE ca.name = 'accessories'
) v
WHERE ca.name IN ('device_condition', 'storage', 'color', 'accessories')
ON CONFLICT (attribute_id, option_value, language_code) DO UPDATE
SET translated_value = EXCLUDED.translated_value;
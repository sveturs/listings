-- Сначала создадим недостающие атрибуты
INSERT INTO category_attributes (name, display_name, attribute_type, options, is_searchable, is_filterable, is_required, sort_order)
VALUES 
    ('model', 'Модель', 'text', '{}', true, true, false, 2),
    ('device_condition', 'Состояние устройства', 'select', 
        '{"values": ["Новое", "Как новое", "Отличное", "Хорошее", "Удовлетворительное"]}', 
        false, true, true, 3),
    ('storage', 'Память', 'select', 
        '{"values": ["16GB", "32GB", "64GB", "128GB", "256GB", "512GB", "1TB", "2TB", "Другое"]}', 
        false, true, false, 4),
    ('warranty', 'Гарантия', 'boolean', '{}', false, true, false, 6),
    ('accessories', 'Комплектация', 'multiselect', 
        '{"values": ["Коробка", "Зарядное устройство", "Кабель", "Наушники", "Чехол", "Защитное стекло", "Документы"], "multiselect": true}', 
        false, false, false, 7)
ON CONFLICT (name) DO NOTHING;

-- Обновим тип атрибута brand на text (сейчас он select)
UPDATE category_attributes SET attribute_type = 'text' WHERE name = 'brand';

-- Обновим опции для color
UPDATE category_attributes 
SET options = '{"values": ["Черный", "Белый", "Серый", "Синий", "Красный", "Зеленый", "Золотой", "Серебристый", "Другой"]}'
WHERE name = 'color';

-- Обновим опции для screen_size
UPDATE category_attributes 
SET options = '{"min": 1, "max": 100, "step": 0.1}'
WHERE name = 'screen_size';

-- Теперь создаем маппинги
DO $$
DECLARE
    cat_id INTEGER := 1;
BEGIN
    -- Добавляем маппинги для всех атрибутов
    INSERT INTO category_attribute_mapping (category_id, attribute_id, is_enabled, is_required, sort_order)
    SELECT cat_id, id, true, 
        CASE 
            WHEN name IN ('brand', 'device_condition') THEN true 
            ELSE false 
        END,
        CASE 
            WHEN name = 'brand' THEN 1
            WHEN name = 'model' THEN 2
            WHEN name = 'device_condition' THEN 3
            WHEN name = 'storage' THEN 4
            WHEN name = 'color' THEN 5
            WHEN name = 'warranty' THEN 6
            WHEN name = 'accessories' THEN 7
            WHEN name = 'screen_size' THEN 8
        END
    FROM category_attributes
    WHERE name IN ('brand', 'model', 'device_condition', 'storage', 'color', 'warranty', 'accessories', 'screen_size')
    ON CONFLICT (category_id, attribute_id) DO UPDATE
    SET is_enabled = EXCLUDED.is_enabled,
        is_required = EXCLUDED.is_required,
        sort_order = EXCLUDED.sort_order;
END $$;
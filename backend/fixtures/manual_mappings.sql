-- Добавляем недостающие атрибуты
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

-- Обновляем существующие атрибуты
UPDATE category_attributes SET 
    display_name = 'Бренд',
    attribute_type = 'text',
    is_searchable = true,
    is_filterable = true
WHERE name = 'brand';

UPDATE category_attributes SET 
    display_name = 'Цвет',
    attribute_type = 'select',
    options = '{"values": ["Черный", "Белый", "Серый", "Синий", "Красный", "Зеленый", "Золотой", "Серебристый", "Другой"]}',
    is_filterable = true
WHERE name = 'color';

UPDATE category_attributes SET 
    display_name = 'Размер экрана (дюймы)',
    options = '{"min": 1, "max": 100, "step": 0.1}',
    is_filterable = true
WHERE name = 'screen_size';

-- Добавляем маппинги
INSERT INTO category_attribute_mapping (category_id, attribute_id, is_enabled, is_required, sort_order)
VALUES 
    (1, (SELECT id FROM category_attributes WHERE name = 'brand'), true, true, 1),
    (1, (SELECT id FROM category_attributes WHERE name = 'model'), true, false, 2),
    (1, (SELECT id FROM category_attributes WHERE name = 'device_condition'), true, true, 3),
    (1, (SELECT id FROM category_attributes WHERE name = 'storage'), true, false, 4),
    (1, (SELECT id FROM category_attributes WHERE name = 'color'), true, false, 5),
    (1, (SELECT id FROM category_attributes WHERE name = 'warranty'), true, false, 6),
    (1, (SELECT id FROM category_attributes WHERE name = 'accessories'), true, false, 7),
    (1, (SELECT id FROM category_attributes WHERE name = 'screen_size'), true, false, 8)
ON CONFLICT (category_id, attribute_id) DO UPDATE
SET is_enabled = EXCLUDED.is_enabled,
    is_required = EXCLUDED.is_required,
    sort_order = EXCLUDED.sort_order;
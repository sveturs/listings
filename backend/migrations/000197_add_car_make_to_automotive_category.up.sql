-- Добавляем атрибут car_make к основной автомобильной категории (1003)
INSERT INTO category_attribute_mapping (category_id, attribute_id, is_enabled, is_required, sort_order)
SELECT 
    1003 as category_id,
    (SELECT id FROM category_attributes WHERE name = 'car_make') as attribute_id,
    true as is_enabled,
    true as is_required,  -- Делаем марку обязательной для автомобилей
    5 as sort_order  -- Ставим перед моделью (которая имеет sort_order = 10)
WHERE 
    EXISTS (SELECT 1 FROM category_attributes WHERE name = 'car_make')
    AND NOT EXISTS (
        SELECT 1 FROM category_attribute_mapping 
        WHERE category_id = 1003 
        AND attribute_id = (SELECT id FROM category_attributes WHERE name = 'car_make')
    );

-- Добавляем атрибут car_make_id к основной автомобильной категории (1003)
INSERT INTO category_attribute_mapping (category_id, attribute_id, is_enabled, is_required, sort_order)
SELECT 
    1003 as category_id,
    (SELECT id FROM category_attributes WHERE name = 'car_make_id') as attribute_id,
    true as is_enabled,
    false as is_required,  -- ID не обязательный, но полезный для связи
    6 as sort_order  -- После марки
WHERE 
    EXISTS (SELECT 1 FROM category_attributes WHERE name = 'car_make_id')
    AND NOT EXISTS (
        SELECT 1 FROM category_attribute_mapping 
        WHERE category_id = 1003 
        AND attribute_id = (SELECT id FROM category_attributes WHERE name = 'car_make_id')
    );

-- Добавляем атрибут car_model_id к основной автомобильной категории (1003) если его еще нет
INSERT INTO category_attribute_mapping (category_id, attribute_id, is_enabled, is_required, sort_order)
SELECT 
    1003 as category_id,
    (SELECT id FROM category_attributes WHERE name = 'car_model_id') as attribute_id,
    true as is_enabled,
    false as is_required,  -- ID не обязательный, но полезный для связи
    11 as sort_order  -- После модели
WHERE 
    EXISTS (SELECT 1 FROM category_attributes WHERE name = 'car_model_id')
    AND NOT EXISTS (
        SELECT 1 FROM category_attribute_mapping 
        WHERE category_id = 1003 
        AND attribute_id = (SELECT id FROM category_attributes WHERE name = 'car_model_id')
    );
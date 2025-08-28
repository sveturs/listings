-- Добавляем атрибут vin_number к основной автомобильной категории (1003)
INSERT INTO category_attribute_mapping (category_id, attribute_id, is_enabled, is_required, sort_order)
SELECT 
    1003 as category_id,
    (SELECT id FROM category_attributes WHERE name = 'vin_number') as attribute_id,
    true as is_enabled,
    false as is_required,  -- VIN не обязательный, но полезный
    13 as sort_order  -- После пробега
WHERE 
    EXISTS (SELECT 1 FROM category_attributes WHERE name = 'vin_number')
    AND NOT EXISTS (
        SELECT 1 FROM category_attribute_mapping 
        WHERE category_id = 1003 
        AND attribute_id = (SELECT id FROM category_attributes WHERE name = 'vin_number')
    );
-- Добавляем недостающие атрибуты в маппинг для категории Automobili (1003)

-- Добавляем car_model (ID: 2201)
INSERT INTO category_attribute_mapping (category_id, attribute_id, is_enabled, is_required, sort_order)
VALUES (1003, 2201, true, true, 10)
ON CONFLICT (category_id, attribute_id) DO UPDATE SET
    is_enabled = EXCLUDED.is_enabled,
    is_required = EXCLUDED.is_required,
    sort_order = EXCLUDED.sort_order;

-- Добавляем year (ID: 2202) 
INSERT INTO category_attribute_mapping (category_id, attribute_id, is_enabled, is_required, sort_order)
VALUES (1003, 2202, true, true, 11)
ON CONFLICT (category_id, attribute_id) DO UPDATE SET
    is_enabled = EXCLUDED.is_enabled,
    is_required = EXCLUDED.is_required,
    sort_order = EXCLUDED.sort_order;

-- Добавляем mileage (ID: 2203)
INSERT INTO category_attribute_mapping (category_id, attribute_id, is_enabled, is_required, sort_order)
VALUES (1003, 2203, true, false, 12)
ON CONFLICT (category_id, attribute_id) DO UPDATE SET
    is_enabled = EXCLUDED.is_enabled,
    is_required = EXCLUDED.is_required,
    sort_order = EXCLUDED.sort_order;

-- Добавляем engine_size (ID: 3001)
INSERT INTO category_attribute_mapping (category_id, attribute_id, is_enabled, is_required, sort_order)
VALUES (1003, 3001, true, false, 13)
ON CONFLICT (category_id, attribute_id) DO UPDATE SET
    is_enabled = EXCLUDED.is_enabled,
    is_required = EXCLUDED.is_required,
    sort_order = EXCLUDED.sort_order;

-- Обновляем mapping_required для fuel_type и transmission (делаем обязательными)
UPDATE category_attribute_mapping 
SET is_required = true 
WHERE category_id = 1003 AND attribute_id IN (2204, 2205);
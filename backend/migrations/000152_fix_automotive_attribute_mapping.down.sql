-- Откат изменений для автомобильных атрибутов

-- Убираем обязательность для fuel_type и transmission
UPDATE category_attribute_mapping 
SET is_required = false 
WHERE category_id = 1003 AND attribute_id IN (2204, 2205);

-- Удаляем добавленные атрибуты из маппинга
DELETE FROM category_attribute_mapping 
WHERE category_id = 1003 AND attribute_id IN (2201, 2202, 2203, 3001);
-- Удаление связей вариативных атрибутов с категорией Cars
DELETE FROM category_variant_attributes 
WHERE variant_attribute_name IN ('engine_type', 'trim_level', 'transmission', 'fuel_type');

-- Удаление переводов для вариативных атрибутов
DELETE FROM translations 
WHERE entity_type = 'product_variant_attribute' 
AND entity_id IN ('engine_type', 'trim_level', 'transmission', 'fuel_type');

-- Удаление вариативных атрибутов
DELETE FROM product_variant_attributes 
WHERE name IN ('engine_type', 'trim_level', 'transmission', 'fuel_type');

-- Удаление или откат атрибутов категории
UPDATE category_attributes
SET is_variant_compatible = false
WHERE name IN ('engine_type', 'trim_level', 'transmission', 'fuel_type');

-- Можно также полностью удалить эти атрибуты, если они были созданы только для вариантов
-- DELETE FROM category_attributes WHERE name IN ('engine_type', 'trim_level', 'transmission', 'fuel_type');
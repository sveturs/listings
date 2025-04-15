-- Добавляем колонку unit в listing_attribute_values
ALTER TABLE listing_attribute_values
ADD COLUMN IF NOT EXISTS unit VARCHAR(20) DEFAULT NULL;

-- Обновляем существующие данные для атрибутов недвижимости
UPDATE listing_attribute_values
SET unit = 'soba'
WHERE attribute_id IN (SELECT id FROM category_attributes WHERE name = 'rooms')
AND numeric_value IS NOT NULL;

UPDATE listing_attribute_values
SET unit = 'm²'
WHERE attribute_id IN (SELECT id FROM category_attributes WHERE name = 'area')
AND numeric_value IS NOT NULL;

-- Убеждаемся, что category_attributes содержит информацию о единицах и типах
ALTER TABLE category_attributes
ALTER COLUMN options SET DEFAULT '{}'::jsonb;

UPDATE category_attributes
SET options = jsonb_build_object('type', 'number', 'unit', 'soba')
WHERE name = 'rooms';

UPDATE category_attributes
SET options = jsonb_build_object('type', 'number', 'unit', 'm²')
WHERE name = 'area';

UPDATE category_attributes
SET options = jsonb_build_object('type', 'text', 'unit', NULL)
WHERE name = 'property_type';
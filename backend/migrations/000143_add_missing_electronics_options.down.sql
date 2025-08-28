-- Возвращаем пустые опции для атрибутов
UPDATE category_attributes 
SET options = '{}'::jsonb
WHERE name IN ('ram', 'storage_type') AND id IN (2104, 2105);
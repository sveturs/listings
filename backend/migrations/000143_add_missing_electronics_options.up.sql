-- Добавляем опции для атрибута RAM
UPDATE category_attributes 
SET options = jsonb_build_object('values', ARRAY['2GB', '4GB', '8GB', '16GB', '32GB', '64GB', '128GB'])
WHERE name = 'ram' AND id = 2104;

-- Добавляем опции для атрибута storage_type
UPDATE category_attributes 
SET options = jsonb_build_object('values', ARRAY['HDD', 'SSD', 'NVMe', 'eMMC', 'Hybrid', 'Memory Card'])
WHERE name = 'storage_type' AND id = 2105;
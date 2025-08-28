-- Откат миграции для поля data_source

-- Удаляем ограничение
ALTER TABLE category_attributes 
DROP CONSTRAINT IF EXISTS check_data_source_values;

-- Удаляем индекс
DROP INDEX IF EXISTS idx_category_attributes_data_source;

-- Удаляем колонку
ALTER TABLE category_attributes 
DROP COLUMN IF EXISTS data_source;
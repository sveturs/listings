-- Откат изменений

-- Удаление представления map_items_view
DROP VIEW IF EXISTS map_items_view CASCADE;

-- Удаление функции
DROP FUNCTION IF EXISTS refresh_map_items_cache();

-- Удаление колонки last_modified_by из translations
ALTER TABLE translations 
DROP COLUMN IF EXISTS last_modified_by;

-- Удаление таблицы listing_attribute_values
DROP TABLE IF EXISTS listing_attribute_values CASCADE;
-- Удаляем индекс
DROP INDEX IF EXISTS idx_listing_attribute_values_unit;

-- Удаляем колонку unit
ALTER TABLE listing_attribute_values
DROP COLUMN IF EXISTS unit;
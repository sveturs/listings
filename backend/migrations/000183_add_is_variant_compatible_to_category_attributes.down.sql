-- Откат миграции: удаление поля is_variant_compatible

-- Удаляем индекс
DROP INDEX IF EXISTS idx_category_attributes_is_variant_compatible;

-- Удаляем поле
ALTER TABLE category_attributes 
DROP COLUMN IF EXISTS is_variant_compatible;
-- Откат миграции 044: удаление полей для вариативных атрибутов

-- Удаление таблицы связей
DROP TABLE IF EXISTS variant_attribute_mappings;

-- Удаление поля из таблицы unified_attributes
ALTER TABLE unified_attributes 
DROP COLUMN IF EXISTS is_variant_compatible;
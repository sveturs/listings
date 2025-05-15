-- Удаление индекса для поля custom_component
DROP INDEX IF EXISTS idx_category_attribute_mapping_custom_component;

-- Удаление поля custom_component из таблицы category_attribute_mapping
ALTER TABLE category_attribute_mapping 
    DROP COLUMN IF EXISTS custom_component;
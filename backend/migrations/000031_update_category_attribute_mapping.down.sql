-- Откат миграции
-- Удаляем представление
DROP VIEW IF EXISTS v_category_attributes;

-- Удаляем триггер
DROP TRIGGER IF EXISTS tr_update_category_attribute_sort_order ON category_attribute_mapping;

-- Удаляем функцию
DROP FUNCTION IF EXISTS update_category_attribute_sort_order();

-- Удаляем индексы
DROP INDEX IF EXISTS idx_category_attribute_map_cat_id;
DROP INDEX IF EXISTS idx_category_attribute_map_attr_id;

-- Удаляем добавленные столбцы
ALTER TABLE category_attribute_mapping 
    DROP COLUMN IF EXISTS sort_order,
    DROP COLUMN IF EXISTS custom_component;
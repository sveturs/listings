-- Удаляем представления
DROP VIEW IF EXISTS v_attribute_groups_with_items;

-- Удаляем триггеры
DROP TRIGGER IF EXISTS update_attribute_groups_updated_at ON attribute_groups;

-- Удаляем функции  
DROP FUNCTION IF EXISTS update_attribute_groups_updated_at();

-- Удаляем индексы
DROP INDEX IF EXISTS idx_category_attribute_groups_component;
DROP INDEX IF EXISTS idx_category_attribute_groups_group;
DROP INDEX IF EXISTS idx_category_attribute_groups_category;
DROP INDEX IF EXISTS idx_attribute_group_items_attribute;
DROP INDEX IF EXISTS idx_attribute_group_items_group;
DROP INDEX IF EXISTS idx_attribute_groups_active;
DROP INDEX IF EXISTS idx_attribute_groups_name;

-- Удаляем таблицы
DROP TABLE IF EXISTS category_attribute_groups;
DROP TABLE IF EXISTS attribute_group_items;
DROP TABLE IF EXISTS attribute_groups;
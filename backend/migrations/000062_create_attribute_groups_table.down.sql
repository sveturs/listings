-- Удаление связи атрибутов с группами
ALTER TABLE attributes DROP COLUMN IF EXISTS group_id;

-- Удаление таблицы связей категорий с группами
DROP TABLE IF EXISTS category_attribute_groups;

-- Удаление таблицы групп атрибутов
DROP TABLE IF EXISTS attribute_groups;
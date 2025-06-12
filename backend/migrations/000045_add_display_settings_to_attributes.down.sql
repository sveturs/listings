-- Откат миграции - удаляем добавленные колонки
ALTER TABLE category_attribute_mapping 
DROP COLUMN IF EXISTS show_in_card,
DROP COLUMN IF EXISTS show_in_list;

ALTER TABLE category_attributes 
DROP COLUMN IF EXISTS show_in_card,
DROP COLUMN IF EXISTS show_in_list;
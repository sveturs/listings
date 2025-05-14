-- Удаление полей для управления кастомными компонентами в категориях
ALTER TABLE marketplace_categories 
DROP COLUMN IF EXISTS has_custom_ui,
DROP COLUMN IF EXISTS custom_ui_component;

-- Удаление поля для кастомного компонента в атрибутах
ALTER TABLE category_attributes 
DROP COLUMN IF EXISTS custom_component;

-- Удаление индексов, созданных в up-миграции
DROP INDEX IF EXISTS idx_category_attributes_name;
DROP INDEX IF EXISTS idx_categories_path;
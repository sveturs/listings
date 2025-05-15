-- Эта миграция исправляет проблему с миграцией 29, поэтому её откат повторяет откат миграции 29

-- Удаление полей для управления кастомными компонентами в категориях
ALTER TABLE marketplace_categories 
DROP COLUMN IF EXISTS has_custom_ui,
DROP COLUMN IF EXISTS custom_ui_component;

-- Удаление поля для кастомного компонента в атрибутах
ALTER TABLE category_attributes 
DROP COLUMN IF EXISTS custom_component;

-- Удаление индексов, созданных в up-миграции
DROP INDEX IF EXISTS idx_category_attributes_name;
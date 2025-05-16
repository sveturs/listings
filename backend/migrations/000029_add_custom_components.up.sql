-- Добавление полей для управления кастомными компонентами в категориях
ALTER TABLE marketplace_categories 
ADD COLUMN IF NOT EXISTS has_custom_ui BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS custom_ui_component VARCHAR(255);

-- Добавление поля для кастомного компонента в атрибутах
ALTER TABLE category_attributes 
ADD COLUMN IF NOT EXISTS custom_component VARCHAR(255);

-- Индексы для оптимизации запросов в админке
CREATE INDEX IF NOT EXISTS idx_category_attributes_name ON category_attributes(name);
-- Убираем индекс для несуществующей колонки category_path
-- CREATE INDEX IF NOT EXISTS idx_categories_path ON marketplace_categories USING GIN (category_path);
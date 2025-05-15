-- Сброс состояния миграции 29
UPDATE schema_migrations SET dirty = false WHERE version = 29;

-- Добавление полей для управления кастомными компонентами в категориях, которые должны были быть добавлены в миграции 29
ALTER TABLE marketplace_categories 
ADD COLUMN IF NOT EXISTS has_custom_ui BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS custom_ui_component VARCHAR(255);

-- Добавление поля для кастомного компонента в атрибутах
ALTER TABLE category_attributes 
ADD COLUMN IF NOT EXISTS custom_component VARCHAR(255);

-- Индексы для оптимизации запросов в админке
CREATE INDEX IF NOT EXISTS idx_category_attributes_name ON category_attributes(name);
-- Не создаем индекс на несуществующем поле category_path
-- CREATE INDEX IF NOT EXISTS idx_categories_path ON marketplace_categories USING GIN (category_path);
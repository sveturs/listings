-- Добавление поля custom_component в таблицу category_attribute_mapping
ALTER TABLE category_attribute_mapping 
    ADD COLUMN IF NOT EXISTS custom_component VARCHAR(255) DEFAULT NULL;

-- Добавляем индекс для более быстрого поиска по кастомным компонентам
CREATE INDEX IF NOT EXISTS idx_category_attribute_mapping_custom_component
    ON category_attribute_mapping (custom_component);

-- Комментарий для колонки
COMMENT ON COLUMN category_attribute_mapping.custom_component IS 'Название пользовательского компонента для отображения атрибута';
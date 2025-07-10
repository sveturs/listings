-- Откат изменений для custom_ui_component_usage
ALTER TABLE custom_ui_component_usage
DROP COLUMN category_id,
DROP COLUMN usage_context,
DROP COLUMN placement,
DROP COLUMN priority,
DROP COLUMN conditions_logic,
DROP COLUMN is_active,
DROP COLUMN created_by,
DROP COLUMN updated_by,
ADD COLUMN entity_type VARCHAR(50) NOT NULL DEFAULT 'category' CHECK (entity_type IN ('category', 'attribute')),
ADD COLUMN entity_id INT NOT NULL DEFAULT 0;

-- Удаляем новые индексы
DROP INDEX IF EXISTS idx_component_usage_category;
DROP INDEX IF EXISTS idx_component_usage_context;
DROP INDEX IF EXISTS idx_component_usage_active;

-- Восстанавливаем старые индексы
CREATE INDEX idx_component_usage_entity ON custom_ui_component_usage(entity_type, entity_id);
CREATE UNIQUE INDEX idx_component_usage_unique ON custom_ui_component_usage(component_id, entity_type, entity_id);

-- Откат изменений для custom_ui_templates
ALTER TABLE custom_ui_templates
DROP COLUMN variables,
DROP COLUMN is_shared,
DROP COLUMN created_by,
DROP COLUMN updated_by,
ADD COLUMN template_type VARCHAR(50) NOT NULL DEFAULT 'category' CHECK (template_type IN ('category', 'attribute', 'filter')),
ADD COLUMN example_configuration JSONB DEFAULT '{}',
ADD COLUMN is_active BOOLEAN DEFAULT true;

-- Удаляем триггер
DROP TRIGGER IF EXISTS set_custom_ui_component_usage_timestamp ON custom_ui_component_usage;
-- Обновляем таблицу custom_ui_component_usage для соответствия бизнес-модели
ALTER TABLE custom_ui_component_usage
DROP COLUMN IF EXISTS entity_type,
DROP COLUMN IF EXISTS entity_id,
ADD COLUMN IF NOT EXISTS category_id INT REFERENCES marketplace_categories(id) ON DELETE CASCADE,
ADD COLUMN IF NOT EXISTS usage_context VARCHAR(100) NOT NULL DEFAULT 'category',
ADD COLUMN IF NOT EXISTS placement VARCHAR(100),
ADD COLUMN IF NOT EXISTS priority INT DEFAULT 0,
ADD COLUMN IF NOT EXISTS conditions_logic JSONB DEFAULT '{}',
ADD COLUMN IF NOT EXISTS is_active BOOLEAN DEFAULT true,
ADD COLUMN IF NOT EXISTS created_by INT REFERENCES users(id),
ADD COLUMN IF NOT EXISTS updated_by INT REFERENCES users(id);

-- Удаляем старые индексы
DROP INDEX IF EXISTS idx_component_usage_entity;
DROP INDEX IF EXISTS idx_component_usage_unique;

-- Создаем новые индексы
CREATE INDEX IF NOT EXISTS idx_component_usage_category ON custom_ui_component_usage(category_id);
CREATE INDEX IF NOT EXISTS idx_component_usage_context ON custom_ui_component_usage(usage_context);
CREATE INDEX IF NOT EXISTS idx_component_usage_active ON custom_ui_component_usage(is_active);

-- Обновляем таблицу custom_ui_templates
ALTER TABLE custom_ui_templates
DROP COLUMN IF EXISTS template_type,
DROP COLUMN IF EXISTS example_configuration,
DROP COLUMN IF EXISTS is_active,
ADD COLUMN IF NOT EXISTS variables JSONB DEFAULT '{}',
ADD COLUMN IF NOT EXISTS is_shared BOOLEAN DEFAULT false,
ADD COLUMN IF NOT EXISTS created_by INT REFERENCES users(id),
ADD COLUMN IF NOT EXISTS updated_by INT REFERENCES users(id);

-- Триггер для updated_at на custom_ui_component_usage
CREATE OR REPLACE FUNCTION update_custom_ui_component_usage_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS set_custom_ui_component_usage_timestamp ON set_custom_ui_component_usage_timestamp;
CREATE TRIGGER set_custom_ui_component_usage_timestamp
BEFORE UPDATE ON custom_ui_component_usage
FOR EACH ROW
EXECUTE FUNCTION update_custom_ui_component_usage_timestamp();
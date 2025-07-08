-- Откат изменений миграции 000037

-- Удаляем созданные таблицы
DROP TABLE IF EXISTS component_templates CASCADE;
DROP TABLE IF EXISTS custom_ui_component_usage CASCADE;

-- Возвращаем старые колонки
ALTER TABLE custom_ui_components ADD COLUMN IF NOT EXISTS display_name VARCHAR(255) NOT NULL DEFAULT '';
ALTER TABLE custom_ui_components ADD COLUMN IF NOT EXISTS component_code TEXT NOT NULL DEFAULT '';
ALTER TABLE custom_ui_components ADD COLUMN IF NOT EXISTS configuration JSONB DEFAULT '{}';
ALTER TABLE custom_ui_components ADD COLUMN IF NOT EXISTS dependencies JSONB DEFAULT '[]';
ALTER TABLE custom_ui_components ADD COLUMN IF NOT EXISTS compiled_code TEXT;
ALTER TABLE custom_ui_components ADD COLUMN IF NOT EXISTS compilation_errors JSONB;
ALTER TABLE custom_ui_components ADD COLUMN IF NOT EXISTS last_compiled_at TIMESTAMP WITH TIME ZONE;

-- Удаляем новые колонки
ALTER TABLE custom_ui_components DROP COLUMN IF EXISTS template_code;
ALTER TABLE custom_ui_components DROP COLUMN IF EXISTS styles;
ALTER TABLE custom_ui_components DROP COLUMN IF EXISTS props_schema;

-- Возвращаем старую структуру custom_ui_component_usage
CREATE TABLE IF NOT EXISTS custom_ui_component_usage (
    id SERIAL PRIMARY KEY,
    component_id INT NOT NULL REFERENCES custom_ui_components(id) ON DELETE CASCADE,
    entity_type VARCHAR(50) NOT NULL CHECK (entity_type IN ('category', 'attribute')),
    entity_id INT NOT NULL,
    configuration JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для старой таблицы
CREATE INDEX IF NOT EXISTS idx_component_usage_component ON custom_ui_component_usage(component_id);
CREATE INDEX IF NOT EXISTS idx_component_usage_entity ON custom_ui_component_usage(entity_type, entity_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_component_usage_unique ON custom_ui_component_usage(component_id, entity_type, entity_id);
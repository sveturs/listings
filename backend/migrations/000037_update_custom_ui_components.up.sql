-- Обновление структуры таблицы custom_ui_components для новой схемы

-- Сначала создаём custom_ui_component_usage если её нет (возможно не применилась миграция 33)
CREATE TABLE IF NOT EXISTS custom_ui_component_usage (
    id SERIAL PRIMARY KEY,
    component_id INT NOT NULL REFERENCES custom_ui_components(id) ON DELETE CASCADE,
    entity_type VARCHAR(50) NOT NULL CHECK (entity_type IN ('category', 'attribute')),
    entity_id INT NOT NULL,
    configuration JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Удаляем старые колонки в правильном порядке
ALTER TABLE custom_ui_components DROP COLUMN IF EXISTS compiled_code CASCADE;
ALTER TABLE custom_ui_components DROP COLUMN IF EXISTS compilation_errors CASCADE;
ALTER TABLE custom_ui_components DROP COLUMN IF EXISTS last_compiled_at CASCADE;
ALTER TABLE custom_ui_components DROP COLUMN IF EXISTS dependencies CASCADE;
ALTER TABLE custom_ui_components DROP COLUMN IF EXISTS configuration CASCADE;
ALTER TABLE custom_ui_components DROP COLUMN IF EXISTS component_code CASCADE;
ALTER TABLE custom_ui_components DROP COLUMN IF EXISTS display_name CASCADE;

-- Добавляем новые колонки
ALTER TABLE custom_ui_components ADD COLUMN IF NOT EXISTS template_code TEXT DEFAULT '';
ALTER TABLE custom_ui_components ADD COLUMN IF NOT EXISTS styles TEXT DEFAULT '';
ALTER TABLE custom_ui_components ADD COLUMN IF NOT EXISTS props_schema JSONB DEFAULT '{}';

-- Делаем template_code NOT NULL после добавления
UPDATE custom_ui_components SET template_code = '' WHERE template_code IS NULL;
ALTER TABLE custom_ui_components ALTER COLUMN template_code SET NOT NULL;

-- Обновляем структуру custom_ui_component_usage
DROP TABLE IF EXISTS custom_ui_component_usage CASCADE;

CREATE TABLE custom_ui_component_usage (
    id SERIAL PRIMARY KEY,
    component_id INT NOT NULL REFERENCES custom_ui_components(id) ON DELETE CASCADE,
    category_id INT REFERENCES marketplace_categories(id) ON DELETE CASCADE,
    usage_context VARCHAR(50) NOT NULL DEFAULT 'listing',
    placement VARCHAR(50) DEFAULT 'default',
    priority INT DEFAULT 0,
    configuration JSONB DEFAULT '{}',
    conditions_logic JSONB DEFAULT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by INT REFERENCES users(id),
    updated_by INT REFERENCES users(id)
);

-- Индексы для usage
CREATE INDEX IF NOT EXISTS idx_ui_comp_usage_component ON custom_ui_component_usage(component_id);
CREATE INDEX IF NOT EXISTS idx_ui_comp_usage_category ON custom_ui_component_usage(category_id);
CREATE INDEX IF NOT EXISTS idx_ui_comp_usage_context ON custom_ui_component_usage(usage_context);
CREATE INDEX IF NOT EXISTS idx_ui_comp_usage_active ON custom_ui_component_usage(is_active);

-- Создание таблицы для шаблонов компонентов
CREATE TABLE IF NOT EXISTS component_templates (
    id SERIAL PRIMARY KEY,
    component_id INT NOT NULL REFERENCES custom_ui_components(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    template_config JSONB DEFAULT '{}',
    preview_image TEXT,
    category_id INT REFERENCES marketplace_categories(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by INT REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_component_templates_comp ON component_templates(component_id);
CREATE INDEX IF NOT EXISTS idx_component_templates_cat ON component_templates(category_id);

-- Триггер для обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Удаляем старый триггер если существует и создаём новый
DROP TRIGGER IF EXISTS update_custom_ui_component_usage_updated_at ON custom_ui_component_usage;
DROP TRIGGER IF EXISTS update_custom_ui_component_usage_updated_at ON custom_ui_component_usage;
CREATE TRIGGER update_custom_ui_component_usage_updated_at 
BEFORE UPDATE ON custom_ui_component_usage 
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
-- Создание таблицы для хранения кастомных UI компонентов
CREATE TABLE IF NOT EXISTS custom_ui_components (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    description TEXT,
    component_type VARCHAR(50) NOT NULL CHECK (component_type IN ('category', 'attribute', 'filter')),
    component_code TEXT NOT NULL,
    configuration JSONB DEFAULT '{}',
    dependencies JSONB DEFAULT '[]',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by INT REFERENCES users(id),
    updated_by INT REFERENCES users(id)
);

-- Индексы
CREATE INDEX IF NOT EXISTS idx_custom_ui_components_name ON custom_ui_components(name);
CREATE INDEX IF NOT EXISTS idx_custom_ui_components_type ON custom_ui_components(component_type);
CREATE INDEX IF NOT EXISTS idx_custom_ui_components_active ON custom_ui_components(is_active);

-- Таблица для отслеживания использования компонентов
CREATE TABLE IF NOT EXISTS custom_ui_component_usage (
    id SERIAL PRIMARY KEY,
    component_id INT NOT NULL REFERENCES custom_ui_components(id) ON DELETE CASCADE,
    entity_type VARCHAR(50) NOT NULL CHECK (entity_type IN ('category', 'attribute')),
    entity_id INT NOT NULL,
    configuration JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для использования
CREATE INDEX IF NOT EXISTS idx_component_usage_component ON custom_ui_component_usage(component_id);
CREATE INDEX IF NOT EXISTS idx_component_usage_entity ON custom_ui_component_usage(entity_type, entity_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_component_usage_unique ON custom_ui_component_usage(component_id, entity_type, entity_id);

-- Таблица для шаблонов компонентов
CREATE TABLE IF NOT EXISTS custom_ui_templates (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    description TEXT,
    template_code TEXT NOT NULL,
    template_type VARCHAR(50) NOT NULL CHECK (template_type IN ('category', 'attribute', 'filter')),
    example_configuration JSONB DEFAULT '{}',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для шаблонов
CREATE INDEX IF NOT EXISTS idx_custom_ui_templates_name ON custom_ui_templates(name);
CREATE INDEX IF NOT EXISTS idx_custom_ui_templates_type ON custom_ui_templates(template_type);

-- Добавляем поле для хранения скомпилированного кода (для кэширования)
ALTER TABLE custom_ui_components 
ADD COLUMN IF NOT EXISTS compiled_code TEXT,
ADD COLUMN IF NOT EXISTS compilation_errors JSONB,
ADD COLUMN IF NOT EXISTS last_compiled_at TIMESTAMP WITH TIME ZONE;

-- Функция для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Триггеры для автоматического обновления updated_at
DROP TRIGGER IF EXISTS update_custom_ui_components_updated_at ON update_custom_ui_components_updated_at;
CREATE TRIGGER update_custom_ui_components_updated_at BEFORE UPDATE
    ON custom_ui_components FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_custom_ui_component_usage_updated_at ON update_custom_ui_component_usage_updated_at;
CREATE TRIGGER update_custom_ui_component_usage_updated_at BEFORE UPDATE
    ON custom_ui_component_usage FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_custom_ui_templates_updated_at ON update_custom_ui_templates_updated_at;
CREATE TRIGGER update_custom_ui_templates_updated_at BEFORE UPDATE
    ON custom_ui_templates FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Комментарии к таблицам
COMMENT ON TABLE custom_ui_components IS 'Хранение кастомных UI компонентов';
COMMENT ON TABLE custom_ui_component_usage IS 'Отслеживание использования компонентов';
COMMENT ON TABLE custom_ui_templates IS 'Шаблоны для создания новых компонентов';

-- Комментарии к полям
COMMENT ON COLUMN custom_ui_components.component_code IS 'JSX код компонента';
COMMENT ON COLUMN custom_ui_components.configuration IS 'Конфигурация компонента (пропсы, настройки)';
COMMENT ON COLUMN custom_ui_components.dependencies IS 'Список внешних зависимостей компонента';
COMMENT ON COLUMN custom_ui_components.compiled_code IS 'Транспилированный код компонента для выполнения';
COMMENT ON COLUMN custom_ui_component_usage.configuration IS 'Специфичная конфигурация для данного использования';
-- Удаление триггеров
DROP TRIGGER IF EXISTS update_custom_ui_components_updated_at ON custom_ui_components;
DROP TRIGGER IF EXISTS update_custom_ui_component_usage_updated_at ON custom_ui_component_usage;
DROP TRIGGER IF EXISTS update_custom_ui_templates_updated_at ON custom_ui_templates;

-- Удаление функции
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Удаление таблиц в обратном порядке зависимостей
DROP TABLE IF EXISTS custom_ui_component_usage;
DROP TABLE IF EXISTS custom_ui_templates;
DROP TABLE IF EXISTS custom_ui_components;
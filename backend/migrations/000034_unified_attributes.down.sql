-- Откат миграции унифицированной системы атрибутов
-- ВАЖНО: Удаляем только новые таблицы, старые остаются нетронутыми

-- Удаление триггеров
DROP TRIGGER IF EXISTS update_unified_attributes_updated_at ON unified_attributes;
DROP TRIGGER IF EXISTS update_unified_category_attributes_updated_at ON unified_category_attributes;
DROP TRIGGER IF EXISTS update_unified_attribute_values_updated_at ON unified_attribute_values;

-- Удаление функции
DROP FUNCTION IF EXISTS update_unified_attributes_updated_at();

-- Удаление таблиц в правильном порядке (из-за foreign keys)
DROP TABLE IF EXISTS unified_attribute_values;
DROP TABLE IF EXISTS unified_category_attributes;
DROP TABLE IF EXISTS unified_attributes;
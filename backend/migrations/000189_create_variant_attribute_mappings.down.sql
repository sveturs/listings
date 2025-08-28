-- Удаление триггера
DROP TRIGGER IF EXISTS update_variant_attribute_mappings_updated_at ON variant_attribute_mappings;

-- Удаление индексов
DROP INDEX IF EXISTS idx_variant_attribute_mappings_variant_id;
DROP INDEX IF EXISTS idx_variant_attribute_mappings_category_id;

-- Удаление таблицы
DROP TABLE IF EXISTS variant_attribute_mappings;
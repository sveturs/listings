-- Удаляем триггер
DROP TRIGGER IF EXISTS update_category_variant_attributes_updated_at ON category_variant_attributes;

-- Удаляем индексы
DROP INDEX IF EXISTS idx_category_variant_attributes_variant_name;
DROP INDEX IF EXISTS idx_category_variant_attributes_category_id;

-- Удаляем таблицу
DROP TABLE IF EXISTS category_variant_attributes;
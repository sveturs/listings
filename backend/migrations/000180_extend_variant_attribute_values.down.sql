-- Удалить переводы значений атрибутов
DELETE FROM translations 
WHERE entity_type = 'product_variant_attribute_value';

-- Удалить добавленные значения атрибутов
DELETE FROM product_variant_attribute_values 
WHERE value IN (
    'xs', 's', 'm', 'l', 'xl', 'xxl', 'xxxl',
    '36', '38', '40', '42', '44', '46', '48', '50', '52',
    'black', 'white', 'gray', 'red', 'blue', 'green', 'yellow', 
    'orange', 'purple', 'pink', 'brown', 'navy', 'beige', 'gold', 'silver'
);

-- Удалить индексы
DROP INDEX IF EXISTS idx_product_variant_attribute_values_metadata_gin;
DROP INDEX IF EXISTS idx_product_variant_attribute_values_popular;

-- Удалить добавленные колонки
ALTER TABLE product_variant_attribute_values
DROP COLUMN IF EXISTS usage_count,
DROP COLUMN IF EXISTS is_popular,
DROP COLUMN IF EXISTS metadata;
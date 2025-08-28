-- Удалить индекс
DROP INDEX IF EXISTS idx_product_variant_attributes_affects_stock;

-- Удалить переводы атрибутов
DELETE FROM translations 
WHERE entity_type = 'product_variant_attribute' 
AND field_name = 'display_name';

-- Удалить поле affects_stock
ALTER TABLE product_variant_attributes 
DROP COLUMN IF EXISTS affects_stock;
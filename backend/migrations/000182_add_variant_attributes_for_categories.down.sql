-- Удаление уникального индекса если он был создан в этой миграции
DROP INDEX IF EXISTS product_variant_attributes_name_unique;

-- Удаление переводов для новых атрибутов
DELETE FROM translations 
WHERE entity_type = 'product_variant_attribute' 
AND entity_id IN (
    SELECT id FROM product_variant_attributes 
    WHERE name IN ('memory', 'storage', 'material', 'capacity', 'power', 'connectivity', 'style', 'pattern', 'weight', 'bundle')
);

-- Удаление новых вариативных атрибутов
DELETE FROM product_variant_attributes 
WHERE name IN ('memory', 'storage', 'material', 'capacity', 'power', 'connectivity', 'style', 'pattern', 'weight', 'bundle');

-- Удаление комментариев
COMMENT ON TABLE product_variant_attributes IS NULL;
COMMENT ON COLUMN product_variant_attributes.affects_stock IS NULL;
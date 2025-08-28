-- Добавляем поле affects_stock в product_variant_attributes
-- Это поле указывает, влияет ли данный вариативный атрибут на учет остатков товара
ALTER TABLE product_variant_attributes 
ADD COLUMN IF NOT EXISTS affects_stock BOOLEAN DEFAULT FALSE;

-- Добавляем комментарий к полю
COMMENT ON COLUMN product_variant_attributes.affects_stock IS 
'Указывает, влияет ли данный вариативный атрибут на учет остатков товара. Например, размер обуви влияет на остатки, а цвет гравировки - нет';

-- Устанавливаем affects_stock = true для атрибутов, которые обычно влияют на остатки
UPDATE product_variant_attributes
SET affects_stock = TRUE
WHERE name IN ('size', 'memory', 'storage', 'capacity', 'volume')
  AND affects_stock IS NOT TRUE;

-- Добавляем индекс для быстрого поиска атрибутов, влияющих на остатки
CREATE INDEX IF NOT EXISTS idx_product_variant_attributes_affects_stock 
ON product_variant_attributes(affects_stock) 
WHERE affects_stock = TRUE;
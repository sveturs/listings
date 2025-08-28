-- Добавляем поле is_variant_compatible в category_attributes
-- Это поле указывает, может ли атрибут категории использоваться как основа для создания вариантов товаров
ALTER TABLE category_attributes 
ADD COLUMN IF NOT EXISTS is_variant_compatible BOOLEAN DEFAULT FALSE;

-- Добавляем комментарий к полю
COMMENT ON COLUMN category_attributes.is_variant_compatible IS 
'Указывает, может ли данный атрибут категории использоваться как основа для создания вариантов товаров';

-- Устанавливаем is_variant_compatible = true для атрибутов, которые уже связаны с вариантами
UPDATE category_attributes ca
SET is_variant_compatible = TRUE
WHERE EXISTS (
    SELECT 1 
    FROM variant_attribute_mappings vam 
    WHERE vam.category_attribute_id = ca.id
);

-- Добавляем индекс для быстрого поиска совместимых атрибутов
CREATE INDEX IF NOT EXISTS idx_category_attributes_variant_compatible 
ON category_attributes(is_variant_compatible) 
WHERE is_variant_compatible = TRUE;
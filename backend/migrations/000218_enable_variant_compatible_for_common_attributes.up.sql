-- Enable variant compatibility for common variant attributes
UPDATE category_attributes 
SET is_variant_compatible = true 
WHERE name IN ('color', 'size', 'material', 'storage', 'memory', 'brand')
  AND is_variant_compatible = false;

-- Добавим комментарий для понимания
COMMENT ON COLUMN category_attributes.is_variant_compatible IS 'Indicates if this attribute can be used for product variants';
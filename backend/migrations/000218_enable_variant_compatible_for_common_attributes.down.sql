-- Disable variant compatibility for common variant attributes
UPDATE category_attributes 
SET is_variant_compatible = false 
WHERE name IN ('color', 'size', 'material', 'storage', 'memory', 'brand');
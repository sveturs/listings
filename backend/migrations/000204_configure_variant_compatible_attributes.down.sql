-- Откат настроек вариативных атрибутов

-- Удаляем индекс
DROP INDEX IF EXISTS idx_category_attributes_variant_compatible;

-- Сбрасываем настройки вариативности всех атрибутов
UPDATE category_attributes 
SET is_variant_compatible = false, affects_stock = false 
WHERE is_variant_compatible = true;
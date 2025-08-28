-- Добавить дополнительные поля для лучшего управления вариантами
ALTER TABLE storefront_product_variants 
ADD COLUMN IF NOT EXISTS reserved_quantity INTEGER NOT NULL DEFAULT 0,
ADD COLUMN IF NOT EXISTS available_quantity INTEGER GENERATED ALWAYS AS (stock_quantity - reserved_quantity) STORED;

-- Добавить составной уникальный индекс для variant_attributes
-- Это предотвратит создание дубликатов вариантов с одинаковыми атрибутами
CREATE UNIQUE INDEX idx_storefront_product_variants_unique_attributes 
ON storefront_product_variants(product_id, variant_attributes)
WHERE is_active = true;

-- Добавить индекс для быстрого поиска вариантов по статусу остатков
CREATE INDEX IF NOT EXISTS idx_storefront_product_variants_stock_status 
ON storefront_product_variants(stock_status, is_active);

-- Добавить индекс для поиска низких остатков
CREATE INDEX idx_storefront_product_variants_low_stock 
ON storefront_product_variants(product_id, available_quantity)
WHERE is_active = true AND available_quantity <= low_stock_threshold;

-- Добавить триггер для автоматического обновления stock_status
CREATE OR REPLACE FUNCTION update_variant_stock_status()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.available_quantity <= 0 THEN
        NEW.stock_status = 'out_of_stock';
    ELSIF NEW.available_quantity <= NEW.low_stock_threshold THEN
        NEW.stock_status = 'low_stock';
    ELSE
        NEW.stock_status = 'in_stock';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_update_variant_stock_status ON storefront_product_variants;
CREATE TRIGGER trigger_update_variant_stock_status
BEFORE INSERT OR UPDATE OF stock_quantity, reserved_quantity, low_stock_threshold
ON storefront_product_variants
FOR EACH ROW
EXECUTE FUNCTION update_variant_stock_status();

-- Обновить существующие записи для применения нового статуса
UPDATE storefront_product_variants 
SET stock_status = CASE
    WHEN stock_quantity - COALESCE(reserved_quantity, 0) <= 0 THEN 'out_of_stock'
    WHEN stock_quantity - COALESCE(reserved_quantity, 0) <= low_stock_threshold THEN 'low_stock'
    ELSE 'in_stock'
END;
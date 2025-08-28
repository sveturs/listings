-- Добавляем поле has_variants к таблице storefront_products
ALTER TABLE storefront_products
ADD COLUMN IF NOT EXISTS has_variants BOOLEAN NOT NULL DEFAULT false;

-- Обновляем существующие товары, у которых есть варианты
UPDATE storefront_products sp
SET has_variants = true
WHERE EXISTS (
    SELECT 1 
    FROM storefront_product_variants spv 
    WHERE spv.product_id = sp.id
);

-- Добавляем индекс для быстрого поиска товаров с вариантами
CREATE INDEX IF NOT EXISTS idx_storefront_products_has_variants 
ON storefront_products(has_variants) 
WHERE has_variants = true;
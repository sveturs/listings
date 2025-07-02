-- Возвращаем тип поля price к numeric(10,2)
ALTER TABLE storefront_products 
ALTER COLUMN price TYPE numeric(10,2);

-- Восстанавливаем ограничение
ALTER TABLE storefront_products 
DROP CONSTRAINT IF EXISTS storefront_products_price_check;

ALTER TABLE storefront_products 
ADD CONSTRAINT storefront_products_price_check CHECK (price >= 0);
-- Изменяем тип поля price для поддержки больших значений
ALTER TABLE storefront_products 
ALTER COLUMN price TYPE numeric(15,2);

-- Проверяем, что ограничение все еще работает
ALTER TABLE storefront_products 
DROP CONSTRAINT IF EXISTS storefront_products_price_check;

ALTER TABLE storefront_products 
ADD CONSTRAINT storefront_products_price_check CHECK (price >= 0);
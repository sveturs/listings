-- Удаляем поле pickup_address из таблицы storefront_orders
ALTER TABLE storefront_orders 
DROP COLUMN IF EXISTS pickup_address;
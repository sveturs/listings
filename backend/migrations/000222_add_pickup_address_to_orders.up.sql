-- Добавляем поле pickup_address в таблицу storefront_orders
ALTER TABLE storefront_orders 
ADD COLUMN IF NOT EXISTS pickup_address JSONB;

-- Комментарий к полю
COMMENT ON COLUMN storefront_orders.pickup_address IS 'Адрес забора товара у продавца (витрины)';
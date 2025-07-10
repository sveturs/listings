-- Откатываем изменения

-- Удаляем view
DROP VIEW IF EXISTS storefront_orders_view;

-- Удаляем добавленные колонки
ALTER TABLE storefront_orders
DROP COLUMN IF EXISTS payment_method,
DROP COLUMN IF EXISTS payment_status,
DROP COLUMN IF EXISTS notes,
DROP COLUMN IF EXISTS metadata,
DROP COLUMN IF EXISTS discount;
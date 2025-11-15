-- Remove customer and admin columns from orders table

DROP INDEX IF EXISTS idx_orders_customer_email;
DROP INDEX IF EXISTS idx_orders_customer_phone;

ALTER TABLE orders
DROP COLUMN IF EXISTS customer_name,
DROP COLUMN IF EXISTS customer_email,
DROP COLUMN IF EXISTS customer_phone,
DROP COLUMN IF EXISTS customer_notes,
DROP COLUMN IF EXISTS admin_notes;

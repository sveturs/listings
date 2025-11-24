-- Remove payment_completed_at column from orders table

DROP INDEX IF EXISTS idx_orders_payment_completed_at;

ALTER TABLE orders
DROP COLUMN IF EXISTS payment_completed_at;

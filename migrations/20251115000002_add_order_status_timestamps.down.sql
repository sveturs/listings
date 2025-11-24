-- Remove status transition timestamps from orders table
DROP INDEX IF EXISTS idx_orders_delivered_at;
DROP INDEX IF EXISTS idx_orders_shipped_at;
DROP INDEX IF EXISTS idx_orders_confirmed_at;

ALTER TABLE orders DROP COLUMN IF EXISTS delivered_at;
ALTER TABLE orders DROP COLUMN IF EXISTS shipped_at;
ALTER TABLE orders DROP COLUMN IF EXISTS confirmed_at;

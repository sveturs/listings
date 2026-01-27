-- Rollback: Remove delivery_method_details column from orders table

-- Drop index first
DROP INDEX IF EXISTS idx_orders_delivery_method_details_provider;

-- Remove column
ALTER TABLE orders DROP COLUMN IF EXISTS delivery_method_details;

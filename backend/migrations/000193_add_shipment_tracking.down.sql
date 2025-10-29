-- Rollback shipment tracking fields

-- Drop indexes
DROP INDEX IF EXISTS idx_c2c_orders_tracking_number;
DROP INDEX IF EXISTS idx_c2c_orders_shipment_id;

-- Revert shipment_id type change
ALTER TABLE c2c_orders
ALTER COLUMN shipment_id TYPE INTEGER;

-- Rename back to original name
ALTER TABLE c2c_orders
RENAME COLUMN shipment_id TO delivery_shipment_id;

-- Remove added column
ALTER TABLE c2c_orders
DROP COLUMN IF EXISTS shipping_provider;

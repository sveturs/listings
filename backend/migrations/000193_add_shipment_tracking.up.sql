-- Add missing shipment tracking fields
-- Following Single Source of Truth principle: only UI-needed fields
-- Full shipment data lives in delivery microservice DB

-- b2c_orders already has tracking_number and shipping_provider
-- c2c_orders has tracking_number but missing shipping_provider
-- Add shipping_provider to c2c_orders for consistency

ALTER TABLE c2c_orders
ADD COLUMN IF NOT EXISTS shipping_provider VARCHAR(50);

-- Rename delivery_shipment_id to shipment_id for consistency
ALTER TABLE c2c_orders
RENAME COLUMN delivery_shipment_id TO shipment_id;

-- Ensure shipment_id is BIGINT (matching microservice DB)
ALTER TABLE c2c_orders
ALTER COLUMN shipment_id TYPE BIGINT;

-- Add indexes for better performance (partial indexes on non-null values)
CREATE INDEX IF NOT EXISTS idx_c2c_orders_tracking_number
ON c2c_orders(tracking_number) WHERE tracking_number IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_c2c_orders_shipment_id
ON c2c_orders(shipment_id) WHERE shipment_id IS NOT NULL;

-- Documentation comments for c2c_orders
COMMENT ON COLUMN c2c_orders.shipping_provider
IS 'Delivery provider code: post_express, bex, aks, d_express, city_express (for UI icons/filters)';

COMMENT ON COLUMN c2c_orders.shipment_id
IS 'Shipment ID in delivery microservice DB (links to microservice shipments table)';

COMMENT ON COLUMN c2c_orders.tracking_number
IS 'Tracking number from delivery microservice (UI only, single source of truth in microservice DB)';

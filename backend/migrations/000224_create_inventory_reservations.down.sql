-- Drop trigger and function
DROP TRIGGER IF EXISTS trigger_update_inventory_reservations_updated_at ON inventory_reservations;
DROP FUNCTION IF EXISTS update_inventory_reservations_updated_at();

-- Drop indexes
DROP INDEX IF EXISTS idx_inventory_reservations_expires_at;
DROP INDEX IF EXISTS idx_inventory_reservations_status;
DROP INDEX IF EXISTS idx_inventory_reservations_order_id;
DROP INDEX IF EXISTS idx_inventory_reservations_variant_id;
DROP INDEX IF EXISTS idx_inventory_reservations_product_id;

-- Drop table (cascades will handle foreign key constraints)
DROP TABLE IF EXISTS inventory_reservations;

-- Drop enum (only if no other tables use it)
-- Note: We don't drop the enum as other tables might be using it
-- DROP TYPE IF EXISTS reservation_status;
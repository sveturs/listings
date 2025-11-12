-- Rollback migration: Remove rollback idempotency support

-- Drop indexes
DROP INDEX IF EXISTS idx_b2c_inventory_movements_rollback_idempotency;
DROP INDEX IF EXISTS idx_b2c_inventory_movements_rollback_variant_idempotency;
DROP INDEX IF EXISTS idx_b2c_inventory_movements_order_id;
DROP INDEX IF EXISTS idx_b2c_inventory_movements_movement_type;

-- Remove columns
ALTER TABLE b2c_inventory_movements
    DROP COLUMN IF EXISTS order_id;

ALTER TABLE b2c_inventory_movements
    DROP COLUMN IF EXISTS movement_type;

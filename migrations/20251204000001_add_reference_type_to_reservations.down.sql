-- Rollback: restore order_id and remove reference_type

-- 1. Drop new indexes
DROP INDEX IF EXISTS idx_inventory_reservations_ref_status;
DROP INDEX IF EXISTS idx_inventory_reservations_reference;

-- 2. Rename reference_id back to order_id
ALTER TABLE inventory_reservations
RENAME COLUMN reference_id TO order_id;

-- 3. Drop CHECK constraint for reference_type
ALTER TABLE inventory_reservations
DROP CONSTRAINT IF EXISTS chk_inventory_reservations_reference_type;

-- 4. Drop reference_type column
ALTER TABLE inventory_reservations
DROP COLUMN IF EXISTS reference_type;

-- 5. Recreate original index
CREATE INDEX IF NOT EXISTS idx_inventory_reservations_order_id
ON inventory_reservations(order_id) WHERE order_id IS NOT NULL;

-- 6. Note: FK to orders table is NOT restored as it depends on orders table existence

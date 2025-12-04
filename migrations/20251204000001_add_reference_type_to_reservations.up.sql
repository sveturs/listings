-- Add reference_type to support different reservation types (orders, transfers, etc.)
-- Rename order_id to reference_id for a more generic approach

-- 1. Add reference_type column
ALTER TABLE inventory_reservations
ADD COLUMN IF NOT EXISTS reference_type VARCHAR(20) NOT NULL DEFAULT 'order';

-- 2. Add CHECK constraint for reference_type
ALTER TABLE inventory_reservations
ADD CONSTRAINT chk_inventory_reservations_reference_type
CHECK (reference_type IN ('order', 'transfer'));

-- 3. Rename order_id to reference_id
ALTER TABLE inventory_reservations
RENAME COLUMN order_id TO reference_id;

-- 4. Drop old FK constraint that references orders table
ALTER TABLE inventory_reservations
DROP CONSTRAINT IF EXISTS fk_inventory_reservations_order;

-- 5. Update index names for clarity
DROP INDEX IF EXISTS idx_inventory_reservations_order_id;
CREATE INDEX IF NOT EXISTS idx_inventory_reservations_reference
ON inventory_reservations(reference_type, reference_id);

-- 6. Add composite index for reference lookup
CREATE INDEX IF NOT EXISTS idx_inventory_reservations_ref_status
ON inventory_reservations(reference_type, reference_id, status)
WHERE status = 'active';

-- Migration: Add order_id and movement_type for rollback idempotency protection
-- Extends b2c_inventory_movements to track rollbacks and prevent double rollback

-- Add order_id column for tracking which order the movement belongs to
ALTER TABLE b2c_inventory_movements
    ADD COLUMN IF NOT EXISTS order_id VARCHAR(255);

-- Add movement_type to distinguish between different operations
-- Types: 'decrement', 'rollback', 'adjustment', 'restock'
ALTER TABLE b2c_inventory_movements
    ADD COLUMN IF NOT EXISTS movement_type VARCHAR(20) DEFAULT 'adjustment'
    CHECK (movement_type IN ('decrement', 'rollback', 'adjustment', 'restock'));

-- Create unique index for idempotency: one rollback per order
-- This prevents duplicate rollbacks for the same order_id
CREATE UNIQUE INDEX IF NOT EXISTS idx_b2c_inventory_movements_rollback_idempotency
    ON b2c_inventory_movements(order_id, storefront_product_id)
    WHERE movement_type = 'rollback' AND order_id IS NOT NULL;

-- Create unique index for variant rollbacks
CREATE UNIQUE INDEX IF NOT EXISTS idx_b2c_inventory_movements_rollback_variant_idempotency
    ON b2c_inventory_movements(order_id, variant_id)
    WHERE movement_type = 'rollback' AND variant_id IS NOT NULL AND order_id IS NOT NULL;

-- Add indexes for performance
CREATE INDEX IF NOT EXISTS idx_b2c_inventory_movements_order_id
    ON b2c_inventory_movements(order_id)
    WHERE order_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_b2c_inventory_movements_movement_type
    ON b2c_inventory_movements(movement_type);

-- Update existing data: mark old records as 'adjustment' (default)
UPDATE b2c_inventory_movements
SET movement_type = 'adjustment'
WHERE movement_type IS NULL;

-- Comments for documentation
COMMENT ON COLUMN b2c_inventory_movements.order_id IS 'External order ID for tracking and idempotency';
COMMENT ON COLUMN b2c_inventory_movements.movement_type IS 'Type: decrement (sale), rollback (cancel), adjustment (manual), restock (add stock)';
COMMENT ON INDEX idx_b2c_inventory_movements_rollback_idempotency IS 'Prevents duplicate rollbacks for same order (idempotency protection)';

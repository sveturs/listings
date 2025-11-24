-- Rollback migration: Remove committed_at and released_at timestamps from inventory_reservations
-- Created: 2025-11-17

-- Drop indexes first
DROP INDEX IF EXISTS idx_inventory_reservations_committed_at;
DROP INDEX IF EXISTS idx_inventory_reservations_released_at;

-- Drop columns
ALTER TABLE inventory_reservations
DROP COLUMN IF EXISTS committed_at;

ALTER TABLE inventory_reservations
DROP COLUMN IF EXISTS released_at;

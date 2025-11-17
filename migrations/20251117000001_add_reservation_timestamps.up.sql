-- Migration: Add committed_at and released_at timestamps to inventory_reservations
-- Created: 2025-11-17
-- Purpose: Support tracking when reservations are committed or released for orders

-- Add committed_at column
ALTER TABLE inventory_reservations
ADD COLUMN committed_at TIMESTAMP WITH TIME ZONE;

-- Add released_at column
ALTER TABLE inventory_reservations
ADD COLUMN released_at TIMESTAMP WITH TIME ZONE;

-- Add index for faster queries on committed reservations
CREATE INDEX idx_inventory_reservations_committed_at
ON inventory_reservations(committed_at)
WHERE committed_at IS NOT NULL;

-- Add index for faster queries on released reservations
CREATE INDEX idx_inventory_reservations_released_at
ON inventory_reservations(released_at)
WHERE released_at IS NOT NULL;

-- Add comment for documentation
COMMENT ON COLUMN inventory_reservations.committed_at IS 'Timestamp when reservation was committed (order confirmed)';
COMMENT ON COLUMN inventory_reservations.released_at IS 'Timestamp when reservation was released (order cancelled or expired)';

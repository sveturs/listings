-- Migration: Add Store Status fields to storefronts table
-- Purpose: Enable store open/closed, vacation mode, accepting orders functionality
-- Date: 2026-01-17

-- Add store status fields
ALTER TABLE storefronts
ADD COLUMN IF NOT EXISTS vacation_mode BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS accepting_orders BOOLEAN DEFAULT TRUE,
ADD COLUMN IF NOT EXISTS vacation_start_date TIMESTAMP,
ADD COLUMN IF NOT EXISTS vacation_end_date TIMESTAMP,
ADD COLUMN IF NOT EXISTS auto_pause_when_out_of_stock BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS status_updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

-- Add index for querying active storefronts
CREATE INDEX IF NOT EXISTS idx_storefronts_vacation_mode ON storefronts(vacation_mode) WHERE vacation_mode = TRUE;
CREATE INDEX IF NOT EXISTS idx_storefronts_accepting_orders ON storefronts(accepting_orders);

-- Add comment
COMMENT ON COLUMN storefronts.vacation_mode IS 'Vacation mode - temporarily pause all orders';
COMMENT ON COLUMN storefronts.accepting_orders IS 'Whether storefront is accepting new orders';
COMMENT ON COLUMN storefronts.vacation_start_date IS 'Vacation period start date';
COMMENT ON COLUMN storefronts.vacation_end_date IS 'Vacation period end date';
COMMENT ON COLUMN storefronts.auto_pause_when_out_of_stock IS 'Automatically pause orders when all products out of stock';
COMMENT ON COLUMN storefronts.status_updated_at IS 'Last time store status was updated';

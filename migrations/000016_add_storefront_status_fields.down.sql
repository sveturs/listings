-- Rollback: Remove Store Status fields from storefronts table

-- Drop indexes
DROP INDEX IF EXISTS idx_storefronts_vacation_mode;
DROP INDEX IF EXISTS idx_storefronts_accepting_orders;

-- Drop columns
ALTER TABLE storefronts
DROP COLUMN IF EXISTS vacation_mode,
DROP COLUMN IF EXISTS accepting_orders,
DROP COLUMN IF EXISTS vacation_start_date,
DROP COLUMN IF EXISTS vacation_end_date,
DROP COLUMN IF EXISTS auto_pause_when_out_of_stock,
DROP COLUMN IF EXISTS status_updated_at;

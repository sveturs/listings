-- Migration rollback: Drop stock_reservations table
-- Author: Claude (Elite Architect)
-- Date: 2025-12-17

BEGIN;

-- Drop function
DROP FUNCTION IF EXISTS cleanup_expired_reservations();

-- Drop triggers
DROP TRIGGER IF EXISTS trigger_reservations_sync_quantity ON stock_reservations;
DROP TRIGGER IF EXISTS trigger_reservations_auto_expire ON stock_reservations;
DROP TRIGGER IF EXISTS trigger_reservations_updated_at ON stock_reservations;

-- Drop functions
DROP FUNCTION IF EXISTS sync_variant_reserved_quantity();
DROP FUNCTION IF EXISTS auto_expire_reservations();

-- Drop table
DROP TABLE IF EXISTS stock_reservations CASCADE;

COMMIT;

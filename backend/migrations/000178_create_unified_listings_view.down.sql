-- Migration: 000178_create_unified_listings_view (ROLLBACK)
-- Description: Удаляет VIEW unified_listings
-- Date: 2025-10-11

BEGIN;

-- Удалить VIEW
DROP VIEW IF EXISTS unified_listings;

COMMIT;

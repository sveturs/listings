-- Rollback Phase 7.5: Remove Storefronts table

-- Remove storefront_id from c2c_listings
ALTER TABLE IF EXISTS c2c_listings DROP COLUMN IF EXISTS storefront_id;

-- Drop trigger and function
DROP TRIGGER IF EXISTS trigger_storefronts_updated_at ON storefronts;
DROP FUNCTION IF EXISTS update_storefronts_updated_at();

-- Drop table
DROP TABLE IF EXISTS storefronts CASCADE;

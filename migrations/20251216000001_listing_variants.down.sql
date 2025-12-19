-- Rollback: Remove Listing Variants Table
-- Date: 2025-12-16

DROP TRIGGER IF EXISTS trigger_listing_variants_updated_at ON listing_variants;
DROP FUNCTION IF EXISTS update_listing_variants_updated_at();
DROP TABLE IF EXISTS listing_variants;

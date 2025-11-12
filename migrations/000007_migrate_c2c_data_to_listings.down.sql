-- ============================================================================
-- Phase 11.2: ROLLBACK - Remove migrated C2C data from listings
-- ============================================================================
-- This rollback migration removes C2C listings that were migrated from
-- c2c_listings table, along with all related data (images, locations,
-- attributes, favorites count updates).
-- ============================================================================
-- WARNING: This does NOT restore data to c2c_listings!
-- Make sure you have a database backup before running this migration.
-- ============================================================================

BEGIN;

-- ============================================================================
-- STEP 1: Find all C2C listings that were migrated in this migration
-- ============================================================================
-- We identify them by source_type='c2c' and created_at >= migration time
-- Since we can't easily determine exact migration time in rollback,
-- we'll use a safer approach: delete ALL c2c source_type listings
-- that don't have existing references in the original c2c_listings table

-- Create temporary table to store IDs to delete
CREATE TEMPORARY TABLE c2c_listings_to_delete AS
SELECT l.id
FROM listings l
WHERE l.source_type = 'c2c'
  AND l.is_deleted = false
  -- Extra safety: only delete if created recently (within last hour)
  -- Adjust this if migration takes longer
  AND l.created_at >= (CURRENT_TIMESTAMP - INTERVAL '1 hour');

-- Log what we're about to delete
DO $$
DECLARE
    delete_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO delete_count FROM c2c_listings_to_delete;
    RAISE NOTICE 'Phase 11.2 Rollback: Preparing to delete % C2C listings', delete_count;
END $$;

-- ============================================================================
-- STEP 2: Delete related data (cascades should handle this, but being explicit)
-- ============================================================================

-- Delete indexing queue entries
DELETE FROM indexing_queue
WHERE listing_id IN (SELECT id FROM c2c_listings_to_delete);

-- Delete listing attributes
DELETE FROM listing_attributes
WHERE listing_id IN (SELECT id FROM c2c_listings_to_delete);

-- Delete listing locations
DELETE FROM listing_locations
WHERE listing_id IN (SELECT id FROM c2c_listings_to_delete);

-- Delete listing images
DELETE FROM listing_images
WHERE listing_id IN (SELECT id FROM c2c_listings_to_delete);

-- Delete listing stats if any
DELETE FROM listing_stats
WHERE listing_id IN (SELECT id FROM c2c_listings_to_delete);

-- Delete listing tags if any
DELETE FROM listing_tags
WHERE listing_id IN (SELECT id FROM c2c_listings_to_delete);

-- ============================================================================
-- STEP 3: Delete the listings themselves
-- ============================================================================
DELETE FROM listings
WHERE id IN (SELECT id FROM c2c_listings_to_delete);

-- ============================================================================
-- STEP 4: Summary and validation
-- ============================================================================
DO $$
DECLARE
    remaining_c2c INTEGER;
BEGIN
    SELECT COUNT(*) INTO remaining_c2c
    FROM listings
    WHERE source_type = 'c2c' AND is_deleted = false;

    RAISE NOTICE '========================================';
    RAISE NOTICE 'Phase 11.2 Rollback Summary:';
    RAISE NOTICE '========================================';
    RAISE NOTICE 'Remaining C2C listings: %', remaining_c2c;
    RAISE NOTICE '========================================';
    RAISE NOTICE 'Phase 11.2: Rollback completed!';
    RAISE NOTICE 'NOTE: Original c2c_listings data was NOT restored!';
    RAISE NOTICE 'Restore from backup if needed.';
    RAISE NOTICE '========================================';
END $$;

COMMIT;

-- ============================================================================
-- IMPORTANT NOTES:
-- ============================================================================
-- 1. This rollback does NOT restore data to c2c_listings table
-- 2. Use database backup to fully restore if needed
-- 3. c2c_images, c2c_favorites still reference c2c_listings (unchanged)
-- 4. Only deletes listings created within last hour (safety measure)
-- 5. All related data is deleted via explicit DELETEs
-- ============================================================================

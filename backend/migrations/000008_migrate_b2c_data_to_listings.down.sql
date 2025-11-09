-- Rollback Migration: Revert B2C data migration from listings back to b2c_products
-- Phase: 11.3 Rollback
-- Date: 2025-11-06

BEGIN;

-- ============================================================================
-- Step 1: Remove indexing_queue entries for B2C listings
-- ============================================================================

DELETE FROM indexing_queue
WHERE listing_id IN (
    SELECT id FROM listings WHERE source_type = 'b2c'
);

-- ============================================================================
-- Step 2: Remove listing_attributes for B2C listings
-- ============================================================================

DELETE FROM listing_attributes
WHERE listing_id IN (
    SELECT id FROM listings WHERE source_type = 'b2c'
);

-- ============================================================================
-- Step 3: Remove listing_locations for B2C listings
-- ============================================================================

DELETE FROM listing_locations
WHERE listing_id IN (
    SELECT id FROM listings WHERE source_type = 'b2c'
);

-- ============================================================================
-- Step 4: Remove B2C listings from listings table
-- ============================================================================

DELETE FROM listings
WHERE source_type = 'b2c';

-- ============================================================================
-- Step 5: Output rollback summary
-- ============================================================================

DO $$
DECLARE
    v_remaining_b2c INTEGER;
BEGIN
    SELECT COUNT(*) INTO v_remaining_b2c FROM listings WHERE source_type = 'b2c';

    RAISE NOTICE '========================================';
    RAISE NOTICE 'Phase 11.3 Rollback Summary';
    RAISE NOTICE '========================================';
    RAISE NOTICE 'Remaining B2C listings: %', v_remaining_b2c;
    RAISE NOTICE 'B2C products table: unchanged (preserved)';
    RAISE NOTICE '========================================';

    IF v_remaining_b2c > 0 THEN
        RAISE EXCEPTION 'Rollback incomplete: % B2C listings still exist', v_remaining_b2c;
    END IF;
END $$;

COMMIT;

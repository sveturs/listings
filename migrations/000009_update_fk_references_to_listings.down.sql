-- Migration Rollback: 000009_update_fk_references_to_listings
-- Purpose: Rollback FK references from listings back to legacy tables (c2c_listings, b2c_products)
-- Phase: 11.4 Rollback
-- Date: 2025-11-06

BEGIN;

-- ============================================================================
-- PART 1: Restore c2c_favorites to original state
-- ============================================================================

-- Drop new FK constraint
ALTER TABLE c2c_favorites
DROP CONSTRAINT IF EXISTS fk_c2c_favorites_listing_id;

-- Restore original data from backup
TRUNCATE TABLE c2c_favorites;
INSERT INTO c2c_favorites
SELECT * FROM c2c_favorites_backup_phase_11_4;

-- Restore original FK constraint pointing to c2c_listings
ALTER TABLE c2c_favorites
ADD CONSTRAINT fk_c2c_favorites_listing_id
FOREIGN KEY (listing_id) REFERENCES c2c_listings(id) ON DELETE CASCADE;

-- Verify restoration
DO $$
DECLARE
    restored_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO restored_count FROM c2c_favorites;
    RAISE NOTICE 'c2c_favorites restored: % records', restored_count;
END $$;

-- ============================================================================
-- PART 2: Drop unified inventory_movements table
-- ============================================================================

-- Drop indexes first
DROP INDEX IF EXISTS idx_inventory_movements_listing_id;
DROP INDEX IF EXISTS idx_inventory_movements_created_at;
DROP INDEX IF EXISTS idx_inventory_movements_movement_type;
DROP INDEX IF EXISTS idx_inventory_movements_user_id;

-- Drop the table (data will be lost, but it's migrated from b2c_inventory_movements)
DROP TABLE IF EXISTS inventory_movements;

-- ============================================================================
-- PART 3: Restore FK constraints for empty legacy tables
-- ============================================================================

-- Restore FK for c2c_listing_variants
ALTER TABLE c2c_listing_variants
ADD CONSTRAINT fk_c2c_listing_variants_listing_id
FOREIGN KEY (listing_id) REFERENCES c2c_listings(id) ON DELETE CASCADE;

-- Restore FK for c2c_orders
ALTER TABLE c2c_orders
ADD CONSTRAINT fk_c2c_orders_listing_id
FOREIGN KEY (listing_id) REFERENCES c2c_listings(id) ON DELETE CASCADE;

-- ============================================================================
-- PART 4: Restore c2c_images FK constraint
-- ============================================================================

-- Restore FK constraint (even though data is in listing_images)
ALTER TABLE c2c_images
ADD CONSTRAINT fk_c2c_images_listing_id
FOREIGN KEY (listing_id) REFERENCES c2c_listings(id) ON DELETE CASCADE;

DO $$
BEGIN
    RAISE NOTICE 'Restored FK constraints to c2c_listing_variants, c2c_orders, c2c_images';
    RAISE NOTICE 'Dropped inventory_movements table';
END $$;

-- ============================================================================
-- PART 5: Drop c2c_favorites indexes created in up migration
-- ============================================================================

-- Note: We keep basic indexes, only drop unique constraint added in up migration
DROP INDEX IF EXISTS idx_c2c_favorites_unique;

-- ============================================================================
-- FINAL VERIFICATION
-- ============================================================================

DO $$
DECLARE
    fk_to_c2c_count INTEGER;
    fk_to_b2c_count INTEGER;
    favorites_count INTEGER;
    inventory_exists BOOLEAN;
BEGIN
    -- Count FK references to c2c_listings (should be restored)
    SELECT COUNT(*) INTO fk_to_c2c_count
    FROM information_schema.table_constraints tc
    JOIN information_schema.constraint_column_usage ccu ON ccu.constraint_name = tc.constraint_name
    WHERE tc.constraint_type = 'FOREIGN KEY'
    AND ccu.table_name = 'c2c_listings';

    -- Count FK references to b2c_products (should still have original FK)
    SELECT COUNT(*) INTO fk_to_b2c_count
    FROM information_schema.table_constraints tc
    JOIN information_schema.constraint_column_usage ccu ON ccu.constraint_name = tc.constraint_name
    WHERE tc.constraint_type = 'FOREIGN KEY'
    AND ccu.table_name = 'b2c_products';

    -- Check c2c_favorites count
    SELECT COUNT(*) INTO favorites_count FROM c2c_favorites;

    -- Check if inventory_movements still exists
    SELECT EXISTS (
        SELECT 1 FROM information_schema.tables
        WHERE table_name = 'inventory_movements'
    ) INTO inventory_exists;

    RAISE NOTICE '=== Rollback Summary ===';
    RAISE NOTICE 'FK references to c2c_listings: % (expected: 4)', fk_to_c2c_count;
    RAISE NOTICE 'FK references to b2c_products: % (expected: 1)', fk_to_b2c_count;
    RAISE NOTICE 'c2c_favorites count: % (expected: 2)', favorites_count;
    RAISE NOTICE 'inventory_movements exists: % (expected: false)', inventory_exists;

    IF fk_to_c2c_count < 4 THEN
        RAISE WARNING 'Expected 4 FK references to c2c_listings, but found %', fk_to_c2c_count;
    END IF;

    IF inventory_exists THEN
        RAISE WARNING 'inventory_movements table still exists after rollback!';
    END IF;

    RAISE NOTICE '=== Phase 11.4 Rollback Complete ===';
END $$;

-- ============================================================================
-- CLEANUP: Drop backup table (optional - comment out if you want to keep it)
-- ============================================================================

-- DROP TABLE IF EXISTS c2c_favorites_backup_phase_11_4;
-- RAISE NOTICE 'Backup table dropped';

COMMIT;

-- Migration: 000010_drop_legacy_tables
-- Description: Drop legacy C2C/B2C tables after successful data migration to unified schema
-- Created: 2025-11-06
-- Phase: 11.5 - Final cleanup

BEGIN;

-- =============================================================================
-- ANALYSIS SUMMARY (2025-11-06)
-- =============================================================================
--
-- SAFELY REMOVABLE (data migrated to unified tables):
-- 1. ✅ c2c_listings (4 records) → listings
-- 2. ✅ b2c_products (7 records) → listings
-- 3. ✅ c2c_images (1 record) → listing_images
-- 4. ✅ b2c_inventory_movements (3 records) → inventory_movements
-- 5. ✅ c2c_listing_variants (0 records, FK already removed)
-- 6. ✅ c2c_favorites_backup_phase_11_4 (temporary backup table)
-- 7. ✅ b2c_product_variants (0 records, empty)
-- 8. ✅ c2c_orders (0 records, empty)
-- 9. ✅ c2c_categories (77 records, only referenced by c2c_listings which is being dropped)
--
-- PRESERVED (actively used by application):
-- 10. ❌ c2c_favorites (2 records) - ACTIVE! Used by listings microservice
-- 11. ❌ c2c_chats (2 records) - ACTIVE! Used by svetu backend
-- 12. ❌ c2c_messages (8 records) - ACTIVE! Related to c2c_chats
--
-- =============================================================================

-- Step 1: Drop backup table from Phase 11.4
DROP TABLE IF EXISTS c2c_favorites_backup_phase_11_4 CASCADE;

-- Step 2: Drop empty variant/order tables
DROP TABLE IF EXISTS c2c_listing_variants CASCADE;
DROP TABLE IF EXISTS c2c_orders CASCADE;
DROP TABLE IF EXISTS b2c_product_variants CASCADE;

-- Step 3: Drop tables with migrated data (order matters due to FKs)
-- First drop tables that reference others
DROP TABLE IF EXISTS c2c_images CASCADE;
DROP TABLE IF EXISTS b2c_inventory_movements CASCADE;

-- Then drop main listing tables
DROP TABLE IF EXISTS c2c_listings CASCADE;
DROP TABLE IF EXISTS b2c_products CASCADE;

-- Finally drop c2c_categories (only referenced by c2c_listings which is now dropped)
DROP TABLE IF EXISTS c2c_categories CASCADE;

-- Step 4: Add migration marker comment
COMMENT ON TABLE listings IS 'Unified listings table (C2C + B2C merged). Legacy tables dropped in Phase 11.5 (2025-11-06)';

-- Step 5: Verify remaining legacy tables
DO $$
DECLARE
    remaining_tables TEXT[];
BEGIN
    SELECT ARRAY_AGG(table_name::TEXT ORDER BY table_name)
    INTO remaining_tables
    FROM information_schema.tables
    WHERE table_schema = 'public'
    AND (table_name LIKE 'c2c_%' OR table_name LIKE 'b2c_%');

    IF remaining_tables IS NOT NULL THEN
        RAISE NOTICE 'Remaining legacy tables (actively used): %', ARRAY_TO_STRING(remaining_tables, ', ');
        -- Expected: c2c_chats, c2c_favorites, c2c_messages
    ELSE
        RAISE NOTICE 'All legacy tables dropped successfully';
    END IF;
END $$;

COMMIT;

-- Post-migration verification queries:
--
-- 1. Check remaining legacy tables:
-- SELECT table_name FROM information_schema.tables
-- WHERE table_schema = 'public'
-- AND (table_name LIKE 'c2c_%' OR table_name LIKE 'b2c_%')
-- ORDER BY table_name;
--
-- 2. Verify unified data:
-- SELECT source_type, COUNT(*) FROM listings GROUP BY source_type;
-- SELECT COUNT(*) FROM listing_images;
-- SELECT COUNT(*) FROM inventory_movements;
--
-- 3. Check FK constraints:
-- SELECT tc.table_name, ccu.table_name AS foreign_table_name
-- FROM information_schema.table_constraints tc
-- JOIN information_schema.constraint_column_usage ccu ON ccu.constraint_name = tc.constraint_name
-- WHERE tc.constraint_type = 'FOREIGN KEY'
-- ORDER BY tc.table_name;

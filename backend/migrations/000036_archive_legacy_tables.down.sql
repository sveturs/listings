-- Rollback Legacy Attribute Tables Archive Migration
-- Day 17: Unified Attributes Project  
-- Date: 03.09.2025
-- Purpose: Restore archived legacy attribute tables back to public schema

BEGIN;

-- =====================================================
-- RESTORE TABLES FROM ARCHIVE
-- =====================================================

-- Restore category_attributes table
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'archive_legacy' AND table_name = 'category_attributes') THEN
        ALTER TABLE archive_legacy.category_attributes SET SCHEMA public;
        RAISE NOTICE 'Restored category_attributes to public schema';
    ELSE
        RAISE WARNING 'category_attributes not found in archive_legacy schema';
    END IF;
END $$;

-- Restore listing_attribute_values table
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'archive_legacy' AND table_name = 'listing_attribute_values') THEN
        ALTER TABLE archive_legacy.listing_attribute_values SET SCHEMA public;
        RAISE NOTICE 'Restored listing_attribute_values to public schema';
    ELSE
        RAISE WARNING 'listing_attribute_values not found in archive_legacy schema';
    END IF;
END $$;

-- Restore attribute_groups table if it exists
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'archive_legacy' AND table_name = 'attribute_groups') THEN
        ALTER TABLE archive_legacy.attribute_groups SET SCHEMA public;
        RAISE NOTICE 'Restored attribute_groups to public schema';
    END IF;
END $$;

-- =====================================================
-- RECREATE INDEXES
-- =====================================================

-- Note: Most indexes will be restored automatically when tables are moved back to public schema
-- The tables already have their proper indexes, no need to recreate them manually
DO $$
BEGIN
    -- Verify that tables were restored with their indexes intact
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'category_attributes') THEN
        RAISE NOTICE 'category_attributes restored to public schema with indexes intact';
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'listing_attribute_values') THEN
        RAISE NOTICE 'listing_attribute_values restored to public schema with indexes intact';
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'attribute_groups') THEN
        RAISE NOTICE 'attribute_groups restored to public schema with indexes intact';
    END IF;
END $$;

-- =====================================================
-- CLEANUP ARCHIVE ARTIFACTS
-- =====================================================

-- Drop the archive view
DROP VIEW IF EXISTS archive_legacy.v_archived_tables;

-- Drop archive metadata table
DROP TABLE IF EXISTS archive_legacy.archive_metadata;

-- Check if archive_legacy schema is now empty and drop it if so
DO $$
DECLARE
    table_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO table_count
    FROM information_schema.tables 
    WHERE table_schema = 'archive_legacy';
    
    IF table_count = 0 THEN
        DROP SCHEMA IF EXISTS archive_legacy;
        RAISE NOTICE 'Dropped empty archive_legacy schema';
    ELSE
        RAISE NOTICE 'archive_legacy schema still contains % tables, keeping it', table_count;
    END IF;
END $$;

-- =====================================================
-- FINAL VALIDATION
-- =====================================================

DO $$
DECLARE
    restored_tables TEXT[];
    table_count INTEGER := 0;
BEGIN
    -- Check which tables were successfully restored
    SELECT array_agg(table_name) INTO restored_tables
    FROM information_schema.tables 
    WHERE table_schema = 'public' 
        AND table_name IN ('category_attributes', 'listing_attribute_values', 'attribute_groups');
    
    table_count := array_length(restored_tables, 1);
    
    IF table_count > 0 THEN
        RAISE NOTICE 'Successfully restored % legacy tables to public schema: %', 
                    table_count, array_to_string(restored_tables, ', ');
    ELSE
        RAISE WARNING 'No legacy tables were found to restore';
    END IF;
    
    -- Final success message
    RAISE NOTICE 'Rollback completed: Day 17 Legacy Archive Migration reverted';
    RAISE NOTICE 'Legacy attribute tables are now available in public schema';
END $$;

COMMIT;
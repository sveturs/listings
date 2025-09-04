-- Archive Legacy Attribute Tables Migration
-- Day 17: Unified Attributes Project
-- Date: 03.09.2025
-- Purpose: Archive legacy attribute system tables that are no longer used after unified attributes migration

BEGIN;

-- Create archive schema if it doesn't exist
CREATE SCHEMA IF NOT EXISTS archive_legacy;

-- Add comment to archive schema
COMMENT ON SCHEMA archive_legacy IS 'Archived legacy attribute tables from Day 17 of unified attributes project (03.09.2025)';

-- =====================================================
-- ARCHIVE LEGACY ATTRIBUTE TABLES
-- =====================================================

-- Archive category_attributes (old marketplace system)
-- This table contained the old category-specific attributes
ALTER TABLE category_attributes SET SCHEMA archive_legacy;
COMMENT ON TABLE archive_legacy.category_attributes IS 'Legacy category attributes table - archived Day 17 (85 records preserved)';

-- Archive listing_attribute_values (old attribute values)
-- This table contained values for the old attribute system
ALTER TABLE listing_attribute_values SET SCHEMA archive_legacy;
COMMENT ON TABLE archive_legacy.listing_attribute_values IS 'Legacy listing attribute values - archived Day 17 (15 records preserved)';

-- Archive attribute_groups if it exists (empty but part of old system)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'attribute_groups') THEN
        ALTER TABLE attribute_groups SET SCHEMA archive_legacy;
        COMMENT ON TABLE archive_legacy.attribute_groups IS 'Legacy attribute groups table - archived Day 17 (empty table)';
    END IF;
END $$;

-- =====================================================
-- CREATE ARCHIVE METADATA TABLE
-- =====================================================

-- Create metadata table to track what was archived and when
CREATE TABLE archive_legacy.archive_metadata (
    id SERIAL PRIMARY KEY,
    table_name VARCHAR(255) NOT NULL,
    original_schema VARCHAR(255) NOT NULL DEFAULT 'public',
    archive_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    record_count INTEGER,
    reason TEXT,
    project_phase VARCHAR(100) DEFAULT 'unified_attributes_day_17'
);

-- Insert metadata for archived tables
INSERT INTO archive_legacy.archive_metadata (table_name, record_count, reason) VALUES
('category_attributes', 85, 'Legacy marketplace category attributes system - replaced by unified_attributes'),
('listing_attribute_values', 15, 'Legacy attribute values - replaced by unified_attribute_values');

-- Insert metadata for attribute_groups if it was archived
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'archive_legacy' AND table_name = 'attribute_groups') THEN
        INSERT INTO archive_legacy.archive_metadata (table_name, record_count, reason) VALUES
        ('attribute_groups', 0, 'Legacy attribute groups - empty table from old system');
    END IF;
END $$;

-- =====================================================
-- CLEANUP RELATED OBJECTS
-- =====================================================

-- Drop any indexes that might reference the archived tables
-- (Most should have been automatically moved with the tables, but check for orphaned ones)
DROP INDEX IF EXISTS idx_category_attributes_category_id;
DROP INDEX IF EXISTS idx_listing_attribute_values_listing_id;
DROP INDEX IF EXISTS idx_attribute_groups_name;

-- =====================================================
-- VERIFICATION QUERIES
-- =====================================================

-- Create a view to easily check archived tables
CREATE OR REPLACE VIEW archive_legacy.v_archived_tables AS
SELECT 
    t.table_name,
    t.table_schema,
    COALESCE(am.record_count, 0) as record_count,
    am.archive_date,
    am.reason
FROM information_schema.tables t
LEFT JOIN archive_legacy.archive_metadata am ON t.table_name = am.table_name
WHERE t.table_schema = 'archive_legacy'
    AND t.table_type = 'BASE TABLE'
    AND t.table_name != 'archive_metadata'
ORDER BY am.archive_date DESC;

-- Grant permissions for archive schema access
-- Allow read-only access to developers for historical data analysis
GRANT USAGE ON SCHEMA archive_legacy TO PUBLIC;
GRANT SELECT ON ALL TABLES IN SCHEMA archive_legacy TO PUBLIC;

-- =====================================================
-- FINAL VALIDATION
-- =====================================================

-- Ensure the archived tables are no longer in public schema
DO $$
DECLARE
    remaining_tables TEXT;
BEGIN
    SELECT string_agg(table_name, ', ') INTO remaining_tables
    FROM information_schema.tables 
    WHERE table_schema = 'public' 
        AND table_name IN ('category_attributes', 'listing_attribute_values', 'attribute_groups');
    
    IF remaining_tables IS NOT NULL THEN
        RAISE EXCEPTION 'Migration failed: The following tables still exist in public schema: %', remaining_tables;
    END IF;
    
    -- Log successful migration
    RAISE NOTICE 'Legacy tables successfully archived to archive_legacy schema';
    RAISE NOTICE 'Archived tables: category_attributes (85 records), listing_attribute_values (15 records)';
    RAISE NOTICE 'Migration completed: Day 17 - Unified Attributes Project';
END $$;

COMMIT;
-- Migration: Update entity_type values from marketplace_/storefronts_ to c2c_/b2c_
-- This ensures consistency across all tables that reference entity types

-- Step 1: Add new enum values to geo_source_type
ALTER TYPE geo_source_type ADD VALUE IF NOT EXISTS 'c2c_listing';
ALTER TYPE geo_source_type ADD VALUE IF NOT EXISTS 'b2c_store';
ALTER TYPE geo_source_type ADD VALUE IF NOT EXISTS 'b2c_product';

-- Step 2: Update unified_geo table
UPDATE unified_geo
SET source_type = 'c2c_listing'
WHERE source_type = 'marketplace_listing';

UPDATE unified_geo
SET source_type = 'b2c_store'
WHERE source_type = 'storefront';

UPDATE unified_geo
SET source_type = 'b2c_product'
WHERE source_type = 'storefront_product';

-- Note: Cannot remove old enum values as PostgreSQL doesn't support removing enum values
-- Old values (marketplace_listing, storefront, storefront_product) will remain in enum definition
-- but won't be used in data

-- Step 3: Update reviews table (if entity_type exists)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'reviews' AND column_name = 'entity_type') THEN
        UPDATE reviews SET entity_type = 'c2c_listing' WHERE entity_type = 'marketplace_listing';
        UPDATE reviews SET entity_type = 'b2c_store' WHERE entity_type = 'storefront';
        UPDATE reviews SET entity_type = 'b2c_product' WHERE entity_type = 'storefront_product';
    END IF;
END $$;

-- Step 4: Update translations table (if entity_type exists)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'translations' AND column_name = 'entity_type') THEN
        UPDATE translations SET entity_type = 'c2c_category' WHERE entity_type = 'marketplace_category';
        UPDATE translations SET entity_type = 'c2c_listing' WHERE entity_type = 'marketplace_listing';
        UPDATE translations SET entity_type = 'b2c_store' WHERE entity_type = 'storefront';
        UPDATE translations SET entity_type = 'b2c_product' WHERE entity_type = 'storefront_product';
    END IF;
END $$;

-- Step 5: Update notifications table (if it exists)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'notifications') THEN
        IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'notifications' AND column_name = 'entity_type') THEN
            UPDATE notifications SET entity_type = 'c2c_listing' WHERE entity_type = 'marketplace_listing';
            UPDATE notifications SET entity_type = 'b2c_store' WHERE entity_type = 'storefront';
            UPDATE notifications SET entity_type = 'b2c_product' WHERE entity_type = 'storefront_product';
        END IF;
    END IF;
END $$;

-- Step 6: Update activity_logs table (if it exists)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'activity_logs') THEN
        IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'activity_logs' AND column_name = 'entity_type') THEN
            UPDATE activity_logs SET entity_type = 'c2c_listing' WHERE entity_type = 'marketplace_listing';
            UPDATE activity_logs SET entity_type = 'b2c_store' WHERE entity_type = 'storefront';
            UPDATE activity_logs SET entity_type = 'b2c_product' WHERE entity_type = 'storefront_product';
        END IF;
    END IF;
END $$;

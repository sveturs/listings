-- ============================================================================
-- Phase 11.2: Migrate C2C data from c2c_listings to unified listings table
-- ============================================================================
-- This migration moves existing C2C listings from the legacy c2c_listings
-- table to the unified listings table, along with all related data
-- (images, locations, attributes, favorites).
-- ============================================================================

BEGIN;

-- ============================================================================
-- STEP 1: Create temporary mapping table for old_id -> new_id
-- ============================================================================
CREATE TEMPORARY TABLE c2c_id_mapping (
    old_c2c_id INTEGER PRIMARY KEY,
    new_listing_id BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- STEP 2: Migrate c2c_listings to listings table
-- ============================================================================
-- Insert c2c_listings into listings with proper field mapping
-- We'll add old c2c_id as a temporary attribute to track mapping
INSERT INTO listings (
    user_id,
    storefront_id,
    title,
    description,
    price,
    currency,
    category_id,
    status,
    visibility,
    quantity,
    sku,
    views_count,
    favorites_count,
    created_at,
    updated_at,
    published_at,
    deleted_at,
    is_deleted,
    source_type
)
SELECT
    -- Map integer IDs to bigint
    c.user_id::bigint,
    c.storefront_id::bigint,

    -- Basic fields
    c.title,
    c.description,
    c.price,

    -- Default currency for legacy C2C
    COALESCE(c.metadata->>'currency', 'RSD') as currency,

    -- Category
    c.category_id::bigint,

    -- Status mapping: 'active' stays 'active', empty/NULL becomes 'draft'
    CASE
        WHEN TRIM(COALESCE(c.status, '')) = '' THEN 'draft'
        WHEN c.status = 'active' THEN 'active'
        WHEN c.status = 'sold' THEN 'sold'
        WHEN c.status = 'inactive' THEN 'inactive'
        ELSE 'draft'
    END as status,

    -- Visibility: default public for legacy C2C
    'public' as visibility,

    -- C2C items are usually single quantity
    1 as quantity,

    -- C2C doesn't have SKU
    NULL as sku,

    -- Views count
    COALESCE(c.views_count, 0),

    -- Favorites count (will be updated later from c2c_favorites)
    0 as favorites_count,

    -- Timestamps
    c.created_at AT TIME ZONE 'UTC' as created_at,
    c.updated_at AT TIME ZONE 'UTC' as updated_at,

    -- Published date: use created_at if status is active
    CASE
        WHEN c.status = 'active' THEN c.created_at AT TIME ZONE 'UTC'
        ELSE NULL
    END as published_at,

    -- Deleted flag
    NULL as deleted_at,
    false as is_deleted,

    -- Source type: CRITICAL - mark as C2C
    'c2c' as source_type
FROM c2c_listings c
ORDER BY c.id;

-- Build ID mapping by matching on unique characteristics
-- We use title + created_at as it should be unique enough
INSERT INTO c2c_id_mapping (old_c2c_id, new_listing_id)
SELECT
    c.id as old_c2c_id,
    l.id as new_listing_id
FROM c2c_listings c
JOIN listings l ON
    l.title = c.title
    AND l.source_type = 'c2c'
    AND l.created_at = c.created_at AT TIME ZONE 'UTC'
    AND l.user_id = c.user_id::bigint
ORDER BY c.id;

-- Log migration results
DO $$
DECLARE
    migrated_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO migrated_count FROM c2c_id_mapping;
    RAISE NOTICE 'Phase 11.2: Migrated % c2c_listings to listings table', migrated_count;
END $$;

-- ============================================================================
-- STEP 3: Migrate location data to listing_locations
-- ============================================================================
INSERT INTO listing_locations (
    listing_id,
    latitude,
    longitude,
    address_line1,
    address_line2,
    city,
    country,
    postal_code,
    created_at,
    updated_at
)
SELECT
    m.new_listing_id,
    c.latitude,
    c.longitude,
    c.location as address_line1,
    NULL as address_line2,
    COALESCE(
        c.address_multilingual->>'city_sr',
        c.address_multilingual->>'city_en',
        c.address_city
    ) as city,
    COALESCE(
        c.address_multilingual->>'country_sr',
        c.address_multilingual->>'country_en',
        c.address_country
    ) as country,
    c.address_multilingual->>'postal_code' as postal_code,
    c.created_at AT TIME ZONE 'UTC',
    c.updated_at AT TIME ZONE 'UTC'
FROM c2c_listings c
JOIN c2c_id_mapping m ON c.id = m.old_c2c_id
WHERE c.latitude IS NOT NULL
   OR c.longitude IS NOT NULL
   OR c.location IS NOT NULL
   OR c.address_city IS NOT NULL
   OR c.address_country IS NOT NULL;

-- Log location migration
DO $$
DECLARE
    locations_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO locations_count
    FROM listing_locations ll
    JOIN c2c_id_mapping m ON ll.listing_id = m.new_listing_id;
    RAISE NOTICE 'Phase 11.2: Migrated % location records', locations_count;
END $$;

-- ============================================================================
-- STEP 4: Migrate custom attributes to listing_attributes
-- ============================================================================
-- Migrate condition attribute
INSERT INTO listing_attributes (
    listing_id,
    attribute_key,
    attribute_value,
    created_at
)
SELECT
    m.new_listing_id,
    'condition' as attribute_key,
    c.condition as attribute_value,
    c.created_at AT TIME ZONE 'UTC'
FROM c2c_listings c
JOIN c2c_id_mapping m ON c.id = m.old_c2c_id
WHERE c.condition IS NOT NULL;

-- Migrate external_id as attribute
INSERT INTO listing_attributes (
    listing_id,
    attribute_key,
    attribute_value,
    created_at
)
SELECT
    m.new_listing_id,
    'external_id' as attribute_key,
    c.external_id as attribute_value,
    c.created_at AT TIME ZONE 'UTC'
FROM c2c_listings c
JOIN c2c_id_mapping m ON c.id = m.old_c2c_id
WHERE c.external_id IS NOT NULL;

-- Migrate original_language as attribute
INSERT INTO listing_attributes (
    listing_id,
    attribute_key,
    attribute_value,
    created_at
)
SELECT
    m.new_listing_id,
    'original_language' as attribute_key,
    COALESCE((c.metadata->>'original_language'), 'sr') as attribute_value,
    c.created_at AT TIME ZONE 'UTC'
FROM c2c_listings c
JOIN c2c_id_mapping m ON c.id = m.old_c2c_id;

-- Migrate show_on_map as attribute
INSERT INTO listing_attributes (
    listing_id,
    attribute_key,
    attribute_value,
    created_at
)
SELECT
    m.new_listing_id,
    'show_on_map' as attribute_key,
    CASE WHEN c.show_on_map THEN 'true' ELSE 'false' END as attribute_value,
    c.created_at AT TIME ZONE 'UTC'
FROM c2c_listings c
JOIN c2c_id_mapping m ON c.id = m.old_c2c_id;

-- Migrate metadata as JSON attribute
INSERT INTO listing_attributes (
    listing_id,
    attribute_key,
    attribute_value,
    created_at
)
SELECT
    m.new_listing_id,
    'legacy_metadata' as attribute_key,
    c.metadata::text as attribute_value,
    c.created_at AT TIME ZONE 'UTC'
FROM c2c_listings c
JOIN c2c_id_mapping m ON c.id = m.old_c2c_id
WHERE c.metadata IS NOT NULL AND c.metadata != 'null'::jsonb;

-- Migrate address_multilingual as JSON attribute
INSERT INTO listing_attributes (
    listing_id,
    attribute_key,
    attribute_value,
    created_at
)
SELECT
    m.new_listing_id,
    'address_multilingual' as attribute_key,
    c.address_multilingual::text as attribute_value,
    c.created_at AT TIME ZONE 'UTC'
FROM c2c_listings c
JOIN c2c_id_mapping m ON c.id = m.old_c2c_id
WHERE c.address_multilingual IS NOT NULL AND c.address_multilingual != 'null'::jsonb;

-- Log attributes migration
DO $$
DECLARE
    attributes_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO attributes_count
    FROM listing_attributes la
    JOIN c2c_id_mapping m ON la.listing_id = m.new_listing_id;
    RAISE NOTICE 'Phase 11.2: Migrated % attribute records', attributes_count;
END $$;

-- ============================================================================
-- STEP 5: Migrate c2c_images to listing_images
-- ============================================================================
-- Temporarily disable FK constraint to allow update
ALTER TABLE c2c_images DROP CONSTRAINT IF EXISTS fk_c2c_images_listing_id;

-- Insert images into listing_images
INSERT INTO listing_images (
    listing_id,
    url,
    storage_path,
    thumbnail_url,
    display_order,
    is_primary,
    width,
    height,
    file_size,
    mime_type,
    created_at,
    updated_at
)
SELECT
    m.new_listing_id,
    COALESCE(ci.public_url, ci.file_path) as url,
    ci.file_path as storage_path,
    NULL as thumbnail_url, -- Will be generated later
    0 as display_order, -- Will be updated if needed
    COALESCE(ci.is_main, false) as is_primary,
    NULL as width, -- Unknown for legacy images
    NULL as height, -- Unknown for legacy images
    ci.file_size::bigint as file_size,
    ci.content_type as mime_type,
    ci.created_at AT TIME ZONE 'UTC',
    ci.created_at AT TIME ZONE 'UTC' as updated_at
FROM c2c_images ci
JOIN c2c_id_mapping m ON ci.listing_id = m.old_c2c_id;

-- Restore FK constraint
ALTER TABLE c2c_images ADD CONSTRAINT fk_c2c_images_listing_id
    FOREIGN KEY (listing_id) REFERENCES c2c_listings(id)
    ON UPDATE CASCADE ON DELETE CASCADE;

-- Log images migration
DO $$
DECLARE
    images_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO images_count
    FROM listing_images li
    JOIN c2c_id_mapping m ON li.listing_id = m.new_listing_id;
    RAISE NOTICE 'Phase 11.2: Migrated % image records', images_count;
END $$;

-- ============================================================================
-- STEP 6: Update favorites_count in listings from c2c_favorites
-- ============================================================================
UPDATE listings l
SET favorites_count = fav_counts.count
FROM (
    SELECT m.new_listing_id, COUNT(*) as count
    FROM c2c_favorites cf
    JOIN c2c_id_mapping m ON cf.listing_id = m.old_c2c_id
    GROUP BY m.new_listing_id
) fav_counts
WHERE l.id = fav_counts.new_listing_id;

-- Log favorites update
DO $$
DECLARE
    favorites_updated INTEGER;
BEGIN
    SELECT COUNT(*) INTO favorites_updated
    FROM listings l
    JOIN c2c_id_mapping m ON l.id = m.new_listing_id
    WHERE l.favorites_count > 0;
    RAISE NOTICE 'Phase 11.2: Updated favorites_count for % listings', favorites_updated;
END $$;

-- ============================================================================
-- STEP 7: Create indexing queue entries for new listings
-- ============================================================================
INSERT INTO indexing_queue (
    listing_id,
    operation,
    status,
    created_at
)
SELECT
    m.new_listing_id,
    'index' as operation,
    'pending' as status,
    CURRENT_TIMESTAMP
FROM c2c_id_mapping m;

-- Log indexing queue
DO $$
DECLARE
    queue_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO queue_count
    FROM indexing_queue iq
    JOIN c2c_id_mapping m ON iq.listing_id = m.new_listing_id;
    RAISE NOTICE 'Phase 11.2: Created % indexing queue entries', queue_count;
END $$;

-- ============================================================================
-- STEP 8: Final validation and summary
-- ============================================================================
DO $$
DECLARE
    total_c2c INTEGER;
    total_migrated INTEGER;
    total_locations INTEGER;
    total_attributes INTEGER;
    total_images INTEGER;
BEGIN
    SELECT COUNT(*) INTO total_c2c FROM c2c_listings;
    SELECT COUNT(*) INTO total_migrated FROM c2c_id_mapping;

    SELECT COUNT(*) INTO total_locations
    FROM listing_locations ll
    JOIN c2c_id_mapping m ON ll.listing_id = m.new_listing_id;

    SELECT COUNT(*) INTO total_attributes
    FROM listing_attributes la
    JOIN c2c_id_mapping m ON la.listing_id = m.new_listing_id;

    SELECT COUNT(*) INTO total_images
    FROM listing_images li
    JOIN c2c_id_mapping m ON li.listing_id = m.new_listing_id;

    RAISE NOTICE '========================================';
    RAISE NOTICE 'Phase 11.2 Migration Summary:';
    RAISE NOTICE '========================================';
    RAISE NOTICE 'C2C listings in source table: %', total_c2c;
    RAISE NOTICE 'C2C listings migrated: %', total_migrated;
    RAISE NOTICE 'Location records created: %', total_locations;
    RAISE NOTICE 'Attribute records created: %', total_attributes;
    RAISE NOTICE 'Image records migrated: %', total_images;
    RAISE NOTICE '========================================';

    -- Validate migration success
    IF total_c2c != total_migrated THEN
        RAISE EXCEPTION 'Migration validation failed: expected % listings, got %',
            total_c2c, total_migrated;
    END IF;

    RAISE NOTICE 'Phase 11.2: Migration completed successfully!';
END $$;

COMMIT;

-- ============================================================================
-- IMPORTANT NOTES:
-- ============================================================================
-- 1. c2c_listings table is NOT dropped - will be dropped in Phase 11.5
-- 2. c2c_images FK still points to c2c_listings (will be handled later)
-- 3. c2c_favorites still reference old listing_id (will be migrated later)
-- 4. c2c_orders are NOT touched (may contain active orders)
-- 5. All migrated listings have source_type='c2c' for identification
-- ============================================================================

-- Migration: Migrate B2C data from b2c_products to unified listings table
-- Phase: 11.3
-- Author: Migration from b2c_products legacy table
-- Date: 2025-11-06

BEGIN;

-- ============================================================================
-- Step 1: Create temporary ID mapping table
-- ============================================================================
CREATE TEMP TABLE b2c_id_mapping (
    old_id INTEGER PRIMARY KEY,
    new_id BIGINT NOT NULL,
    storefront_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL
);

-- ============================================================================
-- Step 2: Migrate b2c_products → listings
-- ============================================================================

-- Insert B2C products into listings with source_type='b2c'
WITH inserted_listings AS (
    INSERT INTO listings (
        storefront_id,
        user_id,
        title,
        description,
        price,
        currency,
        category_id,
        sku,
        quantity,
        status,
        visibility,
        views_count,
        favorites_count,
        source_type,
        created_at,
        updated_at,
        published_at,
        deleted_at,
        is_deleted
    )
    SELECT
        bp.storefront_id::BIGINT,
        -- Use storefront owner as user_id, fallback to 1 if storefront doesn't exist
        COALESCE(s.user_id::BIGINT, 1::BIGINT) as user_id,
        bp.name as title,
        bp.description,
        bp.price,
        bp.currency,
        bp.category_id::BIGINT,
        bp.sku,
        COALESCE(bp.stock_quantity, 0) as quantity,
        CASE
            WHEN bp.is_active THEN 'active'::VARCHAR
            ELSE 'inactive'::VARCHAR
        END as status,
        'public'::VARCHAR as visibility,
        COALESCE(bp.view_count, 0) as views_count,
        0 as favorites_count, -- B2C doesn't have favorites initially
        'b2c'::VARCHAR as source_type,
        bp.created_at,
        bp.updated_at,
        CASE
            WHEN bp.is_active THEN bp.created_at
            ELSE NULL
        END as published_at,
        NULL as deleted_at,
        FALSE as is_deleted
    FROM b2c_products bp
    LEFT JOIN storefronts s ON s.id = bp.storefront_id
    ORDER BY bp.id
    RETURNING id, (SELECT bp.id FROM b2c_products bp WHERE bp.name = listings.title AND bp.storefront_id = listings.storefront_id LIMIT 1) as old_id, storefront_id, user_id
)
-- Save ID mapping
INSERT INTO b2c_id_mapping (old_id, new_id, storefront_id, user_id)
SELECT old_id, id, storefront_id, user_id FROM inserted_listings;

-- ============================================================================
-- Step 3: Migrate listing_locations for products with individual location
-- ============================================================================

INSERT INTO listing_locations (
    listing_id,
    country,
    city,
    postal_code,
    address_line1,
    address_line2,
    latitude,
    longitude,
    created_at,
    updated_at
)
SELECT
    m.new_id as listing_id,
    'RS'::VARCHAR(100) as country, -- Default for Serbia
    NULL as city, -- Not available in b2c_products
    NULL as postal_code,
    bp.individual_address as address_line1,
    NULL as address_line2,
    bp.individual_latitude::NUMERIC(10,8) as latitude,
    bp.individual_longitude::NUMERIC(11,8) as longitude,
    bp.created_at,
    bp.updated_at
FROM b2c_products bp
INNER JOIN b2c_id_mapping m ON m.old_id = bp.id
WHERE bp.has_individual_location = TRUE
  AND bp.individual_address IS NOT NULL;

-- ============================================================================
-- Step 4: Migrate listing_attributes for additional B2C-specific fields
-- ============================================================================

-- Attribute: barcode
INSERT INTO listing_attributes (
    listing_id,
    attribute_key,
    attribute_value,
    created_at
)
SELECT
    m.new_id as listing_id,
    'barcode'::VARCHAR(100) as attribute_key,
    bp.barcode as attribute_value,
    bp.created_at
FROM b2c_products bp
INNER JOIN b2c_id_mapping m ON m.old_id = bp.id
WHERE bp.barcode IS NOT NULL AND bp.barcode != '';

-- Attribute: stock_status
INSERT INTO listing_attributes (
    listing_id,
    attribute_key,
    attribute_value,
    created_at
)
SELECT
    m.new_id as listing_id,
    'stock_status'::VARCHAR(100) as attribute_key,
    bp.stock_status as attribute_value,
    bp.created_at
FROM b2c_products bp
INNER JOIN b2c_id_mapping m ON m.old_id = bp.id
WHERE bp.stock_status IS NOT NULL AND bp.stock_status != '';

-- Attribute: sold_count
INSERT INTO listing_attributes (
    listing_id,
    attribute_key,
    attribute_value,
    created_at
)
SELECT
    m.new_id as listing_id,
    'sold_count'::VARCHAR(100) as attribute_key,
    bp.sold_count::TEXT as attribute_value,
    bp.created_at
FROM b2c_products bp
INNER JOIN b2c_id_mapping m ON m.old_id = bp.id
WHERE bp.sold_count IS NOT NULL AND bp.sold_count > 0;

-- Attribute: has_variants (boolean flag)
INSERT INTO listing_attributes (
    listing_id,
    attribute_key,
    attribute_value,
    created_at
)
SELECT
    m.new_id as listing_id,
    'has_variants'::VARCHAR(100) as attribute_key,
    bp.has_variants::TEXT as attribute_value,
    bp.created_at
FROM b2c_products bp
INNER JOIN b2c_id_mapping m ON m.old_id = bp.id
WHERE bp.has_variants IS NOT NULL;

-- Expand JSONB attributes into individual listing_attributes
-- Only if attributes JSONB is not empty
INSERT INTO listing_attributes (
    listing_id,
    attribute_key,
    attribute_value,
    created_at
)
SELECT
    m.new_id as listing_id,
    jsonb_key::VARCHAR(100) as attribute_key,
    jsonb_value::TEXT as attribute_value,
    bp.created_at
FROM b2c_products bp
INNER JOIN b2c_id_mapping m ON m.old_id = bp.id
CROSS JOIN LATERAL jsonb_each_text(bp.attributes) AS attrs(jsonb_key, jsonb_value)
WHERE bp.attributes IS NOT NULL
  AND jsonb_typeof(bp.attributes) = 'object'
  AND jsonb_key NOT IN ('barcode', 'stock_status', 'sold_count', 'has_variants'); -- Avoid duplicates

-- ============================================================================
-- Step 5: Create indexing_queue entries for all migrated listings
-- ============================================================================

INSERT INTO indexing_queue (
    listing_id,
    operation,
    status,
    retry_count,
    error_message,
    created_at,
    updated_at
)
SELECT
    m.new_id as listing_id,
    'index'::VARCHAR(20) as operation,
    'pending'::VARCHAR(20) as status,
    0 as retry_count,
    NULL as error_message,
    NOW() as created_at,
    NOW() as updated_at
FROM b2c_id_mapping m;

-- ============================================================================
-- Step 6: Output migration summary
-- ============================================================================

DO $$
DECLARE
    v_migrated_count INTEGER;
    v_locations_count INTEGER;
    v_attributes_count INTEGER;
    v_queue_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO v_migrated_count FROM b2c_id_mapping;
    SELECT COUNT(*) INTO v_locations_count
    FROM listing_locations
    WHERE listing_id IN (SELECT new_id FROM b2c_id_mapping);

    SELECT COUNT(*) INTO v_attributes_count
    FROM listing_attributes
    WHERE listing_id IN (SELECT new_id FROM b2c_id_mapping);

    SELECT COUNT(*) INTO v_queue_count
    FROM indexing_queue
    WHERE listing_id IN (SELECT new_id FROM b2c_id_mapping);

    RAISE NOTICE '========================================';
    RAISE NOTICE 'Phase 11.3 Migration Summary';
    RAISE NOTICE '========================================';
    RAISE NOTICE 'Migrated B2C products: %', v_migrated_count;
    RAISE NOTICE 'Created locations: %', v_locations_count;
    RAISE NOTICE 'Created attributes: %', v_attributes_count;
    RAISE NOTICE 'Queued for indexing: %', v_queue_count;
    RAISE NOTICE '========================================';
END $$;

-- Keep temporary table for verification (it will be dropped at session end)
-- SELECT * FROM b2c_id_mapping ORDER BY old_id;

COMMIT;

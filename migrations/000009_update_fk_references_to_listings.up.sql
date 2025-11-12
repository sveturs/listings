-- Migration: 000009_update_fk_references_to_listings
-- Purpose: Update FK references from legacy tables (c2c_listings, b2c_products) to unified listings table
-- Phase: 11.4
-- Date: 2025-11-06

BEGIN;

-- ============================================================================
-- PART 1: Create temporary mapping tables for ID translation
-- ============================================================================

-- Drop temp tables if they exist (from previous migrations in same session)
DROP TABLE IF EXISTS c2c_id_mapping;
DROP TABLE IF EXISTS b2c_id_mapping;

-- C2C mapping table (c2c_listings.id -> listings.id)
CREATE TEMPORARY TABLE c2c_id_mapping AS
SELECT
    cl.id::integer as old_id,
    l.id as new_id
FROM c2c_listings cl
JOIN listings l ON
    l.title = cl.title
    AND l.source_type = 'c2c'
    AND l.created_at::date = cl.created_at::date;

-- B2C mapping table (b2c_products.id -> listings.id)
CREATE TEMPORARY TABLE b2c_id_mapping AS
SELECT
    bp.id::integer as old_id,
    l.id as new_id
FROM b2c_products bp
JOIN listings l ON
    l.title = bp.name
    AND l.source_type = 'b2c'
    AND l.created_at::date = bp.created_at::date;

-- Verify mappings
DO $$
DECLARE
    c2c_count INTEGER;
    b2c_count INTEGER;
    c2c_listings_exist INTEGER;
    b2c_products_exist INTEGER;
BEGIN
    SELECT COUNT(*) INTO c2c_count FROM c2c_id_mapping;
    SELECT COUNT(*) INTO b2c_count FROM b2c_id_mapping;

    -- Check if source tables exist and have data
    SELECT COUNT(*) INTO c2c_listings_exist
    FROM information_schema.tables
    WHERE table_name = 'c2c_listings';

    SELECT COUNT(*) INTO b2c_products_exist
    FROM information_schema.tables
    WHERE table_name = 'b2c_products';

    RAISE NOTICE 'C2C ID mappings created: %', c2c_count;
    RAISE NOTICE 'B2C ID mappings created: %', b2c_count;

    -- Only fail if tables exist but mappings don't
    IF c2c_listings_exist > 0 AND b2c_products_exist > 0 THEN
        IF c2c_count = 0 AND b2c_count = 0 THEN
            -- Check if source tables actually have data
            EXECUTE 'SELECT COUNT(*) FROM c2c_listings' INTO c2c_listings_exist;
            EXECUTE 'SELECT COUNT(*) FROM b2c_products' INTO b2c_products_exist;

            IF c2c_listings_exist > 0 OR b2c_products_exist > 0 THEN
                RAISE EXCEPTION 'Source tables have data (c2c: %, b2c: %) but no ID mappings created!',
                    c2c_listings_exist, b2c_products_exist;
            ELSE
                RAISE NOTICE 'Empty database - migration proceeds without data migration';
            END IF;
        END IF;
    ELSE
        RAISE NOTICE 'Legacy tables do not exist - skipping migration (likely fresh database)';
    END IF;
END $$;

-- ============================================================================
-- PART 2: Migrate c2c_favorites (2 records)
-- ============================================================================

-- Store old data for rollback
CREATE TABLE IF NOT EXISTS c2c_favorites_backup_phase_11_4 AS
SELECT * FROM c2c_favorites;

-- Drop old FK constraint
ALTER TABLE c2c_favorites
DROP CONSTRAINT IF EXISTS fk_c2c_favorites_listing_id;

-- Update listing_id using C2C mapping
UPDATE c2c_favorites f
SET listing_id = m.new_id::integer
FROM c2c_id_mapping m
WHERE f.listing_id = m.old_id;

-- Add new FK constraint to listings table
ALTER TABLE c2c_favorites
ADD CONSTRAINT fk_c2c_favorites_listing_id
FOREIGN KEY (listing_id) REFERENCES listings(id) ON DELETE CASCADE;

-- Verify update
DO $$
DECLARE
    updated_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO updated_count
    FROM c2c_favorites f
    JOIN listings l ON l.id = f.listing_id::bigint
    WHERE l.source_type = 'c2c';

    RAISE NOTICE 'c2c_favorites updated and verified: % records', updated_count;
END $$;

-- ============================================================================
-- PART 3: Create unified inventory_movements table
-- ============================================================================

-- Create new unified inventory movements table
CREATE TABLE IF NOT EXISTS inventory_movements (
    id BIGSERIAL PRIMARY KEY,
    listing_id BIGINT NOT NULL REFERENCES listings(id) ON DELETE CASCADE,
    variant_id BIGINT,
    movement_type VARCHAR(50) NOT NULL, -- 'in', 'out', 'adjustment'
    quantity INTEGER NOT NULL,
    reason VARCHAR(255),
    notes TEXT,
    user_id BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Metadata for audit/tracking
    metadata JSONB,

    CONSTRAINT inventory_movements_quantity_check CHECK (quantity <> 0)
);

-- Create indexes for performance
CREATE INDEX idx_inventory_movements_listing_id ON inventory_movements(listing_id);
CREATE INDEX idx_inventory_movements_created_at ON inventory_movements(created_at DESC);
CREATE INDEX idx_inventory_movements_movement_type ON inventory_movements(movement_type);
CREATE INDEX idx_inventory_movements_user_id ON inventory_movements(user_id) WHERE user_id IS NOT NULL;

-- Migrate data from b2c_inventory_movements (3 records)
INSERT INTO inventory_movements (
    listing_id,
    variant_id,
    movement_type,
    quantity,
    reason,
    notes,
    user_id,
    created_at,
    metadata
)
SELECT
    m.new_id,                    -- new listing_id from mapping
    bim.variant_id,              -- preserve variant_id
    bim.type,                    -- map 'type' to 'movement_type'
    bim.quantity,
    bim.reason,
    bim.notes,
    bim.user_id,
    bim.created_at AT TIME ZONE 'UTC', -- ensure timestamptz
    jsonb_build_object(
        'migrated_from', 'b2c_inventory_movements',
        'old_b2c_product_id', bim.storefront_product_id,
        'migration_date', NOW()
    )
FROM b2c_inventory_movements bim
JOIN b2c_id_mapping m ON m.old_id = bim.storefront_product_id;

-- Verify migration
DO $$
DECLARE
    migrated_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO migrated_count FROM inventory_movements;
    RAISE NOTICE 'inventory_movements created and populated: % records', migrated_count;

    IF migrated_count = 0 THEN
        RAISE WARNING 'No records migrated to inventory_movements!';
    END IF;
END $$;

-- ============================================================================
-- PART 4: Drop FK constraints for empty legacy tables
-- ============================================================================

-- c2c_listing_variants (0 records) - just drop FK
ALTER TABLE c2c_listing_variants
DROP CONSTRAINT IF EXISTS fk_c2c_listing_variants_listing_id;

-- c2c_orders (0 records) - just drop FK
ALTER TABLE c2c_orders
DROP CONSTRAINT IF EXISTS fk_c2c_orders_listing_id;

-- ============================================================================
-- PART 5: Handle c2c_images (already migrated to listing_images)
-- ============================================================================

-- c2c_images data was already migrated to listing_images in Phase 11.2
-- We can safely drop the FK constraint as the table will be removed in Phase 11.5
ALTER TABLE c2c_images
DROP CONSTRAINT IF EXISTS fk_c2c_images_listing_id;

DO $$
BEGIN
    RAISE NOTICE 'Dropped FK constraints from c2c_listing_variants, c2c_orders, c2c_images';
END $$;

-- ============================================================================
-- PART 6: Ensure c2c_favorites index exists
-- ============================================================================

CREATE INDEX IF NOT EXISTS idx_c2c_favorites_listing_id ON c2c_favorites(listing_id);
CREATE INDEX IF NOT EXISTS idx_c2c_favorites_user_id ON c2c_favorites(user_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_c2c_favorites_unique ON c2c_favorites(user_id, listing_id);

-- ============================================================================
-- FINAL VERIFICATION
-- ============================================================================

DO $$
DECLARE
    fk_to_c2c_count INTEGER;
    fk_to_b2c_count INTEGER;
    fk_to_listings_count INTEGER;
BEGIN
    -- Count FK references to legacy tables (should be 1: b2c_inventory_movements)
    SELECT COUNT(*) INTO fk_to_c2c_count
    FROM information_schema.table_constraints tc
    JOIN information_schema.constraint_column_usage ccu ON ccu.constraint_name = tc.constraint_name
    WHERE tc.constraint_type = 'FOREIGN KEY'
    AND ccu.table_name = 'c2c_listings';

    -- Should be 1 (b2c_inventory_movements - we keep it for rollback safety)
    SELECT COUNT(*) INTO fk_to_b2c_count
    FROM information_schema.table_constraints tc
    JOIN information_schema.constraint_column_usage ccu ON ccu.constraint_name = tc.constraint_name
    WHERE tc.constraint_type = 'FOREIGN KEY'
    AND ccu.table_name = 'b2c_products';

    -- Count FK references to listings table (should be >= 7)
    SELECT COUNT(*) INTO fk_to_listings_count
    FROM information_schema.table_constraints tc
    JOIN information_schema.constraint_column_usage ccu ON ccu.constraint_name = tc.constraint_name
    WHERE tc.constraint_type = 'FOREIGN KEY'
    AND ccu.table_name = 'listings';

    RAISE NOTICE '=== Migration Summary ===';
    RAISE NOTICE 'FK references to c2c_listings: %', fk_to_c2c_count;
    RAISE NOTICE 'FK references to b2c_products: %', fk_to_b2c_count;
    RAISE NOTICE 'FK references to listings: %', fk_to_listings_count;

    IF fk_to_c2c_count > 0 THEN
        RAISE WARNING 'Still have % FK references to c2c_listings', fk_to_c2c_count;
    END IF;

    IF fk_to_listings_count < 7 THEN
        RAISE WARNING 'Expected at least 7 FK references to listings, but found %', fk_to_listings_count;
    END IF;

    RAISE NOTICE '=== Phase 11.4 Complete ===';
END $$;

COMMIT;

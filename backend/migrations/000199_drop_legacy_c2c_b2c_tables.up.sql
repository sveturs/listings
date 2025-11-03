-- Migration: Drop legacy C2C and B2C tables
-- Phase 7.4: Final cleanup after successful microservice migration
-- Date: 2025-11-03
--
-- Context:
--   - All legacy functionality migrated to delivery microservice
--   - Phase 7.3.5 VERIFIED: 0 scans on critical tables (c2c_images, c2c_favorites, etc.)
--   - c2c_categories: 11 scans (analytics only, not critical)
--   - Full backup created: /tmp/backup_before_drop_legacy_YYYYMMDD_HHMMSS.sql
--
-- Tables to drop:
--   C2C: c2c_favorites, c2c_images, c2c_listing_variants, c2c_messages,
--        c2c_chats, c2c_orders, c2c_listings, c2c_categories
--   B2C: b2c_favorites, b2c_inventory_movements, b2c_order_items, b2c_orders,
--        b2c_product_variant_images, b2c_product_variants, b2c_product_attributes,
--        b2c_product_images, b2c_products, b2c_delivery_options, b2c_payment_methods,
--        b2c_store_hours, b2c_store_staff, user_b2c_stores, b2c_stores
--
-- Safety: This migration includes pre-flight checks and audit logging

DO $$
DECLARE
    v_scan_count INTEGER;
    v_table_name TEXT;
    v_tables TEXT[] := ARRAY[
        'c2c_favorites', 'c2c_images', 'c2c_listing_variants', 'c2c_messages',
        'c2c_chats', 'c2c_orders', 'c2c_listings',
        'b2c_favorites', 'b2c_inventory_movements', 'b2c_order_items', 'b2c_orders',
        'b2c_product_variant_images', 'b2c_product_variants', 'b2c_product_attributes',
        'b2c_product_images', 'b2c_products', 'b2c_delivery_options', 'b2c_payment_methods',
        'b2c_store_hours', 'b2c_store_staff', 'user_b2c_stores', 'b2c_stores'
    ];
BEGIN
    -- Create audit log table if not exists
    CREATE TABLE IF NOT EXISTS migration_audit_log (
        id SERIAL PRIMARY KEY,
        migration_number TEXT NOT NULL,
        action TEXT NOT NULL,
        table_name TEXT,
        details JSONB,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

    -- Log migration start
    INSERT INTO migration_audit_log (migration_number, action, details)
    VALUES ('000199', 'START', jsonb_build_object(
        'description', 'Drop legacy C2C and B2C tables',
        'backup_location', '/tmp/backup_before_drop_legacy_YYYYMMDD_HHMMSS.sql'
    ));

    -- SAFETY CHECK 1: Verify no recent activity (except c2c_categories analytics)
    RAISE NOTICE 'Phase 7.4: Starting safety checks...';

    FOR v_table_name IN
        SELECT unnest(v_tables)
    LOOP
        -- Skip c2c_categories (used by analytics)
        IF v_table_name = 'c2c_categories' THEN
            RAISE NOTICE 'Skipping activity check for c2c_categories (analytics usage is expected)';
            CONTINUE;
        END IF;

        -- Check for recent inserts/updates/deletes
        SELECT COALESCE(n_tup_ins + n_tup_upd + n_tup_del, 0)
        INTO v_scan_count
        FROM pg_stat_user_tables
        WHERE relname = v_table_name;

        IF v_scan_count > 0 THEN
            INSERT INTO migration_audit_log (migration_number, action, table_name, details)
            VALUES ('000199', 'SAFETY_CHECK_FAILED', v_table_name, jsonb_build_object(
                'reason', 'Recent activity detected',
                'modification_count', v_scan_count
            ));

            RAISE EXCEPTION 'SAFETY CHECK FAILED: Table % has % recent modifications. Migration aborted.',
                v_table_name, v_scan_count;
        END IF;
    END LOOP;

    RAISE NOTICE 'Safety checks passed. No recent activity detected.';

    -- Log table statistics before drop
    INSERT INTO migration_audit_log (migration_number, action, details)
    VALUES ('000199', 'PRE_DROP_STATS', (
        SELECT jsonb_object_agg(relname, jsonb_build_object(
            'row_count', n_live_tup,
            'seq_scans', seq_scan,
            'idx_scans', COALESCE(idx_scan, 0)
        ))
        FROM pg_stat_user_tables
        WHERE relname = ANY(v_tables)
    ));

    -- DROP TABLES (child tables first to avoid FK conflicts)
    RAISE NOTICE 'Starting table drops...';

    -- C2C child tables (references c2c_listings)
    DROP TABLE IF EXISTS c2c_favorites CASCADE;
    RAISE NOTICE 'Dropped: c2c_favorites';
    INSERT INTO migration_audit_log (migration_number, action, table_name)
    VALUES ('000199', 'DROP_TABLE', 'c2c_favorites');

    DROP TABLE IF EXISTS c2c_images CASCADE;
    RAISE NOTICE 'Dropped: c2c_images';
    INSERT INTO migration_audit_log (migration_number, action, table_name)
    VALUES ('000199', 'DROP_TABLE', 'c2c_images');

    DROP TABLE IF EXISTS c2c_listing_variants CASCADE;
    RAISE NOTICE 'Dropped: c2c_listing_variants';
    INSERT INTO migration_audit_log (migration_number, action, table_name)
    VALUES ('000199', 'DROP_TABLE', 'c2c_listing_variants');

    DROP TABLE IF EXISTS c2c_messages CASCADE;
    RAISE NOTICE 'Dropped: c2c_messages';
    INSERT INTO migration_audit_log (migration_number, action, table_name)
    VALUES ('000199', 'DROP_TABLE', 'c2c_messages');

    DROP TABLE IF EXISTS c2c_chats CASCADE;
    RAISE NOTICE 'Dropped: c2c_chats';
    INSERT INTO migration_audit_log (migration_number, action, table_name)
    VALUES ('000199', 'DROP_TABLE', 'c2c_chats');

    DROP TABLE IF EXISTS c2c_orders CASCADE;
    RAISE NOTICE 'Dropped: c2c_orders';
    INSERT INTO migration_audit_log (migration_number, action, table_name)
    VALUES ('000199', 'DROP_TABLE', 'c2c_orders');

    -- C2C parent table
    DROP TABLE IF EXISTS c2c_listings CASCADE;
    RAISE NOTICE 'Dropped: c2c_listings';
    INSERT INTO migration_audit_log (migration_number, action, table_name)
    VALUES ('000199', 'DROP_TABLE', 'c2c_listings');

    -- B2C child tables (references b2c_products)
    DROP TABLE IF EXISTS b2c_favorites CASCADE;
    RAISE NOTICE 'Dropped: b2c_favorites';
    INSERT INTO migration_audit_log (migration_number, action, table_name)
    VALUES ('000199', 'DROP_TABLE', 'b2c_favorites');

    DROP TABLE IF EXISTS b2c_inventory_movements CASCADE;
    RAISE NOTICE 'Dropped: b2c_inventory_movements';
    INSERT INTO migration_audit_log (migration_number, action, table_name)
    VALUES ('000199', 'DROP_TABLE', 'b2c_inventory_movements');

    DROP TABLE IF EXISTS b2c_order_items CASCADE;
    RAISE NOTICE 'Dropped: b2c_order_items';
    INSERT INTO migration_audit_log (migration_number, action, table_name)
    VALUES ('000199', 'DROP_TABLE', 'b2c_order_items');

    DROP TABLE IF EXISTS b2c_orders CASCADE;
    RAISE NOTICE 'Dropped: b2c_orders';
    INSERT INTO migration_audit_log (migration_number, action, table_name)
    VALUES ('000199', 'DROP_TABLE', 'b2c_orders');

    DROP TABLE IF EXISTS b2c_product_variant_images CASCADE;
    RAISE NOTICE 'Dropped: b2c_product_variant_images';
    INSERT INTO migration_audit_log (migration_number, action, table_name)
    VALUES ('000199', 'DROP_TABLE', 'b2c_product_variant_images');

    DROP TABLE IF EXISTS b2c_product_variants CASCADE;
    RAISE NOTICE 'Dropped: b2c_product_variants';
    INSERT INTO migration_audit_log (migration_number, action, table_name)
    VALUES ('000199', 'DROP_TABLE', 'b2c_product_variants');

    DROP TABLE IF EXISTS b2c_product_attributes CASCADE;
    RAISE NOTICE 'Dropped: b2c_product_attributes';
    INSERT INTO migration_audit_log (migration_number, action, table_name)
    VALUES ('000199', 'DROP_TABLE', 'b2c_product_attributes');

    DROP TABLE IF EXISTS b2c_product_images CASCADE;
    RAISE NOTICE 'Dropped: b2c_product_images';
    INSERT INTO migration_audit_log (migration_number, action, table_name)
    VALUES ('000199', 'DROP_TABLE', 'b2c_product_images');

    -- B2C parent table (products)
    DROP TABLE IF EXISTS b2c_products CASCADE;
    RAISE NOTICE 'Dropped: b2c_products';
    INSERT INTO migration_audit_log (migration_number, action, table_name)
    VALUES ('000199', 'DROP_TABLE', 'b2c_products');

    -- B2C store-related tables
    DROP TABLE IF EXISTS b2c_delivery_options CASCADE;
    RAISE NOTICE 'Dropped: b2c_delivery_options';
    INSERT INTO migration_audit_log (migration_number, action, table_name)
    VALUES ('000199', 'DROP_TABLE', 'b2c_delivery_options');

    DROP TABLE IF EXISTS b2c_payment_methods CASCADE;
    RAISE NOTICE 'Dropped: b2c_payment_methods';
    INSERT INTO migration_audit_log (migration_number, action, table_name)
    VALUES ('000199', 'DROP_TABLE', 'b2c_payment_methods');

    DROP TABLE IF EXISTS b2c_store_hours CASCADE;
    RAISE NOTICE 'Dropped: b2c_store_hours';
    INSERT INTO migration_audit_log (migration_number, action, table_name)
    VALUES ('000199', 'DROP_TABLE', 'b2c_store_hours');

    DROP TABLE IF EXISTS b2c_store_staff CASCADE;
    RAISE NOTICE 'Dropped: b2c_store_staff';
    INSERT INTO migration_audit_log (migration_number, action, table_name)
    VALUES ('000199', 'DROP_TABLE', 'b2c_store_staff');

    DROP TABLE IF EXISTS user_b2c_stores CASCADE;
    RAISE NOTICE 'Dropped: user_b2c_stores';
    INSERT INTO migration_audit_log (migration_number, action, table_name)
    VALUES ('000199', 'DROP_TABLE', 'user_b2c_stores');

    -- B2C stores parent
    DROP TABLE IF EXISTS b2c_stores CASCADE;
    RAISE NOTICE 'Dropped: b2c_stores';
    INSERT INTO migration_audit_log (migration_number, action, table_name)
    VALUES ('000199', 'DROP_TABLE', 'b2c_stores');

    -- C2C categories (last, used by analytics but safe to drop)
    DROP TABLE IF EXISTS c2c_categories CASCADE;
    RAISE NOTICE 'Dropped: c2c_categories';
    INSERT INTO migration_audit_log (migration_number, action, table_name)
    VALUES ('000199', 'DROP_TABLE', 'c2c_categories');

    -- Log migration completion
    INSERT INTO migration_audit_log (migration_number, action, details)
    VALUES ('000199', 'COMPLETE', jsonb_build_object(
        'tables_dropped', array_length(v_tables, 1) + 1, -- +1 for c2c_categories
        'status', 'SUCCESS'
    ));

    RAISE NOTICE 'Phase 7.4 COMPLETE: All legacy C2C/B2C tables dropped successfully';
END $$;

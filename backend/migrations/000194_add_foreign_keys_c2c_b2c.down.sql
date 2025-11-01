-- Migration Rollback: 000194_add_foreign_keys_c2c_b2c
-- Description: Remove all Foreign Key constraints added in UP migration
-- Date: 2025-10-30
--
-- This DOWN migration safely removes all FK constraints in reverse order
-- to avoid dependency issues during rollback.

-- ==============================================================================
-- SECTION 7: ADDITIONAL B2C TABLES (reverse order)
-- ==============================================================================

-- Drop b2c_inventory_movements FK (if exists)
ALTER TABLE IF EXISTS b2c_inventory_movements
DROP CONSTRAINT IF EXISTS fk_b2c_inventory_movements_product_id;

-- Drop b2c_product_variant_images FK (if exists)
ALTER TABLE IF EXISTS b2c_product_variant_images
DROP CONSTRAINT IF EXISTS fk_b2c_product_variant_images_variant_id;

-- ==============================================================================
-- SECTION 6: B2C ORDERS AND ORDER ITEMS (reverse order)
-- ==============================================================================

ALTER TABLE b2c_order_items
DROP CONSTRAINT IF EXISTS fk_b2c_order_items_variant_id;

ALTER TABLE b2c_order_items
DROP CONSTRAINT IF EXISTS fk_b2c_order_items_product_id;

ALTER TABLE b2c_order_items
DROP CONSTRAINT IF EXISTS fk_b2c_order_items_order_id;

ALTER TABLE b2c_orders
DROP CONSTRAINT IF EXISTS fk_b2c_orders_storefront_id;

-- ==============================================================================
-- SECTION 5: B2C PRODUCTS - Child Tables
-- ==============================================================================

ALTER TABLE b2c_favorites
DROP CONSTRAINT IF EXISTS fk_b2c_favorites_product_id;

ALTER TABLE b2c_product_variants
DROP CONSTRAINT IF EXISTS fk_b2c_product_variants_product_id;

ALTER TABLE b2c_product_images
DROP CONSTRAINT IF EXISTS fk_b2c_product_images_product_id;

-- ==============================================================================
-- SECTION 4: B2C PRODUCTS - Parent Table Foreign Keys
-- ==============================================================================

ALTER TABLE b2c_products
DROP CONSTRAINT IF EXISTS fk_b2c_products_category_id;

ALTER TABLE b2c_products
DROP CONSTRAINT IF EXISTS fk_b2c_products_storefront_id;

-- ==============================================================================
-- SECTION 3: C2C ORDERS
-- ==============================================================================

ALTER TABLE c2c_orders
DROP CONSTRAINT IF EXISTS fk_c2c_orders_listing_id;

-- ==============================================================================
-- SECTION 2: C2C LISTINGS - Child Tables
-- ==============================================================================

ALTER TABLE c2c_listing_variants
DROP CONSTRAINT IF EXISTS fk_c2c_listing_variants_listing_id;

ALTER TABLE c2c_favorites
DROP CONSTRAINT IF EXISTS fk_c2c_favorites_listing_id;

ALTER TABLE c2c_images
DROP CONSTRAINT IF EXISTS fk_c2c_images_listing_id;

-- ==============================================================================
-- SECTION 1: C2C LISTINGS - Parent Table Foreign Keys
-- ==============================================================================

ALTER TABLE c2c_listings
DROP CONSTRAINT IF EXISTS fk_c2c_listings_storefront_id;

ALTER TABLE c2c_listings
DROP CONSTRAINT IF EXISTS fk_c2c_listings_category_id;

-- ==============================================================================
-- ROLLBACK COMPLETE
-- ==============================================================================
-- All Foreign Key constraints have been removed.
-- Database is back to the state before migration 000194.
-- ==============================================================================

-- Migration: 000179_optimize_unified_listings_indexes (DOWN)
-- Description: Откат оптимизированных индексов для unified listings VIEW

-- ============================================
-- DROP C2C LISTINGS INDEXES
-- ============================================

DROP INDEX IF EXISTS idx_c2c_listings_active_created;
DROP INDEX IF EXISTS idx_c2c_listings_category_active;
DROP INDEX IF EXISTS idx_c2c_listings_price;
DROP INDEX IF EXISTS idx_c2c_listings_location;
DROP INDEX IF EXISTS idx_c2c_listings_text_search;

-- ============================================
-- DROP B2C PRODUCTS INDEXES
-- ============================================

DROP INDEX IF EXISTS idx_b2c_products_active_created;
DROP INDEX IF EXISTS idx_b2c_products_category_active;
DROP INDEX IF EXISTS idx_b2c_products_price;
DROP INDEX IF EXISTS idx_b2c_products_storefront;
DROP INDEX IF EXISTS idx_b2c_products_text_search;

-- ============================================
-- DROP IMAGES INDEXES
-- ============================================

DROP INDEX IF EXISTS idx_c2c_images_listing_main;
DROP INDEX IF EXISTS idx_b2c_images_product_main;

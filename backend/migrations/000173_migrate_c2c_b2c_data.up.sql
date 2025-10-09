-- ============================================================================
-- МИГРАЦИЯ ДАННЫХ: Копирование из marketplace/storefront в c2c/b2c
-- КРИТИЧНО: Сохраняем ID для связей!
-- ============================================================================

BEGIN;

-- ============================================================================
-- C2C DATA MIGRATION
-- ============================================================================

-- 1. Categories (сначала - для FK)
INSERT INTO c2c_categories SELECT * FROM marketplace_categories;

-- 2. Listings (зависит от categories)
INSERT INTO c2c_listings SELECT * FROM marketplace_listings;

-- 3. Images (зависит от listings)
INSERT INTO c2c_images SELECT * FROM marketplace_images;

-- 4. Listing Variants (зависит от listings)
INSERT INTO c2c_listing_variants SELECT * FROM marketplace_listing_variants;

-- 5. Chats (зависит от listings)
INSERT INTO c2c_chats SELECT * FROM marketplace_chats;

-- 6. Messages (зависит от chats)
INSERT INTO c2c_messages SELECT * FROM marketplace_messages;

-- 7. Favorites (зависит от listings)
INSERT INTO c2c_favorites SELECT * FROM marketplace_favorites;

-- 8. Orders (зависит от listings)
INSERT INTO c2c_orders SELECT * FROM marketplace_orders;

-- ============================================================================
-- B2C DATA MIGRATION
-- ============================================================================

-- 1. Stores (основная таблица - сначала)
INSERT INTO b2c_stores SELECT * FROM storefronts;

-- 2. User-Store links (зависит от stores)
INSERT INTO user_b2c_stores SELECT * FROM user_storefronts;

-- 3. Products (зависит от stores)
INSERT INTO b2c_products SELECT * FROM storefront_products;

-- 4. Product Images (зависит от products)
INSERT INTO b2c_product_images SELECT * FROM storefront_product_images;

-- 5. Product Variants (зависит от products)
INSERT INTO b2c_product_variants SELECT * FROM storefront_product_variants;

-- 6. Product Variant Images (зависит от variants)
INSERT INTO b2c_product_variant_images SELECT * FROM storefront_product_variant_images;

-- 7. Product Attributes (зависит от products)
INSERT INTO b2c_product_attributes SELECT * FROM storefront_product_attributes;

-- 8. Orders (зависит от stores)
INSERT INTO b2c_orders SELECT * FROM storefront_orders;

-- 9. Order Items (зависит от orders и products)
INSERT INTO b2c_order_items SELECT * FROM storefront_order_items;

-- 10. Favorites (зависит от products)
INSERT INTO b2c_favorites SELECT * FROM storefront_favorites;

-- 11. Store Hours (зависит от stores)
INSERT INTO b2c_store_hours SELECT * FROM storefront_hours;

-- 12. Store Staff (зависит от stores)
INSERT INTO b2c_store_staff SELECT * FROM storefront_staff;

-- 13. Payment Methods (зависит от stores)
INSERT INTO b2c_payment_methods SELECT * FROM storefront_payment_methods;

-- 14. Delivery Options (зависит от stores)
INSERT INTO b2c_delivery_options SELECT * FROM storefront_delivery_options;

-- 15. Inventory Movements (зависит от variants)
INSERT INTO b2c_inventory_movements SELECT * FROM storefront_inventory_movements;

COMMIT;

-- ============================================================================
-- Проверка целостности данных
-- ============================================================================

DO $$
DECLARE
    c2c_count INTEGER;
    b2c_count INTEGER;
    old_c2c_count INTEGER;
    old_b2c_count INTEGER;
BEGIN
    -- Проверка C2C listings
    SELECT COUNT(*) INTO c2c_count FROM c2c_listings;
    SELECT COUNT(*) INTO old_c2c_count FROM marketplace_listings;

    IF c2c_count != old_c2c_count THEN
        RAISE EXCEPTION 'C2C listings count mismatch! Expected %, got %', old_c2c_count, c2c_count;
    END IF;

    -- Проверка B2C products
    SELECT COUNT(*) INTO b2c_count FROM b2c_products;
    SELECT COUNT(*) INTO old_b2c_count FROM storefront_products;

    IF b2c_count != old_b2c_count THEN
        RAISE EXCEPTION 'B2C products count mismatch! Expected %, got %', old_b2c_count, b2c_count;
    END IF;

    RAISE NOTICE '✅ Data migration validated: C2C=%, B2C=%', c2c_count, b2c_count;
END $$;

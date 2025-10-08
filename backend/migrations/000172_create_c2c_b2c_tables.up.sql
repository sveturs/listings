-- ============================================================================
-- МИГРАЦИЯ: Создание C2C и B2C таблиц
-- Дата: 2025-10-09
-- Описание: Переименование marketplace → c2c, storefronts → b2c
-- Автор: Migration Plan v1.0
-- ============================================================================

BEGIN;

-- ============================================================================
-- C2C TABLES (бывшие marketplace_*)
-- ============================================================================

-- 1. C2C Categories
CREATE TABLE c2c_categories (
    LIKE marketplace_categories INCLUDING ALL
);

-- 2. C2C Listings
CREATE TABLE c2c_listings (
    LIKE marketplace_listings INCLUDING ALL
);

-- 3. C2C Images
CREATE TABLE c2c_images (
    LIKE marketplace_images INCLUDING ALL
);

-- 4. C2C Chats
CREATE TABLE c2c_chats (
    LIKE marketplace_chats INCLUDING ALL
);

-- 5. C2C Messages
CREATE TABLE c2c_messages (
    LIKE marketplace_messages INCLUDING ALL
);

-- 6. C2C Favorites
CREATE TABLE c2c_favorites (
    LIKE marketplace_favorites INCLUDING ALL
);

-- 7. C2C Orders
CREATE TABLE c2c_orders (
    LIKE marketplace_orders INCLUDING ALL
);

-- 8. C2C Listing Variants
CREATE TABLE c2c_listing_variants (
    LIKE marketplace_listing_variants INCLUDING ALL
);

-- ============================================================================
-- B2C TABLES (бывшие storefront_*)
-- ============================================================================

-- 1. B2C Stores (основная таблица магазинов)
CREATE TABLE b2c_stores (
    LIKE storefronts INCLUDING ALL
);

-- 2. B2C Products
CREATE TABLE b2c_products (
    LIKE storefront_products INCLUDING ALL
);

-- 3. B2C Product Images
CREATE TABLE b2c_product_images (
    LIKE storefront_product_images INCLUDING ALL
);

-- 4. B2C Product Variants
CREATE TABLE b2c_product_variants (
    LIKE storefront_product_variants INCLUDING ALL
);

-- 5. B2C Product Attributes
CREATE TABLE b2c_product_attributes (
    LIKE storefront_product_attributes INCLUDING ALL
);

-- 6. B2C Orders
CREATE TABLE b2c_orders (
    LIKE storefront_orders INCLUDING ALL
);

-- 7. B2C Order Items
CREATE TABLE b2c_order_items (
    LIKE storefront_order_items INCLUDING ALL
);

-- 8. B2C Favorites
CREATE TABLE b2c_favorites (
    LIKE storefront_favorites INCLUDING ALL
);

-- 9. B2C Store Hours
CREATE TABLE b2c_store_hours (
    LIKE storefront_hours INCLUDING ALL
);

-- 10. B2C Store Staff
CREATE TABLE b2c_store_staff (
    LIKE storefront_staff INCLUDING ALL
);

-- 11. B2C Payment Methods
CREATE TABLE b2c_payment_methods (
    LIKE storefront_payment_methods INCLUDING ALL
);

-- 12. B2C Delivery Options
CREATE TABLE b2c_delivery_options (
    LIKE storefront_delivery_options INCLUDING ALL
);

-- 13. B2C Inventory Movements
CREATE TABLE b2c_inventory_movements (
    LIKE storefront_inventory_movements INCLUDING ALL
);

-- 14. User B2C Stores (связь пользователей с магазинами)
CREATE TABLE user_b2c_stores (
    LIKE user_storefronts INCLUDING ALL
);

-- 15. B2C Product Variant Images
CREATE TABLE b2c_product_variant_images (
    LIKE storefront_product_variant_images INCLUDING ALL
);

COMMIT;

-- ============================================================================
-- ВАЖНО: Индексы и constraints скопированы через INCLUDING ALL
-- Следующий шаг: миграция данных (отдельная миграция)
-- ============================================================================

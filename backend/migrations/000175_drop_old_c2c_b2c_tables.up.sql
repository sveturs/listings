-- ============================================================================
-- МИГРАЦИЯ: Удаление старых таблиц marketplace_* и storefront_*
-- Дата: 2025-10-11
-- ВНИМАНИЕ: Применять ТОЛЬКО после полного тестирования новых таблиц!
-- ВАЖНО: Эта миграция НЕ ПРИМЕНЯЕТСЯ автоматически!
-- ============================================================================

BEGIN;

-- Удаление C2C таблиц (marketplace_*)
DROP TABLE IF EXISTS marketplace_orders CASCADE;
DROP TABLE IF EXISTS marketplace_messages CASCADE;
DROP TABLE IF EXISTS marketplace_chats CASCADE;
DROP TABLE IF EXISTS marketplace_favorites CASCADE;
DROP TABLE IF EXISTS marketplace_listing_variants CASCADE;
DROP TABLE IF EXISTS marketplace_images CASCADE;
DROP TABLE IF EXISTS marketplace_listings CASCADE;
DROP TABLE IF EXISTS marketplace_categories CASCADE;

-- Удаление B2C таблиц (storefront_*)
DROP TABLE IF EXISTS storefront_cart_items CASCADE;
DROP TABLE IF EXISTS storefront_carts CASCADE;
DROP TABLE IF EXISTS storefront_rating_distribution CASCADE;
DROP TABLE IF EXISTS storefront_ratings CASCADE;
DROP TABLE IF EXISTS storefront_events CASCADE;
DROP TABLE IF EXISTS storefront_order_items CASCADE;
DROP TABLE IF EXISTS storefront_orders CASCADE;
DROP TABLE IF EXISTS storefront_favorites CASCADE;
DROP TABLE IF EXISTS storefront_product_variant_images CASCADE;
DROP TABLE IF EXISTS storefront_product_variants CASCADE;
DROP TABLE IF EXISTS storefront_product_attributes CASCADE;
DROP TABLE IF EXISTS storefront_product_images CASCADE;
DROP TABLE IF EXISTS storefront_products CASCADE;
DROP TABLE IF EXISTS storefront_inventory_movements CASCADE;
DROP TABLE IF EXISTS storefront_delivery_options CASCADE;
DROP TABLE IF EXISTS storefront_payment_methods CASCADE;
DROP TABLE IF EXISTS storefront_staff CASCADE;
DROP TABLE IF EXISTS storefront_hours CASCADE;
DROP TABLE IF EXISTS user_storefronts CASCADE;
DROP TABLE IF EXISTS storefronts CASCADE;

COMMIT;

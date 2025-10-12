-- ============================================================================
-- ОТКАТ: Удаление всех sequences для C2C/B2C таблиц
-- ============================================================================

BEGIN;

-- B2C tables
ALTER TABLE b2c_delivery_options ALTER COLUMN id DROP DEFAULT;
DROP SEQUENCE IF EXISTS b2c_delivery_options_id_seq CASCADE;

ALTER TABLE b2c_inventory_movements ALTER COLUMN id DROP DEFAULT;
DROP SEQUENCE IF EXISTS b2c_inventory_movements_id_seq CASCADE;

ALTER TABLE b2c_order_items ALTER COLUMN id DROP DEFAULT;
DROP SEQUENCE IF EXISTS b2c_order_items_id_seq CASCADE;

ALTER TABLE b2c_orders ALTER COLUMN id DROP DEFAULT;
DROP SEQUENCE IF EXISTS b2c_orders_id_seq CASCADE;

ALTER TABLE b2c_payment_methods ALTER COLUMN id DROP DEFAULT;
DROP SEQUENCE IF EXISTS b2c_payment_methods_id_seq CASCADE;

ALTER TABLE b2c_product_attributes ALTER COLUMN id DROP DEFAULT;
DROP SEQUENCE IF EXISTS b2c_product_attributes_id_seq CASCADE;

ALTER TABLE b2c_product_images ALTER COLUMN id DROP DEFAULT;
DROP SEQUENCE IF EXISTS b2c_product_images_id_seq CASCADE;

ALTER TABLE b2c_product_variant_images ALTER COLUMN id DROP DEFAULT;
DROP SEQUENCE IF EXISTS b2c_product_variant_images_id_seq CASCADE;

ALTER TABLE b2c_product_variants ALTER COLUMN id DROP DEFAULT;
DROP SEQUENCE IF EXISTS b2c_product_variants_id_seq CASCADE;

ALTER TABLE b2c_store_hours ALTER COLUMN id DROP DEFAULT;
DROP SEQUENCE IF EXISTS b2c_store_hours_id_seq CASCADE;

ALTER TABLE b2c_store_staff ALTER COLUMN id DROP DEFAULT;
DROP SEQUENCE IF EXISTS b2c_store_staff_id_seq CASCADE;

ALTER TABLE b2c_stores ALTER COLUMN id DROP DEFAULT;
DROP SEQUENCE IF EXISTS b2c_stores_id_seq CASCADE;

-- C2C tables
ALTER TABLE c2c_categories ALTER COLUMN id DROP DEFAULT;
DROP SEQUENCE IF EXISTS c2c_categories_id_seq CASCADE;

ALTER TABLE c2c_chats ALTER COLUMN id DROP DEFAULT;
DROP SEQUENCE IF EXISTS c2c_chats_id_seq CASCADE;

ALTER TABLE c2c_listing_variants ALTER COLUMN id DROP DEFAULT;
DROP SEQUENCE IF EXISTS c2c_listing_variants_id_seq CASCADE;

ALTER TABLE c2c_messages ALTER COLUMN id DROP DEFAULT;
DROP SEQUENCE IF EXISTS c2c_messages_id_seq CASCADE;

ALTER TABLE c2c_orders ALTER COLUMN id DROP DEFAULT;
DROP SEQUENCE IF EXISTS c2c_orders_id_seq CASCADE;

COMMIT;

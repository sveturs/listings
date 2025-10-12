-- ============================================================================
-- МИГРАЦИЯ: Массовое исправление sequences для всех C2C/B2C таблиц
-- Дата: 2025-10-11
-- Описание: Добавление SERIAL для автоинкремента ID во всех таблицах
-- Проблема: При создании через LIKE не копируются DEFAULT и sequences
-- Количество исправлений: 17 таблиц
-- ============================================================================

BEGIN;

-- ============================================================================
-- B2C TABLES
-- ============================================================================

-- b2c_delivery_options
CREATE SEQUENCE IF NOT EXISTS b2c_delivery_options_id_seq;
SELECT setval('b2c_delivery_options_id_seq', COALESCE((SELECT MAX(id) FROM b2c_delivery_options), 0) + 1, false);
ALTER TABLE b2c_delivery_options ALTER COLUMN id SET DEFAULT nextval('b2c_delivery_options_id_seq');
ALTER SEQUENCE b2c_delivery_options_id_seq OWNED BY b2c_delivery_options.id;

-- b2c_inventory_movements
CREATE SEQUENCE IF NOT EXISTS b2c_inventory_movements_id_seq;
SELECT setval('b2c_inventory_movements_id_seq', COALESCE((SELECT MAX(id) FROM b2c_inventory_movements), 0) + 1, false);
ALTER TABLE b2c_inventory_movements ALTER COLUMN id SET DEFAULT nextval('b2c_inventory_movements_id_seq');
ALTER SEQUENCE b2c_inventory_movements_id_seq OWNED BY b2c_inventory_movements.id;

-- b2c_order_items
CREATE SEQUENCE IF NOT EXISTS b2c_order_items_id_seq;
SELECT setval('b2c_order_items_id_seq', COALESCE((SELECT MAX(id) FROM b2c_order_items), 0) + 1, false);
ALTER TABLE b2c_order_items ALTER COLUMN id SET DEFAULT nextval('b2c_order_items_id_seq');
ALTER SEQUENCE b2c_order_items_id_seq OWNED BY b2c_order_items.id;

-- b2c_orders
CREATE SEQUENCE IF NOT EXISTS b2c_orders_id_seq;
SELECT setval('b2c_orders_id_seq', COALESCE((SELECT MAX(id) FROM b2c_orders), 0) + 1, false);
ALTER TABLE b2c_orders ALTER COLUMN id SET DEFAULT nextval('b2c_orders_id_seq');
ALTER SEQUENCE b2c_orders_id_seq OWNED BY b2c_orders.id;

-- b2c_payment_methods
CREATE SEQUENCE IF NOT EXISTS b2c_payment_methods_id_seq;
SELECT setval('b2c_payment_methods_id_seq', COALESCE((SELECT MAX(id) FROM b2c_payment_methods), 0) + 1, false);
ALTER TABLE b2c_payment_methods ALTER COLUMN id SET DEFAULT nextval('b2c_payment_methods_id_seq');
ALTER SEQUENCE b2c_payment_methods_id_seq OWNED BY b2c_payment_methods.id;

-- b2c_product_attributes
CREATE SEQUENCE IF NOT EXISTS b2c_product_attributes_id_seq;
SELECT setval('b2c_product_attributes_id_seq', COALESCE((SELECT MAX(id) FROM b2c_product_attributes), 0) + 1, false);
ALTER TABLE b2c_product_attributes ALTER COLUMN id SET DEFAULT nextval('b2c_product_attributes_id_seq');
ALTER SEQUENCE b2c_product_attributes_id_seq OWNED BY b2c_product_attributes.id;

-- b2c_product_images
CREATE SEQUENCE IF NOT EXISTS b2c_product_images_id_seq;
SELECT setval('b2c_product_images_id_seq', COALESCE((SELECT MAX(id) FROM b2c_product_images), 0) + 1, false);
ALTER TABLE b2c_product_images ALTER COLUMN id SET DEFAULT nextval('b2c_product_images_id_seq');
ALTER SEQUENCE b2c_product_images_id_seq OWNED BY b2c_product_images.id;

-- b2c_product_variant_images
CREATE SEQUENCE IF NOT EXISTS b2c_product_variant_images_id_seq;
SELECT setval('b2c_product_variant_images_id_seq', COALESCE((SELECT MAX(id) FROM b2c_product_variant_images), 0) + 1, false);
ALTER TABLE b2c_product_variant_images ALTER COLUMN id SET DEFAULT nextval('b2c_product_variant_images_id_seq');
ALTER SEQUENCE b2c_product_variant_images_id_seq OWNED BY b2c_product_variant_images.id;

-- b2c_product_variants
CREATE SEQUENCE IF NOT EXISTS b2c_product_variants_id_seq;
SELECT setval('b2c_product_variants_id_seq', COALESCE((SELECT MAX(id) FROM b2c_product_variants), 0) + 1, false);
ALTER TABLE b2c_product_variants ALTER COLUMN id SET DEFAULT nextval('b2c_product_variants_id_seq');
ALTER SEQUENCE b2c_product_variants_id_seq OWNED BY b2c_product_variants.id;

-- b2c_store_hours
CREATE SEQUENCE IF NOT EXISTS b2c_store_hours_id_seq;
SELECT setval('b2c_store_hours_id_seq', COALESCE((SELECT MAX(id) FROM b2c_store_hours), 0) + 1, false);
ALTER TABLE b2c_store_hours ALTER COLUMN id SET DEFAULT nextval('b2c_store_hours_id_seq');
ALTER SEQUENCE b2c_store_hours_id_seq OWNED BY b2c_store_hours.id;

-- b2c_store_staff
CREATE SEQUENCE IF NOT EXISTS b2c_store_staff_id_seq;
SELECT setval('b2c_store_staff_id_seq', COALESCE((SELECT MAX(id) FROM b2c_store_staff), 0) + 1, false);
ALTER TABLE b2c_store_staff ALTER COLUMN id SET DEFAULT nextval('b2c_store_staff_id_seq');
ALTER SEQUENCE b2c_store_staff_id_seq OWNED BY b2c_store_staff.id;

-- b2c_stores
CREATE SEQUENCE IF NOT EXISTS b2c_stores_id_seq;
SELECT setval('b2c_stores_id_seq', COALESCE((SELECT MAX(id) FROM b2c_stores), 0) + 1, false);
ALTER TABLE b2c_stores ALTER COLUMN id SET DEFAULT nextval('b2c_stores_id_seq');
ALTER SEQUENCE b2c_stores_id_seq OWNED BY b2c_stores.id;

-- ============================================================================
-- C2C TABLES
-- ============================================================================

-- c2c_categories
CREATE SEQUENCE IF NOT EXISTS c2c_categories_id_seq;
SELECT setval('c2c_categories_id_seq', COALESCE((SELECT MAX(id) FROM c2c_categories), 0) + 1, false);
ALTER TABLE c2c_categories ALTER COLUMN id SET DEFAULT nextval('c2c_categories_id_seq');
ALTER SEQUENCE c2c_categories_id_seq OWNED BY c2c_categories.id;

-- c2c_chats
CREATE SEQUENCE IF NOT EXISTS c2c_chats_id_seq;
SELECT setval('c2c_chats_id_seq', COALESCE((SELECT MAX(id) FROM c2c_chats), 0) + 1, false);
ALTER TABLE c2c_chats ALTER COLUMN id SET DEFAULT nextval('c2c_chats_id_seq');
ALTER SEQUENCE c2c_chats_id_seq OWNED BY c2c_chats.id;

-- c2c_listing_variants
CREATE SEQUENCE IF NOT EXISTS c2c_listing_variants_id_seq;
SELECT setval('c2c_listing_variants_id_seq', COALESCE((SELECT MAX(id) FROM c2c_listing_variants), 0) + 1, false);
ALTER TABLE c2c_listing_variants ALTER COLUMN id SET DEFAULT nextval('c2c_listing_variants_id_seq');
ALTER SEQUENCE c2c_listing_variants_id_seq OWNED BY c2c_listing_variants.id;

-- c2c_messages
CREATE SEQUENCE IF NOT EXISTS c2c_messages_id_seq;
SELECT setval('c2c_messages_id_seq', COALESCE((SELECT MAX(id) FROM c2c_messages), 0) + 1, false);
ALTER TABLE c2c_messages ALTER COLUMN id SET DEFAULT nextval('c2c_messages_id_seq');
ALTER SEQUENCE c2c_messages_id_seq OWNED BY c2c_messages.id;

-- c2c_orders
CREATE SEQUENCE IF NOT EXISTS c2c_orders_id_seq;
SELECT setval('c2c_orders_id_seq', COALESCE((SELECT MAX(id) FROM c2c_orders), 0) + 1, false);
ALTER TABLE c2c_orders ALTER COLUMN id SET DEFAULT nextval('c2c_orders_id_seq');
ALTER SEQUENCE c2c_orders_id_seq OWNED BY c2c_orders.id;

COMMIT;

-- ============================================================================
-- РЕЗУЛЬТАТ: 17 таблиц исправлено, автоинкремент ID работает везде
-- ============================================================================

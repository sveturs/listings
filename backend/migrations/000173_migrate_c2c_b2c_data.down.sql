BEGIN;

-- Очистить все данные из новых таблиц (обратная миграция)
TRUNCATE TABLE c2c_messages CASCADE;
TRUNCATE TABLE c2c_chats CASCADE;
TRUNCATE TABLE c2c_favorites CASCADE;
TRUNCATE TABLE c2c_orders CASCADE;
TRUNCATE TABLE c2c_listing_variants CASCADE;
TRUNCATE TABLE c2c_images CASCADE;
TRUNCATE TABLE c2c_listings CASCADE;
TRUNCATE TABLE c2c_categories CASCADE;

TRUNCATE TABLE b2c_inventory_movements CASCADE;
TRUNCATE TABLE b2c_delivery_options CASCADE;
TRUNCATE TABLE b2c_payment_methods CASCADE;
TRUNCATE TABLE b2c_store_staff CASCADE;
TRUNCATE TABLE b2c_store_hours CASCADE;
TRUNCATE TABLE b2c_favorites CASCADE;
TRUNCATE TABLE b2c_order_items CASCADE;
TRUNCATE TABLE b2c_orders CASCADE;
TRUNCATE TABLE b2c_product_attributes CASCADE;
TRUNCATE TABLE b2c_product_variant_images CASCADE;
TRUNCATE TABLE b2c_product_variants CASCADE;
TRUNCATE TABLE b2c_product_images CASCADE;
TRUNCATE TABLE b2c_products CASCADE;
TRUNCATE TABLE user_b2c_stores CASCADE;
TRUNCATE TABLE b2c_stores CASCADE;

COMMIT;

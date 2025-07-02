-- Drop all storefront related tables in reverse order

-- Remove columns from marketplace_listings
ALTER TABLE marketplace_listings 
DROP COLUMN IF EXISTS storefront_id,
DROP COLUMN IF EXISTS is_storefront_featured;

-- Drop import related tables
DROP TABLE IF EXISTS category_mappings CASCADE;
DROP TABLE IF EXISTS import_errors CASCADE;
DROP TABLE IF EXISTS import_jobs CASCADE;

-- Drop storefront events
DROP TABLE IF EXISTS storefront_events CASCADE;

-- Drop inventory and product related tables
DROP TABLE IF EXISTS storefront_inventory_movements CASCADE;
DROP TABLE IF EXISTS storefront_product_variants CASCADE;
DROP TABLE IF EXISTS storefront_product_images CASCADE;
DROP TABLE IF EXISTS storefront_products CASCADE;

-- Drop storefront related tables
DROP TABLE IF EXISTS storefront_analytics CASCADE;
DROP TABLE IF EXISTS storefront_delivery_options CASCADE;
DROP TABLE IF EXISTS storefront_payment_methods CASCADE;
DROP TABLE IF EXISTS storefront_hours CASCADE;
DROP TABLE IF EXISTS storefront_staff CASCADE;
DROP TABLE IF EXISTS storefronts CASCADE;
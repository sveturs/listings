-- Drop triggers
DROP TRIGGER IF EXISTS update_stock_status_trigger ON storefront_products;
DROP TRIGGER IF EXISTS update_storefront_products_updated_at ON storefront_products;
DROP TRIGGER IF EXISTS update_storefront_product_variants_updated_at ON storefront_product_variants;

-- Drop functions
DROP FUNCTION IF EXISTS update_product_stock_status();

-- Drop tables
DROP TABLE IF EXISTS storefront_inventory_movements;
DROP TABLE IF EXISTS storefront_product_variants;
DROP TABLE IF EXISTS storefront_product_images;
DROP TABLE IF EXISTS storefront_products;
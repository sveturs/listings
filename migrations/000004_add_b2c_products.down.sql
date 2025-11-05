-- Phase 9.5.5: Rollback B2C Products and Inventory Tracking

-- Drop tables in reverse order (child tables first)
DROP TABLE IF EXISTS b2c_inventory_movements;
DROP TABLE IF EXISTS b2c_product_variants;
DROP TABLE IF EXISTS b2c_products;

-- Drop functions
DROP FUNCTION IF EXISTS update_b2c_products_updated_at();
DROP FUNCTION IF EXISTS update_b2c_product_variants_updated_at();

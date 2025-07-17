-- Migration down: Drop product variants system

-- Drop triggers
DROP TRIGGER IF EXISTS trigger_update_storefront_product_variants_updated_at ON storefront_product_variants;
DROP TRIGGER IF EXISTS trigger_update_product_variant_attribute_values_updated_at ON product_variant_attribute_values;
DROP TRIGGER IF EXISTS trigger_update_product_variant_attributes_updated_at ON product_variant_attributes;

-- Drop function
DROP FUNCTION IF EXISTS update_product_variants_updated_at();

-- Drop indexes (they will be dropped automatically with tables, but explicit for clarity)
DROP INDEX IF EXISTS idx_storefront_product_variant_images_is_main;
DROP INDEX IF EXISTS idx_storefront_product_variant_images_variant_id;
DROP INDEX IF EXISTS idx_storefront_product_variants_attributes_gin;
DROP INDEX IF EXISTS idx_storefront_product_variants_default_unique;
DROP INDEX IF EXISTS idx_storefront_product_variants_is_default;
DROP INDEX IF EXISTS idx_storefront_product_variants_is_active;
DROP INDEX IF EXISTS idx_storefront_product_variants_stock_status;
DROP INDEX IF EXISTS idx_storefront_product_variants_barcode;
DROP INDEX IF EXISTS idx_storefront_product_variants_sku;
DROP INDEX IF EXISTS idx_storefront_product_variants_parent_id;
DROP INDEX IF EXISTS idx_product_variant_attribute_values_value;
DROP INDEX IF EXISTS idx_product_variant_attribute_values_attribute_id;
DROP INDEX IF EXISTS idx_product_variant_attributes_name;

-- Drop tables in reverse order (due to foreign key constraints)
DROP TABLE IF EXISTS storefront_product_variant_images;
DROP TABLE IF EXISTS storefront_product_variants;
DROP TABLE IF EXISTS product_variant_attribute_values;
DROP TABLE IF EXISTS product_variant_attributes;

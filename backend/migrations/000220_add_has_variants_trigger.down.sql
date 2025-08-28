-- Drop trigger
DROP TRIGGER IF EXISTS update_has_variants_trigger ON storefront_product_variants;

-- Drop function
DROP FUNCTION IF EXISTS update_product_has_variants();
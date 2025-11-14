-- Migration 000029: Create listing_variants VIEW for compatibility with cart_repository queries
--
-- Purpose: Provide backward compatibility layer between cart queries expecting 'listing_variants'
-- and the actual table 'b2c_product_variants'
--
-- Key mappings:
-- - stock_quantity -> stock (alias)
-- - variant_attributes -> attributes (alias)
-- - NULL -> image_url (column doesn't exist in b2c_product_variants)

DROP VIEW IF EXISTS listing_variants;

CREATE OR REPLACE VIEW listing_variants AS
SELECT
    id,
    product_id,
    sku,
    barcode,
    price,
    compare_at_price,
    cost_price,
    stock_quantity AS stock,          -- Alias for cart queries compatibility
    stock_status,
    low_stock_threshold,
    variant_attributes AS attributes,  -- Alias for cart queries compatibility
    NULL::text AS image_url,          -- Missing column, added as NULL
    weight,
    dimensions,
    is_active,
    is_default,
    view_count,
    sold_count,
    created_at,
    updated_at
FROM b2c_product_variants;

-- Add comment
COMMENT ON VIEW listing_variants IS 'View alias for b2c_product_variants table with compatibility mappings for cart queries';

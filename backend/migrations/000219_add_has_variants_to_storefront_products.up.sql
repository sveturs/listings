-- Add has_variants column to storefront_products table
ALTER TABLE storefront_products 
ADD COLUMN has_variants BOOLEAN DEFAULT false;

-- Update has_variants for existing products based on whether they have variants
UPDATE storefront_products p
SET has_variants = EXISTS (
    SELECT 1 
    FROM storefront_product_variants v 
    WHERE v.product_id = p.id
);

-- Add index for faster filtering
CREATE INDEX idx_storefront_products_has_variants ON storefront_products(has_variants);
-- Create function to update has_variants flag
CREATE OR REPLACE FUNCTION update_product_has_variants()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE storefront_products 
        SET has_variants = true 
        WHERE id = NEW.product_id;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE storefront_products 
        SET has_variants = EXISTS (
            SELECT 1 
            FROM storefront_product_variants 
            WHERE product_id = OLD.product_id
        )
        WHERE id = OLD.product_id;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Create trigger on storefront_product_variants table
CREATE TRIGGER update_has_variants_trigger
AFTER INSERT OR DELETE ON storefront_product_variants
FOR EACH ROW EXECUTE FUNCTION update_product_has_variants();
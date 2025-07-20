-- Drop table: product_variant_attributes
DROP SEQUENCE IF EXISTS public.product_variant_attributes_id_seq;
DROP TABLE IF EXISTS public.product_variant_attributes;
DROP INDEX IF EXISTS public.idx_product_variant_attributes_name;
DROP TRIGGER IF EXISTS trigger_update_product_variant_attributes_updated_at ON public.product_variant_attributes;
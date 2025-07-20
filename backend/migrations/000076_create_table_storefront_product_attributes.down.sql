-- Drop table: storefront_product_attributes
DROP SEQUENCE IF EXISTS public.storefront_product_attributes_id_seq;
DROP TABLE IF EXISTS public.storefront_product_attributes;
DROP INDEX IF EXISTS public.idx_storefront_product_attributes_attribute_id;
DROP INDEX IF EXISTS public.idx_storefront_product_attributes_enabled;
DROP INDEX IF EXISTS public.idx_storefront_product_attributes_product_id;
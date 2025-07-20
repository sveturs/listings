-- Drop table: storefront_product_variant_images
DROP SEQUENCE IF EXISTS public.storefront_product_variant_images_id_seq;
DROP TABLE IF EXISTS public.storefront_product_variant_images;
DROP INDEX IF EXISTS public.idx_storefront_product_variant_images_is_main;
DROP INDEX IF EXISTS public.idx_storefront_product_variant_images_variant_id;
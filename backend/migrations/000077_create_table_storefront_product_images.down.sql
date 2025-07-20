-- Drop table: storefront_product_images
DROP SEQUENCE IF EXISTS public.storefront_product_images_id_seq;
DROP TABLE IF EXISTS public.storefront_product_images;
DROP INDEX IF EXISTS public.idx_storefront_product_images_display_order;
DROP INDEX IF EXISTS public.idx_storefront_product_images_product_id;
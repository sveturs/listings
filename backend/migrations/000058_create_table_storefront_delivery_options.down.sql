-- Drop table: storefront_delivery_options
DROP SEQUENCE IF EXISTS public.storefront_delivery_options_id_seq;
DROP TABLE IF EXISTS public.storefront_delivery_options;
DROP INDEX IF EXISTS public.idx_delivery_is_active;
DROP INDEX IF EXISTS public.idx_delivery_storefront_id;
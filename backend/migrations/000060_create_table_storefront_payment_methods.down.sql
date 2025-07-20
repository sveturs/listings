-- Drop table: storefront_payment_methods
DROP SEQUENCE IF EXISTS public.storefront_payment_methods_id_seq;
DROP TABLE IF EXISTS public.storefront_payment_methods;
DROP INDEX IF EXISTS public.idx_payment_is_enabled;
DROP INDEX IF EXISTS public.idx_payment_storefront_id;
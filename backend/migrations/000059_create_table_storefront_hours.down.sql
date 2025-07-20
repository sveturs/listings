-- Drop table: storefront_hours
DROP SEQUENCE IF EXISTS public.storefront_hours_id_seq;
DROP TABLE IF EXISTS public.storefront_hours;
DROP INDEX IF EXISTS public.idx_hours_storefront_id;
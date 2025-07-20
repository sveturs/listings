-- Drop table: storefront_staff
DROP SEQUENCE IF EXISTS public.storefront_staff_id_seq;
DROP TABLE IF EXISTS public.storefront_staff;
DROP INDEX IF EXISTS public.idx_staff_storefront_id;
DROP INDEX IF EXISTS public.idx_staff_user_id;
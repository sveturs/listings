-- Drop table: user_storefronts
DROP SEQUENCE IF EXISTS public.user_storefronts_id_seq;
DROP TABLE IF EXISTS public.user_storefronts;
DROP INDEX IF EXISTS public.idx_user_storefronts_status;
DROP INDEX IF EXISTS public.idx_user_storefronts_user;
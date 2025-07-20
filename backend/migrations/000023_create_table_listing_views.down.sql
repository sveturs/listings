-- Drop table: listing_views
DROP SEQUENCE IF EXISTS public.listing_views_id_seq;
DROP TABLE IF EXISTS public.listing_views;
DROP INDEX IF EXISTS public.idx_listing_views_listing_ip;
DROP INDEX IF EXISTS public.idx_listing_views_listing_user;
DROP INDEX IF EXISTS public.idx_listing_views_time;
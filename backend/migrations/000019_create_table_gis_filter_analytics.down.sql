-- Drop table: gis_filter_analytics
DROP SEQUENCE IF EXISTS public.gis_filter_analytics_id_seq;
DROP TABLE IF EXISTS public.gis_filter_analytics;
DROP INDEX IF EXISTS public.idx_filter_analytics_created;
DROP INDEX IF EXISTS public.idx_filter_analytics_session;
DROP INDEX IF EXISTS public.idx_filter_analytics_type;
DROP INDEX IF EXISTS public.idx_filter_analytics_user;
-- Drop table: unified_geo
DROP SEQUENCE IF EXISTS public.unified_geo_id_seq;
DROP TABLE IF EXISTS public.unified_geo;
DROP INDEX IF EXISTS public.idx_unified_geo_geohash;
DROP INDEX IF EXISTS public.idx_unified_geo_location;
DROP INDEX IF EXISTS public.idx_unified_geo_source_type;
-- Drop table: gis_isochrone_cache
DROP SEQUENCE IF EXISTS public.gis_isochrone_cache_id_seq;
DROP TABLE IF EXISTS public.gis_isochrone_cache;
DROP INDEX IF EXISTS public.idx_isochrone_cache_center;
DROP INDEX IF EXISTS public.idx_isochrone_cache_expires;
DROP INDEX IF EXISTS public.idx_isochrone_cache_lookup;
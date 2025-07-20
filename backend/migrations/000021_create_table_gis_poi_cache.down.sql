-- Drop table: gis_poi_cache
DROP SEQUENCE IF EXISTS public.gis_poi_cache_id_seq;
DROP TABLE IF EXISTS public.gis_poi_cache;
DROP INDEX IF EXISTS public.idx_poi_cache_expires;
DROP INDEX IF EXISTS public.idx_poi_cache_location;
DROP INDEX IF EXISTS public.idx_poi_cache_type;
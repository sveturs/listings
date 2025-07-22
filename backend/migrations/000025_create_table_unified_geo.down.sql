-- Drop table: unified_geo
DROP SEQUENCE IF EXISTS public.unified_geo_id_seq;
DROP TABLE IF EXISTS public.unified_geo;
DROP INDEX IF EXISTS public.idx_unified_geo_composite;
DROP INDEX IF EXISTS public.idx_unified_geo_geohash;
DROP INDEX IF EXISTS public.idx_unified_geo_location;
DROP INDEX IF EXISTS public.idx_unified_geo_location_bounds;
DROP INDEX IF EXISTS public.idx_unified_geo_marketplace_active;
DROP INDEX IF EXISTS public.idx_unified_geo_source_id;
DROP INDEX IF EXISTS public.idx_unified_geo_source_type;
DROP INDEX IF EXISTS public.idx_unified_geo_storefront_active;
DROP TRIGGER IF EXISTS trigger_unified_geo_cache_refresh ON public.unified_geo;
DROP TRIGGER IF EXISTS trigger_update_unified_geo_updated_at ON public.unified_geo;
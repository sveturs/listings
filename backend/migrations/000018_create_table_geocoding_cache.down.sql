-- Drop table: geocoding_cache
DROP SEQUENCE IF EXISTS public.geocoding_cache_id_seq;
DROP TABLE IF EXISTS public.geocoding_cache;
DROP INDEX IF EXISTS public.idx_geocoding_cache_address_components;
DROP INDEX IF EXISTS public.idx_geocoding_cache_cache_hits;
DROP INDEX IF EXISTS public.idx_geocoding_cache_confidence;
DROP INDEX IF EXISTS public.idx_geocoding_cache_country_lang;
DROP INDEX IF EXISTS public.idx_geocoding_cache_expires_at;
DROP INDEX IF EXISTS public.idx_geocoding_cache_input_address;
DROP INDEX IF EXISTS public.idx_geocoding_cache_location;
DROP INDEX IF EXISTS public.idx_geocoding_cache_normalized;
DROP INDEX IF EXISTS public.idx_geocoding_cache_provider;
DROP TRIGGER IF EXISTS trigger_cleanup_geocoding_cache ON public.geocoding_cache;
DROP TRIGGER IF EXISTS trigger_geocoding_cache_updated_at ON public.geocoding_cache;
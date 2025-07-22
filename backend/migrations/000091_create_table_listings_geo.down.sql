-- Drop table: listings_geo
DROP SEQUENCE IF EXISTS public.listings_geo_id_seq;
DROP TABLE IF EXISTS public.listings_geo;
DROP INDEX IF EXISTS public.idx_listings_geo_address_components;
DROP INDEX IF EXISTS public.idx_listings_geo_blurred_location;
DROP INDEX IF EXISTS public.idx_listings_geo_confidence;
DROP INDEX IF EXISTS public.idx_listings_geo_district;
DROP INDEX IF EXISTS public.idx_listings_geo_geohash;
DROP INDEX IF EXISTS public.idx_listings_geo_geohash_precise;
DROP INDEX IF EXISTS public.idx_listings_geo_input_method;
DROP INDEX IF EXISTS public.idx_listings_geo_is_precise;
DROP INDEX IF EXISTS public.idx_listings_geo_location;
DROP INDEX IF EXISTS public.idx_listings_geo_municipality;
DROP INDEX IF EXISTS public.idx_listings_geo_privacy;
DROP INDEX IF EXISTS public.idx_listings_geo_verified;
DROP TRIGGER IF EXISTS trigger_listings_geo_updated_at ON public.listings_geo;
DROP TRIGGER IF EXISTS trigger_update_listings_geo_updated_at ON public.listings_geo;
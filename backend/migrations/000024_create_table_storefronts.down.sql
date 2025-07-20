-- Drop table: storefronts
DROP SEQUENCE IF EXISTS public.storefronts_id_seq;
DROP TABLE IF EXISTS public.storefronts;
DROP INDEX IF EXISTS public.idx_storefronts_city;
DROP INDEX IF EXISTS public.idx_storefronts_coordinates;
DROP INDEX IF EXISTS public.idx_storefronts_geo_strategy;
DROP INDEX IF EXISTS public.idx_storefronts_is_active;
DROP INDEX IF EXISTS public.idx_storefronts_rating;
DROP INDEX IF EXISTS public.idx_storefronts_slug;
DROP INDEX IF EXISTS public.idx_storefronts_user_id;
DROP TRIGGER IF EXISTS trigger_update_storefront_products_geo ON public.storefronts;
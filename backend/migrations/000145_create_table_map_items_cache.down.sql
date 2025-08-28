-- Migration rollback for table: map_items_cache

DROP TABLE IF EXISTS public.map_items_cache CASCADE;
DROP SEQUENCE IF EXISTS public.map_items_cache_id_seq CASCADE;
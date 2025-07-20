-- Drop table: import_sources
DROP SEQUENCE IF EXISTS public.import_sources_id_seq;
DROP TABLE IF EXISTS public.import_sources;
DROP INDEX IF EXISTS public.idx_import_sources_storefront;
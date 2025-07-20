-- Drop table: imported_categories
DROP SEQUENCE IF EXISTS public.imported_categories_id_seq;
DROP TABLE IF EXISTS public.imported_categories;
DROP INDEX IF EXISTS public.idx_imported_categories_source_id;
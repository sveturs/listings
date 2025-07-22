-- Drop table: marketplace_categories
DROP SEQUENCE IF EXISTS public.marketplace_categories_id_seq;
DROP TABLE IF EXISTS public.marketplace_categories;
DROP INDEX IF EXISTS public.idx_marketplace_categories_external_id;
DROP INDEX IF EXISTS public.idx_marketplace_categories_parent;
DROP INDEX IF EXISTS public.idx_marketplace_categories_slug;
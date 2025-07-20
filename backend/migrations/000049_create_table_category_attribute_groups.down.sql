-- Drop table: category_attribute_groups
DROP SEQUENCE IF EXISTS public.category_attribute_groups_id_seq;
DROP TABLE IF EXISTS public.category_attribute_groups;
DROP INDEX IF EXISTS public.idx_category_attribute_groups_category;
DROP INDEX IF EXISTS public.idx_category_attribute_groups_component;
DROP INDEX IF EXISTS public.idx_category_attribute_groups_group;
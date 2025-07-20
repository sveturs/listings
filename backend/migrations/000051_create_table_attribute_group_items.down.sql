-- Drop table: attribute_group_items
DROP SEQUENCE IF EXISTS public.attribute_group_items_id_seq;
DROP TABLE IF EXISTS public.attribute_group_items;
DROP INDEX IF EXISTS public.idx_attribute_group_items_attribute;
DROP INDEX IF EXISTS public.idx_attribute_group_items_group;
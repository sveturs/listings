-- Drop table: attribute_groups
DROP SEQUENCE IF EXISTS public.attribute_groups_id_seq;
DROP TABLE IF EXISTS public.attribute_groups;
DROP INDEX IF EXISTS public.idx_attribute_groups_active;
DROP INDEX IF EXISTS public.idx_attribute_groups_name;
DROP TRIGGER IF EXISTS update_attribute_groups_updated_at ON public.attribute_groups;
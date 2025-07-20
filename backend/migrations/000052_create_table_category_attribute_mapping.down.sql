-- Drop table: category_attribute_mapping
DROP TABLE IF EXISTS public.category_attribute_mapping;
DROP INDEX IF EXISTS public.idx_category_attr_mapping;
DROP INDEX IF EXISTS public.idx_category_attribute_map_attr_id;
DROP INDEX IF EXISTS public.idx_category_attribute_map_cat_id;
DROP INDEX IF EXISTS public.idx_category_attribute_mapping_custom_component;
DROP TRIGGER IF EXISTS tr_update_category_attribute_sort_order ON public.category_attribute_mapping;
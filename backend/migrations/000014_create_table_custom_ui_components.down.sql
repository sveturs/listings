-- Drop table: custom_ui_components
DROP SEQUENCE IF EXISTS public.custom_ui_components_id_seq;
DROP TABLE IF EXISTS public.custom_ui_components;
DROP INDEX IF EXISTS public.idx_custom_ui_components_active;
DROP INDEX IF EXISTS public.idx_custom_ui_components_name;
DROP INDEX IF EXISTS public.idx_custom_ui_components_type;
-- Drop table: custom_ui_component_usage
DROP SEQUENCE IF EXISTS public.custom_ui_component_usage_id_seq;
DROP TABLE IF EXISTS public.custom_ui_component_usage;
DROP INDEX IF EXISTS public.idx_ui_comp_usage_active;
DROP INDEX IF EXISTS public.idx_ui_comp_usage_category;
DROP INDEX IF EXISTS public.idx_ui_comp_usage_component;
DROP INDEX IF EXISTS public.idx_ui_comp_usage_context;
DROP TRIGGER IF EXISTS update_custom_ui_component_usage_updated_at ON public.custom_ui_component_usage;
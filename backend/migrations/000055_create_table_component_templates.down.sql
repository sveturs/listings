-- Drop table: component_templates
DROP SEQUENCE IF EXISTS public.component_templates_id_seq;
DROP TABLE IF EXISTS public.component_templates;
DROP INDEX IF EXISTS public.idx_component_templates_cat;
DROP INDEX IF EXISTS public.idx_component_templates_comp;
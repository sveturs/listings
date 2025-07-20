-- Migration for table: component_templates

CREATE SEQUENCE public.component_templates_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.component_templates (
    id integer NOT NULL,
    component_id integer NOT NULL,
    name character varying(255) NOT NULL,
    description text,
    template_config jsonb DEFAULT '{}'::jsonb,
    preview_image text,
    category_id integer,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    created_by integer
);

ALTER SEQUENCE public.component_templates_id_seq OWNED BY public.component_templates.id;

CREATE INDEX idx_component_templates_cat ON public.component_templates USING btree (category_id);

CREATE INDEX idx_component_templates_comp ON public.component_templates USING btree (component_id);

ALTER TABLE ONLY public.component_templates
    ADD CONSTRAINT component_templates_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.component_templates
    ADD CONSTRAINT component_templates_component_id_fkey FOREIGN KEY (component_id) REFERENCES public.custom_ui_components(id) ON DELETE CASCADE;
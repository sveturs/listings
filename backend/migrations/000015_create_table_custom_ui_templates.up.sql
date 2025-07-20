-- Migration for table: custom_ui_templates

CREATE SEQUENCE public.custom_ui_templates_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.custom_ui_templates (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    display_name character varying(255) NOT NULL,
    description text,
    template_code text NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    variables jsonb DEFAULT '{}'::jsonb,
    is_shared boolean DEFAULT false,
    created_by integer,
    updated_by integer
);

ALTER SEQUENCE public.custom_ui_templates_id_seq OWNED BY public.custom_ui_templates.id;

CREATE INDEX idx_custom_ui_templates_name ON public.custom_ui_templates USING btree (name);

ALTER TABLE ONLY public.custom_ui_templates
    ADD CONSTRAINT custom_ui_templates_name_key UNIQUE (name);

ALTER TABLE ONLY public.custom_ui_templates
    ADD CONSTRAINT custom_ui_templates_pkey PRIMARY KEY (id);
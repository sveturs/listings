-- Migration for table: custom_ui_components

CREATE SEQUENCE public.custom_ui_components_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.custom_ui_components (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    description text,
    component_type character varying(50) NOT NULL,
    is_active boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    created_by integer,
    updated_by integer,
    template_code text DEFAULT ''::text NOT NULL,
    styles text DEFAULT ''::text,
    props_schema jsonb DEFAULT '{}'::jsonb,
    CONSTRAINT custom_ui_components_component_type_check CHECK (((component_type)::text = ANY (ARRAY[('category'::character varying)::text, ('attribute'::character varying)::text, ('filter'::character varying)::text])))
);

ALTER SEQUENCE public.custom_ui_components_id_seq OWNED BY public.custom_ui_components.id;

CREATE INDEX idx_custom_ui_components_active ON public.custom_ui_components USING btree (is_active);

CREATE INDEX idx_custom_ui_components_name ON public.custom_ui_components USING btree (name);

CREATE INDEX idx_custom_ui_components_type ON public.custom_ui_components USING btree (component_type);

ALTER TABLE ONLY public.custom_ui_components
    ADD CONSTRAINT custom_ui_components_name_key UNIQUE (name);

ALTER TABLE ONLY public.custom_ui_components
    ADD CONSTRAINT custom_ui_components_pkey PRIMARY KEY (id);
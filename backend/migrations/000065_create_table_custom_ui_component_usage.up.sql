-- Migration for table: custom_ui_component_usage

CREATE SEQUENCE public.custom_ui_component_usage_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.custom_ui_component_usage (
    id integer NOT NULL,
    component_id integer NOT NULL,
    category_id integer,
    usage_context character varying(50) DEFAULT 'listing'::character varying NOT NULL,
    placement character varying(50) DEFAULT 'default'::character varying,
    priority integer DEFAULT 0,
    configuration jsonb DEFAULT '{}'::jsonb,
    conditions_logic jsonb,
    is_active boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    created_by integer,
    updated_by integer
);

ALTER SEQUENCE public.custom_ui_component_usage_id_seq OWNED BY public.custom_ui_component_usage.id;

CREATE INDEX idx_ui_comp_usage_active ON public.custom_ui_component_usage USING btree (is_active);

CREATE INDEX idx_ui_comp_usage_category ON public.custom_ui_component_usage USING btree (category_id);

CREATE INDEX idx_ui_comp_usage_component ON public.custom_ui_component_usage USING btree (component_id);

CREATE INDEX idx_ui_comp_usage_context ON public.custom_ui_component_usage USING btree (usage_context);

ALTER TABLE ONLY public.custom_ui_component_usage
    ADD CONSTRAINT custom_ui_component_usage_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.custom_ui_component_usage
    ADD CONSTRAINT custom_ui_component_usage_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.marketplace_categories(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.custom_ui_component_usage
    ADD CONSTRAINT custom_ui_component_usage_component_id_fkey FOREIGN KEY (component_id) REFERENCES public.custom_ui_components(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.custom_ui_component_usage
    ADD CONSTRAINT custom_ui_component_usage_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.users(id);

ALTER TABLE ONLY public.custom_ui_component_usage
    ADD CONSTRAINT custom_ui_component_usage_updated_by_fkey FOREIGN KEY (updated_by) REFERENCES public.users(id);

CREATE TRIGGER update_custom_ui_component_usage_updated_at BEFORE UPDATE ON public.custom_ui_component_usage FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
-- Migration for table: search_weights

CREATE SEQUENCE public.search_weights_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.search_weights (
    id bigint NOT NULL,
    field_name character varying(100) NOT NULL,
    weight double precision NOT NULL,
    search_type character varying(20) DEFAULT 'fulltext'::character varying NOT NULL,
    item_type character varying(20) DEFAULT 'global'::character varying NOT NULL,
    category_id integer,
    description text,
    is_active boolean DEFAULT true,
    version integer DEFAULT 1,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    created_by integer,
    updated_by integer,
    CONSTRAINT search_weights_item_type_check CHECK (((item_type)::text = ANY ((ARRAY['marketplace'::character varying, 'storefront'::character varying, 'global'::character varying])::text[]))),
    CONSTRAINT search_weights_search_type_check CHECK (((search_type)::text = ANY ((ARRAY['fulltext'::character varying, 'fuzzy'::character varying, 'exact'::character varying])::text[]))),
    CONSTRAINT search_weights_weight_check CHECK (((weight >= (0.0)::double precision) AND (weight <= (1.0)::double precision)))
);

ALTER SEQUENCE public.search_weights_id_seq OWNED BY public.search_weights.id;

CREATE INDEX idx_search_weights_category_id ON public.search_weights USING btree (category_id) WHERE (category_id IS NOT NULL);

CREATE INDEX idx_search_weights_field_name ON public.search_weights USING btree (field_name);

CREATE INDEX idx_search_weights_is_active ON public.search_weights USING btree (is_active) WHERE (is_active = true);

CREATE INDEX idx_search_weights_item_type ON public.search_weights USING btree (item_type);

CREATE UNIQUE INDEX idx_search_weights_unique_category ON public.search_weights USING btree (field_name, item_type, search_type, category_id) WHERE (category_id IS NOT NULL);

CREATE UNIQUE INDEX idx_search_weights_unique_global ON public.search_weights USING btree (field_name, item_type, search_type) WHERE (category_id IS NULL);

CREATE INDEX idx_search_weights_version ON public.search_weights USING btree (version);

ALTER TABLE ONLY public.search_weights
    ADD CONSTRAINT search_weights_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.search_weights
    ADD CONSTRAINT search_weights_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.admin_users(id) ON DELETE SET NULL;

ALTER TABLE ONLY public.search_weights
    ADD CONSTRAINT search_weights_updated_by_fkey FOREIGN KEY (updated_by) REFERENCES public.admin_users(id) ON DELETE SET NULL;

CREATE TRIGGER trigger_log_search_weight_changes AFTER UPDATE ON public.search_weights FOR EACH ROW EXECUTE FUNCTION public.log_search_weight_changes();

CREATE TRIGGER trigger_update_search_weights_updated_at BEFORE UPDATE ON public.search_weights FOR EACH ROW EXECUTE FUNCTION public.update_search_weights_updated_at();
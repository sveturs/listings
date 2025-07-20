-- Migration for table: product_variant_attributes

CREATE SEQUENCE public.product_variant_attributes_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.product_variant_attributes (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    display_name character varying(255) NOT NULL,
    type character varying(50) DEFAULT 'text'::character varying NOT NULL,
    is_required boolean DEFAULT false NOT NULL,
    sort_order integer DEFAULT 0 NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);

ALTER SEQUENCE public.product_variant_attributes_id_seq OWNED BY public.product_variant_attributes.id;

CREATE INDEX idx_product_variant_attributes_name ON public.product_variant_attributes USING btree (name);

ALTER TABLE ONLY public.product_variant_attributes
    ADD CONSTRAINT product_variant_attributes_pkey PRIMARY KEY (id);

CREATE TRIGGER trigger_update_product_variant_attributes_updated_at BEFORE UPDATE ON public.product_variant_attributes FOR EACH ROW EXECUTE FUNCTION public.update_product_variants_updated_at();
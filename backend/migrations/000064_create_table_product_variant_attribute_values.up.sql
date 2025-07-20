-- Migration for table: product_variant_attribute_values

CREATE SEQUENCE public.product_variant_attribute_values_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.product_variant_attribute_values (
    id integer NOT NULL,
    attribute_id integer NOT NULL,
    value character varying(255) NOT NULL,
    display_name character varying(255) NOT NULL,
    color_hex character varying(7),
    image_url text,
    sort_order integer DEFAULT 0 NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);

ALTER SEQUENCE public.product_variant_attribute_values_id_seq OWNED BY public.product_variant_attribute_values.id;

CREATE INDEX idx_product_variant_attribute_values_attribute_id ON public.product_variant_attribute_values USING btree (attribute_id);

CREATE INDEX idx_product_variant_attribute_values_value ON public.product_variant_attribute_values USING btree (value);

ALTER TABLE ONLY public.product_variant_attribute_values
    ADD CONSTRAINT product_variant_attribute_values_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.product_variant_attribute_values
    ADD CONSTRAINT product_variant_attribute_values_attribute_id_fkey FOREIGN KEY (attribute_id) REFERENCES public.product_variant_attributes(id) ON DELETE CASCADE;

CREATE TRIGGER trigger_update_product_variant_attribute_values_updated_at BEFORE UPDATE ON public.product_variant_attribute_values FOR EACH ROW EXECUTE FUNCTION public.update_product_variants_updated_at();
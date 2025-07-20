-- Migration for table: listing_attribute_values

CREATE SEQUENCE public.listing_attribute_values_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.listing_attribute_values (
    id integer NOT NULL,
    listing_id integer NOT NULL,
    attribute_id integer NOT NULL,
    text_value text,
    numeric_value numeric(20,5),
    boolean_value boolean,
    json_value jsonb,
    unit character varying(20) DEFAULT NULL::character varying
);

ALTER SEQUENCE public.listing_attribute_values_id_seq OWNED BY public.listing_attribute_values.id;

CREATE INDEX idx_attr_name_num_val ON public.listing_attribute_values USING btree (attribute_id, numeric_value) WHERE (numeric_value IS NOT NULL);

CREATE INDEX idx_attr_name_text_val ON public.listing_attribute_values USING btree (attribute_id, text_value) WHERE (text_value IS NOT NULL);

CREATE INDEX idx_attr_unit ON public.listing_attribute_values USING btree (unit) WHERE (unit IS NOT NULL);

CREATE INDEX idx_listing_attr_boolean ON public.listing_attribute_values USING btree (attribute_id, boolean_value) WHERE (boolean_value IS NOT NULL);

CREATE INDEX idx_listing_attr_listing_id ON public.listing_attribute_values USING btree (listing_id);

CREATE INDEX idx_listing_attr_numeric ON public.listing_attribute_values USING btree (attribute_id, numeric_value) WHERE (numeric_value IS NOT NULL);

CREATE INDEX idx_listing_attr_text ON public.listing_attribute_values USING btree (attribute_id, text_value) WHERE (text_value IS NOT NULL);

CREATE UNIQUE INDEX idx_listing_attr_unique ON public.listing_attribute_values USING btree (listing_id, attribute_id);

ALTER TABLE ONLY public.listing_attribute_values
    ADD CONSTRAINT listing_attribute_values_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.listing_attribute_values
    ADD CONSTRAINT listing_attribute_values_attribute_id_fkey FOREIGN KEY (attribute_id) REFERENCES public.category_attributes(id) ON DELETE CASCADE;
-- Migration for table: storefront_product_attributes

CREATE SEQUENCE public.storefront_product_attributes_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.storefront_product_attributes (
    id integer NOT NULL,
    product_id integer NOT NULL,
    attribute_id integer NOT NULL,
    is_enabled boolean DEFAULT true NOT NULL,
    is_required boolean DEFAULT false NOT NULL,
    custom_values jsonb DEFAULT '[]'::jsonb,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);

ALTER SEQUENCE public.storefront_product_attributes_id_seq OWNED BY public.storefront_product_attributes.id;

CREATE INDEX idx_storefront_product_attributes_attribute_id ON public.storefront_product_attributes USING btree (attribute_id);

CREATE INDEX idx_storefront_product_attributes_enabled ON public.storefront_product_attributes USING btree (is_enabled);

CREATE INDEX idx_storefront_product_attributes_product_id ON public.storefront_product_attributes USING btree (product_id);

ALTER TABLE ONLY public.storefront_product_attributes
    ADD CONSTRAINT storefront_product_attributes_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.storefront_product_attributes
    ADD CONSTRAINT storefront_product_attributes_product_id_attribute_id_key UNIQUE (product_id, attribute_id);

ALTER TABLE ONLY public.storefront_product_attributes
    ADD CONSTRAINT storefront_product_attributes_attribute_id_fkey FOREIGN KEY (attribute_id) REFERENCES public.product_variant_attributes(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.storefront_product_attributes
    ADD CONSTRAINT storefront_product_attributes_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.storefront_products(id) ON DELETE CASCADE;
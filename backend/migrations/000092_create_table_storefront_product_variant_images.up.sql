-- Migration for table: storefront_product_variant_images

CREATE SEQUENCE public.storefront_product_variant_images_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.storefront_product_variant_images (
    id integer NOT NULL,
    variant_id integer NOT NULL,
    image_url text NOT NULL,
    thumbnail_url text,
    alt_text character varying(255),
    display_order integer DEFAULT 0 NOT NULL,
    is_main boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);

ALTER SEQUENCE public.storefront_product_variant_images_id_seq OWNED BY public.storefront_product_variant_images.id;

CREATE INDEX idx_storefront_product_variant_images_is_main ON public.storefront_product_variant_images USING btree (is_main);

CREATE INDEX idx_storefront_product_variant_images_variant_id ON public.storefront_product_variant_images USING btree (variant_id);

ALTER TABLE ONLY public.storefront_product_variant_images
    ADD CONSTRAINT storefront_product_variant_images_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.storefront_product_variant_images
    ADD CONSTRAINT storefront_product_variant_images_variant_id_fkey FOREIGN KEY (variant_id) REFERENCES public.storefront_product_variants(id) ON DELETE CASCADE;
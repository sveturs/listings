-- Migration for table: storefront_product_images

CREATE SEQUENCE public.storefront_product_images_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.storefront_product_images (
    id integer NOT NULL,
    storefront_product_id integer NOT NULL,
    image_url text NOT NULL,
    thumbnail_url text NOT NULL,
    display_order integer DEFAULT 0 NOT NULL,
    is_default boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);

ALTER SEQUENCE public.storefront_product_images_id_seq OWNED BY public.storefront_product_images.id;

CREATE INDEX idx_storefront_product_images_display_order ON public.storefront_product_images USING btree (display_order);

CREATE INDEX idx_storefront_product_images_product_id ON public.storefront_product_images USING btree (storefront_product_id);

ALTER TABLE ONLY public.storefront_product_images
    ADD CONSTRAINT storefront_product_images_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.storefront_product_images
    ADD CONSTRAINT storefront_product_images_storefront_product_id_fkey FOREIGN KEY (storefront_product_id) REFERENCES public.storefront_products(id) ON DELETE CASCADE;
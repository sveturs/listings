-- Migration for table: marketplace_images

CREATE SEQUENCE public.marketplace_images_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.marketplace_images (
    id integer NOT NULL,
    listing_id integer,
    file_path character varying(255) NOT NULL,
    file_name character varying(255) NOT NULL,
    file_size integer NOT NULL,
    content_type character varying(100) NOT NULL,
    is_main boolean DEFAULT false,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    storage_type character varying(20) DEFAULT 'local'::character varying,
    storage_bucket character varying(100),
    public_url text
);

ALTER SEQUENCE public.marketplace_images_id_seq OWNED BY public.marketplace_images.id;

CREATE INDEX idx_marketplace_images_listing_main ON public.marketplace_images USING btree (listing_id, is_main) WHERE (is_main = true);

ALTER TABLE ONLY public.marketplace_images
    ADD CONSTRAINT marketplace_images_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.marketplace_images
    ADD CONSTRAINT marketplace_images_listing_id_fkey FOREIGN KEY (listing_id) REFERENCES public.marketplace_listings(id) ON DELETE CASCADE;
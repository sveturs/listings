-- Migration for table: listing_views

CREATE SEQUENCE public.listing_views_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.listing_views (
    id integer NOT NULL,
    listing_id integer NOT NULL,
    user_id integer,
    ip_hash character varying(255),
    view_time timestamp without time zone DEFAULT now(),
    CONSTRAINT at_least_one_identifier CHECK (((user_id IS NOT NULL) OR (ip_hash IS NOT NULL)))
);

ALTER SEQUENCE public.listing_views_id_seq OWNED BY public.listing_views.id;

CREATE INDEX idx_listing_views_listing_ip ON public.listing_views USING btree (listing_id, ip_hash);

CREATE INDEX idx_listing_views_listing_user ON public.listing_views USING btree (listing_id, user_id);

CREATE INDEX idx_listing_views_time ON public.listing_views USING btree (view_time);

ALTER TABLE ONLY public.listing_views
    ADD CONSTRAINT listing_view_uniqueness UNIQUE (listing_id, user_id);

ALTER TABLE ONLY public.listing_views
    ADD CONSTRAINT listing_views_pkey PRIMARY KEY (id);
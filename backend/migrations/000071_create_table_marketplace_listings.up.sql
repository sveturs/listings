-- Migration for table: marketplace_listings

CREATE SEQUENCE public.marketplace_listings_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.marketplace_listings (
    id integer NOT NULL,
    user_id integer,
    category_id integer,
    title character varying(255) NOT NULL,
    description text,
    price numeric(12,2),
    condition character varying(50),
    status character varying(20) DEFAULT 'active'::character varying,
    location character varying(255),
    latitude numeric(10,8),
    longitude numeric(11,8),
    address_city character varying(100),
    address_country character varying(100),
    views_count integer DEFAULT 0,
    show_on_map boolean DEFAULT true NOT NULL,
    original_language character varying(10) DEFAULT 'sr'::character varying,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    storefront_id integer,
    external_id character varying(255),
    metadata jsonb,
    needs_reindex boolean DEFAULT false
);

ALTER SEQUENCE public.marketplace_listings_id_seq OWNED BY public.marketplace_listings.id;

CREATE INDEX idx_listings_metadata_discount ON public.marketplace_listings USING gin (((metadata -> 'discount'::text)));

CREATE INDEX idx_marketplace_listings_category_status ON public.marketplace_listings USING btree (category_id, status) WHERE ((status)::text = 'active'::text);

CREATE INDEX idx_marketplace_listings_city ON public.marketplace_listings USING btree (address_city) WHERE (address_city IS NOT NULL);

CREATE INDEX idx_marketplace_listings_external_id ON public.marketplace_listings USING btree (external_id);

CREATE INDEX idx_marketplace_listings_external_id_storefront_id ON public.marketplace_listings USING btree (external_id, storefront_id);

CREATE INDEX idx_marketplace_listings_location ON public.marketplace_listings USING btree (latitude, longitude) WHERE ((latitude IS NOT NULL) AND (longitude IS NOT NULL));

CREATE INDEX idx_marketplace_listings_price ON public.marketplace_listings USING btree (price) WHERE ((price IS NOT NULL) AND ((status)::text = 'active'::text));

CREATE INDEX idx_marketplace_listings_status ON public.marketplace_listings USING btree (status);

CREATE INDEX idx_marketplace_listings_status_created ON public.marketplace_listings USING btree (status, created_at DESC) WHERE ((status)::text = 'active'::text);

CREATE INDEX idx_marketplace_listings_storefront ON public.marketplace_listings USING btree (storefront_id);

CREATE INDEX idx_marketplace_listings_title_gin ON public.marketplace_listings USING gin (to_tsvector('simple'::regconfig, (title)::text));

CREATE INDEX idx_marketplace_listings_user_status ON public.marketplace_listings USING btree (user_id, status, created_at DESC);

ALTER TABLE ONLY public.marketplace_listings
    ADD CONSTRAINT marketplace_listings_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.marketplace_listings
    ADD CONSTRAINT marketplace_listings_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.marketplace_categories(id);

ALTER TABLE ONLY public.marketplace_listings
    ADD CONSTRAINT marketplace_listings_storefront_id_fkey FOREIGN KEY (storefront_id) REFERENCES public.user_storefronts(id) ON DELETE SET NULL;

ALTER TABLE ONLY public.marketplace_listings
    ADD CONSTRAINT marketplace_listings_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);

CREATE TRIGGER preserve_review_origin_trigger BEFORE DELETE ON public.marketplace_listings FOR EACH ROW EXECUTE FUNCTION public.preserve_review_origin();

CREATE TRIGGER refresh_category_counts_delete AFTER DELETE ON public.marketplace_listings FOR EACH ROW EXECUTE FUNCTION public.refresh_category_listing_counts();

CREATE TRIGGER refresh_category_counts_insert AFTER INSERT ON public.marketplace_listings FOR EACH ROW EXECUTE FUNCTION public.refresh_category_listing_counts();

CREATE TRIGGER refresh_category_counts_update AFTER UPDATE ON public.marketplace_listings FOR EACH ROW WHEN (((old.status)::text IS DISTINCT FROM (new.status)::text)) EXECUTE FUNCTION public.refresh_category_listing_counts();

CREATE TRIGGER trg_new_listing_price_history AFTER INSERT ON public.marketplace_listings FOR EACH ROW EXECUTE FUNCTION public.update_price_history('create');

CREATE TRIGGER trg_update_listing_price_history AFTER UPDATE OF price ON public.marketplace_listings FOR EACH ROW WHEN ((old.price IS DISTINCT FROM new.price)) EXECUTE FUNCTION public.update_price_history('update');

CREATE TRIGGER trigger_marketplace_listings_cache_refresh AFTER INSERT OR DELETE OR UPDATE ON public.marketplace_listings FOR EACH ROW EXECUTE FUNCTION public.trigger_refresh_map_cache();
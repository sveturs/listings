-- Migration for table: listings_geo

CREATE SEQUENCE public.listings_geo_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.listings_geo (
    id bigint NOT NULL,
    listing_id bigint NOT NULL,
    location public.geography(Point,4326) NOT NULL,
    geohash character varying(12) NOT NULL,
    is_precise boolean DEFAULT true NOT NULL,
    blur_radius numeric(10,2) DEFAULT 0,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    address_components jsonb,
    geocoding_confidence numeric(3,2) DEFAULT 0.0,
    address_verified boolean DEFAULT false,
    input_method character varying(50) DEFAULT 'manual'::character varying,
    location_privacy character varying(20) DEFAULT 'exact'::character varying,
    blurred_location public.geography(Point,4326),
    formatted_address text,
    district_id uuid,
    municipality_id uuid,
    CONSTRAINT chk_geocoding_confidence CHECK (((geocoding_confidence >= 0.0) AND (geocoding_confidence <= 1.0))),
    CONSTRAINT chk_input_method CHECK (((input_method)::text = ANY ((ARRAY['manual'::character varying, 'geocoded'::character varying, 'map_click'::character varying, 'current_location'::character varying])::text[]))),
    CONSTRAINT chk_location_privacy CHECK (((location_privacy)::text = ANY ((ARRAY['exact'::character varying, 'street'::character varying, 'district'::character varying, 'city'::character varying])::text[])))
);

ALTER SEQUENCE public.listings_geo_id_seq OWNED BY public.listings_geo.id;

CREATE INDEX idx_listings_geo_address_components ON public.listings_geo USING gin (address_components);

CREATE INDEX idx_listings_geo_blurred_location ON public.listings_geo USING gist (blurred_location);

CREATE INDEX idx_listings_geo_confidence ON public.listings_geo USING btree (geocoding_confidence);

CREATE INDEX idx_listings_geo_district ON public.listings_geo USING btree (district_id);

CREATE INDEX idx_listings_geo_geohash ON public.listings_geo USING btree (geohash);

CREATE INDEX idx_listings_geo_geohash_precise ON public.listings_geo USING btree (geohash, is_precise);

CREATE INDEX idx_listings_geo_input_method ON public.listings_geo USING btree (input_method);

CREATE INDEX idx_listings_geo_is_precise ON public.listings_geo USING btree (is_precise);

CREATE INDEX idx_listings_geo_location ON public.listings_geo USING gist (location);

CREATE INDEX idx_listings_geo_municipality ON public.listings_geo USING btree (municipality_id);

CREATE INDEX idx_listings_geo_privacy ON public.listings_geo USING btree (location_privacy);

CREATE INDEX idx_listings_geo_verified ON public.listings_geo USING btree (address_verified);

ALTER TABLE ONLY public.listings_geo
    ADD CONSTRAINT listings_geo_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.listings_geo
    ADD CONSTRAINT uk_listings_geo_listing_id UNIQUE (listing_id);

ALTER TABLE ONLY public.listings_geo
    ADD CONSTRAINT listings_geo_district_id_fkey FOREIGN KEY (district_id) REFERENCES public.districts(id) ON DELETE SET NULL;

ALTER TABLE ONLY public.listings_geo
    ADD CONSTRAINT listings_geo_municipality_id_fkey FOREIGN KEY (municipality_id) REFERENCES public.municipalities(id) ON DELETE SET NULL;

CREATE TRIGGER trigger_update_listings_geo_updated_at BEFORE UPDATE ON public.listings_geo FOR EACH ROW EXECUTE FUNCTION public.update_listings_geo_updated_at();
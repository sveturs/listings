-- Migration for table: unified_geo

CREATE SEQUENCE public.unified_geo_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.unified_geo (
    id bigint NOT NULL,
    source_type public.geo_source_type NOT NULL,
    source_id bigint NOT NULL,
    location public.geography(Point,4326) NOT NULL,
    geohash character varying(12) NOT NULL,
    formatted_address text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    privacy_level public.location_privacy_level DEFAULT 'exact'::public.location_privacy_level,
    original_location public.geography(Point,4326),
    blur_radius_meters integer DEFAULT 0
);

ALTER SEQUENCE public.unified_geo_id_seq OWNED BY public.unified_geo.id;

CREATE INDEX idx_unified_geo_composite ON public.unified_geo USING btree (source_type, source_id);

CREATE INDEX idx_unified_geo_geohash ON public.unified_geo USING btree (geohash);

CREATE INDEX idx_unified_geo_location ON public.unified_geo USING gist (location);

CREATE INDEX idx_unified_geo_location_bounds ON public.unified_geo USING gist (location) WHERE (location IS NOT NULL);

CREATE INDEX idx_unified_geo_marketplace_active ON public.unified_geo USING btree (source_id) WHERE (source_type = 'marketplace_listing'::public.geo_source_type);

CREATE INDEX idx_unified_geo_source_id ON public.unified_geo USING btree (source_id);

CREATE INDEX idx_unified_geo_source_type ON public.unified_geo USING btree (source_type);

CREATE INDEX idx_unified_geo_storefront_active ON public.unified_geo USING btree (source_id) WHERE (source_type = 'storefront_product'::public.geo_source_type);

ALTER TABLE ONLY public.unified_geo
    ADD CONSTRAINT uk_unified_geo_source UNIQUE (source_type, source_id);

ALTER TABLE ONLY public.unified_geo
    ADD CONSTRAINT unified_geo_pkey PRIMARY KEY (id);

CREATE TRIGGER trigger_unified_geo_cache_refresh AFTER INSERT OR DELETE OR UPDATE ON public.unified_geo FOR EACH ROW EXECUTE FUNCTION public.trigger_refresh_map_cache();

CREATE TRIGGER trigger_update_unified_geo_updated_at BEFORE UPDATE ON public.unified_geo FOR EACH ROW EXECUTE FUNCTION public.update_unified_geo_updated_at();
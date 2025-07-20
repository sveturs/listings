-- Migration for table: geocoding_cache

CREATE SEQUENCE public.geocoding_cache_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.geocoding_cache (
    id bigint NOT NULL,
    input_address text NOT NULL,
    normalized_address text NOT NULL,
    location public.geography(Point,4326) NOT NULL,
    address_components jsonb NOT NULL,
    formatted_address text NOT NULL,
    confidence numeric(3,2) NOT NULL,
    provider character varying(50) DEFAULT 'mapbox'::character varying NOT NULL,
    language character varying(5) DEFAULT 'en'::character varying,
    country_code character varying(2),
    cache_hits bigint DEFAULT 0,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    expires_at timestamp without time zone DEFAULT (CURRENT_TIMESTAMP + '30 days'::interval) NOT NULL
);

ALTER SEQUENCE public.geocoding_cache_id_seq OWNED BY public.geocoding_cache.id;

CREATE INDEX idx_geocoding_cache_address_components ON public.geocoding_cache USING gin (address_components);

CREATE INDEX idx_geocoding_cache_cache_hits ON public.geocoding_cache USING btree (cache_hits);

CREATE INDEX idx_geocoding_cache_confidence ON public.geocoding_cache USING btree (confidence);

CREATE INDEX idx_geocoding_cache_country_lang ON public.geocoding_cache USING btree (country_code, language);

CREATE INDEX idx_geocoding_cache_expires_at ON public.geocoding_cache USING btree (expires_at);

CREATE INDEX idx_geocoding_cache_input_address ON public.geocoding_cache USING btree (input_address);

CREATE INDEX idx_geocoding_cache_location ON public.geocoding_cache USING gist (location);

CREATE INDEX idx_geocoding_cache_normalized ON public.geocoding_cache USING btree (normalized_address);

CREATE INDEX idx_geocoding_cache_provider ON public.geocoding_cache USING btree (provider);

ALTER TABLE ONLY public.geocoding_cache
    ADD CONSTRAINT geocoding_cache_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.geocoding_cache
    ADD CONSTRAINT uk_geocoding_cache_normalized UNIQUE (normalized_address, language, country_code);

CREATE TRIGGER trigger_cleanup_geocoding_cache AFTER INSERT ON public.geocoding_cache FOR EACH STATEMENT EXECUTE FUNCTION public.trigger_cleanup_geocoding_cache();

CREATE TRIGGER trigger_geocoding_cache_updated_at BEFORE UPDATE ON public.geocoding_cache FOR EACH ROW EXECUTE FUNCTION public.update_geocoding_cache_updated_at();
-- Migration for table: gis_poi_cache

CREATE SEQUENCE public.gis_poi_cache_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.gis_poi_cache (
    id integer NOT NULL,
    external_id character varying(255) NOT NULL,
    name character varying(255) NOT NULL,
    poi_type character varying(50) NOT NULL,
    location public.geography(Point,4326) NOT NULL,
    metadata jsonb,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    expires_at timestamp with time zone NOT NULL
);

ALTER SEQUENCE public.gis_poi_cache_id_seq OWNED BY public.gis_poi_cache.id;

CREATE INDEX idx_poi_cache_expires ON public.gis_poi_cache USING btree (expires_at);

CREATE INDEX idx_poi_cache_location ON public.gis_poi_cache USING gist (location);

CREATE INDEX idx_poi_cache_type ON public.gis_poi_cache USING btree (poi_type);

ALTER TABLE ONLY public.gis_poi_cache
    ADD CONSTRAINT gis_poi_cache_external_id_key UNIQUE (external_id);

ALTER TABLE ONLY public.gis_poi_cache
    ADD CONSTRAINT gis_poi_cache_pkey PRIMARY KEY (id);
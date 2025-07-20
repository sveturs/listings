-- Migration for table: gis_isochrone_cache

CREATE SEQUENCE public.gis_isochrone_cache_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.gis_isochrone_cache (
    id integer NOT NULL,
    center_point public.geography(Point,4326) NOT NULL,
    transport_mode character varying(20) NOT NULL,
    max_minutes integer NOT NULL,
    polygon public.geography(Polygon,4326) NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    expires_at timestamp with time zone NOT NULL
);

ALTER SEQUENCE public.gis_isochrone_cache_id_seq OWNED BY public.gis_isochrone_cache.id;

CREATE INDEX idx_isochrone_cache_center ON public.gis_isochrone_cache USING gist (center_point);

CREATE INDEX idx_isochrone_cache_expires ON public.gis_isochrone_cache USING btree (expires_at);

CREATE INDEX idx_isochrone_cache_lookup ON public.gis_isochrone_cache USING btree (transport_mode, max_minutes);

ALTER TABLE ONLY public.gis_isochrone_cache
    ADD CONSTRAINT gis_isochrone_cache_pkey PRIMARY KEY (id);
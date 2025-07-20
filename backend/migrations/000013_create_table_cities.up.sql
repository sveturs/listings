-- Migration for table: cities

CREATE TABLE public.cities (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name character varying(255) NOT NULL,
    slug character varying(255) NOT NULL,
    country_code character varying(2) DEFAULT 'RS'::character varying NOT NULL,
    boundary public.geometry(Polygon,4326),
    center_point public.geometry(Point,4326) NOT NULL,
    population integer,
    area_km2 numeric(10,2),
    postal_codes text[],
    has_districts boolean DEFAULT false NOT NULL,
    priority integer DEFAULT 100 NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);

CREATE INDEX idx_cities_boundary ON public.cities USING gist (boundary);

CREATE INDEX idx_cities_center ON public.cities USING gist (center_point);

CREATE INDEX idx_cities_country ON public.cities USING btree (country_code);

CREATE INDEX idx_cities_has_districts ON public.cities USING btree (has_districts);

CREATE INDEX idx_cities_name ON public.cities USING btree (name);

CREATE INDEX idx_cities_priority ON public.cities USING btree (priority);

CREATE INDEX idx_cities_slug ON public.cities USING btree (slug);

ALTER TABLE ONLY public.cities
    ADD CONSTRAINT cities_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.cities
    ADD CONSTRAINT cities_slug_key UNIQUE (slug);
-- Migration for table: districts

CREATE TABLE public.districts (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name character varying(255) NOT NULL,
    city_id uuid,
    country_code character varying(2) DEFAULT 'RS'::character varying NOT NULL,
    boundary public.geometry(Polygon,4326),
    center_point public.geometry(Point,4326),
    population integer,
    area_km2 numeric(10,2),
    postal_codes text[],
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);

CREATE INDEX idx_districts_boundary ON public.districts USING gist (boundary);

CREATE INDEX idx_districts_center ON public.districts USING gist (center_point);

CREATE INDEX idx_districts_city_id ON public.districts USING btree (city_id);

CREATE INDEX idx_districts_country ON public.districts USING btree (country_code);

CREATE INDEX idx_districts_name ON public.districts USING btree (name);

ALTER TABLE ONLY public.districts
    ADD CONSTRAINT districts_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.districts
    ADD CONSTRAINT fk_districts_city_id FOREIGN KEY (city_id) REFERENCES public.cities(id) ON DELETE SET NULL;
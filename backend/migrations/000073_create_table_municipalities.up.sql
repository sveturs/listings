-- Migration for table: municipalities

CREATE TABLE public.municipalities (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name character varying(255) NOT NULL,
    district_id uuid,
    country_code character varying(2) DEFAULT 'RS'::character varying NOT NULL,
    boundary public.geometry(Polygon,4326),
    center_point public.geometry(Point,4326),
    population integer,
    area_km2 numeric(10,2),
    postal_codes text[],
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);

CREATE INDEX idx_municipalities_boundary ON public.municipalities USING gist (boundary);

CREATE INDEX idx_municipalities_center ON public.municipalities USING gist (center_point);

CREATE INDEX idx_municipalities_district ON public.municipalities USING btree (district_id);

CREATE INDEX idx_municipalities_name ON public.municipalities USING btree (name);

ALTER TABLE ONLY public.municipalities
    ADD CONSTRAINT municipalities_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.municipalities
    ADD CONSTRAINT municipalities_district_id_fkey FOREIGN KEY (district_id) REFERENCES public.districts(id) ON DELETE SET NULL;
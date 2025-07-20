-- Migration for table: districts_leskovac_backup

CREATE TABLE public.districts_leskovac_backup (
    id uuid,
    name character varying(255),
    city_id uuid,
    country_code character varying(2),
    boundary public.geometry(Polygon,4326),
    center_point public.geometry(Point,4326),
    population integer,
    area_km2 numeric(10,2),
    postal_codes text[],
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    city_name character varying(255)
);
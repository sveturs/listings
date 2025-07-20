-- Migration for table: storefront_delivery_options

CREATE SEQUENCE public.storefront_delivery_options_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.storefront_delivery_options (
    id integer NOT NULL,
    storefront_id integer NOT NULL,
    name character varying(255) NOT NULL,
    description text,
    base_price numeric(10,2) DEFAULT 0.00 NOT NULL,
    price_per_km numeric(10,2) DEFAULT 0.00,
    price_per_kg numeric(10,2) DEFAULT 0.00,
    free_above_amount numeric(10,2),
    min_order_amount numeric(10,2),
    max_weight_kg numeric(10,2),
    max_distance_km numeric(10,2),
    estimated_days_min integer DEFAULT 1,
    estimated_days_max integer DEFAULT 3,
    zones jsonb DEFAULT '[]'::jsonb,
    available_days jsonb DEFAULT '[1, 2, 3, 4, 5]'::jsonb,
    cutoff_time time without time zone,
    provider character varying(50),
    provider_config jsonb DEFAULT '{}'::jsonb,
    is_active boolean DEFAULT true,
    display_order integer DEFAULT 0,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

ALTER SEQUENCE public.storefront_delivery_options_id_seq OWNED BY public.storefront_delivery_options.id;

CREATE INDEX idx_delivery_is_active ON public.storefront_delivery_options USING btree (is_active);

CREATE INDEX idx_delivery_storefront_id ON public.storefront_delivery_options USING btree (storefront_id);

ALTER TABLE ONLY public.storefront_delivery_options
    ADD CONSTRAINT storefront_delivery_options_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.storefront_delivery_options
    ADD CONSTRAINT storefront_delivery_options_storefront_id_fkey FOREIGN KEY (storefront_id) REFERENCES public.storefronts(id) ON DELETE CASCADE;
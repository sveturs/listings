-- Migration for table: storefronts

CREATE SEQUENCE public.storefronts_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.storefronts (
    id integer NOT NULL,
    user_id integer NOT NULL,
    slug character varying(100) NOT NULL,
    name character varying(255) NOT NULL,
    description text,
    logo_url character varying(500),
    banner_url character varying(500),
    theme jsonb DEFAULT '{"layout": "grid", "primaryColor": "#1976d2"}'::jsonb,
    phone character varying(50),
    email character varying(255),
    website character varying(255),
    address text,
    city character varying(100),
    postal_code character varying(20),
    country character varying(2) DEFAULT 'RS'::character varying,
    latitude numeric(10,8),
    longitude numeric(11,8),
    settings jsonb DEFAULT '{}'::jsonb,
    seo_meta jsonb DEFAULT '{}'::jsonb,
    is_active boolean DEFAULT false,
    is_verified boolean DEFAULT false,
    verification_date timestamp without time zone,
    rating numeric(3,2) DEFAULT 0.00,
    reviews_count integer DEFAULT 0,
    products_count integer DEFAULT 0,
    sales_count integer DEFAULT 0,
    views_count integer DEFAULT 0,
    subscription_plan character varying(50) DEFAULT 'starter'::character varying,
    subscription_expires_at timestamp without time zone,
    commission_rate numeric(5,2) DEFAULT 3.00,
    ai_agent_enabled boolean DEFAULT false,
    ai_agent_config jsonb DEFAULT '{}'::jsonb,
    live_shopping_enabled boolean DEFAULT false,
    group_buying_enabled boolean DEFAULT false,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    formatted_address text,
    geo_strategy public.storefront_geo_strategy DEFAULT 'storefront_location'::public.storefront_geo_strategy,
    default_privacy_level public.location_privacy_level DEFAULT 'exact'::public.location_privacy_level,
    address_verified boolean DEFAULT false
);

ALTER SEQUENCE public.storefronts_id_seq OWNED BY public.storefronts.id;

CREATE INDEX idx_storefronts_city ON public.storefronts USING btree (city);

CREATE INDEX idx_storefronts_coordinates ON public.storefronts USING btree (latitude, longitude) WHERE ((latitude IS NOT NULL) AND (longitude IS NOT NULL));

CREATE INDEX idx_storefronts_geo_strategy ON public.storefronts USING btree (geo_strategy);

CREATE INDEX idx_storefronts_is_active ON public.storefronts USING btree (is_active);

CREATE INDEX idx_storefronts_rating ON public.storefronts USING btree (rating DESC);

CREATE INDEX idx_storefronts_slug ON public.storefronts USING btree (slug);

CREATE INDEX idx_storefronts_user_id ON public.storefronts USING btree (user_id);

ALTER TABLE ONLY public.storefronts
    ADD CONSTRAINT storefronts_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.storefronts
    ADD CONSTRAINT storefronts_slug_key UNIQUE (slug);

CREATE TRIGGER trigger_update_storefront_products_geo AFTER UPDATE ON public.storefronts FOR EACH ROW EXECUTE FUNCTION public.update_storefront_products_geo();
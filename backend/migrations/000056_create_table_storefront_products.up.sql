-- Migration for table: storefront_products

CREATE SEQUENCE public.storefront_products_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.storefront_products (
    id integer NOT NULL,
    storefront_id integer NOT NULL,
    name character varying(255) NOT NULL,
    description text NOT NULL,
    price numeric(15,2) NOT NULL,
    currency character(3) DEFAULT 'USD'::bpchar NOT NULL,
    category_id integer NOT NULL,
    sku character varying(100),
    barcode character varying(100),
    stock_quantity integer DEFAULT 0 NOT NULL,
    stock_status character varying(20) DEFAULT 'in_stock'::character varying NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    attributes jsonb DEFAULT '{}'::jsonb,
    view_count integer DEFAULT 0 NOT NULL,
    sold_count integer DEFAULT 0 NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    has_individual_location boolean DEFAULT false,
    individual_address text,
    individual_latitude numeric(10,8),
    individual_longitude numeric(11,8),
    location_privacy public.location_privacy_level DEFAULT 'exact'::public.location_privacy_level,
    show_on_map boolean DEFAULT true,
    CONSTRAINT storefront_products_price_check CHECK ((price >= (0)::numeric)),
    CONSTRAINT storefront_products_stock_quantity_check CHECK ((stock_quantity >= 0)),
    CONSTRAINT storefront_products_stock_status_check CHECK (((stock_status)::text = ANY ((ARRAY['in_stock'::character varying, 'low_stock'::character varying, 'out_of_stock'::character varying])::text[])))
);

ALTER SEQUENCE public.storefront_products_id_seq OWNED BY public.storefront_products.id;

CREATE INDEX idx_storefront_products_barcode ON public.storefront_products USING btree (barcode) WHERE (barcode IS NOT NULL);

CREATE INDEX idx_storefront_products_category_id ON public.storefront_products USING btree (category_id);

CREATE INDEX idx_storefront_products_individual_location ON public.storefront_products USING btree (has_individual_location);

CREATE INDEX idx_storefront_products_is_active ON public.storefront_products USING btree (is_active);

CREATE INDEX idx_storefront_products_name_gin ON public.storefront_products USING gin (to_tsvector('simple'::regconfig, (name)::text));

CREATE INDEX idx_storefront_products_privacy ON public.storefront_products USING btree (location_privacy);

CREATE INDEX idx_storefront_products_show_on_map ON public.storefront_products USING btree (show_on_map);

CREATE INDEX idx_storefront_products_sku ON public.storefront_products USING btree (sku) WHERE (sku IS NOT NULL);

CREATE INDEX idx_storefront_products_stock_status ON public.storefront_products USING btree (stock_status);

CREATE INDEX idx_storefront_products_storefront_id ON public.storefront_products USING btree (storefront_id);

CREATE UNIQUE INDEX unique_storefront_product_barcode ON public.storefront_products USING btree (storefront_id, barcode) WHERE (barcode IS NOT NULL);

CREATE UNIQUE INDEX unique_storefront_product_sku ON public.storefront_products USING btree (storefront_id, sku) WHERE (sku IS NOT NULL);

ALTER TABLE ONLY public.storefront_products
    ADD CONSTRAINT storefront_products_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.storefront_products
    ADD CONSTRAINT storefront_products_storefront_id_fkey FOREIGN KEY (storefront_id) REFERENCES public.storefronts(id) ON DELETE CASCADE;

CREATE TRIGGER update_stock_status_trigger BEFORE INSERT OR UPDATE OF stock_quantity ON public.storefront_products FOR EACH ROW EXECUTE FUNCTION public.update_product_stock_status();
-- Migration for table: storefront_product_variants

CREATE SEQUENCE public.storefront_product_variants_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.storefront_product_variants (
    id integer NOT NULL,
    product_id integer NOT NULL,
    sku character varying(100),
    barcode character varying(100),
    price numeric(15,2),
    compare_at_price numeric(15,2),
    cost_price numeric(15,2),
    stock_quantity integer DEFAULT 0 NOT NULL,
    stock_status character varying(20) DEFAULT 'in_stock'::character varying NOT NULL,
    low_stock_threshold integer DEFAULT 5,
    variant_attributes jsonb DEFAULT '{}'::jsonb NOT NULL,
    weight numeric(10,3),
    dimensions jsonb,
    is_active boolean DEFAULT true NOT NULL,
    is_default boolean DEFAULT false NOT NULL,
    view_count integer DEFAULT 0 NOT NULL,
    sold_count integer DEFAULT 0 NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT storefront_product_variants_stock_quantity_check CHECK ((stock_quantity >= 0)),
    CONSTRAINT storefront_product_variants_stock_status_check CHECK (((stock_status)::text = ANY (ARRAY[('in_stock'::character varying)::text, ('low_stock'::character varying)::text, ('out_of_stock'::character varying)::text])))
);

ALTER SEQUENCE public.storefront_product_variants_id_seq OWNED BY public.storefront_product_variants.id;

CREATE INDEX idx_storefront_product_variants_attributes_gin ON public.storefront_product_variants USING gin (variant_attributes);

CREATE INDEX idx_storefront_product_variants_barcode ON public.storefront_product_variants USING btree (barcode) WHERE (barcode IS NOT NULL);

CREATE UNIQUE INDEX idx_storefront_product_variants_default_unique ON public.storefront_product_variants USING btree (product_id) WHERE (is_default = true);

CREATE INDEX idx_storefront_product_variants_is_active ON public.storefront_product_variants USING btree (is_active);

CREATE INDEX idx_storefront_product_variants_is_default ON public.storefront_product_variants USING btree (is_default);

CREATE INDEX idx_storefront_product_variants_product_id ON public.storefront_product_variants USING btree (product_id);

CREATE INDEX idx_storefront_product_variants_sku ON public.storefront_product_variants USING btree (sku) WHERE (sku IS NOT NULL);

CREATE INDEX idx_storefront_product_variants_stock_status ON public.storefront_product_variants USING btree (stock_status);

ALTER TABLE ONLY public.storefront_product_variants
    ADD CONSTRAINT storefront_product_variants_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.storefront_product_variants
    ADD CONSTRAINT storefront_product_variants_sku_key UNIQUE (sku);

ALTER TABLE ONLY public.storefront_product_variants
    ADD CONSTRAINT storefront_product_variants_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.storefront_products(id) ON DELETE CASCADE;

CREATE TRIGGER trigger_update_storefront_product_variants_updated_at BEFORE UPDATE ON public.storefront_product_variants FOR EACH ROW EXECUTE FUNCTION public.update_product_variants_updated_at();
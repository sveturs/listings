-- Migration for table: storefront_orders

CREATE SEQUENCE public.storefront_orders_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.storefront_orders (
    id bigint NOT NULL,
    order_number character varying(32) NOT NULL,
    storefront_id integer,
    customer_id integer,
    payment_transaction_id bigint,
    subtotal_amount numeric(12,2) NOT NULL,
    shipping_amount numeric(12,2) DEFAULT 0,
    tax_amount numeric(12,2) DEFAULT 0,
    total_amount numeric(12,2) NOT NULL,
    commission_amount numeric(12,2) NOT NULL,
    seller_amount numeric(12,2) NOT NULL,
    currency character(3) DEFAULT 'RSD'::bpchar,
    status character varying(50) DEFAULT 'pending'::character varying NOT NULL,
    escrow_release_date date,
    escrow_days integer DEFAULT 3,
    shipping_address jsonb,
    billing_address jsonb,
    shipping_method character varying(100),
    shipping_provider character varying(50),
    tracking_number character varying(100),
    customer_notes text,
    seller_notes text,
    confirmed_at timestamp without time zone,
    shipped_at timestamp without time zone,
    delivered_at timestamp without time zone,
    cancelled_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    payment_method character varying(50) DEFAULT 'allsecure'::character varying NOT NULL,
    payment_status character varying(50) DEFAULT 'pending'::character varying NOT NULL,
    notes text,
    metadata jsonb DEFAULT '{}'::jsonb,
    discount numeric(12,2) DEFAULT 0
);

ALTER SEQUENCE public.storefront_orders_id_seq OWNED BY public.storefront_orders.id;

CREATE INDEX idx_storefront_orders_customer ON public.storefront_orders USING btree (customer_id, created_at DESC);

CREATE INDEX idx_storefront_orders_escrow_date ON public.storefront_orders USING btree (escrow_release_date) WHERE (escrow_release_date IS NOT NULL);

CREATE INDEX idx_storefront_orders_status ON public.storefront_orders USING btree (status);

CREATE INDEX idx_storefront_orders_storefront ON public.storefront_orders USING btree (storefront_id, created_at DESC);

ALTER TABLE ONLY public.storefront_orders
    ADD CONSTRAINT storefront_orders_order_number_key UNIQUE (order_number);

ALTER TABLE ONLY public.storefront_orders
    ADD CONSTRAINT storefront_orders_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.storefront_orders
    ADD CONSTRAINT storefront_orders_payment_transaction_id_fkey FOREIGN KEY (payment_transaction_id) REFERENCES public.payment_transactions(id) ON DELETE SET NULL;

ALTER TABLE ONLY public.storefront_orders
    ADD CONSTRAINT storefront_orders_storefront_id_fkey FOREIGN KEY (storefront_id) REFERENCES public.storefronts(id) ON DELETE RESTRICT;

CREATE TRIGGER calculate_escrow_release_date_trigger BEFORE INSERT OR UPDATE ON public.storefront_orders FOR EACH ROW EXECUTE FUNCTION public.calculate_escrow_release_date();

CREATE TRIGGER set_order_number_trigger BEFORE INSERT ON public.storefront_orders FOR EACH ROW EXECUTE FUNCTION public.set_order_number();
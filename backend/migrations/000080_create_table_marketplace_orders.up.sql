-- Migration for table: marketplace_orders

CREATE SEQUENCE public.marketplace_orders_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.marketplace_orders (
    id integer NOT NULL,
    buyer_id integer NOT NULL,
    seller_id integer NOT NULL,
    listing_id integer NOT NULL,
    item_price numeric(10,2) NOT NULL,
    platform_fee_rate numeric(5,2) DEFAULT 5.00,
    platform_fee_amount numeric(10,2) NOT NULL,
    seller_payout_amount numeric(10,2) NOT NULL,
    payment_transaction_id integer,
    status character varying(50) DEFAULT 'pending'::character varying,
    protection_period_days integer DEFAULT 7,
    protection_expires_at timestamp with time zone,
    shipping_method character varying(100),
    tracking_number character varying(255),
    shipped_at timestamp with time zone,
    delivered_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    CONSTRAINT marketplace_orders_status_check CHECK (((status)::text = ANY (ARRAY[('pending'::character varying)::text, ('paid'::character varying)::text, ('shipped'::character varying)::text, ('delivered'::character varying)::text, ('completed'::character varying)::text, ('disputed'::character varying)::text, ('cancelled'::character varying)::text, ('refunded'::character varying)::text])))
);

ALTER SEQUENCE public.marketplace_orders_id_seq OWNED BY public.marketplace_orders.id;

CREATE INDEX idx_marketplace_orders_buyer ON public.marketplace_orders USING btree (buyer_id);

CREATE INDEX idx_marketplace_orders_buyer_id ON public.marketplace_orders USING btree (buyer_id);

CREATE INDEX idx_marketplace_orders_created_at ON public.marketplace_orders USING btree (created_at DESC);

CREATE INDEX idx_marketplace_orders_listing ON public.marketplace_orders USING btree (listing_id);

CREATE INDEX idx_marketplace_orders_listing_id ON public.marketplace_orders USING btree (listing_id);

CREATE INDEX idx_marketplace_orders_payment_transaction_id ON public.marketplace_orders USING btree (payment_transaction_id);

CREATE INDEX idx_marketplace_orders_protection ON public.marketplace_orders USING btree (protection_expires_at) WHERE ((status)::text = ANY (ARRAY[('delivered'::character varying)::text, ('shipped'::character varying)::text]));

CREATE INDEX idx_marketplace_orders_protection_expires_at ON public.marketplace_orders USING btree (protection_expires_at) WHERE (protection_expires_at IS NOT NULL);

CREATE INDEX idx_marketplace_orders_seller ON public.marketplace_orders USING btree (seller_id);

CREATE INDEX idx_marketplace_orders_seller_id ON public.marketplace_orders USING btree (seller_id);

CREATE INDEX idx_marketplace_orders_status ON public.marketplace_orders USING btree (status);

ALTER TABLE ONLY public.marketplace_orders
    ADD CONSTRAINT marketplace_orders_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.marketplace_orders
    ADD CONSTRAINT marketplace_orders_payment_transaction_id_fkey FOREIGN KEY (payment_transaction_id) REFERENCES public.payment_transactions(id);

CREATE TRIGGER marketplace_orders_updated_at_trigger BEFORE UPDATE ON public.marketplace_orders FOR EACH ROW EXECUTE FUNCTION public.update_marketplace_orders_updated_at();
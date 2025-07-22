-- Migration for table: payment_transactions

CREATE SEQUENCE public.payment_transactions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.payment_transactions (
    id bigint NOT NULL,
    gateway_id integer,
    user_id integer,
    listing_id integer,
    order_reference character varying(255) NOT NULL,
    gateway_transaction_id character varying(255),
    gateway_reference_id character varying(255),
    amount numeric(12,2) NOT NULL,
    currency character varying(3) DEFAULT 'RSD'::character varying,
    marketplace_commission numeric(12,2),
    seller_amount numeric(12,2),
    status character varying(50) DEFAULT 'pending'::character varying,
    gateway_status character varying(50),
    payment_method character varying(50),
    customer_email character varying(255),
    description text,
    gateway_response jsonb,
    error_details jsonb,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    authorized_at timestamp with time zone,
    captured_at timestamp with time zone,
    failed_at timestamp with time zone,
    source_type character varying(20) DEFAULT 'marketplace_listing'::character varying,
    source_id bigint,
    storefront_id integer,
    capture_mode character varying(20) DEFAULT 'manual'::character varying,
    auto_capture_at timestamp with time zone,
    capture_deadline_at timestamp with time zone,
    capture_attempted_at timestamp with time zone,
    capture_attempts integer DEFAULT 0,
    CONSTRAINT payment_transactions_amount_positive CHECK ((amount > (0)::numeric)),
    CONSTRAINT payment_transactions_capture_mode_check CHECK (((capture_mode)::text = ANY (ARRAY[('auto'::character varying)::text, ('manual'::character varying)::text]))),
    CONSTRAINT payment_transactions_status_valid CHECK (((status)::text = ANY (ARRAY[('pending'::character varying)::text, ('authorized'::character varying)::text, ('captured'::character varying)::text, ('failed'::character varying)::text, ('refunded'::character varying)::text, ('voided'::character varying)::text])))
);

ALTER SEQUENCE public.payment_transactions_id_seq OWNED BY public.payment_transactions.id;

CREATE INDEX idx_payment_transactions_auto_capture ON public.payment_transactions USING btree (auto_capture_at, status) WHERE (((status)::text = 'authorized'::text) AND ((capture_mode)::text = 'auto'::text));

CREATE INDEX idx_payment_transactions_created_at ON public.payment_transactions USING btree (created_at);

CREATE INDEX idx_payment_transactions_gateway_transaction_id ON public.payment_transactions USING btree (gateway_transaction_id);

CREATE INDEX idx_payment_transactions_listing_id ON public.payment_transactions USING btree (listing_id);

CREATE INDEX idx_payment_transactions_order_reference ON public.payment_transactions USING btree (order_reference);

CREATE INDEX idx_payment_transactions_source ON public.payment_transactions USING btree (source_type, source_id);

CREATE INDEX idx_payment_transactions_status ON public.payment_transactions USING btree (status);

CREATE INDEX idx_payment_transactions_user_id ON public.payment_transactions USING btree (user_id);

ALTER TABLE ONLY public.payment_transactions
    ADD CONSTRAINT payment_transactions_order_reference_key UNIQUE (order_reference);

ALTER TABLE ONLY public.payment_transactions
    ADD CONSTRAINT payment_transactions_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.payment_transactions
    ADD CONSTRAINT payment_transactions_gateway_id_fkey FOREIGN KEY (gateway_id) REFERENCES public.payment_gateways(id);

ALTER TABLE ONLY public.payment_transactions
    ADD CONSTRAINT payment_transactions_storefront_id_fkey FOREIGN KEY (storefront_id) REFERENCES public.storefronts(id);
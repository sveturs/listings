-- Migration for table: escrow_payments

CREATE SEQUENCE public.escrow_payments_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.escrow_payments (
    id bigint NOT NULL,
    payment_transaction_id bigint,
    seller_id integer,
    buyer_id integer,
    listing_id integer,
    amount numeric(12,2) NOT NULL,
    marketplace_commission numeric(12,2) NOT NULL,
    seller_amount numeric(12,2) NOT NULL,
    status character varying(50) DEFAULT 'held'::character varying,
    release_date timestamp with time zone,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    CONSTRAINT escrow_payments_amount_positive CHECK ((amount > (0)::numeric)),
    CONSTRAINT escrow_payments_amounts_sum CHECK (((marketplace_commission + seller_amount) = amount)),
    CONSTRAINT escrow_payments_status_valid CHECK (((status)::text = ANY ((ARRAY['held'::character varying, 'released'::character varying, 'refunded'::character varying])::text[])))
);

ALTER SEQUENCE public.escrow_payments_id_seq OWNED BY public.escrow_payments.id;

CREATE INDEX idx_escrow_payments_buyer_id ON public.escrow_payments USING btree (buyer_id);

CREATE INDEX idx_escrow_payments_payment_transaction_id ON public.escrow_payments USING btree (payment_transaction_id);

CREATE INDEX idx_escrow_payments_seller_id ON public.escrow_payments USING btree (seller_id);

CREATE INDEX idx_escrow_payments_status ON public.escrow_payments USING btree (status);

ALTER TABLE ONLY public.escrow_payments
    ADD CONSTRAINT escrow_payments_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.escrow_payments
    ADD CONSTRAINT escrow_payments_payment_transaction_id_fkey FOREIGN KEY (payment_transaction_id) REFERENCES public.payment_transactions(id) ON DELETE CASCADE;
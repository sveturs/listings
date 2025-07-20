-- Migration for table: balance_transactions

CREATE SEQUENCE public.balance_transactions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.balance_transactions (
    id integer NOT NULL,
    user_id integer NOT NULL,
    type character varying(20) NOT NULL,
    amount numeric(12,2) NOT NULL,
    currency character varying(3) DEFAULT 'RSD'::character varying NOT NULL,
    status character varying(20) DEFAULT 'pending'::character varying NOT NULL,
    payment_method character varying(50),
    payment_details jsonb,
    description text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    completed_at timestamp without time zone
);

ALTER SEQUENCE public.balance_transactions_id_seq OWNED BY public.balance_transactions.id;

CREATE INDEX idx_transactions_created ON public.balance_transactions USING btree (created_at);

CREATE INDEX idx_transactions_status ON public.balance_transactions USING btree (status);

CREATE INDEX idx_transactions_user ON public.balance_transactions USING btree (user_id);

ALTER TABLE ONLY public.balance_transactions
    ADD CONSTRAINT balance_transactions_pkey PRIMARY KEY (id);
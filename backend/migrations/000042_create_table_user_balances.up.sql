-- Migration for table: user_balances

CREATE TABLE public.user_balances (
    user_id integer NOT NULL,
    balance numeric(12,2) DEFAULT 0 NOT NULL,
    frozen_balance numeric(12,2) DEFAULT 0 NOT NULL,
    currency character varying(3) DEFAULT 'RSD'::character varying NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE ONLY public.user_balances
    ADD CONSTRAINT user_balances_pkey PRIMARY KEY (user_id);
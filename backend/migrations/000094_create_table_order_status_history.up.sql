-- Migration for table: order_status_history

CREATE SEQUENCE public.order_status_history_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.order_status_history (
    id integer NOT NULL,
    order_id integer NOT NULL,
    old_status character varying(50),
    new_status character varying(50) NOT NULL,
    reason text,
    created_by integer,
    created_at timestamp with time zone DEFAULT now()
);

ALTER SEQUENCE public.order_status_history_id_seq OWNED BY public.order_status_history.id;

CREATE INDEX idx_order_status_history_created_at ON public.order_status_history USING btree (created_at DESC);

CREATE INDEX idx_order_status_history_order_id ON public.order_status_history USING btree (order_id);

ALTER TABLE ONLY public.order_status_history
    ADD CONSTRAINT order_status_history_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.order_status_history
    ADD CONSTRAINT order_status_history_order_id_fkey FOREIGN KEY (order_id) REFERENCES public.marketplace_orders(id);
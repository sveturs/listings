-- Migration for table: order_messages

CREATE SEQUENCE public.order_messages_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.order_messages (
    id integer NOT NULL,
    order_id integer NOT NULL,
    sender_id integer NOT NULL,
    message_type character varying(50) DEFAULT 'text'::character varying,
    content text NOT NULL,
    metadata jsonb,
    created_at timestamp with time zone DEFAULT now(),
    CONSTRAINT order_messages_message_type_check CHECK (((message_type)::text = ANY ((ARRAY['text'::character varying, 'shipping_update'::character varying, 'dispute_opened'::character varying, 'dispute_message'::character varying, 'system'::character varying])::text[])))
);

ALTER SEQUENCE public.order_messages_id_seq OWNED BY public.order_messages.id;

CREATE INDEX idx_order_messages_created_at ON public.order_messages USING btree (created_at DESC);

CREATE INDEX idx_order_messages_order_id ON public.order_messages USING btree (order_id);

CREATE INDEX idx_order_messages_sender_id ON public.order_messages USING btree (sender_id);

ALTER TABLE ONLY public.order_messages
    ADD CONSTRAINT order_messages_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.order_messages
    ADD CONSTRAINT order_messages_order_id_fkey FOREIGN KEY (order_id) REFERENCES public.marketplace_orders(id);
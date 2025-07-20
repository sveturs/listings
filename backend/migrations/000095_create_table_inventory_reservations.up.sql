-- Migration for table: inventory_reservations

CREATE SEQUENCE public.inventory_reservations_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.inventory_reservations (
    id bigint NOT NULL,
    order_id bigint,
    product_id bigint,
    variant_id bigint,
    quantity integer NOT NULL,
    status character varying(20) DEFAULT 'active'::character varying,
    expires_at timestamp without time zone DEFAULT (CURRENT_TIMESTAMP + '02:00:00'::interval) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    released_at timestamp without time zone,
    CONSTRAINT inventory_reservations_quantity_check CHECK ((quantity > 0))
);

ALTER SEQUENCE public.inventory_reservations_id_seq OWNED BY public.inventory_reservations.id;

CREATE INDEX idx_inventory_reservations_expires ON public.inventory_reservations USING btree (expires_at) WHERE ((status)::text = 'active'::text);

CREATE INDEX idx_inventory_reservations_product ON public.inventory_reservations USING btree (product_id);

ALTER TABLE ONLY public.inventory_reservations
    ADD CONSTRAINT inventory_reservations_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.inventory_reservations
    ADD CONSTRAINT inventory_reservations_order_id_fkey FOREIGN KEY (order_id) REFERENCES public.storefront_orders(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.inventory_reservations
    ADD CONSTRAINT inventory_reservations_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.storefront_products(id) ON DELETE CASCADE;
-- Migration for table: storefront_order_items

CREATE SEQUENCE public.storefront_order_items_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.storefront_order_items (
    id bigint NOT NULL,
    order_id bigint,
    product_id bigint,
    variant_id bigint,
    product_name character varying(255) NOT NULL,
    product_sku character varying(100),
    variant_name character varying(255),
    quantity integer NOT NULL,
    price_per_unit numeric(12,2) NOT NULL,
    total_price numeric(12,2) NOT NULL,
    product_attributes jsonb,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT storefront_order_items_quantity_check CHECK ((quantity > 0))
);

ALTER SEQUENCE public.storefront_order_items_id_seq OWNED BY public.storefront_order_items.id;

ALTER TABLE ONLY public.storefront_order_items
    ADD CONSTRAINT storefront_order_items_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.storefront_order_items
    ADD CONSTRAINT storefront_order_items_order_id_fkey FOREIGN KEY (order_id) REFERENCES public.storefront_orders(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.storefront_order_items
    ADD CONSTRAINT storefront_order_items_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.storefront_products(id) ON DELETE RESTRICT;
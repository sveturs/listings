-- Migration for table: storefront_inventory_movements

CREATE SEQUENCE public.storefront_inventory_movements_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.storefront_inventory_movements (
    id integer NOT NULL,
    storefront_product_id integer NOT NULL,
    variant_id integer,
    type character varying(20) NOT NULL,
    quantity integer NOT NULL,
    reason character varying(50) NOT NULL,
    order_id integer,
    notes text,
    user_id integer NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT storefront_inventory_movements_type_check CHECK (((type)::text = ANY (ARRAY[('in'::character varying)::text, ('out'::character varying)::text, ('adjustment'::character varying)::text])))
);

ALTER SEQUENCE public.storefront_inventory_movements_id_seq OWNED BY public.storefront_inventory_movements.id;

CREATE INDEX idx_storefront_inventory_movements_created_at ON public.storefront_inventory_movements USING btree (created_at);

CREATE INDEX idx_storefront_inventory_movements_product_id ON public.storefront_inventory_movements USING btree (storefront_product_id);

CREATE INDEX idx_storefront_inventory_movements_type ON public.storefront_inventory_movements USING btree (type);

CREATE INDEX idx_storefront_inventory_movements_variant_id ON public.storefront_inventory_movements USING btree (variant_id) WHERE (variant_id IS NOT NULL);

ALTER TABLE ONLY public.storefront_inventory_movements
    ADD CONSTRAINT storefront_inventory_movements_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.storefront_inventory_movements
    ADD CONSTRAINT storefront_inventory_movements_storefront_product_id_fkey FOREIGN KEY (storefront_product_id) REFERENCES public.storefront_products(id) ON DELETE CASCADE;
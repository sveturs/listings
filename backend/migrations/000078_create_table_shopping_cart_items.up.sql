-- Migration for table: shopping_cart_items

CREATE SEQUENCE public.shopping_cart_items_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.shopping_cart_items (
    id bigint NOT NULL,
    cart_id bigint NOT NULL,
    product_id bigint NOT NULL,
    variant_id bigint,
    quantity integer NOT NULL,
    price_per_unit numeric(10,2) NOT NULL,
    total_price numeric(10,2) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT shopping_cart_items_price_per_unit_check CHECK ((price_per_unit >= (0)::numeric)),
    CONSTRAINT shopping_cart_items_quantity_check CHECK ((quantity > 0)),
    CONSTRAINT shopping_cart_items_total_price_check CHECK ((total_price >= (0)::numeric))
);

ALTER SEQUENCE public.shopping_cart_items_id_seq OWNED BY public.shopping_cart_items.id;

CREATE INDEX idx_shopping_cart_items_cart_id ON public.shopping_cart_items USING btree (cart_id);

CREATE INDEX idx_shopping_cart_items_product_id ON public.shopping_cart_items USING btree (product_id);

ALTER TABLE ONLY public.shopping_cart_items
    ADD CONSTRAINT shopping_cart_items_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.shopping_cart_items
    ADD CONSTRAINT unique_cart_product_variant UNIQUE (cart_id, product_id, variant_id);

ALTER TABLE ONLY public.shopping_cart_items
    ADD CONSTRAINT shopping_cart_items_cart_id_fkey FOREIGN KEY (cart_id) REFERENCES public.shopping_carts(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.shopping_cart_items
    ADD CONSTRAINT shopping_cart_items_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.storefront_products(id) ON DELETE CASCADE;

CREATE TRIGGER trigger_shopping_cart_items_updated_at BEFORE UPDATE ON public.shopping_cart_items FOR EACH ROW EXECUTE FUNCTION public.update_shopping_cart_updated_at();
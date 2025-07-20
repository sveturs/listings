-- Migration for table: shopping_carts

CREATE SEQUENCE public.shopping_carts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.shopping_carts (
    id bigint NOT NULL,
    user_id integer,
    storefront_id integer NOT NULL,
    session_id character varying(255),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT check_cart_owner CHECK ((((user_id IS NOT NULL) AND (session_id IS NULL)) OR ((user_id IS NULL) AND (session_id IS NOT NULL))))
);

ALTER SEQUENCE public.shopping_carts_id_seq OWNED BY public.shopping_carts.id;

CREATE INDEX idx_shopping_carts_session_id ON public.shopping_carts USING btree (session_id);

CREATE INDEX idx_shopping_carts_storefront_id ON public.shopping_carts USING btree (storefront_id);

CREATE INDEX idx_shopping_carts_user_id ON public.shopping_carts USING btree (user_id);

ALTER TABLE ONLY public.shopping_carts
    ADD CONSTRAINT shopping_carts_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.shopping_carts
    ADD CONSTRAINT unique_session_storefront_cart UNIQUE (session_id, storefront_id);

ALTER TABLE ONLY public.shopping_carts
    ADD CONSTRAINT unique_user_storefront_cart UNIQUE (user_id, storefront_id);

ALTER TABLE ONLY public.shopping_carts
    ADD CONSTRAINT shopping_carts_storefront_id_fkey FOREIGN KEY (storefront_id) REFERENCES public.storefronts(id) ON DELETE CASCADE;

CREATE TRIGGER trigger_shopping_carts_updated_at BEFORE UPDATE ON public.shopping_carts FOR EACH ROW EXECUTE FUNCTION public.update_shopping_cart_updated_at();
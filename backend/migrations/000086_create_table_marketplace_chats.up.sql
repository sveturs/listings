-- Migration for table: marketplace_chats

CREATE SEQUENCE public.marketplace_chats_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.marketplace_chats (
    id integer NOT NULL,
    listing_id integer,
    buyer_id integer,
    seller_id integer,
    last_message_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    is_archived boolean DEFAULT false
);

ALTER SEQUENCE public.marketplace_chats_id_seq OWNED BY public.marketplace_chats.id;

CREATE INDEX idx_marketplace_chats_active_sorted ON public.marketplace_chats USING btree (last_message_at DESC) WHERE (NOT is_archived);

CREATE INDEX idx_marketplace_chats_archived ON public.marketplace_chats USING btree (is_archived) WHERE (NOT is_archived);

CREATE INDEX idx_marketplace_chats_buyer ON public.marketplace_chats USING btree (buyer_id);

CREATE INDEX idx_marketplace_chats_listing ON public.marketplace_chats USING btree (listing_id) WHERE (listing_id IS NOT NULL);

CREATE INDEX idx_marketplace_chats_listing_participants ON public.marketplace_chats USING btree (listing_id, buyer_id, seller_id) WHERE (listing_id IS NOT NULL);

CREATE INDEX idx_marketplace_chats_participants ON public.marketplace_chats USING btree (LEAST(buyer_id, seller_id), GREATEST(buyer_id, seller_id));

CREATE INDEX idx_marketplace_chats_seller ON public.marketplace_chats USING btree (seller_id);

CREATE INDEX idx_marketplace_chats_updated ON public.marketplace_chats USING btree (updated_at);

CREATE INDEX idx_marketplace_chats_user_lookup ON public.marketplace_chats USING btree (buyer_id, seller_id, last_message_at DESC);

CREATE UNIQUE INDEX idx_unique_direct_chat ON public.marketplace_chats USING btree (LEAST(buyer_id, seller_id), GREATEST(buyer_id, seller_id)) WHERE (listing_id IS NULL);

ALTER TABLE ONLY public.marketplace_chats
    ADD CONSTRAINT marketplace_chats_listing_id_buyer_id_seller_id_key UNIQUE (listing_id, buyer_id, seller_id);

ALTER TABLE ONLY public.marketplace_chats
    ADD CONSTRAINT marketplace_chats_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.marketplace_chats
    ADD CONSTRAINT marketplace_chats_buyer_id_fkey FOREIGN KEY (buyer_id) REFERENCES public.users(id);

ALTER TABLE ONLY public.marketplace_chats
    ADD CONSTRAINT marketplace_chats_listing_id_fkey FOREIGN KEY (listing_id) REFERENCES public.marketplace_listings(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.marketplace_chats
    ADD CONSTRAINT marketplace_chats_seller_id_fkey FOREIGN KEY (seller_id) REFERENCES public.users(id);

CREATE TRIGGER update_marketplace_chats_timestamp BEFORE UPDATE ON public.marketplace_chats FOR EACH ROW EXECUTE FUNCTION public.update_marketplace_chats_updated_at();
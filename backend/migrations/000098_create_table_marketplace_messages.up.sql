-- Migration for table: marketplace_messages

CREATE SEQUENCE public.marketplace_messages_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.marketplace_messages (
    id integer NOT NULL,
    chat_id integer,
    listing_id integer,
    sender_id integer,
    receiver_id integer,
    content text NOT NULL,
    is_read boolean DEFAULT false,
    original_language character varying(2) DEFAULT 'en'::character varying,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    has_attachments boolean DEFAULT false,
    attachments_count integer DEFAULT 0
);

ALTER SEQUENCE public.marketplace_messages_id_seq OWNED BY public.marketplace_messages.id;

CREATE INDEX idx_marketplace_messages_chat ON public.marketplace_messages USING btree (chat_id);

CREATE INDEX idx_marketplace_messages_chat_last ON public.marketplace_messages USING btree (chat_id, id DESC);

CREATE INDEX idx_marketplace_messages_chat_ordered ON public.marketplace_messages USING btree (chat_id, created_at DESC);

CREATE INDEX idx_marketplace_messages_chat_unread ON public.marketplace_messages USING btree (chat_id, receiver_id) WHERE (NOT is_read);

CREATE INDEX idx_marketplace_messages_created ON public.marketplace_messages USING btree (created_at);

CREATE INDEX idx_marketplace_messages_listing ON public.marketplace_messages USING btree (listing_id);

CREATE INDEX idx_marketplace_messages_receiver ON public.marketplace_messages USING btree (receiver_id);

CREATE INDEX idx_marketplace_messages_receiver_unread_count ON public.marketplace_messages USING btree (receiver_id, chat_id) WHERE (NOT is_read);

CREATE INDEX idx_marketplace_messages_sender ON public.marketplace_messages USING btree (sender_id);

CREATE INDEX idx_marketplace_messages_unread ON public.marketplace_messages USING btree (receiver_id, is_read) WHERE (NOT is_read);

ALTER TABLE ONLY public.marketplace_messages
    ADD CONSTRAINT marketplace_messages_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.marketplace_messages
    ADD CONSTRAINT marketplace_messages_chat_id_fkey FOREIGN KEY (chat_id) REFERENCES public.marketplace_chats(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.marketplace_messages
    ADD CONSTRAINT marketplace_messages_listing_id_fkey FOREIGN KEY (listing_id) REFERENCES public.marketplace_listings(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.marketplace_messages
    ADD CONSTRAINT marketplace_messages_receiver_id_fkey FOREIGN KEY (receiver_id) REFERENCES public.users(id);

ALTER TABLE ONLY public.marketplace_messages
    ADD CONSTRAINT marketplace_messages_sender_id_fkey FOREIGN KEY (sender_id) REFERENCES public.users(id);

CREATE TRIGGER update_marketplace_messages_timestamp BEFORE UPDATE ON public.marketplace_messages FOR EACH ROW EXECUTE FUNCTION public.update_marketplace_chats_updated_at();
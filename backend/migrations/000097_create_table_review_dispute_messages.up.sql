-- Migration for table: review_dispute_messages

CREATE SEQUENCE public.review_dispute_messages_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.review_dispute_messages (
    id integer NOT NULL,
    dispute_id integer NOT NULL,
    user_id integer NOT NULL,
    message text NOT NULL,
    attachments jsonb,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);

ALTER SEQUENCE public.review_dispute_messages_id_seq OWNED BY public.review_dispute_messages.id;

CREATE INDEX idx_dispute_messages_dispute_id ON public.review_dispute_messages USING btree (dispute_id);

ALTER TABLE ONLY public.review_dispute_messages
    ADD CONSTRAINT review_dispute_messages_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.review_dispute_messages
    ADD CONSTRAINT review_dispute_messages_dispute_id_fkey FOREIGN KEY (dispute_id) REFERENCES public.review_disputes(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.review_dispute_messages
    ADD CONSTRAINT review_dispute_messages_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);
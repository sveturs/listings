-- Migration for table: notifications

CREATE SEQUENCE public.notifications_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.notifications (
    id integer NOT NULL,
    user_id integer NOT NULL,
    type character varying(50) NOT NULL,
    title text NOT NULL,
    message text NOT NULL,
    data jsonb,
    is_read boolean DEFAULT false,
    delivered_to jsonb,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

ALTER SEQUENCE public.notifications_id_seq OWNED BY public.notifications.id;

CREATE INDEX idx_notifications_created ON public.notifications USING btree (created_at);

CREATE INDEX idx_notifications_type ON public.notifications USING btree (type);

CREATE INDEX idx_notifications_user ON public.notifications USING btree (user_id);

CREATE INDEX idx_notifications_user_unread ON public.notifications USING btree (user_id, created_at DESC) WHERE (NOT is_read);

ALTER TABLE ONLY public.notifications
    ADD CONSTRAINT notifications_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.notifications
    ADD CONSTRAINT notifications_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);
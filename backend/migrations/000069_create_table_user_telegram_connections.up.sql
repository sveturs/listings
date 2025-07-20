-- Migration for table: user_telegram_connections

CREATE TABLE public.user_telegram_connections (
    user_id integer NOT NULL,
    telegram_chat_id character varying(100) NOT NULL,
    telegram_username character varying(100),
    connected_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE ONLY public.user_telegram_connections
    ADD CONSTRAINT user_telegram_connections_pkey PRIMARY KEY (user_id);

ALTER TABLE ONLY public.user_telegram_connections
    ADD CONSTRAINT user_telegram_connections_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);
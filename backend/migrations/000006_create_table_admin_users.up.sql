-- Migration for table: admin_users

CREATE SEQUENCE public.admin_users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.admin_users (
    id integer NOT NULL,
    email character varying(255) NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    created_by integer,
    notes text
);

ALTER SEQUENCE public.admin_users_id_seq OWNED BY public.admin_users.id;

CREATE INDEX admin_users_email_idx ON public.admin_users USING btree (email);

ALTER TABLE ONLY public.admin_users
    ADD CONSTRAINT admin_users_email_key UNIQUE (email);

ALTER TABLE ONLY public.admin_users
    ADD CONSTRAINT admin_users_pkey PRIMARY KEY (id);
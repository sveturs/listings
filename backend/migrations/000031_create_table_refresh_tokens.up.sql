-- Migration for table: refresh_tokens

CREATE SEQUENCE public.refresh_tokens_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.refresh_tokens (
    id integer NOT NULL,
    user_id integer NOT NULL,
    token text NOT NULL,
    expires_at timestamp without time zone NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    user_agent text,
    ip character varying(45),
    device_name character varying(100),
    is_revoked boolean DEFAULT false,
    revoked_at timestamp without time zone
);

ALTER SEQUENCE public.refresh_tokens_id_seq OWNED BY public.refresh_tokens.id;

CREATE INDEX idx_refresh_tokens_expires_at ON public.refresh_tokens USING btree (expires_at) WHERE (NOT is_revoked);

CREATE INDEX idx_refresh_tokens_token ON public.refresh_tokens USING btree (token) WHERE (NOT is_revoked);

CREATE INDEX idx_refresh_tokens_user_id ON public.refresh_tokens USING btree (user_id) WHERE (NOT is_revoked);

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT refresh_tokens_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT refresh_tokens_token_key UNIQUE (token);
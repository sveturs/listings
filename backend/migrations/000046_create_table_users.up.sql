-- Migration for table: users

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.users (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    email character varying(150) NOT NULL,
    google_id character varying(255),
    picture_url text,
    phone character varying(20),
    bio text,
    notification_email boolean DEFAULT true,
    timezone character varying(50) DEFAULT 'UTC'::character varying,
    last_seen timestamp without time zone,
    account_status character varying(20) DEFAULT 'active'::character varying,
    settings jsonb DEFAULT '{}'::jsonb,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    city character varying(100),
    country character varying(100),
    password character varying(255),
    provider character varying(50) DEFAULT 'email'::character varying,
    CONSTRAINT users_account_status_check CHECK (((account_status)::text = ANY ((ARRAY['active'::character varying, 'inactive'::character varying, 'suspended'::character varying])::text[])))
);

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;

CREATE INDEX idx_users_active ON public.users USING btree (last_seen DESC) WHERE ((account_status)::text = 'active'::text);

CREATE INDEX idx_users_email ON public.users USING btree (email);

CREATE INDEX idx_users_email_lower ON public.users USING btree (lower((email)::text));

CREATE INDEX idx_users_phone ON public.users USING btree (phone);

CREATE INDEX idx_users_provider ON public.users USING btree (provider);

CREATE INDEX idx_users_status ON public.users USING btree (account_status);

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON public.users FOR EACH ROW EXECUTE FUNCTION public.update_user_updated_at();
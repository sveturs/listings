-- Migration for table: user_contacts

CREATE SEQUENCE public.user_contacts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.user_contacts (
    id bigint NOT NULL,
    user_id integer NOT NULL,
    contact_user_id integer NOT NULL,
    status character varying(20) DEFAULT 'pending'::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    added_from_chat_id integer,
    notes text,
    CONSTRAINT user_contacts_check CHECK ((user_id <> contact_user_id)),
    CONSTRAINT user_contacts_status_check CHECK (((status)::text = ANY ((ARRAY['pending'::character varying, 'accepted'::character varying, 'blocked'::character varying])::text[])))
);

ALTER SEQUENCE public.user_contacts_id_seq OWNED BY public.user_contacts.id;

CREATE INDEX idx_user_contacts_contact_user_id ON public.user_contacts USING btree (contact_user_id);

CREATE INDEX idx_user_contacts_created_at ON public.user_contacts USING btree (created_at);

CREATE INDEX idx_user_contacts_status ON public.user_contacts USING btree (status);

CREATE INDEX idx_user_contacts_user_id ON public.user_contacts USING btree (user_id);

ALTER TABLE ONLY public.user_contacts
    ADD CONSTRAINT user_contacts_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.user_contacts
    ADD CONSTRAINT user_contacts_user_id_contact_user_id_key UNIQUE (user_id, contact_user_id);

CREATE TRIGGER update_user_contacts_updated_at BEFORE UPDATE ON public.user_contacts FOR EACH ROW EXECUTE FUNCTION public.update_user_contacts_updated_at();
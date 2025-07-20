-- Migration for table: chat_attachments

CREATE SEQUENCE public.chat_attachments_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.chat_attachments (
    id integer NOT NULL,
    message_id integer NOT NULL,
    file_type character varying(20) NOT NULL,
    file_path character varying(500) NOT NULL,
    file_name character varying(255) NOT NULL,
    file_size bigint NOT NULL,
    content_type character varying(100) NOT NULL,
    storage_type character varying(20) DEFAULT 'minio'::character varying,
    storage_bucket character varying(100) DEFAULT 'chat-files'::character varying,
    public_url text,
    thumbnail_url text,
    metadata jsonb,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chat_attachments_file_type_check CHECK (((file_type)::text = ANY ((ARRAY['image'::character varying, 'video'::character varying, 'document'::character varying])::text[])))
);

ALTER SEQUENCE public.chat_attachments_id_seq OWNED BY public.chat_attachments.id;

CREATE INDEX idx_chat_attachments_created_at ON public.chat_attachments USING btree (created_at);

CREATE INDEX idx_chat_attachments_file_type ON public.chat_attachments USING btree (file_type);

CREATE INDEX idx_chat_attachments_message ON public.chat_attachments USING btree (message_id);

CREATE INDEX idx_chat_attachments_message_id ON public.chat_attachments USING btree (message_id);

ALTER TABLE ONLY public.chat_attachments
    ADD CONSTRAINT chat_attachments_pkey PRIMARY KEY (id);
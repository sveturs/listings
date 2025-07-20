-- Migration for table: translations

CREATE SEQUENCE public.translations_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.translations (
    id integer NOT NULL,
    entity_type character varying(50) NOT NULL,
    entity_id integer NOT NULL,
    language character varying(10) NOT NULL,
    field_name character varying(50) NOT NULL,
    translated_text text NOT NULL,
    is_machine_translated boolean DEFAULT true,
    is_verified boolean DEFAULT false,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    metadata jsonb DEFAULT '{}'::jsonb
);

ALTER SEQUENCE public.translations_id_seq OWNED BY public.translations.id;

CREATE INDEX idx_translations_lookup ON public.translations USING btree (entity_type, entity_id, language);

CREATE INDEX idx_translations_metadata ON public.translations USING gin (metadata);

CREATE INDEX idx_translations_type_lang ON public.translations USING btree (entity_type, language);

ALTER TABLE ONLY public.translations
    ADD CONSTRAINT translations_entity_type_entity_id_language_field_name_key UNIQUE (entity_type, entity_id, language, field_name);

ALTER TABLE ONLY public.translations
    ADD CONSTRAINT translations_pkey PRIMARY KEY (id);

CREATE TRIGGER update_translations_timestamp BEFORE UPDATE ON public.translations FOR EACH ROW EXECUTE FUNCTION public.update_translations_updated_at();
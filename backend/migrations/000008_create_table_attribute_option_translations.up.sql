-- Migration for table: attribute_option_translations

CREATE SEQUENCE public.attribute_option_translations_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.attribute_option_translations (
    id integer NOT NULL,
    attribute_name character varying(100) NOT NULL,
    option_value text NOT NULL,
    ru_translation text NOT NULL,
    sr_translation text NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

ALTER SEQUENCE public.attribute_option_translations_id_seq OWNED BY public.attribute_option_translations.id;

CREATE INDEX idx_attribute_option_translations ON public.attribute_option_translations USING btree (attribute_name, option_value);

ALTER TABLE ONLY public.attribute_option_translations
    ADD CONSTRAINT attribute_option_translations_attribute_name_option_value_key UNIQUE (attribute_name, option_value);

ALTER TABLE ONLY public.attribute_option_translations
    ADD CONSTRAINT attribute_option_translations_pkey PRIMARY KEY (id);

CREATE TRIGGER update_listings_on_attribute_translation_change AFTER INSERT OR UPDATE ON public.attribute_option_translations FOR EACH ROW EXECUTE FUNCTION public.trigger_update_listings_on_attribute_translation_change();
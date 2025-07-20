-- Migration for table: unit_translations

CREATE TABLE public.unit_translations (
    unit character varying(20) NOT NULL,
    language character varying(10) NOT NULL,
    translated_unit character varying(20) NOT NULL,
    display_format character varying(50) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE ONLY public.unit_translations
    ADD CONSTRAINT unit_translations_pkey PRIMARY KEY (unit, language);
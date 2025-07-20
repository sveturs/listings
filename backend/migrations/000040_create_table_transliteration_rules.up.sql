-- Migration for table: transliteration_rules

CREATE SEQUENCE public.transliteration_rules_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.transliteration_rules (
    id integer NOT NULL,
    source_char character varying(10) NOT NULL,
    target_char character varying(20) NOT NULL,
    language character varying(2) NOT NULL,
    enabled boolean DEFAULT true,
    priority integer DEFAULT 0,
    description text,
    rule_type character varying(20) DEFAULT 'custom'::character varying,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);

ALTER SEQUENCE public.transliteration_rules_id_seq OWNED BY public.transliteration_rules.id;

CREATE INDEX idx_transliteration_rules_active ON public.transliteration_rules USING btree (language, enabled, priority DESC);

CREATE INDEX idx_transliteration_rules_enabled ON public.transliteration_rules USING btree (enabled);

CREATE INDEX idx_transliteration_rules_language ON public.transliteration_rules USING btree (language);

CREATE INDEX idx_transliteration_rules_priority ON public.transliteration_rules USING btree (priority DESC);

CREATE INDEX idx_transliteration_rules_type ON public.transliteration_rules USING btree (rule_type);

ALTER TABLE ONLY public.transliteration_rules
    ADD CONSTRAINT transliteration_rules_language_source_char_key UNIQUE (language, source_char);

ALTER TABLE ONLY public.transliteration_rules
    ADD CONSTRAINT transliteration_rules_pkey PRIMARY KEY (id);

CREATE TRIGGER trigger_update_transliteration_rules_updated_at BEFORE UPDATE ON public.transliteration_rules FOR EACH ROW EXECUTE FUNCTION public.update_transliteration_rules_updated_at();
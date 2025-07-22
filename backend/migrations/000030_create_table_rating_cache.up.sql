-- Migration for table: rating_cache

CREATE TABLE public.rating_cache (
    entity_type character varying(50) NOT NULL,
    entity_id integer NOT NULL,
    average_rating numeric(3,2),
    total_reviews integer DEFAULT 0,
    distribution jsonb,
    breakdown jsonb,
    verified_percentage integer DEFAULT 0,
    recent_trend character varying(10),
    calculated_at timestamp without time zone DEFAULT now(),
    CONSTRAINT rating_cache_recent_trend_check CHECK (((recent_trend)::text = ANY (ARRAY[('up'::character varying)::text, ('down'::character varying)::text, ('stable'::character varying)::text])))
);

CREATE INDEX idx_rating_cache_type_rating ON public.rating_cache USING btree (entity_type, average_rating DESC);

CREATE INDEX idx_rating_cache_updated ON public.rating_cache USING btree (calculated_at);

ALTER TABLE ONLY public.rating_cache
    ADD CONSTRAINT rating_cache_pkey PRIMARY KEY (entity_type, entity_id);
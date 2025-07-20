-- Migration for table: gis_filter_analytics

CREATE SEQUENCE public.gis_filter_analytics_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.gis_filter_analytics (
    id integer NOT NULL,
    user_id integer,
    session_id character varying(255) NOT NULL,
    filter_type character varying(50) NOT NULL,
    filter_params jsonb NOT NULL,
    result_count integer NOT NULL,
    response_time_ms integer NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);

ALTER SEQUENCE public.gis_filter_analytics_id_seq OWNED BY public.gis_filter_analytics.id;

CREATE INDEX idx_filter_analytics_created ON public.gis_filter_analytics USING btree (created_at);

CREATE INDEX idx_filter_analytics_session ON public.gis_filter_analytics USING btree (session_id);

CREATE INDEX idx_filter_analytics_type ON public.gis_filter_analytics USING btree (filter_type);

CREATE INDEX idx_filter_analytics_user ON public.gis_filter_analytics USING btree (user_id);

ALTER TABLE ONLY public.gis_filter_analytics
    ADD CONSTRAINT gis_filter_analytics_pkey PRIMARY KEY (id);
-- Migration for table: search_statistics

CREATE SEQUENCE public.search_statistics_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.search_statistics (
    id bigint NOT NULL,
    query text NOT NULL,
    results_count integer DEFAULT 0 NOT NULL,
    search_duration_ms bigint NOT NULL,
    user_id bigint,
    search_filters jsonb,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);

ALTER SEQUENCE public.search_statistics_id_seq OWNED BY public.search_statistics.id;

ALTER TABLE ONLY public.search_statistics
    ADD CONSTRAINT search_statistics_pkey PRIMARY KEY (id);
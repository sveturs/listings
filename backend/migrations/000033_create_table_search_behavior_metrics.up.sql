-- Migration for table: search_behavior_metrics

CREATE SEQUENCE public.search_behavior_metrics_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.search_behavior_metrics (
    id bigint NOT NULL,
    search_query text NOT NULL,
    total_searches integer DEFAULT 0,
    total_clicks integer DEFAULT 0,
    ctr double precision DEFAULT 0,
    avg_click_position double precision DEFAULT 0,
    conversions integer DEFAULT 0,
    conversion_rate double precision DEFAULT 0,
    period_start timestamp with time zone NOT NULL,
    period_end timestamp with time zone NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);

ALTER SEQUENCE public.search_behavior_metrics_id_seq OWNED BY public.search_behavior_metrics.id;

CREATE INDEX idx_search_behavior_metrics_conversions ON public.search_behavior_metrics USING btree (conversions DESC);

CREATE INDEX idx_search_behavior_metrics_ctr ON public.search_behavior_metrics USING btree (ctr DESC);

CREATE INDEX idx_search_behavior_metrics_period ON public.search_behavior_metrics USING btree (period_start, period_end);

CREATE INDEX idx_search_behavior_metrics_query ON public.search_behavior_metrics USING btree (search_query);

CREATE UNIQUE INDEX idx_search_behavior_metrics_unique ON public.search_behavior_metrics USING btree (search_query, period_start);

ALTER TABLE ONLY public.search_behavior_metrics
    ADD CONSTRAINT search_behavior_metrics_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.search_behavior_metrics
    ADD CONSTRAINT search_behavior_metrics_unique UNIQUE (search_query, period_start, period_end);

CREATE TRIGGER trigger_update_search_behavior_metrics_updated_at BEFORE UPDATE ON public.search_behavior_metrics FOR EACH ROW EXECUTE FUNCTION public.update_search_behavior_metrics_updated_at();
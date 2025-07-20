-- Migration for table: item_performance_metrics

CREATE SEQUENCE public.item_performance_metrics_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.item_performance_metrics (
    id bigint NOT NULL,
    item_id character varying(50) NOT NULL,
    item_type character varying(20) NOT NULL,
    impressions integer DEFAULT 0,
    clicks integer DEFAULT 0,
    ctr double precision DEFAULT 0,
    conversions integer DEFAULT 0,
    avg_position double precision DEFAULT 0,
    period_start timestamp with time zone NOT NULL,
    period_end timestamp with time zone NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    CONSTRAINT item_performance_metrics_item_type_check CHECK (((item_type)::text = ANY ((ARRAY['marketplace'::character varying, 'storefront'::character varying])::text[])))
);

ALTER SEQUENCE public.item_performance_metrics_id_seq OWNED BY public.item_performance_metrics.id;

CREATE INDEX idx_item_performance_metrics_conversions ON public.item_performance_metrics USING btree (conversions DESC);

CREATE INDEX idx_item_performance_metrics_ctr ON public.item_performance_metrics USING btree (ctr DESC);

CREATE INDEX idx_item_performance_metrics_impressions ON public.item_performance_metrics USING btree (impressions DESC);

CREATE INDEX idx_item_performance_metrics_item ON public.item_performance_metrics USING btree (item_id, item_type);

CREATE INDEX idx_item_performance_metrics_period ON public.item_performance_metrics USING btree (period_start, period_end);

CREATE UNIQUE INDEX idx_item_performance_metrics_unique ON public.item_performance_metrics USING btree (item_id, period_start);

ALTER TABLE ONLY public.item_performance_metrics
    ADD CONSTRAINT item_performance_metrics_pkey PRIMARY KEY (id);

CREATE TRIGGER trigger_update_item_performance_metrics_updated_at BEFORE UPDATE ON public.item_performance_metrics FOR EACH ROW EXECUTE FUNCTION public.update_item_performance_metrics_updated_at();
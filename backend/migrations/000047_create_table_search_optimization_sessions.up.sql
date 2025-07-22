-- Migration for table: search_optimization_sessions

CREATE SEQUENCE public.search_optimization_sessions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.search_optimization_sessions (
    id bigint NOT NULL,
    status character varying(20) DEFAULT 'running'::character varying NOT NULL,
    start_time timestamp with time zone DEFAULT now() NOT NULL,
    end_time timestamp with time zone,
    total_fields integer DEFAULT 0 NOT NULL,
    processed_fields integer DEFAULT 0 NOT NULL,
    results jsonb,
    error_message text,
    created_by integer NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    CONSTRAINT search_optimization_sessions_status_check CHECK (((status)::text = ANY (ARRAY[('running'::character varying)::text, ('completed'::character varying)::text, ('failed'::character varying)::text, ('cancelled'::character varying)::text])))
);

ALTER SEQUENCE public.search_optimization_sessions_id_seq OWNED BY public.search_optimization_sessions.id;

CREATE INDEX idx_search_optimization_sessions_created_by ON public.search_optimization_sessions USING btree (created_by);

CREATE INDEX idx_search_optimization_sessions_start_time ON public.search_optimization_sessions USING btree (start_time);

CREATE INDEX idx_search_optimization_sessions_status ON public.search_optimization_sessions USING btree (status);

ALTER TABLE ONLY public.search_optimization_sessions
    ADD CONSTRAINT search_optimization_sessions_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.search_optimization_sessions
    ADD CONSTRAINT search_optimization_sessions_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.admin_users(id) ON DELETE RESTRICT;

CREATE TRIGGER trigger_update_search_optimization_sessions_updated_at BEFORE UPDATE ON public.search_optimization_sessions FOR EACH ROW EXECUTE FUNCTION public.update_search_optimization_sessions_updated_at();
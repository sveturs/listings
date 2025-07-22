-- Migration for table: search_weights_history

CREATE SEQUENCE public.search_weights_history_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.search_weights_history (
    id bigint NOT NULL,
    weight_id bigint NOT NULL,
    old_weight double precision NOT NULL,
    new_weight double precision NOT NULL,
    change_reason character varying(50) DEFAULT 'manual'::character varying NOT NULL,
    change_metadata jsonb DEFAULT '{}'::jsonb,
    changed_by integer,
    changed_at timestamp with time zone DEFAULT now(),
    CONSTRAINT search_weights_history_change_reason_check CHECK (((change_reason)::text = ANY (ARRAY[('manual'::character varying)::text, ('optimization'::character varying)::text, ('rollback'::character varying)::text, ('initialization'::character varying)::text])))
);

ALTER SEQUENCE public.search_weights_history_id_seq OWNED BY public.search_weights_history.id;

CREATE INDEX idx_search_weights_history_changed_at ON public.search_weights_history USING btree (changed_at);

CREATE INDEX idx_search_weights_history_changed_by ON public.search_weights_history USING btree (changed_by);

CREATE INDEX idx_search_weights_history_reason ON public.search_weights_history USING btree (change_reason);

CREATE INDEX idx_search_weights_history_weight_id ON public.search_weights_history USING btree (weight_id);

ALTER TABLE ONLY public.search_weights_history
    ADD CONSTRAINT search_weights_history_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.search_weights_history
    ADD CONSTRAINT search_weights_history_changed_by_fkey FOREIGN KEY (changed_by) REFERENCES public.admin_users(id) ON DELETE SET NULL;

ALTER TABLE ONLY public.search_weights_history
    ADD CONSTRAINT search_weights_history_weight_id_fkey FOREIGN KEY (weight_id) REFERENCES public.search_weights(id) ON DELETE CASCADE;
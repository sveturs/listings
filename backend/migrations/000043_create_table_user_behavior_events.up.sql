-- Migration for table: user_behavior_events

CREATE SEQUENCE public.user_behavior_events_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.user_behavior_events (
    id bigint NOT NULL,
    event_type character varying(50) NOT NULL,
    user_id integer,
    session_id character varying(100) NOT NULL,
    search_query text,
    item_id character varying(50),
    item_type character varying(20),
    "position" integer,
    metadata jsonb DEFAULT '{}'::jsonb,
    created_at timestamp with time zone DEFAULT now(),
    CONSTRAINT user_behavior_events_item_type_check CHECK (((item_type)::text = ANY (ARRAY[('marketplace'::character varying)::text, ('storefront'::character varying)::text, (NULL::character varying)::text])))
);

ALTER SEQUENCE public.user_behavior_events_id_seq OWNED BY public.user_behavior_events.id;

CREATE INDEX idx_user_behavior_events_created_at ON public.user_behavior_events USING btree (created_at);

CREATE INDEX idx_user_behavior_events_event_type ON public.user_behavior_events USING btree (event_type);

CREATE INDEX idx_user_behavior_events_item ON public.user_behavior_events USING btree (item_id, item_type) WHERE (item_id IS NOT NULL);

CREATE INDEX idx_user_behavior_events_search_query_type ON public.user_behavior_events USING btree (search_query, event_type) WHERE (search_query IS NOT NULL);

CREATE INDEX idx_user_behavior_events_session_id ON public.user_behavior_events USING btree (session_id);

CREATE INDEX idx_user_behavior_events_user_id ON public.user_behavior_events USING btree (user_id) WHERE (user_id IS NOT NULL);

ALTER TABLE ONLY public.user_behavior_events
    ADD CONSTRAINT user_behavior_events_pkey PRIMARY KEY (id);
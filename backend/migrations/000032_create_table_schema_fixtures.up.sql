-- Migration for table: schema_fixtures

CREATE TABLE public.schema_fixtures (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);

ALTER TABLE ONLY public.schema_fixtures
    ADD CONSTRAINT schema_fixtures_pkey PRIMARY KEY (version);
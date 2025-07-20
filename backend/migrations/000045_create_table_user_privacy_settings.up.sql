-- Migration for table: user_privacy_settings

CREATE TABLE public.user_privacy_settings (
    user_id integer NOT NULL,
    allow_contact_requests boolean DEFAULT true,
    allow_messages_from_contacts_only boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_user_privacy_settings_user_id ON public.user_privacy_settings USING btree (user_id);

ALTER TABLE ONLY public.user_privacy_settings
    ADD CONSTRAINT user_privacy_settings_pkey PRIMARY KEY (user_id);

CREATE TRIGGER update_user_privacy_settings_updated_at BEFORE UPDATE ON public.user_privacy_settings FOR EACH ROW EXECUTE FUNCTION public.update_user_privacy_settings_updated_at();
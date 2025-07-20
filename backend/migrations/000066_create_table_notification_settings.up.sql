-- Migration for table: notification_settings

CREATE TABLE public.notification_settings (
    user_id integer NOT NULL,
    notification_type character varying(50) NOT NULL,
    telegram_enabled boolean DEFAULT false,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    email_enabled boolean DEFAULT false
);

ALTER TABLE ONLY public.notification_settings
    ADD CONSTRAINT notification_settings_pkey PRIMARY KEY (user_id, notification_type);

ALTER TABLE ONLY public.notification_settings
    ADD CONSTRAINT notification_settings_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);

CREATE TRIGGER update_notification_settings_timestamp BEFORE UPDATE ON public.notification_settings FOR EACH ROW EXECUTE FUNCTION public.update_notification_settings_updated_at();
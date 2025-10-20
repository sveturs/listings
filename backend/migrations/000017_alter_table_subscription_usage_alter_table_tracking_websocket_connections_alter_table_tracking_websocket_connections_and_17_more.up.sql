ALTER TABLE ONLY public.subscription_usage
    ADD CONSTRAINT subscription_usage_subscription_id_fkey FOREIGN KEY (subscription_id) REFERENCES public.user_subscriptions(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.tracking_websocket_connections
    ADD CONSTRAINT tracking_websocket_connections_delivery_id_fkey FOREIGN KEY (delivery_id) REFERENCES public.deliveries(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.tracking_websocket_connections
    ADD CONSTRAINT tracking_websocket_connections_viber_user_id_fkey FOREIGN KEY (viber_user_id) REFERENCES public.viber_users(id);
ALTER TABLE ONLY public.translation_quality_metrics
    ADD CONSTRAINT translation_quality_metrics_translation_id_fkey FOREIGN KEY (translation_id) REFERENCES public.translations(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.translation_tasks
    ADD CONSTRAINT translation_tasks_provider_id_fkey FOREIGN KEY (provider_id) REFERENCES public.translation_providers(id);
ALTER TABLE ONLY public.unified_attribute_stats
    ADD CONSTRAINT unified_attribute_stats_attribute_id_fkey FOREIGN KEY (attribute_id) REFERENCES public.unified_attributes(id);
ALTER TABLE ONLY public.unified_attribute_values
    ADD CONSTRAINT unified_attribute_values_attribute_id_fkey FOREIGN KEY (attribute_id) REFERENCES public.unified_attributes(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.unified_category_attributes
    ADD CONSTRAINT unified_category_attributes_attribute_id_fkey FOREIGN KEY (attribute_id) REFERENCES public.unified_attributes(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.user_subscriptions
    ADD CONSTRAINT user_subscriptions_last_payment_id_fkey FOREIGN KEY (last_payment_id) REFERENCES public.payment_transactions(id);
ALTER TABLE ONLY public.user_subscriptions
    ADD CONSTRAINT user_subscriptions_plan_id_fkey FOREIGN KEY (plan_id) REFERENCES public.subscription_plans(id);
ALTER TABLE ONLY public.variant_attribute_mappings
    ADD CONSTRAINT variant_attribute_mappings_variant_attribute_id_fkey FOREIGN KEY (variant_attribute_id) REFERENCES public.unified_attributes(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.viber_messages
    ADD CONSTRAINT viber_messages_session_id_fkey FOREIGN KEY (session_id) REFERENCES public.viber_sessions(id) ON DELETE SET NULL;
ALTER TABLE ONLY public.viber_messages
    ADD CONSTRAINT viber_messages_viber_user_id_fkey FOREIGN KEY (viber_user_id) REFERENCES public.viber_users(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.viber_sessions
    ADD CONSTRAINT viber_sessions_viber_user_id_fkey FOREIGN KEY (viber_user_id) REFERENCES public.viber_users(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.viber_tracking_sessions
    ADD CONSTRAINT viber_tracking_sessions_delivery_id_fkey FOREIGN KEY (delivery_id) REFERENCES public.deliveries(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.viber_tracking_sessions
    ADD CONSTRAINT viber_tracking_sessions_viber_user_id_fkey FOREIGN KEY (viber_user_id) REFERENCES public.viber_users(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.vin_accident_history
    ADD CONSTRAINT vin_accident_history_vin_fkey FOREIGN KEY (vin) REFERENCES public.vin_decode_cache(vin) ON DELETE CASCADE;
ALTER TABLE ONLY public.vin_check_history
    ADD CONSTRAINT vin_check_history_decode_cache_id_fkey FOREIGN KEY (decode_cache_id) REFERENCES public.vin_decode_cache(id);
ALTER TABLE ONLY public.vin_ownership_history
    ADD CONSTRAINT vin_ownership_history_vin_fkey FOREIGN KEY (vin) REFERENCES public.vin_decode_cache(vin) ON DELETE CASCADE;
ALTER TABLE ONLY public.vin_recalls
    ADD CONSTRAINT vin_recalls_vin_fkey FOREIGN KEY (vin) REFERENCES public.vin_decode_cache(vin) ON DELETE CASCADE;

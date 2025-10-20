ALTER SEQUENCE public.price_history_id_seq OWNED BY public.price_history.id;
ALTER SEQUENCE public.product_variant_attribute_values_id_seq OWNED BY public.product_variant_attribute_values.id;
ALTER SEQUENCE public.product_variant_attributes_id_seq OWNED BY public.product_variant_attributes.id;
ALTER SEQUENCE public.query_cache_id_seq OWNED BY public.query_cache.id;
ALTER SEQUENCE public.review_confirmations_id_seq OWNED BY public.review_confirmations.id;
ALTER SEQUENCE public.review_disputes_id_seq OWNED BY public.review_disputes.id;
ALTER SEQUENCE public.review_responses_id_seq OWNED BY public.review_responses.id;
ALTER SEQUENCE public.reviews_id_seq OWNED BY public.reviews.id;
ALTER SEQUENCE public.role_audit_log_id_seq OWNED BY public.role_audit_log.id;
ALTER SEQUENCE public.roles_id_seq OWNED BY public.roles.id;
ALTER SEQUENCE public.saved_search_notifications_id_seq OWNED BY public.saved_search_notifications.id;
ALTER SEQUENCE public.saved_searches_id_seq OWNED BY public.saved_searches.id;
ALTER SEQUENCE public.search_behavior_metrics_id_seq OWNED BY public.search_behavior_metrics.id;
ALTER SEQUENCE public.search_config_id_seq OWNED BY public.search_config.id;
ALTER SEQUENCE public.search_optimization_sessions_id_seq OWNED BY public.search_optimization_sessions.id;
ALTER SEQUENCE public.search_queries_id_seq OWNED BY public.search_queries.id;
ALTER SEQUENCE public.search_statistics_id_seq OWNED BY public.search_statistics.id;
ALTER SEQUENCE public.search_synonyms_config_id_seq OWNED BY public.search_synonyms_config.id;
ALTER SEQUENCE public.search_synonyms_id_seq OWNED BY public.search_synonyms.id;
ALTER SEQUENCE public.search_weights_history_id_seq OWNED BY public.search_weights_history.id;
ALTER SEQUENCE public.search_weights_id_seq OWNED BY public.search_weights.id;
ALTER SEQUENCE public.shopping_cart_items_id_seq OWNED BY public.shopping_cart_items.id;
ALTER SEQUENCE public.shopping_carts_id_seq OWNED BY public.shopping_carts.id;
ALTER SEQUENCE public.storefront_import_errors_id_seq OWNED BY public.import_errors.id;
ALTER SEQUENCE public.storefront_import_jobs_id_seq OWNED BY public.import_jobs.id;
ALTER SEQUENCE public.subscription_history_id_seq OWNED BY public.subscription_history.id;
ALTER SEQUENCE public.subscription_payments_id_seq OWNED BY public.subscription_payments.id;
ALTER SEQUENCE public.subscription_plans_id_seq OWNED BY public.subscription_plans.id;
ALTER SEQUENCE public.subscription_usage_id_seq OWNED BY public.subscription_usage.id;
ALTER SEQUENCE public.tracking_websocket_connections_id_seq OWNED BY public.tracking_websocket_connections.id;
ALTER SEQUENCE public.translation_audit_log_id_seq OWNED BY public.translation_audit_log.id;
ALTER SEQUENCE public.translation_providers_id_seq OWNED BY public.translation_providers.id;
ALTER SEQUENCE public.translation_quality_metrics_id_seq OWNED BY public.translation_quality_metrics.id;
ALTER SEQUENCE public.translation_sync_conflicts_id_seq OWNED BY public.translation_sync_conflicts.id;
ALTER SEQUENCE public.translation_tasks_id_seq OWNED BY public.translation_tasks.id;
ALTER SEQUENCE public.translations_id_seq OWNED BY public.translations.id;
ALTER SEQUENCE public.transliteration_rules_id_seq OWNED BY public.transliteration_rules.id;
ALTER SEQUENCE public.unified_attribute_stats_id_seq OWNED BY public.unified_attribute_stats.id;
ALTER SEQUENCE public.unified_attribute_values_id_seq OWNED BY public.unified_attribute_values.id;
ALTER SEQUENCE public.unified_attributes_id_seq OWNED BY public.unified_attributes.id;
ALTER SEQUENCE public.unified_category_attributes_id_seq OWNED BY public.unified_category_attributes.id;
ALTER SEQUENCE public.unified_geo_id_seq OWNED BY public.unified_geo.id;
ALTER SEQUENCE public.user_behavior_events_id_seq OWNED BY public.user_behavior_events.id;
ALTER SEQUENCE public.user_car_view_history_id_seq OWNED BY public.user_car_view_history.id;
ALTER SEQUENCE public.user_contacts_id_seq OWNED BY public.user_contacts.id;
ALTER SEQUENCE public.user_notification_contacts_id_seq OWNED BY public.user_notification_contacts.id;
ALTER SEQUENCE public.user_notification_preferences_id_seq OWNED BY public.user_notification_preferences.id;
ALTER SEQUENCE public.user_subscriptions_id_seq OWNED BY public.user_subscriptions.id;
ALTER SEQUENCE public.user_view_history_id_seq OWNED BY public.user_view_history.id;
ALTER SEQUENCE public.variant_attribute_mappings_id_seq OWNED BY public.variant_attribute_mappings.id;
ALTER SEQUENCE public.viber_messages_id_seq OWNED BY public.viber_messages.id;
ALTER SEQUENCE public.viber_sessions_id_seq OWNED BY public.viber_sessions.id;
ALTER SEQUENCE public.viber_tracking_sessions_id_seq OWNED BY public.viber_tracking_sessions.id;
ALTER SEQUENCE public.viber_users_id_seq OWNED BY public.viber_users.id;
ALTER SEQUENCE public.view_statistics_id_seq OWNED BY public.view_statistics.id;
ALTER SEQUENCE public.vin_accident_history_id_seq OWNED BY public.vin_accident_history.id;
ALTER SEQUENCE public.vin_check_history_id_seq OWNED BY public.vin_check_history.id;
ALTER SEQUENCE public.vin_decode_cache_id_seq OWNED BY public.vin_decode_cache.id;
ALTER SEQUENCE public.vin_ownership_history_id_seq OWNED BY public.vin_ownership_history.id;
ALTER SEQUENCE public.vin_recalls_id_seq OWNED BY public.vin_recalls.id;
ALTER TABLE ONLY public.bex_tracking_events
    ADD CONSTRAINT bex_tracking_events_shipment_id_fkey FOREIGN KEY (shipment_id) REFERENCES public.bex_shipments(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.car_generations
    ADD CONSTRAINT car_generations_model_id_fkey FOREIGN KEY (model_id) REFERENCES public.car_models(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.car_models
    ADD CONSTRAINT car_models_make_id_fkey FOREIGN KEY (make_id) REFERENCES public.car_makes(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.courier_location_history
    ADD CONSTRAINT courier_location_history_courier_id_fkey FOREIGN KEY (courier_id) REFERENCES public.couriers(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.courier_location_history
    ADD CONSTRAINT courier_location_history_delivery_id_fkey FOREIGN KEY (delivery_id) REFERENCES public.deliveries(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.courier_zones
    ADD CONSTRAINT courier_zones_courier_id_fkey FOREIGN KEY (courier_id) REFERENCES public.couriers(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.deliveries
    ADD CONSTRAINT deliveries_courier_id_fkey FOREIGN KEY (courier_id) REFERENCES public.couriers(id);
ALTER TABLE ONLY public.delivery_notifications
    ADD CONSTRAINT delivery_notifications_shipment_id_fkey FOREIGN KEY (shipment_id) REFERENCES public.delivery_shipments(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.delivery_pricing_rules
    ADD CONSTRAINT delivery_pricing_rules_provider_id_fkey FOREIGN KEY (provider_id) REFERENCES public.delivery_providers(id);
ALTER TABLE ONLY public.delivery_shipments
    ADD CONSTRAINT delivery_shipments_provider_id_fkey FOREIGN KEY (provider_id) REFERENCES public.delivery_providers(id);
ALTER TABLE ONLY public.delivery_tracking_events
    ADD CONSTRAINT delivery_tracking_events_provider_id_fkey FOREIGN KEY (provider_id) REFERENCES public.delivery_providers(id);
ALTER TABLE ONLY public.delivery_tracking_events
    ADD CONSTRAINT delivery_tracking_events_shipment_id_fkey FOREIGN KEY (shipment_id) REFERENCES public.delivery_shipments(id);
ALTER TABLE ONLY public.escrow_payments
    ADD CONSTRAINT escrow_payments_payment_transaction_id_fkey FOREIGN KEY (payment_transaction_id) REFERENCES public.payment_transactions(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.districts
    ADD CONSTRAINT fk_districts_city_id FOREIGN KEY (city_id) REFERENCES public.cities(id) ON DELETE SET NULL;
ALTER TABLE ONLY public.post_express_shipments
    ADD CONSTRAINT fk_recipient_location FOREIGN KEY (recipient_location_id) REFERENCES public.post_express_locations(id) ON DELETE SET NULL;
ALTER TABLE ONLY public.post_express_shipments
    ADD CONSTRAINT fk_sender_location FOREIGN KEY (sender_location_id) REFERENCES public.post_express_locations(id) ON DELETE SET NULL;
ALTER TABLE ONLY public.listing_attribute_values
    ADD CONSTRAINT listing_attribute_values_attribute_id_fkey FOREIGN KEY (attribute_id) REFERENCES public.unified_attributes(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.merchant_payouts
    ADD CONSTRAINT merchant_payouts_gateway_id_fkey FOREIGN KEY (gateway_id) REFERENCES public.payment_gateways(id);
ALTER TABLE ONLY public.municipalities
    ADD CONSTRAINT municipalities_district_id_fkey FOREIGN KEY (district_id) REFERENCES public.districts(id) ON DELETE SET NULL;
ALTER TABLE ONLY public.payment_transactions
    ADD CONSTRAINT payment_transactions_gateway_id_fkey FOREIGN KEY (gateway_id) REFERENCES public.payment_gateways(id);
ALTER TABLE ONLY public.post_express_offices
    ADD CONSTRAINT post_express_offices_location_id_fkey FOREIGN KEY (location_id) REFERENCES public.post_express_locations(id);
ALTER TABLE ONLY public.post_express_tracking_events
    ADD CONSTRAINT post_express_tracking_events_shipment_id_fkey FOREIGN KEY (shipment_id) REFERENCES public.post_express_shipments(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.product_variant_attribute_values
    ADD CONSTRAINT product_variant_attribute_values_attribute_id_fkey FOREIGN KEY (attribute_id) REFERENCES public.product_variant_attributes(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.review_confirmations
    ADD CONSTRAINT review_confirmations_review_id_fkey FOREIGN KEY (review_id) REFERENCES public.reviews(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.review_disputes
    ADD CONSTRAINT review_disputes_review_id_fkey FOREIGN KEY (review_id) REFERENCES public.reviews(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.review_responses
    ADD CONSTRAINT review_responses_review_id_fkey FOREIGN KEY (review_id) REFERENCES public.reviews(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.review_votes
    ADD CONSTRAINT review_votes_review_id_fkey FOREIGN KEY (review_id) REFERENCES public.reviews(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.role_audit_log
    ADD CONSTRAINT role_audit_log_new_role_id_fkey FOREIGN KEY (new_role_id) REFERENCES public.roles(id);
ALTER TABLE ONLY public.role_audit_log
    ADD CONSTRAINT role_audit_log_old_role_id_fkey FOREIGN KEY (old_role_id) REFERENCES public.roles(id);
ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT role_permissions_permission_id_fkey FOREIGN KEY (permission_id) REFERENCES public.permissions(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT role_permissions_role_id_fkey FOREIGN KEY (role_id) REFERENCES public.roles(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.saved_search_notifications
    ADD CONSTRAINT saved_search_notifications_saved_search_id_fkey FOREIGN KEY (saved_search_id) REFERENCES public.saved_searches(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.search_weights_history
    ADD CONSTRAINT search_weights_history_weight_id_fkey FOREIGN KEY (weight_id) REFERENCES public.search_weights(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.shopping_cart_items
    ADD CONSTRAINT shopping_cart_items_cart_id_fkey FOREIGN KEY (cart_id) REFERENCES public.shopping_carts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.import_errors
    ADD CONSTRAINT storefront_import_errors_job_id_fkey FOREIGN KEY (job_id) REFERENCES public.import_jobs(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.subscription_history
    ADD CONSTRAINT subscription_history_from_plan_id_fkey FOREIGN KEY (from_plan_id) REFERENCES public.subscription_plans(id);
ALTER TABLE ONLY public.subscription_history
    ADD CONSTRAINT subscription_history_subscription_id_fkey FOREIGN KEY (subscription_id) REFERENCES public.user_subscriptions(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.subscription_history
    ADD CONSTRAINT subscription_history_to_plan_id_fkey FOREIGN KEY (to_plan_id) REFERENCES public.subscription_plans(id);
ALTER TABLE ONLY public.subscription_payments
    ADD CONSTRAINT subscription_payments_payment_id_fkey FOREIGN KEY (payment_id) REFERENCES public.payment_transactions(id);
ALTER TABLE ONLY public.subscription_payments
    ADD CONSTRAINT subscription_payments_subscription_id_fkey FOREIGN KEY (subscription_id) REFERENCES public.user_subscriptions(id) ON DELETE CASCADE;

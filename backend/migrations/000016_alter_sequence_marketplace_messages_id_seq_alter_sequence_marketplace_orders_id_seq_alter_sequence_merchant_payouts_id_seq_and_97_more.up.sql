ALTER SEQUENCE public.marketplace_messages_id_seq OWNED BY public.marketplace_messages.id;
ALTER SEQUENCE public.marketplace_orders_id_seq OWNED BY public.marketplace_orders.id;
ALTER SEQUENCE public.merchant_payouts_id_seq OWNED BY public.merchant_payouts.id;
ALTER SEQUENCE public.notification_templates_id_seq OWNED BY public.notification_templates.id;
ALTER SEQUENCE public.notifications_id_seq OWNED BY public.notifications.id;
ALTER SEQUENCE public.payment_gateways_id_seq OWNED BY public.payment_gateways.id;
ALTER SEQUENCE public.payment_methods_id_seq OWNED BY public.payment_methods.id;
ALTER SEQUENCE public.payment_transactions_id_seq OWNED BY public.payment_transactions.id;
ALTER SEQUENCE public.permissions_id_seq OWNED BY public.permissions.id;
ALTER SEQUENCE public.post_express_locations_id_seq OWNED BY public.post_express_locations.id;
ALTER SEQUENCE public.post_express_offices_id_seq OWNED BY public.post_express_offices.id;
ALTER SEQUENCE public.post_express_rates_id_seq OWNED BY public.post_express_rates.id;
ALTER SEQUENCE public.post_express_settings_id_seq OWNED BY public.post_express_settings.id;
ALTER SEQUENCE public.post_express_shipments_id_seq OWNED BY public.post_express_shipments.id;
ALTER SEQUENCE public.post_express_tracking_events_id_seq OWNED BY public.post_express_tracking_events.id;
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
ALTER SEQUENCE public.storefront_delivery_options_id_seq OWNED BY public.storefront_delivery_options.id;
ALTER SEQUENCE public.storefront_hours_id_seq OWNED BY public.storefront_hours.id;
ALTER SEQUENCE public.storefront_inventory_movements_id_seq OWNED BY public.storefront_inventory_movements.id;
ALTER SEQUENCE public.storefront_order_items_id_seq OWNED BY public.storefront_order_items.id;
ALTER SEQUENCE public.storefront_orders_id_seq OWNED BY public.storefront_orders.id;
ALTER SEQUENCE public.storefront_payment_methods_id_seq OWNED BY public.storefront_payment_methods.id;
ALTER SEQUENCE public.storefront_product_attributes_id_seq OWNED BY public.storefront_product_attributes.id;
ALTER SEQUENCE public.storefront_product_images_id_seq OWNED BY public.storefront_product_images.id;
ALTER SEQUENCE public.storefront_product_variant_images_id_seq OWNED BY public.storefront_product_variant_images.id;
ALTER SEQUENCE public.storefront_product_variants_id_seq OWNED BY public.storefront_product_variants.id;
ALTER SEQUENCE public.storefront_staff_id_seq OWNED BY public.storefront_staff.id;
ALTER SEQUENCE public.storefronts_id_seq OWNED BY public.storefronts.id;
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
ALTER SEQUENCE public.user_contacts_id_seq OWNED BY public.user_contacts.id;
ALTER SEQUENCE public.user_notification_contacts_id_seq OWNED BY public.user_notification_contacts.id;
ALTER SEQUENCE public.user_notification_preferences_id_seq OWNED BY public.user_notification_preferences.id;
ALTER SEQUENCE public.user_storefronts_id_seq OWNED BY public.user_storefronts.id;
ALTER SEQUENCE public.user_subscriptions_id_seq OWNED BY public.user_subscriptions.id;
ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;
ALTER SEQUENCE public.variant_attribute_mappings_id_seq OWNED BY public.variant_attribute_mappings.id;
ALTER SEQUENCE public.viber_messages_id_seq OWNED BY public.viber_messages.id;
ALTER SEQUENCE public.viber_sessions_id_seq OWNED BY public.viber_sessions.id;
ALTER SEQUENCE public.viber_tracking_sessions_id_seq OWNED BY public.viber_tracking_sessions.id;
ALTER SEQUENCE public.viber_users_id_seq OWNED BY public.viber_users.id;
ALTER TABLE ONLY public.bex_tracking_events
    ADD CONSTRAINT bex_tracking_events_shipment_id_fkey FOREIGN KEY (shipment_id) REFERENCES public.bex_shipments(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.car_generations
    ADD CONSTRAINT car_generations_model_id_fkey FOREIGN KEY (model_id) REFERENCES public.car_models(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.car_models
    ADD CONSTRAINT car_models_make_id_fkey FOREIGN KEY (make_id) REFERENCES public.car_makes(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.category_ai_mappings
    ADD CONSTRAINT category_ai_mappings_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.marketplace_categories(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.category_detection_cache
    ADD CONSTRAINT category_detection_cache_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.marketplace_categories(id);
ALTER TABLE ONLY public.category_detection_feedback
    ADD CONSTRAINT category_detection_feedback_correct_category_id_fkey FOREIGN KEY (correct_category_id) REFERENCES public.marketplace_categories(id);
ALTER TABLE ONLY public.category_detection_feedback
    ADD CONSTRAINT category_detection_feedback_detected_category_id_fkey FOREIGN KEY (detected_category_id) REFERENCES public.marketplace_categories(id);
ALTER TABLE ONLY public.category_detection_feedback
    ADD CONSTRAINT category_detection_feedback_listing_id_fkey FOREIGN KEY (listing_id) REFERENCES public.marketplace_listings(id) ON DELETE SET NULL;
ALTER TABLE ONLY public.category_keyword_weights
    ADD CONSTRAINT category_keyword_weights_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.marketplace_categories(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.category_keywords
    ADD CONSTRAINT category_keywords_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.marketplace_categories(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.category_variant_attributes
    ADD CONSTRAINT category_variant_attributes_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.marketplace_categories(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.component_templates
    ADD CONSTRAINT component_templates_component_id_fkey FOREIGN KEY (component_id) REFERENCES public.custom_ui_components(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.courier_location_history
    ADD CONSTRAINT courier_location_history_courier_id_fkey FOREIGN KEY (courier_id) REFERENCES public.couriers(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.courier_location_history
    ADD CONSTRAINT courier_location_history_delivery_id_fkey FOREIGN KEY (delivery_id) REFERENCES public.deliveries(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.courier_zones
    ADD CONSTRAINT courier_zones_courier_id_fkey FOREIGN KEY (courier_id) REFERENCES public.couriers(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.couriers
    ADD CONSTRAINT couriers_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);
ALTER TABLE ONLY public.custom_ui_component_usage
    ADD CONSTRAINT custom_ui_component_usage_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.marketplace_categories(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.custom_ui_component_usage
    ADD CONSTRAINT custom_ui_component_usage_component_id_fkey FOREIGN KEY (component_id) REFERENCES public.custom_ui_components(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.custom_ui_component_usage
    ADD CONSTRAINT custom_ui_component_usage_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.users(id);
ALTER TABLE ONLY public.custom_ui_component_usage
    ADD CONSTRAINT custom_ui_component_usage_updated_by_fkey FOREIGN KEY (updated_by) REFERENCES public.users(id);
ALTER TABLE ONLY public.deliveries
    ADD CONSTRAINT deliveries_courier_id_fkey FOREIGN KEY (courier_id) REFERENCES public.couriers(id);
ALTER TABLE ONLY public.deliveries
    ADD CONSTRAINT deliveries_order_id_fkey FOREIGN KEY (order_id) REFERENCES public.storefront_orders(id);
ALTER TABLE ONLY public.delivery_category_defaults
    ADD CONSTRAINT delivery_category_defaults_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.marketplace_categories(id);

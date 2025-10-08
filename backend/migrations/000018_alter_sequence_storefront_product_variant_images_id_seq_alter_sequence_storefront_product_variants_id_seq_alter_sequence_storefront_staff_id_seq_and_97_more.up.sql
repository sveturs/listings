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
ALTER SEQUENCE public.user_car_view_history_id_seq OWNED BY public.user_car_view_history.id;
ALTER SEQUENCE public.user_contacts_id_seq OWNED BY public.user_contacts.id;
ALTER SEQUENCE public.user_notification_contacts_id_seq OWNED BY public.user_notification_contacts.id;
ALTER SEQUENCE public.user_notification_preferences_id_seq OWNED BY public.user_notification_preferences.id;
ALTER SEQUENCE public.user_storefronts_id_seq OWNED BY public.user_storefronts.id;
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
ALTER TABLE ONLY public.ai_category_decisions
    ADD CONSTRAINT ai_category_decisions_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.marketplace_categories(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.ai_category_decisions
    ADD CONSTRAINT ai_category_decisions_user_corrected_category_id_fkey FOREIGN KEY (user_corrected_category_id) REFERENCES public.marketplace_categories(id) ON DELETE SET NULL;
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
ALTER TABLE ONLY public.category_proposals
    ADD CONSTRAINT category_proposals_parent_category_id_fkey FOREIGN KEY (parent_category_id) REFERENCES public.marketplace_categories(id) ON DELETE SET NULL;
ALTER TABLE ONLY public.category_proposals
    ADD CONSTRAINT category_proposals_storefront_id_fkey FOREIGN KEY (storefront_id) REFERENCES public.storefronts(id) ON DELETE SET NULL;
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
ALTER TABLE ONLY public.custom_ui_component_usage
    ADD CONSTRAINT custom_ui_component_usage_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.marketplace_categories(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.custom_ui_component_usage
    ADD CONSTRAINT custom_ui_component_usage_component_id_fkey FOREIGN KEY (component_id) REFERENCES public.custom_ui_components(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.deliveries
    ADD CONSTRAINT deliveries_courier_id_fkey FOREIGN KEY (courier_id) REFERENCES public.couriers(id);
ALTER TABLE ONLY public.deliveries
    ADD CONSTRAINT deliveries_order_id_fkey FOREIGN KEY (order_id) REFERENCES public.storefront_orders(id);
ALTER TABLE ONLY public.delivery_category_defaults
    ADD CONSTRAINT delivery_category_defaults_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.marketplace_categories(id);
ALTER TABLE ONLY public.delivery_notifications
    ADD CONSTRAINT delivery_notifications_shipment_id_fkey FOREIGN KEY (shipment_id) REFERENCES public.delivery_shipments(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.delivery_pricing_rules
    ADD CONSTRAINT delivery_pricing_rules_provider_id_fkey FOREIGN KEY (provider_id) REFERENCES public.delivery_providers(id);
ALTER TABLE ONLY public.delivery_shipments
    ADD CONSTRAINT delivery_shipments_order_id_fkey FOREIGN KEY (order_id) REFERENCES public.marketplace_orders(id);
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
ALTER TABLE ONLY public.inventory_reservations
    ADD CONSTRAINT fk_inventory_reservations_order_id FOREIGN KEY (order_id) REFERENCES public.storefront_orders(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.inventory_reservations
    ADD CONSTRAINT fk_inventory_reservations_product_id FOREIGN KEY (product_id) REFERENCES public.storefront_products(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.inventory_reservations
    ADD CONSTRAINT fk_inventory_reservations_variant_id FOREIGN KEY (variant_id) REFERENCES public.storefront_product_variants(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.marketplace_chats
    ADD CONSTRAINT fk_marketplace_chats_storefront_product FOREIGN KEY (storefront_product_id) REFERENCES public.storefront_products(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.marketplace_messages
    ADD CONSTRAINT fk_marketplace_messages_storefront_product FOREIGN KEY (storefront_product_id) REFERENCES public.storefront_products(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.post_express_shipments
    ADD CONSTRAINT fk_recipient_location FOREIGN KEY (recipient_location_id) REFERENCES public.post_express_locations(id) ON DELETE SET NULL;
ALTER TABLE ONLY public.post_express_shipments
    ADD CONSTRAINT fk_sender_location FOREIGN KEY (sender_location_id) REFERENCES public.post_express_locations(id) ON DELETE SET NULL;
ALTER TABLE ONLY public.storefront_category_mappings
    ADD CONSTRAINT fk_storefront_category_mappings_category FOREIGN KEY (target_category_id) REFERENCES public.marketplace_categories(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.storefront_category_mappings
    ADD CONSTRAINT fk_storefront_category_mappings_storefront FOREIGN KEY (storefront_id) REFERENCES public.storefronts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.import_errors
    ADD CONSTRAINT import_errors_job_id_fkey FOREIGN KEY (job_id) REFERENCES public.import_jobs(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.import_history
    ADD CONSTRAINT import_history_source_id_fkey FOREIGN KEY (source_id) REFERENCES public.import_sources(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.import_jobs
    ADD CONSTRAINT import_jobs_storefront_id_fkey FOREIGN KEY (storefront_id) REFERENCES public.storefronts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.import_sources
    ADD CONSTRAINT import_sources_storefront_id_fkey FOREIGN KEY (storefront_id) REFERENCES public.user_storefronts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.imported_categories
    ADD CONSTRAINT imported_categories_source_id_fkey FOREIGN KEY (source_id) REFERENCES public.import_sources(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.listing_attribute_values
    ADD CONSTRAINT listing_attribute_values_attribute_id_fkey FOREIGN KEY (attribute_id) REFERENCES public.unified_attributes(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.listing_attribute_values
    ADD CONSTRAINT listing_attribute_values_listing_id_fkey FOREIGN KEY (listing_id) REFERENCES public.marketplace_listings(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.listings_geo
    ADD CONSTRAINT listings_geo_district_id_fkey FOREIGN KEY (district_id) REFERENCES public.districts(id) ON DELETE SET NULL;
ALTER TABLE ONLY public.listings_geo
    ADD CONSTRAINT listings_geo_municipality_id_fkey FOREIGN KEY (municipality_id) REFERENCES public.municipalities(id) ON DELETE SET NULL;
ALTER TABLE ONLY public.marketplace_categories
    ADD CONSTRAINT marketplace_categories_parent_id_fkey FOREIGN KEY (parent_id) REFERENCES public.marketplace_categories(id);
ALTER TABLE ONLY public.marketplace_chats
    ADD CONSTRAINT marketplace_chats_listing_id_fkey FOREIGN KEY (listing_id) REFERENCES public.marketplace_listings(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.marketplace_favorites
    ADD CONSTRAINT marketplace_favorites_listing_id_fkey FOREIGN KEY (listing_id) REFERENCES public.marketplace_listings(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.marketplace_images
    ADD CONSTRAINT marketplace_images_listing_id_fkey FOREIGN KEY (listing_id) REFERENCES public.marketplace_listings(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.marketplace_listing_variants
    ADD CONSTRAINT marketplace_listing_variants_listing_id_fkey FOREIGN KEY (listing_id) REFERENCES public.marketplace_listings(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.marketplace_listings
    ADD CONSTRAINT marketplace_listings_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.marketplace_categories(id);
ALTER TABLE ONLY public.marketplace_listings
    ADD CONSTRAINT marketplace_listings_storefront_id_fkey FOREIGN KEY (storefront_id) REFERENCES public.user_storefronts(id) ON DELETE SET NULL;
ALTER TABLE ONLY public.marketplace_messages
    ADD CONSTRAINT marketplace_messages_chat_id_fkey FOREIGN KEY (chat_id) REFERENCES public.marketplace_chats(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.marketplace_messages
    ADD CONSTRAINT marketplace_messages_listing_id_fkey FOREIGN KEY (listing_id) REFERENCES public.marketplace_listings(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.marketplace_orders
    ADD CONSTRAINT marketplace_orders_delivery_shipment_id_fkey FOREIGN KEY (delivery_shipment_id) REFERENCES public.delivery_shipments(id);

ALTER TABLE ONLY public.import_jobs
    ADD CONSTRAINT storefront_import_jobs_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.subscription_history
    ADD CONSTRAINT subscription_history_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.subscription_payments
    ADD CONSTRAINT subscription_payments_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.subscription_plans
    ADD CONSTRAINT subscription_plans_code_key UNIQUE (code);
ALTER TABLE ONLY public.subscription_plans
    ADD CONSTRAINT subscription_plans_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.subscription_usage
    ADD CONSTRAINT subscription_usage_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.tracking_websocket_connections
    ADD CONSTRAINT tracking_websocket_connections_connection_id_key UNIQUE (connection_id);
ALTER TABLE ONLY public.tracking_websocket_connections
    ADD CONSTRAINT tracking_websocket_connections_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.translation_audit_log
    ADD CONSTRAINT translation_audit_log_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.translation_providers
    ADD CONSTRAINT translation_providers_name_key UNIQUE (name);
ALTER TABLE ONLY public.translation_providers
    ADD CONSTRAINT translation_providers_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.translation_quality_metrics
    ADD CONSTRAINT translation_quality_metrics_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.translation_sync_conflicts
    ADD CONSTRAINT translation_sync_conflicts_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.translation_tasks
    ADD CONSTRAINT translation_tasks_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.translations
    ADD CONSTRAINT translations_entity_type_entity_id_language_field_name_key UNIQUE (entity_type, entity_id, language, field_name);
ALTER TABLE ONLY public.translations
    ADD CONSTRAINT translations_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.transliteration_rules
    ADD CONSTRAINT transliteration_rules_language_source_char_key UNIQUE (language, source_char);
ALTER TABLE ONLY public.transliteration_rules
    ADD CONSTRAINT transliteration_rules_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.geocoding_cache
    ADD CONSTRAINT uk_geocoding_cache_normalized UNIQUE (normalized_address, language, country_code);
ALTER TABLE ONLY public.unified_geo
    ADD CONSTRAINT uk_unified_geo_source UNIQUE (source_type, source_id);
ALTER TABLE ONLY public.unified_attribute_stats
    ADD CONSTRAINT unified_attribute_stats_attribute_id_category_id_key UNIQUE (attribute_id, category_id);
ALTER TABLE ONLY public.unified_attribute_stats
    ADD CONSTRAINT unified_attribute_stats_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.unified_attribute_values
    ADD CONSTRAINT unified_attribute_values_entity_type_entity_id_attribute_id_key UNIQUE (entity_type, entity_id, attribute_id);
ALTER TABLE ONLY public.unified_attribute_values
    ADD CONSTRAINT unified_attribute_values_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.unified_attributes
    ADD CONSTRAINT unified_attributes_code_key UNIQUE (code);
ALTER TABLE ONLY public.unified_attributes
    ADD CONSTRAINT unified_attributes_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.unified_category_attributes
    ADD CONSTRAINT unified_category_attributes_category_id_attribute_id_key UNIQUE (category_id, attribute_id);
ALTER TABLE ONLY public.unified_category_attributes
    ADD CONSTRAINT unified_category_attributes_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.unified_geo
    ADD CONSTRAINT unified_geo_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.shopping_cart_items
    ADD CONSTRAINT unique_cart_product_variant UNIQUE (cart_id, product_id, variant_id);
ALTER TABLE ONLY public.category_variant_attributes
    ADD CONSTRAINT unique_category_variant_attribute UNIQUE (category_id, variant_attribute_name);
ALTER TABLE ONLY public.search_queries
    ADD CONSTRAINT unique_normalized_query_language UNIQUE (normalized_query, language);
ALTER TABLE ONLY public.shopping_carts
    ADD CONSTRAINT unique_session_storefront_cart UNIQUE (session_id, storefront_id);
ALTER TABLE ONLY public.shopping_carts
    ADD CONSTRAINT unique_user_storefront_cart UNIQUE (user_id, storefront_id);
ALTER TABLE ONLY public.post_express_rates
    ADD CONSTRAINT unique_weight_range UNIQUE (weight_from, weight_to);
ALTER TABLE ONLY public.unit_translations
    ADD CONSTRAINT unit_translations_pkey PRIMARY KEY (unit, language);
ALTER TABLE ONLY public.user_b2c_stores
    ADD CONSTRAINT user_b2c_stores_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.user_b2c_stores
    ADD CONSTRAINT user_b2c_stores_slug_key UNIQUE (slug);
ALTER TABLE ONLY public.user_balances
    ADD CONSTRAINT user_balances_pkey PRIMARY KEY (user_id);
ALTER TABLE ONLY public.user_behavior_events
    ADD CONSTRAINT user_behavior_events_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.user_car_view_history
    ADD CONSTRAINT user_car_view_history_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.user_contacts
    ADD CONSTRAINT user_contacts_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.user_contacts
    ADD CONSTRAINT user_contacts_user_id_contact_user_id_key UNIQUE (user_id, contact_user_id);
ALTER TABLE ONLY public.user_notification_contacts
    ADD CONSTRAINT user_notification_contacts_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.user_notification_contacts
    ADD CONSTRAINT user_notification_contacts_user_id_channel_contact_value_key UNIQUE (user_id, channel, contact_value);
ALTER TABLE ONLY public.user_notification_preferences
    ADD CONSTRAINT user_notification_preferences_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.user_notification_preferences
    ADD CONSTRAINT user_notification_preferences_user_id_channel_key UNIQUE (user_id, channel);
ALTER TABLE ONLY public.user_privacy_settings
    ADD CONSTRAINT user_privacy_settings_pkey PRIMARY KEY (user_id);
ALTER TABLE ONLY public.user_subscriptions
    ADD CONSTRAINT user_subscriptions_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.user_subscriptions
    ADD CONSTRAINT user_subscriptions_user_id_key UNIQUE (user_id);
ALTER TABLE ONLY public.user_telegram_connections
    ADD CONSTRAINT user_telegram_connections_pkey PRIMARY KEY (user_id);
ALTER TABLE ONLY public.user_view_history
    ADD CONSTRAINT user_view_history_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.variant_attribute_mappings
    ADD CONSTRAINT variant_attribute_mappings_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.variant_attribute_mappings
    ADD CONSTRAINT variant_attribute_mappings_variant_attribute_id_category_id_key UNIQUE (variant_attribute_id, category_id);
ALTER TABLE ONLY public.viber_messages
    ADD CONSTRAINT viber_messages_message_token_key UNIQUE (message_token);
ALTER TABLE ONLY public.viber_messages
    ADD CONSTRAINT viber_messages_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.viber_sessions
    ADD CONSTRAINT viber_sessions_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.viber_tracking_sessions
    ADD CONSTRAINT viber_tracking_sessions_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.viber_users
    ADD CONSTRAINT viber_users_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.viber_users
    ADD CONSTRAINT viber_users_viber_id_key UNIQUE (viber_id);
ALTER TABLE ONLY public.view_statistics
    ADD CONSTRAINT view_statistics_listing_id_date_key UNIQUE (listing_id, date);
ALTER TABLE ONLY public.view_statistics
    ADD CONSTRAINT view_statistics_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.vin_accident_history
    ADD CONSTRAINT vin_accident_history_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.vin_check_history
    ADD CONSTRAINT vin_check_history_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.vin_decode_cache
    ADD CONSTRAINT vin_decode_cache_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.vin_decode_cache
    ADD CONSTRAINT vin_decode_cache_vin_key UNIQUE (vin);
ALTER TABLE ONLY public.vin_ownership_history
    ADD CONSTRAINT vin_ownership_history_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.vin_recalls
    ADD CONSTRAINT vin_recalls_pkey PRIMARY KEY (id);
CREATE TRIGGER refresh_rating_summaries_trigger AFTER INSERT OR DELETE OR UPDATE ON public.reviews FOR EACH STATEMENT EXECUTE FUNCTION public.refresh_rating_summaries();
CREATE TRIGGER tg_unified_attributes_search_vector BEFORE INSERT OR UPDATE OF name, code ON public.unified_attributes FOR EACH ROW EXECUTE FUNCTION public.update_unified_attributes_search_vector();
CREATE TRIGGER tr_update_category_attribute_sort_order BEFORE INSERT ON public.category_attribute_mapping FOR EACH ROW EXECUTE FUNCTION public.update_category_attribute_sort_order();
CREATE TRIGGER trg_check_return_visit BEFORE INSERT ON public.user_view_history FOR EACH ROW EXECUTE FUNCTION public.check_return_visit();
CREATE TRIGGER trig_update_metadata_after_price_change AFTER INSERT ON public.price_history FOR EACH ROW EXECUTE FUNCTION public.update_listing_metadata_after_price_change();
CREATE TRIGGER trigger_cleanup_geocoding_cache AFTER INSERT ON public.geocoding_cache FOR EACH STATEMENT EXECUTE FUNCTION public.trigger_cleanup_geocoding_cache();
CREATE TRIGGER trigger_geocoding_cache_updated_at BEFORE UPDATE ON public.geocoding_cache FOR EACH ROW EXECUTE FUNCTION public.update_geocoding_cache_updated_at();
CREATE TRIGGER trigger_log_search_weight_changes AFTER UPDATE ON public.search_weights FOR EACH ROW EXECUTE FUNCTION public.log_search_weight_changes();
CREATE TRIGGER trigger_refresh_rating_distributions AFTER INSERT OR DELETE OR UPDATE ON public.reviews FOR EACH STATEMENT EXECUTE FUNCTION public.refresh_rating_distributions();
CREATE TRIGGER trigger_shopping_cart_items_updated_at BEFORE UPDATE ON public.shopping_cart_items FOR EACH ROW EXECUTE FUNCTION public.update_shopping_cart_updated_at();
CREATE TRIGGER trigger_shopping_carts_updated_at BEFORE UPDATE ON public.shopping_carts FOR EACH ROW EXECUTE FUNCTION public.update_shopping_cart_updated_at();
CREATE TRIGGER trigger_unified_geo_cache_refresh AFTER INSERT OR DELETE OR UPDATE ON public.unified_geo FOR EACH ROW EXECUTE FUNCTION public.trigger_refresh_map_cache();
CREATE TRIGGER trigger_update_ai_category_decisions_updated_at BEFORE UPDATE ON public.ai_category_decisions FOR EACH ROW EXECUTE FUNCTION public.update_ai_category_decisions_updated_at();
CREATE TRIGGER trigger_update_category_proposals_updated_at BEFORE UPDATE ON public.category_proposals FOR EACH ROW EXECUTE FUNCTION public.update_category_proposals_updated_at();
CREATE TRIGGER trigger_update_inventory_reservations_updated_at BEFORE UPDATE ON public.inventory_reservations FOR EACH ROW EXECUTE FUNCTION public.update_inventory_reservations_updated_at();
CREATE TRIGGER trigger_update_item_performance_metrics_updated_at BEFORE UPDATE ON public.item_performance_metrics FOR EACH ROW EXECUTE FUNCTION public.update_item_performance_metrics_updated_at();
CREATE TRIGGER trigger_update_product_variant_attribute_values_updated_at BEFORE UPDATE ON public.product_variant_attribute_values FOR EACH ROW EXECUTE FUNCTION public.update_product_variants_updated_at();
CREATE TRIGGER trigger_update_product_variant_attributes_updated_at BEFORE UPDATE ON public.product_variant_attributes FOR EACH ROW EXECUTE FUNCTION public.update_product_variants_updated_at();
CREATE TRIGGER trigger_update_search_behavior_metrics_updated_at BEFORE UPDATE ON public.search_behavior_metrics FOR EACH ROW EXECUTE FUNCTION public.update_search_behavior_metrics_updated_at();
CREATE TRIGGER trigger_update_search_optimization_sessions_updated_at BEFORE UPDATE ON public.search_optimization_sessions FOR EACH ROW EXECUTE FUNCTION public.update_search_optimization_sessions_updated_at();
CREATE TRIGGER trigger_update_search_synonyms_updated_at BEFORE UPDATE ON public.search_synonyms FOR EACH ROW EXECUTE FUNCTION public.update_search_synonyms_updated_at();
CREATE TRIGGER trigger_update_search_weights_updated_at BEFORE UPDATE ON public.search_weights FOR EACH ROW EXECUTE FUNCTION public.update_search_weights_updated_at();
CREATE TRIGGER trigger_update_transliteration_rules_updated_at BEFORE UPDATE ON public.transliteration_rules FOR EACH ROW EXECUTE FUNCTION public.update_transliteration_rules_updated_at();
CREATE TRIGGER trigger_update_unified_geo_updated_at BEFORE UPDATE ON public.unified_geo FOR EACH ROW EXECUTE FUNCTION public.update_unified_geo_updated_at();
CREATE TRIGGER update_bex_configuration_updated_at BEFORE UPDATE ON public.bex_configuration FOR EACH ROW EXECUTE FUNCTION public.update_bex_updated_at();
CREATE TRIGGER update_bex_shipments_updated_at BEFORE UPDATE ON public.bex_shipments FOR EACH ROW EXECUTE FUNCTION public.update_bex_updated_at();
CREATE TRIGGER update_car_generations_updated_at BEFORE UPDATE ON public.car_generations FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
CREATE TRIGGER update_car_makes_updated_at BEFORE UPDATE ON public.car_makes FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
CREATE TRIGGER update_car_models_updated_at BEFORE UPDATE ON public.car_models FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
CREATE TRIGGER update_category_ai_mappings_updated_at BEFORE UPDATE ON public.category_ai_mappings FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
CREATE TRIGGER update_category_detection_feedback_updated_at BEFORE UPDATE ON public.category_detection_feedback FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
CREATE TRIGGER update_category_keyword_weights_updated_at BEFORE UPDATE ON public.category_keyword_weights FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();

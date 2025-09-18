ALTER TABLE ONLY public.user_balances
    ADD CONSTRAINT user_balances_pkey PRIMARY KEY (user_id);
ALTER TABLE ONLY public.user_behavior_events
    ADD CONSTRAINT user_behavior_events_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.user_contacts
    ADD CONSTRAINT user_contacts_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.user_contacts
    ADD CONSTRAINT user_contacts_user_id_contact_user_id_key UNIQUE (user_id, contact_user_id);
ALTER TABLE ONLY public.user_privacy_settings
    ADD CONSTRAINT user_privacy_settings_pkey PRIMARY KEY (user_id);
ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_pkey PRIMARY KEY (user_id, role_id);
ALTER TABLE ONLY public.user_storefronts
    ADD CONSTRAINT user_storefronts_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.user_storefronts
    ADD CONSTRAINT user_storefronts_slug_key UNIQUE (slug);
ALTER TABLE ONLY public.user_subscriptions
    ADD CONSTRAINT user_subscriptions_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.user_subscriptions
    ADD CONSTRAINT user_subscriptions_user_id_key UNIQUE (user_id);
ALTER TABLE ONLY public.user_telegram_connections
    ADD CONSTRAINT user_telegram_connections_pkey PRIMARY KEY (user_id);
ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);
ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);
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
CREATE TRIGGER calculate_escrow_release_date_trigger BEFORE INSERT OR UPDATE ON public.storefront_orders FOR EACH ROW EXECUTE FUNCTION public.calculate_escrow_release_date();
CREATE TRIGGER marketplace_listing_variants_updated_at BEFORE UPDATE ON public.marketplace_listing_variants FOR EACH ROW EXECUTE FUNCTION public.update_marketplace_listing_variants_updated_at();
CREATE TRIGGER marketplace_orders_updated_at_trigger BEFORE UPDATE ON public.marketplace_orders FOR EACH ROW EXECUTE FUNCTION public.update_marketplace_orders_updated_at();
CREATE TRIGGER preserve_review_origin_trigger BEFORE DELETE ON public.marketplace_listings FOR EACH ROW EXECUTE FUNCTION public.preserve_review_origin();
CREATE TRIGGER refresh_category_counts_delete AFTER DELETE ON public.marketplace_listings FOR EACH ROW EXECUTE FUNCTION public.refresh_category_listing_counts();
CREATE TRIGGER refresh_category_counts_insert AFTER INSERT ON public.marketplace_listings FOR EACH ROW EXECUTE FUNCTION public.refresh_category_listing_counts();
CREATE TRIGGER refresh_category_counts_update AFTER UPDATE ON public.marketplace_listings FOR EACH ROW WHEN (((old.status)::text IS DISTINCT FROM (new.status)::text)) EXECUTE FUNCTION public.refresh_category_listing_counts();
CREATE TRIGGER refresh_rating_summaries_trigger AFTER INSERT OR DELETE OR UPDATE ON public.reviews FOR EACH STATEMENT EXECUTE FUNCTION public.refresh_rating_summaries();
CREATE TRIGGER set_order_number_trigger BEFORE INSERT ON public.storefront_orders FOR EACH ROW EXECUTE FUNCTION public.set_order_number();
CREATE TRIGGER tg_unified_attributes_search_vector BEFORE INSERT OR UPDATE OF name, code ON public.unified_attributes FOR EACH ROW EXECUTE FUNCTION public.update_unified_attributes_search_vector();
CREATE TRIGGER tr_update_category_attribute_sort_order BEFORE INSERT ON public.category_attribute_mapping FOR EACH ROW EXECUTE FUNCTION public.update_category_attribute_sort_order();
CREATE TRIGGER trg_new_listing_price_history AFTER INSERT ON public.marketplace_listings FOR EACH ROW EXECUTE FUNCTION public.update_price_history('create');
CREATE TRIGGER trg_update_listing_price_history AFTER UPDATE OF price ON public.marketplace_listings FOR EACH ROW WHEN ((old.price IS DISTINCT FROM new.price)) EXECUTE FUNCTION public.update_price_history('update');
CREATE TRIGGER trig_update_metadata_after_price_change AFTER INSERT ON public.price_history FOR EACH ROW EXECUTE FUNCTION public.update_listing_metadata_after_price_change();
CREATE TRIGGER trigger_assign_district_municipality BEFORE INSERT OR UPDATE OF location ON public.listings_geo FOR EACH ROW EXECUTE FUNCTION public.assign_district_municipality();
ALTER TABLE public.listings_geo DISABLE TRIGGER trigger_assign_district_municipality;
CREATE TRIGGER trigger_auto_geocode_storefront_product AFTER INSERT OR UPDATE ON public.storefront_products FOR EACH ROW EXECUTE FUNCTION public.auto_geocode_storefront_product();
CREATE TRIGGER trigger_cleanup_geocoding_cache AFTER INSERT ON public.geocoding_cache FOR EACH STATEMENT EXECUTE FUNCTION public.trigger_cleanup_geocoding_cache();
CREATE TRIGGER trigger_cleanup_storefront_product_geo AFTER DELETE ON public.storefront_products FOR EACH ROW EXECUTE FUNCTION public.cleanup_unified_geo();
CREATE TRIGGER trigger_geocoding_cache_updated_at BEFORE UPDATE ON public.geocoding_cache FOR EACH ROW EXECUTE FUNCTION public.update_geocoding_cache_updated_at();
CREATE TRIGGER trigger_listings_geo_updated_at BEFORE UPDATE ON public.listings_geo FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
CREATE TRIGGER trigger_log_role_change AFTER UPDATE ON public.users FOR EACH ROW EXECUTE FUNCTION public.log_role_change();
CREATE TRIGGER trigger_log_search_weight_changes AFTER UPDATE ON public.search_weights FOR EACH ROW EXECUTE FUNCTION public.log_search_weight_changes();
CREATE TRIGGER trigger_marketplace_listings_cache_refresh AFTER INSERT OR DELETE OR UPDATE ON public.marketplace_listings FOR EACH ROW EXECUTE FUNCTION public.trigger_refresh_map_cache();
CREATE TRIGGER trigger_refresh_rating_distributions AFTER INSERT OR DELETE OR UPDATE ON public.reviews FOR EACH STATEMENT EXECUTE FUNCTION public.refresh_rating_distributions();
CREATE TRIGGER trigger_shopping_cart_items_updated_at BEFORE UPDATE ON public.shopping_cart_items FOR EACH ROW EXECUTE FUNCTION public.update_shopping_cart_updated_at();
CREATE TRIGGER trigger_shopping_carts_updated_at BEFORE UPDATE ON public.shopping_carts FOR EACH ROW EXECUTE FUNCTION public.update_shopping_cart_updated_at();
CREATE TRIGGER trigger_storefront_products_cache_refresh AFTER INSERT OR DELETE OR UPDATE ON public.storefront_products FOR EACH ROW EXECUTE FUNCTION public.trigger_refresh_map_cache();
CREATE TRIGGER trigger_unified_geo_cache_refresh AFTER INSERT OR DELETE OR UPDATE ON public.unified_geo FOR EACH ROW EXECUTE FUNCTION public.trigger_refresh_map_cache();
CREATE TRIGGER trigger_update_inventory_reservations_updated_at BEFORE UPDATE ON public.inventory_reservations FOR EACH ROW EXECUTE FUNCTION public.update_inventory_reservations_updated_at();
CREATE TRIGGER trigger_update_item_performance_metrics_updated_at BEFORE UPDATE ON public.item_performance_metrics FOR EACH ROW EXECUTE FUNCTION public.update_item_performance_metrics_updated_at();
CREATE TRIGGER trigger_update_listings_geo_updated_at BEFORE UPDATE ON public.listings_geo FOR EACH ROW EXECUTE FUNCTION public.update_listings_geo_updated_at();
CREATE TRIGGER trigger_update_product_variant_attribute_values_updated_at BEFORE UPDATE ON public.product_variant_attribute_values FOR EACH ROW EXECUTE FUNCTION public.update_product_variants_updated_at();
CREATE TRIGGER trigger_update_product_variant_attributes_updated_at BEFORE UPDATE ON public.product_variant_attributes FOR EACH ROW EXECUTE FUNCTION public.update_product_variants_updated_at();
CREATE TRIGGER trigger_update_search_behavior_metrics_updated_at BEFORE UPDATE ON public.search_behavior_metrics FOR EACH ROW EXECUTE FUNCTION public.update_search_behavior_metrics_updated_at();
CREATE TRIGGER trigger_update_search_optimization_sessions_updated_at BEFORE UPDATE ON public.search_optimization_sessions FOR EACH ROW EXECUTE FUNCTION public.update_search_optimization_sessions_updated_at();
CREATE TRIGGER trigger_update_search_synonyms_updated_at BEFORE UPDATE ON public.search_synonyms FOR EACH ROW EXECUTE FUNCTION public.update_search_synonyms_updated_at();
CREATE TRIGGER trigger_update_search_weights_updated_at BEFORE UPDATE ON public.search_weights FOR EACH ROW EXECUTE FUNCTION public.update_search_weights_updated_at();
CREATE TRIGGER trigger_update_storefront_product_variants_updated_at BEFORE UPDATE ON public.storefront_product_variants FOR EACH ROW EXECUTE FUNCTION public.update_product_variants_updated_at();
CREATE TRIGGER trigger_update_storefront_products_count AFTER INSERT OR DELETE OR UPDATE ON public.storefront_products FOR EACH ROW EXECUTE FUNCTION public.update_storefront_products_count();
CREATE TRIGGER trigger_update_storefront_products_geo AFTER UPDATE ON public.storefronts FOR EACH ROW EXECUTE FUNCTION public.update_storefront_products_geo();
CREATE TRIGGER trigger_update_transliteration_rules_updated_at BEFORE UPDATE ON public.transliteration_rules FOR EACH ROW EXECUTE FUNCTION public.update_transliteration_rules_updated_at();
CREATE TRIGGER trigger_update_unified_geo_updated_at BEFORE UPDATE ON public.unified_geo FOR EACH ROW EXECUTE FUNCTION public.update_unified_geo_updated_at();
CREATE TRIGGER update_car_generations_updated_at BEFORE UPDATE ON public.car_generations FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
CREATE TRIGGER update_car_makes_updated_at BEFORE UPDATE ON public.car_makes FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
CREATE TRIGGER update_car_models_updated_at BEFORE UPDATE ON public.car_models FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
CREATE TRIGGER update_category_keywords_updated_at BEFORE UPDATE ON public.category_keywords FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
CREATE TRIGGER update_category_variant_attributes_updated_at BEFORE UPDATE ON public.category_variant_attributes FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
CREATE TRIGGER update_couriers_updated_at BEFORE UPDATE ON public.couriers FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
CREATE TRIGGER update_custom_ui_component_usage_updated_at BEFORE UPDATE ON public.custom_ui_component_usage FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
CREATE TRIGGER update_deliveries_updated_at BEFORE UPDATE ON public.deliveries FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
CREATE TRIGGER update_has_variants_trigger AFTER INSERT OR DELETE ON public.storefront_product_variants FOR EACH ROW EXECUTE FUNCTION public.update_product_has_variants();
CREATE TRIGGER update_listings_on_attribute_translation_change AFTER INSERT OR UPDATE ON public.attribute_option_translations FOR EACH ROW EXECUTE FUNCTION public.trigger_update_listings_on_attribute_translation_change();
CREATE TRIGGER update_marketplace_chats_timestamp BEFORE UPDATE ON public.marketplace_chats FOR EACH ROW EXECUTE FUNCTION public.update_marketplace_chats_updated_at();
CREATE TRIGGER update_marketplace_messages_timestamp BEFORE UPDATE ON public.marketplace_messages FOR EACH ROW EXECUTE FUNCTION public.update_marketplace_chats_updated_at();
CREATE TRIGGER update_notification_settings_timestamp BEFORE UPDATE ON public.notification_settings FOR EACH ROW EXECUTE FUNCTION public.update_notification_settings_updated_at();
CREATE TRIGGER update_ratings_after_review_change AFTER INSERT OR DELETE OR UPDATE ON public.reviews FOR EACH ROW EXECUTE FUNCTION public.refresh_rating_views();
CREATE TRIGGER update_review_responses_updated_at BEFORE UPDATE ON public.review_responses FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
CREATE TRIGGER update_reviews_updated_at BEFORE UPDATE ON public.reviews FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
CREATE TRIGGER update_search_queries_updated_at_trigger BEFORE UPDATE ON public.search_queries FOR EACH ROW EXECUTE FUNCTION public.update_search_queries_updated_at();
CREATE TRIGGER update_stock_status_trigger BEFORE INSERT OR UPDATE OF stock_quantity ON public.storefront_products FOR EACH ROW EXECUTE FUNCTION public.update_product_stock_status();
CREATE TRIGGER update_storefront_order_items_updated_at BEFORE UPDATE ON public.storefront_order_items FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
CREATE TRIGGER update_storefront_usage AFTER INSERT OR DELETE ON public.storefronts FOR EACH ROW EXECUTE FUNCTION public.update_subscription_usage();
CREATE TRIGGER update_storefront_views_count_trigger AFTER INSERT OR DELETE OR UPDATE OF view_count, storefront_id ON public.storefront_products FOR EACH ROW EXECUTE FUNCTION public.update_storefront_views_count();
CREATE TRIGGER update_translation_providers_timestamp BEFORE UPDATE ON public.translation_providers FOR EACH ROW EXECUTE FUNCTION public.update_translation_providers_updated_at();
CREATE TRIGGER update_translations_timestamp BEFORE UPDATE ON public.translations FOR EACH ROW EXECUTE FUNCTION public.update_translations_updated_at();
CREATE TRIGGER update_unified_attribute_values_updated_at BEFORE UPDATE ON public.unified_attribute_values FOR EACH ROW EXECUTE FUNCTION public.update_unified_attributes_updated_at();
CREATE TRIGGER update_unified_attributes_updated_at BEFORE UPDATE ON public.unified_attributes FOR EACH ROW EXECUTE FUNCTION public.update_unified_attributes_updated_at();
CREATE TRIGGER update_unified_category_attributes_updated_at BEFORE UPDATE ON public.unified_category_attributes FOR EACH ROW EXECUTE FUNCTION public.update_unified_attributes_updated_at();
CREATE TRIGGER update_user_contacts_updated_at BEFORE UPDATE ON public.user_contacts FOR EACH ROW EXECUTE FUNCTION public.update_user_contacts_updated_at();
CREATE TRIGGER update_user_privacy_settings_updated_at BEFORE UPDATE ON public.user_privacy_settings FOR EACH ROW EXECUTE FUNCTION public.update_user_privacy_settings_updated_at();
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON public.users FOR EACH ROW EXECUTE FUNCTION public.update_user_updated_at();
CREATE TRIGGER update_viber_users_updated_at BEFORE UPDATE ON public.viber_users FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
ALTER SEQUENCE public.address_change_log_id_seq OWNED BY public.address_change_log.id;
ALTER SEQUENCE public.admin_users_id_seq OWNED BY public.admin_users.id;
ALTER SEQUENCE public.attribute_group_items_id_seq OWNED BY public.attribute_group_items.id;
ALTER SEQUENCE public.attribute_groups_id_seq OWNED BY public.attribute_groups.id;
ALTER SEQUENCE public.attribute_option_translations_id_seq OWNED BY public.attribute_option_translations.id;
ALTER SEQUENCE public.balance_transactions_id_seq OWNED BY public.balance_transactions.id;

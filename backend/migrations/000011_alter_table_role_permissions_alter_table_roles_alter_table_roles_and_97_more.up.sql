ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT role_permissions_pkey PRIMARY KEY (role_id, permission_id);
ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_name_key UNIQUE (name);
ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.search_behavior_metrics
    ADD CONSTRAINT search_behavior_metrics_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.search_behavior_metrics
    ADD CONSTRAINT search_behavior_metrics_unique UNIQUE (search_query, period_start, period_end);
ALTER TABLE ONLY public.search_config
    ADD CONSTRAINT search_config_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.search_optimization_sessions
    ADD CONSTRAINT search_optimization_sessions_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.search_queries
    ADD CONSTRAINT search_queries_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.search_statistics
    ADD CONSTRAINT search_statistics_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.search_synonyms_config
    ADD CONSTRAINT search_synonyms_config_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.search_synonyms_config
    ADD CONSTRAINT search_synonyms_config_term_language_key UNIQUE (term, language);
ALTER TABLE ONLY public.search_synonyms
    ADD CONSTRAINT search_synonyms_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.search_weights_history
    ADD CONSTRAINT search_weights_history_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.search_weights
    ADD CONSTRAINT search_weights_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.shopping_cart_items
    ADD CONSTRAINT shopping_cart_items_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.shopping_carts
    ADD CONSTRAINT shopping_carts_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.storefront_delivery_options
    ADD CONSTRAINT storefront_delivery_options_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.storefront_hours
    ADD CONSTRAINT storefront_hours_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.storefront_hours
    ADD CONSTRAINT storefront_hours_storefront_id_day_of_week_special_date_key UNIQUE (storefront_id, day_of_week, special_date);
ALTER TABLE ONLY public.storefront_inventory_movements
    ADD CONSTRAINT storefront_inventory_movements_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.storefront_order_items
    ADD CONSTRAINT storefront_order_items_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.storefront_orders
    ADD CONSTRAINT storefront_orders_order_number_key UNIQUE (order_number);
ALTER TABLE ONLY public.storefront_orders
    ADD CONSTRAINT storefront_orders_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.storefront_payment_methods
    ADD CONSTRAINT storefront_payment_methods_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.storefront_product_attributes
    ADD CONSTRAINT storefront_product_attributes_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.storefront_product_attributes
    ADD CONSTRAINT storefront_product_attributes_product_id_attribute_id_key UNIQUE (product_id, attribute_id);
ALTER TABLE ONLY public.storefront_product_images
    ADD CONSTRAINT storefront_product_images_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.storefront_product_variant_images
    ADD CONSTRAINT storefront_product_variant_images_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.storefront_product_variants
    ADD CONSTRAINT storefront_product_variants_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.storefront_product_variants
    ADD CONSTRAINT storefront_product_variants_sku_key UNIQUE (sku);
ALTER TABLE ONLY public.storefront_products
    ADD CONSTRAINT storefront_products_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.storefront_staff
    ADD CONSTRAINT storefront_staff_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.storefront_staff
    ADD CONSTRAINT storefront_staff_storefront_id_user_id_key UNIQUE (storefront_id, user_id);
ALTER TABLE ONLY public.storefronts
    ADD CONSTRAINT storefronts_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.storefronts
    ADD CONSTRAINT storefronts_slug_key UNIQUE (slug);
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
ALTER TABLE ONLY public.listings_geo
    ADD CONSTRAINT uk_listings_geo_listing_id UNIQUE (listing_id);
ALTER TABLE ONLY public.marketplace_listing_variants
    ADD CONSTRAINT uk_marketplace_listing_variants_sku_per_listing UNIQUE (listing_id, sku);
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
ALTER TABLE ONLY public.marketplace_chats
    ADD CONSTRAINT unique_storefront_product_chat UNIQUE (storefront_product_id, buyer_id, seller_id);
ALTER TABLE ONLY public.shopping_carts
    ADD CONSTRAINT unique_user_storefront_cart UNIQUE (user_id, storefront_id);
ALTER TABLE ONLY public.unit_translations
    ADD CONSTRAINT unit_translations_pkey PRIMARY KEY (unit, language);
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

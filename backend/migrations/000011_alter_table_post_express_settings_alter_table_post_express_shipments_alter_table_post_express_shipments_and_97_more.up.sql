ALTER TABLE ONLY public.post_express_settings
    ADD CONSTRAINT post_express_settings_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.post_express_shipments
    ADD CONSTRAINT post_express_shipments_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.post_express_shipments
    ADD CONSTRAINT post_express_shipments_tracking_number_key UNIQUE (tracking_number);
ALTER TABLE ONLY public.post_express_tracking_events
    ADD CONSTRAINT post_express_tracking_events_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.price_history
    ADD CONSTRAINT price_history_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.product_variant_attribute_values
    ADD CONSTRAINT product_variant_attribute_values_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.product_variant_attributes
    ADD CONSTRAINT product_variant_attributes_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.rating_cache
    ADD CONSTRAINT rating_cache_pkey PRIMARY KEY (entity_type, entity_id);
ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT refresh_tokens_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT refresh_tokens_token_key UNIQUE (token);
ALTER TABLE ONLY public.review_confirmations
    ADD CONSTRAINT review_confirmations_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.review_confirmations
    ADD CONSTRAINT review_confirmations_review_id_key UNIQUE (review_id);
ALTER TABLE ONLY public.review_disputes
    ADD CONSTRAINT review_disputes_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.review_responses
    ADD CONSTRAINT review_responses_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.review_votes
    ADD CONSTRAINT review_votes_pkey PRIMARY KEY (review_id, user_id);
ALTER TABLE ONLY public.reviews
    ADD CONSTRAINT reviews_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.role_audit_log
    ADD CONSTRAINT role_audit_log_pkey PRIMARY KEY (id);
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
ALTER TABLE ONLY public.storefront_fbs_settings
    ADD CONSTRAINT storefront_fbs_settings_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.storefront_fbs_settings
    ADD CONSTRAINT storefront_fbs_settings_storefront_id_key UNIQUE (storefront_id);
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
ALTER TABLE ONLY public.unified_geo
    ADD CONSTRAINT uk_unified_geo_source UNIQUE (source_type, source_id);
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
ALTER TABLE ONLY public.variant_attribute_mappings
    ADD CONSTRAINT unique_variant_category_mapping UNIQUE (variant_attribute_id, category_attribute_id);
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
ALTER TABLE ONLY public.user_telegram_connections
    ADD CONSTRAINT user_telegram_connections_pkey PRIMARY KEY (user_id);
ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);
ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.variant_attribute_mappings
    ADD CONSTRAINT variant_attribute_mappings_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.warehouse_inventory
    ADD CONSTRAINT warehouse_inventory_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.warehouse_invoices
    ADD CONSTRAINT warehouse_invoices_invoice_number_key UNIQUE (invoice_number);
ALTER TABLE ONLY public.warehouse_invoices
    ADD CONSTRAINT warehouse_invoices_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.warehouse_movements
    ADD CONSTRAINT warehouse_movements_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.warehouse_pickup_orders
    ADD CONSTRAINT warehouse_pickup_orders_pickup_code_key UNIQUE (pickup_code);
ALTER TABLE ONLY public.warehouse_pickup_orders
    ADD CONSTRAINT warehouse_pickup_orders_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.warehouses
    ADD CONSTRAINT warehouses_code_key UNIQUE (code);
ALTER TABLE ONLY public.warehouses
    ADD CONSTRAINT warehouses_pkey PRIMARY KEY (id);
CREATE TRIGGER calculate_escrow_release_date_trigger BEFORE INSERT OR UPDATE ON public.storefront_orders FOR EACH ROW EXECUTE FUNCTION public.calculate_escrow_release_date();
CREATE TRIGGER log_post_express_shipment_status BEFORE UPDATE ON public.post_express_shipments FOR EACH ROW EXECUTE FUNCTION public.log_shipment_status_change();
CREATE TRIGGER marketplace_orders_updated_at_trigger BEFORE UPDATE ON public.marketplace_orders FOR EACH ROW EXECUTE FUNCTION public.update_marketplace_orders_updated_at();
CREATE TRIGGER preserve_review_origin_trigger BEFORE DELETE ON public.marketplace_listings FOR EACH ROW EXECUTE FUNCTION public.preserve_review_origin();

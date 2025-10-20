ALTER TABLE ONLY public.car_models
    ADD CONSTRAINT car_models_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.category_ai_mappings
    ADD CONSTRAINT category_ai_mappings_ai_domain_product_type_category_id_key UNIQUE (ai_domain, product_type, category_id);
ALTER TABLE ONLY public.category_ai_mappings
    ADD CONSTRAINT category_ai_mappings_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.category_attribute_groups
    ADD CONSTRAINT category_attribute_groups_category_id_group_id_key UNIQUE (category_id, group_id);
ALTER TABLE ONLY public.category_attribute_groups
    ADD CONSTRAINT category_attribute_groups_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.category_attribute_mapping
    ADD CONSTRAINT category_attribute_mapping_pkey PRIMARY KEY (category_id, attribute_id);
ALTER TABLE ONLY public.category_detection_cache
    ADD CONSTRAINT category_detection_cache_cache_key_key UNIQUE (cache_key);
ALTER TABLE ONLY public.category_detection_cache
    ADD CONSTRAINT category_detection_cache_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.category_detection_feedback
    ADD CONSTRAINT category_detection_feedback_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.category_keyword_weights
    ADD CONSTRAINT category_keyword_weights_keyword_category_id_language_key UNIQUE (keyword, category_id, language);
ALTER TABLE ONLY public.category_keyword_weights
    ADD CONSTRAINT category_keyword_weights_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.category_keywords
    ADD CONSTRAINT category_keywords_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.category_proposals
    ADD CONSTRAINT category_proposals_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.category_variant_attributes
    ADD CONSTRAINT category_variant_attributes_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.chat_attachments
    ADD CONSTRAINT chat_attachments_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.cities
    ADD CONSTRAINT cities_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.cities
    ADD CONSTRAINT cities_slug_key UNIQUE (slug);
ALTER TABLE ONLY public.component_templates
    ADD CONSTRAINT component_templates_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.courier_location_history
    ADD CONSTRAINT courier_location_history_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.courier_zones
    ADD CONSTRAINT courier_zones_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.couriers
    ADD CONSTRAINT couriers_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.deliveries
    ADD CONSTRAINT deliveries_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.deliveries
    ADD CONSTRAINT deliveries_tracking_token_key UNIQUE (tracking_token);
ALTER TABLE ONLY public.delivery_category_defaults
    ADD CONSTRAINT delivery_category_defaults_category_id_key UNIQUE (category_id);
ALTER TABLE ONLY public.delivery_category_defaults
    ADD CONSTRAINT delivery_category_defaults_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.delivery_notifications
    ADD CONSTRAINT delivery_notifications_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.delivery_pricing_rules
    ADD CONSTRAINT delivery_pricing_rules_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.delivery_providers
    ADD CONSTRAINT delivery_providers_code_key UNIQUE (code);
ALTER TABLE ONLY public.delivery_providers
    ADD CONSTRAINT delivery_providers_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.delivery_shipments
    ADD CONSTRAINT delivery_shipments_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.delivery_shipments
    ADD CONSTRAINT delivery_shipments_tracking_number_key UNIQUE (tracking_number);
ALTER TABLE ONLY public.delivery_tracking_events
    ADD CONSTRAINT delivery_tracking_events_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.delivery_zones
    ADD CONSTRAINT delivery_zones_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.districts
    ADD CONSTRAINT districts_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.escrow_payments
    ADD CONSTRAINT escrow_payments_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.geocoding_cache
    ADD CONSTRAINT geocoding_cache_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.import_history
    ADD CONSTRAINT import_history_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.imported_categories
    ADD CONSTRAINT imported_categories_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.imported_categories
    ADD CONSTRAINT imported_categories_source_id_source_category_key UNIQUE (source_id, source_category);
ALTER TABLE ONLY public.inventory_reservations
    ADD CONSTRAINT inventory_reservations_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.item_performance_metrics
    ADD CONSTRAINT item_performance_metrics_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.listing_attribute_values
    ADD CONSTRAINT listing_attribute_values_listing_id_attribute_id_key UNIQUE (listing_id, attribute_id);
ALTER TABLE ONLY public.listing_attribute_values
    ADD CONSTRAINT listing_attribute_values_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.listing_views
    ADD CONSTRAINT listing_view_uniqueness UNIQUE (listing_id, user_id);
ALTER TABLE ONLY public.listing_views
    ADD CONSTRAINT listing_views_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.map_items_cache
    ADD CONSTRAINT map_items_cache_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.merchant_payouts
    ADD CONSTRAINT merchant_payouts_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.municipalities
    ADD CONSTRAINT municipalities_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.notification_settings
    ADD CONSTRAINT notification_settings_pkey PRIMARY KEY (user_id, notification_type);
ALTER TABLE ONLY public.notification_templates
    ADD CONSTRAINT notification_templates_code_key UNIQUE (code);
ALTER TABLE ONLY public.notification_templates
    ADD CONSTRAINT notification_templates_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.notifications
    ADD CONSTRAINT notifications_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.payment_gateways
    ADD CONSTRAINT payment_gateways_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.payment_methods
    ADD CONSTRAINT payment_methods_code_key UNIQUE (code);
ALTER TABLE ONLY public.payment_methods
    ADD CONSTRAINT payment_methods_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.payment_transactions
    ADD CONSTRAINT payment_transactions_order_reference_key UNIQUE (order_reference);
ALTER TABLE ONLY public.payment_transactions
    ADD CONSTRAINT payment_transactions_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.permissions
    ADD CONSTRAINT permissions_name_key UNIQUE (name);
ALTER TABLE ONLY public.permissions
    ADD CONSTRAINT permissions_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.post_express_locations
    ADD CONSTRAINT post_express_locations_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.post_express_locations
    ADD CONSTRAINT post_express_locations_post_express_id_key UNIQUE (post_express_id);
ALTER TABLE ONLY public.post_express_offices
    ADD CONSTRAINT post_express_offices_office_code_key UNIQUE (office_code);
ALTER TABLE ONLY public.post_express_offices
    ADD CONSTRAINT post_express_offices_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.post_express_rates
    ADD CONSTRAINT post_express_rates_pkey PRIMARY KEY (id);
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
ALTER TABLE ONLY public.query_cache
    ADD CONSTRAINT query_cache_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.query_cache
    ADD CONSTRAINT query_cache_query_hash_key UNIQUE (query_hash);
ALTER TABLE ONLY public.rating_cache
    ADD CONSTRAINT rating_cache_pkey PRIMARY KEY (entity_type, entity_id);
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
ALTER TABLE ONLY public.saved_search_notifications
    ADD CONSTRAINT saved_search_notifications_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.saved_searches
    ADD CONSTRAINT saved_searches_pkey PRIMARY KEY (id);
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
ALTER TABLE ONLY public.import_errors
    ADD CONSTRAINT storefront_import_errors_pkey PRIMARY KEY (id);

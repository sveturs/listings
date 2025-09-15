SELECT pg_catalog.setval('public.post_express_locations_id_seq', 5, true);
SELECT pg_catalog.setval('public.post_express_offices_id_seq', 1, false);
SELECT pg_catalog.setval('public.post_express_rates_id_seq', 5, true);
SELECT pg_catalog.setval('public.post_express_settings_id_seq', 1, true);
SELECT pg_catalog.setval('public.post_express_shipments_id_seq', 1, true);
SELECT pg_catalog.setval('public.post_express_tracking_events_id_seq', 1, false);
SELECT pg_catalog.setval('public.price_history_id_seq', 180, true);
SELECT pg_catalog.setval('public.product_variant_attribute_values_id_seq', 37, true);
SELECT pg_catalog.setval('public.product_variant_attributes_id_seq', 36, true);
SELECT pg_catalog.setval('public.query_cache_id_seq', 1, false);
SELECT pg_catalog.setval('public.review_confirmations_id_seq', 1, false);
SELECT pg_catalog.setval('public.review_disputes_id_seq', 1, false);
SELECT pg_catalog.setval('public.review_responses_id_seq', 1, false);
SELECT pg_catalog.setval('public.reviews_id_seq', 24, true);
SELECT pg_catalog.setval('public.role_audit_log_id_seq', 3, true);
SELECT pg_catalog.setval('public.roles_id_seq', 30, true);
SELECT pg_catalog.setval('public.search_behavior_metrics_id_seq', 1, false);
SELECT pg_catalog.setval('public.search_config_id_seq', 2, false);
SELECT pg_catalog.setval('public.search_optimization_sessions_id_seq', 1, false);
SELECT pg_catalog.setval('public.search_queries_id_seq', 327, true);
SELECT pg_catalog.setval('public.search_statistics_id_seq', 1, false);
SELECT pg_catalog.setval('public.search_synonyms_config_id_seq', 8, false);
SELECT pg_catalog.setval('public.search_synonyms_id_seq', 41, false);
SELECT pg_catalog.setval('public.search_weights_history_id_seq', 1, false);
SELECT pg_catalog.setval('public.search_weights_id_seq', 17, false);
SELECT pg_catalog.setval('public.shopping_cart_items_id_seq', 64, true);
SELECT pg_catalog.setval('public.shopping_carts_id_seq', 21, true);
SELECT pg_catalog.setval('public.storefront_delivery_options_id_seq', 1, true);
SELECT pg_catalog.setval('public.storefront_hours_id_seq', 55, true);
SELECT pg_catalog.setval('public.storefront_inventory_movements_id_seq', 1, true);
SELECT pg_catalog.setval('public.storefront_order_items_id_seq', 19, true);
SELECT pg_catalog.setval('public.storefront_orders_id_seq', 58, true);
SELECT pg_catalog.setval('public.storefront_payment_methods_id_seq', 33, true);
SELECT pg_catalog.setval('public.storefront_product_attributes_id_seq', 1, true);
SELECT pg_catalog.setval('public.storefront_product_images_id_seq', 70, true);
SELECT pg_catalog.setval('public.storefront_product_variant_images_id_seq', 1, false);
SELECT pg_catalog.setval('public.storefront_product_variants_id_seq', 97, true);
SELECT pg_catalog.setval('public.storefront_staff_id_seq', 19, true);
SELECT pg_catalog.setval('public.storefronts_id_seq', 39, true);
SELECT pg_catalog.setval('public.subscription_history_id_seq', 1, false);
SELECT pg_catalog.setval('public.subscription_payments_id_seq', 1, false);
SELECT pg_catalog.setval('public.subscription_plans_id_seq', 4, true);
SELECT pg_catalog.setval('public.subscription_usage_id_seq', 1, false);
SELECT pg_catalog.setval('public.translation_audit_log_id_seq', 2, true);
SELECT pg_catalog.setval('public.translation_providers_id_seq', 4, true);
SELECT pg_catalog.setval('public.translation_quality_metrics_id_seq', 1, false);
SELECT pg_catalog.setval('public.translation_sync_conflicts_id_seq', 1, false);
SELECT pg_catalog.setval('public.translation_tasks_id_seq', 1, false);
SELECT pg_catalog.setval('public.translations_id_seq', 4999, true);
SELECT pg_catalog.setval('public.transliteration_rules_id_seq', 52, false);
SELECT pg_catalog.setval('public.unified_attribute_stats_id_seq', 1, false);
SELECT pg_catalog.setval('public.unified_attribute_values_id_seq', 31, true);
SELECT pg_catalog.setval('public.unified_attributes_id_seq', 331, true);
SELECT pg_catalog.setval('public.unified_category_attributes_id_seq', 2306, true);
SELECT pg_catalog.setval('public.unified_geo_id_seq', 475, true);
SELECT pg_catalog.setval('public.user_behavior_events_id_seq', 2682, true);
SELECT pg_catalog.setval('public.user_contacts_id_seq', 25, true);
SELECT pg_catalog.setval('public.user_storefronts_id_seq', 1, false);
SELECT pg_catalog.setval('public.user_subscriptions_id_seq', 1, false);
SELECT pg_catalog.setval('public.users_id_seq', 11, true);
SELECT pg_catalog.setval('public.variant_attribute_mappings_id_seq', 1, false);
ALTER TABLE ONLY public.address_change_log
    ADD CONSTRAINT address_change_log_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.admin_users
    ADD CONSTRAINT admin_users_email_key UNIQUE (email);
ALTER TABLE ONLY public.admin_users
    ADD CONSTRAINT admin_users_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.attribute_group_items
    ADD CONSTRAINT attribute_group_items_group_id_attribute_id_key UNIQUE (group_id, attribute_id);
ALTER TABLE ONLY public.attribute_group_items
    ADD CONSTRAINT attribute_group_items_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.attribute_groups
    ADD CONSTRAINT attribute_groups_code_key UNIQUE (code);
ALTER TABLE ONLY public.attribute_groups
    ADD CONSTRAINT attribute_groups_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.attribute_option_translations
    ADD CONSTRAINT attribute_option_translations_attribute_name_option_value_key UNIQUE (attribute_name, option_value);
ALTER TABLE ONLY public.attribute_option_translations
    ADD CONSTRAINT attribute_option_translations_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.balance_transactions
    ADD CONSTRAINT balance_transactions_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.car_generations
    ADD CONSTRAINT car_generations_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.car_makes
    ADD CONSTRAINT car_makes_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.car_makes
    ADD CONSTRAINT car_makes_slug_key UNIQUE (slug);
ALTER TABLE ONLY public.car_market_analysis
    ADD CONSTRAINT car_market_analysis_brand_model_key UNIQUE (brand, model);
ALTER TABLE ONLY public.car_market_analysis
    ADD CONSTRAINT car_market_analysis_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.car_models
    ADD CONSTRAINT car_models_make_id_slug_key UNIQUE (make_id, slug);
ALTER TABLE ONLY public.car_models
    ADD CONSTRAINT car_models_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.category_attribute_groups
    ADD CONSTRAINT category_attribute_groups_category_id_group_id_key UNIQUE (category_id, group_id);
ALTER TABLE ONLY public.category_attribute_groups
    ADD CONSTRAINT category_attribute_groups_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.category_attribute_mapping
    ADD CONSTRAINT category_attribute_mapping_pkey PRIMARY KEY (category_id, attribute_id);
ALTER TABLE ONLY public.category_keywords
    ADD CONSTRAINT category_keywords_pkey PRIMARY KEY (id);
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
ALTER TABLE ONLY public.custom_ui_component_usage
    ADD CONSTRAINT custom_ui_component_usage_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.custom_ui_components
    ADD CONSTRAINT custom_ui_components_name_key UNIQUE (name);
ALTER TABLE ONLY public.custom_ui_components
    ADD CONSTRAINT custom_ui_components_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.custom_ui_templates
    ADD CONSTRAINT custom_ui_templates_name_key UNIQUE (name);
ALTER TABLE ONLY public.custom_ui_templates
    ADD CONSTRAINT custom_ui_templates_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.districts
    ADD CONSTRAINT districts_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.escrow_payments
    ADD CONSTRAINT escrow_payments_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.geocoding_cache
    ADD CONSTRAINT geocoding_cache_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.gis_filter_analytics
    ADD CONSTRAINT gis_filter_analytics_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.gis_isochrone_cache
    ADD CONSTRAINT gis_isochrone_cache_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.gis_poi_cache
    ADD CONSTRAINT gis_poi_cache_external_id_key UNIQUE (external_id);
ALTER TABLE ONLY public.gis_poi_cache
    ADD CONSTRAINT gis_poi_cache_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.import_history
    ADD CONSTRAINT import_history_pkey PRIMARY KEY (id);

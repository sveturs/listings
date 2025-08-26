SELECT pg_catalog.setval('public.custom_ui_components_id_seq', 2, false);
SELECT pg_catalog.setval('public.custom_ui_templates_id_seq', 1, false);
SELECT pg_catalog.setval('public.escrow_payments_id_seq', 1, false);
SELECT pg_catalog.setval('public.geocoding_cache_id_seq', 12, true);
SELECT pg_catalog.setval('public.gis_filter_analytics_id_seq', 1, false);
SELECT pg_catalog.setval('public.gis_isochrone_cache_id_seq', 1, false);
SELECT pg_catalog.setval('public.gis_poi_cache_id_seq', 1, false);
SELECT pg_catalog.setval('public.global_product_id_seq', 136, true);
SELECT pg_catalog.setval('public.import_history_id_seq', 1, false);
SELECT pg_catalog.setval('public.import_sources_id_seq', 1, false);
SELECT pg_catalog.setval('public.imported_categories_id_seq', 1, false);
SELECT pg_catalog.setval('public.inventory_reservations_id_seq', 20, true);
SELECT pg_catalog.setval('public.item_performance_metrics_id_seq', 1, false);
SELECT pg_catalog.setval('public.listing_attribute_values_id_seq', 177, true);
SELECT pg_catalog.setval('public.listing_views_id_seq', 2, true);
SELECT pg_catalog.setval('public.listings_geo_id_seq', 1, false);
SELECT pg_catalog.setval('public.map_items_cache_id_seq', 1, false);
SELECT pg_catalog.setval('public.marketplace_categories_id_seq', 10001, true);
SELECT pg_catalog.setval('public.marketplace_chats_id_seq', 15, true);
SELECT pg_catalog.setval('public.marketplace_images_id_seq', 66, true);
SELECT pg_catalog.setval('public.marketplace_messages_id_seq', 34, true);
SELECT pg_catalog.setval('public.marketplace_orders_id_seq', 1, false);
SELECT pg_catalog.setval('public.merchant_payouts_id_seq', 1, false);
SELECT pg_catalog.setval('public.notifications_id_seq', 40, true);
SELECT pg_catalog.setval('public.payment_gateways_id_seq', 2, false);
SELECT pg_catalog.setval('public.payment_methods_id_seq', 27, false);
SELECT pg_catalog.setval('public.payment_transactions_id_seq', 1, false);
SELECT pg_catalog.setval('public.permissions_id_seq', 86, true);
SELECT pg_catalog.setval('public.price_history_id_seq', 65, true);
SELECT pg_catalog.setval('public.product_variant_attribute_values_id_seq', 37, true);
SELECT pg_catalog.setval('public.product_variant_attributes_id_seq', 36, true);
SELECT pg_catalog.setval('public.refresh_tokens_id_seq', 5211, true);
SELECT pg_catalog.setval('public.review_confirmations_id_seq', 1, false);
SELECT pg_catalog.setval('public.review_disputes_id_seq', 1, false);
SELECT pg_catalog.setval('public.review_responses_id_seq', 1, false);
SELECT pg_catalog.setval('public.reviews_id_seq', 13, true);
SELECT pg_catalog.setval('public.role_audit_log_id_seq', 3, true);
SELECT pg_catalog.setval('public.roles_id_seq', 30, true);
SELECT pg_catalog.setval('public.search_behavior_metrics_id_seq', 1, false);
SELECT pg_catalog.setval('public.search_config_id_seq', 2, false);
SELECT pg_catalog.setval('public.search_optimization_sessions_id_seq', 1, false);
SELECT pg_catalog.setval('public.search_queries_id_seq', 142, true);
SELECT pg_catalog.setval('public.search_statistics_id_seq', 1, false);
SELECT pg_catalog.setval('public.search_synonyms_config_id_seq', 8, false);
SELECT pg_catalog.setval('public.search_synonyms_id_seq', 41, false);
SELECT pg_catalog.setval('public.search_weights_history_id_seq', 1, false);
SELECT pg_catalog.setval('public.search_weights_id_seq', 17, false);
SELECT pg_catalog.setval('public.shopping_cart_items_id_seq', 56, true);
SELECT pg_catalog.setval('public.shopping_carts_id_seq', 10, true);
SELECT pg_catalog.setval('public.storefront_delivery_options_id_seq', 1, true);
SELECT pg_catalog.setval('public.storefront_hours_id_seq', 25, true);
SELECT pg_catalog.setval('public.storefront_inventory_movements_id_seq', 1, true);
SELECT pg_catalog.setval('public.storefront_order_items_id_seq', 19, true);
SELECT pg_catalog.setval('public.storefront_orders_id_seq', 58, true);
SELECT pg_catalog.setval('public.storefront_payment_methods_id_seq', 15, true);
SELECT pg_catalog.setval('public.storefront_product_attributes_id_seq', 1, true);
SELECT pg_catalog.setval('public.storefront_product_images_id_seq', 37, true);
SELECT pg_catalog.setval('public.storefront_product_variant_images_id_seq', 1, false);
SELECT pg_catalog.setval('public.storefront_product_variants_id_seq', 97, true);
SELECT pg_catalog.setval('public.storefront_staff_id_seq', 13, true);
SELECT pg_catalog.setval('public.storefronts_id_seq', 21, true);
SELECT pg_catalog.setval('public.subscription_history_id_seq', 1, false);
SELECT pg_catalog.setval('public.subscription_payments_id_seq', 1, false);
SELECT pg_catalog.setval('public.subscription_plans_id_seq', 4, true);
SELECT pg_catalog.setval('public.subscription_usage_id_seq', 1, false);
SELECT pg_catalog.setval('public.translation_audit_log_id_seq', 2, true);
SELECT pg_catalog.setval('public.translation_providers_id_seq', 4, true);
SELECT pg_catalog.setval('public.translation_quality_metrics_id_seq', 1, false);
SELECT pg_catalog.setval('public.translation_sync_conflicts_id_seq', 1, false);
SELECT pg_catalog.setval('public.translation_tasks_id_seq', 1, false);
SELECT pg_catalog.setval('public.translations_id_seq', 3543, true);
SELECT pg_catalog.setval('public.transliteration_rules_id_seq', 52, false);
SELECT pg_catalog.setval('public.unified_geo_id_seq', 237, true);
SELECT pg_catalog.setval('public.user_behavior_events_id_seq', 918, true);
SELECT pg_catalog.setval('public.user_contacts_id_seq', 25, true);
SELECT pg_catalog.setval('public.user_storefronts_id_seq', 1, false);
SELECT pg_catalog.setval('public.user_subscriptions_id_seq', 1, false);
SELECT pg_catalog.setval('public.users_id_seq', 10, true);
SELECT pg_catalog.setval('public.variant_attribute_mappings_id_seq', 4, true);
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
    ADD CONSTRAINT attribute_groups_name_key UNIQUE (name);
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
ALTER TABLE ONLY public.category_attributes
    ADD CONSTRAINT category_attributes_pkey PRIMARY KEY (id);

SELECT pg_catalog.setval('public.search_behavior_metrics_id_seq', 1, false);
SELECT pg_catalog.setval('public.search_config_id_seq', 2, false);
SELECT pg_catalog.setval('public.search_optimization_sessions_id_seq', 1, false);
SELECT pg_catalog.setval('public.search_queries_id_seq', 368, true);
SELECT pg_catalog.setval('public.search_statistics_id_seq', 1, false);
SELECT pg_catalog.setval('public.search_synonyms_config_id_seq', 8, false);
SELECT pg_catalog.setval('public.search_synonyms_id_seq', 121, true);
SELECT pg_catalog.setval('public.search_weights_history_id_seq', 1, false);
SELECT pg_catalog.setval('public.search_weights_id_seq', 17, false);
SELECT pg_catalog.setval('public.shopping_cart_items_id_seq', 69, true);
SELECT pg_catalog.setval('public.shopping_carts_id_seq', 24, true);
SELECT pg_catalog.setval('public.storefront_import_errors_id_seq', 1, false);
SELECT pg_catalog.setval('public.storefront_import_jobs_id_seq', 6, true);
SELECT pg_catalog.setval('public.subscription_history_id_seq', 1, false);
SELECT pg_catalog.setval('public.subscription_payments_id_seq', 1, false);
SELECT pg_catalog.setval('public.subscription_plans_id_seq', 4, true);
SELECT pg_catalog.setval('public.subscription_usage_id_seq', 1, false);
SELECT pg_catalog.setval('public.tracking_websocket_connections_id_seq', 1, false);
SELECT pg_catalog.setval('public.translation_audit_log_id_seq', 2, true);
SELECT pg_catalog.setval('public.translation_providers_id_seq', 4, true);
SELECT pg_catalog.setval('public.translation_quality_metrics_id_seq', 1, false);
SELECT pg_catalog.setval('public.translation_sync_conflicts_id_seq', 1, false);
SELECT pg_catalog.setval('public.translation_tasks_id_seq', 1, false);
SELECT pg_catalog.setval('public.translations_id_seq', 6079, true);
SELECT pg_catalog.setval('public.transliteration_rules_id_seq', 52, false);
SELECT pg_catalog.setval('public.unified_attribute_stats_id_seq', 1, false);
SELECT pg_catalog.setval('public.unified_attribute_values_id_seq', 192, true);
SELECT pg_catalog.setval('public.unified_attributes_id_seq', 551, true);
SELECT pg_catalog.setval('public.unified_category_attributes_id_seq', 2469, true);
SELECT pg_catalog.setval('public.unified_geo_id_seq', 2025, true);
SELECT pg_catalog.setval('public.user_behavior_events_id_seq', 2882, true);
SELECT pg_catalog.setval('public.user_car_view_history_id_seq', 1, false);
SELECT pg_catalog.setval('public.user_contacts_id_seq', 32, true);
SELECT pg_catalog.setval('public.user_notification_contacts_id_seq', 1, false);
SELECT pg_catalog.setval('public.user_notification_preferences_id_seq', 1, false);
SELECT pg_catalog.setval('public.user_subscriptions_id_seq', 1, false);
SELECT pg_catalog.setval('public.user_view_history_id_seq', 1, false);
SELECT pg_catalog.setval('public.variant_attribute_mappings_id_seq', 13, true);
SELECT pg_catalog.setval('public.viber_messages_id_seq', 1, false);
SELECT pg_catalog.setval('public.viber_sessions_id_seq', 2, true);
SELECT pg_catalog.setval('public.viber_tracking_sessions_id_seq', 1, false);
SELECT pg_catalog.setval('public.viber_users_id_seq', 2, true);
SELECT pg_catalog.setval('public.view_statistics_id_seq', 1, false);
SELECT pg_catalog.setval('public.vin_accident_history_id_seq', 1, false);
SELECT pg_catalog.setval('public.vin_check_history_id_seq', 1, false);
SELECT pg_catalog.setval('public.vin_decode_cache_id_seq', 1, false);
SELECT pg_catalog.setval('public.vin_ownership_history_id_seq', 1, false);
SELECT pg_catalog.setval('public.vin_recalls_id_seq', 1, false);
ALTER TABLE ONLY public.ai_category_decisions
    ADD CONSTRAINT ai_category_decisions_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.ai_category_decisions
    ADD CONSTRAINT ai_category_decisions_unique_hash UNIQUE (title_hash, entity_type);
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
ALTER TABLE ONLY public.b2c_delivery_options
    ADD CONSTRAINT b2c_delivery_options_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.b2c_favorites
    ADD CONSTRAINT b2c_favorites_pkey PRIMARY KEY (user_id, product_id);
ALTER TABLE ONLY public.b2c_inventory_movements
    ADD CONSTRAINT b2c_inventory_movements_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.b2c_order_items
    ADD CONSTRAINT b2c_order_items_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.b2c_orders
    ADD CONSTRAINT b2c_orders_order_number_key UNIQUE (order_number);
ALTER TABLE ONLY public.b2c_orders
    ADD CONSTRAINT b2c_orders_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.b2c_payment_methods
    ADD CONSTRAINT b2c_payment_methods_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.b2c_product_attributes
    ADD CONSTRAINT b2c_product_attributes_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.b2c_product_attributes
    ADD CONSTRAINT b2c_product_attributes_product_id_attribute_id_key UNIQUE (product_id, attribute_id);
ALTER TABLE ONLY public.b2c_product_images
    ADD CONSTRAINT b2c_product_images_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.b2c_product_variant_images
    ADD CONSTRAINT b2c_product_variant_images_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.b2c_product_variants
    ADD CONSTRAINT b2c_product_variants_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.b2c_product_variants
    ADD CONSTRAINT b2c_product_variants_sku_key UNIQUE (sku);
ALTER TABLE ONLY public.b2c_products
    ADD CONSTRAINT b2c_products_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.b2c_store_hours
    ADD CONSTRAINT b2c_store_hours_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.b2c_store_hours
    ADD CONSTRAINT b2c_store_hours_storefront_id_day_of_week_special_date_key UNIQUE (storefront_id, day_of_week, special_date);
ALTER TABLE ONLY public.b2c_store_staff
    ADD CONSTRAINT b2c_store_staff_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.b2c_store_staff
    ADD CONSTRAINT b2c_store_staff_storefront_id_user_id_key UNIQUE (storefront_id, user_id);
ALTER TABLE ONLY public.b2c_stores
    ADD CONSTRAINT b2c_stores_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.b2c_stores
    ADD CONSTRAINT b2c_stores_slug_key UNIQUE (slug);
ALTER TABLE ONLY public.balance_transactions
    ADD CONSTRAINT balance_transactions_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.bex_configuration
    ADD CONSTRAINT bex_configuration_key_key UNIQUE (key);
ALTER TABLE ONLY public.bex_configuration
    ADD CONSTRAINT bex_configuration_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.bex_shipments
    ADD CONSTRAINT bex_shipments_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.bex_shipments
    ADD CONSTRAINT bex_shipments_tracking_number_key UNIQUE (tracking_number);
ALTER TABLE ONLY public.bex_tracking_events
    ADD CONSTRAINT bex_tracking_events_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.c2c_categories
    ADD CONSTRAINT c2c_categories_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.c2c_categories
    ADD CONSTRAINT c2c_categories_slug_key UNIQUE (slug);
ALTER TABLE ONLY public.c2c_chats
    ADD CONSTRAINT c2c_chats_listing_id_buyer_id_seller_id_key UNIQUE (listing_id, buyer_id, seller_id);
ALTER TABLE ONLY public.c2c_chats
    ADD CONSTRAINT c2c_chats_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.c2c_chats
    ADD CONSTRAINT c2c_chats_storefront_product_id_buyer_id_seller_id_key UNIQUE (storefront_product_id, buyer_id, seller_id);
ALTER TABLE ONLY public.c2c_favorites
    ADD CONSTRAINT c2c_favorites_pkey PRIMARY KEY (user_id, listing_id);
ALTER TABLE ONLY public.c2c_images
    ADD CONSTRAINT c2c_images_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.c2c_listing_variants
    ADD CONSTRAINT c2c_listing_variants_listing_id_sku_key UNIQUE (listing_id, sku);
ALTER TABLE ONLY public.c2c_listing_variants
    ADD CONSTRAINT c2c_listing_variants_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.c2c_listings
    ADD CONSTRAINT c2c_listings_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.c2c_messages
    ADD CONSTRAINT c2c_messages_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.c2c_orders
    ADD CONSTRAINT c2c_orders_pkey PRIMARY KEY (id);
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

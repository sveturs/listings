ALTER SEQUENCE public.item_performance_metrics_id_seq OWNED BY public.item_performance_metrics.id;
ALTER SEQUENCE public.listing_attribute_values_id_seq OWNED BY public.listing_attribute_values.id;
ALTER SEQUENCE public.listing_views_id_seq OWNED BY public.listing_views.id;
ALTER SEQUENCE public.listings_geo_id_seq OWNED BY public.listings_geo.id;
ALTER SEQUENCE public.map_items_cache_id_seq OWNED BY public.map_items_cache.id;
ALTER SEQUENCE public.marketplace_categories_id_seq OWNED BY public.marketplace_categories.id;
ALTER SEQUENCE public.marketplace_chats_id_seq OWNED BY public.marketplace_chats.id;
ALTER SEQUENCE public.marketplace_images_id_seq OWNED BY public.marketplace_images.id;
ALTER SEQUENCE public.marketplace_messages_id_seq OWNED BY public.marketplace_messages.id;
ALTER SEQUENCE public.marketplace_orders_id_seq OWNED BY public.marketplace_orders.id;
ALTER SEQUENCE public.merchant_payouts_id_seq OWNED BY public.merchant_payouts.id;
ALTER SEQUENCE public.notifications_id_seq OWNED BY public.notifications.id;
ALTER SEQUENCE public.payment_gateways_id_seq OWNED BY public.payment_gateways.id;
ALTER SEQUENCE public.payment_methods_id_seq OWNED BY public.payment_methods.id;
ALTER SEQUENCE public.payment_transactions_id_seq OWNED BY public.payment_transactions.id;
ALTER SEQUENCE public.permissions_id_seq OWNED BY public.permissions.id;
ALTER SEQUENCE public.post_express_api_logs_id_seq OWNED BY public.post_express_api_logs.id;
ALTER SEQUENCE public.post_express_locations_id_seq OWNED BY public.post_express_locations.id;
ALTER SEQUENCE public.post_express_offices_id_seq OWNED BY public.post_express_offices.id;
ALTER SEQUENCE public.post_express_rates_id_seq OWNED BY public.post_express_rates.id;
ALTER SEQUENCE public.post_express_settings_id_seq OWNED BY public.post_express_settings.id;
ALTER SEQUENCE public.post_express_shipments_id_seq OWNED BY public.post_express_shipments.id;
ALTER SEQUENCE public.post_express_tracking_events_id_seq OWNED BY public.post_express_tracking_events.id;
ALTER SEQUENCE public.price_history_id_seq OWNED BY public.price_history.id;
ALTER SEQUENCE public.product_variant_attribute_values_id_seq OWNED BY public.product_variant_attribute_values.id;
ALTER SEQUENCE public.product_variant_attributes_id_seq OWNED BY public.product_variant_attributes.id;
ALTER SEQUENCE public.refresh_tokens_id_seq OWNED BY public.refresh_tokens.id;
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
ALTER SEQUENCE public.storefront_fbs_settings_id_seq OWNED BY public.storefront_fbs_settings.id;
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
ALTER SEQUENCE public.translation_audit_log_id_seq OWNED BY public.translation_audit_log.id;
ALTER SEQUENCE public.translation_providers_id_seq OWNED BY public.translation_providers.id;
ALTER SEQUENCE public.translation_quality_metrics_id_seq OWNED BY public.translation_quality_metrics.id;
ALTER SEQUENCE public.translation_sync_conflicts_id_seq OWNED BY public.translation_sync_conflicts.id;
ALTER SEQUENCE public.translation_tasks_id_seq OWNED BY public.translation_tasks.id;
ALTER SEQUENCE public.translations_id_seq OWNED BY public.translations.id;
ALTER SEQUENCE public.transliteration_rules_id_seq OWNED BY public.transliteration_rules.id;
ALTER SEQUENCE public.unified_geo_id_seq OWNED BY public.unified_geo.id;
ALTER SEQUENCE public.user_behavior_events_id_seq OWNED BY public.user_behavior_events.id;
ALTER SEQUENCE public.user_contacts_id_seq OWNED BY public.user_contacts.id;
ALTER SEQUENCE public.user_storefronts_id_seq OWNED BY public.user_storefronts.id;
ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;
ALTER SEQUENCE public.variant_attribute_mappings_id_seq OWNED BY public.variant_attribute_mappings.id;
ALTER SEQUENCE public.warehouse_inventory_id_seq OWNED BY public.warehouse_inventory.id;
ALTER SEQUENCE public.warehouse_invoices_id_seq OWNED BY public.warehouse_invoices.id;
ALTER SEQUENCE public.warehouse_movements_id_seq OWNED BY public.warehouse_movements.id;
ALTER SEQUENCE public.warehouse_pickup_orders_id_seq OWNED BY public.warehouse_pickup_orders.id;
ALTER SEQUENCE public.warehouses_id_seq OWNED BY public.warehouses.id;
ALTER TABLE ONLY public.attribute_group_items
    ADD CONSTRAINT attribute_group_items_attribute_id_fkey FOREIGN KEY (attribute_id) REFERENCES public.category_attributes(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.attribute_group_items
    ADD CONSTRAINT attribute_group_items_group_id_fkey FOREIGN KEY (group_id) REFERENCES public.attribute_groups(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.car_generations
    ADD CONSTRAINT car_generations_model_id_fkey FOREIGN KEY (model_id) REFERENCES public.car_models(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.car_models
    ADD CONSTRAINT car_models_make_id_fkey FOREIGN KEY (make_id) REFERENCES public.car_makes(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.category_attribute_groups
    ADD CONSTRAINT category_attribute_groups_group_id_fkey FOREIGN KEY (group_id) REFERENCES public.attribute_groups(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.category_attribute_mapping
    ADD CONSTRAINT category_attribute_mapping_attribute_id_fkey FOREIGN KEY (attribute_id) REFERENCES public.category_attributes(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.category_keywords
    ADD CONSTRAINT category_keywords_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.marketplace_categories(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.category_variant_attributes
    ADD CONSTRAINT category_variant_attributes_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.marketplace_categories(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.component_templates
    ADD CONSTRAINT component_templates_component_id_fkey FOREIGN KEY (component_id) REFERENCES public.custom_ui_components(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.custom_ui_component_usage
    ADD CONSTRAINT custom_ui_component_usage_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.marketplace_categories(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.custom_ui_component_usage
    ADD CONSTRAINT custom_ui_component_usage_component_id_fkey FOREIGN KEY (component_id) REFERENCES public.custom_ui_components(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.custom_ui_component_usage
    ADD CONSTRAINT custom_ui_component_usage_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.users(id);
ALTER TABLE ONLY public.custom_ui_component_usage
    ADD CONSTRAINT custom_ui_component_usage_updated_by_fkey FOREIGN KEY (updated_by) REFERENCES public.users(id);
ALTER TABLE ONLY public.escrow_payments
    ADD CONSTRAINT escrow_payments_payment_transaction_id_fkey FOREIGN KEY (payment_transaction_id) REFERENCES public.payment_transactions(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.variant_attribute_mappings
    ADD CONSTRAINT fk_category_attribute FOREIGN KEY (category_attribute_id) REFERENCES public.category_attributes(id) ON DELETE CASCADE;
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
ALTER TABLE ONLY public.variant_attribute_mappings
    ADD CONSTRAINT fk_variant_attribute FOREIGN KEY (variant_attribute_id) REFERENCES public.product_variant_attributes(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.import_history
    ADD CONSTRAINT import_history_source_id_fkey FOREIGN KEY (source_id) REFERENCES public.import_sources(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.import_sources
    ADD CONSTRAINT import_sources_storefront_id_fkey FOREIGN KEY (storefront_id) REFERENCES public.user_storefronts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.imported_categories
    ADD CONSTRAINT imported_categories_source_id_fkey FOREIGN KEY (source_id) REFERENCES public.import_sources(id) ON DELETE CASCADE;

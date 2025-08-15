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
ALTER TABLE ONLY public.listing_attribute_values
    ADD CONSTRAINT listing_attribute_values_attribute_id_fkey FOREIGN KEY (attribute_id) REFERENCES public.category_attributes(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.listings_geo
    ADD CONSTRAINT listings_geo_district_id_fkey FOREIGN KEY (district_id) REFERENCES public.districts(id) ON DELETE SET NULL;
ALTER TABLE ONLY public.listings_geo
    ADD CONSTRAINT listings_geo_municipality_id_fkey FOREIGN KEY (municipality_id) REFERENCES public.municipalities(id) ON DELETE SET NULL;
ALTER TABLE ONLY public.marketplace_categories
    ADD CONSTRAINT marketplace_categories_parent_id_fkey FOREIGN KEY (parent_id) REFERENCES public.marketplace_categories(id);
ALTER TABLE ONLY public.marketplace_chats
    ADD CONSTRAINT marketplace_chats_buyer_id_fkey FOREIGN KEY (buyer_id) REFERENCES public.users(id);
ALTER TABLE ONLY public.marketplace_chats
    ADD CONSTRAINT marketplace_chats_listing_id_fkey FOREIGN KEY (listing_id) REFERENCES public.marketplace_listings(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.marketplace_chats
    ADD CONSTRAINT marketplace_chats_seller_id_fkey FOREIGN KEY (seller_id) REFERENCES public.users(id);
ALTER TABLE ONLY public.marketplace_favorites
    ADD CONSTRAINT marketplace_favorites_listing_id_fkey FOREIGN KEY (listing_id) REFERENCES public.marketplace_listings(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.marketplace_favorites
    ADD CONSTRAINT marketplace_favorites_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);
ALTER TABLE ONLY public.marketplace_images
    ADD CONSTRAINT marketplace_images_listing_id_fkey FOREIGN KEY (listing_id) REFERENCES public.marketplace_listings(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.marketplace_listings
    ADD CONSTRAINT marketplace_listings_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.marketplace_categories(id);
ALTER TABLE ONLY public.marketplace_listings
    ADD CONSTRAINT marketplace_listings_storefront_id_fkey FOREIGN KEY (storefront_id) REFERENCES public.user_storefronts(id) ON DELETE SET NULL;
ALTER TABLE ONLY public.marketplace_listings
    ADD CONSTRAINT marketplace_listings_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);
ALTER TABLE ONLY public.marketplace_messages
    ADD CONSTRAINT marketplace_messages_chat_id_fkey FOREIGN KEY (chat_id) REFERENCES public.marketplace_chats(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.marketplace_messages
    ADD CONSTRAINT marketplace_messages_listing_id_fkey FOREIGN KEY (listing_id) REFERENCES public.marketplace_listings(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.marketplace_messages
    ADD CONSTRAINT marketplace_messages_receiver_id_fkey FOREIGN KEY (receiver_id) REFERENCES public.users(id);
ALTER TABLE ONLY public.marketplace_messages
    ADD CONSTRAINT marketplace_messages_sender_id_fkey FOREIGN KEY (sender_id) REFERENCES public.users(id);
ALTER TABLE ONLY public.marketplace_orders
    ADD CONSTRAINT marketplace_orders_payment_transaction_id_fkey FOREIGN KEY (payment_transaction_id) REFERENCES public.payment_transactions(id);
ALTER TABLE ONLY public.merchant_payouts
    ADD CONSTRAINT merchant_payouts_gateway_id_fkey FOREIGN KEY (gateway_id) REFERENCES public.payment_gateways(id);
ALTER TABLE ONLY public.municipalities
    ADD CONSTRAINT municipalities_district_id_fkey FOREIGN KEY (district_id) REFERENCES public.districts(id) ON DELETE SET NULL;
ALTER TABLE ONLY public.notification_settings
    ADD CONSTRAINT notification_settings_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);
ALTER TABLE ONLY public.notifications
    ADD CONSTRAINT notifications_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);
ALTER TABLE ONLY public.payment_transactions
    ADD CONSTRAINT payment_transactions_gateway_id_fkey FOREIGN KEY (gateway_id) REFERENCES public.payment_gateways(id);
ALTER TABLE ONLY public.payment_transactions
    ADD CONSTRAINT payment_transactions_storefront_id_fkey FOREIGN KEY (storefront_id) REFERENCES public.storefronts(id);
ALTER TABLE ONLY public.product_variant_attribute_values
    ADD CONSTRAINT product_variant_attribute_values_attribute_id_fkey FOREIGN KEY (attribute_id) REFERENCES public.product_variant_attributes(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.review_confirmations
    ADD CONSTRAINT review_confirmations_confirmed_by_fkey FOREIGN KEY (confirmed_by) REFERENCES public.users(id);
ALTER TABLE ONLY public.review_confirmations
    ADD CONSTRAINT review_confirmations_review_id_fkey FOREIGN KEY (review_id) REFERENCES public.reviews(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.review_disputes
    ADD CONSTRAINT review_disputes_admin_id_fkey FOREIGN KEY (admin_id) REFERENCES public.users(id);
ALTER TABLE ONLY public.review_disputes
    ADD CONSTRAINT review_disputes_disputed_by_fkey FOREIGN KEY (disputed_by) REFERENCES public.users(id);
ALTER TABLE ONLY public.review_disputes
    ADD CONSTRAINT review_disputes_review_id_fkey FOREIGN KEY (review_id) REFERENCES public.reviews(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.review_responses
    ADD CONSTRAINT review_responses_review_id_fkey FOREIGN KEY (review_id) REFERENCES public.reviews(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.review_responses
    ADD CONSTRAINT review_responses_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.review_votes
    ADD CONSTRAINT review_votes_review_id_fkey FOREIGN KEY (review_id) REFERENCES public.reviews(id) ON DELETE CASCADE;

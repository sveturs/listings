ALTER TABLE ONLY public.category_detection_feedback
    ADD CONSTRAINT category_detection_feedback_correct_category_id_fkey FOREIGN KEY (correct_category_id) REFERENCES public.marketplace_categories(id);
ALTER TABLE ONLY public.category_detection_feedback
    ADD CONSTRAINT category_detection_feedback_detected_category_id_fkey FOREIGN KEY (detected_category_id) REFERENCES public.marketplace_categories(id);
ALTER TABLE ONLY public.category_detection_feedback
    ADD CONSTRAINT category_detection_feedback_listing_id_fkey FOREIGN KEY (listing_id) REFERENCES public.marketplace_listings(id) ON DELETE SET NULL;
ALTER TABLE ONLY public.category_keyword_weights
    ADD CONSTRAINT category_keyword_weights_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.marketplace_categories(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.category_keywords
    ADD CONSTRAINT category_keywords_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.marketplace_categories(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.category_variant_attributes
    ADD CONSTRAINT category_variant_attributes_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.marketplace_categories(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.component_templates
    ADD CONSTRAINT component_templates_component_id_fkey FOREIGN KEY (component_id) REFERENCES public.custom_ui_components(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.courier_location_history
    ADD CONSTRAINT courier_location_history_courier_id_fkey FOREIGN KEY (courier_id) REFERENCES public.couriers(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.courier_location_history
    ADD CONSTRAINT courier_location_history_delivery_id_fkey FOREIGN KEY (delivery_id) REFERENCES public.deliveries(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.courier_zones
    ADD CONSTRAINT courier_zones_courier_id_fkey FOREIGN KEY (courier_id) REFERENCES public.couriers(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.couriers
    ADD CONSTRAINT couriers_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);
ALTER TABLE ONLY public.custom_ui_component_usage
    ADD CONSTRAINT custom_ui_component_usage_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.marketplace_categories(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.custom_ui_component_usage
    ADD CONSTRAINT custom_ui_component_usage_component_id_fkey FOREIGN KEY (component_id) REFERENCES public.custom_ui_components(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.custom_ui_component_usage
    ADD CONSTRAINT custom_ui_component_usage_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.users(id);
ALTER TABLE ONLY public.custom_ui_component_usage
    ADD CONSTRAINT custom_ui_component_usage_updated_by_fkey FOREIGN KEY (updated_by) REFERENCES public.users(id);
ALTER TABLE ONLY public.deliveries
    ADD CONSTRAINT deliveries_courier_id_fkey FOREIGN KEY (courier_id) REFERENCES public.couriers(id);
ALTER TABLE ONLY public.deliveries
    ADD CONSTRAINT deliveries_order_id_fkey FOREIGN KEY (order_id) REFERENCES public.storefront_orders(id);
ALTER TABLE ONLY public.delivery_category_defaults
    ADD CONSTRAINT delivery_category_defaults_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.marketplace_categories(id);
ALTER TABLE ONLY public.delivery_notifications
    ADD CONSTRAINT delivery_notifications_shipment_id_fkey FOREIGN KEY (shipment_id) REFERENCES public.delivery_shipments(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.delivery_notifications
    ADD CONSTRAINT delivery_notifications_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);
ALTER TABLE ONLY public.delivery_pricing_rules
    ADD CONSTRAINT delivery_pricing_rules_provider_id_fkey FOREIGN KEY (provider_id) REFERENCES public.delivery_providers(id);
ALTER TABLE ONLY public.delivery_shipments
    ADD CONSTRAINT delivery_shipments_order_id_fkey FOREIGN KEY (order_id) REFERENCES public.marketplace_orders(id);
ALTER TABLE ONLY public.delivery_shipments
    ADD CONSTRAINT delivery_shipments_provider_id_fkey FOREIGN KEY (provider_id) REFERENCES public.delivery_providers(id);
ALTER TABLE ONLY public.delivery_tracking_events
    ADD CONSTRAINT delivery_tracking_events_provider_id_fkey FOREIGN KEY (provider_id) REFERENCES public.delivery_providers(id);
ALTER TABLE ONLY public.delivery_tracking_events
    ADD CONSTRAINT delivery_tracking_events_shipment_id_fkey FOREIGN KEY (shipment_id) REFERENCES public.delivery_shipments(id);
ALTER TABLE ONLY public.escrow_payments
    ADD CONSTRAINT escrow_payments_payment_transaction_id_fkey FOREIGN KEY (payment_transaction_id) REFERENCES public.payment_transactions(id) ON DELETE CASCADE;
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
ALTER TABLE ONLY public.post_express_shipments
    ADD CONSTRAINT fk_recipient_location FOREIGN KEY (recipient_location_id) REFERENCES public.post_express_locations(id) ON DELETE SET NULL;
ALTER TABLE ONLY public.post_express_shipments
    ADD CONSTRAINT fk_sender_location FOREIGN KEY (sender_location_id) REFERENCES public.post_express_locations(id) ON DELETE SET NULL;
ALTER TABLE ONLY public.import_history
    ADD CONSTRAINT import_history_source_id_fkey FOREIGN KEY (source_id) REFERENCES public.import_sources(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.import_sources
    ADD CONSTRAINT import_sources_storefront_id_fkey FOREIGN KEY (storefront_id) REFERENCES public.user_storefronts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.imported_categories
    ADD CONSTRAINT imported_categories_source_id_fkey FOREIGN KEY (source_id) REFERENCES public.import_sources(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.listing_attribute_values
    ADD CONSTRAINT listing_attribute_values_attribute_id_fkey FOREIGN KEY (attribute_id) REFERENCES public.unified_attributes(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.listing_attribute_values
    ADD CONSTRAINT listing_attribute_values_listing_id_fkey FOREIGN KEY (listing_id) REFERENCES public.marketplace_listings(id) ON DELETE CASCADE;
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
ALTER TABLE ONLY public.marketplace_listing_variants
    ADD CONSTRAINT marketplace_listing_variants_listing_id_fkey FOREIGN KEY (listing_id) REFERENCES public.marketplace_listings(id) ON DELETE CASCADE;
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
    ADD CONSTRAINT marketplace_orders_delivery_shipment_id_fkey FOREIGN KEY (delivery_shipment_id) REFERENCES public.delivery_shipments(id);
ALTER TABLE ONLY public.marketplace_orders
    ADD CONSTRAINT marketplace_orders_payment_transaction_id_fkey FOREIGN KEY (payment_transaction_id) REFERENCES public.payment_transactions(id) ON DELETE SET NULL;
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
ALTER TABLE ONLY public.post_express_offices
    ADD CONSTRAINT post_express_offices_location_id_fkey FOREIGN KEY (location_id) REFERENCES public.post_express_locations(id);
ALTER TABLE ONLY public.post_express_tracking_events
    ADD CONSTRAINT post_express_tracking_events_shipment_id_fkey FOREIGN KEY (shipment_id) REFERENCES public.post_express_shipments(id) ON DELETE CASCADE;
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
ALTER TABLE ONLY public.review_votes
    ADD CONSTRAINT review_votes_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.reviews
    ADD CONSTRAINT reviews_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.role_audit_log
    ADD CONSTRAINT role_audit_log_new_role_id_fkey FOREIGN KEY (new_role_id) REFERENCES public.roles(id);
ALTER TABLE ONLY public.role_audit_log
    ADD CONSTRAINT role_audit_log_old_role_id_fkey FOREIGN KEY (old_role_id) REFERENCES public.roles(id);
ALTER TABLE ONLY public.role_audit_log
    ADD CONSTRAINT role_audit_log_target_user_id_fkey FOREIGN KEY (target_user_id) REFERENCES public.users(id);
ALTER TABLE ONLY public.role_audit_log
    ADD CONSTRAINT role_audit_log_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);
ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT role_permissions_permission_id_fkey FOREIGN KEY (permission_id) REFERENCES public.permissions(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT role_permissions_role_id_fkey FOREIGN KEY (role_id) REFERENCES public.roles(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.search_optimization_sessions
    ADD CONSTRAINT search_optimization_sessions_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.admin_users(id) ON DELETE RESTRICT;
ALTER TABLE ONLY public.search_weights
    ADD CONSTRAINT search_weights_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.admin_users(id) ON DELETE SET NULL;
ALTER TABLE ONLY public.search_weights_history
    ADD CONSTRAINT search_weights_history_changed_by_fkey FOREIGN KEY (changed_by) REFERENCES public.admin_users(id) ON DELETE SET NULL;
ALTER TABLE ONLY public.search_weights_history
    ADD CONSTRAINT search_weights_history_weight_id_fkey FOREIGN KEY (weight_id) REFERENCES public.search_weights(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.search_weights
    ADD CONSTRAINT search_weights_updated_by_fkey FOREIGN KEY (updated_by) REFERENCES public.admin_users(id) ON DELETE SET NULL;
ALTER TABLE ONLY public.shopping_cart_items
    ADD CONSTRAINT shopping_cart_items_cart_id_fkey FOREIGN KEY (cart_id) REFERENCES public.shopping_carts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.shopping_cart_items
    ADD CONSTRAINT shopping_cart_items_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.storefront_products(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.shopping_carts
    ADD CONSTRAINT shopping_carts_storefront_id_fkey FOREIGN KEY (storefront_id) REFERENCES public.storefronts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.storefront_delivery_options
    ADD CONSTRAINT storefront_delivery_options_storefront_id_fkey FOREIGN KEY (storefront_id) REFERENCES public.storefronts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.storefront_favorites
    ADD CONSTRAINT storefront_favorites_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.storefront_products(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.storefront_favorites
    ADD CONSTRAINT storefront_favorites_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.storefront_hours
    ADD CONSTRAINT storefront_hours_storefront_id_fkey FOREIGN KEY (storefront_id) REFERENCES public.storefronts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.storefront_inventory_movements
    ADD CONSTRAINT storefront_inventory_movements_storefront_product_id_fkey FOREIGN KEY (storefront_product_id) REFERENCES public.storefront_products(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.storefront_order_items
    ADD CONSTRAINT storefront_order_items_order_id_fkey FOREIGN KEY (order_id) REFERENCES public.storefront_orders(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.storefront_order_items
    ADD CONSTRAINT storefront_order_items_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.storefront_products(id);
ALTER TABLE ONLY public.storefront_order_items
    ADD CONSTRAINT storefront_order_items_variant_id_fkey FOREIGN KEY (variant_id) REFERENCES public.storefront_product_variants(id);
ALTER TABLE ONLY public.storefront_orders
    ADD CONSTRAINT storefront_orders_payment_transaction_id_fkey FOREIGN KEY (payment_transaction_id) REFERENCES public.payment_transactions(id) ON DELETE SET NULL;

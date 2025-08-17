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
ALTER TABLE ONLY public.post_express_api_logs
    ADD CONSTRAINT post_express_api_logs_shipment_id_fkey FOREIGN KEY (shipment_id) REFERENCES public.post_express_shipments(id);
ALTER TABLE ONLY public.post_express_offices
    ADD CONSTRAINT post_express_offices_location_id_fkey FOREIGN KEY (location_id) REFERENCES public.post_express_locations(id);
ALTER TABLE ONLY public.post_express_shipments
    ADD CONSTRAINT post_express_shipments_marketplace_order_id_fkey FOREIGN KEY (marketplace_order_id) REFERENCES public.marketplace_orders(id);
ALTER TABLE ONLY public.post_express_shipments
    ADD CONSTRAINT post_express_shipments_recipient_location_id_fkey FOREIGN KEY (recipient_location_id) REFERENCES public.post_express_locations(id);
ALTER TABLE ONLY public.post_express_shipments
    ADD CONSTRAINT post_express_shipments_sender_location_id_fkey FOREIGN KEY (sender_location_id) REFERENCES public.post_express_locations(id);
ALTER TABLE ONLY public.post_express_shipments
    ADD CONSTRAINT post_express_shipments_storefront_order_id_fkey FOREIGN KEY (storefront_order_id) REFERENCES public.storefront_orders(id);
ALTER TABLE ONLY public.post_express_tracking_events
    ADD CONSTRAINT post_express_tracking_events_shipment_id_fkey FOREIGN KEY (shipment_id) REFERENCES public.post_express_shipments(id);
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
ALTER TABLE ONLY public.storefront_fbs_settings
    ADD CONSTRAINT storefront_fbs_settings_storefront_id_fkey FOREIGN KEY (storefront_id) REFERENCES public.storefronts(id);
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
ALTER TABLE ONLY public.storefront_orders
    ADD CONSTRAINT storefront_orders_storefront_id_fkey FOREIGN KEY (storefront_id) REFERENCES public.storefronts(id) ON DELETE RESTRICT;
ALTER TABLE ONLY public.storefront_payment_methods
    ADD CONSTRAINT storefront_payment_methods_storefront_id_fkey FOREIGN KEY (storefront_id) REFERENCES public.storefronts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.storefront_product_attributes
    ADD CONSTRAINT storefront_product_attributes_attribute_id_fkey FOREIGN KEY (attribute_id) REFERENCES public.product_variant_attributes(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.storefront_product_attributes
    ADD CONSTRAINT storefront_product_attributes_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.storefront_products(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.storefront_product_images
    ADD CONSTRAINT storefront_product_images_storefront_product_id_fkey FOREIGN KEY (storefront_product_id) REFERENCES public.storefront_products(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.storefront_product_variant_images
    ADD CONSTRAINT storefront_product_variant_images_variant_id_fkey FOREIGN KEY (variant_id) REFERENCES public.storefront_product_variants(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.storefront_product_variants
    ADD CONSTRAINT storefront_product_variants_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.storefront_products(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.storefront_products
    ADD CONSTRAINT storefront_products_storefront_id_fkey FOREIGN KEY (storefront_id) REFERENCES public.storefronts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.storefront_staff
    ADD CONSTRAINT storefront_staff_storefront_id_fkey FOREIGN KEY (storefront_id) REFERENCES public.storefronts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.translation_audit_log
    ADD CONSTRAINT translation_audit_log_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE SET NULL;
ALTER TABLE ONLY public.translation_quality_metrics
    ADD CONSTRAINT translation_quality_metrics_translation_id_fkey FOREIGN KEY (translation_id) REFERENCES public.translations(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.translation_sync_conflicts
    ADD CONSTRAINT translation_sync_conflicts_resolved_by_fkey FOREIGN KEY (resolved_by) REFERENCES public.users(id) ON DELETE SET NULL;
ALTER TABLE ONLY public.translation_tasks
    ADD CONSTRAINT translation_tasks_assigned_to_fkey FOREIGN KEY (assigned_to) REFERENCES public.users(id) ON DELETE SET NULL;
ALTER TABLE ONLY public.translation_tasks
    ADD CONSTRAINT translation_tasks_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.users(id) ON DELETE SET NULL;
ALTER TABLE ONLY public.translation_tasks
    ADD CONSTRAINT translation_tasks_provider_id_fkey FOREIGN KEY (provider_id) REFERENCES public.translation_providers(id);
ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_assigned_by_fkey FOREIGN KEY (assigned_by) REFERENCES public.users(id);
ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_role_id_fkey FOREIGN KEY (role_id) REFERENCES public.roles(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.user_storefronts
    ADD CONSTRAINT user_storefronts_creation_transaction_id_fkey FOREIGN KEY (creation_transaction_id) REFERENCES public.balance_transactions(id);
ALTER TABLE ONLY public.user_telegram_connections
    ADD CONSTRAINT user_telegram_connections_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);
ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_role_id_fkey FOREIGN KEY (role_id) REFERENCES public.roles(id);
ALTER TABLE ONLY public.warehouse_inventory
    ADD CONSTRAINT warehouse_inventory_marketplace_listing_id_fkey FOREIGN KEY (marketplace_listing_id) REFERENCES public.marketplace_listings(id);
ALTER TABLE ONLY public.warehouse_inventory
    ADD CONSTRAINT warehouse_inventory_storefront_product_id_fkey FOREIGN KEY (storefront_product_id) REFERENCES public.storefront_products(id);
ALTER TABLE ONLY public.warehouse_inventory
    ADD CONSTRAINT warehouse_inventory_warehouse_id_fkey FOREIGN KEY (warehouse_id) REFERENCES public.warehouses(id);
ALTER TABLE ONLY public.warehouse_invoices
    ADD CONSTRAINT warehouse_invoices_storefront_id_fkey FOREIGN KEY (storefront_id) REFERENCES public.storefronts(id);
ALTER TABLE ONLY public.warehouse_invoices
    ADD CONSTRAINT warehouse_invoices_warehouse_id_fkey FOREIGN KEY (warehouse_id) REFERENCES public.warehouses(id);
ALTER TABLE ONLY public.warehouse_movements
    ADD CONSTRAINT warehouse_movements_inventory_id_fkey FOREIGN KEY (inventory_id) REFERENCES public.warehouse_inventory(id);
ALTER TABLE ONLY public.warehouse_movements
    ADD CONSTRAINT warehouse_movements_order_id_fkey FOREIGN KEY (order_id) REFERENCES public.marketplace_orders(id);
ALTER TABLE ONLY public.warehouse_movements
    ADD CONSTRAINT warehouse_movements_performed_by_fkey FOREIGN KEY (performed_by) REFERENCES public.users(id);
ALTER TABLE ONLY public.warehouse_movements
    ADD CONSTRAINT warehouse_movements_shipment_id_fkey FOREIGN KEY (shipment_id) REFERENCES public.post_express_shipments(id);
ALTER TABLE ONLY public.warehouse_movements
    ADD CONSTRAINT warehouse_movements_storefront_order_id_fkey FOREIGN KEY (storefront_order_id) REFERENCES public.storefront_orders(id);
ALTER TABLE ONLY public.warehouse_movements
    ADD CONSTRAINT warehouse_movements_warehouse_id_fkey FOREIGN KEY (warehouse_id) REFERENCES public.warehouses(id);
ALTER TABLE ONLY public.warehouse_pickup_orders
    ADD CONSTRAINT warehouse_pickup_orders_marketplace_order_id_fkey FOREIGN KEY (marketplace_order_id) REFERENCES public.marketplace_orders(id);
ALTER TABLE ONLY public.warehouse_pickup_orders
    ADD CONSTRAINT warehouse_pickup_orders_storefront_order_id_fkey FOREIGN KEY (storefront_order_id) REFERENCES public.storefront_orders(id);
ALTER TABLE ONLY public.warehouse_pickup_orders
    ADD CONSTRAINT warehouse_pickup_orders_warehouse_id_fkey FOREIGN KEY (warehouse_id) REFERENCES public.warehouses(id);

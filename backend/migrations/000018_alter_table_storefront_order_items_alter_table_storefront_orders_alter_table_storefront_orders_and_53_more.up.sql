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
ALTER TABLE ONLY public.storefronts
    ADD CONSTRAINT storefronts_subscription_id_fkey FOREIGN KEY (subscription_id) REFERENCES public.user_subscriptions(id);
ALTER TABLE ONLY public.subscription_history
    ADD CONSTRAINT subscription_history_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.users(id);
ALTER TABLE ONLY public.subscription_history
    ADD CONSTRAINT subscription_history_from_plan_id_fkey FOREIGN KEY (from_plan_id) REFERENCES public.subscription_plans(id);
ALTER TABLE ONLY public.subscription_history
    ADD CONSTRAINT subscription_history_subscription_id_fkey FOREIGN KEY (subscription_id) REFERENCES public.user_subscriptions(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.subscription_history
    ADD CONSTRAINT subscription_history_to_plan_id_fkey FOREIGN KEY (to_plan_id) REFERENCES public.subscription_plans(id);
ALTER TABLE ONLY public.subscription_history
    ADD CONSTRAINT subscription_history_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.subscription_payments
    ADD CONSTRAINT subscription_payments_payment_id_fkey FOREIGN KEY (payment_id) REFERENCES public.payment_transactions(id);
ALTER TABLE ONLY public.subscription_payments
    ADD CONSTRAINT subscription_payments_subscription_id_fkey FOREIGN KEY (subscription_id) REFERENCES public.user_subscriptions(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.subscription_payments
    ADD CONSTRAINT subscription_payments_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.subscription_usage
    ADD CONSTRAINT subscription_usage_storefront_id_fkey FOREIGN KEY (storefront_id) REFERENCES public.storefronts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.subscription_usage
    ADD CONSTRAINT subscription_usage_subscription_id_fkey FOREIGN KEY (subscription_id) REFERENCES public.user_subscriptions(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.tracking_websocket_connections
    ADD CONSTRAINT tracking_websocket_connections_delivery_id_fkey FOREIGN KEY (delivery_id) REFERENCES public.deliveries(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.tracking_websocket_connections
    ADD CONSTRAINT tracking_websocket_connections_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);
ALTER TABLE ONLY public.tracking_websocket_connections
    ADD CONSTRAINT tracking_websocket_connections_viber_user_id_fkey FOREIGN KEY (viber_user_id) REFERENCES public.viber_users(id);
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
ALTER TABLE ONLY public.translations
    ADD CONSTRAINT translations_last_modified_by_fkey FOREIGN KEY (last_modified_by) REFERENCES public.users(id);
ALTER TABLE ONLY public.unified_attribute_stats
    ADD CONSTRAINT unified_attribute_stats_attribute_id_fkey FOREIGN KEY (attribute_id) REFERENCES public.unified_attributes(id);
ALTER TABLE ONLY public.unified_attribute_stats
    ADD CONSTRAINT unified_attribute_stats_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.marketplace_categories(id);
ALTER TABLE ONLY public.unified_attribute_values
    ADD CONSTRAINT unified_attribute_values_attribute_id_fkey FOREIGN KEY (attribute_id) REFERENCES public.unified_attributes(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.unified_category_attributes
    ADD CONSTRAINT unified_category_attributes_attribute_id_fkey FOREIGN KEY (attribute_id) REFERENCES public.unified_attributes(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.unified_category_attributes
    ADD CONSTRAINT unified_category_attributes_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.marketplace_categories(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.user_notification_contacts
    ADD CONSTRAINT user_notification_contacts_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.user_notification_preferences
    ADD CONSTRAINT user_notification_preferences_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_assigned_by_fkey FOREIGN KEY (assigned_by) REFERENCES public.users(id);
ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_role_id_fkey FOREIGN KEY (role_id) REFERENCES public.roles(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.user_storefronts
    ADD CONSTRAINT user_storefronts_creation_transaction_id_fkey FOREIGN KEY (creation_transaction_id) REFERENCES public.balance_transactions(id);
ALTER TABLE ONLY public.user_subscriptions
    ADD CONSTRAINT user_subscriptions_last_payment_id_fkey FOREIGN KEY (last_payment_id) REFERENCES public.payment_transactions(id);
ALTER TABLE ONLY public.user_subscriptions
    ADD CONSTRAINT user_subscriptions_plan_id_fkey FOREIGN KEY (plan_id) REFERENCES public.subscription_plans(id);
ALTER TABLE ONLY public.user_subscriptions
    ADD CONSTRAINT user_subscriptions_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.user_telegram_connections
    ADD CONSTRAINT user_telegram_connections_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);
ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_role_id_fkey FOREIGN KEY (role_id) REFERENCES public.roles(id);
ALTER TABLE ONLY public.variant_attribute_mappings
    ADD CONSTRAINT variant_attribute_mappings_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.marketplace_categories(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.variant_attribute_mappings
    ADD CONSTRAINT variant_attribute_mappings_variant_attribute_id_fkey FOREIGN KEY (variant_attribute_id) REFERENCES public.unified_attributes(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.viber_messages
    ADD CONSTRAINT viber_messages_session_id_fkey FOREIGN KEY (session_id) REFERENCES public.viber_sessions(id) ON DELETE SET NULL;
ALTER TABLE ONLY public.viber_messages
    ADD CONSTRAINT viber_messages_viber_user_id_fkey FOREIGN KEY (viber_user_id) REFERENCES public.viber_users(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.viber_sessions
    ADD CONSTRAINT viber_sessions_viber_user_id_fkey FOREIGN KEY (viber_user_id) REFERENCES public.viber_users(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.viber_tracking_sessions
    ADD CONSTRAINT viber_tracking_sessions_delivery_id_fkey FOREIGN KEY (delivery_id) REFERENCES public.deliveries(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.viber_tracking_sessions
    ADD CONSTRAINT viber_tracking_sessions_viber_user_id_fkey FOREIGN KEY (viber_user_id) REFERENCES public.viber_users(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.viber_users
    ADD CONSTRAINT viber_users_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);

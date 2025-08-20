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

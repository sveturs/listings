CREATE TRIGGER trigger_variants_auto_status BEFORE UPDATE OF stock_quantity ON public.product_variants FOR EACH ROW EXECUTE FUNCTION public.auto_update_variant_status();
CREATE TRIGGER trigger_variants_single_default BEFORE INSERT OR UPDATE OF is_default ON public.product_variants FOR EACH ROW WHEN ((new.is_default = true)) EXECUTE FUNCTION public.enforce_single_default_variant();
CREATE TRIGGER trigger_variants_updated_at BEFORE UPDATE ON public.product_variants FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
CREATE TRIGGER update_listing_images_updated_at BEFORE UPDATE ON public.listing_images FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
CREATE TRIGGER update_listing_locations_updated_at BEFORE UPDATE ON public.listing_locations FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
CREATE TRIGGER update_listing_stats_updated_at BEFORE UPDATE ON public.listing_stats FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
CREATE TRIGGER update_listings_updated_at BEFORE UPDATE ON public.listings FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
ALTER SEQUENCE public.attribute_options_id_seq OWNED BY public.attribute_options.id;
ALTER SEQUENCE public.attribute_search_cache_id_seq OWNED BY public.attribute_search_cache.id;
ALTER SEQUENCE public.attribute_values_id_seq OWNED BY public.attribute_values.id;
ALTER SEQUENCE public.attributes_id_seq OWNED BY public.attributes.id;
ALTER SEQUENCE public.c2c_chats_id_seq OWNED BY public.c2c_chats.id;
ALTER SEQUENCE public.c2c_messages_id_seq OWNED BY public.c2c_messages.id;
ALTER SEQUENCE public.cart_items_id_seq OWNED BY public.cart_items.id;
ALTER SEQUENCE public.category_attributes_id_seq OWNED BY public.category_attributes.id;
ALTER SEQUENCE public.category_variant_attributes_id_seq OWNED BY public.category_variant_attributes.id;
ALTER SEQUENCE public.chat_attachments_id_seq OWNED BY public.chat_attachments.id;
ALTER SEQUENCE public.chats_id_seq OWNED BY public.chats.id;
ALTER SEQUENCE public.indexing_queue_id_seq OWNED BY public.indexing_queue.id;
ALTER SEQUENCE public.inventory_movements_id_seq OWNED BY public.inventory_movements.id;
ALTER SEQUENCE public.inventory_reservations_id_seq OWNED BY public.inventory_reservations.id;
ALTER SEQUENCE public.listing_attribute_values_id_seq OWNED BY public.listing_attribute_values.id;
ALTER SEQUENCE public.listing_attributes_id_seq OWNED BY public.listing_attributes.id;
ALTER SEQUENCE public.listing_images_id_seq OWNED BY public.listing_images.id;
ALTER SEQUENCE public.listing_locations_id_seq OWNED BY public.listing_locations.id;
ALTER SEQUENCE public.listing_tags_id_seq OWNED BY public.listing_tags.id;
ALTER SEQUENCE public.listing_variants_id_seq OWNED BY public.listing_variants.id;
ALTER SEQUENCE public.listings_id_seq OWNED BY public.listings.id;
ALTER SEQUENCE public.messages_id_seq OWNED BY public.messages.id;
ALTER SEQUENCE public.order_items_id_seq OWNED BY public.order_items.id;
ALTER SEQUENCE public.orders_id_seq OWNED BY public.orders.id;
ALTER SEQUENCE public.search_queries_id_seq OWNED BY public.search_queries.id;
ALTER SEQUENCE public.shopping_carts_id_seq OWNED BY public.shopping_carts.id;
ALTER SEQUENCE public.storefront_delivery_options_id_seq OWNED BY public.storefront_delivery_options.id;
ALTER SEQUENCE public.storefront_events_id_seq OWNED BY public.storefront_events.id;
ALTER SEQUENCE public.storefront_hours_id_seq OWNED BY public.storefront_hours.id;
ALTER SEQUENCE public.storefront_invitations_id_seq OWNED BY public.storefront_invitations.id;
ALTER SEQUENCE public.storefront_payment_methods_id_seq OWNED BY public.storefront_payment_methods.id;
ALTER SEQUENCE public.storefront_staff_id_seq OWNED BY public.storefront_staff.id;
ALTER SEQUENCE public.storefronts_id_seq OWNED BY public.storefronts.id;
ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;
CREATE TRIGGER update_indexing_queue_updated_at BEFORE UPDATE ON public.indexing_queue FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
ALTER TABLE ONLY public.attribute_options
    ADD CONSTRAINT attribute_options_attribute_id_fkey FOREIGN KEY (attribute_id) REFERENCES public.attributes(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.attribute_search_cache
    ADD CONSTRAINT attribute_search_cache_listing_id_fkey FOREIGN KEY (listing_id) REFERENCES public.listings(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.attribute_values
    ADD CONSTRAINT attribute_values_attribute_id_fkey FOREIGN KEY (attribute_id) REFERENCES public.attributes(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_parent_id_fkey FOREIGN KEY (parent_id) REFERENCES public.categories(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.category_ai_mapping
    ADD CONSTRAINT category_ai_mapping_target_category_id_fkey FOREIGN KEY (target_category_id) REFERENCES public.categories(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.category_attributes
    ADD CONSTRAINT category_attributes_attribute_id_fkey FOREIGN KEY (attribute_id) REFERENCES public.attributes(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.category_attributes
    ADD CONSTRAINT category_attributes_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.categories(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.category_detections
    ADD CONSTRAINT category_detections_detected_category_id_fkey FOREIGN KEY (detected_category_id) REFERENCES public.categories(id);
ALTER TABLE ONLY public.category_detections
    ADD CONSTRAINT category_detections_user_selected_category_id_fkey FOREIGN KEY (user_selected_category_id) REFERENCES public.categories(id);
ALTER TABLE ONLY public.category_variant_attributes
    ADD CONSTRAINT category_variant_attributes_attribute_id_fkey FOREIGN KEY (attribute_id) REFERENCES public.attributes(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.chat_attachments
    ADD CONSTRAINT fk_attachments_message FOREIGN KEY (message_id) REFERENCES public.messages(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.cart_items
    ADD CONSTRAINT fk_cart_items_cart FOREIGN KEY (cart_id) REFERENCES public.shopping_carts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.cart_items
    ADD CONSTRAINT fk_cart_items_listing FOREIGN KEY (listing_id) REFERENCES public.listings(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.inventory_reservations
    ADD CONSTRAINT fk_inventory_reservations_listing FOREIGN KEY (listing_id) REFERENCES public.listings(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.listing_favorites
    ADD CONSTRAINT fk_listing_favorites_listing_id FOREIGN KEY (listing_id) REFERENCES public.listings(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.listing_variants
    ADD CONSTRAINT fk_listing_variants_listing FOREIGN KEY (listing_id) REFERENCES public.listings(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.messages
    ADD CONSTRAINT fk_messages_chat FOREIGN KEY (chat_id) REFERENCES public.chats(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT fk_order_items_order FOREIGN KEY (order_id) REFERENCES public.orders(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.orders
    ADD CONSTRAINT fk_orders_storefront FOREIGN KEY (storefront_id) REFERENCES public.storefronts(id) ON DELETE RESTRICT;
ALTER TABLE ONLY public.shopping_carts
    ADD CONSTRAINT fk_shopping_carts_storefront FOREIGN KEY (storefront_id) REFERENCES public.storefronts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.storefront_delivery_options
    ADD CONSTRAINT fk_storefront_delivery_options_storefront FOREIGN KEY (storefront_id) REFERENCES public.storefronts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.storefront_hours
    ADD CONSTRAINT fk_storefront_hours_storefront FOREIGN KEY (storefront_id) REFERENCES public.storefronts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.storefront_invitations
    ADD CONSTRAINT fk_storefront_invitations_storefront FOREIGN KEY (storefront_id) REFERENCES public.storefronts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.storefront_payment_methods
    ADD CONSTRAINT fk_storefront_payment_methods_storefront FOREIGN KEY (storefront_id) REFERENCES public.storefronts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.storefront_staff
    ADD CONSTRAINT fk_storefront_staff_storefront FOREIGN KEY (storefront_id) REFERENCES public.storefronts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.indexing_queue
    ADD CONSTRAINT indexing_queue_listing_id_fkey FOREIGN KEY (listing_id) REFERENCES public.listings(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.inventory_movements
    ADD CONSTRAINT inventory_movements_listing_id_fkey FOREIGN KEY (listing_id) REFERENCES public.listings(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.listing_attribute_values
    ADD CONSTRAINT listing_attribute_values_attribute_id_fkey FOREIGN KEY (attribute_id) REFERENCES public.attributes(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.listing_attribute_values
    ADD CONSTRAINT listing_attribute_values_listing_id_fkey FOREIGN KEY (listing_id) REFERENCES public.listings(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.listing_attributes
    ADD CONSTRAINT listing_attributes_listing_id_fkey FOREIGN KEY (listing_id) REFERENCES public.listings(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.listing_images
    ADD CONSTRAINT listing_images_listing_id_fkey FOREIGN KEY (listing_id) REFERENCES public.listings(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.listing_locations
    ADD CONSTRAINT listing_locations_listing_id_fkey FOREIGN KEY (listing_id) REFERENCES public.listings(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.listing_stats
    ADD CONSTRAINT listing_stats_listing_id_fkey FOREIGN KEY (listing_id) REFERENCES public.listings(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.listing_tags
    ADD CONSTRAINT listing_tags_listing_id_fkey FOREIGN KEY (listing_id) REFERENCES public.listings(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.listings
    ADD CONSTRAINT listings_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.categories(id) ON UPDATE CASCADE ON DELETE SET NULL;
ALTER TABLE ONLY public.storefront_events
    ADD CONSTRAINT storefront_events_storefront_id_fkey FOREIGN KEY (storefront_id) REFERENCES public.storefronts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.storefront_staff
    ADD CONSTRAINT storefront_staff_invitation_id_fkey FOREIGN KEY (invitation_id) REFERENCES public.storefront_invitations(id) ON DELETE SET NULL;
ALTER TABLE ONLY public.variant_attribute_values
    ADD CONSTRAINT variant_attribute_values_attribute_id_fkey FOREIGN KEY (attribute_id) REFERENCES public.attributes(id) ON DELETE CASCADE;
COMMENT ON SCHEMA public IS '';
COMMENT ON EXTENSION cube IS 'data type for multidimensional cubes';
COMMENT ON EXTENSION earthdistance IS 'calculate great-circle distances on the surface of the Earth';
COMMENT ON EXTENSION pg_trgm IS 'text similarity measurement and index searching based on trigrams';
COMMENT ON FUNCTION public.archive_old_analytics_events() IS 'Archive events older than 90 days to manage table size';
COMMENT ON FUNCTION public.cleanup_expired_reservations() IS 'Cleanup expired stock reservations (run as cron job)';
COMMENT ON FUNCTION public.get_category_attributes_with_inheritance(p_category_id integer, p_locale character varying) IS 'Get all attributes for a category including inherited from parents';
COMMENT ON FUNCTION public.get_file_type_from_content_type(content_type text) IS 'Determine file_type from MIME content_type';
COMMENT ON FUNCTION public.is_valid_file_size(file_type text, file_size bigint) IS 'Check if file size is within allowed limits for the file type';
COMMENT ON FUNCTION public.log_analytics_event(p_event_type character varying, p_entity_type character varying, p_entity_id bigint, p_user_id bigint, p_session_id character varying, p_metadata jsonb) IS 'Helper function to log analytics events with validation';
COMMENT ON FUNCTION public.refresh_analytics_trending_cache() IS 'Refresh the trending analytics cache. Run this hourly via cron job or trigger.';
COMMENT ON FUNCTION public.refresh_analytics_views() IS 'Refresh all materialized views concurrently';
COMMENT ON FUNCTION public.update_chat_last_message_at() IS 'Update parent chat last_message_at when new message is inserted';
COMMENT ON FUNCTION public.update_chats_updated_at() IS 'Automatically update updated_at timestamp when chat is modified';
COMMENT ON FUNCTION public.update_message_attachments_count() IS 'Update parent message attachments_count and has_attachments when attachments are added/removed';
COMMENT ON FUNCTION public.update_messages_read_at() IS 'Automatically set read_at timestamp and update status when message is marked as read';
COMMENT ON FUNCTION public.update_messages_updated_at() IS 'Automatically update updated_at timestamp when message is modified';
COMMENT ON TABLE public.listings IS 'Unified listings table (C2C + B2C merged). Legacy tables dropped in Phase 11.5 (2025-11-06)';
COMMENT ON COLUMN public.listings.view_count IS 'Number of times this listing has been viewed (renamed from views_count for b2c compatibility)';
COMMENT ON COLUMN public.listings.source_type IS 'Type of listing: c2c (consumer-to-consumer) or b2c (business-to-consumer)';

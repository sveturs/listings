CREATE INDEX idx_variants_stock_status ON public.product_variants USING btree (product_id, status, stock_quantity) WHERE ((status)::text = 'active'::text);
CREATE INDEX listing_favorites_user_id_created_at_idx ON public.listing_favorites USING btree (user_id, created_at DESC);
ALTER TABLE ONLY public.attribute_options
    ADD CONSTRAINT attribute_options_attribute_id_option_value_key UNIQUE (attribute_id, option_value);
ALTER TABLE ONLY public.attribute_options
    ADD CONSTRAINT attribute_options_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.attribute_search_cache
    ADD CONSTRAINT attribute_search_cache_listing_id_key UNIQUE (listing_id);
ALTER TABLE ONLY public.attribute_search_cache
    ADD CONSTRAINT attribute_search_cache_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.attribute_values
    ADD CONSTRAINT attribute_values_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.attributes
    ADD CONSTRAINT attributes_code_key UNIQUE (code);
ALTER TABLE ONLY public.attributes
    ADD CONSTRAINT attributes_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.brand_category_mapping
    ADD CONSTRAINT brand_category_mapping_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.c2c_chats
    ADD CONSTRAINT c2c_chats_listing_id_buyer_id_seller_id_key UNIQUE (listing_id, buyer_id, seller_id);
ALTER TABLE ONLY public.c2c_chats
    ADD CONSTRAINT c2c_chats_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.c2c_chats
    ADD CONSTRAINT c2c_chats_storefront_product_id_buyer_id_seller_id_key UNIQUE (storefront_product_id, buyer_id, seller_id);
ALTER TABLE ONLY public.c2c_messages
    ADD CONSTRAINT c2c_messages_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.cart_items
    ADD CONSTRAINT cart_items_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_slug_key UNIQUE (slug);
ALTER TABLE ONLY public.category_ai_mapping
    ADD CONSTRAINT category_ai_mapping_ai_category_name_key UNIQUE (ai_category_name);
ALTER TABLE ONLY public.category_ai_mapping
    ADD CONSTRAINT category_ai_mapping_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.category_attributes
    ADD CONSTRAINT category_attributes_category_id_attribute_id_key UNIQUE (category_id, attribute_id);
ALTER TABLE ONLY public.category_attributes
    ADD CONSTRAINT category_attributes_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.category_detections
    ADD CONSTRAINT category_detections_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.category_variant_attributes
    ADD CONSTRAINT category_variant_attributes_category_id_attribute_id_key UNIQUE (category_id, attribute_id);
ALTER TABLE ONLY public.category_variant_attributes
    ADD CONSTRAINT category_variant_attributes_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.chat_attachments
    ADD CONSTRAINT chat_attachments_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.chats
    ADD CONSTRAINT chats_listing_participants_unique UNIQUE (listing_id, buyer_id, seller_id);
ALTER TABLE ONLY public.chats
    ADD CONSTRAINT chats_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.chats
    ADD CONSTRAINT chats_product_participants_unique UNIQUE (storefront_product_id, buyer_id, seller_id);
ALTER TABLE ONLY public.indexing_queue
    ADD CONSTRAINT indexing_queue_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.inventory_movements
    ADD CONSTRAINT inventory_movements_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.inventory_reservations
    ADD CONSTRAINT inventory_reservations_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.listing_attribute_values
    ADD CONSTRAINT listing_attribute_values_listing_id_attribute_id_key UNIQUE (listing_id, attribute_id);
ALTER TABLE ONLY public.listing_attribute_values
    ADD CONSTRAINT listing_attribute_values_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.listing_attributes
    ADD CONSTRAINT listing_attributes_listing_id_attribute_key_key UNIQUE (listing_id, attribute_key);
ALTER TABLE ONLY public.listing_attributes
    ADD CONSTRAINT listing_attributes_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.listing_favorites
    ADD CONSTRAINT listing_favorites_pkey PRIMARY KEY (user_id, listing_id);
ALTER TABLE ONLY public.listing_images
    ADD CONSTRAINT listing_images_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.listing_locations
    ADD CONSTRAINT listing_locations_listing_id_key UNIQUE (listing_id);
ALTER TABLE ONLY public.listing_locations
    ADD CONSTRAINT listing_locations_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.listing_stats
    ADD CONSTRAINT listing_stats_pkey PRIMARY KEY (listing_id);
ALTER TABLE ONLY public.listing_tags
    ADD CONSTRAINT listing_tags_listing_id_tag_key UNIQUE (listing_id, tag);
ALTER TABLE ONLY public.listing_tags
    ADD CONSTRAINT listing_tags_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.listing_variants
    ADD CONSTRAINT listing_variants_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.listings
    ADD CONSTRAINT listings_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.listings
    ADD CONSTRAINT listings_uuid_key UNIQUE (uuid);
ALTER TABLE ONLY public.messages
    ADD CONSTRAINT messages_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT order_items_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_order_number_key UNIQUE (order_number);
ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.product_variants
    ADD CONSTRAINT product_variants_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.search_queries
    ADD CONSTRAINT search_queries_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.shopping_carts
    ADD CONSTRAINT shopping_carts_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.stock_reservations
    ADD CONSTRAINT stock_reservations_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.storefront_delivery_options
    ADD CONSTRAINT storefront_delivery_options_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.storefront_events
    ADD CONSTRAINT storefront_events_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.storefront_hours
    ADD CONSTRAINT storefront_hours_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.storefront_invitations
    ADD CONSTRAINT storefront_invitations_invite_code_key UNIQUE (invite_code);
ALTER TABLE ONLY public.storefront_invitations
    ADD CONSTRAINT storefront_invitations_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.storefront_payment_methods
    ADD CONSTRAINT storefront_payment_methods_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.storefront_staff
    ADD CONSTRAINT storefront_staff_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.storefronts
    ADD CONSTRAINT storefronts_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.storefronts
    ADD CONSTRAINT storefronts_slug_key UNIQUE (slug);
ALTER TABLE ONLY public.brand_category_mapping
    ADD CONSTRAINT unique_brand_category UNIQUE (brand_name, category_slug);
ALTER TABLE ONLY public.product_variants
    ADD CONSTRAINT unique_variant_sku UNIQUE (sku);
ALTER TABLE ONLY public.attribute_values
    ADD CONSTRAINT uq_attribute_value UNIQUE (attribute_id, value);
ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);
ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.variant_attribute_values
    ADD CONSTRAINT variant_attribute_values_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.variant_attribute_values
    ADD CONSTRAINT variant_attribute_values_variant_id_attribute_id_key UNIQUE (variant_id, attribute_id);
CREATE TRIGGER trg_ai_mapping_updated_at BEFORE UPDATE ON public.category_ai_mapping FOR EACH ROW EXECUTE FUNCTION public.update_ai_mapping_updated_at();
CREATE TRIGGER trg_categories_updated_at BEFORE UPDATE ON public.categories FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
CREATE TRIGGER trg_check_category_level BEFORE INSERT OR UPDATE ON public.categories FOR EACH ROW EXECUTE FUNCTION public.check_category_level_constraint();
CREATE TRIGGER trigger_attribute_options_updated_at BEFORE UPDATE ON public.attribute_options FOR EACH ROW EXECUTE FUNCTION public.update_attributes_timestamp();
CREATE TRIGGER trigger_attribute_values_updated_at BEFORE UPDATE ON public.attribute_values FOR EACH ROW EXECUTE FUNCTION public.update_attribute_values_updated_at();
CREATE TRIGGER trigger_attributes_search_vector BEFORE INSERT OR UPDATE OF name, code ON public.attributes FOR EACH ROW EXECUTE FUNCTION public.update_attributes_search_vector();
CREATE TRIGGER trigger_attributes_updated_at BEFORE UPDATE ON public.attributes FOR EACH ROW EXECUTE FUNCTION public.update_attributes_timestamp();
CREATE TRIGGER trigger_cart_items_updated_at BEFORE UPDATE ON public.cart_items FOR EACH ROW EXECUTE FUNCTION public.update_cart_items_updated_at();
CREATE TRIGGER trigger_category_attributes_updated_at BEFORE UPDATE ON public.category_attributes FOR EACH ROW EXECUTE FUNCTION public.update_attributes_timestamp();
CREATE TRIGGER trigger_category_variant_attrs_updated_at BEFORE UPDATE ON public.category_variant_attributes FOR EACH ROW EXECUTE FUNCTION public.update_attributes_timestamp();
CREATE TRIGGER trigger_chats_updated_at BEFORE UPDATE ON public.chats FOR EACH ROW EXECUTE FUNCTION public.update_chats_updated_at();
CREATE TRIGGER trigger_generate_slug BEFORE INSERT OR UPDATE ON public.listings FOR EACH ROW EXECUTE FUNCTION public.generate_slug_from_title();
CREATE TRIGGER trigger_inventory_reservations_updated_at BEFORE UPDATE ON public.inventory_reservations FOR EACH ROW EXECUTE FUNCTION public.update_inventory_reservations_updated_at();
CREATE TRIGGER trigger_listing_attr_values_updated_at BEFORE UPDATE ON public.listing_attribute_values FOR EACH ROW EXECUTE FUNCTION public.update_attributes_timestamp();
CREATE TRIGGER trigger_listing_variants_updated_at BEFORE UPDATE ON public.listing_variants FOR EACH ROW EXECUTE FUNCTION public.update_listing_variants_updated_at();
CREATE TRIGGER trigger_messages_read_at BEFORE UPDATE ON public.messages FOR EACH ROW WHEN ((new.is_read IS DISTINCT FROM old.is_read)) EXECUTE FUNCTION public.update_messages_read_at();
CREATE TRIGGER trigger_messages_updated_at BEFORE UPDATE ON public.messages FOR EACH ROW EXECUTE FUNCTION public.update_messages_updated_at();
CREATE TRIGGER trigger_orders_updated_at BEFORE UPDATE ON public.orders FOR EACH ROW EXECUTE FUNCTION public.update_orders_updated_at();
CREATE TRIGGER trigger_reservations_auto_expire BEFORE UPDATE ON public.stock_reservations FOR EACH ROW WHEN (((new.status)::text = 'active'::text)) EXECUTE FUNCTION public.auto_expire_reservations();
CREATE TRIGGER trigger_reservations_sync_quantity AFTER INSERT OR DELETE OR UPDATE ON public.stock_reservations FOR EACH ROW EXECUTE FUNCTION public.sync_variant_reserved_quantity();
CREATE TRIGGER trigger_reservations_updated_at BEFORE UPDATE ON public.stock_reservations FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
CREATE TRIGGER trigger_shopping_carts_updated_at BEFORE UPDATE ON public.shopping_carts FOR EACH ROW EXECUTE FUNCTION public.update_shopping_carts_updated_at();
CREATE TRIGGER trigger_storefront_delivery_options_updated_at BEFORE UPDATE ON public.storefront_delivery_options FOR EACH ROW EXECUTE FUNCTION public.update_storefront_delivery_options_updated_at();
CREATE TRIGGER trigger_storefront_invitations_updated_at BEFORE UPDATE ON public.storefront_invitations FOR EACH ROW EXECUTE FUNCTION public.update_storefront_invitations_updated_at();
CREATE TRIGGER trigger_storefront_staff_updated_at BEFORE UPDATE ON public.storefront_staff FOR EACH ROW EXECUTE FUNCTION public.update_storefront_staff_updated_at();
CREATE TRIGGER trigger_storefronts_updated_at BEFORE UPDATE ON public.storefronts FOR EACH ROW EXECUTE FUNCTION public.update_storefronts_updated_at();
CREATE TRIGGER trigger_update_chat_last_message AFTER INSERT ON public.messages FOR EACH ROW EXECUTE FUNCTION public.update_chat_last_message_at();
CREATE TRIGGER trigger_update_message_attachments_count_delete AFTER DELETE ON public.chat_attachments FOR EACH ROW EXECUTE FUNCTION public.update_message_attachments_count();
CREATE TRIGGER trigger_update_message_attachments_count_insert AFTER INSERT ON public.chat_attachments FOR EACH ROW EXECUTE FUNCTION public.update_message_attachments_count();
CREATE TRIGGER trigger_validate_variant_attr_value BEFORE INSERT OR UPDATE ON public.variant_attribute_values FOR EACH ROW EXECUTE FUNCTION public.validate_variant_attribute_value();
CREATE TRIGGER trigger_variant_attr_values_updated_at BEFORE UPDATE ON public.variant_attribute_values FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();

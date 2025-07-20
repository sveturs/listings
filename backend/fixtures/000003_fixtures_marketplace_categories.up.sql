-- Fixtures for table: marketplace_categories

INSERT INTO public.marketplace_categories (id, name, slug, parent_id, icon, created_at, has_custom_ui, custom_ui_component, sort_order, level, count, external_id, description, is_active) VALUES ('10411', 'Electronics', 'electronics', NULL, '', '2025-07-19 12:07:52.357167', 'f', '', '0', '0', '0', NULL, 'Electronic devices, gadgets, and accessories', 't');

INSERT INTO public.marketplace_categories (id, name, slug, parent_id, icon, created_at, has_custom_ui, custom_ui_component, sort_order, level, count, external_id, description, is_active) VALUES ('10412', 'homes', 'homes', NULL, '🏠', '2025-07-19 12:52:43.020127', 'f', '', '0', '0', '0', NULL, '', 't');

INSERT INTO public.marketplace_categories (id, name, slug, parent_id, icon, created_at, has_custom_ui, custom_ui_component, sort_order, level, count, external_id, description, is_active) VALUES ('10414', 'rooms', 'rooms', '10412', '🚪', '2025-07-19 14:25:09.365161', 'f', '', '0', '0', '0', NULL, '', 't');

INSERT INTO public.marketplace_categories (id, name, slug, parent_id, icon, created_at, has_custom_ui, custom_ui_component, sort_order, level, count, external_id, description, is_active) VALUES ('10413', 'phones', 'phones', '10411', '📱', '2025-07-19 12:53:42.183188', 'f', '', '0', '0', '0', NULL, '', 't');
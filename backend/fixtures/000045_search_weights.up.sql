INSERT INTO public.search_weights (id, field_name, weight, search_type, item_type, category_id, description, is_active, version, created_at, updated_at, created_by, updated_by) VALUES (1, 'title', 0.9, 'fulltext', 'global', NULL, 'Название/заголовок - самое важное поле для поиска', true, 1, '2025-07-08 15:29:51.062386+00', '2025-07-08 15:29:51.062386+00', NULL, NULL);
INSERT INTO public.search_weights (id, field_name, weight, search_type, item_type, category_id, description, is_active, version, created_at, updated_at, created_by, updated_by) VALUES (2, 'title', 0.85, 'fuzzy', 'global', NULL, 'Название/заголовок для нечеткого поиска', true, 1, '2025-07-08 15:29:51.062386+00', '2025-07-08 15:29:51.062386+00', NULL, NULL);
INSERT INTO public.search_weights (id, field_name, weight, search_type, item_type, category_id, description, is_active, version, created_at, updated_at, created_by, updated_by) VALUES (3, 'description', 0.7, 'fulltext', 'global', NULL, 'Описание - второе по важности поле', true, 1, '2025-07-08 15:29:51.062386+00', '2025-07-08 15:29:51.062386+00', NULL, NULL);
INSERT INTO public.search_weights (id, field_name, weight, search_type, item_type, category_id, description, is_active, version, created_at, updated_at, created_by, updated_by) VALUES (4, 'description', 0.6, 'fuzzy', 'global', NULL, 'Описание для нечеткого поиска', true, 1, '2025-07-08 15:29:51.062386+00', '2025-07-08 15:29:51.062386+00', NULL, NULL);
INSERT INTO public.search_weights (id, field_name, weight, search_type, item_type, category_id, description, is_active, version, created_at, updated_at, created_by, updated_by) VALUES (5, 'category', 0.8, 'exact', 'global', NULL, 'Категория товара/услуги', true, 1, '2025-07-08 15:29:51.062386+00', '2025-07-08 15:29:51.062386+00', NULL, NULL);
INSERT INTO public.search_weights (id, field_name, weight, search_type, item_type, category_id, description, is_active, version, created_at, updated_at, created_by, updated_by) VALUES (6, 'location', 0.75, 'exact', 'global', NULL, 'Местоположение', true, 1, '2025-07-08 15:29:51.062386+00', '2025-07-08 15:29:51.062386+00', NULL, NULL);
INSERT INTO public.search_weights (id, field_name, weight, search_type, item_type, category_id, description, is_active, version, created_at, updated_at, created_by, updated_by) VALUES (7, 'price', 0.5, 'exact', 'global', NULL, 'Цена (важна для фильтрации)', true, 1, '2025-07-08 15:29:51.062386+00', '2025-07-08 15:29:51.062386+00', NULL, NULL);
INSERT INTO public.search_weights (id, field_name, weight, search_type, item_type, category_id, description, is_active, version, created_at, updated_at, created_by, updated_by) VALUES (8, 'brand', 0.6, 'exact', 'global', NULL, 'Бренд/производитель', true, 1, '2025-07-08 15:29:51.062386+00', '2025-07-08 15:29:51.062386+00', NULL, NULL);
INSERT INTO public.search_weights (id, field_name, weight, search_type, item_type, category_id, description, is_active, version, created_at, updated_at, created_by, updated_by) VALUES (9, 'tags', 0.4, 'fulltext', 'global', NULL, 'Теги и ключевые слова', true, 1, '2025-07-08 15:29:51.062386+00', '2025-07-08 15:29:51.062386+00', NULL, NULL);
INSERT INTO public.search_weights (id, field_name, weight, search_type, item_type, category_id, description, is_active, version, created_at, updated_at, created_by, updated_by) VALUES (10, 'title', 0.95, 'fulltext', 'marketplace', NULL, 'Название объявления в маркетплейсе', true, 1, '2025-07-08 15:29:51.062386+00', '2025-07-08 15:29:51.062386+00', NULL, NULL);
INSERT INTO public.search_weights (id, field_name, weight, search_type, item_type, category_id, description, is_active, version, created_at, updated_at, created_by, updated_by) VALUES (11, 'user_name', 0.3, 'fuzzy', 'marketplace', NULL, 'Имя продавца', true, 1, '2025-07-08 15:29:51.062386+00', '2025-07-08 15:29:51.062386+00', NULL, NULL);
INSERT INTO public.search_weights (id, field_name, weight, search_type, item_type, category_id, description, is_active, version, created_at, updated_at, created_by, updated_by) VALUES (12, 'condition', 0.4, 'exact', 'marketplace', NULL, 'Состояние товара', true, 1, '2025-07-08 15:29:51.062386+00', '2025-07-08 15:29:51.062386+00', NULL, NULL);
INSERT INTO public.search_weights (id, field_name, weight, search_type, item_type, category_id, description, is_active, version, created_at, updated_at, created_by, updated_by) VALUES (13, 'title', 0.9, 'fulltext', 'storefront', NULL, 'Название товара в магазине', true, 1, '2025-07-08 15:29:51.062386+00', '2025-07-08 15:29:51.062386+00', NULL, NULL);
INSERT INTO public.search_weights (id, field_name, weight, search_type, item_type, category_id, description, is_active, version, created_at, updated_at, created_by, updated_by) VALUES (14, 'store_name', 0.5, 'fuzzy', 'storefront', NULL, 'Название магазина', true, 1, '2025-07-08 15:29:51.062386+00', '2025-07-08 15:29:51.062386+00', NULL, NULL);
INSERT INTO public.search_weights (id, field_name, weight, search_type, item_type, category_id, description, is_active, version, created_at, updated_at, created_by, updated_by) VALUES (15, 'sku', 0.6, 'exact', 'storefront', NULL, 'Артикул товара', true, 1, '2025-07-08 15:29:51.062386+00', '2025-07-08 15:29:51.062386+00', NULL, NULL);
INSERT INTO public.search_weights (id, field_name, weight, search_type, item_type, category_id, description, is_active, version, created_at, updated_at, created_by, updated_by) VALUES (16, 'availability', 0.3, 'exact', 'storefront', NULL, 'Наличие товара', true, 1, '2025-07-08 15:29:51.062386+00', '2025-07-08 15:29:51.062386+00', NULL, NULL);


--
-- Data for Name: search_weights_history; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: shopping_cart_items; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: shopping_carts; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: spatial_ref_sys; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: storefront_delivery_options; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: storefront_favorites; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: storefront_hours; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: storefront_inventory_movements; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: storefront_order_items; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: storefront_orders; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: storefront_payment_methods; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: storefront_product_attributes; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: storefront_product_images; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: storefront_product_variant_images; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: storefront_product_variants; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: storefront_products; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: storefront_staff; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: storefronts; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: subscription_history; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: subscription_payments; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: subscription_plans; Type: TABLE DATA; Schema: public; Owner: -
--

-- Fixtures for table: search_weights

INSERT INTO public.search_weights (id, field_name, weight, search_type, item_type, category_id, description, is_active, version, created_at, updated_at, created_by, updated_by) VALUES ('1', 'title', '0.9', 'fulltext', 'global', NULL, 'Название/заголовок - самое важное поле для поиска', 't', '1', '2025-07-08 15:29:51.062386+00', '2025-07-08 15:29:51.062386+00', NULL, NULL);

INSERT INTO public.search_weights (id, field_name, weight, search_type, item_type, category_id, description, is_active, version, created_at, updated_at, created_by, updated_by) VALUES ('2', 'title', '0.85', 'fuzzy', 'global', NULL, 'Название/заголовок для нечеткого поиска', 't', '1', '2025-07-08 15:29:51.062386+00', '2025-07-08 15:29:51.062386+00', NULL, NULL);

INSERT INTO public.search_weights (id, field_name, weight, search_type, item_type, category_id, description, is_active, version, created_at, updated_at, created_by, updated_by) VALUES ('3', 'description', '0.7', 'fulltext', 'global', NULL, 'Описание - второе по важности поле', 't', '1', '2025-07-08 15:29:51.062386+00', '2025-07-08 15:29:51.062386+00', NULL, NULL);

INSERT INTO public.search_weights (id, field_name, weight, search_type, item_type, category_id, description, is_active, version, created_at, updated_at, created_by, updated_by) VALUES ('4', 'description', '0.6', 'fuzzy', 'global', NULL, 'Описание для нечеткого поиска', 't', '1', '2025-07-08 15:29:51.062386+00', '2025-07-08 15:29:51.062386+00', NULL, NULL);

INSERT INTO public.search_weights (id, field_name, weight, search_type, item_type, category_id, description, is_active, version, created_at, updated_at, created_by, updated_by) VALUES ('5', 'category', '0.8', 'exact', 'global', NULL, 'Категория товара/услуги', 't', '1', '2025-07-08 15:29:51.062386+00', '2025-07-08 15:29:51.062386+00', NULL, NULL);

INSERT INTO public.search_weights (id, field_name, weight, search_type, item_type, category_id, description, is_active, version, created_at, updated_at, created_by, updated_by) VALUES ('6', 'location', '0.75', 'exact', 'global', NULL, 'Местоположение', 't', '1', '2025-07-08 15:29:51.062386+00', '2025-07-08 15:29:51.062386+00', NULL, NULL);

INSERT INTO public.search_weights (id, field_name, weight, search_type, item_type, category_id, description, is_active, version, created_at, updated_at, created_by, updated_by) VALUES ('7', 'price', '0.5', 'exact', 'global', NULL, 'Цена (важна для фильтрации)', 't', '1', '2025-07-08 15:29:51.062386+00', '2025-07-08 15:29:51.062386+00', NULL, NULL);

INSERT INTO public.search_weights (id, field_name, weight, search_type, item_type, category_id, description, is_active, version, created_at, updated_at, created_by, updated_by) VALUES ('8', 'brand', '0.6', 'exact', 'global', NULL, 'Бренд/производитель', 't', '1', '2025-07-08 15:29:51.062386+00', '2025-07-08 15:29:51.062386+00', NULL, NULL);

INSERT INTO public.search_weights (id, field_name, weight, search_type, item_type, category_id, description, is_active, version, created_at, updated_at, created_by, updated_by) VALUES ('9', 'tags', '0.4', 'fulltext', 'global', NULL, 'Теги и ключевые слова', 't', '1', '2025-07-08 15:29:51.062386+00', '2025-07-08 15:29:51.062386+00', NULL, NULL);

INSERT INTO public.search_weights (id, field_name, weight, search_type, item_type, category_id, description, is_active, version, created_at, updated_at, created_by, updated_by) VALUES ('10', 'title', '0.95', 'fulltext', 'marketplace', NULL, 'Название объявления в маркетплейсе', 't', '1', '2025-07-08 15:29:51.062386+00', '2025-07-08 15:29:51.062386+00', NULL, NULL);

INSERT INTO public.search_weights (id, field_name, weight, search_type, item_type, category_id, description, is_active, version, created_at, updated_at, created_by, updated_by) VALUES ('11', 'user_name', '0.3', 'fuzzy', 'marketplace', NULL, 'Имя продавца', 't', '1', '2025-07-08 15:29:51.062386+00', '2025-07-08 15:29:51.062386+00', NULL, NULL);

INSERT INTO public.search_weights (id, field_name, weight, search_type, item_type, category_id, description, is_active, version, created_at, updated_at, created_by, updated_by) VALUES ('12', 'condition', '0.4', 'exact', 'marketplace', NULL, 'Состояние товара', 't', '1', '2025-07-08 15:29:51.062386+00', '2025-07-08 15:29:51.062386+00', NULL, NULL);

INSERT INTO public.search_weights (id, field_name, weight, search_type, item_type, category_id, description, is_active, version, created_at, updated_at, created_by, updated_by) VALUES ('13', 'title', '0.9', 'fulltext', 'storefront', NULL, 'Название товара в магазине', 't', '1', '2025-07-08 15:29:51.062386+00', '2025-07-08 15:29:51.062386+00', NULL, NULL);

INSERT INTO public.search_weights (id, field_name, weight, search_type, item_type, category_id, description, is_active, version, created_at, updated_at, created_by, updated_by) VALUES ('14', 'store_name', '0.5', 'fuzzy', 'storefront', NULL, 'Название магазина', 't', '1', '2025-07-08 15:29:51.062386+00', '2025-07-08 15:29:51.062386+00', NULL, NULL);

INSERT INTO public.search_weights (id, field_name, weight, search_type, item_type, category_id, description, is_active, version, created_at, updated_at, created_by, updated_by) VALUES ('15', 'sku', '0.6', 'exact', 'storefront', NULL, 'Артикул товара', 't', '1', '2025-07-08 15:29:51.062386+00', '2025-07-08 15:29:51.062386+00', NULL, NULL);

INSERT INTO public.search_weights (id, field_name, weight, search_type, item_type, category_id, description, is_active, version, created_at, updated_at, created_by, updated_by) VALUES ('16', 'availability', '0.3', 'exact', 'storefront', NULL, 'Наличие товара', 't', '1', '2025-07-08 15:29:51.062386+00', '2025-07-08 15:29:51.062386+00', NULL, NULL);
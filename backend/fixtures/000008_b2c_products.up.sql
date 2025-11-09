-- Migrated from b2c_products to unified listings table
-- Phase 11.5: Fixture migration to unified schema
-- All B2C products now use listings table with source_type='b2c'

INSERT INTO public.listings (id, storefront_id, user_id, title, description, price, currency, category_id, sku, quantity, status, visibility, views_count, favorites_count, source_type, created_at, updated_at, published_at, deleted_at, is_deleted) VALUES (1073, 43, 1, 'Baterija za LG B2050 950 mAh', '<p>Standard zamenska baterija za LG.</p><p>Garancija: 6 meseci</p>', 590.00, 'RSD', 1001, '3023', 1, 'active', 'public', 0, 0, 'b2c', '2025-10-11 18:45:12.333036+00', '2025-10-11 18:45:57.904856+00', '2025-10-11 18:45:12.333036+00', NULL, false);

INSERT INTO public.listings (id, storefront_id, user_id, title, description, price, currency, category_id, sku, quantity, status, visibility, views_count, favorites_count, source_type, created_at, updated_at, published_at, deleted_at, is_deleted) VALUES (1074, 43, 1, 'Baterija za LG KU800 900 mAh.', '<p>Standard zamenska baterija za LG.</p><p>Garancija: 6 meseci</p>', 390.00, 'RSD', 1001, '3037', 1, 'active', 'public', 0, 0, 'b2c', '2025-10-11 18:45:12.333036+00', '2025-10-11 18:45:59.005523+00', '2025-10-11 18:45:12.333036+00', NULL, false);

INSERT INTO public.listings (id, storefront_id, user_id, title, description, price, currency, category_id, sku, quantity, status, visibility, views_count, favorites_count, source_type, created_at, updated_at, published_at, deleted_at, is_deleted) VALUES (1075, 43, 1, 'Baterija za Mot E1000 1100 mAh.', '<p>Standard zamenska baterija za Motorolu.</p><p>Garancija: 6 meseci</p>', 390.00, 'RSD', 1001, '3045', 1, 'active', 'public', 0, 0, 'b2c', '2025-10-11 18:45:12.333036+00', '2025-10-11 18:46:00.233717+00', '2025-10-11 18:45:12.333036+00', NULL, false);

INSERT INTO public.listings (id, storefront_id, user_id, title, description, price, currency, category_id, sku, quantity, status, visibility, views_count, favorites_count, source_type, created_at, updated_at, published_at, deleted_at, is_deleted) VALUES (1076, 43, 1, 'Baterija za Nokia BL-5F (N95/E65) 900 mAh', '<p>Standard zamenska baterija za Nokiu.</p><p>Garancija: 6 meseci</p>', 490.00, 'RSD', 1001, '3058', 1, 'active', 'public', 24, 0, 'b2c', '2025-10-11 18:45:12.333036+00', '2025-10-11 18:45:55.931739+00', '2025-10-11 18:45:12.333036+00', NULL, false);

INSERT INTO public.listings (id, storefront_id, user_id, title, description, price, currency, category_id, sku, quantity, status, visibility, views_count, favorites_count, source_type, created_at, updated_at, published_at, deleted_at, is_deleted) VALUES (1071, 43, 1, 'МФУ Canon G3420', 'Многофункциональный принтер Canon G3420 серого цвета. Оснащен беспроводным подключением, ЖК-дисплеем, лотком для бумаги и функцией сканирования. Находится в хорошем рабочем состоянии с минимальными следами использования.', 15000.00, 'RSD', 1106, NULL, 1, 'active', 'public', 6, 0, 'b2c', '2025-10-11 17:40:08.186825+00', '2025-10-11 17:40:08.186825+00', '2025-10-11 17:40:08.186825+00', NULL, false);

INSERT INTO public.listings (id, storefront_id, user_id, title, description, price, currency, category_id, sku, quantity, status, visibility, views_count, favorites_count, source_type, created_at, updated_at, published_at, deleted_at, is_deleted) VALUES (1072, 43, 1, 'Baterija za Nokia BL-6F (N95 8GB) 1000 mAh', '<p>Standard zamenska baterija za Nokiu.</p><p>Garancija: 6 meseci</p>', 390.00, 'RSD', 1001, '3062', 1, 'active', 'public', 16, 0, 'b2c', '2025-10-11 18:45:12.333036+00', '2025-10-11 18:45:56.8797+00', '2025-10-11 18:45:12.333036+00', NULL, false);

-- Migrate JSONB attributes to listing_attributes table
-- Listing 1073 attributes
INSERT INTO public.listing_attributes (listing_id, attribute_key, attribute_value, created_at) VALUES (1073, 'uvoznik', 'Digital Vision doo', '2025-10-11 18:45:12.333036+00');
INSERT INTO public.listing_attributes (listing_id, attribute_key, attribute_value, created_at) VALUES (1073, 'kategorija1', 'OPREMA ZA MOBILNI', '2025-10-11 18:45:12.333036+00');
INSERT INTO public.listing_attributes (listing_id, attribute_key, attribute_value, created_at) VALUES (1073, 'kategorija2', 'BATERIJE', '2025-10-11 18:45:12.333036+00');
INSERT INTO public.listing_attributes (listing_id, attribute_key, attribute_value, created_at) VALUES (1073, 'kategorija3', 'BATERIJE OUTLET', '2025-10-11 18:45:12.333036+00');
INSERT INTO public.listing_attributes (listing_id, attribute_key, attribute_value, created_at) VALUES (1073, 'godina_uvoza', '2025.', '2025-10-11 18:45:12.333036+00');
INSERT INTO public.listing_attributes (listing_id, attribute_key, attribute_value, created_at) VALUES (1073, 'zemlja_porekla', 'Kina', '2025-10-11 18:45:12.333036+00');

-- Listing 1074 attributes
INSERT INTO public.listing_attributes (listing_id, attribute_key, attribute_value, created_at) VALUES (1074, 'uvoznik', 'Digital Vision doo', '2025-10-11 18:45:12.333036+00');
INSERT INTO public.listing_attributes (listing_id, attribute_key, attribute_value, created_at) VALUES (1074, 'kategorija1', 'OPREMA ZA MOBILNI', '2025-10-11 18:45:12.333036+00');
INSERT INTO public.listing_attributes (listing_id, attribute_key, attribute_value, created_at) VALUES (1074, 'kategorija2', 'BATERIJE', '2025-10-11 18:45:12.333036+00');
INSERT INTO public.listing_attributes (listing_id, attribute_key, attribute_value, created_at) VALUES (1074, 'kategorija3', 'BATERIJE OUTLET', '2025-10-11 18:45:12.333036+00');
INSERT INTO public.listing_attributes (listing_id, attribute_key, attribute_value, created_at) VALUES (1074, 'godina_uvoza', '2025.', '2025-10-11 18:45:12.333036+00');
INSERT INTO public.listing_attributes (listing_id, attribute_key, attribute_value, created_at) VALUES (1074, 'zemlja_porekla', 'Kina', '2025-10-11 18:45:12.333036+00');

-- Listing 1075 attributes
INSERT INTO public.listing_attributes (listing_id, attribute_key, attribute_value, created_at) VALUES (1075, 'uvoznik', 'Digital Vision doo', '2025-10-11 18:45:12.333036+00');
INSERT INTO public.listing_attributes (listing_id, attribute_key, attribute_value, created_at) VALUES (1075, 'kategorija1', 'OPREMA ZA MOBILNI', '2025-10-11 18:45:12.333036+00');
INSERT INTO public.listing_attributes (listing_id, attribute_key, attribute_value, created_at) VALUES (1075, 'kategorija2', 'BATERIJE', '2025-10-11 18:45:12.333036+00');
INSERT INTO public.listing_attributes (listing_id, attribute_key, attribute_value, created_at) VALUES (1075, 'kategorija3', 'BATERIJE OUTLET', '2025-10-11 18:45:12.333036+00');
INSERT INTO public.listing_attributes (listing_id, attribute_key, attribute_value, created_at) VALUES (1075, 'godina_uvoza', '2025.', '2025-10-11 18:45:12.333036+00');
INSERT INTO public.listing_attributes (listing_id, attribute_key, attribute_value, created_at) VALUES (1075, 'zemlja_porekla', 'Kina', '2025-10-11 18:45:12.333036+00');

-- Listing 1076 attributes
INSERT INTO public.listing_attributes (listing_id, attribute_key, attribute_value, created_at) VALUES (1076, 'uvoznik', 'Digital Vision doo', '2025-10-11 18:45:12.333036+00');
INSERT INTO public.listing_attributes (listing_id, attribute_key, attribute_value, created_at) VALUES (1076, 'kategorija1', 'OPREMA ZA MOBILNI', '2025-10-11 18:45:12.333036+00');
INSERT INTO public.listing_attributes (listing_id, attribute_key, attribute_value, created_at) VALUES (1076, 'kategorija2', 'BATERIJE', '2025-10-11 18:45:12.333036+00');
INSERT INTO public.listing_attributes (listing_id, attribute_key, attribute_value, created_at) VALUES (1076, 'kategorija3', 'BATERIJE OUTLET', '2025-10-11 18:45:12.333036+00');
INSERT INTO public.listing_attributes (listing_id, attribute_key, attribute_value, created_at) VALUES (1076, 'godina_uvoza', '2025.', '2025-10-11 18:45:12.333036+00');
INSERT INTO public.listing_attributes (listing_id, attribute_key, attribute_value, created_at) VALUES (1076, 'zemlja_porekla', 'Kina', '2025-10-11 18:45:12.333036+00');

-- Listing 1072 attributes
INSERT INTO public.listing_attributes (listing_id, attribute_key, attribute_value, created_at) VALUES (1072, 'uvoznik', 'Digital Vision doo', '2025-10-11 18:45:12.333036+00');
INSERT INTO public.listing_attributes (listing_id, attribute_key, attribute_value, created_at) VALUES (1072, 'kategorija1', 'OPREMA ZA MOBILNI', '2025-10-11 18:45:12.333036+00');
INSERT INTO public.listing_attributes (listing_id, attribute_key, attribute_value, created_at) VALUES (1072, 'kategorija2', 'BATERIJE', '2025-10-11 18:45:12.333036+00');
INSERT INTO public.listing_attributes (listing_id, attribute_key, attribute_value, created_at) VALUES (1072, 'kategorija3', 'BATERIJE OUTLET', '2025-10-11 18:45:12.333036+00');
INSERT INTO public.listing_attributes (listing_id, attribute_key, attribute_value, created_at) VALUES (1072, 'godina_uvoza', '2025.', '2025-10-11 18:45:12.333036+00');
INSERT INTO public.listing_attributes (listing_id, attribute_key, attribute_value, created_at) VALUES (1072, 'zemlja_porekla', 'Kina', '2025-10-11 18:45:12.333036+00');

-- Listing 1071 attributes (Canon Printer)
INSERT INTO public.listing_attributes (listing_id, attribute_key, attribute_value, created_at) VALUES (1071, 'size', 'Standard', '2025-10-11 17:40:08.186825+00');
INSERT INTO public.listing_attributes (listing_id, attribute_key, attribute_value, created_at) VALUES (1071, 'brand', 'Canon', '2025-10-11 17:40:08.186825+00');
INSERT INTO public.listing_attributes (listing_id, attribute_key, attribute_value, created_at) VALUES (1071, 'color', 'Grey', '2025-10-11 17:40:08.186825+00');
INSERT INTO public.listing_attributes (listing_id, attribute_key, attribute_value, created_at) VALUES (1071, 'model', 'G3420', '2025-10-11 17:40:08.186825+00');
INSERT INTO public.listing_attributes (listing_id, attribute_key, attribute_value, created_at) VALUES (1071, 'material', 'Plastic', '2025-10-11 17:40:08.186825+00');

-- Migrate individual location data for listing 1071 (has_individual_location = true)
INSERT INTO public.listing_locations (listing_id, country, city, postal_code, address_line1, latitude, longitude, created_at, updated_at)
VALUES (1071, 'RS', NULL, NULL, 'Васе Стајића 18, Нови-Сад 21101, Южно-Бачский округ, Сербия', 45.25126266, 19.84361076, '2025-10-11 17:40:08.186825+00', '2025-10-11 17:40:08.186825+00');


--
-- Data for Name: b2c_store_hours; Type: TABLE DATA; Schema: public; Owner: -
--

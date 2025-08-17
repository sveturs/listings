INSERT INTO public.warehouses (id, code, name, type, address, city, postal_code, country, phone, email, manager_name, manager_phone, latitude, longitude, working_hours, total_area_m2, storage_area_m2, max_capacity_m3, current_occupancy_m3, supports_fbs, supports_pickup, has_refrigeration, has_loading_dock, is_active, created_at, updated_at) VALUES (1, 'NS-MAIN-01', 'Главный склад Sve Tu - Нови Сад', 'main', 'Улица Микија Манојловића 53', 'Нови Сад', '21000', 'RS', '+381 21 XXX-XXXX', NULL, NULL, NULL, NULL, NULL, '{"friday": {"open": "09:00", "close": "19:00"}, "monday": {"open": "09:00", "close": "19:00"}, "sunday": null, "tuesday": {"open": "09:00", "close": "19:00"}, "saturday": {"open": "10:00", "close": "16:00"}, "thursday": {"open": "09:00", "close": "19:00"}, "wednesday": {"open": "09:00", "close": "19:00"}}', NULL, NULL, NULL, 0.00, true, true, false, true, true, '2025-08-15 10:48:45.782268+00', '2025-08-15 10:48:45.782268+00');


--
-- Name: address_change_log_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.address_change_log_id_seq', 1, false);

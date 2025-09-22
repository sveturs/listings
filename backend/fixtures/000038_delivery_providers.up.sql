INSERT INTO public.delivery_providers (id, code, name, logo_url, is_active, supports_cod, supports_insurance, supports_tracking, api_config, capabilities, created_at, updated_at) VALUES (1, 'post_express', 'Post Express', '/images/providers/post-express.png', true, true, true, true, NULL, '{"max_volume_m3": 0.5, "max_weight_kg": 30, "delivery_types": ["standard", "express", "office_pickup"], "delivery_zones": ["serbia", "montenegro", "bosnia"]}', '2025-09-20 13:46:44.21752+00', '2025-09-20 13:46:44.21752+00');
INSERT INTO public.delivery_providers (id, code, name, logo_url, is_active, supports_cod, supports_insurance, supports_tracking, api_config, capabilities, created_at, updated_at) VALUES (2, 'bex_express', 'BEX Express', '/images/providers/bex-express.png', false, true, true, true, NULL, '{"max_volume_m3": 1.0, "max_weight_kg": 50, "delivery_types": ["standard", "express"], "delivery_zones": ["serbia"]}', '2025-09-20 13:46:44.21752+00', '2025-09-20 13:46:44.21752+00');
INSERT INTO public.delivery_providers (id, code, name, logo_url, is_active, supports_cod, supports_insurance, supports_tracking, api_config, capabilities, created_at, updated_at) VALUES (3, 'aks_express', 'AKS Express', '/images/providers/aks-express.png', false, true, true, true, NULL, '{"max_volume_m3": 0.8, "max_weight_kg": 40, "delivery_types": ["standard"], "delivery_zones": ["serbia", "region"]}', '2025-09-20 13:46:44.21752+00', '2025-09-20 13:46:44.21752+00');
INSERT INTO public.delivery_providers (id, code, name, logo_url, is_active, supports_cod, supports_insurance, supports_tracking, api_config, capabilities, created_at, updated_at) VALUES (4, 'd_express', 'D Express', '/images/providers/d-express.png', false, true, false, true, NULL, '{"max_volume_m3": 0.4, "max_weight_kg": 25, "delivery_types": ["standard", "same_day"], "delivery_zones": ["serbia"]}', '2025-09-20 13:46:44.21752+00', '2025-09-20 13:46:44.21752+00');
INSERT INTO public.delivery_providers (id, code, name, logo_url, is_active, supports_cod, supports_insurance, supports_tracking, api_config, capabilities, created_at, updated_at) VALUES (5, 'city_express', 'City Express', '/images/providers/city-express.png', false, true, true, true, NULL, '{"max_volume_m3": 0.6, "max_weight_kg": 35, "delivery_types": ["standard", "express"], "delivery_zones": ["serbia", "montenegro"]}', '2025-09-20 13:46:44.21752+00', '2025-09-20 13:46:44.21752+00');
INSERT INTO public.delivery_providers (id, code, name, logo_url, is_active, supports_cod, supports_insurance, supports_tracking, api_config, capabilities, created_at, updated_at) VALUES (6, 'dhl_express', 'DHL Express', '/images/providers/dhl-express.png', false, false, true, true, NULL, '{"max_volume_m3": 2.0, "max_weight_kg": 70, "delivery_types": ["express", "international"], "delivery_zones": ["worldwide"]}', '2025-09-20 13:46:44.21752+00', '2025-09-20 13:46:44.21752+00');


--
-- Data for Name: delivery_shipments; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: delivery_tracking_events; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: delivery_zones; Type: TABLE DATA; Schema: public; Owner: -
--

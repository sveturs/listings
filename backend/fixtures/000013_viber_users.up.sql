INSERT INTO public.viber_users (id, viber_id, user_id, name, avatar_url, language, country_code, api_version, subscribed, subscribed_at, last_session_at, conversation_started_at, created_at, updated_at) VALUES (1, '381604485063', NULL, 'Test User', '', '', '', 0, false, '2025-09-17 22:56:10.980104+00', '2025-09-17 22:58:10.728101+00', NULL, '2025-09-17 22:56:10.980104+00', '2025-09-17 23:14:38.067168+00');
INSERT INTO public.viber_users (id, viber_id, user_id, name, avatar_url, language, country_code, api_version, subscribed, subscribed_at, last_session_at, conversation_started_at, created_at, updated_at) VALUES (2, 'test-viber-id', NULL, 'Test User', '', '', '', 0, false, '2025-09-18 13:51:23.263168+00', '2025-09-18 13:53:00.782258+00', NULL, '2025-09-18 13:51:23.263168+00', '2025-09-18 13:54:27.643174+00');


--
-- Data for Name: view_statistics; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: vin_accident_history; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: vin_check_history; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: vin_decode_cache; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: vin_ownership_history; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: vin_recalls; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Name: address_change_log_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.address_change_log_id_seq', 1, false);

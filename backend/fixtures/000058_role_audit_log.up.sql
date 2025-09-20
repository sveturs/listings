INSERT INTO public.role_audit_log (id, user_id, target_user_id, action, old_role_id, new_role_id, details, created_at) VALUES (3, NULL, 4, 'role_changed', 5, 14, '{"new_role": "warehouse_manager", "old_role": "user", "timestamp": "2025-08-14T21:52:51.065245+00:00"}', '2025-08-14 21:52:51.065245+00');
INSERT INTO public.role_audit_log (id, user_id, target_user_id, action, old_role_id, new_role_id, details, created_at) VALUES (1, NULL, 6, 'role_changed', 2, 1, '{"new_role": "super_admin", "old_role": "admin", "timestamp": "2025-08-14T20:11:54.367852+00:00"}', '2025-08-14 20:11:54.367852+00');
INSERT INTO public.role_audit_log (id, user_id, target_user_id, action, old_role_id, new_role_id, details, created_at) VALUES (2, NULL, 8, 'role_changed', 5, 10, '{"new_role": "vendor_manager", "old_role": "user", "timestamp": "2025-08-14T21:48:21.315572+00:00"}', '2025-08-14 21:48:21.315572+00');
INSERT INTO public.role_audit_log (id, user_id, target_user_id, action, old_role_id, new_role_id, details, created_at) VALUES (4, NULL, 6, 'role_changed', NULL, 2, '{"new_role": "admin", "old_role": null, "timestamp": "2025-09-18T21:55:57.322829+00:00"}', '2025-09-18 21:55:57.322829+00');
INSERT INTO public.role_audit_log (id, user_id, target_user_id, action, old_role_id, new_role_id, details, created_at) VALUES (5, NULL, 3, 'role_changed', 2, 15, '{"new_role": "warehouse_worker", "old_role": "admin", "timestamp": "2025-09-20T11:48:17.876946+00:00"}', '2025-09-20 11:48:17.876946+00');


--
-- Data for Name: role_permissions; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.payment_methods (id, name, code, type, is_active, minimum_amount, maximum_amount, fee_percentage, fixed_fee, credentials, created_at) VALUES (1, 'Bank transfer', 'bank_transfer', 'bank', true, 1000.00, 10000000.00, 0.00, 100.00, NULL, '2025-04-17 10:15:18.418515');
INSERT INTO public.payment_methods (id, name, code, type, is_active, minimum_amount, maximum_amount, fee_percentage, fixed_fee, credentials, created_at) VALUES (2, 'Post office', 'post_office', 'cash', true, 500.00, 500000.00, 1.50, 50.00, NULL, '2025-04-17 10:15:18.418515');
INSERT INTO public.payment_methods (id, name, code, type, is_active, minimum_amount, maximum_amount, fee_percentage, fixed_fee, credentials, created_at) VALUES (3, 'IPS QR code', 'ips_qr', 'digital', true, 100.00, 1000000.00, 0.80, 0.00, NULL, '2025-04-17 10:15:18.418515');
INSERT INTO public.payment_methods (id, name, code, type, is_active, minimum_amount, maximum_amount, fee_percentage, fixed_fee, credentials, created_at) VALUES (26, 'Mock Payment (для тестирования)', 'mock_payment', 'card', true, 100.00, 1000000.00, 0.00, 0.00, NULL, '2025-07-01 11:55:07.694232');


--
-- Data for Name: payment_transactions; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: permissions; Type: TABLE DATA; Schema: public; Owner: -
--

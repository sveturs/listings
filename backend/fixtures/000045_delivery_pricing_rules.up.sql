INSERT INTO public.delivery_pricing_rules (id, provider_id, rule_type, weight_ranges, volume_ranges, zone_multipliers, fragile_surcharge, oversized_surcharge, special_handling_surcharge, min_price, max_price, custom_formula, priority, is_active, created_at, updated_at) VALUES (1, 1, 'weight_based', '[{"to": 1, "from": 0, "base_price": 300, "price_per_kg": 0}, {"to": 5, "from": 1, "base_price": 300, "price_per_kg": 50}, {"to": 10, "from": 5, "base_price": 500, "price_per_kg": 40}, {"to": 20, "from": 10, "base_price": 700, "price_per_kg": 35}, {"to": 30, "from": 20, "base_price": 1000, "price_per_kg": 30}]', NULL, NULL, 50.00, 0.00, 0.00, 250.00, 5000.00, NULL, 0, true, '2025-09-20 13:46:44.21752+00', '2025-09-20 13:46:44.21752+00');


--
-- Data for Name: delivery_providers; Type: TABLE DATA; Schema: public; Owner: -
--

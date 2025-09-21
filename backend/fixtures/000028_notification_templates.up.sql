INSERT INTO public.notification_templates (id, code, channel, name, subject, body_template, is_active, created_at, updated_at) VALUES (1, 'delivery_confirmed', 'email', 'Заказ подтвержден', 'Заказ {{.TrackingNumber}} подтвержден', 'Здравствуйте, {{.RecipientName}}!

Ваш заказ {{.TrackingNumber}} подтвержден и будет передан в службу доставки.

С уважением,
Команда Sve Tu', true, '2025-09-20 15:17:32.182386+00', '2025-09-20 15:17:32.182386+00');
INSERT INTO public.notification_templates (id, code, channel, name, subject, body_template, is_active, created_at, updated_at) VALUES (2, 'delivery_picked_up', 'email', 'Передан в службу доставки', 'Заказ {{.TrackingNumber}} передан в службу доставки', 'Здравствуйте, {{.RecipientName}}!

Ваш заказ {{.TrackingNumber}} передан в службу доставки.
Ожидаемая дата доставки: {{.EstimatedDelivery}}

Отследить заказ: https://svetu.rs/tracking/{{.TrackingNumber}}

С уважением,
Команда Sve Tu', true, '2025-09-20 15:17:32.182386+00', '2025-09-20 15:17:32.182386+00');
INSERT INTO public.notification_templates (id, code, channel, name, subject, body_template, is_active, created_at, updated_at) VALUES (3, 'delivery_out_for_delivery', 'email', 'Передан курьеру', 'Заказ {{.TrackingNumber}} передан курьеру', 'Здравствуйте, {{.RecipientName}}!

Ваш заказ {{.TrackingNumber}} передан курьеру для доставки.
Ожидайте доставку сегодня.

Курьер свяжется с вами перед доставкой.

С уважением,
Команда Sve Tu', true, '2025-09-20 15:17:32.182386+00', '2025-09-20 15:17:32.182386+00');
INSERT INTO public.notification_templates (id, code, channel, name, subject, body_template, is_active, created_at, updated_at) VALUES (4, 'delivery_delivered', 'email', 'Заказ доставлен', 'Заказ {{.TrackingNumber}} доставлен', 'Здравствуйте, {{.RecipientName}}!

Ваш заказ {{.TrackingNumber}} успешно доставлен.

Спасибо за покупку!
Оставьте отзыв: https://svetu.rs/orders/{{.OrderID}}/review

С уважением,
Команда Sve Tu', true, '2025-09-20 15:17:32.182386+00', '2025-09-20 15:17:32.182386+00');
INSERT INTO public.notification_templates (id, code, channel, name, subject, body_template, is_active, created_at, updated_at) VALUES (5, 'delivery_failed', 'email', 'Проблема с доставкой', 'Проблема с доставкой заказа {{.TrackingNumber}}', 'Здравствуйте, {{.RecipientName}}!

К сожалению, возникла проблема с доставкой заказа {{.TrackingNumber}}.
Причина: {{.Reason}}

Пожалуйста, свяжитесь с нами для решения вопроса.

С уважением,
Команда Sve Tu', true, '2025-09-20 15:17:32.182386+00', '2025-09-20 15:17:32.182386+00');
INSERT INTO public.notification_templates (id, code, channel, name, subject, body_template, is_active, created_at, updated_at) VALUES (6, 'delivery_out_for_delivery_sms', 'sms', 'Передан курьеру', NULL, 'Заказ {{.TrackingNumber}} передан курьеру. Ожидайте доставку сегодня.', true, '2025-09-20 15:17:32.182386+00', '2025-09-20 15:17:32.182386+00');
INSERT INTO public.notification_templates (id, code, channel, name, subject, body_template, is_active, created_at, updated_at) VALUES (7, 'delivery_delivered_sms', 'sms', 'Доставлен', NULL, 'Заказ {{.TrackingNumber}} доставлен. Спасибо за покупку!', true, '2025-09-20 15:17:32.182386+00', '2025-09-20 15:17:32.182386+00');
INSERT INTO public.notification_templates (id, code, channel, name, subject, body_template, is_active, created_at, updated_at) VALUES (8, 'delivery_failed_sms', 'sms', 'Проблема с доставкой', NULL, 'Проблема с доставкой {{.TrackingNumber}}. Свяжитесь с нами.', true, '2025-09-20 15:17:32.182386+00', '2025-09-20 15:17:32.182386+00');


--
-- Data for Name: notifications; Type: TABLE DATA; Schema: public; Owner: -
--

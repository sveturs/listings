-- Добавление тестовых данных для демонстрации системы мониторинга логистики

-- Добавляем тестовые данные в BEX shipments
INSERT INTO bex_shipments (
    marketplace_order_id,
    tracking_number,
    sender_name,
    sender_address,
    sender_city,
    sender_postal_code,
    sender_phone,
    sender_email,
    recipient_name,
    recipient_address,
    recipient_city,
    recipient_postal_code,
    recipient_phone,
    recipient_email,
    shipment_type,
    shipment_category,
    shipment_contents,
    weight_kg,
    total_packages,
    pay_type,
    cod_amount,
    insurance_amount,
    personal_delivery,
    status,
    status_text,
    registered_at,
    delivered_at,
    status_history,
    created_at,
    updated_at
) VALUES
    -- Доставленные отправления
    (NULL, 'BEX2024110001', 'SveTu Marketplace', 'Bulevar Oslobođenja 100', 'Novi Sad', '21000', '+381601234567', 'orders@svetu.rs',
     'Милан Петровић', 'Улица Лазе Костића 15', 'Београд', '11000', '+381641234567', 'milan.petrovic@example.com',
     1, 1, 1, 2.5, 1, 1, 250.00, 0, false, 3, 'Delivered', 
     NOW() - INTERVAL '5 days', NOW() - INTERVAL '2 days',
     '[{"status": "registered", "timestamp": "2024-11-05T10:00:00Z"}, {"status": "in_transit", "timestamp": "2024-11-06T14:00:00Z"}, {"status": "delivered", "timestamp": "2024-11-08T16:00:00Z"}]'::jsonb,
     NOW() - INTERVAL '5 days', NOW() - INTERVAL '2 days'),
     
    (NULL, 'BEX2024110002', 'Electronics Store', 'Zmaj Jovina 25', 'Novi Sad', '21000', '+381601234568', 'shop@electronics.rs',
     'Јована Стојановић', 'Кнез Михаилова 50', 'Београд', '11000', '+381641234568', 'jovana@example.com',
     1, 2, 2, 1.2, 1, 2, 0, 150.00, false, 3, 'Delivered',
     NOW() - INTERVAL '4 days', NOW() - INTERVAL '1 day',
     '[{"status": "registered", "timestamp": "2024-11-06T09:00:00Z"}, {"status": "delivered", "timestamp": "2024-11-09T15:00:00Z"}]'::jsonb,
     NOW() - INTERVAL '4 days', NOW() - INTERVAL '1 day'),
     
    -- В пути
    (NULL, 'BEX2024110003', 'Fashion Boutique', 'Dunavska 15', 'Novi Sad', '21000', '+381601234569', 'info@fashion.rs',
     'Марко Николић', 'Булевар Краља Александра 120', 'Београд', '11000', '+381641234569', 'marko.n@example.com',
     1, 1, 3, 0.8, 1, 1, 180.00, 0, true, 2, 'In Transit',
     NOW() - INTERVAL '2 days', NULL,
     '[{"status": "registered", "timestamp": "2024-11-08T11:00:00Z"}, {"status": "in_transit", "timestamp": "2024-11-09T08:00:00Z"}]'::jsonb,
     NOW() - INTERVAL '2 days', NOW() - INTERVAL '1 hour'),
     
    (NULL, 'BEX2024110004', 'Sports Equipment', 'Futoška 45', 'Novi Sad', '21000', '+381601234570', 'orders@sports.rs',
     'Ана Јовановић', 'Народног фронта 23', 'Нови Сад', '21000', '+381641234570', 'ana.j@example.com',
     1, 1, 4, 3.5, 2, 1, 420.00, 50.00, false, 2, 'In Transit',
     NOW() - INTERVAL '1 day', NULL,
     '[{"status": "registered", "timestamp": "2024-11-09T14:00:00Z"}, {"status": "in_transit", "timestamp": "2024-11-09T18:00:00Z"}]'::jsonb,
     NOW() - INTERVAL '1 day', NOW() - INTERVAL '30 minutes'),
     
    -- Ожидает отправки
    (NULL, 'BEX2024110005', 'Book Store', 'Modene 5', 'Novi Sad', '21000', '+381601234571', 'books@store.rs',
     'Петар Поповић', 'Светозара Марковића 85', 'Крагујевац', '34000', '+381641234571', 'petar.popovic@example.com',
     1, 1, 5, 0.5, 1, 2, 0, 0, false, 1, 'Pending',
     NOW() - INTERVAL '3 hours', NULL,
     '[{"status": "registered", "timestamp": "2024-11-10T07:00:00Z"}]'::jsonb,
     NOW() - INTERVAL '3 hours', NOW() - INTERVAL '3 hours'),
     
    -- Проблемное (не доставлено)
    (NULL, 'BEX2024110006', 'Home Decor', 'Jevrejska 10', 'Novi Sad', '21000', '+381601234572', 'decor@home.rs',
     'Милица Савић', 'Косовска 45', 'Ниш', '18000', '+381641234572', 'milica.s@example.com',
     1, 1, 6, 5.0, 3, 1, 890.00, 100.00, true, 4, 'Failed',
     NOW() - INTERVAL '7 days', NULL,
     '[{"status": "registered", "timestamp": "2024-11-03T10:00:00Z"}, {"status": "in_transit", "timestamp": "2024-11-04T09:00:00Z"}, {"status": "failed", "timestamp": "2024-11-05T16:00:00Z", "reason": "Recipient not available"}]'::jsonb,
     NOW() - INTERVAL '7 days', NOW() - INTERVAL '5 days');

-- Добавляем тестовые данные в Post Express shipments
INSERT INTO post_express_shipments (
    marketplace_order_id,
    tracking_number,
    barcode,
    post_express_id,
    sender_name,
    sender_address,
    sender_city,
    sender_postal_code,
    sender_phone,
    sender_email,
    recipient_name,
    recipient_address,
    recipient_city,
    recipient_postal_code,
    recipient_phone,
    recipient_email,
    weight_kg,
    length_cm,
    width_cm,
    height_cm,
    declared_value,
    service_type,
    cod_amount,
    insurance_amount,
    base_price,
    total_price,
    status,
    delivery_status,
    registered_at,
    delivered_at,
    status_history,
    created_at,
    updated_at
) VALUES
    -- Доставленные
    (NULL, 'PE2024110001', '1234567890123', 'PE-001', 'Organic Food Store', 'Danila Kiša 20', 'Novi Sad', '21000', 
     '+381601234573', 'organic@food.rs',
     'Драган Миловановић', 'Војводе Степе 250', 'Београд', '11000', '+381641234573', 'dragan.m@example.com',
     2.0, 30, 20, 15, 150.00, 'standard', 150.00, 0, 25.00, 30.00, 'delivered', 'Delivered successfully',
     NOW() - INTERVAL '6 days', NOW() - INTERVAL '3 days',
     '[{"status": "pending", "timestamp": "2024-11-04T10:00:00Z"}, {"status": "in_transit", "timestamp": "2024-11-05T08:00:00Z"}, {"status": "delivered", "timestamp": "2024-11-07T14:00:00Z"}]'::jsonb,
     NOW() - INTERVAL '6 days', NOW() - INTERVAL '3 days'),
     
    (NULL, 'PE2024110002', '1234567890124', 'PE-002', 'Tech Gadgets', 'Bulevar Evrope 50', 'Novi Sad', '21000',
     '+381601234574', 'gadgets@tech.rs',
     'Тамара Ђорђевић', 'Цара Душана 100', 'Суботица', '24000', '+381641234574', 'tamara.dj@example.com',
     0.5, 20, 15, 10, 350.00, 'express', 0, 50.00, 35.00, 45.00, 'delivered', 'Delivered to neighbor',
     NOW() - INTERVAL '3 days', NOW() - INTERVAL '1 day',
     '[{"status": "pending", "timestamp": "2024-11-07T11:00:00Z"}, {"status": "delivered", "timestamp": "2024-11-09T13:00:00Z"}]'::jsonb,
     NOW() - INTERVAL '3 days', NOW() - INTERVAL '1 day'),
     
    -- В пути
    (NULL, 'PE2024110003', '1234567890125', 'PE-003', 'Jewelry Store', 'Zmaj Jovina 30', 'Novi Sad', '21000',
     '+381601234575', 'jewelry@store.rs',
     'Никола Станковић', 'Партизанска 25', 'Панчево', '26000', '+381641234575', 'nikola.s@example.com',
     0.2, 15, 10, 5, 500.00, 'express', 500.00, 100.00, 40.00, 55.00, 'in_transit', 'Out for delivery',
     NOW() - INTERVAL '1 day', NULL,
     '[{"status": "pending", "timestamp": "2024-11-09T09:00:00Z"}, {"status": "in_transit", "timestamp": "2024-11-09T14:00:00Z"}]'::jsonb,
     NOW() - INTERVAL '1 day', NOW() - INTERVAL '2 hours'),
     
    -- Ожидает
    (NULL, 'PE2024110004', '1234567890126', 'PE-004', 'Pet Supplies', 'Kisačka 55', 'Novi Sad', '21000',
     '+381601234576', 'pets@supplies.rs',
     'Јелена Радовановић', 'Југ Богданова 10', 'Зрењанин', '23000', '+381641234576', 'jelena.r@example.com',
     4.0, 40, 30, 25, 120.00, 'standard', 120.00, 0, 30.00, 35.00, 'pending', NULL,
     NOW() - INTERVAL '2 hours', NULL,
     '[{"status": "pending", "timestamp": "2024-11-10T08:00:00Z"}]'::jsonb,
     NOW() - INTERVAL '2 hours', NOW() - INTERVAL '2 hours'),
     
    -- Проблемное (задержка)
    (NULL, 'PE2024110005', 'Cosmetics Shop', '1234567890127', 'PE-005', 'Cara Lazara 15', 'Novi Sad', '21000',
     '+381601234577', 'beauty@shop.rs',
     'Стефан Илић', 'Синђелићева 50', 'Чачак', '32000', '+381641234577', 'stefan.i@example.com',
     1.5, 25, 20, 15, 200.00, 'standard', 200.00, 30.00, 28.00, 35.00, 'in_transit', 'Delayed - weather conditions',
     NOW() - INTERVAL '8 days', NULL,
     '[{"status": "pending", "timestamp": "2024-11-02T10:00:00Z"}, {"status": "in_transit", "timestamp": "2024-11-03T08:00:00Z"}]'::jsonb,
     NOW() - INTERVAL '8 days', NOW() - INTERVAL '7 days'),
     
    -- Неудачная доставка
    (NULL, 'PE2024110006', '1234567890128', 'PE-006', 'Furniture Store', 'Rumenačka 100', 'Novi Sad', '21000',
     '+381601234578', 'furniture@store.rs',
     'Марија Костић', 'Хајдук Вељкова 75', 'Сомбор', '25000', '+381641234578', 'marija.k@example.com',
     15.0, 60, 40, 30, 1500.00, 'standard', 1500.00, 200.00, 50.00, 75.00, 'failed', 'Address not found',
     NOW() - INTERVAL '10 days', NULL,
     '[{"status": "pending", "timestamp": "2024-10-31T09:00:00Z"}, {"status": "in_transit", "timestamp": "2024-11-01T10:00:00Z"}, {"status": "failed", "timestamp": "2024-11-02T15:00:00Z", "reason": "Incorrect address"}]'::jsonb,
     NOW() - INTERVAL '10 days', NOW() - INTERVAL '8 days');

-- Таблица problem_shipments пока не создана, добавим позже
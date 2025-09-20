-- Система уведомлений о доставке

-- Удаляем старую таблицу если есть
DROP TABLE IF EXISTS delivery_notifications CASCADE;

-- Создаем новую таблицу уведомлений
CREATE TABLE delivery_notifications (
    id SERIAL PRIMARY KEY,
    shipment_id INTEGER REFERENCES delivery_shipments(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id),
    channel VARCHAR(20) NOT NULL, -- email, sms, viber, telegram, push
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending, sent, failed
    template VARCHAR(50),
    data JSONB, -- данные для шаблона
    sent_at TIMESTAMPTZ,
    error_message TEXT,
    retry_count INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Шаблоны уведомлений
CREATE TABLE IF NOT EXISTS notification_templates (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    channel VARCHAR(20) NOT NULL,
    name VARCHAR(255) NOT NULL,
    subject TEXT, -- для email
    body_template TEXT NOT NULL, -- шаблон сообщения
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Настройки уведомлений пользователя
CREATE TABLE user_notification_preferences (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    channel VARCHAR(20) NOT NULL,
    is_enabled BOOLEAN DEFAULT true,

    -- Какие события хочет получать
    notify_on_confirmed BOOLEAN DEFAULT true,
    notify_on_picked_up BOOLEAN DEFAULT true,
    notify_on_in_transit BOOLEAN DEFAULT false,
    notify_on_out_for_delivery BOOLEAN DEFAULT true,
    notify_on_delivered BOOLEAN DEFAULT true,
    notify_on_failed BOOLEAN DEFAULT true,
    notify_on_returned BOOLEAN DEFAULT true,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE(user_id, channel)
);

-- Контакты для уведомлений
CREATE TABLE user_notification_contacts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    channel VARCHAR(20) NOT NULL,
    contact_value VARCHAR(255) NOT NULL, -- email, phone, chat_id
    is_verified BOOLEAN DEFAULT false,
    is_primary BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE(user_id, channel, contact_value)
);

-- Индексы для производительности
CREATE INDEX idx_delivery_notifications_shipment ON delivery_notifications(shipment_id);
CREATE INDEX idx_delivery_notifications_user ON delivery_notifications(user_id);
CREATE INDEX idx_delivery_notifications_status ON delivery_notifications(status);
CREATE INDEX idx_delivery_notifications_channel ON delivery_notifications(channel);
CREATE INDEX idx_delivery_notifications_created ON delivery_notifications(created_at);
CREATE INDEX idx_user_notification_preferences_user ON user_notification_preferences(user_id);
CREATE INDEX idx_user_notification_contacts_user ON user_notification_contacts(user_id);

-- Вставляем дефолтные шаблоны (если еще не существуют)
INSERT INTO notification_templates (code, channel, name, subject, body_template)
SELECT * FROM (VALUES
-- Email шаблоны
('delivery_confirmed', 'email', 'Заказ подтвержден', 'Заказ {{.TrackingNumber}} подтвержден',
'Здравствуйте, {{.RecipientName}}!

Ваш заказ {{.TrackingNumber}} подтвержден и будет передан в службу доставки.

С уважением,
Команда Sve Tu'),

('delivery_picked_up', 'email', 'Передан в службу доставки', 'Заказ {{.TrackingNumber}} передан в службу доставки',
'Здравствуйте, {{.RecipientName}}!

Ваш заказ {{.TrackingNumber}} передан в службу доставки.
Ожидаемая дата доставки: {{.EstimatedDelivery}}

Отследить заказ: https://svetu.rs/tracking/{{.TrackingNumber}}

С уважением,
Команда Sve Tu'),

('delivery_out_for_delivery', 'email', 'Передан курьеру', 'Заказ {{.TrackingNumber}} передан курьеру',
'Здравствуйте, {{.RecipientName}}!

Ваш заказ {{.TrackingNumber}} передан курьеру для доставки.
Ожидайте доставку сегодня.

Курьер свяжется с вами перед доставкой.

С уважением,
Команда Sve Tu'),

('delivery_delivered', 'email', 'Заказ доставлен', 'Заказ {{.TrackingNumber}} доставлен',
'Здравствуйте, {{.RecipientName}}!

Ваш заказ {{.TrackingNumber}} успешно доставлен.

Спасибо за покупку!
Оставьте отзыв: https://svetu.rs/orders/{{.OrderID}}/review

С уважением,
Команда Sve Tu'),

('delivery_failed', 'email', 'Проблема с доставкой', 'Проблема с доставкой заказа {{.TrackingNumber}}',
'Здравствуйте, {{.RecipientName}}!

К сожалению, возникла проблема с доставкой заказа {{.TrackingNumber}}.
Причина: {{.Reason}}

Пожалуйста, свяжитесь с нами для решения вопроса.

С уважением,
Команда Sve Tu'),

-- SMS шаблоны
('delivery_out_for_delivery_sms', 'sms', 'Передан курьеру', NULL,
'Заказ {{.TrackingNumber}} передан курьеру. Ожидайте доставку сегодня.'),

('delivery_delivered_sms', 'sms', 'Доставлен', NULL,
'Заказ {{.TrackingNumber}} доставлен. Спасибо за покупку!'),

('delivery_failed_sms', 'sms', 'Проблема с доставкой', NULL,
'Проблема с доставкой {{.TrackingNumber}}. Свяжитесь с нами.')
) AS t(code, channel, name, subject, body_template)
WHERE NOT EXISTS (
    SELECT 1 FROM notification_templates WHERE code = t.code
);

-- Триггер для обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_delivery_notifications_updated_at
    BEFORE UPDATE ON delivery_notifications
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_notification_templates_updated_at
    BEFORE UPDATE ON notification_templates
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_notification_preferences_updated_at
    BEFORE UPDATE ON user_notification_preferences
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
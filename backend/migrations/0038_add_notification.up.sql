-- Таблица настроек уведомлений пользователя
CREATE TABLE notification_settings (
    user_id INT NOT NULL REFERENCES users(id),
    notification_type VARCHAR(50) NOT NULL, -- new_message, new_review, review_vote, review_response, listing_status, favorite_price
    telegram_enabled BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, notification_type)
);

-- Таблица для хранения Telegram данных пользователей
CREATE TABLE user_telegram_connections (
    user_id INT PRIMARY KEY REFERENCES users(id),
    telegram_chat_id VARCHAR(100) NOT NULL,
    telegram_username VARCHAR(100),
    connected_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);



-- Таблица для истории уведомлений
CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    type VARCHAR(50) NOT NULL,
    title TEXT NOT NULL,
    message TEXT NOT NULL,
    data JSONB,
    is_read BOOLEAN DEFAULT false,
    delivered_to JSONB, -- {telegram: true}
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индексы
CREATE INDEX idx_notifications_user ON notifications(user_id);
CREATE INDEX idx_notifications_type ON notifications(type);
CREATE INDEX idx_notifications_created ON notifications(created_at);

-- Триггер для обновления updated_at
CREATE OR REPLACE FUNCTION update_notification_settings_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_notification_settings_timestamp
    BEFORE UPDATE ON notification_settings
    FOR EACH ROW
    EXECUTE FUNCTION update_notification_settings_updated_at();
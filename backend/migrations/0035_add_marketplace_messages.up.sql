-- /backend/migrations/0035_add_marketplace_messages.up.sql

-- Сначала создаем таблицу чатов, так как на нее будут ссылаться сообщения
CREATE TABLE marketplace_chats (
    id SERIAL PRIMARY KEY,
    listing_id INT REFERENCES marketplace_listings(id) ON DELETE CASCADE,
    buyer_id INT REFERENCES users(id),
    seller_id INT REFERENCES users(id),
    last_message_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_archived BOOLEAN DEFAULT false,
    UNIQUE(listing_id, buyer_id, seller_id)
);

-- Затем создаем таблицу сообщений
CREATE TABLE marketplace_messages (
    id SERIAL PRIMARY KEY,
    chat_id INT REFERENCES marketplace_chats(id) ON DELETE CASCADE,
    listing_id INT REFERENCES marketplace_listings(id) ON DELETE CASCADE,
    sender_id INT REFERENCES users(id),
    receiver_id INT REFERENCES users(id),
    content TEXT NOT NULL,
    is_read BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для оптимизации запросов
CREATE INDEX idx_marketplace_messages_chat ON marketplace_messages(chat_id);
CREATE INDEX idx_marketplace_messages_listing ON marketplace_messages(listing_id);
CREATE INDEX idx_marketplace_messages_sender ON marketplace_messages(sender_id);
CREATE INDEX idx_marketplace_messages_receiver ON marketplace_messages(receiver_id);
CREATE INDEX idx_marketplace_messages_created ON marketplace_messages(created_at);

CREATE INDEX idx_marketplace_chats_buyer ON marketplace_chats(buyer_id);
CREATE INDEX idx_marketplace_chats_seller ON marketplace_chats(seller_id);
CREATE INDEX idx_marketplace_chats_updated ON marketplace_chats(updated_at);

-- Функция для обновления updated_at
CREATE OR REPLACE FUNCTION update_marketplace_chats_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Триггеры для автоматического обновления updated_at
CREATE TRIGGER update_marketplace_chats_timestamp
    BEFORE UPDATE ON marketplace_chats
    FOR EACH ROW
    EXECUTE FUNCTION update_marketplace_chats_updated_at();

CREATE TRIGGER update_marketplace_messages_timestamp
    BEFORE UPDATE ON marketplace_messages
    FOR EACH ROW
    EXECUTE FUNCTION update_marketplace_chats_updated_at();
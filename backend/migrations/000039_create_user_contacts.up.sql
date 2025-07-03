-- Создание таблицы контактов пользователей
CREATE TABLE IF NOT EXISTS user_contacts (
    id BIGSERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    contact_user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending, accepted, blocked
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Дополнительные поля
    added_from_chat_id INTEGER REFERENCES marketplace_chats(id), -- Откуда добавлен контакт
    notes TEXT, -- Заметки о контакте
    
    -- Индексы и ограничения
    UNIQUE(user_id, contact_user_id),
    CHECK (user_id != contact_user_id), -- Нельзя добавить себя в контакты
    CHECK (status IN ('pending', 'accepted', 'blocked'))
);

-- Создание индексов для быстрых запросов
CREATE INDEX idx_user_contacts_user_id ON user_contacts(user_id);
CREATE INDEX idx_user_contacts_contact_user_id ON user_contacts(contact_user_id);
CREATE INDEX idx_user_contacts_status ON user_contacts(status);
CREATE INDEX idx_user_contacts_created_at ON user_contacts(created_at);

-- Создание триггера для обновления updated_at
CREATE OR REPLACE FUNCTION update_user_contacts_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_user_contacts_updated_at
    BEFORE UPDATE ON user_contacts
    FOR EACH ROW EXECUTE FUNCTION update_user_contacts_updated_at();

-- Создание таблицы настроек приватности
CREATE TABLE IF NOT EXISTS user_privacy_settings (
    user_id INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    allow_contact_requests BOOLEAN DEFAULT true,
    allow_messages_from_contacts_only BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Триггер для настроек приватности
CREATE TRIGGER update_user_privacy_settings_updated_at
    BEFORE UPDATE ON user_privacy_settings
    FOR EACH ROW EXECUTE FUNCTION update_user_contacts_updated_at();
    
    
-- Создание таблицы для хранения вложений в чатах
CREATE TABLE IF NOT EXISTS chat_attachments (
    id SERIAL PRIMARY KEY,
    message_id INT NOT NULL REFERENCES marketplace_messages(id) ON DELETE CASCADE,
    file_type VARCHAR(20) NOT NULL CHECK (file_type IN ('image', 'video', 'document')),
    file_path VARCHAR(500) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    storage_type VARCHAR(20) DEFAULT 'minio',
    storage_bucket VARCHAR(100) DEFAULT 'chat-files',
    public_url TEXT,
    thumbnail_url TEXT, -- для превью видео
    metadata JSONB, -- duration для видео, pages для PDF и т.д.
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_message FOREIGN KEY (message_id) REFERENCES marketplace_messages(id) ON DELETE CASCADE
);

-- Индексы для производительности
CREATE INDEX idx_chat_attachments_message_id ON chat_attachments(message_id);
CREATE INDEX idx_chat_attachments_created_at ON chat_attachments(created_at);
CREATE INDEX idx_chat_attachments_file_type ON chat_attachments(file_type);

-- Обновление таблицы сообщений для поддержки вложений
ALTER TABLE marketplace_messages 
ADD COLUMN IF NOT EXISTS has_attachments BOOLEAN DEFAULT false,
ADD COLUMN IF NOT EXISTS attachments_count INT DEFAULT 0;

-- Комментарии для документации
COMMENT ON TABLE chat_attachments IS 'Хранение файлов, прикрепленных к сообщениям чата';
COMMENT ON COLUMN chat_attachments.file_type IS 'Тип файла: image, video, document';
COMMENT ON COLUMN chat_attachments.thumbnail_url IS 'URL превью для видео файлов';
COMMENT ON COLUMN chat_attachments.metadata IS 'Дополнительные метаданные файла (размеры изображения, длительность видео, количество страниц PDF и т.д.)';

-- Создаем частичный уникальный индекс для прямых сообщений (где listing_id IS NULL)
-- Этот индекс гарантирует, что между двумя пользователями может быть только один прямой чат
CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_direct_chat 
ON marketplace_chats (
    LEAST(buyer_id, seller_id), 
    GREATEST(buyer_id, seller_id)
) 
WHERE listing_id IS NULL;

-- Удаляем дубликаты прямых чатов, оставляя только самый новый
DELETE FROM marketplace_chats
WHERE id IN (
    SELECT id
    FROM (
        SELECT id,
                ROW_NUMBER() OVER (
                    PARTITION BY 
                        LEAST(buyer_id, seller_id), 
                        GREATEST(buyer_id, seller_id)
                    ORDER BY last_message_at DESC
                ) AS rn
        FROM marketplace_chats
        WHERE listing_id IS NULL
    ) t
    WHERE t.rn > 1
);

-- Добавление индексов для улучшения производительности
-- Все индексы создаются без CONCURRENTLY для совместимости с миграциями

-- ===== ИНДЕКСЫ ДЛЯ ЧАТА =====

-- Индекс для быстрого поиска непрочитанных сообщений
CREATE INDEX IF NOT EXISTS idx_marketplace_messages_unread 
ON marketplace_messages(receiver_id, is_read) 
WHERE NOT is_read;

-- Составной индекс для поиска чатов пользователя
CREATE INDEX IF NOT EXISTS idx_marketplace_chats_user_lookup 
ON marketplace_chats(buyer_id, seller_id, last_message_at DESC);

-- Индекс для фильтрации неархивированных чатов
CREATE INDEX IF NOT EXISTS idx_marketplace_chats_archived 
ON marketplace_chats(is_archived) 
WHERE NOT is_archived;

-- Индекс для быстрой выборки сообщений чата с сортировкой
CREATE INDEX IF NOT EXISTS idx_marketplace_messages_chat_ordered 
ON marketplace_messages(chat_id, created_at DESC);

-- Индекс для поиска вложений по сообщению
CREATE INDEX IF NOT EXISTS idx_chat_attachments_message 
ON chat_attachments(message_id);

-- Составной индекс для поддержки запросов с listing_id
CREATE INDEX IF NOT EXISTS idx_marketplace_chats_listing 
ON marketplace_chats(listing_id)
WHERE listing_id IS NOT NULL;

-- Индекс для подсчета непрочитанных по чату
CREATE INDEX IF NOT EXISTS idx_marketplace_messages_chat_unread
ON marketplace_messages(chat_id, receiver_id)
WHERE NOT is_read;

-- Составной индекс для подсчета непрочитанных сообщений по получателю
CREATE INDEX IF NOT EXISTS idx_marketplace_messages_receiver_unread_count 
ON marketplace_messages(receiver_id, chat_id)
WHERE NOT is_read;

-- Индекс для поиска последнего сообщения в чате
CREATE INDEX IF NOT EXISTS idx_marketplace_messages_chat_last 
ON marketplace_messages(chat_id, id DESC);

-- Составной индекс для поиска чатов между двумя конкретными пользователями
CREATE INDEX IF NOT EXISTS idx_marketplace_chats_participants 
ON marketplace_chats(
    LEAST(buyer_id, seller_id), 
    GREATEST(buyer_id, seller_id)
);

-- Индекс для поиска активных чатов с сортировкой
CREATE INDEX IF NOT EXISTS idx_marketplace_chats_active_sorted 
ON marketplace_chats(last_message_at DESC)
WHERE NOT is_archived;

-- Составной индекс для эффективного JOIN с listings
CREATE INDEX IF NOT EXISTS idx_marketplace_chats_listing_participants 
ON marketplace_chats(listing_id, buyer_id, seller_id)
WHERE listing_id IS NOT NULL;

-- ===== ИНДЕКСЫ ДЛЯ ОБЪЯВЛЕНИЙ =====

-- Составной индекс для поиска активных объявлений с сортировкой по дате создания
CREATE INDEX IF NOT EXISTS idx_marketplace_listings_status_created 
ON marketplace_listings(status, created_at DESC)
WHERE status = 'active';

-- Индекс для географического поиска
CREATE INDEX IF NOT EXISTS idx_marketplace_listings_location 
ON marketplace_listings(latitude, longitude)
WHERE latitude IS NOT NULL AND longitude IS NOT NULL;

-- Индекс для поиска по категории и статусу
CREATE INDEX IF NOT EXISTS idx_marketplace_listings_category_status 
ON marketplace_listings(category_id, status)
WHERE status = 'active';

-- Индекс для поиска по пользователю и статусу
CREATE INDEX IF NOT EXISTS idx_marketplace_listings_user_status 
ON marketplace_listings(user_id, status, created_at DESC);

-- Индекс для поиска по городу
CREATE INDEX IF NOT EXISTS idx_marketplace_listings_city 
ON marketplace_listings(address_city)
WHERE address_city IS NOT NULL;

-- Индекс для полнотекстового поиска по заголовку
CREATE INDEX IF NOT EXISTS idx_marketplace_listings_title_gin 
ON marketplace_listings USING gin(to_tsvector('simple', title));

-- Индекс для поиска по ценовому диапазону
CREATE INDEX IF NOT EXISTS idx_marketplace_listings_price 
ON marketplace_listings(price)
WHERE price IS NOT NULL AND status = 'active';

-- ===== ИНДЕКСЫ ДЛЯ ПОЛЬЗОВАТЕЛЕЙ =====

-- Индекс для поиска по email
CREATE INDEX IF NOT EXISTS idx_users_email_lower 
ON users(LOWER(email));

-- Индекс для активных пользователей
CREATE INDEX IF NOT EXISTS idx_users_active 
ON users(last_seen DESC)
WHERE account_status = 'active';

-- ===== ИНДЕКСЫ ДЛЯ ИЗБРАННОГО =====

-- Составной индекс для подсчета избранных по пользователю
CREATE INDEX IF NOT EXISTS idx_marketplace_favorites_user_count 
ON marketplace_favorites(user_id, created_at DESC);

-- Индекс для поиска избранных объявлений
CREATE INDEX IF NOT EXISTS idx_marketplace_favorites_listing 
ON marketplace_favorites(listing_id);

-- ===== ИНДЕКСЫ ДЛЯ УВЕДОМЛЕНИЙ =====

-- Составной индекс для непрочитанных уведомлений
CREATE INDEX IF NOT EXISTS idx_notifications_user_unread 
ON notifications(user_id, created_at DESC)
WHERE NOT is_read;

-- ===== ИНДЕКСЫ ДЛЯ ПЕРЕВОДОВ =====

-- Индекс для быстрого поиска переводов по типу и языку
CREATE INDEX IF NOT EXISTS idx_translations_type_lang 
ON translations(entity_type, language);

-- ===== ИНДЕКСЫ ДЛЯ КАТЕГОРИЙ =====

-- Индекс для поиска по slug
CREATE INDEX IF NOT EXISTS idx_marketplace_categories_slug 
ON marketplace_categories(slug);

-- ===== ДОПОЛНИТЕЛЬНЫЕ ИНДЕКСЫ =====

-- Индексы для изображений
CREATE INDEX IF NOT EXISTS idx_marketplace_images_listing_main 
ON marketplace_images(listing_id, is_main)
WHERE is_main = true;

-- Индекс для review_responses
CREATE INDEX IF NOT EXISTS idx_review_responses_review 
ON review_responses(review_id);

-- Индекс для price_history для поиска текущих цен
CREATE INDEX IF NOT EXISTS idx_price_history_current 
ON price_history(listing_id, effective_from DESC)
WHERE effective_to IS NULL;

-- ===== ОБНОВЛЕНИЕ СТАТИСТИКИ =====

-- Анализируем основные таблицы для обновления статистики планировщика
ANALYZE marketplace_listings;
ANALYZE marketplace_messages;
ANALYZE marketplace_chats;
ANALYZE users;
ANALYZE marketplace_favorites;
ANALYZE marketplace_categories;
ANALYZE notifications;
ANALYZE translations;
ANALYZE chat_attachments;
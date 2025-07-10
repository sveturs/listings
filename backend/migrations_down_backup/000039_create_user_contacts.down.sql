-- Откат миграции контактов
DROP TRIGGER IF EXISTS update_user_privacy_settings_updated_at ON user_privacy_settings;
DROP TABLE IF EXISTS user_privacy_settings;

DROP TRIGGER IF EXISTS update_user_contacts_updated_at ON user_contacts;
DROP FUNCTION IF EXISTS update_user_contacts_updated_at();
DROP TABLE IF EXISTS user_contacts;

-- Удаление изменений в таблице сообщений
ALTER TABLE marketplace_messages 
DROP COLUMN IF EXISTS has_attachments,
DROP COLUMN IF EXISTS attachments_count;

-- Удаление таблицы вложений
DROP TABLE IF EXISTS chat_attachments;

-- Удаляем частичный уникальный индекс для прямых сообщений
DROP INDEX IF EXISTS idx_unique_direct_chat;

-- Откат миграции 000041_add_performance_indexes

-- Удаление индексов для чата
DROP INDEX IF EXISTS idx_marketplace_messages_unread;
DROP INDEX IF EXISTS idx_marketplace_chats_user_lookup;
DROP INDEX IF EXISTS idx_marketplace_chats_archived;
DROP INDEX IF EXISTS idx_marketplace_messages_chat_ordered;
DROP INDEX IF EXISTS idx_chat_attachments_message;
DROP INDEX IF EXISTS idx_marketplace_chats_listing;
DROP INDEX IF EXISTS idx_marketplace_messages_chat_unread;
DROP INDEX IF EXISTS idx_marketplace_messages_receiver_unread_count;
DROP INDEX IF EXISTS idx_marketplace_messages_chat_last;
DROP INDEX IF EXISTS idx_marketplace_chats_participants;
DROP INDEX IF EXISTS idx_marketplace_chats_active_sorted;
DROP INDEX IF EXISTS idx_marketplace_chats_listing_participants;

-- Удаление индексов для объявлений
DROP INDEX IF EXISTS idx_marketplace_listings_status_created;
DROP INDEX IF EXISTS idx_marketplace_listings_location;
DROP INDEX IF EXISTS idx_marketplace_listings_category_status;
DROP INDEX IF EXISTS idx_marketplace_listings_user_status;
DROP INDEX IF EXISTS idx_marketplace_listings_city;
DROP INDEX IF EXISTS idx_marketplace_listings_title_gin;
DROP INDEX IF EXISTS idx_marketplace_listings_price;

-- Удаление индексов для пользователей
DROP INDEX IF EXISTS idx_users_email_lower;
DROP INDEX IF EXISTS idx_users_active;

-- Удаление индексов для избранного
DROP INDEX IF EXISTS idx_marketplace_favorites_user_count;
DROP INDEX IF EXISTS idx_marketplace_favorites_listing;

-- Удаление индексов для уведомлений
DROP INDEX IF EXISTS idx_notifications_user_unread;

-- Удаление индексов для переводов
DROP INDEX IF EXISTS idx_translations_type_lang;

-- Удаление индексов для категорий
DROP INDEX IF EXISTS idx_marketplace_categories_slug;

-- Удаление дополнительных индексов
DROP INDEX IF EXISTS idx_marketplace_images_listing_main;
DROP INDEX IF EXISTS idx_review_responses_review;
DROP INDEX IF EXISTS idx_price_history_current;
-- Complete migration that combines all previous migrations
-- First, drop all existing tables if they exist (in correct order due to dependencies)
DROP TABLE IF EXISTS marketplace_messages CASCADE;
DROP TABLE IF EXISTS marketplace_chats CASCADE;
DROP TABLE IF EXISTS marketplace_favorites CASCADE;
DROP TABLE IF EXISTS marketplace_images CASCADE;
DROP TABLE IF EXISTS translations CASCADE;
DROP TABLE IF EXISTS notification_settings CASCADE;
DROP TABLE IF EXISTS notifications CASCADE;
DROP TABLE IF EXISTS review_votes CASCADE;
DROP TABLE IF EXISTS review_responses CASCADE;
DROP TABLE IF EXISTS reviews CASCADE;
DROP TABLE IF EXISTS marketplace_listings CASCADE;
DROP TABLE IF EXISTS marketplace_categories CASCADE;
DROP TABLE IF EXISTS user_telegram_connections CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS schema_migrations CASCADE;

-- Drop existing functions
DROP FUNCTION IF EXISTS update_marketplace_chats_updated_at CASCADE;
DROP FUNCTION IF EXISTS update_updated_at_column CASCADE;
DROP FUNCTION IF EXISTS update_user_updated_at CASCADE;
DROP FUNCTION IF EXISTS update_notification_settings_updated_at CASCADE;
DROP FUNCTION IF EXISTS update_translations_updated_at CASCADE;
DROP FUNCTION IF EXISTS calculate_entity_rating CASCADE;

-- Create base tables
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(150) NOT NULL UNIQUE,
    google_id VARCHAR(255) UNIQUE,
    picture_url TEXT,
    phone VARCHAR(20),
    bio TEXT,
    notification_email BOOLEAN DEFAULT true,
    timezone VARCHAR(50) DEFAULT 'UTC',
    last_seen TIMESTAMP,
    account_status VARCHAR(20) DEFAULT 'active',
    settings JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT users_account_status_check CHECK (account_status IN ('active', 'inactive', 'suspended'))
);

CREATE TABLE marketplace_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    parent_id INT REFERENCES marketplace_categories(id),
    icon VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE marketplace_listings (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    category_id INT REFERENCES marketplace_categories(id),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(12,2),
    condition VARCHAR(50),
    status VARCHAR(20) DEFAULT 'active',
    location VARCHAR(255),
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8),
    address_city VARCHAR(100),
    address_country VARCHAR(100),
    views_count INT DEFAULT 0,
    show_on_map BOOLEAN NOT NULL DEFAULT true,
    original_language VARCHAR(10) DEFAULT 'sr',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE marketplace_images (
    id SERIAL PRIMARY KEY,
    listing_id INT REFERENCES marketplace_listings(id) ON DELETE CASCADE,
    file_path VARCHAR(255) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_size INT NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    is_main BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE marketplace_favorites (
    user_id INT REFERENCES users(id),
    listing_id INT REFERENCES marketplace_listings(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, listing_id)
);

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

CREATE TABLE marketplace_messages (
    id SERIAL PRIMARY KEY,
    chat_id INT REFERENCES marketplace_chats(id) ON DELETE CASCADE,
    listing_id INT REFERENCES marketplace_listings(id) ON DELETE CASCADE,
    sender_id INT REFERENCES users(id),
    receiver_id INT REFERENCES users(id),
    content TEXT NOT NULL,
    is_read BOOLEAN DEFAULT false,
    original_language VARCHAR(2) DEFAULT 'en',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE reviews (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    entity_type VARCHAR(50) NOT NULL,
    entity_id INT NOT NULL,
    rating INT NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,
    pros TEXT,
    cons TEXT,
    photos TEXT[],
    likes_count INT DEFAULT 0,
    helpful_votes INT DEFAULT 0,
    not_helpful_votes INT DEFAULT 0,
    is_verified_purchase BOOLEAN DEFAULT false,
    status VARCHAR(20) DEFAULT 'published' CHECK (status IN ('draft', 'published', 'hidden')),
    original_language VARCHAR(2) DEFAULT 'en',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE review_responses (
    id SERIAL PRIMARY KEY,
    review_id INT NOT NULL REFERENCES reviews(id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    response TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE review_votes (
    review_id INT NOT NULL REFERENCES reviews(id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    vote_type VARCHAR(20) NOT NULL CHECK (vote_type IN ('helpful', 'not_helpful')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (review_id, user_id)
);

CREATE TABLE translations (
    id SERIAL PRIMARY KEY,
    entity_type VARCHAR(50) NOT NULL,
    entity_id INTEGER NOT NULL,
    language VARCHAR(10) NOT NULL,
    field_name VARCHAR(50) NOT NULL,
    translated_text TEXT NOT NULL,
    is_machine_translated BOOLEAN DEFAULT true,
    is_verified BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(entity_type, entity_id, language, field_name)
);

CREATE TABLE user_telegram_connections (
    user_id INT PRIMARY KEY REFERENCES users(id),
    telegram_chat_id VARCHAR(100) NOT NULL,
    telegram_username VARCHAR(100),
    connected_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE notification_settings (
    user_id INT NOT NULL REFERENCES users(id),
    notification_type VARCHAR(50) NOT NULL,
    telegram_enabled BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, notification_type)
);

CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    type VARCHAR(50) NOT NULL,
    title TEXT NOT NULL,
    message TEXT NOT NULL,
    data JSONB,
    is_read BOOLEAN DEFAULT false,
    delivered_to JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE schema_migrations (
    version bigint NOT NULL PRIMARY KEY,
    dirty boolean NOT NULL
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone);
CREATE INDEX IF NOT EXISTS idx_users_status ON users(account_status);

CREATE INDEX IF NOT EXISTS idx_marketplace_messages_chat ON marketplace_messages(chat_id);
CREATE INDEX IF NOT EXISTS idx_marketplace_messages_listing ON marketplace_messages(listing_id);
CREATE INDEX IF NOT EXISTS idx_marketplace_messages_sender ON marketplace_messages(sender_id);
CREATE INDEX IF NOT EXISTS idx_marketplace_messages_receiver ON marketplace_messages(receiver_id);
CREATE INDEX IF NOT EXISTS idx_marketplace_messages_created ON marketplace_messages(created_at);

CREATE INDEX IF NOT EXISTS idx_marketplace_chats_buyer ON marketplace_chats(buyer_id);
CREATE INDEX IF NOT EXISTS idx_marketplace_chats_seller ON marketplace_chats(seller_id);
CREATE INDEX IF NOT EXISTS idx_marketplace_chats_updated ON marketplace_chats(updated_at);

CREATE INDEX IF NOT EXISTS idx_marketplace_listings_status ON marketplace_listings(status);

CREATE INDEX IF NOT EXISTS idx_translations_lookup ON translations(entity_type, entity_id, language);

CREATE INDEX IF NOT EXISTS idx_reviews_entity ON reviews(entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_reviews_user ON reviews(user_id);
CREATE INDEX IF NOT EXISTS idx_reviews_rating ON reviews(rating);
CREATE INDEX IF NOT EXISTS idx_reviews_status ON reviews(status);

CREATE INDEX IF NOT EXISTS idx_notifications_user ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_type ON notifications(type);
CREATE INDEX IF NOT EXISTS idx_notifications_created ON notifications(created_at);

-- Create trigger functions
CREATE OR REPLACE FUNCTION update_marketplace_chats_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE OR REPLACE FUNCTION update_user_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE OR REPLACE FUNCTION update_notification_settings_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE OR REPLACE FUNCTION update_translations_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE OR REPLACE FUNCTION calculate_entity_rating(p_entity_type VARCHAR, p_entity_id INT)
RETURNS NUMERIC AS $$
DECLARE
    avg_rating NUMERIC;
BEGIN
    SELECT COALESCE(AVG(rating)::NUMERIC(3,2), 0)
    INTO avg_rating
    FROM reviews
    WHERE entity_type = p_entity_type 
    AND entity_id = p_entity_id 
    AND status = 'published';
    RETURN avg_rating;
END;
$$ LANGUAGE plpgsql;

-- Create triggers
CREATE TRIGGER update_marketplace_chats_timestamp
    BEFORE UPDATE ON marketplace_chats
    FOR EACH ROW
    EXECUTE FUNCTION update_marketplace_chats_updated_at();

CREATE TRIGGER update_marketplace_messages_timestamp
    BEFORE UPDATE ON marketplace_messages
    FOR EACH ROW
    EXECUTE FUNCTION update_marketplace_chats_updated_at();

CREATE TRIGGER update_reviews_updated_at
    BEFORE UPDATE ON reviews
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_review_responses_updated_at
    BEFORE UPDATE ON review_responses
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_user_updated_at();

CREATE TRIGGER update_notification_settings_timestamp
    BEFORE UPDATE ON notification_settings
    FOR EACH ROW
    EXECUTE FUNCTION update_notification_settings_updated_at();

CREATE TRIGGER update_translations_timestamp
    BEFORE UPDATE ON translations
    FOR EACH ROW
    EXECUTE FUNCTION update_translations_updated_at();

-- Insert initial data
INSERT INTO users (id, name, email, created_at, google_id, picture_url, phone, bio, notification_email, timezone, last_seen, account_status, settings, updated_at) VALUES
(1, 'Demo User', 'test@example.com', '2025-02-07 07:13:52.804162', NULL, NULL, NULL, NULL, true, 'UTC', NULL, 'active', '{}', '2025-02-07 07:13:52.887303'),
(2, 'Dmitry Voroshilov', 'voroshilovdo@gmail.com', '2025-02-07 07:14:48.873609', '102440686443518161778', 'https://lh3.googleusercontent.com/a/ACg8ocI7sRI0jzgpBiZmqITMtqT_BavTMXgiBJcZ5Qy1GlqyvN3BdlI=s96-c', NULL, NULL, true, 'UTC', NULL, 'active', '{}', '2025-02-07 07:14:48.873609'),
(3, 'Mail Box', 'boxmail386@gmail.com', '2025-02-07 07:48:08.691308', '105865404629114315097', 'https://lh3.googleusercontent.com/a/ACg8ocKlutGBMU-dy0DBOKpu78fL5o3b-PxMVaOwtlyNqEuEXNJRfg=s96-c', NULL, NULL, true, 'UTC', NULL, 'active', '{}', '2025-02-07 07:48:08.691308'),
(4, 'Azamat Salakhov', 'azamat777salakhov@gmail.com', '2025-02-08 11:50:47.665263', '112906345377302904578', 'https://lh3.googleusercontent.com/a/ACg8ocKQ_MwCSDPuSujNEC7UVdGWfEO6CG0D2RtaxLJ80JJcVrDKfA=s96-c', NULL, NULL, true, 'UTC', NULL, 'active', '{}', '2025-02-08 11:50:47.665263');
SELECT setval('users_id_seq', 4, true);

INSERT INTO user_telegram_connections (user_id, telegram_chat_id, telegram_username, connected_at) VALUES
(2, '158107689', 'Dmvool', '2025-02-07 07:34:33.511245'),
(4, '888561089', 'A777A1982', '2025-02-08 11:31:36.516123');


INSERT INTO notification_settings (user_id, notification_type, telegram_enabled, created_at, updated_at) VALUES
(2, 'new_message', true, '2025-02-07 07:14:49.238702', '2025-02-07 07:34:33.513959'),
(2, 'new_review', true, '2025-02-07 07:14:49.240252', '2025-02-07 07:34:33.519069'),
(2, 'review_vote', true, '2025-02-07 07:14:49.240911', '2025-02-07 07:34:33.519634'),
(2, 'review_response', true, '2025-02-07 07:14:49.241442', '2025-02-07 07:34:33.522901'),
(2, 'listing_status', true, '2025-02-07 07:14:49.241945', '2025-02-07 07:34:33.523478'),
(2, 'favorite_price', true, '2025-02-07 07:14:49.243803', '2025-02-07 07:34:33.523998'),
(4, 'new_message', true, '2025-02-08 11:38:09.274778', '2025-02-07 07:48:09.274778'),
(4, 'new_review', true, '2025-02-08 11:38:09.27687', '2025-02-07 07:48:09.27687'),
(4, 'review_vote', true, '2025-02-08 11:38:09.277562', '2025-02-07 07:48:09.277562'),
(4, 'review_response', true, '2025-02-08 11:38:09.278246', '2025-02-07 07:48:09.278246'),
(4, 'listing_status', true, '2025-02-08 11:38:09.279014', '2025-02-07 07:48:09.279014'),
(4, 'favorite_price', true, '2025-02-08 11:38:09.279632', '2025-02-07 07:48:09.279632');




-- Set initial schema version
INSERT INTO schema_migrations (version, dirty) VALUES (1, false);

-- Create final indexes for optimization
CREATE INDEX IF NOT EXISTS idx_marketplace_listings_status ON marketplace_listings(status);
CREATE INDEX IF NOT EXISTS idx_translations_lookup ON translations(entity_type, entity_id, language);
CREATE INDEX IF NOT EXISTS idx_reviews_entity ON reviews(entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_reviews_user ON reviews(user_id);
CREATE INDEX IF NOT EXISTS idx_reviews_rating ON reviews(rating);
CREATE INDEX IF NOT EXISTS idx_reviews_status ON reviews(status);
CREATE INDEX IF NOT EXISTS idx_notifications_user ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_type ON notifications(type);
CREATE INDEX IF NOT EXISTS idx_notifications_created ON notifications(created_at);

CREATE INDEX IF NOT EXISTS idx_marketplace_categories_parent ON marketplace_categories(parent_id);
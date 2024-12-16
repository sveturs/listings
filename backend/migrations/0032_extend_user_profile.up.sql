-- backend/migrations/0032_extend_user_profile.up.sql
ALTER TABLE users
    ADD COLUMN phone VARCHAR(20),
    ADD COLUMN bio TEXT,
    ADD COLUMN notification_email BOOLEAN DEFAULT true,
    ADD COLUMN notification_push BOOLEAN DEFAULT true,
    ADD COLUMN timezone VARCHAR(50) DEFAULT 'UTC',
    ADD COLUMN last_seen TIMESTAMP,
    ADD COLUMN account_status VARCHAR(20) DEFAULT 'active' 
        CHECK (account_status IN ('active', 'inactive', 'suspended')),
    ADD COLUMN settings JSONB DEFAULT '{}';

-- Индекс для поиска по телефону
CREATE INDEX idx_users_phone ON users(phone);

-- Индекс для статуса аккаунта
CREATE INDEX idx_users_status ON users(account_status);
-- Добавляем поле email_enabled в таблицу notification_settings
ALTER TABLE notification_settings ADD COLUMN IF NOT EXISTS email_enabled BOOLEAN DEFAULT FALSE;

-- Обновляем существующие записи, устанавливая email_enabled = true для всех пользователей
UPDATE notification_settings SET email_enabled = true;
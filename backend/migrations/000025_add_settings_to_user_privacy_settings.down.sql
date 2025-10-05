-- Remove settings column from user_privacy_settings table
DROP INDEX IF EXISTS idx_user_privacy_settings_settings;
ALTER TABLE user_privacy_settings DROP COLUMN IF EXISTS settings;

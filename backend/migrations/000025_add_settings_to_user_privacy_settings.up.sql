-- Add settings JSONB column to user_privacy_settings table for chat and other extensible settings
ALTER TABLE user_privacy_settings
ADD COLUMN IF NOT EXISTS settings JSONB DEFAULT '{}'::jsonb;

-- Create index for faster JSONB queries
CREATE INDEX IF NOT EXISTS idx_user_privacy_settings_settings ON user_privacy_settings USING gin(settings);

-- Add comment
COMMENT ON COLUMN user_privacy_settings.settings IS 'JSONB field for extensible user settings (chat preferences, etc)';

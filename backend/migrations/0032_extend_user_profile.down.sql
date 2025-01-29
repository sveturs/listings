ALTER TABLE users
    DROP COLUMN IF EXISTS phone,
    DROP COLUMN IF EXISTS bio,
    DROP COLUMN IF EXISTS notification_email,
    DROP COLUMN IF EXISTS timezone,
    DROP COLUMN IF EXISTS last_seen,
    DROP COLUMN IF EXISTS account_status,
    DROP COLUMN IF EXISTS settings;

DROP INDEX IF EXISTS idx_users_phone;
DROP INDEX IF EXISTS idx_users_status;
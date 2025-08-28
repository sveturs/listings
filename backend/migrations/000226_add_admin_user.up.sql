-- Add test admin user to admin_users table
INSERT INTO admin_users (email, notes, created_at)
VALUES ('test@example.com', 'Test admin user for translation panel', CURRENT_TIMESTAMP)
ON CONFLICT (email) DO UPDATE
SET notes = 'Test admin user for translation panel';
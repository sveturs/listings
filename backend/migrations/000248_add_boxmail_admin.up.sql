-- Add boxmail386@gmail.com to admin_users table for full admin access
INSERT INTO admin_users (email, notes, created_at)
VALUES ('boxmail386@gmail.com', 'Primary admin user', CURRENT_TIMESTAMP)
ON CONFLICT (email) DO UPDATE
SET notes = 'Primary admin user';
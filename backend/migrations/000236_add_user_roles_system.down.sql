-- Drop function
DROP FUNCTION IF EXISTS check_user_permission(INTEGER, VARCHAR);

-- Note: update_updated_at trigger removal would be handled if it exists

-- Remove role_id from users table
ALTER TABLE users DROP COLUMN IF EXISTS role_id;

-- Drop tables in correct order
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS roles;
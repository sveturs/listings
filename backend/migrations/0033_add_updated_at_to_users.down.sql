-- backend/migrations/0033_add_updated_at_to_users.down.sql
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP FUNCTION IF EXISTS update_user_updated_at();
ALTER TABLE users DROP COLUMN IF EXISTS updated_at;
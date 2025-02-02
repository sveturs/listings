--backend/migrations/0021_add_google_id_to_users.up.sql
ALTER TABLE users 
ADD COLUMN google_id VARCHAR(255) UNIQUE,
ADD COLUMN picture_url TEXT;
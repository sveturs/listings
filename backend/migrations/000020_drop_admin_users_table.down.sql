-- Recreate admin_users table (rollback)
CREATE TABLE IF NOT EXISTS admin_users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Add some common indexes
CREATE INDEX IF NOT EXISTS idx_admin_users_email ON admin_users(email);
